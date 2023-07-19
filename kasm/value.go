// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"math/big"
)

type ValueType int

const (
	IntegerValueType ValueType = iota
	StringValueType
)

type Value interface {
	GetValueType() ValueType
}

// Integer value -------------------------------------------------------------------------------------------------------

type ValueOffset struct {
	symbol     string
	isNegative bool
}

type ValueComponent struct {
	bitCount int
	value    *big.Int
	offsets  []ValueOffset
}

type IntegerValue struct {
	components []ValueComponent
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

func NewBigIntegerValue(value *big.Int) *IntegerValue {
	vc := ValueComponent{
		bitCount: 36,
		value:    value,
		offsets:  nil,
	}

	return &IntegerValue{
		components: []ValueComponent{vc},
	}
}

func NewSimpleIntegerValue(value int64) *IntegerValue {
	vc := ValueComponent{
		bitCount: 36,
		value:    big.NewInt(value),
		offsets:  nil,
	}

	return &IntegerValue{
		components: []ValueComponent{vc},
	}
}

// String value --------------------------------------------------------------------------------------------------------

type StringValue struct {
	value string
}

func (v *StringValue) GetValueType() ValueType {
	return StringValueType
}

func NewStringValue(value string) *StringValue {
	return &StringValue{
		value: value,
	}
}
