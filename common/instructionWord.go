// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package common

const (
	B0 = iota
	B1
	B2
	B3
	B4
	B5
	B6
	B7
	B8
	B9
	B10
	B11
	B12
	B13
	B14
	B15
	B16
	B17
	B18
	B19
	B20
	B21
	B22
	B23
	B24
	B25
	B26
	B27
	B28
	B29
	B30
	B31
)

const (
	JFieldW   = 0
	JFieldH2  = 1
	JFieldH1  = 2
	JFieldXH2 = 3
	JFieldXH1 = 4
	JFieldQ2  = 4
	JFieldT3  = 5
	JFieldQ4  = 5
	JFieldT2  = 6
	JFieldQ3  = 6
	JFieldT1  = 7
	JFieldQ1  = 7
	JFieldS6  = 8
	JFieldS5  = 9
	JFieldS4  = 10
	JFieldS3  = 11
	JFieldS2  = 12
	JFieldS1  = 13
	JFieldU   = 14
	JFieldXU  = 15
)

type InstructionWord Word36

func (iw *InstructionWord) GetF() uint64 {
	return uint64(*iw) >> 30
}

func (iw *InstructionWord) GetJ() uint64 {
	return (uint64(*iw) >> 26) & 017
}

func (iw *InstructionWord) GetA() uint64 {
	return (uint64(*iw) >> 22) & 017
}

func (iw *InstructionWord) GetX() uint64 {
	return (uint64(*iw) >> 18) & 017
}

func (iw *InstructionWord) GetHIU() uint64 {
	return uint64(*iw) & 0777777
}

func (iw *InstructionWord) GetH() uint64 {
	return (uint64(*iw) >> 17) & 01
}

func (iw *InstructionWord) GetI() uint64 {
	return (uint64(*iw) >> 16) & 01
}

func (iw *InstructionWord) GetIB() uint64 {
	return (uint64(*iw) >> 12) & 037
}

func (iw *InstructionWord) GetU() uint64 {
	return uint64(*iw) & 0177777
}

func (iw *InstructionWord) GetB() uint64 {
	return (uint64(*iw) >> 12) & 017
}

func (iw *InstructionWord) GetD() uint64 {
	return uint64(*iw) & 07777
}

func (iw *InstructionWord) GetW() uint64 {
	return uint64(*iw)
}

func (iw *InstructionWord) SetW(value uint64) {
	*iw = InstructionWord(value & 0_777777_777777)
}

func (iw *InstructionWord) SetXHIU(value uint64) {
	res := uint64(*iw) & 0_777760_000000
	res |= value & 017_777777
	*iw = InstructionWord(res)
}
