// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
)

type Parser struct {
	text  string
	index int
	mark  int
}

func NewParser(text string) *Parser {
	return &Parser{
		text:  text,
		index: 0,
		mark:  0,
	}
}

func (pi *Parser) IsAtEnd() bool {
	return pi.index >= len(pi.text)
}

func (pi *Parser) MarkPosition() {
	pi.mark = pi.index
}

func (pi *Parser) ParseDecimalDigit() (result *byte) {
	result = nil
	if !pi.IsAtEnd() && (pi.text[pi.index] >= '0' && pi.text[pi.index] <= '9') {
		ch := pi.text[pi.index]
		result = &ch
		pi.index++
	}
	return
}

func (pi *Parser) ParseUnsignedInteger() (value *uint64, err error) {
	value = nil
	err = nil

	val := uint64(0)
	dig := pi.ParseDecimalDigit()
	found := false
	for dig != nil {
		found = true
		newVal := val*10 + uint64(*dig-'0')
		if newVal < val {
			err = fmt.Errorf("overflow")
			return
		}
		val = newVal
		dig = pi.ParseDecimalDigit()
	}

	if found {
		value = &val
	}
	return
}

func (pi *Parser) ParseNextCharacter() (byte, bool) {
	if !pi.IsAtEnd() {
		ch := pi.text[pi.index]
		pi.index++
		return ch, true
	} else {
		return 0, false
	}
}

func (pi *Parser) ParseSpecificCharacter(ch byte) bool {
	if !pi.IsAtEnd() && pi.text[pi.index] == ch {
		pi.index++
		return true
	} else {
		return false
	}
}

func (pi *Parser) ResetPosition() {
	pi.index = pi.mark
	pi.mark = 0
}

func (pi *Parser) SkipSpaces() int {
	count := 0
	for !pi.IsAtEnd() {
		if pi.text[pi.index] != ' ' {
			break
		}
		count++
		pi.index++
	}
	return count
}
