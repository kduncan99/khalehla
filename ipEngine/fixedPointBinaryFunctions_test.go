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
