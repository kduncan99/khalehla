// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

var partialwordloadsBasicThirdWord = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data1", "w", []string{"0556677001122"}),
	tasm.NewSourceItem("data2", "hw", []string{"0506070", "0507777"}),
	tasm.NewSourceItem("data3", "hw", []string{"0223344", "0221100"}),
	tasm.NewSourceItem("data4", "tw", []string{"01111", "02222", "03333"}),
	tasm.NewSourceItem("data5", "tw", []string{"05500", "06600", "07700"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jW, rA0, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jH1, rA1, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jH2, rA2, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jXH1, rA3, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jXH1, rA4, zero, zero, zero, "data3"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jXH2, rA5, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jXH2, rA6, zero, zero, zero, "data3"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT1, rA7, zero, zero, zero, "data4"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT1, rA8, zero, zero, zero, "data5"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT2, rA9, zero, zero, zero, "data4"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT2, rA10, zero, zero, zero, "data5"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT3, rA11, zero, zero, zero, "data4"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLA, jT3, rA12, zero, zero, zero, "data5"}),
	IARSourceItem("", "0"),
}

var partialwordloadsBasicQuarterWord = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("data1", "qw", []string{"0400", "0501", "0677", "0777"}),
	tasm.NewSourceItem("data2", "sw", []string{"012", "034", "056", "075", "042", "010"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jQ1, rA0, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jQ2, rA1, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jQ3, rA2, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jQ4, rA3, zero, zero, zero, "data1"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS1, rA4, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS2, rA5, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS3, rA6, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS4, rA7, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS5, rA8, zero, zero, zero, "data2"}),
	tasm.NewSourceItem("", "fjaxhiu", []string{fLR, jS6, rA9, zero, zero, zero, "data2"}),
	IARSourceItem("", "0"),
}

func Test_PartialWordLoads_BasicThirdWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", partialwordloadsBasicThirdWord)
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

	grs := ute.GetEngine().generalRegisterSet
	res := grs.GetRegister(A0).GetW()
	exp := uint64(0556677001122)
	if res != exp {
		t.Fatalf("Register A0 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A1).GetW()
	exp = uint64(0506070)
	if res != exp {
		t.Fatalf("Register A1 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A2).GetW()
	exp = uint64(0507777)
	if res != exp {
		t.Fatalf("Register A2 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A3).GetW()
	exp = uint64(0777777506070)
	if res != exp {
		t.Fatalf("Register A3 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A4).GetW()
	exp = uint64(0223344)
	if res != exp {
		t.Fatalf("Register A4 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A5).GetW()
	exp = uint64(0777777507777)
	if res != exp {
		t.Fatalf("Register A5 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A6).GetW()
	exp = uint64(0221100)
	if res != exp {
		t.Fatalf("Register A6 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A7).GetW()
	exp = uint64(01111)
	if res != exp {
		t.Fatalf("Register A7 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A8).GetW()
	exp = uint64(0_777777_775500)
	if res != exp {
		t.Fatalf("Register A8 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A9).GetW()
	exp = uint64(02222)
	if res != exp {
		t.Fatalf("Register A9 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A10).GetW()
	exp = uint64(0_777777_776600)
	if res != exp {
		t.Fatalf("Register A10 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A11).GetW()
	exp = uint64(03333)
	if res != exp {
		t.Fatalf("Register A11 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(A12).GetW()
	exp = uint64(0_777777_777700)
	if res != exp {
		t.Fatalf("Register A12 is %012o, expected %012o", res, exp)
	}
}

func Test_PartialWordLoads_BasicQuarterWord(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", partialwordloadsBasicQuarterWord)
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

	grs := ute.GetEngine().generalRegisterSet
	res := grs.GetRegister(R0).GetW()
	exp := uint64(0400)
	if res != exp {
		t.Fatalf("Register R0 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(R1).GetW()
	exp = uint64(0501)
	if res != exp {
		t.Fatalf("Register R1 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(R2).GetW()
	exp = uint64(0677)
	if res != exp {
		t.Fatalf("Register R2 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(R3).GetW()
	exp = uint64(0777)
	if res != exp {
		t.Fatalf("Register R3 is %012o, expected %012o", res, exp)
	}

	grs = ute.GetEngine().generalRegisterSet
	res = grs.GetRegister(R4).GetW()
	exp = uint64(012)
	if res != exp {
		t.Fatalf("Register R4 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(R5).GetW()
	exp = uint64(034)
	if res != exp {
		t.Fatalf("Register R5 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(R6).GetW()
	exp = uint64(056)
	if res != exp {
		t.Fatalf("Register R6 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(R7).GetW()
	exp = uint64(075)
	if res != exp {
		t.Fatalf("Register R7 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(R8).GetW()
	exp = uint64(042)
	if res != exp {
		t.Fatalf("Register R8 is %012o, expected %012o", res, exp)
	}

	res = grs.GetRegister(R9).GetW()
	exp = uint64(010)
	if res != exp {
		t.Fatalf("Register R9 is %012o, expected %012o", res, exp)
	}

}

//	TODO basic mode partial word testing, stores

//	TODO basic mode index register handling

//	TODO basic mode address determination

//	TODO extended mode partial word testing, loads

//	TODO extended mode partial word testing, stores

//	TODO extended mode index register handling

//	TODO extended mode addressing across multiple banks
