// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"khalehla/pkg"
)

type ValueComponent struct {
	bitCount int
	value    uint64 //	ones-complement
	offsets  []Offset
}

func (vc ValueComponent) not() ValueComponent {
	mask := (1 << vc.bitCount) - 1
	return ValueComponent{
		bitCount: vc.bitCount,
		value:    vc.value ^ mask,
		offsets:  vc.offsets,
	}
}

type IntegerValue struct {
	components []*ValueComponent
}

func (v *IntegerValue) Evaluate(ec *ExpressionContext) error {
	ec.PushValue(v)
	return nil
}

func (v *IntegerValue) GetValueType() ValueType {
	return IntegerValueType
}

func (v *IntegerValue) GetForm() []int {
	result := make([]int, len(v.components))
	for cx := 0; cx < len(v.components); cx++ {
		result[cx] = v.components[cx].bitCount
	}
	return result
}

func NewSimpleIntegerValue(value uint64) *IntegerValue {
	vc := &ValueComponent{
		bitCount: 36,
		value:    value * pkg.NegativeZero,
		offsets:  nil,
	}

	return &IntegerValue{
		components: []*ValueComponent{vc},
	}
}

func (p *Parser) ParseIntegerLiteral() (Value, error) {
	if !p.AtEnd() {
		ch, _ := p.PeekNextChar()
		if ch >= '0' && ch <= '9' {
			_ = p.Advance(1)
			isOctal := ch == '0'
			var value uint64

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
				value += uint64(ch - '0')
			}

			if value&pkg.NegativeZero != value {
				return nil, truncationError
			}

			return NewSimpleIntegerValue(value), nil
		}
	}

	return nil, nil
}
