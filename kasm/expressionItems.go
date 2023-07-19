// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type expressionItemType int

const (
	ValueItemType expressionItemType = iota
	OperatorItemType
)

type expressionItem interface {
	Evaluate(expression *Expression) Value
	GetExpressionItemType() expressionItemType
}
