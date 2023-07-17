// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

// AddSimple takes two numbers which are 36-bit signed values packed into uint64's,
// and adds them according to ones-complement rules.
func AddSimple(operand1 uint64, operand2 uint64) uint64 {
	if (operand1 == NegativeZero) && (operand2 == NegativeZero) {
		return NegativeZero
	} else {
		native1 := GetTwosComplement(operand1)
		native2 := GetTwosComplement(operand2)
		return GetOnesComplement(native1 + native2)
	}
}

// GetOnesComplement takes a standard twos-complement value and converts it to a
// 36-bit ones-complement value packed in a uint64.
func GetOnesComplement(operand uint64) uint64 {
	if operand < 0 {
		return Negate(-operand)
	} else {
		return operand
	}
}

// GetSignExtended12 sign-extends an 12-bit value to 36 bits
func GetSignExtended12(value uint64) uint64 {
	if (value & 04000) == 0 {
		return value
	} else {
		return value | 0_777777_770000
	}
}

// GetSignExtended18 sign-extends an 18-bit value to 36 bits
func GetSignExtended18(value uint64) (result uint64) {
	result = value & 0_777777
	if (result & 0_400000) != 0 {
		result |= 0_777777_000000
	}
	return
}

// GetSignExtended24 sign-extends a 24-bit value to 36 bits
func GetSignExtended24(value uint64) (result uint64) {
	result = value & 077_777777
	if (result & 040_000000) != 0 {
		result |= 0_777700_000000
	}
	return
}

// GetTwosComplement takes a number which is a 36-bit signed value packed into a uint64,
// and converts it to twos-complement.
func GetTwosComplement(operand uint64) uint64 {
	if IsNegative(operand) {
		return -Negate(operand)
	} else {
		return operand
	}
}

// Negate returns the additive inverse of a given 36-bit signed value packed into a uint64
func Negate(operand uint64) uint64 {
	return operand ^ NegativeZero
}
