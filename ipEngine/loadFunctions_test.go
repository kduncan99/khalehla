package ipEngine

import (
	"khalehla/tasm"
	"testing"
)

var f_la = "010"

var j_w = "0"
var j_s1 = "015"
var j_s2 = "014"
var j_s3 = "013"
var j_s4 = "012"
var j_s5 = "011"
var j_s6 = "010"
var j_u = "016"
var j_xu = "017"

var a0 = "0"
var a1 = "01"
var a2 = "02"
var a3 = "03"
var a4 = "04"
var a5 = "05"

var b0 = "0"
var b1 = "01"
var b2 = "02"
var b3 = "03"

var x0 = "0"
var x1 = "01"
var x2 = "02"
var x3 = "03"

var z = "0"

var laFunctionCode = []*tasm.SourceItem{
	tasm.NewSourceItem("", ".SEG", []string{"077"}),
	tasm.NewSourceItem("a1value", "sw", []string{"01", "02", "03", "04", "05", "06"}),

	tasm.NewSourceItem("", ".SEG", []string{"000"}),
	tasm.NewSourceItem("", "fjaxu", []string{f_la, j_u, a0, x0, "0123"}),
	tasm.NewSourceItem("", "fjaxhibd", []string{f_la, j_w, a1, x0, z, z, b0, "a1value"}),
	tasm.NewSourceItem("", "w", []string{"0"}), //	cause an illop interrupt
}

func Test_LA_Basic(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", laFunctionCode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), false)
	e.Show()

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
		ute.GetEngine().GetDesignatorRegister().SetBasicModeEnabled(true)
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
		t.Fatalf("Register A0 is %012o, expected %012o", res, exp)
	}
}

func Test_LA_Extended(t *testing.T) {
	sourceSet := tasm.NewSourceSet("Test", laFunctionCode)
	a := tasm.NewTinyAssembler()
	a.Assemble(sourceSet)

	e := tasm.Executable{}
	e.LinkSimple(a.GetSegments(), true)
	e.Show()

	ute := NewUnitTestExecutor()
	err := ute.Load(&e)
	if err == nil {
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
		t.Fatalf("Register A0 is %012o, expected %012o", res, exp)
	}
}
