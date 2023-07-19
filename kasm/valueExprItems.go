// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

type valueItemType int
type literalType int

const (
	LiteralType valueItemType = iota
	ReferenceType
)

const (
	FloatLiteralType literalType = iota
	IntegerLiteralType
	StringLiteralType
)

type valueExpressionItem interface {
	Evaluate(expression *Expression) Value
	GetExpressionItemType() expressionItemType
	GetValueItemType() valueItemType
}

type reference struct {
	symbol     string
	parameters []*Expression
}

type literal interface {
	Evaluate(expression *Expression) Value
	GetExpressionItemType() expressionItemType
	GetValueItemType() valueItemType
	GetLiteralType() literalType
	GetText() string
	GetValue() Value
	IsDouble() bool
	IsSingle() bool
	IsLeft() bool
	IsRight() bool
}

type floatLiteral struct {
	text     string
	isDouble bool
	isSingle bool
	isLeft   bool
	isRight  bool
}

type integerLiteral struct {
	text     string
	isDouble bool
	isSingle bool
	isLeft   bool
	isRight  bool
}

type stringLiteral struct {
	text     string
	isDouble bool
	isSingle bool
	isLeft   bool
	isRight  bool
}

// reference -----------------------------------------------------------------------------------------------------------

func (r *reference) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (r *reference) GetExpressionItemType() expressionItemType {
	return ValueItemType
}

func (r *reference) GetValueItemType() valueItemType {
	return ReferenceType
}

// float literal -------------------------------------------------------------------------------------------------------

func (l *floatLiteral) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (l *floatLiteral) GetExpressionItemType() expressionItemType {
	return ValueItemType
}

func (l *floatLiteral) GetValueItemType() valueItemType {
	return LiteralType
}

func (l *floatLiteral) GetLiteralType() literalType {
	return FloatLiteralType
}

func (l *floatLiteral) GetText() string {
	return l.text
}

func (l *floatLiteral) GetValue() Value {
	return nil //	TODO
}

func (l *floatLiteral) IsDouble() bool {
	return l.isDouble
}

func (l *floatLiteral) IsSingle() bool {
	return l.isSingle
}

func (l *floatLiteral) IsLeft() bool {
	return l.isLeft
}

func (l *floatLiteral) IsRight() bool {
	return l.isRight
}

// integer literal -----------------------------------------------------------------------------------------------------

func (l *integerLiteral) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (l *integerLiteral) GetExpressionItemType() expressionItemType {
	return ValueItemType
}

func (l *integerLiteral) GetValueItemType() valueItemType {
	return LiteralType
}

func (l *integerLiteral) GetLiteralType() literalType {
	return IntegerLiteralType
}

func (l *integerLiteral) GetText() string {
	return l.text
}

func (l *integerLiteral) GetValue() Value {
	return nil //	TODO
}

func (l *integerLiteral) IsDouble() bool {
	return l.isDouble
}

func (l *integerLiteral) IsSingle() bool {
	return l.isSingle
}

func (l *integerLiteral) IsLeft() bool {
	return l.isLeft
}

func (l *integerLiteral) IsRight() bool {
	return l.isRight
}

// string literal ------------------------------------------------------------------------------------------------------

func (l *stringLiteral) Evaluate(expression *Expression) Value {
	return nil // TODO
}

func (l *stringLiteral) GetExpressionItemType() expressionItemType {
	return ValueItemType
}

func (l *stringLiteral) GetValueItemType() valueItemType {
	return LiteralType
}

func (l *stringLiteral) GetLiteralType() literalType {
	return StringLiteralType
}

func (l *stringLiteral) GetText() string {
	return l.text
}

func (l *stringLiteral) GetValue() Value {
	return nil //	TODO
}

func (l *stringLiteral) IsDouble() bool {
	return l.isDouble
}

func (l *stringLiteral) IsSingle() bool {
	return l.isSingle
}

func (l *stringLiteral) IsLeft() bool {
	return l.isLeft
}

func (l *stringLiteral) IsRight() bool {
	return l.isRight
}
