// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
)

type LabelSpecification struct {
	symbol     string
	selectors  []*Expression
	levelCount int
}

func NewLabelSpecification(symbol string, selectors []*Expression, levelCount int) *LabelSpecification {
	return &LabelSpecification{
		symbol:     symbol,
		selectors:  selectors,
		levelCount: levelCount,
	}
}

// ParseLabelSpecification parses a label subfield consisting of
//
//	symbol [ selectors ] [ levelers ]
func (p *parser.Parser) ParseLabelSpecification(context *Context) (*LabelSpecification, error) {
	p.SkipWhiteSpace()
	symbol, err := p.ParseSymbol()
	if err != nil {
		return nil, err
	} else if symbol == nil {
		return nil, nil
	}

	if (*symbol)[0] == '$' {
		return nil, fmt.Errorf("label specification cannot begin with a $")
	}

	//	parse optional arg list in parentheses
	expList, err := p.ParseExpressionList(context)
	if err != nil {
		return nil, err
	}

	//	parse optional leveler asterisks
	var levelCount int
	for p.ParseCharacter('*') {
		levelCount++
	}

	return NewLabelSpecification(*symbol, expList, levelCount), nil
}
