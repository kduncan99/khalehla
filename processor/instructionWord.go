package processor

import "kalehla/types"

type InstructionWord types.Word36

func (iw *InstructionWord) GetF() uint {
	return uint(*iw) >> 30
}

func (iw *InstructionWord) GetJ() uint {
	return (uint(*iw) >> 26) & 0xF
}

func (iw *InstructionWord) GetA() uint {
	return (uint(*iw) >> 22) & 0xF
}

func (iw *InstructionWord) GetX() uint {
	return (uint(*iw) >> 18) & 0xF
}

func (iw *InstructionWord) GetHIU() uint {
	return uint(*iw) & 0777777
}

func (iw *InstructionWord) GetH() uint {
	return (uint(*iw) >> 17) & 01
}

func (iw *InstructionWord) GetI() uint {
	return (uint(*iw) >> 16) & 01
}

func (iw *InstructionWord) GetU() uint {
	return uint(*iw) & 0177777
}

func (iw *InstructionWord) GetB() uint {
	return (uint(*iw) >> 12) & 0xF
}

func (iw *InstructionWord) GetD() uint {
	return uint(*iw) & 07777
}
