// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type Offset interface {
	GetStartBit() int
	GetBitLength() int
	GetSymbol() string
}

type LocationCounterOffset struct {
	startBit  int
	bitLength int
	symbol    string
}

type UndefinedReferenceOffset struct {
	startBit  int
	bitLength int
	symbol    string
}

// CodeWord represents a single word of generated code, with potential relocation and/or undefined reference
// information attached.
type CodeWord struct {
	baseValue                 uint64
	locationCounterOffset     *LocationCounterOffset
	UndefinedReferenceOffsets []*UndefinedReferenceOffset
}
