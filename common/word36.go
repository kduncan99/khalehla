// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package common

import "math/big"

const (
	PositiveOne  = 01
	PositiveZero = 0
	Mask36       = 0_777777_777777
	NegativeOne  = 0_777777_777776
	NegativeZero = 0_777777_777777
)

type Word36 uint64

var DoubleNegativeZero = []uint64{NegativeZero, NegativeZero}
var DoublePositiveZero = []uint64{PositiveZero, PositiveZero}

var bigMask = big.NewInt(0_777777_777777)
var bigMaxQuotient = big.NewInt(0_377777_777777)

func uint64Min(v1 uint64, v2 uint64) uint64 {
	if v1 < v2 {
		return v1
	} else {
		return v2
	}
}

// Partial word getters --------------------------------------------------------

// GetH1 retrieves H1 of the given value as an unsigned integer
func (w *Word36) GetH1() uint64 {
	return GetH1(uint64(*w))
}

func GetH1(value uint64) uint64 {
	return (value >> 18) & 0777777
}

// GetH2 retrieves H2 of the given value as an unsigned integer
func (w *Word36) GetH2() uint64 {
	return GetH1(uint64(*w))
}

func GetH2(value uint64) uint64 {
	return value & 0777777
}

// GetQ1 retrieves Q1 of the given value as an unsigned integer
func (w *Word36) GetQ1() uint64 {
	return GetQ1(uint64(*w))
}

func GetQ1(value uint64) uint64 {
	return (value >> 27) & 0777
}

// GetQ2 retrieves Q2 of the given value as an unsigned integer
func (w *Word36) GetQ2() uint64 {
	return GetQ2(uint64(*w))
}

func GetQ2(value uint64) uint64 {
	return (value >> 18) & 0777
}

// GetQ3 retrieves Q3 of the given value as an unsigned integer
func (w *Word36) GetQ3() uint64 {
	return GetQ3(uint64(*w))
}

func GetQ3(value uint64) uint64 {
	return (value >> 9) & 0777
}

// GetQ4 retrieves Q4 of the given value as an unsigned integer
func (w *Word36) GetQ4() uint64 {
	return GetQ4(uint64(*w))
}

func GetQ4(value uint64) uint64 {
	return value & 0777
}

// GetS1 retrieves S1 of the given value as an unsigned integer
func (w *Word36) GetS1() uint64 {
	return GetS1(uint64(*w))
}

func GetS1(value uint64) uint64 {
	return (value >> 30) & 077
}

// GetS2 retrieves S2 of the given value as an unsigned integer
func (w *Word36) GetS2() uint64 {
	return GetS2(uint64(*w))
}

func GetS2(value uint64) uint64 {
	return (value >> 24) & 077
}

// GetS3 retrieves S3 of the given value as an unsigned integer
func (w *Word36) GetS3() uint64 {
	return GetS3(uint64(*w))
}

func GetS3(value uint64) uint64 {
	return (value >> 18) & 077
}

// GetS4 retrieves S4 of the given value as an unsigned integer
func (w *Word36) GetS4() uint64 {
	return GetS4(uint64(*w))
}

func GetS4(value uint64) uint64 {
	return (value >> 12) & 077
}

// GetS5 retrieves S5 of the given value as an unsigned integer
func (w *Word36) GetS5() uint64 {
	return GetS5(uint64(*w))
}

func GetS5(value uint64) uint64 {
	return (value >> 6) & 077
}

// GetS6 retrieves S6 of the given value as an unsigned integer
func (w *Word36) GetS6() uint64 {
	return GetS6(uint64(*w))
}

func GetS6(value uint64) uint64 {
	return value & 077
}

// GetT1 retrieves T1 of the given value as an unsigned integer
func (w *Word36) GetT1() uint64 {
	return GetT1(uint64(*w))
}

func GetT1(value uint64) uint64 {
	return (value >> 24) & 07777
}

// GetT2 retrieves T2 of the given value as an unsigned integer
func (w *Word36) GetT2() uint64 {
	return GetT2(uint64(*w))
}

func GetT2(value uint64) uint64 {
	return (value >> 12) & 07777
}

// GetT3 retrieves T3 of the given value as an unsigned integer
func (w *Word36) GetT3() uint64 {
	return GetT3(uint64(*w))
}

func GetT3(value uint64) uint64 {
	return value & 07777
}

// GetXH1 retrieves H1 of the given value and returns it, sign-extended to 36-bits
func (w *Word36) GetXH1() uint64 {
	return GetXH1(uint64(*w))
}

func GetXH1(value uint64) uint64 {
	res := (value >> 18) & 0_777777
	if (res & 0_400000) != 0 {
		res |= 0_777777_000000
	}
	return res
}

// GetXH2 retrieves H2 of the given value and returns it, sign-extended to 36-bits
func (w *Word36) GetXH2() uint64 {
	return GetXH2(uint64(*w))
}

func GetXH2(value uint64) uint64 {
	res := value & 0_777777
	if (res & 0_400000) != 0 {
		res |= 0_777777_000000
	}
	return res
}

// GetXT1 retrieves T1 of the given value and returns it, sign-extended to 36-bits
func (w *Word36) GetXT1() uint64 {
	return GetXT1(uint64(*w))
}

func GetXT1(value uint64) uint64 {
	res := (value >> 24) & 0_7777
	if (res & 004000) != 0 {
		res |= 0_777777_770000
	}
	return res
}

// GetXT2 retrieves T2 of the given value and returns it, sign-extended to 36-bits
func (w *Word36) GetXT2() uint64 {
	return GetXT2(uint64(*w))
}

func GetXT2(value uint64) uint64 {
	res := (value >> 12) & 0_7777
	if (res & 004000) != 0 {
		res |= 0_777777_770000
	}
	return res
}

// GetXT3 retrieves T3 of the given value and returns it, sign-extended to 36-bits
func (w *Word36) GetXT3() uint64 {
	return GetXT3(uint64(*w))
}

func GetXT3(value uint64) uint64 {
	res := value & 0_7777
	if (res & 004000) != 0 {
		res |= 0_777777_770000
	}
	return res
}

// GetW retrieves the given value, masked to the right-most 36 bits
func (w *Word36) GetW() uint64 {
	return uint64(*w)
}

func GetW(value uint64) uint64 {
	return value & Mask36
}

// Partial word setters --------------------------------------------------------

// SetH1 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetH1(new uint64) *Word36 {
	*w = Word36(SetH1(uint64(*w), new))
	return w
}

// SetH1 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetH1(orig uint64, new uint64) uint64 {
	return (orig & 0_777777) | ((new & 0_777777) << 18)
}

// SetH2 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetH2(new uint64) *Word36 {
	*w = Word36(SetH2(uint64(*w), new))
	return w
}

// SetH2 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetH2(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_000000) | (new & 0_777777)
}

// SetQ1 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetQ1(new uint64) *Word36 {
	*w = Word36(SetQ1(uint64(*w), new))
	return w
}

// SetQ1 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetQ1(orig uint64, new uint64) uint64 {
	return (orig & 0_000777_777777) | ((new & 0_777) << 27)
}

// SetQ2 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetQ2(new uint64) *Word36 {
	*w = Word36(SetQ2(uint64(*w), new))
	return w
}

// SetQ2 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetQ2(orig uint64, new uint64) uint64 {
	return (orig & 0_777000_777777) | ((new & 0_777) << 18)
}

// SetQ3 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetQ3(new uint64) *Word36 {
	*w = Word36(SetQ3(uint64(*w), new))
	return w
}

// SetQ3 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetQ3(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_000777) | ((new & 0_777) << 9)
}

// SetQ4 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetQ4(new uint64) *Word36 {
	*w = Word36(SetQ4(uint64(*w), new))
	return w
}

// SetQ4 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetQ4(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_777000) | (new & 0_777)
}

// SetS1 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetS1(new uint64) *Word36 {
	*w = Word36(SetS1(uint64(*w), new))
	return w
}

// SetS1 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetS1(orig uint64, new uint64) uint64 {
	return (orig & 0_007777_777777) | ((new & 077) << 30)
}

// SetS2 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetS2(new uint64) *Word36 {
	*w = Word36(SetS2(uint64(*w), new))
	return w
}

// SetS2 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetS2(orig uint64, new uint64) uint64 {
	return (orig & 0_770077_777777) | ((new & 077) << 24)
}

// SetS3 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetS3(new uint64) *Word36 {
	*w = Word36(SetS3(uint64(*w), new))
	return w
}

// SetS3 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetS3(orig uint64, new uint64) uint64 {
	return (orig & 0_777700_777777) | ((new & 077) << 18)
}

// SetS4 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetS4(new uint64) *Word36 {
	*w = Word36(SetS4(uint64(*w), new))
	return w
}

// SetS4 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetS4(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_007777) | ((new & 077) << 12)
}

// SetS5 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetS5(new uint64) *Word36 {
	*w = Word36(SetS5(uint64(*w), new))
	return w
}

// SetS5 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetS5(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_770077) | ((new & 077) << 6)
}

// SetS6 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetS6(new uint64) *Word36 {
	*w = Word36(SetS6(uint64(*w), new))
	return w
}

// SetS6 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetS6(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_777700) | (new & 077)
}

// SetT1 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetT1(new uint64) *Word36 {
	*w = Word36(SetT1(uint64(*w), new))
	return w
}

// SetT1 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetT1(orig uint64, new uint64) uint64 {
	return (orig & 0_000077_777777) | ((new & 0_7777) << 24)
}

// SetT2 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetT2(new uint64) *Word36 {
	*w = Word36(SetT2(uint64(*w), new))
	return w
}

// SetT2 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetT2(orig uint64, new uint64) uint64 {
	return (orig & 0_777700_007777) | ((new & 0_7777) << 12)
}

// SetT3 masks the new value into the appropriate partial-word of the original value
func (w *Word36) SetT3(new uint64) *Word36 {
	*w = Word36(SetT3(uint64(*w), new))
	return w
}

// SetT3 masks the new value into the appropriate partial-word of the original value, returning the expectedResult
func SetT3(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_770000) | (new & 0_7777)
}

// SetW sets the value to the given new value
func (w *Word36) SetW(new uint64) *Word36 {
	*w = Word36(new)
	return w
}

// SetW returns the input value, masked to 36 bits
// For consistency, we accept an orig value, but we dont really use it
func SetW(orig uint64, new uint64) uint64 {
	return new & Mask36
}

// Logical functions -----------------------------------------------------------

// And performs a logical AND with this value and the new value
func (w *Word36) And(new uint64) *Word36 {
	*w = *w & Word36(new)
	return w
}

// And returns the 36-bit logical AND of the two input values
func And(value1 uint64, value2 uint64) uint64 {
	return value1 & value2
}

// Not performs a logical inverse on this value
func (w *Word36) Not() *Word36 {
	*w ^= Mask36
	return w
}

// Not returns the 36-bit logical inverse of the input value
func Not(value uint64) uint64 {
	return value ^ Mask36
}

// Or performs a logical OR with this value and the new value
func (w *Word36) Or(new uint64) *Word36 {
	*w = *w | Word36(new)
	return w
}

// Or returns the 36-bit logical OR of the two input values
func Or(value1 uint64, value2 uint64) uint64 {
	return value1 | value2
}

// Xor performs a logical OR with this value and the new value
func (w *Word36) Xor(new uint64) *Word36 {
	*w = *w ^ Word36(new)
	return w
}

// Xor returns the 36-bit logical XOR of the two input values
func Xor(value1 uint64, value2 uint64) uint64 {
	return value1 ^ value2
}

// Arithmetic ------------------------------------------------------------------

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

// Divide performs integer division
func Divide(dividend []uint64, divisor uint64) (quotient uint64, remainder uint64, divByZero bool, overflow bool) {
	var div0 big.Int
	var div1 big.Int
	var bigDividend big.Int
	var bigDivisor big.Int
	var bigRemainder big.Int
	var bigQuotient big.Int

	div0.SetUint64(Magnitude(dividend[0]))
	div1.SetUint64(Magnitude(dividend[1]))
	bigDividend.Lsh(&div0, 36)
	bigDividend.Or(&bigDividend, &div1)

	if divisor == 0 {
		divByZero = true
		return
	}
	bigDivisor.SetUint64(Magnitude(divisor))

	bigQuotient.DivMod(&bigDividend, &bigDivisor, &bigRemainder)
	if bigQuotient.Cmp(bigMaxQuotient) == 1 {
		overflow = true
		return
	}

	quotient = bigQuotient.Uint64()
	remainder = bigRemainder.Uint64()
	return
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

// Magnitude converts this value to its absolute value
func (w *Word36) Magnitude() *Word36 {
	if w.IsNegative() {
		*w = *w ^ Mask36
	}
	return w
}

// Magnitude returns the absolute value of the given operand
func Magnitude(operand uint64) uint64 {
	if IsPositive(operand) {
		return operand
	} else {
		return Negate(operand)
	}
}

// Multiply returns the product of two 36-bit signed factors as a 72-bit signed integer
// packed into two uint64's.
func Multiply(factor1 uint64, factor2 uint64) []uint64 {
	var mag1 big.Int
	var mag2 big.Int
	mag1.SetUint64(Magnitude(factor1))
	mag2.SetUint64(Magnitude(factor2))
	neg := IsNegative(factor1) != IsNegative(factor2)

	mag1.Mul(&mag1, &mag2)
	mag2.Set(&mag1)

	magResult := []uint64{
		mag1.Rsh(&mag1, 36).Uint64(),
		mag2.And(&mag2, bigMask).Uint64(),
	}

	if neg {
		magResult = NegateDouble(magResult)
	}

	return magResult
}

// Negate converts this value to its additive inverse
func (w *Word36) Negate() *Word36 {
	*w ^= Mask36
	return w
}

// Negate returns the additive inverse of a given 36-bit signed value packed into a uint64
func Negate(op uint64) uint64 {
	return (op ^ NegativeZero) & NegativeZero
}

// Bit Shifting ----------------------------------------------------------------

// LeftShiftCircular shifts the 36-bit word to the left by the given count value,
// where every bit shifted out of bit 0 is end-around shifted into bit 35.
func (w *Word36) LeftShiftCircular(count uint64) *Word36 {
	*w = Word36(LeftShiftCircular(uint64(*w), count))
	return w
}

func LeftShiftCircular(operand uint64, count uint64) uint64 {
	result := operand
	count %= 36
	for count > 0 {
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
func (w *Word36) LeftShiftLogical(count uint64) *Word36 {
	*w = Word36(LeftShiftLogical(uint64(*w), count))
	return w
}

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

// RightShiftAlgebraic shifts the 72-bit word to the left by the given count value,
// Bits shifted out of bit 35 are lost while bit 0 is propagated to the right.
func (w *Word36) RightShiftAlgebraic(count uint64) *Word36 {
	*w = Word36(RightShiftAlgebraic(uint64(*w), count))
	return w
}

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
func (w *Word36) RightShiftCircular(count uint64) *Word36 {
	*w = Word36(RightShiftCircular(uint64(*w), count))
	return w
}

func RightShiftCircular(operand uint64, count uint64) uint64 {
	result := operand
	count %= 36
	for count > 0 {
		shift := uint64Min(27, count)
		mask := uint64(1<<shift) - 1
		partial := (result & mask) << 36
		result = (partial | result) >> shift
		count -= shift
	}

	return result
}

// RightShiftLogical shifts the 36-bit word to the left by the given count value.
// Bits shifted out of bit 35 are lost, and zeroes are shift into bit 0.
func (w *Word36) RightShiftLogical(count uint64) *Word36 {
	*w = Word36(RightShiftLogical(uint64(*w), count))
	return w
}

func RightShiftLogical(operand uint64, count uint64) uint64 {
	result := operand
	if count >= 36 {
		result = PositiveZero
	} else {
		result >>= count
	}

	return result
}

// Double-word functions -------------------------------------------------------

func AddDouble(operand1 []uint64, operand2 []uint64) []uint64 {
	if IsDoubleNegativeZero(operand1) && IsDoubleNegativeZero(operand2) {
		return DoubleNegativeZero
	} else {
		var op1MSW big.Int
		var op1LSW big.Int
		var bigAddend1 big.Int
		var op2MSW big.Int
		var op2LSW big.Int
		var bigAddend2 big.Int
		var bigSum big.Int
		var bigSumCopy big.Int

		op1Mag := MagnitudeDouble(operand1)
		op1MSW.SetUint64(op1Mag[0])
		op1LSW.SetUint64(op1Mag[1])
		bigAddend1.Lsh(&op1MSW, 36)
		bigAddend1.Or(&bigAddend1, &op1LSW)
		if IsNegativeDouble(operand1) {
			bigAddend1.Neg(&bigAddend1)
		}

		op2Mag := MagnitudeDouble(operand2)
		op2MSW.SetUint64(op2Mag[0])
		op2LSW.SetUint64(op2Mag[1])
		bigAddend2.Lsh(&op2MSW, 36)
		bigAddend2.Or(&bigAddend2, &op2LSW)
		if IsNegativeDouble(operand2) {
			bigAddend2.Neg(&bigAddend2)
		}

		bigSum.Add(&bigAddend1, &bigAddend2)
		negSumFlag := bigSum.Cmp(big.NewInt(0)) < 0
		bigSum.Abs(&bigSum)
		bigSumCopy.Set(&bigSum)

		sum := []uint64{
			bigSum.Rsh(&bigSum, 36).Uint64(),
			bigSumCopy.And(&bigSumCopy, bigMask).Uint64(),
		}
		if negSumFlag {
			sum = NegateDouble(sum)
		}
		return sum
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

func IsNegativeDouble(value []uint64) bool {
	return (value[0] & 0_400000_000000) != 0
}

func IsPositiveDouble(value []uint64) bool {
	return (value[0] & 0_400000_000000) == 0
}

func IsDoubleZero(operand []uint64) bool {
	return IsZero(operand[0]) && operand[0] == operand[1]
}

func IsDoubleNegativeZero(operand []uint64) bool {
	return operand[0] == NegativeZero && operand[1] == NegativeZero
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

// MagnitudeDouble returns the absolute value of the given 72-bit operand
func MagnitudeDouble(operand []uint64) []uint64 {
	if IsPositiveDouble(operand) {
		return operand
	} else {
		return NegateDouble(operand)
	}
}

// NegateDouble returns the additive inverse of the given 72-bit signed value packed into uint64's
func NegateDouble(op []uint64) []uint64 {
	return []uint64{
		Negate(op[0]),
		Negate(op[1]),
	}
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

// Miscellaneous functions -----------------------------------------------------

// CountBits returns a count of the number of bits set in the right-most 36-bits of the given value.
func (w *Word36) CountBits() uint64 {
	return CountBits(uint64(*w))
}

func CountBits(value uint64) uint64 {
	v := value
	var count uint64
	for v > 0 {
		if v&01 == 01 {
			count++
		}
		v >>= 1
	}
	return count
}

// EliminateNegativeZero returns the given value, converting it to positive zero if it is negative zero.
func (w *Word36) EliminateNegativeZero() *Word36 {
	if *w == NegativeZero {
		*w = PositiveZero
	}
	return w
}

func EliminateNegativeZero(value uint64) uint64 {
	if value == NegativeZero {
		return PositiveZero
	}
	return value
}

// IsNegative returns true if the given value represents a negative 36-bit value
func (w *Word36) IsNegative() bool {
	return IsNegative(uint64(*w))
}

func IsNegative(value uint64) bool {
	return (value & 0_400000_000000) != 0
}

// IsPositive returns true if the given value represents a positive 36-bit value
func (w *Word36) IsPositive() bool {
	return IsPositive(uint64(*w))
}

func IsPositive(value uint64) bool {
	return (value & 0_400000_000000) == 0
}

// IsZero returns true if the given value represents a positive or negative 36-bit zero value
func (w *Word36) IsZero() bool {
	return IsZero(uint64(*w))
}

func IsZero(value uint64) bool {
	return value == PositiveZero || value == NegativeZero
}

// ExtractPartialWord pulls the partial word indicated by the partialWordIndicator and the quarterWordMode flag
// from the given 36-bit source value.
func ExtractPartialWord(source uint64, partialWordIndicator uint, quarterWordMode bool) uint64 {
	switch partialWordIndicator {
	case JFieldW:
		return GetW(source)
	case JFieldH2:
		return GetH2(source)
	case JFieldH1:
		return GetH1(source)
	case JFieldXH2:
		return GetXH2(source)
	case JFieldXH1: // XH1 or Q2
		if quarterWordMode {
			return GetQ2(source)
		} else {
			return GetXH1(source)
		}
	case JFieldT3: // T3 or Q4
		if quarterWordMode {
			return GetQ4(source)
		} else {
			return GetXT3(source)
		}
	case JFieldT2: // T2 or Q3
		if quarterWordMode {
			return GetQ3(source)
		} else {
			return GetXT2(source)
		}
	case JFieldT1: // T1 or Q1
		if quarterWordMode {
			return GetQ1(source)
		} else {
			return GetXT1(source)
		}
	case JFieldS6:
		return GetS6(source)
	case JFieldS5:
		return GetS5(source)
	case JFieldS4:
		return GetS4(source)
	case JFieldS3:
		return GetS3(source)
	case JFieldS2:
		return GetS2(source)
	case JFieldS1:
		return GetS1(source)
	}

	return source
}

// InjectPartialWord creates a value comprised of an original value and a new value inserted there-in under j-field control.
func InjectPartialWord(originalValue uint64, newValue uint64, jField uint, quarterWordMode bool) uint64 {
	switch jField {
	case JFieldW:
		return newValue
	case JFieldH2:
		return SetH2(originalValue, newValue)
	case JFieldXH2:
		return SetH2(originalValue, newValue)
	case JFieldH1:
		return SetH1(originalValue, newValue)
	case JFieldXH1: // XH1 or Q2
		if quarterWordMode {
			return SetQ2(originalValue, newValue)
		} else {
			return SetH1(originalValue, newValue)
		}
	case JFieldT3: // T3 or Q4
		if quarterWordMode {
			return SetQ4(originalValue, newValue)
		} else {
			return SetT3(originalValue, newValue)
		}
	case JFieldT2: // T2 or Q3
		if quarterWordMode {
			return SetQ3(originalValue, newValue)
		} else {
			return SetT2(originalValue, newValue)
		}
	case JFieldT1: // T1 or Q1
		if quarterWordMode {
			return SetQ1(originalValue, newValue)
		} else {
			return SetT1(originalValue, newValue)
		}
	case JFieldS6:
		return SetS6(originalValue, newValue)
	case JFieldS5:
		return SetS5(originalValue, newValue)
	case JFieldS4:
		return SetS4(originalValue, newValue)
	case JFieldS3:
		return SetS3(originalValue, newValue)
	case JFieldS2:
		return SetS2(originalValue, newValue)
	case JFieldS1:
		return SetS1(originalValue, newValue)
	}

	return originalValue
}
