// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

// MFDDisableFlags is a 6-bit field in main item sector 0 describing disable states of a file
type MFDDisableFlags struct {
	directoryError                 bool
	assignedAndWrittenAtSystemStop bool
	inaccessibleBackup             bool
	cacheDrainFailure              bool
}

func (df *MFDDisableFlags) Compose() uint64 {
	value := uint64(0)
	if df.directoryError {
		value |= 0_60
	}
	if df.assignedAndWrittenAtSystemStop {
		value |= 0_50
	}
	if df.inaccessibleBackup {
		value |= 0_44
	}
	if df.cacheDrainFailure {
		value |= 0_42
	}
	return value
}

func (df *MFDDisableFlags) ExtractFrom(field uint64) {
	df.directoryError = field&0_20 != 0
	df.assignedAndWrittenAtSystemStop = field&0_10 != 0
	df.inaccessibleBackup = field&0_04 != 0
	df.cacheDrainFailure = field&0_02 != 0
}
