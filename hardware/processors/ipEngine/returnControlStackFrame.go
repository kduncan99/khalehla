// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/common"
)

type ReturnControlStackFrame struct {
	bankLevel             uint64
	bankDescriptorIndex   uint64
	offset                uint64
	trapFlag              bool
	basicModeBaseRegister uint64                     // For return to basic mode, add 12 to this to find the base register for return
	designatorRegister    *common.DesignatorRegister //	only bits 12-17 are significant
	accessKey             *common.AccessKey
}

func (rcsf *ReturnControlStackFrame) SetBankLevel(value uint64) *ReturnControlStackFrame {
	rcsf.bankLevel = value & 07
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetBankDescriptorIndex(value uint64) *ReturnControlStackFrame {
	rcsf.bankDescriptorIndex = value & 077777
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetOffset(value uint64) *ReturnControlStackFrame {
	rcsf.offset = value & 0_777777
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetTrapFlag(value bool) *ReturnControlStackFrame {
	rcsf.trapFlag = value
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetBasicModeBaseRegister(value uint64) *ReturnControlStackFrame {
	rcsf.basicModeBaseRegister = value & 03
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetDesignatorRegister(value *common.DesignatorRegister) *ReturnControlStackFrame {
	rcsf.designatorRegister = value
	return rcsf
}

func (rcsf *ReturnControlStackFrame) SetAccessKey(value *common.AccessKey) *ReturnControlStackFrame {
	rcsf.accessKey = value
	return rcsf
}

func (rcsf *ReturnControlStackFrame) WriteToBuffer(buffer []common.Word36) {
	w0 := uint64(rcsf.bankLevel) << 34
	w0 |= uint64(rcsf.bankDescriptorIndex) << 18
	w0 |= uint64(rcsf.offset)
	buffer[0] = common.Word36(w0)

	w1 := uint64(0)
	if rcsf.trapFlag {
		w1 |= 0_400000_000000
	}
	w1 |= uint64(rcsf.basicModeBaseRegister << 24)
	w1 |= uint64(rcsf.bankDescriptorIndex & 0_000077_000000)
	w1 |= uint64(rcsf.accessKey.GetComposite())
	buffer[1] = common.Word36(w1)
}

func NewReturnControlStackFrameFromComponents(
	bankLevel uint64,
	bankDescriptorIndex uint64,
	offset uint64,
	trapFlag bool,
	basicModeBaseRegister uint64,
	designatorRegister *common.DesignatorRegister,
	accessKey *common.AccessKey) *ReturnControlStackFrame {
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

func NewReturnControlStackFrameFromBuffer(source []common.Word36) *ReturnControlStackFrame {
	return &ReturnControlStackFrame{
		bankLevel:             uint64(source[0] >> 33),
		bankDescriptorIndex:   uint64(source[0]>>18) & 077777,
		offset:                uint64(source[0] & 0777777),
		trapFlag:              (source[1] >> 35) == 01,
		basicModeBaseRegister: uint64(source[1]>>24) & 03,
		designatorRegister:    common.NewDesignatorRegisterFromComposite(uint64(source[1] & 0_000077_000000)),
		accessKey:             common.NewAccessKeyFromComposite(uint64(source[1])),
	}
}
