// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"

	"khalehla/parser"
)

type ExpressionContext struct {
	context   *Context
	values    []Value
	operators []Operator
}

func NewExpressionContext(context *Context) *ExpressionContext {
	return &ExpressionContext{
		context:   context,
		values:    make([]Value, 0),
		operators: make([]Operator, 0),
	}
}

func (ec *ExpressionContext) PeekOperator() (Operator, error) {
	l := len(ec.operators)
	if l == 0 {
		return nil, parser.outOfData
	} else {
		op := ec.operators[l-1]
		return op, nil
	}
}

func (ec *ExpressionContext) PeekValue() (Value, error) {
	l := len(ec.values)
	if l == 0 {
		return nil, parser.outOfData
	} else {
		v := ec.values[l-1]
		return v, nil
	}
}

func (ec *ExpressionContext) PopOperator() (Operator, error) {
	l := len(ec.operators)
	if l == 0 {
		return nil, fmt.Errorf("internal error - value stack is empty")
	} else {
		op := ec.operators[l-1]
		ec.operators = ec.operators[:l-1]
		return op, nil
	}
}

func (ec *ExpressionContext) PopValue() (Value, error) {
	l := len(ec.values)
	if l == 0 {
		return nil, parser.outOfData
	} else {
		v := ec.values[l-1]
		ec.values = ec.values[:l-1]
		return v, nil
	}
}

func (ec *ExpressionContext) PopVariableParameterList() ([]Value, error) {
	count, err := ec.PopValue()
	if err != nil || count.GetValueType() != IntegerValueType {
		return nil, err
	}

	iCount := count.(*IntegerValue)
	if len(iCount.components) > 1 || len(iCount.components[0].offsets) > 0 {
		return nil, fmt.Errorf("data type or relocation error popping function parameter list")
	}

	pCount := iCount.components[0].value
	values := make([]Value, pCount)
	for vx := pCount; vx > 1; vx-- {
		values[vx-1], err = ec.PopValue()
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func (ec *ExpressionContext) PushOperator(op Operator) {
	ec.operators = append(ec.operators, op)
}

func (ec *ExpressionContext) PushValue(value Value) {
	ec.values = append(ec.values, value)
}

func (ec *ExpressionContext) PushVariableParameterList(values []Value) {
	// push the values left-to-right, then push the count
	for _, v := range values {
		ec.PushValue(v)
	}

	ec.PushValue(NewSimpleIntegerValue(uint64(len(values))))
}
