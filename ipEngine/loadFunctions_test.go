// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
	"khalehla/tasm"
	"testing"
)

const fLA = 010
const fLMA = 012
const fLNA = 011
const fLNMA = 013
const fDL = 071
const fDLM = 071
const fDLN = 071
const fLR = 023
const fLX = 027
const fLXI = 046
const fLXM = 026

const jDL = 013
const jDLM = 015
const jDLN = 014

// ---------------------------------------------------
// LA

func laSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fLA, j, a, x, u)
}

func laSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fLA, j, a, x, h, i, ref)
}

func laSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fLA, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// LMA

func lmaSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fLMA, j, a, x, u)
}

func lmaSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fLMA, j, a, x, h, i, ref)
}

func lmaSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fLMA, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// LNA

func lnaSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fLNA, j, a, x, u)
}

func lnaSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fLNA, j, a, x, h, i, ref)
}

func lnaSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fLNA, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// LNMA

func lnmaSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fLNMA, j, a, x, u)
}

func lnmaSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fLNMA, j, a, x, h, i, ref)
}

func lnmaSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fLNMA, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// DL

func dlSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fDL, jDL, a, x, h, i, ref)
}

func dlSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fDL, jDL, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// DLM

func dlmSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fDLM, jDLM, a, x, h, i, ref)
}

func dlmSourceItemHIBDRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fDLM, jDLM, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// DLN

func dlnSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fDLN, jDLN, a, x, h, i, ref)
}

func dlnSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fDLN, jDLN, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// LR

func lrSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fLR, j, a, x, u)
}

func lrSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fLR, j, a, x, h, i, ref)
}

func lrSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fLR, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// LX

func lxSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fLX, j, a, x, u)
}

func lxSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fLX, j, a, x, h, i, ref)
}

func lxSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fLX, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// LXI

func lxiSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fLXI, j, a, x, u)
}

func lxiSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fLXI, j, a, x, h, i, ref)
}

func lxiSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fLXI, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// LXM

func lxmSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fLXM, j, a, x, u)
}

func lxmSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fLXM, j, a, x, h, i, ref)
}

func lxmSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fLXM, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------------------------------------------------------------------------

//	TODO LXLM
//	TODO LXSI
//	TODO LRS
//	TODO LAQW
//	TODO LSBO
//	TODO LSBL

var laBasicMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("data", "", []string{}),
	tasm.NewSourceItem("a1value", "sw", []string{"01", "02", "03", "04", "05", "06"}),
	tasm.NewSourceItem("a2value", "qw", []string{"0101", "0102", "0103", "0104"}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{"07777"}),

	segSourceItem(0),
	laSourceItemU(jU, 0, 0, 0123),
	laSourceItemHIRef(jW, 1, 0, 0, 0, "a1Value"),
	laSourceItemHIRef(jQ2, 2, 0, 0, 0, "a2value"),
	lxSourceItemU(jU, 4, 0, 5),
	laSourceItemHIRef(jW, 3, 4, 0, 0, "data"),
	iarSourceItem(0),
}

func Test_LA_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", laBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.A0, 0123)
	checkRegister(t, engine, pkg.A1, 0_010203_040506)
	checkRegister(t, engine, pkg.A2, 0102)
	checkRegister(t, engine, pkg.A3, 07777)
}

var laExtendedMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("data", "", []string{}),
	tasm.NewSourceItem("a1value", "sw", []string{"01", "02", "03", "04", "05", "06"}),
	tasm.NewSourceItem("a2value", "qw", []string{"0101", "0102", "0103", "0104"}),
	sourceItem("", "w", []int{0}),
	sourceItem("", "w", []int{0}),
	sourceItem("", "w", []int{0}),
	sourceItem("", "w", []int{07777}),

	segSourceItem(0),
	laSourceItemU(jU, 0, 0, 0123),
	laSourceItemHIBRef(jW, 1, 0, 0, 0, 0, "a1Value"),
	laSourceItemHIBRef(jQ2, 2, 0, 0, 0, 0, "a2Value"),
	lxmSourceItemU(jU, 4, 0, 05),
	laSourceItemHIBRef(jW, 3, 4, 0, 0, 0, "data"),
	iarSourceItem(0),
}

func Test_LA_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", laExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.A0, 0123)
	checkRegister(t, engine, pkg.A1, 0_010203_040506)
	checkRegister(t, engine, pkg.A2, 0102)
	checkRegister(t, engine, pkg.A3, 07777)
}

var lmaExtendedMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("posValue", "w", []string{"0300000123456"}),
	tasm.NewSourceItem("negValue", "w", []string{"0400000000001"}),
	tasm.NewSourceItem("partValue", "w", []string{"0555577664444"}),

	segSourceItem(0),
	lmaSourceItemU(jU, 0, 0, 0_377777),
	lmaSourceItemU(jU, 1, 0, 0_477777),
	lmaSourceItemU(jXU, 2, 0, 0_377777),
	lmaSourceItemU(jXU, 3, 0, 0_477777),
	lmaSourceItemHIBRef(jW, 4, 0, 0, 0, 0, "posValue"),
	lmaSourceItemHIBRef(jW, 5, 0, 0, 0, 0, "negValue"),
	lmaSourceItemHIBRef(jT2, 6, 0, 0, 0, 0, "partValue"),
	lmaSourceItemHIBRef(jS5, 7, 0, 0, 0, 0, "partValue"),
	iarSourceItem(0),
}

func Test_LMA_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lmaExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.A0, 0_377777)
	checkRegister(t, engine, pkg.A1, 0_477777)
	checkRegister(t, engine, pkg.A2, 0_377777)
	checkRegister(t, engine, pkg.A3, 0_300000)
	checkRegister(t, engine, pkg.A4, 0_300000_123456)
	checkRegister(t, engine, pkg.A5, 0_377777_777776)
	checkRegister(t, engine, pkg.A6, 011)
	checkRegister(t, engine, pkg.A7, 044)
}

var lnaExtendedMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("posValue", "w", []string{"0300000123456"}),
	tasm.NewSourceItem("negValue", "w", []string{"0400000000001"}),
	tasm.NewSourceItem("partValue", "w", []string{"0555577664444"}),

	segSourceItem(0),
	lnaSourceItemU(jU, 0, 0, 0_377777),
	lnaSourceItemU(jU, 1, 0, 0_477777),
	lnaSourceItemU(jXU, 2, 0, 0_377777),
	lnaSourceItemU(jXU, 3, 0, 0_477777),
	lnaSourceItemHIBRef(jW, 4, 0, 0, 0, 0, "posValue"),
	lnaSourceItemHIBRef(jW, 5, 0, 0, 0, 0, "negValue"),
	lnaSourceItemHIBRef(jT2, 6, 0, 0, 0, 0, "partValue"),
	lnaSourceItemHIBRef(jS5, 7, 0, 0, 0, 0, "partValue"),
	iarSourceItem(0),
}

func Test_LNA_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lnaExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.A0, 0_777777_400000)
	checkRegister(t, engine, pkg.A1, 0_777777_300000)
	checkRegister(t, engine, pkg.A2, 0_777777_400000)
	checkRegister(t, engine, pkg.A3, 0_300000)
	checkRegister(t, engine, pkg.A4, 0_477777_654321)
	checkRegister(t, engine, pkg.A5, 0_377777_777776)
	checkRegister(t, engine, pkg.A6, 011)
	checkRegister(t, engine, pkg.A7, 0_777777_777733)
}

var lnmaExtendedMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("posValue", "w", []string{"0300000123456"}),
	tasm.NewSourceItem("negValue", "w", []string{"0400000000001"}),
	tasm.NewSourceItem("partValue", "w", []string{"0555577664444"}),

	segSourceItem(0),
	lnmaSourceItemU(jU, 0, 0, 0_377777),
	lnmaSourceItemU(jU, 1, 0, 0_477777),
	lnmaSourceItemU(jXU, 2, 0, 0_377777),
	lnmaSourceItemU(jXU, 3, 0, 0_477777),
	lnmaSourceItemHIBRef(jW, 4, 0, 0, 0, 0, "posValue"),
	lnmaSourceItemHIBRef(jW, 5, 0, 0, 0, 0, "negValue"),
	lnmaSourceItemHIBRef(jT2, 6, 0, 0, 0, 0, "partValue"),
	lnmaSourceItemHIBRef(jS5, 7, 0, 0, 0, 0, "partValue"),
	iarSourceItem(0),
}

func Test_LNMA_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lnmaExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.A0, 0_777777_400000)
	checkRegister(t, engine, pkg.A1, 0_777777_300000)
	checkRegister(t, engine, pkg.A2, 0_777777_400000)
	checkRegister(t, engine, pkg.A3, 0_777777_477777)
	checkRegister(t, engine, pkg.A4, 0_477777_654321)
	checkRegister(t, engine, pkg.A5, 0_400000_000001)
	checkRegister(t, engine, pkg.A6, 0_777777_777766)
	checkRegister(t, engine, pkg.A7, 0_777777_777733)
}

var lrBasicMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("r7value", "qw", []string{"061", "062", "063", "064"}),
	tasm.NewSourceItem("r8value", "sw", []string{"01", "02", "03", "04", "05", "06"}),

	segSourceItem(0),
	lrSourceItemHIRef(jQ3, 7, 0, 0, 0, "r7value"),
	lrSourceItemHIRef(jXH2, 8, 0, 0, 0, "r8value"),
	iarSourceItem(0),
}

func Test_LR_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lrBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.R7, 063)
	checkRegister(t, engine, pkg.R8, 040506)
}

var lrExtendedMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("r5value", "tw", []string{"03210", "04000", "0123"}),
	tasm.NewSourceItem("r4value", "sw", []string{"01", "02", "03", "04", "05", "06"}),

	segSourceItem(0),
	lrSourceItemHIBRef(jT2, 5, 0, 0, 0, 0, "r5value"),
	lrSourceItemHIBRef(jXH2, 4, 0, 0, 0, 0, "r4value"),
	iarSourceItem(0),
}

func Test_LR_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lrExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.R4, 0_000000_040506)
	checkRegister(t, engine, pkg.R5, 0_777777_774000)
}

var lxBasicMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("data", "w", []string{"0112233445566"}),

	segSourceItem(0),
	lxSourceItemU(jU, 1, 0, 0_377777),
	lxSourceItemHIRef(jW, 15, 0, 0, 0, "data"),
	iarSourceItem(0),
}

func Test_LX_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lxBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.X1, 0_377777)
	checkRegister(t, engine, pkg.X15, 0_112233_445566)
	checkRegister(t, engine, pkg.A3, 0_112233_445566)
}

var lxExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	lxSourceItemU(jU, 1, 0, 05),
	iarSourceItem(0),
}

func Test_LX_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lxExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.X1, 05)
}

var dlBasicMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("posValue", "w", []string{"0100200300400"}),
	tasm.NewSourceItem("negValue", "w", []string{"0500600700777"}),
	tasm.NewSourceItem("", "w", []string{"05"}),
	tasm.NewSourceItem("indAddr1", "w", []string{"0200000+indAddr2"}),
	tasm.NewSourceItem("indAddr2", "w", []string{"posValue"}),

	segSourceItem(0),
	dlSourceItemHIRef(4, 0, 0, 0, "posValue"),
	dlSourceItemHIRef(0, 0, 0, 1, "indAddr1"),
	dlnSourceItemHIRef(2, 0, 0, 0, "posValue"),
	dlnSourceItemHIRef(6, 0, 0, 0, "negValue"),
	dlmSourceItemHIRef(10, 0, 0, 0, "posValue"),
	dlmSourceItemHIRef(12, 0, 0, 0, "negValue"),
	iarSourceItem(0),
}

func Test_DL_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", dlBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetProcessorPrivilege(2)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	//	don't check stopped - we execute at PP=2 and expect to get Invalid Interrupt on IAR instruction
	checkRegister(t, engine, pkg.A4, 0_100200_300400)
	checkRegister(t, engine, pkg.A5, 0_500600_700777)
	checkRegister(t, engine, pkg.A0, 0_100200_300400)
	checkRegister(t, engine, pkg.A1, 0_500600_700777)
	checkRegister(t, engine, pkg.A2, 0_677577_477377)
	checkRegister(t, engine, pkg.A3, 0_277177_077000)
	checkRegister(t, engine, pkg.A6, 0_277177_077000)
	checkRegister(t, engine, pkg.A7, 0_777777_777772)
	checkRegister(t, engine, pkg.A10, 0_100200_300400)
	checkRegister(t, engine, pkg.A11, 0_500600_700777)
	checkRegister(t, engine, pkg.A12, 0_277177_077000)
	checkRegister(t, engine, pkg.A13, 0_777777_777772)
}
