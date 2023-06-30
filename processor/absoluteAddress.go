package processor

// AbsoluteAddress structs are defined architecturally as a composite value generally not exceeding 54 bits,
//
//	which allows the underlying hardware to determine the real-life location of a word in storage.
//	Given the following architectural specifications:
//		Small_Bank contains not more than 2^18 words
//		Large_Bank contains not more than 2^24 words
//		Very_Large_Bank is (usually) conceptual, being described as a set of consecutive large banks,
//			but in any case contains not more than 2^33 words.
//		Newer structures include the data expanse and the postern banks, which we do not implement,
//			at least for now.
//
//	Therefore, we need to implement an addressing scheme which accommodates these bank sizes.
//	We expect that all banks will be contained in real-word allocated blocks of memory which are sized
//	according to the logical bank size, where one logical word (Word36) struct is right-justified and
//	guaranteed zero-filled in a space of 64 bits.
//
//	In this light, our absolute address consists of the concatenation of a 21-bit segment identifier
//	which allows us to describe a maximum of 2^21 (4194304) defined banks, with a offset value of 33 bits
//	which describes a particular address within the bank identified by the segment identifier.
//
//	This should be sufficient for our purposes, as we are focused on scale-out, not scale-up.
//
//	To ensure that the non-significant bit fields are always zero,
//	all code setting these values must use the Set* methods.
type AbsoluteAddress struct {
	segment uint //	21 bits significant
	offset  uint //	33 bits significant
}

func (aa *AbsoluteAddress) GetComposite() uint64 {
	return uint64(aa.segment<<33) | uint64(aa.offset)
}

func (aa *AbsoluteAddress) SetComposite(value uint64) *AbsoluteAddress {
	aa.segment = uint(value>>33) | 07_777777
	aa.offset = uint(value & 0_077777_777777)
	return aa
}

func (aa *AbsoluteAddress) SetSegment(value uint) *AbsoluteAddress {
	aa.segment = value & 07_777777
	return aa
}

func (aa *AbsoluteAddress) SegOffset(value uint) *AbsoluteAddress {
	aa.offset = value & 0_077777_7777777
	return aa
}

func NewAbsoluteAddressFromComposite(value uint64) *AbsoluteAddress {
	aa := AbsoluteAddress{}
	aa.SetComposite(value)
	return &aa
}
