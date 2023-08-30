// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

const (
	B0 = iota
	B1
	B2
	B3
	B4
	B5
	B6
	B7
	B8
	B9
	B10
	B11
	B12
	B13
	B14
	B15
	B16
	B17
	B18
	B19
	B20
	B21
	B22
	B23
	B24
	B25
	B26
	B27
	B28
	B29
	B30
	B31
)

type BaseRegister struct {
	storage        []Word36        // slice representing the storage for this bank - nil for void bank
	bankDescriptor *BankDescriptor // reference to BankDescriptor for this base register - nil for void bank
	subsetting     uint64          // offset from start of real bank
}

// CheckAccessLimits verifies that the given relative address is within the limits defined
// by the lower and upper normalized limits.
// If successful, return nil - otherwise, returns an interrupt to be posted
func (reg *BaseRegister) CheckAccessLimits(relativeAddress uint64, fetchFlag bool) (interrupt Interrupt) {
	interrupt = nil
	if reg.bankDescriptor == nil {
		interrupt = NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, fetchFlag)
	} else {
		if (relativeAddress < reg.bankDescriptor.GetLowerLimitNormalized()) ||
			(relativeAddress > reg.bankDescriptor.GetUpperLimitNormalized()) {
			interrupt = NewReferenceViolationInterrupt(ReferenceViolationStorageLimits, fetchFlag)
		}
	}
	return
}

// ConvertRelativeAddress converts a relative address to a new absolute address in the context of this base register
// TODO obsolete?
// func (reg *BaseRegister) ConvertRelativeAddress(relAddr uint64) *AbsoluteAddress {
// 	actualOffset := relAddr - reg.lowerLimitNormalized
// 	offset := reg.GetBaseAddress().GetOffset() + actualOffset
// 	return NewAbsoluteAddress(reg.GetBaseAddress().GetSegment(), offset)
// }

func (reg *BaseRegister) GetBankDescriptor() *BankDescriptor {
	return reg.bankDescriptor
}

// GetEffectivePermissions returns either special or general access permissions, depending upon the
// combination of the lock for this base register, and the given key.
func (reg *BaseRegister) GetEffectivePermissions(key *AccessKey) *AccessPermissions {
	lock := reg.bankDescriptor.GetAccessLock()
	spec := reg.bankDescriptor.GetSpecialAccessPermissions()
	gen := reg.bankDescriptor.GetGeneralAccessPermissions()
	return lock.GetEffectivePermissions(key, spec, gen)
}

func (reg *BaseRegister) GetStorage() []Word36 {
	return reg.storage
}

func (reg *BaseRegister) GetSubsetting() uint64 {
	return reg.subsetting
}

func (reg *BaseRegister) IsVoid() bool {
	return reg.bankDescriptor == nil
}

// FromBankDescriptor loads the fields of a BaseRegister struct based on the contents of the
// given bank descriptor, and the given storage slice.
// Needs to be a method so we can avoid alloc/dealloc lots of these things.
func (reg *BaseRegister) FromBankDescriptor(bankDescriptor *BankDescriptor, storage []Word36) {
	reg.bankDescriptor = bankDescriptor
	reg.storage = storage
	reg.subsetting = 0
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
	reg.bankDescriptor = bd
	reg.storage = storage
	reg.subsetting = offset
}

func (reg *BaseRegister) MakeVoid() {
	reg.bankDescriptor = nil
	reg.storage = nil
	reg.subsetting = 0
}

// PopulateAbsoluteAddress converts a relative address to an absolute address.
// The AbsoluteAddress is passed as a reference so that we can avoid leaving little structs all over the heap.
// func (reg *BaseRegister) PopulateAbsoluteAddress(relativeAddress uint64, addr *AbsoluteAddress) {
// 	addr.SetSegment(reg.GetBaseAddress().GetSegment())
// 	actualOffset := relativeAddress - reg.GetLowerLimitNormalized()
// 	addr.SetOffset(reg.GetBaseAddress().GetOffset() + actualOffset)
// }

// NewBaseRegisterFromBankDescriptor is a convenience wrapper for the method above
func NewBaseRegisterFromBankDescriptor(bd *BankDescriptor, storage []Word36) *BaseRegister {
	br := BaseRegister{}
	br.FromBankDescriptor(bd, storage)
	return &br
}

// NewBaseRegisterFromBankDescriptorWithSubsetting is a convenience wrapper for the method above
func NewBaseRegisterFromBankDescriptorWithSubsetting(bd *BankDescriptor, offset uint64, storage []Word36) *BaseRegister {
	br := BaseRegister{}
	br.FromBankDescriptorWithSubsetting(bd, offset, storage)
	return &br
}

// NewBaseRegisterFromBuffer produces a struct given the 4-word slice buffer,
// and given a slice which represents the content of the bank.
// func NewBaseRegisterFromBuffer(buffer []Word36, storage []Word36) *BaseRegister {
// 	sizeFlag := buffer[0]&0_000004_000000 != 0
// 	llNormal := uint64(buffer[1] >> 27)
// 	ulNormal := uint64(buffer[1] & 0777777)
// 	if sizeFlag {
// 		llNormal <<= 15
// 	} else {
// 		llNormal <<= 9
// 		ulNormal <<= 6
// 	}
//
// 	reg := BaseRegister{
// 		accessLock:  NewAccessLock((buffer[0].GetW()>>16)&03, buffer[0].GetW()&0xFFFF),
// 		baseAddress: NewAbsoluteAddress(0, 0).SetCompositeFromWord36(buffer[2:4]),
// 		generalAccessPermissions: NewAccessPermissions(
// 			false,
// 			buffer[0]&0_200000_000000 != 0,
// 			buffer[0]&0_100000_000000 != 0),
// 		specialAccessPermissions: NewAccessPermissions(
// 			false,
// 			buffer[0]&0_020000_000000 != 0,
// 			buffer[0]&0_020000_000000 != 0),
// 		largeSizeFlag:        sizeFlag,
// 		lowerLimitNormalized: llNormal,
// 		upperLimitNormalized: ulNormal,
// 		voidFlag:             buffer[0]&0_000200_000000 != 0,
// 		storage:              storage,
// 	}
//
// 	return &reg
// }

func NewVoidBaseRegister() *BaseRegister {
	reg := BaseRegister{}
	reg.MakeVoid()
	return &reg
}

// TODO obsolete
// func NewBaseRegister(
// 	baseAddr *AbsoluteAddress,
// 	lock *AccessLock,
// 	gap *AccessPermissions,
// 	sap *AccessPermissions,
// 	lowerNormal uint64,
// 	upperNormal uint64,
// 	largeSize bool,
// 	storage []Word36) *BaseRegister {
//
// 	return &BaseRegister{
// 		baseAddress:              baseAddr,
// 		accessLock:               lock,
// 		generalAccessPermissions: gap,
// 		specialAccessPermissions: sap,
// 		lowerLimitNormalized:     lowerNormal,
// 		upperLimitNormalized:     upperNormal,
// 		largeSizeFlag:            largeSize,
// 		voidFlag:                 false,
// 		storage:                  storage,
// 	}
// }
