package pkg

import "testing"

func Test_AddSimple_2(t *testing.T) {
	value1 := uint64(0_300000_000000)
	value2 := uint64(0_077777_777777)
	expected := uint64(0_377777_777777)
	result := AddSimple(value1, value2)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_AddSimple_1(t *testing.T) {
	value1 := uint64(0_777777_777722)
	value2 := uint64(0_000000_000055)
	expected := uint64(0)
	result := AddSimple(value1, value2)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetOnesComplement_1(t *testing.T) {
	value := uint64(100234)
	expected := uint64(100234)
	result := GetOnesComplement(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetOnesComplement_2(t *testing.T) {
	// -17dec is -021oct
	value := uint64(0xFFFFFFFF_FFFFFFEF)
	expected := uint64(0_777777_777756)
	result := GetOnesComplement(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended12_1(t *testing.T) {
	value := uint64(0_776644_011111)
	expected := uint64(0_000000_001111)
	result := GetSignExtended12(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended12_2(t *testing.T) {
	value := uint64(0_776644_004111)
	expected := uint64(0_777777_774111)
	result := GetSignExtended12(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended18_1(t *testing.T) {
	value := uint64(0_776644_311111)
	expected := uint64(0_000000_311111)
	result := GetSignExtended18(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended18_2(t *testing.T) {
	value := uint64(0_000004_404111)
	expected := uint64(0_777777_404111)
	result := GetSignExtended18(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended24_1(t *testing.T) {
	value := uint64(0_776637_311111)
	expected := uint64(0_000037_311111)
	result := GetSignExtended24(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetSignExtended24_2(t *testing.T) {
	value := uint64(0_0066_4440_4111)
	expected := uint64(0_7777_4440_4111)
	result := GetSignExtended24(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetTwosComplement_1(t *testing.T) {
	value := uint64(100000)
	expected := uint64(100000)
	result := GetTwosComplement(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetTwosComplement_2(t *testing.T) {
	value := uint64(0_777777_777770)
	expected := uint64(0xFFFFFFFF_FFFFFFF9)
	result := GetTwosComplement(uint64(value))
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func Test_GetTwosComplement_3(t *testing.T) {
	value := uint64(0_777777_777777)
	expected := uint64(0)
	result := GetTwosComplement(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}

func TestNegate(t *testing.T) {
	value := uint64(0377_123456)
	expected := uint64(0_777400_654321)
	result := Negate(value)
	if result != expected {
		t.Fatalf("Error expected %12o, got %12o", expected, result)
	}
}
