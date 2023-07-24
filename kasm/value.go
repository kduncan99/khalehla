// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
)

type ValueType int

const (
	IntegerValueType ValueType = iota + 1
	FloatValueType
	StringValueType
	NodeValueType
	InternalNameValueType          // (NAME line label)
	ProcedureValueType             // (label on PROC line)
	FunctionValueType              // (label on FUNC line)
	MasmDirectiveValueType         // (including instruction mnemonics and forms)
	MasmIntrinsicFunctionValueType // (build-in function)
)

type Value interface {
	Evaluate(ec *ExpressionContext) error
	GetValueType() ValueType
}

type BasicValue interface {
	ClearFlags(flags ValueFlags)
	Copy() Value
	Evaluate(ec *ExpressionContext) error
	GetFlags() ValueFlags
	GetValueType() ValueType
	SetFlags(flags ValueFlags)
}

var nonOctalDigit = fmt.Errorf("non-octal digit in numeric literal")
var syntaxError = fmt.Errorf("syntax error in numeric literal")
var truncationError = fmt.Errorf("truncation in numeric literal")

func (p *Parser) ParseLiteral(context *Context) (Value, error) {
	value, err := p.ParseFloatLiteral()
	if value == nil && err == nil {
		value, err = p.ParseIntegerLiteral()
	}
	if value == nil && err == nil {
		value, err = p.ParseStringLiteral(context)
	}

	return value, err
}
