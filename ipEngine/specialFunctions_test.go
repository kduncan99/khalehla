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

var exrExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXM, jU, rX7, zero, "02"}), //	EXR 2 past target
	tasm.NewSourceItem("", "fjaxu", []string{fLXI, jU, rX8, zero, "01"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXM, jU, rX8, zero, "00"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLR, jU, rR1, zero, "010"}), //	repeat 8 times
	tasm.NewSourceItem("", "fjaxhibd", []string{fEXR, jEXR, aEXR, rX7, zero, zero, rB2, "target"}),
	iarSourceItem("end", "0"),

	tasm.NewSourceItem("", ".SEG", []string{"002"}),
	tasm.NewSourceItem("target", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSAS, jH1, aSAS, rX8, "1", zero, rB3, "data"}),

	tasm.NewSourceItem("", ".SEG", []string{"003"}),
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
	checkRegister(t, engine, X8, 0_000001_000010, "X8")

	for ox := uint64(0); ox < 8; ox++ {
		checkMemory(t, engine, engine.baseRegisters[3].GetBankDescriptor().GetBaseAddress(), ox, 0_040040_777777)
	}

	for ox := uint64(8); ox < 10; ox++ {
		checkMemory(t, engine, engine.baseRegisters[3].GetBankDescriptor().GetBaseAddress(), ox, 0_000000_777777)
	}
}

var exrExtendedModeInvalidInstruction = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fEXR, jEXR, aEXR, rX7, zero, zero, rB0, "target"}),
	iarSourceItem("end", "1"),
	tasm.NewSourceItem("target", "fjaxu", []string{fLA, jU, rA3, zero, "0177777"}),
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
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXI, jU, rX8, zero, "01"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXM, jU, rX8, zero, "00"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLR, jU, rR1, zero, "020"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fEXR, jEXR, aEXR, zero, zero, zero, rB0, "target"}),
	iarSourceItem("badend", "1"),
	iarSourceItem("goodend", "0"),

	tasm.NewSourceItem("target", "fjaxhibd", []string{fTZ, jH2, aTZExtended, rX8, "1", zero, rB0, "data"}),

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
	checkRegister(t, engine, R1, 013, "R1")
}

// TODO NOP

var dcbExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fDCB, jDCB, rA3, zero, zero, zero, zero, "data"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fDCB, jDCB, rA4, zero, zero, zero, zero, "data+1"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fDCB, jDCB, rA5, zero, zero, zero, zero, "data+2"}),
	iarSourceItem("end", "0"),
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
	checkRegister(t, engine, A3, 46, "A3")
	checkRegister(t, engine, A4, 36, "A4")
	checkRegister(t, engine, A5, 1, "A5")
}

var rngbExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"002"}),
	tasm.NewSourceItem("data", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fRNGB, jRNGB, aRNGB, zero, zero, zero, rB2, "data"}),
	iarSourceItem("end", "0"),
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
	tasm.NewSourceItem("", ".SEG", []string{"002"}),
	tasm.NewSourceItem("data", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),
	tasm.NewSourceItem("", "w", []string{"0"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fRNGI, jRNGI, aRNGI, zero, zero, zero, rB2, "data"}),
	iarSourceItem("end", "0"),
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
