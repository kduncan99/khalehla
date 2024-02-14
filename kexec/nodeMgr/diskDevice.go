// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
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
var bytesPerBlockMap = map[types.PrepFactor]uint{
	28:   128,
	56:   256,
	112:  512,
	224:  1024,
	448:  2048,
	896:  4096,
	1792: 8192,
}

type DiskDevice struct {
	fileName         *string
	file             *os.File
	isReady          bool
	isWriteProtected bool
	packName         string
	geometry         *types.DiskPackGeometry
	mutex            sync.Mutex
	buffer           []byte
}

func NewDiskDevice(initialFileName *string) *DiskDevice {
	dd := &DiskDevice{
		fileName:         initialFileName,
		isWriteProtected: true,
	}

	if initialFileName != nil {
		pkt := NewDiskIoPacketMount(0, *initialFileName, false)
		dd.doMount(pkt)
		dd.isReady = pkt.GetIoStatus() == types.IosComplete
	}

	return dd
}

func (disk *DiskDevice) GetNodeType() types.NodeType {
	return types.NodeTypeDisk
}

func (disk *DiskDevice) GetGeometry() *types.DiskPackGeometry {
	return disk.geometry
}

func (disk *DiskDevice) IsMounted() bool {
	return disk.file != nil
}

func (disk *DiskDevice) IsPrepped() bool {
	return disk.geometry != nil
}

func (disk *DiskDevice) IsReady() bool {
	return disk.isReady
}

func (disk *DiskDevice) IsWriteProtected() bool {
	return disk.isWriteProtected
}

func (disk *DiskDevice) SetIsReady(flag bool) {
	disk.isReady = flag
}

func (disk *DiskDevice) SetIsWriteProtected(flag bool) {
	disk.isWriteProtected = flag
}

func (disk *DiskDevice) StartIo(pkt types.IoPacket) {
	pkt.SetIoStatus(types.IosInProgress)

	if pkt.GetNodeType() != disk.GetNodeType() {
		pkt.SetIoStatus(types.IosInvalidNodeType)
	}

	switch pkt.GetIoFunction() {
	case types.IofMount:
		disk.doMount(pkt.(*DiskIoPacket))
	case types.IofPrep:
		disk.doPrep(pkt.(*DiskIoPacket))
	case types.IofRead:
		disk.doRead(pkt.(*DiskIoPacket))
	case types.IofReadLabel:
		disk.doReadLabel(pkt.(*DiskIoPacket))
	case types.IofReset:
		disk.doReset(pkt.(*DiskIoPacket))
	case types.IofUnmount:
		disk.doUnmount(pkt.(*DiskIoPacket))
	case types.IofWrite:
		disk.doWrite(pkt.(*DiskIoPacket))
	case types.IofWriteLabel:
		disk.doWriteLabel(pkt.(*DiskIoPacket))
	default:
		pkt.SetIoStatus(types.IosInvalidFunction)
	}
}

func (disk *DiskDevice) doMount(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if disk.IsMounted() {
		pkt.SetIoStatus(types.IosMediaAlreadyMounted)
		return
	}

	f, err := os.OpenFile(pkt.fileName, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0755)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(types.IosSystemError)
		return
	}

	disk.isReady = true

	// At this point, the pack is now mounted. It may not be prepped, but that is okay.
	disk.file = f
	disk.isWriteProtected = pkt.writeProtected
	pkt.SetIoStatus(types.IosComplete)

	err = disk.probeGeometry()
	if err != nil {
		log.Printf("%v\n", err)
	}
}

func (disk *DiskDevice) doPrep(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(types.IosMediaNotMounted)
		return
	}

	if !IsValidPrepFactor(pkt.prepFactor) {
		pkt.SetIoStatus(types.IosInvalidPrepFactor)
		return
	}

	if pkt.trackCount < 10000 {
		pkt.SetIoStatus(types.IosInvalidTrackCount)
		return
	}

	if !IsValidPackName(pkt.packName) {
		pkt.SetIoStatus(types.IosInvalidPackName)
		return
	}

	// basic geometry - some of these values exist simply so we do not have to constantly cast them to do math
	recordLength := uint64(pkt.prepFactor)
	trackCount := uint64(pkt.trackCount)
	blocksPerTrack := 1792 / recordLength
	blockCount := trackCount * blocksPerTrack
	dirTrackAddr := uint64(1792) // we set this to the device-relative word address of the initial directory track

	// create initial label and write it
	label := make([]pkg.Word36, recordLength)
	pkg.FromStringToAscii("VOL1", label[0:1])
	pkg.FromStringToAscii(pkt.packName, label[1:3])
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

	err := writeBlock(disk.file, 0, label, uint(recordLength))
	if err != nil {
		log.Printf("Error writing label:%v\n", err)
		pkt.SetIoStatus(types.IosSystemError)
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
	err = writeTrack(disk.file, dirTrackAddr, dirTrack, uint(recordLength))
	if err != nil {
		log.Printf("Error writing label:%v\n", err)
		pkt.SetIoStatus(types.IosSystemError)
		return
	}

	err = disk.probeGeometry()
	if err != nil {
		log.Printf("%v\n", err)
	}

	pkt.ioStatus = types.IosComplete
}

func (disk *DiskDevice) doRead(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(types.IosMediaNotMounted)
		return
	}

	if !disk.IsPrepped() {
		pkt.SetIoStatus(types.IosPackNotPrepped)
		return
	}

	if pkt.buffer == nil {
		pkt.SetIoStatus(types.IosNilBuffer)
		return
	}

	if uint(len(pkt.buffer)) != uint(disk.geometry.PrepFactor) {
		pkt.SetIoStatus(types.IosInvalidBufferSize)
		return
	}

	if uint64(pkt.blockId) >= uint64(disk.geometry.BlockCount) {
		pkt.SetIoStatus(types.IosInvalidBlockId)
		return
	}

	offset := int64(pkt.blockId) * int64(bytesPerBlockMap[disk.geometry.PrepFactor])
	_, err := disk.file.ReadAt(disk.buffer, offset)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(types.IosSystemError)
		return
	}
	pkg.UnpackWord36(disk.buffer, pkt.buffer)

	pkt.ioStatus = types.IosComplete
}

func (disk *DiskDevice) doReadLabel(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(types.IosMediaNotMounted)
		return
	}

	if !disk.IsPrepped() {
		pkt.SetIoStatus(types.IosPackNotPrepped)
		return
	}

	if pkt.buffer == nil {
		pkt.SetIoStatus(types.IosNilBuffer)
		return
	}

	if uint(len(pkt.buffer)) != 28 {
		pkt.SetIoStatus(types.IosInvalidBufferSize)
		return
	}

	_, err := disk.file.ReadAt(disk.buffer, 0)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(types.IosSystemError)
		return
	}
	pkg.UnpackWord36(disk.buffer[:126], pkt.buffer)

	pkt.ioStatus = types.IosComplete
}

// doReset cancels any pending IOs. It is a NOP for us.
func (disk *DiskDevice) doReset(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	// nothing to do for now
	pkt.ioStatus = types.IosComplete
}

func (disk *DiskDevice) doUnmount(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(types.IosMediaNotMounted)
		return
	}

	err := disk.file.Close()
	if err != nil {
		log.Printf("%v\n", err)
	}

	disk.geometry = nil
	disk.file = nil
	pkt.SetIoStatus(types.IosComplete)
}

func (disk *DiskDevice) doWrite(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(types.IosMediaNotMounted)
		return
	}

	if !disk.IsPrepped() {
		pkt.SetIoStatus(types.IosPackNotPrepped)
		return
	}

	if pkt.buffer == nil {
		pkt.SetIoStatus(types.IosNilBuffer)
		return
	}

	if disk.isWriteProtected {
		pkt.SetIoStatus(types.IosWriteProtected)
		return
	}

	if uint(len(pkt.buffer)) != uint(disk.geometry.PrepFactor) {
		pkt.SetIoStatus(types.IosInvalidBufferSize)
		return
	}

	if uint64(pkt.blockId) >= uint64(disk.geometry.BlockCount) {
		pkt.SetIoStatus(types.IosInvalidBlockId)
		return
	}

	pkg.PackWord36(pkt.buffer, disk.buffer)
	offset := int64(pkt.blockId) * int64(bytesPerBlockMap[disk.geometry.PrepFactor])
	_, err := disk.file.WriteAt(disk.buffer, offset)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(types.IosSystemError)
		return
	}

	pkt.ioStatus = types.IosComplete
}

func (disk *DiskDevice) doWriteLabel(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	if !disk.IsMounted() {
		pkt.SetIoStatus(types.IosMediaNotMounted)
		return
	}

	if !disk.IsPrepped() {
		pkt.SetIoStatus(types.IosPackNotPrepped)
		return
	}

	if pkt.buffer == nil {
		pkt.SetIoStatus(types.IosNilBuffer)
		return
	}

	if disk.isWriteProtected {
		pkt.SetIoStatus(types.IosWriteProtected)
		return
	}

	if uint(len(pkt.buffer)) != 28 {
		pkt.SetIoStatus(types.IosInvalidBufferSize)
		return
	}

	pkg.PackWord36(pkt.buffer, disk.buffer)
	_, err := disk.file.WriteAt(disk.buffer, 0)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(types.IosSystemError)
		return
	}

	pkt.ioStatus = types.IosComplete
}

// do this any time we need to read the geometry from a (hopefully) prepped pack.
// we will pretend the prep factor is 28 - this will work for block 0
func (disk *DiskDevice) probeGeometry() error {
	disk.geometry = nil

	label := make([]pkg.Word36, 28)
	err := readBlock(disk.file, 0, label, 28)
	if err != nil {
		log.Printf("Cannot read disk label - assuming pack is not prepped\n")
		return err
	}

	str := label[0].ToStringAsAscii()
	if str != "VOL1" {
		// invalid label - pack is not prepped
		return fmt.Errorf("invalid label (not VOL1)")
	}

	packName := label[1].ToStringAsAscii() + label[2].ToStringAsAscii()[:2]
	if !IsValidPackName(packName) {
		return fmt.Errorf("invalid pack name '%v'", packName)
	}

	prepFactor := types.PrepFactor(label[4].GetH2())
	if !IsValidPrepFactor(prepFactor) {
		return fmt.Errorf("invalid prep factor %v", prepFactor)
	}

	blockCount := types.BlockCount(label[021].GetW())
	blocksPerTrack := uint(1792 / prepFactor)
	trackCount := types.TrackCount(blockCount / types.BlockCount(blocksPerTrack))
	sectorsPerBlock := uint(prepFactor / 28)
	bytesPerBlock := bytesPerBlockMap[prepFactor]

	disk.packName = packName
	disk.geometry = &types.DiskPackGeometry{
		PrepFactor:           prepFactor,
		BlockCount:           blockCount,
		BlocksPerTrack:       blocksPerTrack,
		TrackCount:           trackCount,
		SectorsPerBlock:      sectorsPerBlock,
		BytesPerBlock:        bytesPerBlock,
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

func readBlock(file *os.File, wordAddress uint64, data []pkg.Word36, recordLength uint) error {
	bufLen := recordLength * 9 / 2
	buf := make([]byte, bufLen)
	blockId := wordAddress / uint64(recordLength)
	offset := int64(blockId * uint64(bytesPerBlockMap[types.PrepFactor(recordLength)]))
	//fmt.Printf("readBlock blkId:%v addr:%v\n", blockId, wordAddress) // TODO remove
	_, err := file.ReadAt(buf, offset)
	if err != nil {
		return err
	}
	//dumpBuffer(buf) // TODO remove

	pkg.UnpackWord36(buf, data)
	//pkg.DumpWord36Buffer(data, 7) // TODO remove
	return nil
}

func writeBlock(file *os.File, wordAddress uint64, data []pkg.Word36, recordLength uint) error {
	bufLen := recordLength * 9 / 2
	buf := make([]byte, bufLen)
	blockId := wordAddress / uint64(recordLength)
	offset := int64(blockId * uint64(bytesPerBlockMap[types.PrepFactor(recordLength)]))
	//fmt.Printf("writeBlock blkId:%v addr:%v\n", blockId, wordAddress) // TODO remove
	//pkg.DumpWord36Buffer(data, 7)                                     // TODO remove
	pkg.PackWord36(data, buf)
	//dumpBuffer(buf) // TODO remove
	_, err := file.WriteAt(buf, offset)
	if err != nil {
		return err
	}
	return nil
}

func writeTrack(file *os.File, wordAddress uint64, data []pkg.Word36, recordLength uint) error {
	blocksPerTrack := 1792 / recordLength
	addr := wordAddress
	dx := 0
	for bx := 0; bx < int(blocksPerTrack); bx++ {
		err := writeBlock(file, addr, data[dx:dx+int(recordLength)], recordLength)
		if err != nil {
			return err
		}
		addr += uint64(recordLength)
		dx += int(recordLength)
	}

	return nil
}

func (disk *DiskDevice) Dump(dest io.Writer, indent string) {
	str := fmt.Sprintf("Rdy:%v WProt:%v pack:%v file:%v\n",
		disk.isReady, disk.isWriteProtected, disk.packName, disk.fileName)
	_, _ = fmt.Fprintf(dest, "%v%v", indent, str)

	if disk.geometry != nil {
		str := fmt.Sprintf("  prep:%v trks:%v blks:%v sec/blk:%v blk/trk:%v bytes/blk:%v\n",
			disk.geometry.PrepFactor,
			disk.geometry.TrackCount,
			disk.geometry.BlockCount,
			disk.geometry.SectorsPerBlock,
			disk.geometry.BlocksPerTrack,
			disk.geometry.BytesPerBlock)
		_, _ = fmt.Fprintf(dest, "%v%v", indent, str)
	}
}
