package storage

import (
	"khalehla/pkg"
)

// TemporaryBlockDevice implements BlockDevice, storing logical blocks in memory.
// All storage is lost when the device is closed
type TemporaryBlockDevice struct {
	geometry BlockGeometry
	storage  map[pkg.BlockId][]pkg.Word36 // index is logical block id, value is the actual block of data
	isOpen   bool
}

func (bd *TemporaryBlockDevice) AllocateBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount) DeviceResult {
	bid := blockId
	for bx := pkg.BlockCount(0); bx < blockCount; bx++ {
		_, ok := bd.storage[bid]
		if !ok {
			bd.storage[bid] = make([]pkg.Word36, bd.geometry.wordsPerBlock)
		}
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *TemporaryBlockDevice) Close() DeviceResult {
	if !bd.IsOpen() {
		return DeviceResult{DeviceStatusNotOpen, nil}
	}

	bd.isOpen = false
	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *TemporaryBlockDevice) GetDeviceType() pkg.DeviceType {
	return DeviceTypeTemporaryBlock
}

func (bd *TemporaryBlockDevice) GetGeometry() (BlockGeometry, DeviceResult) {
	return bd.geometry, DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *TemporaryBlockDevice) IsOpen() bool {
	return bd.isOpen
}

func (bd *TemporaryBlockDevice) IsWriteProtected() bool {
	return false
}

func (bd *TemporaryBlockDevice) Open(writeProtected bool, writeThrough bool) DeviceResult {
	// writeThrough is ignored - we are effectively always and never write-through.
	if bd.IsOpen() {
		return DeviceResult{DeviceStatusAlreadyOpen, nil}
	}

	if writeProtected {
		return DeviceResult{DeviceStatusCannotSetWriteProtect, nil}
	}

	bd.storage = make(map[pkg.BlockId][]pkg.Word36)
	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *TemporaryBlockDevice) readBlock(blockId pkg.BlockId, buffer []pkg.Word36) {
	storage, ok := bd.storage[blockId]
	if ok {
		copy(buffer, storage)
	} else {
		for bx := 0; bx < len(buffer); bx++ {
			buffer[bx] = 0
		}
	}
}

func (bd *TemporaryBlockDevice) ReadBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount, buffer []pkg.Word36) DeviceResult {
	if len(buffer) != int(blockCount)*int(bd.geometry.wordsPerBlock) {
		return DeviceResult{DeviceStatusInvalidBufferSize, nil}
	}

	bid := blockId
	bx := 0
	for bc := pkg.BlockCount(0); bc < blockCount; bc++ {
		bd.readBlock(bid, buffer[bx:bx+int(bd.geometry.wordsPerBlock)])
		bid += 1
		bx += int(bd.geometry.wordsPerBlock)
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *TemporaryBlockDevice) ReleaseBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount) DeviceResult {
	bid := blockId
	for bx := 0; bx < int(blockCount); bx++ {
		delete(bd.storage, bid)
		bid++
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

func (bd *TemporaryBlockDevice) writeBlock(blockId pkg.BlockId, buffer []pkg.Word36) {
	storage, ok := bd.storage[blockId]
	if !ok {
		storage = make([]pkg.Word36, bd.geometry.wordsPerBlock)
		bd.storage[blockId] = storage
	}

	copy(storage, buffer)
}

func (bd *TemporaryBlockDevice) WriteBlocks(blockId pkg.BlockId, blockCount pkg.BlockCount, buffer []pkg.Word36) DeviceResult {
	if len(buffer) != int(blockCount)*int(bd.geometry.wordsPerBlock) {
		return DeviceResult{DeviceStatusInvalidBufferSize, nil}
	}

	if bd.IsWriteProtected() {
		return DeviceResult{DeviceStatusWriteProtected, nil}
	}

	bid := blockId
	bx := 0
	for bc := pkg.BlockCount(0); bc < blockCount; bc++ {
		bd.writeBlock(bid, buffer[bx:bx+int(bd.geometry.wordsPerBlock)])
		bid += 1
		bx += int(bd.geometry.wordsPerBlock)
	}

	return DeviceResult{DeviceStatusSuccessful, nil}
}

// NewTemporaryBlockDevice creates an in-memory block device.
// wordsPerBlock is always 1792.
func NewTemporaryBlockDevice(label string, blockCount pkg.BlockCount) (*TemporaryBlockDevice, DeviceResult) {
	if !IsLabelValid(label) {
		return nil, DeviceResult{DeviceStatusInvalidLabel, nil}
	}

	g := BlockGeometry{
		blockCount:     blockCount,
		blocksPerTrack: 1,
		bytesPerBlock:  1792 * 8,
		label:          label,
		wordsPerBlock:  1792,
	}

	bd := TemporaryBlockDevice{
		geometry: g,
		storage:  make(map[pkg.BlockId][]pkg.Word36),
		isOpen:   false,
	}

	return &bd, DeviceResult{DeviceStatusSuccessful, nil}
}
