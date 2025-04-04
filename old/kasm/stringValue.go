// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"khalehla/old/parser"
)

type StringValue struct {
	value    string
	flags    ValueFlags
	codeType StringCodeType
}

func (v *StringValue) Evaluate(ec *ExpressionContext) error {
	ec.PushValue(v)
	return nil
}

func (v *StringValue) ClearFlags(flags ValueFlags) {
	v.flags &= flags ^ ValueFlags(-1)
}

func (v *StringValue) GetValueType() ValueType {
	return StringValueType
}

func (v *StringValue) Copy() Value {
	return NewStringValue(v.value, v.codeType, v.flags)
}

func (v *StringValue) GetFlags() ValueFlags {
	return v.flags
}

func (v *StringValue) SetFlags(flags ValueFlags) {
	v.flags |= flags
}

func NewStringValue(value string, codeType StringCodeType, flags ValueFlags) *StringValue {
	return &StringValue{
		value:    value,
		codeType: codeType,
		flags:    flags,
	}
}

func ParseStringLiteral(p *parser.Parser, context *Context) (Value, error) {
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

		return NewStringValue(str, context.currentStringCodeType, 0), nil
	}

	return nil, nil
}
