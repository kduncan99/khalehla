package storage

import (
	"khalehla/pkg"
	"os"
	"unsafe"
)

type fileBlockZero struct {
	ident         pkg.Word36
	label         pkg.Word36
	wordsPerBlock pkg.Word36
	blockCount    pkg.Word36
}

var fileBlockZeroSize = pkg.RawBytesPerBlockFromWords[4]

const fileIdentifierConstant = "BLKDVF"

// A FileBlockDevice persists data to an underlying system file.
// All data is written in contiguous blocks (but with random access) where the blocks are in order by block id.
// There is considerable waste, as we persist the Word36 objects (which have 28 bits of slop per word)
// as 8-byte entities. We do NOT pad the blocks out to the next 4k physical block, so that might be an issue.
type FileBlockDevice struct {
	fileName       string
	geometry       BlockGeometry
	file           *os.File
	writeProtected bool
	writeThrough   bool
}

func (bd *FileBlockDevice) AllocateBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount) DeviceResult {
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

func (bd *FileBlockDevice) Close() DeviceResult {
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

func (bd *FileBlockDevice) GetDeviceType() pkg.DeviceType {
	return DeviceTypeFileBlock
}

func (bd *FileBlockDevice) GetGeometry() (BlockGeometry, DeviceResult) {
	if !bd.IsOpen() {
		return BlockGeometry{}, DeviceResult{DeviceStatusNotOpen, nil}
	}

	return bd.geometry, DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *FileBlockDevice) IsOpen() bool {
	return bd.file != nil
}

func (bd *FileBlockDevice) IsWriteProtected() bool {
	return bd.writeProtected
}

func (bd *FileBlockDevice) Open(writeProtected bool, writeThrough bool) DeviceResult {
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

	bz := fileBlockZero{}
	_, err = f.ReadAt(unsafe.Slice((*byte)(unsafe.Pointer(&bz)), fileBlockZeroSize), 0)
	if err != nil {
		_ = f.Close()
		return DeviceResult{DeviceStatusSystemError, err}
	}

	ident := bz.ident.ToStringAsFieldata()
	if ident != fileIdentifierConstant {
		return DeviceResult{DeviceStatusInvalidIdentifierConstant, nil}
	}

	bd.geometry.label = bz.label.ToStringAsFieldata()
	bd.geometry.blockCount = pkg.BlockCount(bz.blockCount.GetW())
	bd.geometry.wordsPerBlock = pkg.BlockSize(bz.wordsPerBlock.GetW())
	bd.geometry.bytesPerBlock = pkg.RawBytesPerBlockFromWords[bd.geometry.wordsPerBlock]
	bd.geometry.blocksPerTrack = pkg.BlockCount(1792 / bd.geometry.wordsPerBlock)

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *FileBlockDevice) ReadBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount, buffer []pkg.Word36) DeviceResult {
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

	pos := int64(blockId) * int64(bd.geometry.wordsPerBlock)
	_, err := bd.file.ReadAt(unsafe.Slice((*byte)(unsafe.Pointer(&buffer[0])), bd.geometry.bytesPerBlock), pos)
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

// ReleaseBlocks releases the indicated blocks by truncating the file, but ONLY if the indicated extent reaches
// or surpasses the physical end of the file (so that we do not remove data which follows the indicated extent),
// and ONLY if the beginning of the indicated area is less than the physical end of the file (so that the truncate
// operation does not add space which wasn't there to begin with)
func (bd *FileBlockDevice) ReleaseBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount) DeviceResult {
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

func (bd *FileBlockDevice) WriteBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount, buffer []pkg.Word36) DeviceResult {
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

	pos := int64(blockId) * int64(bd.geometry.wordsPerBlock)
	_, err := bd.file.WriteAt(unsafe.Slice((*byte)(unsafe.Pointer(&buffer[0])), bd.geometry.bytesPerBlock), pos)
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func CreateFileBlockDevice(fileName string, label string, wordsPerBlock pkg.BlockSize, blockCount pkg.BlockCount, preallocate bool) DeviceResult {
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
		bytesPerBlock:  pkg.RawBytesPerBlockFromWords[wordsPerBlock],
		wordsPerBlock:  wordsPerBlock,
	}

	//	Write the block zero content
	ident := []byte(fileIdentifierConstant)
	wIdent := pkg.Word36(0)
	wIdent.FromStringToFieldata(ident)

	wLabel := pkg.Word36(0)
	wLabel.FromStringToFieldata([]byte(label + "     "))

	bz := fileBlockZero{
		ident:         wIdent,
		label:         wLabel,
		wordsPerBlock: pkg.Word36(wordsPerBlock),
		blockCount:    pkg.Word36(blockCount),
	}

	f, err := os.OpenFile(bd.fileName, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return DeviceResult{DeviceStatusSystemError, err}
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = f.WriteAt(unsafe.Slice((*byte)(unsafe.Pointer(&bz)), fileBlockZeroSize), 0)
	if err != nil {
		_ = f.Close()
		return DeviceResult{DeviceStatusSystemError, err}
	}

	//	For a standard system block, it should be sufficient to write the last byte of the last block.
	//  This *should* cause allocation of all the bytes from the beginning of the file to the end.
	if preallocate {
		offset := int64(int64(blockCount)*int64(pkg.RawBytesPerBlockFromWords[wordsPerBlock])) - 1
		b := []byte{0}
		_, err = f.WriteAt(b, offset)
		if err != nil {
			_ = f.Close()
			return DeviceResult{DeviceStatusSystemError, err}
		}
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func NewFileBlockDevice(fileName string) *FileBlockDevice {
	return &FileBlockDevice{
		file:           nil,
		fileName:       fileName,
		geometry:       BlockGeometry{},
		writeProtected: true,
		writeThrough:   false,
	}
}
