// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
	"khalehla/tasm"
	"testing"
)

var sscCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 0, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 1, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 2, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 3, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 5, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 6, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 7, 0, 0, 0, 2, "data0"),
	sscSourceItemU("", 0, 0, 0, 0, 0),
	sscSourceItemU("", 1, 0, 0, 0, 36),
	sscSourceItemU("", 2, 0, 0, 0, 72),
	sscSourceItemU("", 3, 0, 0, 0, 18),
	sscSourceItemU("", 4, 0, 0, 0, 1),
	sscSourceItemU("", 5, 0, 0, 0, 35),
	sscSourceItemU("", 6, 0, 0, 0, 0206), // should just be 06, the 0200 gets stripped out
	iarSourceItem("end", 0),

	segSourceItem(2),
	tasm.NewSourceItem("data0", "w", []string{"0112233445566"}),
}

func Test_SSC(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", sscCode)
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
	checkRegister(t, engine, pkg.A0, 0_112233_445566, "A0")
	checkRegister(t, engine, pkg.A1, 0_112233_445566, "A1")
	checkRegister(t, engine, pkg.A2, 0_112233_445566, "A2")
	checkRegister(t, engine, pkg.A3, 0_445566_112233, "A3")
	checkRegister(t, engine, pkg.A4, 0_045115_622673, "A4")
	checkRegister(t, engine, pkg.A5, 0_224467_113354, "A5")
	checkRegister(t, engine, pkg.A6, 0_661122_334455, "A6")
}

// TODO DSC

var sslCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 0, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 1, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 2, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 3, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data1"),
	laSourceItemHIBDRef("", jW, 5, 0, 0, 0, 2, "data1"),
	laSourceItemHIBDRef("", jW, 6, 0, 0, 0, 2, "data2"),
	laSourceItemHIBDRef("", jW, 7, 0, 0, 0, 2, "data2"),
	sslSourceItemU("", 0, 0, 0, 0, 0),
	sslSourceItemU("", 1, 0, 0, 0, 36),
	sslSourceItemU("", 2, 0, 0, 0, 72),
	sslSourceItemU("", 3, 0, 0, 0, 18),
	sslSourceItemU("", 4, 0, 0, 0, 1),
	sslSourceItemU("", 5, 0, 0, 0, 35),
	sslSourceItemU("", 6, 0, 0, 0, 1),
	sslSourceItemU("", 7, 0, 0, 0, 35),
	iarSourceItem("end", 0),

	segSourceItem(2),
	tasm.NewSourceItem("data0", "w", []string{"0112233445566"}),
	tasm.NewSourceItem("data1", "w", []string{"0555555555555"}),
	tasm.NewSourceItem("data2", "w", []string{"0111111666666"}),
}

func Test_SSL(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", sslCode)
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
	checkRegister(t, engine, pkg.A0, 0_112233_445566, "A0")
	checkRegister(t, engine, pkg.A1, 0, "A1")
	checkRegister(t, engine, pkg.A2, 0, "A2")
	checkRegister(t, engine, pkg.A3, 0_112233, "A3")
	checkRegister(t, engine, pkg.A4, 0_266666_666666, "A4")
	checkRegister(t, engine, pkg.A5, 01, "A5")
	checkRegister(t, engine, pkg.A6, 0_044444_733333, "A6")
	checkRegister(t, engine, pkg.A7, 0, "A7")
}

// TODO DSL
// TODO SSA
// TODO DSA
// TODO LSC
// TODO DLSC

var lsscCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 0, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 1, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 2, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 3, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 5, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 6, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 7, 0, 0, 0, 2, "data0"),
	lsscSourceItemU("", 0, 0, 0, 0, 0),
	lsscSourceItemU("", 1, 0, 0, 0, 36),
	lsscSourceItemU("", 2, 0, 0, 0, 72),
	lsscSourceItemU("", 3, 0, 0, 0, 18),
	lsscSourceItemU("", 4, 0, 0, 0, 1),
	lsscSourceItemU("", 5, 0, 0, 0, 35),
	lsscSourceItemU("", 6, 0, 0, 0, 0206), // should just be 06, the 0200 gets stripped out
	iarSourceItem("end", 0),

	segSourceItem(2),
	tasm.NewSourceItem("data0", "w", []string{"0112233445566"}),
}

func Test_LSSC(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lsscCode)
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
	checkRegister(t, engine, pkg.A0, 0_112233_445566, "A0")
	checkRegister(t, engine, pkg.A1, 0_112233_445566, "A1")
	checkRegister(t, engine, pkg.A2, 0_112233_445566, "A2")
	checkRegister(t, engine, pkg.A3, 0_445566_112233, "A3")
	checkRegister(t, engine, pkg.A4, 0_224467_113354, "A4")
	checkRegister(t, engine, pkg.A5, 0_045115_622673, "A5")
	checkRegister(t, engine, pkg.A6, 0_223344_556611, "A6")
}

// TODO LDSC
// TODO LSSL
// TODO LDSL
