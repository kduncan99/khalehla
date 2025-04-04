// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
)

const fOR = 040
const fXOR = 041
const fAND = 042
const fMLU = 043

// ---------------------------------------------------
// OR

func orSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fOR, j, a, x, h, i, ref)
}

func orSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fOR, j, a, x, h, i, b, ref)
}

func orSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fOR, j, a, x, u)
}

// ---------------------------------------------------
// XOR

func xorSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fXOR, j, a, x, h, i, ref)
}

func xorSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fXOR, j, a, x, h, i, b, ref)
}

func xorSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fXOR, j, a, x, u)
}

// ---------------------------------------------------
// AND

func andSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fAND, j, a, x, h, i, ref)
}

func andSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fAND, j, a, x, h, i, b, ref)
}

func andSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fAND, j, a, x, u)
}

// ---------------------------------------------------
// MLU

func mluSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fMLU, j, a, x, h, i, ref)
}

func mluSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fMLU, j, a, x, h, i, b, ref)
}

func mluSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fMLU, j, a, x, u)
}

// ---------------------------------------------------------------------------------------------------------------------

// TODO OR
// TODO XOR
// TODO AND
// TODO MLU
