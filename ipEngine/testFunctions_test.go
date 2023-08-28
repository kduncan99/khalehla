// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

var tepCode = []*tasm.SourceItem{
	segSourceItem(0),
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data1"),
	tepSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data2"),
	jSourceItemExtended("", 0, 0, 0, "tag"),
	iarSourceItem("badend1", 1),
	tepSourceItemHIBDRef("tag", jQ4, 4, 0, 0, 0, 2, "data2"),
	iarSourceItem("badend2", 2),
	iarSourceItem("goodend", 0),

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
	laSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data1"),
	topSourceItemHIBDRef("", jW, 4, 0, 0, 0, 2, "data2"),
	jSourceItemExtended("", 0, 0, 0, "tag"),
	iarSourceItem("badend1", 1),
	topSourceItemHIBDRef("tag", jQ4, 4, 0, 0, 0, 2, "data2"),
	iarSourceItem("badend2", 2),
	iarSourceItem("goodend", 0),

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
