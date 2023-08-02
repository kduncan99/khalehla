// Khalehla Project
// disassembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

import (
	"khalehla/pkg"
	"khalehla/tasm"
)

// DisAssembler is a simple disassembler which assists in Khalehla development
type DisAssembler struct {
}

func NewDisAssembler() *DisAssembler {
	return &DisAssembler{}
}

func (da *DisAssembler) DisAssemble(storage *pkg.MainStorage, executable *tasm.Executable) {

}

func (da *DisAssembler) DisAssembleLine(asp *pkg.ActivityStatePacket, instruction []pkg.Word36) {

}
