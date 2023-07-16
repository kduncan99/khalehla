// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

type InstructionWord pkg.Word36

func (iw *InstructionWord) GetF() uint64 {
	return uint64(*iw) >> 30
}

func (iw *InstructionWord) GetJ() uint64 {
	return (uint64(*iw) >> 26) & 0xF
}

func (iw *InstructionWord) GetA() uint64 {
	return (uint64(*iw) >> 22) & 0xF
}

func (iw *InstructionWord) GetX() uint64 {
	return (uint64(*iw) >> 18) & 0xF
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

func (iw *InstructionWord) GetU() uint64 {
	return uint64(*iw) & 0177777
}

func (iw *InstructionWord) GetB() uint64 {
	return (uint64(*iw) >> 12) & 0xF
}

func (iw *InstructionWord) GetD() uint64 {
	return uint64(*iw) & 07777
}
