// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/pkg"
	"strings"
)

type DiskPackEntry struct {
	PackName     string
	MainItemLink uint64
}

type FileInfo interface {
	GetQualifier() string
	GetFilename() string
	GetProjectId() string
	GetAccountId() string
	GetAbsoluteFileCycle() uint
	GetAssignMnemonic() string
}

// -----------------------------------------------------------------------------

type FixedFileInfo struct {
	Qualifier                string
	Filename                 string
	ProjectId                string
	AccountId                string
	AbsoluteFileCycle        uint
	TimeOfFirstWriteOrUnload uint64
	DisableFlags             DisableFlags
	DescriptorFlags          DescriptorFlags
	FileFlags                FileFlags
	PCHARFlags               PCHARFlags
	AssignMnemonic           string
	InitialSmoqueLink        uint64
	NumberOfTimesAssigned    uint64
	InhibitFlags             InhibitFlags
	AssignedIndicator        bool
	AbsoluteFCycle           uint64
	TimeOfLastReference      uint64
	TimeCataloged            uint64
	InitialGranulesReserved  uint64
	MaxGranules              uint64
	HighestGranuleAssigned   uint64
	HighestTrackWritten      uint64
	// unit selection indicators (for fixed only) - yet another POGO
	QuotaGroupGranules []uint64
	BackupInfo         BackupInfo
	DiskPackEntries    []DiskPackEntry
	FileAllocations    []FileAllocation
}

func (fi *FixedFileInfo) GetQualifier() string {
	return fi.Qualifier
}

func (fi *FixedFileInfo) GetFilename() string {
	return fi.Filename
}

func (fi *FixedFileInfo) GetProjectId() string {
	return fi.ProjectId
}

func (fi *FixedFileInfo) GetAccountId() string {
	return fi.AccountId
}

func (fi *FixedFileInfo) GetAbsoluteFileCycle() uint {
	return fi.AbsoluteFileCycle
}

func (fi *FixedFileInfo) GetAssignMnemonic() string {
	return fi.AssignMnemonic
}

func (fi *FixedFileInfo) PopulateFromMainItems(mainItem0 []pkg.Word36, mainItem1 []pkg.Word36) {
	fi.Qualifier = strings.TrimRight(mainItem0[1].ToStringAsFieldata()+mainItem0[2].ToStringAsFieldata(), " ")
	fi.Filename = strings.TrimRight(mainItem0[3].ToStringAsFieldata()+mainItem0[4].ToStringAsFieldata(), " ")
	fi.ProjectId = strings.TrimRight(mainItem0[5].ToStringAsFieldata()+mainItem0[6].ToStringAsFieldata(), " ")
	fi.AccountId = strings.TrimRight(mainItem0[7].ToStringAsFieldata()+mainItem0[010].ToStringAsFieldata(), " ")
	fi.TimeCataloged = mainItem0[012].GetW()
	fi.DisableFlags.ExtractFrom(mainItem0[013].GetS1())
	fi.DescriptorFlags.ExtractFrom(mainItem0[014].GetT1())
	fi.FileFlags.ExtractFrom(mainItem0[014].GetS3())
	fi.PCHARFlags.ExtractFrom(mainItem0[015].GetS1())
	fi.AssignMnemonic = strings.TrimRight(mainItem0[016].ToStringAsFieldata(), " ")
	fi.InitialSmoqueLink = mainItem0[017].GetH1()
	fi.NumberOfTimesAssigned = mainItem0[017].GetH2()
	fi.InhibitFlags.ExtractFrom(mainItem0[021].GetS2())
	fi.AssignedIndicator = mainItem0[021].GetT2() != 0
	fi.AbsoluteFCycle = mainItem0[021].GetT3()
	fi.TimeOfLastReference = mainItem0[022].GetW()
	fi.TimeCataloged = mainItem0[023].GetW()
	fi.InitialGranulesReserved = mainItem0[024].GetH1()
	fi.MaxGranules = mainItem0[025].GetH1()
	fi.HighestGranuleAssigned = mainItem0[026].GetH1()
	fi.HighestTrackWritten = mainItem0[027].GetH1()
	// TODO the rest of it
}

// -----------------------------------------------------------------------------

type RemovableFileInfo struct {
	Qualifier                string
	Filename                 string
	ProjectId                string
	AccountId                string
	AbsoluteFileCycle        uint
	TimeOfFirstWriteOrUnload uint64
	DisableFlags             DisableFlags
	DescriptorFlags          DescriptorFlags
	FileFlags                FileFlags
	PCHARFlags               PCHARFlags
	AssignMnemonic           string
	InitialSmoqueLink        uint64
	NumberOfTimesAssigned    uint64
	InhibitFlags             InhibitFlags
	AssignedIndicator        bool
	AbsoluteFCycle           uint64
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

func (fi *RemovableFileInfo) GetQualifier() string {
	return fi.Qualifier
}

func (fi *RemovableFileInfo) GetFilename() string {
	return fi.Filename
}

func (fi *RemovableFileInfo) GetProjectId() string {
	return fi.ProjectId
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
	Qualifier              string
	Filename               string
	ProjectId              string
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

func (fi *TapeFileInfo) GetQualifier() string {
	return fi.Qualifier
}

func (fi *TapeFileInfo) GetFilename() string {
	return fi.Filename
}

func (fi *TapeFileInfo) GetProjectId() string {
	return fi.ProjectId
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
