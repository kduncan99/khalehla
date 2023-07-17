// Khalehla Project
// testing assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

type parsedLine struct {
	sourceName  *string
	lineNumber  int
	fields      [][]string
	result      bool
	diagnostics []string
}

func parse(sourceName *string, lineNumber int, source string) *parsedLine {
	pl := parsedLine{
		sourceName:  sourceName,
		lineNumber:  lineNumber,
		fields:      make([][]string, 0),
		result:      true,
		diagnostics: make([]string, 0),
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
