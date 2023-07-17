package ipEngine

import "testing"

func Test_DecrementModifier_1(t *testing.T) {
	x := IndexRegister(0_000100_100000)
	x.DecrementModifier()
	expected := uint64(0_000100_077700)
	result := x.GetW()
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_DecrementModifier_2(t *testing.T) {
	x := IndexRegister(0_000011_000000)
	x.DecrementModifier()
	expected := uint64(0_000011_777766)
	result := x.GetW()
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_DecrementModifier24_1(t *testing.T) {
	x := IndexRegister(0_0100_0010_0000)
	x.DecrementModifier24()
	expected := uint64(0_0100_0007_7700)
	result := x.GetW()
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_DecrementModifier24_2(t *testing.T) {
	x := IndexRegister(0_0011_0000_0000)
	x.DecrementModifier24()
	expected := uint64(0_0011_7777_7766)
	result := x.GetW()
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_IncrementModifier_1(t *testing.T) {
	x := IndexRegister(0_000100_100000)
	x.IncrementModifier()
	expected := uint64(0_000100_100100)
	result := x.GetW()
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_IncrementModifier_2(t *testing.T) {
	x := IndexRegister(0_000011_777777)
	x.IncrementModifier()
	expected := uint64(0_000011_000011)
	result := x.GetW()
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_IncrementModifier24_1(t *testing.T) {
	x := IndexRegister(0_0100_0010_0000)
	x.IncrementModifier24()
	expected := uint64(0_0100_0010_0100)
	result := x.GetW()
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_IncrementModifier24_2(t *testing.T) {
	x := IndexRegister(0_0010_3777_7777)
	x.IncrementModifier24()
	expected := uint64(0_0010_4000_0007)
	result := x.GetW()
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}
