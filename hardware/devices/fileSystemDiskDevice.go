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

/*
TODO this is obsolete
// This is a very simple pseudo disk device
// We depart from the conventional disk layout in the following manner:
//   There is no booting from disk, so there is no bootstrap in physical blocks 0 or 1
//   The label is always located in physical block 0 instead of physical block 2
//   The first directory track is always located at the second track-aligned physical block
//   We do not store HMBT / SMBT in the directory - MFD can derive that and manage it more efficiently
//   This means that the very first directory track DOES have a DAS in sector 0
// This means that we can always determine whether a pseudo-pack is prepped, and at what prep factor,
// and how many tracks it holds, by simply reading the label from the first block.
// Even though the label is only 28 words, we're okay because we don't actually *have* to read a whole block.

// VOL1 disk label - canonically the third physical record of the device, but always the first block for us
// +000     "VOL1" - ASCII
// +001     pack-id - ASCII LJSF
// +002,H1  pack-id continued - ASCII LJSF
// +003     device-relative address of first directory track
// +004,H1  records per track (1792 / prep_factor)
// +004,H2  words per record (prep_factor)
// +005,H2  reserved size in tracks (for DRS) - we don't do DRS
// +011     normally contains SMBT, HMBT info - we just zero it out here
// +014,S1  Prepped-by: 010:Workstation Utility 020:TPREP, 040:DPREP
// +014,S2  Vol1 Version (we use 1 which is wrong, but so are 0, 2, and 3)
// +014,H2  Heads per cylinder
// +016     Disk capacity in tracks
// +017,H1  Words per physical record (also prep_factor)
// +020,H2  Attributes - we just set these all to zeros
// +021     (non-canonical) total number of blocks on pack
*/

// Simple lookup table - the key is words per block, the value is bytes per block padded to power of two
var bytesPerBlockMap = map[hardware.PrepFactor]uint{
	28:   128,
	56:   256,
	112:  512,
	224:  1024,
	448:  2048,
	896:  4096,
	1792: 8192,
}

// fsDiskHeader is the first 16 bytes in the file.
// It is not visible to the host.
var fsIdentifierBytes = []byte{'F', 'S', 'T', 'P'}
var fsIdentifier = pkg.DeserializeUint32FromBuffer(fsIdentifierBytes)

type fsDiskHeader struct {
	identifier uint32 // "FSTP"
	blockSize  uint32
	prepFactor uint32
	trackCount uint32
}

func (h *fsDiskHeader) deserializeFrom(buffer []byte) {
	h.identifier = pkg.DeserializeUint32FromBuffer(buffer[0:4])
	h.blockSize = pkg.DeserializeUint32FromBuffer(buffer[4:8])
	h.prepFactor = pkg.DeserializeUint32FromBuffer(buffer[8:12])
	h.trackCount = pkg.DeserializeUint32FromBuffer(buffer[12:16])
}

func (h *fsDiskHeader) serializeInto(buffer []byte) {
	pkg.SerializeUint32IntoBuffer(h.identifier, buffer[0:4])
	pkg.SerializeUint32IntoBuffer(h.blockSize, buffer[4:8])
	pkg.SerializeUint32IntoBuffer(h.prepFactor, buffer[8:12])
	pkg.SerializeUint32IntoBuffer(h.trackCount, buffer[12:16])
}

type FileSystemDiskDevice struct {
	fileName         *string
	file             *os.File
	diskHeader       fsDiskHeader
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
		pkt := ioPackets.NewDiskIoPacketMount(0, *initialFileName, false)
		dd.doMount(pkt)
		dd.isReady = pkt.GetIoStatus() == ioPackets.IosComplete
	}

	return dd
}

func (disk *FileSystemDiskDevice) FormatFile(
	filename string,
	prepFactor hardware.PrepFactor,
	trackCount hardware.TrackCount,
) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0666)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	fsDiskHeader := fsDiskHeader{
		identifier: fsIdentifier,
		blockSize:  uint32(bytesPerBlockMap[prepFactor]),
		prepFactor: uint32(prepFactor),
		trackCount: uint32(trackCount),
	}

	buffer := make([]byte, 16)
	fsDiskHeader.serializeInto(buffer)
	_, err = f.Write(buffer)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	err = f.Close()
	return err
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

	if pkt.GetNodeDeviceType() != disk.GetNodeDeviceType() {
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
	pkt.GetListener().IoComplete(pkt)
}

func (disk *FileSystemDiskDevice) doMount(pkt *ioPackets.DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if disk.IsMounted() {
		pkt.SetIoStatus(ioPackets.IosMediaAlreadyMounted)
		return
	}

	f, err := os.OpenFile(pkt.Filename, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	buffer := make([]byte, 16)
	bytes, err := f.ReadAt(buffer, 0)
	if err != nil {
		log.Printf("FSDISK:%v", err)
		pkt.SetIoStatus(ioPackets.IosSystemError)
		return
	}

	if bytes != 16 {
		log.Printf("FSDISK:Cannot read fs header")
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	disk.diskHeader.deserializeFrom(buffer)
	if disk.diskHeader.identifier != fsIdentifier {
		log.Printf("FSDISK:Not formatted properly")
		pkt.SetIoStatus(ioPackets.IosPackNotPrepped)
		return
	}

	disk.isReady = true
	disk.file = f
	disk.isWriteProtected = pkt.WriteProtected
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
		pkt.SetIoStatus(ioPackets.IosNilBuffer)
		return
	}

	if uint(len(pkt.Buffer)) != uint(disk.diskHeader.blockSize) {
		pkt.SetIoStatus(ioPackets.IosInvalidBufferSize)
		return
	}

	blocksPerTrack := 1792 / disk.diskHeader.prepFactor
	blockCount := uint64(disk.diskHeader.trackCount) * uint64(blocksPerTrack)
	if uint64(pkt.BlockId) >= blockCount {
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
		pkt.SetIoStatus(ioPackets.IosNilBuffer)
		return
	}

	if uint(len(pkt.Buffer)) != uint(disk.diskHeader.blockSize) {
		pkt.SetIoStatus(ioPackets.IosInvalidBufferSize)
		return
	}

	blocksPerTrack := 1792 / disk.diskHeader.prepFactor
	blockCount := uint64(disk.diskHeader.trackCount) * uint64(blocksPerTrack)
	if uint64(pkt.BlockId) >= blockCount {
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
