// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"

	"khalehla/old/parser"
)

// LocationCounterSpecification is used to house a parsed but not-yet-executed expression
// describing a location counter.
type LocationCounterSpecification struct {
	expression *Expression
}

// NewLocationCounterSpecification creates a new LocationCounterSpecification struct
// given the text from source code which is in the following format:
//
//	'$(' expression ')'
//
// If successful, we return a reference to a LocationCounterSpecification struct.
// If the given text does not fit the format specification, we return nil with no error.
// If it *does* fit the general format but an error exists in the syntax of the expression,
// we return nil with an error.
func NewLocationCounterSpecification(context *Context, text string) (*LocationCounterSpecification, error) {
	p := parser.NewParser(text)
	if p.ParseToken("$(") {
		p.SkipWhiteSpace()
		exp, err := ParseExpression(p, context)
		if err != nil {
			return nil, err
		}

		p.SkipWhiteSpace()
		if !p.ParseCharacter(')') {
			return nil, fmt.Errorf("unterminated location counter specification")
		}

		lcn := &LocationCounterSpecification{
			expression: exp,
		}

		return lcn, nil
	}

	return nil, nil
}

func (lcs *LocationCounterSpecification) Evaluate(context *Context) (int, error) {
	ec := NewExpressionContext(context)
	err := lcs.expression.Evaluate(ec)
	if err != nil {
		return 0, err
	}

	val, err := ec.PopValue()
	if err != nil {
		return 0, err
	}

	if val.GetValueType() != IntegerValueType {
		return 0, fmt.Errorf("wrong value type for location counter specification")
	}

	iVal := val.(*IntegerValue)
	if !iVal.form.Equals(SimpleForm) || len(iVal.componentValues) != 1 {
		return 0, fmt.Errorf("invalid value for location counter specification")
	}

	comp := iVal.componentValues[0]
	if comp > 63 {
		return 0, fmt.Errorf("invalid value for location counter specification")
	}

	return int(comp), nil
}
