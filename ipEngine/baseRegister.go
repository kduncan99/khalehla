// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

type BaseRegister struct {
	//	ring and domain for the bank described by this base register
	accessLock *pkg.AccessLock

	//	physical location of the described bank
	baseAddress *AbsoluteAddress

	//	ERW permissions for access by a key of lower privilege
	generalAccessPermissions *pkg.AccessPermissions

	//	ERW permissions for access by a key of equal or greater privilege
	specialAccessPermissions *pkg.AccessPermissions

	//	If true, area does not exceed 2^24 bytes - if false, area does not exceed 2^18 bytes
	largeSizeFlag bool

	//	Relative address, lower limit - 24 bits significant
	//	Corresponds to the first word/value in the storage subset
	//	one-word-granularity normalized form of the lower limit,
	//	accounting for the large size flag
	lowerLimitNormalized uint

	//	Relative address, upper limit - 24 bits significant
	//	One-word-granularity normalized form of the upper limit,
	//	accounting for the large size flag
	upperLimitNormalized uint

	//	If true, this is a void bank (no storage)
	voidFlag bool

	//	slice representing the storage for this bank
	storage []pkg.Word36
}

// CheckAccessLimits verifies that the given relative address is within the limits defined
// by the lower and upper normalized limits.
// If successful, return nil - otherwise, returns an interrupt to be posted
func (reg *BaseRegister) CheckAccessLimits(relativeAddress uint, fetchFlag bool) Interrupt {
	if (relativeAddress < reg.lowerLimitNormalized) || (relativeAddress > reg.upperLimitNormalized) {
		return NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, fetchFlag)
	} else {
		return nil
	}
}

// GetEffectivePermissions returns either special or general access permissions, depending upon the
// combination of the lock for this base register, and the given key.
func (reg *BaseRegister) GetEffectivePermissions(key *pkg.AccessKey) *pkg.AccessPermissions {
	return reg.accessLock.GetEffectivePermissions(key, reg.specialAccessPermissions, reg.generalAccessPermissions)
}

func (reg *BaseRegister) GetLowerLimitAdjusted() uint64 {
	if reg.largeSizeFlag {
		return uint64(reg.lowerLimitNormalized) << 15
	} else {
		return uint64(reg.lowerLimitNormalized) << 9
	}
}

func (reg *BaseRegister) GetUpperLimitAdjusted() uint64 {
	if reg.largeSizeFlag {
		return uint64(reg.upperLimitNormalized) << 6
	} else {
		return uint64(reg.upperLimitNormalized)
	}
}

// FromBankDescriptor loads the fields of a BaseRegister struct based on the contents of the
// given bank descriptor, and the given storage slice
func (reg *BaseRegister) FromBankDescriptor(bd *BankDescriptor, storage []pkg.Word36) *BaseRegister {
	reg.accessLock = bd.accessLock
	reg.baseAddress = bd.baseAddress
	reg.generalAccessPermissions = bd.generalAccessPermissions
	reg.specialAccessPermissions = bd.specialAccessPermissions
	reg.largeSizeFlag = bd.largeBankSize
	reg.lowerLimitNormalized = bd.GetLowerLimitNormalized()
	reg.upperLimitNormalized = bd.GetUpperLimitNormalized()
	reg.voidFlag = false
	reg.storage = storage
	return reg
}

// FromBankDescriptorWithSubsetting loads the fields of a base register from the given bank descriptor,
// using the given offset for subsetting. We get into this mess when the caller wishes to access a bank larger
// than the D-field allows, by accessing consecutive sections of said bank by basing those segments on consecutive
// base registers.
// In this case, we add the given offset to the base offset from the BD, and adjust the lower and upper
// limits accordingly.  Subsequent accesses proceed as desired by virtue of the fact that we've set
// the base address in the bank register, along with the limits, in this fashion.
func (reg *BaseRegister) FromBankDescriptorWithSubsetting(bd *BankDescriptor, offset uint, storage []pkg.Word36) {
	reg.accessLock = bd.accessLock
	reg.baseAddress = bd.baseAddress
	reg.generalAccessPermissions =
		pkg.NewAccessPermissions(false,
			bd.generalAccessPermissions.CanRead(),
			bd.generalAccessPermissions.CanWrite())
	reg.specialAccessPermissions =
		pkg.NewAccessPermissions(false,
			bd.specialAccessPermissions.CanRead(),
			bd.specialAccessPermissions.CanWrite())
	reg.largeSizeFlag = bd.largeBankSize

	reg.lowerLimitNormalized = 0
	bdLowerNorm := bd.GetLowerLimitNormalized()
	if bdLowerNorm > offset {
		reg.lowerLimitNormalized = bdLowerNorm - offset
	}

	reg.upperLimitNormalized = bd.GetUpperLimitNormalized() - offset
	reg.voidFlag = (reg.upperLimitNormalized < 0) || (reg.lowerLimitNormalized > reg.upperLimitNormalized)
	reg.storage = storage
}

// NewBaseRegisterFromBuffer produces a struct given the 4-word slice buffer,
// and given a slice which represents the content of the bank.
func NewBaseRegisterFromBuffer(buffer []pkg.Word36, storage []pkg.Word36) *BaseRegister {
	sizeFlag := buffer[0]&0_000004_000000 != 0
	llNormal := uint(buffer[1] >> 27)
	ulNormal := uint(buffer[1] & 0777777)
	if sizeFlag {
		llNormal <<= 15
	} else {
		llNormal <<= 9
		ulNormal <<= 6
	}

	reg := BaseRegister{
		accessLock:  pkg.NewAccessLock(uint(buffer[0]>>16)&03, uint(buffer[0]&0xFFFF)),
		baseAddress: NewAbsoluteAddressFromComposite(((uint64(buffer[2]) << 36) & 0777777) | uint64(buffer[3])),
		generalAccessPermissions: pkg.NewAccessPermissions(
			false,
			buffer[0]&0_200000_000000 != 0,
			buffer[0]&0_100000_000000 != 0),
		specialAccessPermissions: pkg.NewAccessPermissions(
			false,
			buffer[0]&0_020000_000000 != 0,
			buffer[0]&0_020000_000000 != 0),
		largeSizeFlag:        sizeFlag,
		lowerLimitNormalized: llNormal,
		upperLimitNormalized: ulNormal,
		voidFlag:             buffer[0]&0_000200_000000 != 0,
		storage:              storage,
	}

	return &reg
}

func NewVoidBaseRegister() *BaseRegister {
	reg := BaseRegister{
		voidFlag: true,
	}
	return &reg
}
