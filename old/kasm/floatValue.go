// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"

	"khalehla/old/parser"
)

type FloatValue struct {
	value float64
	flags ValueFlags
}

func (v *FloatValue) Evaluate(ec *ExpressionContext) error {
	ec.PushValue(v)
	return nil
}

func (v *FloatValue) ClearFlags(flags ValueFlags) {
	v.flags &= flags ^ ValueFlags(-1)
}

func (v *FloatValue) GetValueType() ValueType {
	return FloatValueType
}

func (v *FloatValue) Copy() Value {
	return NewFloatValue(v.value)
}

func (v *FloatValue) GetFlags() ValueFlags {
	return v.flags
}

func (v *FloatValue) SetFlags(flags ValueFlags) {
	v.flags |= flags
}

func NewFloatValue(value float64) *FloatValue {
	return &FloatValue{
		value: value,
	}
}

func NewFloatValueFromInteger(value *IntegerValue) (*FloatValue, error) {
	if value.form.Equals(SimpleForm) && len(value.offsets) == 0 {
		return &FloatValue{value: float64(value.componentValues[0])}, nil
	} else {
		return nil, fmt.Errorf("cannot convert composite integer or integer with offsets to float")
	}
}

func ParseFloatLiteral(p *parser.Parser) (Value, error) {
	if !p.AtEnd() {
		ch, _ := p.PeekNextChar()
		if ch >= '0' && ch <= '9' {
			position := p.GetPosition()
			_ = p.Advance(1)

			isOctal := ch == '0'
			hasDecimal := false
			var value float64
			var mult = 1.0

			for !p.AtEnd() {
				ch, _ := p.PeekNextChar()
				if ch == '.' {
					if hasDecimal {
						return nil, syntaxError
					}

					hasDecimal = true
					_ = p.Advance(1)
				} else if ch >= '0' && ch <= '9' {
					if isOctal && ch >= '8' {
						return nil, nonOctalDigit
					}

					if hasDecimal {
						if isOctal {
							mult /= 8.0
						} else {
							mult /= 10.0
						}

						value += float64(ch-'0') * mult
					} else {
						if isOctal {
							value *= 8
						} else {
							value *= 10
						}
						value += float64(ch - '0')
					}

					_ = p.Advance(1)
				} else {
					break
				}
			}

			if !hasDecimal {
				_ = p.SetPosition(position)
				return nil, nil
			}

			return NewFloatValue(value), nil
		}
	}

	return nil, nil
}
