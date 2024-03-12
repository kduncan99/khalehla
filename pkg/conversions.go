// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

import (
	"log"
)

var PackedBytesPerBlockFromWords = map[BlockSize]BlockSize{
	28:   128,  // slop 2 bytes
	56:   256,  // slop 4 bytes
	112:  512,  // slop 8 bytes
	224:  1024, // slop 16 bytes
	448:  2048, // slop 32 bytes
	896:  4096, // slop 64 bytes
	1792: 8192, // slop 128 bytes
}

var RawBytesPerBlockFromWords = map[BlockSize]BlockSize{
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

func Word36ToByteArrayPacked(
	source []Word36,
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
	source []Word36,
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
	source []Word36,
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
		destination[dx] = byte(source[sx].GetQ1())
		destination[dx+1] = byte(source[sx].GetQ2())
		destination[dx+2] = byte(source[sx].GetQ3())
		destination[dx+3] = byte(source[sx].GetQ4())
		sx++
		dx += 4
		byteCount += 4
	}

	return
}

func Word36ToByteArray8BitReversed(
	source []Word36,
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
		destination[dx] = byte(source[sx].GetQ4())
		destination[dx+1] = byte(source[sx].GetQ3())
		destination[dx+2] = byte(source[sx].GetQ2())
		destination[dx+3] = byte(source[sx].GetQ1())
		dx += 4
		byteCount += 4
	}

	return
}

func Word36ToByteArray6Bit(
	source []Word36,
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
		destination[dx] = byte(source[sx].GetS1())
		destination[dx+1] = byte(source[sx].GetS2())
		destination[dx+2] = byte(source[sx].GetS3())
		destination[dx+3] = byte(source[sx].GetS4())
		destination[dx+4] = byte(source[sx].GetS5())
		destination[dx+5] = byte(source[sx].GetS6())
		sx++
		dx += 6
		byteCount += 6
	}

	return
}

func Word36ToByteArray6BitReversed(
	source []Word36,
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
		destination[dx] = byte(source[sx].GetS6())
		destination[dx+1] = byte(source[sx].GetS5())
		destination[dx+2] = byte(source[sx].GetS4())
		destination[dx+3] = byte(source[sx].GetS3())
		destination[dx+3] = byte(source[sx].GetS2())
		destination[dx+3] = byte(source[sx].GetS1())
		dx += 6
		byteCount += 6
	}

	return
}

func ByteArrayPackedToWord36(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []Word36,
	destinationOffset uint,
) (nonIntegral bool, wordCount uint) {
	nonIntegral = false
	wordCount = 0

	sx := sourceOffset
	dx := destinationOffset
	for sy := 0; sx < sourceOffset+sourceLength; sy++ {
		switch sy % 9 {
		case 0:
			destination[dx].SetW(uint64(source[sx]) << 28)
			sx++
			nonIntegral = true
			wordCount++
			break

		case 1:
			destination[dx].Or(uint64(source[sx]) << 20)
			sx++
			break

		case 2:
			destination[dx].Or(uint64(source[sx]) << 12)
			sx++
			break

		case 3:
			destination[dx].Or(uint64(source[sx]) << 4)
			sx++
			break

		case 4:
			destination[dx].Or(uint64(source[sx] >> 4))
			dx++
			destination[dx].SetW(uint64(source[sx]&0x0F) << 32)
			sx++
			wordCount++
			break

		case 5:
			destination[dx].Or(uint64(source[sx]) << 24)
			sx++
			break

		case 6:
			destination[dx].Or(uint64(source[sx]) << 16)
			sx++
			break

		case 7:
			destination[dx].Or(uint64(source[sx]) << 8)
			sx++
			break

		case 8:
			destination[dx].Or(uint64(source[sx]))
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
	destination []Word36,
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
			destination[dx].SetW(uint64(source[sx]) << 28)
			nonIntegral = true
			wordCount++
			break

		case 1:
			sx--
			destination[dx].Or(uint64(source[sx]) << 20)
			break

		case 2:
			sx--
			destination[dx].Or(uint64(source[sx]) << 12)
			break

		case 3:
			sx--
			destination[dx].Or(uint64(source[sx]) << 4)
			break

		case 4:
			sx--
			destination[dx].Or(uint64(source[sx] >> 4))
			dx++
			destination[dx].SetW(uint64(source[sx]&0x0F) << 32)
			wordCount++
			break

		case 5:
			sx--
			destination[dx].Or(uint64(source[sx]) << 24)
			break

		case 6:
			sx--
			destination[dx].Or(uint64(source[sx]) << 16)
			break

		case 7:
			sx--
			destination[dx].Or(uint64(source[sx]) << 8)
			break

		case 8:
			sx--
			destination[dx].Or(uint64(source[sx]))
			nonIntegral = false
		}
	}

	return
}

func ByteArray8BitToWord36(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []Word36,
	destinationOffset uint,
) (nonIntegral bool, wordCount uint) {
	nonIntegral = false
	wordCount = 0

	sx := sourceOffset
	dx := destinationOffset
	for sy := 0; sx < sourceOffset+sourceLength; sy++ {
		switch sy % 4 {
		case 0:
			destination[dx].SetQ1(uint64(source[sx]))
			sx++
			nonIntegral = true
			wordCount++
			break

		case 1:
			destination[dx].SetQ2(uint64(source[sx]))
			sx++
			break

		case 2:
			destination[dx].SetQ3(uint64(source[sx]))
			sx++
			break

		case 3:
			destination[dx].SetQ4(uint64(source[sx]))
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
	destination []Word36,
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
			destination[dx].SetQ1(uint64(source[sx]))
			break

		case 2:
			destination[dx].SetQ2(uint64(source[sx]))
			break

		case 1:
			destination[dx].SetQ3(uint64(source[sx]))
			break

		case 0:
			destination[dx].SetQ4(uint64(source[sx]))
			break
		}
	}

	return
}

func ByteArray6BitToWord36(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []Word36,
	destinationOffset uint,
) (nonIntegral bool, wordCount uint) {
	nonIntegral = false
	wordCount = 0

	sx := sourceOffset
	dx := destinationOffset
	for sy := 0; sx < sourceOffset+sourceLength; sy++ {
		switch sy % 6 {
		case 0:
			destination[dx].SetS1(uint64(source[sx]))
			sx++
			nonIntegral = true
			wordCount++
			break

		case 1:
			destination[dx].SetS2(uint64(source[sx]))
			sx++
			break

		case 2:
			destination[dx].SetS3(uint64(source[sx]))
			sx++
			break

		case 3:
			destination[dx].SetS4(uint64(source[sx]))
			sx++
			break

		case 4:
			destination[dx].SetS5(uint64(source[sx]))
			sx++
			break

		case 5:
			destination[dx].SetS6(uint64(source[sx]))
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
	destination []Word36,
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
			destination[dx].SetS1(uint64(source[sx]))
			break

		case 4:
			destination[dx].SetS2(uint64(source[sx]))
			break

		case 3:
			destination[dx].SetS3(uint64(source[sx]))
			break

		case 2:
			destination[dx].SetS4(uint64(source[sx]))
			break

		case 1:
			destination[dx].SetS5(uint64(source[sx]))
			break

		case 0:
			destination[dx].SetS6(uint64(source[sx]))
			break
		}
	}

	return
}

// PackWord36Strict packs pairs of word36 structs into 9-byte sequences
func PackWord36Strict(source []Word36, destination []byte) {
	sl := len(source)
	if sl%1 != 0 {
		log.Panic("source buffer does not contain an even number of words")
	}

	if sl*9/2 > len(destination) {
		log.Panic("destination buffer insufficient size")
	}

	_, _ = Word36ToByteArrayPacked(source, 0, uint(len(source)), destination, 0)
}

// UnpackWord36Strict unpacks 9-byte groups into pairs of Word36 structs
func UnpackWord36Strict(source []byte, destination []Word36) {
	srcLen := len(source)
	if srcLen%9 != 0 {
		log.Panic("source buffer length is not a multiple of 9 bytes")
	}

	if srcLen*2/9 > len(destination) {
		log.Panic("destination buffer insufficient size")
	}

	_, _ = ByteArrayPackedToWord36(source, 0, uint(len(source)), destination, 0)
}
