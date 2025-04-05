// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"

	"khalehla/common"
	"khalehla/dasm"
	"khalehla/hardware"
)

type BreakpointComparison uint

const (
	BreakpointFetch BreakpointComparison = 1
	BreakpointRead  BreakpointComparison = 2
	BreakpointWrite BreakpointComparison = 3
)

//	TODO need to revisit breakpoint detection- I think we are supposed to detect it, but not let it
//   impede the completion of the instruction - which would be a bit of a change to how we do stuffs...

type InstructionPoint uint

const (
	BetweenInstructions InstructionPoint = 1 //	an instruction is not in F0, previous instruction has completed.
	ResolvingAddress    InstructionPoint = 2 //	we have not yet begun processing the instruction in F0, or we are still resolving indirect addressing for an operand.
	MidInstruction      InstructionPoint = 3 //  we are processing an instruction and have reached a mid-instruction interrupt point (e.g., for EXR).
)

type StopReason uint

const (
	NotStopped StopReason = iota

	InitialStop
	ClearedStop
	DebugStop
	DevelopmentStop
	BreakpointStop
	HaltJumpExecutedStop
	ICSBaseRegisterInvalidStop
	ICSOverflowStop
	InitiateAutoRecoveryStop
	L0BaseRegisterInvalidStop
	PanelHaltStop

	//  Interrupt Handler initiated stops...

	InterruptHandlerHardwareFailureStop
	InterruptHandlerOffsetOutOfRangeStop
	InterruptHandlerInvalidBankTypeStop
	InterruptHandlerInvalidLevelBDIStop
)

const L0BDTBaseRegister = common.B16
const ICSBaseRegister = common.B26
const ICSIndexRegister = common.EX1
const RCSBaseRegister = common.B25
const RCSIndexRegister = common.EX0

type GetOperandResult struct {
	complete                bool
	operand                 uint64
	source                  *common.Word36
	sourceIsGRS             bool
	sourceBaseRegisterIndex uint
	sourceRelativeAddress   uint64
	sourceVirtualAddress    common.VirtualAddress
	sourceAbsoluteAddress   *common.AbsoluteAddress
	interrupt               common.Interrupt
}

func (gor *GetOperandResult) GetString() string {
	return fmt.Sprintf(
		"comp=%v op=%012o src=%012o grs=%v brx=%d relAddr=%012o virtAddr=%012o absAddr=%s int=%s",
		gor.complete,
		gor.operand,
		gor.source.GetW(),
		gor.sourceIsGRS,
		gor.sourceBaseRegisterIndex,
		gor.sourceRelativeAddress,
		gor.sourceVirtualAddress.GetComposite(),
		gor.sourceAbsoluteAddress.GetString(),
		common.GetInterruptString(gor.interrupt))
}

type ConsecutiveOperandsResult struct {
	complete                bool
	source                  []common.Word36
	sourceBaseRegisterIndex uint
	sourceRelativeAddress   uint64
	sourceVirtualAddress    common.VirtualAddress
	sourceAbsoluteAddress   *common.AbsoluteAddress
	interrupt               common.Interrupt
}

func (or *ConsecutiveOperandsResult) GetString() string {
	return fmt.Sprintf(
		"comp=%v src=[?] brx=%d relAddr=%012o virtAddr=%012o absAddr=%s int=%s",
		or.complete,
		or.sourceBaseRegisterIndex,
		or.sourceRelativeAddress,
		or.sourceVirtualAddress.GetComposite(),
		or.sourceAbsoluteAddress.GetString(),
		common.GetInterruptString(or.interrupt))
}

// InstructionEngine implements the basic functionality required to execute 36-bit code.
// It does not handle any actual hardware considerations such as interrupts, etc.
// It does track stop reasons, as there are certain generic processes that require this ability
// whether we are emulating hardware or an operating system with hardware.
// It must be wrapped by additional code which does this, either as an IP emulator or an OS emulator.
type InstructionEngine struct {
	name        string                // unique name of this engine - must be set externally
	mainStorage *hardware.MainStorage // must be set externally

	activeBaseTable           [16]*ActiveBaseTableEntry // [0] is unused
	activityStatePacket       *common.ActivityStatePacket
	cachedInstructionHandler  func(*InstructionEngine) (completed bool)
	baseRegisters             [32]*common.BaseRegister
	baseRegisterIndexForFetch uint // only applies to basic mode - if 0, it is not valid; otherwise it is 12:15
	generalRegisterSet        *common.GeneralRegisterSet

	//	If not nil, describes an interrupt which needs to be handled as soon as possible
	pendingInterrupts *InterruptStack
	jumpHistory       *JumpHistory

	logInstructions bool
	logInterrupts   bool

	//	If true, the current (or most recent) instruction has set the PAR.PC the way it wants,
	//	and we should not increment it for the next instruction
	preventPCUpdate bool

	breakpointAddress *common.AbsoluteAddress
	breakpointHalt    bool
	breakpointFetch   bool
	breakpointRead    bool
	breakpointWrite   bool

	isStopped        bool
	stopReason       StopReason
	stopDetail       common.Word36
	instructionPoint InstructionPoint
}

// Order of base register selection for Basic Mode address resolution
// when the Basic Mode Base Register Selection Designator Register bit is false
var baseRegisterCandidatesFalse = []uint{12, 14, 13, 15}

// Order of base register selection for Basic Mode address resolution
// when the Basic Mode Base Register Selection Designator Register bit is true
var baseRegisterCandidatesTrue = []uint{13, 15, 12, 14}

// And the composite of the above...
var baseRegisterCandidates = map[bool][]uint{
	false: baseRegisterCandidatesFalse,
	true:  baseRegisterCandidatesTrue,
}

// TODO NOTE THIS: we need to deal with this eventually...
//  For iterative instructions, U is recalculated for every iteration of the Repeat_Count_Register (R1).
//  While U is recalculated for each iteration of R1, the determination for U < 0200 is made only for
//  the initial address of an iterative operand. Some instructions have both source and destination
//  operands; the determination for U < 0200 is made once for the initial address of the iterativ
//  source operand and once for the initial address of the iterative destination operand. The iterative
//  instructions are:
//  • Block Transfer (BT; see 6.12.1)
//  • Block Add Octets (BAO; see 6.3.28)
//  • All search instructions (see 6.6)
//  • For Execute Repeated (EXR; see 6.27.2), the formation of the instruction operand of the target instruction is iterative.
//  Architecturally_Undefined: Operation is undefined if a subsequent address of the iterative
//  operand is on the other side of the GRS/Storage boundary from the initial address.

//	external stuffs ----------------------------------------------------------------------------------------------------

func NewEngine(name string, mainStorage *hardware.MainStorage) *InstructionEngine {
	e := &InstructionEngine{}
	e.name = name
	e.mainStorage = mainStorage
	e.Clear()
	return e
}

func (e *InstructionEngine) Clear() {
	e.pendingInterrupts = NewInterruptStack()
	e.jumpHistory = NewJumpHistory()

	for ax := 0; ax < 16; ax++ {
		e.activeBaseTable[ax] = NewActiveBaseTableEntryFromComposite(0)
	}

	for bx := 0; bx < 32; bx++ {
		e.baseRegisters[bx] = common.NewVoidBaseRegister()
	}
	e.baseRegisterIndexForFetch = 0

	e.generalRegisterSet = common.NewGeneralRegisterSet()
	e.activityStatePacket = common.NewActivityStatePacket()
	e.breakpointAddress = nil
	e.breakpointHalt = false
	e.breakpointFetch = false
	e.breakpointRead = false
	e.breakpointWrite = false

	e.isStopped = true
	e.stopReason = InitialStop
	e.stopDetail = 0

	e.preventPCUpdate = false
	e.instructionPoint = BetweenInstructions
}

func (e *InstructionEngine) ClearAllInterrupts() {
	e.pendingInterrupts.Clear()
}

func (e *InstructionEngine) ClearJumpHistory() {
	e.jumpHistory.Clear()
}

func (e *InstructionEngine) ClearStop() {
	e.isStopped = false
}

// DoCycle executes one cycle
// caller should disposition any pending interrupts before invoking this...
// Since the engine is not specifically hardware (could be an executor for a native mode OS),
// we don't actually know how to handle the interrupts.
// In any event, we are driven by the following two flags in the indicator key register:
//
//	INF: Indicates that a valid instruction is in F0.
//			If we are returning from an interrupt, the instruction in F0 was interrupted.
//			Otherwise, we are at a mid-execution point for the instruction in F0, or are doing address resolution.
//			PAR.PC is the address of that instruction, or of an EX or EXR instruction which invoked the instruciton
//			in F0.
//
//	EXRF: Indicates that F0 contains the target of an EXR instruction.
//		if zero, PAR.PC contains the address of the next instruction to be loaded.
//			if non-zero:
//				if we are returning from an interrupt, then the instruction in F0 was interrupted.
//				otherwise we are at a mid-execution point for the instruction in F0, or are doing address resolution.
//				PAR.PC will be the address of the instruction (or of an EX or EXR instruction which invoked the
//				instruction in F0).
//
// With INF == 0, we fetch the instruction referenced by PAR.PC. If that fails, INF and EXRF are still zero,
// and interrupt is posted, and all we have to do is return to the caller so they can manage the interrupt.
//
// With INF == 1 and EXRF == 0, we hand off to the instruction handler for processing.
//
// With INF and EXRF == 1, we check R1 and terminate EXRF processing if R1 is zero, or else we hand off to the
// instruction handler for processing if R1 is non-zero.
//
// When the instruction handler returns, we are in one of several possible situations:
//
//	The instruction never started because of an interrupt. The interrupt will be posted by the time we get here.
//		INF and EXRF will both be clear, and we just return to the caller and let him sort it out (see next item).
//	The instruction is not complete, and it posted an interrupt. INF is set, EXRF may be set or clear.
//		We should service the interrupt, and let the interrupt handler decide whether we proceed with the interrupted
//		instruction (by restoring the activity state packet and GRS), or to abandon it (by simply never 'returning')
//		It doesn't matter to us; we just let the caller handle the interrupts and manipulate our internals as it
//		wishes - eventually we'll get back here ready to do the next right thing.
//		We do not change INF or EXRF before returning to the caller.
//		Note that, if EXRF is set, we'll get here at least once after the target instruction has been placed in F0
//		but its U has not yet been developed. We don't care; we just keep cycling.
//	The instruction is complete, and it was NOT the target of an EXR instruction - INF is set, EXRF is clear.
//		It will return complete==true, interrupt==nil, and e.preventPCUpdate will be true if the instruction was
//		a successful jump or a test that skipped NI - in both of these cases, PAR.PC has already been set
//		appropriately. We don't have to worry about preserving the state of e.preventPCUpdate because it will only
//		be set by the JUMP/TEST instructions just before returning complete, and we'll increment PAR.PC (or not)
//		and clear INF before returning to the caller for potential interrupt processing.
//	The instruction is complete, and it WAS the target of an EXR instruction (INF and EXRF are set).
//		If the repeat register (R1) is zero, EXR processing is complete (even if we haven't yet executed the target
//		instruction - i.e., EXR invoked with any valid target instruction with R1==0... We don't care at this point.
//		In this case, we clear INF and EXRF, and increment PAR.PC *if* we are not prevented from doing so.
//		If R1 is NOT zero, we simply return to the caller *without* incrementing PAR.PC even if we could have done so
//		otherwise. Having INF and EXRF already set, we won't waste time re-evaluating the EXR instruction, we'll
//		just (re-)execute the target instruction.
func (e *InstructionEngine) DoCycle() {
	ikr := e.activityStatePacket.GetIndicatorKeyRegister()
	if !ikr.IsInstructionInF0() {
		e.fetchInstructionWord()
		return
	}

	complete := false
	isEXRF := ikr.IsExecuteRepeatedInstruction()
	if ikr.IsExecuteRepeatedInstruction() {
		rReg := e.GetExecOrUserRRegister(1)
		if rReg.IsZero() {
			complete = true
		}
	}

	if !complete {
		wasEXRF := isEXRF
		complete = e.executeCurrentInstruction()
		if ikr.IsExecuteRepeatedInstruction() {
			if wasEXRF {
				rReg := e.GetExecOrUserRRegister(1)
				rReg.SetW(rReg.GetW() - 1)
				complete = e.preventPCUpdate
			}
		}
	}

	if complete {
		e.SetInstructionPoint(BetweenInstructions)
		e.clearStorageLocks()
		e.cachedInstructionHandler = nil
		ikr.SetInstructionInF0(false)
		ikr.SetExecuteRepeatedInstruction(false)
		if !e.preventPCUpdate {
			e.GetProgramAddressRegister().IncrementProgramCounter()
		}
	}
}

func (e *InstructionEngine) Dump() {
	fmt.Printf("Instruction Engine Dump for %s ---------------------------------------------------------------------\n", e.name)

	if !e.pendingInterrupts.IsClear() {
		fmt.Printf("  Pending Interrupts:\n")
		e.pendingInterrupts.Dump()
	}

	var f0String string
	if e.activityStatePacket.GetIndicatorKeyRegister().IsInstructionInF0() {
		f0String = fmt.Sprintf("%012o : %s",
			e.activityStatePacket.GetCurrentInstruction(),
			dasm.DisassembleInstruction(e.activityStatePacket))
	} else {
		f0String = "invalid"
	}
	fmt.Printf("  F0: %s\n", f0String)

	par := e.activityStatePacket.GetProgramAddressRegister()
	fmt.Printf("  PAR.PC L:%o BDI:%05o PC:%06o\n",
		par.GetLevel(), par.GetBankDescriptorIndex(), par.GetProgramCounter())

	ikr := e.activityStatePacket.GetIndicatorKeyRegister()
	fmt.Printf("  Indicator Key Register: %012o\n", ikr.GetComposite())
	fmt.Printf("    Access Key:        %s\n", ikr.GetAccessKey().GetString())
	fmt.Printf("    SSF:               %03o\n", ikr.GetShortStatusField())
	fmt.Printf("    Interrupt Class:   %03o\n", ikr.GetInterruptClassField())
	fmt.Printf("    EXR Instruction:   %v\n", ikr.IsExecuteRepeatedInstruction())
	fmt.Printf("    Breakpoint Match:  %v\n", ikr.IsBreakpointRegisterMatchCondition())
	fmt.Printf("    Software Break:    %v\n", ikr.IsSoftwareBreak())
	fmt.Printf("    Instruction in F0: %v\n", ikr.IsInstructionInF0())

	dr := e.activityStatePacket.GetDesignatorRegister()
	fmt.Printf("  Designator Register: %012o\n", dr.GetComposite())
	fmt.Printf("    FHIP:                        %v\n", dr.IsFaultHandlingInProgress())
	fmt.Printf("    Executive 24-bit Indexing:   %v\n", dr.IsExecutive24BitIndexingSet())
	fmt.Printf("    Quantum Timer Enable:        %v\n", dr.IsQuantumTimerEnabled())
	fmt.Printf("    Deferrable Interrupt Enable: %v\n", dr.IsDeferrableInterruptEnabled())
	fmt.Printf("    Processor Privilege:         %v\n", dr.GetProcessorPrivilege())
	fmt.Printf("    Basic Mode:                  %v\n", dr.IsBasicModeEnabled())
	fmt.Printf("    Exec Register Set Selection: %v\n", dr.IsExecRegisterSetSelected())
	fmt.Printf("    Carry:                       %v\n", dr.IsCarrySet())
	fmt.Printf("    Overflow:                    %v\n", dr.IsOverflowSet())
	fmt.Printf("    Characteristic Underflow:    %v\n", dr.IsCharacteristicUnderflowSet())
	fmt.Printf("    Characteristic Overflow:     %v\n", dr.IsCharacteristicOverflowSet())
	fmt.Printf("    Divide Check:                %v\n", dr.IsDivideCheckSet())
	fmt.Printf("    Operation Trap Enable:       %v\n", dr.IsOperationTrapEnabled())
	fmt.Printf("    Arithmetic Exception Enable: %v\n", dr.IsArithmeticExceptionEnabled())
	fmt.Printf("    Basic Mode Base Reg Sel:     %v\n", dr.GetBasicModeBaseRegisterSelection())
	fmt.Printf("    Quarter Word Selection:      %v\n", dr.IsQuarterWordModeEnabled())

	e.generalRegisterSet.Dump()

	fmt.Printf("  Base Register Set\n")
	for bx := 0; bx < 32; bx++ {
		br := e.baseRegisters[bx]
		if !br.IsVoid() {
			bd := br.GetBankDescriptor()
			fmt.Printf("    B%-2d: addr:%s lower:%012o upper:%012o large:%v subset:%012o\n",
				bx,
				bd.GetBaseAddress().GetString(),
				bd.GetLowerLimitNormalized(),
				bd.GetUpperLimitNormalized(),
				bd.IsLargeBank(),
				br.GetSubsetting())
		}
	}

	fmt.Printf("  Active Base Table Entries\n")
	for bx := 1; bx < 16; bx++ {
		abte := e.activeBaseTable[bx]
		fmt.Printf("    %2d:%s\n", bx, abte.GetString())
	}
}

// FindBasicModeBank takes a relative address and determines which (if any) of the basic mode banks
// currently based on BDR12-15 is to be selected for that address.
// Returns the bank descriptor index (from 12 to 15) for the proper bank descriptor.
// Returns zero if the address is not within any of the based bank limits.
func (e *InstructionEngine) FindBasicModeBank(relativeAddress uint64) uint {
	db31 := e.activityStatePacket.GetDesignatorRegister().GetBasicModeBaseRegisterSelection()
	for tx := 0; tx < 4; tx++ {
		//  See IP PRM 4.4.5 - select the base register from the selection table.
		//  If the bank is void, skip it.
		//  If the program counter is outside the bank limits, skip it.
		//  Otherwise, we found the BDR we want to use.
		brIndex := baseRegisterCandidates[db31][tx]
		bReg := e.baseRegisters[brIndex]
		if e.isWithinLimits(bReg, relativeAddress) {
			return brIndex
		}
	}

	return 0
}

// GetActiveBaseTableEntry retrieves a pointer to the ABET for the indicated base register 0 to 15
func (e *InstructionEngine) GetActiveBaseTableEntry(index uint64) *ActiveBaseTableEntry {
	return e.activeBaseTable[index]
}

// GetBaseRegister retrieves a pointer to the indicated base register
func (e *InstructionEngine) GetBaseRegister(index uint64) *common.BaseRegister {
	return e.baseRegisters[index]
}

// GetConsecutiveOperands retrieves one or more word values (for double- or multi-word transfer operations).
// The assumption is that this call is made for a single iteration of an instruction.
// Per doc 9.2, effective relative address (U) will be calculated only once; however, access checks must succeed
// for all accesses.
// We presume we are retrieving from TRS of from storage - i.e., NOT allowing immediate addressing.
// Also, we presume that we are doing full-word transfers - not partial word.
// grsCheck: if true, we should check U to see if it is a GRS location
// count: number of consecutive words to be returned
// forUpdate: if true, we perform access checks for read and for write; otherwise, only for read.
// this is done for SYSC, which both reads and updates consecutive operands.
// Returns requested operands or an interrupt which should be posted.
// Returns complete == false if we are in the middle of resolving addresses.
func (e *InstructionEngine) GetConsecutiveOperands(grsCheck bool, count uint64, forUpdate bool) (result ConsecutiveOperandsResult) {
	result.complete = true

	//  Get the relative address so we can do a grsCheck
	result.sourceRelativeAddress, result.complete, result.interrupt = e.resolveRelativeAddress(false)
	if !result.complete || result.interrupt != nil {
		return
	}

	e.incrementIndexRegisterInF0()

	//  If this is a GRS reference - we do not need to look for containing banks or validate storage limits.
	asp := e.activityStatePacket
	dr := asp.GetDesignatorRegister()
	if (grsCheck) &&
		(dr.IsBasicModeEnabled() || (asp.GetCurrentInstruction().GetB() == 0)) &&
		(result.sourceRelativeAddress < 0200) {

		//  For multiple accesses, advancing beyond GRS 0177 throws a limits violation
		//  Do accessibility check for each GRS access
		grsIndex := result.sourceRelativeAddress
		for ox := uint64(0); ox < count; ox++ {
			if grsIndex == 0200 {
				result.interrupt = common.NewReferenceViolationInterrupt(common.ReferenceViolationGRS, false)
				return
			}

			if !e.isGRSAccessAllowed(grsIndex, dr.GetProcessorPrivilege(), false) {
				result.interrupt = common.NewReferenceViolationInterrupt(common.ReferenceViolationReadAccess, false)
				return
			}

			grsIndex++
		}

		result.source = e.generalRegisterSet.GetConsecutiveRegisters(grsIndex, count)
		return
	}

	//  Get base register and check storage and access limits
	result.sourceBaseRegisterIndex, result.interrupt = e.findBaseRegisterIndex(result.sourceRelativeAddress)
	if result.interrupt != nil {
		return
	}

	bReg := e.baseRegisters[result.sourceBaseRegisterIndex]
	ikr := e.activityStatePacket.GetIndicatorKeyRegister()
	result.interrupt = e.checkAccessLimitsRange(bReg, result.sourceRelativeAddress, count, true, forUpdate, ikr.GetAccessKey())
	if result.interrupt != nil {
		return
	}

	result.sourceVirtualAddress, result.sourceAbsoluteAddress, result.interrupt =
		e.translateAddress(result.sourceBaseRegisterIndex, result.sourceRelativeAddress)
	if result.interrupt != nil {
		return
	}

	result.source, result.interrupt = e.mainStorage.GetSliceFromAddress(result.sourceAbsoluteAddress, count)

	_, result.interrupt = e.checkBreakpointRange(BreakpointRead, result.sourceAbsoluteAddress, count)
	return
}

func (e *InstructionEngine) GetCurrentInstruction() *common.InstructionWord {
	return e.activityStatePacket.GetCurrentInstruction()
}

func (e *InstructionEngine) GetDesignatorRegister() *common.DesignatorRegister {
	return e.activityStatePacket.GetDesignatorRegister()
}

// GetExecOrUserARegister retrieves either the EA{index} or A{index} register
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserARegister(registerIndex uint64) *common.Word36 {
	return e.generalRegisterSet.GetRegister(e.GetExecOrUserARegisterIndex(registerIndex))
}

// GetExecOrUserARegisterIndex retrieves the GRS index of either EA{index} or A{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserARegisterIndex(registerIndex uint64) uint64 {
	if e.activityStatePacket.GetDesignatorRegister().IsExecRegisterSetSelected() {
		return common.EA0 + registerIndex
	} else {
		return common.A0 + registerIndex
	}
}

// GetExecOrUserRRegister retrieves either the ER{index} or R{index} register
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserRRegister(registerIndex uint64) *common.Word36 {
	return e.generalRegisterSet.GetRegister(e.GetExecOrUserRRegisterIndex(registerIndex))
}

// GetExecOrUserRRegisterIndex retrieves the GRS index of either ER{index} or R{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserRRegisterIndex(registerIndex uint64) uint64 {
	if e.activityStatePacket.GetDesignatorRegister().IsExecRegisterSetSelected() {
		return common.ER0 + registerIndex
	} else {
		return common.R0 + registerIndex
	}
}

// GetExecOrUserXRegister retrieves a pointer to the index register which corresponds to
// the given register index (0 to 15), and based upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserXRegister(registerIndex uint64) *common.IndexRegister {
	index := e.GetExecOrUserXRegisterIndex(registerIndex)
	return (*common.IndexRegister)(e.generalRegisterSet.GetRegister(index))
}

// GetExecOrUserXRegisterIndex retrieves the GRS index of either EX{index} or X{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserXRegisterIndex(registerIndex uint64) uint64 {
	if e.activityStatePacket.GetDesignatorRegister().IsExecRegisterSetSelected() {
		return common.EX0 + registerIndex
	} else {
		return common.X0 + registerIndex
	}
}

// GetGeneralRegisterSet retrieves a pointer to the GRS
func (e *InstructionEngine) GetGeneralRegisterSet() *common.GeneralRegisterSet {
	return e.generalRegisterSet
}

// GetImmediateOperand retrieves an operand in the case where the u (and possibly h and i) fields
// comprise the requested data.  This is NOT for jump instructions, which have slightly different rules.
// Load the value indicated in F0 as follows:
//
//	For Processor Privilege 0,1
//		value is 24 bits for DR.11 (kexec 24bit indexing enabled) true, else 18 bits
//	For Processor Privilege 2,3
//		value is 24 bits for FO.i set, else 18 bits
//
// If F0.x is zero, the immediate value is taken from the h,i, and u fields (unsigned), and negative zero is eliminated.
// For F0.x nonzero, the immediate value is the sum of the u field (unsigned) with the F0.x(mod) signed field.
//
//	For Extended Mode, with Processor Privilege 0,1 and DR.11 set, index modifiers are 24 bits;
//		otherwise, they are 18 bits.
//	For Basic Mode, index modifiers are always 18 bits.
//
// In either case, the value will be left alone for j-field=016, and sign-extended for j-field=017.
func (e *InstructionEngine) GetImmediateOperand() (operand uint64, interrupt common.Interrupt) {
	operand = 0
	interrupt = nil

	ci := e.activityStatePacket.GetCurrentInstruction()
	dr := e.activityStatePacket.GetDesignatorRegister()

	exec24Index := dr.IsExecutive24BitIndexingSet()
	privilege := dr.GetProcessorPrivilege()
	valueIs24Bits := ((privilege < 2) && exec24Index) || ((privilege > 1) && (ci.GetI() != 0))

	if ci.GetX() == 0 {
		//  No indexing (x-field is zero).  Value is derived from h, i, and u fields.
		//  Get the value from h,i,u, and eliminate negative zero.
		operand = ci.GetHIU()
		if operand == 0777777 {
			operand = 0
		}

		if (ci.GetJ() == 017) && ((operand & 0400000) != 0) {
			operand |= 0_777777_000000
		}
	} else {
		//  Value is taken only from the u field, and we eliminate negative zero at this point.
		operand = ci.GetU()
		if operand == 0177777 {
			operand = 0
		}

		//  Add the contents of Xx(m) if F0.x is non-zero
		if ci.GetX() != 0 {
			xReg := e.GetExecOrUserXRegister(ci.GetX())
			if !dr.IsBasicModeEnabled() && (privilege < 2) && exec24Index {
				operand = common.AddSimple(operand, xReg.GetXM24())
			} else {
				operand = common.AddSimple(operand, xReg.GetXM())
			}
		}

		e.incrementIndexRegisterInF0()
	}

	//  Truncate the result to the proper size, then sign-extend if appropriate to do so.
	extend := ci.GetJ() == 017
	if valueIs24Bits {
		operand &= 077_777777
		if extend && (operand&040_000000) != 0 {
			operand |= 0_777700_000000
		}
	} else {
		operand &= 0_777777
		if extend && (operand&0_400000) != 0 {
			operand |= 0_777777_000000
		}
	}

	return
}

func (e *InstructionEngine) GetInstructionPoint() InstructionPoint {
	return e.instructionPoint
}

// GetJumpHistory retrieves all cached jump history entries, and clears the history
func (e *InstructionEngine) GetJumpHistory() []common.VirtualAddress {
	return e.jumpHistory.GetEntries()
}

// GetJumpOperand is similar to getImmediateOperand()
// However the calculated U field is only ever 16 or 18 bits, and is never sign-extended.
// Also, we do not rely upon j-field for anything, as that has no meaning for conditionalJump instructions.
// in the designator register if necessary
// Returns requested operand or an interrupt which should be posted.
// Returns flip31==true if designator register bit 31 should be flipped if/when the jump is actually taken.
// Returns complete == false if we are in the middle of resolving addresses.
func (e *InstructionEngine) GetJumpOperand() (operand uint64, flip31 bool, completed bool, interrupt common.Interrupt) {
	flip31 = false
	operand, completed, interrupt = e.resolveRelativeAddress(true)
	if completed && interrupt == nil {
		dr := e.GetDesignatorRegister()
		if dr.IsBasicModeEnabled() {
			var brx uint
			brx, interrupt = e.findBaseRegisterIndexBasicMode(operand)
			// base register paris are ordered according LSB... the pair B12/B14 have the bit clear,
			// while the pair B13/B15 have the bit set. It can be shown that two base register indices with
			// matching LSBs are in the same pair, while having differing values indicates they are in opposite
			// pairs. If we can show that the current fetch base register is in the opposite pair to the
			// base register developed for the jump operand, then we can say that DB31 needs to be flipped.
			flip31 = (e.baseRegisterIndexForFetch % 01) != (brx % 01)
		}
	}

	return
}

// GetOperand implements the general case of retrieving an operand, including all forms of addressing
// and partial word access. Instructions which use the j-field as part of the function code will likely set
// allowImmediate and allowPartial false.
//
// grsDest: true if we are going to put this value into a GRS location
// grsCheck: true if we should consider GRS for addresses < 0200 for our source
// allowImm: true if we should allow immediate addressing
// allowPartial: true if we should do partial word transfers (presuming we are not in a GRS address)
// Returns requested operand or an interrupt which should be posted and a reference to the word in storage, if the
// operand came from storage.
// Returns complete == false if we are in the middle of resolving addresses.
func (e *InstructionEngine) GetOperand(
	grsDest bool,
	grsCheck bool,
	allowImm bool,
	allowPartial bool,
	lockStorage bool) (result GetOperandResult) {

	result.complete = true

	// immediate operand?
	jField := uint(e.activityStatePacket.GetCurrentInstruction().GetJ())
	if allowImm && ((jField == common.JFieldU) || (jField == common.JFieldXU)) {
		result.operand, result.interrupt = e.GetImmediateOperand()
		return
	}

	// get relative address and handle indirect addressing
	result.sourceRelativeAddress, result.complete, result.interrupt = e.resolveRelativeAddress(false)
	if !result.complete || result.interrupt != nil {
		return
	}

	asp := e.activityStatePacket
	dReg := asp.GetDesignatorRegister()
	basicMode := dReg.IsBasicModeEnabled()
	privilege := dReg.GetProcessorPrivilege()
	grs := e.generalRegisterSet

	// using kexec base registers?
	if !basicMode {
		result.sourceBaseRegisterIndex = e.getEffectiveBaseRegisterIndex()
	}

	e.incrementIndexRegisterInF0()

	//  Loading from GRS?  If so, go get the value.
	//  If grsDest is true, get the full value. Otherwise, honor j-field for partial-word transfer.
	//  (Any GRS-to-GRS transfer is full-word, regardless of j-field)
	if (grsCheck) && (basicMode || (result.sourceBaseRegisterIndex == 0)) && (result.sourceRelativeAddress < 0200) {
		//  First, do accessibility checks
		if !e.isGRSAccessAllowed(result.sourceRelativeAddress, privilege, false) {
			result.interrupt = common.NewReferenceViolationInterrupt(common.ReferenceViolationReadAccess, true)
			return
		}

		//  If we are GRS or not allowing partial word transfers, do a full word.
		//  Otherwise, honor partial word transferring.
		if grsDest || !allowPartial {
			result.operand = grs.GetRegister(result.sourceRelativeAddress).GetW()
		} else {
			qWordMode := dReg.IsQuarterWordModeEnabled()
			result.operand = common.ExtractPartialWord(grs.GetRegister(result.sourceRelativeAddress).GetW(), jField, qWordMode)
		}

		result.sourceIsGRS = true
	} else {
		//  Loading from storage.  Do so, then (maybe) honor partial word handling.
		if basicMode {
			result.sourceBaseRegisterIndex, result.interrupt = e.findBaseRegisterIndex(result.sourceRelativeAddress)
			if result.interrupt != nil {
				return
			}
		}

		bReg := e.baseRegisters[result.sourceBaseRegisterIndex]
		key := asp.GetIndicatorKeyRegister().GetAccessKey()
		result.interrupt = e.checkAccessLimitsAndAccessibility(basicMode, result.sourceBaseRegisterIndex, result.sourceRelativeAddress, false, true, false, key)
		if result.interrupt != nil {
			return
		}

		result.sourceVirtualAddress, result.sourceAbsoluteAddress, result.interrupt =
			e.translateAddress(result.sourceBaseRegisterIndex, result.sourceRelativeAddress)
		if result.interrupt != nil {
			return
		}

		if lockStorage {
			e.mainStorage.Lock(result.sourceVirtualAddress, e)
		}

		readOffset := result.sourceRelativeAddress - bReg.GetBankDescriptor().GetLowerLimitNormalized()
		result.source = &bReg.GetStorage()[readOffset]
		if allowPartial {
			qWordMode := dReg.IsQuarterWordModeEnabled()
			result.operand = common.ExtractPartialWord(result.source.GetW(), jField, qWordMode)
		} else {
			result.operand = result.source.GetW()
		}

		_, result.interrupt = e.checkBreakpoint(BreakpointRead, result.sourceAbsoluteAddress)
	}

	return
}

func (e *InstructionEngine) GetProgramAddressRegister() *common.ProgramAddressRegister {
	return e.activityStatePacket.GetProgramAddressRegister()
}

func (e *InstructionEngine) GetStopReason() (StopReason, uint64) {
	return e.stopReason, e.stopDetail.GetW()
}

func (e *InstructionEngine) GetStorageLockClientName() string {
	return e.name
}

func (e *InstructionEngine) HasPendingInterrupt() bool {
	return !e.pendingInterrupts.IsClear()
}

// IgnoreOperand is specifically for the NOP instruction.  We go through the process of developing U,
// but we do not retrieve the operand therefrom. This means that we do no access checks except in the case of checking
// for read access during indirect address resolution.
// Returns complete == false if we are in the middle of resolving addresses,
// or an interrupt if we fail some sort of limits or access checking.
func (e *InstructionEngine) IgnoreOperand() (complete bool, interrupt common.Interrupt) {
	var relAddr uint64
	relAddr, complete, interrupt = e.resolveRelativeAddress(false)
	if !complete || interrupt != nil {
		return
	}

	asp := e.activityStatePacket
	dReg := asp.GetDesignatorRegister()
	basicMode := dReg.IsBasicModeEnabled()

	var brx uint
	if !basicMode {
		brx = e.getEffectiveBaseRegisterIndex()
	}

	if relAddr > 0177 || (!basicMode && brx > 0) {
		// we do the following simply for access checks - not necessary for GRS locations
		brx, interrupt = e.findBaseRegisterIndex(relAddr)
		if interrupt != nil {
			return
		}
		key := asp.GetIndicatorKeyRegister().GetAccessKey()
		interrupt = e.checkAccessLimitsAndAccessibility(basicMode, brx, relAddr, false, false, false, key)
	}

	return
}

func (e *InstructionEngine) IsLoggingInstructions() bool {
	return e.logInstructions
}

func (e *InstructionEngine) IsLoggingInterrupts() bool {
	return e.logInterrupts
}

func (e *InstructionEngine) IsStopped() bool {
	return e.isStopped
}

// PostInterrupt posts a new interrupt, provided that no higher-priority interrupt is already pending.
// Interrupts are posted top-down, in order of priority.
// Synchronous interrupts of a lower priority than the new interrupt are discarded.
func (e *InstructionEngine) PostInterrupt(i common.Interrupt) {
	fmt.Printf("===Posting %s\n", common.GetInterruptString(i)) // TODO remove later
	e.pendingInterrupts.Post(i)
	e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
}

// SetBaseRegister sets the base register identified by brIndex (0 to 15) to the given register
func (e *InstructionEngine) SetBaseRegister(brIndex uint64, register *common.BaseRegister) {
	e.baseRegisters[brIndex] = register
}

func (e *InstructionEngine) SetExecOrUserARegister(regIndex uint64, value uint64) {
	e.generalRegisterSet.SetRegisterValue(e.GetExecOrUserARegisterIndex(regIndex), value)
}

func (e *InstructionEngine) SetExecOrUserRRegister(regIndex uint64, value uint64) {
	e.generalRegisterSet.SetRegisterValue(e.GetExecOrUserRRegisterIndex(regIndex), value)
}

func (e *InstructionEngine) SetExecOrUserXRegister(regIndex uint64, value uint64) {
	e.generalRegisterSet.SetRegisterValue(e.GetExecOrUserXRegisterIndex(regIndex), value)
}

func (e *InstructionEngine) SetLogInstructions(flag bool) {
	e.logInstructions = flag
}

func (e *InstructionEngine) SetLogInterrupts(flag bool) {
	e.logInterrupts = flag
}

func (e *InstructionEngine) SetInstructionPoint(value InstructionPoint) {
	e.instructionPoint = value
}

// SetProgramCounter sets the program counter in the PAR as well as the preventPCUpdate (aka prevent increment) flag
func (e *InstructionEngine) SetProgramCounter(counter uint64, preventIncrement bool) {
	e.activityStatePacket.GetProgramAddressRegister().SetProgramCounter(counter)
	e.preventPCUpdate = preventIncrement
}

// Stop posts a system stop, providing a reason and optionally some detail.
// This does not actually stop anything - it is up to whoever is managing the engine
// to make some sense of this and do something appropriate.
func (e *InstructionEngine) Stop(reason StopReason, detail common.Word36) {
	fmt.Printf("Stopping Processor: Reason=%d, Detail=%012o\n", reason, detail)
	e.isStopped = true
	e.stopReason = reason
	e.stopDetail = detail
}

// StoreConsecutiveOperands handles the general case of storing operands either to consecutive locations
// in storage or in the GRS.
//
// grsCheck: true if relative addresses < 0200 should be considered GRS locations
// operands: values to be stored
// grsWrap:  true indicates that, if the grs index exceeds 0177, that it will wrap to 0 instead of resulting
// in a potential reference violation.
// Returns complete == false if we are in the middle of resolving addresses, or an interrupt if one needs to be posted.
func (e *InstructionEngine) StoreConsecutiveOperands(
	grsCheck bool,
	operands []uint64) (complete bool, interrupt common.Interrupt) {

	complete = true
	interrupt = nil

	count := uint64(len(operands))
	var relAddr uint64
	relAddr, complete, interrupt = e.resolveRelativeAddress(false)
	if !complete || interrupt != nil {
		return
	}

	e.incrementIndexRegisterInF0()

	dr := e.activityStatePacket.GetDesignatorRegister()
	basicMode := dr.IsBasicModeEnabled()

	var brx uint
	if !basicMode {
		brx = e.getEffectiveBaseRegisterIndex()
	}

	if (grsCheck) && (basicMode || (brx == 0)) && (relAddr < 0200) {
		grsIndex := relAddr
		for ox := uint64(0); ox < count; ox++ {
			if grsIndex == 0200 {
				interrupt = common.NewReferenceViolationInterrupt(common.ReferenceViolationGRS, false)
				return
			}

			if !e.isGRSAccessAllowed(grsIndex, dr.GetProcessorPrivilege(), true) {
				interrupt = common.NewReferenceViolationInterrupt(common.ReferenceViolationReadAccess, false)
				return
			}

			e.generalRegisterSet.SetRegisterValue(relAddr, operands[ox])
			grsIndex++
		}
	} else {
		//  This is going to be a storage thing...
		//  Get base register and check storage and access limits
		var brx uint
		brx, interrupt = e.findBaseRegisterIndex(relAddr)
		if interrupt != nil {
			return
		}

		bReg := e.baseRegisters[brx]
		ikr := e.activityStatePacket.GetIndicatorKeyRegister()
		interrupt = e.checkAccessLimitsRange(bReg, relAddr, count, false, true, ikr.GetAccessKey())
		if interrupt != nil {
			return
		}

		var absAddr *common.AbsoluteAddress
		_, absAddr, interrupt = e.translateAddress(brx, relAddr)
		if interrupt != nil {
			return
		}

		var dest []common.Word36
		dest, interrupt = e.mainStorage.GetSliceFromAddress(absAddr, count)
		if interrupt != nil {
			return
		}

		for dx := uint64(0); dx < count; dx++ {
			dest[dx].SetW(operands[dx])
		}

		_, interrupt = e.checkBreakpointRange(BreakpointWrite, absAddr, count)
	}

	return
}

// StoreOperand handles the general case of storing an operand either to storage or to a GRS location
//
// grsSource: true if the value came from a register, so we know whether to ignore partial-word transfers
// grsCheck: true if relative addresses < 0200 should be considered GRS locations
// checkImmediate: true if we should consider j-fields 016 and 017 as immediate addressing (and throw away the operand)
// allowPartial: true if we should allow partial-word transfers (subject to GRS-GRS transfers)
// operand: value to be stored
// Returns complete == false if we are in the middle of resolving addresses, or an interrupt if one needs to be posted.
func (e *InstructionEngine) StoreOperand(
	grsSource bool,
	grsCheck bool,
	checkImmediate bool,
	allowPartial bool,
	operand uint64) (complete bool, interrupt common.Interrupt) {

	complete = true
	interrupt = nil

	//  If we allow immediate addressing mode and j-field is U or XU... we do mostly nothing.
	ci := e.activityStatePacket.GetCurrentInstruction()
	jField := uint(ci.GetJ())
	if (checkImmediate) && (jField >= 016) {
		e.incrementIndexRegisterInF0()
		return
	}

	var relAddr uint64
	relAddr, complete, interrupt = e.resolveRelativeAddress(false)
	if !complete || interrupt != nil {
		return
	}

	e.incrementIndexRegisterInF0()

	dr := e.activityStatePacket.GetDesignatorRegister()
	basicMode := dr.IsBasicModeEnabled()
	privilege := dr.GetProcessorPrivilege()

	var brx uint
	if !basicMode {
		brx = e.getEffectiveBaseRegisterIndex()
	}

	if (grsCheck) && (basicMode || (brx == 0)) && (relAddr < 0200) {
		// We're storing into the GRS... First, do accessibility checks
		if !e.isGRSAccessAllowed(relAddr, privilege, true) {
			interrupt = common.NewReferenceViolationInterrupt(common.ReferenceViolationWriteAccess, false)
			return
		}

		//  If we are GRS or not allowing partial word transfers, do a full word.
		//  Otherwise, honor partial word transfer.
		if !grsSource && allowPartial {
			qWordMode := dr.IsQuarterWordModeEnabled()
			originalValue := e.generalRegisterSet.GetRegister(relAddr).GetW()
			newValue := common.InjectPartialWord(originalValue, operand, jField, qWordMode)
			e.generalRegisterSet.GetRegister(relAddr).SetW(newValue)
		} else {
			e.generalRegisterSet.GetRegister(relAddr).SetW(operand)
		}
	} else {
		//  This is going to be a storage thing...
		if basicMode {
			brx, interrupt = e.findBaseRegisterIndexBasicMode(relAddr)
			if interrupt != nil {
				return
			}
		}

		ikr := e.activityStatePacket.GetIndicatorKeyRegister()
		interrupt = e.checkAccessLimitsAndAccessibility(basicMode, brx, relAddr, false, false, true, ikr.GetAccessKey())
		if interrupt != nil {
			return
		}

		var absAddr *common.AbsoluteAddress
		_, absAddr, interrupt = e.translateAddress(brx, relAddr)
		if interrupt != nil {
			return
		}

		var found bool
		found, interrupt = e.checkBreakpoint(BreakpointWrite, absAddr)
		if found || interrupt != nil {
			return
		}

		bReg := e.baseRegisters[brx]
		offset := relAddr - bReg.GetBankDescriptor().GetLowerLimitNormalized()
		if allowPartial {
			qWordMode := dr.IsQuarterWordModeEnabled()
			originalValue := bReg.GetStorage()[offset].GetW()
			newValue := common.InjectPartialWord(originalValue, operand, jField, qWordMode)
			bReg.GetStorage()[offset].SetW(newValue)
		} else {
			bReg.GetStorage()[offset].SetW(operand)
		}
	}

	return
}

//	Internal stuffs ----------------------------------------------------------------------------------------------------

// checkAccessibility compares the given key to the lock for this base register, and determines whether
// the requested access (fetch, read, and/or write) are allowed.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessibility(
	bReg *common.BaseRegister,
	fetchFlag bool,
	readFlag bool,
	writeFlag bool,
	accessKey *common.AccessKey) common.Interrupt {

	perms := bReg.GetEffectivePermissions(accessKey)
	if e.activityStatePacket.GetDesignatorRegister().IsBasicModeEnabled() && fetchFlag && !perms.CanEnter() {
		ssf := uint(040)
		if fetchFlag {
			ssf |= 01
		}
		return common.NewReferenceViolationInterrupt(common.ReferenceViolationReadAccess, fetchFlag)
	} else if readFlag && !perms.CanRead() {
		return common.NewReferenceViolationInterrupt(common.ReferenceViolationReadAccess, fetchFlag)
	} else if writeFlag && !perms.CanWrite() {
		return common.NewReferenceViolationInterrupt(common.ReferenceViolationWriteAccess, fetchFlag)
	}

	return nil
}

// checkAccessLimitsAndAccessibility checks the accessibility of a given relative address in the bank described by this
// base register for the given flags, using the given key.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessLimitsAndAccessibility(
	basicMode bool,
	baseRegisterIndex uint,
	relativeAddress uint64,
	fetchFlag bool,
	readFlag bool,
	writeFlag bool,
	accessKey *common.AccessKey) common.Interrupt {

	bReg := e.baseRegisters[baseRegisterIndex]
	i := e.checkAccessLimitsForAddress(basicMode, baseRegisterIndex, relativeAddress, fetchFlag)
	if i != nil {
		return i
	}

	i = e.checkAccessibility(bReg, fetchFlag, readFlag, writeFlag, accessKey)
	return i
}

// checkAccessLimitsForAddress checks whether the relative address is within the limits of the bank
// described by this base register. We only need the fetch flag for posting an interrupt.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessLimitsForAddress(
	basicMode bool,
	baseRegisterIndex uint,
	relativeAddress uint64,
	fetchFlag bool) (interrupt common.Interrupt) {

	interrupt = nil
	if fetchFlag && relativeAddress < 0200 {
		if basicMode || baseRegisterIndex == 0 {
			interrupt = common.NewReferenceViolationInterrupt(common.ReferenceViolationStorageLimits, true)
			return
		}
	}

	bReg := e.baseRegisters[baseRegisterIndex]
	if bReg.IsVoid() {
		return common.NewReferenceViolationInterrupt(common.ReferenceViolationStorageLimits, fetchFlag)
	}

	bDesc := bReg.GetBankDescriptor()
	if (relativeAddress < bDesc.GetLowerLimitNormalized()) ||
		(relativeAddress > bDesc.GetUpperLimitNormalized()) {
		return common.NewReferenceViolationInterrupt(common.ReferenceViolationStorageLimits, fetchFlag)
	}

	return nil
}

// checkAccessLimitsRange checks the access limits for a consecutive range of addresses, starting at the given
// relativeAddress, for the number of addresses. Checks for read and/or write access according to the values given
// for readFlag and writeFlag. Uses the given access key for the determination.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessLimitsRange(
	bReg *common.BaseRegister,
	relativeAddress uint64,
	addressCount uint64,
	readFlag bool,
	writeFlag bool,
	accessKey *common.AccessKey) common.Interrupt {

	bDesc := bReg.GetBankDescriptor()
	if (relativeAddress < bDesc.GetLowerLimitNormalized()) ||
		((relativeAddress + addressCount - 1) > bDesc.GetUpperLimitNormalized()) {
		return common.NewReferenceViolationInterrupt(common.ReferenceViolationStorageLimits, false)
	}

	return e.checkAccessibility(bReg, false, readFlag, writeFlag, accessKey)
}

func (e *InstructionEngine) checkBreakpoint(
	comparison BreakpointComparison,
	absAddr *common.AbsoluteAddress) (found bool, interrupt common.Interrupt) {

	found = false
	interrupt = nil

	if e.breakpointAddress != nil && e.breakpointAddress.Equals(absAddr) {
		if (comparison == BreakpointFetch && e.breakpointFetch) ||
			(comparison == BreakpointRead && e.breakpointRead) ||
			(comparison == BreakpointWrite && e.breakpointWrite) {

			e.activityStatePacket.GetIndicatorKeyRegister().SetBreakpointRegisterMatchCondition(true)
			found = true
			if e.breakpointHalt {
				e.Stop(BreakpointStop, 0)
			} else {
				interrupt = common.NewBreakpointInterrupt()
			}
		}
	}

	return
}

func (e *InstructionEngine) checkBreakpointRange(
	comparison BreakpointComparison,
	absAddr *common.AbsoluteAddress,
	count uint64) (found bool, interrupt common.Interrupt) {

	found = false
	interrupt = nil

	if e.breakpointAddress != nil && e.breakpointAddress.GetSegment() == absAddr.GetSegment() {
		if (comparison == BreakpointFetch && e.breakpointFetch) ||
			(comparison == BreakpointRead && e.breakpointRead) ||
			(comparison == BreakpointWrite && e.breakpointWrite) {

			brkOffset := e.breakpointAddress.GetOffset()
			addrOffset := absAddr.GetOffset()
			for x := uint64(0); x < count; x++ {
				if addrOffset == brkOffset {
					e.activityStatePacket.GetIndicatorKeyRegister().SetBreakpointRegisterMatchCondition(true)
					found = true
					if e.breakpointHalt {
						e.Stop(BreakpointStop, 0)
					} else {
						interrupt = common.NewBreakpointInterrupt()
					}
					return
				}
				addrOffset++
			}
		}
	}

	return
}

func (e *InstructionEngine) clearStorageLocks() {
	e.mainStorage.ReleaseAllLocks(e)
}

func (e *InstructionEngine) createJumpHistoryEntry(address common.VirtualAddress) {
	interrupt := e.jumpHistory.StoreEntry(address)
	if interrupt != nil {
		e.PostInterrupt(interrupt)
	}
}

// executeCurrentInstruction executes the instruction in F0 (which we cache to save some cycles)
// Functions return true if the have completed normally. They will return false in the following conditions:
//
//	Interrupt Mid-point
//	Still resolving addressing (indirect addressing always does this)
//	An interrupt was posted - this is so that the instruction can be retried for certain interrupts.
//
// Returns true if the instruction was complete, else false
func (e *InstructionEngine) executeCurrentInstruction() (completed bool) {
	if e.logInstructions {
		code := dasm.DisassembleInstruction(e.activityStatePacket)
		fmt.Printf("--[%012o  %s]\n", e.activityStatePacket.GetProgramAddressRegister().GetComposite(), code)
	}

	dr := e.activityStatePacket.GetDesignatorRegister()
	ci := e.activityStatePacket.GetCurrentInstruction()
	e.preventPCUpdate = false

	// Find the instruction handler for the instruction if it is not cached
	if e.cachedInstructionHandler == nil {
		fTable := FunctionTable[dr.IsBasicModeEnabled()]
		var found bool
		e.cachedInstructionHandler, found = fTable[uint(ci.GetF())]
		if !found {
			// illegal instruction - post an interrupt, then note that we are between instructions.
			e.PostInterrupt(common.NewInvalidInstructionInterrupt(common.InvalidInstructionBadFunctionCode))
			e.SetInstructionPoint(BetweenInstructions)
			return false
		}
	}

	return e.cachedInstructionHandler(e)
}

// fetchInstructionWord retrieves the next instruction word from the appropriate bank.
// For extended mode this is straight-forward.
// For basic mode, we have to hunt around a bit to make sure we pull it from the most appropriate bank.
// If something bad happens, an interrupt is posted and we return false
func (e *InstructionEngine) fetchInstructionWord() bool {
	basicMode := e.activityStatePacket.GetDesignatorRegister().IsBasicModeEnabled()
	programCounter := e.activityStatePacket.GetProgramAddressRegister().GetProgramCounter()

	var bReg *common.BaseRegister
	var brx uint
	if basicMode {
		// If we don't know the index of the current basic mode instruction bank,
		// find it and set DB31 accordingly.
		if e.baseRegisterIndexForFetch == 0 {
			brx = e.FindBasicModeBank(programCounter)
			if brx == 0 {
				interrupt := common.NewReferenceViolationInterrupt(common.ReferenceViolationStorageLimits, false)
				e.PostInterrupt(interrupt)
				return false
			}

			e.baseRegisterIndexForFetch = brx
			e.GetDesignatorRegister().SetBasicModeBaseRegisterSelection(brx == 13 || brx == 15)
		}

		bReg = e.baseRegisters[e.baseRegisterIndexForFetch]
		if !e.isReadAllowed(bReg) {
			interrupt := common.NewReferenceViolationInterrupt(common.ReferenceViolationStorageLimits, false)
			e.PostInterrupt(interrupt)
			return false
		}
	} else {
		brx = 0
		bReg = e.baseRegisters[0]
		ikr := e.activityStatePacket.GetIndicatorKeyRegister()
		interrupt := e.checkAccessLimitsAndAccessibility(basicMode, 0, programCounter, true, false, false, ikr.GetAccessKey())
		if interrupt != nil {
			e.PostInterrupt(interrupt)
			return false
		}
	}

	if bReg.IsVoid() || bReg.GetBankDescriptor().IsLargeBank() {
		interrupt := common.NewReferenceViolationInterrupt(common.ReferenceViolationStorageLimits, false)
		e.PostInterrupt(interrupt)
		return false
	}

	pcOffset := programCounter - bReg.GetBankDescriptor().GetLowerLimitNormalized()
	iw := common.InstructionWord(bReg.GetStorage()[pcOffset])
	asp := e.activityStatePacket
	asp.SetCurrentInstruction(&iw)
	asp.GetIndicatorKeyRegister().SetInstructionInF0(true)
	asp.GetIndicatorKeyRegister().SetExecuteRepeatedInstruction(false)

	_, absAddr, _ := e.translateAddress(brx, programCounter)
	e.checkBreakpoint(BreakpointFetch, absAddr)

	return true
}

// findBankDescriptor retrieves a struct to describe the given named bank.
//
//	This is for interrupt handling.
//	The bank name is in L,BDI format.
//	bankLevel level of the bank, 0:7
//	bankDescriptorIndex BDI of the bank 0:077777
func (e *InstructionEngine) findBankDescriptor(bankLevel uint64, bankDescriptorIndex uint64) (*common.BankDescriptor, bool) {
	// The bank descriptor tables for bank levels 0 through 7 are described by the banks based on B16 through B23.
	// The bank descriptor will be the {n}th bank descriptor in the particular bank descriptor table,
	// where {n} is the bank descriptor index.
	bdRegIndex := bankLevel + 16
	if e.baseRegisters[bdRegIndex].IsVoid() {
		e.PostInterrupt(common.NewAddressingExceptionInterrupt(common.AddressingExceptionFatal, bankLevel, bankDescriptorIndex))
		return nil, false
	}

	//  bdStorage contains the BDT for the given bank_name level
	//  bdTableOffset indicates the offset into the BDT, where the bank descriptor is to be found.
	bdStorage := e.baseRegisters[bdRegIndex].GetStorage()
	bdTableOffset := bankDescriptorIndex + 8
	if bdTableOffset+8 > uint64(len(bdStorage)) {
		e.PostInterrupt(common.NewAddressingExceptionInterrupt(common.AddressingExceptionFatal, bankLevel, bankDescriptorIndex))
		return nil, false
	}

	//  Create and return a BankDescriptor object
	bd := common.NewBankDescriptorFromStorage(bdStorage[bdTableOffset : bdTableOffset+8])
	return bd, true
}

// findBaseRegisterIndexBasicMode locates the index of the base register which represents the bank which contains
// the given relative address for basic mode. Does appropriate limits checking.
// relAddr: relative address to be considered (basic mode only)
func (e *InstructionEngine) findBaseRegisterIndexBasicMode(relAddr uint64) (baseRegisterIndex uint, interrupt common.Interrupt) {
	interrupt = nil

	baseRegisterIndex = e.FindBasicModeBank(relAddr)
	if baseRegisterIndex == 0 {
		interrupt = common.NewReferenceViolationInterrupt(common.ReferenceViolationStorageLimits, false)
	}

	return
}

// findBaseRegisterIndex checks the execution mode and returns the base register which should be used for any
// operand fetch or store, given the current F0 and other processor states.
// relAddr: only relevant for basic mode
func (e *InstructionEngine) findBaseRegisterIndex(relAddr uint64) (baseRegisterIndex uint, interrupt common.Interrupt) {
	baseRegisterIndex = 0
	interrupt = nil

	if e.GetDesignatorRegister().IsBasicModeEnabled() {
		baseRegisterIndex, interrupt = e.findBaseRegisterIndexBasicMode(relAddr)
	} else {
		baseRegisterIndex = e.getEffectiveBaseRegisterIndex()
	}

	return
}

func (e *InstructionEngine) flipDesignatorRegisterBit31() {
	dr := e.GetDesignatorRegister()
	dr.SetBasicModeBaseRegisterSelection(!dr.GetBasicModeBaseRegisterSelection())
}

func (e *InstructionEngine) getCurrentVirtualAddress() common.VirtualAddress {
	dr := e.GetDesignatorRegister()
	if dr.IsBasicModeEnabled() {
		brx := e.baseRegisterIndexForFetch
		if brx == 0 {
			brx, _ = e.findBaseRegisterIndexBasicMode(e.GetProgramAddressRegister().GetProgramCounter())
		}

		abte := e.activeBaseTable[brx]
		return common.TranslateToBasicMode(abte.bankLevel, abte.bankDescriptorIndex, abte.subsetSpecification)
	} else {
		brx := e.getEffectiveBaseRegisterIndex() & 017
		abte := e.activeBaseTable[brx]
		return common.NewExtendedModeVirtualAddress(abte.bankLevel, abte.bankDescriptorIndex, abte.subsetSpecification)
	}
}

// getEffectiveBaseRegisterIndex determines (for extended mode only) the effective base register to use,
// based on the current processor privilege and the values in F0.B and possibly F0.I.
func (e *InstructionEngine) getEffectiveBaseRegisterIndex() uint {
	if e.activityStatePacket.GetDesignatorRegister().GetProcessorPrivilege() < 2 {
		return uint(e.activityStatePacket.GetCurrentInstruction().GetIB())
	} else {
		return uint(e.activityStatePacket.GetCurrentInstruction().GetB())
	}
}

// incrementIndexRegisterInF0 checks the instruction and current modes to determine whether register Xx
// should be incremented, and if so it performs the appropriate incrementation.
func (e *InstructionEngine) incrementIndexRegisterInF0() {
	ci := e.GetCurrentInstruction()
	if ci.GetX() > 0 && ci.GetH() > 0 {
		xReg := e.GetExecOrUserXRegister(ci.GetX())
		dr := e.GetDesignatorRegister()
		if !dr.IsBasicModeEnabled() && (dr.GetProcessorPrivilege() < 2) && dr.IsExecutive24BitIndexingSet() {
			xReg.IncrementModifier24()
		} else {
			xReg.IncrementModifier()
		}
	}
}

func (e *InstructionEngine) isGRSAccessAllowed(registerIndex uint64, processorPrivilege uint64, writeAccess bool) bool {
	if registerIndex < 040 {
		return true
	} else if registerIndex < 0100 {
		return false
	} else if registerIndex < 0120 {
		return true
	} else {
		return (writeAccess && (processorPrivilege == 0)) || (!writeAccess && (processorPrivilege <= 2))
	}
}

func (e *InstructionEngine) isReadAllowed(bReg *common.BaseRegister) bool {
	permissions := bReg.GetEffectivePermissions(e.activityStatePacket.GetIndicatorKeyRegister().GetAccessKey())
	return permissions.CanRead()
}

// isWithinLimits evaluates the given offset within the constraints of the given base register,
// returning true if the offset is within those constraints, else false
func (e *InstructionEngine) isWithinLimits(bReg *common.BaseRegister, offset uint64) bool {
	return !bReg.IsVoid() &&
		(offset >= bReg.GetBankDescriptor().GetLowerLimitNormalized()) &&
		(offset <= bReg.GetBankDescriptor().GetUpperLimitNormalized())
}

// resolveRelativeAddress reads the instruction in F0, and in conjunction with the current ASP environment,
// develops the relative address as a function of the unsigned 16-bit U or the 12-bit D field,
// added with the signed modifier portion of the index register indicated by F0.x (presuming that field is not zero).
// useU: for Extended Mode Jump instructions which use the entire U (or HIU) fields for the relative address.
// Basic mode always uses the u field.
// If we handle an iteration of indirect addressing, we return with complete == false
// If, during an iteration of indirect addressing we hit an access or limits check we return with interrupt == the
// appropriate interrupt.
// Otherwise, we return with complete==true, interrupt == nil, and the relative address in relAddr
func (e *InstructionEngine) resolveRelativeAddress(useU bool) (relAddr uint64, complete bool, interrupt common.Interrupt) {
	relAddr = 0
	complete = false
	interrupt = nil

	e.SetInstructionPoint(ResolvingAddress)

	ci := e.activityStatePacket.GetCurrentInstruction()
	dr := e.activityStatePacket.GetDesignatorRegister()

	var base uint64
	if dr.IsBasicModeEnabled() || useU {
		base = ci.GetU()
	} else {
		base = ci.GetD()
	}

	x := ci.GetX()
	var addend uint64
	if x != 0 {
		xReg := e.GetExecOrUserXRegister(x)
		if dr.IsExecutive24BitIndexingSet() && dr.GetProcessorPrivilege() < 2 {
			addend = xReg.GetSignedXM24()
		} else {
			addend = xReg.GetSignedXM()
		}
	}

	relAddr = common.AddSimple(base, addend)

	basicMode := dr.IsBasicModeEnabled()
	if ci.GetI() != 0 && basicMode && dr.GetProcessorPrivilege() > 1 {
		//	indirect addressing specified
		var brx uint
		brx, interrupt = e.findBaseRegisterIndex(relAddr)
		if interrupt != nil {
			return
		}

		key := e.activityStatePacket.GetIndicatorKeyRegister().GetAccessKey()
		interrupt = e.checkAccessLimitsAndAccessibility(basicMode, brx, relAddr, true, false, false, key)
		if interrupt != nil {
			return
		}

		var absAddr *common.AbsoluteAddress
		_, absAddr, interrupt = e.translateAddress(brx, relAddr)
		if interrupt != nil {
			return
		}

		var found bool
		found, interrupt = e.checkBreakpoint(BreakpointRead, absAddr)
		if found || interrupt != nil {
			// TODO the processor has been stopped - immediately stop
			return
		}

		var word *common.Word36
		word, interrupt = e.mainStorage.GetWordFromAddress(absAddr)
		if interrupt != nil {
			return
		}

		e.GetCurrentInstruction().SetXHIU(word.GetW())
		return
	}

	e.SetInstructionPoint(MidInstruction)
	complete = true
	return
}

func (e *InstructionEngine) translateAddress(
	baseRegisterIndex uint,
	relativeAddress uint64) (virAddr common.VirtualAddress, absAddr *common.AbsoluteAddress, interrupt common.Interrupt) {

	var level uint64
	var bdi uint64
	var offset uint64

	if baseRegisterIndex == 0 {
		par := e.GetProgramAddressRegister()
		level = par.GetLevel()
		bdi = par.GetBankDescriptorIndex()
		offset = 0
	} else {
		abte := e.activeBaseTable[baseRegisterIndex]
		level = abte.bankLevel
		bdi = abte.bankDescriptorIndex
		offset = abte.subsetSpecification
	}

	bReg := e.baseRegisters[baseRegisterIndex]
	if bReg.IsVoid() {
		interrupt = common.NewAddressingExceptionInterrupt(common.AddressingExceptionFatal, level, bdi)
		return
	}

	bDesc := bReg.GetBankDescriptor()
	offset += relativeAddress - bDesc.GetLowerLimitNormalized()
	if bDesc.GetBankType() == common.BasicModeBankDescriptor {
		virAddr = common.TranslateToBasicMode(level, bdi, offset)
	} else if bDesc.GetBankType() == common.ExtendedModeBankDescriptor {
		virAddr = common.NewExtendedModeVirtualAddress(level, bdi, offset)
	} else {
		interrupt = common.NewAddressingExceptionInterrupt(common.AddressingExceptionFatal, level, bdi)
		return
	}

	offset += bDesc.GetBaseAddress().GetOffset()
	absAddr = common.NewAbsoluteAddress(bDesc.GetBaseAddress().GetSegment(), offset)

	return
}
