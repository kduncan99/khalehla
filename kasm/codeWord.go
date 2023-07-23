// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type OffsetType int

const (
	LocationCounterOffsetType = iota
	UndefinedReferenceOffsetType
)

type Offset interface {
	GetStartBit() int
	GetBitLength() int
	IsNegative() bool
}

type LocationCounterOffset struct {
	locationCounter int
	startBit        int
	bitLength       int
	isNegative      bool
}

type UndefinedReferenceOffset struct {
	startBit   int
	bitLength  int
	symbol     string
	isNegative bool
}

// CodeWord represents a single word of generated code, with potential relocation and/or undefined
// referenceExpressionItem information attached.
type CodeWord struct {
	baseValue uint64
	offsets   []Offset
}

func (lco *LocationCounterOffset) GetBitLength() int {
	return lco.bitLength
}

func (lco *LocationCounterOffset) GetStartBit() int {
	return lco.startBit
}

func (lco *LocationCounterOffset) GetLocationCounter() int {
	return lco.bitLength
}

func (lco *LocationCounterOffset) IsNegative() bool {
	return lco.isNegative
}

func (uro *UndefinedReferenceOffset) GetBitLength() int {
	return uro.bitLength
}

func (uro *UndefinedReferenceOffset) GetStartBit() int {
	return uro.startBit
}

func (uro *UndefinedReferenceOffset) GetSymbol() string {
	return uro.symbol
}

func (uro *UndefinedReferenceOffset) IsNegative() bool {
	return uro.isNegative
}
