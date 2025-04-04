// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package common

import (
	"fmt"
)

// AbsoluteAddress structs are defined architecturally as a composite value generally not exceeding 54 bits,
// which allows the underlying hardware to determine the real-life location of a word in storage.
//
//	Given the following architectural specifications:
//		Small_Bank contains not more than 2^18 words
//		Large_Bank contains not more than 2^24 words
//		Very_Large_Bank is (usually) conceptual, being described as a set of consecutive large banks,
//			but in any case contains not more than 2^33 words.
//		Newer structures include the data expanse and the postern banks, which we do not implement,
//			at least for now.
//
// We need to implement an addressing scheme which accommodates these bank sizes.
// We expect that all virtual address banks will be contained in individual real-word allocated blocks of memory
// which are sized according to the logical bank size, where one logical word is right-justified and zero-filled in a
// space of 64 bits, which we represent as a uint64 (and thus, the bank is a slice of uint64 entities).
// Whenever the host program / operating system needs to carve out a new logical bank from machine memory,
// it will invoke a service instruction which does a malloc() to get a new block of real memory, and then the
// operating system is responsible for updating some existing bank descriptor table accordingly.
//
// Thus, our absolute address consists of the concatenation of a 21-bit segment identifier
// which allows us to describe a maximum of 2^21 (4194304) defined banks, with a offset value of 33 bits
// which describes a particular address within the bank identified by the segment identifier.
// This should be sufficient for our purposes, as we are focused on scale-out, not scale-up.
// To ensure that the non-significant bit fields are always zero, all code setting these values must use the Set* methods.
type AbsoluteAddress struct {
	segment uint   // 21 bits significant
	offset  uint64 // 33 bits significant
}

func (aa *AbsoluteAddress) Equals(comp *AbsoluteAddress) bool {
	return aa.segment == comp.segment && aa.offset == comp.offset
}

func (aa *AbsoluteAddress) GetComposite() []uint64 {
	return []uint64{uint64(aa.segment), aa.offset}
}

func (aa *AbsoluteAddress) GetOffset() uint64 {
	return aa.offset
}

func (aa *AbsoluteAddress) GetSegment() uint {
	return aa.segment
}

func (aa *AbsoluteAddress) GetString() string {
	return fmt.Sprintf("%012o:%012o", aa.segment, aa.offset)
}

func (aa *AbsoluteAddress) SetComposite(value []uint64) *AbsoluteAddress {
	aa.segment = uint(value[0])
	aa.offset = value[1]
	return aa
}

func (aa *AbsoluteAddress) SetCompositeFromWord36(value []Word36) *AbsoluteAddress {
	aa.segment = uint(value[0].GetW())
	aa.offset = value[1].GetW()
	return aa
}

func (aa *AbsoluteAddress) SetSegment(value uint) *AbsoluteAddress {
	aa.segment = value & 07_777777
	return aa
}

func (aa *AbsoluteAddress) SetOffset(value uint64) *AbsoluteAddress {
	aa.offset = value & 0_077777_7777777
	return aa
}

func NewAbsoluteAddress(segmentIndex uint, offset uint64) *AbsoluteAddress {
	return &AbsoluteAddress{
		segment: segmentIndex,
		offset:  offset,
	}
}
