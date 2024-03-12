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

// This is a very simple pseudo disk device.
// We store a header in the first {n} bytes which describe the basic geometry of the
// virtual disk we are modeling. Then we perform I/O in a manner consistent with the model.

// fileSystemDiskHeader is the first {n} bytes in the file.
// It is not visible to the host.
var fsIdentifierBytes = []byte{'*', 'F', 'S', 'D', 'I', 'S', 'K', '*'}
var fsIdentifier = pkg.DeserializeUint64FromBuffer(fsIdentifierBytes)

const fsDiskHeaderLength = 32

type fileSystemDiskHeader struct {
	identifier uint64
	blockSize  uint32
	prepFactor uint32
	blockCount uint64
	trackCount uint64
}

func deserializeFileSystemDiskHeader(buffer []byte) *fileSystemDiskHeader {
	result := &fileSystemDiskHeader{}
	result.deserializeFrom(buffer)
	return result
}

func (h *fileSystemDiskHeader) deserializeFrom(buffer []byte) {
	h.identifier = pkg.DeserializeUint64FromBuffer(buffer[0:8])
	h.blockSize = pkg.DeserializeUint32FromBuffer(buffer[8:12])
	h.prepFactor = pkg.DeserializeUint32FromBuffer(buffer[12:16])
	h.blockCount = pkg.DeserializeUint64FromBuffer(buffer[16:24])
	h.trackCount = pkg.DeserializeUint64FromBuffer(buffer[24:32])
}

func (h *fileSystemDiskHeader) serializeInto(buffer []byte) {
	pkg.SerializeUint64IntoBuffer(h.identifier, buffer[0:8])
	pkg.SerializeUint32IntoBuffer(h.blockSize, buffer[8:12])
	pkg.SerializeUint32IntoBuffer(h.prepFactor, buffer[12:16])
	pkg.SerializeUint64IntoBuffer(h.blockCount, buffer[16:24])
	pkg.SerializeUint64IntoBuffer(h.trackCount, buffer[24:32])
}

type FileSystemDiskDevice struct {
	fileName         *string
	file             *os.File
	diskHeader       *fileSystemDiskHeader
	isReady          bool
	isWriteProtected bool
	packName         string
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
	trackCount hardware.TrackCount,
) {
	if disk.diskHeader != nil {
		blockSize = hardware.BlockSize(disk.diskHeader.blockSize)
		blockCount = hardware.BlockCount(disk.diskHeader.blockCount)
		trackCount = hardware.TrackCount(disk.diskHeader.trackCount)
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

	disk.diskHeader = nil
	disk.isReady = true
	disk.file = f
	disk.isWriteProtected = pkt.MountInfo.WriteProtect

	buffer := make([]byte, fsDiskHeaderLength)
	err = readExact(disk, buffer, fsDiskHeaderLength, 0)
	if err != nil {
		log.Printf("FSDISK:%v", err)
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	hdr := deserializeFileSystemDiskHeader(buffer)
	if hdr.identifier != fsIdentifier {
		log.Printf("FSDISK:Not formatted properly")
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	disk.diskHeader = hdr
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

	blocksPerTrack := 1792 / uint64(pkt.PrepInfo.PrepFactor)
	disk.diskHeader = &fileSystemDiskHeader{
		identifier: fsIdentifier,
		blockSize:  uint32(hardware.BlockSizeFromPrepFactor[pkt.PrepInfo.PrepFactor]),
		prepFactor: uint32(pkt.PrepInfo.PrepFactor),
		blockCount: blocksPerTrack * uint64(pkt.PrepInfo.TrackCount),
		trackCount: uint64(pkt.PrepInfo.TrackCount),
	}

	buffer := make([]byte, fsDiskHeaderLength)
	disk.diskHeader.serializeInto(buffer)
	err := writeExact(disk, buffer, fsDiskHeaderLength, 0)
	if err != nil {
		log.Printf("Write Error:%v\n", err)
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	pkt.SetIoStatus(ioPackets.IosComplete)
}

/*
TODO saving this because we need it somewhere... just not here
func (disk *FileSystemDiskDevice) doPrep(pkt *ioPackets.DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsReady() {
		pkt.SetIoStatus(ioPackets.IosDeviceIsNotReady)
		return
	}

	if !hardware.IsValidPrepFactor(pkt.PrepFactor) {
		pkt.SetIoStatus(ioPackets.IosInvalidPrepFactor)
		return
	}

	if pkt.TrackCount < 10000 {
		pkt.SetIoStatus(ioPackets.IosInvalidTrackCount)
		return
	}

	if !hardware.IsValidPackName(pkt.PackName) {
		pkt.SetIoStatus(ioPackets.IosInvalidPackName)
		return
	}

	// basic geometry - some of these values exist simply so we do not have to constantly cast them to do math
	recordLength := uint64(pkt.PrepFactor)
	trackCount := uint64(pkt.TrackCount)
	blocksPerTrack := 1792 / recordLength
	blockCount := trackCount * blocksPerTrack
	dirTrackAddr := uint64(1792) // we set this to the device-relative word address of the initial directory track

	// create initial label and write it
	label := make([]pkg.Word36, 28)
	pkg.FromStringToAsciiWithOffset("VOL1", label, 0, 1)
	pkg.FromStringToAsciiWithOffset(pkt.PackName, label, 1, 2)
	label[2].SetH2(0)
	label[3].SetW(dirTrackAddr)
	label[4].SetH1(blocksPerTrack)
	label[4].SetH2(recordLength)
	label[5].SetW(0)      // no DRS tracks
	label[014].SetS1(010) // Pretend we are a workstation utility
	label[014].SetS2(1)   // VOL1 version
	label[014].SetH2(10)  // heads per cylinder - make up something
	label[016].SetW(trackCount)
	label[017].SetH1(recordLength)
	label[021].SetW(blockCount)

	buffer := make([]byte, 128)
	pkg.PackWord36Strict(label, buffer[:126])
	if disk.verbose {
		log.Printf("  WriteAt offset=%v len=%v", 0, len(disk.buffer))
	}

	_, err := disk.file.WriteAt(buffer, 0)
	if err != nil {
		log.Printf("Error writing label:%v\n", err)
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	// initial directory
	dirTrack := make([]pkg.Word36, 1792)
	availableTracks := trackCount - 2 // subtract label track and first directory track

	// sector 0
	das := dirTrack[0:28]
	das[1].SetW(0_600000_000000) // first 2 sectors are allocated
	for dx := 3; dx < 27; dx += 3 {
		das[dx].SetW(0_400000_000000)
	}
	das[27].SetW(0_400000_000000)

	// sector 1
	s1 := dirTrack[28:56]
	// leave +0 and +1 alone (We aren't doing HMBT/SMBT so we don't need the addresses)
	s1[2].SetW(availableTracks)
	s1[3].SetW(availableTracks)
	s1[4].FromStringToFieldata(disk.packName)
	if !pkt.Removable {
		s1[5].SetH1(0_400000)
	}
	s1[010].SetT1(blocksPerTrack)
	s1[010].SetS3(1) // Sector 1 version
	s1[010].SetT3(recordLength)

	// Write the initial directory track
	wx := 0
	wLen := int(recordLength)
	offset := int64(8192)
	ioLen := int64(bytesPerBlockMap[hardware.PrepFactor(recordLength)])
	buffer = make([]byte, ioLen)
	byteCount := wLen * 9 / 2
	for wx < 1792 {
		pkg.PackWord36Strict(dirTrack[wx:wx+wLen], buffer[0:byteCount])
		if disk.verbose {
			log.Printf("  WriteAt offset=%v len=%v", offset, len(disk.buffer))
		}

		_, err := disk.file.WriteAt(buffer, offset)
		if err != nil {
			log.Printf("Error writing directory track:%v\n", err)
			pkt.SetIoStatus(ioPackets.IosSystemError)
			return
		}

		wx += wLen
		offset += ioLen
	}

	err = disk.probeGeometry()
	if err != nil {
		log.Printf("%v\n", err)
	}

	pkt.SetIoStatus(ioPackets.IosComplete)
}
*/

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

	if disk.diskHeader == nil {
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	if uint(len(pkt.Buffer)) != uint(disk.diskHeader.blockSize) {
		pkt.SetIoStatus(ioPackets.IosInvalidBufferSize)
		return
	}

	if uint64(pkt.BlockId) >= disk.diskHeader.blockCount {
		pkt.SetIoStatus(ioPackets.IosInvalidBlockId)
		return
	}

	offset := int64(disk.diskHeader.blockSize) * int64(pkt.BlockId+1)
	if disk.verbose {
		log.Printf("  ReadAt offset=%v len=%v", offset, disk.diskHeader.blockSize)
	}

	index := uint(0)
	remaining := uint(disk.diskHeader.blockSize)
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

	if disk.diskHeader == nil {
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	if uint(len(pkt.Buffer)) != uint(disk.diskHeader.blockSize) {
		pkt.SetIoStatus(ioPackets.IosInvalidBufferSize)
		return
	}

	if uint64(pkt.BlockId) >= disk.diskHeader.blockCount {
		pkt.SetIoStatus(ioPackets.IosInvalidBlockId)
		return
	}

	offset := int64(disk.diskHeader.blockSize) * int64(pkt.BlockId+1)
	if disk.verbose {
		log.Printf("  WriteAt offset=%v len=%v", offset, disk.diskHeader.blockSize)
	}

	index := uint(0)
	remaining := uint(disk.diskHeader.blockSize)
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

// TODO obsolete
//// do this any time we need to read the geometry from a (hopefully) prepped pack.
//// we will pretend the prep factor is 28 - this will work for block 0
//func (disk *FileSystemDiskDevice) probeGeometry() error {
//	disk.geometry = nil
//
//	buffer := make([]byte, 128)
//	if disk.verbose {
//		log.Printf("  ReadAt offset=%v len=%v", 0, len(disk.buffer))
//	}
//
//	_, err := disk.file.ReadAt(buffer, 0)
//	if err != nil {
//		log.Printf("Cannot read disk label - assuming pack is not prepped\n")
//		return err
//	}
//
//	label := make([]pkg.Word36, 28)
//	pkg.UnpackWord36Strict(buffer[:126], label)
//	if disk.verbose {
//		pkg.DumpWord36Buffer(label, 7)
//	}
//
//	str := label[0].ToStringAsAscii()
//	if str != "VOL1" {
//		// invalid label - pack is not prepped
//		return fmt.Errorf("invalid label (not VOL1)")
//	}
//
//	packName := label[1].ToStringAsAscii() + label[2].ToStringAsAscii()[:2]
//	if !hardware.IsValidPackName(packName) {
//		return fmt.Errorf("invalid pack name '%v'", packName)
//	}
//
//	prepFactor := hardware.PrepFactor(label[4].GetH2())
//	if !hardware.IsValidPrepFactor(prepFactor) {
//		return fmt.Errorf("invalid prep factor %v", prepFactor)
//	}
//
//	blockCount := hardware.BlockCount(label[021].GetW())
//	blocksPerTrack := uint(1792 / prepFactor)
//	trackCount := hardware.TrackCount(blockCount / hardware.BlockCount(blocksPerTrack))
//	sectorsPerBlock := uint(prepFactor / 28)
//	bytesPerBlock := uint(prepFactor) * 9 / 2
//	paddedBytesPerBlock := bytesPerBlockMap[prepFactor]
//
//	disk.packName = packName
//	disk.geometry = &ioPackets.DiskPackGeometry{
//		PrepFactor:           prepFactor,
//		BlockCount:           blockCount,
//		BlocksPerTrack:       blocksPerTrack,
//		TrackCount:           trackCount,
//		SectorsPerBlock:      sectorsPerBlock,
//		BytesPerBlock:        bytesPerBlock,
//		PaddedBytesPerBlock:  paddedBytesPerBlock,
//		FirstDirTrackBlockId: 1792 / blocksPerTrack,
//	}
//	disk.buffer = make([]byte, bytesPerBlockMap[prepFactor])
//
//	return nil
//}

func (disk *FileSystemDiskDevice) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vRdy:%v WProt:%v pack:%v file:%v\n",
		indent,
		disk.isReady,
		disk.isWriteProtected,
		disk.packName,
		*disk.fileName)

	str := fmt.Sprintf("  prep:%v trks:%v\n",
		disk.diskHeader.prepFactor,
		disk.diskHeader.trackCount)
	_, _ = fmt.Fprintf(dest, "%v%v", indent, str)
}
