// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

//	TODO LXM, LXLM
//	TODO LXI, LXSI

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
	laSourceItemU("", jU, 0, 0, 0123),
	laSourceItemHIURef("", jW, 1, 0, 0, 0, "a1Value"),
	laSourceItemHIURef("", jQ2, 2, 0, 0, 0, "a2value"),
	lxSourceItemU("", jU, 4, 0, 5),
	laSourceItemHIURef("", jW, 3, 4, 0, 0, "data"),
	iarSourceItem("", 0),
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
	segSourceItem(077),
	tasm.NewSourceItem("data", "", []string{}),
	tasm.NewSourceItem("a1value", "sw", []string{"01", "02", "03", "04", "05", "06"}),
	tasm.NewSourceItem("a2value", "qw", []string{"0101", "0102", "0103", "0104"}),
	sourceItem("", "w", []int{0}),
	sourceItem("", "w", []int{0}),
	sourceItem("", "w", []int{0}),
	sourceItem("", "w", []int{0}),

	segSourceItem(0),
	laSourceItemU("", jU, 0, 0, 0123),
	laSourceItemHIBDRef("", jW, 1, 0, 0, 0, 0, "a1Value"),
	laSourceItemHIBDRef("", jQ2, 2, 0, 0, 0, 0, "a2Value"),
	lxSourceItemU("", jU, 4, 0, 05),
	laSourceItemHIBDRef("", jW, 3, 4, 0, 0, 0, "data"),
	iarSourceItem("", 0),
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

var lmaExtendedMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("posValue", "w", []string{"0300000123456"}),
	tasm.NewSourceItem("negValue", "w", []string{"0400000000001"}),
	tasm.NewSourceItem("partValue", "w", []string{"0555577664444"}),

	segSourceItem(0),
	lmaSourceItemU("", jU, 0, 0, 0_377777),
	lmaSourceItemU("", jU, 1, 0, 0_477777),
	lmaSourceItemU("", jXU, 2, 0, 0_377777),
	lmaSourceItemU("", jXU, 3, 0, 0_477777),
	lmaSourceItemHIBDRef("", jW, 4, 0, 0, 0, 0, "posValue"),
	lmaSourceItemHIBDRef("", jW, 5, 0, 0, 0, 0, "negValue"),
	lmaSourceItemHIBDRef("", jT2, 6, 0, 0, 0, 0, "partValue"),
	lmaSourceItemHIBDRef("", jS5, 7, 0, 0, 0, 0, "partValue"),
	iarSourceItem("", 0),
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
	checkStopped(t, engine)
	checkRegister(t, engine, A0, 0_377777, "A0")
	checkRegister(t, engine, A1, 0_477777, "A1")
	checkRegister(t, engine, A2, 0_377777, "A2")
	checkRegister(t, engine, A3, 0_300000, "A3")
	checkRegister(t, engine, A4, 0_300000_123456, "A4")
	checkRegister(t, engine, A5, 0_377777_777776, "A5")
	checkRegister(t, engine, A6, 011, "A6")
	checkRegister(t, engine, A7, 044, "A7")
}

var lnaExtendedMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("posValue", "w", []string{"0300000123456"}),
	tasm.NewSourceItem("negValue", "w", []string{"0400000000001"}),
	tasm.NewSourceItem("partValue", "w", []string{"0555577664444"}),

	segSourceItem(0),
	lnaSourceItemU("", jU, 0, 0, 0_377777),
	lnaSourceItemU("", jU, 1, 0, 0_477777),
	lnaSourceItemU("", jXU, 2, 0, 0_377777),
	lnaSourceItemU("", jXU, 3, 0, 0_477777),
	lnaSourceItemHIBDRef("", jW, 4, 0, 0, 0, 0, "posValue"),
	lnaSourceItemHIBDRef("", jW, 5, 0, 0, 0, 0, "negValue"),
	lnaSourceItemHIBDRef("", jT2, 6, 0, 0, 0, 0, "partValue"),
	lnaSourceItemHIBDRef("", jS5, 7, 0, 0, 0, 0, "partValue"),
	iarSourceItem("", 0),
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
	checkStopped(t, engine)
	checkRegister(t, engine, A0, 0_777777_400000, "A0")
	checkRegister(t, engine, A1, 0_777777_300000, "A1")
	checkRegister(t, engine, A2, 0_777777_400000, "A2")
	checkRegister(t, engine, A3, 0_300000, "A3")
	checkRegister(t, engine, A4, 0_477777_654321, "A4")
	checkRegister(t, engine, A5, 0_377777_777776, "A5")
	checkRegister(t, engine, A6, 011, "A6")
	checkRegister(t, engine, A7, 0_777777_777733, "A7")
}

var lnmaExtendedMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("posValue", "w", []string{"0300000123456"}),
	tasm.NewSourceItem("negValue", "w", []string{"0400000000001"}),
	tasm.NewSourceItem("partValue", "w", []string{"0555577664444"}),

	segSourceItem(0),
	lnmaSourceItemU("", jU, 0, 0, 0_377777),
	lnmaSourceItemU("", jU, 1, 0, 0_477777),
	lnmaSourceItemU("", jXU, 2, 0, 0_377777),
	lnmaSourceItemU("", jXU, 3, 0, 0_477777),
	lnmaSourceItemHIBDRef("", jW, 4, 0, 0, 0, 0, "posValue"),
	lnmaSourceItemHIBDRef("", jW, 5, 0, 0, 0, 0, "negValue"),
	lnmaSourceItemHIBDRef("", jT2, 6, 0, 0, 0, 0, "partValue"),
	lnmaSourceItemHIBDRef("", jS5, 7, 0, 0, 0, 0, "partValue"),
	iarSourceItem("", 0),
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
	checkStopped(t, engine)
	checkRegister(t, engine, A0, 0_777777_400000, "A0")
	checkRegister(t, engine, A1, 0_777777_300000, "A1")
	checkRegister(t, engine, A2, 0_777777_400000, "A2")
	checkRegister(t, engine, A3, 0_777777_477777, "A3")
	checkRegister(t, engine, A4, 0_477777_654321, "A4")
	checkRegister(t, engine, A5, 0_400000_000001, "A5")
	checkRegister(t, engine, A6, 0_777777_777766, "A6")
	checkRegister(t, engine, A7, 0_777777_777733, "A7")
}

var lrBasicMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("r7value", "qw", []string{"061", "062", "063", "064"}),
	tasm.NewSourceItem("r8value", "sw", []string{"01", "02", "03", "04", "05", "06"}),

	segSourceItem(0),
	lrSourceItemHIURef("", jQ3, 7, 0, 0, 0, "r7value"),
	lrSourceItemHIURef("", jXH2, 8, 0, 0, 0, "r8value"),
	iarSourceItem("", 0),
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
	segSourceItem(077),
	tasm.NewSourceItem("r5value", "tw", []string{"03210", "04000", "0123"}),
	tasm.NewSourceItem("r4value", "sw", []string{"01", "02", "03", "04", "05", "06"}),

	segSourceItem(0),
	lrSourceItemHIBDRef("", jT2, 5, 0, 0, 0, 0, "r5value"),
	lrSourceItemHIBDRef("", jXH2, 4, 0, 0, 0, 0, "r4value"),
	iarSourceItem("", 0),
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
	segSourceItem(077),
	tasm.NewSourceItem("data", "w", []string{"0112233445566"}),

	segSourceItem(0),
	lxSourceItemU("", jU, 1, 0, 0_377777),
	lxSourceItemHIURef("", jW, 15, 0, 0, 0, "data"),
	iarSourceItem("", 0),
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
	segSourceItem(0),
	lxSourceItemU("", jU, 1, 0, 05),
	iarSourceItem("", 0),
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
	segSourceItem(077),
	tasm.NewSourceItem("posValue", "w", []string{"0100200300400"}),
	tasm.NewSourceItem("negValue", "w", []string{"0500600700777"}),
	tasm.NewSourceItem("", "w", []string{"05"}),
	tasm.NewSourceItem("indAddr1", "w", []string{"0200000+indAddr2"}),
	tasm.NewSourceItem("indAddr2", "w", []string{"posValue"}),

	segSourceItem(0),
	dlSourceItemHIURef("", 4, 0, 0, 0, "posValue"),
	dlSourceItemHIURef("", 0, 0, 0, 1, "indAddr1"),
	dlSourceItemHIURef("", 2, 0, 0, 0, "posValue"),
	dlSourceItemHIURef("", 6, 0, 0, 0, "negValue"),
	dlSourceItemHIURef("", 10, 0, 0, 0, "posValue"),
	dlSourceItemHIURef("", 12, 0, 0, 0, "negValue"),
	iarSourceItem("", 0),
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
