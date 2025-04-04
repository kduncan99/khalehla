// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

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
