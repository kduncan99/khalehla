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
	a.context.currentLineNumber = 1
	for a.context.currentLineIndex < len(sourceCode) {
		fields, ok := a.parseLine(a.context.currentLineIndex+1, sourceCode[a.context.currentLineIndex])
		a.context.currentLineIndex++

		if !ok {
			continue
		}

		a.interpretLine(fields)
	}

	errors, _, _ := a.context.diagnostics.GetDiagnosticCounters()
	return errors == 0
}

func (a *Assembler) interpretLine(fields [][]string) {
	var labelField []string
	var operationField []string
	var operandField []string

	if len(fields) > 0 {
		labelField = fields[0]
		if len(fields) > 1 {
			operationField = fields[1]
			if len(fields) > 2 {
				operandField = fields[2]
			}
		}
	}

	dir, err := InterpretDirective(a.context, labelField, operationField, operandField)
	if err != nil {
		a.context.diagnostics.AppendError(a.context.currentLineNumber, err.Error())
	} else if dir != nil {
		//	TODO do something with the directive we got back
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
		a.context.diagnostics.AppendError(lineNumber, "unterminated string literalExpressionItem")
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
