// Khalehla Project
// disassembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package dasm

import (
	"fmt"
	"khalehla/pkg"
)

type AFieldUsage int
type JFieldUsage int
type IndexField int

const (
	ARegister AFieldUsage = iota
	BRegister
	RRegister
	XRegister
	AFunctionDiscriminator
)

const (
	PartialWordDesignator JFieldUsage = iota
	JFunctionDiscriminator
)

const (
	IndexByF IndexField = iota
	IndexByJ
	IndexByA
)

type Interpreter interface {
	Interpret(word *pkg.InstructionWord, basicMode bool, quarterWordMode bool) (string, bool)
	IsInstruction() bool
}

type FunctionTable struct {
	table   map[int]Interpreter
	indexBy IndexField
}

var aFieldPrefix = map[AFieldUsage]string{
	ARegister: "A",
	BRegister: "B",
	RRegister: "R",
	XRegister: "X",
}

func (ft *FunctionTable) Interpret(iw *pkg.InstructionWord, basicMode bool, quarterWordMode bool) (string, bool) {
	var function Interpreter
	var ok bool

	if ft.indexBy == IndexByF {
		function, ok = ft.table[int(iw.GetF())]
	} else if ft.indexBy == IndexByJ {
		function, ok = ft.table[int(iw.GetJ())]
	} else if ft.indexBy == IndexByA {
		function, ok = ft.table[int(iw.GetA())]
	}

	if !ok {
		return "", false
	}

	if function.IsInstruction() {
		i := function.(*Instruction)
		return i.Interpret(iw, basicMode, quarterWordMode)
	} else {
		ft := function.(*FunctionTable)
		return ft.Interpret(iw, basicMode, quarterWordMode)
	}
}

func (ft *FunctionTable) IsInstruction() bool {
	return false
}

type Instruction struct {
	mnemonic string
	mode     Mode
	aField   AFieldUsage
	jField   JFieldUsage
}

func (i *Instruction) Interpret(iw *pkg.InstructionWord, basicMode bool, quarterWordMode bool) (string, bool) {
	str := i.mnemonic
	if i.jField == PartialWordDesignator {
		str += ","
		if quarterWordMode {
			str += jFieldQuarterWord[iw.GetJ()]
		} else {
			str += jFieldThirdWord[iw.GetJ()]
		}
	}
	str = fmt.Sprintf("%-10s", str)

	if i.aField != AFunctionDiscriminator {
		str += fmt.Sprintf("%s%d", aFieldPrefix[i.aField], iw.GetA()) + ","
	}

	if basicMode {
		if iw.GetI() > 0 {
			str += "*"
		}
		str += fmt.Sprintf("0%o", iw.GetU())
	} else {
		str += fmt.Sprintf("0%o", iw.GetD())
	}

	str += ","
	if iw.GetX() > 0 {
		if iw.GetH() > 0 {
			str += "*"
		}
		str += fmt.Sprintf("X%d", iw.GetX())
	}

	if !basicMode {
		str += fmt.Sprintf(",B%d", iw.GetB())
	}

	return str, true
}

func (i *Instruction) IsInstruction() bool {
	return true
}

var BasicFunctionTable = FunctionTable{
	indexBy: IndexByF,
	table: map[int]Interpreter{
		001: &Instruction{mnemonic: "SA", aField: ARegister, mode: BOTH},
		002: &Instruction{mnemonic: "SNA", aField: ARegister, mode: BOTH},
		003: &Instruction{mnemonic: "SMA", aField: ARegister, mode: BOTH},
		004: &Instruction{mnemonic: "SR", aField: RRegister, mode: BOTH},
		005: &function005InterpreterBasic,
		006: &Instruction{mnemonic: "SX", aField: XRegister, mode: BOTH},
		007: &function007InterpreterBasic,
		010: &Instruction{mnemonic: "LA", aField: ARegister, mode: BOTH},
		011: &Instruction{mnemonic: "LNA", aField: ARegister, mode: BOTH},
		012: &Instruction{mnemonic: "LMA", aField: ARegister, mode: BOTH},
		013: &Instruction{mnemonic: "LNMA", aField: ARegister, mode: BOTH},
		023: &Instruction{mnemonic: "LR", aField: RRegister, mode: BOTH},
		026: &Instruction{mnemonic: "LXM", aField: XRegister, mode: BOTH},
		027: &Instruction{mnemonic: "LX", aField: XRegister, mode: BOTH},
		046: &Instruction{mnemonic: "LXI", aField: XRegister, mode: BOTH},
		071: &function071InterpreterBasic,
		072: &function072InterpreterBasic,
		073: &function073InterpreterBasic,
		075: &function075InterpreterBasic,
	},
}

var ExtendedFunctionTable = FunctionTable{
	indexBy: IndexByF,
	table: map[int]Interpreter{
		001: &Instruction{mnemonic: "SA", aField: ARegister, mode: BOTH},
		002: &Instruction{mnemonic: "SNA", aField: ARegister, mode: BOTH},
		003: &Instruction{mnemonic: "SMA", aField: ARegister, mode: BOTH},
		004: &Instruction{mnemonic: "SR", aField: RRegister, mode: BOTH},
		005: &function005InterpreterExtended,
		006: &Instruction{mnemonic: "SX", aField: XRegister, mode: BOTH},
		007: &function007InterpreterExtended,
		010: &Instruction{mnemonic: "LA", aField: ARegister, mode: BOTH},
		011: &Instruction{mnemonic: "LNA", aField: ARegister, mode: BOTH},
		012: &Instruction{mnemonic: "LMA", aField: ARegister, mode: BOTH},
		013: &Instruction{mnemonic: "LNMA", aField: ARegister, mode: BOTH},
		023: &Instruction{mnemonic: "LR", aField: RRegister, mode: BOTH},
		026: &Instruction{mnemonic: "LXM", aField: XRegister, mode: BOTH},
		027: &Instruction{mnemonic: "LX", aField: XRegister, mode: BOTH},
		046: &Instruction{mnemonic: "LXI", aField: XRegister, mode: BOTH},
		051: &Instruction{mnemonic: "LXSI", aField: XRegister, mode: EXTENDED},
		060: &Instruction{mnemonic: "LSBO", aField: XRegister, mode: EXTENDED},
		061: &Instruction{mnemonic: "LSBL", aField: XRegister, mode: EXTENDED},
		071: &function071InterpreterExtended,
		072: &function072InterpreterExtended,
		073: &function073InterpreterExtended,
		075: &function075InterpreterExtended,
	},
}

var function005InterpreterBasic = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "SZ", aField: AFunctionDiscriminator, mode: BOTH},
		001: &Instruction{mnemonic: "SNZ", aField: AFunctionDiscriminator, mode: BOTH},
		002: &Instruction{mnemonic: "SP1", aField: AFunctionDiscriminator, mode: BOTH},
		003: &Instruction{mnemonic: "SN1", aField: AFunctionDiscriminator, mode: BOTH},
		004: &Instruction{mnemonic: "SFS", aField: AFunctionDiscriminator, mode: BOTH},
		005: &Instruction{mnemonic: "SFZ", aField: AFunctionDiscriminator, mode: BOTH},
		006: &Instruction{mnemonic: "SAS", aField: AFunctionDiscriminator, mode: BOTH},
		007: &Instruction{mnemonic: "SAZ", aField: AFunctionDiscriminator, mode: BOTH},
	},
}

var function005InterpreterExtended = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "SZ", aField: AFunctionDiscriminator, mode: BOTH},
		001: &Instruction{mnemonic: "SNZ", aField: AFunctionDiscriminator, mode: BOTH},
		002: &Instruction{mnemonic: "SP1", aField: AFunctionDiscriminator, mode: BOTH},
		003: &Instruction{mnemonic: "SN1", aField: AFunctionDiscriminator, mode: BOTH},
		004: &Instruction{mnemonic: "SFS", aField: AFunctionDiscriminator, mode: BOTH},
		005: &Instruction{mnemonic: "SFZ", aField: AFunctionDiscriminator, mode: BOTH},
		006: &Instruction{mnemonic: "SAS", aField: AFunctionDiscriminator, mode: BOTH},
		007: &Instruction{mnemonic: "SAZ", aField: AFunctionDiscriminator, mode: BOTH},
	},
}

var function007InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		004: &Instruction{mnemonic: "LAQW", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		005: &Instruction{mnemonic: "SAQW", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var function007InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		004: &Instruction{mnemonic: "LAQW", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		005: &Instruction{mnemonic: "SAQW", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var function071InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		012: &Instruction{mnemonic: "DS", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		013: &Instruction{mnemonic: "DL", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		014: &Instruction{mnemonic: "DLN", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		015: &Instruction{mnemonic: "DLM", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var function071InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		012: &Instruction{mnemonic: "DS", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		013: &Instruction{mnemonic: "DL", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		014: &Instruction{mnemonic: "DLN", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		015: &Instruction{mnemonic: "DLM", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var function072InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		016: &Instruction{mnemonic: "SRS", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		017: &Instruction{mnemonic: "LRS", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var function072InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		016: &Instruction{mnemonic: "SRS", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
		017: &Instruction{mnemonic: "LRS", aField: ARegister, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var function073InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		015: &function07315InterpreterBasic,
	},
}

var function073InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		015: &function07315InterpreterExtended,
	},
}

var function07315InterpreterBasic = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		014: &Instruction{mnemonic: "LD", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, mode: BOTH},
		015: &Instruction{mnemonic: "SD", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var function07315InterpreterExtended = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		014: &Instruction{mnemonic: "LD", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, mode: BOTH},
		015: &Instruction{mnemonic: "SD", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var function075InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		013: &Instruction{mnemonic: "LXLM", aField: XRegister, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var function075InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		013: &Instruction{mnemonic: "LXLM", aField: XRegister, jField: JFunctionDiscriminator, mode: BOTH},
	},
}

var jFieldThirdWord = []string{
	"W", "H2", "H1", "XH2", "XH1", "T3", "T2", "T1", "S6", "S5", "S4", "S3", "S2", "S1", "U", "XU",
}

var jFieldQuarterWord = []string{
	"W", "H2", "H1", "XH2", "Q2", "Q4", "Q3", "Q1", "S6", "S5", "S4", "S3", "S2", "S1", "U", "XU",
}
