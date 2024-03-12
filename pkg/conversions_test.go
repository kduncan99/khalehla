// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

import (
	"bytes"
	"testing"
)

func wordsEqual(a1 []Word36, a2 []Word36) bool {
	if len(a1) != len(a2) {
		return false
	}

	for ax := 0; ax < len(a1); ax++ {
		if a1[ax].GetW() != a2[ax].GetW() {
			return false
		}
	}

	return true
}

func Test_Word36ToByteArrayPacked_1(t *testing.T) {
	source := make([]Word36, 28)
	expected := make([]byte, 28*9/2)
	destination := make([]byte, 1024)

	nonInt, byteCount := Word36ToByteArrayPacked(source, 0, 28, destination, 0)
	chk := bytes.Equal(expected, destination[0:28*9/2])
	if !chk {
		t.Errorf("Error:Unequal arrays")
	}

	expBC := 28 * 9 / 2
	if byteCount != 28*9/2 {
		t.Errorf("Error:incorrect byte count:%v expected:%v", byteCount, expBC)
	}

	if nonInt {
		t.Errorf("Error:non-integral was true")
	}
}

func Test_Word36ToByteArrayPacked_2(t *testing.T) {
	source := []Word36{0, 0_111222_333444, 0_555666777000, 0}
	expected := []byte{0x24, 0xA4, 0x9B, 0x72, 0x4B, 0x6E, 0xDB, 0xFE, 0x00}
	destination := make([]byte, 16)

	nonInt, byteCount := Word36ToByteArrayPacked(source, 1, 2, destination, 2)
	chk := bytes.Equal(expected, destination[2:11])
	if !chk {
		t.Errorf("Error:Unequal arrays")
	}

	if byteCount != 9 {
		t.Errorf("Error:incorrect byte count")
	}

	if nonInt {
		t.Errorf("Error:non-integral was true")
	}
}

func Test_Word36ToByteArrayPacked_3(t *testing.T) {
	source := []Word36{0, 0_444444_444444, 0}
	expected := []byte{0x92, 0x49, 0x24, 0x92, 0x40}
	destination := make([]byte, 16)

	nonInt, byteCount := Word36ToByteArrayPacked(source, 1, 1, destination, 0)
	chk := bytes.Equal(expected, destination[:5])
	if !chk {
		t.Errorf("Error:Unequal arrays")
	}

	if byteCount != 5 {
		t.Errorf("Error:incorrect byte count")
	}

	if !nonInt {
		t.Errorf("Error:non-integral was false")
	}
}

func Test_ByteArrayPackedToWord36(t *testing.T) {
	source := []byte{0x00, 0x00, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0x00, 0x00}
	expected := []Word36{0_104315_042526, 0_316742_114652}
	destination := make([]Word36, 28)

	nonInt, wordCount := ByteArrayPackedToWord36(source, 2, 9, destination, 2)
	chk := wordsEqual(expected, destination[2:4])
	if !chk {
		t.Errorf("Error:Unequal arrays")
	}

	if wordCount != 2 {
		t.Errorf("Error:incorrect byte count")
	}

	if nonInt {
		t.Errorf("Error:non-integral was true")
	}
}

func Test_Word36ToByteArray8Bit_1(t *testing.T) {
	source := []Word36{0_777_777_777_777, 0_012_345_234_567}
	expected := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x0A, 0xE5, 0x9C, 0x77}
	destination := make([]byte, 16)

	nonInt, byteCount := Word36ToByteArray8Bit(source, 0, 2, destination, 3)
	chk := bytes.Equal(expected, destination[3:11])
	if !chk {
		t.Errorf("Error:Unequal arrays")
	}

	if byteCount != 8 {
		t.Errorf("Error:incorrect byte count:%v expected:%v", byteCount, 8)
	}

	if nonInt {
		t.Errorf("Error:non-integral was true")
	}
}

func Test_ByteArray8BitToWord36(t *testing.T) {
	source := []byte{0x00, 0x00, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC}
	expected := []Word36{0_000000_042063, 0_104125_146167, 0_210231_252273, 0_314000_000000}
	destination := make([]Word36, 28)

	nonInt, wordCount := ByteArray8BitToWord36(source, 0, 13, destination, 0)
	chk := wordsEqual(expected, destination[0:4])
	if !chk {
		t.Errorf("Error:Unequal arrays")
	}

	if wordCount != 4 {
		t.Errorf("Error:incorrect byte count:%v expected:%v", wordCount, 4)
	}

	if !nonInt {
		t.Errorf("Error:non-integral was false")
	}
}

func Test_Word36ToByteArray6Bit_1(t *testing.T) {
	source := []Word36{0_77_77_77_33_11_00, 0_01_23_45_23_45_67}
	expected := []byte{0x3F, 0x3F, 0x3F, 0x1B, 0x09, 0x00, 0x01, 0x13, 0x25, 0x13, 0x25, 0x37}
	destination := make([]byte, 16)

	nonInt, byteCount := Word36ToByteArray6Bit(source, 0, 2, destination, 3)
	chk := bytes.Equal(expected, destination[3:15])
	if !chk {
		t.Errorf("Error:Unequal arrays")
	}

	if byteCount != 12 {
		t.Errorf("Error:incorrect byte count:%v expected:%v", byteCount, 12)
	}

	if nonInt {
		t.Errorf("Error:non-integral was true")
	}
}

func Test_ByteArray6BitToWord36(t *testing.T) {
	source := []byte{0x00, 0x00, 0x22, 0x33, 0x04, 0x15, 0x26, 0x37, 0x0F, 0x38}
	expected := []Word36{0_000042_630425, 0_466717_700000}
	destination := make([]Word36, 28)

	nonInt, wordCount := ByteArray6BitToWord36(source, 0, 10, destination, 0)
	chk := wordsEqual(expected, destination[0:2])
	if !chk {
		t.Errorf("Error:Unequal arrays")
	}

	if wordCount != 2 {
		t.Errorf("Error:incorrect byte count:%v expected:%v", wordCount, 2)
	}

	if !nonInt {
		t.Errorf("Error:non-integral was false")
	}
}
