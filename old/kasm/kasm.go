// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type StringCodeType int

const (
	FieldataString = iota
	AsciiString
)

type ValueFlags int

const (
	SingleFlag         = 1 << 0
	DoubleFlag         = 1 << 1
	LeftJustifiedFlag  = 1 << 2
	RightJustifiedFlag = 1 << 3
	FlaggedFlag        = 1 << 4
)
