// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/kexec/types"
)

type DiskPackEntry struct {
	packName     string
	mainItemLink uint64
}

type FileInfo interface {
	GetAccountId() string
	GetAbsoluteFileCycle() uint
}

// -----------------------------------------------------------------------------

type FixedFileInfo struct {
	accountId                string
	absoluteFileCycle        uint
	timeOfFirstWriteOrUnload uint64
	descriptorFlags          DescriptorFlags
	writtenTo                bool
	granularity              types.Granularity
	wordAddressable          bool
	assignMnemonic           string
	hasSmoqueEntry           bool
	numberOfTimesAssigned    uint64
	inhibitFlags             InhibitFlags
	timeOfLastReference      uint64
	timeCataloged            uint64
	initialGranulesReserved  uint64
	maxGranules              uint64
	highestGranuleAssigned   uint64
	highestTrackWritten      uint64
	quotaGroupGranules       []uint64
	backupInfo               BackupInfo
	diskPackEntries          []DiskPackEntry
	fileAllocations          []FileAllocation
}

func (fi *FixedFileInfo) GetAccountId() string {
	return fi.accountId
}

func (fi *FixedFileInfo) GetAbsoluteFileCycle() uint {
	return fi.absoluteFileCycle
}

// -----------------------------------------------------------------------------

type RemovableFileInfo struct {
	accountId                string
	absoluteFileCycle        uint
	timeOfFirstWriteOrUnload uint64
	descriptorFlags          DescriptorFlags
	writtenTo                bool
	granularity              types.Granularity
	wordAddressable          bool
	assignMnemonic           string
	hasSmoqueEntry           bool
	numberOfTimesAssigned    uint64
	inhibitFlags             InhibitFlags
	timeOfLastReference      uint64
	timeCataloged            uint64
	initialGranulesReserved  uint64
	maxGranules              uint64
	highestGranuleAssigned   uint64
	highestTrackWritten      uint64
	readKey                  string
	writeKey                 string
	quotaGroupGranules       []uint64
	backupInfo               BackupInfo
	diskPackEntries          []DiskPackEntry
	fileAllocations          []FileAllocation
}

func (fi *RemovableFileInfo) GetAccountId() string {
	return fi.accountId
}

func (fi *RemovableFileInfo) GetAbsoluteFileCycle() uint {
	return fi.absoluteFileCycle
}

// -----------------------------------------------------------------------------

type TapeFileInfo struct {
	accountId              string
	absoluteFileCycle      uint
	descriptorFlags        DescriptorFlags
	assignMnemonic         string
	numberOfTimesAssigned  uint64
	inhibitFlags           InhibitFlags
	currentAssignCount     uint64
	timeOfLastReference    uint64
	timeCataloged          uint64
	density                uint64
	format                 uint64
	features               uint64
	featuresExtension      uint64
	featuresExtension1     uint64
	numberOfReelsCataloged uint64
	mtaPop                 uint64
	noiseConstant          uint64
	translatorMnemonics    []string
	tapeLibraryPool        string
	reelNumber             []string
}

func (fi *TapeFileInfo) GetAccountId() string {
	return fi.accountId
}

func (fi *TapeFileInfo) GetAbsoluteFileCycle() uint {
	return fi.absoluteFileCycle
}
