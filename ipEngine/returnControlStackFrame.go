// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

type ReturnControlStackFrame struct {
	bankLevel             uint
	bankDescriptorIndex   uint
	offset                uint
	trapFlag              bool
	basicModeBaseRegister uint                // For return to basic mode, add 12 to this to find the base register for return
	designatorRegister    *DesignatorRegister //	only bits 12-17 are significant
	accessKey             *AccessKey
}

func (rcsf *ReturnControlStackFrame) SetBankLevel(value uint) *ReturnControlStackFrame {
	rcsf.bankLevel = value & 07
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetBankDescriptorIndex(value uint) *ReturnControlStackFrame {
	rcsf.bankDescriptorIndex = value & 077777
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetOffset(value uint) *ReturnControlStackFrame {
	rcsf.offset = value & 0_777777
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetTrapFlag(value bool) *ReturnControlStackFrame {
	rcsf.trapFlag = value
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetBasicModeBaseRegister(value uint) *ReturnControlStackFrame {
	rcsf.basicModeBaseRegister = value & 03
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetDesignatorRegister(value *DesignatorRegister) *ReturnControlStackFrame {
	rcsf.designatorRegister = value
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetAccessKey(value *AccessKey) *ReturnControlStackFrame {
	rcsf.accessKey = value
	return rcsf
}

func (rcsf *ReturnControlStackFrame) WriteToBuffer(buffer []pkg.Word36) {
	w0 := uint64(rcsf.bankLevel) << 34
	w0 |= uint64(rcsf.bankDescriptorIndex) << 18
	w0 |= uint64(rcsf.offset)
	buffer[0] = pkg.Word36(w0)

	w1 := uint64(0)
	if rcsf.trapFlag {
		w1 |= 0_400000_000000
	}
	w1 |= uint64(rcsf.basicModeBaseRegister << 24)
	w1 |= uint64(rcsf.bankDescriptorIndex & 0_000077_000000)
	w1 |= uint64(rcsf.accessKey.GetComposite())
	buffer[1] = pkg.Word36(w1)
}

func NewReturnControlStackFrameFromComponents(
	bankLevel uint,
	bankDescriptorIndex uint,
	offset uint,
	trapFlag bool,
	basicModeBaseRegister uint,
	designatorRegister *DesignatorRegister,
	accessKey *AccessKey) *ReturnControlStackFrame {
	rcsf := ReturnControlStackFrame{}
	rcsf.SetBankLevel(bankLevel).
		SetBankDescriptorIndex(bankDescriptorIndex).
		SetOffset(offset).
		SetTrapFlag(trapFlag).
		SetBasicModeBaseRegister(basicModeBaseRegister).
		SetDesignatorRegister(designatorRegister).
		SetAccessKey(accessKey)
	return &rcsf
}

func NewReturnControlStackFrameFromBuffer(source []pkg.Word36) *ReturnControlStackFrame {
	return &ReturnControlStackFrame{
		bankLevel:             uint(source[0] >> 33),
		bankDescriptorIndex:   uint(source[0]>>18) & 077777,
		offset:                uint(source[0] & 0777777),
		trapFlag:              (source[1] >> 35) == 01,
		basicModeBaseRegister: uint(source[1]>>24) & 03,
		designatorRegister:    NewDesignatorRegisterFromComposite(uint64(source[1] & 0_000077_000000)),
		accessKey:             NewAccessKeyFromComposite(uint(source[1])),
	}
}
