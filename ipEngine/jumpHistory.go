// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

const JumpHistoryStackSize = 512
const JumpHistoryThreshold = 480

var jhInterrupt = pkg.NewJumpHistoryFullInterrupt()

type JumpHistory struct {
	stack            []pkg.VirtualAddress
	stackIndex       int // index of the next stack entry to be written
	interruptPending bool
	overflow         bool
}

func (jh *JumpHistory) Clear() {
	jh.stackIndex = 0
	jh.overflow = false
	jh.stack = make([]pkg.VirtualAddress, 0)
}

func (jh *JumpHistory) GetEntries() (result []pkg.VirtualAddress) {
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

func (jh *JumpHistory) StoreEntry(address pkg.VirtualAddress) (interrupt pkg.Interrupt) {
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
