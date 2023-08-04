// Khalehla Project
// tiny assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

type Reference struct {
	symbol      string
	startingBit uint64
	bitCount    uint64
	offset      uint64
}

func NewReference(symbol string, startingBit uint64, bitCount uint64, offset uint64) *Reference {
	return &Reference{
		symbol:      symbol,
		startingBit: startingBit,
		bitCount:    bitCount,
		offset:      offset,
	}
}
