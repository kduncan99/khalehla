// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
	"sync"
)

type BreakpointComparison uint

const (
	BreakpointNone  BreakpointComparison = 0
	BreakpointFetch BreakpointComparison = 1 << 0
	BreakpointRead  BreakpointComparison = 1 << 1
	BreakpointWrite BreakpointComparison = 1 << 2
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

const JumpHistoryTableThreshold = 120 //	Raise interrupt when this many entries exist
const JumpHistoryTableSize = 128      //	Size of the table

var storageLocks = map[uint64]*InstructionEngine{}
var storageLocksMutex sync.Mutex

// InstructionEngine implements the basic functionality required to execute 36-bit code.
// It does not handle any actual hardware considerations such as interrupts, etc.
// It does track stop reasons, as there are certain generic processes that require this ability
// whether we are emulating hardware or an operating system with hardware.
// It must be wrapped by additional code which does this, either as an IP emulator or an OS emulator.
type InstructionEngine struct {
	mainStorage *pkg.MainStorage // must be set externally

	activeBaseTable     [16]*ActiveBaseTableEntry //	[0] is unused
	activityStatePacket ActivityStatePacket
	baseRegisters       [32]*BaseRegister
	generalRegisterSet  GeneralRegisterSet

	//	If not nil, describes an interrupt which needs to be handled as soon as possible
	pendingInterrupt Interrupt

	//	See 2.4.2
	//	Should this be saved off somewhere during an interrupt?
	jumpHistory                 [JumpHistoryTableThreshold]pkg.Word36
	jumpHistoryIndex            int
	jumpHistoryThresholdReached bool

	//	If true, the current (or most recent) instructionType has set the PAR.PC the way it wants,
	//	and we should not increment it for the next instruction
	preventPCUpdate bool

	mutex sync.Mutex

	breakpointAddress  AbsoluteAddress
	breakpointRegister BreakpointComparison

	stopReason StopReason
	stopDetail pkg.Word36
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

func NewEngine(mainStorage *pkg.MainStorage) *InstructionEngine {
	e := &InstructionEngine{}
	e.mainStorage = mainStorage
	for ax := 0; ax < 16; ax++ {
		e.activeBaseTable[ax] = NewActiveBaseTableEntryFromComposite(0)
	}
	for bx := 0; bx < 32; bx++ {
		e.baseRegisters[bx] = NewVoidBaseRegister()
	}

	e.breakpointRegister = BreakpointNone
	e.stopReason = NotStopped

	return e
}

// FindBasicModeBank takes a relative address and determines which (if any) of the basic mode banks
// currently based on BDR12-15 is to be selected for that address.
// Set updatedDB31 true if you want the code to update designator register Bit31 in the event we cross
// primary/secondary bank pairs.
// Returns the bank descriptor index (from 12 to 15) for the proper bank descriptor.
// Returns zero if the address is not within any of the based bank limits.
func (e *InstructionEngine) FindBasicModeBank(relativeAddress uint64, updateDB31 bool) uint {
	db31 := e.activityStatePacket.designatorRegister.BasicModeBaseRegisterSelection
	for tx := 0; tx < 4; tx++ {
		//  See IP PRM 4.4.5 - select the base register from the selection table.
		//  If the bank is void, skip it.
		//  If the program counter is outside the bank limits, skip it.
		//  Otherwise, we found the BDR we want to use.
		brIndex := baseRegisterCandidates[db31][tx]
		bReg := e.baseRegisters[brIndex]
		if e.isWithinLimits(bReg, relativeAddress) {
			if updateDB31 && (tx >= 2) {
				//  address is found in a secondary bank, so we need to flip DB31
				e.activityStatePacket.designatorRegister.BasicModeBaseRegisterSelection = !db31
			}

			return brIndex
		}
	}

	return 0
}

// GetActiveBaseTableEntry retrieves a pointer to the ABET for the indicated base register 0 to 15
func (e *InstructionEngine) GetActiveBaseTableEntry(index uint) *ActiveBaseTableEntry {
	return e.activeBaseTable[index]
}

// GetBaseRegister retrieves a pointer to the indicated base register
func (e *InstructionEngine) GetBaseRegister(index uint) *BaseRegister {
	return e.baseRegisters[index]
}

func (e *InstructionEngine) GetCurrentInstruction() InstructionWord {
	return e.activityStatePacket.currentInstruction
}

// GetExecOrUserARegister retrieves either the EA{index} or A{index} register
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserARegister(registerIndex uint) *pkg.Word36 {
	return e.generalRegisterSet.GetRegister(e.GetExecOrUserARegisterIndex(registerIndex))
}

// GetExecOrUserARegisterIndex retrieves the GRS index of either EA{index} or A{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserARegisterIndex(registerIndex uint) uint {
	if e.activityStatePacket.designatorRegister.ExecRegisterSetSelected {
		return EA0 + registerIndex
	} else {
		return A0 + registerIndex
	}
}

// GetExecOrUserRRegister retrieves either the ER{index} or R{index} register
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserRRegister(registerIndex uint) *pkg.Word36 {
	return e.generalRegisterSet.GetRegister(e.GetExecOrUserRRegisterIndex(registerIndex))
}

// GetExecOrUserRRegisterIndex retrieves the GRS index of either ER{index} or R{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserRRegisterIndex(registerIndex uint) uint {
	if e.activityStatePacket.designatorRegister.ExecRegisterSetSelected {
		return ER0 + registerIndex
	} else {
		return R0 + registerIndex
	}
}

// GetExecOrUserXRegister retrieves a pointer to the index register which corresponds to
// the given register index (0 to 15), and based upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserXRegister(registerIndex uint) *IndexRegister {
	index := e.GetExecOrUserXRegisterIndex(registerIndex)
	return (*IndexRegister)(e.generalRegisterSet.GetRegister(index))
}

// GetExecOrUserXRegisterIndex retrieves the GRS index of either EX{index} or X{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserXRegisterIndex(registerIndex uint) uint {
	if e.activityStatePacket.designatorRegister.ExecRegisterSetSelected {
		return EX0 + registerIndex
	} else {
		return X0 + registerIndex
	}
}

// GetGeneralRegisterSet retrieves a pointer to the GRS
func (e *InstructionEngine) GetGeneralRegisterSet() *GeneralRegisterSet {
	return &e.generalRegisterSet
}

// GetJumpOperand is similar to getImmediateOperand()
// However the calculated U field is only ever 16 or 18 bits, and is never sign-extended.
// Also, we do not rely upon j-field for anything, as that has no meaning for conditionalJump instructions.
//
// updateDesignatorRegister:
//
//	if true and if we are in basic mode, we update the basic mode bank selection bit
//	in the designator register if necessary
//
// returns complete==true and jump operand value if complete and successful
// returns Interrupt if an interrupt needs to be raised - caller should post the interrupt
// returns complete==false if an address is not fully resolved (basic mode indirect address only)
func (e *InstructionEngine) GetJumpOperand(updateDesignatorRegister bool) (complete bool, operand uint64, interrupt Interrupt) {
	complete = true
	interrupt = nil
	operand = e.calculateRelativeAddressForJump()

	//  The following bit is how we deal with indirect addressing for basic mode.
	//  If we are doing that, it will update the U portion of the current instruction with new address information,
	//  then throw UnresolvedAddressException which will eventually route us back through here again, but this
	//  time with new address info (in reladdress), and we keep doing this until we're not doing indirect addressing.
	asp := e.activityStatePacket
	if asp.designatorRegister.BasicModeEnabled && (asp.currentInstruction.GetI() != 0) {
		complete, _, interrupt = e.findBaseRegisterIndex(operand, updateDesignatorRegister)
	} else {
		e.incrementIndexRegisterInF0()
	}

	return
}

// GetOperand implements the general case of retrieving an operand, including all forms of addressing
// and partial word access. Instructions which use the j-field as part of the function code will likely set
// allowImmediate and allowPartial false.
//
// grsDestination: true if we are going to put this value into a GRS location
// grsCheck: true if we should consider GRS for addresses < 0200 for our source
// allowImmediate: true if we should allow immediate addressing
// allowPartial: true if we should do partial word transfers (presuming we are not in a GRS address)
//
// returns complete == true and the operand value is successful
// returns complete == false if address resolution is unfinished (such as can happen in Basic Mode with
// indirect addressing). In this case, caller should call back again after checking for any pending interrupts.
// returns an Interrupt if any interrupt needs to be raised - caller should post the interrupt
func (e *InstructionEngine) GetOperand(grsDestination bool, grsCheck bool, allowImmediate bool, allowPartial bool) (complete bool, operand uint64, interrupt Interrupt) {
	complete = false
	operand = 0
	interrupt = nil

	jField := uint(e.activityStatePacket.currentInstruction.GetJ())
	if allowImmediate {
		//  j-field is U or XU? If so, get the value from the instruction itself (immediate addressing)
		if jField >= 016 {
			operand, interrupt = e.getImmediateOperand()
			complete = true
			return
		}
	}

	relAddress := e.calculateRelativeAddressForGRSOrStorage()

	asp := e.activityStatePacket
	ci := asp.currentInstruction
	dReg := asp.designatorRegister
	basicMode := dReg.BasicModeEnabled
	pPriv := dReg.processorPrivilege
	grs := e.generalRegisterSet

	var baseRegisterIndex uint
	if !basicMode {
		baseRegisterIndex = uint(ci.GetB())
		if (pPriv < 2) && (ci.GetI() != 0) {
			baseRegisterIndex += 16
		}
	}

	//  Loading from GRS?  If so, go get the value.
	//  If grsDestination is true, get the full value. Otherwise, honor j-field for partial-word transfer.
	//  See hardware guide section 4.3.2 - any GRS-to-GRS transfer is full-word, regardless of j-field.
	if (grsCheck) && (basicMode || (baseRegisterIndex == 0)) && (relAddress < 0200) {
		e.incrementIndexRegisterInF0()

		//  First, do accessibility checks
		if e.isGRSAccessAllowed(relAddress, pPriv, false) {
			interrupt = NewReferenceViolationInterrupt(ReferenceViolationReadAccess, true)
			return
		}

		//  If we are GRS or not allowing partial word transfers, do a full word.
		//  Otherwise, honor partial word transferring.
		if grsDestination || !allowPartial {
			operand = grs.GetRegister(uint(relAddress)).GetW()
		} else {
			qWordMode := dReg.QuarterWordModeEnabled
			operand = e.extractPartialWord(grs.GetRegister(uint(relAddress)).GetW(), jField, qWordMode)
			return
		}
	}

	//  Loading from storage.  Do so, then (maybe) honor partial word handling.
	if basicMode {
		complete, baseRegisterIndex, interrupt = e.findBaseRegisterIndex(relAddress, false)
		if !complete || interrupt != nil {
			return
		}
	}

	baseRegister := e.baseRegisters[baseRegisterIndex]
	interrupt = e.checkAccessLimits(baseRegister, relAddress, false, true, false, asp.indicatorKeyRegister.accessKey)
	if interrupt != nil {
		return
	}

	e.incrementIndexRegisterInF0()

	var absAddress AbsoluteAddress
	e.getAbsoluteAddress(baseRegister, relAddress, &absAddress)
	e.checkBreakpoint(BreakpointRead, &absAddress)

	readOffset := relAddress - uint64(baseRegister.lowerLimitNormalized)
	operand = baseRegister.storage[readOffset].GetW()
	if allowPartial {
		qWordMode := dReg.QuarterWordModeEnabled
		operand = e.extractPartialWord(operand, jField, qWordMode)
	}

	complete = true
	return
}

// PostInterrupt posts a new interrupt, provided that no higher-priority interrupt is already pending.
func (e *InstructionEngine) PostInterrupt(i Interrupt) {
	if e.pendingInterrupt != nil {
		if i.GetClass() < e.pendingInterrupt.GetClass() {
			e.pendingInterrupt = i
		}
	}
}

// SetBaseRegister sets the base register identified by brIndex (0 to 15) to the given register
func (e *InstructionEngine) SetBaseRegister(brIndex uint, register *BaseRegister) {
	e.baseRegisters[brIndex] = register
}

func (e *InstructionEngine) SetExecOrUserARegister(regIndex uint, value pkg.Word36) {
	e.generalRegisterSet.SetRegisterValue(e.GetExecOrUserARegisterIndex(regIndex), value)
}

func (e *InstructionEngine) SetExecOrUserRRegister(regIndex uint, value pkg.Word36) {
	e.generalRegisterSet.SetRegisterValue(e.GetExecOrUserRRegisterIndex(regIndex), value)
}

func (e *InstructionEngine) SetExecOrUserXRegister(regIndex uint, value IndexRegister) {
	e.generalRegisterSet.SetRegisterValue(e.GetExecOrUserXRegisterIndex(regIndex), pkg.Word36(value))
}

// SetProgramCounter sets the program counter in the PAR as well as the preventPCUpdate (aka prevent increment) flag
func (e *InstructionEngine) SetProgramCounter(counter uint, preventIncrement bool) {
	e.activityStatePacket.programAddressRegister.SetProgramCounter(counter)
	e.preventPCUpdate = preventIncrement
}

// Stop posts a system stop, providing a reason and optionally some detail.
// This does not actually stop anything - it is up to whoever is managing the engine
// to make some sense of this and do something appropriate.
func (e *InstructionEngine) Stop(reason StopReason, detail pkg.Word36) {
	e.stopReason = reason
	e.stopDetail = detail
}

//	Internal stuffs ----------------------------------------------------------------------------------------------------

// calculateRelativeAddressForGRSOrStorage calculates the raw relative address (the U) for the current instruction.
// Does NOT increment any x registers, even if their content contributes to the result.
// Returns the relative address for the current instruction
func (e *InstructionEngine) calculateRelativeAddressForGRSOrStorage() uint64 {
	ci := e.activityStatePacket.currentInstruction
	dr := e.activityStatePacket.designatorRegister

	var xReg *IndexRegister
	xx := uint(ci.GetX())
	if xx != 0 {
		xReg = e.GetExecOrUserXRegister(xx)
	}

	var addend1 uint64
	var addend2 uint64
	if dr.BasicModeEnabled {
		addend1 = ci.GetU()
		if xReg != nil {
			addend2 = xReg.GetSignedXM()
		}
	} else {
		addend1 = ci.GetD()
		if xReg != nil {
			if dr.Executive24BitIndexingEnabled && dr.processorPrivilege < 2 {
				//  Exec 24-bit indexing is requested
				addend2 = xReg.GetSignedXM24()
			} else {
				addend2 = xReg.GetSignedXM()
			}
		}
	}

	return pkg.AddSimple(addend1, addend2)
}

// calculateRelativeAddressForJump calculates the raw relative address (the U) for the current instruction.
// Does NOT increment any x registers, even if their content contributes to the result.
// returns the relative address for the current instruction
func (e *InstructionEngine) calculateRelativeAddressForJump() uint64 {
	ci := e.activityStatePacket.currentInstruction
	dr := e.activityStatePacket.designatorRegister

	var xReg *IndexRegister
	xx := uint(ci.GetX())
	if xx != 0 {
		xReg = e.GetExecOrUserXRegister(xx)
	}

	addend1 := ci.GetU()
	var addend2 uint64
	if dr.BasicModeEnabled {
		if xReg != nil {
			addend2 = xReg.GetSignedXM()
		}
	} else {
		addend1 = ci.GetU()
		if xReg != nil {
			if dr.Executive24BitIndexingEnabled && dr.processorPrivilege < 2 {
				//  Exec 24-bit indexing is requested
				addend2 = xReg.GetSignedXM24()
			} else {
				addend2 = xReg.GetSignedXM()
			}
		}
	}

	return pkg.AddSimple(addend1, addend2)
}

// checkAccessibility compares the given key to the lock for this base register, and determines whether
// the requested access (fetch, read, and/or write) are allowed.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessibility(
	bReg *BaseRegister,
	fetchFlag bool,
	readFlag bool,
	writeFlag bool,
	accessKey *AccessKey) Interrupt {
	perms := bReg.GetEffectivePermissions(accessKey)
	if e.activityStatePacket.designatorRegister.BasicModeEnabled && fetchFlag && !perms.CanEnter() {
		ssf := uint(040)
		if fetchFlag {
			ssf |= 01
		}
		return NewReferenceViolationInterrupt(ReferenceViolationReadAccess, fetchFlag)
	} else if readFlag && !perms.CanRead() {
		return NewReferenceViolationInterrupt(ReferenceViolationReadAccess, fetchFlag)
	} else if writeFlag && !perms.CanWrite() {
		return NewReferenceViolationInterrupt(ReferenceViolationWriteAccess, fetchFlag)
	}

	return nil
}

// checkAccessLimits checks the accessibility of a given relative address in the bank described by this
// base register for the given flags, using the given key.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessLimits(
	bReg *BaseRegister,
	relativeAddress uint64,
	fetchFlag bool,
	readFlag bool,
	writeFlag bool,
	accessKey *AccessKey) Interrupt {

	i := e.checkAccessLimitsForAddress(bReg, relativeAddress, fetchFlag)
	if i != nil {
		return i
	}

	i = e.checkAccessibility(bReg, fetchFlag, readFlag, writeFlag, accessKey)
	return i
}

// checkAccessLimitsForAddress checks whether the relative address is within the limits of the bank
// described by this base register. We only need the fetch flag for posting an interrupt.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessLimitsForAddress(bReg *BaseRegister, relativeAddress uint64, fetchFlag bool) Interrupt {

	// TODO if we try to execute something in GRS - we take ReferenceViolationInterruptClass, 01, 0, 0

	if (relativeAddress < uint64(bReg.lowerLimitNormalized)) ||
		(relativeAddress > uint64(bReg.upperLimitNormalized)) {
		return NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, fetchFlag)
	}

	return nil
}

// checkAccessLimitsRange checks the access limits for a consecutive range of addresses, starting at the given
// relativeAddress, for the number of addresses. Checks for read and/or write access according to the values given
// for readFlag and writeFlag. Uses the given access key for the determination.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessLimitsRange(
	bReg *BaseRegister,
	relativeAddress uint,
	addressCount uint,
	readFlag bool,
	writeFlag bool,
	accessKey *AccessKey) Interrupt {
	if (relativeAddress < bReg.lowerLimitNormalized) ||
		(relativeAddress+addressCount-1 > bReg.upperLimitNormalized) {
		return NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, false)
	}

	return e.checkAccessibility(bReg, false, readFlag, writeFlag, accessKey)
}

func (e *InstructionEngine) checkBreakpoint(comp BreakpointComparison, addr *AbsoluteAddress) {
	//  TODO Per doc, 2.4.1.2 Breakpoint_Register - we need to halt if Halt Enable is set
	//      which means Stop Right Now... how do we do that for all callers of this code?
	if e.breakpointAddress.Equals(addr) {
		if comp&e.breakpointRegister != 0 {
			e.activityStatePacket.indicatorKeyRegister.breakpointRegisterMatchCondition = true
		}
	}
}

func (e *InstructionEngine) clearStorageLocks() {
	storageLocksMutex.Lock()
	for key, value := range storageLocks {
		if value == e {
			delete(storageLocks, key)
		}
	}
	storageLocksMutex.Unlock()
}

// createJumpHistoryTableENtry puts a new entry into the jump history table.
//
//	If we cross the interrupt threshold, set the threshold-reached flag
func (e *InstructionEngine) createJumpHistoryTableEntry(absoluteAddress pkg.Word36) {
	e.jumpHistory[e.jumpHistoryIndex] = absoluteAddress

	if e.jumpHistoryIndex > JumpHistoryTableThreshold {
		e.jumpHistoryThresholdReached = true
	}

	e.jumpHistoryIndex++
	if e.jumpHistoryIndex == JumpHistoryTableSize {
		e.jumpHistoryIndex = 0
	}
}

// doCycle executes one cycle (which may or may not correspond to one instruction
// do not invoke if there is a non-deferrable interrupt pending
func (e *InstructionEngine) doCycle() {
	//	Are we continuing an interrupted instructionType?
	// See 5.1.3
	//	INF EXRF Action on User Return from interrupt
	//   0   0   Fetch and execute the instructionType addressed by PAR.
	//   1   0   Obtain the instructionType from F0 (rather than using PAR).
	//   1   1   EXR mid-execution. Enter normal EXR logic at the point where the target instructionType has
	//             just been fetched (but not decoded), using F0 as the target instructionType.
	// Note: In the special case where EXR is itself the target of an EX instructionType, mid-execution state will have
	// EXRF clear until the first interrupt point after the EXR instructionType has been fetched.

	// TODO have we handled all the cases above?
	if e.activityStatePacket.indicatorKeyRegister.instructionInF0 {
		e.executeCurrentInstruction()
		return
	}

	//	Fetch an instructionType and execute it
	if e.fetchInstructionWord() {
		e.executeCurrentInstruction()
	}

	return
}

func (e *InstructionEngine) executeCurrentInstruction() {
	//	functions return true if they have completed (normally, or by posting an interrupt).
	//	They return false if they return before completion, but without posting an interrupt.
	//	Generally, a false return results either from a repeated execution instructionType giving
	//	up the ipEngine, or by an indirect basic mode instructionType giving up the ipEngine
	//	before completely developing the operand address.
	dr := e.activityStatePacket.designatorRegister
	ci := e.activityStatePacket.currentInstruction

	e.preventPCUpdate = false
	fTable := FunctionTable[dr.BasicModeEnabled]
	if inst, found := fTable[uint(ci.GetF())]; found {
		completed, interrupt := inst(e)
		if interrupt != nil {
			e.PostInterrupt(interrupt)
		} else if completed {
			e.clearStorageLocks()
			e.activityStatePacket.indicatorKeyRegister.instructionInF0 = false
			if !e.preventPCUpdate {
				e.activityStatePacket.programAddressRegister.IncrementProgramCounter()
			}
		}
	} else {
		e.PostInterrupt(NewInvalidInstructionInterrupt(InvalidInstructionBadFunctionCode))
	}
}

// extractPartialWord pulls the partial word indicated by the partialWordIndicator and the quarterWordMode flag
// from the given 36-bit source value.
func (e *InstructionEngine) extractPartialWord(source uint64, partialWordIndicator uint, quarterWordMode bool) uint64 {
	switch partialWordIndicator {
	case pkg.JFIELD_W:
		return pkg.GetW(source)
	case pkg.JFIELD_H2:
		return pkg.GetH2(source)
	case pkg.JFIELD_H1:
		return pkg.GetH1(source)
	case pkg.JFIELD_XH2:
		return pkg.GetXH2(source)
	case pkg.JFIELD_XH1: // XH1 or Q2
		if quarterWordMode {
			return pkg.GetQ2(source)
		} else {
			return pkg.GetXH1(source)
		}
	case pkg.JFIELD_T3: // T3 or Q4
		if quarterWordMode {
			return pkg.GetQ4(source)
		} else {
			return pkg.GetXT3(source)
		}
	case pkg.JFIELD_T2: // T2 or Q3
		if quarterWordMode {
			return pkg.GetQ3(source)
		} else {
			return pkg.GetXT2(source)
		}
	case pkg.JFIELD_T1: // T1 or Q1
		if quarterWordMode {
			return pkg.GetQ1(source)
		} else {
			return pkg.GetXT1(source)
		}
	case pkg.JFIELD_S6:
		return pkg.GetS6(source)
	case pkg.JFIELD_S5:
		return pkg.GetS5(source)
	case pkg.JFIELD_S4:
		return pkg.GetS4(source)
	case pkg.JFIELD_S3:
		return pkg.GetS3(source)
	case pkg.JFIELD_S2:
		return pkg.GetS2(source)
	case pkg.JFIELD_S1:
		return pkg.GetS1(source)
	}

	return source
}

// fetchInstructionWord retrieves the next instruction word from the appropriate bank.
// For extended mode this is straight-forward.
// For basic mode, we have to hunt around a bit to make sure we pull it from the most appropriate bank.
// If something bad happens, an interrupt is posted and we return false
func (e *InstructionEngine) fetchInstructionWord() bool {
	basicMode := e.activityStatePacket.designatorRegister.BasicModeEnabled
	programCounter := uint64(e.activityStatePacket.programAddressRegister.GetProgramCounter())

	var bReg *BaseRegister
	if basicMode {
		brIndex := e.FindBasicModeBank(programCounter, true)
		if brIndex == 0 {
			e.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, false))
			return false
		}

		bReg = e.baseRegisters[brIndex]
		if !e.isReadAllowed(bReg) {
			e.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, false))
			return false
		}
	} else {
		bReg = e.baseRegisters[0]
		intp := e.checkAccessLimits(bReg, programCounter, true, false, false, e.activityStatePacket.indicatorKeyRegister.accessKey)
		if intp != nil {
			return false
		}
	}

	if bReg.voidFlag || bReg.largeSizeFlag {
		e.PostInterrupt(NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, false))
		return false
	}

	pcOffset := programCounter - uint64(bReg.lowerLimitNormalized)
	e.activityStatePacket.currentInstruction = InstructionWord(bReg.storage[pcOffset])
	e.activityStatePacket.indicatorKeyRegister.instructionInF0 = true

	return true
}

// findBankDescriptor retrieves a struct to describe the given named bank.
//
//	This is for interrupt handling.
//	The bank name is in L,BDI format.
//	bankLevel level of the bank, 0:7
//	bankDescriptorIndex BDI of the bank 0:077777
func (e *InstructionEngine) findBankDescriptor(bankLevel uint, bankDescriptorIndex uint) (*BankDescriptor, bool) {
	// The bank descriptor tables for bank levels 0 through 7 are described by the banks based on B16 through B23.
	// The bank descriptor will be the {n}th bank descriptor in the particular bank descriptor table,
	// where {n} is the bank descriptor index.
	bdRegIndex := bankLevel + 16
	if e.baseRegisters[bdRegIndex].voidFlag {
		e.PostInterrupt(NewAddressingExceptionInterrupt(AddressingExceptionFatal, bankLevel, bankDescriptorIndex))
		return nil, false
	}

	//  bdStorage contains the BDT for the given bank_name level
	//  bdTableOffset indicates the offset into the BDT, where the bank descriptor is to be found.
	bdStorage := e.baseRegisters[bdRegIndex].storage
	bdTableOffset := bankDescriptorIndex + 8
	if bdTableOffset+8 > uint(len(bdStorage)) {
		e.PostInterrupt(NewAddressingExceptionInterrupt(AddressingExceptionFatal, bankLevel, bankDescriptorIndex))
		return nil, false
	}

	//  Create and return a BankDescriptor object
	bd := NewBankDescriptorFromStorage(bdStorage[bdTableOffset : bdTableOffset+8])
	return bd, true
}

// getAbsoluteAddress converts a relative address to an absolute address.
// The AbsoluteAddress is passed as a reference so that we can avoid leaving little structs all over the heap.
func (e *InstructionEngine) getAbsoluteAddress(baseRegister *BaseRegister, relativeAddress uint64, addr *AbsoluteAddress) {
	addr.segment = baseRegister.baseAddress.segment
	actualOffset := relativeAddress - uint64(baseRegister.lowerLimitNormalized)
	addr.offset = uint(uint64(baseRegister.baseAddress.offset) + actualOffset)
}

// findBaseRegisterIndex locates the index of the base register which represents the bank which contains the given
// relative address. Does appropriate limits checking.  Delegates to the appropriate basic or extended mode implementation.
//
// relativeAddress: relative address to be considered
// updateDesignatorRegister: if true and if we are in basic mode, we update the basic mode bank selection bit in the designator register if necessary
//
// Returns complete==true and the base register index if successful
//
// Returns complete==false if address resolution is unfinished (such as can happen in Basic Mode with Indirect Addressing).
// In this case, caller should call back here again after checking for any pending interrupts.
//
// Returns an interrupt if any interrupt needs to be raised. In this case, the instruction is incomplete and should
// be retried if appropriate. Caller should post the interrupt.
func (e *InstructionEngine) findBaseRegisterIndex(relativeAddress uint64, updateDesignatorRegister bool) (complete bool, index uint, interrupt Interrupt) {
	complete = false
	index = 0
	interrupt = nil

	dr := e.activityStatePacket.designatorRegister
	if dr.BasicModeEnabled {
		//  Find the bank containing the current offset.
		//  We don't need to check for storage limits, since this is done for us by findBasicModeBank() in terms of
		//  returning a zero.
		brIndex := e.FindBasicModeBank(relativeAddress, updateDesignatorRegister)
		if brIndex == 0 {
			interrupt = NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, false)
			return
		}

		//  Are we doing indirect addressing?
		ci := e.activityStatePacket.currentInstruction
		if ci.GetI() != 0 {
			//  Increment the X register (if any) indicated by F0 (if H bit is set, of course)
			e.incrementIndexRegisterInF0()
			bReg := e.baseRegisters[brIndex]

			//  Ensure we can read from the selected bank
			if !e.isReadAllowed(bReg) {
				interrupt = NewReferenceViolationInterrupt(ReferenceViolationReadAccess, false)
				return
			}

			ikr := e.activityStatePacket.indicatorKeyRegister
			interrupt = e.checkAccessLimits(bReg, relativeAddress, false, true, false, ikr.accessKey)
			if interrupt != nil {
				return
			}

			//  Get xhiu fields from the referenced word, and place them into _currentInstruction,
			//  then throw UnresolvedAddressException so the caller knows we're not done here.
			wx := relativeAddress - uint64(bReg.lowerLimitNormalized)
			e.activityStatePacket.currentInstruction.SetXHIU(bReg.storage[wx].GetW())
			return
		}

		//  We're at our final destination
		complete = true
		index = brIndex
	} else {
		index = e.getEffectiveBaseRegisterIndex()
	}

	return
}

func (e *InstructionEngine) getEffectiveBaseRegisterIndex() uint {
	//  If PP < 2, we use the i-bit and the b-field to select the base registers from B0 to B31.
	//  For PP >= 2, we only use the b-field, to select base registers from B0 to B15 (See IP PRM 4.3.7).
	if e.activityStatePacket.designatorRegister.processorPrivilege < 2 {
		return uint(e.activityStatePacket.currentInstruction.GetIB())
	} else {
		return uint(e.activityStatePacket.currentInstruction.GetB())
	}
}

// getImmediateOperand retrieves an operand in the case where the u (and possibly h and i) fields
// comprise the requested data ... e.g., for immediate operands.
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
func (e *InstructionEngine) getImmediateOperand() (operand uint64, interrupt Interrupt) {
	operand = 0
	interrupt = nil

	ci := e.activityStatePacket.currentInstruction
	dr := e.activityStatePacket.designatorRegister

	exec24Index := dr.Executive24BitIndexingEnabled
	privilege := dr.processorPrivilege
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

		//  Add the contents of Xx(m), and do index register incrementation if appropriate.
		xReg := e.GetExecOrUserXRegister(uint(ci.GetX()))

		//  24-bit indexing?
		if !dr.BasicModeEnabled && (privilege < 2) && exec24Index {
			//  Add the 24-bit modifier
			operand = pkg.AddSimple(operand, xReg.GetXM24())
			if ci.GetH() != 0 {
				e.GetExecOrUserXRegister(uint(ci.GetX())).IncrementModifier24()
			}
		} else {
			//  Add the 18-bit modifier
			operand = pkg.AddSimple(operand, xReg.GetXM())
			if ci.GetH() != 0 {
				e.GetExecOrUserXRegister(uint(ci.GetX())).IncrementModifier()
			}
		}
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

func (e *InstructionEngine) incrementIndexRegisterInF0() {
	ci := e.activityStatePacket.currentInstruction
	if (ci.GetX() != 0) && (ci.GetH() != 0) {
		iReg := e.GetExecOrUserXRegister(uint(ci.GetX()))
		dReg := e.activityStatePacket.designatorRegister
		if !dReg.BasicModeEnabled && dReg.Executive24BitIndexingEnabled && (dReg.processorPrivilege < 2) {
			iReg.IncrementModifier24()
		} else {
			iReg.IncrementModifier()
		}
	}
}

func (e *InstructionEngine) isGRSAccessAllowed(registerIndex uint64, processorPrivilege uint, writeAccess bool) bool {
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

func (e *InstructionEngine) isReadAllowed(bReg *BaseRegister) bool {
	permissions := bReg.GetEffectivePermissions(e.activityStatePacket.indicatorKeyRegister.accessKey)
	return permissions.CanRead()
}

// isWithinLimits evaluates the given offset within the constraints of the given base register,
// returning true if the offset is within those constraints, else false
func (e *InstructionEngine) isWithinLimits(bReg *BaseRegister, offset uint64) bool {
	return !bReg.voidFlag &&
		(offset >= uint64(bReg.lowerLimitNormalized)) &&
		(offset <= uint64(bReg.upperLimitNormalized))
}
