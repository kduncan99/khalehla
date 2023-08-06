// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

type BaseRegister struct {
	//	ring and domain for the bank described by this base register
	accessLock *AccessLock

	//	physical location of the described bank
	baseAddress *AbsoluteAddress

	//	ERW permissions for access by a key of lower privilege
	generalAccessPermissions *AccessPermissions

	//	ERW permissions for access by a key of equal or greater privilege
	specialAccessPermissions *AccessPermissions

	//	If true, area does not exceed 2^24 bytes - if false, area does not exceed 2^18 bytes
	largeSizeFlag bool

	//	Relative address, lower limit - 24 bits significant
	//	Corresponds to the first word/value in the storage subset
	//	one-word-granularity normalized form of the lower limit,
	//	accounting for the large size flag
	lowerLimitNormalized uint64

	//	Relative address, upper limit - 24 bits significant
	//	One-word-granularity normalized form of the upper limit,
	//	accounting for the large size flag
	upperLimitNormalized uint64

	//	If true, this is a void bank (no storage)
	voidFlag bool

	//	slice representing the storage for this bank
	storage []Word36
}

// CheckAccessLimits verifies that the given relative address is within the limits defined
// by the lower and upper normalized limits.
// If successful, return nil - otherwise, returns an interrupt to be posted
func (reg *BaseRegister) CheckAccessLimits(relativeAddress uint64, fetchFlag bool) Interrupt {
	if (relativeAddress < reg.lowerLimitNormalized) ||
		(relativeAddress > reg.upperLimitNormalized) {
		return NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, fetchFlag)
	} else {
		return nil
	}
}

// ConvertRelativeAddress converts a relative address to a new absolute address in the context of this base register
func (reg *BaseRegister) ConvertRelativeAddress(relAddr uint64) *AbsoluteAddress {
	actualOffset := relAddr - reg.lowerLimitNormalized
	offset := reg.baseAddress.GetOffset() + actualOffset
	return NewAbsoluteAddress(reg.baseAddress.GetSegment(), offset)
}

func (reg *BaseRegister) GetBaseAddress() *AbsoluteAddress {
	return reg.baseAddress
}

// GetEffectivePermissions returns either special or general access permissions, depending upon the
// combination of the lock for this base register, and the given key.
func (reg *BaseRegister) GetEffectivePermissions(key *AccessKey) *AccessPermissions {
	return reg.accessLock.GetEffectivePermissions(key, reg.specialAccessPermissions, reg.generalAccessPermissions)
}

// GetLowerLimitAdjusted returns the real lower limit shifted according to the size flag
func (reg *BaseRegister) GetLowerLimitAdjusted() uint64 {
	if reg.largeSizeFlag {
		return reg.lowerLimitNormalized >> 15
	} else {
		return reg.lowerLimitNormalized >> 9
	}
}

// GetLowerLimitNormalized returns the real upper limit for the bank
func (reg *BaseRegister) GetLowerLimitNormalized() uint64 {
	return reg.lowerLimitNormalized
}

func (reg *BaseRegister) GetStorage() []Word36 {
	return reg.storage
}

// GetUpperLimitAdjusted returns the real upper limit shifted according to the size flag
func (reg *BaseRegister) GetUpperLimitAdjusted() uint64 {
	if reg.largeSizeFlag {
		return reg.upperLimitNormalized >> 6
	} else {
		return reg.upperLimitNormalized
	}
}

// GetUpperLimitNormalized returns the real upper limit for the bank
func (reg *BaseRegister) GetUpperLimitNormalized() uint64 {
	return reg.upperLimitNormalized
}

func (reg *BaseRegister) IsLargeSize() bool {
	return reg.largeSizeFlag
}

func (reg *BaseRegister) IsVoid() bool {
	return reg.voidFlag
}

// FromBankDescriptor loads the fields of a BaseRegister struct based on the contents of the
// given bank descriptor, and the given storage slice.
// Needs to be a method so we can avoid alloc/dealloc lots of these things.
func (reg *BaseRegister) FromBankDescriptor(bd *BankDescriptor, storage []Word36) {
	reg.accessLock = bd.GetAccessLock()
	reg.baseAddress = bd.GetBaseAddress()
	reg.generalAccessPermissions = bd.GetGeneralAccessPermissions()
	reg.specialAccessPermissions = bd.GetSpecialAccessPermissions()
	reg.largeSizeFlag = bd.IsLargeBank()
	reg.lowerLimitNormalized = bd.GetLowerLimitNormalized()
	reg.upperLimitNormalized = bd.GetUpperLimitNormalized()
	reg.voidFlag = false
	reg.storage = storage
}

// FromBankDescriptorWithSubsetting loads the fields of a base register from the given bank descriptor,
// using the given offset for subsetting. We get into this mess when the caller wishes to access a bank larger
// than the D-field allows, by accessing consecutive sections of said bank by basing those segments on consecutive
// base registers.
// In this case, we add the given offset to the base offset from the BD, and adjust the lower and upper
// limits accordingly.  Subsequent accesses proceed as desired by virtue of the fact that we've set
// the base address in the bank register, along with the limits, in this fashion.
// Needs to be a method so we can avoid alloc/dealloc lots of these things.
func (reg *BaseRegister) FromBankDescriptorWithSubsetting(bd *BankDescriptor, offset uint64, storage []Word36) {
	reg.accessLock = bd.GetAccessLock()
	reg.baseAddress = bd.GetBaseAddress()
	reg.generalAccessPermissions =
		NewAccessPermissions(false,
			bd.GetGeneralAccessPermissions().CanRead(),
			bd.GetGeneralAccessPermissions().CanWrite())
	reg.specialAccessPermissions =
		NewAccessPermissions(false,
			bd.GetSpecialAccessPermissions().CanRead(),
			bd.GetSpecialAccessPermissions().CanWrite())
	reg.largeSizeFlag = bd.IsLargeBank()

	reg.lowerLimitNormalized = 0
	bdLowerNorm := bd.GetLowerLimitNormalized()
	if bdLowerNorm > offset {
		reg.lowerLimitNormalized = bdLowerNorm - offset
	}

	reg.upperLimitNormalized = bd.GetUpperLimitNormalized() - offset
	reg.voidFlag = (reg.upperLimitNormalized < 0) || (reg.lowerLimitNormalized > reg.upperLimitNormalized)
	reg.storage = storage
}

func (reg *BaseRegister) MakeVoid() {
	reg.accessLock.SetDomain(0).SetRing(0)
	reg.baseAddress.SetSegment(0).SetOffset(0)
	reg.generalAccessPermissions.Clear()
	reg.specialAccessPermissions.Clear()
	reg.largeSizeFlag = false
	reg.lowerLimitNormalized = 0
	reg.upperLimitNormalized = 0
	reg.voidFlag = true
	reg.storage = nil
}

// NewBaseRegisterFromBankDescriptor is a convenience wrapper for the method above
func NewBaseRegisterFromBankDescriptor(bd *BankDescriptor, storage []Word36) *BaseRegister {
	br := BaseRegister{}
	br.FromBankDescriptor(bd, storage)
	return &br
}

// NewBaseRegisterFromBuffer produces a struct given the 4-word slice buffer,
// and given a slice which represents the content of the bank.
func NewBaseRegisterFromBuffer(buffer []Word36, storage []Word36) *BaseRegister {
	sizeFlag := buffer[0]&0_000004_000000 != 0
	llNormal := uint64(buffer[1] >> 27)
	ulNormal := uint64(buffer[1] & 0777777)
	if sizeFlag {
		llNormal <<= 15
	} else {
		llNormal <<= 9
		ulNormal <<= 6
	}

	reg := BaseRegister{
		accessLock:  NewAccessLock((buffer[0].GetW()>>16)&03, buffer[0].GetW()&0xFFFF),
		baseAddress: NewAbsoluteAddress(0, 0).SetCompositeFromWord36(buffer[2:4]),
		generalAccessPermissions: NewAccessPermissions(
			false,
			buffer[0]&0_200000_000000 != 0,
			buffer[0]&0_100000_000000 != 0),
		specialAccessPermissions: NewAccessPermissions(
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

func NewBaseRegister(
	baseAddr *AbsoluteAddress,
	lock *AccessLock,
	gap *AccessPermissions,
	sap *AccessPermissions,
	lowerNormal uint64,
	upperNormal uint64,
	largeSize bool,
	storage []Word36) *BaseRegister {

	return &BaseRegister{
		baseAddress:              baseAddr,
		accessLock:               lock,
		generalAccessPermissions: gap,
		specialAccessPermissions: sap,
		lowerLimitNormalized:     lowerNormal,
		upperLimitNormalized:     upperNormal,
		largeSizeFlag:            largeSize,
		voidFlag:                 false,
		storage:                  storage,
	}
}
