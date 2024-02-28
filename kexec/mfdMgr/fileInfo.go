// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/kexec"
)

type DiskPackEntry struct {
	packName     string
	mainItemLink uint64
}

type FileInfo interface {
	GetAccountId() string
	GetAbsoluteFileCycle() uint
	GetAssignMnemonic() string
}

// -----------------------------------------------------------------------------

type FixedFileInfo struct {
	accountId                string
	AbsoluteFileCycle        uint
	TimeOfFirstWriteOrUnload uint64
	DescriptorFlags          DescriptorFlags
	WrittenTo                bool
	Granularity              kexec.Granularity
	WordAddressable          bool
	AssignMnemonic           string
	HasSmoqueEntry           bool
	NumberOfTimesAssigned    uint64
	InhibitFlags             InhibitFlags
	TimeOfLastReference      uint64
	TimeCataloged            uint64
	InitialGranulesReserved  uint64
	MaxGranules              uint64
	HighestGranuleAssigned   uint64
	HighestTrackWritten      uint64
	QuotaGroupGranules       []uint64
	BackupInfo               BackupInfo
	DiskPackEntries          []DiskPackEntry
	FileAllocations          []FileAllocation
}

func (fi *FixedFileInfo) GetAccountId() string {
	return fi.accountId
}

func (fi *FixedFileInfo) GetAbsoluteFileCycle() uint {
	return fi.AbsoluteFileCycle
}

func (fi *FixedFileInfo) GetAssignMnemonic() string {
	return fi.AssignMnemonic
}

// -----------------------------------------------------------------------------

type RemovableFileInfo struct {
	AccountId                string
	AbsoluteFileCycle        uint
	TimeOfFirstWriteOrUnload uint64
	DescriptorFlags          DescriptorFlags
	WrittenTo                bool
	Granularity              kexec.Granularity
	WordAddressable          bool
	AssignMnemonic           string
	HasSmoqueEntry           bool
	NumberOfTimesAssigned    uint64
	InhibitFlags             InhibitFlags
	TimeOfLastReference      uint64
	TimeCataloged            uint64
	InitialGranulesReserved  uint64
	MaxGranules              uint64
	HighestGranuleAssigned   uint64
	HighestTrackWritten      uint64
	ReadKey                  string
	WriteKey                 string
	QuotaGroupGranules       []uint64
	BackupInfo               BackupInfo
	DiskPackEntries          []DiskPackEntry
	FileAllocations          []FileAllocation
}

func (fi *RemovableFileInfo) GetAccountId() string {
	return fi.AccountId
}

func (fi *RemovableFileInfo) GetAbsoluteFileCycle() uint {
	return fi.AbsoluteFileCycle
}

func (fi *RemovableFileInfo) GetAssignMnemonic() string {
	return fi.AssignMnemonic
}

// -----------------------------------------------------------------------------

type TapeFileInfo struct {
	AccountId              string
	AbsoluteFileCycle      uint
	DescriptorFlags        DescriptorFlags
	AssignMnemonic         string
	NumberOfTimesAssigned  uint64
	InhibitFlags           InhibitFlags
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

func (fi *TapeFileInfo) GetAccountId() string {
	return fi.AccountId
}

func (fi *TapeFileInfo) GetAbsoluteFileCycle() uint {
	return fi.AbsoluteFileCycle
}

func (fi *TapeFileInfo) GetAssignMnemonic() string {
	return fi.AssignMnemonic
}
