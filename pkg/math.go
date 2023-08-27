// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

var DoubleNegativeZero = []uint64{NegativeZero, NegativeZero}
var DoublePositiveZero = []uint64{PositiveZero, PositiveZero}

func uint64Min(v1 uint64, v2 uint64) uint64 {
	if v1 < v2 {
		return v1
	} else {
		return v2
	}
}

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

// And calculates the logical AND of two 36-bit values
func And(operand1 uint64, operand2 uint64) uint64 {
	return (operand1 & operand2) & NegativeZero
}

// Compare indicates whether operand1 is less than, equal to, or greater than operand2.
// the result negative if less than, zero if equal, or positive if greater than.
// For our purposes, negative zero is less than (and thus, not equal to) positive zero.
func Compare(operand1 uint64, operand2 uint64) int {
	if operand1 == operand2 {
		return 0
	}

	pos1 := IsPositive(operand1)
	pos2 := IsPositive(operand2)
	if pos1 && pos2 {
		if operand1 < operand2 {
			return -1
		} else {
			return 1
		}
	} else if !pos1 && !pos2 {
		if operand1 > operand2 {
			return -1
		} else {
			return 1
		}
	} else if pos1 {
		return 1
	} else {
		return -1
	}
}

// CompareDouble indicates whether operand1 is less than, equal to, or greater than operand2.
// the result negative if less than, zero if equal, or positive if greater than.
// For our purposes, negative zero is less than (and thus, not equal to) positive zero.
// both operands consist of a 72-bit value, stored as two consecutive 36-bit values wrapped
// in uint64's.
func CompareDouble(operand1 []uint64, operand2 []uint64) int {
	if (operand1[0] == operand2[0]) && (operand1[1] == operand2[1]) {
		return 0
	}

	pos1 := IsPositive(operand1[0])
	pos2 := IsPositive(operand2[0])
	if pos1 != pos2 {
		if pos1 {
			return 1
		} else {
			return -1
		}
	} else if pos1 {
		if operand1[0] > operand2[0] {
			return 1
		} else if operand1[0] < operand2[0] {
			return -1
		} else {
			if operand1[1] > operand2[1] {
				return 1
			} else if operand1[1] < operand2[1] {
				return -1
			} else {
				return 0
			}
		}
	} else {
		if operand1[0] > operand2[0] {
			return -1
		} else if operand1[0] < operand2[0] {
			return 1
		} else {
			if operand1[1] > operand2[1] {
				return -1
			} else if operand1[1] < operand2[1] {
				return 1
			} else {
				return 0
			}
		}
	}
}

// GetOnesComplement takes a standard twos-complement value and converts it to a
// 36-bit ones-complement value packed in a uint64.
func GetOnesComplement(operand uint64) uint64 {
	if int64(operand) < 0 {
		return Negate(-operand)
	} else {
		return operand
	}
}

// GetSignExtended12 sign-extends an 12-bit value to 36 bits
func GetSignExtended12(value uint64) (result uint64) {
	result = value & 0_7777
	if (result & 0_04000) != 0 {
		result |= 0_777777_770000
	}
	return
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
	result = value & 0_7777_7777
	if (result & 0_4000_0000) != 0 {
		result |= 0_7777_0000_0000
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

func IsNegative(value uint64) bool {
	return (value & 0_400000_000000) != 0
}

func IsPositive(value uint64) bool {
	return (value & 0_400000_000000) == 0
}

func IsZero(operand uint64) bool {
	return operand == PositiveZero || operand == NegativeZero
}

func IsDoubleZero(operand []uint64) bool {
	return IsZero(operand[0]) && operand[0] == operand[1]
}

// LeftDoubleShiftCircular shifts the 72-bit word stored in two consecutive uint64's (MSW first)
// to the left by the given count value, where every bit shifted out of bit 0 is end-around shifted into bit 71.
func LeftDoubleShiftCircular(operand []uint64, count uint64) []uint64 {
	result := []uint64{operand[0], operand[1]}
	count %= 72
	for count > 0 {
		if count >= 36 {
			r := result[0]
			result[0] = result[1]
			result[1] = r
			count -= 36
		} else {
			shift := uint64Min(27, count)
			result[0] <<= shift
			result[1] <<= shift
			result[1] |= result[0] >> 36
			result[0] |= result[1] >> 36
			result[0] &= NegativeZero
			result[1] &= NegativeZero
			count -= shift
		}
	}

	return result
}

// LeftDoubleShiftLogical shifts the 72-bit word stored in two consecutive uint64's (MSW first)
// to the left by the given count value. Bits shifted out of bit 0 are lost, and zeroes are shift into bit 71.
func LeftDoubleShiftLogical(operand []uint64, count uint64) []uint64 {
	result := []uint64{operand[0], operand[1]}
	count %= 72
	if count > 0 {
		if count >= 36 {
			result[0] = result[1]
			result[1] = PositiveZero
			count -= 36
		}

		for count > 0 {
			shift := uint64Min(27, count)
			result[0] <<= shift
			result[1] <<= shift
			result[0] |= result[1] >> 36
			result[1] &= NegativeZero
			count -= shift
		}
		result[0] &= NegativeZero
	}

	return result
}

// LeftShiftCircular shifts the 36-bit word to the left by the given count value,
// where every bit shifted out of bit 0 is end-around shifted into bit 35.
func LeftShiftCircular(operand uint64, count uint64) uint64 {
	result := operand
	count %= 36
	if count > 0 {
		shift := uint64Min(27, count)
		result <<= shift
		result |= result >> 36
		result &= NegativeZero
		count -= shift
	}

	return result
}

// LeftShiftLogical shifts the 36-bit word to the left by the given count value.
// Bits shifted out of bit 0 are lost, and zeroes are shift into bit 35.
func LeftShiftLogical(operand uint64, count uint64) uint64 {
	result := operand
	if count > 0 {
		if count >= 36 {
			result = PositiveZero
		} else {
			result = (operand << count) & NegativeZero
		}
	}

	return result
}

// Magnitude returns the absolute value of the given operand
func Magnitude(operand uint64) uint64 {
	if IsPositive(operand) {
		return operand
	} else {
		return Not(operand)
	}
}

// MagnitudeDouble returns the absolute value of the given 72-bit operand
func MagnitudeDouble(operand uint64) uint64 {
	if IsPositive(operand) {
		return operand
	} else {
		return Not(operand)
	}
}

// Negate returns the additive inverse of a given 36-bit signed value packed into a uint64
func Negate(operand uint64) uint64 {
	return operand ^ NegativeZero
}

func Not(op uint64) uint64 {
	return (op ^ NegativeZero) & NegativeZero
}

func Or(lhs uint64, rhs uint64) uint64 {
	return (lhs | rhs) & NegativeZero
}

func Xor(lhs uint64, rhs uint64) uint64 {
	return (lhs ^ rhs) & NegativeZero
}

// RightDoubleShiftAlgebraic shifts the 72-bit word stored in two consecutive uint64's (MSW first)
// to the right. Bits shifted out of bit 71 are lost while bit 0 is propagated to the right.
func RightDoubleShiftAlgebraic(operand []uint64, count uint64) []uint64 {
	var result []uint64
	if count > 71 {
		if IsNegative(operand[0]) {
			result = DoubleNegativeZero
		} else {
			result = DoublePositiveZero
		}
	} else {
		if count > 0 {
			neg := IsNegative(operand[0])
			if count > 36 {
				if neg {
					result[0] = NegativeZero
				} else {
					result[0] = PositiveZero
				}
				result[1] = operand[0]
				count -= 36
			} else {
				result = []uint64{operand[0], operand[1]}
			}

			for count > 0 {
				shift := uint64Min(27, count)
				mask := uint64(1<<shift) - 1
				partial := result[0] & mask
				result[0] >>= shift
				if neg {
					bits := mask << (36 - count)
					result[0] |= bits
				}
				result[1] |= partial << 36
				result[1] >>= shift
				count -= shift
			}
		} else {
			result = operand
		}
	}
	return result
}

// RightDoubleShiftCircular shifts the 72-bit word stored in two consecutive uint64's (MSW first)
// where every bit shifted out of bit 35 is end-around shifted into bit 0.
func RightDoubleShiftCircular(operand []uint64, count uint64) []uint64 {
	result := []uint64{operand[0], operand[1]}
	count %= 72
	for count > 0 {
		if count >= 36 {
			r := result[0]
			result[0] = result[1]
			result[1] = r
			count -= 36
		} else {
			shift := uint64Min(27, count)
			mask := uint64(1<<shift) - 1
			result[0] |= (result[1] & mask) << 36
			result[1] |= (result[0] & mask) << 36
			result[0] >>= shift
			result[1] >>= shift
			count -= shift
		}
	}
	return result
}

// RightDoubleShiftLogical shifts the 72-bit word stored in two consecutive uint64's (MSW first)
// Bits shifted out of bit 35 are lost, and zeroes are shift into bit 0.
func RightDoubleShiftLogical(operand []uint64, count uint64) []uint64 {
	var result []uint64
	if count >= 72 {
		result = DoublePositiveZero
	} else {
		result = []uint64{operand[0], operand[1]}
		if count > 36 {
			result[1] = result[0]
			result[0] = PositiveZero
			count -= 36
		}

		for count > 0 {
			shift := uint64Min(count, 27)
			mask := uint64(1<<shift) - 1
			partial := (result[0] & mask) << (36 - count)
			result[0] >>= count

			result[1] >>= shift
			result[1] |= partial
			count -= shift
		}
	}
	return result
}

// RightShiftAlgebraic shifts the 72-bit word to the left by the given count value,
// Bits shifted out of bit 35 are lost while bit 0 is propagated to the right.
func RightShiftAlgebraic(operand uint64, count uint64) uint64 {
	var result uint64
	if count >= 35 {
		if IsNegative(operand) {
			result = NegativeZero
		} else {
			result = PositiveZero
		}
	} else {
		if count > 0 {
			result = operand >> count
			if IsNegative(operand) {
				propMask := uint64((1<<count)-1) << (36 - count)
				return result | propMask
			}
		} else {
			result = operand
		}
	}
	return result
}

// RightShiftCircular shifts the 36-bit word to the right by the given count value,
// where every bit shifted out of bit 35 is end-around shifted into bit 0.
func RightShiftCircular(operand uint64, count uint64) uint64 {
	result := operand
	count %= 36
	for count > 0 {
		shift := uint64Min(27, count)
		mask := uint64(1<<shift) - 1
		partial := (result & mask) << 36
		result = (partial | result) >> count
		count -= shift
	}

	return result
}

// RightShiftLogical shifts the 36-bit word to the left by the given count value.
// Bits shifted out of bit 35 are lost, and zeroes are shift into bit 0.
func RightShiftLogical(operand uint64, count uint64) uint64 {
	result := operand
	if count >= 36 {
		result = PositiveZero
	} else {
		result >>= count
	}

	return result
}
