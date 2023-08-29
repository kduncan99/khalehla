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

// The following instructions update DB18 (carry) and DB19 (overflow) in the following conditions:
//
//	Input Signs   Output Sign   DB18 DB19
//	    +/+             +         0    0
//	    +/+             -         0    1
//	    +/-             +         1    0
//	    +/-             -         0    0
//	    -/-             +         1    1
//	    -/-             -         1    0
var tcoData = [][]bool{
	//	add1Pos, add2Pos, sumPos, DB18, DB19
	{true, true, true, false, false},
	{true, true, false, false, true},
	{true, false, true, true, false},
	{true, false, false, false, false},
	{false, false, true, true, true},
	{false, false, false, true, false},
}

func Test_Carry_Overflow(t *testing.T) {
	engine := NewEngine("TEST", nil, nil)
	dr := engine.GetDesignatorRegister()
	for opTrap := 0; opTrap < 2; opTrap++ {
		for _, tcoEntry := range tcoData {
			engine.ClearAllInterrupts()
			dr.Clear()
			dr.SetOperationTrapEnabled(opTrap == 1)

			add1Pos := tcoEntry[0]
			add2Pos := tcoEntry[1]
			sumPos := tcoEntry[2]
			updateDesignatorRegister(engine, add1Pos, add2Pos, sumPos)

			prefix := fmt.Sprintf("With addend1Positive = %v and addend2Positive = %v and SumPositive = %v:", add1Pos, add2Pos, sumPos)
			if dr.IsCarrySet() != tcoEntry[3] {
				t.Errorf("%s Carry Flag was incorrectly = %v", prefix, dr.IsCarrySet())
			}
			if dr.IsOverflowSet() != tcoEntry[4] {
				t.Errorf("%s Overflow Flag was incorrectly = %v", prefix, dr.IsOverflowSet())
			}

			prefix = fmt.Sprintf("With Operation Trap Enabled = %v and overflow = %v:", opTrap, dr.IsOverflowSet())
			interruptWanted := dr.IsOperationTrapEnabled() && dr.IsOverflowSet()
			if interruptWanted != engine.HasPendingInterrupt() {
				if interruptWanted {
					t.Errorf("%s An interrupt was expected but not posted", prefix)
				} else {
					t.Errorf("%s An interrupt was not expected but was posted", prefix)
				}
			}
		}
	}
}

var aaCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIURef("", jW, 2, 0, 0, 0, "a2data"),
	aaSourceItemHIURef("", jW, 2, 0, 0, 0, "data"),
	iarSourceItem("", 0),

	segSourceItem(077),
	tasm.NewSourceItem("a2data", "w", []string{"0700000001314"}),
	tasm.NewSourceItem("data", "w", []string{"000212273555"}),
}

func Test_AA_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", aaCode)
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
	checkRegister(t, engine, pkg.A2, 0_700212_275071, "A2")
}

var anaCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIURef("", jW, 2, 0, 0, 0, "a2data"),
	anaSourceItemHIURef("", jW, 2, 0, 0, 0, "data"),
	iarSourceItem("", 0),

	segSourceItem(077),
	tasm.NewSourceItem("a2data", "w", []string{"0344072777"}),
	tasm.NewSourceItem("data", "w", []string{"02227"}),
}

func Test_ANA_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", anaCode)
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
	checkRegister(t, engine, pkg.A2, 0_000344_070550, "A2")
}

var amaCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIURef("", jW, 2, 0, 0, 0, "a2data"),
	amaSourceItemHIURef("", jW, 2, 0, 0, 0, "data"),
	iarSourceItem("", 0),

	segSourceItem(077),
	tasm.NewSourceItem("a2data", "w", []string{"0427031272"}),
	tasm.NewSourceItem("data", "w", []string{"0703247006666"}),
}

func Test_AMA_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", amaCode)
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
	checkRegister(t, engine, pkg.A2, 0_075160_022403, "A2")
}

var anmaCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIURef("", jW, 2, 0, 0, 0, "a2data"),
	anmaSourceItemHIURef("", jW, 2, 0, 0, 0, "data"),
	iarSourceItem("", 0),

	segSourceItem(077),
	tasm.NewSourceItem("a2data", "w", []string{"0300004000000"}),
	tasm.NewSourceItem("data", "w", []string{"032222123223"}),
}

func Test_ANMA_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", anmaCode)
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
	checkRegister(t, engine, pkg.A2, 0_245561_654555, "A2")
}

var anuCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIURef("", jW, 2, 0, 0, 0, "a2data"),
	anuSourceItemHIURef("", jW, 2, 0, 0, 0, "data"),
	iarSourceItem("", 0),

	segSourceItem(077),
	tasm.NewSourceItem("a2data", "w", []string{"0_000000_372117"}),
	tasm.NewSourceItem("a3data", "w", []string{"0_400377_777777"}),
	tasm.NewSourceItem("data", "w", []string{"0_500374_120000"}),
}

func Test_ANU_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", anuCode)
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
	checkRegister(t, engine, pkg.A2, 0_000000_372117, "A2")
	checkRegister(t, engine, pkg.A3, 0_277404_252116, "A3")
}

var axCode = []*tasm.SourceItem{
	segSourceItem(0),
	lxSourceItemHIURef("", jW, 2, 0, 0, 0, "a2data"),
	axSourceItemHIURef("", jW, 2, 0, 0, 0, "data"),
	iarSourceItem("", 0),

	segSourceItem(077),
	tasm.NewSourceItem("a2data", "w", []string{"0700000001314"}),
	tasm.NewSourceItem("data", "w", []string{"000212273555"}),
}

func Test_AX_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", axCode)
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
	checkRegister(t, engine, pkg.X2, 0_700212_275071, "X2")
}

var anxCode = []*tasm.SourceItem{
	segSourceItem(0),
	lxSourceItemHIURef("", jW, 2, 0, 0, 0, "a2data"),
	anxSourceItemHIURef("", jW, 2, 0, 0, 0, "data"),
	iarSourceItem("", 0),

	segSourceItem(077),
	tasm.NewSourceItem("a2data", "w", []string{"0344072777"}),
	tasm.NewSourceItem("data", "w", []string{"02227"}),
}

func Test_ANX_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", anxCode)
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
	checkRegister(t, engine, pkg.X2, 0_000344_070550, "X2")
}

var miCode = []*tasm.SourceItem{
	segSourceItem(0),
	dlSourceItemHIBDRef("", 2, 0, 0, 0, 2, "a2data"),
	miSourceItemHIBDRef("", jW, 2, 0, 0, 0, 2, "data"),
	iarSourceItem("", 0),

	segSourceItem(2),
	tasm.NewSourceItem("a2data", "w", []string{"0_000001_612175"}),
	tasm.NewSourceItem("a3data", "w", []string{"0_437777_700000"}),
	tasm.NewSourceItem("data", "w", []string{"0_000000053746"}),
}

func Test_MI_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", miCode)
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
	checkRegister(t, engine, pkg.A2, 0_000000_000000, "A2")
	checkRegister(t, engine, pkg.A3, 0_115624_561516, "A3")
}

var msiCodeOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 2, 0, 0, 0, 2, "a2data"),
	msiSourceItemHIBDRef("", jW, 2, 0, 0, 0, 2, "data"),
	iarSourceItem("", 1),

	segSourceItem(2),
	tasm.NewSourceItem("a2data", "w", []string{"0_377777_777777"}),
	tasm.NewSourceItem("data", "w", []string{"02"}),
}

func Test_MSI_Overflow(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", msiCodeOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkBankPerSegment(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetOperationTrapEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkInterruptAndSSF(t, engine, pkg.OperationTrapInterruptClass, pkg.OperationTrapMultiplySingleIntegerOverflow)
}

var msiCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "a4data"),
	msiSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data"),
	iarSourceItem("", 0),

	segSourceItem(2),
	tasm.NewSourceItem("a4data", "w", []string{"0_000000_000312"}),
	tasm.NewSourceItem("data", "w", []string{"0_000000_000041"}),
}

func Test_MSI_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", msiCode)
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
	checkRegister(t, engine, pkg.A4, 0_000000_015012, "A4")
}

var mfCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 3, 0, 0, 0, 2, "a3data"),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "a4data"),
	mfSourceItemHIBDRef("", jW, 3, 0, 0, 0, 2, "data"),
	iarSourceItem("", 0),

	segSourceItem(2),
	tasm.NewSourceItem("a3data", "w", []string{"0_200000_000002"}),
	tasm.NewSourceItem("a4data", "w", []string{"0_777777_777777"}),
	tasm.NewSourceItem("data", "w", []string{"0_111111_111111"}),
}

func Test_MF_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", mfCode)
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
	checkRegister(t, engine, pkg.A3, 0_044444_444445, "A3")
	checkRegister(t, engine, pkg.A4, 0_044444_444444, "A4")
}

var diCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 2, 0, 0, 0, 2, "a2data"),
	laSourceItemHIBDRef("", jW, 3, 0, 0, 0, 2, "a3data"),
	diSourceItemHIBDRef("", jW, 2, 0, 0, 0, 2, "data"),
	iarSourceItem("", 0),

	segSourceItem(2),
	tasm.NewSourceItem("a2data", "w", []string{"0_000000_011416"}),
	tasm.NewSourceItem("a3data", "w", []string{"0_110621_672145"}),
	tasm.NewSourceItem("data", "w", []string{"0_000001_635035"}),
}

func Test_DI_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", diCode)
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
	checkRegister(t, engine, pkg.A2, 0_005213_747442, "A2")
	checkRegister(t, engine, pkg.A3, 0_000000_244613, "A3")
}

var diCodeDivideCheck = []*tasm.SourceItem{
	segSourceItem(0),
	dlSourceItemHIBDRef("", 2, 0, 0, 0, 2, "a2data"),
	laSourceItemU("", jU, 10, 0, 0),
	diSourceItemHIBD("", jW, 2, 0, 0, 0, 0, grsA10),
	iarSourceItem("", 1),

	segSourceItem(2),
	tasm.NewSourceItem("a2data", "w", []string{"0_000001_612175"}),
	tasm.NewSourceItem("a3data", "w", []string{"0_437777_700000"}),
	tasm.NewSourceItem("data", "w", []string{"0_000000053746"}),
}

func Test_DI_DivideCheck(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", diCodeDivideCheck)
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
	checkInterruptAndSSF(t, engine, pkg.ArithmeticExceptionInterruptClass, pkg.ArithmeticExceptionDivideCheck)
	checkRegister(t, engine, pkg.A2, 0_000001_612175, "A2")
	checkRegister(t, engine, pkg.A3, 0_437777_700000, "A3")
}

var dsfCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 3, 0, 0, 0, 2, "a3data"),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "a4data"),
	dsfSourceItemHIBDRef("", jW, 3, 0, 0, 0, 2, "data"),
	iarSourceItem("", 0),

	segSourceItem(2),
	tasm.NewSourceItem("a3data", "w", []string{"0_000000_007236"}),
	tasm.NewSourceItem("a4data", "w", []string{"0_743464_241454"}),
	tasm.NewSourceItem("data", "w", []string{"0_000001_711467"}),
}

func Test_DSF_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", dsfCode)
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
	checkRegister(t, engine, pkg.A3, 0_000000_007236, "A3")
	checkRegister(t, engine, pkg.A4, 0_001733_765274, "A4")
}

var dfCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "a4data"),
	laSourceItemHIBDRef("", jW, 5, 0, 0, 0, 2, "a5data"),
	dfSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data"),
	iarSourceItem("", 0),

	segSourceItem(2),
	tasm.NewSourceItem("a4data", "w", []string{"0_000000_000000"}),
	tasm.NewSourceItem("a5data", "w", []string{"0_000061_026335"}),
	tasm.NewSourceItem("data", "w", []string{"0_000000_001300"}),
}

func Test_DF_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", dfCode)
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
	checkRegister(t, engine, pkg.A4, 0_000000_021653, "A4")
	checkRegister(t, engine, pkg.A5, 0_000000_000056, "A5")
}

var daCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "a4data"),
	laSourceItemHIBDRef("", jW, 5, 0, 0, 0, 2, "a5data"),
	daSourceItemHIBDRef("", 4, 0, 0, 0, 2, "data1"),
	iarSourceItem("", 0),

	segSourceItem(2),
	tasm.NewSourceItem("a4data", "w", []string{"0_123001_230121"}),
	tasm.NewSourceItem("a5data", "w", []string{"0_400002_321021"}),
	tasm.NewSourceItem("data1", "w", []string{"0_000011_112431"}),
	tasm.NewSourceItem("data2", "w", []string{"0_456321_000105"}),
}

func Test_DA_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", daCode)
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
	checkRegister(t, engine, pkg.A4, 0_123012_342553, "A4")
	checkRegister(t, engine, pkg.A5, 0_056323_321126, "A5")
}

var danCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "a4data"),
	laSourceItemHIBDRef("", jW, 5, 0, 0, 0, 2, "a5data"),
	danSourceItemHIBDRef("", 4, 0, 0, 0, 2, "data1"),
	iarSourceItem("", 0),

	segSourceItem(2),
	tasm.NewSourceItem("a4data", "w", []string{"0_000000_543210"}),
	tasm.NewSourceItem("a5data", "w", []string{"0_210056_523004"}),
	tasm.NewSourceItem("data1", "w", []string{"0_000000_430100"}),
	tasm.NewSourceItem("data2", "w", []string{"0_000042_110002"}),
}

func Test_DAN_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", danCode)
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
	checkRegister(t, engine, pkg.A4, 0_000000_113110, "A4")
	checkRegister(t, engine, pkg.A5, 0_210014_413002, "A5")
}

//	TODO AH
//	TODO ANH
//	TODO AT
//	TODO ANT
//	TODO ADD1
//	TODO SUB1
//	TODO INC
//	TODO INC2
//	TODO DEC
//	TODO DEC2
//	TODO ENZ
//	TODO BAO
