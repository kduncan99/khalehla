// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
	"strings"
)

// Expression represents an evaluable expression
type Expression struct {
	text  string
	items []expressionItem

	textIndex int
	itemIndex int
	context   *Context
	values    []Value
	operators []*operatorExpressionItem
}

func NewExpression(text string) *Expression {
	return &Expression{
		text:  text,
		items: make([]expressionItem, 0),
	}
}

func (e *Expression) evaluate(context *Context) (Value, error) {
	e.context = context
	e.itemIndex = 0
	e.values = make([]Value, 0)
	e.operators = make([]*operatorExpressionItem, 0)

	for e.itemIndex < len(e.items) {
		if e.items[e.itemIndex].GetExpressionItemType() == OperatorItemType {

		} else if e.items[e.itemIndex].GetExpressionItemType() == ValueItemType {
			vi := e.items[e.itemIndex].(valueExpressionItem)
			e.values = append(e.values, vi.Evaluate(e))
		}
	}

	if len(e.values) != 1 {
		return nil, fmt.Errorf("internal expression evaluation error")
	}

	val := e.values[0]
	e.values = nil
	return val, nil
}

func (e *Expression) skipWhitespace() {
	for e.textIndex < len(e.text) {
		if e.text[e.textIndex] != ' ' {
			break
		}
	}
}

func (e *Expression) parse() error {
	wantBinaryOperator := false
	wantUnaryPostfixOperator := false
	wantUnaryPrefixOperator := true
	wantValue := true

	for e.textIndex < len(e.text) {
		e.skipWhitespace()

		if wantUnaryPostfixOperator {
			result, err := e.parseUnaryPostfixOperator()
			if err != nil {
				return err
			} else if result {
				continue
			}
		}

		if wantUnaryPrefixOperator {
			result, err := e.parseUnaryPrefixOperator()
			if err != nil {
				return err
			} else if result {
				continue
			}
		}

		if wantBinaryOperator {
			result, err := e.parseBinaryOperator()
			if err != nil {
				return err
			} else if result {
				wantBinaryOperator = false
				wantUnaryPostfixOperator = false
				wantUnaryPrefixOperator = true
				wantValue = true
				continue
			}
		}

		if wantValue {
			result, err := e.parseValue()
			if err != nil {
				return err
			} else if result {
				wantBinaryOperator = true
				wantUnaryPostfixOperator = true
				wantUnaryPrefixOperator = false
				wantValue = false
				continue
			}
		}

		return fmt.Errorf("syntax error in expresion")
	}

	return nil
}

func (e *Expression) parseToken(token string) bool {
	remaining := len(e.text) - e.textIndex
	if remaining >= len(token) {
		textx := e.textIndex
		for tokx := 0; tokx < len(token); tokx++ {
			tokChar := strings.ToUpper(token[tokx : tokx+1])
			textChar := strings.ToUpper(e.text[textx : textx+1])
			if tokChar != textChar {
				return false
			}
		}

		e.textIndex += len(token)
		return true
	}

	return false
}

func (e *Expression) parseBinaryOperator() (bool, error) {
	for _, op := range Operators {
		if op.GetOperatorPosition() == BinaryOperator && e.parseToken(op.GetToken()) {
			e.items = append(e.items, op)
			return true, nil
		}
	}
	return false, nil
}

func (e *Expression) parseUnaryPostfixOperator() (bool, error) {
	for _, op := range Operators {
		if op.GetOperatorPosition() == UnaryPostfixOperator && e.parseToken(op.GetToken()) {
			e.items = append(e.items, op)
			return true, nil
		}
	}
	return false, nil
}

func (e *Expression) parseUnaryPrefixOperator() (bool, error) {
	for _, op := range Operators {
		if op.GetOperatorPosition() == UnaryPrefixOperator && e.parseToken(op.GetToken()) {
			e.items = append(e.items, op)
			return true, nil
		}
	}
	return false, nil
}

func (e *Expression) parseValue() (bool, error) {
	result, err := e.parseLiteralValue()
	if !result && err == nil {
		result, err = e.parseReference()
	}
	return result, err
}

func (e *Expression) parseLiteralValue() (bool, error) {
	result, err := e.parseFloatLiteralValue()
	if !result && err == nil {
		result, err = e.parseIntegerLiteralValue()
	}
	if !result && err == nil {
		result, err = e.parseStringLiteralValue()
	}

	return result, err
}

func (e *Expression) parseFloatLiteralValue() (bool, error) {
	return false, nil // TODO
}

func (e *Expression) parseIntegerLiteralValue() (bool, error) {
	return false, nil // TODO
}

func (e *Expression) parseStringLiteralValue() (bool, error) {
	return false, nil // TODO
}

func (e *Expression) parseReference() (bool, error) {
	return false, nil // TODO
}
