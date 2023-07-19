// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type operatorItemType int
type operatorPosition int

const (
	AddOperator operatorItemType = iota
	AndOperator
	DivideOperator
	MultiplyOperator
	NotOperator
	OrOperator
	SubtractOperator
)

const (
	BinaryOperator operatorPosition = iota
	UnaryPrefixOperator
	UnaryPostfixOperator
)

type operatorExpressionItem interface {
	Evaluate(expression *Expression) Value
	GetExpressionItemType() expressionItemType
	GetPrecedence() int
	GetOperatorPosition() operatorPosition
	GetToken() string
}

var Operators = []operatorExpressionItem{
	// order is important for parsing
	&andOperator{},
	&xorOperator{},
	&orOperator{},

	&addOperator{},
	// &concatenationOperator{},
	// &divideRemainderOperator{},
	// &divideCoveredQuotientOperator{},
	// &divideOperator{},

	// &nodeIdentityOperator{},
	// &nodeNonIdentityOperator{},

	// &equalOperator{},
	// &notEqualOperator{},
	// &greaterThanOrEqualOperator{},
	// &lessThanOrEqualOperator{},
	// &lessThanOperator{},
	// &greaterThanOperator{},

	// &multiplyOperator{},
	&negativeOperator{},
	&notOperator{},
	&subtractOperator{},
	&positiveOperator{},

	&doublePrecisionOperator{},
	&singlePrecisionOperator{},
	// &leftJustifyOperator{},
	// &rightJustifyOperator{},
}

type addOperator struct{}
type andOperator struct{}
type concatenationOperator struct{}
type divideOperator struct{}
type divideCoveredQuotientOperator struct{}
type divideRemainderOperator struct{}
type doublePrecisionOperator struct{}
type equalOperator struct{}
type greaterThanOperator struct{}
type greaterThanOrEqualOperator struct{}
type leftJustifyOperator struct{}
type lessThanOperator struct{}
type lessThanOrEqualOperator struct{}
type multiplyOperator struct{}
type negativeOperator struct{}
type nodeIdentityOperator struct{}
type nodeNonIdentityOperator struct{}
type notOperator struct{}
type notEqualOperator struct{}
type orOperator struct{}
type rightJustifyOperator struct{}
type singlePrecisionOperator struct{}
type subtractOperator struct{}
type positiveOperator struct{}
type xorOperator struct{}

//	TODO lots more operator functions below

// add operator --------------------------------------------------------------------------------------------------------

func (op *addOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *addOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *addOperator) GetOperatorPosition() operatorPosition {
	return BinaryOperator
}

func (op *addOperator) GetPrecedence() int {
	return 6
}

func (op *addOperator) GetToken() string {
	return "+"
}

// and operator --------------------------------------------------------------------------------------------------------

func (op *andOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *andOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *andOperator) GetOperatorPosition() operatorPosition {
	return BinaryOperator
}

func (op *andOperator) GetPrecedence() int {
	return 5
}

func (op *andOperator) GetToken() string {
	return "**"
}

// double operator -----------------------------------------------------------------------------------------------------

func (op *doublePrecisionOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *doublePrecisionOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *doublePrecisionOperator) GetOperatorPosition() operatorPosition {
	return UnaryPostfixOperator
}

func (op *doublePrecisionOperator) GetPrecedence() int {
	return 10
}

func (op *doublePrecisionOperator) GetToken() string {
	return "D"
}

// negative operator ---------------------------------------------------------------------------------------------------

func (op *negativeOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *negativeOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *negativeOperator) GetOperatorPosition() operatorPosition {
	return UnaryPrefixOperator
}

func (op *negativeOperator) GetPrecedence() int {
	return 9
}

func (op *negativeOperator) GetToken() string {
	return "-"
}

// not operator --------------------------------------------------------------------------------------------------------

func (op *notOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *notOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *notOperator) GetOperatorPosition() operatorPosition {
	return UnaryPrefixOperator
}

func (op *notOperator) GetPrecedence() int {
	return 1
}

func (op *notOperator) GetToken() string {
	return "\\"
}

// or operator ---------------------------------------------------------------------------------------------------------

func (op *orOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *orOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *orOperator) GetOperatorPosition() operatorPosition {
	return BinaryOperator
}

func (op *orOperator) GetPrecedence() int {
	return 4
}

func (op *orOperator) GetToken() string {
	return "++"
}

// positive operator ---------------------------------------------------------------------------------------------------

func (op *positiveOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *positiveOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *positiveOperator) GetOperatorPosition() operatorPosition {
	return UnaryPrefixOperator
}

func (op *positiveOperator) GetPrecedence() int {
	return 9
}

func (op *positiveOperator) GetToken() string {
	return "+"
}

// single operator -----------------------------------------------------------------------------------------------------

func (op *singlePrecisionOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *singlePrecisionOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *singlePrecisionOperator) GetOperatorPosition() operatorPosition {
	return UnaryPostfixOperator
}

func (op *singlePrecisionOperator) GetPrecedence() int {
	return 10
}

func (op *singlePrecisionOperator) GetToken() string {
	return "S"
}

// subtract operator ---------------------------------------------------------------------------------------------------

func (op *subtractOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *subtractOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *subtractOperator) GetOperatorPosition() operatorPosition {
	return BinaryOperator
}

func (op *subtractOperator) GetPrecedence() int {
	return 6
}

func (op *subtractOperator) GetToken() string {
	return "-"
}

// xor operator --------------------------------------------------------------------------------------------------------

func (op *xorOperator) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (op *xorOperator) GetExpressionItemType() expressionItemType {
	return OperatorItemType
}

func (op *xorOperator) GetOperatorPosition() operatorPosition {
	return BinaryOperator
}

func (op *xorOperator) GetPrecedence() int {
	return 4
}

func (op *xorOperator) GetToken() string {
	return "--"
}
