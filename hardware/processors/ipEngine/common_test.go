// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"testing"

	"khalehla/common"
	"khalehla/tasm"
)

// partial word designators for j-field specification
const (
	jH1  = 002
	jH2  = 001
	jQ1  = 007
	jQ2  = 004
	jQ3  = 006
	jQ4  = 005
	jS1  = 015
	jS2  = 014
	jS3  = 013
	jS4  = 012
	jS5  = 011
	jS6  = 010
	jT1  = 007
	jT2  = 006
	jT3  = 005
	jU   = 016
	jW   = 000
	jXH2 = 003
	jXH1 = 004
	jXU  = 017
)

const zero = "0"

// A-field/X-field values for registers
const (
	regX0 = iota
	regX1
	regX2
	regX3
	regX4
	regX5
	regX6
	regX7
	regX8
	regX9
	regX10
	regX11
	regX12
	regX13
	regX14
	regX15
)

const (
	regA0 = iota
	regA1
	regA2
	regA3
	regA4
	regA5
	regA6
	regA7
	regA8
	regA9
	regA10
	regA11
	regA12
	regA13
	regA14
	regA15
)

const (
	regR0 = iota
	regR1
	regR2
	regR3
	regR4
	regR5
	regR6
	regR7
	regR8
	regR9
	regR10
	regR11
	regR12
	regR13
	regR14
	regR15
)

// ---------------------------------------------------------------------------------------------------------------------

func grsRef(grsRegIndex uint64) string {
	return fmt.Sprintf("0%o", grsRegIndex)
}

func sourceItem(label string, operator string, operands []int) *tasm.SourceItem {
	strOps := make([]string, len(operands))
	for ox := 0; ox < len(operands); ox++ {
		strOps[ox] = fmt.Sprintf("0%o", operands[ox])
	}

	return tasm.NewSourceItem(label, operator, strOps)
}

func labelSourceItem(label string) *tasm.SourceItem {
	return tasm.NewSourceItem(label, "", []string{})
}

func labelDataSourceItem(label string, values []uint64) *tasm.SourceItem {
	var operator string
	if len(values) == 1 {
		operator = "w"
	} else if len(values) == 2 {
		operator = "hw"
	} else if len(values) == 3 {
		operator = "tw"
	} else if len(values) == 4 {
		operator = "qw"
	} else if len(values) == 6 {
		operator = "sw"
	} else {
		operator = "?"
	}

	strValues := make([]string, len(values))
	for vx := 0; vx < len(values); vx++ {
		strValues[vx] = fmt.Sprintf("0%o", values[vx])
	}

	return tasm.NewSourceItem(label, operator, strValues)
}

func dataSourceItem(values []uint64) *tasm.SourceItem {
	return labelDataSourceItem("", values)
}

func fjaxuSourceItem(f uint64, j uint64, a uint64, x uint64, u uint64) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("0%o", f),
		fmt.Sprintf("0%o", j),
		fmt.Sprintf("0%o", a),
		fmt.Sprintf("0%o", x),
		fmt.Sprintf("0%o", u),
	}
	return tasm.NewSourceItem("", "fjaxu", ops)
}

// This is for 18-bit U fields, for which x is *usually* zero.
// Nonetheless, we still require specification of x-field just inc ase.
func fjaxRefSourceItem(f uint64, j uint64, a uint64, x uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("0%o", f),
		fmt.Sprintf("0%o", j),
		fmt.Sprintf("0%o", a),
		fmt.Sprintf("0%o", x),
		ref,
	}
	return tasm.NewSourceItem("", "fjaxu", ops)
}

func fjaxhibRefSourceItem(f uint64, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("0%o", f),
		fmt.Sprintf("0%o", j),
		fmt.Sprintf("0%o", a),
		fmt.Sprintf("0%o", x),
		fmt.Sprintf("0%o", h),
		fmt.Sprintf("0%o", i),
		fmt.Sprintf("0%o", b),
		ref,
	}
	return tasm.NewSourceItem("", "fjaxhibd", ops)
}

func fjaxhiRefSourceItem(f uint64, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("0%o", f),
		fmt.Sprintf("0%o", j),
		fmt.Sprintf("0%o", a),
		fmt.Sprintf("0%o", x),
		fmt.Sprintf("0%o", h),
		fmt.Sprintf("0%o", i),
		ref,
	}
	return tasm.NewSourceItem("", "fjaxhiu", ops)
}

func segSourceItem(segIndex int) *tasm.SourceItem {
	return tasm.NewSourceItem("", ".SEG", []string{fmt.Sprintf("%d", segIndex)})
}

// ---------------------------------------------------------------------------------------------------------------------

func checkInterrupt(t *testing.T, engine *InstructionEngine, interruptClass common.InterruptClass) {
	for _, i := range engine.pendingInterrupts.stack {
		if i.GetClass() == interruptClass {
			return
		}
	}

	engine.pendingInterrupts.Dump()
	t.Errorf("Error:Expected interrupt class %d to be posted", interruptClass)
}

func checkInterruptAndSSF(
	t *testing.T,
	engine *InstructionEngine,
	interruptClass common.InterruptClass,
	shortStatusField common.InterruptShortStatus) {

	for _, i := range engine.pendingInterrupts.stack {
		if i.GetClass() == interruptClass {
			if i.GetShortStatusField() != shortStatusField {
				t.Errorf("Error:Found interrupt class %d but SSF was %d and we expected %d",
					interruptClass, i.GetShortStatusField(), shortStatusField)
			}
			return
		}
	}

	engine.pendingInterrupts.Dump()
	t.Errorf("Error:Expected interrupt class %d to be posted", interruptClass)
}

func checkProgramAddress(t *testing.T, engine *InstructionEngine, expectedAddress uint64) {
	actual := engine.GetProgramAddressRegister().GetProgramCounter()
	if actual != expectedAddress {
		t.Errorf("Error:Expected PAR.PC is %06o but we expected it to be %06o", actual, expectedAddress)
	}
}

func checkMemory(t *testing.T, engine *InstructionEngine, addr *common.AbsoluteAddress, offset uint64, expected uint64) {
	seg, interrupt := engine.mainStorage.GetSegment(addr.GetSegment())
	if interrupt != nil {
		engine.mainStorage.Dump()
		t.Errorf("Error:%s", common.GetInterruptString(interrupt))
	}

	if addr.GetOffset() >= uint64(len(seg)) {
		engine.mainStorage.Dump()
		t.Errorf("Error:offset is out of range for address %s - segment size is %012o", addr.GetString(), len(seg))
	}

	result := seg[addr.GetOffset()+offset]
	if result.GetW() != expected {
		engine.mainStorage.Dump()
		t.Errorf("Storage at (%s)+0%o is %012o, expected %012o", addr.GetString(), offset, result, expected)
	}
}

func checkRegister(t *testing.T, engine *InstructionEngine, regIndex uint64, expected uint64) {
	result := engine.generalRegisterSet.GetRegister(regIndex).GetW()
	if result != expected {
		engine.generalRegisterSet.Dump()
		t.Errorf("Register %s is %012o, expected %012o", common.RegisterNames[regIndex], result, expected)
	}
}

func checkStoppedReason(t *testing.T, engine *InstructionEngine, reason StopReason, detail uint64) {
	if engine.HasPendingInterrupt() {
		engine.Dump()
		t.Errorf("Engine has unexpected pending interrupts")
	}

	if !engine.IsStopped() {
		engine.Dump()
		t.Errorf("Expected engine to be stopped; it is not")
	}

	actualReason, actualDetail := engine.GetStopReason()
	if actualReason != reason {
		engine.Dump()
		t.Errorf("Engine stopped for reason %d; expected reason %d", actualReason, reason)
	}

	if actualDetail != detail {
		engine.Dump()
		t.Errorf("Engine stopped for detail %d; expected detail %d", actualDetail, detail)
	}
}
