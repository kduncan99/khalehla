// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
	"regexp"
)

type context struct {
	currentLineIndex       int
	currentLocationCounter int
	currentLiteralPool     int
	diagnostics            diagnostics
}

func interpretLocationCounter(text string) (int, error) {
	exprText := (text)[2 : len(text)-1]
	expr := NewExpression(exprText)
	value, err := expr.evaluate()
	if err != nil {
		return 0, err
	}

	if value.GetValueType() != IntegerValueType {
		return 0, fmt.Errorf("bad Value Type for Location Counter")
	}

	iv := value.(*IntegerValue)
	if len(iv.components) > 1 {
		return 0, fmt.Errorf("cannot use component value for Location Counter")
	}

	if len(iv.components[0].offsets) > 0 {
		return 0, fmt.Errorf("cannot use relocatable value for Location Counter")
	}

	lcn := iv.components[0].value.Int64()
	if (lcn < 0) || (lcn > 63) {
		return 0, fmt.Errorf("value out of range for location counter")
	}

	return int(lcn), nil
}

// isValidLabelReference checks the given string to see if it is formatted as a valid label reference or specification.
// A valid label reference consists of 1 alphabetic character or dollar sign,
// followed by zero or more alphabetic characters, dollar signs, underscores, or digits.
func isValidLabelReference(s string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z\\$][a-zA-Z\\d\\$_]{0,11}$", s)
	return match
}

// isValidLabelSpecification checks the given string to see if it is formatted as a valid label specification.
// A valid label specification consists of 1 alphabetic character or dollar sign,
// followed by zero or more alphabetic characters, dollar signs, underscores, or digits,
// followed by zero or more asterisks.
func isValidLabelSpecification(s string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z\\$][a-zA-Z\\d\\$_]{0,11}\\**$", s)
	return match
}

// isValidLocationCounter checks the given string to see if it is formatted as an LCN.
// format is $(x) where x is an expression which we do not evaluate in any way here.
func isValidLocationCounter(s string) bool {
	match, _ := regexp.MatchString("^\\$\\(.+\\)$", s)
	return match
}
