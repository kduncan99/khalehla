// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

type Symbol struct {
	symbol     string
	bankNumber int
	offset     uint64
}

type UnresolvedValue struct {
	bankNumber int
	offset     uint64
}

type TestBank struct {
	accessLock       pkg.AccessLock
	bankNumber       int
	code             []pkg.Word36
	symbolTable      []*Symbol
	unresolvedValues []*UnresolvedValue
}

// TestExecutable is a small collection of test banks used for unit testing the IP engine.
// It contains enough of an environment to allow the engine to execute code.
type TestExecutable struct {
	initiallyBased   map[uint]int // maps base register to bank number
	code             map[uint64]*TestBank
	symbolTable      []*Symbol
	unresolvedValues []*UnresolvedValue
}
