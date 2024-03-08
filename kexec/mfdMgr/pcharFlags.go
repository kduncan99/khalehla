// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/kexec"
)

// PCHARFlags is a 6-bit field found in main item sector 0.
// It contains certain flags regarding inhibits which apply to a file.
type PCHARFlags struct {
	Granularity       kexec.Granularity
	IsWordAddressable bool
}

func (pcf *PCHARFlags) Compose() uint64 {
	value := uint64(0)
	if pcf.Granularity == kexec.PositionGranularity {
		value |= 040
	}
	if pcf.IsWordAddressable {
		value |= 010
	}
	return value
}

func (pcf *PCHARFlags) ExtractFrom(field uint64) {
	if field&040 != 0 {
		pcf.Granularity = kexec.PositionGranularity
	} else {
		pcf.Granularity = kexec.TrackGranularity
	}
	pcf.IsWordAddressable = field&010 != 0
}

func ExtractNewPCHARFlags(field uint64) *PCHARFlags {
	pcf := &PCHARFlags{}
	pcf.ExtractFrom(field)
	return pcf
}
