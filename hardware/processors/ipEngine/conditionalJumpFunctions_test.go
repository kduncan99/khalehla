// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"testing"

	"khalehla/common"
	"khalehla/tasm"
)

const fJZ = 074
const fDJZ = 071
const fJNZ = 074
const fJP = 074
const fJPS = 072
const fJN = 074
const fJNS = 072
const fJB = 074
const fJNB = 074
const fJGD = 070
const fJMGI = 074
const fJO = 074
const fJNO = 074
const fJC = 074
const fJNC = 074
const fJDF = 074
const fJNDF = 074
const fJFO = 074
const fJNFO = 074
const fJFU = 074
const fJNFU = 074

const jJZ = 000
const jDJZ = 016
const jJNZ = 001
const jJP = 002
const jJPS = 002
const jJN = 003
const jJNS = 003
const jJB = 011
const jJNB = 010
const jJMGI = 012
const jJO = 014
const jJNO = 015
const jJCBasic = 016
const jJCExtended = 014
const jJNCBasic = 017
const jJNCExtended = 014
const jJDF = 014
const jJNDF = 015
const jJFO = 014
const jJNFO = 015
const jJFU = 014
const jJNFU = 015

const aJO = 000
const aJNO = 000
const aJCBasic = 000
const aJCExtended = 004
const aJNCBasic = 000
const aJNCExtended = 005
const aJDF = 003
const aJNDF = 003
const aJFO = 002
const aJNFO = 002
const aJFU = 001
const aJNFU = 001

// ---------------------------------------------------
// JZ

func jzSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJZ, jJZ, a, x, h, i, ref)
}

func jzSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJZ, jJZ, a, 0, ref)
}

// ---------------------------------------------------
// DJZ

func djzSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fDJZ, jDJZ, a, x, h, i, ref)
}

func djzSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fDJZ, jDJZ, a, 0, ref)
}

func djzSourceItemHIURef(label string, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fDJZ),
		fmt.Sprintf("%03o", jDJZ),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// ---------------------------------------------------
// JNZ

func jnzSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJNZ, jJNZ, a, x, h, i, ref)
}

func jnzSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJNZ, jJNZ, a, 0, ref)
}

// ---------------------------------------------------
// JP

func jpSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJP, jJP, a, x, h, i, ref)
}

func jpSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJP, jJP, a, 0, ref)
}

// ---------------------------------------------------
// JPS

func jpsSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJPS, jJPS, a, x, h, i, ref)
}

func jpsSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJPS, jJPS, a, 0, ref)
}

// ---------------------------------------------------
// JN

func jnSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJN, jJN, a, x, h, i, ref)
}

func jnSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJN, jJN, a, 0, ref)
}

// ---------------------------------------------------
// JNS

func jnsSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJNS, jJNS, a, x, h, i, ref)
}

func jnsSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJNS, jJNS, a, 0, ref)
}

// ---------------------------------------------------
// JB

func jbSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJB, jJB, a, x, h, i, ref)
}

func jbSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJB, jJB, a, 0, ref)
}

// ---------------------------------------------------
// JNB

func jnbSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJNB, jJNB, a, x, h, i, ref)
}

func jnbSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJNB, jJNB, a, 0, ref)
}

// ---------------------------------------------------
// JGD - this is atypical

func jgdSourceItemHIRef(grsIndex uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJGD, grsIndex>>4, grsIndex&017, x, h, i, ref)
}

func jgdSourceItemRef(grsIndex uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJGD, grsIndex>>4, grsIndex&017, 0, ref)
}

// ---------------------------------------------------
// JMGI

func jmgiSourceItemHIRef(a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJMGI, jJMGI, a, x, h, i, ref)
}

func jmgiSourceItemRef(a uint64, ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJMGI, jJMGI, a, 0, ref)
}

// ---------------------------------------------------
// JO

func joSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJO, jJO, aJO, x, h, i, ref)
}

func joSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJO, jJO, aJO, 0, ref)
}

// ---------------------------------------------------
// JNO

func jnoSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJNO, jJNO, aJNO, x, h, i, ref)
}

func jnoSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJNO, jJNO, aJNO, 0, ref)
}

// ---------------------------------------------------
// JC

func jcSourceItemHIRefBasic(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJC, jJCBasic, aJCBasic, x, h, i, ref)
}

func jcSourceItemRefBasic(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJC, jJCBasic, aJCBasic, 0, ref)
}

func jcSourceItemHIRefExtended(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJC, jJCExtended, aJCExtended, x, h, i, ref)
}

func jcSourceItemRefExtended(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJC, jJCExtended, aJCExtended, 0, ref)
}

// ---------------------------------------------------
// JNC

func jncSourceItemHIRefBasic(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJNC, jJNCBasic, aJNCBasic, x, h, i, ref)
}

func jncSourceItemRefBasic(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJNC, jJNCBasic, aJNCBasic, 0, ref)
}

func jncSourceItemHIRefExtended(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJNC, jJNCExtended, aJNCExtended, x, h, i, ref)
}

func jncSourceItemRefExtended(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJNC, jJNCExtended, aJNCExtended, 0, ref)
}

// ---------------------------------------------------
// JDF

func jdfSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJDF, jJDF, aJDF, x, h, i, ref)
}

func jdfSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJDF, jJDF, aJDF, 0, ref)
}

// ---------------------------------------------------
// JNDF

func jndfSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJNDF, jJNDF, aJNDF, x, h, i, ref)
}

func jndfSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJNDF, jJNDF, aJNDF, 0, ref)
}

// ---------------------------------------------------
// JFO

func jfoSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJFO, jJFO, aJFO, x, h, i, ref)
}

func jfoSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJFO, jJFO, aJFO, 0, ref)
}

// ---------------------------------------------------
// JNFO

func jnfoSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJNFO, jJNFO, aJNFO, x, h, i, ref)
}

func jnfoSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJNFO, jJNFO, aJNFO, 0, ref)
}

// ---------------------------------------------------
// JFU

func jfuSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJFU, jJFU, aJFU, x, h, i, ref)
}

func jfuSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJFU, jJFU, aJFU, 0, ref)
}

// ---------------------------------------------------
// JNFU

func jnfuSourceItemHIRef(x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	return fjaxhiRefSourceItem(fJNFU, jJNFU, aJNFU, x, h, i, ref)
}

func jnfuSourceItemRef(ref string) *tasm.SourceItem {
	return fjaxRefSourceItem(fJNFU, jJNFU, aJNFU, 0, ref)
}

// ---------------------------------------------------------------------------------------------------------------------

var jumpZeroExtendedPosZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA5, 0, 0),
	jzSourceItemRef(regA5, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JZ_Extended_PosZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpZeroExtendedPosZero)
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

var jumpZeroExtendedNegZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jXU, regA5, 0, 0_777777),
	jzSourceItemRef(regA5, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JZ_Extended_NegZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpZeroExtendedNegZero)
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

var jumpZeroExtendedNotZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA5, 0, 01),
	jzSourceItemRef(regA5, "target"),
	iarSourceItem(1),
	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JZ_Extended_NotZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpZeroExtendedNotZero)
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
}

var doubleJumpZeroExtendedMode = []*tasm.SourceItem{
	segSourceItem(3),
	labelSourceItem("posZero"),
	dataSourceItem([]uint64{common.PositiveZero}),
	dataSourceItem([]uint64{common.PositiveZero}),

	labelSourceItem("negZero"),
	dataSourceItem([]uint64{common.NegativeZero}),
	dataSourceItem([]uint64{common.NegativeZero}),

	labelSourceItem("notZero1"),
	dataSourceItem([]uint64{common.PositiveZero}),
	dataSourceItem([]uint64{common.PositiveOne}),

	labelSourceItem("notZero2"),
	dataSourceItem([]uint64{common.NegativeZero}),
	dataSourceItem([]uint64{common.PositiveZero}),

	segSourceItem(0),
	dlSourceItemHIBRef(regA0, 0, 0, 0, common.B3, "posZero"),
	djzSourceItemRef(regA0, "target1"),
	iarSourceItem(1),

	labelSourceItem("target1"),
	dlSourceItemHIBRef(regA2, 0, 0, 0, common.B3, "negZero"),
	djzSourceItemRef(regA2, "target2"),
	iarSourceItem(2),

	labelSourceItem("target2"),
	dlSourceItemHIBRef(regA4, 0, 0, 0, common.B3, "notZero1"),
	djzSourceItemRef(regA4, "bad3"),

	labelSourceItem("target3"),
	dlSourceItemHIBRef(regA6, 0, 0, 0, common.B3, "notZero2"),
	djzSourceItemRef(regA6, "bad4"),
	jSourceItemRefExtended("end"),

	labelSourceItem("bad3"),
	iarSourceItem(3),
	labelSourceItem("bad4"),
	iarSourceItem(4),

	labelSourceItem("end"),
	iarSourceItem(0),
}

func Test_DJZ_Extended_PosZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", doubleJumpZeroExtendedMode)
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

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

var jumpNonZeroExtendedPosZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA5, 0, 0),
	jnzSourceItemRef(regA5, "target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNZ_Extended_PosZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNonZeroExtendedPosZero)
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
}

var jumpNonZeroExtendedNegZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jXU, regA5, 0, 0_777777),
	jnzSourceItemRef(regA5, "target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNZ_Extended_NegZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNonZeroExtendedNegZero)
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
}

var jumpNonZeroExtendedNotZero = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA5, 0, 1),
	jnzSourceItemRef(regA5, "target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNZ_Extended_NotZero(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNonZeroExtendedNotZero)
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

var jumpPosNegExtended = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemU(jU, regA10, 0, 0),
	jpSourceItemRef(regA10, "target1"),
	iarSourceItem(1),

	labelSourceItem("target1"),
	jnSourceItemRef(regA10, "bad2"),

	nopExtended(),
	laSourceItemU(jXU, regA10, 0, 0_444444),
	jpSourceItemRef(regA10, "bad3"),
	jnSourceItemRef(regA10, "end"),

	labelSourceItem("bad4"),
	iarSourceItem(4),

	labelSourceItem("bad2"),
	iarSourceItem(2),

	labelSourceItem("bad3"),
	iarSourceItem(3),

	labelSourceItem("end"),
	iarSourceItem(0),
}

func Test_JP_JN_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpPosNegExtended)
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
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
}

//	TODO JPS, JNS
//	TODO JB, JNB
//	TODO JGD, JMGI

var jumpCarryBasic = []*tasm.SourceItem{
	segSourceItem(0),
	jcSourceItemRefBasic("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JC_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpCarryBasic)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCarry(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JC_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpCarryBasic)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpCarryExtended = []*tasm.SourceItem{
	segSourceItem(0),
	jcSourceItemRefExtended("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JC_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpCarryExtended)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCarry(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

var jumpDivideFault = []*tasm.SourceItem{
	segSourceItem(0),
	jdfSourceItemRef("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JDF_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JDF_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JDF_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JDF_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpFloatingOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	jfoSourceItemRef("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JFO_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JFO_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JFO_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JFO_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpFloatingUnderflow = []*tasm.SourceItem{
	segSourceItem(0),
	jfuSourceItemRef("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JFU_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JFU_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JFU_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JFU_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoCarryBasic = []*tasm.SourceItem{
	segSourceItem(0),
	jncSourceItemRefBasic("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNC_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoCarryBasic)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCarry(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNC_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoCarryExtended = []*tasm.SourceItem{
	segSourceItem(0),
	jncSourceItemRefExtended("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNC_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoCarryExtended)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCarry(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

var jumpNoDivideFault = []*tasm.SourceItem{
	segSourceItem(0),
	jndfSourceItemRef("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNDF_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNDF_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JNDF_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNDF_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoDivideFault)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetDivideCheck(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoFloatingOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	jnfoSourceItemRef("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNFO_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNFO_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JNFO_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNFO_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoFloatingUnderflow = []*tasm.SourceItem{
	segSourceItem(0),
	jnfuSourceItemRef("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNFU_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNFU_Basic_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

func Test_JNFU_Extended_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNFU_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoFloatingUnderflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetCharacteristicUnderflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpNoOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	jnoSourceItemRef("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JNO_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JNO_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpNoOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}

var jumpOverflow = []*tasm.SourceItem{
	segSourceItem(0),
	joSourceItemRef("target"),
	iarSourceItem(1),

	labelSourceItem("target"),
	iarSourceItem(0),
}

func Test_JO_Basic_Pos(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 0)
	checkProgramAddress(t, engine, 01003)
}

func Test_JO_Extended_Neg(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", jumpOverflow)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(false)
		ute.GetEngine().GetDesignatorRegister().SetOverflow(false)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	engine := ute.GetEngine()
	checkStoppedReason(t, engine, InitiateAutoRecoveryStop, 1)
	checkProgramAddress(t, engine, 01002)
}
