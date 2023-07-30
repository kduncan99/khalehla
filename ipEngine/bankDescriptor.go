// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

// BankType values
const (
	ExtendedModeBankDescriptor    uint = 00
	BasicModeBankDescriptor       uint = 01
	GateBankDescriptor            uint = 02
	IndirectBankDescriptor        uint = 03
	QueueBankDescriptor           uint = 04
	PosternBankDescriptor         uint = 05
	QueueRepositoryBankDescriptor uint = 06
	DataExpanseBankDescriptor     uint = 07
)

type BankDescriptor struct {
	generalAccessPermissions     *pkg.AccessPermissions
	specialAccessPermissions     *pkg.AccessPermissions
	bankType                     uint
	generalFault                 bool
	largeBankSize                bool
	upperLimitSuppressionControl bool
	accessLock                   *pkg.AccessLock
	indirectLevelAndBDI          uint
	lowerLimit                   uint
	upperLimit                   uint
	inactiveFlag                 bool
	displacement                 uint
	baseAddress                  *AbsoluteAddress
	inactiveQBDListNextPointer   uint64
}

func (bd *BankDescriptor) GetLowerLimitNormalized() uint {
	if bd.largeBankSize {
		return bd.lowerLimit << 15
	} else {
		return bd.lowerLimit << 9
	}
}

func (bd *BankDescriptor) GetUpperLimitNormalized() uint {
	if bd.largeBankSize {
		return bd.upperLimit << 6
	} else {
		return bd.upperLimit
	}
}

func NewBankDescriptorFromStorage(buffer []pkg.Word36) *BankDescriptor {
	gap := pkg.NewAccessPermissions(
		buffer[0]&0_400000_000000 != 0,
		buffer[0]&0_200000_000000 != 0,
		buffer[0]&0_100000_000000 != 0)
	sap := pkg.NewAccessPermissions(
		buffer[0]&0_0400000_000000 != 0,
		buffer[0]&0_0200000_000000 != 0,
		buffer[0]&0_0100000_000000 != 0)
	typ := uint(buffer[0]>>24) & 0x0F
	gBit := buffer[0]&0_000020_000000 != 0
	sBit := buffer[0]&0_000004_000000 != 0
	uBit := buffer[0]&0_000002_000000 != 0
	lock := pkg.NewAccessLock(uint(buffer[0]>>16)&03, uint(buffer[0]&0xFFFF))

	ilBDI := uint(0)
	lLimit := uint(0)
	uLimit := uint(0)
	if typ == IndirectBankDescriptor {
		ilBDI = uint(buffer[1]>>18) & 0_777777
	} else {
		lLimit = uint(buffer[1]>>27) & 0777
		uLimit = uint(buffer[1] & 0_777777777)
	}

	ina := buffer[2].IsNegative()
	disp := uint(buffer[2]>>18) & 077777
	addr := NewAbsoluteAddressFromComposite(uint64(buffer[2]&0_777777) | uint64(buffer[3]))
	inQBD := uint64(buffer[3])

	bd := BankDescriptor{
		generalAccessPermissions:     gap,
		specialAccessPermissions:     sap,
		bankType:                     typ,
		generalFault:                 gBit,
		largeBankSize:                sBit,
		upperLimitSuppressionControl: uBit,
		accessLock:                   lock,
		indirectLevelAndBDI:          ilBDI,
		lowerLimit:                   lLimit,
		upperLimit:                   uLimit,
		inactiveFlag:                 ina,
		displacement:                 disp,
		baseAddress:                  addr,
		inactiveQBDListNextPointer:   inQBD,
	}

	return &bd
}
