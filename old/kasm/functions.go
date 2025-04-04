// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import "fmt"

// Function defines a function
type Function interface {
	Evaluate(*ExpressionContext) error
	GetValueType() ValueType
}

type UserFunction struct {
	code          []string
	externalNames map[string]int // maps a $NAME to the textIndex of the line which contains it
}

type LitFunction struct {
	locationCounter int
}

type APFunction struct{}
type BAFunction struct{}
type BREGFunction struct{}
type CASFunction struct{}
type CBFunction struct{}
type CDFunction struct{}
type CFSFunction struct{}
type CSFunction struct{}
type DATEFunction struct{}
type FNFunction struct{}
type IBITSFunction struct{}
type L0Function struct{}
type L1Function struct{}
type LCBFunction struct{}
type LCFVFunction struct{}
type LCNFunction struct{}
type LCVFunction struct{}
type LEVFunction struct{}
type LINESFunction struct{}
type NODEFunction struct{}
type NSFunction struct{}
type SLFunction struct{}
type SNFunction struct{}
type SRFunction struct{}
type SSFunction struct{}
type SSSFunction struct{}
type TYPEFunction struct{}

var Functions = map[string]Function{
	"$":    &LCVFunction{},
	"$CAS": &CASFunction{},
	"$CFS": &CASFunction{},
	"$LCN": &LCNFunction{},
	"$LCV": &LCVFunction{},
	"$SL":  &SLFunction{},
	"$SR":  &SRFunction{},
	"$SS":  &SSFunction{},
}

var invalidParameterError = fmt.Errorf("invalid parameter in function parameter list")
var parameterTypeError = fmt.Errorf("wrong parameter type in function parameter list")
var wrongNumberOfParameters = fmt.Errorf("wrong number of parameters supplied to function call")

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

//	TODO  $AP
//	TODO  $BA
//	TODO  $BREG

//	$CAS ---------------------------------------------------------------------------------------------------------------

func (f *CASFunction) Evaluate(ec *ExpressionContext) error {
	params, err := ec.PopVariableParameterList()
	if err != nil {
		return err
	} else if len(params) != 1 {
		return wrongNumberOfParameters
	}

	if params[0].GetValueType() != StringValueType {
		return parameterTypeError
	}

	sv := params[0].(*StringValue)
	nsv := NewStringValue(sv.value, AsciiString, sv.flags)
	ec.PushValue(nsv)
	return nil
}

func (f *CASFunction) GetValueType() ValueType {
	return FunctionValueType
}

//	TODO  $CB
//	TODO  $CD

//	$CFS ---------------------------------------------------------------------------------------------------------------

func (f *CFSFunction) Evaluate(ec *ExpressionContext) error {
	params, err := ec.PopVariableParameterList()
	if err != nil {
		return err
	} else if len(params) != 1 {
		return wrongNumberOfParameters
	}

	if params[0].GetValueType() != StringValueType {
		return parameterTypeError
	}

	sv := params[0].(*StringValue)
	nsv := NewStringValue(sv.value, FieldataString, sv.flags)
	ec.PushValue(nsv)
	return nil
}

func (f *CFSFunction) GetValueType() ValueType {
	return FunctionValueType
}

//	TODO  $CS
//	TODO  $DATE
//	TODO  $FN
//	TODO  $IBITS
//	TODO  $L0
//	TODO  $L1
//	TODO  $LCB
//	TODO  $LCFV

// $LCN ----------------------------------------------------------------------------------------------------------------

func (f *LCNFunction) Evaluate(ec *ExpressionContext) error {
	params, err := ec.PopVariableParameterList()
	if err != nil {
		return err
	}

	var lcn int
	if len(params) == 0 {
		lcn = ec.context.currentLocationCounter
	} else if len(params) == 1 {
		if params[0].GetValueType() != IntegerValueType {
			return parameterTypeError
		}

		iv := params[0].(*IntegerValue)
		if !iv.form.Equals(SimpleForm) || len(iv.offsets) > 0 {
			return invalidParameterError
		}

		lcn = int(iv.componentValues[0])
	} else {
		return wrongNumberOfParameters
	}

	iVal := int64(len(ec.context.code[lcn]))
	ec.PushValue(NewSimpleIntegerValue(iVal))
	return nil
}

func (f *LCNFunction) GetValueType() ValueType {
	return FunctionValueType
}

//	$LCV ---------------------------------------------------------------------------------------------------------------

func (f *LCVFunction) Evaluate(ec *ExpressionContext) error {
	params, err := ec.PopVariableParameterList()
	if err != nil {
		return err
	} else if len(params) != 1 {
		return wrongNumberOfParameters
	}

	if params[0].GetValueType() != StringValueType {
		return parameterTypeError
	}

	sv := params[0].(*StringValue)
	niv := NewSimpleIntegerValue(int64(len(sv.value)))
	ec.PushValue(niv)
	return nil
}

func (f *LCVFunction) GetValueType() ValueType {
	return FunctionValueType
}

//	TODO  $LEV
//	TODO  $LINES
//	TODO  $NODE
//	TODO  $NS

//	$SL ----------------------------------------------------------------------------------------------------------------

func (f *SLFunction) Evaluate(ec *ExpressionContext) error {
	params, err := ec.PopVariableParameterList()
	if err != nil {
		return err
	} else if len(params) != 1 {
		return wrongNumberOfParameters
	}

	if params[0].GetValueType() != StringValueType {
		return parameterTypeError
	}

	sv := params[0].(*StringValue)
	niv := NewSimpleIntegerValue(int64(len(sv.value)))
	ec.PushValue(niv)
	return nil
}

func (f *SLFunction) GetValueType() ValueType {
	return FunctionValueType
}

//	TODO  $SN

//	$SR ---------------------------------------------------------------------------------------------------------------0

func (f *SRFunction) Evaluate(ec *ExpressionContext) error {
	params, err := ec.PopVariableParameterList()
	if err != nil {
		return err
	} else if len(params) != 2 {
		return wrongNumberOfParameters
	}

	if params[0].GetValueType() != StringValueType || params[1].GetValueType() != IntegerValueType {
		return parameterTypeError
	}

	sv := params[0].(*StringValue)
	iv := params[1].(*IntegerValue)
	if !iv.form.Equals(SimpleForm) || len(iv.offsets) > 0 {
		return invalidParameterError
	}

	iVal := int(iv.componentValues[0])
	var str string
	for x := 0; x < iVal; x++ {
		str += sv.value
	}

	nsv := NewStringValue(str, AsciiString, 0)
	ec.PushValue(nsv)
	return nil
}

func (f *SRFunction) GetValueType() ValueType {
	return FunctionValueType
}

//	$SS ----------------------------------------------------------------------------------------------------------------

func (f *SSFunction) Evaluate(ec *ExpressionContext) error {
	params, err := ec.PopVariableParameterList()
	if err != nil {
		return err
	} else if len(params) < 2 || len(params) > 3 {
		return wrongNumberOfParameters
	}

	if params[0].GetValueType() != StringValueType ||
		params[1].GetValueType() != IntegerValueType ||
		(len(params) == 3 && params[2].GetValueType() != IntegerValueType) {
		return parameterTypeError
	}

	sv := params[0].(*StringValue)
	iv1 := params[1].(*IntegerValue)
	if !iv1.form.Equals(SimpleForm) || len(iv1.offsets) > 0 {
		return invalidParameterError
	}
	int1 := int(iv1.componentValues[0])

	int2 := 1
	if len(params) == 3 {
		iv2 := params[2].(*IntegerValue)
		if !iv2.form.Equals(SimpleForm) || len(iv2.offsets) > 0 {
			return invalidParameterError
		}
		int2 = int(iv2.componentValues[0])
	}

	str := sv.value[int1 : int1+int2]
	for len(str) < int2 {
		str += " "
	}

	nsv := NewStringValue(str, AsciiString, 0)
	ec.PushValue(nsv)
	return nil
}

func (f *SSFunction) GetValueType() ValueType {
	return FunctionValueType
}

//	TODO  $SSS
//	TODO  $TYPE
