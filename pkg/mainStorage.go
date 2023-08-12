// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

import (
	"fmt"
)

type MainStorage struct {
	segmentMap         map[uint64][]Word36
	freeSegmentIndices []uint64
	maxIndices         uint64
}

// Allocate obtains a storage segment of the indicated type, returning the index of the segment
func (ms *MainStorage) Allocate(length uint64) (uint64, error) {
	if uint64(len(ms.segmentMap)) == ms.maxIndices {
		return 0, fmt.Errorf("main storage segment table is full")
	}

	var seg uint64
	if len(ms.freeSegmentIndices) == 0 {
		seg = uint64(len(ms.segmentMap))
	} else {
		ix := len(ms.freeSegmentIndices) - 1
		seg = ms.freeSegmentIndices[ix]
		ms.freeSegmentIndices = ms.freeSegmentIndices[:]
	}

	ms.segmentMap[seg] = make([]Word36, length)
	return seg, nil
}

func (ms *MainStorage) Clear() {
	ms.freeSegmentIndices = make([]uint64, 0)
	ms.segmentMap = make(map[uint64][]Word36)
}

func (ms *MainStorage) Dump() {
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
}

func (ms *MainStorage) GetSegment(segmentIndex uint64) (segment []Word36, interrupt Interrupt) {
	var ok bool
	segment, ok = ms.segmentMap[segmentIndex]
	if !ok {
		interrupt = NewHardwareCheckInterrupt(NewAbsoluteAddress(segmentIndex, 0))
	}

	return
}

func (ms *MainStorage) GetSlice(segmentIndex uint64, offset uint64, length uint64) (slice []Word36, interrupt Interrupt) {
	segment, ok := ms.segmentMap[segmentIndex]
	if !ok {
		interrupt = NewHardwareCheckInterrupt(NewAbsoluteAddress(segmentIndex, offset))
		return
	}

	if offset+length >= uint64(len(segment)) {
		interrupt = NewHardwareCheckInterrupt(NewAbsoluteAddress(segmentIndex, offset))
		return
	}

	slice = segment[offset : offset+length]
	return
}

func (ms *MainStorage) GetSliceFromAddress(absAddr *AbsoluteAddress, length uint64) (slice []Word36, interrupt Interrupt) {
	return ms.GetSlice(absAddr.GetSegment(), absAddr.GetOffset(), length)
}

var zero = Word36(0)

func (ms *MainStorage) GetWordFromAddress(absAddr *AbsoluteAddress) (word *Word36, interrupt Interrupt) {
	word = &zero
	interrupt = nil

	var segment []Word36
	segment, interrupt = ms.GetSegment(absAddr.GetSegment())
	if interrupt == nil {
		return
	}

	offset := absAddr.GetOffset()
	if offset >= uint64(len(segment)) {
		interrupt = NewHardwareCheckInterrupt(absAddr)
		return
	}

	word = &segment[offset]
	return
}

func (ms *MainStorage) Release(segmentIndex uint64) (interrupt Interrupt) {
	_, ok := ms.segmentMap[segmentIndex]
	if ok {
		delete(ms.segmentMap, segmentIndex)
		ms.freeSegmentIndices = append(ms.freeSegmentIndices, segmentIndex)
	} else {
		interrupt = NewHardwareCheckInterrupt(NewAbsoluteAddress(segmentIndex, 0))
	}

	return
}

func (ms *MainStorage) Resize(segmentIndex uint64, length uint64) (interrupt Interrupt) {
	slice, ok := ms.segmentMap[segmentIndex]
	if !ok {
		interrupt = NewHardwareCheckInterrupt(NewAbsoluteAddress(segmentIndex, 0))
		return
	}

	if length > uint64(len(slice)) {
		extension := make([]Word36, length-uint64(len(slice)))
		ms.segmentMap[segmentIndex] = append(slice, extension...)
	} else if length < uint64(len(slice)) {
		dst := make([]Word36, length)
		copy(dst, slice)
		ms.segmentMap[segmentIndex] = dst
	}

	return
}

func NewMainStorage(maxIndices uint64) *MainStorage {
	ms := MainStorage{
		segmentMap:         make(map[uint64][]Word36),
		freeSegmentIndices: make([]uint64, 0),
		maxIndices:         maxIndices,
	}
	return &ms
}
