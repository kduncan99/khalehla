// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
	"khalehla/parser"
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
	NotOperator
	NotEqualOperator
	OrOperator
	RightJustifyOperator
	SinglePrecisionOperator
	SubtractOperator
	PositiveOperator
	XorOperator
)

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

// TODO fixedIntegerScalingOperator
type greaterThanOperator struct{}
type greaterThanOrEqualOperator struct{}

type leadingAsteriskOperator struct{}
type leftJustifyOperator struct{}
type lessThanOperator struct{}
type lessThanOrEqualOperator struct{}
type multiplyOperator struct{}
type negativeOperator struct{}
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
	notEqualOperator{},
	lessThanOrEqualOperator{},
	greaterThanOrEqualOperator{},
	orOperator{},
	xorOperator{},
	andOperator{},
	divideCoveredQuotientOperator{},
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
	leadingAsteriskOperator{},
	positiveOperator{},
	negativeOperator{},
	notOperator{},
}

var divideByZeroError = fmt.Errorf("divide by zero error")

func ParseBinaryOperator(p *parser.Parser) Operator {
	for _, op := range binaryOperators {
		result := p.ParseTokenCaseInsensitive(op.GetToken())
		if result {
			return op
		}
	}
	return nil
}

func ParseUnaryPostfixOperator(p *parser.Parser) Operator {
	for _, op := range unaryPostfixOperators {
		result := p.ParseTokenCaseInsensitive(op.GetToken())
		if result {
			return op
		}
	}
	return nil
}

func ParseUnaryPrefixOperator(p *parser.Parser) Operator {
	for _, op := range unaryPrefixOperators {
		result := p.ParseTokenCaseInsensitive(op.GetToken())
		if result {
			return op
		}
	}
	return nil
}

// popArithmeticOperand pops one operand from the value stack, to be used with a unary arithmetic operator.
// The requirements are that the operand is either an integer or a float.
func popArithmeticOperand(context *ExpressionContext) (BasicValue, error) {
	op, err := context.PopValue()
	if err != nil {
		return nil, err
	}

	if op.GetValueType() == IntegerValueType || op.GetValueType() == FloatValueType {
		return op.(BasicValue), nil
	}

	return nil, fmt.Errorf("unary arithmetic operator must be an integer or a float")
}

// popArithmeticOperands pops left and right hand operands, and ensures that they are suitable for
// a binary arithmetic operator. The requirements are:
//
//	They must both be float values, *or*
//	They must both be integer values with equal forms, *or*
//	One must be a float, and one must be an integer with a simple form and no offsets
func popArithmeticOperands(context *ExpressionContext) (BasicValue, BasicValue, error) {
	rhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	lhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	rhsType := rhs.GetValueType()
	lhsType := lhs.GetValueType()
	if lhsType == FloatValueType && rhsType == FloatValueType {
		return lhs.(BasicValue), rhs.(BasicValue), nil
	}

	if lhsType == IntegerValueType && rhsType == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		if !lhsInt.form.Equals(rhsInt.form) {
			return nil, nil, fmt.Errorf("cannot compare integers with differing forms")
		}
		return lhs.(BasicValue), rhs.(BasicValue), nil
	}

	if lhsType == FloatValueType && rhsType == IntegerValueType {
		rhsFloat, err := NewFloatValueFromInteger(rhs.(*IntegerValue))
		if err != nil {
			return nil, nil, fmt.Errorf("integer operand for comparison is composite or has offsets")
		} else {
			return lhs.(BasicValue), rhsFloat, nil
		}
	} else if lhsType == IntegerValueType && rhsType == FloatValueType {
		lhsFloat, err := NewFloatValueFromInteger(lhs.(*IntegerValue))
		if err != nil {
			return nil, nil, fmt.Errorf("integer operand for comparison is composite or has offsets")
		} else {
			return lhsFloat, rhs.(BasicValue), nil
		}
	}

	return nil, nil, fmt.Errorf("operands are not comparable")
}

// popBasicOperand pops a single operand from the value stack, and ensures it is a BasicValue...
// that is, an integer, float, or string.
func popBasicOperand(context *ExpressionContext) (BasicValue, error) {
	op, err := context.PopValue()
	if err != nil {
		return nil, err
	}

	if op.GetValueType() == IntegerValueType || op.GetValueType() == FloatValueType || op.GetValueType() == StringValueType {
		return op.(BasicValue), nil
	} else {
		return nil, fmt.Errorf("invalid basic value operand - must be integer, float, or string")
	}
}

// popComparableOperands pops left and right hand operands, and ensures that they can be
// compared to each other in terms of equal or not equal.
//
//	They must both be string values, *or*
//	They must both be float values, *or*
//	They must both be integer values with equal forms and offsets, *or*
//	One must be a float, and one must be an integer with a simple form and no offsets
func popComparableOperands(context *ExpressionContext) (BasicValue, BasicValue, error) {
	rhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	lhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	rhsType := rhs.GetValueType()
	lhsType := lhs.GetValueType()
	if (lhsType == StringValueType && rhsType == StringValueType) ||
		lhsType == FloatValueType && rhsType == FloatValueType {
		return lhs.(BasicValue), rhs.(BasicValue), nil
	}

	if lhsType == IntegerValueType && rhsType == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		if !lhsInt.form.Equals(rhsInt.form) {
			return nil, nil, fmt.Errorf("cannot compare integers with differing forms")
		}
		if !OffsetListsAreEqual(lhsInt.offsets, rhsInt.offsets) {
			return nil, nil, fmt.Errorf("cannot compare integers with differing offsets")
		}

		return lhsInt, rhsInt, nil
	}

	if lhsType == FloatValueType && rhsType == IntegerValueType {
		rhsFloat, err := NewFloatValueFromInteger(rhs.(*IntegerValue))
		if err != nil {
			return nil, nil, fmt.Errorf("integer operand for comparison is composite or has offsets")
		} else {
			return lhs.(BasicValue), rhsFloat, nil
		}
	} else if lhsType == IntegerValueType && rhsType == FloatValueType {
		lhsFloat, err := NewFloatValueFromInteger(lhs.(*IntegerValue))
		if err != nil {
			return nil, nil, fmt.Errorf("integer operand for comparison is composite or has offsets")
		} else {
			return lhsFloat, rhs.(BasicValue), nil
		}
	}

	return nil, nil, fmt.Errorf("operands are not comparable")
}

// popStrictlyComparableOperands pops left and right hand operands, and ensures that they can be
// compared to each other in terms of greater, less, or equal. The requirements are:
//
//	They must both be string values, *or*
//	They must both be float values, *or*
//	They must both be integer values, or one integer and one float value.
//
// Integer values must have only a simple form (one component of 36 bits) and cannot have undefined offsets.
// If one is an integer and the other is a float, the integer is converted to a float when returned.
func popStrictlyComparableOperands(context *ExpressionContext) (BasicValue, BasicValue, error) {
	rhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	lhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	rhsType := rhs.GetValueType()
	lhsType := lhs.GetValueType()
	if (lhsType == StringValueType && rhsType == StringValueType) ||
		lhsType == FloatValueType && rhsType == FloatValueType {
		return lhs.(BasicValue), rhs.(BasicValue), nil
	}

	if lhsType == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		if !lhsInt.form.Equals(SimpleForm) || len(lhsInt.offsets) > 0 {
			return nil, nil, fmt.Errorf("LH operand is not comparable")
		}
	}

	if rhsType == IntegerValueType {
		rhsInt := rhs.(*IntegerValue)
		if !rhsInt.form.Equals(SimpleForm) || len(rhsInt.offsets) > 0 {
			return nil, nil, fmt.Errorf("RH operand is not comparable")
		}
	}

	if lhsType == IntegerValueType && rhsType == IntegerValueType {
		return lhs.(BasicValue), rhs.(BasicValue), nil
	}

	if lhsType == IntegerValueType && rhsType == FloatValueType {
		lhsNew, _ := NewFloatValueFromInteger(lhs.(*IntegerValue))
		return lhsNew, rhs.(BasicValue), nil
	} else if lhsType == FloatValueType && rhsType == IntegerValueType {
		rhsNew, _ := NewFloatValueFromInteger(rhs.(*IntegerValue))
		return lhs.(BasicValue), rhsNew, nil
	}

	return nil, nil, fmt.Errorf("operands are not strictly comparable")
}

// popLogicalOperand pops a value from the context and ensures that it can participate as a unary logical operand.
// The requirements are that it must be an integer value.
func popLogicalOperand(context *ExpressionContext) (*IntegerValue, error) {
	op, err := context.PopValue()
	if err != nil {
		return nil, err
	}

	if op.GetValueType() != IntegerValueType {
		return nil, fmt.Errorf("logical operand must be an integer")
	}

	return op.(*IntegerValue), nil
}

// popLogicalOperands pops two values from the context and ensures that they can participate as the left and
// right hand operands for a logical binary operator. The requirements are:
//
//	They must both be integer values
//	They may have multiple components, but their forms must be equal
//	At least one of them must not have any undefined offsets
func popLogicalOperands(context *ExpressionContext) (*IntegerValue, *IntegerValue, error) {
	rhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	lhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	if lhs.GetValueType() != IntegerValueType || rhs.GetValueType() != IntegerValueType {
		return nil, nil, fmt.Errorf("logical operands must be integers")
	}

	lhsInt := lhs.(*IntegerValue)
	rhsInt := rhs.(*IntegerValue)
	if (lhsInt.form == nil && rhsInt.form == nil) ||
		(lhsInt.form != nil && rhsInt.form != nil && lhsInt.form.Equals(rhsInt.form)) {
	} else {
		return nil, nil, fmt.Errorf("logical operands have unequal forms")
	}

	if len(lhsInt.offsets) > 0 && len(rhsInt.offsets) > 0 {
		return nil, nil, fmt.Errorf("at least one logical operand must have no offsets attached")
	}

	return lhsInt, rhsInt, nil
}

func popStringOperands(context *ExpressionContext) (*StringValue, *StringValue, error) {
	rhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	lhs, err := context.PopValue()
	if err != nil {
		return nil, nil, err
	}

	if lhs.GetValueType() != StringValueType || rhs.GetValueType() != StringValueType {
		return nil, nil, fmt.Errorf("expected string operands")
	}

	return lhs.(*StringValue), rhs.(*StringValue), nil
}

var boolToInt = map[bool]int64{
	true:  1,
	false: 0,
}

func mergeOffsets(off1 []Offset, off2 []Offset) []Offset {
	offsets := make([]Offset, len(off1)+len(off2))
	ox := 0
	for _, offset := range off1 {
		offsets[ox] = offset
		ox++
	}
	for _, offset := range off2 {
		offsets[ox] = offset
		ox++
	}

	return CollapseOffsetList(offsets)
}

// add operator --------------------------------------------------------------------------------------------------------

func (op addOperator) Evaluate(context *ExpressionContext) error {
	lhs, rhs, err := popArithmeticOperands(context)
	if err != nil {
		return err
	}

	if lhs.GetValueType() == FloatValueType {
		value := NewFloatValue(lhs.(*FloatValue).value + rhs.(*FloatValue).value)
		context.PushValue(value)
		return nil
	} else if lhs.GetValueType() == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		values := make([]int64, len(lhsInt.componentValues))
		for vx := 0; vx < len(lhsInt.componentValues); vx++ {
			values[vx] = lhsInt.componentValues[vx] + rhsInt.componentValues[vx]
		}

		offsets := CollapseOffsetList(append(lhsInt.offsets, rhsInt.offsets...))
		lcCount := 0
		for _, offset := range offsets {
			if offset.GetOffsetType() == LocationCounterOffsetType {
				lcCount++
			}
		}

		if lcCount > 1 {
			return fmt.Errorf("result of operation produces incompatible LC offset references")
		}

		newValue, _ := NewIntegerValue(values, lhsInt.form, offsets, 0)
		context.PushValue(newValue)
		return nil
	} else {
		return fmt.Errorf("internal error")
	}
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
	lhs, rhs, err := popLogicalOperands(context)
	if err != nil {
		return err
	}

	values := make([]int64, len(lhs.componentValues))
	for vx := 0; vx < len(lhs.componentValues); vx++ {
		values[vx] = lhs.componentValues[vx] & rhs.componentValues[vx]
	}

	value, err := NewIntegerValue(values, lhs.form, mergeOffsets(lhs.offsets, rhs.offsets), 0)
	if err != nil {
		return err
	}
	context.PushValue(value)
	return nil
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
	lhs, rhs, err := popStringOperands(context)
	if err != nil {
		return err
	}

	var codeType StringCodeType
	if lhs.codeType == AsciiString || rhs.codeType == AsciiString {
		codeType = AsciiString
	} else {
		codeType = FieldataString
	}

	value := NewStringValue(lhs.value+rhs.value, codeType, 0)
	context.PushValue(value)
	return nil
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
	lhs, rhs, err := popArithmeticOperands(context)
	if err != nil {
		return err
	}

	if lhs.GetValueType() == FloatValueType {
		rhsFloat := rhs.(*FloatValue)
		if rhsFloat.value == 0.0 {
			return divideByZeroError
		}

		value := NewFloatValue(lhs.(*FloatValue).value / rhs.(*FloatValue).value)
		context.PushValue(value)
		return nil
	} else if lhs.GetValueType() == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		if !lhsInt.form.Equals(SimpleForm) {
			return fmt.Errorf("cannot divide operands with composite values")
		}

		if len(lhsInt.offsets) > 0 || len(rhsInt.offsets) > 0 {
			return fmt.Errorf("cannot divide operands with undefined offsets")
		}

		if rhsInt.componentValues[0] == 0 {
			return divideByZeroError
		}

		value := NewSimpleIntegerValue(lhsInt.componentValues[0] / rhsInt.componentValues[0])
		context.PushValue(value)
		return nil
	} else {
		return fmt.Errorf("internal error")
	}
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
	lhs, rhs, err := popArithmeticOperands(context)
	if err != nil {
		return err
	}

	if lhs.GetValueType() == FloatValueType {
		return fmt.Errorf("invalid operand type for divide covered quotient")
	} else if lhs.GetValueType() == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		if !lhsInt.form.Equals(SimpleForm) {
			return fmt.Errorf("cannot divide operands with composite values")
		}

		if len(lhsInt.offsets) > 0 || len(rhsInt.offsets) > 0 {
			return fmt.Errorf("cannot divide operands with undefined offsets")
		}

		if rhsInt.componentValues[0] == 0 {
			return divideByZeroError
		}

		result := lhsInt.componentValues[0] / rhsInt.componentValues[0]
		if lhsInt.componentValues[0]%rhsInt.componentValues[0] != 0 {
			result++
		}

		value := NewSimpleIntegerValue(result)
		context.PushValue(value)
		return nil
	} else {
		return fmt.Errorf("internal error")
	}
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
	lhs, rhs, err := popArithmeticOperands(context)
	if err != nil {
		return err
	}

	if lhs.GetValueType() == FloatValueType {
		return fmt.Errorf("invalid operand type for divide remainder")
	} else if lhs.GetValueType() == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		if !lhsInt.form.Equals(SimpleForm) {
			return fmt.Errorf("cannot divide operands with composite values")
		}

		if len(lhsInt.offsets) > 0 || len(rhsInt.offsets) > 0 {
			return fmt.Errorf("cannot divide operands with undefined offsets")
		}

		if rhsInt.componentValues[0] == 0 {
			return divideByZeroError
		}

		result := lhsInt.componentValues[0] % rhsInt.componentValues[0]
		value := NewSimpleIntegerValue(result)
		context.PushValue(value)
		return nil
	} else {
		return fmt.Errorf("internal error")
	}
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
	val, err := popBasicOperand(context)
	if err != nil {
		return err
	}

	var basic BasicValue
	if val.GetValueType() == IntegerValueType {
		basic = val.Copy().(*IntegerValue)
	} else if val.GetValueType() == StringValueType {
		basic = val.Copy().(*StringValue)
	} else if val.GetValueType() == FloatValueType {
		basic = val.Copy().(*FloatValue)
	}

	basic.SetFlags(DoubleFlag)
	context.PushValue(basic)
	return nil
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
	lhs, rhs, err := popComparableOperands(context)
	if err != nil {
		return nil
	}

	var result int64
	if lhs.GetValueType() == StringValueType {
		result = boolToInt[lhs.(*StringValue).value == rhs.(*StringValue).value]
	} else if lhs.GetValueType() == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		result = 1
		for x := 0; x < len(lhsInt.componentValues); x++ {
			if lhsInt.componentValues[x] != rhsInt.componentValues[x] {
				result = 0
				break
			}
		}
	} else if lhs.GetValueType() == FloatValueType {
		result = boolToInt[lhs.(*FloatValue).value == lhs.(*FloatValue).value]
	}

	context.PushValue(NewSimpleIntegerValue(result))
	return nil
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
	lhs, rhs, err := popStrictlyComparableOperands(context)
	if err != nil {
		return err
	}

	var result int64
	if lhs.GetValueType() == StringValueType {
		result = boolToInt[lhs.(*StringValue).value > rhs.(*StringValue).value]
	} else if lhs.GetValueType() == IntegerValueType {
		result = boolToInt[lhs.(*IntegerValue).componentValues[0] > rhs.(*IntegerValue).componentValues[0]]
	} else if lhs.GetValueType() == FloatValueType {
		result = boolToInt[lhs.(*FloatValue).value > lhs.(*FloatValue).value]
	}

	context.PushValue(NewSimpleIntegerValue(result))
	return nil
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
	lhs, rhs, err := popStrictlyComparableOperands(context)
	if err != nil {
		return err
	}

	var result int64
	if lhs.GetValueType() == StringValueType {
		result = boolToInt[lhs.(*StringValue).value >= rhs.(*StringValue).value]
	} else if lhs.GetValueType() == IntegerValueType {
		result = boolToInt[lhs.(*IntegerValue).componentValues[0] >= rhs.(*IntegerValue).componentValues[0]]
	} else if lhs.GetValueType() == FloatValueType {
		result = boolToInt[lhs.(*FloatValue).value >= lhs.(*FloatValue).value]
	}

	context.PushValue(NewSimpleIntegerValue(result))
	return nil
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

// leading asterisk operator -------------------------------------------------------------------------------------------

func (op leadingAsteriskOperator) Evaluate(context *ExpressionContext) error {
	val, err := popBasicOperand(context)
	if err != nil {
		return err
	}

	var basic BasicValue
	if val.GetValueType() == IntegerValueType {
		basic = val.Copy().(*IntegerValue)
	} else if val.GetValueType() == StringValueType {
		basic = val.Copy().(*StringValue)
	} else if val.GetValueType() == FloatValueType {
		basic = val.Copy().(*FloatValue)
	}

	basic.SetFlags(FlaggedFlag)
	context.PushValue(basic)
	return nil
}

func (op leadingAsteriskOperator) GetOperatorType() operatorType {
	return SinglePrecisionOperator
}

func (op leadingAsteriskOperator) GetPrecedence() int {
	return 10
}

func (op leadingAsteriskOperator) GetToken() string {
	return "*"
}

// left justify operator -----------------------------------------------------------------------------------------------

func (op leftJustifyOperator) Evaluate(context *ExpressionContext) error {
	val, err := popBasicOperand(context)
	if err != nil {
		return err
	}

	var basic BasicValue
	if val.GetValueType() == IntegerValueType {
		basic = val.Copy().(*IntegerValue)
	} else if val.GetValueType() == StringValueType {
		basic = val.Copy().(*StringValue)
	} else if val.GetValueType() == FloatValueType {
		basic = val.Copy().(*FloatValue)
	}

	basic.SetFlags(LeftJustifiedFlag)
	context.PushValue(basic)
	return nil
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
	lhs, rhs, err := popStrictlyComparableOperands(context)
	if err != nil {
		return err
	}

	var result int64
	if lhs.GetValueType() == StringValueType {
		result = boolToInt[lhs.(*StringValue).value < rhs.(*StringValue).value]
	} else if lhs.GetValueType() == IntegerValueType {
		result = boolToInt[lhs.(*IntegerValue).componentValues[0] < rhs.(*IntegerValue).componentValues[0]]
	} else if lhs.GetValueType() == FloatValueType {
		result = boolToInt[lhs.(*FloatValue).value < lhs.(*FloatValue).value]
	}

	context.PushValue(NewSimpleIntegerValue(result))
	return nil
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
	lhs, rhs, err := popStrictlyComparableOperands(context)
	if err != nil {
		return err
	}

	var result int64
	if lhs.GetValueType() == StringValueType {
		result = boolToInt[lhs.(*StringValue).value <= rhs.(*StringValue).value]
	} else if lhs.GetValueType() == IntegerValueType {
		// comparable operands are guaranteed to be simple and with no undefined offsets
		result = boolToInt[lhs.(*IntegerValue).componentValues[0] <= rhs.(*IntegerValue).componentValues[0]]
	} else if lhs.GetValueType() == FloatValueType {
		result = boolToInt[lhs.(*FloatValue).value <= lhs.(*FloatValue).value]
	}

	context.PushValue(NewSimpleIntegerValue(result))
	return nil
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
	lhs, rhs, err := popArithmeticOperands(context)
	if err != nil {
		return err
	}

	if lhs.GetValueType() == FloatValueType {
		value := NewFloatValue(lhs.(*FloatValue).value * rhs.(*FloatValue).value)
		context.PushValue(value)
		return nil
	} else if lhs.GetValueType() == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		if !lhsInt.form.Equals(SimpleForm) {
			return fmt.Errorf("cannot multiply operands with composite values")
		}

		if len(lhsInt.offsets) > 0 || len(rhsInt.offsets) > 0 {
			return fmt.Errorf("cannot multiply operands with undefined offsets")
		}

		value := NewSimpleIntegerValue(lhsInt.componentValues[0] * rhsInt.componentValues[0])
		context.PushValue(value)
		return nil
	} else {
		return fmt.Errorf("internal error")
	}
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
	rhs, err := popArithmeticOperand(context)
	if err != nil {
		return err
	}

	if rhs.GetValueType() == FloatValueType {
		context.PushValue(NewFloatValue(0 - rhs.(*FloatValue).value))
	} else if rhs.GetValueType() == IntegerValueType {
		value := rhs.(*IntegerValue).Copy().(*IntegerValue)
		for vx := 0; vx < len(value.componentValues); vx++ {
			value.componentValues[vx] = -value.componentValues[vx]
		}
		context.PushValue(value)
	} else {
		return fmt.Errorf("internal error")
	}

	return nil
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

// not operator --------------------------------------------------------------------------------------------------------

func (op notOperator) Evaluate(context *ExpressionContext) error {
	rhs, err := popLogicalOperand(context)
	if err != nil {
		return err
	}

	value := rhs.Copy().(*IntegerValue)
	for vx := 0; vx < len(value.componentValues); vx++ {
		// At first glance, this seems tricky. We are applying a not in a ones-complement semantic sense,
		// but we are storing the value internally as twos-complement.
		// However, ones-complement logical NOT is equivalent to arithmetic negation, so all we have to do
		// is take the negative value of the components.
		value.componentValues[vx] = -value.componentValues[vx]
	}

	context.PushValue(value)
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
	lhs, rhs, err := popComparableOperands(context)
	if err != nil {
		return nil
	}

	var result int64
	if lhs.GetValueType() == StringValueType {
		result = boolToInt[lhs.(*StringValue).value != rhs.(*StringValue).value]
	} else if lhs.GetValueType() == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		result = 0
		for x := 0; x < len(lhsInt.componentValues); x++ {
			if lhsInt.componentValues[x] != rhsInt.componentValues[x] {
				result = 1
				break
			}
		}
	} else if lhs.GetValueType() == FloatValueType {
		result = boolToInt[lhs.(*FloatValue).value != lhs.(*FloatValue).value]
	}

	context.PushValue(NewSimpleIntegerValue(result))
	return nil
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
	lhs, rhs, err := popLogicalOperands(context)
	if err != nil {
		return err
	}

	values := make([]int64, len(lhs.componentValues))
	for vx := 0; vx < len(lhs.componentValues); vx++ {
		values[vx] = lhs.componentValues[vx] | rhs.componentValues[vx]
	}

	value, err := NewIntegerValue(values, lhs.form, mergeOffsets(lhs.offsets, rhs.offsets), 0)
	if err != nil {
		return err
	}
	context.PushValue(value)
	return nil
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
	rhs, err := popArithmeticOperand(context)
	if err != nil {
		return err
	}

	if rhs.GetValueType() == FloatValueType {
		context.PushValue(rhs.(*FloatValue).Copy())
	} else if rhs.GetValueType() == IntegerValueType {
		context.PushValue(rhs.(*IntegerValue).Copy())
	} else {
		return fmt.Errorf("internal error")
	}

	return nil
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
	val, err := popBasicOperand(context)
	if err != nil {
		return err
	}

	var basic BasicValue
	if val.GetValueType() == IntegerValueType {
		basic = val.Copy().(*IntegerValue)
	} else if val.GetValueType() == StringValueType {
		basic = val.Copy().(*StringValue)
	} else if val.GetValueType() == FloatValueType {
		basic = val.Copy().(*FloatValue)
	}

	basic.SetFlags(RightJustifiedFlag)
	context.PushValue(basic)
	return nil
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
	val, err := popBasicOperand(context)
	if err != nil {
		return err
	}

	var basic BasicValue
	if val.GetValueType() == IntegerValueType {
		basic = val.Copy().(*IntegerValue)
	} else if val.GetValueType() == StringValueType {
		basic = val.Copy().(*StringValue)
	} else if val.GetValueType() == FloatValueType {
		basic = val.Copy().(*FloatValue)
	}

	basic.SetFlags(SingleFlag)
	context.PushValue(basic)
	return nil
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
	lhs, rhs, err := popArithmeticOperands(context)
	if err != nil {
		return err
	}

	if lhs.GetValueType() == FloatValueType {
		value := NewFloatValue(lhs.(*FloatValue).value + rhs.(*FloatValue).value)
		context.PushValue(value)
		return nil
	} else if lhs.GetValueType() == IntegerValueType {
		lhsInt := lhs.(*IntegerValue)
		rhsInt := rhs.(*IntegerValue)
		values := make([]int64, len(lhsInt.componentValues))
		for vx := 0; vx < len(lhsInt.componentValues); vx++ {
			values[vx] = lhsInt.componentValues[vx] - rhsInt.componentValues[vx]
		}

		offsets := CollapseOffsetList(append(lhsInt.offsets, rhsInt.offsets...))
		lcCount := 0
		for _, offset := range offsets {
			if offset.GetOffsetType() == LocationCounterOffsetType {
				lcCount++
			}
		}

		if lcCount > 1 {
			return fmt.Errorf("result of operation produces incompatible LC offset references")
		}

		newValue, _ := NewIntegerValue(values, lhsInt.form, offsets, 0)
		context.PushValue(newValue)
		return nil
	} else {
		return fmt.Errorf("internal error")
	}
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
	lhs, rhs, err := popLogicalOperands(context)
	if err != nil {
		return err
	}

	values := make([]int64, len(lhs.componentValues))
	for vx := 0; vx < len(lhs.componentValues); vx++ {
		values[vx] = lhs.componentValues[vx] ^ rhs.componentValues[vx]
	}

	value, err := NewIntegerValue(values, lhs.form, mergeOffsets(lhs.offsets, rhs.offsets), 0)
	if err != nil {
		return err
	}
	context.PushValue(value)
	return nil
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
