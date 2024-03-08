// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"khalehla/kexec"
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

func (tape *FileSystemTapeDevice) GetNodeCategoryType() kexec.NodeCategoryType {
	return kexec.NodeCategoryDevice
}

func (tape *FileSystemTapeDevice) GetNodeDeviceType() kexec.NodeDeviceType {
	return kexec.NodeDeviceTape
}

func (tape *FileSystemTapeDevice) GetNodeModelType() NodeModelType {
	return NodeModelFileSystemTapeDevice
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

func (tape *FileSystemTapeDevice) StartIo(pkt IoPacket) {
	if tape.verbose {
		log.Printf("FSTAPE:%v", pkt.GetString())
	}
	pkt.SetIoStatus(IosInProgress)

	if pkt.GetNodeDeviceType() != tape.GetNodeDeviceType() {
		pkt.SetIoStatus(IosInvalidNodeType)
	} else {
		switch pkt.GetIoFunction() {
		case IofMount:
			tape.doMount(pkt.(*TapeIoPacket))
		case IofRead:
			tape.doRead(pkt.(*TapeIoPacket))
		case IofReset:
			tape.doUnmount(pkt.(*TapeIoPacket))
		case IofRewind:
			tape.doRewind(pkt.(*TapeIoPacket))
		case IofRewindAndUnload:
			tape.doRewind(pkt.(*TapeIoPacket))
			tape.doUnmount(pkt.(*TapeIoPacket))
		case IofUnmount:
			tape.doUnmount(pkt.(*TapeIoPacket))
		case IofWrite:
			tape.doWrite(pkt.(*TapeIoPacket))
		case IofWriteTapeMark:
			tape.doWriteTapeMark(pkt.(*TapeIoPacket))
		default:
			pkt.SetIoStatus(IosInvalidFunction)
		}
	}

	if tape.verbose {
		log.Printf("FSTAPE:ioStatus=%v", pkt.GetIoStatus())
	}
}

func (tape *FileSystemTapeDevice) doMount(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if tape.IsMounted() {
		pkt.SetIoStatus(IosMediaAlreadyMounted)
		return
	}

	f, err := os.OpenFile(pkt.fileName, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0755)
	if err != nil {
		pkt.SetIoStatus(IosSystemError)
		return
	}

	tape.isReady = true

	// At this point, the tape is mounted. The device may not yet be ready.
	tape.file = f
	tape.isWriteProtected = pkt.writeProtected
	tape.currentOffset = 0
	tape.canRead = true
	pkt.SetIoStatus(IosComplete)
}

func (tape *FileSystemTapeDevice) doRead(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(IosDeviceIsNotReady)
	} else if !tape.canRead {
		pkt.SetIoStatus(IosReadNotAllowed)
		return
	}

	header := make([]byte, 8)
	_, err := tape.file.Read(header)
	if err != nil {
		pkt.SetIoStatus(IosSystemError)
		return
	}

	payloadLength := (uint32(header[0]))<<24 | (uint32(header[1]))<<16 | (uint32(header[2]))<<8 | (uint32(header[3]))
	if payloadLength == 0xFFFFFFFF {
		pkt.SetIoStatus(IosEndOfFile)
		return
	} else if payloadLength%9 > 0 {
		pkt.SetIoStatus(IosInvalidTapeBlock)
		return
	}

	payload := make([]byte, payloadLength)
	_, err = tape.file.Read(header)
	if err != nil {
		pkt.SetIoStatus(IosSystemError)
		return
	}

	pkt.buffer = make([]pkg.Word36, payloadLength*2/9)
	pkg.UnpackWord36(payload, pkt.buffer)
	pkt.SetIoStatus(IosComplete)
}

// doRewind effective rewinds the volume to the tape mark
func (tape *FileSystemTapeDevice) doRewind(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(IosDeviceIsNotReady)
	}

	tape.currentOffset = 0
	tape.canRead = true
	pkt.SetIoStatus(IosComplete)
}

// doUnmount unmounts the virtual volume from the device
func (tape *FileSystemTapeDevice) doUnmount(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsMounted() {
		pkt.SetIoStatus(IosMediaNotMounted)
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
	pkt.SetIoStatus(IosComplete)
}

func (tape *FileSystemTapeDevice) doWrite(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(IosDeviceIsNotReady)
	} else if tape.isWriteProtected {
		pkt.SetIoStatus(IosWriteProtected)
		return
	} else if pkt.buffer == nil {
		pkt.SetIoStatus(IosNilBuffer)
		return
	}

	dataLength := len(pkt.buffer)
	if dataLength == 0 || dataLength&0x01 == 0x01 {
		pkt.SetIoStatus(IosInvalidTapeBlock)
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

	pkg.PackWord36(pkt.buffer, bytes[8:])
	_, err := tape.file.Write(bytes)
	if err != nil {
		pkt.ioStatus = IosSystemError
		return
	}

	tape.previousPayloadLength = payloadLength
	tape.canRead = false
	pkt.SetIoStatus(IosComplete)
}

func (tape *FileSystemTapeDevice) doWriteTapeMark(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(IosDeviceIsNotReady)
	} else if tape.isWriteProtected {
		pkt.SetIoStatus(IosWriteProtected)
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
		pkt.SetIoStatus(IosSystemError)
		return
	}

	tape.previousPayloadLength = 0
	tape.canRead = false
	pkt.SetIoStatus(IosComplete)
}

func (tape *FileSystemTapeDevice) Dump(destination io.Writer, indent string) {
	// TODO Dump()
}
