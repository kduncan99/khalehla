// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
	"khalehla/tasm"
	"testing"
)

// f, j, and a fields for the various instructions
var fDL = "071"
var fDLM = "071"
var fDLN = "071"
var fJ = "074"
var fHLTJ = "074"
var fIAR = "073"
var fLA = "010"
var fLD = "073"
var fLMA = "012"
var fLNA = "011"
var fLNMA = "013"
var fLR = "023"
var fLX = "027"
var fLXI = "046"
var fLXM = "026"
var fSA = "001"
var fSD = "073"

var jDL = "013"
var jDLM = "015"
var jDLN = "014"
var jHLTJ = "015"
var jIAR = "017"
var jJBasic = "04"
var jJExtended = "015"

var aHLTJ = "05"
var aIAR = "006"
var aJBasic = "0"
var aJExtended = "04"

// partial word designators for j-field specification
var jH1 = "002"
var jH2 = "001"
var jQ1 = "007"
var jQ2 = "004"
var jQ3 = "006"
var jQ4 = "005"
var jS1 = "015"
var jS2 = "014"
var jS3 = "013"
var jS4 = "012"
var jS5 = "011"
var jS6 = "010"
var jT1 = "007"
var jT2 = "006"
var jT3 = "005"
var jU = "016"
var jW = "000"
var jXH2 = "003"
var jXH1 = "004"
var jXU = "017"

// various register values for the a-field and b-field
var rA0 = "0"
var rA1 = "1"
var rA2 = "2"
var rA3 = "3"
var rA4 = "4"
var rA5 = "5"
var rA6 = "6"
var rA7 = "7"
var rA8 = "8"
var rA9 = "9"
var rA10 = "10"
var rA11 = "11"
var rA12 = "12"
var rA13 = "13"
var rA14 = "14"
var rA15 = "15"

var rB0 = "0"
var rB1 = "1"
var rB2 = "2"
var rB3 = "3"
var rB4 = "4"
var rB5 = "5"
var rB6 = "6"
var rB8 = "8"
var rB9 = "9"
var rB10 = "10"
var rB12 = "12"
var rB13 = "13"
var rB14 = "14"
var rB15 = "15"

var rR0 = "0"
var rR1 = "1"
var rR2 = "2"
var rR3 = "3"
var rR4 = "4"
var rR5 = "5"
var rR6 = "6"
var rR7 = "7"
var rR8 = "8"
var rR9 = "9"
var rR10 = "10"
var rR11 = "11"
var rR12 = "12"
var rR13 = "13"
var rR14 = "14"
var rR15 = "15"

var rX0 = "0"
var rX1 = "1"
var rX2 = "2"
var rX3 = "3"
var rX4 = "4"
var rX5 = "5"
var rX6 = "6"
var rX7 = "7"
var rX8 = "8"
var rX9 = "9"
var rX10 = "10"
var rX11 = "11"
var rX12 = "12"
var rX13 = "13"
var rX14 = "14"
var rX15 = "15"

var zero = "0"

var grsx0 = "000"
var grsx1 = "001"
var grsx2 = "002"
var grsx3 = "003"
var grsx4 = "004"
var grsx5 = "005"
var grsx6 = "006"
var grsx7 = "007"
var grsx8 = "010"
var grsx9 = "011"
var grsx10 = "012"
var grsx11 = "013"
var grsx12 = "014"
var grsx13 = "015"
var grsx14 = "016"
var grsx15 = "017"

// GRS locations for registers
var grsa0 = "014"
var grsa1 = "015"
var grsa2 = "016"
var grsa3 = "017"
var grsa4 = "020"
var grsa5 = "021"
var grsa6 = "022"
var grsa7 = "023"
var grsa8 = "024"
var grsa9 = "025"
var grsa10 = "026"
var grsa11 = "027"
var grsa12 = "030"
var grsa13 = "031"
var grsa14 = "032"
var grsa15 = "033"
var grsr0 = "0100"
var grsr1 = "0101"
var grsr2 = "0102"
var grsr3 = "0103"
var grsr4 = "0104"
var grsr5 = "0105"
var grsr6 = "0106"
var grsr7 = "0107"
var grsr8 = "0110"
var grsr9 = "0111"
var grsr10 = "0112"
var grsr11 = "0113"
var grsr12 = "0114"
var grsr13 = "0115"
var grsr14 = "0116"
var grsr15 = "0117"

// iarSourceItem creates an instruction to perform an IAR - this is the preferred way to end a unit test
func iarSourceItem(label string, uField string) *tasm.SourceItem {
	return tasm.NewSourceItem(label, "fjaxu", []string{fIAR, jIAR, aIAR, zero, uField})
}

func checkProgramAddress(t *testing.T, engine *InstructionEngine, expectedAddress uint64) {
	actual := engine.GetProgramAddressRegister().GetProgramCounter()
	if actual != expectedAddress {
		t.Fatalf("Error:Expected PAR.PC is %06o but we expected it to be %06o", actual, expectedAddress)
	}
}

func checkMemory(t *testing.T, engine *InstructionEngine, addr *pkg.AbsoluteAddress, offset uint64, expected uint64) {
	seg, ok := engine.mainStorage.GetSegment(addr.GetSegment())
	if !ok {
		t.Fatalf("Error:segment does not exist for address %s", addr.GetString())
	}

	if addr.GetOffset() >= uint64(len(seg)) {
		t.Fatalf("Error:offset is out of range for address %s - segment size is %012o",
			addr.GetString(), len(seg))
	}

	result := seg[addr.GetOffset()+offset]
	if result.GetW() != expected {
		t.Fatalf("Storage at %s+0%o is %012o, expected %012o", addr.GetString(), offset, result, expected)
	}
}

func checkRegister(t *testing.T, engine *InstructionEngine, register uint64, expected uint64, name string) {
	result := engine.generalRegisterSet.GetRegister(register).GetW()
	if result != expected {
		t.Fatalf("Register %s is %012o, expected %012o", name, result, expected)
	}
}

func checkStopped(t *testing.T, engine *InstructionEngine) {
	if engine.HasPendingInterrupt() {
		t.Fatalf("Engine has unexpected pending interrupt:%s", pkg.GetInterruptString(engine.PopInterrupt()))
	}

	if !engine.IsStopped() {
		t.Fatalf("Expected engine to be stopped; it is not")
	}
}
