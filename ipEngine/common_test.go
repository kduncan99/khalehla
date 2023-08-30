// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"khalehla/pkg"
	"khalehla/tasm"
	"testing"
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
		fmt.Sprintf("%03o", f),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", u),
	}
	return tasm.NewSourceItem("", "fjaxu", ops)
}

// This is for 18-bit U fields, for which x is *usually* zero.
// Nonetheless, we still require specification of x-field just inc ase.
func fjaxRefSourceItem(f uint64, j uint64, a uint64, x uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", f),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		ref,
	}
	return tasm.NewSourceItem("", "fjaxu", ops)
}

func fjaxhibRefSourceItem(f uint64, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", f),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem("", "fjaxhibd", ops)
}

func fjaxhiRefSourceItem(f uint64, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", f),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem("", "fjaxhiu", ops)
}

func segSourceItem(segIndex int) *tasm.SourceItem {
	return tasm.NewSourceItem("", ".SEG", []string{fmt.Sprintf("%d", segIndex)})
}

// ---------------------------------------------------------------------------------------------------------------------
// Logical functions
// ---------------------------------------------------------------------------------------------------------------------

// OR ------------------------------------------------------------------------------------------------------------------

const fOR = 040

func orSourceItemHIBD(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fOR, j, a, x, h, i, b, d})
}

func orSourceItemHIBDRef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fOR),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func orSourceItemHIU(label string, j uint64, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fOR, j, a, x, h, i, u})
}

func orSourceItemHIURef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fOR),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func orSourceItemU(label string, j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fOR, j, a, x, u})
}

// XOR -----------------------------------------------------------------------------------------------------------------

const fXOR = 041

func xorSourceItemHIBD(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fXOR, j, a, x, h, i, b, d})
}

func xorSourceItemHIBDRef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fXOR),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func xorSourceItemHIU(label string, j uint64, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fXOR, j, a, x, h, i, u})
}

func xorSourceItemHIURef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fXOR),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func xorSourceItemU(label string, j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fXOR, j, a, x, u})
}

// AND -----------------------------------------------------------------------------------------------------------------

const fAND = 042

func andSourceItemHIBD(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fAND, j, a, x, h, i, b, d})
}

func andSourceItemHIBDRef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fAND),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func andSourceItemHIU(label string, j uint64, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fAND, j, a, x, h, i, u})
}

func andSourceItemHIURef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fAND),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func andSourceItemU(label string, j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fAND, j, a, x, u})
}

// MLU -----------------------------------------------------------------------------------------------------------------

const fMLU = 043

func mluSourceItemHIBD(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fMLU, j, a, x, h, i, b, d})
}

func mluSourceItemHIBDRef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fMLU),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func mluSourceItemHIU(label string, j uint64, a uint64, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fMLU, j, a, x, h, i, u})
}

func mluSourceItemHIURef(label string, j uint64, a uint64, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fMLU),
		fmt.Sprintf("%03o", j),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func mluSourceItemU(label string, j uint64, a uint64, x uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fMLU, j, a, x, u})
}

// ---------------------------------------------------------------------------------------------------------------------
// Storage-to-storage functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO BT
//	TODO BIM
//	TODO BIC
//	TODO BIMT
//	TODO BICL
//	TODO BIML
//	TODO BN
//	TODO BBN

// ---------------------------------------------------------------------------------------------------------------------
// String functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LS
//	TODO LSA
//	TODO SS
//	TODO TES
//	TODO TNES

// ---------------------------------------------------------------------------------------------------------------------
// Address Space Management functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LBU
//	TODO LBE
//	TODO LBUD
//	TODO LBED
//	TODO SBUD
//	TODO SBED
//	TODO SBU
//	TODO LBN
//	TODO TRA
//	TODO TVA
//	TODO DABT
//	TODO TRARS

// ---------------------------------------------------------------------------------------------------------------------
// Procedure Control functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO GOTO
//	TODO CALL
//	TODO LOCL
//	TODO RTN
//	TODO LBJ
//	TODO LIJ
//	TODO LDJ

// ---------------------------------------------------------------------------------------------------------------------
// Queuing functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO ENQ
//	TODO ENQF
//	TODO DEQ
//	TODO DEQW
//	TODO DEPOSITQB
//	TODO WITHDRAWQB

// ---------------------------------------------------------------------------------------------------------------------
// Activity Control functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LD
//	TODO SD
//	TODO LPD
//	TODO SPD
//	TODO LUD
//	TODO SUD
//	TODO LAE
//	TODO UR
//	TODO ACEL
//	TODO DCEL
//	TODO SKQT
//	TODO KCHG

// ---------------------------------------------------------------------------------------------------------------------
// Stack functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO BUY
//	TODO SELL

// ---------------------------------------------------------------------------------------------------------------------
// Interrupt Control functions
// ---------------------------------------------------------------------------------------------------------------------

// ER ------------------------------------------------------------------------------------------------------------------

const fER = 072
const jER = 011

func erSourceItemHIRef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fER),
		fmt.Sprintf("%03o", jER),
		"0",
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func erSourceItemU(label string, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fER, jER, 0, u})
}

// SGNL ----------------------------------------------------------------------------------------------------------------

const fSGNL = 073
const jSGNL = 015
const aSGNL = 017

func sgnlSourceItemHIRef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fSGNL),
		fmt.Sprintf("%03o", jSGNL),
		fmt.Sprintf("%03o", aSGNL),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func sgnlSourceItemU(label string, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxu", []int{fSGNL, jSGNL, aSGNL, u})
}

// PAIJ ----------------------------------------------------------------------------------------------------------------

const fPAIJ = 074
const jPAIJ = 014
const aPAIJ = 007

func paijSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fPAIJ),
		fmt.Sprintf("%03o", jPAIJ),
		fmt.Sprintf("%03o", aPAIJ),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func paijSourceItemRef(label string, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fPAIJ),
		fmt.Sprintf("%03o", jPAIJ),
		fmt.Sprintf("%03o", aPAIJ),
		"0",
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxu", ops)
}

// AAIJ ----------------------------------------------------------------------------------------------------------------

const fAAIJ = 074
const jAAIJExtended = 014
const jAAIJBasic = 007
const aAAIJExtended = 006
const aAAIJBasic = 000

func aaijSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fAAIJ),
		fmt.Sprintf("%03o", jAAIJExtended),
		fmt.Sprintf("%03o", aAAIJExtended),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func aaijSourceItemRef(label string, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fAAIJ),
		fmt.Sprintf("%03o", jAAIJBasic),
		fmt.Sprintf("%03o", aAAIJBasic),
		"0",
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxu", ops)
}

// ---------------------------------------------------------------------------------------------------------------------
// System Control functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO SPID
//	TODO IPC
//	TODO SPC

// IAR -----------------------------------------------------------------------------------------------------------------

const fIAR = 073
const jIAR = 017
const aIAR = 006

func iarSourceItem(uField int) *tasm.SourceItem {
	return fjaxuSourceItem(fIAR, jIAR, aIAR, 0, uField)
}

// ---------------------------------------------------------------------------------------------------------------------
// Dayclock functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LRD
//	TODO SMD
//	TODO RMD
//	TODO LMC
//	TODO SDMN
//	TODO SDMS
//	TODO SDMF
//	TODO RDC

// ---------------------------------------------------------------------------------------------------------------------
// UPI functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO SEND
//	TODO ACK

// ---------------------------------------------------------------------------------------------------------------------
// System Instrumentation functions
// ---------------------------------------------------------------------------------------------------------------------

//	TODO LBRX
//	TODO CJHE
//	TODO SJH

// ---------------------------------------------------------------------------------------------------------------------
// Special functions
// ---------------------------------------------------------------------------------------------------------------------

// DCB -----------------------------------------------------------------------------------------------------------------

const fDCB = 033
const jDCB = 015

func dcbSourceItemHIBDRef(label string, a uint64, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fDCB),
		fmt.Sprintf("%03o", jDCB),
		fmt.Sprintf("%03o", a),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// EX ------------------------------------------------------------------------------------------------------------------

const fEXBasic = 072
const fEXExtended = 073
const jEXBasic = 010
const jEXExtended = 014
const aEXExtended = 005

func exSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fEXExtended),
		fmt.Sprintf("%03o", jEXExtended),
		fmt.Sprintf("%03o", aEXExtended),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

func exSourceItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fEXBasic),
		fmt.Sprintf("%03o", jEXBasic),
		fmt.Sprintf("%03o", 0),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

// EXR -----------------------------------------------------------------------------------------------------------------

const fEXR = 073
const jEXR = 014
const aEXR = 006

func exrSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fEXR),
		fmt.Sprintf("%03o", jEXR),
		fmt.Sprintf("%03o", aEXR),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// NOP ----------------------------------------------------------------------------------------------------------------0

const fNOPBasic = 074
const fNOPExtended = 073
const jNOPBasic = 006
const jNOPExtended = 014
const aNOP = 000

func nopItemHIU(label string, x uint64, h uint64, i uint64, u int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhiu", []int{fNOPBasic, jNOPBasic, 0, x, h, i, u})
}

func nopItemHIURef(label string, x uint64, h uint64, i uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fNOPBasic),
		fmt.Sprintf("%03o", jNOPBasic),
		"0",
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhiu", ops)
}

func nopItemHIBD(label string, x uint64, h uint64, i uint64, b uint64, d int) *tasm.SourceItem {
	return sourceItem(label, "fjaxhibd", []int{fNOPExtended, jNOPExtended, aNOP, x, h, i, b, d})
}

func nopItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fNOPBasic),
		fmt.Sprintf("%03o", jNOPBasic),
		"0",
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// RNGB ----------------------------------------------------------------------------------------------------------------

const fRNGB = 037
const jRNGB = 004
const aRNGB = 006

func rngbSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fRNGB),
		fmt.Sprintf("%03o", jRNGB),
		fmt.Sprintf("%03o", aRNGB),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// RNGI ----------------------------------------------------------------------------------------------------------------

const fRNGI = 037
const jRNGI = 004
const aRNGI = 005

func rngiSourceItemHIBDRef(label string, x uint64, h uint64, i uint64, b uint64, ref string) *tasm.SourceItem {
	ops := []string{
		fmt.Sprintf("%03o", fRNGI),
		fmt.Sprintf("%03o", jRNGI),
		fmt.Sprintf("%03o", aRNGI),
		fmt.Sprintf("%03o", x),
		fmt.Sprintf("%03o", h),
		fmt.Sprintf("%03o", i),
		fmt.Sprintf("%03o", b),
		ref,
	}
	return tasm.NewSourceItem(label, "fjaxhibd", ops)
}

// ---------------------------------------------------------------------------------------------------------------------
// Convenience methods
// ---------------------------------------------------------------------------------------------------------------------

func checkInterrupt(t *testing.T, engine *InstructionEngine, interruptClass pkg.InterruptClass) {
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
	interruptClass pkg.InterruptClass,
	shortStatusField pkg.InterruptShortStatus) {

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

func checkMemory(t *testing.T, engine *InstructionEngine, addr *pkg.AbsoluteAddress, offset uint64, expected uint64) {
	seg, interrupt := engine.mainStorage.GetSegment(addr.GetSegment())
	if interrupt != nil {
		engine.mainStorage.Dump()
		t.Errorf("Error:%s", pkg.GetInterruptString(interrupt))
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
		t.Errorf("Register %s is %012o, expected %012o", pkg.RegisterNames[regIndex], result, expected)
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
