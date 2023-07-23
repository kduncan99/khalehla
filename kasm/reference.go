// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
	"strings"
)

type Reference struct {
	symbol     string
	arguments  []*Expression
	levelCount int
}

func (r *Reference) Evaluate(ec *ExpressionContext) error {
	selectors := make([]Value, len(r.arguments))
	for ax := 0; ax < len(r.arguments); ax++ {
		err := r.arguments[ax].Evaluate(ec)
		if err != nil {
			return err
		}

		val, err := ec.PopValue()
		if err != nil {
			return err
		}

		selectors[ax] = val
	}

	value, err := ec.context.dictionary.Lookup(r.symbol)
	if err != nil {
		if len(selectors) > 0 {
			return fmt.Errorf("undefined reference with arguments or selectors cannot be resolved")
		}

		offset := ValueOffset{
			symbol:     strings.ToUpper(r.symbol),
			isNegative: false,
		}
		component := ValueComponent{
			bitCount: 36,
			value:    0,
			offsets:  []ValueOffset{offset},
		}
		iv := &IntegerValue{
			components: []ValueComponent{component},
		}
		ec.PushValue(iv)
		return nil
	}

	if value.GetValueType() == FloatValueType {
		if len(selectors) > 0 {
			return fmt.Errorf("symbol %v does not resolve to a node", r.symbol)
		}

		ec.PushValue(value)
		return nil
	}

	if value.GetValueType() == NodeValueType {
		n := value.(*NodeValue)
		return n.Evaluate(ec)
	}

	if value.GetValueType() == FunctionValueType {
		f := value.(FunctionValue)
		return f.Evaluate(ec)
	}

	//	User referred to a value which isn't allowed in an expression
	return fmt.Errorf("symbol value cannot be used in an expression")
}

func (p *Parser) ParseReference(allowLeadingDollar bool, allowLevelers bool) (*Reference, error) {
	p.SkipWhiteSpace()
	symbol, err := p.ParseSymbol(allowLeadingDollar)
	if err != nil {
		return nil, err
	}

	//	parse optional arg list in parentheses
	args, err := p.ParseExpressionList()
	if err != nil {
		return nil, err
	}

	//	parse optional leveler asterisks
	var levelCount int
	if allowLevelers {
		for p.ParseCharacter('*') {
			levelCount++
		}
	}

	ref := &Reference{
		symbol:     *symbol,
		arguments:  args,
		levelCount: levelCount,
	}

	return ref, nil
}
