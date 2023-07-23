// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type StringValue struct {
	value string
}

func (v *StringValue) Evaluate(ec *ExpressionContext) error {
	ec.PushValue(v)
	return nil
}

func (v *StringValue) GetValueType() ValueType {
	return StringValueType
}

func NewStringValue(value string) *StringValue {
	return &StringValue{
		value: value,
	}
}

func (p *Parser) ParseStringLiteral() (Value, error) {
	if p.ParseCharacter('\'') {
		var str string
		for !p.AtEnd() {
			if p.ParseCharacter('\'') {
				if p.ParseCharacter('\'') {
					str += "'"
				} else {
					break
				}
			} else {
				ch, _ := p.NextChar()
				str += string(ch)
			}
		}

		return NewStringValue(str), nil
	}

	return nil, nil
}
