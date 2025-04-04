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

// TapeChannel routes IOs to the appropriate deviceInfos which it manages.
// Some day in the future we may add caching, perhaps in a CacheTapeChannel.
type TapeChannel struct {
	identifier hardware.NodeIdentifier
	logName    string
	devices    map[hardware.NodeIdentifier]devices.TapeDevice
	cpChannel  chan *ChannelProgram
	ioChannel  chan *ioPackets2.TapeIoPacket
	packetMap  map[ioPackets2.IoPacket]*ChannelProgram
	resetIos   bool
	verbose    bool
	mutex      sync.Mutex
}

func NewTapeChannel() *TapeChannel {
	ch := &TapeChannel{
		identifier: hardware2.GetNextNodeIdentifier(),
		devices:    make(map[hardware.NodeIdentifier]devices.TapeDevice),
		cpChannel:  make(chan *ChannelProgram, 10),
		ioChannel:  make(chan *ioPackets2.TapeIoPacket, 10),
		packetMap:  make(map[ioPackets2.IoPacket]*ChannelProgram),
	}

	ch.logName = fmt.Sprintf("CHTAPE[%v]", ch.identifier)
	go ch.goRoutine()
	return ch
}

func (ch *TapeChannel) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryChannel
}

func (ch *TapeChannel) GetNodeDeviceType() hardware2.NodeDeviceType {
	return hardware2.NodeDeviceTape
}

func (ch *TapeChannel) GetNodeIdentifier() hardware.NodeIdentifier {
	return ch.identifier
}

func (ch *TapeChannel) GetNodeModelType() hardware2.NodeModelType {
	return hardware2.NodeModelTapeLibraryChannel
}

func (ch *TapeChannel) SetVerbose(flag bool) {
	ch.verbose = flag
}

func (ch *TapeChannel) AssignDevice(nodeIdentifier hardware.NodeIdentifier, device devices2.Device) error {
	if device.GetNodeDeviceType() != hardware2.NodeDeviceTape {
		return fmt.Errorf("device is not a tape")
	}

	ch.devices[nodeIdentifier] = device.(devices.TapeDevice)
	return nil
}

func (ch *TapeChannel) Reset() {
	ch.resetIos = true
}

func (ch *TapeChannel) StartIo(cp *ChannelProgram) {
	if ch.verbose {
		logger.LogInfoF(ch.logName, "StartIo:%v", cp.GetString())
	}
	cp.IoStatus = ioPackets2.IosInProgress
	ch.cpChannel <- cp
}

func (ch *TapeChannel) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vTapeChannel %v connections\n", indent, ch.identifier)
	for id := range ch.devices {
		_, _ = fmt.Fprintf(dest, "%v  %v\n", indent, id)
	}

	_, _ = fmt.Fprintf(dest, "%v  Inflight IOs:\n", indent)
	for iop, chp := range ch.packetMap {
		_, _ = fmt.Fprintf(dest, "%v    %v -> %v\n", indent, iop.GetString(), chp)
	}
}

func (ch *TapeChannel) IoComplete(ioPkt ioPackets2.IoPacket) {
	ch.ioChannel <- ioPkt.(*ioPackets2.TapeIoPacket)
}

func (ch *TapeChannel) prepareIoPacket(chProg *ChannelProgram) (*ioPackets2.TapeIoPacket, bool) {
	// Verify the channel program is not too-badly buggered
	for _, cw := range chProg.ControlWords {
		// packed format must have an even number of words
		if cw.Format == TransferPacked && cw.Length%2 != 0 {
			if ch.verbose {
				logger.LogInfoF(ch.logName, "Invalid length in control word:%v", cw.Length)
			}
			return nil, false
		}
		if cw.Direction == DirectionBackward || cw.Direction == DirectionForward {
			// Fwd/Back transfer must be within the limits of the given buffer
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
		}
	}

	pkt := &ioPackets2.TapeIoPacket{}
	pkt.IoFunction = chProg.IoFunction
	pkt.IoStatus = ioPackets2.IosNotStarted
	pkt.Listener = ch

	if chProg.IoFunction == ioPackets.IofWrite {
		// If this is a data transfer to the device,
		// we have to gather and translate caller's Word36 data into a byte buffer.
		wordCount := uint(0)
		byteCount := uint(0)
		for _, cw := range chProg.ControlWords {
			wordCount += cw.Length
			switch cw.Format {
			case Transfer6Bit:
				byteCount += 6 * cw.Length
			case Transfer8Bit:
				byteCount += 4 * cw.Length
			case TransferPacked:
				byteCount += cw.Length * 9 / 2
			}
		}

		pkt.Buffer = make([]byte, byteCount)
		pkt.DataLength = uint32(byteCount)

		bx := uint(0)
		for _, cw := range chProg.ControlWords {
			switch cw.Format {
			case Transfer6Bit:
				byteCount += 6 * cw.Length
			case Transfer8Bit:
				byteCount += 4 * cw.Length
			case TransferPacked:
				byteCount += cw.Length * 9 / 2
			}
			transferFromWords(cw.Buffer, cw.Offset, cw.Length, pkt.Buffer[bx:bx+byteCount], bx, cw.Direction, cw.Format)
		}

		chProg.BytesTransferred = byteCount
		chProg.WordsTransferred = wordCount
	} else if chProg.IoFunction == ioPackets.IofMount {
		pkt.MountInfo = chProg.MountInfo
	}

	return pkt, true
}

func (ch *TapeChannel) resolveIoPacket(chProg *ChannelProgram, ioPacket ioPackets2.IoPacket) {
	nonIntegral := false
	overRun := false
	if chProg.IoFunction == ioPackets.IofRead || chProg.IoFunction == ioPackets.IofReadBackward {
		if chProg.IoStatus == ioPackets2.IosComplete ||
			chProg.IoStatus == ioPackets2.IosEndOfTape ||
			chProg.IoStatus == ioPackets2.IosAtLoadPoint {

			chProg.BytesTransferred = uint(ioPacket.(*ioPackets2.TapeIoPacket).DataLength)
			chProg.WordsTransferred = 0

			byteBuffer := ioPacket.(*ioPackets2.TapeIoPacket).Buffer
			bytesRemaining := chProg.BytesTransferred
			bufferIndex := uint(0)

			for _, cw := range chProg.ControlWords {
				var byteCount uint
				switch cw.Format {
				case Transfer6Bit:
					byteCount += 6 * cw.Length
				case Transfer8Bit:
					byteCount += 4 * cw.Length
				case TransferPacked:
					byteCount += cw.Length * 9 / 2
				}

				subSize := byteCount
				if subSize > bytesRemaining {
					subSize = bytesRemaining
				}

				subBuffer := byteBuffer[bufferIndex : bufferIndex+subSize]
				niFlag, wordCount :=
					transferFromBytes(subBuffer, 0, subSize, cw.Buffer, 0, cw.Direction, cw.Format)
				chProg.WordsTransferred += wordCount
				if niFlag {
					nonIntegral = true
				}
			}

			overRun = bytesRemaining > 0
		}
	}

	if nonIntegral {
		chProg.IoStatus = ioPackets2.IosNonIntegralRead
	} else if overRun {
		chProg.IoStatus = ioPackets2.IosReadOverrun
	} else {
		chProg.IoStatus = ioPacket.GetIoStatus()
	}
}

func (ch *TapeChannel) goRoutine() {
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
