// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

var exBasicMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data", "hw", []string{"0123456", "0654321"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXI, jU, rX5, zero, "01"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXM, jU, rX5, zero, "04"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fEXBasic, jEXBasic, zero, rX5, "01", zero, "target"}),
	iarSourceItem("end", "0"),

	tasm.NewSourceItem("target", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jH1, rR3, zero, zero, zero, "data"}),
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
	checkStopped(t, engine)
	checkRegister(t, engine, X5, 0_000001_000005, "X5")
	checkRegister(t, engine, R3, 0_123456, "R3")
}

var exBasicModeIndirect = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"12"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXI, jU, rX5, zero, "01"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXM, jU, rX5, zero, "04"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fEXBasic, jEXBasic, zero, zero, zero, "1", "ind1"}),
	iarSourceItem("end", "0"),

	tasm.NewSourceItem("", ".SEG", []string{"15"}),
	tasm.NewSourceItem("data1", "fjaxhiu", []string{zero, zero, zero, zero, zero, "1", "data2"}),
	tasm.NewSourceItem("data2", "fjaxhiu", []string{zero, zero, zero, zero, zero, zero, "data3"}),
	tasm.NewSourceItem("data3", "hw", []string{"0123456", "0654321"}),

	tasm.NewSourceItem("", ".SEG", []string{"14"}),
	tasm.NewSourceItem("ind1", "fjaxhiu", []string{zero, zero, zero, zero, zero, "1", "ind2"}),
	tasm.NewSourceItem("ind2", "fjaxhiu", []string{zero, zero, zero, zero, zero, zero, "target"}),
	tasm.NewSourceItem("target", "fjaxhiu", []string{fLR, jH1, rR3, zero, zero, "1", "data1"}),
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
	checkRegister(t, engine, R3, 0_123456, "R3")
}

var exExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"02"}),
	tasm.NewSourceItem("data", "hw", []string{"0123456", "0654321"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXI, jU, rX5, zero, "01"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXM, jU, rX5, zero, "04"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fEXExtended, jEXExtended, aEXExtended, rX5, "01", zero, zero, "target"}),
	iarSourceItem("end", "0"),
	tasm.NewSourceItem("target", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jH1, rR3, zero, zero, zero, zero, "data"}),
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
	checkStopped(t, engine)
	checkRegister(t, engine, X5, 0_000001_000005, "X5")
	checkRegister(t, engine, R3, 0_123456, "R3")
}

var exExtendedModeCascade = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"02"}),
	tasm.NewSourceItem("data", "hw", []string{"0123456", "0654321"}),

	tasm.NewSourceItem("", ".SEG", []string{"00"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fEXExtended, jEXExtended, aEXExtended, zero, zero, zero, rB4, "target1"}),
	iarSourceItem("end", "0"),

	tasm.NewSourceItem("", ".SEG", []string{"04"}),
	tasm.NewSourceItem("target1", "fjaxhibd", []string{fEXExtended, jEXExtended, aEXExtended, zero, zero, zero, rB4, "target2"}),
	tasm.NewSourceItem("target2", "fjaxhibd", []string{fEXExtended, jEXExtended, aEXExtended, zero, zero, zero, rB4, "target3"}),
	tasm.NewSourceItem("target3", "fjaxhibd", []string{fEXExtended, jEXExtended, aEXExtended, zero, zero, zero, rB5, "target4"}),

	tasm.NewSourceItem("", ".SEG", []string{"05"}),
	tasm.NewSourceItem("target4", "fjaxhibd", []string{fLR, jH1, rR7, zero, zero, zero, rB2, "data"}),
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
	checkStopped(t, engine)
	checkRegister(t, engine, R7, 0_123456, "R7")
}

var exExtendedModeJump = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fEXExtended, jEXExtended, aEXExtended, zero, zero, zero, zero, "target"}),
	iarSourceItem("badend", "1"),
	tasm.NewSourceItem("target", "fjaxu", []string{fJ, jJExtended, aJExtended, zero, "goodend"}),
	iarSourceItem("goodend", "0"),
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
	checkStopped(t, engine)
}

var exExtendedModeTest = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fEXExtended, jEXExtended, aEXExtended, zero, zero, zero, zero, "target"}),
	iarSourceItem("badend", "1"),
	iarSourceItem("goodend", "0"),
	tasm.NewSourceItem("target", "fjaxu", []string{fTSKP, zero, aTSKP, zero, "target"}),
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
	checkStopped(t, engine)
}

// TODO EXR
// TODO NOP
// TODO DCB
// TODO RNGI
// TODO RNGB
