// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package channels

import (
	"fmt"
	"io"

	"khalehla/hardware"
	"khalehla/hardware/devices"
	"khalehla/hardware/ioPackets"
	"khalehla/logger"
	hardware2 "khalehla/old/hardware"
	devices2 "khalehla/old/hardware/devices"
	ioPackets2 "khalehla/old/hardware/ioPackets"

	"sync"
	"time"
)

// DiskChannel routes IOs to the appropriate deviceInfos which it manages.
// Some day in the future we may add caching, perhaps in a CacheDiskChannel.
// DiskChannel ChannelPrograms can have only one control word for reads and writes.
type DiskChannel struct {
	identifier hardware.NodeIdentifier
	logName    string
	devices    map[hardware.NodeIdentifier]devices.DiskDevice
	cpChannel  chan *ChannelProgram
	ioChannel  chan *ioPackets2.DiskIoPacket
	packetMap  map[ioPackets2.IoPacket]*ChannelProgram
	resetIos   bool
	verbose    bool
	mutex      sync.Mutex
}

func NewDiskChannel() *DiskChannel {
	ch := &DiskChannel{
		identifier: hardware2.GetNextNodeIdentifier(),
		devices:    make(map[hardware.NodeIdentifier]devices.DiskDevice),
		cpChannel:  make(chan *ChannelProgram, 10),
		ioChannel:  make(chan *ioPackets2.DiskIoPacket, 10),
		packetMap:  make(map[ioPackets2.IoPacket]*ChannelProgram),
	}

	ch.logName = fmt.Sprintf("CHDISK[%v]", ch.identifier)
	go ch.goRoutine()
	return ch
}

func (ch *DiskChannel) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryChannel
}

func (ch *DiskChannel) GetNodeDeviceType() hardware2.NodeDeviceType {
	return hardware2.NodeDeviceDisk
}

func (ch *DiskChannel) GetNodeIdentifier() hardware.NodeIdentifier {
	return ch.identifier
}

func (ch *DiskChannel) GetNodeModelType() hardware2.NodeModelType {
	return hardware2.NodeModelDiskChannel
}

func (ch *DiskChannel) SetVerbose(flag bool) {
	ch.verbose = flag
}

func (ch *DiskChannel) AssignDevice(nodeIdentifier hardware.NodeIdentifier, device devices2.Device) error {
	if device.GetNodeDeviceType() != hardware2.NodeDeviceDisk {
		return fmt.Errorf("device is not a disk")
	}

	ch.devices[nodeIdentifier] = device.(*devices2.FileSystemDiskDevice)
	return nil
}

func (ch *DiskChannel) Reset() {
	ch.resetIos = true
}

func (ch *DiskChannel) StartIo(cp *ChannelProgram) {
	if ch.verbose {
		logger.LogInfoF(ch.logName, "StartIo:%v", cp.GetString())
	}
	cp.IoStatus = ioPackets2.IosInProgress
	ch.cpChannel <- cp
}

func (ch *DiskChannel) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vDiskChannel %v connections\n", indent, ch.identifier)
	for id := range ch.devices {
		_, _ = fmt.Fprintf(dest, "%v  %v\n", indent, id)
	}

	_, _ = fmt.Fprintf(dest, "%v  Inflight IOs:\n", indent)
	for iop, chp := range ch.packetMap {
		_, _ = fmt.Fprintf(dest, "%v    %v -> %v\n", indent, iop.GetString(), chp)
	}
}

func (ch *DiskChannel) IoComplete(ioPkt ioPackets2.IoPacket) {
	ch.ioChannel <- ioPkt.(*ioPackets2.DiskIoPacket)
}

func (ch *DiskChannel) prepareIoPacket(chProg *ChannelProgram) (*ioPackets2.DiskIoPacket, bool) {
	// Verify the channel program is not too-badly buggered.
	// Note that for us, there must be exactly one control word if we are transferring data.
	if chProg.IoFunction == ioPackets.IofRead || chProg.IoFunction == ioPackets.IofWrite {
		if len(chProg.ControlWords) != 1 {
			if ch.verbose {
				logger.LogInfoF(ch.logName, "Invalid number of control words:%v", len(chProg.ControlWords))
			}
			return nil, false
		}
		for _, cw := range chProg.ControlWords {
			if cw.Format != TransferPacked {
				if ch.verbose {
					logger.LogInfoF(ch.logName, "Invalid transfer format in control word:%v", cw.Format)
				}
				return nil, false
			}
			if cw.Length%2 != 0 {
				if ch.verbose {
					logger.LogInfoF(ch.logName, "Invalid length in control word:%v", cw.Length)
				}
				return nil, false
			}
			if cw.Direction == DirectionForward {
				// transfer must be within the limits of the given buffer
				if cw.Offset+cw.Length > uint(len(cw.Buffer)) {
					if ch.verbose {
						logger.LogInfoF(ch.logName, "Invalid control word offset:%v or length:%v for buffer length:%v",
							cw.Offset, cw.Length, len(cw.Buffer))
					}
					return nil, false
				}
			} else if cw.Direction == DirectionStatic {
				// Static transfer start word must be within the limits of the given buffer
				if cw.Offset > uint(len(cw.Buffer)) {
					if ch.verbose {
						logger.LogInfoF(ch.logName, "Invalid control word offset:%v or length:%v for buffer length:%v",
							cw.Offset, cw.Length, len(cw.Buffer))
					}
					return nil, false
				}
			} else {
				if ch.verbose {
					logger.LogInfoF(ch.logName, "Invalid transfer direction:%v", cw.Direction)
				}
				return nil, false
			}
		}
	} else if chProg.IoFunction == ioPackets.IofPrep {
		if chProg.PrepInfo == nil || len(chProg.ControlWords) != 0 {
			return nil, false
		}
	} else if chProg.IoFunction == ioPackets.IofMount {
		if chProg.MountInfo == nil || len(chProg.ControlWords) != 0 {
			return nil, false
		}
	} else if chProg.IoFunction == ioPackets.IofUnmount {
		if len(chProg.ControlWords) != 0 {
			return nil, false
		}
	} else {
		return nil, false
	}

	pkt := &ioPackets2.DiskIoPacket{}
	pkt.IoFunction = chProg.IoFunction
	pkt.IoStatus = ioPackets2.IosNotStarted
	pkt.Listener = ch

	byteCount := uint(0)
	wordCount := uint(0)
	if chProg.IoFunction == ioPackets.IofRead || chProg.IoFunction == ioPackets.IofWrite {
		cw := chProg.ControlWords[0]
		pkt.BlockId = chProg.BlockId
		wordCount += cw.Length
		bytes, ok := hardware.BlockSizeFromPrepFactor[hardware.PrepFactor(wordCount)]
		if !ok {
			return nil, false
		}

		byteCount = uint(bytes)
		pkt.Buffer = make([]byte, byteCount)

		if chProg.IoFunction == ioPackets.IofWrite {
			// If this is a data transfer to the device, we have to translate caller's Word36 data into a byte buffer.
			transferFromWords(cw.Buffer, cw.Offset, cw.Length, pkt.Buffer, 0, cw.Direction, cw.Format)
			chProg.BytesTransferred = byteCount
			chProg.WordsTransferred = wordCount
		}
	} else if chProg.IoFunction == ioPackets.IofRead {
		pkt.BlockId = chProg.BlockId
	} else if chProg.IoFunction == ioPackets.IofMount {
		pkt.MountInfo = chProg.MountInfo
	} else if chProg.IoFunction == ioPackets.IofPrep {
		pkt.PrepInfo = chProg.PrepInfo
	}

	return pkt, true
}

func (ch *DiskChannel) resolveIoPacket(chProg *ChannelProgram, ioPacket ioPackets2.IoPacket) {
	if chProg.IoFunction == ioPackets.IofRead {
		if ioPacket.GetIoStatus() == ioPackets2.IosComplete {
			chProg.BytesTransferred = uint(len(ioPacket.(*ioPackets2.DiskIoPacket).Buffer))
			chProg.WordsTransferred = 0
			buffer := ioPacket.(*ioPackets2.DiskIoPacket).Buffer
			cw := chProg.ControlWords[0]
			byteCount := cw.Length * 9 / 2
			transferFromBytes(buffer, 0, byteCount, cw.Buffer, 0, cw.Direction, cw.Format)
			chProg.WordsTransferred += cw.Length
		}
	}

	chProg.IoStatus = ioPacket.GetIoStatus()
}

func (ch *DiskChannel) goRoutine() {
	logger.LogTrace(ch.logName, "goRoutine started")
	for {
		select {
		case channelProgram := <-ch.cpChannel:
			ch.mutex.Lock()
			dev, ok := ch.devices[channelProgram.NodeIdentifier]
			if !ok {
				channelProgram.IoStatus = ioPackets2.IosDeviceDoesNotExist
				if ch.verbose {
					logger.LogErrorF(ch.logName, "RejectIo:%v", channelProgram.GetString())
				}
				if channelProgram.Listener != nil {
					channelProgram.Listener.ChannelProgramComplete(channelProgram)
				}
				ch.mutex.Unlock()
				break
			}

			ioPkt, ok := ch.prepareIoPacket(channelProgram)
			if !ok {
				channelProgram.IoStatus = ioPackets2.IosInvalidChannelProgram
				if ch.verbose {
					logger.LogErrorF(ch.logName, "RejectIo:%v", channelProgram.GetString())
				}
				if channelProgram.Listener != nil {
					channelProgram.Listener.ChannelProgramComplete(channelProgram)
				}
				ch.mutex.Unlock()
				break
			}

			dev.StartIo(ioPkt)
			ch.packetMap[ioPkt] = channelProgram
			ch.mutex.Unlock()

		case ioPacket := <-ch.ioChannel:
			ch.mutex.Lock()
			channelProgram, ok := ch.packetMap[ioPacket]
			if ok {
				ch.resolveIoPacket(channelProgram, ioPacket)
				if channelProgram.Listener != nil {
					channelProgram.Listener.ChannelProgramComplete(channelProgram)
				}
				delete(ch.packetMap, ioPacket)
			}
			ch.mutex.Unlock()
			if ch.verbose && channelProgram != nil {
				logger.LogInfoF(ch.logName, "EndIo:%v", channelProgram.GetString())
			}

		case <-time.After(100 * time.Millisecond):
			if ch.resetIos {
				logger.LogInfo(ch.logName, "Resetting IOs")
				ch.mutex.Lock()
				for _, chProg := range ch.packetMap {
					chProg.IoStatus = ioPackets2.IosCanceled
					if chProg.Listener != nil {
						chProg.Listener.ChannelProgramComplete(chProg)
					}
				}
				ch.packetMap = make(map[ioPackets2.IoPacket]*ChannelProgram)
				ch.mutex.Unlock()
				ch.resetIos = false
			}
			break
		}
	}
}
