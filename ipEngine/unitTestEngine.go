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
	bankAddresses map[uint64]*pkg.AbsoluteAddress

	//	BDTable lookup into storage - key is level (0 to 7), and value is the absolute address of the table in storage
	//  If there is not an entry for a level, then there is not a BDT for that level.
	//  For now, we have one main storage segment for each BDT, so we can easily tell the size of the BDT from the
	//  size of the segment, and the offset of the table (from the start of the segment) is always zero.
	bankDescriptorTableAddresses map[uint64]*pkg.AbsoluteAddress

	engine *InstructionEngine
}

func NewUnitTestExecutor() *UnitTestEngine {
	e := &UnitTestEngine{
		storage: pkg.NewMainStorage(100),
	}

	return e
}

func (ute *UnitTestEngine) Clear() {
	ute.storage.Clear()
	ute.bankAddresses = make(map[uint64]*pkg.AbsoluteAddress)
	ute.bankDescriptorTableAddresses = make(map[uint64]*pkg.AbsoluteAddress)
	ute.executable = nil
	ute.engine = nil
}

func (ute *UnitTestEngine) GetEngine() *InstructionEngine {
	return ute.engine
}

func (ute *UnitTestEngine) Load(executable *tasm.Executable) error {
	fmt.Printf("\nLoading Executable...\n")
	ute.Clear()
	ute.executable = executable

	// maps bdi to segment index of segment containing the bank
	segIndexMap := make(map[uint64]uint64)

	for _, bank := range executable.GetBanks() {
		//	allocate a segment from storage and copy the bank to the segment
		//  there are no fix-ups required; all our binaries are self-contained.

		lbdi := bank.GetBankDescriptorIndex()
		bankSegIndex, err := ute.storage.Allocate(bank.GetCodeLength())
		if err != nil {
			return err
		}
		seg, _ := ute.storage.GetSegment(bankSegIndex)
		segIndexMap[lbdi] = bankSegIndex

		code := bank.GetCode()
		for cx := 0; cx < len(code); cx++ {
			seg[cx].SetW(code[cx])
		}

		bankAddress := pkg.NewAbsoluteAddress(bankSegIndex, 0)
		ute.bankAddresses[lbdi] = bankAddress

		//	now build a bank descriptor for the bank.
		//	if this is the first bank with its level, then we have to allocate space for a bdt for the level.
		level := lbdi >> 15
		bdi := lbdi & 077777
		newBDTLen := (bdi + 1) * 8

		//	For this algorithm, each BDT gets its own segment and thus is always at offset 0 of the segment
		absAddr, ok := ute.bankDescriptorTableAddresses[level]
		var bdtSegment uint64
		var bdTable []pkg.Word36
		if ok {
			bdtSegment := absAddr.GetSegment()
			bdTable, _ = ute.storage.GetSegment(bdtSegment)
			if uint64(len(bdTable)) < newBDTLen {
				interrupt := ute.storage.Resize(bdtSegment, newBDTLen)
				if interrupt != nil {
					return fmt.Errorf("interrupt:%s\n", pkg.GetInterruptString(interrupt))
				}
				bdTable, _ = ute.storage.GetSegment(bdtSegment)
			}
		} else {
			bdtSegment, err = ute.storage.Allocate(newBDTLen)
			if err != nil {
				return err
			}
			bdTable, _ = ute.storage.GetSegment(bdtSegment)
			absAddr = &pkg.AbsoluteAddress{}
			absAddr.SetSegment(bdtSegment)
			ute.bankDescriptorTableAddresses[level] = absAddr
		}
		bank.GetBankDescriptor().SetBaseAddress(bankAddress)

		bdOffset := bdi * 8
		bank.GetBankDescriptor().Serialize(bdTable[bdOffset : bdOffset+8])
	}

	//	Set up base registers
	ute.engine = NewEngine("IPTEST", ute.storage, NewStorageLocks())
	for brx := uint64(0); brx < 16; brx++ {
		ute.engine.SetBaseRegister(brx, pkg.NewVoidBaseRegister())
	}

	//	Load BDT base registers B16 to B23. BDTable level 0 -> B16, 1 -> B17, etc.
	//  For any level which does not have a BDT, we make the corresponding base register void.
	fmt.Printf("  Loading base registers for bank descriptor tables...\n")
	for level, address := range ute.bankDescriptorTableAddresses {
		table, _ := ute.storage.GetSegment(address.GetSegment())
		brx := level + 16
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
	for brx, bdi := range ute.executable.GetInitiallyBasedBanks() {
		//	Find the bank descriptor (in memory) for this bank
		level := bdi >> 15
		index := bdi & 077777
		bdAddr := ute.bankDescriptorTableAddresses[level]
		bdOffset := index * 8
		bdSlice, interrupt := ute.storage.GetSlice(bdAddr.GetSegment(), bdAddr.GetOffset()+bdOffset, 8)
		if interrupt != nil {
			return fmt.Errorf("interrupt:%s", pkg.GetInterruptString(interrupt))
		}

		//	Create an in-memory bank descriptor based on the bd in memory,
		//	then we're ready to create the bank register.
		bd := pkg.NewBankDescriptorFromStorage(bdSlice)

		bankAddr := bd.GetBaseAddress()
		bankSlice, interrupt := ute.storage.GetSegment(bankAddr.GetSegment())
		if interrupt != nil {
			return fmt.Errorf("interrupt:%s", pkg.GetInterruptString(interrupt))
		}

		br := pkg.NewBaseRegisterFromBankDescriptor(bd, bankSlice)
		ute.engine.SetBaseRegister(brx, br)
		fmt.Printf("    Set B%d -> %06o\n", brx, bdi)

		//	Update PAR.BDI or an appropriate ABTE. We are not expecting to handle large banks.
		par := ute.engine.GetProgramAddressRegister()
		if brx == 0 {
			par.SetLevel(level).SetBankDescriptorIndex(index)
		} else {
			abte := ute.engine.GetActiveBaseTableEntry(brx)
			abte.SetBankLevel(level).SetBankDescriptorIndex(index).SetSubsetSpecification(0)
			fmt.Printf("    Set ABTE[%d]\n", brx)
		}

		//	Finally, set program counter to initial address
		par.SetProgramCounter(executable.GetStartingAddress())
	}

	//	Initialize designator register
	dr := ute.engine.GetDesignatorRegister()
	dr.Clear()
	dr.SetProcessorPrivilege(ute.executable.GetProcessorPrivilege())
	dr.SetArithmeticExceptionEnabled(ute.executable.IsArithmeticExceptionEnabled())
	dr.SetBasicModeBaseRegisterSelection(ute.executable.GetBaseRegisterSelection())
	dr.SetBasicModeEnabled(ute.executable.IsBasicMode())
	dr.SetExecutive24BitIndexingEnabled(ute.executable.IsExec24BitIndexingEnabled())
	dr.SetExecRegisterSetSelected(ute.executable.IsExecRegisterSetEnabled())
	dr.SetOperationTrapEnabled(ute.executable.IsOperationTrapEnabled())
	dr.SetQuarterWordModeEnabled(ute.executable.IsQuarterWordMode())

	ute.engine.SetLogInstructions(true)
	ute.engine.SetLogInterrupts(true)
	return nil
}

func (ute *UnitTestEngine) Run() error {
	fmt.Printf("\nSetting up to run...\n")
	if ute.executable == nil {
		return fmt.Errorf("no executable has been loaded")
	}

	//	Reset GRS, clear interrupts and jump history
	ute.engine.GetGeneralRegisterSet().Clear()
	ute.engine.ClearStop()
	ute.engine.ClearAllInterrupts()

	//	TODO clear jump history

	//	Now Iterate until an interrupt is posted
	for !ute.engine.HasPendingInterrupt() && !ute.engine.IsStopped() {
		ute.engine.doCycle()
	}

	if ute.engine.HasPendingInterrupt() {
		fmt.Printf("Execution Interrupted\n")
	} else {
		fmt.Printf("Processor Stopped\n")
	}

	return nil
}
