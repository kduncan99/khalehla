// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type MFDFileAllocation struct {
	FileRegion    TrackRegion
	LDATIndex     LDATIndex
	DeviceTrackId TrackId
}

func NewMFDFileAllocation(
	fileTrackId TrackId,
	trackCount TrackCount,
	ldatIndex LDATIndex,
	deviceTrackId TrackId) *MFDFileAllocation {
	return &MFDFileAllocation{
		FileRegion: TrackRegion{
			TrackId:    fileTrackId,
			TrackCount: trackCount,
		},
		LDATIndex:     ldatIndex,
		DeviceTrackId: deviceTrackId,
	}
}
