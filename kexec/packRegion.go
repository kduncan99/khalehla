// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

// PackRegion represents a particular region (a track id and a track count) for a specific pack.
type PackRegion struct {
	LDATIndex  LDATIndex
	TrackId    TrackId
	TrackCount TrackCount
}

func NewPackRegion(ldatIndex LDATIndex, trackId TrackId, trackCount TrackCount) *PackRegion {
	return &PackRegion{
		LDATIndex:  ldatIndex,
		TrackId:    trackId,
		TrackCount: trackCount,
	}
}
