package types

import (
	"fmt"
	"log"
)

// Wraps a regular 64-bit unsigned int with a type name which we must always use
// when we refer to a 36-bit word in Kalehla world.
// The top 28 bits are *always* expected to be zero.

type Word36 uint64

const PositiveZero = 0
const NegativeZero = 0_777777_777777

func (w *Word36) EliminateNegativeZero() *Word36 {
	if *w == NegativeZero {
		*w = 0
	}
	return w
}

func (w *Word36) FromBytesToAscii(inp []byte) *Word36 {
	var value uint
	for bx := 0; bx < 4; bx++ {
		value <<= 9
		if bx < len(inp) {
			value |= uint(inp[bx])
		} else {
			value |= ' '
		}
	}
	return w.SetW(uint64(value))
}

func (w *Word36) FromBytesToFieldata(inp []byte) *Word36 {
	var value uint
	for bx := 0; bx < 6; bx++ {
		value <<= 6
		if bx < len(inp) {
			value |= uint(FieldataFromAscii[inp[bx]&0177])
		} else {
			value |= 05
		}
	}
	return w.SetW(uint64(value))
}

func (w *Word36) FromStringToAscii(inp []byte) *Word36 {
	temp := fmt.Sprintf("-%4s", inp)
	w.FromBytesToAscii([]byte(temp)[0:4])
	return w
}

func (w *Word36) FromStringToFieldata(inp []byte) *Word36 {
	temp := fmt.Sprintf("-%6s", inp)
	w.FromBytesToFieldata([]byte(temp)[0:6])
	return w
}

func (w *Word36) GetW() uint64 {
	return uint64(*w)
}

func (w *Word36) GetH1() uint {
	return uint(*w >> 18)
}

func (w *Word36) GetH2() uint {
	return uint(*w & 0_777777)
}

func (w *Word36) GetQ1() uint {
	return uint(*w >> 27)
}

func (w *Word36) GetQ2() uint {
	return uint((*w >> 18) & 0777)
}

func (w *Word36) GetQ3() uint {
	return uint((*w >> 9) & 0777)
}

func (w *Word36) GetQ4() uint {
	return uint(*w & 0777)
}

func (w *Word36) GetS1() uint {
	return uint(*w >> 30)
}

func (w *Word36) GetS2() uint {
	return uint((*w >> 24) & 077)
}

func (w *Word36) GetS3() uint {
	return uint((*w >> 18) & 077)
}

func (w *Word36) GetS4() uint {
	return uint((*w >> 12) & 077)
}

func (w *Word36) GetS5() uint {
	return uint((*w >> 6) & 077)
}

func (w *Word36) GetS6() uint {
	return uint(*w & 077)
}

func (w *Word36) IsNegative() bool {
	return ((*w) & 0_400000_000000) != 0
}

func (w *Word36) IsZero() bool {
	return (*w == 0) || (*w == 0_777777_777777)
}

func (w *Word36) SetW(op uint64) *Word36 {
	*w = Word36(op) & 0777777777777
	return w
}

func (w *Word36) SetH1(op uint) *Word36 {
	*w = (*w & 0777777) | ((Word36(op) & 0777777) << 18)
	return w
}

func (w *Word36) SetH2(op uint) *Word36 {
	*w = (*w & 0777777000000) | (Word36(op) & 0777777)
	return w
}

func (w *Word36) SetS1(op uint) *Word36 {
	*w = (*w & 0007777777777) | ((Word36(op) & 077) << 30)
	return w
}

func (w *Word36) SetS2(op uint) *Word36 {
	*w = (*w & 0770077777777) | ((Word36(op) & 077) << 24)
	return w
}

func (w *Word36) SetS3(op uint) *Word36 {
	*w = (*w & 0777700777777) | ((Word36(op) & 077) << 18)
	return w
}

func (w *Word36) SetS4(op uint) *Word36 {
	*w = (*w & 0777777007777) | ((Word36(op) & 077) << 12)
	return w
}

func (w *Word36) SetS5(op uint) *Word36 {
	*w = (*w & 0777777770077) | ((Word36(op) & 077) << 6)
	return w
}

func (w *Word36) SetS6(op uint) *Word36 {
	*w = (*w & 0777777777700) | (Word36(op) & 077)
	return w
}

func (w *Word36) ToStringAsAscii() string {
	tempVal := uint64(*w)
	temp := make([]byte, 4)
	temp[3] = byte(tempVal & 0377)
	tempVal >>= 9
	temp[2] = byte(tempVal & 0377)
	tempVal >>= 9
	temp[1] = byte(tempVal & 0377)
	tempVal >>= 9
	temp[0] = byte(tempVal)
	return string(temp)
}

func (w *Word36) ToStringAsFieldata() string {
	tempVal := uint64(*w)
	temp := make([]byte, 8)
	temp[5] = AsciiFromFieldata[tempVal&077]
	tempVal >>= 6
	temp[4] = AsciiFromFieldata[tempVal&077]
	tempVal >>= 6
	temp[3] = AsciiFromFieldata[tempVal&077]
	tempVal >>= 6
	temp[2] = AsciiFromFieldata[tempVal&077]
	tempVal >>= 6
	temp[1] = AsciiFromFieldata[tempVal&077]
	tempVal >>= 6
	temp[0] = AsciiFromFieldata[tempVal&077]
	return string(temp)
}

func (w *Word36) And(op uint) {
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

// AddSimple adds two ones-complement values
func AddSimple(operand1 Word36, operand2 Word36) Word36 {
	if (operand1 == NegativeZero) && (operand2 == NegativeZero) {
		return NegativeZero
	}

	native1 := ConvertOnesComplementToNative(operand1)
	native2 := ConvertOnesComplementToNative(operand2)
	return ConvertNativeToOnesComplement(native1 + native2)
}

// ConvertNativeToOnesComplement converts 2's complement to 1's complement
func ConvertNativeToOnesComplement(operand int64) Word36 {
	if operand < 0 {
		op := (-operand) & 0377777777777
		return Word36(op ^ NegativeZero)
	} else {
		return Word36(operand & NegativeZero)
	}
}

// ConvertOnesComplementToNative converts 1's complement to 2's complement
func ConvertOnesComplementToNative(operand Word36) int64 {
	if operand.IsNegative() {
		return -(int64(operand) ^ NegativeZero)
	} else {
		return int64(operand)
	}
}

// GetSignExtended12 sign-extends an 12-bit value to 36 bits
func GetSignExtended12(value uint64) Word36 {
	if (value & 04000) == 0 {
		return Word36(value)
	} else {
		return Word36(value | 0_777777_770000)
	}
}

// GetSignExtended18 sign-extends an 18-bit value to 36 bits
func GetSignExtended18(value uint64) Word36 {
	if (value & 0_400000) == 0 {
		return Word36(value)
	} else {
		return Word36(value | 0_777777_000000)
	}
}

// GetSignExtended24 sign-extends a 24-bit value to 36 bits
func GetSignExtended24(value uint64) Word36 {
	if (value & 040_000000) == 0 {
		return Word36(value)
	} else {
		return Word36(value | 0_777700_000000)
	}
}

// Negate returns the additive inverse of the given ones-complement value, in ones-complement
// Note that this is the same thing as taking the logical not of the operand
func Negate(value Word36) Word36 {
	return value ^ NegativeZero
}

func FromStringToAsciiWords(inp string, buffer []Word36) {
	arr := []byte(fmt.Sprintf("%-*s", len(buffer)*4, inp))
	ax := 0
	for wx := 0; wx < len(buffer); wx++ {
		buffer[wx].FromStringToAscii(arr[ax : ax+4])
		ax += 4
	}
}

func FromStringToFieldataWords(inp string, buffer []Word36) {
	arr := []byte(fmt.Sprintf("%-*s", len(buffer)*6, inp))
	ax := 0
	for wx := 0; wx < len(buffer); wx++ {
		buffer[wx].FromStringToFieldata(arr[ax : ax+6])
		ax += 6
	}
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
	sl := len(source)
	if sl%9 != 0 {
		log.Panic("source buffer length is not a multiple of 9 bytes")
	}

	if sl*2/9 > len(destination) {
		log.Panic("destination buffer insufficient size")
	}

	dx := 0
	for sx := 0; sx < sl; {
		w0 := Word36(source[sx])
		w0 <<= 8
		sx += 1
		w0 |= Word36(source[sx])
		w0 <<= 8
		sx += 1
		w0 |= Word36(source[sx])
		w0 <<= 8
		sx += 1
		w0 |= Word36(source[sx])
		w0 <<= 8
		sx += 1
		w0 |= Word36(source[sx] >> 4)

		w1 := Word36(source[sx] & 0xF)
		w1 <<= 4
		sx += 1
		w1 |= Word36(source[sx])
		w1 <<= 8
		sx += 1
		w1 |= Word36(source[sx])
		w1 <<= 8
		sx += 1
		w1 |= Word36(source[sx])
		w1 <<= 8
		sx += 1
		w1 |= Word36(source[sx])

		destination[dx] = w0
		dx += 1
		destination[dx] = w1
		dx += 1
	}
}
