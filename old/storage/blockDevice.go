package storage

import (
	pkg2 "khalehla/old/pkg"
	"khalehla/pkg"
)

// A BlockGeometry struct allows a particular block device to report the included information
// regarding the actual storage medium which it controls. The implementation must be able to
// derive the information from whatever underlying storage medium it uses, such that the
// client needs only construct the medium once, providing the necessary parameters, and then
// can cause the medium to be opened and read without knowing this information ahead of time.
// Generally, the device stores this information in some convenient format in block zero of
// the storage - and would then reject any attempt to overwrite block 0.
type BlockGeometry struct {
	bytesPerBlock  pkg.BlockSize
	wordsPerBlock  pkg.BlockSize
	blocksPerTrack pkg.BlockCount
	blockCount     pkg.BlockCount
	label          string // always 8 characters, LJSF
}

// A BlockDevice is used for storing fixed-size blocked data.
// The data must be read and written in fixed block sizes.
type BlockDevice interface {
	// AllocateBlocks causes the device to allocate space without writing to it, or to write zeroes,
	// according to the particularities of the implementation. Any indicated block which is already
	// allocated is left alone.
	// Some devices may not be able to effectively allocate or release blocks.
	// For those devices, this method is a no-op.
	// blockId is the id of the first block to be allocated.
	// blockCount is the number of consecutive blocks to be allocated.
	AllocateBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount) DeviceResult

	// Close causes the device to release the storage. For temporary devices, the storage is lost.
	Close() DeviceResult

	// GetDeviceType retrieves the device type from any device even if it is not open
	GetDeviceType() pkg.DeviceType

	// GetGeometry retrieves the geometry of the medium managed by the device.
	// The device must be open.
	GetGeometry() (BlockGeometry, DeviceResult)

	// IsOpen indicates whether the device has been opened
	IsOpen() bool

	// IsWriteProtected indicates whether write, allocate, and release operations are prohibited
	IsWriteProtected() bool

	// Open prepares the device for I/O.
	// writeProtected indicates whether write, allocate, and release operations should be prohibited
	// writeThrough indicates that, where appropriate, all writes should immediately be persisted to the
	// underlying medium.
	Open(writeProtected bool, writeThrough bool) DeviceResult

	// ReadBlocks reads one or more blocks of Word32 structs.
	// blockId is the id of the first block to be read.
	// blockCount is the number of consecutive blocks to be read.
	// buffer is a buffer of Word36 structs, exactly large enough to contain the given number of blocks.
	ReadBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount, buffer []pkg2.Word36) DeviceResult

	// ReleaseBlocks releases one or more blocks of storage.
	// Any block which is not allocated is left alone.
	// Some devices may not be able to effectively allocate or release blocks.
	// For those devices, this method is a no-op.
	// blockId is the id of the first block to be released.
	// blockCount is the number of consecutive blocks to be released.
	ReleaseBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount) DeviceResult

	// WriteBlocks writes one or more blocks of Word32 structs.
	// blockId is the id of the first block to be written.
	// blockCount is the number of consecutive blocks to be written.
	// buffer is a buffer of Word36 structs, exactly large enough to contain the given number of blocks.
	WriteBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount, buffer []pkg2.Word36) DeviceResult
}

func IsLabelValid(label string) bool {
	//	TODO
	return true
}

func IsWordsPerBlockValid(wordsPerBlock pkg.BlockSize) bool {
	_, ok := pkg.BlockSizeFromPrepFactor[wordsPerBlock]
	return ok
}
