// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import (
	"fmt"
	"io"
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
	"khalehla/pkg"
	"log"
	"os"
	"sync"
)

// FileSystemTapeDevice stores tape blocks in a lightly-formatted manner
// in a host filesystem file.
//
// We store tape blocks and tape marks, with no care for what the content of the blocks might be.
// Especially, we do not recognize nor care about tape labels - that is for higher-level code
// to deal with.
//
// A data block consists of a data block header consisting of:
//
//		32-bit length of the data payload (0 to 0xFFFFFFFE bytes)
//	 32-bit length of the previous payload, unless this is the first data block on the volume,
//	     or the first data block after a tape mark
//
// which is then followed by {n} bytes of data, where {n} is an even number corresponding to
// {k} * 9 / 2 where {k} is an even number of Word36 structs, the content of which are packed
// into the payload.
//
// A tape mark consists of 8 bytes formatted as:
//
//	32-bit 0xFFFFFFFF
//	32-bit length of the previous payload, unless this tape mark is at the beginning of the volume,
//	    or follows another tape mark
type FileSystemTapeDevice struct {
	fileName              *string
	file                  *os.File
	isReady               bool
	isWriteProtected      bool
	mutex                 sync.Mutex
	currentOffset         int64
	canRead               bool
	previousPayloadLength uint32
	verbose               bool
}

func NewFileSystemTapeDevice() *FileSystemTapeDevice {
	return &FileSystemTapeDevice{
		isWriteProtected: true,
	}
}

func (tape *FileSystemTapeDevice) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryDevice
}

func (tape *FileSystemTapeDevice) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceTape
}

func (tape *FileSystemTapeDevice) GetNodeModelType() hardware.NodeModelType {
	return hardware.NodeModelFileSystemTapeDevice
}

func (tape *FileSystemTapeDevice) IsMounted() bool {
	return tape.file != nil
}

func (tape *FileSystemTapeDevice) IsReady() bool {
	return tape.isReady
}

func (tape *FileSystemTapeDevice) IsWriteProtected() bool {
	return tape.isWriteProtected
}

func (tape *FileSystemTapeDevice) SetIsReady(flag bool) {
	tape.isReady = flag
}

func (tape *FileSystemTapeDevice) SetIsWriteProtected(flag bool) {
	tape.isWriteProtected = flag
}

func (tape *FileSystemTapeDevice) SetVerbose(flag bool) {
	tape.verbose = flag
}

func (tape *FileSystemTapeDevice) StartIo(pkt ioPackets.IoPacket) {
	if tape.verbose {
		log.Printf("FSTAPE:%v", pkt.GetString())
	}
	pkt.SetIoStatus(ioPackets.IosInProgress)

	if pkt.GetNodeDeviceType() != tape.GetNodeDeviceType() {
		pkt.SetIoStatus(ioPackets.IosInvalidNodeType)
	} else {
		switch pkt.GetIoFunction() {
		case ioPackets.IofMount:
			tape.doMount(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofRead:
			tape.doRead(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofReset:
			tape.doUnmount(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofRewind:
			tape.doRewind(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofRewindAndUnload:
			tape.doRewind(pkt.(*ioPackets.TapeIoPacket))
			tape.doUnmount(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofUnmount:
			tape.doUnmount(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofWrite:
			tape.doWrite(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofWriteTapeMark:
			tape.doWriteTapeMark(pkt.(*ioPackets.TapeIoPacket))
		default:
			pkt.SetIoStatus(ioPackets.IosInvalidFunction)
		}
	}

	if tape.verbose {
		log.Printf("FSTAPE:ioStatus=%v", pkt.GetIoStatus())
	}
}

func (tape *FileSystemTapeDevice) doMount(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if tape.IsMounted() {
		pkt.SetIoStatus(ioPackets.IosMediaAlreadyMounted)
		return
	}

	f, err := os.OpenFile(pkt.Filename, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	tape.isReady = true

	// At this point, the tape is mounted. The device may not yet be ready.
	tape.file = f
	tape.isWriteProtected = pkt.WriteProtected
	tape.currentOffset = 0
	tape.canRead = true
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (tape *FileSystemTapeDevice) doRead(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
	} else if !tape.canRead {
		pkt.SetIoStatus(ioPackets.IosReadNotAllowed)
		return
	}

	header := make([]byte, 8)
	_, err := tape.file.Read(header)
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	payloadLength := (uint32(header[0]))<<24 | (uint32(header[1]))<<16 | (uint32(header[2]))<<8 | (uint32(header[3]))
	if payloadLength == 0xFFFFFFFF {
		pkt.SetIoStatus(ioPackets.IosEndOfFile)
		return
	} else if payloadLength%9 > 0 {
		pkt.SetIoStatus(ioPackets.IosInvalidTapeBlock)
		return
	}

	payload := make([]byte, payloadLength)
	_, err = tape.file.Read(header)
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	pkt.Buffer = make([]pkg.Word36, payloadLength*2/9)
	pkg.UnpackWord36(payload, pkt.Buffer)
	pkt.SetIoStatus(ioPackets.IosComplete)
}

// doRewind effective rewinds the volume to the tape mark
func (tape *FileSystemTapeDevice) doRewind(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
	}

	tape.currentOffset = 0
	tape.canRead = true
	pkt.SetIoStatus(ioPackets.IosComplete)
}

// doUnmount unmounts the virtual volume from the device
func (tape *FileSystemTapeDevice) doUnmount(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsMounted() {
		pkt.SetIoStatus(ioPackets.IosMediaNotMounted)
		return
	}

	err := tape.file.Close()
	if err != nil {
		log.Printf("%v\n", err)
	}

	tape.isReady = false
	tape.file = nil
	tape.currentOffset = 0
	tape.canRead = false
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (tape *FileSystemTapeDevice) doWrite(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
	} else if tape.isWriteProtected {
		pkt.SetIoStatus(ioPackets.IosWriteProtected)
		return
	} else if pkt.Buffer == nil {
		pkt.SetIoStatus(ioPackets.IosNilBuffer)
		return
	}

	dataLength := len(pkt.Buffer)
	if dataLength == 0 || dataLength&0x01 == 0x01 {
		pkt.SetIoStatus(ioPackets.IosInvalidTapeBlock)
		return
	}
	payloadLength := uint32(dataLength * 9 / 2)

	bytes := make([]byte, 8+payloadLength)
	bytes[0] = byte(payloadLength >> 24)
	bytes[1] = byte(payloadLength >> 16)
	bytes[2] = byte(payloadLength >> 8)
	bytes[3] = byte(payloadLength)
	bytes[4] = byte(tape.previousPayloadLength >> 24)
	bytes[5] = byte(tape.previousPayloadLength >> 16)
	bytes[6] = byte(tape.previousPayloadLength >> 8)
	bytes[7] = byte(tape.previousPayloadLength)

	pkg.PackWord36(pkt.Buffer, bytes[8:])
	_, err := tape.file.Write(bytes)
	if err != nil {
		pkt.IoStatus = ioPackets.IosSystemError
		return
	}

	tape.previousPayloadLength = payloadLength
	tape.canRead = false
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (tape *FileSystemTapeDevice) doWriteTapeMark(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
	} else if tape.isWriteProtected {
		pkt.SetIoStatus(ioPackets.IosWriteProtected)
		return
	}

	bytes := make([]byte, 8)
	bytes[0] = 0xff
	bytes[1] = 0xff
	bytes[2] = 0xff
	bytes[3] = 0xff
	bytes[4] = byte(tape.previousPayloadLength >> 24)
	bytes[5] = byte(tape.previousPayloadLength >> 16)
	bytes[6] = byte(tape.previousPayloadLength >> 8)
	bytes[7] = byte(tape.previousPayloadLength)

	_, err := tape.file.Write(bytes)
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	tape.previousPayloadLength = 0
	tape.canRead = false
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (tape *FileSystemTapeDevice) Dump(dest io.Writer, indent string) {
	fnstr := "<none>"
	if tape.fileName != nil {
		fnstr = *tape.fileName
	}

	_, _ = fmt.Fprintf(dest, "%vRdy:%v WProt:%v file:%v, pos:%v\n",
		indent,
		tape.isReady,
		tape.isWriteProtected,
		fnstr,
		tape.currentOffset)
}
