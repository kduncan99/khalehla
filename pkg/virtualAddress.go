// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

type VirtualAddress interface {
	GetBankDescriptorIndex() uint64
	GetComposite() uint64
	GetOffset() uint64
	SetComposite(uint64)
}

type BasicModeVirtualAddress struct {
	execFlag            bool
	levelFlag           bool
	bankDescriptorIndex uint64
	offset              uint64 //	offset from start of bank
}

func (addr *BasicModeVirtualAddress) GetComposite() uint64 {
	var value uint64
	if addr.execFlag {
		value |= 0_400000_000000
	}
	if addr.levelFlag {
		value |= 0_040000_000000
	}
	value |= addr.bankDescriptorIndex << 18
	value |= addr.offset
	return value
}

func (addr *BasicModeVirtualAddress) GetExecFlag() bool {
	return addr.execFlag
}

func (addr *BasicModeVirtualAddress) GetLevelFlag() bool {
	return addr.levelFlag
}

func (addr *BasicModeVirtualAddress) GetBankDescriptorIndex() uint64 {
	return addr.bankDescriptorIndex
}

func (addr *BasicModeVirtualAddress) GetOffset() uint64 {
	return addr.offset
}

func (addr *BasicModeVirtualAddress) SetComposite(value uint64) {
	addr.execFlag = value&0_400000_000000 != 0
	addr.levelFlag = value&0_040000_000000 != 0
	addr.bankDescriptorIndex = (value >> 18) & 07777
	addr.offset = value & 0_777777
}

func (addr *BasicModeVirtualAddress) SetExecFlag(value bool) *BasicModeVirtualAddress {
	addr.execFlag = value
	return addr
}

func (addr *BasicModeVirtualAddress) SetLevelFlag(value bool) *BasicModeVirtualAddress {
	addr.levelFlag = value
	return addr
}

func (addr *BasicModeVirtualAddress) SetBankDescriptorIndex(value uint64) *BasicModeVirtualAddress {
	addr.bankDescriptorIndex = value & 07777
	return addr
}

func (addr *BasicModeVirtualAddress) SetOffset(value uint64) *BasicModeVirtualAddress {
	addr.offset = value & 0777777
	return addr
}

// TranslateToBasicMode translates extended mode semantic level, BDI, and offset to basic mode semantics
func TranslateToBasicMode(bankLevel uint64, bankDescriptorIndex uint64, offset uint64) *BasicModeVirtualAddress {
	var bdi uint64
	var execFlag bool
	var levelFlag bool
	if (bankDescriptorIndex >= 0) && (bankDescriptorIndex <= 07777) && (bankLevel%1 == 0) {
		bdi = bankDescriptorIndex & 07777
		execFlag = bankLevel&04 == 0
		levelFlag = (bankLevel&06 == 0) || (bankLevel == 6)
	} else {
		bdi = 0
		execFlag = true
		levelFlag = true
	}

	va := BasicModeVirtualAddress{
		execFlag:            execFlag,
		levelFlag:           levelFlag,
		bankDescriptorIndex: bdi,
		offset:              offset & 0_777777,
	}
	return &va
}

func NewBasicModeVirtualAddress(
	execFlag bool,
	levelFlag bool,
	bankDescriptorIndex uint64,
	offset uint64) *BasicModeVirtualAddress {

	addr := BasicModeVirtualAddress{}
	addr.SetExecFlag(execFlag).
		SetLevelFlag(levelFlag).
		SetBankDescriptorIndex(bankDescriptorIndex).
		SetOffset(offset)
	return &addr
}

type ExtendedModeVirtualAddress struct {
	level               uint64
	bankDescriptorIndex uint64
	offset              uint64 //	offset from start of bank
}

func (addr *ExtendedModeVirtualAddress) GetComposite() uint64 {
	value := addr.level << 33
	value |= addr.bankDescriptorIndex << 18
	value |= addr.offset
	return value
}

func (addr *ExtendedModeVirtualAddress) GetLevel() uint64 {
	return addr.level
}

func (addr *ExtendedModeVirtualAddress) GetBankDescriptorIndex() uint64 {
	return addr.bankDescriptorIndex
}

func (addr *ExtendedModeVirtualAddress) GetOffset() uint64 {
	return addr.offset
}

func (addr *ExtendedModeVirtualAddress) SetComposite(value uint64) {
	addr.level = (value >> 33) & 07
	addr.bankDescriptorIndex = (value >> 18) & 077777
	addr.offset = value & 0_777777
}

func (addr *ExtendedModeVirtualAddress) SetLevel(value uint64) *ExtendedModeVirtualAddress {
	addr.level = value & 07
	return addr
}

func (addr *ExtendedModeVirtualAddress) SetBankDescriptorIndex(value uint64) *ExtendedModeVirtualAddress {
	addr.bankDescriptorIndex = value & 077777
	return addr
}

func (addr *ExtendedModeVirtualAddress) SetOffset(value uint64) *ExtendedModeVirtualAddress {
	addr.offset = value & 0777777
	return addr
}

func NewExtendedModeVirtualAddress(level uint64, bankDescriptorIndex uint64, offset uint64) *ExtendedModeVirtualAddress {
	addr := ExtendedModeVirtualAddress{}
	addr.SetLevel(level).
		SetBankDescriptorIndex(bankDescriptorIndex).
		SetOffset(offset)
	return &addr
}
