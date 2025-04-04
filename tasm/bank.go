// khalehla Project
// tiny assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tasm

import (
	"khalehla/common"
)

type Bank struct {
	bankDescriptor      *common.BankDescriptor
	bankDescriptorIndex uint64
	code                []uint64 //	key is L,BDI (18 bits), value is the table of binary 36-bit values to comprise the bank
}

func (b *Bank) GetBankDescriptor() *common.BankDescriptor {
	return b.bankDescriptor
}

func (b *Bank) GetBankDescriptorIndex() uint64 {
	return b.bankDescriptorIndex
}

func (b *Bank) GetCode() []uint64 {
	return b.code
}

func (b *Bank) GetCodeLength() uint64 {
	return uint64(len(b.code))
}

func NewBank(bd *common.BankDescriptor, bdi uint64, code []uint64) *Bank {
	return &Bank{
		bankDescriptor:      bd,
		bankDescriptorIndex: bdi,
		code:                code,
	}
}
