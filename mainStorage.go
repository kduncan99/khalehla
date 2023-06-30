package kalehla

import (
	"fmt"
	"kalehla/types"
)

type MainStorage struct {
	segmentMap         map[uint][]types.Word36
	freeSegmentIndices []uint
	maxIndices         uint
}

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

	ms.segmentMap[seg] = make([]types.Word36, length)
	return seg, nil
}

func (ms *MainStorage) GetBlock(segment uint) ([]types.Word36, bool) {
	sli, ok := ms.segmentMap[segment]
	return sli, ok
}

func (ms *MainStorage) GetSlice(segment uint, offset uint, length uint) ([]types.Word36, bool) {
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

func NewMainStorage(maxIndices uint) *MainStorage {
	ms := MainStorage{
		segmentMap:         make(map[uint][]types.Word36),
		freeSegmentIndices: make([]uint, 0),
		maxIndices:         maxIndices,
	}
	return &ms
}
