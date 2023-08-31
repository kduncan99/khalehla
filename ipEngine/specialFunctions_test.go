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

const fEXBasic = 072
const fEXExtended = 073
const fEXR = 073
const fNOPBasic = 074
const fNOPExtended = 073
const fDCB = 033
const fRNGI = 037
const fRNGB = 037

const jEXBasic = 010
const jEXExtended = 014
const jEXR = 014
const jNOPBasic = 006
const jNOPExtended = 014
const jDCB = 015
const jRNGI = 004
const jRNGB = 004

const aEXBasic = 000
const aEXExtended = 005
const aEXR = 006
const aNOPBasic = 000
const aNOPExtended = 000
const aRNGI = 005
const aRNGB = 006

// ---------------------------------------------------
// EX

func exSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fEXBasic, jEXBasic, aEXBasic, x, h, i, ref)
}

func exSourceItemHIBRef(x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fEXExtended, jEXExtended, aEXExtended, x, h, i, b, ref)
}

// ---------------------------------------------------
// EXR - extended mode only

func exrSourceItemHIBRef(x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fEXR, jEXR, aEXR, x, h, i, b, ref)
}

// ---------------------------------------------------
// NOP

func nopBasic() *tasm.SourceItem {
	return fjaxuSourceItem(fNOPBasic, jNOPBasic, aNOPBasic, 0, 0)
}

func nopExtended() *tasm.SourceItem {
	return fjaxuSourceItem(fNOPExtended, jNOPExtended, aNOPExtended, 0, 0)
}

func nopItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fNOPBasic, jNOPBasic, aNOPBasic, x, h, i, ref)
}

func nopItemHIBRef(x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fNOPExtended, jNOPExtended, aNOPExtended, x, h, i, b, ref)
}

// ---------------------------------------------------
// DCB - extended mode only

func dcbSourceItemU(a uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fDCB, jDCB, a, 0, u)
}

func dcbSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fDCB, jDCB, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// RNGI - extended mode only

func rngiSourceItemHIBRef(x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fRNGI, jRNGI, aRNGI, x, h, i, b, ref)
}

// ---------------------------------------------------
// RNGB - extended mode only

func rngbSourceItemHIBRef(x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fRNGB, jRNGB, aRNGB, x, h, i, b, ref)
}

// ---------------------------------------------------------------------------------------------------------------------

var exBasicMode = []*tasm.SourceItem{
	segSourceItem(077),
	labelDataSourceItem("data", []uint64{0123456, 0654321}),

	segSourceItem(0),
	lxiSourceItemU(jU, regX5, 0, 01),
	lxmSourceItemU(jU, regX5, 0, 04),
	exSourceItemHIRef(regX5, 1, 0, "target"),
	iarSourceItem(0),

	labelDataSourceItem("target", []uint64{0}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),
	lrSourceItemHIRef(jH1, regR3, 0, 0, 0, "data"),
}

func Test_EX_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", exBasicMode)
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
	checkRegister(t, engine, pkg.X5, 0_000001_000005)
	checkRegister(t, engine, pkg.R3, 0_123456)
}

var exBasicModeIndirect = []*tasm.SourceItem{
	segSourceItem(12),
	lxiSourceItemU(jU, regX5, 0, 01),
	lxmSourceItemU(jU, regX5, 0, 04),
	exSourceItemHIRef(0, 0, 1, "ind1"),
	iarSourceItem(0),

	segSourceItem(15),
	tasm.NewSourceItem("data1", "fjaxhiu", []string{zero, zero, zero, zero, zero, "1", "data2"}),
	tasm.NewSourceItem("data2", "fjaxhiu", []string{zero, zero, zero, zero, zero, zero, "data3"}),
	tasm.NewSourceItem("data3", "hw", []string{"0123456", "0654321"}),

	segSourceItem(14),
	tasm.NewSourceItem("ind1", "fjaxhiu", []string{zero, zero, zero, zero, zero, "1", "ind2"}),
	tasm.NewSourceItem("ind2", "fjaxhiu", []string{zero, zero, zero, zero, zero, zero, "target"}),
	labelSourceItem("target"),
	lrSourceItemHIRef(jH1, regR3, 0, 0, 1, "data1"),
}

func Test_EX_BasicModeIndirect(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", exBasicModeIndirect)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetProcessorPrivilege(2) // for indirect to work
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkRegister(t, engine, pkg.R3, 0_123456)
}

var exExtendedMode = []*tasm.SourceItem{
	segSourceItem(2),
	tasm.NewSourceItem("data", "hw", []string{"0123456", "0654321"}),

	segSourceItem(0),
	lxiSourceItemU(jU, regX5, 0, 01),
	lxmSourceItemU(jU, regX5, 0, 04),
	exSourceItemHIBRef(regX5, 1, 0, pkg.B0, "target"),
	iarSourceItem(0),

	labelDataSourceItem("target", []uint64{0}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),
	lrSourceItemHIBRef(jH1, regR3, 0, 0, 0, 0, "data"),
}

func Test_EX_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", exExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.X5, 0_000001_000005)
	checkRegister(t, engine, pkg.R3, 0_123456)
}

var exExtendedModeCascade = []*tasm.SourceItem{
	segSourceItem(02),
	tasm.NewSourceItem("data", "hw", []string{"0123456", "0654321"}),

	segSourceItem(0),
	exSourceItemHIBRef(0, 0, 0, pkg.B4, "target1"),
	iarSourceItem(0),

	segSourceItem(04),
	labelSourceItem("target1"),
	exSourceItemHIBRef(0, 0, 0, pkg.B4, "target2"),
	labelSourceItem("target2"),
	exSourceItemHIBRef(0, 0, 0, pkg.B4, "target3"),
	labelSourceItem("target3"),
	exSourceItemHIBRef(0, 0, 0, pkg.B5, "target4"),

	segSourceItem(05),
	labelSourceItem("target4"),
	lrSourceItemHIBRef(jH1, regR7, 0, 0, 0, pkg.B2, "data"),
}

func Test_EX_ExtendedCascade(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", exExtendedModeCascade)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.R7, 0_123456)
}

var exExtendedModeJump = []*tasm.SourceItem{
	segSourceItem(0),
	exSourceItemHIBRef(0, 0, 0, 0, "target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	jSourceItemRefExtended("goodend"),

	labelSourceItem("goodend"),
	iarSourceItem(0),
}

func Test_EX_ExtendedJump(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", exExtendedModeJump)
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

var exExtendedModeTest = []*tasm.SourceItem{
	segSourceItem(0),
	exSourceItemHIBRef(0, 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem(0),

	labelSourceItem("target"),
	tskpSourceItemHIBRef(jW, 0, 0, 0, 0, "target"),
}

func Test_EX_ExtendedTest(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", exExtendedModeTest)
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

var exrExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	lxmSourceItemU(jU, regX7, 0, 02),
	lxiSourceItemU(jU, regX8, 0, 01),
	lxmSourceItemU(jU, regX8, 0, 0),
	lrSourceItemU(jU, regR1, 0, 010),
	exrSourceItemHIBRef(regX7, 0, 0, pkg.B2, "target"),
	iarSourceItem(0),

	segSourceItem(02),
	labelSourceItem("target"),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),
	sasSourceItemHIBRef(jH1, regX8, 1, 0, pkg.B3, "data"),

	segSourceItem(03),
	labelSourceItem("data"),
	dataSourceItem([]uint64{0, 0_777777}),
	dataSourceItem([]uint64{0, 0_777777}),
	dataSourceItem([]uint64{0, 0_777777}),
	dataSourceItem([]uint64{0, 0_777777}),
	dataSourceItem([]uint64{0, 0_777777}),
	dataSourceItem([]uint64{0, 0_777777}),
	dataSourceItem([]uint64{0, 0_777777}),
	dataSourceItem([]uint64{0, 0_777777}),
	dataSourceItem([]uint64{0, 0_777777}),
	dataSourceItem([]uint64{0, 0_777777}),
}

func Test_EXR_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", exrExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.X8, 0_000001_000010)

	for ox := uint64(0); ox < 8; ox++ {
		checkMemory(t, engine, engine.baseRegisters[3].GetBankDescriptor().GetBaseAddress(), ox, 0_040040_777777)
	}

	for ox := uint64(8); ox < 10; ox++ {
		checkMemory(t, engine, engine.baseRegisters[3].GetBankDescriptor().GetBaseAddress(), ox, 0_000000_777777)
	}
}

var exrExtendedModeInvalidInstruction = []*tasm.SourceItem{
	segSourceItem(0),
	exrSourceItemHIBRef(regX7, 0, 0, 0, "target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	laSourceItemU(jU, regA3, 0, 0_177777),
}

func Test_EXR_ExtendedInvalidInstruction(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", exrExtendedModeInvalidInstruction)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkInterrupt(t, engine, pkg.InvalidInstructionInterruptClass)
}

var exrExtendedModeTZ = []*tasm.SourceItem{
	segSourceItem(0),
	lxiSourceItemU(jU, regX8, 0, 01),
	lxmSourceItemU(jU, regX8, 0, 0),
	lrSourceItemU(jU, regR1, 0, 020),
	exrSourceItemHIBRef(0, 0, 0, 0, "target"),
	iarSourceItem(1),
	iarSourceItem(0),

	labelSourceItem("target"),
	tzSourceItemHIBDRef(jH2, regX8, 1, 0, 0, "data"),

	labelSourceItem("data"),
	dataSourceItem([]uint64{04}),
	dataSourceItem([]uint64{03}),
	dataSourceItem([]uint64{02}),
	dataSourceItem([]uint64{01}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{022}),
	dataSourceItem([]uint64{023}),
	dataSourceItem([]uint64{024}),
	dataSourceItem([]uint64{025}),
	dataSourceItem([]uint64{0}),
}

func Test_EXR_ExtendedTZ(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", exrExtendedModeTZ)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.R1, 013)
}

// TODO NOP

var dcbExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	dcbSourceItemHIBRef(regA3, 0, 0, 0, 0, "data"),
	dcbSourceItemHIBRef(regA4, 0, 0, 0, 0, "data+1"),
	dcbSourceItemHIBRef(regA5, 0, 0, 0, 0, "data+2"),
	iarSourceItem(0),

	labelSourceItem("data"),
	dataSourceItem([]uint64{0_030405_030405}),
	dataSourceItem([]uint64{0_777777_777777}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{01}),
}

func Test_DCB_ExtendedTest(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", dcbExtendedMode)
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
	checkRegister(t, engine, pkg.A3, 46)
	checkRegister(t, engine, pkg.A4, 36)
	checkRegister(t, engine, pkg.A5, 1)
}

var rngbExtendedMode = []*tasm.SourceItem{
	segSourceItem(02),
	labelSourceItem("data"),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),

	segSourceItem(0),
	rngbSourceItemHIBRef(0, 0, 0, pkg.B2, "data"),
	iarSourceItem(0),
}

func Test_RNGB_ExtendedTest(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", rngbExtendedMode)
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

	//	not much we can do except verify bits 0, 9, 18, and 27 of storage are all zero.
	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	storage := engine.baseRegisters[2].GetStorage()
	for mx := 0; mx < 4; mx++ {
		value := storage[mx]
		fmt.Printf("%04o: %012o\n", mx, value)
		if value&0_400400_400400 != 0 {
			t.Fatalf("Expected MSB of each quarter-word to be zero")
		}
	}
}

var rngiExtendedMode = []*tasm.SourceItem{
	segSourceItem(02),
	labelSourceItem("data"),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),
	dataSourceItem([]uint64{0}),

	segSourceItem(0),
	rngiSourceItemHIBRef(0, 0, 0, pkg.B2, "data"),
	iarSourceItem(0),
}

func Test_RNGI_ExtendedTest(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", rngiExtendedMode)
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

	//	not much we can do except verify bits 0-3 of storage are all zero.
	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	storage := engine.baseRegisters[2].GetStorage()
	for mx := 0; mx < 4; mx++ {
		value := storage[mx]
		fmt.Printf("%04o: %012o\n", mx, value)
		if value&0_740000_000000 != 0 {
			t.Fatalf("Expected bits 0-3 to be zero")
		}
	}
}
