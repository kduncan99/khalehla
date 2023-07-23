// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
)

type operatorType int

const (
	AddOperator operatorType = iota
	AndOperator
	ConcatenationOperator
	DivideOperator
	DivideCoveredQuotientOperator
	DivideRemainderOperator
	DoublePrecisionOperator
	EqualOperator
	GreaterThanOperator
	GreaterThanOrEqualOperator
	LeftJustifyOperator
	LessThanOperator
	LessThanOrEqualOperator
	MultiplyOperator
	NegativeOperator
	NodeIdentityOperator
	NodeNonIdentityOperator
	NotOperator
	NotEqualOperator
	OrOperator
	RightJustifyOperator
	SinglePrecisionOperator
	SubtractOperator
	PositiveOperator
	XorOperator
)

var insufficientOperands = fmt.Errorf("insufficient operands for operator")
var wrongOperandType = fmt.Errorf("wrong operand type for unary operator")

type Operator interface {
	Evaluate(context *ExpressionContext) error
	GetPrecedence() int
	GetToken() string
	GetOperatorType() operatorType
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

var binaryOperators = []Operator{
	divideRemainderOperator{},
	nodeNonIdentityOperator{},
	notEqualOperator{},
	lessThanOrEqualOperator{},
	greaterThanOrEqualOperator{},
	orOperator{},
	xorOperator{},
	andOperator{},
	divideCoveredQuotientOperator{},
	nodeIdentityOperator{},
	equalOperator{},
	lessThanOperator{},
	greaterThanOperator{},
	addOperator{},
	multiplyOperator{},
	subtractOperator{},
	divideOperator{},
	concatenationOperator{},
}

var unaryPostfixOperators = []Operator{
	doublePrecisionOperator{},
	leftJustifyOperator{},
	rightJustifyOperator{},
	singlePrecisionOperator{},
}

var unaryPrefixOperators = []Operator{
	positiveOperator{},
	negativeOperator{},
	notOperator{},
}

func (p *Parser) ParseBinaryOperator() Operator {
	for _, op := range binaryOperators {
		result := p.ParseTokenCaseInsensitive(op.GetToken())
		if result {
			return op
		}
	}
	return nil
}

func (p *Parser) ParseUnaryPostfixOperator() Operator {
	for _, op := range unaryPostfixOperators {
		result := p.ParseTokenCaseInsensitive(op.GetToken())
		if result {
			return op
		}
	}
	return nil
}

func (p *Parser) ParseUnaryPrefixOperator() Operator {
	for _, op := range unaryPrefixOperators {
		result := p.ParseTokenCaseInsensitive(op.GetToken())
		if result {
			return op
		}
	}
	return nil
}

// add operator --------------------------------------------------------------------------------------------------------

func (op addOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op addOperator) GetOperatorType() operatorType {
	return AddOperator
}

func (op addOperator) GetPrecedence() int {
	return 6
}

func (op addOperator) GetToken() string {
	return "+"
}

// and operator --------------------------------------------------------------------------------------------------------

func (op andOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op andOperator) GetOperatorType() operatorType {
	return AndOperator
}

func (op andOperator) GetPrecedence() int {
	return 5
}

func (op andOperator) GetToken() string {
	return "**"
}

// concatenation operator ----------------------------------------------------------------------------------------------

func (op concatenationOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op concatenationOperator) GetOperatorType() operatorType {
	return ConcatenationOperator
}

func (op concatenationOperator) GetPrecedence() int {
	return 3
}

func (op concatenationOperator) GetToken() string {
	return ":"
}

// divide operator -----------------------------------------------------------------------------------------------------

func (op divideOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op divideOperator) GetOperatorType() operatorType {
	return DivideOperator
}

func (op divideOperator) GetPrecedence() int {
	return 7
}

func (op divideOperator) GetToken() string {
	return "/"
}

// divideCoveredQuotient operator --------------------------------------------------------------------------------------

func (op divideCoveredQuotientOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op divideCoveredQuotientOperator) GetOperatorType() operatorType {
	return DivideCoveredQuotientOperator
}

func (op divideCoveredQuotientOperator) GetPrecedence() int {
	return 7
}

func (op divideCoveredQuotientOperator) GetToken() string {
	return "//"
}

// divideRemainder operator --------------------------------------------------------------------------------------------

func (op divideRemainderOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op divideRemainderOperator) GetOperatorType() operatorType {
	return DivideRemainderOperator
}

func (op divideRemainderOperator) GetPrecedence() int {
	return 7
}

func (op divideRemainderOperator) GetToken() string {
	return "///"
}

// double precision operator -------------------------------------------------------------------------------------------

func (op doublePrecisionOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op doublePrecisionOperator) GetOperatorType() operatorType {
	return DoublePrecisionOperator
}

func (op doublePrecisionOperator) GetPrecedence() int {
	return 10
}

func (op doublePrecisionOperator) GetToken() string {
	return "D"
}

// equal operator ------------------------------------------------------------------------------------------------------

func (op equalOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op equalOperator) GetOperatorType() operatorType {
	return EqualOperator
}

func (op equalOperator) GetPrecedence() int {
	return 6
}

func (op equalOperator) GetToken() string {
	return "="
}

// greater than operator -----------------------------------------------------------------------------------------------

func (op greaterThanOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op greaterThanOperator) GetOperatorType() operatorType {
	return GreaterThanOperator
}

func (op greaterThanOperator) GetPrecedence() int {
	return 6
}

func (op greaterThanOperator) GetToken() string {
	return ">"
}

// greater than or equal operator --------------------------------------------------------------------------------------

func (op greaterThanOrEqualOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op greaterThanOrEqualOperator) GetOperatorType() operatorType {
	return GreaterThanOrEqualOperator
}

func (op greaterThanOrEqualOperator) GetPrecedence() int {
	return 6
}

func (op greaterThanOrEqualOperator) GetToken() string {
	return ">="
}

// left justify operator -----------------------------------------------------------------------------------------------

func (op leftJustifyOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op leftJustifyOperator) GetOperatorType() operatorType {
	return LeftJustifyOperator
}

func (op leftJustifyOperator) GetPrecedence() int {
	return 10
}

func (op leftJustifyOperator) GetToken() string {
	return "L"
}

// less than operator --------------------------------------------------------------------------------------------------

func (op lessThanOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op lessThanOperator) GetOperatorType() operatorType {
	return LessThanOperator
}

func (op lessThanOperator) GetPrecedence() int {
	return 6
}

func (op lessThanOperator) GetToken() string {
	return "<"
}

// less than or equal operator -----------------------------------------------------------------------------------------

func (op lessThanOrEqualOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op lessThanOrEqualOperator) GetOperatorType() operatorType {
	return LessThanOrEqualOperator
}

func (op lessThanOrEqualOperator) GetPrecedence() int {
	return 6
}

func (op lessThanOrEqualOperator) GetToken() string {
	return "<="
}

// multiply operator ---------------------------------------------------------------------------------------------------

func (op multiplyOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op multiplyOperator) GetOperatorType() operatorType {
	return MultiplyOperator
}

func (op multiplyOperator) GetPrecedence() int {
	return 7
}

func (op multiplyOperator) GetToken() string {
	return "*"
}

// negative operator ---------------------------------------------------------------------------------------------------

func (op negativeOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op negativeOperator) GetOperatorType() operatorType {
	return NegativeOperator
}

func (op negativeOperator) GetPrecedence() int {
	return 9
}

func (op negativeOperator) GetToken() string {
	return "-"
}

// node identity operator ----------------------------------------------------------------------------------------------

func (op nodeIdentityOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op nodeIdentityOperator) GetOperatorType() operatorType {
	return NodeIdentityOperator
}

func (op nodeIdentityOperator) GetPrecedence() int {
	return 6
}

func (op nodeIdentityOperator) GetToken() string {
	return "=="
}

// node non-identity operator ------------------------------------------------------------------------------------------

func (op nodeNonIdentityOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op nodeNonIdentityOperator) GetOperatorType() operatorType {
	return NodeNonIdentityOperator
}

func (op nodeNonIdentityOperator) GetPrecedence() int {
	return 6
}

func (op nodeNonIdentityOperator) GetToken() string {
	return "=/="
}

// not operator --------------------------------------------------------------------------------------------------------

func (op notOperator) Evaluate(context *ExpressionContext) error {
	v, err := context.PopValue()
	if err != nil {
		return insufficientOperands
	}

	if v.GetValueType() != IntegerValueType {
		return wrongOperandType
	}

	//	make a new integer value, keep the components but the component values must all be
	//	bit-flipped according to their lengths
	iv := v.(*IntegerValue)
	comps := make([]ValueComponent, len(iv.components))
	for cx := 0; cx < len(iv.components); cx++ {
		comps[cx] = iv.components[cx].not()
	}

	newValue := &IntegerValue{
		components: comps,
	}
	context.PushValue(newValue)
	return nil
}

func (op notOperator) GetOperatorType() operatorType {
	return NotOperator
}

func (op notOperator) GetPrecedence() int {
	return 1
}

func (op notOperator) GetToken() string {
	return "\\"
}

// not equal operator --------------------------------------------------------------------------------------------------

func (op notEqualOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op notEqualOperator) GetOperatorType() operatorType {
	return NotEqualOperator
}

func (op notEqualOperator) GetPrecedence() int {
	return 6
}

func (op notEqualOperator) GetToken() string {
	return "<>"
}

// or operator ---------------------------------------------------------------------------------------------------------

func (op orOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op orOperator) GetOperatorType() operatorType {
	return OrOperator
}

func (op orOperator) GetPrecedence() int {
	return 4
}

func (op orOperator) GetToken() string {
	return "++"
}

// positive operator ---------------------------------------------------------------------------------------------------

func (op positiveOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op positiveOperator) GetOperatorType() operatorType {
	return PositiveOperator
}

func (op positiveOperator) GetPrecedence() int {
	return 9
}

func (op positiveOperator) GetToken() string {
	return "+"
}

// right justify operator ----------------------------------------------------------------------------------------------

func (op rightJustifyOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op rightJustifyOperator) GetOperatorType() operatorType {
	return RightJustifyOperator
}

func (op rightJustifyOperator) GetPrecedence() int {
	return 10
}

func (op rightJustifyOperator) GetToken() string {
	return "R"
}

// single precision operator -------------------------------------------------------------------------------------------

func (op singlePrecisionOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op singlePrecisionOperator) GetOperatorType() operatorType {
	return SinglePrecisionOperator
}

func (op singlePrecisionOperator) GetPrecedence() int {
	return 10
}

func (op singlePrecisionOperator) GetToken() string {
	return "S"
}

// subtract operator ---------------------------------------------------------------------------------------------------

func (op subtractOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op subtractOperator) GetOperatorType() operatorType {
	return SubtractOperator
}

func (op subtractOperator) GetPrecedence() int {
	return 6
}

func (op subtractOperator) GetToken() string {
	return "-"
}

// xor operator --------------------------------------------------------------------------------------------------------

func (op xorOperator) Evaluate(context *ExpressionContext) error {
	return nil // TODO
}

func (op xorOperator) GetOperatorType() operatorType {
	return XorOperator
}

func (op xorOperator) GetPrecedence() int {
	return 4
}

func (op xorOperator) GetToken() string {
	return "--"
}
