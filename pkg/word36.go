// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

import (
	"fmt"
	"log"
)

// Wraps a regular 64-bit unsigned int with a type name which we must always use
// when we refer to a 36-bit word in Kalehla world.
// The top 28 bits are *always* expected to be zero.

type Word36 uint64

const PositiveOne = 01
const PositiveZero = 0
const NegativeOne = 0_777777_777776
const NegativeZero = 0_777777_777777

func (w *Word36) CountBits() uint64 {
	return CountBits(w.GetW())
}

func CountBits(value uint64) uint64 {
	v := value & NegativeZero
	var count uint64
	for v > 0 {
		if v&01 == 01 {
			count++
		}
		v >>= 1
	}
	return count
}

func (w *Word36) EliminateNegativeZero() *Word36 {
	if *w == NegativeZero {
		*w = 0
	}
	return w
}

// FromBytesToAscii converts an array of exactly 4 bytes to a Word36 containing 4 ASCII characters.
// The input array should be LJSF, but we do not enforce that.
func (w *Word36) FromBytesToAscii(inp []byte) *Word36 {
	var value uint
	for bx := 0; bx < 4; bx++ {
		value <<= 9
		value |= uint(inp[bx])
	}
	return w.SetW(uint64(value))
}

// FromBytesToAscii converts an input array of a length which is at least the length of the output array
// multiplied by six, to the output array of Word36 structs such that each byte is wrapped in a consecutive
// ASCII character in the corresponding output word. The input should be LJSF.
func FromBytesToAscii(input []byte, output []Word36) {
	bx := 0
	wx := 0
	for wx < len(output) {
		output[wx].FromBytesToAscii(input[bx : bx+4])
		bx += 4
		wx++
	}
}

// FromBytesToFieldata converts an array of exactly 6 bytes to a Word36 containing 6 Fieldata characters.
// The input array should be LJSF, but we do not enforce that.
func (w *Word36) FromBytesToFieldata(inp []byte) *Word36 {
	var value uint
	for bx := 0; bx < 6; bx++ {
		value <<= 6
		value |= uint(FieldataFromAscii[inp[bx]&0177])
	}
	return w.SetW(uint64(value))
}

// FromBytesToFieldata converts an input array of zero or more bytes, into an output array of a fixed size
// where the output words consist of 6 Fieldata characters. The entire output array is LJSF.
func FromBytesToFieldata(input []byte, output []Word36) {
	bx := 0
	wx := 0
	for wx < len(output) {
		output[wx].FromBytesToFieldata(input[bx : bx+6])
		bx += 6
		wx++
	}
}

// FromStringToAscii converts a string of up to 4 characters to a Word36 containing 4 ASCII characters LJSF.
func (w *Word36) FromStringToAscii(inp string) *Word36 {
	temp := fmt.Sprintf("%-4s", inp)
	w.FromBytesToAscii([]byte(temp))
	return w
}

// FromStringToFieldata converts an array of up to 6 characters to a Word36 containing 6 Fieldata characters LJSF.
func (w *Word36) FromStringToFieldata(inp string) *Word36 {
	temp := fmt.Sprintf("%-6s", inp)
	w.FromBytesToFieldata([]byte(temp))
	return w
}

// FromStringToAscii converts a string of any number of characters to a Word36 array where-in the output array
// consists of Word36 structs of 4 ASCII bytes per struct. The entire output is LJSF.
func FromStringToAscii(input string, output []Word36) {
	tempLen := len(output) * 4
	if tempLen&03 != 0 {
		tempLen = (tempLen &^ 03) + 4
	}
	temp := fmt.Sprintf("%-*s", tempLen, input)
	FromBytesToAscii([]byte(temp), output)
}

// FromStringToFieldata converts a string of any number of characters to a Word36 array where-in the output array
// consists of Word36 structs of 6 Fieldata bytes per struct. The entire output is LJSF.
func FromStringToFieldata(input string, output []Word36) {
	tempLen := len(output) * 4
	tempMod := tempLen % 6
	if tempMod != 0 {
		tempLen += 6 - tempMod
	}
	temp := fmt.Sprintf("%-*s", tempLen, input)
	FromBytesToFieldata([]byte(temp), output)
}

// NewFromStringToAscii converts a string of any number of characters to a Word36 array where-in the output array
// consists of Word36 structs of 4 ASCII bytes per struct. The entire output is LJSF.
func NewFromStringToAscii(input string, outputWordCount int) []Word36 {
	temp := fmt.Sprintf("%-*s", outputWordCount*4, input)
	result := make([]Word36, outputWordCount)
	FromBytesToAscii([]byte(temp), result)
	return result
}

// NewFromStringToFieldata converts a string of any number of characters to a Word36 array where-in the output array
// consists of Word36 structs of 6 Fieldata bytes per struct. The entire output is LJSF.
func NewFromStringToFieldata(input string, outputWordCount int) []Word36 {
	temp := fmt.Sprintf("%-*s", outputWordCount*6, input)
	result := make([]Word36, outputWordCount)
	FromBytesToFieldata([]byte(temp), result)
	return result
}

func (w *Word36) GetW() uint64 {
	return uint64(*w)
}

func GetW(value uint64) uint64 {
	return value & 0_777777_777777
}

func (w *Word36) GetH1() uint64 {
	return uint64(*w >> 18)
}

func GetH1(value uint64) uint64 {
	return (value >> 18) & 0777777
}

func (w *Word36) GetH2() uint64 {
	return uint64(*w & 0_777777)
}

func GetH2(value uint64) uint64 {
	return value & 0777777
}

func (w *Word36) GetQ1() uint64 {
	return uint64(*w >> 27)
}

func GetQ1(value uint64) uint64 {
	return (value >> 27) & 0777
}

func (w *Word36) GetQ2() uint64 {
	return uint64((*w >> 18) & 0777)
}

func GetQ2(value uint64) uint64 {
	return (value >> 18) & 0777
}

func (w *Word36) GetQ3() uint64 {
	return uint64((*w >> 9) & 0777)
}

func GetQ3(value uint64) uint64 {
	return (value >> 9) & 0777
}

func (w *Word36) GetQ4() uint64 {
	return uint64(*w & 0777)
}

func GetQ4(value uint64) uint64 {
	return value & 0777
}

func (w *Word36) GetS1() uint64 {
	return uint64(*w >> 30)
}

func GetS1(value uint64) uint64 {
	return (value >> 30) & 077
}

func (w *Word36) GetS2() uint64 {
	return uint64((*w >> 24) & 077)
}

func GetS2(value uint64) uint64 {
	return (value >> 24) & 077
}

func (w *Word36) GetS3() uint64 {
	return uint64((*w >> 18) & 077)
}

func GetS3(value uint64) uint64 {
	return (value >> 18) & 077
}

func (w *Word36) GetS4() uint64 {
	return uint64((*w >> 12) & 077)
}

func GetS4(value uint64) uint64 {
	return (value >> 12) & 077
}

func (w *Word36) GetS5() uint64 {
	return uint64((*w >> 6) & 077)
}

func GetS5(value uint64) uint64 {
	return (value >> 6) & 077
}

func (w *Word36) GetS6() uint64 {
	return uint64(*w & 077)
}

func GetS6(value uint64) uint64 {
	return value & 077
}

func (w *Word36) GetXH1() uint64 {
	value := uint64((*w >> 18) & 0_777777)
	if (value & 0_400000) != 0 {
		value |= 0_777777_000000
	}
	return value
}

func GetXH1(value uint64) uint64 {
	res := (value >> 18) & 0_777777
	if (res & 0_400000) != 0 {
		res |= 0_777777_000000
	}
	return res
}

func (w *Word36) GetXH2() uint64 {
	value := uint64(*w & 0_777777)
	if (value & 0_400000) != 0 {
		value |= 0_777777_000000
	}
	return value
}

func GetXH2(value uint64) uint64 {
	res := value & 0_777777
	if (res & 0_400000) != 0 {
		res |= 0_777777_000000
	}
	return res
}

func (w *Word36) GetXT1() uint64 {
	value := uint64((*w >> 24) & 0_7777)
	if (value & 004000) != 0 {
		value |= 0_777777_770000
	}
	return value
}

func GetXT1(value uint64) uint64 {
	res := (value >> 24) & 0_7777
	if (res & 004000) != 0 {
		res |= 0_777777_770000
	}
	return res
}

func (w *Word36) GetXT2() uint64 {
	value := uint64((*w >> 12) & 0_7777)
	if (value & 004000) != 0 {
		value |= 0_777777_770000
	}
	return value
}

func GetXT2(value uint64) uint64 {
	res := (value >> 12) & 0_7777
	if (res & 004000) != 0 {
		res |= 0_777777_770000
	}
	return res
}

func (w *Word36) GetXT3() uint64 {
	value := uint64(*w & 0_7777)
	if (value & 004000) != 0 {
		value |= 0_777777_770000
	}
	return value
}

func GetXT3(value uint64) uint64 {
	res := value & 0_7777
	if (res & 004000) != 0 {
		res |= 0_777777_770000
	}
	return res
}

func (w *Word36) IsNegative() bool {
	return ((*w) & 0_400000_000000) != 0
}

func (w *Word36) IsPositive() bool {
	return ((*w) & 0_400000_000000) == 0
}

func (w *Word36) IsZero() bool {
	return IsZero(uint64(*w))
}

func (w *Word36) SetW(op uint64) *Word36 {
	*w = Word36(op) & 0777777777777
	return w
}

func (w *Word36) SetH1(op uint64) *Word36 {
	*w = (*w & 0_000000_777777) | ((Word36(op) & 0_777777) << 18)
	return w
}

func SetH1(orig uint64, new uint64) uint64 {
	return (orig & 0_777777) | ((new & 0_777777) << 18)
}

func (w *Word36) SetH2(op uint64) *Word36 {
	*w = (*w & 0_777777_000000) | (Word36(op) & 0_777777)
	return w
}

func SetH2(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_000000) | (new & 0_777777)
}

func (w *Word36) SetQ1(op uint64) *Word36 {
	*w = (*w & 0_000777_777777) | ((Word36(op) & 0777) << 27)
	return w
}

func SetQ1(orig uint64, new uint64) uint64 {
	return (orig & 0_000777_777777) | ((new & 0_777) << 27)
}

func (w *Word36) SetQ2(op uint64) *Word36 {
	*w = (*w & 0_777000_777777) | ((Word36(op) & 0777) << 18)
	return w
}

func SetQ2(orig uint64, new uint64) uint64 {
	return (orig & 0_777000_777777) | ((new & 0_777) << 18)
}

func (w *Word36) SetQ3(op uint64) *Word36 {
	*w = (*w & 0_777777_000777) | ((Word36(op) & 0777) << 9)
	return w
}

func SetQ3(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_000777) | ((new & 0_777) << 9)
}

func (w *Word36) SetQ4(op uint64) *Word36 {
	*w = (*w & 0_777777_777000) | (Word36(op) & 0777)
	return w
}

func SetQ4(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_777000) | (new & 0_777)
}

func (w *Word36) SetS1(op uint64) *Word36 {
	*w = (*w & 0_007777_777777) | ((Word36(op) & 077) << 30)
	return w
}

func SetS1(orig uint64, new uint64) uint64 {
	return (orig & 0_007777_777777) | ((new & 077) << 30)
}

func (w *Word36) SetS2(op uint64) *Word36 {
	*w = (*w & 0_770077_777777) | ((Word36(op) & 077) << 24)
	return w
}

func SetS2(orig uint64, new uint64) uint64 {
	return (orig & 0_770077_777777) | ((new & 077) << 24)
}

func (w *Word36) SetS3(op uint64) *Word36 {
	*w = (*w & 0_777700_777777) | ((Word36(op) & 077) << 18)
	return w
}

func SetS3(orig uint64, new uint64) uint64 {
	return (orig & 0_777700_777777) | ((new & 077) << 18)
}

func (w *Word36) SetS4(op uint64) *Word36 {
	*w = (*w & 0_777777_007777) | ((Word36(op) & 077) << 12)
	return w
}

func SetS4(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_007777) | ((new & 077) << 12)
}

func (w *Word36) SetS5(op uint64) *Word36 {
	*w = (*w & 0_777777_770077) | ((Word36(op) & 077) << 6)
	return w
}

func SetS5(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_770077) | ((new & 077) << 6)
}

func (w *Word36) SetS6(op uint64) *Word36 {
	*w = (*w & 0_777777_777700) | (Word36(op) & 077)
	return w
}

func SetS6(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_777700) | (new & 077)
}

func (w *Word36) SetT1(op uint64) *Word36 {
	*w = (*w & 0_000077_777777) | ((Word36(op) & 07777) << 24)
	return w
}

func SetT1(orig uint64, new uint64) uint64 {
	return (orig & 0_000077_777777) | ((new & 0_7777) << 24)
}

func (w *Word36) SetT2(op uint64) *Word36 {
	*w = (*w & 0_777700_007777) | ((Word36(op) & 07777) << 12)
	return w
}

func SetT2(orig uint64, new uint64) uint64 {
	return (orig & 0_777700_007777) | ((new & 0_7777) << 12)
}

func (w *Word36) SetT3(op uint64) *Word36 {
	*w = (*w & 0_777777_770000) | (Word36(op) & 07777)
	return w
}

func SetT3(orig uint64, new uint64) uint64 {
	return (orig & 0_777777_770000) | (new & 0_7777)
}

func (w *Word36) ShiftLeftCircular(count uint64) {
	*w = Word36(ShiftLeftCircular(uint64(*w), count))
}

func ShiftLeftCircular(orig uint64, count uint64) uint64 {
	mod := count % 36
	result := orig
	if mod >= 18 {
		result <<= 18
		mod -= 18
		result = (result & 0_777777_000000) | ((result & 0_777777_000000_000000) >> 36)
	}

	for mod > 0 {
		result <<= 1
		if result&01_000000_000000 > 0 {
			result |= 01
		}
		mod -= 1
	}

	return result & NegativeZero
}

// DumpWord36Buffer displays the indicated buffer to os.Stdout
func DumpWord36Buffer(buffer []Word36, wordsPerLine int) {
	for wx := 0; wx < len(buffer); wx += wordsPerLine {
		offsetStr := fmt.Sprintf("%06o", wx)

		octalStr := ""
		fdStr := ""
		asciiStr := ""
		for wy := 0; wy < wordsPerLine; wy++ {
			wz := wx + wy
			if wz < len(buffer) {
				word := buffer[wz]
				octalStr += word.ToStringAsOctal() + " "
				fdStr += word.ToStringAsFieldata() + " "
				asciiStr += word.ToStringAsAsciiWithReplacementChar('.') + " "
			} else {
				octalStr += "             "
				fdStr += "       "
				asciiStr += "     "
			}
		}

		fmt.Printf("%s:  %s %s %s\n", offsetStr, octalStr, fdStr, asciiStr)
	}
}

// ToStringAsAscii converts one Word36 presumed to contain ASCII characters, to a 4-character string
func (w *Word36) ToStringAsAscii() string {
	temp := make([]byte, 4)
	temp[0] = byte(w.GetQ1())
	temp[1] = byte(w.GetQ2())
	temp[2] = byte(w.GetQ3())
	temp[3] = byte(w.GetQ4())
	return string(temp)
}

// ToStringAsAsciiWithReplacementChar converts one Word36 presumed to contain ASCII characters, to a 4-character string
// replacing any non-printing bytes with the given byte
func (w *Word36) ToStringAsAsciiWithReplacementChar(char byte) string {
	temp := make([]byte, 4)
	temp[0] = byte(w.GetQ1())
	temp[1] = byte(w.GetQ2())
	temp[2] = byte(w.GetQ3())
	temp[3] = byte(w.GetQ4())

	for tx := 0; tx < 4; tx++ {
		if temp[tx] < 32 || temp[tx] >= 127 {
			temp[tx] = char
		}
	}

	return string(temp)
}

// ToStringAsFieldata converts one Word36 presumed ton contain FIELDATA characters, to a 6-character string
func (w *Word36) ToStringAsFieldata() string {
	temp := make([]byte, 6)
	temp[0] = AsciiFromFieldata[w.GetS1()]
	temp[1] = AsciiFromFieldata[w.GetS2()]
	temp[2] = AsciiFromFieldata[w.GetS3()]
	temp[3] = AsciiFromFieldata[w.GetS4()]
	temp[4] = AsciiFromFieldata[w.GetS5()]
	temp[5] = AsciiFromFieldata[w.GetS6()]

	return string(temp)
}

func (w *Word36) ToStringAsOctal() string {
	return fmt.Sprintf("%012o", uint64(*w))
}

func (w *Word36) And(op uint64) {
	*w &= Word36(op & 0777777777777)
}

func (w *Word36) Not() {
	*w ^= 0777777777777
}

func (w *Word36) Or(op uint64) {
	*w |= Word36(op & 0777777777777)
}

func (w *Word36) Xor(op uint64) {
	*w ^= Word36(op & 0777777777777)
}

// PackWord36 packs pairs of word36 structs into 9-byte sequences
func PackWord36(source []Word36, destination []byte) {
	sl := len(source)
	if sl%1 != 0 {
		log.Panic("source buffer does not contain an even number of words")
	}

	if sl*9/2 > len(destination) {
		log.Panic("destination buffer insufficient size")
	}

	dx := 0
	for wx := 0; wx < sl; wx += 2 {
		val0 := source[wx].GetW()
		val1 := source[wx+1].GetW()

		destination[dx+8] = byte(val1)
		val1 >>= 8
		destination[dx+7] = byte(val1)
		val1 >>= 8
		destination[dx+6] = byte(val1)
		val1 >>= 8
		destination[dx+5] = byte(val1)
		val1 >>= 8
		destination[dx+4] = byte(val0<<4) | byte(val1)
		val0 >>= 4
		destination[dx+3] = byte(val0)
		val0 >>= 8
		destination[dx+2] = byte(val0)
		val0 >>= 8
		destination[dx+1] = byte(val0)
		val0 >>= 8
		destination[dx] = byte(val0)

		dx += 9
	}
}

// UnpackWord36 unpacks 9-byte groups into pairs of Word36 structs
func UnpackWord36(source []byte, destination []Word36) {
	srcLen := len(source)
	if srcLen%9 != 0 {
		log.Panic("source buffer length is not a multiple of 9 bytes")
	}

	if srcLen*2/9 > len(destination) {
		log.Panic("destination buffer insufficient size")
	}

	dx := 0
	for sx := 0; sx < srcLen; {
		w0 := uint64(source[sx])
		w0 <<= 8
		sx += 1
		w0 |= uint64(source[sx])
		w0 <<= 8
		sx += 1
		w0 |= uint64(source[sx])
		w0 <<= 8
		sx += 1
		w0 |= uint64(source[sx])
		w0 <<= 4
		sx += 1
		w0 |= uint64(source[sx] >> 4)

		w1 := uint64(source[sx] & 0xF)
		w1 <<= 8
		sx += 1
		w1 |= uint64(source[sx])
		w1 <<= 8
		sx += 1
		w1 |= uint64(source[sx])
		w1 <<= 8
		sx += 1
		w1 |= uint64(source[sx])
		w1 <<= 8
		sx += 1
		w1 |= uint64(source[sx])
		sx += 1

		destination[dx].SetW(w0)
		destination[dx+1].SetW(w1)
		dx += 2
	}
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
