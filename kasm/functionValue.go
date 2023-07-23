// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import "fmt"

// FunctionValue defines a function
type FunctionValue interface {
	Evaluate(*ExpressionContext) error
	GetValueType() ValueType
}

type UserFunction struct {
	code          []string
	externalNames map[string]int // maps a $NAME to the textIndex of the line which contains it
}

type CurrentLocationCounterNumber struct{} //	$LCN
type LitFunction struct {
	locationCounter int
}

//	TODO many more built-in functions
//	TODO need to insert these into the top-level dictionary

var tooManyParameters = fmt.Errorf("too many parameters in function call")

// User function -------------------------------------------------------------------------------------------------------

func (f *UserFunction) Evaluate(context *ExpressionContext) error {
	return nil //	TODO
}

func (f *UserFunction) GetValueType() ValueType {
	return FunctionValueType
}

// Literal function ----------------------------------------------------------------------------------------------------

func (f *LitFunction) Evaluate(ec *ExpressionContext) error {
	// TODO the parameters are portions of a word which is created in the literal pool defined by
	//		f.locationCounter
	// and the resulting value is the LC offset address of that literal.
	return nil
}

func (f *LitFunction) GetValueType() ValueType {
	return FunctionValueType
}

// $LCN ----------------------------------------------------------------------------------------------------------------

func (f *CurrentLocationCounterNumber) Evaluate(ec *ExpressionContext) error {
	params, err := ec.PopVariableParameterList()
	if err != nil {
		return err
	} else if len(params) > 0 {
		return tooManyParameters
	}

	ec.PushValue(NewSimpleIntegerValue(uint64(ec.context.currentLocationCounter)))
	return nil
}

func (f *CurrentLocationCounterNumber) GetValueType() ValueType {
	return FunctionValueType
}
