// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"khalehla/pkg"
	"khalehla/tasm"
	"testing"
)

// partial word designators for j-field specification
const (
	jH1  = 002
	jH2  = 001
	jQ1  = 007
	jQ2  = 004
	jQ3  = 006
	jQ4  = 005
	jS1  = 015
	jS2  = 014
	jS3  = 013
	jS4  = 012
	jS5  = 011
	jS6  = 010
	jT1  = 007
	jT2  = 006
	jT3  = 005
	jU   = 016
	jW   = 000
	jXH2 = 003
	jXH1 = 004
	jXU  = 017
)

const zero = "0"

// GRS locations for registers
const (
	grsX0  = 000
	grsX1  = 001
	grsX2  = 002
	grsX3  = 003
	grsX4  = 004
	grsX5  = 005
	grsX6  = 006
	grsX7  = 007
	grsX8  = 010
	grsX9  = 011
	grsX10 = 012
	grsX11 = 013
	grsX12 = 014
	grsX13 = 015
	grsX14 = 016
	grsX15 = 017
)

const (
	grsA0  = 014
	grsA1  = 015
	grsA2  = 016
	grsA3  = 017
	grsA4  = 020
	grsA5  = 021
	grsA6  = 022
	grsA7  = 023
	grsA8  = 024
	grsA9  = 025
	grsA10 = 026
	grsA11 = 027
	grsA12 = 030
	grsA13 = 031
	grsA14 = 032
	grsA15 = 033
)

const (
	grsR0  = 0100
	grsR1  = 0101
	grsR2  = 0102
	grsR3  = 0103
	grsR4  = 0104
	grsR5  = 0105
	grsR6  = 0106
	grsR7  = 0107
	grsR8  = 0110
	grsR9  = 0111
	grsR10 = 0112
	grsR11 = 0113
	grsR12 = 0114
	grsR13 = 0115
	grsR14 = 0116
	grsR15 = 0117
)

// ---------------------------------------------------------------------------------------------------------------------

func sourceItem(label string, operator string, operands []int) *tasm.SourceItem {
	strOps := make([]string, len(operands))
	for ox := 0; ox < len(operands); ox++ {
		strOps[ox] = fmt.Sprintf("0%o", operands[ox])
	}

	return tasm.NewSourceItem(label, operator, strOps)
}

// ---------------------------------------------------------------------------------------------------------------------
// Load functions
// ---------------------------------------------------------------------------------------------------------------------

// DL ------------------------------------------------------------------------------------------------------------------

const fDL = 071
const jDL = 013

func dlSourceItemHIBD(label string, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fDL, jDL, a, x, h, i, b, d})
}

func dlSourceItemHIBDRef(label string, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLA),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func dlSourceItemHIU(label string, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLA, a, x, h, i, u})
}

func dlSourceItemHIURef(label string, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLX),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// LA ------------------------------------------------------------------------------------------------------------------

const fLA = 010

func laSourceItemHIBD(label string, j int, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLA, j, a, x, h, i, b, d})
}

func laSourceItemHIBDRef(label string, j int, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func laSourceItemHIU(label string, j int, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLA, j, a, x, h, i, u})
}

func laSourceItemHIURef(label string, j int, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func laSourceItemU(label string, j int, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLA, j, a, x, u})
}

// LMA -----------------------------------------------------------------------------------------------------------------

const fLMA = 012

func lmaSourceItemHIBD(label string, j int, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLMA, j, a, x, h, i, b, d})
}

func lmaSourceItemHIBDRef(label string, j int, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLMA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func lmaSourceItemHIU(label string, j int, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLMA, j, a, x, h, i, u})
}

func lmaSourceItemHIURef(label string, j int, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLMA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func lmaSourceItemU(label string, j int, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLMA, j, a, x, u})
}

// LNA -----------------------------------------------------------------------------------------------------------------

const fLNA = 011

func lnaSourceItemHIBD(label string, j int, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLNA, j, a, x, h, i, b, d})
}

func lnaSourceItemHIBDRef(label string, j int, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLNA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func lnaSourceItemHIU(label string, j int, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLNA, j, a, x, h, i, u})
}

func lnaSourceItemHIURef(label string, j int, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLNA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func lnaSourceItemU(label string, j int, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLNA, j, a, x, u})
}

// LNMA ----------------------------------------------------------------------------------------------------------------

const fLNMA = 013

func lnmaSourceItemHIBD(label string, j int, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLNMA, j, a, x, h, i, b, d})
}

func lnmaSourceItemHIBDRef(label string, j int, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLNMA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func lnmaSourceItemHIU(label string, j int, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLNMA, j, a, x, h, i, u})
}

func lnmaSourceItemHIURef(label string, j int, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLNMA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func lnmaSourceItemU(label string, j int, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLNMA, j, a, x, u})
}

// LR ------------------------------------------------------------------------------------------------------------------

const fLR = 023

func lrSourceItemHIBD(label string, j int, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLR, j, a, x, h, i, b, d})
}

func lrSourceItemHIBDRef(label string, j int, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLR),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func lrSourceItemHIU(label string, j int, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLR, j, a, x, h, i, u})
}

func lrSourceItemHIURef(label string, j int, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLR),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func lrSourceItemU(label string, j int, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLR, j, a, x, u})
}

// LX ------------------------------------------------------------------------------------------------------------------

const fLX = 027

func lxSourceItemHIBD(label string, j int, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLX, j, a, x, h, i, b, d})
}

func lxSourceItemHIBDRef(label string, j int, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLX),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func lxSourceItemHIU(label string, j int, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLX, j, a, x, h, i, u})
}

func lxSourceItemHIURef(label string, j int, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLX),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func lxSourceItemU(label string, j int, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLX, j, a, x, u})
}

// LXI -----------------------------------------------------------------------------------------------------------------

const fLXI = 046

func lxiSourceItemHIBD(label string, j int, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLXI, j, a, x, h, i, b, d})
}

func lxiSourceItemHIBDRef(label string, j int, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLXI),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func lxiSourceItemHIU(label string, j int, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLXI, j, a, x, h, i, u})
}

func lxiSourceItemHIURef(label string, j int, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLXI),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func lxiSourceItemU(label string, j int, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLXI, j, a, x, u})
}

// LXM -----------------------------------------------------------------------------------------------------------------

const fLXM = 026

func lxmSourceItemHIBD(label string, j int, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLXM, j, a, x, h, i, b, d})
}

func lxmSourceItemHIBDRef(label string, j int, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLXM),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func lxmSourceItemHIU(label string, j int, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLXM, j, a, x, h, i, u})
}

func lxmSourceItemHIURef(label string, j int, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLXM),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func lxmSourceItemU(label string, j int, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLXM, j, a, x, u})
}

// ---------------------------------------------------------------------------------------------------------------------
// Store functions
// ---------------------------------------------------------------------------------------------------------------------

// SA ------------------------------------------------------------------------------------------------------------------

const fSA = 010

func saSourceItemHIBD(label string, j int, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fSA, j, a, x, h, i, b, d})
}

func saSourceItemHIBDRef(label string, j int, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fSA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func saSourceItemHIU(label string, j int, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fSA, j, a, x, h, i, u})
}

func saSourceItemHIURef(label string, j int, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fSA),
		fmt.Sprintf("%o", j),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func saSourceItemU(label string, j int, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fSA, j, a, x, u})
}

// ---------------------------------------------------------------------------------------------------------------------
// Fixed-point Binary functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Floating-point Binary functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Search functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Test functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Shift functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Jump functions
// ---------------------------------------------------------------------------------------------------------------------

// DJZ -----------------------------------------------------------------------------------------------------------------

const fDJZ = 071
const jDJZ = 016

func djzSourceItemHIBD(label string, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fDJZ, jDJZ, a, x, h, i, b, d})
}

func djzSourceItemHIBDRef(label string, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fDJZ),
		fmt.Sprintf("%o", jDJZ),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func djzSourceItemHIU(label string, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fDJZ, jDJZ, a, x, h, i, u})
}

func djzSourceItemHIURef(label string, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fDJZ),
		fmt.Sprintf("%o", jDJZ),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// HKJ -----------------------------------------------------------------------------------------------------------------

const fHKJ = 074
const jHKJ = 005

func hkjSourceItemHIU(label string, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fHKJ, jHKJ, a, x, h, i, u})
}

func hkjSourceItemHIURef(label string, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fHKJ),
		fmt.Sprintf("%o", jHKJ),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// HLTJ ----------------------------------------------------------------------------------------------------------------

const fHLTJ = 074
const jHLTJ = 015
const aHLTJ = 005

func hltjSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fHLTJ, jHLTJ, aHLTJ, x, h, i, b, d})
}

func hltjSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fHLTJ),
		fmt.Sprintf("%o", jHLTJ),
		fmt.Sprintf("%o", aHLTJ),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func hltjSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fHLTJ, jHLTJ, aHLTJ, x, h, i, u})
}

func hltjSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fHLTJ),
		fmt.Sprintf("%o", jHLTJ),
		fmt.Sprintf("%o", aHLTJ),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// J -------------------------------------------------------------------------------------------------------------------

const fJ = 074
const jJBasic = 004
const jJExtended = 015
const aJBasic = 000
const aJExtended = 004

func jSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJ, jJExtended, aJExtended, x, h, i, b, d})
}

func jSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJ),
		fmt.Sprintf("%o", jJExtended),
		fmt.Sprintf("%o", aJExtended),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJ, jJBasic, aJBasic, x, h, i, u})
}

func jSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJ),
		fmt.Sprintf("%o", jJBasic),
		fmt.Sprintf("%o", aJBasic),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JC ------------------------------------------------------------------------------------------------------------------

const fJC = 074
const jJCBasic = 016
const jJCExtended = 014
const aJCExtended = 004

func jcSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJC, jJCExtended, aJCExtended, x, h, i, b, d})
}

func jcSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJC),
		fmt.Sprintf("%o", jJCExtended),
		fmt.Sprintf("%o", aJCExtended),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jcSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJC, jJCBasic, 0, x, h, i, u})
}

func jcSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJC),
		fmt.Sprintf("%o", jJCBasic),
		fmt.Sprintf("%o", "0"),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JDF -----------------------------------------------------------------------------------------------------------------

const fJDF = 074
const jJDF = 014
const aJDF = 003

func jdfSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJDF, jJDF, aJDF, x, h, i, b, d})
}

func jdfSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJDF),
		fmt.Sprintf("%o", jJDF),
		fmt.Sprintf("%o", aJDF),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jdfSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJDF, jJDF, aJDF, x, h, i, u})
}

func jdfSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJDF),
		fmt.Sprintf("%o", jJDF),
		fmt.Sprintf("%o", aJDF),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JFO -----------------------------------------------------------------------------------------------------------------

const fJFO = 074
const jJFO = 014
const aJFO = 002

func jfoSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJFO, jJFO, aJFO, x, h, i, b, d})
}

func jfoSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJFO),
		fmt.Sprintf("%o", jJFO),
		fmt.Sprintf("%o", aJFO),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jfoSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJFO, jJFO, aJFO, x, h, i, u})
}

func jfoSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJFO),
		fmt.Sprintf("%o", jJFO),
		fmt.Sprintf("%o", aJFO),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JFU -----------------------------------------------------------------------------------------------------------------

const fJFU = 074
const jJFU = 014
const aJFU = 001

func jfuSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJFU, jJFU, aJFU, x, h, i, b, d})
}

func jfuSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJFU),
		fmt.Sprintf("%o", jJFU),
		fmt.Sprintf("%o", aJFU),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jfuSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJFU, jJFU, aJFU, x, h, i, u})
}

func jfuSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJFU),
		fmt.Sprintf("%o", jJFU),
		fmt.Sprintf("%o", aJFU),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JK ------------------------------------------------------------------------------------------------------------------

const fJK = 074
const jJK = 004

func jkSourceItemHIU(label string, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJK, jJK, a, x, h, i, u})
}

func jkSourceItemHIURef(label string, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJK),
		fmt.Sprintf("%o", jJK),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JN ------------------------------------------------------------------------------------------------------------------

const fJN = 074
const jJN = 003

func jnSourceItemHIBD(label string, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJN, jJN, a, x, h, i, b, d})
}

func jnSourceItemHIBDRef(label string, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJN),
		fmt.Sprintf("%o", jJN),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnSourceItemHIU(label string, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJN, jJN, a, x, h, i, u})
}

func jnSourceItemHIURef(label string, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJN),
		fmt.Sprintf("%o", jJN),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNC -----------------------------------------------------------------------------------------------------------------

const fJNC = 074
const jJNCBasic = 017
const jJNCExtended = 014
const aJNCExtended = 005

func jncSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNC, jJNCExtended, aJNCExtended, x, h, i, b, d})
}

func jncSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNC),
		fmt.Sprintf("%o", jJNCExtended),
		fmt.Sprintf("%o", aJNCExtended),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jncSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNC, jJNCBasic, 0, x, h, i, u})
}

func jncSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNC),
		fmt.Sprintf("%o", jJNCBasic),
		fmt.Sprintf("%o", "0"),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNDF ----------------------------------------------------------------------------------------------------------------

const fJNDF = 074
const jJNDF = 015
const aJNDF = 003

func jndfSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNDF, jJNDF, aJNDF, x, h, i, b, d})
}

func jndfSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNDF),
		fmt.Sprintf("%o", jJNDF),
		fmt.Sprintf("%o", aJNDF),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jndfSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNDF, jJNDF, aJNDF, x, h, i, u})
}

func jndfSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNDF),
		fmt.Sprintf("%o", jJNDF),
		fmt.Sprintf("%o", aJNDF),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNFO ----------------------------------------------------------------------------------------------------------------

const fJNFO = 074
const jJNFO = 015
const aJNFO = 002

func jnfoSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNFO, jJNFO, aJNFO, x, h, i, b, d})
}

func jnfoSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNFO),
		fmt.Sprintf("%o", jJNFO),
		fmt.Sprintf("%o", aJNFO),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnfoSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNFO, jJNFO, aJNFO, x, h, i, u})
}

func jnfoSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNFO),
		fmt.Sprintf("%o", jJNFO),
		fmt.Sprintf("%o", aJNFO),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNFU ----------------------------------------------------------------------------------------------------------------

const fJNFU = 074
const jJNFU = 015
const aJNFU = 001

func jnfuSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNFU, jJNFU, aJNFU, x, h, i, b, d})
}

func jnfuSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNFU),
		fmt.Sprintf("%o", jJNFU),
		fmt.Sprintf("%o", aJNFU),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnfuSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNFU, jJNFU, aJNFU, x, h, i, u})
}

func jnfuSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNFU),
		fmt.Sprintf("%o", jJNFU),
		fmt.Sprintf("%o", aJNFU),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNO -----------------------------------------------------------------------------------------------------------------

const fJNO = 074
const jJNO = 015
const aJNO = 000

func jnoSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNO, jJNO, aJNO, x, h, i, b, d})
}

func jnoSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNO),
		fmt.Sprintf("%o", jJNO),
		fmt.Sprintf("%o", aJNO),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnoSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNO, jJNO, aJNO, x, h, i, u})
}

func jnoSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNO),
		fmt.Sprintf("%o", jJNO),
		fmt.Sprintf("%o", aJNO),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNZ -----------------------------------------------------------------------------------------------------------------

const fJNZ = 074
const jJNZ = 000

func jnzSourceItemHIBD(label string, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNZ, jJNZ, a, x, h, i, b, d})
}

func jnzSourceItemHIBDRef(label string, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNZ),
		fmt.Sprintf("%o", jJNZ),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnzSourceItemHIU(label string, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNZ, jJNZ, a, x, h, i, u})
}

func jnzSourceItemHIURef(label string, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJNZ),
		fmt.Sprintf("%o", jJNZ),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JO ------------------------------------------------------------------------------------------------------------------

const fJO = 074
const jJO = 014
const aJO = 000

func joSourceItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJO, jJO, aJO, x, h, i, b, d})
}

func joSourceItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJO),
		fmt.Sprintf("%o", jJO),
		fmt.Sprintf("%o", aJO),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func joSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJO, jJO, aJO, x, h, i, u})
}

func joSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJO),
		fmt.Sprintf("%o", jJO),
		fmt.Sprintf("%o", aJO),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JP ------------------------------------------------------------------------------------------------------------------

const fJP = 074
const jJP = 002

func jpSourceItemHIBD(label string, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJP, jJP, a, x, h, i, b, d})
}

func jpSourceItemHIBDRef(label string, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJP),
		fmt.Sprintf("%o", jJP),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jpSourceItemHIU(label string, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJP, jJP, a, x, h, i, u})
}

func jpSourceItemHIURef(label string, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJP),
		fmt.Sprintf("%o", jJP),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JZ ------------------------------------------------------------------------------------------------------------------

const fJZ = 074
const jJZ = 000

func jzSourceItemHIBD(label string, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJZ, jJZ, a, x, h, i, b, d})
}

func jzSourceItemHIBDRef(label string, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJZ),
		fmt.Sprintf("%o", jJZ),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jzSourceItemHIU(label string, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJZ, jJZ, a, x, h, i, u})
}

func jzSourceItemHIURef(label string, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fJZ),
		fmt.Sprintf("%o", jJZ),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// LMJ -----------------------------------------------------------------------------------------------------------------

const fLMJ = 074
const jLMJ = 013

func lmjSourceItemHIBD(label string, a int, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLMJ, jLMJ, a, x, h, i, b, d})
}

func lmjSourceItemHIBDRef(label string, a int, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLMJ),
		fmt.Sprintf("%o", jLMJ),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func lmjSourceItemHIU(label string, a int, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLMJ, jLMJ, a, x, h, i, u})
}

func lmjSourceItemHIURef(label string, a int, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fLMJ),
		fmt.Sprintf("%o", jLMJ),
		fmt.Sprintf("%o", a),
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func lmjSourceItemU(label string, a int, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLMJ, jLMJ, a, x, u})
}

// SLJ -----------------------------------------------------------------------------------------------------------------

const fSLJ = 072
const jSLJ = 001

func sljSourceItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fSLJ, jLMJ, 0, x, h, i, u})
}

func sljSourceItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fSLJ),
		fmt.Sprintf("%o", jLMJ),
		"0",
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func sljSourceItemU(label string, x int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fSLJ, jLMJ, 0, x, u})
}

// ---------------------------------------------------------------------------------------------------------------------
// Logical functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Storage-to-storage functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// String functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Address Space Management functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Procedure Control functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Queuing functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Activity Control functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Stack functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Interrupt functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// System Control functions
// ---------------------------------------------------------------------------------------------------------------------

// IAR -----------------------------------------------------------------------------------------------------------------

const fIAR = 073
const jIAR = 017
const aIAR = 006

func iarSourceItem(label string, uField int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fIAR, jIAR, aIAR, 0, uField})
}

func segSourceItem(segIndex int) *tasm.SourceItem {
	return tasm.NewSourceItem("", ".SEG", []string{fmt.Sprintf("%d", segIndex)})
}

// ---------------------------------------------------------------------------------------------------------------------
// Dayclock functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// UPI functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// System Instrumentation functions
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// Special functions
// ---------------------------------------------------------------------------------------------------------------------

// NOP ----------------------------------------------------------------------------------------------------------------0

const fNOPBasic = 074
const fNOPExtended = 703
const jNOPBasic = 006
const jNOPExtended = 014
const aNOP = 000

func nopItemHIU(label string, x int, h int, i int, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fNOPBasic, jNOPBasic, 0, x, h, i, u})
}

func nopItemHIURef(label string, x int, h int, i int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fNOPBasic),
		fmt.Sprintf("%o", jNOPBasic),
		"0",
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func nopItemHIBD(label string, x int, h int, i int, b int, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fNOPExtended, jNOPExtended, aNOP, x, h, i, b, d})
}

func nopItemHIBDRef(label string, x int, h int, i int, b int, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%o", fNOPBasic),
		fmt.Sprintf("%o", jNOPBasic),
		"0",
		fmt.Sprintf("%o", x),
		fmt.Sprintf("%o", h),
		fmt.Sprintf("%o", i),
		fmt.Sprintf("%o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// ---------------------------------------------------------------------------------------------------------------------
// Convenience methods
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------

func checkInterrupt(t *testing.T, engine *InstructionEngine, interruptClass pkg.InterruptClass) {
	for _, i := range engine.pendingInterrupts.stack {
		if i.GetClass() == interruptClass {
			return
		}
	}

	engine.pendingInterrupts.Dump()
	t.Fatalf("Error:Expected interrupt class %d to be posted", interruptClass)
}

func checkProgramAddress(t *testing.T, engine *InstructionEngine, expectedAddress uint64) {
	actual := engine.GetProgramAddressRegister().GetProgramCounter()
	if actual != expectedAddress {
		t.Fatalf("Error:Expected PAR.PC is %06o but we expected it to be %06o", actual, expectedAddress)
	}
}

func checkMemory(t *testing.T, engine *InstructionEngine, addr *pkg.AbsoluteAddress, offset uint64, expected uint64) {
	seg, interrupt := engine.mainStorage.GetSegment(addr.GetSegment())
	if interrupt != nil {
		engine.mainStorage.Dump()
		t.Fatalf("Error:%s", pkg.GetInterruptString(interrupt))
	}

	if addr.GetOffset() >= uint64(len(seg)) {
		engine.mainStorage.Dump()
		t.Fatalf("Error:offset is out of range for address %s - segment size is %012o", addr.GetString(), len(seg))
	}

	result := seg[addr.GetOffset()+offset]
	if result.GetW() != expected {
		engine.mainStorage.Dump()
		t.Fatalf("Storage at %s+0%o is %012o, expected %012o", addr.GetString(), offset, result, expected)
	}
}

func checkRegister(t *testing.T, engine *InstructionEngine, register uint64, expected uint64, name string) {
	result := engine.generalRegisterSet.GetRegister(register).GetW()
	if result != expected {
		engine.generalRegisterSet.Dump()
		t.Fatalf("Register %s is %012o, expected %012o", name, result, expected)
	}
}

// TODO following is deprecated, use checkStoppedReason() instead, so we don't have to check the PAR.PC
func checkStopped(t *testing.T, engine *InstructionEngine) {
	if engine.HasPendingInterrupt() {
		engine.Dump()
		t.Fatalf("Engine has unexpected pending interrupts")
	}

	if !engine.IsStopped() {
		engine.Dump()
		t.Fatalf("Expected engine to be stopped; it is not")
	}
}

func checkStoppedReason(t *testing.T, engine *InstructionEngine, reason StopReason, detail uint64) {
	if engine.HasPendingInterrupt() {
		engine.Dump()
		t.Fatalf("Engine has unexpected pending interrupts")
	}

	if !engine.IsStopped() {
		engine.Dump()
		t.Fatalf("Expected engine to be stopped; it is not")
	}

	actualReason, actualDetail := engine.GetStopReason()
	if actualReason != reason {
		engine.Dump()
		t.Fatalf("Engine stopped for reason %d; expected reason %d", actualReason, reason)
	}

	if actualDetail != detail {
		engine.Dump()
		t.Fatalf("Engine stopped for detail %d; expected detail %d", actualDetail, detail)
	}
}
