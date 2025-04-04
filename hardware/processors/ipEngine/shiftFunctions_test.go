// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"testing"

	"khalehla/common"
	"khalehla/tasm"
)

const fSSC = 073
const fDSC = 073
const fSSL = 073
const fDSL = 073
const fSSA = 073
const fDSA = 073
const fLSC = 073
const fLSSC = 073
const fDLSC = 073
const fLDSC = 073
const fLSSL = 073
const fLDSL = 073

const jSSC = 000
const jDSC = 001
const jSSL = 002
const jDSL = 003
const jSSA = 004
const jDSA = 005
const jLSC = 006
const jDLSC = 007
const jLSSC = 010
const jLDSC = 011
const jLSSL = 012
const jLDSL = 013

// ---------------------------------------------------
// SSC

func sscSourceItemU(a uint64, x uint64, h uint64, i uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fSSC, jSSC, a, x, u)
}

// ---------------------------------------------------
// DSC

func dscSourceItemU(a uint64, x uint64, h uint64, i uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fDSC, jDSC, a, x, u)
}

// ---------------------------------------------------
// SSL

func sslSourceItemU(a uint64, x uint64, h uint64, i uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fSSL, jSSL, a, x, u)
}

// ---------------------------------------------------
// DSL

func dslSourceItemU(a uint64, x uint64, h uint64, i uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fDSL, jDSL, a, x, u)
}

// ---------------------------------------------------
// SSA

func ssaSourceItemU(a uint64, x uint64, h uint64, i uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fSSA, jSSA, a, x, u)
}

// ---------------------------------------------------
// DSA

func dsaSourceItemU(a uint64, x uint64, h uint64, i uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fSSA, jSSA, a, x, u)
}

// ---------------------------------------------------
// LSSC

func lsscSourceItemU(a uint64, x uint64, h uint64, i uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fLSSC, jLSSC, a, x, u)
}

// ---------------------------------------------------------------------------------------------------------------------

var sscCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 0, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 1, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 2, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 3, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 5, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 6, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 7, 0, 0, 0, 2, "data0"),
	sscSourceItemU(0, 0, 0, 0, 0),
	sscSourceItemU(1, 0, 0, 0, 36),
	sscSourceItemU(2, 0, 0, 0, 72),
	sscSourceItemU(3, 0, 0, 0, 18),
	sscSourceItemU(4, 0, 0, 0, 1),
	sscSourceItemU(5, 0, 0, 0, 35),
	sscSourceItemU(6, 0, 0, 0, 0206), // should just be 06, the 0200 gets stripped out
	iarSourceItem(0),

	segSourceItem(2),
	tasm.NewSourceItem("data0", "w", []string{"0112233445566"}),
}

func Test_SSC(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", sscCode)
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
	checkRegister(t, engine, common.A0, 0_112233_445566)
	checkRegister(t, engine, common.A1, 0_112233_445566)
	checkRegister(t, engine, common.A2, 0_112233_445566)
	checkRegister(t, engine, common.A3, 0_445566_112233)
	checkRegister(t, engine, common.A4, 0_045115_622673)
	checkRegister(t, engine, common.A5, 0_224467_113354)
	checkRegister(t, engine, common.A6, 0_661122_334455)
}

// TODO DSC

var sslCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 0, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 1, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 2, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 3, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "data1"),
	laSourceItemHIBRef(jW, 5, 0, 0, 0, 2, "data1"),
	laSourceItemHIBRef(jW, 6, 0, 0, 0, 2, "data2"),
	laSourceItemHIBRef(jW, 7, 0, 0, 0, 2, "data2"),
	sslSourceItemU(0, 0, 0, 0, 0),
	sslSourceItemU(1, 0, 0, 0, 36),
	sslSourceItemU(2, 0, 0, 0, 72),
	sslSourceItemU(3, 0, 0, 0, 18),
	sslSourceItemU(4, 0, 0, 0, 1),
	sslSourceItemU(5, 0, 0, 0, 35),
	sslSourceItemU(6, 0, 0, 0, 1),
	sslSourceItemU(7, 0, 0, 0, 35),
	iarSourceItem(0),

	segSourceItem(2),
	tasm.NewSourceItem("data0", "w", []string{"0112233445566"}),
	tasm.NewSourceItem("data1", "w", []string{"0555555555555"}),
	tasm.NewSourceItem("data2", "w", []string{"0111111666666"}),
}

func Test_SSL(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", sslCode)
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
	checkRegister(t, engine, common.A0, 0_112233_445566)
	checkRegister(t, engine, common.A1, 0)
	checkRegister(t, engine, common.A2, 0)
	checkRegister(t, engine, common.A3, 0_112233)
	checkRegister(t, engine, common.A4, 0_266666_666666)
	checkRegister(t, engine, common.A5, 01)
	checkRegister(t, engine, common.A6, 0_044444_733333)
	checkRegister(t, engine, common.A7, 0)
}

// TODO DSL
// TODO SSA
// TODO DSA
// TODO LSC
// TODO DLSC

var lsscCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 0, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 1, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 2, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 3, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 5, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 6, 0, 0, 0, 2, "data0"),
	laSourceItemHIBRef(jW, 7, 0, 0, 0, 2, "data0"),
	lsscSourceItemU(regA0, 0, 0, 0, 0),
	lsscSourceItemU(regA1, 0, 0, 0, 36),
	lsscSourceItemU(regA2, 0, 0, 0, 72),
	lsscSourceItemU(regA3, 0, 0, 0, 18),
	lsscSourceItemU(regA4, 0, 0, 0, 1),
	lsscSourceItemU(regA5, 0, 0, 0, 35),
	lsscSourceItemU(regA6, 0, 0, 0, 0206), // should just be 06, the 0200 gets stripped out
	iarSourceItem(0),

	segSourceItem(2),
	labelDataSourceItem("data0", []uint64{0_112233_445566}),
}

func Test_LSSC(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lsscCode)
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
	checkRegister(t, engine, common.A0, 0_112233_445566)
	checkRegister(t, engine, common.A1, 0_112233_445566)
	checkRegister(t, engine, common.A2, 0_112233_445566)
	checkRegister(t, engine, common.A3, 0_445566_112233)
	checkRegister(t, engine, common.A4, 0_224467_113354)
	checkRegister(t, engine, common.A5, 0_045115_622673)
	checkRegister(t, engine, common.A6, 0_223344_556611)
}

// TODO LDSC
// TODO LSSL
// TODO LDSL
