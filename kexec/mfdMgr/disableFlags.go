// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// DisableFlags is a 6-bit field in main item sector 0 describing disable states of a file
// It is contained within one or more of the more-specific file cycle info structs.
type DisableFlags struct {
	DirectoryError                 bool
	AssignedAndWrittenAtSystemStop bool
	InaccessibleBackup             bool
	CacheDrainFailure              bool
}

func (dif *DisableFlags) Compose() uint64 {
	value := uint64(0)
	if dif.DirectoryError {
		value |= 0_60
	}
	if dif.AssignedAndWrittenAtSystemStop {
		value |= 0_50
	}
	if dif.InaccessibleBackup {
		value |= 0_44
	}
	if dif.CacheDrainFailure {
		value |= 0_42
	}
	return value
}

func (dif *DisableFlags) ExtractFrom(field uint64) {
	dif.DirectoryError = field&0_20 != 0
	dif.AssignedAndWrittenAtSystemStop = field&0_10 != 0
	dif.InaccessibleBackup = field&0_04 != 0
	dif.CacheDrainFailure = field&0_02 != 0
}

func ExtractNewDisableFlags(field uint64) *DisableFlags {
	dif := &DisableFlags{}
	dif.ExtractFrom(field)
	return dif
}
