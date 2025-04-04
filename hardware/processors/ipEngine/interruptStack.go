// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"

	"khalehla/common"
)

type InterruptStack struct {
	stack []common.Interrupt
}

func NewInterruptStack() *InterruptStack {
	return &InterruptStack{
		stack: make([]common.Interrupt, 0),
	}
}

// Clear removes all interrupts from the stack
func (is *InterruptStack) Clear() {
	is.stack = make([]common.Interrupt, 0)
}

func (is *InterruptStack) Dump() {
	for ix := 0; ix < len(is.stack); ix++ {
		fmt.Printf("    %s\n", common.GetInterruptString(is.stack[ix]))
	}
}

// IsClear returns true if there are no interrupts on the stack
func (is *InterruptStack) IsClear() bool {
	return len(is.stack) == 0
}

// Pop pops an interrupt from the top of the stack (or at some proper point in the stack)
// while observing the following rules:
// All other things being equal, the highest priority interrupt is popped first.
// Since the interrupts are stored in the stack in decreasing order of priority, we are guaranteed that
// higher priority interrupts precede lower priority interrupts.
// We will not pop any instruction which is interrupt-able between instructions if we are not between instructions.
// We will not pop any instruction which is interrupt-able mid-execution if we are still resolving indirect addressing.
// If the deferred flag is set, no interrupt which is deferrable, will be popped from the stack.
func (is *InterruptStack) Pop(midExecution bool, resolvingAddress bool, deferred bool) (interrupt common.Interrupt) {
	interrupt = nil

	isLen := len(is.stack)
	for ix := 0; ix < isLen; ix++ {
		i := is.stack[ix]
		if deferred && i.IsDeferrable() {
			continue
		}

		if i.GetInterruptPoint() == common.InterruptBetweenInstruction && (midExecution || resolvingAddress) {
			continue
		}

		if i.GetInterruptPoint() == common.InterruptMidExecution && resolvingAddress {
			continue
		}

		interrupt = i
		is.removeAt(ix)
		break
	}

	return
}

func (is *InterruptStack) PopAll() (result []common.Interrupt) {
	result = is.stack
	is.Clear()
	return
}

// Post posts a new interrupt, provided that no higher-priority interrupt is already pending.
// Interrupts are posted top-down, in order of priority.
// Synchronous interrupts of a lower priority than the new interrupt are discarded.
func (is *InterruptStack) Post(i common.Interrupt) {
	isLen := len(is.stack)
	found := false
	for ix := 0; ix < isLen; {
		if !found {
			if i.GetClass() < is.stack[ix].GetClass() {
				is.insertAt(i, ix)
				found = true
			}
			ix++
		} else {
			if i.GetSynchrony() == common.InterruptSynchronous {
				is.removeAt(ix)
			} else {
				ix++
			}
		}
	}

	if !found {
		is.insertAt(i, isLen)
	}
}

func (is *InterruptStack) insertAt(i common.Interrupt, index int) {
	if len(is.stack) == 0 {
		is.stack = []common.Interrupt{i}
	} else if index == 0 {
		left := []common.Interrupt{i}
		is.stack = append(left, is.stack...)
	} else if index >= len(is.stack) {
		is.stack = append(is.stack, i)
	} else {
		left := append(is.stack[:index], i)
		right := is.stack[index:]
		is.stack = append(left, right...)
	}
}

func (is *InterruptStack) removeAt(index int) {
	if index < len(is.stack) {
		if len(is.stack) == 1 {
			is.stack = make([]common.Interrupt, 0)
		} else {
			is.stack = append(is.stack[:index], is.stack[index:]...)
		}
	}
}
