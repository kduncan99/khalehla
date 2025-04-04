// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package hardware

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
