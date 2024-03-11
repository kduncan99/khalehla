// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

type IoFunction uint

const (
	_ IoFunction = iota
	IofMount
	IofMoveBackward
	IofMoveForward
	IofPrep
	IofRead
	IofReadBackward
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
	IofMoveBackward:    "MoveBack",
	IofMoveForward:     "MoveFwd",
	IofPrep:            "Prep",
	IofRead:            "Read",
	IofReadBackward:    "ReadBack",
	IofReadLabel:       "ReadLabel",
	IofReset:           "Reset",
	IofRewind:          "Rewind",
	IofRewindAndUnload: "RewindUnload",
	IofUnmount:         "Unmount",
	IofWrite:           "Write",
	IofWriteLabel:      "WriteLabel",
	IofWriteTapeMark:   "WriteMark",
}
