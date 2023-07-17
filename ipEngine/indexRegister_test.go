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
