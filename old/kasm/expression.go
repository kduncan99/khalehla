// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"

	"khalehla/old/parser"
)

// ExpressionItem is an Operator or a Value
type ExpressionItem interface {
	Evaluate(context *ExpressionContext) error
}

// Expression represents an evaluable expression
type Expression struct {
	items []ExpressionItem
}

func NewExpression() *Expression {
	return &Expression{
		items: make([]ExpressionItem, 0),
	}
}

func (e *Expression) pushItem(item ExpressionItem) {
	e.items = append(e.items, item)
}

func (e *Expression) Evaluate(context *ExpressionContext) error {
	ix := 0
	for ix < len(e.items) {
		err := e.items[ix].Evaluate(context)
		if err != nil {
			return err
		}
	}

	if len(context.values) != 1 {
		return fmt.Errorf("internal expression evaluation error")
	}

	return nil
}

func ParseExpression(p *parser.Parser, context *Context) (*Expression, error) {
	e := NewExpression()

	// TODO This is a thing, which must be thing'ed somewhere (if not here):
	//	To detect a line item, MASM evaluates the first expression following a left parenthesis.
	//	If the next character after the expression is not a right parenthesis, the character must be a comma or a space.
	//	If it is a comma and no significant space (not following an operator) is found before the right parenthesis,
	//	an implicit call to $GEN is made. Otherwise, the first expression after the left parenthesis must be a directive,
	//	instruction, or procedure call. This means MASM recognizes the format (LA,U A0,1) as a valid line item.
	// So, yeah...

	wantBinaryOperator := false
	wantUnaryPostfixOperator := false
	wantUnaryPrefixOperator := true
	wantTerm := true

	p.SkipWhiteSpace()
	for !p.AtEnd() {

		if wantUnaryPostfixOperator {
			op := ParseUnaryPostfixOperator(p)
			if op != nil {
				e.pushItem(op)
				p.SkipWhiteSpace()
				continue
			}
		}

		if wantUnaryPrefixOperator {
			op := ParseUnaryPrefixOperator(p)
			if op != nil {
				e.pushItem(op)
				p.SkipWhiteSpace()
				continue
			}
		}

		if wantBinaryOperator {
			op := ParseBinaryOperator(p)
			if op != nil {
				wantBinaryOperator = false
				wantUnaryPostfixOperator = false
				wantUnaryPrefixOperator = true
				wantTerm = true
				e.pushItem(op)
				p.SkipWhiteSpace()
				continue
			}
		}

		if wantTerm {
			term, err := ParseTerm(p, context)
			if err != nil {
				return nil, err
			} else if term != nil {
				wantBinaryOperator = true
				wantUnaryPostfixOperator = true
				wantUnaryPrefixOperator = false
				wantTerm = false
				e.pushItem(term)
				p.SkipWhiteSpace()
				continue
			}
		}

		return nil, fmt.Errorf("syntax error in expression")
	}

	return e, nil
}

// ParseExpressionList parses the following into a slice of references to Expression structs:
//
//	'(' [ expr [ ',' expr ]* ] ')'
//
// If we do not find any expression, we return nil as the result.
// This is for function references and node selectors
func ParseExpressionList(p *parser.Parser, context *Context) ([]*Expression, error) {
	p.SkipWhiteSpace()
	if !p.ParseCharacter('(') {
		return nil, nil
	}

	expList := make([]*Expression, 0)
	p.SkipWhiteSpace()
	if p.ParseCharacter(')') {
		return expList, nil
	}

	exp, err := ParseExpression(p, context)
	if err != nil {
		return nil, err
	} else if exp == nil {
		return nil, fmt.Errorf("syntax error in function arguments or node selectors")
	}

	p.SkipWhiteSpace()
	for p.ParseCharacter(',') {
		p.SkipWhiteSpace()
		exp, err = ParseExpression(p, context)
		if err != nil {
			return nil, err
		} else if exp == nil {
			return nil, fmt.Errorf("syntax error in function arguments or node selectors")
		}

		expList = append(expList, exp)
		p.SkipWhiteSpace()
	}

	if !p.ParseCharacter('(') {
		return nil, fmt.Errorf("unterminated function arguments or node selectors")
	}

	return expList, nil
}

// ParseTerm parses a term from the input text. A term is anything which is not an operator.
func ParseTerm(p *parser.Parser, context *Context) (ExpressionItem, error) {
	result, err := ParseLiteral(p, context)

	if result == nil && err == nil {
		result, err = ParseReference(p, context)
	}

	//	TODO more alternatives...

	return result, err
}
