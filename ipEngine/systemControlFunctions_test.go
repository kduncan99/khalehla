// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/tasm"

const fSPID = 073
const fIPC = 073
const fSYSC = 073
const fIAR = 073

const jSPID = 015
const jIPC = 017
const jSYSC = 017
const jIAR = 017

const aSPID = 005
const aIPC = 010
const aSYSC = 012
const aIAR = 006

// ---------------------------------------------------
// SPID

//	TODO

// ---------------------------------------------------
// IPC - extended mode only

//	TODO

// ---------------------------------------------------
// SYSC - extended mode only

//	TODO

// ---------------------------------------------------
// IAR - extended mode only (although we may implement basic mode as well)

func iarSourceItem(uField uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fIAR, jIAR, aIAR, 0, uField)
}

// ---------------------------------------------------------------------------------------------------------------------

// TODO SPID
// TODO IPC
// TODO SYSC
// TODO IAR
