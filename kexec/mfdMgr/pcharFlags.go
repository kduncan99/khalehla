// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import "khalehla/kexec"

// PCHARFlags is a 6-bit field found in main item sector 0.
// It contains certain flags regarding inhibits which apply to a file.
type PCHARFlags struct {
	Granularity       kexec.Granularity
	IsWordAddressable bool
}

func (pf *PCHARFlags) Compose() uint64 {
	value := uint64(0)
	if pf.Granularity == kexec.PositionGranularity {
		value |= 040
	}
	if pf.IsWordAddressable {
		value |= 010
	}
	return value
}

func (pf *PCHARFlags) ExtractFrom(field uint64) {
	if field&040 != 0 {
		pf.Granularity = kexec.PositionGranularity
	} else {
		pf.Granularity = kexec.TrackGranularity
	}
	pf.IsWordAddressable = field&010 != 0
}
