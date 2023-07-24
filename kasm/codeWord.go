// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

// CodeWord represents a single word of generated code, with potential relocation and/or undefined
// referenceExpressionItem information attached.
type CodeWord struct {
	baseValue uint64
	offsets   []Offset
}
