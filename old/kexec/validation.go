// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/common"
)

func IsValidFilenameChar(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || ch == '-' || ch == '$'
}

func IsValidNodeName(name string) bool {
	if len(name) < 1 || len(name) > 6 {
		return false
	}

	if name[0] < 'A' || name[0] > 'Z' {
		return false
	}

	for nx := 1; nx < len(name); nx++ {
		if (name[nx] < 'A' || name[nx] > 'Z') && (name[nx] < '0' || name[nx] > '9') {
			return false
		}
	}

	return true
}

func IsValidQualifierChar(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || ch == '-' || ch == '$'
}

// IsValidFilename tests a given string to ensure it is a valid filename.
// The string must be 1 to 12 character in length, containing any combination of
// upper-case letters, digits, hyphens, and dollar signs.
func IsValidFilename(filename string) bool {
	if len(filename) < 1 || len(filename) > 12 {
		return false
	}

	for chx := range filename {
		if !IsValidFilenameChar(filename[chx]) {
			return false
		}
	}

	return true
}

// IsValidQualifier test a given string to ensure it is a valid qualifier.
// The string must be 1 to 12 character in length, containing any combination of
// upper-case letters, digits, hyphens, and dollar signs.
func IsValidQualifier(qualifier string) bool {
	if len(qualifier) < 1 || len(qualifier) > 12 {
		return false
	}

	for chx := range qualifier {
		if !IsValidQualifierChar(qualifier[chx]) {
			return false
		}
	}

	return true
}

// IsValidReadWriteKey examines a string to see if it is a valid read or write key
// empty strings are vacuously valid.
func IsValidReadWriteKey(str string) bool {
	for sx := 0; sx < len(str); sx++ {
		if !IsValidReadWriteKeyChar(str[sx]) {
			return false
		}
	}
	return true
}

// IsValidReadWriteKeyChar examines a character to see whether it is a valid character
// for a read or write key.
// Any fieldata character is allowed exception period, comma, semicolon, slash, and blank.
func IsValidReadWriteKeyChar(ch uint8) bool {
	if ch > 127 || common.FieldataFromAscii[ch] == 005 {
		return false
	}
	return ch != '.' && ch != ',' && ch != ';' && ch != '/'
}

func IsValidReelNumber(name string) bool {
	if len(name) < 1 || len(name) > 6 {
		return false
	}

	for nx := 0; nx < len(name); nx++ {
		if (name[nx] < 'A' || name[nx] > 'Z') && (name[nx] < '0' || name[nx] > '9') {
			return false
		}
	}

	return true
}
