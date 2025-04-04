// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/old/kexec"
)

// UnitSelectionIndicators is an 18-bit field found in fixed main item 0.
// It contains information regarding the location of the file text.
type UnitSelectionIndicators struct {
	CreatedViaDevicePlacement      bool
	CreatedViaControlUnitPlacement bool
	CreatedViaLogicalPlacement     bool
	MultipleDevices                bool            // file text is distributed across multiple devices
	InitialLDATIndex               kexec.LDATIndex // LDAT index of initially selected device
}

func (usi *UnitSelectionIndicators) Compose() uint64 {
	value := uint64(0)
	if usi.CreatedViaDevicePlacement {
		value |= 0400000
	}
	if usi.CreatedViaControlUnitPlacement {
		value |= 0200000
	}
	if usi.CreatedViaLogicalPlacement {
		value |= 0100000
	}
	value |= uint64(usi.InitialLDATIndex)
	return value
}

func (usi *UnitSelectionIndicators) ExtractFrom(field uint64) {
	usi.CreatedViaDevicePlacement = field&0400000 != 0
	usi.CreatedViaControlUnitPlacement = field&0200000 != 0
	usi.CreatedViaLogicalPlacement = field&0100000 != 0
	usi.MultipleDevices = field&0040000 != 0
	usi.InitialLDATIndex = kexec.LDATIndex(field & 07777)
}
