// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

var laBasicMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data", "", []string{}),
	tasm.NewSourceItem("a1value", "sw", []string{"01", "02", "03", "04", "05", "06"}),
	tasm.NewSourceItem("a2value", "qw", []string{"0101", "0102", "0103", "0104"}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{"07777"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLA, jU, rA0, rX0, "0123"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jW, rA1, rX0, zero, zero, "a1value"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jQ2, rA2, rX0, zero, zero, "a2value"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLX, jU, rX4, rX0, "05"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jW, rA3, rX4, zero, zero, "data"}),
	tasm.NewSourceItem("", "w", []string{"0"}), //	cause an illop interrupt
}

var laExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data", "", []string{}),
	tasm.NewSourceItem("a1value", "sw", []string{"01", "02", "03", "04", "05", "06"}),
	tasm.NewSourceItem("a2value", "qw", []string{"0101", "0102", "0103", "0104"}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{zero}),
	tasm.NewSourceItem("", "w", []string{"07777"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLA, jU, rA0, rX0, "0123"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jW, rA1, rX0, zero, zero, rB0, "a1value"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jQ2, rA2, rX0, zero, zero, rB0, "a2value"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLX, jU, rX4, rX0, "05"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{fLA, jW, rA3, rX4, zero, zero, rB0, "data"}),
	tasm.NewSourceItem("", "w", []string{"0"}), //	cause an illop interrupt
}

var lxBasicMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLX, jU, rX1, rX0, "05"}),
	tasm.NewSourceItem("", "w", []string{"0"}), //	cause an illop interrupt
}

var lxExtendedMode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{fLX, jU, rX1, rX0, "05"}),
	tasm.NewSourceItem("", "w", []string{"0"}), //	cause an illop interrupt
}

func Test_LA_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", laBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)
	e.Show()

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

	grs := ute.GetEngine().generalRegisterSet
	res := grs.GetRegister(A0).GetW()
	exp := uint64(0123)
	if res != exp {
		t.Fatalf("Register A0 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A1).GetW()
	exp = uint64(0_010203_040506)
	if res != exp {
		t.Fatalf("Register A1 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A2).GetW()
	exp = uint64(0102)
	if res != exp {
		t.Fatalf("Register A2 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A3).GetW()
	exp = uint64(07777)
	if res != exp {
		t.Fatalf("Register A3 is %012o, expected %012o", res, exp)
	}
}

func Test_LA_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", laExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)
	e.Show()

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetQuarterWordModeEnabled(true)
		err = ute.Run()
	}

	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	grs := ute.GetEngine().generalRegisterSet
	res := grs.GetRegister(A0).GetW()
	exp := uint64(0123)
	if res != exp {
		t.Fatalf("Register A0 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A1).GetW()
	exp = uint64(0_010203_040506)
	if res != exp {
		t.Fatalf("Register A1 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A2).GetW()
	exp = uint64(0102)
	if res != exp {
		t.Fatalf("Register A2 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A3).GetW()
	exp = uint64(07777)
	if res != exp {
		t.Fatalf("Register A3 is %012o, expected %012o", res, exp)
	}
}

func Test_LX_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lxBasicMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)
	e.Show()

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

	grs := ute.GetEngine().generalRegisterSet
	res := grs.GetRegister(X1).GetW()
	exp := uint64(05)
	if res != exp {
		t.Fatalf("Register X1 is %012o, expected %012o", res, exp)
	}
}

func Test_LX_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", lxExtendedMode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)
	e.Show()

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

	grs := ute.GetEngine().generalRegisterSet
	res := grs.GetRegister(X1).GetW()
	exp := uint64(05)
	if res != exp {
		t.Fatalf("Register X1 is %012o, expected %012o", res, exp)
	}
}
