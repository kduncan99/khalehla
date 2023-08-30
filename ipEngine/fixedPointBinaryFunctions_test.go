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

const fAA = 014
const fANA = 015
const fAMA = 016
const fANMA = 017
const fAU = 020
const fANU = 021
const fAX = 024
const fANX = 025
const fMI = 030
const fMSI = 031
const fMF = 032
const fDI = 034
const fDSF = 035
const fDF = 036
const fDA = 071
const fDAN = 071
const fAH = 072
const fANH = 072
const fAT = 072
const fANT = 072
const fADD1 = 005
const fSUB1 = 005
const fINC = 005
const fINC2 = 005
const fDEC = 005
const fDEC2 = 005
const fENZ = 005
const BAO = 072

const jDA = 010
const jDAN = 011
const jAH = 004
const jANH = 005
const jAT = 006
const jANT = 007
const jBAO = 013

const aADD1 = 015
const aSUB1 = 016
const aINC = 010
const aINC2 = 012
const aDEC = 011
const aDEC2 = 013
const aENZ = 014

// ---------------------------------------------------
// AA

func aaSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fAA, j, a, x, u)
}

func aaSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fAA, j, a, x, h, i, ref)
}

func aaSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fAA, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// ANA

func anaSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fANA, j, a, x, u)
}

func anaSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fANA, j, a, x, h, i, ref)
}

func anaSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fANA, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// AMA

func amaSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fAMA, j, a, x, u)
}

func amaSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fAMA, j, a, x, h, i, ref)
}

func amaSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fAMA, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// ANMA

func anmaSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fANMA, j, a, x, u)
}

func anmaSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fANMA, j, a, x, h, i, ref)
}

func anmaSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fANMA, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// AU

func auSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fAU, j, a, x, u)
}

func auSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fAU, j, a, x, h, i, ref)
}

func auSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fAU, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// ANU

func anuSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fANU, j, a, x, u)
}

func anuSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fANU, j, a, x, h, i, ref)
}

func anuSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fANU, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// AX

func axSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fAX, j, a, x, u)
}

func axSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fAX, j, a, x, h, i, ref)
}

func axSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fAX, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// ANX

func anxSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fANX, j, a, x, u)
}

func anxSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fANX, j, a, x, h, i, ref)
}

func anxSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fANX, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// MI

func miSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fMI, j, a, x, u)
}

func miaSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fMI, j, a, x, h, i, ref)
}

func miSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fMI, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// MSI

func msiSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fMSI, j, a, x, u)
}

func msiSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fMSI, j, a, x, h, i, ref)
}

func msiSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fMSI, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// MF

func mfSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fMF, j, a, x, u)
}

func mfSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fMF, j, a, x, h, i, ref)
}

func mfSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fMF, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// DI

func diSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fDI, j, a, x, u)
}

func diSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fDI, j, a, x, h, i, ref)
}

func diSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fDI, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// DSF

func dsfSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fDSF, j, a, x, u)
}

func dsfSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fDSF, j, a, x, h, i, ref)
}

func dsfSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fDSF, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// DF

func dfSourceItemU(j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fDF, j, a, x, u)
}

func dfSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fDF, j, a, x, h, i, ref)
}

func dfSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fDF, j, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// DA

func daSourceItemU(a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fDA, jDA, a, x, u)
}

func daSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fDA, jDA, a, x, h, i, ref)
}

func daSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fDA, jDA, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// DAN

func danSourceItemU(a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fDAN, jDAN, a, x, u)
}

func danSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fDAN, jDAN, a, x, h, i, ref)
}

func danSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fDAN, jDAN, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// AH

func ahSourceItemU(a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fAH, jAH, a, x, u)
}

func ahSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fAH, jAH, a, x, h, i, ref)
}

func ahSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fAH, jAH, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// ANH

func anhSourceItemU(a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fANH, jANH, a, x, u)
}

func anhSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fANH, jANH, a, x, h, i, ref)
}

func anhSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fANH, jANH, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// AT

func atSourceItemU(a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fAT, jAT, a, x, u)
}

func atSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fAT, jAT, a, x, h, i, ref)
}

func atSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fAT, jAT, a, x, h, i, b, ref)
}

// ---------------------------------------------------
// ANT

func antSourceItemU(a uint64, x uint64, u int) *tasm.SourceItem {
	return fjaxuSourceItem(fANT, jANT, a, x, u)
}

func antSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fANT, jANT, a, x, h, i, ref)
}

func antSourceItemHIBRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fANT, jANT, a, x, h, i, b, ref)
}

//	TODO ADD1
//	TODO SUB1
//	TODO INC
//	TODO INC2
//	TODO DEC
//	TODO DEC2
//	TODO ENZ
//	TODO BAO

// ---------------------------------------------------------------------------------------------------------------------

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
	laSourceItemHIRef(jW, 2, 0, 0, 0, "a2data"),
	aaSourceItemHIRef(jW, 2, 0, 0, 0, "data"),
	iarSourceItem(0),

	segSourceItem(077),
	labelDataSourceItem("a2data", []uint64{0_700000_001314}),
	labelDataSourceItem("data", []uint64{0_000212_273555}),
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
	checkRegister(t, engine, pkg.A2, 0_700212_275071)
}

var anaCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIRef(jW, 2, 0, 0, 0, "a2data"),
	anaSourceItemHIRef(jW, 2, 0, 0, 0, "data"),
	iarSourceItem(0),

	segSourceItem(077),
	labelDataSourceItem("a2data", []uint64{0344_072777}),
	labelDataSourceItem("data", []uint64{02227}),
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
	checkRegister(t, engine, pkg.A2, 0_000344_070550)
}

var amaCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIRef(jW, 2, 0, 0, 0, "a2data"),
	amaSourceItemHIRef(jW, 2, 0, 0, 0, "data"),
	iarSourceItem(0),

	segSourceItem(077),
	labelDataSourceItem("a2data", []uint64{0_000427_031272}),
	labelDataSourceItem("data", []uint64{0_703247_006666}),
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
	checkRegister(t, engine, pkg.A2, 0_075160_022403)
}

var anmaCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIRef(jW, 2, 0, 0, 0, "a2data"),
	anmaSourceItemHIRef(jW, 2, 0, 0, 0, "data"),
	iarSourceItem(0),

	segSourceItem(077),
	labelDataSourceItem("a2data", []uint64{0_300004_000000}),
	labelDataSourceItem("data", []uint64{0_032222_123223}),
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
	checkRegister(t, engine, pkg.A2, 0_245561_654555)
}

var anuCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIRef(jW, 2, 0, 0, 0, "a2data"),
	anuSourceItemHIRef(jW, 2, 0, 0, 0, "data"),
	iarSourceItem(0),

	segSourceItem(077),
	labelSourceItem("a2data"),
	dataSourceItem([]uint64{0_000000_372117}),
	dataSourceItem([]uint64{0_400377_777777}),
	labelDataSourceItem("data", []uint64{0_500374_120000}),
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
	checkRegister(t, engine, pkg.A2, 0_000000_372117)
	checkRegister(t, engine, pkg.A3, 0_277404_252116)
}

var axCode = []*tasm.SourceItem{
	segSourceItem(0),
	lxSourceItemHIRef(jW, 2, 0, 0, 0, "a2data"),
	axSourceItemHIRef(jW, 2, 0, 0, 0, "data"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.X2, 0_700212_275071)
}

var anxCode = []*tasm.SourceItem{
	segSourceItem(0),
	lxSourceItemHIRef(jW, 2, 0, 0, 0, "a2data"),
	anxSourceItemHIRef(jW, 2, 0, 0, 0, "data"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.X2, 0_000344_070550)
}

var miCode = []*tasm.SourceItem{
	segSourceItem(0),
	dlSourceItemHIBRef(2, 0, 0, 0, 2, "a2data"),
	miSourceItemHIBRef(jW, 2, 0, 0, 0, 2, "data"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.A2, 0_000000_000000)
	checkRegister(t, engine, pkg.A3, 0_115624_561516)
}

var msiCodeOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 2, 0, 0, 0, 2, "a2data"),
	msiSourceItemHIBRef(jW, 2, 0, 0, 0, 2, "data"),
	iarSourceItem(1),

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
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "a4data"),
	msiSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "data"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.A4, 0_000000_015012)
}

var mfCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 3, 0, 0, 0, 2, "a3data"),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "a4data"),
	mfSourceItemHIBRef(jW, 3, 0, 0, 0, 2, "data"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.A3, 0_044444_444445)
	checkRegister(t, engine, pkg.A4, 0_044444_444444)
}

var diCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 2, 0, 0, 0, 2, "a2data"),
	laSourceItemHIBRef(jW, 3, 0, 0, 0, 2, "a3data"),
	diSourceItemHIBRef(jW, 2, 0, 0, 0, 2, "data"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.A2, 0_005213_747442)
	checkRegister(t, engine, pkg.A3, 0_000000_244613)
}

var diCodeDivideCheck = []*tasm.SourceItem{
	segSourceItem(0),
	dlSourceItemHIBRef(2, 0, 0, 0, 2, "a2data"),
	laSourceItemU(jU, 10, 0, 0),
	diSourceItemHIBRef(jW, 2, 0, 0, 0, 0, grsRef(pkg.A10)),
	iarSourceItem(1),

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
	checkRegister(t, engine, pkg.A2, 0_000001_612175)
	checkRegister(t, engine, pkg.A3, 0_437777_700000)
}

var dsfCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 3, 0, 0, 0, 2, "a3data"),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "a4data"),
	dsfSourceItemHIBRef(jW, 3, 0, 0, 0, 2, "data"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.A3, 0_000000_007236)
	checkRegister(t, engine, pkg.A4, 0_001733_765274)
}

var dfCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "a4data"),
	laSourceItemHIBRef(jW, 5, 0, 0, 0, 2, "a5data"),
	dfSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "data"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.A4, 0_000000_021653)
	checkRegister(t, engine, pkg.A5, 0_000000_000056)
}

var daCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "a4data"),
	laSourceItemHIBRef(jW, 5, 0, 0, 0, 2, "a5data"),
	daSourceItemHIBRef(regA4, 0, 0, 0, 2, "data1"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.A4, 0_123012_342553)
	checkRegister(t, engine, pkg.A5, 0_056323_321126)
}

var danCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "a4data"),
	laSourceItemHIBRef(jW, 5, 0, 0, 0, 2, "a5data"),
	danSourceItemHIBRef(regA4, 0, 0, 0, 2, "data1"),
	iarSourceItem(0),

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
	checkRegister(t, engine, pkg.A4, 0_000000_113110)
	checkRegister(t, engine, pkg.A5, 0_210014_413002)
}

var ahCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, regA5, 0, 0, 0, 2, "a5data"),
	ahSourceItemHIBRef(regA5, 0, 0, 0, 2, "data"),
	iarSourceItem(0),

	segSourceItem(2),
	labelDataSourceItem("a5data", []uint64{0_000123_555123}),
	labelDataSourceItem("data", []uint64{0_000001_223000}),
}

func Test_AH_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", ahCode)
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
	checkRegister(t, engine, pkg.A5, 0_000124_000124)
}

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
