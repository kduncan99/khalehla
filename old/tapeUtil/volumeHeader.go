// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tapeUtil

type VolumeHeader struct {
	reelNumber         string
	siteId             string
	accountId          string
	numberOfFiles      uint64
	numberOfDataBlocks uint64
}

func NewVolumeHeader() *VolumeHeader {
	return &VolumeHeader{}
}

func (vh *VolumeHeader) SetAccountId(accountId string) *VolumeHeader {
	vh.accountId = accountId
	return vh
}

func (vh *VolumeHeader) SetNumberOfDataBlocks(numberOfDataBlocks uint64) *VolumeHeader {
	vh.numberOfDataBlocks = numberOfDataBlocks
	return vh
}

func (vh *VolumeHeader) SetNumberOfFiles(numberOfFiles uint64) *VolumeHeader {
	vh.numberOfFiles = numberOfFiles
	return vh
}

func (vh *VolumeHeader) SetReelNumber(reelNumber string) *VolumeHeader {
	vh.reelNumber = reelNumber
	return vh
}

func (vh *VolumeHeader) SetSiteId(siteId string) *VolumeHeader {
	vh.siteId = siteId
	return vh
}
