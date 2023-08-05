// Khalehla Project
// tiny assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

import (
	"fmt"
	"khalehla/pkg"
	"strings"
)

type Executable struct {
	//	map of BDIs to the bank for that BDI
	banks map[uint64]*Bank

	//	map of BDIs to the base register upon which the bank should be registered at run time.
	//  key is base register index (0 to 15) and the value is the BDI of the bank.
	initiallyBasedBanks map[uint64]uint64

	//  stuff for setting the designator register
	arithmeticExceptionEnable bool
	baseRegisterSelection     bool
	basicMode                 bool
	execRegisterSet           bool
	exec24BitIndexing         bool
	operationTrapEnable       bool
	processorPrivilege        uint64
	quarterWordMode           bool
	startingAddress           uint64
}

func (e *Executable) GetBanks() map[uint64]*Bank {
	return e.banks
}

func (e *Executable) GetBaseRegisterSelection() bool {
	return e.baseRegisterSelection
}

func (e *Executable) GetInitiallyBasedBanks() map[uint64]uint64 {
	return e.initiallyBasedBanks
}

func (e *Executable) GetProcessorPrivilege() uint64 {
	return e.processorPrivilege

}
func (e *Executable) GetStartingAddress() uint64 {
	return e.startingAddress
}

func (e *Executable) IsArithmeticExceptionEnabled() bool {
	return e.arithmeticExceptionEnable
}

func (e *Executable) IsBasicMode() bool {
	return e.basicMode
}

func (e *Executable) IsExecRegisterSetEnabled() bool {
	return e.execRegisterSet
}

func (e *Executable) IsExec24BitIndexingEnabled() bool {
	return e.exec24BitIndexing
}

func (e *Executable) IsOperationTrapEnabled() bool {
	return e.operationTrapEnable
}

func (e *Executable) IsQuarterWordMode() bool {
	return e.quarterWordMode
}

// LinkSimple links the given segments into a single bank, all accessLock allowed, ring/domain == 0.
// the BDI for the bank will be 0_600004 (level 6, BDI 00004)
func (e *Executable) LinkSimple(segments map[uint64]*Segment, extendedMode bool) {
	fmt.Printf("\nLink Simple...\n")

	bdi := uint64(0_600004)
	e.banks = make(map[uint64]*Bank)
	e.initiallyBasedBanks = make(map[uint64]uint64)
	orderedSegmentNumbers := getOrderedSegmentNumbers(segments)

	//	Find the offsets of all the segments relative to the start of the bank
	//	key is the segment number, value is the offset
	offsets := make(map[uint64]uint64)
	var offset uint64
	var bankLength uint64
	for _, segmentNumber := range orderedSegmentNumbers {
		segment := segments[segmentNumber]
		offsets[segmentNumber] = offset
		for _, codeBlock := range segment.generatedCode {
			blockLen := uint64(len(codeBlock.code))
			offset += blockLen
			bankLength += blockLen
		}
	}

	fmt.Printf("  Segment Table:\n")
	for segmentNumber, offset := range offsets {
		fmt.Printf("    Seg %03o is in bank %06o at offset %08o\n", segmentNumber, bdi, offset)
	}

	bankCode := make([]uint64, bankLength)
	lowerLimit := uint64(01000)

	//	Resolve undefined references for the segments
	resolved := make(map[string]uint64)
	for segmentNumber, segment := range segments {
		for symbol, offset := range segment.labels {
			//	offset is from the start of the segment -
			//  we need to also include the offset of the segment from the start of the bank,
			//  and the lower limit (base address) of the bank.
			resolved[symbol] = offset + offsets[segmentNumber] + lowerLimit
		}
	}

	fmt.Printf("  Label Values:\n")
	for symbol, value := range resolved {
		fmt.Printf("    %-12s: %012o\n", symbol, value)
	}

	//	Load code one segment at a time (unresolved)
	for segmentNumber, segment := range segments {
		cx := offsets[segmentNumber]
		for _, codeBlock := range segment.generatedCode {
			for _, code := range codeBlock.code {
				bankCode[cx] = code
				cx++
			}
		}
	}

	//	Now resolve references
	for segNumber, segment := range segments {
		segOffset := offsets[segNumber]
		for _, ref := range segment.references {
			newValue := resolved[strings.ToUpper(ref.symbol)]
			targetIndex := segOffset + ref.offset
			baseValue := bankCode[targetIndex]
			var err error
			bankCode[targetIndex], err = addFractional(baseValue, newValue, ref.startingBit, ref.bitCount)
			if err != nil {
				fmt.Printf("E: BDI:%06o Offset:%012o: %s\n", bdi, targetIndex, err.Error())
			}
		}
	}

	bd := pkg.NewExtendedModeBankDescriptor(
		pkg.NewAccessLock(0, 0),
		pkg.NewAccessPermissions(true, true, true),
		pkg.NewAccessPermissions(true, true, true),
		nil, // this has to be filled in when the bank is loaded
		false,
		lowerLimit,
		lowerLimit+bankLength,
		0)
	e.banks[bdi] = &Bank{
		bankDescriptor:      bd,
		bankDescriptorIndex: bdi,
		code:                bankCode,
	}

	if extendedMode {
		e.initiallyBasedBanks[0] = bdi // the bank should be based on B0
	} else {
		e.initiallyBasedBanks[12] = bdi // the bank should be based on B12
	}
	e.startingAddress = 01000 // TODO pull this from .OPT command
}

// LinkBankPerSegment creates individual banks, one per input segment.
// The BDI will be 6010xx where xx is the segment number.
//
// For Extended Mode:
// The banks for segments 0 through 15 will be initially based on B0 through B15.
// This requires segment 0 to be the initial code bank.
// By convention, segment 1 will be the RCS stack, and segments 2 through 15 will be data banks.
// Segments do not have to be contiguous; you may have segments 0, 2, and 12, for example.
//
// For Basic Mode:
// The banks for segments 0 through 3 will be initially based on B12 through B15.
// The lower/upper limits for these banks will be set to avoid overlapping addresses, with B12 starting at 01000.
// Thus, we expect the initial code address to be in segment 12, at 01000.
// The lower/upper limits for non-initially-based banks will all be set to 01000.
func (e *Executable) LinkBankPerSegment(segments map[uint64]*Segment, extendedMode bool) {
	fmt.Printf("\nLink Bank-Per-Segment...\n")

	e.banks = make(map[uint64]*Bank)
	e.initiallyBasedBanks = make(map[uint64]uint64)
	orderedSegmentNumbers := getOrderedSegmentNumbers(segments)

	//  Create the bank descriptors here.
	//  We have to do this in segment number order so that we can properly place basic mode banks
	basicModeOffset := uint64(01000)
	for _, segmentNumber := range orderedSegmentNumbers {
		segment := segments[segmentNumber]
		lbdi := 0601000 + segmentNumber
		canEnter := (extendedMode && segmentNumber == 0) || !extendedMode
		canWrite := (extendedMode && segmentNumber != 0) || !extendedMode
		var lowerLimit uint64
		if extendedMode && segmentNumber == 0 {
			lowerLimit = 01000
		} else {
			lowerLimit = basicModeOffset
			basicModeOffset += segment.currentLength
		}

		bd := pkg.NewExtendedModeBankDescriptor(
			pkg.NewAccessLock(0, 0),
			pkg.NewAccessPermissions(canEnter, true, canWrite),
			pkg.NewAccessPermissions(canEnter, true, canWrite),
			nil, // this has to be filled in when the bank is loaded
			false,
			lowerLimit,
			lowerLimit+segment.currentLength,
			0)

		e.banks[lbdi] = &Bank{
			bankDescriptor:      bd,
			bankDescriptorIndex: lbdi,
			code:                make([]uint64, segment.currentLength),
		}

		if (extendedMode && segmentNumber < 16) || (!extendedMode && (segmentNumber >= 12 && segmentNumber <= 15)) {
			e.initiallyBasedBanks[segmentNumber] = lbdi
		}
	}

	//	Resolve undefined references for the segments
	resolved := make(map[string]uint64)
	for segmentNumber, segment := range segments {
		lbdi := 0601000 + segmentNumber
		bank := e.banks[lbdi]
		for symbol, offset := range segment.labels {
			resolved[symbol] = offset + bank.bankDescriptor.GetLowerLimitNormalized()
		}
	}

	fmt.Printf("  Label Values:\n")
	for symbol, value := range resolved {
		fmt.Printf("    %-12s: %012o\n", symbol, value)
	}

	fmt.Printf("  Bank Table:\n")
	fmt.Printf("    L,BDI  Lower  Upper        Seg\n")
	fmt.Printf("    ------ ------ ------------ ---\n")
	for lbdi, bank := range e.banks {
		fmt.Printf("    %06o %06o %012o %03o\n",
			lbdi,
			bank.bankDescriptor.GetLowerLimitNormalized(),
			bank.bankDescriptor.GetUpperLimitNormalized(),
			lbdi&0777)
	}

	//	Load code
	for segmentNumber, segment := range segments {
		lbdi := 0601000 + segmentNumber
		bank := e.banks[lbdi]
		cx := 0
		for _, codeBlock := range segment.generatedCode {
			for _, code := range codeBlock.code {
				bank.code[cx] = code
				cx++
			}
		}
	}

	//	Now resolve references
	for segmentNumber, segment := range segments {
		for _, ref := range segment.references {
			//	L,BDI and bank descriptor for the bank in which the reference exists
			lbdi := 0601000 + segmentNumber
			bank := e.banks[lbdi]

			newValue := resolved[strings.ToUpper(ref.symbol)]
			baseValue := bank.code[ref.offset]
			var err error
			bank.code[ref.offset], err = addFractional(baseValue, newValue, ref.startingBit, ref.bitCount)
			if err != nil {
				fmt.Printf("E: BDI:%06o Offset:%012o: %s\n", lbdi, ref.offset, err.Error())
			}
		}
	}

	e.startingAddress = 01000
}

func (e *Executable) Dump() {
	for _, bank := range e.banks {
		bd := bank.GetBankDescriptor()
		fmt.Printf("  Bank BDI:%06o  Access:%v  GAP %s  SAP %s  Lower:%012o\n",
			bank.bankDescriptorIndex,
			bd.GetAccessLock().GetString(),
			bd.GetGeneralAccessPermissions().GetString(),
			bd.GetSpecialAccessPermissions().GetString(),
			bd.GetLowerLimitNormalized())
		addr := bd.GetLowerLimitNormalized()
		for cx := 0; cx < len(bank.code); cx++ {
			fmt.Printf("    %08o: %012o\n", addr+uint64(cx), bank.code[cx])
		}
	}
}

func addFractional(baseValue uint64, addend2 uint64, startingBit uint64, bitCount uint64) (uint64, error) {
	mask := uint64(1<<bitCount) - 1
	shift := 36 - startingBit - bitCount
	shiftedMask := mask << shift
	shiftedNotMask := (^shiftedMask) & pkg.NegativeZero

	addend1 := (baseValue & shiftedMask) >> shift
	sum := addend1 + addend2
	if (sum & mask) != sum {
		return 0, fmt.Errorf("value %012o truncated startingBit:%v length:%v", sum, startingBit, bitCount)
	}

	shiftedSum := sum << shift
	return (baseValue & shiftedNotMask) | (shiftedSum & shiftedMask), nil
}

// getOrderedSegmentNumbers reates a list in ascending order, of the existing segment numbers
func getOrderedSegmentNumbers(segments map[uint64]*Segment) []uint64 {
	orderedSegmentNumbers := make([]uint64, 0)
	for segNum, _ := range segments {
		ox := 0
		done := false
		for ox < len(orderedSegmentNumbers) && !done {
			if segNum < orderedSegmentNumbers[ox] {
				orderedSegmentNumbers = append([]uint64{segNum}, orderedSegmentNumbers...)
				done = true
			} else {
				ox++
			}
		}
		if !done {
			orderedSegmentNumbers = append(orderedSegmentNumbers, segNum)
		}
	}

	return orderedSegmentNumbers
}
