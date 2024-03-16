// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import (
	"fmt"
	"io"
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
	"khalehla/klog"
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
	identifier       hardware.NodeIdentifier
	logName          string
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
	dev := &FileSystemTapeDevice{
		identifier:       hardware.GetNextNodeIdentifier(),
		buffer:           make([]byte, 8192),
		isWriteProtected: true,
	}

	dev.logName = fmt.Sprintf("FSTAPE[%v]", dev.identifier)
	return dev
}

func (dev *FileSystemTapeDevice) GetBlocksExtended() int {
	return dev.blocksExtended
}

func (dev *FileSystemTapeDevice) GetFile() *os.File {
	return dev.file
}

func (dev *FileSystemTapeDevice) GetFilesExtended() uint {
	return dev.filesExtended
}

func (dev *FileSystemTapeDevice) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryDevice
}

func (dev *FileSystemTapeDevice) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceTape
}

func (dev *FileSystemTapeDevice) GetNodeIdentifier() hardware.NodeIdentifier {
	return dev.identifier
}

func (dev *FileSystemTapeDevice) GetNodeModelType() hardware.NodeModelType {
	return hardware.NodeModelFileSystemTapeDevice
}

func (dev *FileSystemTapeDevice) IsAtLoadPoint() bool {
	return dev.atLoadPoint
}

func (dev *FileSystemTapeDevice) IsMounted() bool {
	return dev.file != nil
}

func (dev *FileSystemTapeDevice) IsReady() bool {
	return dev.isReady
}

func (dev *FileSystemTapeDevice) IsWriteProtected() bool {
	return dev.isWriteProtected
}

func (dev *FileSystemTapeDevice) Reset() {
	// nothing to do
}

func (dev *FileSystemTapeDevice) SetIsReady(flag bool) {
	dev.isReady = flag
}

func (dev *FileSystemTapeDevice) SetIsWriteProtected(flag bool) {
	dev.isWriteProtected = flag
}

func (dev *FileSystemTapeDevice) SetVerbose(flag bool) {
	dev.verbose = flag
}

func (dev *FileSystemTapeDevice) StartIo(pkt ioPackets.IoPacket) {
	if dev.verbose {
		klog.LogInfo(dev.logName, pkt.GetString())
	}
	pkt.SetIoStatus(ioPackets.IosInProgress)

	if pkt.GetPacketType() != ioPackets.TapePacketType {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
	} else {
		switch pkt.GetIoFunction() {
		case ioPackets.IofMount:
			dev.doMount(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofMoveBackward:
			dev.doMoveBackward(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofMoveForward:
			dev.doMoveForward(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofRead:
			dev.doRead(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofReadBackward:
			dev.doReadBackward(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofReset:
			dev.doUnmount(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofRewind:
			dev.doRewind(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofRewindAndUnload:
			dev.doRewind(pkt.(*ioPackets.TapeIoPacket))
			dev.doUnmount(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofUnmount:
			dev.doUnmount(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofWrite:
			dev.doWrite(pkt.(*ioPackets.TapeIoPacket))
		case ioPackets.IofWriteTapeMark:
			dev.doWriteTapeMark(pkt.(*ioPackets.TapeIoPacket))
		default:
			pkt.SetIoStatus(ioPackets.IosInvalidFunction)
		}
	}

	if dev.verbose {
		klog.LogInfoF(dev.logName, "ioStatus:%v", ioPackets.IoStatusTable[pkt.GetIoStatus()])
	}
	if pkt.GetListener() != nil {
		pkt.GetListener().IoComplete(pkt)
	}
}

func (dev *FileSystemTapeDevice) doMount(pkt *ioPackets.TapeIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if pkt.MountInfo == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if dev.IsMounted() {
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

	dev.isReady = true

	// At this point, the dev is mounted. The device may not yet be ready.
	dev.file = f
	dev.isWriteProtected = pkt.MountInfo.WriteProtect
	dev.currentOffset = 0
	dev.canRead = true
	dev.atLoadPoint = true
	dev.atEndOfTape = false
	dev.positionLost = false
	dev.filesExtended = 0
	dev.blocksExtended = 0

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (dev *FileSystemTapeDevice) doMoveBackward(pkt *ioPackets.TapeIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if dev.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if !dev.canRead {
		pkt.SetIoStatus(ioPackets.IosReadNotAllowed)
		return
	}

	for {
		if dev.atLoadPoint {
			pkt.SetIoStatus(ioPackets.IosAtLoadPoint)
			return
		}

		dev.currentOffset -= 4
		if dev.currentOffset < 0 {
			dev.positionLost = true
			pkt.SetIoStatus(ioPackets.IosLostPosition)
			return
		}

		cw, err := dev.readControlWord()
		if err != nil {
			pkt.SetIoStatus(ioPackets.IosSystemError)
			dev.positionLost = true
			return
		}

		if cw == 0xFFFFFFFF {
			pkt.SetIoStatus(ioPackets.IosEndOfFile)
			dev.filesExtended--
			dev.blocksExtended = 0
			break
		}

		dev.currentOffset -= int64(cw + 4)
		dev.blocksExtended--
		if dev.currentOffset < 0 {
			dev.positionLost = true
			pkt.SetIoStatus(ioPackets.IosLostPosition)
			return
		}
	}
}

func (dev *FileSystemTapeDevice) doMoveForward(pkt *ioPackets.TapeIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if dev.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if !dev.canRead {
		pkt.SetIoStatus(ioPackets.IosReadNotAllowed)
		return
	}

	for {
		cw, err := dev.readControlWord()
		if err != nil {
			pkt.SetIoStatus(ioPackets.IosSystemError)
			dev.positionLost = true
			return
		}

		dev.currentOffset += 4
		if cw == 0xFFFFFFFF {
			pkt.SetIoStatus(ioPackets.IosEndOfFile)
			dev.filesExtended++
			dev.blocksExtended = 0
			break
		}

		dev.currentOffset += int64(cw + 4)
		dev.blocksExtended++
	}
}

func (dev *FileSystemTapeDevice) doRead(pkt *ioPackets.TapeIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if dev.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if !dev.canRead {
		pkt.SetIoStatus(ioPackets.IosReadNotAllowed)
		return
	}

	// read control word header
	cw, err := dev.readControlWord()
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		dev.positionLost = true
		return
	}
	dev.atLoadPoint = false

	// dev mark?
	if cw == 0xFFFFFFFF {
		dev.filesExtended++
		dev.blocksExtended = 0
		pkt.SetIoStatus(ioPackets.IosEndOfFile)
		return
	}

	// read the payload
	dev.expandBuffer(cw)
	err = readExact(dev, dev.buffer, cw, dev.currentOffset)
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		dev.positionLost = true
		return
	}

	// update current offset to beyond this payload and the subsequent end-of-payload control word
	dev.currentOffset += int64(cw + 4)
	dev.blocksExtended++

	pkt.DataLength = cw
	if dev.atEndOfTape {
		pkt.SetIoStatus(ioPackets.IosEndOfTape)
	} else {
		pkt.SetIoStatus(ioPackets.IosComplete)
	}
}

func (dev *FileSystemTapeDevice) doReadBackward(pkt *ioPackets.TapeIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if !dev.atLoadPoint {
		pkt.SetIoStatus(ioPackets.IosAtLoadPoint)
		return
	}

	if dev.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if !dev.canRead {
		pkt.SetIoStatus(ioPackets.IosReadNotAllowed)
		return
	}

	// read previous control word header
	dev.currentOffset -= int64(4)
	if dev.currentOffset < 0 {
		dev.positionLost = true
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	cw, err := dev.readControlWord()
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		dev.positionLost = true
		return
	}

	if cw == 0xFFFFFFFF {
		dev.filesExtended--
		dev.blocksExtended = 0
		pkt.SetIoStatus(ioPackets.IosEndOfFile)
		return
	}

	// position currentOffset to the beginning of the previous block's payload
	dev.currentOffset -= int64(cw)
	if dev.currentOffset < 0 {
		dev.positionLost = true
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	dev.expandBuffer(cw)
	err = readExact(dev, dev.buffer, cw, dev.currentOffset)
	if err != nil {
		pkt.SetIoStatus(ioPackets.IosSystemError)
		dev.positionLost = true
		return
	}

	// fix current offset to point to the control word which precedes the payload we just read
	dev.currentOffset -= int64(4)
	if dev.currentOffset == 0 {
		dev.atLoadPoint = true
	}
	dev.blocksExtended--

	pkt.DataLength = cw
	if dev.atEndOfTape {
		pkt.SetIoStatus(ioPackets.IosEndOfTape)
	} else if dev.atLoadPoint {
		pkt.SetIoStatus(ioPackets.IosAtLoadPoint)
	} else {
		pkt.SetIoStatus(ioPackets.IosComplete)
	}
}

// doRewind effective rewinds the volume to the tape mark
func (dev *FileSystemTapeDevice) doRewind(pkt *ioPackets.TapeIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	dev.currentOffset = 0
	dev.canRead = true
	dev.positionLost = false
	dev.atLoadPoint = true
	pkt.SetIoStatus(ioPackets.IosAtLoadPoint)
}

// doUnmount unmounts the virtual volume from the device
func (dev *FileSystemTapeDevice) doUnmount(pkt *ioPackets.TapeIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsMounted() {
		pkt.SetIoStatus(ioPackets.IosMediaNotMounted)
		return
	}

	err := dev.file.Close()
	if err != nil {
		klog.LogErrorF(dev.logName, "Error closing file:%v", err)
	}

	dev.isReady = false
	dev.file = nil
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (dev *FileSystemTapeDevice) doWrite(pkt *ioPackets.TapeIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if pkt.Buffer == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if dev.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if dev.isWriteProtected {
		pkt.SetIoStatus(ioPackets.IosWriteProtected)
		return
	}

	if pkt.Buffer == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	payloadLength := uint32(len(pkt.Buffer))
	err := dev.writeControlWord(payloadLength)
	dev.canRead = false
	if err != nil {
		pkt.IoStatus = ioPackets.IosSystemError
		dev.positionLost = true
		return
	}
	dev.currentOffset += 4

	err = writeExact(dev, pkt.Buffer, payloadLength, dev.currentOffset)
	if err != nil {
		pkt.IoStatus = ioPackets.IosSystemError
		dev.positionLost = true
		return
	}
	dev.currentOffset += int64(payloadLength)

	err = dev.writeControlWord(payloadLength)
	if err != nil {
		pkt.IoStatus = ioPackets.IosSystemError
		dev.positionLost = true
		return
	}
	dev.currentOffset += 4
	dev.blocksExtended++

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (dev *FileSystemTapeDevice) doWriteTapeMark(pkt *ioPackets.TapeIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
	}

	if dev.positionLost {
		pkt.SetIoStatus(ioPackets.IosLostPosition)
		return
	}

	if dev.isWriteProtected {
		pkt.SetIoStatus(ioPackets.IosWriteProtected)
		return
	}

	err := dev.writeControlWord(0xFFFFFFFF)
	dev.canRead = false
	if err != nil {
		pkt.IoStatus = ioPackets.IosSystemError
		dev.positionLost = true
		return
	}
	dev.currentOffset += 4
	dev.filesExtended++
	dev.blocksExtended = 0

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (dev *FileSystemTapeDevice) Dump(dest io.Writer, indent string) {
	fnStr := "<none>"
	if dev.fileName != nil {
		fnStr = *dev.fileName
	}

	_, _ = fmt.Fprintf(dest, "%vRdy:%v WProt:%v file:%v ldpt:%v eot:%v lost:%v pos:%v fExt:%v blkExt:%v\n",
		indent,
		dev.isReady,
		dev.isWriteProtected,
		fnStr,
		dev.atLoadPoint,
		dev.atEndOfTape,
		dev.positionLost,
		dev.currentOffset,
		dev.filesExtended,
		dev.blocksExtended)
}

func (dev *FileSystemTapeDevice) expandBuffer(requiredSize uint32) {
	newSize := uint32(len(dev.buffer))
	for newSize < requiredSize {
		newSize += 8192
	}
	dev.buffer = nil
	dev.buffer = make([]byte, newSize)
}

// readControlWord uses readExact to read a 4-byte control word from the device current offset.
// does not update current offset.
func (dev *FileSystemTapeDevice) readControlWord() (uint32, error) {
	err := readExact(dev, dev.buffer, 4, dev.currentOffset)
	if err != nil {
		return 0, err
	}

	result :=
		(uint32(dev.buffer[0]))<<24 |
			(uint32(dev.buffer[1]))<<16 |
			(uint32(dev.buffer[2]))<<8 |
			(uint32(dev.buffer[3]))
	return result, nil
}

// writeControlWord uses writeExact to write a 4-byte control word to the device current offset.
// Does not update current offset.
func (dev *FileSystemTapeDevice) writeControlWord(value uint32) error {
	dev.buffer[0] = byte(value >> 24)
	dev.buffer[1] = byte(value >> 16)
	dev.buffer[2] = byte(value >> 8)
	dev.buffer[3] = byte(value)

	return writeExact(dev, dev.buffer, 4, dev.currentOffset)
}
