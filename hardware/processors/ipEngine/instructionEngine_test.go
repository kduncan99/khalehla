// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"testing"

	"khalehla/common"
	"khalehla/tasm"
)

var partialWordLoadsBasicThirdWord = []*tasm.SourceItem{
	segSourceItem(077),
	labelDataSourceItem("data1", []uint64{0_556677_001122}),
	labelDataSourceItem("data2", []uint64{0_506070, 0_507777}),
	labelDataSourceItem("data3", []uint64{0_223344, 0_221100}),
	labelDataSourceItem("data4", []uint64{01111, 02222, 03333}),
	labelDataSourceItem("data5", []uint64{05500, 06600, 07700}),

	segSourceItem(0),
	laSourceItemHIRef(jW, regA0, 0, 0, 0, "data1"),
	laSourceItemHIRef(jH1, regA1, 0, 0, 0, "data2"),
	laSourceItemHIRef(jH2, regA2, 0, 0, 0, "data2"),
	laSourceItemHIRef(jXH1, regA3, 0, 0, 0, "data2"),
	laSourceItemHIRef(jXH1, regA4, 0, 0, 0, "data3"),
	laSourceItemHIRef(jXH2, regA5, 0, 0, 0, "data2"),
	laSourceItemHIRef(jXH2, regA6, 0, 0, 0, "data3"),
	laSourceItemHIRef(jT1, regA7, 0, 0, 0, "data4"),
	laSourceItemHIRef(jT1, regA8, 0, 0, 0, "data5"),
	laSourceItemHIRef(jT2, regA9, 0, 0, 0, "data4"),
	laSourceItemHIRef(jT2, regA10, 0, 0, 0, "data5"),
	laSourceItemHIRef(jT3, regA11, 0, 0, 0, "data4"),
	laSourceItemHIRef(jT3, regA12, 0, 0, 0, "data5"),
	iarSourceItem(0),
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
	checkRegister(t, engine, common.A0, 0_556677_001122)
	checkRegister(t, engine, common.A1, 0_506070)
	checkRegister(t, engine, common.A2, 0_507777)
	checkRegister(t, engine, common.A3, 0_777777_506070)
	checkRegister(t, engine, common.A4, 0_223344)
	checkRegister(t, engine, common.A5, 0_777777_507777)
	checkRegister(t, engine, common.A6, 0_221100)
	checkRegister(t, engine, common.A7, 01111)
	checkRegister(t, engine, common.A8, 0_777777_775500)
	checkRegister(t, engine, common.A9, 02222)
	checkRegister(t, engine, common.A10, 0_777777_776600)
	checkRegister(t, engine, common.A11, 03333)
	checkRegister(t, engine, common.A12, 0_777777_777700)
}

var PartialWordLoadsBasicQuarterWord = []*tasm.SourceItem{
	segSourceItem(077),
	tasm.NewSourceItem("data1", "qw", []string{"0400", "0501", "0677", "0777"}),
	tasm.NewSourceItem("data2", "sw", []string{"012", "034", "056", "075", "042", "010"}),

	segSourceItem(0),
	lrSourceItemHIRef(jQ1, regR0, 0, 0, 0, "data1"),
	lrSourceItemHIRef(jQ2, regR1, 0, 0, 0, "data1"),
	lrSourceItemHIRef(jQ3, regR2, 0, 0, 0, "data1"),
	lrSourceItemHIRef(jQ4, regR3, 0, 0, 0, "data1"),
	lrSourceItemHIRef(jS1, regR4, 0, 0, 0, "data2"),
	lrSourceItemHIRef(jS2, regR5, 0, 0, 0, "data2"),
	lrSourceItemHIRef(jS3, regR6, 0, 0, 0, "data2"),
	lrSourceItemHIRef(jS4, regR7, 0, 0, 0, "data2"),
	lrSourceItemHIRef(jS5, regR8, 0, 0, 0, "data2"),
	lrSourceItemHIRef(jS6, regR9, 0, 0, 0, "data2"),
	iarSourceItem(0),
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
	checkRegister(t, engine, common.R0, 0400)
	checkRegister(t, engine, common.R1, 0501)
	checkRegister(t, engine, common.R2, 0677)
	checkRegister(t, engine, common.R3, 0777)
	checkRegister(t, engine, common.R4, 012)
	checkRegister(t, engine, common.R5, 034)
	checkRegister(t, engine, common.R6, 056)
	checkRegister(t, engine, common.R7, 075)
	checkRegister(t, engine, common.R8, 042)
	checkRegister(t, engine, common.R9, 010)
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
	laSourceItemHIRef(jW, 0, 0, 0, 0, "data1"),
	lxiSourceItemU(jU, regX1, 0, 1),
	lxmSourceItemU(jU, regX1, 0, 0),
	laSourceItemHIRef(jW, regA0, 0, 0, 0, "data1"),
	saSourceItemHIRef(jW, regA0, 1, 1, 0, "data2"),
	saSourceItemHIRef(jH1, regA0, 1, 1, 0, "data2"),
	saSourceItemHIRef(jH2, regA0, 1, 1, 0, "data2"),
	saSourceItemHIRef(jXH1, regA0, 1, 1, 0, "data2"),
	saSourceItemHIRef(jXH2, regA0, 1, 1, 0, "data2"),
	saSourceItemHIRef(jT1, regA0, 1, 1, 0, "data2"),
	saSourceItemHIRef(jT2, regA0, 1, 1, 0, "data2"),
	saSourceItemHIRef(jT3, regA0, 1, 1, 0, "data2"),
	iarSourceItem(0),
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
	labelDataSourceItem("data1", []uint64{0_444444_444444}),
	labelSourceItem("data2"),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),

	segSourceItem(12),
	laSourceItemHIRef(jW, regA0, 0, 0, 0, "data1"),
	lxiSourceItemU(jU, regX1, 0, 1),
	lxmSourceItemU(jU, regX1, 0, 0),
	laSourceItemHIRef(jW, regA0, regX0, 0, 0, "data1"),
	saSourceItemHIRef(jQ1, regA0, regX1, 1, 0, "data2"),
	saSourceItemHIRef(jQ2, regA0, regX1, 1, 0, "data2"),
	saSourceItemHIRef(jQ3, regA0, regX1, 1, 0, "data2"),
	saSourceItemHIRef(jQ4, regA0, regX1, 1, 0, "data2"),
	saSourceItemHIRef(jS1, regA0, regX1, 1, 0, "data2"),
	saSourceItemHIRef(jS2, regA0, regX1, 1, 0, "data2"),
	saSourceItemHIRef(jS3, regA0, regX1, 1, 0, "data2"),
	saSourceItemHIRef(jS4, regA0, regX1, 1, 0, "data2"),
	saSourceItemHIRef(jS5, regA0, regX1, 1, 0, "data2"),
	saSourceItemHIRef(jS6, regA0, regX1, 1, 0, "data2"),
	iarSourceItem(0),
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
	labelDataSourceItem("data1", []uint64{0556677001122}),
	labelDataSourceItem("data2", []uint64{0506070, 0507777}),
	labelDataSourceItem("data3", []uint64{0223344, 0221100}),
	labelDataSourceItem("data4", []uint64{01111, 02222, 03333}),
	labelDataSourceItem("data5", []uint64{05500, 06600, 07700}),

	segSourceItem(0),
	laSourceItemHIBRef(jW, regA0, 0, 0, 0, 0, "data1"),
	laSourceItemHIBRef(jH1, regA1, 0, 0, 0, 0, "data2"),
	laSourceItemHIBRef(jH2, regA2, 0, 0, 0, 0, "data2"),
	laSourceItemHIBRef(jXH1, regA3, 0, 0, 0, 0, "data2"),
	laSourceItemHIBRef(jXH1, regA4, 0, 0, 0, 0, "data3"),
	laSourceItemHIBRef(jXH2, regA5, 0, 0, 0, 0, "data2"),
	laSourceItemHIBRef(jXH2, regA6, 0, 0, 0, 0, "data3"),
	laSourceItemHIBRef(jT1, regA7, 0, 0, 0, 0, "data4"),
	laSourceItemHIBRef(jT1, regA8, 0, 0, 0, 0, "data5"),
	laSourceItemHIBRef(jT2, regA9, 0, 0, 0, 0, "data4"),
	laSourceItemHIBRef(jT2, regA10, 0, 0, 0, 0, "data5"),
	laSourceItemHIBRef(jT3, regA11, 0, 0, 0, 0, "data4"),
	laSourceItemHIBRef(jT3, regA12, 0, 0, 0, 0, "data5"),
	iarSourceItem(0),
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
	checkRegister(t, engine, common.A0, 0_556677_001122)
	checkRegister(t, engine, common.A1, 0_506070)
	checkRegister(t, engine, common.A2, 0_507777)
	checkRegister(t, engine, common.A3, 0_777777_506070)
	checkRegister(t, engine, common.A4, 0_223344)
	checkRegister(t, engine, common.A5, 0_777777_507777)
	checkRegister(t, engine, common.A6, 0_221100)
	checkRegister(t, engine, common.A7, 01111)
	checkRegister(t, engine, common.A8, 0_777777_775500)
	checkRegister(t, engine, common.A9, 02222)
	checkRegister(t, engine, common.A10, 0_777777_776600)
	checkRegister(t, engine, common.A11, 03333)
	checkRegister(t, engine, common.A12, 0_777777_777700)
}

var PartialWordLoadsExtendedQuarterWord = []*tasm.SourceItem{
	segSourceItem(077),
	labelDataSourceItem("data1", []uint64{0400, 0501, 0677, 0777}),
	labelDataSourceItem("data2", []uint64{012, 034, 056, 075, 042, 010}),

	segSourceItem(0),
	lrSourceItemHIBRef(jQ1, regR0, 0, 0, 0, 0, "data1"),
	lrSourceItemHIBRef(jQ2, regR1, 0, 0, 0, 0, "data1"),
	lrSourceItemHIBRef(jQ3, regR2, 0, 0, 0, 0, "data1"),
	lrSourceItemHIBRef(jQ4, regR3, 0, 0, 0, 0, "data1"),
	lrSourceItemHIBRef(jS1, regR4, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBRef(jS2, regR5, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBRef(jS3, regR6, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBRef(jS4, regR7, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBRef(jS5, regR8, 0, 0, 0, 0, "data2"),
	lrSourceItemHIBRef(jS6, regR9, 0, 0, 0, 0, "data2"),
	iarSourceItem(0),
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
	checkRegister(t, engine, common.R0, 0400)
	checkRegister(t, engine, common.R1, 0501)
	checkRegister(t, engine, common.R2, 0677)
	checkRegister(t, engine, common.R3, 0777)
	checkRegister(t, engine, common.R4, 012)
	checkRegister(t, engine, common.R5, 034)
	checkRegister(t, engine, common.R6, 056)
	checkRegister(t, engine, common.R7, 075)
	checkRegister(t, engine, common.R8, 042)
	checkRegister(t, engine, common.R9, 010)
}

var partialWordStoresExtendedThirdWord = []*tasm.SourceItem{
	segSourceItem(2),
	labelDataSourceItem("data1", []uint64{0_444444_444444}),
	labelSourceItem("data2"),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),

	segSourceItem(0),
	laSourceItemHIBRef(jW, regA0, 0, 0, 0, common.B2, "data1"),
	lxiSourceItemU(jU, regX1, 0, 1),
	lxmSourceItemU(jU, regX1, 0, 0),
	saSourceItemHIBRef(jW, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jH1, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jH2, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jXH1, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jXH2, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jT1, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jT2, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jT3, regA0, regX1, 1, 0, common.B2, "data2"),
	iarSourceItem(0),
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
	labelDataSourceItem("data1", []uint64{0_444444_444444}),
	labelSourceItem("data2"),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),
	dataSourceItem([]uint64{0_333333_333333}),

	segSourceItem(0),
	laSourceItemHIBRef(jW, regA0, 0, 0, 0, common.B2, "data1"),
	lxiSourceItemU(jU, regX1, 0, 1),
	lxmSourceItemU(jU, regX1, 0, 0),
	laSourceItemHIBRef(jW, regA0, 0, 0, 0, common.B2, "data1"),
	saSourceItemHIBRef(jQ1, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jQ2, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jQ3, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jQ4, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jS1, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jS2, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jS3, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jS4, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jS5, regA0, regX1, 1, 0, common.B2, "data2"),
	saSourceItemHIBRef(jS6, regA0, regX1, 1, 0, common.B2, "data2"),
	iarSourceItem(0),
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
	segSourceItem(0),
	lxSourceItemU(jU, regX5, 0, 42),
	laSourceItemU(jW, regA5, 0, common.X5),
	laSourceItemU(jS3, regA6, 0, common.X5), // partial word ignored, we're register-to-register
	iarSourceItem(0),
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
	checkRegister(t, engine, common.A5, 42)
	checkRegister(t, engine, common.A6, 42)
}

//	TODO extended mode index register handling

//	TODO extended mode addressing across multiple banks
