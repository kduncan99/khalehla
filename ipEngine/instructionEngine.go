// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"khalehla/dasm"
	"khalehla/pkg"
)

type BreakpointComparison uint

const (
	BreakpointFetch BreakpointComparison = 1
	BreakpointRead  BreakpointComparison = 2
	BreakpointWrite BreakpointComparison = 3
)

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

const L0BDTBaseRegister = 16
const ICSBaseRegister = 26
const ICSIndexRegister = EX1
const RCSBaseRegister = 25
const RCSIndexRegister = EX0

// InstructionEngine implements the basic functionality required to execute 36-bit code.
// It does not handle any actual hardware considerations such as interrupts, etc.
// It does track stop reasons, as there are certain generic processes that require this ability
// whether we are emulating hardware or an operating system with hardware.
// It must be wrapped by additional code which does this, either as an IP emulator or an OS emulator.
type InstructionEngine struct {
	name         string           // unique name of this engine - must be set externally
	mainStorage  *pkg.MainStorage // must be set externally
	storageLocks *StorageLocks    // must be set externally

	activeBaseTable           [16]*ActiveBaseTableEntry // [0] is unused
	activityStatePacket       *pkg.ActivityStatePacket
	baseRegisters             [32]*pkg.BaseRegister
	baseRegisterIndexForFetch uint // only applies to basic mode - if 0, it is not valid; otherwise it is 12:15
	generalRegisterSet        *GeneralRegisterSet

	//	If not nil, describes an interrupt which needs to be handled as soon as possible
	pendingInterrupts *InterruptStack
	jumpHistory       *JumpHistory

	logInstructions bool
	logInterrupts   bool

	//	If true, the current (or most recent) instruction has set the PAR.PC the way it wants,
	//	and we should not increment it for the next instruction
	preventPCUpdate bool

	breakpointAddress *pkg.AbsoluteAddress
	breakpointHalt    bool
	breakpointFetch   bool
	breakpointRead    bool
	breakpointWrite   bool

	isStopped        bool
	stopReason       StopReason
	stopDetail       pkg.Word36
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

//	external stuffs ----------------------------------------------------------------------------------------------------

func NewEngine(name string, mainStorage *pkg.MainStorage, storageLocks *StorageLocks) *InstructionEngine {
	e := &InstructionEngine{}
	e.name = name
	e.mainStorage = mainStorage
	e.storageLocks = storageLocks
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
		e.baseRegisters[bx] = pkg.NewVoidBaseRegister()
	}
	e.baseRegisterIndexForFetch = 0

	e.generalRegisterSet = NewGeneralRegisterSet()
	e.activityStatePacket = pkg.NewActivityStatePacket()
	e.breakpointAddress = nil
	e.breakpointHalt = false
	e.breakpointFetch = false
	e.breakpointRead = false
	e.breakpointWrite = false

	e.isStopped = true
	e.stopReason = NotStopped
	e.stopDetail = 0

	e.preventPCUpdate = false
	e.instructionPoint = BetweenInstructions
}

func (e *InstructionEngine) ClearAllInterrupts() {
	e.pendingInterrupts.Clear()
}

func (e *InstructionEngine) ClearStop() {
	e.isStopped = false
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
			fmt.Printf("    B%-2d: addr:%08o.%012o lower:%012o upper:%012o large:%v\n",
				bx,
				br.GetBaseAddress().GetSegment(),
				br.GetBaseAddress().GetOffset(),
				br.GetLowerLimitNormalized(),
				br.GetUpperLimitNormalized(),
				br.IsLargeSize())
		}
	}

	fmt.Printf("  Storage Locks:\n")
	e.storageLocks.Dump()
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
func (e *InstructionEngine) GetBaseRegister(index uint64) *pkg.BaseRegister {
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
func (e *InstructionEngine) GetConsecutiveOperands(grsCheck bool, count uint64, forUpdate bool) (operands []pkg.Word36, complete bool, interrupt pkg.Interrupt) {
	operands = nil
	complete = true
	interrupt = nil

	//  Get the relative address so we can do a grsCheck
	var relAddr uint64
	relAddr, complete, interrupt = e.resolveRelativeAddress(false)
	if !complete || interrupt != nil {
		return
	}

	e.incrementIndexRegisterInF0()

	//  If this is a GRS reference - we do not need to look for containing banks or validate storage limits.
	asp := e.activityStatePacket
	dr := asp.GetDesignatorRegister()
	if (grsCheck) &&
		(dr.IsBasicModeEnabled() || (asp.GetCurrentInstruction().GetB() == 0)) &&
		(relAddr < 0200) {

		//  For multiple accesses, advancing beyond GRS 0177 throws a limits violation
		//  Do accessibility check for each GRS access
		grsIndex := relAddr
		for ox := uint64(0); ox < count; ox++ {
			if grsIndex == 0200 {
				interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationGRS, false)
				return
			}

			if !e.isGRSAccessAllowed(grsIndex, dr.GetProcessorPrivilege(), false) {
				interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationReadAccess, false)
				return
			}

			grsIndex++
		}

		operands = e.generalRegisterSet.GetConsecutiveRegisters(grsIndex, count)
		return
	}

	//  Get base register and check storage and access limits
	var brIndex uint
	brIndex, interrupt = e.findBaseRegisterIndex(relAddr)
	if interrupt != nil {
		return
	}

	bReg := e.baseRegisters[brIndex]
	ikr := e.activityStatePacket.GetIndicatorKeyRegister()
	interrupt = e.checkAccessLimitsRange(bReg, relAddr, count, true, forUpdate, ikr.GetAccessKey())
	if interrupt != nil {
		return
	}

	var found bool
	absAddr := bReg.ConvertRelativeAddress(relAddr)
	found, interrupt = e.checkBreakpointRange(BreakpointRead, absAddr, count)
	if found {
		return
	}

	operands, _ = e.mainStorage.GetSlice(absAddr.GetSegment(), absAddr.GetOffset(), count)
	return
}

func (e *InstructionEngine) GetCurrentInstruction() *pkg.InstructionWord {
	return e.activityStatePacket.GetCurrentInstruction()
}

func (e *InstructionEngine) GetDesignatorRegister() *pkg.DesignatorRegister {
	return e.activityStatePacket.GetDesignatorRegister()
}

// GetExecOrUserARegister retrieves either the EA{index} or A{index} register
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserARegister(registerIndex uint64) *pkg.Word36 {
	return e.generalRegisterSet.GetRegister(e.GetExecOrUserARegisterIndex(registerIndex))
}

// GetExecOrUserARegisterIndex retrieves the GRS index of either EA{index} or A{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserARegisterIndex(registerIndex uint64) uint64 {
	if e.activityStatePacket.GetDesignatorRegister().IsExecRegisterSetSelected() {
		return EA0 + registerIndex
	} else {
		return A0 + registerIndex
	}
}

// GetExecOrUserRRegister retrieves either the ER{index} or R{index} register
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserRRegister(registerIndex uint64) *pkg.Word36 {
	return e.generalRegisterSet.GetRegister(e.GetExecOrUserRRegisterIndex(registerIndex))
}

// GetExecOrUserRRegisterIndex retrieves the GRS index of either ER{index} or R{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserRRegisterIndex(registerIndex uint64) uint64 {
	if e.activityStatePacket.GetDesignatorRegister().IsExecRegisterSetSelected() {
		return ER0 + registerIndex
	} else {
		return R0 + registerIndex
	}
}

// GetExecOrUserXRegister retrieves a pointer to the index register which corresponds to
// the given register index (0 to 15), and based upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserXRegister(registerIndex uint64) *IndexRegister {
	index := e.GetExecOrUserXRegisterIndex(registerIndex)
	return (*IndexRegister)(e.generalRegisterSet.GetRegister(index))
}

// GetExecOrUserXRegisterIndex retrieves the GRS index of either EX{index} or X{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserXRegisterIndex(registerIndex uint64) uint64 {
	if e.activityStatePacket.GetDesignatorRegister().IsExecRegisterSetSelected() {
		return EX0 + registerIndex
	} else {
		return X0 + registerIndex
	}
}

// GetGeneralRegisterSet retrieves a pointer to the GRS
func (e *InstructionEngine) GetGeneralRegisterSet() *GeneralRegisterSet {
	return e.generalRegisterSet
}

// GetImmediateOperand retrieves an operand in the case where the u (and possibly h and i) fields
// comprise the requested data.  This is NOT for jump instructions, which have slightly different rules.
// Load the value indicated in F0 as follows:
//
//	For Processor Privilege 0,1
//		value is 24 bits for DR.11 (exec 24bit indexing enabled) true, else 18 bits
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
func (e *InstructionEngine) GetImmediateOperand() (operand uint64, interrupt pkg.Interrupt) {
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
				operand = pkg.AddSimple(operand, xReg.GetXM24())
			} else {
				operand = pkg.AddSimple(operand, xReg.GetXM())
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
func (e *InstructionEngine) GetJumpHistory() []pkg.VirtualAddress {
	return e.jumpHistory.GetEntries()
}

// GetJumpOperand is similar to getImmediateOperand()
// However the calculated U field is only ever 16 or 18 bits, and is never sign-extended.
// Also, we do not rely upon j-field for anything, as that has no meaning for conditionalJump instructions.
// in the designator register if necessary
// Returns requested operand or an interrupt which should be posted.
// Returns flip31==true if designator register bit 31 should be flipped if/when the jump is actually taken.
// Returns complete == false if we are in the middle of resolving addresses.
func (e *InstructionEngine) GetJumpOperand() (operand uint64, flip31 bool, completed bool, interrupt pkg.Interrupt) {
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
// Returns requested operand or an interrupt which should be posted.
// Returns complete == false if we are in the middle of resolving addresses.
func (e *InstructionEngine) GetOperand(
	grsDest bool,
	grsCheck bool,
	allowImm bool,
	allowPartial bool) (operand uint64, completed bool, interrupt pkg.Interrupt) {

	operand = 0
	completed = true
	interrupt = nil

	// immediate operand?
	jField := uint(e.activityStatePacket.GetCurrentInstruction().GetJ())
	if allowImm && ((jField == pkg.JFieldU) || (jField == pkg.JFieldXU)) {
		operand, interrupt = e.GetImmediateOperand()
		return
	}

	// get relative address and handle indirect addressing
	var relAddr uint64
	relAddr, completed, interrupt = e.resolveRelativeAddress(false)
	if !completed || interrupt != nil {
		return
	}

	asp := e.activityStatePacket
	ci := asp.GetCurrentInstruction()
	dReg := asp.GetDesignatorRegister()
	basicMode := dReg.IsBasicModeEnabled()
	privilege := dReg.GetProcessorPrivilege()
	grs := e.generalRegisterSet

	// using exec base registers?
	var brx uint
	if !basicMode {
		brx = uint(ci.GetB())
		if (privilege < 2) && (ci.GetI() != 0) {
			brx += 16
		}
	}

	e.incrementIndexRegisterInF0()

	//  Loading from GRS?  If so, go get the value.
	//  If grsDest is true, get the full value. Otherwise, honor j-field for partial-word transfer.
	//  (Any GRS-to-GRS transfer is full-word, regardless of j-field)
	if (grsCheck) && (basicMode || (brx == 0)) && (relAddr < 0200) {
		//  First, do accessibility checks
		if e.isGRSAccessAllowed(relAddr, privilege, false) {
			interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationReadAccess, true)
			return
		}

		//  If we are GRS or not allowing partial word transfers, do a full word.
		//  Otherwise, honor partial word transferring.
		if grsDest || !allowPartial {
			operand = grs.GetRegister(relAddr).GetW()
		} else {
			qWordMode := dReg.IsQuarterWordModeEnabled()
			operand = pkg.ExtractPartialWord(grs.GetRegister(relAddr).GetW(), jField, qWordMode)
		}
	} else {
		//  Loading from storage.  Do so, then (maybe) honor partial word handling.
		if basicMode {
			brx, interrupt = e.findBaseRegisterIndex(relAddr)
			if interrupt != nil {
				return
			}
		}

		bReg := e.baseRegisters[brx]
		key := asp.GetIndicatorKeyRegister().GetAccessKey()
		interrupt = e.checkAccessLimitsAndAccessibility(basicMode, brx, relAddr, false, true, false, key)
		if interrupt != nil {
			return
		}

		var absAddress pkg.AbsoluteAddress
		bReg.PopulateAbsoluteAddress(relAddr, &absAddress)
		e.checkBreakpoint(BreakpointRead, &absAddress)

		readOffset := relAddr - bReg.GetLowerLimitNormalized()
		operand = bReg.GetStorage()[readOffset].GetW()
		if allowPartial {
			qWordMode := dReg.IsQuarterWordModeEnabled()
			operand = pkg.ExtractPartialWord(operand, jField, qWordMode)
		}
	}

	return
}

func (e *InstructionEngine) GetProgramAddressRegister() *pkg.ProgramAddressRegister {
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
func (e *InstructionEngine) IgnoreOperand() (complete bool, interrupt pkg.Interrupt) {
	var relAddr uint64
	relAddr, complete, interrupt = e.resolveRelativeAddress(false)
	if !complete || interrupt != nil {
		return
	}

	asp := e.activityStatePacket
	ci := asp.GetCurrentInstruction()
	dReg := asp.GetDesignatorRegister()
	basicMode := dReg.IsBasicModeEnabled()
	privilege := dReg.GetProcessorPrivilege()

	var brx uint
	if !basicMode {
		brx = uint(ci.GetB())
		if (privilege < 2) && (ci.GetI() != 0) {
			brx += 16
		}
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
func (e *InstructionEngine) PostInterrupt(i pkg.Interrupt) {
	fmt.Printf("===Posting %s\n", pkg.GetInterruptString(i)) // TODO remove
	e.pendingInterrupts.Post(i)
	e.createJumpHistoryEntry(e.getCurrentVirtualAddress())
}

// SetBaseRegister sets the base register identified by brIndex (0 to 15) to the given register
func (e *InstructionEngine) SetBaseRegister(brIndex uint64, register *pkg.BaseRegister) {
	e.baseRegisters[brIndex] = register
}

func (e *InstructionEngine) SetExecOrUserARegister(regIndex uint64, value pkg.Word36) {
	e.generalRegisterSet.SetRegisterValue(e.GetExecOrUserARegisterIndex(regIndex), value)
}

func (e *InstructionEngine) SetExecOrUserRRegister(regIndex uint64, value pkg.Word36) {
	e.generalRegisterSet.SetRegisterValue(e.GetExecOrUserRRegisterIndex(regIndex), value)
}

func (e *InstructionEngine) SetExecOrUserXRegister(regIndex uint64, value IndexRegister) {
	e.generalRegisterSet.SetRegisterValue(e.GetExecOrUserXRegisterIndex(regIndex), pkg.Word36(value))
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
func (e *InstructionEngine) Stop(reason StopReason, detail pkg.Word36) {
	fmt.Printf("Stopping Processor: Reason=%d, Detail=%012o\n", reason, detail)
	e.isStopped = true
	e.stopReason = reason
	e.stopDetail = detail
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
	operand uint64) (complete bool, interrupt pkg.Interrupt) {

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
		brx = uint(ci.GetB())
		if (privilege < 2) && (ci.GetI() != 0) {
			brx += 16
		}
	}

	if (grsCheck) && (basicMode || (brx == 0)) && (relAddr < 0200) {
		// We're storing into the GRS... First, do accessibility checks
		if !e.isGRSAccessAllowed(relAddr, privilege, true) {
			interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationWriteAccess, false)
			return
		}

		//  If we are GRS or not allowing partial word transfers, do a full word.
		//  Otherwise, honor partial word transfer.
		if !grsSource && allowPartial {
			qWordMode := dr.IsQuarterWordModeEnabled()
			originalValue := e.generalRegisterSet.GetRegister(relAddr).GetW()
			newValue := pkg.InjectPartialWord(originalValue, operand, jField, qWordMode)
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

		bReg := e.baseRegisters[brx]
		ikr := e.activityStatePacket.GetIndicatorKeyRegister()
		interrupt = e.checkAccessLimitsAndAccessibility(basicMode, brx, relAddr, false, false, true, ikr.GetAccessKey())
		if interrupt != nil {
			return
		}

		var absAddr pkg.AbsoluteAddress
		bReg.PopulateAbsoluteAddress(relAddr, &absAddr)
		e.checkBreakpoint(BreakpointWrite, &absAddr)

		offset := relAddr - bReg.GetLowerLimitNormalized()
		if allowPartial {
			qWordMode := dr.IsQuarterWordModeEnabled()
			originalValue := bReg.GetStorage()[offset].GetW()
			newValue := pkg.InjectPartialWord(originalValue, operand, jField, qWordMode)
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
	bReg *pkg.BaseRegister,
	fetchFlag bool,
	readFlag bool,
	writeFlag bool,
	accessKey *pkg.AccessKey) pkg.Interrupt {

	perms := bReg.GetEffectivePermissions(accessKey)
	if e.activityStatePacket.GetDesignatorRegister().IsBasicModeEnabled() && fetchFlag && !perms.CanEnter() {
		ssf := uint(040)
		if fetchFlag {
			ssf |= 01
		}
		return pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationReadAccess, fetchFlag)
	} else if readFlag && !perms.CanRead() {
		return pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationReadAccess, fetchFlag)
	} else if writeFlag && !perms.CanWrite() {
		return pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationWriteAccess, fetchFlag)
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
	accessKey *pkg.AccessKey) pkg.Interrupt {

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
	fetchFlag bool) (interrupt pkg.Interrupt) {

	interrupt = nil
	if fetchFlag && relativeAddress < 0200 {
		if basicMode || baseRegisterIndex == 0 {
			interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, true)
			return
		}
	}

	bReg := e.baseRegisters[baseRegisterIndex]
	if (relativeAddress < bReg.GetLowerLimitNormalized()) ||
		(relativeAddress > bReg.GetUpperLimitNormalized()) {
		return pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, fetchFlag)
	}

	return nil
}

// checkAccessLimitsRange checks the access limits for a consecutive range of addresses, starting at the given
// relativeAddress, for the number of addresses. Checks for read and/or write access according to the values given
// for readFlag and writeFlag. Uses the given access key for the determination.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessLimitsRange(
	bReg *pkg.BaseRegister,
	relativeAddress uint64,
	addressCount uint64,
	readFlag bool,
	writeFlag bool,
	accessKey *pkg.AccessKey) pkg.Interrupt {

	if (relativeAddress < bReg.GetLowerLimitNormalized()) ||
		((relativeAddress + addressCount - 1) > bReg.GetUpperLimitNormalized()) {
		return pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, false)
	}

	return e.checkAccessibility(bReg, false, readFlag, writeFlag, accessKey)
}

func (e *InstructionEngine) checkBreakpoint(comparison BreakpointComparison, absAddr *pkg.AbsoluteAddress) (found bool, interrupt pkg.Interrupt) {
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
				interrupt = pkg.NewBreakpointInterrupt()
			}
		}
	}

	return
}

func (e *InstructionEngine) checkBreakpointRange(
	comparison BreakpointComparison,
	absAddr *pkg.AbsoluteAddress,
	count uint64) (found bool, interrupt pkg.Interrupt) {

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
						interrupt = pkg.NewBreakpointInterrupt()
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
	e.storageLocks.ReleaseAll(e)
}

func (e *InstructionEngine) createJumpHistoryEntry(address pkg.VirtualAddress) {
	interrupt := e.jumpHistory.StoreEntry(address)
	if interrupt != nil {
		e.PostInterrupt(interrupt)
	}
}

// doCycle executes one cycle
// caller should disposition any pending interrupts before invoking this...
// since the engine is not specifically hardware (could be an executor for a native mode OS),
// we don't actually know how to handle the interrupts.
func (e *InstructionEngine) doCycle() {
	if !e.activityStatePacket.GetIndicatorKeyRegister().IsInstructionInF0() {
		// If we manage to fetch an instruction, then we transition to resolving address...
		// this is so that invoking code can handle any interrupts which can be handled at any
		// interrupt point.
		// If we do NOT fetch an instruction, then we are between instructions.
		// There is an interrupt posted ('cause we failed), and the invoker needs to know
		// that it can handle any interrupts which can be handled between instructions
		// (which is essentially all of them).
		if e.fetchInstructionWord() {
			e.SetInstructionPoint(ResolvingAddress)
		} else {
			e.SetInstructionPoint(BetweenInstructions)
		}
	} else {
		// There is already an instruction in F0.
		// We might be still resolving its address, or maybe we've not started that yet
		// (although the flag for resolving addresses would still be set), or maybe
		// we have finished resolving the address, and had been interrupted during an
		// interrupt-able instruction (such as EXR), in which case midInstruction flag is set.
		// We don't care - we just execute the instruction as it is stored in F0.
		e.executeCurrentInstruction()
	}
}

func (e *InstructionEngine) executeCurrentInstruction() {
	//	functions return true if they have completed (normally, or by posting an interrupt).
	//	They return false if they return before completion, but without posting an interrupt.
	//	Generally, a false return results either from a repeated execution instructionType giving
	//	up the ipEngine, or by an indirect basic mode instructionType giving up the ipEngine
	//	before completely developing the operand address.

	if e.logInstructions {
		code := dasm.DisassembleInstruction(e.activityStatePacket)
		fmt.Printf("--[%012o  %s]\n", e.activityStatePacket.GetProgramAddressRegister().GetComposite(), code)
	}

	dr := e.activityStatePacket.GetDesignatorRegister()
	ci := e.activityStatePacket.GetCurrentInstruction()

	e.preventPCUpdate = false

	// Find the instruction handler for the instruction
	fTable := FunctionTable[dr.IsBasicModeEnabled()]
	inst, found := fTable[uint(ci.GetF())]
	if !found {
		// illegal instruction - post an interrupt, then note that we are between instructions.
		e.PostInterrupt(pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode))
		e.SetInstructionPoint(BetweenInstructions)
		return
	}

	// Execute the instruction - if it throws an interrupt, don't change any of the instruction flags...
	// the interrupt handler will want to know where we were in the process of handling the instruction.
	completed, interrupt := inst(e)
	if interrupt != nil {
		e.PostInterrupt(interrupt)
		return
	}

	// If the instruction did not complete, then we are either still resolving the address or we are at a
	// mid-instruction point for something like EXR. Just return to the invoker so it can check for interrupts.
	if completed {
		e.SetInstructionPoint(BetweenInstructions)
		e.clearStorageLocks()
		e.activityStatePacket.GetIndicatorKeyRegister().SetInstructionInF0(false)
		if !e.preventPCUpdate {
			e.activityStatePacket.GetProgramAddressRegister().IncrementProgramCounter()
		}
	}
}

// fetchInstructionWord retrieves the next instruction word from the appropriate bank.
// For extended mode this is straight-forward.
// For basic mode, we have to hunt around a bit to make sure we pull it from the most appropriate bank.
// If something bad happens, an interrupt is posted and we return false
func (e *InstructionEngine) fetchInstructionWord() bool {
	basicMode := e.activityStatePacket.GetDesignatorRegister().IsBasicModeEnabled()
	programCounter := e.activityStatePacket.GetProgramAddressRegister().GetProgramCounter()

	var bReg *pkg.BaseRegister
	if basicMode {
		// If we don't know the index of the current basic mode instruction bank,
		// find it and set DB31 accordingly.
		if e.baseRegisterIndexForFetch == 0 {
			brx := e.FindBasicModeBank(programCounter)
			if brx == 0 {
				e.PostInterrupt(pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, false))
				return false
			}

			e.baseRegisterIndexForFetch = brx
			e.GetDesignatorRegister().SetBasicModeBaseRegisterSelection(brx == 13 || brx == 15)
		}

		bReg = e.baseRegisters[e.baseRegisterIndexForFetch]
		if !e.isReadAllowed(bReg) {
			e.PostInterrupt(pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, false))
			return false
		}
	} else {
		bReg = e.baseRegisters[0]
		ikr := e.activityStatePacket.GetIndicatorKeyRegister()
		intp := e.checkAccessLimitsAndAccessibility(basicMode, 0, programCounter, true, false, false, ikr.GetAccessKey())
		if intp != nil {
			e.PostInterrupt(intp)
			return false
		}
	}

	if bReg.IsVoid() || bReg.IsLargeSize() {
		e.PostInterrupt(pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, false))
		return false
	}

	pcOffset := programCounter - bReg.GetLowerLimitNormalized()
	a := pkg.InstructionWord(bReg.GetStorage()[pcOffset])
	e.activityStatePacket.SetCurrentInstruction(&a)
	e.activityStatePacket.GetIndicatorKeyRegister().SetInstructionInF0(true)

	return true
}

// findBankDescriptor retrieves a struct to describe the given named bank.
//
//	This is for interrupt handling.
//	The bank name is in L,BDI format.
//	bankLevel level of the bank, 0:7
//	bankDescriptorIndex BDI of the bank 0:077777
func (e *InstructionEngine) findBankDescriptor(bankLevel uint64, bankDescriptorIndex uint64) (*pkg.BankDescriptor, bool) {
	// The bank descriptor tables for bank levels 0 through 7 are described by the banks based on B16 through B23.
	// The bank descriptor will be the {n}th bank descriptor in the particular bank descriptor table,
	// where {n} is the bank descriptor index.
	bdRegIndex := bankLevel + 16
	if e.baseRegisters[bdRegIndex].IsVoid() {
		e.PostInterrupt(pkg.NewAddressingExceptionInterrupt(pkg.AddressingExceptionFatal, bankLevel, bankDescriptorIndex))
		return nil, false
	}

	//  bdStorage contains the BDT for the given bank_name level
	//  bdTableOffset indicates the offset into the BDT, where the bank descriptor is to be found.
	bdStorage := e.baseRegisters[bdRegIndex].GetStorage()
	bdTableOffset := bankDescriptorIndex + 8
	if bdTableOffset+8 > uint64(len(bdStorage)) {
		e.PostInterrupt(pkg.NewAddressingExceptionInterrupt(pkg.AddressingExceptionFatal, bankLevel, bankDescriptorIndex))
		return nil, false
	}

	//  Create and return a BankDescriptor object
	bd := pkg.NewBankDescriptorFromStorage(bdStorage[bdTableOffset : bdTableOffset+8])
	return bd, true
}

// findBaseRegisterIndexBasicMode locates the index of the base register which represents the bank which contains
// the given relative address for basic mode. Does appropriate limits checking.
// relAddr: relative address to be considered (basic mode only)
func (e *InstructionEngine) findBaseRegisterIndexBasicMode(relAddr uint64) (baseRegisterIndex uint, interrupt pkg.Interrupt) {
	interrupt = nil

	baseRegisterIndex = e.FindBasicModeBank(relAddr)
	if baseRegisterIndex == 0 {
		interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, false)
	}

	return
}

// findBaseRegisterIndex checks the execution mode and returns the base register which should be used for any
// operand fetch or store, given the current F0 and other processor states.
// relAddr: only relevant for basic mode
func (e *InstructionEngine) findBaseRegisterIndex(relAddr uint64) (baseRegisterIndex uint, interrupt pkg.Interrupt) {
	baseRegisterIndex = 0
	interrupt = nil

	if e.GetDesignatorRegister().IsBasicModeEnabled() {
		baseRegisterIndex, interrupt = e.findBaseRegisterIndexBasicMode(relAddr)
	} else {
		baseRegisterIndex = e.getEffectiveBaseRegisterIndex()
	}

	return
}

func (e *InstructionEngine) getCurrentVirtualAddress() pkg.VirtualAddress {
	dr := e.GetDesignatorRegister()
	if dr.IsBasicModeEnabled() {
		brx := e.baseRegisterIndexForFetch
		if brx == 0 {
			brx, _ = e.findBaseRegisterIndexBasicMode(e.GetProgramAddressRegister().GetProgramCounter())
		}

		abte := e.activeBaseTable[brx]
		return pkg.TranslateToBasicMode(abte.bankLevel, abte.bankDescriptorIndex, abte.subsetSpecification)
	} else {
		brx := e.getEffectiveBaseRegisterIndex()
		abte := e.activeBaseTable[brx]
		return pkg.NewExtendedModeVirtualAddress(abte.bankLevel, abte.bankDescriptorIndex, abte.subsetSpecification)
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

func (e *InstructionEngine) isReadAllowed(bReg *pkg.BaseRegister) bool {
	permissions := bReg.GetEffectivePermissions(e.activityStatePacket.GetIndicatorKeyRegister().GetAccessKey())
	return permissions.CanRead()
}

// isWithinLimits evaluates the given offset within the constraints of the given base register,
// returning true if the offset is within those constraints, else false
func (e *InstructionEngine) isWithinLimits(bReg *pkg.BaseRegister, offset uint64) bool {
	return !bReg.IsVoid() &&
		(offset >= bReg.GetLowerLimitNormalized()) &&
		(offset <= bReg.GetUpperLimitNormalized())
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
func (e *InstructionEngine) resolveRelativeAddress(useU bool) (relAddr uint64, complete bool, interrupt pkg.Interrupt) {
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

	relAddr = pkg.AddSimple(base, addend)

	basicMode := dr.IsBasicModeEnabled()
	if ci.GetI() != 0 && basicMode && dr.GetProcessorPrivilege() > 1 {
		//	indirect addressing specified
		var brx uint
		brx, interrupt = e.findBaseRegisterIndex(relAddr)
		if interrupt != nil {
			return
		}

		bReg := e.baseRegisters[brx]
		key := e.activityStatePacket.GetIndicatorKeyRegister().GetAccessKey()
		interrupt = e.checkAccessLimitsAndAccessibility(basicMode, brx, relAddr, true, false, false, key)
		if interrupt != nil {
			return
		}

		var absAddr pkg.AbsoluteAddress
		bReg.PopulateAbsoluteAddress(relAddr, &absAddr)

		var found bool
		found, interrupt = e.checkBreakpoint(BreakpointRead, &absAddr)
		if found {
			// the processor has been stopped - immediately stop
			return
		}

		var word *pkg.Word36
		word, interrupt = e.mainStorage.GetWordFromAddress(&absAddr)
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
