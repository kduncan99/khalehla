// Khalehla Project
// tiny assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

type Reference struct {
	symbol      string
	startingBit int
	bitCount    int
	offset      int
}

func NewReference(symbol string, startingBit int, bitCount int, offset int) *Reference {
	return &Reference{
		symbol:      symbol,
		startingBit: startingBit,
		bitCount:    bitCount,
		offset:      offset,
	}
}
