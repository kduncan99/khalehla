// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// DescriptorFlags is a 12-bit field found in main item sector 0.
// It contains certain flags regarding a particular file.
type DescriptorFlags struct {
	unloaded            bool
	backedUp            bool
	saveOnCheckpoint    bool
	toBeCataloged       bool
	isTapeFile          bool
	isRemovableDiskFile bool
	toBeWriteOnly       bool
	toBeReadOnly        bool
	toBeDropped         bool
}

func (df *DescriptorFlags) Compose() uint64 {
	value := uint64(0)
	if df.unloaded {
		value |= 0_4000
	}
	if df.backedUp {
		value |= 0_2000
	}
	if df.saveOnCheckpoint {
		value |= 0_1000
	}
	if df.toBeCataloged {
		value |= 0_0100
	}
	if df.isTapeFile {
		value |= 0_0040
	}
	if df.isRemovableDiskFile {
		value |= 0_0010
	}
	if df.toBeWriteOnly {
		value |= 0_0004
	}
	if df.toBeReadOnly {
		value |= 0_0002
	}
	if df.toBeDropped {
		value |= 0_0001
	}
	return value
}

func (df *DescriptorFlags) ExtractFrom(field uint64) {
	df.unloaded = field&0_4000 != 0
	df.backedUp = field&0_2000 != 0
	df.saveOnCheckpoint = field&0_1000 != 0
	df.toBeCataloged = field&0_0100 != 0
	df.isTapeFile = field&0_0040 != 0
	df.isRemovableDiskFile = field&0_0010 != 0
	df.toBeWriteOnly = field&0_0004 != 0
	df.toBeReadOnly = field&0_0002 != 0
	df.toBeDropped = field&0_0001 != 0
}
