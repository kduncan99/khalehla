// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
	"khalehla/parser"
	"khalehla/pkg"
)

type IntegerValue struct {
	componentValues []int64 //	each value is 2-s complement native integer
	form            *Form
	offsets         []Offset
	flags           ValueFlags
}

func (v *IntegerValue) Evaluate(ec *ExpressionContext) error {
	ec.PushValue(v)
	return nil
}

func (v *IntegerValue) GetValueType() ValueType {
	return IntegerValueType
}

func (v *IntegerValue) ClearFlags(flags ValueFlags) {
	v.flags &= flags ^ ValueFlags(-1)
}

func (v *IntegerValue) Copy() Value {
	result, _ := NewIntegerValue(v.componentValues, v.form, v.offsets, v.flags)
	return result
}

func (v *IntegerValue) GetFlags() ValueFlags {
	return v.flags
}

func (v *IntegerValue) SetFlags(flags ValueFlags) {
	v.flags |= flags
}

func NewSimpleIntegerValue(value int64) *IntegerValue {
	return &IntegerValue{
		componentValues: []int64{value},
		form:            SimpleForm,
		offsets:         make([]Offset, 0),
		flags:           0,
	}
}

func NewIntegerValue(values []int64, form *Form, offsets []Offset, flags ValueFlags) (*IntegerValue, error) {
	if len(values) != len(form.bitSizes) {
		return nil, fmt.Errorf("number of values does not correspond to number of form fields")
	}

	return &IntegerValue{
		componentValues: values,
		form:            form,
		offsets:         offsets,
		flags:           flags,
	}, nil
}

func ParseIntegerLiteral(p *parser.Parser) (Value, error) {
	if !p.AtEnd() {
		ch, _ := p.PeekNextChar()
		if ch >= '0' && ch <= '9' {
			_ = p.Advance(1)
			isOctal := ch == '0'
			var value int64

			for !p.AtEnd() {
				ch, _ := p.PeekNextChar()
				if ch < '0' || ch > '9' {
					break
				}
				_ = p.Advance(1)

				if isOctal && ch >= '8' {
					return nil, nonOctalDigit
				}

				if isOctal {
					value *= 8
				} else {
					value *= 10
				}
				value += int64(ch - '0')
			}

			if value&pkg.NegativeZero != value {
				return nil, truncationError
			}

			return NewSimpleIntegerValue(value), nil
		}
	}

	return nil, nil
}
