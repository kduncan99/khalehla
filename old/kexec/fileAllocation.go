// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/hardware"
)

type FileAllocation struct {
	FileRegion    TrackRegion
	LDATIndex     LDATIndex
	DeviceTrackId hardware.TrackId
}

func NewFileAllocation(
	fileTrackId hardware.TrackId,
	trackCount hardware.TrackCount,
	ldatIndex LDATIndex,
	deviceTrackId hardware.TrackId) *FileAllocation {
	return &FileAllocation{
		FileRegion: TrackRegion{
			TrackId:    fileTrackId,
			TrackCount: trackCount,
		},
		LDATIndex:     ldatIndex,
		DeviceTrackId: deviceTrackId,
	}
}
