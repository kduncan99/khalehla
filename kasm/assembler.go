// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type assembledLine struct {
	sourceName *string
	sourceCode *string
	lineNumber int
}

type Assembler struct {
	assembledLines []assembledLine
	context        *Context
}

func (a *Assembler) Assemble(sourceName string, sourceCode []string) bool {
	a.context.diagnostics.Clear()
	a.context.currentLineIndex = 0
	a.context.currentLocationCounter = 0
	a.context.currentLiteralPool = 0
	lineNumber := 1
	for a.context.currentLineIndex < len(sourceCode) {
		fields, ok := a.parseLine(lineNumber, sourceCode[a.context.currentLineIndex])
		if !ok {
			a.context.currentLineIndex++
			continue
		}

		a.interpretFields(lineNumber, fields)
		//	TODO
	}

	errors, _, _ := a.context.diagnostics.GetDiagnosticCounters()
	return errors == 0
}

func (a *Assembler) interpretFields(lineNumber int, fields [][]string) {
	var locationCounter *string
	var labelSpecifications []string

	//	Field 0 contains 0 or more subfields.
	//  The first subfield is either a location counter or a label specification.
	//  All subsequent subfields are label specifications.
	f0 := fields[0]
	if f0 != nil && len(f0) > 0 {
		sfx := 0
		if isValidLocationCounter(f0[sfx]) {
			locationCounter = &f0[sfx]
			sfx++
		}

		for sfx < len(f0) {
			if !isValidLabelSpecification(f0[sfx]) {
				a.context.diagnostics.AppendError(lineNumber, "Invalid label specification:"+f0[sfx])
			} else {
				labelSpecifications = append(labelSpecifications, f0[sfx])
			}
		}
	}

	//	All else depends upon the content of field 1, subfield 0 - if there is such a thing
	f1 := fields[1]
	if f1 == nil || len(f1) == 0 {
		a.context.currentLineIndex++
		return
	}

	//	Is this a directive invocation?
	if a.interpretDirective(lineNumber, fields, locationCounter, labelSpecifications) {
		return
	}

	//	Is this a PROC call?
	//	TODO

	//	None of the above... treat it as an expression and try to generate code
	//	TODO
}

func (a *Assembler) interpretDirective(lineNumber int, fields [][]string, locCounter *string, labelSpecs []string) bool {
	dir := fields[1][0]
	if dir == "$EQU" {
		//	TODO
	} else if dir == "$EQUF" {
		//	TODO
	} else if dir == "$END" {
		//	TODO
	} else if dir == "$FUNC" {
		//	TODO
	} else if dir == "$LIT" {
		a.interpretLiteralDirective(lineNumber, fields, locCounter, labelSpecs)
		a.context.currentLineIndex++
		return true
	} else if dir == "$PROC" {
		//	TODO
	}

	return false
}

func (a *Assembler) interpretLiteralDirective(lineNumber int, fields [][]string, locCounter *string, labelSpecs []string) {
	if len(labelSpecs) > 0 {
		//	we have labels, we'd better have a location counter spec
		if locCounter == nil {
			a.context.diagnostics.AppendWarning(lineNumber, "ignoring label on $LIT directive")
			a.context.currentLiteralPool = a.context.currentLocationCounter
		} else {
			//	create functions out of all the label names per MASM spec
			//	TODO
		}
	} else if locCounter != nil {
		//	no label, just a location counter - set the literal pool to the lcn but do not update the current lc
		lcn, err := interpretLocationCounter(*locCounter)
		if err != nil {
			a.context.diagnostics.AppendError(lineNumber, err.Error())
		}
		a.context.currentLiteralPool = lcn
	} else {
		//	no lc, no label, just $LIT - set the literal pool to the current location counter
		a.context.currentLiteralPool = a.context.currentLocationCounter
	}
}

func (a *Assembler) parseLine(lineNumber int, source string) ([][]string, bool) {
	fields := make([][]string, 0)
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
			a.context.diagnostics.AppendError(lineNumber, "unexpected close-parenthesis")
			return fields, false
		}

		if char == "(" {
			staging += char
			parenDepth++
			postComma = false
			prevSpace = false
			continue
		}

		if char == "," {
			if fields[fieldNumber] == nil {
				fields[fieldNumber] = make([]string, 0)
			}
			fields[fieldNumber] = append(fields[fieldNumber], staging)
			staging = ""
			postComma = true
			continue
		}

		if char == " " {
			prevSpace = true
			if !postComma {
				if fields[fieldNumber] == nil {
					fields[fieldNumber] = make([]string, 0)
				}
				fields[fieldNumber] = append(fields[fieldNumber], staging)
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
		a.context.diagnostics.AppendError(lineNumber, "unterminated string literal")
		return fields, false
	}

	if parenDepth > 0 {
		a.context.diagnostics.AppendError(lineNumber, "unterminated grouping")
		return fields, false
	}

	if len(staging) > 0 {
		if fields[fieldNumber] == nil {
			fields[fieldNumber] = make([]string, 0)
		}
		fields[fieldNumber] = append(fields[fieldNumber], staging)
	}

	return fields, true
}

func (a *Assembler) display(allSource bool, undefinedReferences bool) {

}
