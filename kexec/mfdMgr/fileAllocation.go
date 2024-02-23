// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/kexec/types"
)

type FileAllocation struct {
	fileRegion    TrackRegion
	ldatIndex     types.LDATIndex
	deviceTrackId types.TrackId
}

func newFileAllocation(
	fileTrackId types.TrackId,
	trackCount types.TrackCount,
	ldatIndex types.LDATIndex,
	deviceTrackId types.TrackId) *FileAllocation {
	return &FileAllocation{
		fileRegion: TrackRegion{
			trackId:    fileTrackId,
			trackCount: trackCount,
		},
		ldatIndex:     ldatIndex,
		deviceTrackId: deviceTrackId,
	}
}
