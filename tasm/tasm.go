// Khalehla Project
// tiny assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

import (
	"fmt"
	"khalehla/parser"
	"strconv"
)

// TinyAssembler is a very tiny assembler which assists in unit tests
type TinyAssembler struct {
	currentSegmentNumber uint64
	forms                map[string][]uint64
	segments             map[uint64]*Segment
}

func NewTinyAssembler() *TinyAssembler {
	ta := &TinyAssembler{}
	ta.currentSegmentNumber = 0
	ta.forms = map[string][]uint64{
		"W":        {36},
		"HW":       {18, 18},
		"TW":       {12, 12, 12},
		"QW":       {9, 9, 9, 9},
		"SW":       {6, 6, 6, 6, 6, 6},
		"FJAXHIU":  {6, 4, 4, 4, 1, 1, 16},
		"FJAXHIBD": {6, 4, 4, 4, 1, 1, 4, 12},
		"FJAXU":    {6, 4, 4, 4, 18},
	}
	ta.segments = make(map[uint64]*Segment)
	ta.establishSegment(ta.currentSegmentNumber)
	return ta
}

func (a *TinyAssembler) establishSegment(segmentNumber uint64) {
	_, ok := a.segments[segmentNumber]
	if !ok {
		a.segments[segmentNumber] = NewSegment()
	}
}

func evaluate(expression string) (uint64, []string, error) {
	p := parser.NewParser(expression)
	var value uint64
	var references []string
	wantOperator := false
	wantTerm := true
	for !p.AtEnd() {
		if wantOperator {
			p.SkipWhiteSpace()
			if p.ParseCharacter('+') {
				wantTerm = true
				wantOperator = false
				continue
			}

			return 0, nil, fmt.Errorf("expected an operator")
		}

		if wantTerm {
			p.SkipWhiteSpace()
			iVal, ok := p.ParseInteger(true, true)
			if ok {
				value += iVal
				wantTerm = false
				wantOperator = true
				continue
			}

			sVal, err := p.ParseSymbol() // TODO we need to parse our own symbol
			if err != nil {
				return 0, nil, err
			} else if sVal != nil {
				references = append(references, *sVal)
				wantTerm = false
				wantOperator = true
				continue
			}
		}

		return 0, nil, fmt.Errorf("syntax error in expression")
	}

	if wantTerm {
		return 0, nil, fmt.Errorf("incomplete expression")
	}

	return value, references, nil
}

func (a *TinyAssembler) processCommand(cb *CodeBlock) {
	command := cb.sourceItem.command
	operands := cb.sourceItem.operands

	if command == nil || len(*command) == 0 {
		if operands != nil && len(operands) > 0 {
			cb.diagnostics.NewWarning(cb.sourceSet, cb.lineNumber, "operands ignored - no operator specified")
		}
		return
	}

	form, ok := a.forms[*command]
	if ok {
		a.processDataGeneration(cb, form)
		return
	}

	switch *command {
	case ".ASC":
		a.processDataGenerationAscii(cb)
		break

	case ".FD":
		a.processDataGenerationFieldata(cb)
		break

	case ".FORM":
		a.processForm(cb)
		break

	case ".RES":
		//	TODO reserve space
		break

	case ".SEG":
		a.processSegment(cb)
		return
	}

	cb.diagnostics.NewError(cb.sourceSet, cb.lineNumber, "operator not recognized")
}

func (a *TinyAssembler) processDataGeneration(cb *CodeBlock, form []uint64) {
	if len(form) != len(cb.sourceItem.operands) {
		cb.diagnostics.NewError(cb.sourceSet, cb.lineNumber, "Wrong number of operands for form")
		return
	}

	var bit uint64
	var compositeValue uint64

	for fx := 0; fx < len(form); fx++ {
		bitCount := form[fx]
		iVal, symbols, err := evaluate(cb.sourceItem.operands[fx])
		if err != nil {
			cb.diagnostics.NewError(cb.sourceSet, cb.lineNumber, err.Error())
			continue
		}

		mask := uint64((1 << bitCount) - 1)
		if iVal&mask != iVal {
			cb.diagnostics.NewWarning(cb.sourceSet, cb.lineNumber, fmt.Sprintf("truncated value at bit %v", bit))
			iVal &= mask
		}

		compositeValue <<= bitCount
		compositeValue |= iVal

		offset := a.segments[a.currentSegmentNumber].currentLength
		for _, sym := range symbols {
			cb.references = append(cb.references, NewReference(sym, bit, bitCount, offset))
		}

		bit += bitCount
	}

	cb.code = append(cb.code, compositeValue)
}

func (a *TinyAssembler) processDataGenerationAscii(cb *CodeBlock) {
	//	TODO
}

func (a *TinyAssembler) processDataGenerationFieldata(cb *CodeBlock) {
	//	TODO
}

func (a *TinyAssembler) processForm(cb *CodeBlock) {
	//	TODO
}

func (a *TinyAssembler) processLabel(cb *CodeBlock) {
	if cb.sourceItem.label != nil && len(*cb.sourceItem.label) > 0 {
		_, ok := a.segments[a.currentSegmentNumber].labels[*cb.sourceItem.label]
		if ok {
			cb.diagnostics.NewWarning(cb.sourceSet, cb.lineNumber, "label overridden")
		}
		a.segments[a.currentSegmentNumber].labels[*cb.sourceItem.label] = a.segments[a.currentSegmentNumber].currentLength
	}
}

func (a *TinyAssembler) processSegment(cb *CodeBlock) {
	if len(cb.sourceItem.operands) != 1 {
		cb.diagnostics.NewError(cb.sourceSet, cb.lineNumber, "Too many operands")
		return
	}

	var err error
	oper := cb.sourceItem.operands[0]
	radix := 10
	if oper[0:1] == "0" {
		radix = 8
	}

	segNum, err := strconv.ParseInt(oper, radix, 64)
	if err != nil {
		cb.diagnostics.NewError(cb.sourceSet, cb.lineNumber, "Bad operand")
	}

	if segNum < 0 || segNum > 077 {
		cb.diagnostics.NewError(cb.sourceSet, cb.lineNumber, "Invalid segment number")
	} else {
		a.establishSegment(uint64(segNum))
		a.currentSegmentNumber = uint64(segNum)
	}
}

func (a *TinyAssembler) Assemble(source *SourceSet) {
	fmt.Printf("\nAssembling %s...\n", source.name)
	codeBlocks := make([]*CodeBlock, len(source.sourceItems))
	for sx, item := range source.sourceItems {
		lineNumber := uint64(sx + 1)

		seg := a.segments[a.currentSegmentNumber]
		offset := seg.currentLength
		codeBlocks[sx] = NewCodeBlock(source, lineNumber, a.currentSegmentNumber, offset)
		a.processLabel(codeBlocks[sx])

		if item.command != nil {
			a.processCommand(codeBlocks[sx])
		}

		seg.AppendCodeBlock(codeBlocks[sx])
		seg.references = append(seg.references, codeBlocks[sx].references...)
	}

	for _, cb := range codeBlocks {
		cb.Emit()
	}

	fmt.Printf("  Labels:\n")
	for segNumber, segment := range a.segments {
		for symbol, value := range segment.labels {
			fmt.Printf("    %-12s = %03o:%08o\n", symbol, segNumber, value)
		}
	}
}

func (a *TinyAssembler) GetSegments() map[uint64]*Segment {
	return a.segments
}
