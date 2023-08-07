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
const (
	fDL   = "071"
	fDLM  = "071"
	fDLN  = "071"
	fHKJ  = "074"
	fHLTJ = "074"
	fIAR  = "073"
	fJ    = "074"
	fJC   = "074"
	fJDF  = "074"
	fJFO  = "074"
	fJFU  = "074"
	fJK   = "074"
	fJNC  = "074"
	fJNDF = "074"
	fJNFO = "074"
	fJNFU = "074"
	fJNO  = "074"
	fJO   = "074"
	fLA   = "010"
	fLD   = "073"
	fLMA  = "012"
	fLNA  = "011"
	fLNMA = "013"
	fLR   = "023"
	fLX   = "027"
	fLXI  = "046"
	fLXM  = "026"
	fSA   = "001"
	fSD   = "073"
)

const (
	jDL          = "013"
	jDLM         = "015"
	jDLN         = "014"
	jHKJ         = "005"
	jHLTJ        = "015"
	jIAR         = "017"
	jJBasic      = "004"
	jJExtended   = "015"
	jJCBasic     = "016"
	jJCExtended  = "014"
	jJDF         = "014"
	jJFO         = "014"
	jJFU         = "014"
	jJK          = "004"
	jJNCBasic    = "017"
	jJNCExtended = "014"
	jJNDF        = "015"
	jJNFO        = "015"
	jJNFU        = "015"
	jJNO         = "015"
	jJO          = "014"
)

const (
	aHLTJ        = "005"
	aIAR         = "006"
	aJBasic      = "000"
	aJExtended   = "004"
	aJCBasic     = "000"
	aJCExtended  = "004"
	aJDF         = "003"
	aJFO         = "002"
	aJFU         = "001"
	aJNCBasic    = "000"
	aJNCExtended = "005"
	aJNDF        = "003"
	aJNFO        = "002"
	aJNFU        = "001"
	aJNO         = "000"
	aJO          = "000"
)

// partial word designators for j-field specification
const (
	jH1  = "002"
	jH2  = "001"
	jQ1  = "007"
	jQ2  = "004"
	jQ3  = "006"
	jQ4  = "005"
	jS1  = "015"
	jS2  = "014"
	jS3  = "013"
	jS4  = "012"
	jS5  = "011"
	jS6  = "010"
	jT1  = "007"
	jT2  = "006"
	jT3  = "005"
	jU   = "016"
	jW   = "000"
	jXH2 = "003"
	jXH1 = "004"
	jXU  = "017"
)

const zero = "0"

// various register values for the a-field and b-field
const (
	rA0  = "0"
	rA1  = "1"
	rA2  = "2"
	rA3  = "3"
	rA4  = "4"
	rA5  = "5"
	rA6  = "6"
	rA7  = "7"
	rA8  = "8"
	rA9  = "9"
	rA10 = "10"
	rA11 = "11"
	rA12 = "12"
	rA13 = "13"
	rA14 = "14"
	rA15 = "15"
)

const (
	rB0  = "0"
	rB1  = "1"
	rB2  = "2"
	rB3  = "3"
	rB4  = "4"
	rB5  = "5"
	rB6  = "6"
	rB8  = "8"
	rB9  = "9"
	rB10 = "10"
	rB12 = "12"
	rB13 = "13"
	rB14 = "14"
	rB15 = "15"
)

const (
	rR0  = "0"
	rR1  = "1"
	rR2  = "2"
	rR3  = "3"
	rR4  = "4"
	rR5  = "5"
	rR6  = "6"
	rR7  = "7"
	rR8  = "8"
	rR9  = "9"
	rR10 = "10"
	rR11 = "11"
	rR12 = "12"
	rR13 = "13"
	rR14 = "14"
	rR15 = "15"
)

const (
	rX0  = "0"
	rX1  = "1"
	rX2  = "2"
	rX3  = "3"
	rX4  = "4"
	rX5  = "5"
	rX6  = "6"
	rX7  = "7"
	rX8  = "8"
	rX9  = "9"
	rX10 = "10"
	rX11 = "11"
	rX12 = "12"
	rX13 = "13"
	rX14 = "14"
	rX15 = "15"
)

// GRS locations for registers
const (
	grsX0  = "000"
	grsX1  = "001"
	grsX2  = "002"
	grsX3  = "003"
	grsX4  = "004"
	grsX5  = "005"
	grsX6  = "006"
	grsX7  = "007"
	grsX8  = "010"
	grsX9  = "011"
	grsX10 = "012"
	grsX11 = "013"
	grsX12 = "014"
	grsX13 = "015"
	grsX14 = "016"
	grsX15 = "017"
)

const (
	grsA0  = "014"
	grsA1  = "015"
	grsA2  = "016"
	grsA3  = "017"
	grsA4  = "020"
	grsA5  = "021"
	grsA6  = "022"
	grsA7  = "023"
	grsA8  = "024"
	grsA9  = "025"
	grsA10 = "026"
	grsA11 = "027"
	grsA12 = "030"
	grsA13 = "031"
	grsA14 = "032"
	grsA15 = "033"
)

const (
	grsR0  = "0100"
	grsR1  = "0101"
	grsR2  = "0102"
	grsR3  = "0103"
	grsR4  = "0104"
	grsR5  = "0105"
	grsR6  = "0106"
	grsR7  = "0107"
	grsR8  = "0110"
	grsR9  = "0111"
	grsR10 = "0112"
	grsR11 = "0113"
	grsR12 = "0114"
	grsR13 = "0115"
	grsR14 = "0116"
	grsR15 = "0117"
)

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
