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

const fTEP = 044
const fTOP = 045
const fTLEM = 047
const fTNOP = 050
const fTGZ = 050
const fTPZ = 050
const fTPBasic = 060
const fTPExtended = 050
const fTMZ = 050
const fTMZG = 050
const fTZ = 050
const fTNLZ = 050
const fTLZ = 050
const fTNZBasic = 051
const fTNZExtended = 050
const fTPZL = 050
const fTNMZ = 050
const fTNBasic = 061
const fTNExtended = 050
const fTNPZ = 050
const fTNGZ = 050
const fTSKP = 050
const fTE = 052
const fDTE = 071
const fTNE = 053
const fTLE = 054
const fTG = 055
const fTGM = 033
const fDTGM = 033
const fTW = 056
const fTNW = 057
const fMTE = 071
const fMTNE = 071
const fMTLE = 071
const fMTG = 071
const fMTW = 071
const fMTNW = 071
const fMATL = 071
const fMATG = 071
const fTS = 073
const fTSS = 073
const fTCS = 073
const fCR = 075
const fUNLK = 073

const jDTE = 017
const jTGM = 013
const jDTGM = 014
const jMTE = 000
const jMTNE = 001
const jMTLE = 002
const jMTG = 003
const jMTW = 004
const jMTNW = 005
const jMATL = 006
const jMATG = 007
const jTS = 017
const jTSS = 017
const jTCS = 017
const jCR = 015
const jUNLK = 014

const aTNOP = 000
const aTGZ = 001
const aTPZ = 002
const aTPBasic = 000
const aTPExtended = 003
const aTMZ = 004
const aTMZG = 005
const aTZBasic = 000
const aTZExtended = 006
const aTNLZ = 007
const aTLZ = 010
const aTNZBasic = 000
const aTNZExtended = 011
const aTPZL = 012
const aTNMZ = 013
const aTNBasic = 000
const aTNExtended = 014
const aTNPZ = 015
const aTNGZ = 016
const aTSKP = 017
const aTS = 000
const aTSS = 001
const aTCS = 002
const aUNLK = 004

// ---------------------------------------------------
// TEP

func tepSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fTEP, j, a, x, h, i, ref)
}

func tepSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTEP, j, a, x, h, i, b, ref)
}

func tepSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTEP, j, a, x, u)
}

// ---------------------------------------------------
// TOP

func topSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fTOP, j, a, x, h, i, ref)
}

func topSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTOP, j, a, x, h, i, b, ref)
}

func topSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTOP, j, a, x, u)
}

// ---------------------------------------------------
// TLEM

func tlemSourceItemHIRef(j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fTLEM, j, a, x, h, i, ref)
}

func tlemSourceItemHIBRef(j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTLEM, j, a, x, h, i, b, ref)
}

func tlemSourceItemU(j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTLEM, j, a, x, u)
}

// ---------------------------------------------------
// TNOP - extended mode only

func tnopSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTNOP, j, aTNOP, x, u)
}

func tnopSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTNOP, j, aTNOP, x, h, i, b, ref)
}

// ---------------------------------------------------
// TGZ - extended mode only

func tgzSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTGZ, j, aTGZ, x, u)
}

func tgzSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTGZ, j, aTGZ, x, h, i, b, ref)
}

// ---------------------------------------------------
// TPZ - extended mode only

func tpzSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTPZ, j, aTPZ, x, u)
}

func tpzSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTPZ, j, aTPZ, x, h, i, b, ref)
}

// ---------------------------------------------------
// TP

func tpSourceItemUBasic(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTPBasic, j, aTPBasic, x, u)
}

func tpSourceItemUExtended(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTPExtended, j, aTPExtended, x, u)
}

func tpSourceItemHIRef(j uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fTPBasic, j, aTPBasic, x, h, i, ref)
}

func tpSourceItemHIBDRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTPExtended, j, aTPExtended, x, h, i, b, ref)
}

// ---------------------------------------------------
// TMZ - extended mode only

func tmzSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTMZ, j, aTMZ, x, u)
}

func tmzSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTMZ, j, aTMZ, x, h, i, b, ref)
}

// ---------------------------------------------------
// TMZG - extended mode only

func tmzgSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTMZG, j, aTMZG, x, u)
}

func tmzgSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTMZG, j, aTMZG, x, h, i, b, ref)
}

// ---------------------------------------------------
// TZ

func tzSourceItemUBasic(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTZ, j, aTZBasic, x, u)
}

func tzSourceItemUExtended(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTZ, j, aTZExtended, x, u)
}

func tzSourceItemHIRef(j uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fTZ, j, aTZBasic, x, h, i, ref)
}

func tzSourceItemHIBDRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTZ, j, aTZExtended, x, h, i, b, ref)
}

// ---------------------------------------------------
// TNLZ - extended mode only

func tnlzSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTNLZ, j, aTNLZ, x, u)
}

func tnlzSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTNLZ, j, aTNLZ, x, h, i, b, ref)
}

// ---------------------------------------------------
// TLZ - extended mode only

func tlzSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTLZ, j, aTLZ, x, u)
}

func tlzSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTLZ, j, aTLZ, x, h, i, b, ref)
}

// ---------------------------------------------------
// TNZ

func tnzSourceItemUBasic(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTNZBasic, j, aTNZBasic, x, u)
}

func tnzSourceItemUExtended(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTNZExtended, j, aTNZExtended, x, u)
}

func tnzSourceItemHIRef(j uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fTNZBasic, j, aTNZBasic, x, h, i, ref)
}

func tnzSourceItemHIBDRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTNZExtended, j, aTNZExtended, x, h, i, b, ref)
}

// ---------------------------------------------------
// TPZL - extended mode only

func tpzlSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTPZL, j, aTPZL, x, u)
}

func tpzlSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTPZL, j, aTPZL, x, h, i, b, ref)
}

// ---------------------------------------------------
// TNMZ - extended mode only

func tnmzSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTNMZ, j, aTNMZ, x, u)
}

func tnmzSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTNMZ, j, aTNMZ, x, h, i, b, ref)
}

// ---------------------------------------------------
// TN

func tnSourceItemUBasic(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTNBasic, j, aTNBasic, x, u)
}

func tnSourceItemUExtended(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTNExtended, j, aTNExtended, x, u)
}

func tnSourceItemHIRef(j uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fTNBasic, j, aTNBasic, x, h, i, ref)
}

func tnSourceItemHIBDRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTNExtended, j, aTNExtended, x, h, i, b, ref)
}

// ---------------------------------------------------
// TNPZ - extended mode only

func tnpzSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTNPZ, j, aTNPZ, x, u)
}

func tnpzSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTNPZ, j, aTNPZ, x, h, i, b, ref)
}

// ---------------------------------------------------
// TNGZ - extended mode only

func tngzSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTNGZ, j, aTNGZ, x, u)
}

func tngzSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTNGZ, j, aTNGZ, x, h, i, b, ref)
}

// ---------------------------------------------------
// TSKP - extended mode only

func tskpSourceItemU(j uint64, x uint64, u uint64) *tasm.SourceItem {
	return fjaxuSourceItem(fTSKP, j, aTSKP, x, u)
}

func tskpSourceItemHIBRef(j uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhibRefSourceItem(fTSKP, j, aTSKP, x, h, i, b, ref)
}

// ---------------------------------------------------------------------------------------------------------------------

var tepCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "data1"),
	tepSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "data2"),
	jSourceItemExtended("", 0, 0, 0, "tag"),
	iarSourceItem(1),

	labelSourceItem("tag"),
	tepSourceItemHIBRef(jQ4, 4, 0, 0, 0, 2, "data2"),
	iarSourceItem(2),
	iarSourceItem(0),

	segSourceItem(2),
	tasm.NewSourceItem("data1", "w", []string{"0123456543210"}),
	tasm.NewSourceItem("data2", "w", []string{"0777777"}),
}

func Test_TEP(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", tepCode)
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
}

var topCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "data1"),
	topSourceItemHIBRef(jW, 4, 0, 0, 0, 2, "data2"),
	jSourceItemExtended("", 0, 0, 0, "tag"),
	iarSourceItem(1),
	labelSourceItem("tag"),
	topSourceItemHIBRef(jQ4, 4, 0, 0, 0, 2, "data2"),
	iarSourceItem(2),
	iarSourceItem(0),

	segSourceItem(2),
	tasm.NewSourceItem("data1", "w", []string{"0123456543211"}),
	tasm.NewSourceItem("data2", "w", []string{"0777777"}),
}

func Test_TOP(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", topCode)
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
}

var tlemCode = []*tasm.SourceItem{
	segSourceItem(0),
	lxSourceItemHIBRef(jW, 5, 0, 0, 0, 2, "x5content"),
	tlemSourceItemHIBRef(jW, 5, 0, 0, 0, 2, "arm"),
	jSourceItemExtended("", 0, 0, 0, "tag1"),
	iarSourceItem(1),

	labelSourceItem("tag1"),
	lxSourceItemHIBRef(jW, 6, 0, 0, 0, 0, fmt.Sprintf("%d", pkg.X5)),
	tlemSourceItemHIBRef(jS5, 6, 0, 0, 0, 2, "arm"),
	iarSourceItem(2),
	iarSourceItem(0),

	segSourceItem(2),
	tasm.NewSourceItem("arm", "w", []string{"0135471234"}),
	tasm.NewSourceItem("x5content", "w", []string{"02061234"}),
}

func Test_TLEM(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", tlemCode)
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
}

var tnopCode = []*tasm.SourceItem{
	segSourceItem(0),
	tnopSourceItemHIBRef(jW, 0, 0, 0, 2, "data1"),
	iarSourceItem(0),
	iarSourceItem(1),

	segSourceItem(2),
	tasm.NewSourceItem("data1", "w", []string{"0123456543210"}),
}

func Test_TNOP(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", tnopCode)
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
}

//	TODO TGZ
//	TODO TPZ
//	TODO TP
//	TODO TMZ
//	TODO TMZG
//	TODO TZ
//	TODO TNLZ
//	TODO TLZ
//	TODO TNZ
//	TODO TPZL
//	TODO TNMZ
//	TODO TN
//	TODO TNPZ
//	TODO TNGZ

var tskpCode = []*tasm.SourceItem{
	segSourceItem(0),
	tskpSourceItemHIBRef(jW, 0, 0, 0, 2, "data1"),
	iarSourceItem(1),
	iarSourceItem(0),

	segSourceItem(2),
	tasm.NewSourceItem("data1", "w", []string{"0123456543210"}),
}

func Test_TSKP(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", tskpCode)
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
}

//	TODO TE
//	TODO DTE
//	TODO TNE
//	TODO TLE
//	TODO TG
//	TODO TGM
//	TODO DTGM
//	TODO TW
//	TODO TNW
//	TODO MTE
//	TODO MTNE
//	TODO MTLE
//	TODO MTG
//	TODO MTW
//	TODO MTNW
//	TODO MATL
//	TODO MATG
//	TODO TS
//	TODO TSS
//	TODO TCS
//	TODO CR
//	TODO UNLK
