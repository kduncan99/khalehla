// khalehla Project
// disassembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package dasm

import (
	"fmt"

	"khalehla/common"
	"khalehla/hardware"
	"khalehla/tasm"
)

// Disassembler is a simple disassembler which assists in khalehla development

type Disassembler struct {
}

func NewDisassembler() *Disassembler {
	return &Disassembler{}
}

func (da *Disassembler) DisassembleStorage(storage *hardware.MainStorage, executable *tasm.Executable) {
	//	TODO I think we need BDTables as a parameter as well...?
}

func DisassembleInstruction(asp *common.ActivityStatePacket) string {
	var s string
	var ok bool
	iw := asp.GetCurrentInstruction()
	dr := asp.GetDesignatorRegister()
	if asp.GetDesignatorRegister().IsBasicModeEnabled() {
		s, ok = BasicFunctionTable.Interpret(iw, dr.IsBasicModeEnabled(), dr.IsQuarterWordModeEnabled())
		if !ok {
			s = fmt.Sprintf("%012o", *iw)
		}
	} else {
		s, ok = ExtendedFunctionTable.Interpret(iw, dr.IsBasicModeEnabled(), dr.IsQuarterWordModeEnabled())
		if !ok {
			s = fmt.Sprintf("%012o", *iw)
		}
	}

	return s
}

//	TODO WiP... this will definitely get more complicated in the near future

type Mode int

const (
	BOTH Mode = iota
	BASIC
	EXTENDED
)
