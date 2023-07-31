// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

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
	generalAccessPermissions *AccessPermissions
	specialAccessPermissions *AccessPermissions
	bankType                 uint

	// If set, an addressing exception interrupt is raised if this BD is resolved via the bank manipulation alogirhtm
	generalFault bool

	// If false, the storage area is a single bank not exceeding 2^18 words,
	//	lowerLimit has a granularity of 512 words, and upperLimit has a granularity of 1 word.
	// If true, the storage area is a portion of a large bank not exceeding 2^24 words,
	//  or of a very large bank not exceeding 2^33 words.
	//  The lowerLimit has granularity of a granularity of 32768 words, and upperLimit has a granularity of 64 words.
	// If bank type is not extended, this must be false
	largeBankSize bool

	// Must be false for banks other than extended type
	upperLimitSuppressionControl bool

	accessLock                 *AccessLock
	indirectLevelAndBDI        uint
	lowerLimit                 uint
	upperLimit                 uint
	inactiveFlag               bool
	displacement               uint
	baseAddress                *AbsoluteAddress
	inactiveQBDListNextPointer uint64
}

func (bd *BankDescriptor) GetAccessLock() *AccessLock {
	return bd.accessLock
}

func (bd *BankDescriptor) GetBankType() uint {
	return bd.bankType
}

func (bd *BankDescriptor) GetBaseAddress() *AbsoluteAddress {
	return bd.baseAddress
}

func (bd *BankDescriptor) GetGeneralAccessPermissions() *AccessPermissions {
	return bd.generalAccessPermissions
}

func (bd *BankDescriptor) GetIndirectLevelAndBDI() uint {
	return bd.indirectLevelAndBDI
}

func (bd *BankDescriptor) GetLowerLimitNormalized() uint64 {
	if bd.largeBankSize {
		return uint64(bd.lowerLimit) << 15
	} else {
		return uint64(bd.lowerLimit) << 9
	}
}

func (bd *BankDescriptor) GetSpecialAccessPermissions() *AccessPermissions {
	return bd.specialAccessPermissions
}

func (bd *BankDescriptor) GetUpperLimitNormalized() uint64 {
	if bd.largeBankSize {
		return uint64(bd.upperLimit) << 6
	} else {
		return uint64(bd.upperLimit)
	}
}

func (bd *BankDescriptor) IsGeneralFault() bool {
	return bd.generalFault
}

func (bd *BankDescriptor) IsLargeBank() bool {
	return bd.largeBankSize
}

func (bd *BankDescriptor) SetBaseAddress(baseAddress *AbsoluteAddress) *BankDescriptor {
	bd.baseAddress = baseAddress
	return bd
}

func NewExtendedModeBankDescriptor(
	lock *AccessLock,
	general *AccessPermissions,
	special *AccessPermissions,
	baseAddress *AbsoluteAddress,
	largeBank bool,
	actualLowerLimit uint,
	actualUpperLimit uint,
	displacement uint) *BankDescriptor {

	bd := &BankDescriptor{}
	bd.bankType = ExtendedModeBankDescriptor
	bd.generalAccessPermissions = general
	bd.specialAccessPermissions = special
	bd.generalFault = false
	bd.largeBankSize = largeBank
	bd.upperLimitSuppressionControl = false
	bd.accessLock = lock

	bd.baseAddress = baseAddress

	ll := actualLowerLimit
	ul := actualUpperLimit
	if largeBank {
		ll >>= 15
		if actualLowerLimit&077777 != 0 {
			ll += 1
		}

		ul >>= 6
		if actualUpperLimit&077 != 0 {
			ul += 1
		}

	} else {
		ll >>= 9
		if actualLowerLimit&0777 != 0 {
			ll += 1
		}
	}
	bd.lowerLimit = ll
	bd.upperLimit = ul
	bd.inactiveFlag = false
	bd.inactiveQBDListNextPointer = 0

	bd.displacement = displacement

	return bd
}

func NewBankDescriptorFromStorage(buffer []Word36) *BankDescriptor {
	gap := NewAccessPermissions(
		buffer[0]&0_400000_000000 != 0,
		buffer[0]&0_200000_000000 != 0,
		buffer[0]&0_100000_000000 != 0)
	sap := NewAccessPermissions(
		buffer[0]&0_0400000_000000 != 0,
		buffer[0]&0_0200000_000000 != 0,
		buffer[0]&0_0100000_000000 != 0)
	typ := uint(buffer[0]>>24) & 0x0F
	gBit := buffer[0]&0_000020_000000 != 0
	sBit := buffer[0]&0_000004_000000 != 0
	uBit := buffer[0]&0_000002_000000 != 0
	lock := NewAccessLock(uint(buffer[0]>>16)&03, uint(buffer[0]&0xFFFF))

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

func (bd *BankDescriptor) Serialize(buffer []Word36) {
	var value0 uint64
	var value1 uint64
	var value2 uint64
	var value3 uint64
	var value4 uint64

	value0 |= uint64(bd.generalAccessPermissions.GetComposite()) << 33
	value0 |= uint64(bd.specialAccessPermissions.GetComposite()) << 30
	value0 |= uint64(bd.bankType) << 26
	if bd.generalFault {
		value0 |= 0_000020_000000
	}
	if bd.largeBankSize {
		value0 |= 0_000004_000000
	}
	if bd.upperLimitSuppressionControl {
		value0 |= 0_000002_000000
	}
	value0 |= uint64(bd.accessLock.GetComposite())

	if bd.bankType == IndirectBankDescriptor {
		value1 |= uint64(bd.indirectLevelAndBDI) << 18
	} else {
		value1 |= uint64(bd.lowerLimit) << 27
		value1 |= uint64(bd.upperLimit)
	}

	baseAddr := bd.baseAddress.GetComposite()
	value2 = baseAddr >> 36
	if bd.bankType == QueueBankDescriptor && bd.inactiveFlag {
		value3 = bd.inactiveQBDListNextPointer
	} else {
		value3 = baseAddr & NegativeZero
	}

	if bd.inactiveFlag {
		value4 |= 0_400000_000000
	}
	value4 |= uint64(bd.displacement&077777) << 18

	buffer[0].SetW(value0)
	buffer[1].SetW(value1)
	buffer[2].SetW(value2)
	buffer[3].SetW(value3)
	buffer[4].SetW(value4)
	buffer[5].SetW(0)
	buffer[4].SetW(0)
	buffer[5].SetW(0)
}
