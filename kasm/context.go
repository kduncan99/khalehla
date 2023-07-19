// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type Context struct {
	currentLineIndex       int
	currentLocationCounter int
	currentLiteralPool     int
	diagnostics            diagnostics
	dictionary             Dictionary
}
