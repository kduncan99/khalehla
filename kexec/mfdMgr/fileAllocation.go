// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/hardware"
	"khalehla/kexec"
)

type FileAllocation struct {
	FileRegion    kexec.TrackRegion
	LDATIndex     kexec.LDATIndex
	DeviceTrackId hardware.TrackId
}

func NewFileAllocation(
	fileTrackId hardware.TrackId,
	trackCount hardware.TrackCount,
	ldatIndex kexec.LDATIndex,
	deviceTrackId hardware.TrackId) *FileAllocation {
	return &FileAllocation{
		FileRegion: kexec.TrackRegion{
			TrackId:    fileTrackId,
			TrackCount: trackCount,
		},
		LDATIndex:     ldatIndex,
		DeviceTrackId: deviceTrackId,
	}
}
