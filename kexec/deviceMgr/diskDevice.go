// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import (
	"fmt"
	"khalehla/pkg"
	"log"
	"math"
	"os"
	"sync"
)

// This is a very simple pseudo disk device
// We depart from the conventional disk layout in the following manner:
//   There is no booting from disk, so there is no bootstrap in physical blocks 0 or 1
//   The label is always located in physical block 0 instead of physical block 2
//   The first directory track is always located at the second track-aligned physical block
// This means that we can always determine whether a pseudo-pack is prepped, and at what prep factor,
// and how many tracks it holds, by simply reading the label from the first block.
// Even though the label is only 28 words, we're okay because we don't actually *have* to read a whole block.

// VOL1 disk label - canonically the third physical record of the device, but always the first block for us
// +000     "VOL1" - ASCII
// +001     pack-id - ASCII LJSF
// +002,H1  pack-id continued - ASCII LJSF
// +003     sector address of first directory track
// +004,H1  records per track (1792 / prep_factor)
// +004,H2  words per record (prep_factor)
// +005,H2  reserved size in tracks (for DRS) - we don't do DRS
// +011,H1  Length of S0 + S1 + HMBT + pad to the next physical record boundary in words
// +011,H2  Master bit table length (HMBT, SMBT) in words
// +014,S1  020
// +014,S2  Vol1 Version (we use 1)
// +014,H2  Heads per cylinder
// +016     Disk capacity in tracks, not including label or initial directory allocation
// +017,H1  Words per physical record (also prep_factor)
// +020,H2  Attributes - we just set these all to zero
// +021     (non-canonical) total number of blocks on pack

// Simple lookup table - the key is words per block, the value is bytes per block padded to power of two
var bytesPerBlockMap = map[PrepFactor]uint{
	28:   128,
	56:   256,
	112:  512,
	224:  1024,
	448:  2048,
	896:  4096,
	1792: 8192,
}

type DiskDevice struct {
	fileName     *string
	file         *os.File
	writeProtect bool
	packName     string
	geometry     *DiskPackGeometry
	mutex        sync.Mutex
	buffer       []byte
}

func NewDiskDevice(initialFileName *string) *DiskDevice {
	return &DiskDevice{
		fileName: initialFileName,
	}
}

func (disk *DiskDevice) getNodeType() NodeType {
	return NodeTypeDisk
}

func (disk *DiskDevice) GetGeometry() *DiskPackGeometry {
	return disk.geometry
}

func (disk *DiskDevice) IsMounted() bool {
	return disk.file != nil
}

func (disk *DiskDevice) IsPrepped() bool {
	return disk.geometry != nil
}

func (disk *DiskDevice) startIo(pkt IoPacket) {
	pkt.SetIoStatus(IosInProgress)

	if pkt.GetNodeType() != disk.getNodeType() {
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
	default:
		pkt.SetIoStatus(IosInvalidFunction)
	}
}

func (disk *DiskDevice) doMount(pkt *DiskIoPacket) {
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

	// At this point, the pack is now mounted. It may not be prepped, but that is okay.
	disk.file = f
	disk.writeProtect = pkt.writeProtected
	pkt.SetIoStatus(IosComplete)

	err = disk.probeGeometry()
	if err != nil {
		log.Printf("%v\n", err)
	}
}

func (disk *DiskDevice) doPrep(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()
	if !disk.IsMounted() {
		pkt.SetIoStatus(IosMediaNotMounted)
		return
	}

	if !IsValidPrepFactor(pkt.prepFactor) {
		pkt.SetIoStatus(IosInvalidPrepFactor)
		return
	}

	if pkt.trackCount < 10000 {
		pkt.SetIoStatus(IosInvalidTrackCount)
		return
	}

	if !IsValidPackName(pkt.packName) {
		pkt.SetIoStatus(IosInvalidPackName)
		return
	}

	// basic geometry - some of these values exist simply so we do not have to constantly cast them to do math
	recordLength := uint64(pkt.prepFactor)
	trackCount := uint64(pkt.trackCount)
	blocksPerTrack := 1792 / recordLength
	blockCount := trackCount * blocksPerTrack
	dirTrackAddr := uint64(1792) // we set this to the device-relative word address of the initial directory track

	rawMBTWordCount := trackCount / 32
	if trackCount%32 > 0 {
		rawMBTWordCount++
	}
	mbtWordCount := rawMBTWordCount + 2 // includes control word and checksum

	paddedMBTWordCount := 28 + 28 + mbtWordCount // s0 + s1 + HMBT + pad to physical record boundary
	mod := paddedMBTWordCount % recordLength
	paddedMBTWordCount += recordLength - mod

	// how many tracks do we need for the initial directory? this includes s0, s1, HMBT, and SMBT
	initWords := paddedMBTWordCount + mbtWordCount
	initTracks := initWords / 1792
	mod = initWords % 1792
	if mod > 0 {
		initTracks++
	}

	// Where are the HMBT and SMBT?
	hmbtAddress := uint64(1792 + 56)
	smbtAddress := 1792 + paddedMBTWordCount

	// Where is the first DAS? It is NOT in the initial directory track...
	// Rather, it is the first directory track following the MBT nonsense.
	// We will place it at the next software track following the MBT initial stuffs.
	firstDasTrackOffset := 1 + initTracks
	firstDasWordOffset := firstDasTrackOffset * 1792

	// Available tracks - total tracks, minus VOL1 space, minus initial tracks, minus first DAS directory track.
	availableTracks := trackCount - 1 - initTracks - 1

	// DAS sector offset
	// This is the sector (28-word) offset into the directory for this pack, which contains the first DAS.
	// It must follow all initial directory tracks, and be on a 9-word boundary
	dasOffset := initTracks
	if dasOffset%9 > 0 {
		dasOffset += 9 - (dasOffset % 9)
	}
	dasOffset *= 64

	// create initial label and write it
	label := make([]pkg.Word36, recordLength)
	pkg.FromStringToAscii("VOL1", label[0:1])
	pkg.FromStringToAscii(pkt.packName, label[1:3])
	label[2].SetH2(0)
	label[3].SetW(dirTrackAddr)
	label[4].SetH1(blocksPerTrack)
	label[4].SetH2(recordLength)
	label[5].SetW(0) // no DRS tracks
	label[011].SetH1(paddedMBTWordCount)
	label[011].SetH2(mbtWordCount)
	label[014].SetS1(010) // Pretend we are DPREP
	label[014].SetS2(1)   // VOL1 version
	label[014].SetH2(10)  // heads per cylinder - make up something
	label[016].SetW(availableTracks)
	label[017].SetH1(recordLength)
	label[021].SetW(blockCount)

	err := writeBlock(disk.file, 0, label, uint(recordLength))
	if err != nil {
		log.Printf("Error writing label:%v\n", err)
		pkt.SetIoStatus(IosSystemError)
		return
	}

	// initial directory
	initDir := make([]pkg.Word36, initTracks*1792)

	// sector 0
	s0 := initDir[0:28]
	s0[0].SetS1(initTracks)         // number of software tracks allocated according to the first DAS
	s0[27].SetW(firstDasWordOffset) // device-relative word address of the first DAS, which is the next directory track

	// sector 1
	s1 := initDir[28:56]
	s1[0].SetW(hmbtAddress) // device-relative address of hardware MBT
	s1[1].SetW(smbtAddress) // device-relative address of software MBT
	s1[2].SetW(availableTracks)
	s1[3].SetW(availableTracks)
	s1[4].FromStringToFieldata(disk.packName)
	if !pkt.removable {
		s1[5].SetS1(040)
	}
	s1[5].SetH2(mbtWordCount)
	s1[7].SetW(0) // Normally we'd set this to timestamp, but it is TDATE$, and we don't want to do that
	s1[010].SetT1(blocksPerTrack)
	s1[010].SetS3(1) // Sector 1 version
	s1[010].SetT3(recordLength)
	s1[020].SetH2(dasOffset) // normally set by EXEC...

	// populate the HMBT, accounting for the VOL1 space, the initial track allocation, and the first DAS track
	hmbt := initDir[56 : 56+mbtWordCount]
	counter := (1 + initTracks + 1) * 2 // counter is number of half-tracks
	for hx := 1; hx < int(mbtWordCount)-1; hx++ {
		if counter > 32 {
			hmbt[hx].SetW(0777777_777777)
			counter -= 32
		} else if counter > 0 {
			value := uint64(0)
			for counter > 0 {
				value = (value >> 1) | 0_400000_000000
				counter--
			}
			hmbt[hx].SetW(value | 017)
		} else {
			hmbt[hx].SetW(017)
		}
	}

	// mark anything off the end as allocated so kexec doesn't try to use tracks which don't exist
	mod = trackCount % 32
	if mod > 0 {
		val := (uint64(math.Pow(2, float64(mod))) - 1) << 4
		hmbt[mbtWordCount-2].Or(val)
	}
	hmbt[mbtWordCount-1].FromStringToFieldata("NOCKSM")

	// copy the HMBT to the SMBT
	smbt := initDir[paddedMBTWordCount : paddedMBTWordCount+mbtWordCount]
	for sx := 0; sx < int(mbtWordCount); sx++ {
		smbt[sx] = hmbt[sx]
	}

	// Write the initial directory tracks
	addr := dirTrackAddr
	wx := 0
	for ix := 0; ix < int(initTracks); ix++ {
		err = writeTrack(disk.file, addr, initDir[wx:wx+1792], uint(recordLength))
		if err != nil {
			log.Printf("Error writing label:%v\n", err)
			pkt.SetIoStatus(IosSystemError)
			return
		}
		addr += 1792
	}

	// Create and write the initial DAS track
	dasTrack := make([]pkg.Word36, 1792)
	das := dasTrack[0:28]
	das[0].SetW(firstDasWordOffset)
	if pkt.removable {
		das[0].SetS1(02)
	} else {
		das[0].SetS1(03)
	}
	for dx := 3; dx < 27; dx += 3 {
		das[dx].SetW(0_400000_000000)
	}
	das[27].SetW(0_400000_000000)

	err = writeTrack(disk.file, firstDasWordOffset, dasTrack, uint(recordLength))
	if err != nil {
		log.Printf("Error writing label:%v\n", err)
		pkt.SetIoStatus(IosSystemError)
		return
	}

	err = disk.probeGeometry()
	if err != nil {
		log.Printf("%v\n", err)
	}

	pkt.ioStatus = IosComplete
}

func (disk *DiskDevice) doRead(pkt *DiskIoPacket) {
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

	offset := int64(pkt.blockId) * int64(disk.geometry.BytesPerBlock)
	_, err := disk.file.ReadAt(disk.buffer, offset)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(IosSystemError)
		return
	}
	pkg.UnpackWord36(disk.buffer, pkt.buffer)

	pkt.ioStatus = IosComplete
}

func (disk *DiskDevice) doReadLabel(pkt *DiskIoPacket) {
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
func (disk *DiskDevice) doReset(pkt *DiskIoPacket) {
	disk.mutex.Lock()
	defer disk.mutex.Unlock()

	// nothing to do for now
	pkt.ioStatus = IosComplete
}

func (disk *DiskDevice) doUnmount(pkt *DiskIoPacket) {
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

func (disk *DiskDevice) doWrite(pkt *DiskIoPacket) {
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

	if disk.writeProtect {
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

	offset := int64(pkt.blockId) * int64(disk.geometry.BytesPerBlock)
	_, err := disk.file.ReadAt(disk.buffer, offset)
	if err != nil {
		log.Printf("%v\n", err)
		pkt.SetIoStatus(IosSystemError)
		return
	}
	pkg.UnpackWord36(disk.buffer, pkt.buffer)

	pkt.ioStatus = IosComplete
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

	prepFactor := PrepFactor(label[4].GetH2())
	if !IsValidPrepFactor(prepFactor) {
		return fmt.Errorf("invalid prep factor %v", prepFactor)
	}

	blockCount := BlockCount(label[021].GetW())
	blocksPerTrack := uint(1792 / prepFactor)
	trackCount := TrackCount(blockCount / BlockCount(blocksPerTrack))
	sectorsPerBlock := uint(prepFactor / 28)
	bytesPerBlock := bytesPerBlockMap[prepFactor]

	disk.packName = packName
	disk.geometry = &DiskPackGeometry{
		PrepFactor:      prepFactor,
		BlockCount:      blockCount,
		BlocksPerTrack:  blocksPerTrack,
		TrackCount:      trackCount,
		SectorsPerBlock: sectorsPerBlock,
		BytesPerBlock:   bytesPerBlock,
	}
	disk.buffer = make([]byte, bytesPerBlockMap[prepFactor])

	return nil
}

func dump(buffer []byte) {
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
	buf := make([]byte, recordLength*9/2)
	blockId := wordAddress / uint64(recordLength)
	//	fmt.Printf("readBlock blkId:%v addr:%v\n", blockId, wordAddress) // TODO remove
	_, err := file.ReadAt(buf, int64(blockId*uint64(recordLength)))
	if err != nil {
		return err
	}
	//	dump(buf) // TODO remove

	pkg.UnpackWord36(buf, data)
	//	pkg.DumpWord36Buffer(data, 7) // TODO remove
	return nil
}

func writeBlock(file *os.File, wordAddress uint64, data []pkg.Word36, recordLength uint) error {
	buf := make([]byte, recordLength*9/2)
	blockId := wordAddress / uint64(recordLength)
	//	fmt.Printf("writeBlock blkId:%v addr:%v\n", blockId, wordAddress) // TODO remove
	//	pkg.DumpWord36Buffer(data, 7)                                     // TODO remove
	pkg.PackWord36(data, buf)
	//	dump(buf) // TODO remove
	_, err := file.WriteAt(buf, int64(blockId*uint64(recordLength)))
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
