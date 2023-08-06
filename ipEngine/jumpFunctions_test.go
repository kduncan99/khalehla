// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

var jumpBasicMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fJ, jJBasic, aJBasic, zero, "target"}),
	iarSourceItem("", "1"),
	iarSourceItem("target", "0"),
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

var jumpExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fJ, jJExtended, aJExtended, zero, "target"}),
	iarSourceItem("", "1"),
	iarSourceItem("target", "0"),
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

var haltJumpExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fHLTJ, jHLTJ, aHLTJ, zero, "target"}),
	iarSourceItem("", "1"),
	iarSourceItem("target", "0"),
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
