// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

type IoFunction uint

const (
	_ IoFunction = iota
	IofMount
	IofPrep
	IofRead
	IofReadLabel
	IofReset
	IofRewind
	IofRewindAndUnload
	IofUnmount
	IofWrite
	IofWriteLabel
	IofWriteTapeMark
)

var IoFunctionTable = map[IoFunction]string{
	IofMount:           "Mount",
	IofPrep:            "Prep",
	IofRead:            "Read",
	IofReadLabel:       "ReadLabel",
	IofReset:           "Reset",
	IofRewind:          "Rewind",
	IofRewindAndUnload: "RewindUnload",
	IofUnmount:         "Unmount",
	IofWrite:           "Write",
	IofWriteLabel:      "WriteLabel",
	IofWriteTapeMark:   "WriteMark",
}
