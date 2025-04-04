// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

// locationSegment represents the code, external refs, and relocations for a chunk of code
// corresponding to a location counter of a module.
type locationSegment struct {
	code []CodeWord
	//  TODO external references...?
}
