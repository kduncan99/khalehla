// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// DescriptorFlags is a 12-bit field found in main item sector 0.
// It contains certain flags regarding a particular file.
type DescriptorFlags struct {
	Unloaded            bool
	BackedUp            bool
	SaveOnCheckpoint    bool
	ToBeCataloged       bool
	IsTapeFile          bool
	IsRemovableDiskFile bool
	ToBeWriteOnly       bool
	ToBeReadOnly        bool
	ToBeDropped         bool
}

func (df *DescriptorFlags) Compose() uint64 {
	value := uint64(0)
	if df.Unloaded {
		value |= 0_4000
	}
	if df.BackedUp {
		value |= 0_2000
	}
	if df.SaveOnCheckpoint {
		value |= 0_1000
	}
	if df.ToBeCataloged {
		value |= 0_0100
	}
	if df.IsTapeFile {
		value |= 0_0040
	}
	if df.IsRemovableDiskFile {
		value |= 0_0010
	}
	if df.ToBeWriteOnly {
		value |= 0_0004
	}
	if df.ToBeReadOnly {
		value |= 0_0002
	}
	if df.ToBeDropped {
		value |= 0_0001
	}
	return value
}

func (df *DescriptorFlags) ExtractFrom(field uint64) {
	df.Unloaded = field&0_4000 != 0
	df.BackedUp = field&0_2000 != 0
	df.SaveOnCheckpoint = field&0_1000 != 0
	df.ToBeCataloged = field&0_0100 != 0
	df.IsTapeFile = field&0_0040 != 0
	df.IsRemovableDiskFile = field&0_0010 != 0
	df.ToBeWriteOnly = field&0_0004 != 0
	df.ToBeReadOnly = field&0_0002 != 0
	df.ToBeDropped = field&0_0001 != 0
}
