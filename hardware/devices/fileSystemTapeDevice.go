// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import (
	"fmt"
	"io"
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
	"log"
	"os"
	"sync"
)

// FileSystemTapeDevice stores tape blocks in a lightly-formatted manner in a host filesystem file.
//
// We store tape blocks and tape marks, with no care for what the content of the blocks might be.
// Especially, we do not recognize nor care about tape labels - that is for higher-level code
// to deal with.
//
// A data block consists of a data block header consisting of
//   - 32-bit length of the data payload (0 to 0xFFFFFFFE bytes)
//   - the actual payload
//   - 32-bit length of the data payload again
//
// A tape mark consists of 8 bytes formatted as
//   - 32-bit 0xFFFFFFFF
type FileSystemTapeDevice struct {
	fileName         *string
	file             *os.File
	buffer           []byte
	isReady          bool
	isWriteProtected bool
	mutex            sync.Mutex
	currentOffset    int64
	canRead          bool
	atLoadPoint      bool
	atEndOfTape      bool
	positionLost     bool
	blocksExtended   int
	filesExtended    uint
	verbose          bool
}

func NewFileSystemTapeDevice() *FileSystemTapeDevice {
	return &FileSystemTapeDevice{
		buffer:           make([]byte, 8192),
		isWriteProtected: true,
	}
}

func (tape *FileSystemTapeDevice) GetBlocksExtended() int {
	return tape.blocksExtended
}

func (tape *FileSystemTapeDevice) GetFile() *os.File {
	return tape.file
}

func (tape *FileSystemTapeDevice) GetFilesExtended() uint {
	return tape.filesExtended
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

func (tape *FileSystemTapeDevice) IsAtLoadPoint() bool {
	return tape.atLoadPoint
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

func (tape *FileSystemTapeDevice) Reset() {
	// nothing to do
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

	if pkt.GetPacketType() != ioPackets.TapePacketType {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
	} else {
		switch pkt.GetIoFunction() {
		case ioPackets.IofMount:
			tape.doMount(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofMoveBackward:
			tape.doMoveBackward(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofMoveForward:
			tape.doMoveForward(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofRead:
			tape.doRead(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofReadBackward:
			tape.doReadBackward(pkt.(*ioPackets.TapeIoPacket))
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
	if pkt.GetListener() != nil {
		pkt.GetListener().IoComplete(pkt)
	}
}

func (tape *FileSystemTapeDevice) doMount(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if pkt.MountInfo == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if tape.IsMounted() {
		pkt.SetIoStatus(ioPackets.IosMediaAlreadyMounted)
		return
	}

	flags := os.O_CREATE | os.O_SYNC
	if !pkt.MountInfo.WriteProtect {
		flags |= os.O_RDWR
	} else {
		flags |= os.O_RDONLY
	}

	f, err := os.OpenFile(pkt.MountInfo.Filename, flags, 0666)
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	tape.isReady = true

	// At this point, the tape is mounted. The device may not yet be ready.
	tape.file = f
	tape.isWriteProtected = pkt.MountInfo.WriteProtect
	tape.currentOffset = 0
	tape.canRead = true
	tape.atLoadPoint = true
	tape.atEndOfTape = false
	tape.positionLost = false
	tape.filesExtended = 0
	tape.blocksExtended = 0

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (tape *FileSystemTapeDevice) doMoveBackward(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if tape.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if !tape.canRead {
		pkt.SetIoStatus(ioPackets.IosReadNotAllowed)
		return
	}

	for {
		if tape.atLoadPoint {
			pkt.SetIoStatus(ioPackets.IosAtLoadPoint)
			return
		}

		tape.currentOffset -= 4
		if tape.currentOffset < 0 {
			tape.positionLost = true
			pkt.SetIoStatus(ioPackets.IosLostPosition)
			return
		}

		cw, err := tape.readControlWord()
		if err != nil {
			pkt.SetIoStatus(ioPackets.IosSystemError)
			tape.positionLost = true
			return
		}

		if cw == 0xFFFFFFFF {
			pkt.SetIoStatus(ioPackets.IosEndOfFile)
			tape.filesExtended--
			tape.blocksExtended = 0
			break
		}

		tape.currentOffset -= int64(cw + 4)
		tape.blocksExtended--
		if tape.currentOffset < 0 {
			tape.positionLost = true
			pkt.SetIoStatus(ioPackets.IosLostPosition)
			return
		}
	}
}

func (tape *FileSystemTapeDevice) doMoveForward(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if tape.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if !tape.canRead {
		pkt.SetIoStatus(ioPackets.IosReadNotAllowed)
		return
	}

	for {
		cw, err := tape.readControlWord()
		if err != nil {
			pkt.SetIoStatus(ioPackets.IosSystemError)
			tape.positionLost = true
			return
		}

		tape.currentOffset += 4
		if cw == 0xFFFFFFFF {
			pkt.SetIoStatus(ioPackets.IosEndOfFile)
			tape.filesExtended++
			tape.blocksExtended = 0
			break
		}

		tape.currentOffset += int64(cw + 4)
		tape.blocksExtended++
	}
}

func (tape *FileSystemTapeDevice) doRead(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if tape.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if !tape.canRead {
		pkt.SetIoStatus(ioPackets.IosReadNotAllowed)
		return
	}

	// read control word header
	cw, err := tape.readControlWord()
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		tape.positionLost = true
		return
	}
	tape.atLoadPoint = false

	// tape mark?
	if cw == 0xFFFFFFFF {
		tape.filesExtended++
		tape.blocksExtended = 0
		pkt.SetIoStatus(ioPackets.IosEndOfFile)
		return
	}

	// read the payload
	tape.expandBuffer(cw)
	err = readExact(tape, tape.buffer, cw, tape.currentOffset)
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		tape.positionLost = true
		return
	}

	// update current offset to beyond this payload and the subsequent end-of-payload control word
	tape.currentOffset += int64(cw + 4)
	tape.blocksExtended++

	pkt.DataLength = cw
	if tape.atEndOfTape {
		pkt.SetIoStatus(ioPackets.IosEndOfTape)
	} else {
		pkt.SetIoStatus(ioPackets.IosComplete)
	}
}

func (tape *FileSystemTapeDevice) doReadBackward(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if !tape.atLoadPoint {
		pkt.SetIoStatus(ioPackets.IosAtLoadPoint)
		return
	}

	if tape.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if !tape.canRead {
		pkt.SetIoStatus(ioPackets.IosReadNotAllowed)
		return
	}

	// read previous control word header
	tape.currentOffset -= int64(4)
	if tape.currentOffset < 0 {
		tape.positionLost = true
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	cw, err := tape.readControlWord()
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		tape.positionLost = true
		return
	}

	if cw == 0xFFFFFFFF {
		tape.filesExtended--
		tape.blocksExtended = 0
		pkt.SetIoStatus(ioPackets.IosEndOfFile)
		return
	}

	// position currentOffset to the beginning of the previous block's payload
	tape.currentOffset -= int64(cw)
	if tape.currentOffset < 0 {
		tape.positionLost = true
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	tape.expandBuffer(cw)
	err = readExact(tape, tape.buffer, cw, tape.currentOffset)
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		tape.positionLost = true
		return
	}

	// fix current offset to point to the control word which precedes the payload we just read
	tape.currentOffset -= int64(4)
	if tape.currentOffset == 0 {
		tape.atLoadPoint = true
	}
	tape.blocksExtended--

	pkt.DataLength = cw
	if tape.atEndOfTape {
		pkt.SetIoStatus(ioPackets.IosEndOfTape)
	} else if tape.atLoadPoint {
		pkt.SetIoStatus(ioPackets.IosAtLoadPoint)
	} else {
		pkt.SetIoStatus(ioPackets.IosComplete)
	}
}

// doRewind effective rewinds the volume to the tape mark
func (tape *FileSystemTapeDevice) doRewind(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	tape.currentOffset = 0
	tape.canRead = true
	tape.positionLost = false
	tape.atLoadPoint = true
	pkt.SetIoStatus(ioPackets.IosAtLoadPoint)
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
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (tape *FileSystemTapeDevice) doWrite(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if pkt.Buffer == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if tape.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if tape.isWriteProtected {
		pkt.SetIoStatus(ioPackets.IosWriteProtected)
		return
	}

	if pkt.Buffer == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	payloadLength := uint32(len(pkt.Buffer))
	err := tape.writeControlWord(payloadLength)
	tape.canRead = false
	if err != nil {
		pkt.IoStatus = ioPackets.IosSystemError
		tape.positionLost = true
		return
	}
	tape.currentOffset += 4

	err = writeExact(tape, pkt.Buffer, payloadLength, tape.currentOffset)
	if err != nil {
		pkt.IoStatus = ioPackets.IosSystemError
		tape.positionLost = true
		return
	}
	tape.currentOffset += int64(payloadLength)

	err = tape.writeControlWord(payloadLength)
	if err != nil {
		pkt.IoStatus = ioPackets.IosSystemError
		tape.positionLost = true
		return
	}
	tape.currentOffset += 4
	tape.blocksExtended++

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (tape *FileSystemTapeDevice) doWriteTapeMark(pkt *ioPackets.TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()

	if !tape.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
	}

	if tape.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if tape.isWriteProtected {
		pkt.SetIoStatus(ioPackets.IosWriteProtected)
		return
	}

	err := tape.writeControlWord(0xFFFFFFFF)
	tape.canRead = false
	if err != nil {
		pkt.IoStatus = ioPackets.IosSystemError
		tape.positionLost = true
		return
	}
	tape.currentOffset += 4
	tape.filesExtended++
	tape.blocksExtended = 0

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (tape *FileSystemTapeDevice) Dump(dest io.Writer, indent string) {
	fnStr := "<none>"
	if tape.fileName != nil {
		fnStr = *tape.fileName
	}

	_, _ = fmt.Fprintf(dest, "%vRdy:%v WProt:%v file:%v ldpt:%v eot:%v lost:%v pos:%v fExt:%v blkExt:%v\n",
		indent,
		tape.isReady,
		tape.isWriteProtected,
		fnStr,
		tape.atLoadPoint,
		tape.atEndOfTape,
		tape.positionLost,
		tape.currentOffset,
		tape.filesExtended,
		tape.blocksExtended)
}

func (tape *FileSystemTapeDevice) expandBuffer(requiredSize uint32) {
	newSize := uint32(len(tape.buffer))
	for newSize < requiredSize {
		newSize += 8192
	}
	tape.buffer = nil
	tape.buffer = make([]byte, newSize)
}

// readControlWord uses readExact to read a 4-byte control word from the device current offset.
// does not update current offset.
func (tape *FileSystemTapeDevice) readControlWord() (uint32, error) {
	err := readExact(tape, tape.buffer, 4, tape.currentOffset)
	if err != nil {
		return 0, err
	}

	result :=
		(uint32(tape.buffer[0]))<<24 |
			(uint32(tape.buffer[1]))<<16 |
			(uint32(tape.buffer[2]))<<8 |
			(uint32(tape.buffer[3]))
	return result, nil
}

// writeControlWord uses writeExact to write a 4-byte control word to the device current offset.
// Does not update current offset.
func (tape *FileSystemTapeDevice) writeControlWord(value uint32) error {
	tape.buffer[0] = byte(value >> 24)
	tape.buffer[1] = byte(value >> 16)
	tape.buffer[2] = byte(value >> 8)
	tape.buffer[3] = byte(value)

	return writeExact(tape, tape.buffer, 4, tape.currentOffset)
}
