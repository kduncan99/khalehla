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

func IsNegative(value uint64) bool {
	return (value & 0_400000_000000) != 0
}

func (w *Word36) IsZero() bool {
	return (*w == PositiveZero) || (*w == NegativeZero)
}

func IsZero(value uint64) bool {
	return (value == PositiveZero) || (value == NegativeZero)
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
