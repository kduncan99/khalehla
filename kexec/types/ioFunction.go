// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

type IoFunction uint

const (
	_ IoFunction = iota
	IofMount
	IofPrep
	IofReset
	IofRead
	IofReadLabel
	IofUnmount
	IofWrite
	IofWriteLabel
)

var IoFunctionTable = map[IoFunction]string{
	IofMount:      "Mount",
	IofPrep:       "Prep",
	IofRead:       "Read",
	IofReadLabel:  "ReadLabel",
	IofReset:      "Reset",
	IofUnmount:    "Unmount",
	IofWrite:      "Write",
	IofWriteLabel: "WriteLabel",
}
