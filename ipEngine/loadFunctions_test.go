// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

//	TODO other loads that we are not yet testing

var laBasicMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data", "", []string{}),
	tasm.NewSourceItem("a1value", "sw", []string{"01", "02", "03", "04", "05", "06"}),
	tasm.NewSourceItem("a2value", "qw", []string{"0101", "0102", "0103", "0104"}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{"07777"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLA, jU, rA0, rX0, "0123"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jW, rA1, rX0, zero, zero, "a1value"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jQ2, rA2, rX0, zero, zero, "a2value"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLX, jU, rX4, rX0, "05"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jW, rA3, rX4, zero, zero, "data"}),
	IARSourceItem("", "0"),
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
	checkStopped(t, engine)
	checkRegister(t, engine, A0, 0123, "A0")
	checkRegister(t, engine, A1, 0_010203_040506, "A1")
	checkRegister(t, engine, A2, 0102, "A2")
	checkRegister(t, engine, A3, 07777, "A3")
}

var laExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data", "", []string{}),
	tasm.NewSourceItem("a1value", "sw", []string{"01", "02", "03", "04", "05", "06"}),
	tasm.NewSourceItem("a2value", "qw", []string{"0101", "0102", "0103", "0104"}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{"07777"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLA, jU, rA0, rX0, "0123"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jW, rA1, rX0, zero, zero, rB0, "a1value"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jQ2, rA2, rX0, zero, zero, rB0, "a2value"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLX, jU, rX4, rX0, "05"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jW, rA3, rX4, zero, zero, rB0, "data"}),
	IARSourceItem("", "0"),
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
	checkStopped(t, engine)
	checkRegister(t, engine, A0, 0123, "A0")
	checkRegister(t, engine, A1, 0_010203_040506, "A1")
	checkRegister(t, engine, A2, 0102, "A2")
	checkRegister(t, engine, A3, 07777, "A3")
}

var lrBasicMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("r7value", "qw", []string{"061", "062", "063", "064"}),
	tasm.NewSourceItem("r8value", "sw", []string{"01", "02", "03", "04", "05", "06"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jQ3, rR7, rX0, zero, zero, "r7value"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jXH2, rR8, rX0, zero, zero, "r8value"}),
	IARSourceItem("", "0"),
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
	checkStopped(t, engine)
	checkRegister(t, engine, R7, 063, "R7")
	checkRegister(t, engine, R8, 040506, "R8")
}

var lrExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("r5value", "tw", []string{"03210", "04000", "0123"}),
	tasm.NewSourceItem("r4value", "sw", []string{"01", "02", "03", "04", "05", "06"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jT2, rR5, rX0, zero, zero, rB0, "r5value"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jXH2, rR4, rX0, zero, zero, rB0, "r4value"}),
	IARSourceItem("", "0"),
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
	checkStopped(t, engine)
	checkRegister(t, engine, R4, 0_000000_040506, "R4")
	checkRegister(t, engine, R5, 0_777777_774000, "R5")
}

var lxBasicMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data", "w", []string{"0112233445566"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLX, jU, rX1, rX0, "0377777"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLX, jW, rX15, rX0, zero, zero, "data"}),
	IARSourceItem("", "0"),
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
	checkStopped(t, engine)
	checkRegister(t, engine, X1, 0_377777, "X1")
	checkRegister(t, engine, X15, 0_112233_445566, "X15")
	checkRegister(t, engine, A3, 0_112233_445566, "A3")
}

var lxExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLX, jU, rX1, rX0, "05"}),
	IARSourceItem("", "0"),
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
	checkStopped(t, engine)
	checkRegister(t, engine, X1, 05, "X1")
}

var dlBasicMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("posValue", "w", []string{"0100200300400"}),
	tasm.NewSourceItem("negValue", "w", []string{"0500600700777"}),
	tasm.NewSourceItem("", "w", []string{"05"}),
	tasm.NewSourceItem("indAddr1", "w", []string{"0200000+indAddr2"}),
	tasm.NewSourceItem("indAddr2", "w", []string{"posValue"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fDL, jDL, rA4, zero, zero, zero, "posValue"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fDL, jDL, rA0, zero, zero, "1", "indAddr1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fDLN, jDLN, rA2, zero, zero, zero, "posValue"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fDLN, jDLN, rA6, zero, zero, zero, "negValue"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fDLM, jDLM, rA10, zero, zero, zero, "posValue"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fDLM, jDLM, rA12, zero, zero, zero, "negValue"}),
	IARSourceItem("", "0"),
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
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStopped(t, engine)
	checkRegister(t, engine, A4, 0_100200_300400, "A4")
	checkRegister(t, engine, A5, 0_500600_700777, "A5")
	checkRegister(t, engine, A0, 0_100200_300400, "A0")
	checkRegister(t, engine, A1, 0_500600_700777, "A1")
	checkRegister(t, engine, A2, 0_677577_477377, "A2")
	checkRegister(t, engine, A3, 0_277177_077000, "A3")
	checkRegister(t, engine, A6, 0_277177_077000, "A6")
	checkRegister(t, engine, A7, 0_777777_777772, "A7")
	checkRegister(t, engine, A10, 0_100200_300400, "A10")
	checkRegister(t, engine, A11, 0_500600_700777, "A11")
	checkRegister(t, engine, A12, 0_277177_077000, "A12")
	checkRegister(t, engine, A13, 0_777777_777772, "A13")
}
