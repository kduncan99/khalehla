// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

var partialWordLoadsBasicThirdWord = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data1", "w", []string{"0556677001122"}),
	tasm.NewSourceItem("data2", "hw", []string{"0506070", "0507777"}),
	tasm.NewSourceItem("data3", "hw", []string{"0223344", "0221100"}),
	tasm.NewSourceItem("data4", "tw", []string{"01111", "02222", "03333"}),
	tasm.NewSourceItem("data5", "tw", []string{"05500", "06600", "07700"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jW, rA0, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jH1, rA1, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jH2, rA2, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jXH1, rA3, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jXH1, rA4, zero, zero, zero, "data3"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jXH2, rA5, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jXH2, rA6, zero, zero, zero, "data3"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT1, rA7, zero, zero, zero, "data4"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT1, rA8, zero, zero, zero, "data5"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT2, rA9, zero, zero, zero, "data4"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT2, rA10, zero, zero, zero, "data5"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT3, rA11, zero, zero, zero, "data4"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT3, rA12, zero, zero, zero, "data5"}),
	IARSourceItem("", "0"),
}

func Test_PartialWordLoads_BasicThirdWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", partialWordLoadsBasicThirdWord)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStopped(t, engine)
	checkRegister(t, engine, A0, 0_556677_001122, "A0")
	checkRegister(t, engine, A1, 0_506070, "A1")
	checkRegister(t, engine, A2, 0_507777, "A2")
	checkRegister(t, engine, A3, 0_777777_506070, "A3")
	checkRegister(t, engine, A4, 0_223344, "A4")
	checkRegister(t, engine, A5, 0_777777_507777, "A5")
	checkRegister(t, engine, A6, 0_221100, "A6")
	checkRegister(t, engine, A7, 01111, "A7")
	checkRegister(t, engine, A8, 0_777777_775500, "A8")
	checkRegister(t, engine, A9, 02222, "A9")
	checkRegister(t, engine, A10, 0_777777_776600, "A10")
	checkRegister(t, engine, A11, 03333, "A11")
	checkRegister(t, engine, A12, 0_777777_777700, "A12")
}

var PartialWordLoadsBasicQuarterWord = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data1", "qw", []string{"0400", "0501", "0677", "0777"}),
	tasm.NewSourceItem("data2", "sw", []string{"012", "034", "056", "075", "042", "010"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jQ1, rA0, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jQ2, rA1, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jQ3, rA2, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jQ4, rA3, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS1, rA4, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS2, rA5, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS3, rA6, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS4, rA7, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS5, rA8, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS6, rA9, zero, zero, zero, "data2"}),
	IARSourceItem("", "0"),
}

func Test_PartialWordLoads_BasicQuarterWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", PartialWordLoadsBasicQuarterWord)
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
	checkRegister(t, engine, R0, 0400, "R0")
	checkRegister(t, engine, R1, 0501, "R1")
	checkRegister(t, engine, R2, 0677, "R2")
	checkRegister(t, engine, R3, 0777, "R3")
	checkRegister(t, engine, R4, 012, "R4")
	checkRegister(t, engine, R5, 034, "R5")
	checkRegister(t, engine, R6, 056, "R6")
	checkRegister(t, engine, R7, 075, "R7")
	checkRegister(t, engine, R8, 042, "R8")
	checkRegister(t, engine, R9, 010, "R9")
}

//	TODO basic mode partial word testing, stores

//	TODO basic mode index register handling

//	TODO basic mode addressing topics (incl. indirect)

var partialWordLoadsExtendedThirdWord = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data1", "w", []string{"0556677001122"}),
	tasm.NewSourceItem("data2", "hw", []string{"0506070", "0507777"}),
	tasm.NewSourceItem("data3", "hw", []string{"0223344", "0221100"}),
	tasm.NewSourceItem("data4", "tw", []string{"01111", "02222", "03333"}),
	tasm.NewSourceItem("data5", "tw", []string{"05500", "06600", "07700"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jW, rA0, zero, zero, zero, rB0, "data1"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jH1, rA1, zero, zero, zero, rB0, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jH2, rA2, zero, zero, zero, rB0, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jXH1, rA3, zero, zero, zero, rB0, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jXH1, rA4, zero, zero, zero, rB0, "data3"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jXH2, rA5, zero, zero, zero, rB0, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jXH2, rA6, zero, zero, zero, rB0, "data3"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jT1, rA7, zero, zero, zero, rB0, "data4"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jT1, rA8, zero, zero, zero, rB0, "data5"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jT2, rA9, zero, zero, zero, rB0, "data4"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jT2, rA10, zero, zero, zero, rB0, "data5"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jT3, rA11, zero, zero, zero, rB0, "data4"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jT3, rA12, zero, zero, zero, rB0, "data5"}),
	IARSourceItem("", "0"),
}

func Test_PartialWordLoads_ExtendedThirdWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", partialWordLoadsExtendedThirdWord)
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
	checkRegister(t, engine, A0, 0_556677_001122, "A0")
	checkRegister(t, engine, A1, 0_506070, "A1")
	checkRegister(t, engine, A2, 0_507777, "A2")
	checkRegister(t, engine, A3, 0_777777_506070, "A3")
	checkRegister(t, engine, A4, 0_223344, "A4")
	checkRegister(t, engine, A5, 0_777777_507777, "A5")
	checkRegister(t, engine, A6, 0_221100, "A6")
	checkRegister(t, engine, A7, 01111, "A7")
	checkRegister(t, engine, A8, 0_777777_775500, "A8")
	checkRegister(t, engine, A9, 02222, "A9")
	checkRegister(t, engine, A10, 0_777777_776600, "A10")
	checkRegister(t, engine, A11, 03333, "A11")
	checkRegister(t, engine, A12, 0_777777_777700, "A12")
}

var PartialWordLoadsExtendedQuarterWord = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data1", "qw", []string{"0400", "0501", "0677", "0777"}),
	tasm.NewSourceItem("data2", "sw", []string{"012", "034", "056", "075", "042", "010"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jQ1, rA0, zero, zero, zero, rB0, "data1"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jQ2, rA1, zero, zero, zero, rB0, "data1"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jQ3, rA2, zero, zero, zero, rB0, "data1"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jQ4, rA3, zero, zero, zero, rB0, "data1"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jS1, rA4, zero, zero, zero, rB0, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jS2, rA5, zero, zero, zero, rB0, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jS3, rA6, zero, zero, zero, rB0, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jS4, rA7, zero, zero, zero, rB0, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jS5, rA8, zero, zero, zero, rB0, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLR, jS6, rA9, zero, zero, zero, rB0, "data2"}),
	IARSourceItem("", "0"),
}

func Test_PartialWordLoads_ExtendedQuarterWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", PartialWordLoadsExtendedQuarterWord)
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
	checkRegister(t, engine, R0, 0400, "R0")
	checkRegister(t, engine, R1, 0501, "R1")
	checkRegister(t, engine, R2, 0677, "R2")
	checkRegister(t, engine, R3, 0777, "R3")
	checkRegister(t, engine, R4, 012, "R4")
	checkRegister(t, engine, R5, 034, "R5")
	checkRegister(t, engine, R6, 056, "R6")
	checkRegister(t, engine, R7, 075, "R7")
	checkRegister(t, engine, R8, 042, "R8")
	checkRegister(t, engine, R9, 010, "R9")
}

var partialWordStoresExtendedThirdWord = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"002"}),
	tasm.NewSourceItem("data1", "w", []string{"0444444444444"}),
	tasm.NewSourceItem("data2", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jW, rA0, zero, zero, zero, rB2, "data1"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXI, jU, rX1, zero, "1"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXM, jU, rX1, zero, "0"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jW, rA0, zero, zero, zero, rB2, "data1"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jW, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jH1, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jH2, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jXH1, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jXH2, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jT1, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jT2, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jT3, rA0, rX1, "1", zero, rB2, "data2"}),
	IARSourceItem("", "0"),
}

func Test_PartialWordStores_ExtendedThirdWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", partialWordStoresExtendedThirdWord)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), true)

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
	dataBankAddr := e.GetBanks()[0601002].GetBankDescriptor().GetBaseAddress()
	checkMemory(t, engine, dataBankAddr, 1, 0_444444_444444)
	checkMemory(t, engine, dataBankAddr, 2, 0_444444_333333)
	checkMemory(t, engine, dataBankAddr, 3, 0_333333_444444)
	checkMemory(t, engine, dataBankAddr, 4, 0_444444_333333)
	checkMemory(t, engine, dataBankAddr, 5, 0_333333_444444)
	checkMemory(t, engine, dataBankAddr, 6, 0_4444_3333_3333)
	checkMemory(t, engine, dataBankAddr, 7, 0_3333_4444_3333)
	checkMemory(t, engine, dataBankAddr, 8, 0_3333_3333_4444)
}

var partialWordStoresExtendedQuarterWord = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"002"}),
	tasm.NewSourceItem("data1", "w", []string{"0444444444444"}),
	tasm.NewSourceItem("data2", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),
	tasm.NewSourceItem("", "w", []string{"0333333333333"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jW, rA0, zero, zero, zero, rB2, "data1"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXI, jU, rX1, zero, "1"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLXM, jU, rX1, zero, "0"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jW, rA0, zero, zero, zero, rB2, "data1"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jQ1, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jQ2, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jQ3, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jQ4, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jS1, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jS2, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jS3, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jS4, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jS5, rA0, rX1, "1", zero, rB2, "data2"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fSA, jS6, rA0, rX1, "1", zero, rB2, "data2"}),
	IARSourceItem("", "0"),
}

func Test_PartialWordStores_ExtendedQuarterWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", partialWordStoresExtendedQuarterWord)
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
	dataBankAddr := e.GetBanks()[0601002].GetBankDescriptor().GetBaseAddress()
	checkMemory(t, engine, dataBankAddr, 1, 0_444_333_333_333)
	checkMemory(t, engine, dataBankAddr, 2, 0_333_444_333_333)
	checkMemory(t, engine, dataBankAddr, 3, 0_333_333_444_333)
	checkMemory(t, engine, dataBankAddr, 4, 0_333_333_333_444)
	checkMemory(t, engine, dataBankAddr, 5, 0_443333_333333)
	checkMemory(t, engine, dataBankAddr, 6, 0_334433_333333)
	checkMemory(t, engine, dataBankAddr, 7, 0_333344_333333)
	checkMemory(t, engine, dataBankAddr, 8, 0_333333_443333)
	checkMemory(t, engine, dataBankAddr, 9, 0_333333_334433)
	checkMemory(t, engine, dataBankAddr, 10, 0_333333_333344)
}

//	TODO extended mode index register handling

//	TODO extended mode addressing across multiple banks
