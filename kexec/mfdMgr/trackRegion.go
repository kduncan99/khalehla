// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/kexec/types"
)

type TrackRegion struct {
	trackId    types.TrackId
	trackCount types.TrackCount
}

func newTrackRegion(trackId types.TrackId, trackCount types.TrackCount) *TrackRegion {
	return &TrackRegion{
		trackId:    trackId,
		trackCount: trackCount,
	}
}
