// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"strings"
)

type Directive interface {
	GetToken() string
	Interpret(context *Context, labels []string, operation []string, operands []string)
}

var extraneousOperationSubfields = "Ignoring extraneous operation subfields"
var extraneousOperandSubfields = "Ignoring extraneous operand subfields"

type AsciiDirective struct{}
type BaseDirective struct{}
type BasicDirective struct{}
type DeleteDirective struct{}
type EqufDirective struct{}
type ElseDirective struct{}
type ElsfDirective struct{}
type EndDirective struct{}
type EndfDirective struct{}
type ExtendDirective struct{}
type FDataDirective struct{}
type FormDirective struct{}
type GenDirective struct{}
type GoDirective struct{}
type IfDirective struct{}
type IncludeDirective struct{}
type InfoDirective struct{}
type LitDirective struct{}
type NameDirective struct{}
type ProcDirective struct{}
type ResDirective struct{}
type UseDirective struct{}

var directives = []Directive{
	&EquDirective{},
	&EqufDirective{},
	&EndDirective{},
	&GenDirective{},
	&LitDirective{},
}

// InterpretDirective processes the location counter subfield (if found) from the label field,
// then passes all the other label subfields, the operation field, and the operand field to
// successive literal Interpret methods until one of them succeeds.
// Returns true if we process a literal, either successfully or otherwise.
func InterpretDirective(context *Context, labels []string, operation []string, operands []string) bool {
	// If there is a label field and the first subfield contains a location counter, process it and
	// strip the subfield from the label field.
	if labels != nil && len(labels) > 0 {
		lcs, err := NewLocationCounterSpecification(labels[0])
		if err != nil {
			context.AppendErr(err)
		} else if lcs != nil {
			labels = labels[1:]
			val, err := lcs.Evaluate(context)
			if err != nil {
				context.AppendErr(err)
			} else {
				context.currentLocationCounter = val
			}
		}
	}

	if len(operation) > 0 {
		tag := strings.ToUpper(operation[0])
		for _, dir := range directives {
			if tag == strings.ToUpper(dir.GetToken()) {
				dir.Interpret(context, labels, operation, operands)
				return true
			}
		}
	}

	return false
}

func getValuesFromExpressions(context *Context, expressions []*Expression) []Value {
	values := make([]Value, len(expressions))
	for ex, expr := range expressions {
		ec := NewExpressionContext(context)
		err := expr.Evaluate(ec)
		if err != nil {
			context.AppendErr(err)
			values[ex] = NewSimpleIntegerValue(0)
		} else {
			values[ex], err = ec.PopValue()
			if err != nil {
				context.AppendErr(err)
				values[ex] = NewSimpleIntegerValue(0)
			}
		}
	}

	return values
}

// processLabels processes the given labels, setting each of them to the current location counter
func processLabels(context *Context, labels []string) {
	offset := &LocationCounterOffset{
		locationCounter: context.currentLocationCounter,
		startBit:        0,
		bitLength:       36,
		isNegative:      false,
	}
	comp := &ValueComponent{
		value:   uint64(len(context.code[context.currentLocationCounter])),
		offsets: []Offset{offset},
	}
	lcValue := &IntegerValue{
		components: []*ValueComponent{comp},
	}

	for _, label := range labels {
		p := NewParser(label)
		ref, err := p.ParseReference(false, true)
		if err != nil {
			context.AppendErr(err)
		} else if ref != nil {
			values := getValuesFromExpressions(context, ref.arguments)
			err = context.dictionary.Establish(ref.symbol, values, ref.levelCount, lcValue)
			if err != nil {
				context.AppendErr(err)
			}
		} else {
			//	Not a label reference, complain about it
			context.AppendWarning("Non-label reference found in label field ignored")
		}
	}
}

//	$EQU ---------------------------------------------------------------------------------------------------------------

type EquDirective struct {
}

func (d *EquDirective) GetToken() string {
	return "$EQU"
}

func (d *EquDirective) Interpret(context *Context, labels []string, operation []string, operands []string) {
	//	TODO
}

//	$EQUF --------------------------------------------------------------------------------------------------------------

func (d *EqufDirective) GetToken() string {
	return "$EQUF"
}

func (d *EqufDirective) Interpret(context *Context, labels []string, operation []string, operands []string) {
	//	TODO
}

//	$END ---------------------------------------------------------------------------------------------------------------

func (d *EndDirective) GetToken() string {
	return "$END"
}

func (d *EndDirective) Interpret(context *Context, labels []string, operation []string, operands []string) {
	//	TODO
}

//	$GEN ---------------------------------------------------------------------------------------------------------------

func (d *GenDirective) GetToken() string {
	return "$GEN"
}

func (d *GenDirective) Interpret(context *Context, labels []string, operation []string, operands []string) {
	//	TODO
}

//	$LIT ---------------------------------------------------------------------------------------------------------------

func (d *LitDirective) GetToken() string {
	return "$LIT"
}

func (d *LitDirective) Interpret(context *Context, labels []string, operation []string, operands []string) {
	if len(operation) > 0 {
		context.AppendWarning(extraneousOperationSubfields)
	}

	if operands != nil && len(operands) > 0 {
		context.AppendWarning(extraneousOperandSubfields)
	}

	if len(labels) > 1 {
		subLabels := labels[:len(labels)-2]
		labels = labels[len(labels)-1:]
		processLabels(context, subLabels)
	}

	if len(labels) == 1 {
		p := NewParser(labels[0])
		ref, err := p.ParseReference(false, true)
		if err != nil {
			context.AppendErr(err)
		} else {
			values := getValuesFromExpressions(context, ref.arguments)
			fn := &LitFunction{locationCounter: context.currentLocationCounter}
			err = context.dictionary.Establish(ref.symbol, values, ref.levelCount, fn)
			if err != nil {
				context.AppendErr(err)
			}
		}
	}

	context.currentLiteralPool = context.currentLocationCounter
}
