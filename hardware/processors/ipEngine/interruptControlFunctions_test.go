// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
)

const fER = 072
const fSGNL = 073
const fPAIJ = 074
const fAAIJ = 074

const jER = 011
const jSGNL = 015
const jPAIJ = 014
const jAAIJExtended = 014
const jAAIJBasic = 007

const aER = 000
const aSGNL = 017
const aPAIJ = 007
const aAAIJExtended = 006
const aAAIJBasic = 000

// ---------------------------------------------------
// ER

func erSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fER, jER, aER, 0, ref)
}

func erSourceItemU(u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fER, jER, aER, 0, u)
}

// ---------------------------------------------------
// SGNL

func sgnlSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fSGNL, jSGNL, aSGNL, 0, ref)
}

func sgnlSourceItemU(u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fSGNL, jSGNL, aSGNL, 0, u)
}

// ---------------------------------------------------
// PAIJ

func paijSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fPAIJ, jPAIJ, aPAIJ, x, h, i, ref)
}

func paijSourceItemHIBRef(x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fPAIJ, jPAIJ, aPAIJ, x, h, i, b, ref)
}

func paijSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fPAIJ, jPAIJ, aPAIJ, 0, ref)
}

// ---------------------------------------------------
// AAIJ

func aaijSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fAAIJ, jAAIJBasic, aAAIJBasic, x, h, i, ref)
}

func aaijSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fAAIJ, jAAIJExtended, aAAIJExtended, x, h, i, b, ref)
}

func aaijSourceItemRefBasic(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fAAIJ, jAAIJBasic, aAAIJBasic, 0, ref)
}

func aaijSourceItemRefExtended(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fAAIJ, jAAIJExtended, aAAIJExtended, 0, ref)
}

// ---------------------------------------------------------------------------------------------------------------------

// TODO ER
// TODO SGNL
// TODO PAIJ Basic
// TODO PAIJ Extended
// TODO AAIJ Basic
// TODO AAIJ Extended
