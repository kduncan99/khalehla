// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
)

const fSA = 001
const fSNA = 002
const fSMA = 003
const fSR = 004
const fSX = 006
const fDS = 071
const fSZ = 005
const fSNZ = 005
const fSP1 = 005
const fSN1 = 005
const fSFS = 005
const fSFZ = 005
const fSAS = 005
const fSAZ = 005
const fSRS = 072
const fSAQW = 072

const jDS = 012
const jSRS = 016
const jSAQW = 005

const aSZ = 000
const aSNZ = 001
const aSP1 = 002
const aSN1 = 003
const aSFS = 004
const aSFZ = 005
const aSAS = 006
const aSAZ = 007

// ---------------------------------------------------
// SA

func saSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fSA, j, a, x, h, i, ref)
}

func saSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fSA, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// SFS

func sfsSourceItemHIRef(j uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fSFS, j, aSFS, x, h, i, ref)
}

func sfsSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fSFS, j, aSFS, x, h, i, b, ref)
}

// ---------------------------------------------------
// SFZ

func sfzSourceItemHIRef(j uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fSFZ, j, aSFZ, x, h, i, ref)
}

func sfzSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fSFZ, j, aSFZ, x, h, i, b, ref)
}

// ---------------------------------------------------
// SAS

func sasSourceItemHIRef(j uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fSAS, j, aSAS, x, h, i, ref)
}

func sasSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fSAS, j, aSAS, x, h, i, b, ref)
}

// ---------------------------------------------------
// SAZ

func sazSourceItemHIRef(j uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fSAZ, j, aSAZ, x, h, i, ref)
}

func sazSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fSAZ, j, aSAZ, x, h, i, b, ref)
}

// ---------------------------------------------------------------------------------------------------------------------

//	TODO everything
