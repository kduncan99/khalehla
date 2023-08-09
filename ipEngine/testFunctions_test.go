// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

// Unconditional -------------------------------------------------------------------------------------------------------

var tzBasicMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fNOPBasic, jNOPBasic, aNOP, zero, zero}),
	tasm.NewSourceItem("", "fjaxu", []string{fNOPBasic, jNOPBasic, aNOP, zero, zero}),
	tasm.NewSourceItem("", "fjaxu", []string{fNOPBasic, jNOPBasic, aNOP, zero, zero}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXM, jU, rX10, zero, "03"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLX, jU, rX11, zero, "0"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLMJ, jLMJ, rX11, rX10, zero, zero, "label"}),
	tasm.NewSourceItem("", "fjaxu", []string{fNOPBasic, jNOPBasic, aNOP, zero, zero}),
	tasm.NewSourceItem("", "fjaxu", []string{fNOPBasic, jNOPBasic, aNOP, zero, zero}),
	tasm.NewSourceItem("", "fjaxu", []string{fNOPBasic, jNOPBasic, aNOP, zero, zero}),
	tasm.NewSourceItem("label", "fjaxu", []string{fNOPBasic, jNOPBasic, aNOP, zero, zero}),
	iarSourceItem("", "1"),
	iarSourceItem("", "2"),
	iarSourceItem("target", "0"),
}

func Test_TZ_Basic(t *testing.T) {
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
