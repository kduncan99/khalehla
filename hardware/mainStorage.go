// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package hardware

import (
	"fmt"
	"sync"
	"time"

	"khalehla/common"
)

type MainStorageClient interface{}

const waitTime = time.Millisecond

type MainStorage struct {
	segmentMap         map[uint][]common.Word36
	freeSegmentIndices []uint
	maxIndices         uint
	locks              map[storageLockKey]StorageLockClient
	mutex              sync.Mutex
}

type StorageLockClient interface {
	GetStorageLockClientName() string
}

type storageLockKey uint64

func newStorageLockKey(address VirtualAddress) storageLockKey {
	return storageLockKey(address.GetComposite())
}

// NewMainStorage creates one of these entities.
// There should be exactly one, shared among all the InstructionEngine instances.
// This struct serves two purposes.
// 1) It acts as main storage for the entire processing complex.
// It provides memory in relatively small chunks, intended to be one-for-one with each logical bank of storage.
// Each chunk is allocated either outside of processing by some operating system loader, or inside processing
// via operating system calls to a service instruction intended for the purpose.
// 2) The storage lock table protects the integrity of the following instructions:
// ADD1, CR, DEC, DEC2, ENZ, INC, INC2, SUB1, TCS, TS, TSS
func NewMainStorage(maxIndices uint) *MainStorage {
	ms := MainStorage{
		segmentMap:         make(map[uint][]common.Word36),
		freeSegmentIndices: make([]uint, 0),
		maxIndices:         maxIndices,
		locks:              make(map[storageLockKey]StorageLockClient),
		mutex:              sync.Mutex{},
	}
	return &ms
}

// Allocate obtains a storage segment of the indicated type, returning the index of the segment.
// May be invoked by a service processor for pre-loading memory prior to booting the system,
// or by an instruction processor as part of executing a service instruction designed for that purpose.
func (ms *MainStorage) Allocate(length uint64) (uint, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if uint(len(ms.segmentMap)) == ms.maxIndices {
		return 0, fmt.Errorf("main storage segment table is full")
	}

	var seg uint
	if len(ms.freeSegmentIndices) == 0 {
		seg = uint(len(ms.segmentMap))
	} else {
		ix := len(ms.freeSegmentIndices) - 1
		seg = ms.freeSegmentIndices[ix]
		ms.freeSegmentIndices = ms.freeSegmentIndices[:]
	}

	ms.segmentMap[seg] = make([]common.Word36, length)
	return seg, nil
}

// Clear will ensure the entire storage is removed
func (ms *MainStorage) Clear() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.freeSegmentIndices = make([]uint, 0)
	ms.segmentMap = make(map[uint][]common.Word36)
}

// Dump will display the content of memory to stdout - used only for debugging
func (ms *MainStorage) Dump() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	fmt.Printf("Main Storage Dump ----------------------\n")

	fmt.Printf("  Free Segments:\n")
	if len(ms.freeSegmentIndices) > 0 {
		for _, idx := range ms.freeSegmentIndices {
			fmt.Printf("    %d ", idx)
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("    none\n")
	}

	for index, slice := range ms.segmentMap {
		fmt.Printf("  Segment %d:\n", index)
		for ix := 0; ix < len(slice); ix += 8 {
			fmt.Printf("    %08o:  ", ix)
			yLimit := ix + 8
			if yLimit > len(slice) {
				yLimit = len(slice)
			}
			for iy := ix; iy < yLimit; iy++ {
				fmt.Printf("%012o ", slice[iy])
			}
			fmt.Printf("\n")
		}
	}

	for value, client := range ms.locks {
		fmt.Printf("    %012o:%s\n", value, client.GetStorageLockClientName())
	}
}

func (ms *MainStorage) getSegmentWorker(segmentIndex uint) (segment []common.Word36, interrupt common.Interrupt) {
	var ok bool
	segment, ok = ms.segmentMap[segmentIndex]
	if !ok {
		interrupt = common.NewHardwareCheckInterrupt(common.NewAbsoluteAddress(segmentIndex, 0))
	}
	return
}

func (ms *MainStorage) GetSegment(segmentIndex uint) (segment []common.Word36, interrupt common.Interrupt) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	return ms.getSegmentWorker(segmentIndex)
}

func (ms *MainStorage) GetSlice(segmentIndex uint, offset uint64, length uint64) (slice []common.Word36, interrupt common.Interrupt) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	segment, ok := ms.segmentMap[segmentIndex]
	if !ok {
		interrupt = common.NewHardwareCheckInterrupt(common.NewAbsoluteAddress(segmentIndex, offset))
		return
	}

	if offset+length > uint64(len(segment)) {
		interrupt = common.NewHardwareCheckInterrupt(common.NewAbsoluteAddress(segmentIndex, offset))
		return
	}

	slice = segment[offset : offset+length]
	return
}

func (ms *MainStorage) GetSliceFromAddress(absAddr *common.AbsoluteAddress, length uint64) (slice []common.Word36, interrupt common.Interrupt) {
	return ms.GetSlice(absAddr.GetSegment(), absAddr.GetOffset(), length)
}

var zero = common.Word36(0)

func (ms *MainStorage) GetWordFromAddress(absAddr *common.AbsoluteAddress) (word *common.Word36, interrupt common.Interrupt) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	word = &zero
	interrupt = nil

	var segment []common.Word36
	segment, interrupt = ms.getSegmentWorker(absAddr.GetSegment())
	if interrupt != nil {
		return
	}

	offset := absAddr.GetOffset()
	if offset >= uint64(len(segment)) {
		interrupt = common.NewHardwareCheckInterrupt(absAddr)
		return
	}

	word = &segment[offset]
	return
}

func (ms *MainStorage) Release(segmentIndex uint) (interrupt common.Interrupt) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	_, ok := ms.segmentMap[segmentIndex]
	if ok {
		delete(ms.segmentMap, segmentIndex)
		ms.freeSegmentIndices = append(ms.freeSegmentIndices, segmentIndex)
	} else {
		interrupt = common.NewHardwareCheckInterrupt(common.NewAbsoluteAddress(segmentIndex, 0))
	}

	return
}

func (ms *MainStorage) Resize(segmentIndex uint, length uint64) (interrupt common.Interrupt) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	slice, ok := ms.segmentMap[segmentIndex]
	if !ok {
		interrupt = common.NewHardwareCheckInterrupt(common.NewAbsoluteAddress(segmentIndex, 0))
		return
	}

	if length > uint64(len(slice)) {
		extension := make([]common.Word36, length-uint64(len(slice)))
		ms.segmentMap[segmentIndex] = append(slice, extension...)
	} else if length < uint64(len(slice)) {
		dst := make([]common.Word36, length)
		copy(dst, slice)
		ms.segmentMap[segmentIndex] = dst
	}

	return
}

func (ms *MainStorage) Lock(address VirtualAddress, client StorageLockClient) bool {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	key := newStorageLockKey(address)
	_, ok := ms.locks[key]
	if ok {
		return false
	}

	ms.locks[key] = client
	return true
}

func (ms *MainStorage) LockWait(address VirtualAddress, client StorageLockClient) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	for true {
		key := newStorageLockKey(address)
		_, ok := ms.locks[key]
		if ok {
			ms.mutex.Unlock()
			time.Sleep(waitTime)
			ms.mutex.Lock()
		} else {
			ms.locks[key] = client
			return
		}
	}
}

func (ms *MainStorage) ReleaseLocks(address VirtualAddress, client StorageLockClient) bool {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	key := newStorageLockKey(address)
	lockClient, ok := ms.locks[key]
	if ok && lockClient == client {
		delete(ms.locks, key)
		return true
	} else {
		return false
	}
}

func (ms *MainStorage) ReleaseAllLocks(client StorageLockClient) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	for key, lockClient := range ms.locks {
		if lockClient == client {
			delete(ms.locks, key)
		}
	}
}
