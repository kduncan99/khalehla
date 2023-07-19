// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

// Procedure defines a proc
type Procedure struct {
	isBasicMode    bool
	isExtendedMode bool
	code           []string
	externalNames  map[string]int // maps a $NAME to the textIndex of the line which contains it
}
