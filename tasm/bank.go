// Khalehla Project
// tiny assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

import (
	"khalehla/pkg"
)

type Bank struct {
	accessLock          *pkg.AccessLock
	extendedMode        bool
	generalPermissions  *pkg.AccessPermissions
	specialPermissions  *pkg.AccessPermissions
	lowerLimit          uint
	bankDescriptorIndex uint
	code                []uint64 //	key is L,BDI (18 bits), value is the table of binary 36-bit values to comprise the bank

	//	List of initially-based banks
	//	key is the base register index (0 for B0, etc), value is the BDI of the bank to be based.
	initiallyBased map[int]uint
}

func (b *Bank) GetAccessLock() *pkg.AccessLock {
	return b.accessLock
}

func (b *Bank) GetBankDescriptorIndex() uint {
	return b.bankDescriptorIndex
}

func (b *Bank) GetCode() []uint64 {
	return b.code
}

func (b *Bank) GetCodeLength() uint {
	return uint(len(b.code))
}

func (b *Bank) GetGeneralPermissions() *pkg.AccessPermissions {
	return b.generalPermissions
}

func (b *Bank) GetInitiallyBasedMap() map[int]uint {
	return b.initiallyBased
}

func (b *Bank) GetLowerLimit() uint {
	return b.lowerLimit
}

func (b *Bank) GetSpecialPermissions() *pkg.AccessPermissions {
	return b.specialPermissions
}

func (b *Bank) IsExtendedMode() bool {
	return b.extendedMode
}
