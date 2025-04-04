// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"testing"

	"khalehla/common"
	"khalehla/tasm"
)

const fSLJ = 072
const fLMJ = 074
const fJ = 074
const fJK = 074
const fHKJ = 074
const fHLTJ = 074

const jSLJ = 001
const jLMJ = 013
const jJBasic = 004
const jJExtended = 015
const jJK = 004
const jHKJ = 005
const jHLTJ = 015

const aSLJ = 000
const aJBasic = 000
const aJExtended = 004
const aHLTJ = 005

// ---------------------------------------------------
// SLJ - basic mode only

func sljSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fSLJ, jSLJ, aSLJ, 0, ref)
}

func sljSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fSLJ, jSLJ, aSLJ, x, h, i, ref)
}

// ---------------------------------------------------
// LMJ

func lmjSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fLMJ, jLMJ, a, 0, ref)
}

func lmjSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fLMJ, jLMJ, a, x, h, i, ref)
}

func lmjSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fLMJ, jLMJ, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// J

func jSourceItemRefBasic(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJ, jJBasic, aJBasic, 0, ref)
}

func jSourceItemRefExtended(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJ, jJExtended, aJExtended, 0, ref)
}

func jSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJ, jJBasic, aJBasic, x, h, i, ref)
}

func jSourceItemHIBRef(x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fJ, jJExtended, aJExtended, x, h, i, b, ref)
}

// ---------------------------------------------------
// JK - basic mode only

func jkSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJK, jJK, a, 0, ref)
}

func jkSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJK, jJK, a, x, h, i, ref)
}

// ---------------------------------------------------
// HKJ - basic mode only

func hkjSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fHKJ, jHKJ, a, 0, ref)
}

func hkjSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fHKJ, jHKJ, a, x, h, i, ref)
}

// ---------------------------------------------------
// HLTJ

func hltjSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJ, jHLTJ, aHLTJ, 0, ref)
}

func hltjSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJ, jHLTJ, aHLTJ, x, h, i, ref)
}

func hltjSourceItemHIBRef(x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fHLTJ, jHLTJ, aHLTJ, x, h, i, b, ref)
}

// ---------------------------------------------------------------------------------------------------------------------

var lmjBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	nopBasic(),
	nopBasic(),
	nopBasic(),
	lxmSourceItemU(jU, regX10, 0, 03),
	lxSourceItemU(jU, regX11, 0, 0),
	lmjSourceItemHIRef(regX11, regX10, 0, 0, "label"),
	nopBasic(),
	nopBasic(),
	nopBasic(),
	labelSourceItem("label"),
	nopBasic(),
	iarSourceItem(1),
	iarSourceItem(2),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_LMJ_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lmjBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, common.X11, 0_000000_01006)
}

var sljBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	sljSourceItemRef("target"),
	iarSourceItem(1),
	iarSourceItem(1),
	iarSourceItem(1),

	tasm.NewSourceItem("target", "w", []string{"0"}),
	iarSourceItem(0),
}

func Test_SLJ_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", sljBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	codeBankAddr := e.GetBanks()[0600004].GetBankDescriptor().GetBaseAddress()
	checkMemory(t, engine, codeBankAddr, 04, 01001)
}

var lmjExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	nopExtended(),
	nopExtended(),
	nopExtended(),
	lxmSourceItemU(jU, regX10, 0, 03),
	lxSourceItemU(jU, regX11, 0, 0),
	lmjSourceItemHIBRef(regX11, regX10, 0, 0, 0, "label"),
	nopExtended(),
	nopExtended(),
	nopExtended(),
	labelSourceItem("label"),
	nopExtended(),
	iarSourceItem(1),
	iarSourceItem(2),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_LMJ_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lmjExtendedMode)
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
	checkRegister(t, engine, common.X11, 0_000000_001006)
}

var jumpBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	jSourceItemRefBasic("target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_J_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

var jumpKeyBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	jkSourceItemRef(regA1, "target"),
	iarSourceItem(0),
	labelSourceItem("target"),
	iarSourceItem(1),
}

func Test_JK_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpKeyBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01002)
}

var jumpExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	jSourceItemRefExtended("target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_J_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpExtendedMode)
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
	checkProgramAddress(t, engine, 01003)
}

var haltKeysAndJumpBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	hkjSourceItemRef(0, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_HKJ_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", haltKeysAndJumpBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

var haltJumpExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	hltjSourceItemRef("target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_HLTJ_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", haltJumpExtendedMode)
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
	checkStoppedReason(t, engine, HaltJumpExecutedStop, 0)
	checkProgramAddress(t, engine, 01002)
}
