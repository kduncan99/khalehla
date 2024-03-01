// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

type TapeFileCycleInfo struct {
	FileSetIdentifier      FileSetIdentifier
	FileCycleIdentifier    FileCycleIdentifier
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

func (fci *TapeFileCycleInfo) GetFileSetIdentifier() FileSetIdentifier {
	return fci.FileSetIdentifier
}

func (fci *TapeFileCycleInfo) GetFileCycleIdentifier() FileCycleIdentifier {
	return fci.FileCycleIdentifier
}

func (fci *TapeFileCycleInfo) GetQualifier() string {
	return fci.Qualifier
}

func (fci *TapeFileCycleInfo) GetFilename() string {
	return fci.Filename
}

func (fci *TapeFileCycleInfo) GetProjectId() string {
	return fci.ProjectId
}

func (fci *TapeFileCycleInfo) GetAccountId() string {
	return fci.AccountId
}

func (fci *TapeFileCycleInfo) GetAbsoluteFileCycle() uint {
	return fci.AbsoluteFileCycle
}

func (fci *TapeFileCycleInfo) GetAssignMnemonic() string {
	return fci.AssignMnemonic
}
