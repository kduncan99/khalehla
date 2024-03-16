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
	"khalehla/pkg"
	"os"
	"strings"
	"sync"
)

// FileSystemDiskDevice is a very simple pseudo disk device.
// Any such information can be considered to be the basic geometry of the device.
// That information includes the block size (in bytes, for this handler), and the number of blocks
// which comprise the extent of the storage.
// Since it is a virtual device, there aren't physical records to help us determine the actual geometry.
// However, the data which is stored in a conventional disk VOL1 label does contain that information.
// Thus, a component of our virtual format process will be to write such a label into block zero.
// A component of our virtual mount process will be to attempt to read such a label from block zero.
// A further consideration is that the label is written in packed 36-bit mode - that is, two 36-bit words
// stored tightly in 9 consecutive bytes. The label is comprised of 28 words, which corresponds to
// 126 consecutive bytes.
// If the label cannot be read when media is mounted, no IO will be permitted.
type FileSystemDiskDevice struct {
	identifier       hardware.NodeIdentifier
	logName          string
	fileName         *string
	file             *os.File
	isReady          bool
	isWriteProtected bool
	blockGeometry    *BlockGeometry
	mutex            sync.Mutex
	verbose          bool
}

func NewFileSystemDiskDevice(initialFileName *string) *FileSystemDiskDevice {
	dev := &FileSystemDiskDevice{
		identifier:       hardware.GetNextNodeIdentifier(),
		fileName:         initialFileName,
		isWriteProtected: true,
	}

	dev.logName = fmt.Sprintf("FSDISK[%v]", dev.identifier)

	if initialFileName != nil {
		mi := &ioPackets.IoMountInfo{
			Filename:     *initialFileName,
			WriteProtect: false,
		}
		pkt := &ioPackets.DiskIoPacket{
			Listener:   nil,
			IoFunction: ioPackets.IofMount,
			MountInfo:  mi,
		}
		dev.doMount(pkt)
		dev.isReady = pkt.GetIoStatus() == ioPackets.IosComplete
	}

	return dev
}

func (dev *FileSystemDiskDevice) GetDiskGeometry() (
	blockSize hardware.BlockSize,
	blockCount hardware.BlockCount,
	prepFactor hardware.PrepFactor,
	trackCount hardware.TrackCount,
) {
	if dev.blockGeometry != nil {
		blockSize = dev.blockGeometry.BytesPerBlock
		blockCount = dev.blockGeometry.BlockCount
		prepFactor = dev.blockGeometry.WordsPerBlock
		trackCount = dev.blockGeometry.TrackCount
	}
	return
}

func (dev *FileSystemDiskDevice) GetFile() *os.File {
	return dev.file
}

func (dev *FileSystemDiskDevice) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryDevice
}

func (dev *FileSystemDiskDevice) GetNodeIdentifier() hardware.NodeIdentifier {
	return dev.identifier
}

func (dev *FileSystemDiskDevice) GetNodeModelType() hardware.NodeModelType {
	return hardware.NodeModelFileSystemDiskDevice
}

func (dev *FileSystemDiskDevice) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceDisk
}

func (dev *FileSystemDiskDevice) IsMounted() bool {
	return dev.file != nil
}

func (dev *FileSystemDiskDevice) IsReady() bool {
	return dev.isReady
}

func (dev *FileSystemDiskDevice) IsWriteProtected() bool {
	return dev.isWriteProtected
}

func (dev *FileSystemDiskDevice) Reset() {
	// nothing to do
}

func (dev *FileSystemDiskDevice) SetIsReady(flag bool) {
	dev.isReady = flag
}

func (dev *FileSystemDiskDevice) SetIsWriteProtected(flag bool) {
	dev.isWriteProtected = flag
}

func (dev *FileSystemDiskDevice) SetVerbose(flag bool) {
	dev.verbose = flag
}

func (dev *FileSystemDiskDevice) StartIo(pkt ioPackets.IoPacket) {
	if dev.verbose {
		klog.LogInfo(dev.logName, pkt.GetString())
	}
	pkt.SetIoStatus(ioPackets.IosInProgress)

	if pkt.GetPacketType() != ioPackets.DiskPacketType {
		pkt.SetIoStatus(ioPackets.IosInvalidNodeType)
	} else {
		switch pkt.GetIoFunction() {
		case ioPackets.IofMount:
			dev.doMount(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofPrep:
			dev.doPrep(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofRead:
			dev.doRead(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofReset:
			dev.doReset(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofUnmount:
			dev.doUnmount(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofWrite:
			dev.doWrite(pkt.(*ioPackets.DiskIoPacket))
		default:
			pkt.SetIoStatus(ioPackets.IosInvalidFunction)
		}
	}

	if dev.verbose {
		klog.LogInfoF(dev.logName, "ioStatus:%v", pkt.GetIoStatus())
	}
	if pkt.GetListener() != nil {
		pkt.GetListener().IoComplete(pkt)
	}
}

func (dev *FileSystemDiskDevice) doMount(pkt *ioPackets.DiskIoPacket) {
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
		klog.LogErrorF(dev.logName, "Error opening file %v:%v", pkt.MountInfo.Filename, err.Error())
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	// pack is now mounted - io status shall be either IoComplete or IoPackNotPrepped
	dev.blockGeometry = nil
	dev.isReady = true
	dev.file = f
	dev.isWriteProtected = pkt.MountInfo.WriteProtect

	// is the pack prepped?
	// We do not know the prep factor (indeed, we won't know that until we read the label)
	// so we cannot read a block... but we do know that the information will be stored
	// Word36-packed in the first 28 words / 126 bytes in the file (if it exists at all).
	buffer := make([]byte, 126)
	err = readExact(dev, buffer, 126, 0)
	if err != nil {
		klog.LogErrorF(dev.logName, "Cannot read label:%v", err)
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	label := make([]pkg.Word36, 28)
	pkg.ByteArrayPackedToWord36(buffer, 0, 0, label, 0)

	if label[0].ToStringAsAscii() != "VOL1" {
		klog.LogError(dev.logName, "No VOL1 label")
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	prepFactor := hardware.PrepFactor(label[04].GetH2())
	if !hardware.IsValidPrepFactor(prepFactor) {
		klog.LogErrorF(dev.logName, "VOL1 label contains invalid prep factor:%v", prepFactor)
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	label[2].SetH2(040040)
	packName := strings.TrimRight(label[1].ToStringAsAscii()+label[2].ToStringAsAscii(), " ")
	if !hardware.IsValidPackName(packName) {
		klog.LogErrorF(dev.logName, "VOL1 label contains invalid pack name:%v", packName)
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	trackCount := hardware.TrackCount(label[016].GetW())
	blocksPerTrack := hardware.BlockCount(1792 / prepFactor)
	blockCount := hardware.BlockCount(uint64(blocksPerTrack) * uint64(trackCount))

	dev.blockGeometry = &BlockGeometry{
		BytesPerBlock:  hardware.BlockSizeFromPrepFactor[prepFactor],
		WordsPerBlock:  prepFactor,
		BlocksPerTrack: blocksPerTrack,
		BlockCount:     blockCount,
		TrackCount:     trackCount,
		Label:          packName,
	}

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (dev *FileSystemDiskDevice) doPrep(pkt *ioPackets.DiskIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if pkt.PrepInfo == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if !hardware.IsValidPrepFactor(pkt.PrepInfo.PrepFactor) {
		pkt.SetIoStatus(ioPackets.IosInvalidPrepFactor)
		return
	}

	if pkt.PrepInfo.TrackCount < 10000 {
		pkt.SetIoStatus(ioPackets.IosInvalidTrackCount)
		return
	}

	if !hardware.IsValidPackName(pkt.PrepInfo.PackName) {
		pkt.SetIoStatus(ioPackets.IosInvalidPackName)
		return
	}

	// Create geometry information
	blocksPerTrack := hardware.BlockCount(1792 / uint64(pkt.PrepInfo.PrepFactor))
	dev.blockGeometry = &BlockGeometry{
		BytesPerBlock:  hardware.BlockSizeFromPrepFactor[pkt.PrepInfo.PrepFactor],
		WordsPerBlock:  pkt.PrepInfo.PrepFactor,
		BlocksPerTrack: blocksPerTrack,
		BlockCount:     hardware.BlockCount(uint64(pkt.PrepInfo.TrackCount) * uint64(blocksPerTrack)),
		TrackCount:     pkt.PrepInfo.TrackCount,
		Label:          pkt.PrepInfo.PackName,
	}

	// Create label record
	label := make([]pkg.Word36, pkt.PrepInfo.PrepFactor)
	for bx := 0; bx < len(label); bx++ {
		label[bx].SetW(0)
	}

	firstDirTrackDRWA := uint64(1792)
	label[0].FromStringToAscii("VOL1")
	label[1].FromStringToAscii(pkt.PrepInfo.PackName[0:4])
	label[2].FromStringToAscii(pkt.PrepInfo.PackName[4:6])
	label[2].SetH2(0)
	label[3].SetW(firstDirTrackDRWA)
	label[4].SetH1(uint64(blocksPerTrack))
	label[4].SetH2(uint64(pkt.PrepInfo.PrepFactor))
	label[5].SetW(0) // no DRS tracks
	// We leave 011 set to zero, because we don't do MBTs
	label[014].SetS1(010) // Pretend we are a workstation utility
	label[014].SetS2(1)   // VOL1 version
	label[014].SetH2(10)  // heads per cylinder - make up something
	label[016].SetW(uint64(pkt.PrepInfo.TrackCount))
	label[017].SetH1(uint64(pkt.PrepInfo.PrepFactor))
	label[021].SetW(uint64(dev.blockGeometry.TrackCount))

	buffer := make([]byte, dev.blockGeometry.BytesPerBlock)
	pkg.Word36ToByteArrayPacked(label, 0, 0, buffer, 0)
	err := readExact(dev, buffer, 126, 0)
	if err != nil {
		klog.LogErrorF(dev.logName, "Cannot write label:%v", err)
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (dev *FileSystemDiskDevice) doRead(pkt *ioPackets.DiskIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if pkt.Buffer == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if dev.blockGeometry == nil {
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	if uint(len(pkt.Buffer)) != uint(dev.blockGeometry.BytesPerBlock) {
		pkt.SetIoStatus(ioPackets.IosInvalidBufferSize)
		return
	}

	if uint64(pkt.BlockId) >= uint64(dev.blockGeometry.BlockCount) {
		pkt.SetIoStatus(ioPackets.IosInvalidBlockId)
		return
	}

	offset := int64(dev.blockGeometry.BytesPerBlock) * int64(pkt.BlockId+1)
	if dev.verbose {
		klog.LogInfoF(dev.logName, "ReadAt offset=%v len=%v", offset, dev.blockGeometry.BytesPerBlock)
	}

	index := uint(0)
	remaining := uint(dev.blockGeometry.BytesPerBlock)
	for remaining > 0 {
		count, err := dev.file.ReadAt(pkt.Buffer[index:], offset)
		if err != nil {
			klog.LogErrorF(dev.logName, "Read Error:%v", err)
			pkt.SetIoStatus(ioPackets.IosSystemError)
			return
		}

		remaining -= uint(count)
		index += uint(count)
		offset += int64(count)
	}

	pkt.SetIoStatus(ioPackets.IosComplete)
}

// doReset cancels any pending IOs. It is a NOP for us.
func (dev *FileSystemDiskDevice) doReset(pkt *ioPackets.DiskIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	// nothing to do for now
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (dev *FileSystemDiskDevice) doUnmount(pkt *ioPackets.DiskIoPacket) {
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

	dev.file = nil
	dev.isReady = false
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (dev *FileSystemDiskDevice) doWrite(pkt *ioPackets.DiskIoPacket) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	if !dev.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
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

	if dev.blockGeometry == nil {
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	if uint(len(pkt.Buffer)) != uint(dev.blockGeometry.BytesPerBlock) {
		pkt.SetIoStatus(ioPackets.IosInvalidBufferSize)
		return
	}

	if uint64(pkt.BlockId) >= uint64(dev.blockGeometry.BlockCount) {
		pkt.SetIoStatus(ioPackets.IosInvalidBlockId)
		return
	}

	offset := int64(dev.blockGeometry.BytesPerBlock) * int64(pkt.BlockId)
	if dev.verbose {
		klog.LogInfoF(dev.logName, "WriteAt offset=%v len=%v", offset, dev.blockGeometry.BytesPerBlock)
	}

	index := uint(0)
	remaining := uint(dev.blockGeometry.BytesPerBlock)
	for remaining > 0 {
		count, err := dev.file.WriteAt(pkt.Buffer[index:], offset)
		if err != nil {
			klog.LogErrorF(dev.logName, "Write Error:%v\n", err)
			pkt.SetIoStatus(ioPackets.IosSystemError)
			return
		}

		remaining -= uint(count)
		index += uint(count)
		offset += int64(count)
	}

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (dev *FileSystemDiskDevice) Dump(dest io.Writer, indent string) {
	str := fmt.Sprintf("Rdy:%v WProt:%v file:%v prepped:%v\n",
		dev.isReady, dev.isWriteProtected, *dev.fileName, dev.blockGeometry != nil)
	if dev.blockGeometry != nil {
		str += fmt.Sprintf("bytes/Blk:%v blks/Trk:%v wrds/Blk:%v blks:%v tracks:%v",
			dev.blockGeometry.BytesPerBlock,
			dev.blockGeometry.BlocksPerTrack,
			dev.blockGeometry.WordsPerBlock,
			dev.blockGeometry.BlockCount,
			dev.blockGeometry.TrackCount)
	}

	_, _ = fmt.Fprintf(dest, "%v%v", indent, str)
}
