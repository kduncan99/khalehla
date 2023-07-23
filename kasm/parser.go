// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
	"strings"
)

var invalidPosition = fmt.Errorf("invalid position")
var outOfData = fmt.Errorf("out of data")

type Parser struct {
	index int
	text  string
}

func NewParser(text string) *Parser {
	return &Parser{
		text:  text,
		index: 0,
	}
}

func (p *Parser) Advance(skipCount int) error {
	if (skipCount < 0) || (p.index+skipCount > len(p.text)) {
		return invalidPosition
	} else {
		p.index += skipCount
		return nil
	}
}

func (p *Parser) AtEnd() bool {
	return p.index >= len(p.text)
}

func (p *Parser) GetPosition() int {
	return p.index
}

func IsAlphabetic(char uint8) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
}

func IsDecimalDigit(char uint8) bool {
	return char >= '0' && char <= '9'
}

func IsWhiteSpace(char uint8) bool {
	return char == ' ' || char == '\t' || char == '\r' || char == '\n'
}

func (p *Parser) NextChar() (uint8, error) {
	if p.AtEnd() {
		return 0, outOfData
	} else {
		res := p.text[p.index]
		p.index++
		return res, nil
	}
}

func (p *Parser) ParseCharacter(char uint8) bool {
	if !p.AtEnd() && p.text[p.index] == char {
		p.index++
		return true
	} else {
		return false
	}
}

func (p *Parser) ParseToken(token string) bool {
	if p.Remaining() >= len(token) {
		px := p.index
		tx := 0
		for tx < len(token) {
			if token[tx] != p.text[px] {
				return false
			} else {
				px++
				tx++
			}
		}
		return true
	} else {
		return false
	}
}

func (p *Parser) ParseSymbol(allowLeadingDollar bool) (*string, error) {
	if !p.AtEnd() {
		ch, _ := p.PeekNextChar()
		if IsAlphabetic(ch) || (allowLeadingDollar && ch == '$') {
			_ = p.Advance(1)
			symbol := string(ch)
			for !p.AtEnd() {
				ch, _ = p.PeekNextChar()
				if !IsAlphabetic(ch) && !IsDecimalDigit(ch) && ch != '_' && ch != '$' {
					break
				}

				_ = p.Advance(1)
				symbol += string(ch)
			}

			if len(symbol) > 12 {
				return nil, fmt.Errorf("symbol too long")
			}
			return &symbol, nil
		}
	}

	return nil, nil
}

func (p *Parser) ParseTokenCaseInsensitive(token string) bool {
	if p.Remaining() >= len(token) {
		px := p.index
		tx := 0
		for tx < len(token) {
			tUpper := strings.ToUpper(token[tx : tx+1])
			pUpper := strings.ToUpper(p.text[px : px+1])
			if tUpper != pUpper {
				return false
			} else {
				px++
				tx++
			}
		}
		return true
	} else {
		return false
	}
}

func (p *Parser) PeekNextChar() (uint8, error) {
	if p.AtEnd() {
		return 0, outOfData
	} else {
		return p.text[p.index], nil
	}
}

func (p *Parser) Remaining() int {
	return len(p.text) - p.index
}

func (p *Parser) SetPosition(index int) error {
	if index >= len(p.text) {
		return invalidPosition
	} else {
		p.index = index
		return nil
	}
}

func (p *Parser) SkipWhiteSpace() int {
	var count int
	for p.index < len(p.text) {
		if !IsWhiteSpace(p.text[p.index]) {
			break
		} else {
			count++
			p.index++
		}
	}
	return count
}
