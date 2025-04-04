// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/common"
	"khalehla/hardware"
)

const (
	JumpHistoryStackSize = 512
	JumpHistoryThreshold = 480
)

var jhInterrupt = common.NewJumpHistoryFullInterrupt()

// JumpHistory tracks the most recent jump-from addresses for a task.
// It is part of the instruction processor context.
type JumpHistory struct {
	stack            []hardware.VirtualAddress
	stackIndex       int // index of the next stack entry to be written
	interruptPending bool
	overflow         bool
}

func (jh *JumpHistory) Clear() {
	jh.stackIndex = 0
	jh.overflow = false
	jh.stack = make([]hardware.VirtualAddress, JumpHistoryStackSize)
}

func (jh *JumpHistory) GetEntries() (result []hardware.VirtualAddress) {
	if jh.overflow {
		left := jh.stack[jh.stackIndex+1:]
		right := jh.stack[:jh.stackIndex]
		result = append(left, right...)
	} else {
		result = jh.stack[:jh.stackIndex]
	}

	jh.Clear()
	return
}

func (jh *JumpHistory) StoreEntry(address hardware.VirtualAddress) (interrupt common.Interrupt) {
	interrupt = nil

	jh.stack[jh.stackIndex] = address
	jh.stackIndex++
	if jh.stackIndex == JumpHistoryStackSize {
		jh.stackIndex = 0
		jh.overflow = true
	}

	if jh.stackIndex >= JumpHistoryThreshold && !jh.interruptPending {
		interrupt = jhInterrupt
		jh.interruptPending = true
	}

	return
}

func NewJumpHistory() *JumpHistory {
	jh := JumpHistory{}
	jh.Clear()
	return &jh
}
