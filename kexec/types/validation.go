// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

// IsValidFilename tests a given string to ensure it is a valid filename.
// The string must be 1 to 12 character in length, containing any combination of
// upper-case letters, digits, hyphens, and dollar signs.
func IsValidFilename(filename string) bool {
	if len(filename) < 1 || len(filename) > 12 {
		return false
	}

	for _, ch := range filename {
		if (ch < 'A' || ch > 'Z') && ch != '-' && ch != '$' {
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

	for _, ch := range qualifier {
		if (ch < 'A' || ch > 'Z') && ch != '-' && ch != '$' {
			return false
		}
	}

	return true
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

func IsValidPackName(name string) bool {
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

func IsValidPrepFactor(prepFactor PrepFactor) bool {
	return prepFactor == 28 || prepFactor == 56 || prepFactor == 112 || prepFactor == 224 ||
		prepFactor == 448 || prepFactor == 896 || prepFactor == 1792
}
