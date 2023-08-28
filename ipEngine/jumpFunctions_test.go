// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

// Unconditional -------------------------------------------------------------------------------------------------------

var lmjBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	nopItemHIU("", 0, 0, 0, 0),
	nopItemHIU("", 0, 0, 0, 0),
	nopItemHIU("", 0, 0, 0, 0),
	lxmSourceItemU("", jU, 10, 0, 03),
	lxSourceItemU("", jU, 11, 0, 0),
	lmjSourceItemHIURef("", 11, 10, 0, 0, "label"),
	nopItemHIU("", 0, 0, 0, 0),
	nopItemHIU("", 0, 0, 0, 0),
	nopItemHIU("", 0, 0, 0, 0),
	nopItemHIU("label", 0, 0, 0, 0),
	iarSourceItem("", 1),
	iarSourceItem("", 2),
	iarSourceItem("target", 0),
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
	checkRegister(t, engine, X11, 0_000000_01006, "X11")
}

var sljBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	sljSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("", 1),
	iarSourceItem("", 1),
	tasm.NewSourceItem("target", "w", []string{"0"}),
	iarSourceItem("", 0),
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
	nopItemHIBD("", 0, 0, 0, 0, 0),
	nopItemHIBD("", 0, 0, 0, 0, 0),
	nopItemHIBD("", 0, 0, 0, 0, 0),
	lxmSourceItemU("", jU, 10, 0, 03),
	lxSourceItemU("", jU, 11, 0, 0),
	lmjSourceItemHIBDRef("", 11, 10, 0, 0, 0, "label"),
	nopItemHIBD("", 0, 0, 0, 0, 0),
	nopItemHIBD("", 0, 0, 0, 0, 0),
	nopItemHIBD("", 0, 0, 0, 0, 0),
	nopItemHIBD("label", 0, 0, 0, 0, 0),
	iarSourceItem("", 1),
	iarSourceItem("", 2),
	iarSourceItem("target", 0),
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
	checkRegister(t, engine, X11, 0_000000_01006, "X11")
}

var jumpBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	jSourceItemBasic("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01003)
}

var jumpKeyBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	jkSourceItemHIURef("", 1, 0, 0, 0, "target"),
	iarSourceItem("", 0),
	iarSourceItem("target", 1),
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

var jumpExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	jSourceItemExtended("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01003)
}

var haltKeysAndJumpBasicMode = []*tasm.SourceItem{
	segSourceItem(0),
	hkjSourceItemHIURef("", 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01003)
}

var haltJumpExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	hltjSourceItemHIBDRef("", 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	checkStopped(t, engine)
	reason, detail := engine.GetStopReason()
	if reason != HaltJumpExecutedStop {
		t.Fatalf("Processor stopped for wrong reason: %d detail: %012o", reason, detail)
	}

	checkProgramAddress(t, engine, 01002)
}

// Conditional based on register ---------------------------------------------------------------------------------------

var jumpZeroExtendedPosZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIU("", jU, 5, 0, 0, 0, 0),
	jzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	laSourceItemU("", jXU, 5, 0, 0_777777),
	jzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	laSourceItemHIU("", jU, 5, 0, 0, 0, 01),
	jzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	dlSourceItemHIBDRef("", 0, 0, 0, 0, 3, "posZero"),
	djzSourceItemHIBDRef("", 0, 0, 0, 0, 0, "target1"),
	iarSourceItem("", 1),

	dlSourceItemHIBDRef("target1", 2, 0, 0, 0, 3, "negZero"),
	djzSourceItemHIBDRef("", 2, 0, 0, 0, 0, "target2"),
	iarSourceItem("", 2),

	dlSourceItemHIBDRef("target2", 4, 0, 0, 0, 3, "notZero1"),
	djzSourceItemHIBDRef("", 4, 0, 0, 0, 0, "bad3"),

	dlSourceItemHIBDRef("target3", 6, 0, 0, 0, 3, "notZero2"),
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
	laSourceItemU("", jU, 5, 0, 0),
	jnzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	laSourceItemU("", jXU, 5, 0, 0_777777),
	jnzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	laSourceItemU("", jU, 5, 0, 1),
	jnzSourceItemHIBDRef("", 5, 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
	iarSourceItem("target", 0),
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
	laSourceItemHIU("", jU, 10, 0, 0, 0, 0),
	jpSourceItemHIBDRef("", 10, 0, 0, 0, 0, "target1"),
	iarSourceItem("bad1", 1),
	jnSourceItemHIBDRef("target1", 10, 0, 0, 0, 0, "bad2"),

	nopItemHIBD("", 0, 0, 0, 0, 0),
	laSourceItemU("", jXU, 10, 0, 0_444444),
	jpSourceItemHIBDRef("", 10, 0, 0, 0, 0, "bad3"),
	jnSourceItemHIBDRef("", 10, 0, 0, 0, 0, "end"),

	iarSourceItem("bad4", 4),
	iarSourceItem("bad2", 2),
	iarSourceItem("bad3", 3),
	iarSourceItem("end", 0),
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
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

var jumpCarryExtended = []*tasm.SourceItem{
	segSourceItem(0),
	jcSourceItemHIBDRef("", 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01003)
}

var jumpDivideFault = []*tasm.SourceItem{
	segSourceItem(0),
	jdfSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

var jumpFloatingOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	jfoSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

var jumpFloatingUnderflow = []*tasm.SourceItem{
	segSourceItem(0),
	jfuSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

//	--------------------------------------------------------------------------------------------------------------------

var jumpNoCarryBasic = []*tasm.SourceItem{
	segSourceItem(0),
	jncSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoCarryExtended = []*tasm.SourceItem{
	segSourceItem(0),
	jncSourceItemHIBDRef("", 0, 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01003)
}

var jumpNoDivideFault = []*tasm.SourceItem{
	segSourceItem(0),
	jndfSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoFloatingOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	jnfoSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoFloatingUnderflow = []*tasm.SourceItem{
	segSourceItem(0),
	jnfuSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	jnoSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}

var jumpOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	joSourceItemHIURef("", 0, 0, 0, "target"),
	iarSourceItem("", 1),
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
	checkStopped(t, engine)
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
	checkStopped(t, engine)
	checkProgramAddress(t, engine, 01002)
}
