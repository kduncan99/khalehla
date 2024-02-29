// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// DisableFlags is a 6-bit field in main item sector 0 describing disable states of a file
type DisableFlags struct {
	DirectoryError                 bool
	AssignedAndWrittenAtSystemStop bool
	InaccessibleBackup             bool
	CacheDrainFailure              bool
}

func (df *DisableFlags) Compose() uint64 {
	value := uint64(0)
	if df.DirectoryError {
		value |= 0_60
	}
	if df.AssignedAndWrittenAtSystemStop {
		value |= 0_50
	}
	if df.InaccessibleBackup {
		value |= 0_44
	}
	if df.CacheDrainFailure {
		value |= 0_42
	}
	return value
}

func (df *DisableFlags) ExtractFrom(field uint64) {
	df.DirectoryError = field&0_20 != 0
	df.AssignedAndWrittenAtSystemStop = field&0_10 != 0
	df.InaccessibleBackup = field&0_04 != 0
	df.CacheDrainFailure = field&0_02 != 0
}
