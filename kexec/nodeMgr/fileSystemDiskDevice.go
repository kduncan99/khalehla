// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/kexec"
	"khalehla/kexec/pkg"
	"khalehla/pkg"
	"log"
	"os"
	"sync"
)

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

// Simple lookup table - the key is words per block, the value is bytes per block padded to power of two
var bytesPerBlockMap = map[kexec.PrepFactor]uint{
	28:   128,
	56:   256,
	112:  512,
	224:  1024,
	448:  2048,
	896:  4096,
	1792: 8192,
}

type FileSystemDiskDevice struct {
	fileName         *string
	file             *os.File
	isReady          bool
	isWriteProtected bool
	packName         string
	geometry         *kexec.DiskPackGeometry
	mutex            sync.Mutex
	buffer           []byte
	Verbose          bool
}

func NewFileSystemDiskDevice(initialFileName *string) *FileSystemDiskDevice {
	dd := &FileSystemDiskDevice{
		fileName:         initialFileName,
		isWriteProtected: true,
	}

	if initialFileName != nil {
		pkt := NewDiskIoPacketMount(0, *initialFileName, false)
		dd.doMount(pkt)
		dd.isReady = pkt.GetIoStatus() == IosComplete
	}

	return dd
}

func (disk *FileSystemDiskDevice) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryDevice
}

func (disk *FileSystemDiskDevice) GetNodeModelType() NodeModelType {
	return NodeModelFileSystemDiskDevice
}

func (disk *FileSystemDiskDevice) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceDisk
}

func (disk *FileSystemDiskDevice) GetGeometry() *kexec.DiskPackGeometry {
	return disk.geometry
}

func (disk *FileSystemDiskDevice) IsMounted() bool {
	return disk.file != nil
}

func (disk *FileSystemDiskDevice) IsPrepped() bool {
	return disk.geometry != nil
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

func (disk *FileSystemDiskDevice) StartIo(pkt IoPacket) {
	pkt.SetIoStatus(IosInProgress)

	if pkt.GetNodeDeviceType() != disk.GetNodeDeviceType() {
		pkt.SetIoStatus(IosInvalidNodeType)
	}

	switch pkt.GetIoFunction() {
	case IofMount:
		disk.doMount(pkt.(*DiskIoPacket))
	case IofPrep:
		disk.doPrep(pkt.(*DiskIoPacket))
	case IofRead:
		disk.doRead(pkt.(*DiskIoPacket))
	case IofReadLabel:
		disk.doReadLabel(pkt.(*DiskIoPacket))
	case IofReset:
		disk.doReset(pkt.(*DiskIoPacket))
	case IofUnmount:
		disk.doUnmount(pkt.(*DiskIoPacket))
	case IofWrite:
		disk.doWrite(pkt.(*DiskIoPacket))
	case IofWriteLabel:
		disk.doWriteLabel(pkt.(*DiskIoPacket))
	default:
		pkt.SetIoStatus(IosInvalidFunction)
	}

	if disk.Verbose {
		log.Printf("  ioStatus:%v", pkt.GetIoStatus())
	}
}

func (disk *FileSystemDiskDevice) doMount(pkt *DiskIoPacket) {
	if disk.Verbose {
		log.Printf("doMount fName:%v", pkt.fileName)
	}

	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if disk.IsMounted() {
		pkt.SetIoStatus(IosMediaAlreadyMounted)
		return
	}

	f, err := os.OpenFile(pkt.fileName, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0755)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(IosSystemError)
		return
	}

	disk.isReady = true

	// At this point, the pack is now mounted. It may not be prepped, but that is okay.
	disk.file = f
	disk.isWriteProtected = pkt.writeProtected
	pkt.SetIoStatus(IosComplete)

	err = disk.probeGeometry()
	if err != nil {
		log.Printf("%v\n", err)
	}
}

func (disk *FileSystemDiskDevice) doPrep(pkt *DiskIoPacket) {
	if disk.Verbose {
		log.Printf("doPrep prepF:%v tracks:%v", pkt.prepFactor, pkt.trackCount)
	}

	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(IosMediaNotMounted)
		return
	}

	if !kexec.IsValidPrepFactor(pkt.prepFactor) {
		pkt.SetIoStatus(IosInvalidPrepFactor)
		return
	}

	if pkt.trackCount < 10000 {
		pkt.SetIoStatus(IosInvalidTrackCount)
		return
	}

	if !kexec.IsValidPackName(pkt.packName) {
		pkt.SetIoStatus(IosInvalidPackName)
		return
	}

	// basic geometry - some of these values exist simply so we do not have to constantly cast them to do math
	recordLength := uint64(pkt.prepFactor)
	trackCount := uint64(pkt.trackCount)
	blocksPerTrack := 1792 / recordLength
	blockCount := trackCount * blocksPerTrack
	dirTrackAddr := uint64(1792) // we set this to the device-relative word address of the initial directory track

	// create initial label and write it
	label := make([]pkg.Word36, 28)
	pkg.FromStringToAsciiWithOffset("VOL1", label, 0, 1)
	pkg.FromStringToAsciiWithOffset(pkt.packName, label, 1, 2)
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
	pkg.PackWord36(label, buffer[:126])
	if disk.Verbose {
		pkg.DumpWord36Buffer(label, 7)
		log.Printf("  WriteAt offset=%v len=%v", 0, len(disk.buffer))
	}

	_, err := disk.file.WriteAt(buffer, 0)
	if err != nil {
		log.Printf("Error writing label:%v\n", err)
		pkt.SetIoStatus(IosSystemError)
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
	if !pkt.removable {
		s1[5].SetH1(0_400000)
	}
	s1[010].SetT1(blocksPerTrack)
	s1[010].SetS3(1) // Sector 1 version
	s1[010].SetT3(recordLength)

	// Write the initial directory track
	wx := 0
	wLen := int(recordLength)
	offset := int64(8192)
	ioLen := int64(bytesPerBlockMap[kexec.PrepFactor(recordLength)])
	buffer = make([]byte, ioLen)
	byteCount := wLen * 9 / 2
	for wx < 1792 {
		pkg.PackWord36(dirTrack[wx:wx+wLen], buffer[0:byteCount])
		if disk.Verbose {
			log.Printf("  WriteAt offset=%v len=%v", offset, len(disk.buffer))
		}

		_, err := disk.file.WriteAt(buffer, offset)
		if err != nil {
			log.Printf("Error writing directory track:%v\n", err)
			pkt.SetIoStatus(IosSystemError)
			return
		}

		wx += wLen
		offset += ioLen
	}

	err = disk.probeGeometry()
	if err != nil {
		log.Printf("%v\n", err)
	}

	pkt.ioStatus = IosComplete
}

func (disk *FileSystemDiskDevice) doRead(pkt *DiskIoPacket) {
	if disk.Verbose {
		log.Printf("doRead blkId:%v", pkt.blockId)
	}

	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(IosMediaNotMounted)
		return
	}

	if !disk.IsPrepped() {
		pkt.SetIoStatus(IosPackNotPrepped)
		return
	}

	if pkt.buffer == nil {
		pkt.SetIoStatus(IosNilBuffer)
		return
	}

	if uint(len(pkt.buffer)) != uint(disk.geometry.PrepFactor) {
		pkt.SetIoStatus(IosInvalidBufferSize)
		return
	}

	if uint64(pkt.blockId) >= uint64(disk.geometry.BlockCount) {
		pkt.SetIoStatus(IosInvalidBlockId)
		return
	}

	offset := int64(pkt.blockId) * int64(disk.geometry.PaddedBytesPerBlock)
	if disk.Verbose {
		log.Printf("  ReadAt offset=%v len=%v", offset, len(disk.buffer))
	}

	_, err := disk.file.ReadAt(disk.buffer, offset)
	if err != nil {
		log.Printf("Read Error:%v\n", err)
		pkt.SetIoStatus(IosSystemError)
		return
	}
	pkg.UnpackWord36(disk.buffer[:disk.geometry.BytesPerBlock], pkt.buffer)
	if disk.Verbose {
		pkg.DumpWord36Buffer(pkt.buffer, 7)
	}

	pkt.ioStatus = IosComplete
}

func (disk *FileSystemDiskDevice) doReadLabel(pkt *DiskIoPacket) {
	if disk.Verbose {
		log.Printf("doReadLabel blkId:%v", pkt.blockId)
	}

	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(IosMediaNotMounted)
		return
	}

	if !disk.IsPrepped() {
		pkt.SetIoStatus(IosPackNotPrepped)
		return
	}

	if pkt.buffer == nil {
		pkt.SetIoStatus(IosNilBuffer)
		return
	}

	if uint(len(pkt.buffer)) != 28 {
		pkt.SetIoStatus(IosInvalidBufferSize)
		return
	}

	if disk.Verbose {
		log.Printf("  ReadAt offset=%v len=%v", 0, len(disk.buffer))
	}

	_, err := disk.file.ReadAt(disk.buffer, 0)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(IosSystemError)
		return
	}
	pkg.UnpackWord36(disk.buffer[:126], pkt.buffer)

	pkt.ioStatus = IosComplete
}

// doReset cancels any pending IOs. It is a NOP for us.
func (disk *FileSystemDiskDevice) doReset(pkt *DiskIoPacket) {
	if disk.Verbose {
		log.Printf("doReset")
	}

	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	// nothing to do for now
	pkt.ioStatus = IosComplete
}

func (disk *FileSystemDiskDevice) doUnmount(pkt *DiskIoPacket) {
	if disk.Verbose {
		log.Printf("doRead Unmount")
	}

	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(IosMediaNotMounted)
		return
	}

	err := disk.file.Close()
	if err != nil {
		log.Printf("%v\n", err)
	}

	disk.geometry = nil
	disk.file = nil
	pkt.SetIoStatus(IosComplete)
}

func (disk *FileSystemDiskDevice) doWrite(pkt *DiskIoPacket) {
	if disk.Verbose {
		log.Printf("doWrite blkId:%v", pkt.blockId)
	}

	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(IosMediaNotMounted)
		return
	}

	if !disk.IsPrepped() {
		pkt.SetIoStatus(IosPackNotPrepped)
		return
	}

	if pkt.buffer == nil {
		pkt.SetIoStatus(IosNilBuffer)
		return
	}

	if disk.isWriteProtected {
		pkt.SetIoStatus(IosWriteProtected)
		return
	}

	if uint(len(pkt.buffer)) != uint(disk.geometry.PrepFactor) {
		pkt.SetIoStatus(IosInvalidBufferSize)
		return
	}

	if uint64(pkt.blockId) >= uint64(disk.geometry.BlockCount) {
		pkt.SetIoStatus(IosInvalidBlockId)
		return
	}

	pkg.PackWord36(pkt.buffer, disk.buffer[:disk.geometry.BytesPerBlock])
	offset := int64(pkt.blockId) * int64(bytesPerBlockMap[disk.geometry.PrepFactor])
	_, err := disk.file.WriteAt(disk.buffer, offset)
	if err != nil {
		log.Printf("Write Error:%v\n", err)
		pkt.SetIoStatus(IosSystemError)
		return
	}

	pkt.ioStatus = IosComplete
}

func (disk *FileSystemDiskDevice) doWriteLabel(pkt *DiskIoPacket) {
	if disk.Verbose {
		log.Printf("doWriteLabel")
	}

	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(IosMediaNotMounted)
		return
	}

	if !disk.IsPrepped() {
		pkt.SetIoStatus(IosPackNotPrepped)
		return
	}

	if pkt.buffer == nil {
		pkt.SetIoStatus(IosNilBuffer)
		return
	}

	if disk.isWriteProtected {
		pkt.SetIoStatus(IosWriteProtected)
		return
	}

	if uint(len(pkt.buffer)) != 28 {
		pkt.SetIoStatus(IosInvalidBufferSize)
		return
	}

	if disk.Verbose {
		pkg.DumpWord36Buffer(pkt.buffer, 7)
	}

	pkg.PackWord36(pkt.buffer, disk.buffer)
	_, err := disk.file.WriteAt(disk.buffer, 0)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(IosSystemError)
		return
	}

	pkt.ioStatus = IosComplete
}

// do this any time we need to read the geometry from a (hopefully) prepped pack.
// we will pretend the prep factor is 28 - this will work for block 0
func (disk *FileSystemDiskDevice) probeGeometry() error {
	if disk.Verbose {
		log.Printf("probeGeometry()")
	}

	disk.geometry = nil

	buffer := make([]byte, 128)
	if disk.Verbose {
		log.Printf("  ReadAt offset=%v len=%v", 0, len(disk.buffer))
	}

	_, err := disk.file.ReadAt(buffer, 0)
	if err != nil {
		log.Printf("Cannot read disk label - assuming pack is not prepped\n")
		return err
	}

	label := make([]pkg.Word36, 28)
	pkg.UnpackWord36(buffer[:126], label)
	if disk.Verbose {
		pkg.DumpWord36Buffer(label, 7)
	}

	str := label[0].ToStringAsAscii()
	if str != "VOL1" {
		// invalid label - pack is not prepped
		return fmt.Errorf("invalid label (not VOL1)")
	}

	packName := label[1].ToStringAsAscii() + label[2].ToStringAsAscii()[:2]
	if !kexec.IsValidPackName(packName) {
		return fmt.Errorf("invalid pack name '%v'", packName)
	}

	prepFactor := kexec.PrepFactor(label[4].GetH2())
	if !kexec.IsValidPrepFactor(prepFactor) {
		return fmt.Errorf("invalid prep factor %v", prepFactor)
	}

	blockCount := kexec.BlockCount(label[021].GetW())
	blocksPerTrack := uint(1792 / prepFactor)
	trackCount := kexec.TrackCount(blockCount / kexec.BlockCount(blocksPerTrack))
	sectorsPerBlock := uint(prepFactor / 28)
	bytesPerBlock := uint(prepFactor) * 9 / 2
	paddedBytesPerBlock := bytesPerBlockMap[prepFactor]

	disk.packName = packName
	disk.geometry = &kexec.DiskPackGeometry{
		PrepFactor:           prepFactor,
		BlockCount:           blockCount,
		BlocksPerTrack:       blocksPerTrack,
		TrackCount:           trackCount,
		SectorsPerBlock:      sectorsPerBlock,
		BytesPerBlock:        bytesPerBlock,
		PaddedBytesPerBlock:  paddedBytesPerBlock,
		FirstDirTrackBlockId: 1792 / blocksPerTrack,
	}
	disk.buffer = make([]byte, bytesPerBlockMap[prepFactor])

	return nil
}

func dumpBuffer(buffer []byte) {
	fmt.Println("Byte Buffer:")
	incr := 32
	for bx := 0; bx < len(buffer); bx += incr {
		str := ""
		for by := 0; by < incr; by++ {
			bz := bx + by
			if bz >= len(buffer) {
				break
			} else {
				str += fmt.Sprintf("%02X ", buffer[bz])
			}
		}
		fmt.Println(str)
	}
}

func (disk *FileSystemDiskDevice) Dump(dest io.Writer, indent string) {
	str := fmt.Sprintf("Rdy:%v WProt:%v pack:%v file:%v\n",
		disk.isReady, disk.isWriteProtected, disk.packName, *disk.fileName)
	_, _ = fmt.Fprintf(dest, "%v%v", indent, str)

	if disk.geometry != nil {
		str := fmt.Sprintf("  prep:%v trks:%v blks:%v sec/blk:%v blk/trk:%v bytes/blk:%v padded:%v\n",
			disk.geometry.PrepFactor,
			disk.geometry.TrackCount,
			disk.geometry.BlockCount,
			disk.geometry.SectorsPerBlock,
			disk.geometry.BlocksPerTrack,
			disk.geometry.BytesPerBlock,
			disk.geometry.PaddedBytesPerBlock)
		_, _ = fmt.Fprintf(dest, "%v%v", indent, str)
	}
}
