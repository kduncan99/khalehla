// Khalehla Project
// testing assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

import "regexp"

// isValidLabel checks the given string to see if it is formatted as a valid label reference or specification.
// A valid label reference consists of 1 alphabetic character or dollar sign,
// followed by zero or more alphabetic characters, dollar signs, underscores, or digits.
func isValidLabelReference(s string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z\\$][a-zA-Z\\d\\$_]{0,11}$", s)
	return match
}

// isValidLabelSpec checks the given string to see if it is formatted as a valid label specification.
// A valid label specification consists of 1 alphabetic character or dollar sign,
// followed by zero or more alphabetic characters, dollar signs, underscores, or digits,
// followed by zero or more asterisks.
func isValidLabelSpecification(s string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z\\$][a-zA-Z\\d\\$_]{0,11}\\**$", s)
	return match
}

// isValidLocationCounter checks the given string to see if it is formatted as an LCN.
// format is $(n) where n is a 1-to-3 digit decimal or octal literal.
func isValidLocationCounter(s string) bool {
	match, _ := regexp.MatchString("^\\$\\(\\d{1,2}\\)$", s)
	return match
}
