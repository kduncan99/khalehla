// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import "khalehla/hardware"

type TrackRegion struct {
	TrackId    hardware.TrackId
	TrackCount hardware.TrackCount
}

func NewTrackRegion(trackId hardware.TrackId, trackCount hardware.TrackCount) *TrackRegion {
	return &TrackRegion{
		TrackId:    trackId,
		TrackCount: trackCount,
	}
}
