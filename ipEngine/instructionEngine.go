// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
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
	baseRegisters       [32]*pkg.BaseRegister
	generalRegisterSet  *GeneralRegisterSet

	//	If not nil, describes an interrupt which needs to be handled as soon as possible
	pendingInterrupt pkg.Interrupt

	//	See 2.4.2
	//	Should this be saved off somewhere during an interrupt?
	jumpHistory                 [JumpHistoryTableThreshold]pkg.Word36
	jumpHistoryIndex            int
	jumpHistoryThresholdReached bool

	//	If true, the current (or most recent) instructionType has set the PAR.PC the way it wants,
	//	and we should not increment it for the next instruction
	preventPCUpdate bool

	mutex sync.Mutex

	breakpointAddress  pkg.AbsoluteAddress
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
		e.baseRegisters[bx] = pkg.NewVoidBaseRegister()
	}

	e.generalRegisterSet = NewGeneralRegisterSet()
	e.activityStatePacket.designatorRegister = &DesignatorRegister{}
	e.activityStatePacket.indicatorKeyRegister = &IndicatorKeyRegister{}
	e.activityStatePacket.indicatorKeyRegister.accessKey = pkg.NewAccessKeyFromComponents(0, 0)
	e.activityStatePacket.programAddressRegister = &ProgramAddressRegister{}

	e.breakpointRegister = BreakpointNone
	e.stopReason = NotStopped

	return e
}

func (e *InstructionEngine) Dump() {
	fmt.Printf("Instruction Engine Dump ---------------------------------------------------------------------\n")

	if e.HasPendingInterrupt() {
		//	TODO move the following to interrupt.go as GetString(i Interrupt) and include text description
		//	 of (at least) the interrupt class
		fmt.Printf("Pending Interrupt Class:%v SSF:%v ISW0:%012o ISW1:%012o\n",
			e.pendingInterrupt.GetClass(),
			e.pendingInterrupt.GetShortStatusField(),
			e.pendingInterrupt.GetStatusWord0(),
			e.pendingInterrupt.GetStatusWord1())
	}

	var f0String string
	if e.activityStatePacket.indicatorKeyRegister.instructionInF0 {
		// TODO disassemble the instruction as well as showing the octal word
		f0String = fmt.Sprintf("%012o", e.activityStatePacket.currentInstruction)
	} else {
		f0String = "invalid"
	}
	fmt.Printf("  F0: %s\n", f0String)

	par := e.activityStatePacket.programAddressRegister
	fmt.Printf("  PAR.PC L:%o BDI:%05o PC:%06o\n", par.level, par.bankDescriptorIndex, par.programCounter)

	ikr := e.activityStatePacket.indicatorKeyRegister
	fmt.Printf("  Indicator Key Register: %012o\n", ikr.GetComposite())
	fmt.Printf("    Access Key:        %s\n", ikr.accessKey.GetString())
	fmt.Printf("    SSF:               %03o\n", ikr.shortStatusField)
	fmt.Printf("    Interrupt Class:   %03o\n", ikr.interruptClassField)
	fmt.Printf("    EXR Instruction:   %v\n", ikr.executeRepeatedInstruction)
	fmt.Printf("    Breakpoint Match:  %v\n", ikr.breakpointRegisterMatchCondition)
	fmt.Printf("    Software Break:    %v\n", ikr.softwareBreak)
	fmt.Printf("    Instruction in F0: %v\n", ikr.instructionInF0)

	dr := e.activityStatePacket.designatorRegister
	fmt.Printf("  Designator Register: %012o\n", dr.GetComposite())
	fmt.Printf("    FHIP:                        %v\n", dr.faultHandlingInProgress)
	fmt.Printf("    Executive 24-bit Indexing:   %v\n", dr.executive24BitIndexingEnabled)
	fmt.Printf("    Quantum Timer Enable:        %v\n", dr.quantumTimerEnabled)
	fmt.Printf("    Deferrable Interrupt Enable: %v\n", dr.deferrableInterruptEnabled)
	fmt.Printf("    Processor Privilege:         %v\n", dr.processorPrivilege)
	fmt.Printf("    Basic Mode:                  %v\n", dr.basicModeEnabled)
	fmt.Printf("    Exec Register Set Selection: %v\n", dr.execRegisterSetSelected)
	fmt.Printf("    Carry:                       %v\n", dr.carry)
	fmt.Printf("    Overflow:                    %v\n", dr.overflow)
	fmt.Printf("    Characteristic Underflow:    %v\n", dr.characteristicUnderflow)
	fmt.Printf("    Characteristic Overflow:     %v\n", dr.characteristicOverflow)
	fmt.Printf("    Divide Check:                %v\n", dr.divideCheck)
	fmt.Printf("    Operation Trap Enable:       %v\n", dr.operationTrapEnabled)
	fmt.Printf("    Arithmetic Exception Enable: %v\n", dr.arithmeticExceptionEnabled)
	fmt.Printf("    Basic Mode Base Reg Sel:     %v\n", dr.basicModeBaseRegisterSelection)
	fmt.Printf("    Quarter Word Selection:      %v\n", dr.quarterWordModeEnabled)

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
}

// FindBasicModeBank takes a relative address and determines which (if any) of the basic mode banks
// currently based on BDR12-15 is to be selected for that address.
// Set updatedDB31 true if you want the code to update designator register Bit31 in the event we cross
// primary/secondary bank pairs.
// Returns the bank descriptor index (from 12 to 15) for the proper bank descriptor.
// Returns zero if the address is not within any of the based bank limits.
func (e *InstructionEngine) FindBasicModeBank(relativeAddress uint64, updateDB31 bool) uint {
	db31 := e.activityStatePacket.designatorRegister.basicModeBaseRegisterSelection
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
				e.activityStatePacket.designatorRegister.basicModeBaseRegisterSelection = !db31
			}

			return brIndex
		}
	}

	return 0
}

func (e *InstructionEngine) ClearInterrupt() {
	e.pendingInterrupt = nil
}

// GetActiveBaseTableEntry retrieves a pointer to the ABET for the indicated base register 0 to 15
func (e *InstructionEngine) GetActiveBaseTableEntry(index uint) *ActiveBaseTableEntry {
	return e.activeBaseTable[index]
}

// GetBaseRegister retrieves a pointer to the indicated base register
func (e *InstructionEngine) GetBaseRegister(index uint) *pkg.BaseRegister {
	return e.baseRegisters[index]
}

func (e *InstructionEngine) GetCurrentInstruction() InstructionWord {
	return e.activityStatePacket.currentInstruction
}

func (e *InstructionEngine) GetDesignatorRegister() *DesignatorRegister {
	return e.activityStatePacket.designatorRegister
}

// GetExecOrUserARegister retrieves either the EA{index} or A{index} register
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserARegister(registerIndex uint) *pkg.Word36 {
	return e.generalRegisterSet.GetRegister(e.GetExecOrUserARegisterIndex(registerIndex))
}

// GetExecOrUserARegisterIndex retrieves the GRS index of either EA{index} or A{index}
// depending upon the setting of designator register ExecRegisterSetSelected
func (e *InstructionEngine) GetExecOrUserARegisterIndex(registerIndex uint) uint {
	if e.activityStatePacket.designatorRegister.execRegisterSetSelected {
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
	if e.activityStatePacket.designatorRegister.execRegisterSetSelected {
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
	if e.activityStatePacket.designatorRegister.execRegisterSetSelected {
		return EX0 + registerIndex
	} else {
		return X0 + registerIndex
	}
}

// GetGeneralRegisterSet retrieves a pointer to the GRS
func (e *InstructionEngine) GetGeneralRegisterSet() *GeneralRegisterSet {
	return e.generalRegisterSet
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
func (e *InstructionEngine) GetJumpOperand(updateDesignatorRegister bool) (complete bool, operand uint64, interrupt pkg.Interrupt) {
	complete = true
	interrupt = nil
	operand = e.calculateRelativeAddressForJump()

	//  The following bit is how we deal with indirect addressing for basic mode.
	//  If we are doing that, it will update the U portion of the current instruction with new address information,
	//  then throw UnresolvedAddressException which will eventually route us back through here again, but this
	//  time with new address info (in reladdress), and we keep doing this until we're not doing indirect addressing.
	asp := e.activityStatePacket
	if asp.designatorRegister.basicModeEnabled && (asp.currentInstruction.GetI() != 0) {
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
func (e *InstructionEngine) GetOperand(grsDestination bool, grsCheck bool, allowImmediate bool, allowPartial bool) (complete bool, operand uint64, interrupt pkg.Interrupt) {
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
	basicMode := dReg.basicModeEnabled
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
			interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationReadAccess, true)
			return
		}

		//  If we are GRS or not allowing partial word transfers, do a full word.
		//  Otherwise, honor partial word transferring.
		if grsDestination || !allowPartial {
			operand = grs.GetRegister(uint(relAddress)).GetW()
		} else {
			qWordMode := dReg.quarterWordModeEnabled
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

	var absAddress pkg.AbsoluteAddress
	e.getAbsoluteAddress(baseRegister, relAddress, &absAddress)
	e.checkBreakpoint(BreakpointRead, &absAddress)

	readOffset := relAddress - baseRegister.GetLowerLimitNormalized()
	operand = baseRegister.GetStorage()[readOffset].GetW()
	if allowPartial {
		qWordMode := dReg.quarterWordModeEnabled
		operand = e.extractPartialWord(operand, jField, qWordMode)
	}

	complete = true
	return
}

func (e *InstructionEngine) GetPARPC() uint64 {
	return uint64(e.activityStatePacket.programAddressRegister.GetComposite())
}

func (e *InstructionEngine) HasPendingInterrupt() bool {
	return e.pendingInterrupt != nil
}

func (e *InstructionEngine) PopInterrupt() pkg.Interrupt {
	i := e.pendingInterrupt
	e.pendingInterrupt = nil
	return i
}

// PostInterrupt posts a new interrupt, provided that no higher-priority interrupt is already pending.
func (e *InstructionEngine) PostInterrupt(i pkg.Interrupt) {
	if e.pendingInterrupt == nil || i.GetClass() < e.pendingInterrupt.GetClass() {
		e.pendingInterrupt = i
	}
}

// SetBaseRegister sets the base register identified by brIndex (0 to 15) to the given register
func (e *InstructionEngine) SetBaseRegister(brIndex uint, register *pkg.BaseRegister) {
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

// SetPARPC sets the individual components of PAR.PC
func (e *InstructionEngine) SetPARPC(level uint, index uint, counter uint) {
	par := e.activityStatePacket.programAddressRegister
	par.SetLevel(level)
	par.SetBankDescriptorIndex(index)
	par.SetProgramCounter(counter)
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

// StoreOperand handles the general case of storing an operand either to storage or to a GRS location
//
// grsSource: true if the value came from a register, so we know whether to ignore partial-word transfers
// grsCheck: true if relative addresses < 0200 should be considered GRS locations
// checkImmediate: true if we should consider j-fields 016 and 017 as immediate addressing (and throw away the operand)
// allowPartial: true if we should allow partial-word transfers (subject to GRS-GRS transfers)
// operand: value to be stored
//
// returns complete==true indicates the operation is complete - if false, address resolution is unfinished and the
// caller should try again
// returns an interrupt if anything goes wrong, which the caller should post
func (e *InstructionEngine) StoreOperand(grsSource bool, grsCheck bool, checkImmediate bool, allowPartial bool, operand uint64) (complete bool, interrupt pkg.Interrupt) {
	complete = false
	interrupt = nil

	//  If we allow immediate addressing mode and j-field is U or XU... we do nothing.
	ci := e.activityStatePacket.currentInstruction

	jField := uint(ci.GetJ())
	if (checkImmediate) && (jField >= 016) {
		complete = true
		return
	}

	dr := e.activityStatePacket.designatorRegister
	relAddress := e.calculateRelativeAddressForGRSOrStorage()
	basicMode := dr.basicModeEnabled
	pPriv := dr.processorPrivilege

	var baseRegisterIndex uint
	if !basicMode {
		baseRegisterIndex = uint(ci.GetB())
		if (pPriv < 2) && (ci.GetI() != 0) {
			baseRegisterIndex += 16
		}
	}

	if (grsCheck) && (basicMode || (baseRegisterIndex == 0)) && (relAddress < 0200) {
		e.incrementIndexRegisterInF0()

		//  First, do accessibility checks
		if !e.isGRSAccessAllowed(relAddress, pPriv, true) {
			interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationWriteAccess, false)
			return
		}

		//  If we are GRS or not allowing partial word transfers, do a full word.
		//  Otherwise, honor partial word transfer.
		if !grsSource && allowPartial {
			qWordMode := dr.quarterWordModeEnabled
			originalValue := e.generalRegisterSet.GetRegister(uint(relAddress)).GetW()
			newValue := e.injectPartialWord(originalValue, operand, jField, qWordMode)
			e.generalRegisterSet.GetRegister(uint(relAddress)).SetW(newValue)
		} else {
			e.generalRegisterSet.GetRegister(uint(relAddress)).SetW(operand)
		}

		complete = true
		return
	}

	//  This is going to be a storage thing...
	if basicMode {
		complete, baseRegisterIndex, interrupt = e.findBaseRegisterIndex(relAddress, false)
		if !complete || (interrupt != nil) {
			return
		}
	}

	bReg := e.baseRegisters[baseRegisterIndex]
	ikr := e.activityStatePacket.indicatorKeyRegister
	interrupt = e.checkAccessLimits(bReg, relAddress, false, false, true, ikr.accessKey)
	if interrupt != nil {
		return
	}

	e.incrementIndexRegisterInF0()

	var absAddr pkg.AbsoluteAddress
	e.getAbsoluteAddress(bReg, relAddress, &absAddr)
	e.checkBreakpoint(BreakpointWrite, &absAddr)

	offset := relAddress - bReg.GetLowerLimitNormalized()
	if allowPartial {
		qWordMode := dr.quarterWordModeEnabled
		originalValue := bReg.GetStorage()[offset].GetW()
		newValue := e.injectPartialWord(originalValue, operand, jField, qWordMode)
		bReg.GetStorage()[offset].SetW(newValue)
	} else {
		bReg.GetStorage()[offset].SetW(operand)
	}

	complete = true
	return
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
	if dr.basicModeEnabled {
		addend1 = ci.GetU()
		if xReg != nil {
			addend2 = xReg.GetSignedXM()
		}
	} else {
		addend1 = ci.GetD()
		if xReg != nil {
			if dr.executive24BitIndexingEnabled && dr.processorPrivilege < 2 {
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
	if dr.basicModeEnabled {
		if xReg != nil {
			addend2 = xReg.GetSignedXM()
		}
	} else {
		addend1 = ci.GetU()
		if xReg != nil {
			if dr.executive24BitIndexingEnabled && dr.processorPrivilege < 2 {
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
	bReg *pkg.BaseRegister,
	fetchFlag bool,
	readFlag bool,
	writeFlag bool,
	accessKey *pkg.AccessKey) pkg.Interrupt {
	perms := bReg.GetEffectivePermissions(accessKey)
	if e.activityStatePacket.designatorRegister.basicModeEnabled && fetchFlag && !perms.CanEnter() {
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

// checkAccessLimits checks the accessibility of a given relative address in the bank described by this
// base register for the given flags, using the given key.
// If the check fails, we return an interrupt which the caller should post
func (e *InstructionEngine) checkAccessLimits(
	bReg *pkg.BaseRegister,
	relativeAddress uint64,
	fetchFlag bool,
	readFlag bool,
	writeFlag bool,
	accessKey *pkg.AccessKey) pkg.Interrupt {

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
func (e *InstructionEngine) checkAccessLimitsForAddress(bReg *pkg.BaseRegister, relativeAddress uint64, fetchFlag bool) pkg.Interrupt {

	// TODO if we try to execute something in GRS - we take ReferenceViolationInterruptClass, 01, 0, 0

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
	relativeAddress uint,
	addressCount uint,
	readFlag bool,
	writeFlag bool,
	accessKey *pkg.AccessKey) pkg.Interrupt {
	if (uint64(relativeAddress) < bReg.GetLowerLimitNormalized()) ||
		(uint64(relativeAddress+addressCount-1) > bReg.GetUpperLimitNormalized()) {
		return pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, false)
	}

	return e.checkAccessibility(bReg, false, readFlag, writeFlag, accessKey)
}

func (e *InstructionEngine) checkBreakpoint(comp BreakpointComparison, addr *pkg.AbsoluteAddress) {
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
	fTable := FunctionTable[dr.basicModeEnabled]
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
		e.PostInterrupt(pkg.NewInvalidInstructionInterrupt(pkg.InvalidInstructionBadFunctionCode))
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
	basicMode := e.activityStatePacket.designatorRegister.basicModeEnabled
	programCounter := uint64(e.activityStatePacket.programAddressRegister.GetProgramCounter())

	var bReg *pkg.BaseRegister
	if basicMode {
		brIndex := e.FindBasicModeBank(programCounter, true)
		if brIndex == 0 {
			e.PostInterrupt(pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, false))
			return false
		}

		bReg = e.baseRegisters[brIndex]
		if !e.isReadAllowed(bReg) {
			e.PostInterrupt(pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, false))
			return false
		}
	} else {
		bReg = e.baseRegisters[0]
		ikr := e.activityStatePacket.indicatorKeyRegister
		intp := e.checkAccessLimits(bReg, programCounter, true, false, false, ikr.accessKey)
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
	e.activityStatePacket.currentInstruction = InstructionWord(bReg.GetStorage()[pcOffset])
	e.activityStatePacket.indicatorKeyRegister.instructionInF0 = true

	return true
}

// findBankDescriptor retrieves a struct to describe the given named bank.
//
//	This is for interrupt handling.
//	The bank name is in L,BDI format.
//	bankLevel level of the bank, 0:7
//	bankDescriptorIndex BDI of the bank 0:077777
func (e *InstructionEngine) findBankDescriptor(bankLevel uint, bankDescriptorIndex uint) (*pkg.BankDescriptor, bool) {
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
	if bdTableOffset+8 > uint(len(bdStorage)) {
		e.PostInterrupt(pkg.NewAddressingExceptionInterrupt(pkg.AddressingExceptionFatal, bankLevel, bankDescriptorIndex))
		return nil, false
	}

	//  Create and return a BankDescriptor object
	bd := pkg.NewBankDescriptorFromStorage(bdStorage[bdTableOffset : bdTableOffset+8])
	return bd, true
}

// getAbsoluteAddress converts a relative address to an absolute address.
// The AbsoluteAddress is passed as a reference so that we can avoid leaving little structs all over the heap.
func (e *InstructionEngine) getAbsoluteAddress(baseRegister *pkg.BaseRegister, relativeAddress uint64, addr *pkg.AbsoluteAddress) {
	addr.SetSegment(baseRegister.GetBaseAddress().GetSegment())
	actualOffset := relativeAddress - baseRegister.GetLowerLimitNormalized()
	addr.SetOffset(baseRegister.GetBaseAddress().GetOffset() + uint(actualOffset))
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
func (e *InstructionEngine) findBaseRegisterIndex(relativeAddress uint64, updateDesignatorRegister bool) (complete bool, index uint, interrupt pkg.Interrupt) {
	complete = false
	index = 0
	interrupt = nil

	dr := e.activityStatePacket.designatorRegister
	if dr.basicModeEnabled {
		//  Find the bank containing the current offset.
		//  We don't need to check for storage limits, since this is done for us by findBasicModeBank() in terms of
		//  returning a zero.
		brIndex := e.FindBasicModeBank(relativeAddress, updateDesignatorRegister)
		if brIndex == 0 {
			interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationStorageLimits, false)
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
				interrupt = pkg.NewReferenceViolationInterrupt(pkg.ReferenceViolationReadAccess, false)
				return
			}

			ikr := e.activityStatePacket.indicatorKeyRegister
			interrupt = e.checkAccessLimits(bReg, relativeAddress, false, true, false, ikr.accessKey)
			if interrupt != nil {
				return
			}

			//  Get xhiu fields from the referenced word, and place them into _currentInstruction,
			//  then throw UnresolvedAddressException so the caller knows we're not done here.
			wx := relativeAddress - bReg.GetLowerLimitNormalized()
			e.activityStatePacket.currentInstruction.SetXHIU(bReg.GetStorage()[wx].GetW())
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
func (e *InstructionEngine) getImmediateOperand() (operand uint64, interrupt pkg.Interrupt) {
	operand = 0
	interrupt = nil

	ci := e.activityStatePacket.currentInstruction
	dr := e.activityStatePacket.designatorRegister

	exec24Index := dr.executive24BitIndexingEnabled
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
		if !dr.basicModeEnabled && (privilege < 2) && exec24Index {
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
		if !dReg.basicModeEnabled && dReg.executive24BitIndexingEnabled && (dReg.processorPrivilege < 2) {
			iReg.IncrementModifier24()
		} else {
			iReg.IncrementModifier()
		}
	}
}

// injectPartialWord creates a value which is comprised of an original value, with a new value inserted there-in
// under j-field control.
func (e *InstructionEngine) injectPartialWord(originalValue uint64, newValue uint64, jField uint, quarterWordMode bool) uint64 {
	switch jField {
	case pkg.JFIELD_W:
		return newValue
	case pkg.JFIELD_H2:
	case pkg.JFIELD_XH2:
		return pkg.SetH2(originalValue, newValue)
	case pkg.JFIELD_H1:
		return pkg.SetH1(originalValue, newValue)
	case pkg.JFIELD_XH1: // XH1 or Q2
		if quarterWordMode {
			return pkg.SetQ2(originalValue, newValue)
		} else {
			return pkg.SetH1(originalValue, newValue)
		}
	case pkg.JFIELD_T3: // T3 or Q4
		if quarterWordMode {
			return pkg.SetQ4(originalValue, newValue)
		} else {
			return pkg.SetT3(originalValue, newValue)
		}
	case pkg.JFIELD_T2: // T2 or Q3
		if quarterWordMode {
			return pkg.SetQ3(originalValue, newValue)
		} else {
			return pkg.SetT2(originalValue, newValue)
		}
	case pkg.JFIELD_T1: // T1 or Q1
		if quarterWordMode {
			return pkg.SetQ1(originalValue, newValue)
		} else {
			return pkg.SetT1(originalValue, newValue)
		}
	case pkg.JFIELD_S6:
		return pkg.SetS6(originalValue, newValue)
	case pkg.JFIELD_S5:
		return pkg.SetS5(originalValue, newValue)
	case pkg.JFIELD_S4:
		return pkg.SetS4(originalValue, newValue)
	case pkg.JFIELD_S3:
		return pkg.SetS3(originalValue, newValue)
	case pkg.JFIELD_S2:
		return pkg.SetS2(originalValue, newValue)
	case pkg.JFIELD_S1:
		return pkg.SetS1(originalValue, newValue)
	}

	return originalValue
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

func (e *InstructionEngine) isReadAllowed(bReg *pkg.BaseRegister) bool {
	permissions := bReg.GetEffectivePermissions(e.activityStatePacket.indicatorKeyRegister.accessKey)
	return permissions.CanRead()
}

// isWithinLimits evaluates the given offset within the constraints of the given base register,
// returning true if the offset is within those constraints, else false
func (e *InstructionEngine) isWithinLimits(bReg *pkg.BaseRegister, offset uint64) bool {
	return !bReg.IsVoid() &&
		(offset >= bReg.GetLowerLimitNormalized()) &&
		(offset <= bReg.GetUpperLimitNormalized())
}
