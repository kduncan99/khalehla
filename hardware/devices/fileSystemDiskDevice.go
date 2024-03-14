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
	fileName         *string
	file             *os.File
	isReady          bool
	isWriteProtected bool
	blockGeometry    *BlockGeometry
	mutex            sync.Mutex
	verbose          bool
}

func NewFileSystemDiskDevice(initialFileName *string) *FileSystemDiskDevice {
	dd := &FileSystemDiskDevice{
		fileName:         initialFileName,
		isWriteProtected: true,
	}

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
		dd.doMount(pkt)
		dd.isReady = pkt.GetIoStatus() == ioPackets.IosComplete
	}

	return dd
}

func (disk *FileSystemDiskDevice) GetDiskGeometry() (
	blockSize hardware.BlockSize,
	blockCount hardware.BlockCount,
	prepFactor hardware.PrepFactor,
	trackCount hardware.TrackCount,
) {
	if disk.blockGeometry != nil {
		blockSize = disk.blockGeometry.BytesPerBlock
		blockCount = disk.blockGeometry.BlockCount
		prepFactor = disk.blockGeometry.WordsPerBlock
		trackCount = disk.blockGeometry.TrackCount
	}
	return
}

func (disk *FileSystemDiskDevice) GetFile() *os.File {
	return disk.file
}

func (disk *FileSystemDiskDevice) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryDevice
}

func (disk *FileSystemDiskDevice) GetNodeModelType() hardware.NodeModelType {
	return hardware.NodeModelFileSystemDiskDevice
}

func (disk *FileSystemDiskDevice) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceDisk
}

func (disk *FileSystemDiskDevice) IsMounted() bool {
	return disk.file != nil
}

func (disk *FileSystemDiskDevice) IsReady() bool {
	return disk.isReady
}

func (disk *FileSystemDiskDevice) IsWriteProtected() bool {
	return disk.isWriteProtected
}

func (disk *FileSystemDiskDevice) Reset() {
	// nothing to do
}

func (disk *FileSystemDiskDevice) SetIsReady(flag bool) {
	disk.isReady = flag
}

func (disk *FileSystemDiskDevice) SetIsWriteProtected(flag bool) {
	disk.isWriteProtected = flag
}

func (disk *FileSystemDiskDevice) SetVerbose(flag bool) {
	disk.verbose = flag
}

func (disk *FileSystemDiskDevice) StartIo(pkt ioPackets.IoPacket) {
	if disk.verbose {
		log.Printf("FSDISK:%v", pkt.GetString())
	}
	pkt.SetIoStatus(ioPackets.IosInProgress)

	if pkt.GetPacketType() != ioPackets.DiskPacketType {
		pkt.SetIoStatus(ioPackets.IosInvalidNodeType)
	} else {
		switch pkt.GetIoFunction() {
		case ioPackets.IofMount:
			disk.doMount(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofPrep:
			disk.doPrep(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofRead:
			disk.doRead(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofReset:
			disk.doReset(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofUnmount:
			disk.doUnmount(pkt.(*ioPackets.DiskIoPacket))
		case ioPackets.IofWrite:
			disk.doWrite(pkt.(*ioPackets.DiskIoPacket))
		default:
			pkt.SetIoStatus(ioPackets.IosInvalidFunction)
		}
	}

	if disk.verbose {
		log.Printf("  ioStatus:%v", pkt.GetIoStatus())
	}
	if pkt.GetListener() != nil {
		pkt.GetListener().IoComplete(pkt)
	}
}

func (disk *FileSystemDiskDevice) doMount(pkt *ioPackets.DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if pkt.MountInfo == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if disk.IsMounted() {
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
		log.Printf("%v\n", err)
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	// pack is now mounted - io status shall be either IoComplete or IoPackNotPrepped
	disk.blockGeometry = nil
	disk.isReady = true
	disk.file = f
	disk.isWriteProtected = pkt.MountInfo.WriteProtect

	// is the pack prepped?
	// We do not know the prep factor (indeed, we won't know that until we read the label)
	// so we cannot read a block... but we do know that the information will be stored
	// Word36-packed in the first 28 words / 126 bytes in the file (if it exists at all).
	buffer := make([]byte, 126)
	err = readExact(disk, buffer, 126, 0)
	if err != nil {
		log.Printf("FSDISK:Cannot read label:%v", err)
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	label := make([]pkg.Word36, 28)
	pkg.ByteArrayPackedToWord36(buffer, 0, 0, label, 0)

	if label[0].ToStringAsAscii() != "VOL1" {
		log.Printf("FSDISK:No VOL1 label:%v", err)
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	prepFactor := hardware.PrepFactor(label[04].GetH2())
	if !hardware.IsValidPrepFactor(prepFactor) {
		log.Printf("FSDISK:VOL1 label contains invalid prep factor:%v", prepFactor)
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	label[2].SetH2(040040)
	packName := strings.TrimRight(label[1].ToStringAsAscii()+label[2].ToStringAsAscii(), " ")
	if !hardware.IsValidPackName(packName) {
		log.Printf("FSDISK:VOL1 label contains invalid pack name:%v", packName)
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	trackCount := hardware.TrackCount(label[016].GetW())
	blocksPerTrack := hardware.BlockCount(1792 / prepFactor)
	blockCount := hardware.BlockCount(uint64(blocksPerTrack) * uint64(trackCount))

	disk.blockGeometry = &BlockGeometry{
		BytesPerBlock:  hardware.BlockSizeFromPrepFactor[prepFactor],
		WordsPerBlock:  prepFactor,
		BlocksPerTrack: blocksPerTrack,
		BlockCount:     blockCount,
		TrackCount:     trackCount,
		Label:          packName,
	}

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (disk *FileSystemDiskDevice) doPrep(pkt *ioPackets.DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if pkt.PrepInfo == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if !disk.IsReady() {
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
	disk.blockGeometry = &BlockGeometry{
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
	label[021].SetW(uint64(disk.blockGeometry.TrackCount))

	buffer := make([]byte, disk.blockGeometry.BytesPerBlock)
	pkg.Word36ToByteArrayPacked(label, 0, 0, buffer, 0)
	err := readExact(disk, buffer, 126, 0)
	if err != nil {
		log.Printf("FSDISK:Cannot write label:%v", err)
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (disk *FileSystemDiskDevice) doRead(pkt *ioPackets.DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if pkt.Buffer == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if disk.blockGeometry == nil {
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	if uint(len(pkt.Buffer)) != uint(disk.blockGeometry.BytesPerBlock) {
		pkt.SetIoStatus(ioPackets.IosInvalidBufferSize)
		return
	}

	if uint64(pkt.BlockId) >= uint64(disk.blockGeometry.BlockCount) {
		pkt.SetIoStatus(ioPackets.IosInvalidBlockId)
		return
	}

	offset := int64(disk.blockGeometry.BytesPerBlock) * int64(pkt.BlockId+1)
	if disk.verbose {
		log.Printf("  ReadAt offset=%v len=%v", offset, disk.blockGeometry.BytesPerBlock)
	}

	index := uint(0)
	remaining := uint(disk.blockGeometry.BytesPerBlock)
	for remaining > 0 {
		count, err := disk.file.ReadAt(pkt.Buffer[index:], offset)
		if err != nil {
			log.Printf("Read Error:%v\n", err)
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
func (disk *FileSystemDiskDevice) doReset(pkt *ioPackets.DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	// nothing to do for now
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (disk *FileSystemDiskDevice) doUnmount(pkt *ioPackets.DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(ioPackets.IosMediaNotMounted)
		return
	}

	err := disk.file.Close()
	if err != nil {
		log.Printf("%v\n", err)
	}

	disk.file = nil
	disk.isReady = false
	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (disk *FileSystemDiskDevice) doWrite(pkt *ioPackets.DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if disk.isWriteProtected {
		pkt.SetIoStatus(ioPackets.IosWriteProtected)
		return
	}

	if pkt.Buffer == nil {
		pkt.SetIoStatus(ioPackets.IosInvalidPacket)
		return
	}

	if disk.blockGeometry == nil {
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	if uint(len(pkt.Buffer)) != uint(disk.blockGeometry.BytesPerBlock) {
		pkt.SetIoStatus(ioPackets.IosInvalidBufferSize)
		return
	}

	if uint64(pkt.BlockId) >= uint64(disk.blockGeometry.BlockCount) {
		pkt.SetIoStatus(ioPackets.IosInvalidBlockId)
		return
	}

	offset := int64(disk.blockGeometry.BytesPerBlock) * int64(pkt.BlockId)
	if disk.verbose {
		log.Printf("  WriteAt offset=%v len=%v", offset, disk.blockGeometry.BytesPerBlock)
	}

	index := uint(0)
	remaining := uint(disk.blockGeometry.BytesPerBlock)
	for remaining > 0 {
		count, err := disk.file.WriteAt(pkt.Buffer[index:], offset)
		if err != nil {
			log.Printf("Write Error:%v\n", err)
			pkt.SetIoStatus(ioPackets.IosSystemError)
			return
		}

		remaining -= uint(count)
		index += uint(count)
		offset += int64(count)
	}

	pkt.SetIoStatus(ioPackets.IosComplete)
}

func (disk *FileSystemDiskDevice) Dump(dest io.Writer, indent string) {
	str := fmt.Sprintf("Rdy:%v WProt:%v file:%v prepped:%v\n",
		disk.isReady, disk.isWriteProtected, *disk.fileName, disk.blockGeometry != nil)
	if disk.blockGeometry != nil {
		str += fmt.Sprintf("bytes/Blk:%v blks/Trk:%v wrds/Blk:%v blks:%v tracks:%v",
			disk.blockGeometry.BytesPerBlock,
			disk.blockGeometry.BlocksPerTrack,
			disk.blockGeometry.WordsPerBlock,
			disk.blockGeometry.BlockCount,
			disk.blockGeometry.TrackCount)
	}

	_, _ = fmt.Fprintf(dest, "%v%v", indent, str)
}
