// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import "khalehla/hardware"

// PackRegion represents a particular region (a track id and a track count) for a specific pack.
type PackRegion struct {
	LDATIndex LDATIndex
	TrackId   hardware.TrackId
	TrackCount hardware.TrackCount
}

func NewPackRegion(ldatIndex LDATIndex, trackId hardware.TrackId, trackCount hardware.TrackCount) *PackRegion {
	return &PackRegion{
		LDATIndex:  ldatIndex,
		TrackId:    trackId,
		TrackCount: trackCount,
	}
}
