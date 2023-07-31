// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

import (
	"fmt"
)

type MainStorage struct {
	segmentMap         map[uint][]Word36
	freeSegmentIndices []uint
	maxIndices         uint
}

// Allocate obtains a storage segment of the indicated type, returning the index of the segment
func (ms *MainStorage) Allocate(length uint) (uint, error) {
	if uint(len(ms.segmentMap)) == ms.maxIndices {
		return 0, fmt.Errorf("main storage segment table is full")
	}

	seg := uint(0)
	if len(ms.freeSegmentIndices) == 0 {
		seg = uint(len(ms.segmentMap))
	} else {
		ix := len(ms.freeSegmentIndices) - 1
		seg = ms.freeSegmentIndices[ix]
		ms.freeSegmentIndices = ms.freeSegmentIndices[:]
	}

	ms.segmentMap[seg] = make([]Word36, length)
	return seg, nil
}

func (ms *MainStorage) Clear() {
	ms.freeSegmentIndices = make([]uint, 0)
	ms.segmentMap = make(map[uint][]Word36)
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

func (ms *MainStorage) GetBlock(segment uint) ([]Word36, bool) {
	sli, ok := ms.segmentMap[segment]
	return sli, ok
}

func (ms *MainStorage) GetSegment(segment uint) ([]Word36, bool) {
	sli, ok := ms.segmentMap[segment]
	if !ok {
		return nil, false
	} else {
		return sli, true
	}
}

func (ms *MainStorage) GetSlice(segment uint, offset uint, length uint) ([]Word36, bool) {
	sli, ok := ms.segmentMap[segment]
	if !ok {
		return nil, false
	} else {
		return sli[offset : offset+length], true
	}
}

func (ms *MainStorage) Release(segment uint) bool {
	_, ok := ms.segmentMap[segment]
	if ok {
		delete(ms.segmentMap, segment)
		ms.freeSegmentIndices = append(ms.freeSegmentIndices, segment)
	}
	return ok
}

func (ms *MainStorage) Resize(segment uint, length uint) error {
	slice, ok := ms.segmentMap[segment]
	if !ok {
		return fmt.Errorf("no such segment has been allocated")
	} else {
		if length > uint(len(slice)) {
			extension := make([]Word36, length-uint(len(slice)))
			ms.segmentMap[segment] = append(slice, extension...)
		} else if length < uint(len(slice)) {
			dst := make([]Word36, length)
			copy(dst, slice)
			ms.segmentMap[segment] = dst
		}
		return nil
	}
}

func NewMainStorage(maxIndices uint) *MainStorage {
	ms := MainStorage{
		segmentMap:         make(map[uint][]Word36),
		freeSegmentIndices: make([]uint, 0),
		maxIndices:         maxIndices,
	}
	return &ms
}
