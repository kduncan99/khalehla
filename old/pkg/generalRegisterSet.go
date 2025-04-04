// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

import (
	"fmt"
)

const (
	X0 uint64 = iota
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
	A0 uint64 = iota + 12
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
	R0 uint64 = iota + 64
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
	ER0 uint64 = iota + 80
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
	EX0 uint64 = iota + 96
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
	EA0 uint64 = iota + 108
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

var RegisterNames = map[uint64]string{
	X0:   "X0",
	X1:   "X1",
	X2:   "X2",
	X3:   "X3",
	X4:   "X4",
	X5:   "X5",
	X6:   "X6",
	X7:   "X7",
	X8:   "X8",
	X9:   "X9",
	X10:  "X10",
	X11:  "X11",
	A0:   "A0",
	A1:   "A1",
	A2:   "A2",
	A3:   "A3",
	A4:   "A4",
	A5:   "A5",
	A6:   "A6",
	A7:   "A7",
	A8:   "A0",
	A9:   "A9",
	A10:  "A10",
	A11:  "A11",
	A12:  "A12",
	A13:  "A13",
	A14:  "A14",
	A15:  "A15",
	R0:   "R0",
	R1:   "R1",
	R2:   "R2",
	R3:   "R3",
	R4:   "R4",
	R5:   "R5",
	R6:   "R6",
	R7:   "R7",
	R8:   "R0",
	R9:   "R9",
	R10:  "R10",
	R11:  "R11",
	R12:  "R12",
	R13:  "R13",
	R14:  "R14",
	R15:  "R15",
	EX0:  "EX0",
	EX1:  "EX1",
	EX2:  "EX2",
	EX3:  "EX3",
	EX4:  "EX4",
	EX5:  "EX5",
	EX6:  "EX6",
	EX7:  "EX7",
	EX8:  "EX8",
	EX9:  "EX9",
	EX10: "EX10",
	EX11: "EX11",
	EA0:  "EA0",
	EA1:  "EA1",
	EA2:  "EA2",
	EA3:  "EA3",
	EA4:  "EA4",
	EA5:  "EA5",
	EA6:  "EA6",
	EA7:  "EA7",
	EA8:  "EA0",
	EA9:  "EA9",
	EA10: "EA10",
	EA11: "EA11",
	EA12: "EA12",
	EA13: "EA13",
	EA14: "EA14",
	EA15: "EA15",
	ER0:  "ER0",
	ER1:  "ER1",
	ER2:  "ER2",
	ER3:  "ER3",
	ER4:  "ER4",
	ER5:  "ER5",
	ER6:  "ER6",
	ER7:  "ER7",
	ER8:  "ER0",
	ER9:  "ER9",
	ER10: "ER10",
	ER11: "ER11",
	ER12: "ER12",
	ER13: "ER13",
	ER14: "ER14",
	ER15: "ER15",
}

type GeneralRegisterSet struct {
	registers []Word36
}

func NewGeneralRegisterSet() *GeneralRegisterSet {
	GeneralRegisterSet := GeneralRegisterSet{}
	GeneralRegisterSet.registers = make([]Word36, 128)
	return &GeneralRegisterSet
}

func (grs *GeneralRegisterSet) Clear() {
	for rx := 0; rx < len(grs.registers); rx++ {
		grs.registers[rx] = 0
	}
}

func (grs *GeneralRegisterSet) Dump() {
	fmt.Printf("  General Register Set\n")
	for gx := 0; gx < 128; gx += 8 {
		str := fmt.Sprintf("    %04o:  ", gx)
		for gy := gx; gy < gx+8; gy++ {
			str += fmt.Sprintf("%012o ", grs.registers[gy])
		}
		fmt.Printf("%s\n", str)
	}
}

func (grs *GeneralRegisterSet) GetRegister(regName uint64) *Word36 {
	return &grs.registers[regName]
}

func (grs *GeneralRegisterSet) GetConsecutiveRegisters(firstRegName uint64, count uint64) []Word36 {
	return grs.registers[firstRegName : firstRegName+count]
}

func (grs *GeneralRegisterSet) GetRegisterValue(regName uint64) uint64 {
	return grs.registers[regName].GetW()
}

func (grs *GeneralRegisterSet) SetRegisterValue(regName uint64, value uint64) {
	grs.registers[regName].SetW(value)
}
