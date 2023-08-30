// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"khalehla/tasm"
	"testing"
)

const fDJZ = 071
const jDJZ = 016
const fJC = 074
const jJCBasic = 016
const jJCExtended = 014
const aJCExtended = 004
const fJDF = 074
const jJDF = 014
const aJDF = 003
const fJFO = 074
const jJFO = 014
const aJFO = 002
const fJFU = 074
const jJFU = 014
const aJFU = 001
const fJN = 074
const jJN = 003
const fJNC = 074
const jJNCBasic = 017
const jJNCExtended = 014
const aJNCExtended = 005
const fJNDF = 074
const jJNDF = 015
const aJNDF = 003
const fJNFO = 074
const jJNFO = 015
const aJNFO = 002
const fJNFU = 074
const jJNFU = 015
const aJNFU = 001
const fJNO = 074
const jJNO = 015
const aJNO = 000
const fJNZ = 074
const jJNZ = 001
const fJO = 074
const jJO = 014
const aJO = 000
const fJP = 074
const jJP = 002
const fJZ = 074
const jJZ = 000

// ---------------------------------------------------
// DJZ

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

// ---------------------------------------------------
// JC

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

// ---------------------------------------------------
// JDF

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

// ---------------------------------------------------
// JFO

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

// ---------------------------------------------------
// JFU

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

// ---------------------------------------------------
// JN

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

// ---------------------------------------------------
// JNC

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

// ---------------------------------------------------
// JNDF

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

// ---------------------------------------------------
// JNFO

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

// ---------------------------------------------------
// JNFU

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

// ---------------------------------------------------
// JNO

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

// ---------------------------------------------------
// JNZ

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

// ---------------------------------------------------
// JO

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

// ---------------------------------------------------
// JP

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

// ---------------------------------------------------
// JZ

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

// ---------------------------------------------------------------------------------------------------------------------

var jumpZeroExtendedPosZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA5, 0, 0),
	jzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JZ_Extended_PosZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpZeroExtendedPosZero)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

var jumpZeroExtendedNegZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jXU, regA5, 0, 0_777777),
	jzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JZ_Extended_NegZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpZeroExtendedNegZero)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

var jumpZeroExtendedNotZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA5, 0, 01),
	jzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JZ_Extended_NotZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpZeroExtendedNotZero)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
}

var doubleJumpZeroExtendedMode = []*tasm.SourceItem{
	segSourceItem(3),
	tasm.NewSourceItem("posZero", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("negZero", "w", []string{"0777777777777"}),
	tasm.NewSourceItem("", "w", []string{"0777777777777"}),
	tasm.NewSourceItem("notZero1", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"011"}),
	tasm.NewSourceItem("notZero2", "w", []string{"0777777777777"}),
	tasm.NewSourceItem("", "w", []string{"0"}),

	segSourceItem(0),
	dlSourceItemHIBRef("", 0, 0, 0, 0, 3, "posZero"),
	djzSourceItemHIBDRef("", 0, 0, 0, 0, 0, "target1"),
	iarSourceItem(1),

	dlSourceItemHIBRef("target1", 2, 0, 0, 0, 3, "negZero"),
	djzSourceItemHIBDRef("", 2, 0, 0, 0, 0, "target2"),
	iarSourceItem(2),

	dlSourceItemHIBRef("target2", 4, 0, 0, 0, 3, "notZero1"),
	djzSourceItemHIBDRef("", 4, 0, 0, 0, 0, "bad3"),

	dlSourceItemHIBRef("target3", 6, 0, 0, 0, 3, "notZero2"),
	djzSourceItemHIBDRef("", 6, 0, 0, 0, 0, "bad4"),
	jSourceItemExtended("", 0, 0, 0, "end"),

	iarSourceItem("bad3", 3),
	iarSourceItem("bad4", 4),
	iarSourceItem("end", 0),
}

func Test_DJZ_Extended_PosZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", doubleJumpZeroExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

var jumpNonZeroExtendedPosZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA5, 0, 0),
	jnzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNZ_Extended_PosZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNonZeroExtendedPosZero)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
}

var jumpNonZeroExtendedNegZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jXU, regA5, 0, 0_777777),
	jnzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNZ_Extended_NegZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNonZeroExtendedNegZero)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
}

var jumpNonZeroExtendedNotZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA5, 0, 1),
	jnzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNZ_Extended_NotZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNonZeroExtendedNotZero)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

var jumpPosNegExtended = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA10, 0, 0),
	jpSourceItemHIBDRef("", 10, 0, 0, 0, 0, "target1"),
	iarSourceItem("bad1", 1),
	jnSourceItemHIBDRef("target1", 10, 0, 0, 0, 0, "bad2"),

	nopItemHIBD("", 0, 0, 0, 0, 0),
	laSourceItemU(jXU, regA10, 0, 0_444444),
	jpSourceItemHIBDRef("", 10, 0, 0, 0, 0, "bad3"),
	jnSourceItemHIBDRef("", 10, 0, 0, 0, 0, "end"),

	iarSourceItem("bad4", 4),
	iarSourceItem("bad2", 2),
	iarSourceItem("bad3", 3),

	labelSourceItem("end"),
	iarSourceItem(0),
}

func Test_JP_JN_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpPosNegExtended)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

//	TODO JPS, JNS
//	TODO JB, JNB
//	TODO JGD, JMGI

// Conditional based on designator register bits -----------------------------------------------------------------------

var jumpCarryBasic = []*tasm.SourceItem{
	segSourceItem(0),
	jcSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JC_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpCarryBasic)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCarry(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JC_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpCarryBasic)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpCarryExtended = []*tasm.SourceItem{
	segSourceItem(0),
	jcSourceItemHIBDRef("", 0, 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JC_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpCarryExtended)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCarry(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

var jumpDivideFault = []*tasm.SourceItem{
	segSourceItem(0),
	jdfSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JDF_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JDF_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JDF_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JDF_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpFloatingOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	jfoSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JFO_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JFO_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JFO_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JFO_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpFloatingUnderflow = []*tasm.SourceItem{
	segSourceItem(0),
	jfuSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JFU_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JFU_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JFU_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JFU_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

//	--------------------------------------------------------------------------------------------------------------------

var jumpNoCarryBasic = []*tasm.SourceItem{
	segSourceItem(0),
	jncSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JNC_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoCarryBasic)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCarry(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNC_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoCarryExtended = []*tasm.SourceItem{
	segSourceItem(0),
	jncSourceItemHIBDRef("", 0, 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JNC_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoCarryExtended)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCarry(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

var jumpNoDivideFault = []*tasm.SourceItem{
	segSourceItem(0),
	jndfSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JNDF_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNDF_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JNDF_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNDF_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoFloatingOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	jnfoSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JNFO_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNFO_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JNFO_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNFO_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoFloatingUnderflow = []*tasm.SourceItem{
	segSourceItem(0),
	jnfuSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JNFU_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNFU_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JNFU_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNFU_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	jnoSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JNO_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNO_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	joSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem("target", 0),
}

func Test_JO_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JO_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}
