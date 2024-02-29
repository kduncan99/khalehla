// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
	"strings"
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

func (pi *Parser) GetRemainder() string {
	return pi.text[pi.index:]
}

func (pi *Parser) IsAtEnd() bool {
	return pi.index >= len(pi.text)
}

func (pi *Parser) MarkPosition() {
	pi.mark = pi.index
}

func (pi *Parser) ParseDecimalDigit() (result byte, found bool) {
	result = 0
	found = false

	if !pi.IsAtEnd() && (pi.text[pi.index] >= '0' && pi.text[pi.index] <= '9') {
		result = pi.text[pi.index]
		pi.index++
		found = true
	}

	return
}

func (pi *Parser) ParseIdentifier() (result string, found bool, ok bool) {
	result = ""
	found = false
	ok = true

	if pi.IsAtEnd() || !pi.isAlphabeticCharacter(pi.text[pi.index]) {
		return
	}

	ch, _ := pi.ParseNextCharacter()
	result = string(ch)
	found = true
	for !pi.IsAtEnd() && pi.isIdentifierCharacter(pi.text[pi.index]) {
		if len(result) == 6 {
			ok = false
			return
		}

		ch, _ := pi.ParseNextCharacter()
		result += string(ch)
	}

	result = strings.ToUpper(result)
	return
}

func (pi *Parser) ParseNextCharacter() (result byte, found bool) {
	result = 0
	found = false

	if !pi.IsAtEnd() {
		ch := pi.text[pi.index]
		pi.index++
		result = ch
		found = true
	}

	return
}

func (pi *Parser) ParseUntil(cutSet string) (result string, terminator uint8) {
	result = ""
	terminator = 0

	for !pi.IsAtEnd() && !strings.ContainsRune(cutSet, rune(pi.text[pi.index])) {
		result += string(pi.text[pi.index])
		pi.index++
	}
	if !pi.IsAtEnd() {
		terminator = pi.text[pi.index]
	}

	return
}

func (pi *Parser) ParseSpecificCharacter(ch byte) (found bool) {
	found = false

	if !pi.IsAtEnd() && pi.text[pi.index] == ch {
		pi.index++
		found = true
	}

	return
}

func (pi *Parser) ParseUnsignedInteger() (value uint64, found bool, err error) {
	value = 0
	err = nil

	dig, subFound := pi.ParseDecimalDigit()
	found = subFound
	for subFound {
		newValue := value*10 + uint64(dig-'0')
		if newValue < value {
			err = fmt.Errorf("overflow")
			return
		}

		value = newValue
		dig, subFound = pi.ParseDecimalDigit()
	}

	return
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

func (pi *Parser) isAlphabeticCharacter(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func (pi *Parser) isDecimalCharacter(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (pi *Parser) isIdentifierCharacter(ch byte) bool {
	return pi.isAlphabeticCharacter(ch) || pi.isDecimalCharacter(ch)
}
