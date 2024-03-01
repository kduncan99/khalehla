// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

type RemovableFileCycleInfo struct {
	FileCycleIdentifier      FileCycleIdentifier
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

func (fci *RemovableFileCycleInfo) GetQualifier() string {
	return fci.Qualifier
}

func (fci *RemovableFileCycleInfo) GetFilename() string {
	return fci.Filename
}

func (fci *RemovableFileCycleInfo) GetProjectId() string {
	return fci.ProjectId
}

func (fci *RemovableFileCycleInfo) GetAccountId() string {
	return fci.AccountId
}

func (fci *RemovableFileCycleInfo) GetAbsoluteFileCycle() uint {
	return fci.AbsoluteFileCycle
}

func (fci *RemovableFileCycleInfo) GetAssignMnemonic() string {
	return fci.AssignMnemonic
}
