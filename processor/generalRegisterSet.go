package processor

import "kalehla/types"

//	TODO important - review 2.3.6 reserved for hardware, and 2.3.7 accessing GeneralRegisterSet locations
//
// register names
const (
	X0 uint = iota
	X1
	X2
	X3
	X4
	X5
	X6
	X7
	X8
	X9
	X10
	X11
	X12
	X13
	X14
	X15
)

const (
	A0 uint = iota + 12
	A1
	A2
	A3
	A4
	A5
	A6
	A7
	A8
	A9
	A10
	A11
	A12
	A13
	A14
	A15
)

const (
	R0 uint = iota + 64
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	R11
	R12
	R13
	R14
	R15
)

const (
	ER0 uint = iota + 80
	ER1
	ER2
	ER3
	ER4
	ER5
	ER6
	ER7
	ER8
	ER9
	ER10
	ER11
	ER12
	ER13
	ER14
	ER15
)

const (
	EX0 uint = iota + 96
	EX1
	EX2
	EX3
	EX4
	EX5
	EX6
	EX7
	EX8
	EX9
	EX10
	EX11
	EX12
	EX13
	EX14
	EX15
)

const (
	EA0 uint = iota + 108
	EA1
	EA2
	EA3
	EA4
	EA5
	EA6
	EA7
	EA8
	EA9
	EA10
	EA11
	EA12
	EA13
	EA14
	EA15
)

type GeneralRegisterSet struct {
	Registers []types.Word36
}

func NewGeneralRegisterSet() *GeneralRegisterSet {
	GeneralRegisterSet := GeneralRegisterSet{}
	GeneralRegisterSet.Registers = make([]types.Word36, 128)
	return &GeneralRegisterSet
}

func (grs *GeneralRegisterSet) GetRegister(regName uint) *types.Word36 {
	return &grs.Registers[regName]
}

func (grs *GeneralRegisterSet) GetValueOfRegister(regName uint) types.Word36 {
	return grs.Registers[regName]
}

func (grs *GeneralRegisterSet) SetRegisterValue(regName uint, value types.Word36) {
	grs.Registers[regName] = value
}
