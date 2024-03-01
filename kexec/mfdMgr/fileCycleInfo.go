// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/kexec"
)

// FileCycleInfo is an interface for the more-specific fixed, removable, and tape file cycle info structs
type FileCycleInfo interface {
	GetFileSetIdentifier() FileSetIdentifier
	GetFileCycleIdentifier() FileCycleIdentifier
	GetQualifier() string
	GetFilename() string
	GetProjectId() string
	GetAccountId() string
	GetAbsoluteFileCycle() uint
	GetAssignMnemonic() string
	setFileCycleIdentifier(fcIdentifier FileCycleIdentifier)
	setFileSetIdentifier(fsIdentifier FileSetIdentifier)
}

// FileCycleIdentifier is a unique opaque identifier allowing clients to refer to a file cycle
// without using qualifier, filename, and file cycle. Internally it is the main item sector 0 address
// for the file cycle - but clients should not be concerned with, nor rely on, that.
type FileCycleIdentifier uint64

// BackupInfo describes information from the file cycle pertaining to the backup state of the file.
// It is contained within one or more of the more-specific file cycle info structs.
// We do not (currently) support multiple backup levels, so don't look for that information.
type BackupInfo struct {
	TimeBackupCreated    uint64
	FASBits              uint64
	NumberOfTextBlocks   uint64
	StartingFilePosition uint64
	BackupReelNumbers    []string
}

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

// DisableFlags is a 6-bit field in main item sector 0 describing disable states of a file
// It is contained within one or more of the more-specific file cycle info structs.
type DisableFlags struct {
	DirectoryError                 bool
	AssignedAndWrittenAtSystemStop bool
	InaccessibleBackup             bool
	CacheDrainFailure              bool
}

// DiskPackEntry describes a particular disk pack as it is referred to by a particular file cycle entity.
// It is contained within one or more of the more-specific file cycle info structs.
type DiskPackEntry struct {
	PackName     string
	MainItemLink uint64
}

// FileFlags contain miscellaneous information about a file cycle
// It is contained within one or more of the more-specific file cycle info structs.
type FileFlags struct {
	IsLargeFile            bool
	AssignmentAcceleration bool
	IsWrittenTo            bool
	StoreThrough           bool
}

// InhibitFlags is a 6-bit field found in main item sector 0.
// It contains certain flags regarding inhibits which apply to a file.
type InhibitFlags struct {
	IsGuarded           bool
	IsUnloadInhibited   bool
	IsPrivate           bool
	isAssignedExclusive bool
	IsWriteOnly         bool
	IsReadOnly          bool
}

// PCHARFlags is a 6-bit field found in main item sector 0.
// It contains certain flags regarding inhibits which apply to a file.
type PCHARFlags struct {
	Granularity       kexec.Granularity
	IsWordAddressable bool
}

// UnitSelectionIndicators is an 18-bit field found in fixed main item 0.
// It contains information regarding the location of the file text.
type UnitSelectionIndicators struct {
	CreatedViaDevicePlacement      bool
	CreatedViaControlUnitPlacement bool
	CreatedViaLogicalPlacement     bool
	MultipleDevices                bool            // file text is distributed across multiple devices
	InitialLDATIndex               kexec.LDATIndex // LDAT index of initially selected device
}

// DescriptorFlags functions -----------------------------------------

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

// DisableFlags functions -----------------------------------------

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

// FileFlags functions -----------------------------------------

func (fif *FileFlags) Compose() uint64 {
	value := uint64(0)
	if fif.IsLargeFile {
		value |= 040
	}
	if fif.AssignmentAcceleration {
		value |= 004
	}
	if fif.IsWrittenTo {
		value |= 002
	}
	if fif.StoreThrough {
		value |= 001
	}
	return value
}

func (fif *FileFlags) ExtractFrom(field uint64) {
	fif.IsLargeFile = field&040 != 0
	fif.AssignmentAcceleration = field&004 != 0
	fif.IsWrittenTo = field&002 != 0
	fif.StoreThrough = field&001 != 0
}

// InhibitFlags functions -----------------------------------------

func (inf *InhibitFlags) Compose() uint64 {
	value := uint64(0)
	if inf.IsGuarded {
		value |= 040
	}
	if inf.IsUnloadInhibited {
		value |= 020
	}
	if inf.IsPrivate {
		value |= 010
	}
	if inf.isAssignedExclusive {
		value |= 004
	}
	if inf.IsWriteOnly {
		value |= 002
	}
	if inf.IsReadOnly {
		value |= 001
	}
	return value
}

func (inf *InhibitFlags) ExtractFrom(field uint64) {
	inf.IsGuarded = field&040 != 0
	inf.IsUnloadInhibited = field&020 != 0
	inf.IsPrivate = field&010 != 0
	inf.isAssignedExclusive = field&004 != 0
	inf.IsWriteOnly = field&002 != 0
	inf.IsReadOnly = field&001 != 0
}

// PCHARFlags functions -----------------------------------------

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

// UnitSelectionIndicators functions -----------------------------------------

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
