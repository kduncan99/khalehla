// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

var tepCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data1"),
	tepSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data2"),
	jSourceItemExtended("", 0, 0, 0, "tag"),
	iarSourceItem("badend1", 1),
	tepSourceItemHIBDRef("tag", jQ4, 4, 0, 0, 0, 2, "data2"),
	iarSourceItem("badend2", 2),
	iarSourceItem("goodend", 0),

	segSourceItem(2),
	tasm.NewSourceItem("data1", "w", []string{"0123456543210"}),
	tasm.NewSourceItem("data2", "w", []string{"0777777"}),
}

func Test_TEP(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", tepCode)
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
}

var topCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data1"),
	topSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data2"),
	jSourceItemExtended("", 0, 0, 0, "tag"),
	iarSourceItem("badend1", 1),
	topSourceItemHIBDRef("tag", jQ4, 4, 0, 0, 0, 2, "data2"),
	iarSourceItem("badend2", 2),
	iarSourceItem("goodend", 0),

	segSourceItem(2),
	tasm.NewSourceItem("data1", "w", []string{"0123456543211"}),
	tasm.NewSourceItem("data2", "w", []string{"0777777"}),
}

func Test_TOP(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", topCode)
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
}

var tlemCode = []*tasm.SourceItem{
	segSourceItem(0),
	lxSourceItemHIBDRef("", jW, 5, 0, 0, 0, 2, "x5content"),
	tlemSourceItemHIBDRef("", jW, 5, 0, 0, 0, 2, "arm"),
	jSourceItemExtended("", 0, 0, 0, "tag1"),
	iarSourceItem("badend1", 1),

	labelSourceItem("tag1"),
	lxSourceItemHIBD("", jW, 6, 0, 0, 0, 0, grsX5),
	tlemSourceItemHIBDRef("", jS5, 6, 0, 0, 0, 2, "arm"),
	iarSourceItem("badend2", 2),
	iarSourceItem("goodend", 0),

	segSourceItem(2),
	tasm.NewSourceItem("arm", "w", []string{"0135471234"}),
	tasm.NewSourceItem("x5content", "w", []string{"02061234"}),
}

func Test_TLEM(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", tlemCode)
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
}

var tnopCode = []*tasm.SourceItem{
	segSourceItem(0),
	tnopSourceItemHIBDRef("", jW, 0, 0, 0, 2, "data1"),
	iarSourceItem("goodend", 0),
	iarSourceItem("badend1", 1),

	segSourceItem(2),
	tasm.NewSourceItem("data1", "w", []string{"0123456543210"}),
}

func Test_TNOP(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", tnopCode)
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
}

//	TODO TGZ
//	TODO TPZ
//	TODO TP
//	TODO TMZ
//	TODO TMZG
//	TODO TZ
//	TODO TNLZ
//	TODO TLZ
//	TODO TNZ
//	TODO TPZL
//	TODO TNMZ
//	TODO TN
//	TODO TNPZ
//	TODO TNGZ

var tskpCode = []*tasm.SourceItem{
	segSourceItem(0),
	tskpSourceItemHIBDRef("", jW, 0, 0, 0, 2, "data1"),
	iarSourceItem("badend1", 1),
	iarSourceItem("goodend", 0),

	segSourceItem(2),
	tasm.NewSourceItem("data1", "w", []string{"0123456543210"}),
}

func Test_TSKP(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", tskpCode)
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
}

//	TODO TE
//	TODO DTE
//	TODO TNE
//	TODO TLE
//	TODO TG
//	TODO TGM
//	TODO DTGM
//	TODO TW
//	TODO TNW
//	TODO MTE
//	TODO MTNE
//	TODO MTLE
//	TODO MTG
//	TODO MTW
//	TODO MTNW
//	TODO MATL
//	TODO MATG
//	TODO TS
//	TODO TSS
//	TODO TCS
//	TODO CR
//	TODO UNLK
