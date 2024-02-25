// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/kexec"
)

type FileAllocation struct {
	fileRegion    TrackRegion
	ldatIndex     kexec.LDATIndex
	deviceTrackId kexec.TrackId
}

func newFileAllocation(
	fileTrackId kexec.TrackId,
	trackCount kexec.TrackCount,
	ldatIndex kexec.LDATIndex,
	deviceTrackId kexec.TrackId) *FileAllocation {
	return &FileAllocation{
		fileRegion: TrackRegion{
			trackId:    fileTrackId,
			trackCount: trackCount,
		},
		ldatIndex:     ldatIndex,
		deviceTrackId: deviceTrackId,
	}
}
