// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/old/kexec"
)

type RemovableFileCycleInfo struct {
	FileSetIdentifier        FileSetIdentifier
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
	InitialSMOQUELink        uint64
	NumberOfTimesAssigned    uint64
	InhibitFlags             InhibitFlags
	AssignedIndicator        uint64
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
	FileAllocations          []kexec.FileAllocation
}

func (fci *RemovableFileCycleInfo) GetFileSetIdentifier() FileSetIdentifier {
	return fci.FileSetIdentifier
}

func (fci *RemovableFileCycleInfo) GetFileCycleIdentifier() FileCycleIdentifier {
	return fci.FileCycleIdentifier
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

func (fci *RemovableFileCycleInfo) GetInhibitFlags() InhibitFlags {
	return fci.InhibitFlags
}

func (fci *RemovableFileCycleInfo) IsAssigned() bool {
	return fci.AssignedIndicator > 0
}

func (fci *RemovableFileCycleInfo) SetFileCycleIdentifier(fcIdentifier FileCycleIdentifier) {
	fci.FileCycleIdentifier = fcIdentifier
}

func (fci *RemovableFileCycleInfo) SetFileSetIdentifier(fsIdentifier FileSetIdentifier) {
	fci.FileSetIdentifier = fsIdentifier
}
