// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"khalehla/pkg"
	"khalehla/tasm"
)

type UnitTestEngine struct {
	executable *tasm.Executable

	//  Our own private main storage entity
	storage *pkg.MainStorage

	//	Lookup for absolute addresses of loaded banks.
	//	key is BDI, value is the absolute address (base address) of the bank.
	bankAddresses map[uint]*pkg.AbsoluteAddress

	//	BDTable lookup into storage - key is level (0 to 7), and value is the absolute address of the table in storage
	//  If there is not an entry for a level, then there is not a BDT for that level.
	//  For now, we have one main storage segment for each BDT, so we can easily tell the size of the BDT from the
	//  size of the segment, and the offset of the table (from the start of the segment) is always zero.
	bankDescriptorTableAddresses map[int]*pkg.AbsoluteAddress

	engine *InstructionEngine
}

func NewExecutor() *UnitTestEngine {
	e := &UnitTestEngine{
		storage: pkg.NewMainStorage(100),
	}

	return e
}

func (ute *UnitTestEngine) Clear() {
	ute.storage.Clear()
	ute.bankAddresses = make(map[uint]*pkg.AbsoluteAddress)
	ute.bankDescriptorTableAddresses = make(map[int]*pkg.AbsoluteAddress)
	ute.executable = nil
	ute.engine = nil
}

func (ute *UnitTestEngine) Load(executable *tasm.Executable) error {
	fmt.Printf("\nLoading Executable...\n")
	ute.Clear()
	ute.executable = executable

	// maps bdi to segment index of segment containing the bank
	segIndexMap := make(map[int]uint)

	for _, bank := range executable.GetBanks() {
		//	allocate a segment from storage and copy the bank to the segment
		//  there are no fix-ups required; all our binaries are self-contained.

		bdi := bank.GetBankDescriptorIndex()
		segIndex, err := ute.storage.Allocate(bank.GetCodeLength())
		if err != nil {
			return err
		}
		seg, _ := ute.storage.GetSegment(segIndex)
		segIndexMap[int(bdi)] = segIndex

		code := bank.GetCode()
		for cx := 0; cx < len(code); cx++ {
			seg[cx].SetW(code[cx])
		}

		ute.bankAddresses[bdi] = pkg.NewAbsoluteAddress(segIndex, 0)
		//	now build a bank descriptor for the bank.
		//	if this is the first bank with its level, then we have to allocate space for a bdt for the level.
		bdiLevel := bdi >> 15
		bdiOffset := bdi & 077777
		newBDTLen := (bdiOffset + 1) << 3

		//	For this algorithm, each BDT gets its own segment and thus is always at offset 0 of the segment
		absAddr, ok := ute.bankDescriptorTableAddresses[int(bdiLevel)]
		var bdtSegment uint
		var bdTable []pkg.Word36
		if ok {
			bdtSegment := absAddr.GetSegment()
			bdTable, _ = ute.storage.GetSegment(bdtSegment)
			if uint(len(bdTable)) < newBDTLen {
				err := ute.storage.Resize(bdtSegment, newBDTLen)
				if err != nil {
					return err
				}
			}
		} else {
			bdtSegment, err = ute.storage.Allocate(newBDTLen)
			if err != nil {
				return err
			}
			bdTable, _ = ute.storage.GetSegment(bdtSegment)
			absAddr = &pkg.AbsoluteAddress{}
			absAddr.SetSegment(bdtSegment)
			ute.bankDescriptorTableAddresses[int(bdiLevel)] = absAddr
		}
		bank.GetBankDescriptor().SetBaseAddress(absAddr)

		bdOffset := bdiOffset * 8
		bank.GetBankDescriptor().Serialize(bdTable[bdOffset : bdOffset+8])
	}

	ute.storage.Dump() // TODO remove
	return nil
}

func (ute *UnitTestEngine) getBankDescriptor(bdi uint) (*pkg.BankDescriptor, error) {
	//	TODO need to do translation for extended mode
	level := int(bdi >> 15)
	index := bdi & 077777

	bdtAddress, ok := ute.bankDescriptorTableAddresses[level]
	if !ok {
		return nil, fmt.Errorf("no BDT for level %d for BDI %06o", level, bdi)
	}

	offset := bdtAddress.GetOffset() + (index * 8)
	slice, ok := ute.storage.GetSlice(bdtAddress.GetSegment(), offset, offset+8)
	if !ok {
		return nil, fmt.Errorf("cannot retrieve BD for BDI %06o", bdi)
	}

	return pkg.NewBankDescriptorFromStorage(slice), nil
}

func (ute *UnitTestEngine) Run() error {
	fmt.Printf("\nSetting up to run...\n")
	if ute.executable == nil {
		return fmt.Errorf("no executable has been loaded")
	}

	ute.engine = NewEngine(ute.storage)
	for brx := uint(0); brx < 16; brx++ {
		ute.engine.SetBaseRegister(brx, pkg.NewVoidBaseRegister())
	}

	//	Load BDT base registers B16 to B23. BDTable level 0 -> B16, 1 -> B17, etc.
	//  For any level which does not have a BDT, we make the corresponding base register void.
	fmt.Printf("  Loading base registers for bank descriptor tables...\n")
	for level, address := range ute.bankDescriptorTableAddresses {
		table, _ := ute.storage.GetSegment(address.GetSegment())
		brx := uint(level + 16)
		br := pkg.NewBaseRegister(
			address,
			pkg.NewAccessLock(0, 0),
			pkg.NewAccessPermissions(false, true, false),
			pkg.NewAccessPermissions(false, true, false),
			0,
			0,
			false,
			table)

		ute.engine.SetBaseRegister(brx, br)
		fmt.Printf("    Set B%d -> BDT for level %d\n", brx, level)
	}

	//	Now load the lower base registers according to the executable instructions.
	//	This is also a convenient place to set PAR.PC
	fmt.Printf("  Loading base registers and ABTEs with program bank information...\n")
	banks := ute.executable.GetBanks()
	for brx, bdi := range ute.executable.GetInitiallyBasedBanks() {
		bank, _ := banks[bdi]
		baseAddr := ute.bankAddresses[bdi]
		offset := baseAddr.GetOffset()
		limit := offset + bank.GetCodeLength()
		storage, _ := ute.storage.GetSlice(baseAddr.GetSegment(), offset, limit)
		br := pkg.NewBaseRegisterFromBankDescriptor(bank.GetBankDescriptor(), storage)
		ute.engine.SetBaseRegister(brx, br)
		fmt.Printf("    Set B%d -> %06o\n", brx, bdi)

		level := bdi >> 15
		index := bdi & 077777
		if brx == 0 {
			ute.engine.SetPARPC(level, index, ute.executable.GetStartingAddress())
		} else {
			abte := ute.engine.GetActiveBaseTableEntry(brx)
			abte.SetBankLevel(level).SetBankDescriptorIndex(index).SetSubsetSpecification(0)
			fmt.Printf("    Set ABTE[%d]\n", brx)
		}

	}

	//	Initialize designator register
	dr := ute.engine.GetDesignatorRegister()
	dr.Clear()
	dr.SetProcessorPrivilege(ute.executable.GetProcessorPrivilege())
	dr.arithmeticExceptionEnabled = ute.executable.IsArithmeticExceptionEnabled()
	dr.basicModeBaseRegisterSelection = ute.executable.GetBaseRegisterSelection()
	dr.basicModeEnabled = ute.executable.IsBasicMode()
	dr.executive24BitIndexingEnabled = ute.executable.IsExec24BitIndexingEnabled()
	dr.execRegisterSetSelected = ute.executable.IsExecRegisterSetEnabled()
	dr.operationTrapEnabled = ute.executable.IsOperationTrapEnabled()
	dr.quarterWordModeEnabled = ute.executable.IsQuarterWordMode()

	//	Reset GRS, clear interrupts and jump history
	//  TODO maybe we should have a Clear() on the engine?
	ute.engine.GetGeneralRegisterSet().Clear()
	ute.engine.ClearInterrupt()
	//	TODO

	//	Now Iterate until an interrupt is posted
	for !ute.engine.HasPendingInterrupt() {
		ute.engine.doCycle()
	}
	fmt.Printf("Execution Interrupted\n")

	ute.engine.Dump()
	return nil
}

func GetInterruptString(i pkg.Interrupt) string {
	return fmt.Sprintf("Class:%v SSF:%v ISW0:%012o ISW1:%012o",
		i.GetClass(), i.GetShortStatusField(), i.GetStatusWord0(), i.GetStatusWord1())
}
