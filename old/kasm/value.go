// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"

	"khalehla/old/parser"
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
	MasmIntrinsicFunctionValueType // (built-in function)
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

func ParseLiteral(p *parser.Parser, context *Context) (Value, error) {
	value, err := ParseFloatLiteral(p)
	if value == nil && err == nil {
		value, err = ParseIntegerLiteral(p)
	}
	if value == nil && err == nil {
		value, err = ParseStringLiteral(p, context)
	}

	return value, err
}
