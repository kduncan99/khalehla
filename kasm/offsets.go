// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type OffsetType int

const (
	LocationCounterOffsetType = iota
	UndefinedReferenceOffsetType
)

type Offset interface {
	Equals(comp Offset) bool
	EqualsNegative(comp Offset) bool
	GetStartBit() int
	GetBitLength() int
	GetOffsetType() OffsetType
	IsNegative() bool
}

type LocationCounterOffset struct {
	locationCounter int
	startBit        int
	bitLength       int
	isNegative      bool
}

func NewLocationCounterOffset(lc int) *LocationCounterOffset {
	return &LocationCounterOffset{
		locationCounter: lc,
		startBit:        0,
		bitLength:       36,
		isNegative:      false,
	}
}

type UndefinedReferenceOffset struct {
	symbol     string
	startBit   int
	bitLength  int
	isNegative bool
}

func NewUndefinedReferenceOffset(symbol string) *UndefinedReferenceOffset {
	return &UndefinedReferenceOffset{
		symbol:     symbol,
		startBit:   0,
		bitLength:  36,
		isNegative: false,
	}
}

func (lco *LocationCounterOffset) Equals(comp Offset) bool {
	if comp.GetOffsetType() == LocationCounterOffsetType {
		lcoComp := comp.(*LocationCounterOffset)
		return lco.locationCounter == lcoComp.locationCounter &&
			lco.bitLength == lcoComp.bitLength &&
			lco.startBit == lcoComp.startBit &&
			lco.isNegative == lcoComp.isNegative
	}
	return false
}

func (lco *LocationCounterOffset) EqualsNegative(comp Offset) bool {
	if comp.GetOffsetType() == LocationCounterOffsetType {
		lcoComp := comp.(*LocationCounterOffset)
		return lco.locationCounter == lcoComp.locationCounter &&
			lco.bitLength == lcoComp.bitLength &&
			lco.startBit == lcoComp.startBit &&
			lco.isNegative != lcoComp.isNegative
	}
	return false
}

func (lco *LocationCounterOffset) GetBitLength() int {
	return lco.bitLength
}

func (lco *LocationCounterOffset) GetOffsetType() OffsetType {
	return LocationCounterOffsetType
}

func (lco *LocationCounterOffset) GetStartBit() int {
	return lco.startBit
}

func (lco *LocationCounterOffset) IsNegative() bool {
	return lco.isNegative
}

func (uro *UndefinedReferenceOffset) Equals(comp Offset) bool {
	if comp.GetOffsetType() == UndefinedReferenceOffsetType {
		uroComp := comp.(*UndefinedReferenceOffset)
		return uro.symbol == uroComp.symbol &&
			uro.bitLength == uroComp.bitLength &&
			uro.startBit == uroComp.startBit &&
			uro.isNegative == uroComp.isNegative
	}
	return false
}

func (uro *UndefinedReferenceOffset) EqualsNegative(comp Offset) bool {
	if comp.GetOffsetType() == UndefinedReferenceOffsetType {
		uroComp := comp.(*UndefinedReferenceOffset)
		return uro.symbol == uroComp.symbol &&
			uro.bitLength == uroComp.bitLength &&
			uro.startBit == uroComp.startBit &&
			uro.isNegative != uroComp.isNegative
	}
	return false
}

func (uro *UndefinedReferenceOffset) GetBitLength() int {
	return uro.bitLength
}

func (uro *UndefinedReferenceOffset) GetOffsetType() OffsetType {
	return UndefinedReferenceOffsetType
}

func (uro *UndefinedReferenceOffset) GetStartBit() int {
	return uro.startBit
}

func (uro *UndefinedReferenceOffset) IsNegative() bool {
	return uro.isNegative
}

// CollapseOffsetList returns a copy of a given offset list with arithmetic inverse items removed.
// This allows us to make sense of, for example, TAG2-TAG1 as the distance between two symbols which represent
// locations in a common location counter pool
func CollapseOffsetList(offsetList []Offset) []Offset {
	temp := offsetList[:]
	for tx := 0; tx < len(temp)-1; {
		ty := tx + 1
		found := false
		for ty < len(temp) {
			if temp[tx].EqualsNegative(temp[ty]) {
				found = true
				temp = append(temp[:tx], temp[tx+1:]...)
				ty--
				temp = append(temp[:ty], temp[ty+1:]...)
				break
			}
		}

		if !found {
			tx++
		}
	}

	return temp
}

// OffsetListsAreEqual compares two lists of offsets to see whether the collective one is equal to the
// collective other.
func OffsetListsAreEqual(offsetList1 []Offset, offsetList2 []Offset) bool {
	//	This is made more interesting by the fact that the offsets do not have to be in equal order.
	temp1 := offsetList1[:]
	temp2 := offsetList2[:]
	if len(temp1) == len(temp2) {
		for len(temp1) > 0 {
			off1 := temp1[0]
			temp1 = temp1[1:]

			found := false
			for x := 0; x < len(temp2); x++ {
				if off1.Equals(temp2[x]) {
					found = true
					temp2 = append(temp2[:x], temp2[x+1:]...)
					break
				}
			}

			if !found {
				return false
			}
		}

		return true
	}

	return false
}
