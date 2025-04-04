// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package common

import (
	"fmt"
)

var RawBytesPerBlock = map[uint]uint{
	28:   28 * 8,
	56:   56 * 8,
	112:  112 * 8,
	224:  224 * 8,
	448:  448 * 8,
	896:  896 * 8,
	1792: 1792 * 8,
}

var AsciiFromFieldata = []byte{
	'@', '[', ']', '#', '^', ' ', 'A', 'B',
	'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R',
	'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	')', '-', '+', '<', '=', '>', '&', '$',
	'*', '(', '%', ':', '?', '!', ',', '\\',
	'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', '\'', ';', '/', '.', '"', '_',
}

var FieldataFromAscii = []int{
	005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005,
	005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005, 005,
	005, 055, 076, 003, 047, 052, 046, 072, 051, 040, 050, 042, 056, 041, 075, 074,
	060, 061, 062, 063, 064, 065, 066, 067, 070, 071, 053, 073, 043, 044, 045, 054,
	000, 006, 007, 010, 011, 012, 013, 014, 015, 016, 017, 020, 021, 022, 023, 024,
	025, 026, 027, 030, 031, 032, 033, 034, 035, 036, 037, 001, 057, 060, 004, 077,
	000, 006, 007, 010, 011, 012, 013, 014, 015, 016, 017, 020, 021, 022, 023, 024,
	025, 026, 027, 030, 031, 032, 033, 034, 035, 036, 037, 054, 057, 055, 004, 077,
}

func SerializeUint32IntoBuffer(value uint32, buffer []byte) {
	buffer[0] = byte(value >> 24)
	buffer[1] = byte(value >> 16)
	buffer[2] = byte(value >> 8)
	buffer[3] = byte(value)
}

func SerializeUint64IntoBuffer(value uint64, buffer []byte) {
	buffer[0] = byte(value >> 56)
	buffer[1] = byte(value >> 48)
	buffer[2] = byte(value >> 40)
	buffer[3] = byte(value >> 32)
	buffer[4] = byte(value >> 24)
	buffer[5] = byte(value >> 16)
	buffer[6] = byte(value >> 8)
	buffer[7] = byte(value)
}

func DeserializeUint32FromBuffer(buffer []byte) uint32 {
	return (uint32(buffer[0]) << 24) |
		(uint32(buffer[1]) << 16) |
		(uint32(buffer[2]) << 8) |
		uint32(buffer[3])
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

// GetTwosComplement takes a number which is a 36-bit signed value packed into a uint64,
// and converts it to twos-complement.
func GetTwosComplement(operand uint64) uint64 {
	if IsNegative(operand) {
		return -Negate(operand)
	} else {
		return operand
	}
}

func DeserializeUint64FromBuffer(buffer []byte) uint64 {
	return (uint64(buffer[0]) << 56) |
		(uint64(buffer[1]) << 48) |
		(uint64(buffer[2]) << 40) |
		(uint64(buffer[3]) << 32) |
		(uint64(buffer[4]) << 24) |
		(uint64(buffer[5]) << 16) |
		(uint64(buffer[6]) << 8) |
		uint64(buffer[7])
}

func Word36ToByteArrayPacked(
	source []uint64,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
) (nonIntegral bool, byteCount uint) {
	nonIntegral = false
	byteCount = 0

	sourceLimit := sourceOffset + sourceLength
	sx := sourceOffset
	dx := destinationOffset
	for sx < sourceLimit {
		destination[dx] = byte(source[sx] >> 28)
		destination[dx+1] = byte(source[sx] >> 20)
		destination[dx+2] = byte(source[sx] >> 12)
		destination[dx+3] = byte(source[sx] >> 4)
		if sx == sourceLimit-1 {
			destination[dx+4] = byte(source[sx] << 4)
			sx++
			nonIntegral = true
			byteCount += 5
			break
		}

		destination[dx+4] = (byte(source[sx] << 4)) | (byte(source[sx+1]>>32) & 0x0F)
		sx++

		destination[dx+5] = byte(source[sx] >> 24)
		destination[dx+6] = byte(source[sx] >> 16)
		destination[dx+7] = byte(source[sx] >> 8)
		destination[dx+8] = byte(source[sx])

		sx++
		dx += 9
		byteCount += 9
	}

	return
}

func Word36ToByteArrayPackedReversed(
	source []uint64,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
) (nonIntegral bool, byteCount uint) {
	nonIntegral = false
	byteCount = 0

	sx := sourceOffset + sourceLength
	dx := destinationOffset
	for sx > sourceOffset {
		sx--
		destination[dx] = byte(source[sx])
		destination[dx+1] = byte(source[sx] >> 8)
		destination[dx+2] = byte(source[sx] >> 16)
		destination[dx+3] = byte(source[sx] >> 24)

		if sx == sourceOffset+1 {
			destination[dx+4] = byte((source[sx] >> 28) & 0xF0)
			nonIntegral = true
			byteCount += 5
			break
		}

		destination[dx+4] = byte(source[sx]>>32) | byte(source[sx-1]&0x0F)
		sx--

		destination[dx+5] = byte(source[sx] >> 4)
		destination[dx+6] = byte(source[sx] >> 12)
		destination[dx+7] = byte(source[sx] >> 20)
		destination[dx+8] = byte(source[sx] >> 28)

		dx += 9
		byteCount += 9
	}

	return
}

func Word36ToByteArray8Bit(
	source []uint64,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
) (nonIntegral bool, byteCount uint) {
	nonIntegral = false
	byteCount = 0

	sourceLimit := sourceOffset + sourceLength
	sx := sourceOffset
	dx := destinationOffset
	for sx < sourceLimit {
		destination[dx] = byte(GetQ1(source[sx]))
		destination[dx+1] = byte(GetQ2(source[sx]))
		destination[dx+2] = byte(GetQ3(source[sx]))
		destination[dx+3] = byte(GetQ4(source[sx]))
		sx++
		dx += 4
		byteCount += 4
	}

	return
}

func Word36ToByteArray8BitReversed(
	source []uint64,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
) (nonIntegral bool, byteCount uint) {
	nonIntegral = false
	byteCount = 0

	sx := sourceOffset + sourceLength
	dx := destinationOffset
	for sx > sourceOffset {
		sx--
		destination[dx] = byte(GetQ4(source[sx]))
		destination[dx+1] = byte(GetQ3(source[sx]))
		destination[dx+2] = byte(GetQ2(source[sx]))
		destination[dx+3] = byte(GetQ1(source[sx]))
		dx += 4
		byteCount += 4
	}

	return
}

func Word36ToByteArray6Bit(
	source []uint64,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
) (nonIntegral bool, byteCount uint) {
	nonIntegral = false
	byteCount = 0

	sourceLimit := sourceOffset + sourceLength
	sx := sourceOffset
	dx := destinationOffset
	for sx < sourceLimit {
		destination[dx] = byte(GetS1(source[sx]))
		destination[dx+1] = byte(GetS2(source[sx]))
		destination[dx+2] = byte(GetS3(source[sx]))
		destination[dx+3] = byte(GetS4(source[sx]))
		destination[dx+4] = byte(GetS5(source[sx]))
		destination[dx+5] = byte(GetS6(source[sx]))
		sx++
		dx += 6
		byteCount += 6
	}

	return
}

func Word36ToByteArray6BitReversed(
	source []uint64,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
) (nonIntegral bool, byteCount uint) {
	nonIntegral = false
	byteCount = 0

	sx := sourceOffset + sourceLength
	dx := destinationOffset
	for sx > sourceOffset {
		sx--
		destination[dx] = byte(GetS6(source[sx]))
		destination[dx+1] = byte(GetS5(source[sx]))
		destination[dx+2] = byte(GetS4(source[sx]))
		destination[dx+3] = byte(GetS3(source[sx]))
		destination[dx+3] = byte(GetS2(source[sx]))
		destination[dx+3] = byte(GetS1(source[sx]))
		dx += 6
		byteCount += 6
	}

	return
}

func ByteArrayPackedToWord36(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []uint64,
	destinationOffset uint,
) (nonIntegral bool, wordCount uint) {
	nonIntegral = false
	wordCount = 0

	sx := sourceOffset
	dx := destinationOffset
	for sy := 0; sx < sourceOffset+sourceLength; sy++ {
		switch sy % 9 {
		case 0:
			destination[dx] = SetW(0, uint64(source[sx])<<28)
			sx++
			nonIntegral = true
			wordCount++
			break

		case 1:
			destination[dx] = Or(destination[dx], uint64(source[sx])<<20)
			sx++
			break

		case 2:
			destination[dx] = Or(destination[dx], uint64(source[sx])<<12)
			sx++
			break

		case 3:
			destination[dx] = Or(destination[dx], uint64(source[sx])<<4)
			sx++
			break

		case 4:
			destination[dx] = Or(destination[dx], uint64(source[sx]>>4))
			dx++
			destination[dx] = SetW(0, uint64(source[sx]&0x0F)<<32)
			sx++
			wordCount++
			break

		case 5:
			destination[dx] = Or(destination[dx], uint64(source[sx])<<24)
			sx++
			break

		case 6:
			destination[dx] = Or(destination[dx], uint64(source[sx])<<16)
			sx++
			break

		case 7:
			destination[dx] = Or(destination[dx], uint64(source[sx])<<8)
			sx++
			break

		case 8:
			destination[dx] = Or(destination[dx], uint64(source[sx]))
			dx++
			sx++
			nonIntegral = false
		}
	}

	return
}

func ByteArrayPackedToWord36Reversed(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []uint64,
	destinationOffset uint,
) (nonIntegral bool, wordCount uint) {
	nonIntegral = false
	wordCount = 0

	sx := sourceOffset + sourceLength
	dx := destinationOffset
	for sx > sourceOffset {
		switch sx % 9 {
		case 0:
			sx--
			destination[dx] = SetW(0, uint64(source[sx])<<28)
			nonIntegral = true
			wordCount++
			break

		case 1:
			sx--
			destination[dx] = Or(destination[dx], uint64(source[sx])<<20)
			break

		case 2:
			sx--
			destination[dx] = Or(destination[dx], uint64(source[sx])<<12)
			break

		case 3:
			sx--
			destination[dx] = Or(destination[dx], uint64(source[sx])<<4)
			break

		case 4:
			sx--
			destination[dx] = Or(destination[dx], uint64(source[sx]>>4))
			dx++
			destination[dx] = SetW(0, uint64(source[sx]&0x0F)<<32)
			wordCount++
			break

		case 5:
			sx--
			destination[dx] = Or(destination[dx], uint64(source[sx])<<24)
			break

		case 6:
			sx--
			destination[dx] = Or(destination[dx], uint64(source[sx])<<16)
			break

		case 7:
			sx--
			destination[dx] = Or(destination[dx], uint64(source[sx])<<8)
			break

		case 8:
			sx--
			destination[dx] = Or(destination[dx], uint64(source[sx]))
			nonIntegral = false
		}
	}

	return
}

func ByteArray8BitToWord36(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []uint64,
	destinationOffset uint,
) (nonIntegral bool, wordCount uint) {
	nonIntegral = false
	wordCount = 0

	sx := sourceOffset
	dx := destinationOffset
	for sy := 0; sx < sourceOffset+sourceLength; sy++ {
		switch sy % 4 {
		case 0:
			destination[dx] = SetQ1(destination[dx], uint64(source[sx]))
			sx++
			nonIntegral = true
			wordCount++
			break

		case 1:
			destination[dx] = SetQ2(destination[dx], uint64(source[sx]))
			sx++
			break

		case 2:
			destination[dx] = SetQ3(destination[dx], uint64(source[sx]))
			sx++
			break

		case 3:
			destination[dx] = SetQ4(destination[dx], uint64(source[sx]))
			sx++
			dx++
			nonIntegral = false
		}
	}

	return
}

func ByteArray8BitToWord36Reversed(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []uint64,
	destinationOffset uint,
) (nonIntegral bool, wordCount uint) {
	nonIntegral = false
	wordCount = sourceLength / 4
	if sourceLength%4 > 0 {
		wordCount++
		nonIntegral = true
	}

	sx := sourceOffset + sourceLength
	dx := destinationOffset
	for sx > sourceOffset {
		sx--
		switch sx % 4 {
		case 3:
			destination[dx] = SetQ1(destination[dx], uint64(source[sx]))
			break

		case 2:
			destination[dx] = SetQ2(destination[dx], uint64(source[sx]))
			break

		case 1:
			destination[dx] = SetQ3(destination[dx], uint64(source[sx]))
			break

		case 0:
			destination[dx] = SetQ4(destination[dx], uint64(source[sx]))
			break
		}
	}

	return
}

func ByteArray6BitToWord36(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []uint64,
	destinationOffset uint,
) (nonIntegral bool, wordCount uint) {
	nonIntegral = false
	wordCount = 0

	sx := sourceOffset
	dx := destinationOffset
	for sy := 0; sx < sourceOffset+sourceLength; sy++ {
		switch sy % 6 {
		case 0:
			destination[dx] = SetS1(destination[dx], uint64(source[sx]))
			sx++
			nonIntegral = true
			wordCount++
			break

		case 1:
			destination[dx] = SetS2(destination[dx], uint64(source[sx]))
			sx++
			break

		case 2:
			destination[dx] = SetS3(destination[dx], uint64(source[sx]))
			sx++
			break

		case 3:
			destination[dx] = SetS4(destination[dx], uint64(source[sx]))
			sx++
			break

		case 4:
			destination[dx] = SetS5(destination[dx], uint64(source[sx]))
			sx++
			break

		case 5:
			destination[dx] = SetS6(destination[dx], uint64(source[sx]))
			sx++
			dx++
			nonIntegral = false
			break
		}
	}

	return
}

func ByteArray6BitToWord36Reversed(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []uint64,
	destinationOffset uint,
) (nonIntegral bool, wordCount uint) {
	nonIntegral = false
	wordCount = sourceLength / 6
	if sourceLength%6 > 0 {
		wordCount++
		nonIntegral = true
	}

	sx := sourceOffset + sourceLength
	dx := destinationOffset
	for sx > sourceOffset {
		sx--
		switch sx % 6 {
		case 5:
			destination[dx] = SetS1(destination[dx], uint64(source[sx]))
			break

		case 4:
			destination[dx] = SetS2(destination[dx], uint64(source[sx]))
			break

		case 3:
			destination[dx] = SetS3(destination[dx], uint64(source[sx]))
			break

		case 2:
			destination[dx] = SetS4(destination[dx], uint64(source[sx]))
			break

		case 1:
			destination[dx] = SetS5(destination[dx], uint64(source[sx]))
			break

		case 0:
			destination[dx] = SetS6(destination[dx], uint64(source[sx]))
			break
		}
	}

	return
}

// PackWord36Strict packs pairs of word36 structs into 9-byte sequences
func PackWord36Strict(source []uint64, destination []byte) error {
	sl := len(source)
	if sl%1 != 0 {
		return fmt.Errorf("source buffer does not contain an even number of words (%v)", len(source))
	}

	if sl*9/2 > len(destination) {
		return fmt.Errorf("destination buffer insufficient size (%v)", len(destination))
	}

	_, _ = Word36ToByteArrayPacked(source, 0, uint(len(source)), destination, 0)
	return nil
}

// UnpackWord36Strict unpacks 9-byte groups into pairs of Word36 structs
func UnpackWord36Strict(source []byte, destination []uint64) error {
	srcLen := len(source)
	if srcLen%9 != 0 {
		return fmt.Errorf("source buffer length %v is not a multiple of 9 bytes", len(source))
	}

	if srcLen*2/9 > len(destination) {
		return fmt.Errorf("destination buffer insufficient size (%v)", len(destination))
	}

	_, _ = ByteArrayPackedToWord36(source, 0, uint(len(source)), destination, 0)
	return nil
}
