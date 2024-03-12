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
