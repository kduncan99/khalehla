// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/pkg"
	"strings"
)

type FixedFileCycleInfo struct {
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
	UnitSelectionIndicators  UnitSelectionIndicators
	QuotaGroupGranules       []uint64
	BackupInfo               BackupInfo
	DiskPackEntries          []DiskPackEntry
	FileAllocations          []FileAllocation
}

func (fci *FixedFileCycleInfo) GetFileSetIdentifier() FileSetIdentifier {
	return fci.FileSetIdentifier
}

func (fci *FixedFileCycleInfo) GetFileCycleIdentifier() FileCycleIdentifier {
	return fci.FileCycleIdentifier
}

func (fci *FixedFileCycleInfo) GetQualifier() string {
	return fci.Qualifier
}

func (fci *FixedFileCycleInfo) GetFilename() string {
	return fci.Filename
}

func (fci *FixedFileCycleInfo) GetProjectId() string {
	return fci.ProjectId
}

func (fci *FixedFileCycleInfo) GetAccountId() string {
	return fci.AccountId
}

func (fci *FixedFileCycleInfo) GetAbsoluteFileCycle() uint {
	return fci.AbsoluteFileCycle
}

func (fci *FixedFileCycleInfo) GetAssignMnemonic() string {
	return fci.AssignMnemonic
}

func (fci *FixedFileCycleInfo) GetInhibitFlags() InhibitFlags {
	return fci.InhibitFlags
}

func (fci *FixedFileCycleInfo) setFileCycleIdentifier(fcIdentifier FileCycleIdentifier) {
	fci.FileCycleIdentifier = fcIdentifier
}

func (fci *FixedFileCycleInfo) setFileSetIdentifier(fsIdentifier FileSetIdentifier) {
	fci.FileSetIdentifier = fsIdentifier
}

// populateFromMainItems populates a FixedFileCycleInfo struct with information derived from the
// provided main items. There must always be at least two main items (sector 0 and sector 1)
// and there may be more, depending upon the forward-links in the previous sectors.
// Any additional main items are going to contain backup reel entries for backup reels beyond
// the first two, or disk pack entries beyond the first five.
func (fci *FixedFileCycleInfo) populateFromMainItems(mainItems [][]pkg.Word36) {
	fci.Qualifier = strings.TrimRight(mainItems[0][1].ToStringAsFieldata()+mainItems[0][2].ToStringAsFieldata(), " ")
	fci.Filename = strings.TrimRight(mainItems[0][3].ToStringAsFieldata()+mainItems[0][4].ToStringAsFieldata(), " ")
	fci.ProjectId = strings.TrimRight(mainItems[0][5].ToStringAsFieldata()+mainItems[0][6].ToStringAsFieldata(), " ")
	fci.AccountId = strings.TrimRight(mainItems[0][7].ToStringAsFieldata()+mainItems[0][010].ToStringAsFieldata(), " ")
	fci.TimeCataloged = mainItems[0][012].GetW()
	fci.DisableFlags.ExtractFrom(mainItems[0][013].GetS1())
	fci.DescriptorFlags.ExtractFrom(mainItems[0][014].GetT1())
	fci.FileFlags.ExtractFrom(mainItems[0][014].GetS3())
	fci.PCHARFlags.ExtractFrom(mainItems[0][015].GetS1())
	fci.AssignMnemonic = strings.TrimRight(mainItems[0][016].ToStringAsFieldata(), " ")
	fci.InitialSmoqueLink = mainItems[0][017].GetH1()
	fci.NumberOfTimesAssigned = mainItems[0][017].GetH2()
	fci.InhibitFlags.ExtractFrom(mainItems[0][021].GetS2())
	fci.AssignedIndicator = mainItems[0][021].GetT2() != 0
	fci.AbsoluteFCycle = mainItems[0][021].GetT3()
	fci.TimeOfLastReference = mainItems[0][022].GetW()
	fci.TimeCataloged = mainItems[0][023].GetW()
	fci.InitialGranulesReserved = mainItems[0][024].GetH1()
	fci.MaxGranules = mainItems[0][025].GetH1()
	fci.HighestGranuleAssigned = mainItems[0][026].GetH1()
	fci.HighestTrackWritten = mainItems[0][027].GetH1()
	fci.UnitSelectionIndicators.ExtractFrom(mainItems[0][033].GetH1())

	fci.AbsoluteFCycle = mainItems[1][07].GetT3()
	fci.BackupInfo.TimeBackupCreated = mainItems[1][010].GetW()
	fci.BackupInfo.FASBits = mainItems[1][011].GetS2()
	fci.BackupInfo.NumberOfTextBlocks = mainItems[1][011].GetH2()
	fci.BackupInfo.StartingFilePosition = mainItems[1][012].GetW()
	fci.BackupInfo.BackupReelNumbers = make([]string, 0)
	if mainItems[1][013].GetW() != 0 {
		reelNumber := mainItems[1][013].ToStringAsFieldata()
		fci.BackupInfo.BackupReelNumbers = append(fci.BackupInfo.BackupReelNumbers, reelNumber)
		if mainItems[1][014].GetW() != 0 {
			reelNumber = mainItems[1][014].ToStringAsFieldata()
			fci.BackupInfo.BackupReelNumbers = append(fci.BackupInfo.BackupReelNumbers, reelNumber)
		}
	}

	fci.DiskPackEntries = make([]DiskPackEntry, 0)
	for wx := 022; wx < 033; wx += 2 {
		if mainItems[1][wx].GetW() != 0 {
			packName := mainItems[1][wx].ToStringAsFieldata()
			link := mainItems[1][wx+1].GetW()
			dpe := DiskPackEntry{
				PackName:     packName,
				MainItemLink: link,
			}
			fci.DiskPackEntries = append(fci.DiskPackEntries, dpe)
		}
	}

	sectorNum := 2
	for sectorNum < len(mainItems) {
		entryType := mainItems[sectorNum][07].GetS1()
		if entryType == 00 {
			for wx := 010; wx < 033; wx += 2 {
				if mainItems[sectorNum][wx].GetW() != 0 {
					packName := mainItems[1][wx].ToStringAsFieldata()
					link := mainItems[1][wx+1].GetW()
					dpe := DiskPackEntry{
						PackName:     packName,
						MainItemLink: link,
					}
					fci.DiskPackEntries = append(fci.DiskPackEntries, dpe)
				}
			}
		} else if entryType == 01 {
			// backup reel entries
			for wx := 010; wx < 033; wx++ {
				if mainItems[sectorNum][wx].GetW() != 0 {
					reelNumber := mainItems[sectorNum][wx].ToStringAsFieldata()
					fci.BackupInfo.BackupReelNumbers = append(fci.BackupInfo.BackupReelNumbers, reelNumber)
				}
			}
		}
	}
}
