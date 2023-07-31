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
	banks []*Bank
}

func (e *Executable) GetBanks() []*Bank {
	return e.banks
}

// LinkSimple links the given segments into a single bank, all accessLock allowed, ring/domain == 0.
// the BDI for the bank will be 0_600004 (level 6, BDI 00004)
func (e *Executable) LinkSimple(segments map[int]*Segment) {
	fmt.Printf("\nLink Simple...\n")

	e.banks = make([]*Bank, 1)
	e.banks[0] = &Bank{
		accessLock:          pkg.NewAccessLock(0, 0),
		generalPermissions:  pkg.NewAccessPermissions(true, true, true),
		specialPermissions:  pkg.NewAccessPermissions(true, true, true),
		bankDescriptorIndex: 0_600004,
		lowerLimit:          01000,
	}

	//	Find the offsets of all the segments relative to the start of the bank
	//	key is the segment number, value is the offset
	offsets := make(map[int]uint64)
	var offset uint64
	var bankLength uint64
	for segmentNumber, segment := range segments {
		offsets[segmentNumber] = offset
		for _, codeBlock := range segment.generatedCode {
			blockLen := uint64(len(codeBlock.code))
			offset += blockLen
			bankLength += blockLen
		}
	}

	fmt.Printf("  Segment Table:\n")
	for segmentNumber, offset := range offsets {
		fmt.Printf("    Seg %03o is at offset %08o\n", segmentNumber, offset)
	}

	e.banks[0].code = make([]uint64, bankLength)

	//	Resolve undefined references for the segments
	resolved := make(map[string]uint64)
	for segmentNumber, segment := range segments {
		for symbol, offset := range segment.labels {
			//	offset is from the start of the segment -
			//  we need to also include the offset of the segment from the start of the bank,
			//  and the lower limit (base address) of the bank.
			resolved[symbol] = uint64(offset) + offsets[segmentNumber] + uint64(e.banks[0].lowerLimit)
		}
	}

	fmt.Printf("  References:\n")
	for symbol, value := range resolved {
		fmt.Printf("    %-12s: %012o\n", symbol, value)
	}

	//	Load code one segment at a time (unresolved)
	cx := 0
	for _, segment := range segments {
		for _, codeBlock := range segment.generatedCode {
			for _, code := range codeBlock.code {
				e.banks[0].code[cx] = code
				cx++
			}
		}
	}

	//	Now resolve references
	for segNumber, segment := range segments {
		segOffset := offsets[segNumber]
		for _, ref := range segment.references {
			newValue := resolved[strings.ToUpper(ref.symbol)]
			targetIndex := int(segOffset) + ref.offset
			baseValue := e.banks[0].code[targetIndex]
			var err error
			e.banks[0].code[targetIndex], err = addFractional(baseValue, newValue, ref.startingBit, ref.bitCount)
			if err != nil {
				fmt.Printf("E: BDI:%06o Offset:%012o: %s\n", e.banks[0].bankDescriptorIndex, targetIndex, err.Error())
			}
		}
	}
}

func addFractional(baseValue uint64, addend2 uint64, startingBit int, bitCount int) (uint64, error) {
	mask := uint64(1<<bitCount) - 1
	shift := 36 - startingBit - bitCount
	shiftedMask := uint64(mask << shift)
	shiftedNotMask := (^shiftedMask) & pkg.NegativeZero

	addend1 := (baseValue & shiftedMask) >> shift
	sum := addend1 + addend2
	if (sum & mask) != sum {
		return 0, fmt.Errorf("value %012o truncated startingBit:%v length:%v", sum, startingBit, bitCount)
	}

	shiftedSum := sum << shift
	return (baseValue & shiftedNotMask) | (shiftedSum & shiftedMask), nil
}

func (e *Executable) Show() {
	for _, bank := range e.banks {
		fmt.Printf("  Bank BDI:%06o  Access:%v  GAP %s  SAP %s  Lower:%012o\n",
			bank.bankDescriptorIndex, bank.accessLock.GetString(),
			bank.generalPermissions.GetString(), bank.specialPermissions.GetString(), bank.lowerLimit)
		addr := bank.lowerLimit
		for cx := 0; cx < len(bank.code); cx++ {
			fmt.Printf("    %08o: %012o\n", addr+uint(cx), bank.code[cx])
		}
	}
}
