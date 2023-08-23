// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

var sscCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 0, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 1, 0, 0, 0, 2, "data0"),
	laSourceItemHIBDRef("", jW, 2, 0, 0, 0, 2, "data0"),
	tasm.NewSourceItem("", "fjaxu", []string{fSSC, jSSC, rA0, zero, "0"}),
	tasm.NewSourceItem("", "fjaxu", []string{fSSC, jSSC, rA1, zero, "36"}),
	tasm.NewSourceItem("", "fjaxu", []string{fSSC, jSSC, rA2, zero, "72"}),
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
	checkStopped(t, engine)
	checkRegister(t, engine, A0, 0_112233_445566, "A0")
	checkRegister(t, engine, A1, 0_112233_445566, "A1")
	checkRegister(t, engine, A2, 0_112233_445566, "A2")
}
