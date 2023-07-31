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

	//	BDT lookup into storage - key is level (0 to 7), and value is the absolute address of the table in storage
	//  If there is not an entry for a level, then there is not a BDT for that level.
	bankDescriptorTables map[int]*pkg.AbsoluteAddress

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
	ute.bankDescriptorTables = make(map[int]*pkg.AbsoluteAddress)
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

		//	now build a bank descriptor for the bank.
		//	if this is the first bank with its level, then we have to allocate space for a bdt for the level.
		bdiLevel := bdi >> 15
		bdiOffset := bdi & 077777
		newBDTLen := (bdiOffset + 1) << 3

		//	For this algorithm, each BDT gets its own segment and thus is always at offset 0 of the segment
		absAddr, ok := ute.bankDescriptorTables[int(bdiLevel)]
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
		}

		upperLimit := bank.GetLowerLimit() + bank.GetCodeLength() - 1
		bd := pkg.NewExtendedModeBankDescriptor(
			bank.GetAccessLock(), bank.GetGeneralPermissions(), bank.GetSpecialPermissions(),
			absAddr, false, bank.GetLowerLimit(), upperLimit, 0)
		bdOffset := bdiOffset * 8
		bd.Serialize(bdTable[bdOffset : bdOffset+8])
	}

	ute.storage.Dump() // TODO remove
	return nil
}

func (ute *UnitTestEngine) Run() error {
	if ute.executable == nil {
		return fmt.Errorf("no executable has been loaded")
	}

	ute.engine = NewEngine(ute.storage)

	//	Load BDT base registers B16 to B23. BDTable level 0 -> B16, 1 -> B17, etc.
	//  For any level which does not have a BDT, we make the corresponding base register void.
	// for bRegIndex := 16; bRegIndex < 24; bRegIndex++ {
	// 	level := bRegIndex - 16
	// 	bdi, ok := ute.bankDescriptorTables[level]
	// 	if ok {
	//
	// 	} else {
	//
	// 	}
	// }

	//	Now load the lower base registers according to the executable instructions.
	//	TODO

	//	Load PAR.PC and DR, reset GRS, and clear interrupts and jump history
	//	TODO

	//	Now Iterate until an interrupt is posted
	//	TODO

	return nil
}
