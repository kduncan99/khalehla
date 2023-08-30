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

// A-field/X-field values for registers
const (
	regX0 = iota
	regX1
	regX2
	regX3
	regX4
	regX5
	regX6
	regX7
	regX8
	regX9
	regX10
	regX11
	regX12
	regX13
	regX14
	regX15
)

const (
	regA0 = iota
	regA1
	regA2
	regA3
	regA4
	regA5
	regA6
	regA7
	regA8
	regA9
	regA10
	regA11
	regA12
	regA13
	regA14
	regA15
)

const (
	regR0 = iota
	regR1
	regR2
	regR3
	regR4
	regR5
	regR6
	regR7
	regR8
	regR9
	regR10
	regR11
	regR12
	regR13
	regR14
	regR15
)

// ---------------------------------------------------------------------------------------------------------------------

func grsRef(grsRegIndex uint64) string {
	return fmt.Sprintf("0%o", grsRegIndex)
}

func sourceItem(label string, operator string, operands []int) *tasm.SourceItem {
	strOps := make([]string, len(operands))
	for ox := 0; ox < len(operands); ox++ {
		strOps[ox] = fmt.Sprintf("0%o", operands[ox])
	}

	return tasm.NewSourceItem(label, operator, strOps)
}

func labelSourceItem(label string) *tasm.SourceItem {
	return tasm.NewSourceItem(label, "", []string{})
}

func labelDataSourceItem(label string, values []uint64) *tasm.SourceItem {
	var operator string
	if len(values) == 1 {
		operator = "w"
	} else if len(values) == 2 {
		operator = "hw"
	} else if len(values) == 3 {
		operator = "tw"
	} else if len(values) == 4 {
		operator = "qw"
	} else if len(values) == 6 {
		operator = "sw"
	} else {
		operator = "?"
	}

	strValues := make([]string, len(values))
	for vx := 0; vx < len(values); vx++ {
		strValues[vx] = fmt.Sprintf("0%o", values[vx])
	}

	return tasm.NewSourceItem(label, operator, strValues)
}

func dataSourceItem(values []uint64) *tasm.SourceItem {
	return labelDataSourceItem("", values)
}

func fjaxuSourceItem(f uint64, j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", f),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", u),
	}
	return tasm.NewSourceItem("", "fjaxu", ops)
}

func fjaxRefSourceItem(f uint64, j uint64, a uint64, x uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", f),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		ref,
	}
	return tasm.NewSourceItem("", "fjaxu", ops)
}

func fjaxhibRefSourceItem(f uint64, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", f),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem("", "fjaxhibd", ops)
}

func fjaxhiRefSourceItem(f uint64, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", f),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem("", "fjaxhiu", ops)
}

func segSourceItem(segIndex int) *tasm.SourceItem {
	return tasm.NewSourceItem("", ".SEG", []string{fmt.Sprintf("%d", segIndex)})
}

// ---------------------------------------------------------------------------------------------------------------------
// Jump functions
// ---------------------------------------------------------------------------------------------------------------------

// DJZ -----------------------------------------------------------------------------------------------------------------

const fDJZ = 071
const jDJZ = 016

func djzSourceItemHIBD(label string, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fDJZ, jDJZ, a, x, h, i, b, d})
}

func djzSourceItemHIBDRef(label string, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fDJZ),
		fmt.Sprintf("%03o", jDJZ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func djzSourceItemHIU(label string, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fDJZ, jDJZ, a, x, h, i, u})
}

func djzSourceItemHIURef(label string, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fDJZ),
		fmt.Sprintf("%03o", jDJZ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// HKJ -----------------------------------------------------------------------------------------------------------------

const fHKJ = 074
const jHKJ = 005

func hkjSourceItemHIU(label string, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fHKJ, jHKJ, a, x, h, i, u})
}

func hkjSourceItemHIURef(label string, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fHKJ),
		fmt.Sprintf("%03o", jHKJ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// HLTJ ----------------------------------------------------------------------------------------------------------------

const fHLTJ = 074
const jHLTJ = 015
const aHLTJ = 005

func hltjSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fHLTJ, jHLTJ, aHLTJ, x, h, i, b, d})
}

func hltjSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fHLTJ),
		fmt.Sprintf("%03o", jHLTJ),
		fmt.Sprintf("%03o", aHLTJ),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func hltjSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fHLTJ, jHLTJ, aHLTJ, x, h, i, u})
}

func hltjSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fHLTJ),
		fmt.Sprintf("%03o", jHLTJ),
		fmt.Sprintf("%03o", aHLTJ),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
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

func jSourceItemBasic(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJ),
		fmt.Sprintf("%03o", jJBasic),
		fmt.Sprintf("%03o", aJBasic),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func jSourceItemExtended(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJ),
		fmt.Sprintf("%03o", jJExtended),
		fmt.Sprintf("%03o", aJExtended),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JC ------------------------------------------------------------------------------------------------------------------

const fJC = 074
const jJCBasic = 016
const jJCExtended = 014
const aJCExtended = 004

func jcSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJC, jJCExtended, aJCExtended, x, h, i, b, d})
}

func jcSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJC),
		fmt.Sprintf("%03o", jJCExtended),
		fmt.Sprintf("%03o", aJCExtended),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jcSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJC, jJCBasic, 0, x, h, i, u})
}

func jcSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJC),
		fmt.Sprintf("%03o", jJCBasic),
		"0",
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JDF -----------------------------------------------------------------------------------------------------------------

const fJDF = 074
const jJDF = 014
const aJDF = 003

func jdfSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJDF, jJDF, aJDF, x, h, i, b, d})
}

func jdfSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJDF),
		fmt.Sprintf("%03o", jJDF),
		fmt.Sprintf("%03o", aJDF),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jdfSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJDF, jJDF, aJDF, x, h, i, u})
}

func jdfSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJDF),
		fmt.Sprintf("%03o", jJDF),
		fmt.Sprintf("%03o", aJDF),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JFO -----------------------------------------------------------------------------------------------------------------

const fJFO = 074
const jJFO = 014
const aJFO = 002

func jfoSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJFO, jJFO, aJFO, x, h, i, b, d})
}

func jfoSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJFO),
		fmt.Sprintf("%03o", jJFO),
		fmt.Sprintf("%03o", aJFO),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jfoSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJFO, jJFO, aJFO, x, h, i, u})
}

func jfoSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJFO),
		fmt.Sprintf("%03o", jJFO),
		fmt.Sprintf("%03o", aJFO),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JFU -----------------------------------------------------------------------------------------------------------------

const fJFU = 074
const jJFU = 014
const aJFU = 001

func jfuSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJFU, jJFU, aJFU, x, h, i, b, d})
}

func jfuSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJFU),
		fmt.Sprintf("%03o", jJFU),
		fmt.Sprintf("%03o", aJFU),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jfuSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJFU, jJFU, aJFU, x, h, i, u})
}

func jfuSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJFU),
		fmt.Sprintf("%03o", jJFU),
		fmt.Sprintf("%03o", aJFU),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JK ------------------------------------------------------------------------------------------------------------------

const fJK = 074
const jJK = 004

func jkSourceItemHIU(label string, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJK, jJK, a, x, h, i, u})
}

func jkSourceItemHIURef(label string, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJK),
		fmt.Sprintf("%03o", jJK),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JN ------------------------------------------------------------------------------------------------------------------

const fJN = 074
const jJN = 003

func jnSourceItemHIBD(label string, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJN, jJN, a, x, h, i, b, d})
}

func jnSourceItemHIBDRef(label string, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJN),
		fmt.Sprintf("%03o", jJN),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnSourceItemHIU(label string, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJN, jJN, a, x, h, i, u})
}

func jnSourceItemHIURef(label string, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJN),
		fmt.Sprintf("%03o", jJN),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNC -----------------------------------------------------------------------------------------------------------------

const fJNC = 074
const jJNCBasic = 017
const jJNCExtended = 014
const aJNCExtended = 005

func jncSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNC, jJNCExtended, aJNCExtended, x, h, i, b, d})
}

func jncSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNC),
		fmt.Sprintf("%03o", jJNCExtended),
		fmt.Sprintf("%03o", aJNCExtended),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jncSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNC, jJNCBasic, 0, x, h, i, u})
}

func jncSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNC),
		fmt.Sprintf("%03o", jJNCBasic),
		"0",
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNDF ----------------------------------------------------------------------------------------------------------------

const fJNDF = 074
const jJNDF = 015
const aJNDF = 003

func jndfSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNDF, jJNDF, aJNDF, x, h, i, b, d})
}

func jndfSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNDF),
		fmt.Sprintf("%03o", jJNDF),
		fmt.Sprintf("%03o", aJNDF),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jndfSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNDF, jJNDF, aJNDF, x, h, i, u})
}

func jndfSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNDF),
		fmt.Sprintf("%03o", jJNDF),
		fmt.Sprintf("%03o", aJNDF),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNFO ----------------------------------------------------------------------------------------------------------------

const fJNFO = 074
const jJNFO = 015
const aJNFO = 002

func jnfoSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNFO, jJNFO, aJNFO, x, h, i, b, d})
}

func jnfoSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNFO),
		fmt.Sprintf("%03o", jJNFO),
		fmt.Sprintf("%03o", aJNFO),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnfoSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNFO, jJNFO, aJNFO, x, h, i, u})
}

func jnfoSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNFO),
		fmt.Sprintf("%03o", jJNFO),
		fmt.Sprintf("%03o", aJNFO),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNFU ----------------------------------------------------------------------------------------------------------------

const fJNFU = 074
const jJNFU = 015
const aJNFU = 001

func jnfuSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNFU, jJNFU, aJNFU, x, h, i, b, d})
}

func jnfuSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNFU),
		fmt.Sprintf("%03o", jJNFU),
		fmt.Sprintf("%03o", aJNFU),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnfuSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNFU, jJNFU, aJNFU, x, h, i, u})
}

func jnfuSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNFU),
		fmt.Sprintf("%03o", jJNFU),
		fmt.Sprintf("%03o", aJNFU),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNO -----------------------------------------------------------------------------------------------------------------

const fJNO = 074
const jJNO = 015
const aJNO = 000

func jnoSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNO, jJNO, aJNO, x, h, i, b, d})
}

func jnoSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNO),
		fmt.Sprintf("%03o", jJNO),
		fmt.Sprintf("%03o", aJNO),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnoSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNO, jJNO, aJNO, x, h, i, u})
}

func jnoSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNO),
		fmt.Sprintf("%03o", jJNO),
		fmt.Sprintf("%03o", aJNO),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JNZ -----------------------------------------------------------------------------------------------------------------

const fJNZ = 074
const jJNZ = 001

func jnzSourceItemHIBD(label string, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJNZ, jJNZ, a, x, h, i, b, d})
}

func jnzSourceItemHIBDRef(label string, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNZ),
		fmt.Sprintf("%03o", jJNZ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jnzSourceItemHIU(label string, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJNZ, jJNZ, a, x, h, i, u})
}

func jnzSourceItemHIURef(label string, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJNZ),
		fmt.Sprintf("%03o", jJNZ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JO ------------------------------------------------------------------------------------------------------------------

const fJO = 074
const jJO = 014
const aJO = 000

func joSourceItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJO, jJO, aJO, x, h, i, b, d})
}

func joSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJO),
		fmt.Sprintf("%03o", jJO),
		fmt.Sprintf("%03o", aJO),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func joSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJO, jJO, aJO, x, h, i, u})
}

func joSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJO),
		fmt.Sprintf("%03o", jJO),
		fmt.Sprintf("%03o", aJO),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JP ------------------------------------------------------------------------------------------------------------------

const fJP = 074
const jJP = 002

func jpSourceItemHIBD(label string, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJP, jJP, a, x, h, i, b, d})
}

func jpSourceItemHIBDRef(label string, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJP),
		fmt.Sprintf("%03o", jJP),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jpSourceItemHIU(label string, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJP, jJP, a, x, h, i, u})
}

func jpSourceItemHIURef(label string, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJP),
		fmt.Sprintf("%03o", jJP),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// JZ ------------------------------------------------------------------------------------------------------------------

const fJZ = 074
const jJZ = 000

func jzSourceItemHIBD(label string, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fJZ, jJZ, a, x, h, i, b, d})
}

func jzSourceItemHIBDRef(label string, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJZ),
		fmt.Sprintf("%03o", jJZ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func jzSourceItemHIU(label string, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fJZ, jJZ, a, x, h, i, u})
}

func jzSourceItemHIURef(label string, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fJZ),
		fmt.Sprintf("%03o", jJZ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// LMJ -----------------------------------------------------------------------------------------------------------------

const fLMJ = 074
const jLMJ = 013

func lmjSourceItemHIBD(label string, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fLMJ, jLMJ, a, x, h, i, b, d})
}

func lmjSourceItemHIBDRef(label string, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fLMJ),
		fmt.Sprintf("%03o", jLMJ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func lmjSourceItemHIU(label string, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fLMJ, jLMJ, a, x, h, i, u})
}

func lmjSourceItemHIURef(label string, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fLMJ),
		fmt.Sprintf("%03o", jLMJ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func lmjSourceItemU(label string, a uint64, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fLMJ, jLMJ, a, x, u})
}

// SLJ -----------------------------------------------------------------------------------------------------------------

const fSLJ = 072
const jSLJ = 001

func sljSourceItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fSLJ, jSLJ, 0, x, h, i, u})
}

func sljSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fSLJ),
		fmt.Sprintf("%03o", jSLJ),
		"0",
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func sljSourceItemU(label string, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fSLJ, jLMJ, 0, x, u})
}

// ---------------------------------------------------------------------------------------------------------------------
// Logical functions
// ---------------------------------------------------------------------------------------------------------------------

// OR ------------------------------------------------------------------------------------------------------------------

const fOR = 040

func orSourceItemHIBD(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fOR, j, a, x, h, i, b, d})
}

func orSourceItemHIBDRef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fOR),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func orSourceItemHIU(label string, j uint64, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fOR, j, a, x, h, i, u})
}

func orSourceItemHIURef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fOR),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func orSourceItemU(label string, j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fOR, j, a, x, u})
}

// XOR -----------------------------------------------------------------------------------------------------------------

const fXOR = 041

func xorSourceItemHIBD(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fXOR, j, a, x, h, i, b, d})
}

func xorSourceItemHIBDRef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fXOR),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func xorSourceItemHIU(label string, j uint64, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fXOR, j, a, x, h, i, u})
}

func xorSourceItemHIURef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fXOR),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func xorSourceItemU(label string, j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fXOR, j, a, x, u})
}

// AND -----------------------------------------------------------------------------------------------------------------

const fAND = 042

func andSourceItemHIBD(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fAND, j, a, x, h, i, b, d})
}

func andSourceItemHIBDRef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fAND),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func andSourceItemHIU(label string, j uint64, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fAND, j, a, x, h, i, u})
}

func andSourceItemHIURef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fAND),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func andSourceItemU(label string, j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fAND, j, a, x, u})
}

// MLU -----------------------------------------------------------------------------------------------------------------

const fMLU = 043

func mluSourceItemHIBD(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fMLU, j, a, x, h, i, b, d})
}

func mluSourceItemHIBDRef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fMLU),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func mluSourceItemHIU(label string, j uint64, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fMLU, j, a, x, h, i, u})
}

func mluSourceItemHIURef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fMLU),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func mluSourceItemU(label string, j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fMLU, j, a, x, u})
}

// ---------------------------------------------------------------------------------------------------------------------
// Storage-to-storage functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO BT
//	TODO BIM
//	TODO BIC
//	TODO BIMT
//	TODO BICL
//	TODO BIML
//	TODO BN
//	TODO BBN

// ---------------------------------------------------------------------------------------------------------------------
// String functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LS
//	TODO LSA
//	TODO SS
//	TODO TES
//	TODO TNES

// ---------------------------------------------------------------------------------------------------------------------
// Address Space Management functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LBU
//	TODO LBE
//	TODO LBUD
//	TODO LBED
//	TODO SBUD
//	TODO SBED
//	TODO SBU
//	TODO LBN
//	TODO TRA
//	TODO TVA
//	TODO DABT
//	TODO TRARS

// ---------------------------------------------------------------------------------------------------------------------
// Procedure Control functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO GOTO
//	TODO CALL
//	TODO LOCL
//	TODO RTN
//	TODO LBJ
//	TODO LIJ
//	TODO LDJ

// ---------------------------------------------------------------------------------------------------------------------
// Queuing functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO ENQ
//	TODO ENQF
//	TODO DEQ
//	TODO DEQW
//	TODO DEPOSITQB
//	TODO WITHDRAWQB

// ---------------------------------------------------------------------------------------------------------------------
// Activity Control functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LD
//	TODO SD
//	TODO LPD
//	TODO SPD
//	TODO LUD
//	TODO SUD
//	TODO LAE
//	TODO UR
//	TODO ACEL
//	TODO DCEL
//	TODO SKQT
//	TODO KCHG

// ---------------------------------------------------------------------------------------------------------------------
// Stack functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO BUY
//	TODO SELL

// ---------------------------------------------------------------------------------------------------------------------
// Interrupt Control functions
// ---------------------------------------------------------------------------------------------------------------------

// ER ------------------------------------------------------------------------------------------------------------------

const fER = 072
const jER = 011

func erSourceItemHIRef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fER),
		fmt.Sprintf("%03o", jER),
		"0",
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func erSourceItemU(label string, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fER, jER, 0, u})
}

// SGNL ----------------------------------------------------------------------------------------------------------------

const fSGNL = 073
const jSGNL = 015
const aSGNL = 017

func sgnlSourceItemHIRef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fSGNL),
		fmt.Sprintf("%03o", jSGNL),
		fmt.Sprintf("%03o", aSGNL),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func sgnlSourceItemU(label string, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fSGNL, jSGNL, aSGNL, u})
}

// PAIJ ----------------------------------------------------------------------------------------------------------------

const fPAIJ = 074
const jPAIJ = 014
const aPAIJ = 007

func paijSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fPAIJ),
		fmt.Sprintf("%03o", jPAIJ),
		fmt.Sprintf("%03o", aPAIJ),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func paijSourceItemRef(label string, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fPAIJ),
		fmt.Sprintf("%03o", jPAIJ),
		fmt.Sprintf("%03o", aPAIJ),
		"0",
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxu", ops)
}

// AAIJ ----------------------------------------------------------------------------------------------------------------

const fAAIJ = 074
const jAAIJExtended = 014
const jAAIJBasic = 007
const aAAIJExtended = 006
const aAAIJBasic = 000

func aaijSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fAAIJ),
		fmt.Sprintf("%03o", jAAIJExtended),
		fmt.Sprintf("%03o", aAAIJExtended),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func aaijSourceItemRef(label string, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fAAIJ),
		fmt.Sprintf("%03o", jAAIJBasic),
		fmt.Sprintf("%03o", aAAIJBasic),
		"0",
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxu", ops)
}

// ---------------------------------------------------------------------------------------------------------------------
// System Control functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO SPID
//	TODO IPC
//	TODO SPC

// IAR -----------------------------------------------------------------------------------------------------------------

const fIAR = 073
const jIAR = 017
const aIAR = 006

func iarSourceItem(uField int) *tasm.SourceItem {
	return fjaxuSourceItem(fIAR, jIAR, aIAR, 0, uField)
}

// ---------------------------------------------------------------------------------------------------------------------
// Dayclock functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LRD
//	TODO SMD
//	TODO RMD
//	TODO LMC
//	TODO SDMN
//	TODO SDMS
//	TODO SDMF
//	TODO RDC

// ---------------------------------------------------------------------------------------------------------------------
// UPI functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO SEND
//	TODO ACK

// ---------------------------------------------------------------------------------------------------------------------
// System Instrumentation functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LBRX
//	TODO CJHE
//	TODO SJH

// ---------------------------------------------------------------------------------------------------------------------
// Special functions
// ---------------------------------------------------------------------------------------------------------------------

// DCB -----------------------------------------------------------------------------------------------------------------

const fDCB = 033
const jDCB = 015

func dcbSourceItemHIBDRef(label string, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fDCB),
		fmt.Sprintf("%03o", jDCB),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// EX ------------------------------------------------------------------------------------------------------------------

const fEXBasic = 072
const fEXExtended = 073
const jEXBasic = 010
const jEXExtended = 014
const aEXExtended = 005

func exSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fEXExtended),
		fmt.Sprintf("%03o", jEXExtended),
		fmt.Sprintf("%03o", aEXExtended),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func exSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fEXBasic),
		fmt.Sprintf("%03o", jEXBasic),
		fmt.Sprintf("%03o", 0),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// EXR -----------------------------------------------------------------------------------------------------------------

const fEXR = 073
const jEXR = 014
const aEXR = 006

func exrSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fEXR),
		fmt.Sprintf("%03o", jEXR),
		fmt.Sprintf("%03o", aEXR),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// NOP ----------------------------------------------------------------------------------------------------------------0

const fNOPBasic = 074
const fNOPExtended = 073
const jNOPBasic = 006
const jNOPExtended = 014
const aNOP = 000

func nopItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fNOPBasic, jNOPBasic, 0, x, h, i, u})
}

func nopItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fNOPBasic),
		fmt.Sprintf("%03o", jNOPBasic),
		"0",
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func nopItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fNOPExtended, jNOPExtended, aNOP, x, h, i, b, d})
}

func nopItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fNOPBasic),
		fmt.Sprintf("%03o", jNOPBasic),
		"0",
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// RNGB ----------------------------------------------------------------------------------------------------------------

const fRNGB = 037
const jRNGB = 004
const aRNGB = 006

func rngbSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fRNGB),
		fmt.Sprintf("%03o", jRNGB),
		fmt.Sprintf("%03o", aRNGB),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// RNGI ----------------------------------------------------------------------------------------------------------------

const fRNGI = 037
const jRNGI = 004
const aRNGI = 005

func rngiSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fRNGI),
		fmt.Sprintf("%03o", jRNGI),
		fmt.Sprintf("%03o", aRNGI),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// ---------------------------------------------------------------------------------------------------------------------
// Convenience methods
// ---------------------------------------------------------------------------------------------------------------------

func checkInterrupt(t *testing.T, engine *InstructionEngine, interruptClass pkg.InterruptClass) {
	for _, i := range engine.pendingInterrupts.stack {
		if i.GetClass() == interruptClass {
			return
		}
	}

	engine.pendingInterrupts.Dump()
	t.Errorf("Error:Expected interrupt class %d to be posted", interruptClass)
}

func checkInterruptAndSSF(
	t *testing.T,
	engine *InstructionEngine,
	interruptClass pkg.InterruptClass,
	shortStatusField pkg.InterruptShortStatus) {

	for _, i := range engine.pendingInterrupts.stack {
		if i.GetClass() == interruptClass {
			if i.GetShortStatusField() != shortStatusField {
				t.Errorf("Error:Found interrupt class %d but SSF was %d and we expected %d",
					interruptClass, i.GetShortStatusField(), shortStatusField)
			}
			return
		}
	}

	engine.pendingInterrupts.Dump()
	t.Errorf("Error:Expected interrupt class %d to be posted", interruptClass)
}

func checkProgramAddress(t *testing.T, engine *InstructionEngine, expectedAddress uint64) {
	actual := engine.GetProgramAddressRegister().GetProgramCounter()
	if actual != expectedAddress {
		t.Errorf("Error:Expected PAR.PC is %06o but we expected it to be %06o", actual, expectedAddress)
	}
}

func checkMemory(t *testing.T, engine *InstructionEngine, addr *pkg.AbsoluteAddress, offset uint64, expected uint64) {
	seg, interrupt := engine.mainStorage.GetSegment(addr.GetSegment())
	if interrupt != nil {
		engine.mainStorage.Dump()
		t.Errorf("Error:%s", pkg.GetInterruptString(interrupt))
	}

	if addr.GetOffset() >= uint64(len(seg)) {
		engine.mainStorage.Dump()
		t.Errorf("Error:offset is out of range for address %s - segment size is %012o", addr.GetString(), len(seg))
	}

	result := seg[addr.GetOffset()+offset]
	if result.GetW() != expected {
		engine.mainStorage.Dump()
		t.Errorf("Storage at (%s)+0%o is %012o, expected %012o", addr.GetString(), offset, result, expected)
	}
}

func checkRegister(t *testing.T, engine *InstructionEngine, regIndex uint64, expected uint64) {
	result := engine.generalRegisterSet.GetRegister(regIndex).GetW()
	if result != expected {
		engine.generalRegisterSet.Dump()
		t.Errorf("Register %s is %012o, expected %012o", pkg.RegisterNames[regIndex], result, expected)
	}
}

func checkStoppedReason(t *testing.T, engine *InstructionEngine, reason StopReason, detail uint64) {
	if engine.HasPendingInterrupt() {
		engine.Dump()
		t.Errorf("Engine has unexpected pending interrupts")
	}

	if !engine.IsStopped() {
		engine.Dump()
		t.Errorf("Expected engine to be stopped; it is not")
	}

	actualReason, actualDetail := engine.GetStopReason()
	if actualReason != reason {
		engine.Dump()
		t.Errorf("Engine stopped for reason %d; expected reason %d", actualReason, reason)
	}

	if actualDetail != detail {
		engine.Dump()
		t.Errorf("Engine stopped for detail %d; expected detail %d", actualDetail, detail)
	}
}
