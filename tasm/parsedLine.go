// Khalehla Project
// testing assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

type parsedLine struct {
	sourceName *string
	lineNumber int

	//	parsed fields
	fields [][]string

	//	interpreted things
	locationCounter     *string
	labelSpecifications []string

	//	diagnostics
	result      bool
	diagnostics []string
}

func parse(sourceName *string, lineNumber int, source string) *parsedLine {
	pl := parsedLine{
		sourceName:          sourceName,
		lineNumber:          lineNumber,
		fields:              make([][]string, 0),
		locationCounter:     nil,
		labelSpecifications: make([]string, 0),
		result:              true,
		diagnostics:         make([]string, 0),
	}

	fieldNumber := 0
	sx := 0
	postComma := false // cannot be true if inQuote is true
	inQuote := false
	parenDepth := 0
	prevSpace := false
	staging := ""

	for sx < len(source) {
		char := source[sx : sx+1]
		sx++
		nextChar := " "
		isLastChar := sx == len(source)
		if !isLastChar {
			nextChar = source[sx : sx+1]
		}

		if inQuote {
			staging += char
			if char == "'" {
				if !isLastChar && nextChar == "'" {
					sx++
				} else {
					inQuote = false
				}
			}
			continue
		}

		if char == "'" {
			staging += char
			inQuote = true
			prevSpace = false
			postComma = false
			continue
		}

		if parenDepth > 0 {
			staging += char
			if char == ")" {
				parenDepth--
			}
			continue
		}

		if char == ")" {
			pl.diagnostics = append(pl.diagnostics, "Unexpected close-parenthesis ignored")
			pl.result = false
			continue
		}

		if char == "(" {
			staging += char
			parenDepth++
			postComma = false
			prevSpace = false
			continue
		}

		if char == "," {
			if pl.fields[fieldNumber] == nil {
				pl.fields[fieldNumber] = make([]string, 0)
			}
			pl.fields[fieldNumber] = append(pl.fields[fieldNumber], staging)
			staging = ""
			postComma = true
			continue
		}

		if char == " " {
			prevSpace = true
			if !postComma {
				if pl.fields[fieldNumber] == nil {
					pl.fields[fieldNumber] = make([]string, 0)
				}
				pl.fields[fieldNumber] = append(pl.fields[fieldNumber], staging)
				fieldNumber++
				staging = ""
			}
			continue
		}

		if char == "." {
			if sx == 0 || prevSpace {
				prevSpace = false
				if isLastChar || nextChar == " " {
					break
				}
			}
		}
	}

	if inQuote {
		pl.diagnostics = append(pl.diagnostics, "Unterminated string literal")
		pl.result = false
	}

	if parenDepth > 0 {
		pl.diagnostics = append(pl.diagnostics, "Unterminated group")
		pl.result = false
	}

	if len(staging) > 0 {
		if pl.fields[fieldNumber] == nil {
			pl.fields[fieldNumber] = make([]string, 0)
		}
		pl.fields[fieldNumber] = append(pl.fields[fieldNumber], staging)
	}

	return &pl
}

// interpret interprets the parsed line
func (pl *parsedLine) interpret() {
	//	Field 0 contains 0 or more subfields.
	//  The first subfield is either a location counter or a label specification.
	//  All subsequent subfields are label specifications.
	f0 := pl.fields[0]
	if f0 != nil && len(f0) > 0 {
		sfx := 0
		if isValidLocationCounter(f0[sfx]) {
			pl.locationCounter = &f0[sfx]
			sfx++
		}

		for sfx < len(f0) {
			if !isValidLabelSpecification(f0[sfx]) {
				pl.diagnostics = append(pl.diagnostics, "Invalid label specification:"+f0[sfx])
				pl.result = false
			} else {
				pl.labelSpecifications = append(pl.labelSpecifications, f0[sfx])
			}
		}
	}

	//	All else depends upon the content of field 1, subfield 0 - if there is such a thing
	f1 := pl.fields[1]
	if f1 == nil || len(f1) == 0 {
		return
	}

	//	Is this a PROC call?
	//	TODO

	//	Is this a directive invocation?
	//	TODO

	//	None of the above... treat it as an expression and try to generate code
	//	TODO
}
