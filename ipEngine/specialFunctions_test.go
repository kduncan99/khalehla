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

var exBasicMode = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("data", "hw", []string{"0123456", "0654321"}),

	segSourceItem(0),
	lxiSourceItemU("", jU, 5, 0, 01),
	lxmSourceItemU("", jU, 5, 0, 04),
	exSourceItemHIURef("", 5, 1, 0, "target"),
	iarSourceItem("end", 0),

	tasm.NewSourceItem("target", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	lrSourceItemHIURef("", jH1, 3, 0, 0, 0, "data"),
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
	checkRegister(t, engine, pkg.X5, 0_000001_000005, "X5")
	checkRegister(t, engine, pkg.R3, 0_123456, "R3")
}

var exBasicModeIndirect = []*tasm.SourceItem{
	segSourceItem(12),
	lxiSourceItemU("", jU, 5, 0, 01),
	lxmSourceItemU("", jU, 5, 0, 04),
	exSourceItemHIURef("", 0, 0, 1, "ind1"),
	iarSourceItem("end", 0),

	segSourceItem(15),
	tasm.NewSourceItem("data1", "fjaxhiu", []string{zero, zero, zero, zero, zero, "1", "data2"}),
	tasm.NewSourceItem("data2", "fjaxhiu", []string{zero, zero, zero, zero, zero, zero, "data3"}),
	tasm.NewSourceItem("data3", "hw", []string{"0123456", "0654321"}),

	tasm.NewSourceItem("", ".SEG", []string{"14"}),
	tasm.NewSourceItem("ind1", "fjaxhiu", []string{zero, zero, zero, zero, zero, "1", "ind2"}),
	tasm.NewSourceItem("ind2", "fjaxhiu", []string{zero, zero, zero, zero, zero, zero, "target"}),
	lrSourceItemHIURef("target", jH1, 3, 0, 0, 1, "data1"),
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
	checkRegister(t, engine, pkg.R3, 0_123456, "R3")
}

var exExtendedMode = []*tasm.SourceItem{
	segSourceItem(2),
	tasm.NewSourceItem("data", "hw", []string{"0123456", "0654321"}),

	segSourceItem(0),
	lxiSourceItemU("", jU, 5, 0, 01),
	lxmSourceItemU("", jU, 5, 0, 04),
	exSourceItemHIBDRef("", 5, 1, 0, 0, "target"),
	iarSourceItem("end", 0),

	tasm.NewSourceItem("target", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	lrSourceItemHIBDRef("", jH1, 3, 0, 0, 0, 0, "data"),
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
	checkRegister(t, engine, pkg.X5, 0_000001_000005, "X5")
	checkRegister(t, engine, pkg.R3, 0_123456, "R3")
}

var exExtendedModeCascade = []*tasm.SourceItem{
	segSourceItem(02),
	tasm.NewSourceItem("data", "hw", []string{"0123456", "0654321"}),

	segSourceItem(0),
	exSourceItemHIBDRef("", 0, 0, 0, 4, "target1"),
	iarSourceItem("end", 0),

	tasm.NewSourceItem("", ".SEG", []string{"04"}),
	exSourceItemHIBDRef("target1", 0, 0, 0, 4, "target2"),
	exSourceItemHIBDRef("target1", 0, 0, 0, 4, "target3"),
	exSourceItemHIBDRef("target1", 0, 0, 0, 5, "target4"),

	tasm.NewSourceItem("", ".SEG", []string{"05"}),
	lrSourceItemHIBDRef("target4", jH1, 7, 0, 0, 0, 2, "data"),
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
	checkRegister(t, engine, pkg.R7, 0_123456, "R7")
}

var exExtendedModeJump = []*tasm.SourceItem{
	segSourceItem(0),
	exSourceItemHIBDRef("", 0, 0, 0, 0, "target"),
	iarSourceItem("badend", 1),
	jcSourceItemHIBDRef("target", 0, 0, 0, 0, "goodend"),
	iarSourceItem("goodend", 0),
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
	segSourceItem(0),
	exSourceItemHIBDRef("", 0, 0, 0, 0, "target"),
	iarSourceItem("badend", 1),
	iarSourceItem("goodend", 0),
	tskpSourceItemHIBDRef("target", jW, 0, 0, 0, 0, "target"),
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

var exrExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	lxmSourceItemU("", jU, 7, 0, 02),
	lxiSourceItemU("", jU, 8, 0, 01),
	lxmSourceItemU("", jU, 8, 0, 0),
	lrSourceItemU("", jU, 1, 0, 010),
	exrSourceItemHIBDRef("", 7, 0, 0, 2, "target"),
	iarSourceItem("end", 0),

	segSourceItem(02),
	tasm.NewSourceItem("target", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	sasSourceItemHIBDRef("", jH1, 8, 1, 0, 3, "data"),

	segSourceItem(03),
	tasm.NewSourceItem("data", "hw", []string{"0", "0777777"}),
	tasm.NewSourceItem("", "hw", []string{"0", "0777777"}),
	tasm.NewSourceItem("", "hw", []string{"0", "0777777"}),
	tasm.NewSourceItem("", "hw", []string{"0", "0777777"}),
	tasm.NewSourceItem("", "hw", []string{"0", "0777777"}),
	tasm.NewSourceItem("", "hw", []string{"0", "0777777"}),
	tasm.NewSourceItem("", "hw", []string{"0", "0777777"}),
	tasm.NewSourceItem("", "hw", []string{"0", "0777777"}),
	tasm.NewSourceItem("", "hw", []string{"0", "0777777"}),
	tasm.NewSourceItem("", "hw", []string{"0", "0777777"}),
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
	checkStopped(t, engine)
	checkRegister(t, engine, pkg.X8, 0_000001_000010, "X8")

	for ox := uint64(0); ox < 8; ox++ {
		checkMemory(t, engine, engine.baseRegisters[3].GetBankDescriptor().GetBaseAddress(), ox, 0_040040_777777)
	}

	for ox := uint64(8); ox < 10; ox++ {
		checkMemory(t, engine, engine.baseRegisters[3].GetBankDescriptor().GetBaseAddress(), ox, 0_000000_777777)
	}
}

var exrExtendedModeInvalidInstruction = []*tasm.SourceItem{
	segSourceItem(0),
	exrSourceItemHIBDRef("", 7, 0, 0, 0, "target"),
	iarSourceItem("end", 1),
	laSourceItemU("target", jU, 3, 0, 0_177777),
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
	lxiSourceItemU("", jU, 8, 0, 01),
	lxmSourceItemU("", jU, 8, 0, 0),
	lrSourceItemU("", jU, 1, 0, 020),
	exrSourceItemHIBDRef("", 0, 0, 0, 0, "target"),
	iarSourceItem("badend", 1),
	iarSourceItem("goodend", 0),

	tzSourceItemHIBDRef("target", jH2, 8, 1, 0, 0, "data"),

	tasm.NewSourceItem("data", "hw", []string{"01", "04"}),
	tasm.NewSourceItem("", "hw", []string{"01", "03"}),
	tasm.NewSourceItem("", "hw", []string{"01", "02"}),
	tasm.NewSourceItem("", "hw", []string{"01", "01"}),
	tasm.NewSourceItem("", "hw", []string{"01", "0"}),
	tasm.NewSourceItem("", "hw", []string{"01", "022"}),
	tasm.NewSourceItem("", "hw", []string{"01", "023"}),
	tasm.NewSourceItem("", "hw", []string{"01", "024"}),
	tasm.NewSourceItem("", "hw", []string{"01", "025"}),
	tasm.NewSourceItem("", "hw", []string{"01", "0"}),
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
	checkStopped(t, engine)
	checkRegister(t, engine, pkg.R1, 013, "R1")
}

// TODO NOP

var dcbExtendedMode = []*tasm.SourceItem{
	segSourceItem(0),
	dcbSourceItemHIBDRef("", 3, 0, 0, 0, 0, "data"),
	dcbSourceItemHIBDRef("", 4, 0, 0, 0, 0, "data+1"),
	dcbSourceItemHIBDRef("", 5, 0, 0, 0, 0, "data+2"),
	iarSourceItem("end", 0),
	tasm.NewSourceItem("data", "w", []string{"030405030405"}),
	tasm.NewSourceItem("", "w", []string{"0777777777777"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"01"}),
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
	checkStopped(t, engine)
	checkRegister(t, engine, pkg.A3, 46, "A3")
	checkRegister(t, engine, pkg.A4, 36, "A4")
	checkRegister(t, engine, pkg.A5, 1, "A5")
}

var rngbExtendedMode = []*tasm.SourceItem{
	segSourceItem(02),
	tasm.NewSourceItem("data", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),

	segSourceItem(0),
	rngbSourceItemHIBDRef("", 0, 0, 0, 2, "data"),
	iarSourceItem("end", 0),
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
	checkStopped(t, engine)
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
	tasm.NewSourceItem("data", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),

	segSourceItem(0),
	rngiSourceItemHIBDRef("", 0, 0, 0, 2, "data"),
	iarSourceItem("end", 0),
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
	checkStopped(t, engine)
	storage := engine.baseRegisters[2].GetStorage()
	for mx := 0; mx < 4; mx++ {
		value := storage[mx]
		fmt.Printf("%04o: %012o\n", mx, value)
		if value&0_740000_000000 != 0 {
			t.Fatalf("Expected bits 0-3 to be zero")
		}
	}
}
