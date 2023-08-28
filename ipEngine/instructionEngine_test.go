// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
	"khalehla/tasm"
	"testing"
)

var partialWordLoadsBasicThirdWord = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("data1", "w", []string{"0556677001122"}),
	tasm.NewSourceItem("data2", "hw", []string{"0506070", "0507777"}),
	tasm.NewSourceItem("data3", "hw", []string{"0223344", "0221100"}),
	tasm.NewSourceItem("data4", "tw", []string{"01111", "02222", "03333"}),
	tasm.NewSourceItem("data5", "tw", []string{"05500", "06600", "07700"}),

	segSourceItem(0),
	laSourceItemHIURef("", jW, 0, 0, 0, 0, "data1"),
	laSourceItemHIURef("", jH1, 1, 0, 0, 0, "data2"),
	laSourceItemHIURef("", jH2, 2, 0, 0, 0, "data2"),
	laSourceItemHIURef("", jXH1, 3, 0, 0, 0, "data2"),
	laSourceItemHIURef("", jXH1, 4, 0, 0, 0, "data3"),
	laSourceItemHIURef("", jXH2, 5, 0, 0, 0, "data2"),
	laSourceItemHIURef("", jXH2, 6, 0, 0, 0, "data3"),
	laSourceItemHIURef("", jT1, 7, 0, 0, 0, "data4"),
	laSourceItemHIURef("", jT1, 8, 0, 0, 0, "data5"),
	laSourceItemHIURef("", jT2, 9, 0, 0, 0, "data4"),
	laSourceItemHIURef("", jT2, 10, 0, 0, 0, "data5"),
	laSourceItemHIURef("", jT3, 11, 0, 0, 0, "data4"),
	laSourceItemHIURef("", jT3, 12, 0, 0, 0, "data5"),
	iarSourceItem("", 0),
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.A0, 0_556677_001122, "A0")
	checkRegister(t, engine, pkg.A1, 0_506070, "A1")
	checkRegister(t, engine, pkg.A2, 0_507777, "A2")
	checkRegister(t, engine, pkg.A3, 0_777777_506070, "A3")
	checkRegister(t, engine, pkg.A4, 0_223344, "A4")
	checkRegister(t, engine, pkg.A5, 0_777777_507777, "A5")
	checkRegister(t, engine, pkg.A6, 0_221100, "A6")
	checkRegister(t, engine, pkg.A7, 01111, "A7")
	checkRegister(t, engine, pkg.A8, 0_777777_775500, "A8")
	checkRegister(t, engine, pkg.A9, 02222, "A9")
	checkRegister(t, engine, pkg.A10, 0_777777_776600, "A10")
	checkRegister(t, engine, pkg.A11, 03333, "A11")
	checkRegister(t, engine, pkg.A12, 0_777777_777700, "A12")
}

var PartialWordLoadsBasicQuarterWord = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("data1", "qw", []string{"0400", "0501", "0677", "0777"}),
	tasm.NewSourceItem("data2", "sw", []string{"012", "034", "056", "075", "042", "010"}),

	segSourceItem(0),
	lrSourceItemHIURef("", jQ1, 0, 0, 0, 0, "data1"),
	lrSourceItemHIURef("", jQ2, 1, 0, 0, 0, "data1"),
	lrSourceItemHIURef("", jQ3, 2, 0, 0, 0, "data1"),
	lrSourceItemHIURef("", jQ4, 3, 0, 0, 0, "data1"),
	lrSourceItemHIURef("", jS1, 4, 0, 0, 0, "data2"),
	lrSourceItemHIURef("", jS2, 5, 0, 0, 0, "data2"),
	lrSourceItemHIURef("", jS3, 6, 0, 0, 0, "data2"),
	lrSourceItemHIURef("", jS4, 7, 0, 0, 0, "data2"),
	lrSourceItemHIURef("", jS5, 8, 0, 0, 0, "data2"),
	lrSourceItemHIURef("", jS6, 9, 0, 0, 0, "data2"),
	iarSourceItem("", 0),
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.R0, 0400, "R0")
	checkRegister(t, engine, pkg.R1, 0501, "R1")
	checkRegister(t, engine, pkg.R2, 0677, "R2")
	checkRegister(t, engine, pkg.R3, 0777, "R3")
	checkRegister(t, engine, pkg.R4, 012, "R4")
	checkRegister(t, engine, pkg.R5, 034, "R5")
	checkRegister(t, engine, pkg.R6, 056, "R6")
	checkRegister(t, engine, pkg.R7, 075, "R7")
	checkRegister(t, engine, pkg.R8, 042, "R8")
	checkRegister(t, engine, pkg.R9, 010, "R9")
}

var partialWordStoresBasicThirdWord = []*tasm.SourceItem{
	segSourceItem(13),
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

	segSourceItem(12),
	laSourceItemHIURef("", jW, 0, 0, 0, 0, "data1"),
	lxiSourceItemU("", jU, 1, 0, 1),
	lxmSourceItemU("", jU, 1, 0, 0),
	laSourceItemHIURef("", jW, 0, 0, 0, 0, "data1"),
	saSourceItemHIURef("", jW, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jH1, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jH2, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jXH1, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jXH2, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jT1, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jT2, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jT3, 0, 1, 1, 0, "data2"),
	iarSourceItem("", 0),
}

func Test_PartialWordStores_BasicThirdWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", partialWordStoresBasicThirdWord)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), false)

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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	dataBankAddr := e.GetBanks()[0601015].GetBankDescriptor().GetBaseAddress()
	checkMemory(t, engine, dataBankAddr, 1, 0_444444_444444)
	checkMemory(t, engine, dataBankAddr, 2, 0_444444_333333)
	checkMemory(t, engine, dataBankAddr, 3, 0_333333_444444)
	checkMemory(t, engine, dataBankAddr, 4, 0_444444_333333)
	checkMemory(t, engine, dataBankAddr, 5, 0_333333_444444)
	checkMemory(t, engine, dataBankAddr, 6, 0_4444_3333_3333)
	checkMemory(t, engine, dataBankAddr, 7, 0_3333_4444_3333)
	checkMemory(t, engine, dataBankAddr, 8, 0_3333_3333_4444)
}

var partialWordStoresBasicQuarterWord = []*tasm.SourceItem{
	segSourceItem(13),
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

	segSourceItem(12),
	laSourceItemHIURef("", jW, 0, 0, 0, 0, "data1"),
	lxiSourceItemU("", jU, 1, 0, 1),
	lxmSourceItemU("", jU, 1, 0, 0),
	laSourceItemHIURef("", jW, 0, 0, 0, 0, "data1"),
	saSourceItemHIURef("", jQ1, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jQ2, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jQ3, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jQ4, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jS1, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jS2, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jS3, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jS4, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jS5, 0, 1, 1, 0, "data2"),
	saSourceItemHIURef("", jS6, 0, 1, 1, 0, "data2"),
	iarSourceItem("", 0),
}

func Test_PartialWordStores_BasicQuarterWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", partialWordStoresBasicQuarterWord)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), false)

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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	dataBankAddr := e.GetBanks()[0601015].GetBankDescriptor().GetBaseAddress()
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

//	TODO basic mode GRS addressing

//	TODO basic mode index register handling

//	TODO basic mode addressing topics (incl. indirect)

var partialWordLoadsExtendedThirdWord = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("data1", "w", []string{"0556677001122"}),
	tasm.NewSourceItem("data2", "hw", []string{"0506070", "0507777"}),
	tasm.NewSourceItem("data3", "hw", []string{"0223344", "0221100"}),
	tasm.NewSourceItem("data4", "tw", []string{"01111", "02222", "03333"}),
	tasm.NewSourceItem("data5", "tw", []string{"05500", "06600", "07700"}),

	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 0, 0, 0, 0, 0, "data1"),
	laSourceItemHIBDRef("", jH1, 1, 0, 0, 0, 0, "data2"),
	laSourceItemHIBDRef("", jH2, 2, 0, 0, 0, 0, "data2"),
	laSourceItemHIBDRef("", jXH1, 3, 0, 0, 0, 0, "data2"),
	laSourceItemHIBDRef("", jXH1, 4, 0, 0, 0, 0, "data3"),
	laSourceItemHIBDRef("", jXH2, 5, 0, 0, 0, 0, "data2"),
	laSourceItemHIBDRef("", jXH2, 6, 0, 0, 0, 0, "data3"),
	laSourceItemHIBDRef("", jT1, 7, 0, 0, 0, 0, "data4"),
	laSourceItemHIBDRef("", jT1, 8, 0, 0, 0, 0, "data5"),
	laSourceItemHIBDRef("", jT2, 9, 0, 0, 0, 0, "data4"),
	laSourceItemHIBDRef("", jT2, 10, 0, 0, 0, 0, "data5"),
	laSourceItemHIBDRef("", jT3, 11, 0, 0, 0, 0, "data4"),
	laSourceItemHIBDRef("", jT3, 12, 0, 0, 0, 0, "data5"),
	iarSourceItem("", 0),
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.A0, 0_556677_001122, "A0")
	checkRegister(t, engine, pkg.A1, 0_506070, "A1")
	checkRegister(t, engine, pkg.A2, 0_507777, "A2")
	checkRegister(t, engine, pkg.A3, 0_777777_506070, "A3")
	checkRegister(t, engine, pkg.A4, 0_223344, "A4")
	checkRegister(t, engine, pkg.A5, 0_777777_507777, "A5")
	checkRegister(t, engine, pkg.A6, 0_221100, "A6")
	checkRegister(t, engine, pkg.A7, 01111, "A7")
	checkRegister(t, engine, pkg.A8, 0_777777_775500, "A8")
	checkRegister(t, engine, pkg.A9, 02222, "A9")
	checkRegister(t, engine, pkg.A10, 0_777777_776600, "A10")
	checkRegister(t, engine, pkg.A11, 03333, "A11")
	checkRegister(t, engine, pkg.A12, 0_777777_777700, "A12")
}

var PartialWordLoadsExtendedQuarterWord = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("data1", "qw", []string{"0400", "0501", "0677", "0777"}),
	tasm.NewSourceItem("data2", "sw", []string{"012", "034", "056", "075", "042", "010"}),

	segSourceItem(0),
	lrSourceItemHIBDRef("", jQ1, 0, 0, 0, 0, 0, "data1"),
	lrSourceItemHIBDRef("", jQ2, 1, 0, 0, 0, 0, "data1"),
	lrSourceItemHIBDRef("", jQ3, 2, 0, 0, 0, 0, "data1"),
	lrSourceItemHIBDRef("", jQ4, 3, 0, 0, 0, 0, "data1"),
	lrSourceItemHIBDRef("", jS1, 4, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBDRef("", jS2, 5, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBDRef("", jS3, 6, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBDRef("", jS4, 7, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBDRef("", jS5, 8, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBDRef("", jS6, 9, 0, 0, 0, 0, "data2"),
	iarSourceItem("", 0),
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkRegister(t, engine, pkg.R0, 0400, "R0")
	checkRegister(t, engine, pkg.R1, 0501, "R1")
	checkRegister(t, engine, pkg.R2, 0677, "R2")
	checkRegister(t, engine, pkg.R3, 0777, "R3")
	checkRegister(t, engine, pkg.R4, 012, "R4")
	checkRegister(t, engine, pkg.R5, 034, "R5")
	checkRegister(t, engine, pkg.R6, 056, "R6")
	checkRegister(t, engine, pkg.R7, 075, "R7")
	checkRegister(t, engine, pkg.R8, 042, "R8")
	checkRegister(t, engine, pkg.R9, 010, "R9")
}

var partialWordStoresExtendedThirdWord = []*tasm.SourceItem{
	segSourceItem(2),
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

	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 0, 0, 0, 0, 2, "data1"),
	lxiSourceItemU("", jU, 1, 0, 1),
	lxmSourceItemU("", jU, 1, 0, 0),
	saSourceItemHIBDRef("", jW, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jH1, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jH2, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jXH1, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jXH2, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jT1, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jT2, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jT3, 0, 1, 1, 0, 2, "data2"),
	iarSourceItem("", 0),
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
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
	segSourceItem(2),
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

	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 0, 0, 0, 0, 2, "data1"),
	lxiSourceItemU("", jU, 1, 0, 1),
	lxmSourceItemU("", jU, 1, 0, 0),
	laSourceItemHIBDRef("", jW, 0, 0, 0, 0, 2, "data1"),
	saSourceItemHIBDRef("", jQ1, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jQ2, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jQ3, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jQ4, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jS1, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jS2, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jS3, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jS4, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jS5, 0, 1, 1, 0, 2, "data2"),
	saSourceItemHIBDRef("", jS6, 0, 1, 1, 0, 2, "data2"),
	iarSourceItem("", 0),
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
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

// TODO extended mode GRS addressing
var grsAddressingExtended = []*tasm.SourceItem{
	segSourceItem(2),

	segSourceItem(0),
	lxSourceItemU("", jU, 5, 0, 42),
	laSourceItemU("", jW, 5, 0, grsX5),
	laSourceItemU("", jS3, 6, 0, grsX5), // partial word ignored, we're register-to-register
	iarSourceItem("", 0),
}

func Test_GRSAddressing_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", grsAddressingExtended)
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
	checkRegister(t, engine, pkg.A5, 42, "A5")
	checkRegister(t, engine, pkg.A6, 42, "A5")
}

//	TODO extended mode index register handling

//	TODO extended mode addressing across multiple banks
