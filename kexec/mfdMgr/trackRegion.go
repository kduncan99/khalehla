// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/kexec"
)

type TrackRegion struct {
	trackId    kexec.TrackId
	trackCount kexec.TrackCount
}

func newTrackRegion(trackId kexec.TrackId, trackCount kexec.TrackCount) *TrackRegion {
	return &TrackRegion{
		trackId:    trackId,
		trackCount: trackCount,
	}
}
