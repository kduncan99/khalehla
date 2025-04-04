package storage

import (
	"os"
	"sync"
	"unsafe"

	"khalehla/common"
	pkg2 "khalehla/old/pkg"
	"khalehla/pkg"
)

type packedBlockZero struct {
	ident         pkg2.Word36
	label         pkg2.Word36
	wordsPerBlock pkg2.Word36
	blockCount    pkg2.Word36
}

var packedBlockZeroSize = common.RawBytesPerBlock[4]

const packedIdentifierConstant = "BLKDVP"

// TODO implement varying number of midBuffer structs, so that we can have multiple IOs in progress concurrently
//	this would only limit Open(), ReadBlocks(), and WriteBlocks()

// A PackedBlockDevice persists data to an underlying system file, packing 2 words to 9 bytes.
// All data is written in contiguous blocks (but with random access) where the blocks are in order by block id.
// Not much different than the FileBlockDevice, this one will be a little slower due to packing/unpacking,
// but that can be mitigated with a cache aggregator. It will save roughly 43% of storage footprint.
type PackedBlockDevice struct {
	fileName       string
	geometry       BlockGeometry
	file           *os.File
	writeProtected bool
	writeThrough   bool
	mutex          sync.Mutex // Protects midBuffer
	midBuffer      []byte
}

func (bd *PackedBlockDevice) AllocateBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount) DeviceResult {
	bd.mutex.Lock()
	defer bd.mutex.Unlock()

	if !bd.IsOpen() {
		return DeviceResult{DeviceStatusNotOpen, nil}
	}

	if bd.IsWriteProtected() {
		return DeviceResult{DeviceStatusWriteProtected, nil}
	}

	if int(blockId) >= int(bd.geometry.blockCount) {
		return DeviceResult{DeviceStatusInvalidBlockId, nil}
	}

	if int(blockId)+int(blockCount) > int(bd.geometry.blockCount) {
		return DeviceResult{DeviceStatusMaxBlocksExceeded, nil}
	}

	limitBlockId := int64(blockId) + int64(blockCount)
	limitOffset := limitBlockId * int64(bd.geometry.bytesPerBlock)

	fi, err := bd.file.Stat()
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	if fi.Size() < limitOffset {
		err := bd.file.Truncate(limitOffset)
		if err != nil {
			return DeviceResult{DeviceStatusSystemError, err}
		}
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *PackedBlockDevice) Close() DeviceResult {
	bd.mutex.Lock()
	defer bd.mutex.Unlock()

	if !bd.IsOpen() {
		return DeviceResult{DeviceStatusNotOpen, nil}
	}

	err := bd.file.Close()
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	bd.file = nil
	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *PackedBlockDevice) GetDeviceType() pkg.DeviceType {
	return DeviceTypeFileBlock
}

func (bd *PackedBlockDevice) GetGeometry() (BlockGeometry, DeviceResult) {
	if !bd.IsOpen() {
		return BlockGeometry{}, DeviceResult{DeviceStatusNotOpen, nil}
	}

	return bd.geometry, DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *PackedBlockDevice) IsOpen() bool {
	return bd.file != nil
}

func (bd *PackedBlockDevice) IsWriteProtected() bool {
	return bd.writeProtected
}

func (bd *PackedBlockDevice) Open(writeProtect bool, writeThrough bool) DeviceResult {
	if bd.IsOpen() {
		return DeviceResult{DeviceStatusAlreadyOpen, nil}
	}

	openFlag := os.O_CREATE
	if bd.writeProtected {
		openFlag |= os.O_RDONLY
	} else {
		openFlag |= os.O_RDWR
	}
	if bd.writeThrough {
		openFlag |= os.O_SYNC
	}

	var err error
	bd.file, err = os.OpenFile(bd.fileName, openFlag, 0755)
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	// Read geometry
	f, err := os.OpenFile(bd.fileName, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	bb := make([]byte, packedBlockZeroSize)
	_, err = f.ReadAt(bb, 0)
	if err != nil {
		_ = f.Close()
		return DeviceResult{DeviceStatusSystemError, err}
	}

	bz := packedBlockZero{}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&bz)), packedBlockZeroSize), bb)
	ident := bz.ident.ToStringAsFieldata()
	if ident != packedIdentifierConstant {
		return DeviceResult{DeviceStatusInvalidIdentifierConstant, nil}
	}

	bd.geometry.label = bz.label.ToStringAsFieldata()
	bd.geometry.blockCount = pkg.BlockCount(bz.blockCount.GetW())
	bd.geometry.wordsPerBlock = pkg.BlockSize(bz.wordsPerBlock.GetW())
	bd.geometry.bytesPerBlock = pkg.BlockSizeFromPrepFactor[bd.geometry.wordsPerBlock]
	bd.geometry.blocksPerTrack = pkg.BlockCount(1792 / bd.geometry.wordsPerBlock)

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *PackedBlockDevice) readBlock(blockId pkg.BlockId, buffer []pkg2.Word36) error {
	var err error
	pos := int64(blockId) * int64(bd.geometry.bytesPerBlock)
	_, err = bd.file.ReadAt(bd.midBuffer, pos)
	if err == nil {
		common.UnpackWord36Strict(bd.midBuffer, buffer)
	}
	return err
}

func (bd *PackedBlockDevice) ReadBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount, buffer []pkg2.Word36) DeviceResult {
	bd.mutex.Lock()
	defer bd.mutex.Unlock()

	if !bd.IsOpen() {
		return DeviceResult{DeviceStatusNotOpen, nil}
	}

	if len(buffer) != int(blockCount)*int(bd.geometry.wordsPerBlock) {
		return DeviceResult{DeviceStatusInvalidBufferSize, nil}
	}

	if int(blockId) >= int(bd.geometry.blockCount) {
		return DeviceResult{DeviceStatusInvalidBlockId, nil}
	}

	if int(blockId)+int(blockCount) > int(bd.geometry.blockCount) {
		return DeviceResult{DeviceStatusMaxBlocksExceeded, nil}
	}

	bid := blockId
	bx := 0
	for bc := pkg.BlockCount(0); bc < blockCount; bc++ {
		err := bd.readBlock(bid, buffer[bx:bx+int(bd.geometry.wordsPerBlock)])
		if err != nil {
			return DeviceResult{DeviceStatusSystemError, err}
		}
		bid += 1
		bx += int(bd.geometry.wordsPerBlock)
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *PackedBlockDevice) ReleaseBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount) DeviceResult {
	bd.mutex.Lock()
	defer bd.mutex.Unlock()

	if !bd.IsOpen() {
		return DeviceResult{DeviceStatusNotOpen, nil}
	}

	if bd.IsWriteProtected() {
		return DeviceResult{DeviceStatusWriteProtected, nil}
	}

	if int(blockId) >= int(bd.geometry.blockCount) {
		return DeviceResult{DeviceStatusInvalidBlockId, nil}
	}

	if int(blockId)+int(blockCount) > int(bd.geometry.blockCount) {
		return DeviceResult{DeviceStatusMaxBlocksExceeded, nil}
	}

	firstOffset := int64(blockId) * int64(bd.geometry.bytesPerBlock)
	limitOffset := firstOffset + (int64(blockCount) * int64(bd.geometry.bytesPerBlock))

	fi, err := bd.file.Stat()
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	if (firstOffset < fi.Size()) && (fi.Size() < limitOffset) {
		err := bd.file.Truncate(limitOffset)
		if err != nil {
			return DeviceResult{DeviceStatusSystemError, err}
		}
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *PackedBlockDevice) writeBlock(blockId pkg.BlockId, buffer []pkg2.Word36) error {
	var err error
	pos := int64(blockId) * int64(bd.geometry.bytesPerBlock)
	_, err = bd.file.ReadAt(bd.midBuffer, pos)
	if err == nil {
		common.UnpackWord36Strict(bd.midBuffer, buffer)
	}
	return err
}

func (bd *PackedBlockDevice) WriteBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount, buffer []pkg2.Word36) DeviceResult {
	bd.mutex.Lock()
	defer bd.mutex.Unlock()

	if !bd.IsOpen() {
		return DeviceResult{DeviceStatusNotOpen, nil}
	}

	if bd.IsWriteProtected() {
		return DeviceResult{DeviceStatusWriteProtected, nil}
	}

	if len(buffer) != int(blockCount)*int(bd.geometry.wordsPerBlock) {
		return DeviceResult{DeviceStatusInvalidBufferSize, nil}
	}

	if int(blockId) >= int(bd.geometry.blockCount) {
		return DeviceResult{DeviceStatusInvalidBlockId, nil}
	}

	if int(blockId)+int(blockCount) > int(bd.geometry.blockCount) {
		return DeviceResult{DeviceStatusMaxBlocksExceeded, nil}
	}

	bid := blockId
	bx := 0
	for bc := pkg.BlockCount(0); bc < blockCount; bc++ {
		err := bd.readBlock(bid, buffer[bx:bx+int(bd.geometry.wordsPerBlock)])
		if err != nil {
			return DeviceResult{DeviceStatusSystemError, err}
		}
		bid += 1
		bx += int(bd.geometry.wordsPerBlock)
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func CreatePackedBlockDevice(fileName string, label string, wordsPerBlock pkg.BlockSize, blockCount pkg.BlockCount, preallocate bool) DeviceResult {
	if !IsLabelValid(label) {
		return DeviceResult{DeviceStatusInvalidLabel, nil}
	}

	if !IsWordsPerBlockValid(wordsPerBlock) {
		return DeviceResult{DeviceStatusInvalidBlockSize, nil}
	}

	bd := NewFileBlockDevice(fileName)
	var err error
	bd.file, err = os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(bd.file)

	bd.geometry = BlockGeometry{
		blocksPerTrack: pkg.BlockCount(1792 / wordsPerBlock),
		blockCount:     blockCount,
		bytesPerBlock:  pkg.BlockSizeFromPrepFactor[wordsPerBlock],
		wordsPerBlock:  wordsPerBlock,
	}

	//	Write the block zero content
	ident := []byte(fileIdentifierConstant)
	wIdent := pkg2.Word36(0)
	wIdent.FromStringToFieldata(ident)

	wLabel := pkg2.Word36(0)
	wLabel.FromStringToFieldata([]byte(label + "     "))

	bz := fileBlockZero{
		ident:         wIdent,
		label:         wLabel,
		wordsPerBlock: pkg2.Word36(wordsPerBlock),
		blockCount:    pkg2.Word36(blockCount),
	}

	f, err := os.OpenFile(bd.fileName, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	bb := make([]byte, packedBlockZeroSize)
	copy(bb, unsafe.Slice((*byte)(unsafe.Pointer(&bz)), fileBlockZeroSize))
	_, err = f.WriteAt(bb, 0)
	if err != nil {
		_ = f.Close()
		return DeviceResult{DeviceStatusSystemError, err}
	}

	//	For a standard system block, it should be sufficient to write the last block.
	//  This *should* cause allocation of all the bytes from the beginning of the file to the end.
	if preallocate {
		offset := int64(int64(blockCount)*int64(pkg.BlockSizeFromPrepFactor[wordsPerBlock])) - 1
		b := []byte{0}
		_, err = f.WriteAt(b, offset)
		if err != nil {
			_ = f.Close()
			return DeviceResult{DeviceStatusSystemError, err}
		}
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func NewPackedBlockDevice(fileName string) *FileBlockDevice {
	return &FileBlockDevice{
		file:           nil,
		fileName:       fileName,
		geometry:       BlockGeometry{},
		writeProtected: true,
		writeThrough:   false,
	}
}
