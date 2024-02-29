// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type TrackRegion struct {
	TrackId    TrackId
	TrackCount TrackCount
}

func NewTrackRegion(trackId TrackId, trackCount TrackCount) *TrackRegion {
	return &TrackRegion{
		TrackId:    trackId,
		TrackCount: trackCount,
	}
}
