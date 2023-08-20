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
	AGRSComponent
	AFunctionDiscriminator
	AUnused
)

const (
	JPartialWordDesignator JFieldUsage = iota
	JGRSComponent
	JFunctionDiscriminator
	JUnused
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
	mnemonic  string
	aField    AFieldUsage
	jField    JFieldUsage
	uIs18Bits bool
}

func (i *Instruction) Interpret(iw *pkg.InstructionWord, basicMode bool, quarterWordMode bool) (string, bool) {
	str := i.mnemonic
	var immediate bool

	if i.jField == JPartialWordDesignator {
		str += ","
		if quarterWordMode {
			str += jFieldQuarterWord[iw.GetJ()]
		} else {
			str += jFieldThirdWord[iw.GetJ()]
		}
		immediate = (iw.GetX() == 0) && ((iw.GetJ() == pkg.JFieldU) || (iw.GetJ() == pkg.JFieldXU))
	}
	str = fmt.Sprintf("%-10s", str)

	if i.aField != AFunctionDiscriminator && i.aField != AUnused {
		str += fmt.Sprintf("%s%d", aFieldPrefix[i.aField], iw.GetA()) + ","
	}

	if basicMode {
		if immediate {
			str += fmt.Sprintf("0%o", iw.GetHIU())
		} else {
			if iw.GetI() > 0 {
				str += "*"
			}
			str += fmt.Sprintf("0%o", iw.GetU())
		}
	} else /* !basicMode */ {
		if immediate {
			str += fmt.Sprintf("0%o", iw.GetHIU())
		} else {
			str += fmt.Sprintf("0%o", iw.GetD())
		}
	}

	str += ","
	if iw.GetX() > 0 {
		if (iw.GetH() > 0) && !immediate {
			str += "*"
		}
		str += fmt.Sprintf("X%d", iw.GetX())
	}

	if !basicMode && !i.uIs18Bits && !immediate {
		str += fmt.Sprintf(",B%d", iw.GetB())
	}

	lastX := len(str) - 1
	if str[lastX:lastX+1] == "," {
		str = str[0:lastX]
	}
	return str, true
}

func (i *Instruction) IsInstruction() bool {
	return true
}

//	Basic --------------------------------------------------------------------------------------------------------------

var BasicFunctionTable = FunctionTable{
	indexBy: IndexByF,
	table: map[int]Interpreter{
		001: &Instruction{mnemonic: "SA", aField: ARegister, jField: JPartialWordDesignator},
		002: &Instruction{mnemonic: "SNA", aField: ARegister, jField: JPartialWordDesignator},
		003: &Instruction{mnemonic: "SMA", aField: ARegister, jField: JPartialWordDesignator},
		004: &Instruction{mnemonic: "SR", aField: RRegister, jField: JPartialWordDesignator},
		005: &function005InterpreterBasic,
		006: &Instruction{mnemonic: "SX", aField: XRegister, jField: JPartialWordDesignator},
		007: &function007InterpreterBasic,
		010: &Instruction{mnemonic: "LA", aField: ARegister, jField: JPartialWordDesignator},
		011: &Instruction{mnemonic: "LNA", aField: ARegister, jField: JPartialWordDesignator},
		012: &Instruction{mnemonic: "LMA", aField: ARegister, jField: JPartialWordDesignator},
		013: &Instruction{mnemonic: "LNMA", aField: ARegister, jField: JPartialWordDesignator},
		023: &Instruction{mnemonic: "LR", aField: RRegister, jField: JPartialWordDesignator},
		026: &Instruction{mnemonic: "LXM", aField: XRegister, jField: JPartialWordDesignator},
		027: &Instruction{mnemonic: "LX", aField: XRegister, jField: JPartialWordDesignator},
		044: &Instruction{mnemonic: "TEP", aField: RRegister, jField: JPartialWordDesignator},
		045: &Instruction{mnemonic: "TOP", aField: RRegister, jField: JPartialWordDesignator},
		046: &Instruction{mnemonic: "LXI", aField: XRegister},
		070: &Instruction{mnemonic: "JGD", aField: AGRSComponent, jField: JGRSComponent, uIs18Bits: true},
		071: &function071InterpreterBasic,
		072: &function072InterpreterBasic,
		073: &function073InterpreterBasic,
		074: &function074InterpreterBasic,
		075: &function075InterpreterBasic,
	},
}

var function005InterpreterBasic = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "SZ", aField: AFunctionDiscriminator},
		001: &Instruction{mnemonic: "SNZ", aField: AFunctionDiscriminator},
		002: &Instruction{mnemonic: "SP1", aField: AFunctionDiscriminator},
		003: &Instruction{mnemonic: "SN1", aField: AFunctionDiscriminator},
		004: &Instruction{mnemonic: "SFS", aField: AFunctionDiscriminator},
		005: &Instruction{mnemonic: "SFZ", aField: AFunctionDiscriminator},
		006: &Instruction{mnemonic: "SAS", aField: AFunctionDiscriminator},
		007: &Instruction{mnemonic: "SAZ", aField: AFunctionDiscriminator},
	},
}

var function007InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		004: &Instruction{mnemonic: "LAQW", aField: ARegister, jField: JFunctionDiscriminator},
		005: &Instruction{mnemonic: "SAQW", aField: ARegister, jField: JFunctionDiscriminator},
	},
}

var function071InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		012: &Instruction{mnemonic: "DS", aField: ARegister, jField: JFunctionDiscriminator},
		013: &Instruction{mnemonic: "DL", aField: ARegister, jField: JFunctionDiscriminator},
		014: &Instruction{mnemonic: "DLN", aField: ARegister, jField: JFunctionDiscriminator},
		015: &Instruction{mnemonic: "DLM", aField: ARegister, jField: JFunctionDiscriminator},
		016: &Instruction{mnemonic: "DJZ", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
	},
}

var function072InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		001: &Instruction{mnemonic: "SLJ", aField: AUnused, jField: JFunctionDiscriminator},
		002: &Instruction{mnemonic: "JPS", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		003: &Instruction{mnemonic: "JNS", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		010: &Instruction{mnemonic: "EX", aField: AUnused, jField: JFunctionDiscriminator},
		016: &Instruction{mnemonic: "SRS", aField: ARegister, jField: JFunctionDiscriminator},
		017: &Instruction{mnemonic: "LRS", aField: ARegister, jField: JFunctionDiscriminator},
	},
}

var function073InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		015: &function07315InterpreterBasic,
		017: &function07317InterpreterBasic,
	},
}

var function07315InterpreterBasic = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		014: &Instruction{mnemonic: "LD", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator},
		015: &Instruction{mnemonic: "SD", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator},
	},
}

var function07317InterpreterBasic = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		006: &Instruction{mnemonic: "IAR", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator},
	},
}

var function074InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "JZ", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		001: &Instruction{mnemonic: "JNZ", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		002: &Instruction{mnemonic: "JP", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		003: &Instruction{mnemonic: "JN", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		004: &function07404InterpreterBasic,
		005: &Instruction{mnemonic: "HKJ", aField: AUnused, jField: JFunctionDiscriminator, uIs18Bits: true},
		006: &Instruction{mnemonic: "NOP", aField: AUnused, jField: JFunctionDiscriminator},
		010: &Instruction{mnemonic: "JNLB", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		011: &Instruction{mnemonic: "JLB", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		012: &Instruction{mnemonic: "JMGI", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		013: &Instruction{mnemonic: "LMJ", aField: XRegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		014: &function07414InterpreterBasic,
		015: &function07415InterpreterBasic,
		016: &Instruction{mnemonic: "JC", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		017: &Instruction{mnemonic: "JNC", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
	},
}

var function07404InterpreterBasic = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "J", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		001: &Instruction{mnemonic: "JK01", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		002: &Instruction{mnemonic: "JK02", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		003: &Instruction{mnemonic: "JK03", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		004: &Instruction{mnemonic: "JK04", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		005: &Instruction{mnemonic: "JK05", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		006: &Instruction{mnemonic: "JK06", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		007: &Instruction{mnemonic: "JK07", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		010: &Instruction{mnemonic: "JK10", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		011: &Instruction{mnemonic: "JK11", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		012: &Instruction{mnemonic: "JK12", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		013: &Instruction{mnemonic: "JK13", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		014: &Instruction{mnemonic: "JK14", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		015: &Instruction{mnemonic: "JK15", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		016: &Instruction{mnemonic: "JK16", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
		017: &Instruction{mnemonic: "JK17", aField: AFunctionDiscriminator, jField: JUnused, uIs18Bits: true},
	},
}

var function07414InterpreterBasic = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "JO", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		001: &Instruction{mnemonic: "JFU", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		002: &Instruction{mnemonic: "JFO", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		003: &Instruction{mnemonic: "JDF", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
	},
}

var function07415InterpreterBasic = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "JNO", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		001: &Instruction{mnemonic: "JNFU", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		002: &Instruction{mnemonic: "JNFO", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		003: &Instruction{mnemonic: "JNDF", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		005: &Instruction{mnemonic: "HLTJ", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
	},
}

var function075InterpreterBasic = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		013: &Instruction{mnemonic: "LXLM", aField: XRegister, jField: JFunctionDiscriminator},
	},
}

//	Extended -----------------------------------------------------------------------------------------------------------

var ExtendedFunctionTable = FunctionTable{
	indexBy: IndexByF,
	table: map[int]Interpreter{
		001: &Instruction{mnemonic: "SA", aField: ARegister, jField: JPartialWordDesignator},
		002: &Instruction{mnemonic: "SNA", aField: ARegister, jField: JPartialWordDesignator},
		003: &Instruction{mnemonic: "SMA", aField: ARegister, jField: JPartialWordDesignator},
		004: &Instruction{mnemonic: "SR", aField: RRegister, jField: JPartialWordDesignator},
		005: &function005InterpreterExtended,
		006: &Instruction{mnemonic: "SX", aField: XRegister, jField: JPartialWordDesignator},
		007: &function007InterpreterExtended,
		010: &Instruction{mnemonic: "LA", aField: ARegister, jField: JPartialWordDesignator},
		011: &Instruction{mnemonic: "LNA", aField: ARegister, jField: JPartialWordDesignator},
		012: &Instruction{mnemonic: "LMA", aField: ARegister, jField: JPartialWordDesignator},
		013: &Instruction{mnemonic: "LNMA", aField: ARegister, jField: JPartialWordDesignator},
		023: &Instruction{mnemonic: "LR", aField: RRegister, jField: JPartialWordDesignator},
		026: &Instruction{mnemonic: "LXM", aField: XRegister, jField: JPartialWordDesignator},
		027: &Instruction{mnemonic: "LX", aField: XRegister, jField: JPartialWordDesignator},
		044: &Instruction{mnemonic: "TEP", aField: RRegister, jField: JPartialWordDesignator},
		045: &Instruction{mnemonic: "TOP", aField: RRegister, jField: JPartialWordDesignator},
		046: &Instruction{mnemonic: "LXI", aField: XRegister, jField: JPartialWordDesignator},
		047: &Instruction{mnemonic: "TLEM", aField: XRegister, jField: JPartialWordDesignator},
		050: &function050InterpreterExtended,
		051: &Instruction{mnemonic: "LXSI", aField: XRegister, jField: JPartialWordDesignator},
		060: &Instruction{mnemonic: "LSBO", aField: XRegister, jField: JPartialWordDesignator},
		061: &Instruction{mnemonic: "LSBL", aField: XRegister, jField: JPartialWordDesignator},
		070: &Instruction{mnemonic: "JGD", aField: AGRSComponent, jField: JGRSComponent, uIs18Bits: true},
		071: &function071InterpreterExtended,
		072: &function072InterpreterExtended,
		073: &function073InterpreterExtended,
		074: &function074InterpreterExtended,
		075: &function075InterpreterExtended,
	},
}

var function005InterpreterExtended = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "SZ", aField: AFunctionDiscriminator},
		001: &Instruction{mnemonic: "SNZ", aField: AFunctionDiscriminator},
		002: &Instruction{mnemonic: "SP1", aField: AFunctionDiscriminator},
		003: &Instruction{mnemonic: "SN1", aField: AFunctionDiscriminator},
		004: &Instruction{mnemonic: "SFS", aField: AFunctionDiscriminator},
		005: &Instruction{mnemonic: "SFZ", aField: AFunctionDiscriminator},
		006: &Instruction{mnemonic: "SAS", aField: AFunctionDiscriminator},
		007: &Instruction{mnemonic: "SAZ", aField: AFunctionDiscriminator},
	},
}

var function007InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		004: &Instruction{mnemonic: "LAQW", aField: ARegister, jField: JFunctionDiscriminator},
		005: &Instruction{mnemonic: "SAQW", aField: ARegister, jField: JFunctionDiscriminator},
	},
}

var function050InterpreterExtended = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		017: &Instruction{mnemonic: "TSKP", aField: AUnused, jField: JFunctionDiscriminator},
	},
}

var function071InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		012: &Instruction{mnemonic: "DS", aField: ARegister, jField: JFunctionDiscriminator},
		013: &Instruction{mnemonic: "DL", aField: ARegister, jField: JFunctionDiscriminator},
		014: &Instruction{mnemonic: "DLN", aField: ARegister, jField: JFunctionDiscriminator},
		015: &Instruction{mnemonic: "DLM", aField: ARegister, jField: JFunctionDiscriminator},
		016: &Instruction{mnemonic: "DJZ", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
	},
}

var function072InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		002: &Instruction{mnemonic: "JPS", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		003: &Instruction{mnemonic: "JNS", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		016: &Instruction{mnemonic: "SRS", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		017: &Instruction{mnemonic: "LRS", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
	},
}

var function073InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		014: &function07314InterpreterExtended,
		015: &function07315InterpreterExtended,
		017: &function07317InterpreterExtended,
	},
}

var function07314InterpreterExtended = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "NOP", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator},
		005: &Instruction{mnemonic: "EX", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator},
	},
}

var function07315InterpreterExtended = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		014: &Instruction{mnemonic: "LD", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator},
		015: &Instruction{mnemonic: "SD", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator},
	},
}

var function07317InterpreterExtended = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		006: &Instruction{mnemonic: "IAR", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator},
	},
}

var function074InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "JZ", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		001: &Instruction{mnemonic: "JNZ", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		002: &Instruction{mnemonic: "JP", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		003: &Instruction{mnemonic: "JN", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		010: &Instruction{mnemonic: "JNLB", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		011: &Instruction{mnemonic: "JLB", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		012: &Instruction{mnemonic: "JMGI", aField: ARegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		013: &Instruction{mnemonic: "LMJ", aField: XRegister, jField: JFunctionDiscriminator, uIs18Bits: true},
		014: &function07414InterpreterExtended,
		015: &function07415InterpreterExtended,
	},
}

var function07414InterpreterExtended = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "JO", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		001: &Instruction{mnemonic: "JFU", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		002: &Instruction{mnemonic: "JFO", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		003: &Instruction{mnemonic: "JDF", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		004: &Instruction{mnemonic: "JC", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		005: &Instruction{mnemonic: "JNC", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
	},
}

var function07415InterpreterExtended = FunctionTable{
	indexBy: IndexByA,
	table: map[int]Interpreter{
		000: &Instruction{mnemonic: "JNO", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		001: &Instruction{mnemonic: "JNFU", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		002: &Instruction{mnemonic: "JNFO", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		003: &Instruction{mnemonic: "JNDF", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		004: &Instruction{mnemonic: "J", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
		005: &Instruction{mnemonic: "HLTJ", aField: AFunctionDiscriminator, jField: JFunctionDiscriminator, uIs18Bits: true},
	},
}

var function075InterpreterExtended = FunctionTable{
	indexBy: IndexByJ,
	table: map[int]Interpreter{
		013: &Instruction{mnemonic: "LXLM", aField: XRegister, jField: JFunctionDiscriminator},
	},
}

//	Other stuff --------------------------------------------------------------------------------------------------------

var jFieldThirdWord = []string{
	"W", "H2", "H1", "XH2", "XH1", "T3", "T2", "T1", "S6", "S5", "S4", "S3", "S2", "S1", "U", "XU",
}

var jFieldQuarterWord = []string{
	"W", "H2", "H1", "XH2", "Q2", "Q4", "Q3", "Q1", "S6", "S5", "S4", "S3", "S2", "S1", "U", "XU",
}
