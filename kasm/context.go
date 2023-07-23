// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type Context struct {
	code                   map[int][]CodeWord
	currentLineIndex       int
	currentLineNumber      int
	currentLocationCounter int
	currentLiteralPool     int
	diagnostics            *Diagnostics
	dictionary             *Dictionary
}

func NewContext(dictionary *Dictionary) *Context {
	ctx := &Context{
		code:                   make(map[int][]CodeWord),
		currentLineIndex:       0,
		currentLineNumber:      1,
		currentLocationCounter: 0,
		currentLiteralPool:     1,
		dictionary:             dictionary,
		diagnostics:            &Diagnostics{},
	}

	return ctx
}

func (c *Context) AppendErr(err error) {
	c.diagnostics.AppendError(c.currentLineNumber, err.Error())
}

func (c *Context) AppendError(text string) {
	c.diagnostics.AppendError(c.currentLineNumber, text)
}

func (c *Context) AppendInfo(text string) {
	c.diagnostics.AppendInfo(c.currentLineNumber, text)
}

func (c *Context) AppendWarning(text string) {
	c.diagnostics.AppendWarning(c.currentLineNumber, text)
}
