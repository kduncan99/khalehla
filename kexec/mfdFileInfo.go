// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type MFDDiskPackEntry struct {
	packName     string
	mainItemLink uint64
}

type MFDFileInfo interface {
	GetAccountId() string
	GetAbsoluteFileCycle() uint
	GetAssignMnemonic() string
}

// -----------------------------------------------------------------------------

type MFDFixedFileInfo struct {
	accountId                string
	AbsoluteFileCycle        uint
	TimeOfFirstWriteOrUnload uint64
	DescriptorFlags          MFDDescriptorFlags
	WrittenTo                bool
	Granularity              Granularity
	WordAddressable          bool
	AssignMnemonic           string
	HasSmoqueEntry           bool
	NumberOfTimesAssigned    uint64
	InhibitFlags             MFDInhibitFlags
	TimeOfLastReference      uint64
	TimeCataloged            uint64
	InitialGranulesReserved  uint64
	MaxGranules              uint64
	HighestGranuleAssigned   uint64
	HighestTrackWritten      uint64
	QuotaGroupGranules       []uint64
	BackupInfo               MFDBackupInfo
	DiskPackEntries          []MFDDiskPackEntry
	FileAllocations          []MFDFileAllocation
}

func (fi *MFDFixedFileInfo) GetAccountId() string {
	return fi.accountId
}

func (fi *MFDFixedFileInfo) GetAbsoluteFileCycle() uint {
	return fi.AbsoluteFileCycle
}

func (fi *MFDFixedFileInfo) GetAssignMnemonic() string {
	return fi.AssignMnemonic
}

// -----------------------------------------------------------------------------

type MFDRemovableFileInfo struct {
	AccountId                string
	AbsoluteFileCycle        uint
	TimeOfFirstWriteOrUnload uint64
	DescriptorFlags          MFDDescriptorFlags
	WrittenTo                bool
	Granularity              Granularity
	WordAddressable          bool
	AssignMnemonic           string
	HasSmoqueEntry           bool
	NumberOfTimesAssigned    uint64
	InhibitFlags             MFDInhibitFlags
	TimeOfLastReference      uint64
	TimeCataloged            uint64
	InitialGranulesReserved  uint64
	MaxGranules              uint64
	HighestGranuleAssigned   uint64
	HighestTrackWritten      uint64
	ReadKey                  string
	WriteKey                 string
	QuotaGroupGranules       []uint64
	BackupInfo               MFDBackupInfo
	DiskPackEntries          []MFDDiskPackEntry
	FileAllocations          []MFDFileAllocation
}

func (fi *MFDRemovableFileInfo) GetAccountId() string {
	return fi.AccountId
}

func (fi *MFDRemovableFileInfo) GetAbsoluteFileCycle() uint {
	return fi.AbsoluteFileCycle
}

func (fi *MFDRemovableFileInfo) GetAssignMnemonic() string {
	return fi.AssignMnemonic
}

// -----------------------------------------------------------------------------

type MFDTapeFileInfo struct {
	AccountId              string
	AbsoluteFileCycle      uint
	DescriptorFlags        MFDDescriptorFlags
	AssignMnemonic         string
	NumberOfTimesAssigned  uint64
	InhibitFlags           MFDInhibitFlags
	CurrentAssignCount     uint64
	TimeOfLastReference    uint64
	TimeCataloged          uint64
	Density                uint64
	Format                 uint64
	Features               uint64
	FeaturesExtension      uint64
	FeaturesExtension1     uint64
	NumberOfReelsCataloged uint64
	MtaPop                 uint64
	NoiseConstant          uint64
	TranslatorMnemonics    []string
	TapeLibraryPool        string
	ReelNumber             []string
}

func (fi *MFDTapeFileInfo) GetAccountId() string {
	return fi.AccountId
}

func (fi *MFDTapeFileInfo) GetAbsoluteFileCycle() uint {
	return fi.AbsoluteFileCycle
}

func (fi *MFDTapeFileInfo) GetAssignMnemonic() string {
	return fi.AssignMnemonic
}
