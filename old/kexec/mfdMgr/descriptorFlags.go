// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// DescriptorFlags is a 12-bit field found in main item sector 0.
// It contains certain flags regarding a particular file.
// It is contained within one or more of the more-specific file cycle info structs.
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

func (def *DescriptorFlags) Compose() uint64 {
	value := uint64(0)
	if def.Unloaded {
		value |= 0_4000
	}
	if def.BackedUp {
		value |= 0_2000
	}
	if def.SaveOnCheckpoint {
		value |= 0_1000
	}
	if def.ToBeCataloged {
		value |= 0_0100
	}
	if def.IsTapeFile {
		value |= 0_0040
	}
	if def.IsRemovableDiskFile {
		value |= 0_0010
	}
	if def.ToBeWriteOnly {
		value |= 0_0004
	}
	if def.ToBeReadOnly {
		value |= 0_0002
	}
	if def.ToBeDropped {
		value |= 0_0001
	}
	return value
}

func (def *DescriptorFlags) ExtractFrom(field uint64) {
	def.Unloaded = field&0_4000 != 0
	def.BackedUp = field&0_2000 != 0
	def.SaveOnCheckpoint = field&0_1000 != 0
	def.ToBeCataloged = field&0_0100 != 0
	def.IsTapeFile = field&0_0040 != 0
	def.IsRemovableDiskFile = field&0_0010 != 0
	def.ToBeWriteOnly = field&0_0004 != 0
	def.ToBeReadOnly = field&0_0002 != 0
	def.ToBeDropped = field&0_0001 != 0
}

func ExtractNewDescriptorFlags(field uint64) *DescriptorFlags {
	df := &DescriptorFlags{}
	df.ExtractFrom(field)
	return df
}
