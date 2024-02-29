// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import "khalehla/kexec"

type FileAllocation struct {
	FileRegion    kexec.TrackRegion
	LDATIndex     kexec.LDATIndex
	DeviceTrackId kexec.TrackId
}

func NewFileAllocation(
	fileTrackId kexec.TrackId,
	trackCount kexec.TrackCount,
	ldatIndex kexec.LDATIndex,
	deviceTrackId kexec.TrackId) *FileAllocation {
	return &FileAllocation{
		FileRegion: kexec.TrackRegion{
			TrackId:    fileTrackId,
			TrackCount: trackCount,
		},
		LDATIndex:     ldatIndex,
		DeviceTrackId: deviceTrackId,
	}
}
