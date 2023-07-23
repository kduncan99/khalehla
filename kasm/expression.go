// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
)

// ExpressionItem is an Operator or a Value
type ExpressionItem interface {
	Evaluate(context *ExpressionContext) error
}

// Expression represents an evaluable expression
type Expression struct {
	items []ExpressionItem
}

type ExpressionList struct {
	expressions []*Expression
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

func (p *Parser) ParseExpression() (*Expression, error) {
	e := NewExpression()

	wantBinaryOperator := false
	wantUnaryPostfixOperator := false
	wantUnaryPrefixOperator := true
	wantTerm := true

	p.SkipWhiteSpace()
	for !p.AtEnd() {

		if wantUnaryPostfixOperator {
			op := p.ParseUnaryPostfixOperator()
			if op != nil {
				e.pushItem(op)
				p.SkipWhiteSpace()
				continue
			}
		}

		if wantUnaryPrefixOperator {
			op := p.ParseUnaryPrefixOperator()
			if op != nil {
				e.pushItem(op)
				p.SkipWhiteSpace()
				continue
			}
		}

		if wantBinaryOperator {
			op := p.ParseBinaryOperator()
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
			term, err := p.ParseTerm()
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
//	expr [ ',' expr ]*
//
// If we do not find any expression, we return nil as the result.
// Note that we DO NOT parse nor deal with enclosing parenthesis.
func (p *Parser) ParseExpressionList() (*ExpressionList, error) {
	pos := p.GetPosition()
	p.SkipWhiteSpace()
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	} else if expr == nil {
		_ = p.SetPosition(pos)
		return nil, nil
	}

	expList := []*Expression{expr}
	p.SkipWhiteSpace()
	for p.ParseCharacter(',') {
		p.SkipWhiteSpace()
		expr, err = p.ParseExpression()
		if err != nil {
			return nil, err
		} else if expr == nil {
			return nil, fmt.Errorf("expected another expression in expression list")
		}

		expList = append(expList, expr)
		p.SkipWhiteSpace()
	}

	result := &ExpressionList{
		expressions: expList,
	}
	return result, nil
}

func (el *ExpressionList) Evaluate(ec *ExpressionContext) error {
	// TODO an interesting conundrum - bare list gives us fields equally divided into 36, but what about forms?
}

// ParseTerm parses a term from the input text. A term is anything which is not an operator.
func (p *Parser) ParseTerm() (ExpressionItem, error) {
	result, err := p.ParseLiteral()
	// if result == nil && err == nil {
	// 	result, err = p.ParseReference()
	// }
	return result, err
}
