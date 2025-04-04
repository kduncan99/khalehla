// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"math/rand"

	"khalehla/common"
)

// Instructions allowed as targets for EXR - by f field and optionally by j field and/or a field
// only applies to extended mode, since EXR is not a basic mode thing

type exrHandler interface {
	isAllowed(f uint64, j uint64, a uint64) bool
}

type exrFHandler struct {
	exrMap map[uint64]exrHandler
}

type exrJHandler struct {
	exrMap map[uint64]exrHandler
}

type exrAHandler struct {
	exrMap map[uint64]exrHandler
}

func (eh *exrFHandler) isAllowed(f uint64, j uint64, a uint64) bool {
	sub, ok := eh.exrMap[f]
	if ok {
		if sub != nil {
			return sub.isAllowed(f, j, a)
		} else {
			return true
		}
	}
	return false
}

func (eh *exrJHandler) isAllowed(f uint64, j uint64, a uint64) bool {
	sub, ok := eh.exrMap[j]
	if ok {
		if sub != nil {
			return sub.isAllowed(f, j, a)
		} else {
			return true
		}
	}
	return false
}

func (eh *exrAHandler) isAllowed(f uint64, j uint64, a uint64) bool {
	sub, ok := eh.exrMap[a]
	if ok {
		if sub != nil {
			return sub.isAllowed(f, j, a)
		} else {
			return true
		}
	}
	return false
}

var exrMainHandler = &exrFHandler{exrMap: exrFMap}
var exrF005Handler = &exrAHandler{exrMap: exrF005Map}
var exrF007Handler = &exrJHandler{exrMap: exrF007Map}
var exrF033Handler = &exrJHandler{exrMap: exrF033Map}
var exrF037Handler = &exrJHandler{exrMap: exrF037Map}
var exrF037J004Handler = &exrAHandler{exrMap: exrF037J004Map}
var exrF050Handler = &exrJHandler{exrMap: exrF050Map}
var exrF071Handler = &exrJHandler{exrMap: exrF071Map}
var exrF072Handler = &exrJHandler{exrMap: exrF072Map}
var exrF076Handler = &exrJHandler{exrMap: exrF076Map}

var exrFMap = map[uint64]exrHandler{ // indexed by f
	001: nil, // SA
	002: nil, // SNA
	003: nil, // SMA
	004: nil, // SR
	005: exrF005Handler,
	006: nil, // SX
	007: exrF007Handler,
	014: nil, // AA
	015: nil, // ANA
	016: nil, // AMA
	017: nil, // ANMA
	024: nil, // AX
	025: nil, // ANX
	031: nil, // MSI
	033: exrF033Handler,
	037: exrF037Handler,
	044: nil, // TEP
	045: nil, // TOP
	050: exrF050Handler,
	052: nil, // TE
	053: nil, // TNE
	054: nil, // TLE
	055: nil, // TG
	056: nil, // TW
	057: nil, // TNW
	071: exrF071Handler,
	072: exrF072Handler,
	076: exrF076Handler,
}

var exrF005Map = map[uint64]exrHandler{ // indexed by a
	000: nil, // SZ
	001: nil, // SNZ
	002: nil, // SP1
	003: nil, // SN1
	004: nil, // SFS
	005: nil, // SFZ
	006: nil, // SAS
	007: nil, // SAZ
}

var exrF007Map = map[uint64]exrHandler{ // indexed by j
	000: nil, // ADE
	001: nil, // DADE
	002: nil, // SDE
	003: nil, // DSDE
}

var exrF033Map = map[uint64]exrHandler{ // indexed by j
	012: nil, // SS
	013: nil, // TGM
	014: nil, // DTGM
	016: nil, // TES
	017: nil, // TNES
}

var exrF037Map = map[uint64]exrHandler{ // indexed by j
	004: exrF037J004Handler,
}

var exrF037J004Map = map[uint64]exrHandler{ // indexed by a
	005: nil, // RNGI
	006: nil, // RNGB
}

var exrF050Map = map[uint64]exrHandler{ // indexed by j
	000: nil, // TNOP
	001: nil, // TGZ
	002: nil, // TPZ
	003: nil, // TP
	004: nil, // TMZ
	005: nil, // TMZG
	006: nil, // TZ
	007: nil, // TNLZ
	010: nil, // TLZ
	011: nil, // TNZ
	012: nil, // TPZL
	013: nil, // TNMZ
	014: nil, // TN
	015: nil, // TNPZ
	016: nil, // TNGZ
	017: nil, // TSKP
}

var exrF071Map = map[uint64]exrHandler{ // indexed by j
	000: nil, // MTE
	001: nil, // MTNE
	002: nil, // MTLE
	003: nil, // MTG
	004: nil, // MTW
	005: nil, // MTNW
	006: nil, // MATL
	007: nil, // MATG
	010: nil, // DA
	011: nil, // DAN
	012: nil, // DS
	017: nil, // DTE
}

var exrF072Map = map[uint64]exrHandler{ // indexed by j
	004: nil, // AH
	005: nil, // ANH
	006: nil, // AT
	007: nil, // ANT
}

var exrF076Map = map[uint64]exrHandler{ // indexed by j
	000: nil, // FA
	001: nil, // FAN
	002: nil, // FM
	003: nil, // FD
	010: nil, // DFA
	011: nil, // DFAN
	012: nil, // DFM
	013: nil, // DFD
}

// DoubleCountBits (DCB) counts the number of bits set in U and U+1, storing the result in Aa
func DoubleCountBits(e *InstructionEngine) (completed bool) {
	result := e.GetConsecutiveOperands(true, 2, false)
	completed = result.complete

	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		count := result.source[0].CountBits() + result.source[1].CountBits()
		ci := e.GetCurrentInstruction()
		e.GetExecOrUserARegister(ci.GetA()).SetW(count)
	}

	return
}

// Execute (EX) fetches the instruction at U (non-GRS) and executes it as if it were at the current program address.
func Execute(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, false, false, false, false)
	completed = result.complete

	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		e.activityStatePacket.GetCurrentInstruction().SetW(result.operand)
		e.cachedInstructionHandler = nil
		completed = false // this forces another cycle which will now execute the new instruction in F0.
	}

	return
}

// ExecuteRepeated (EXR) executes the instruction at U (storage only, not GRS) the number of times specified by R1
func ExecuteRepeated(e *InstructionEngine) (completed bool) {
	result := e.GetOperand(false, false, false, false, false)
	completed = result.complete

	if result.interrupt != nil {
		e.PostInterrupt(result.interrupt)
	} else if result.complete {
		e.activityStatePacket.GetCurrentInstruction().SetW(result.operand)
		e.cachedInstructionHandler = nil

		ci := common.InstructionWord(result.operand)
		if !exrMainHandler.isAllowed(ci.GetF(), ci.GetJ(), ci.GetA()) {
			i := common.NewInvalidInstructionInterrupt(common.InvalidInstructionEXRInvalidTarget)
			e.PostInterrupt(i)
			return false
		}

		e.activityStatePacket.GetIndicatorKeyRegister().SetExecuteRepeatedInstruction(true)
		rReg := e.GetExecOrUserRRegister(1)
		if !rReg.IsZero() {
			completed = false // this forces another cycle which will now execute the new instruction in F0.
		}
	}

	return
}

// NoOperation (NOP) evaluates the HIU field, but takes no other action (it does do x-register incrementation)
func NoOperation(e *InstructionEngine) (completed bool) {
	complete, i := e.IgnoreOperand()
	if i != nil {
		e.PostInterrupt(i)
	}

	return complete
}

// RandomNumberGeneratorInteger (RNGI) stores a 32-bit random integer in bits 4-35 of 4 consecutive
// U locations (GRS or storage).
func RandomNumberGeneratorInteger(e *InstructionEngine) (completed bool) {
	ops := []uint64{uint64(rand.Uint32()), uint64(rand.Uint32()), uint64(rand.Uint32()), uint64(rand.Uint32())}
	comp, i := e.StoreConsecutiveOperands(true, ops)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}

// RandomNumberGeneratorByte (RNGB) stores 16 8-bit random integers in successive quarter words of
// locations U through U+3
func RandomNumberGeneratorByte(e *InstructionEngine) (completed bool) {
	ops := make([]uint64, 4)
	for ox := 0; ox < 4; ox++ {
		v := uint64(rand.Uint32() & 0377)
		v = (v << 9) | uint64(rand.Uint32()&0377)
		v = (v << 9) | uint64(rand.Uint32()&0377)
		ops[ox] = (v << 9) | uint64(rand.Uint32()&0377)
	}

	comp, i := e.StoreConsecutiveOperands(true, ops)
	if i != nil {
		e.PostInterrupt(i)
	}

	return comp
}
