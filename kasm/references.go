// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
	"khalehla/parser"
)

type Reference interface {
	Evaluate(ec *ExpressionContext) error
	GetValueType() ValueType
}

type FunctionReference struct {
	arguments []*Expression
	function  Function
}

type NodeReference struct {
	selectors []*Expression
	node      *NodeValue
}

type ValueReference struct {
	value BasicValue
}

type UndefinedReference struct {
	symbol string
	value  *IntegerValue
}

func NewFunctionReference(arguments []*Expression, function Function) *FunctionReference {
	return &FunctionReference{
		arguments: arguments,
		function:  function,
	}
}

func NewNodeReference(selectors []*Expression, node *NodeValue) *NodeReference {
	return &NodeReference{
		selectors: selectors,
		node:      node,
	}
}

func NewUndefinedReference(symbol string) *UndefinedReference {
	offset := NewUndefinedReferenceOffset(symbol)
	value, _ := NewIntegerValue([]int64{0}, SimpleForm, []Offset{offset}, 0)
	return &UndefinedReference{
		symbol: symbol,
		value:  value,
	}
}

func NewValueReference(value BasicValue) *ValueReference {
	return &ValueReference{
		value: value,
	}
}

func (r *FunctionReference) Evaluate(ec *ExpressionContext) error {
	return r.function.Evaluate(ec)
}

func (r *FunctionReference) GetValueType() ValueType {
	return FunctionValueType
}

func (r *NodeReference) Evaluate(ec *ExpressionContext) error {
	return r.node.Evaluate(ec)
}

func (r *NodeReference) GetValueType() ValueType {
	return NodeValueType
}

func (r *UndefinedReference) Evaluate(ec *ExpressionContext) error {
	return r.value.Evaluate(ec)
}

func (r *UndefinedReference) GetValueType() ValueType {
	return IntegerValueType
}

func (r *ValueReference) Evaluate(ec *ExpressionContext) error {
	return r.value.Evaluate(ec)
}

func (r *ValueReference) GetValueType() ValueType {
	return r.value.GetValueType()
}

func ParseReference(p *parser.Parser, context *Context) (Reference, error) {
	pos := p.GetPosition()
	p.SkipWhiteSpace()
	symbol, err := p.ParseSymbol()
	if err != nil {
		return nil, err
	} else if symbol == nil {
		_ = p.SetPosition(pos)
		return nil, nil
	}

	entry, _ := context.dictionary.Lookup(*symbol)
	if entry == nil {
		return NewUndefinedReference(*symbol), nil
	}

	if entry.GetValueType() == FunctionValueType {
		expList, err := ParseExpressionList(p, context)
		if err != nil {
			return nil, err
		}
		return NewFunctionReference(expList, entry.(Function)), nil
	} else if entry.GetValueType() == NodeValueType {
		expList, err := ParseExpressionList(p, context)
		if err != nil {
			return nil, err
		}
		return NewNodeReference(expList, entry.(*NodeValue)), nil
	} else if entry.GetValueType() == IntegerValueType ||
		entry.GetValueType() == StringValueType ||
		entry.GetValueType() == FloatValueType {
		return NewValueReference(entry.(BasicValue)), nil
	} else {
		return nil, fmt.Errorf("improper reference item for expression")
	}
}
