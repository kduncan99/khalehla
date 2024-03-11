// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/hardware"
	"khalehla/kexec"
)

// Keeps track of things which pertain to a specific disk pack (an internal MFD struct)
type packDescriptor struct {
	nodeId                     hardware.NodeIdentifier
	prepFactor                 hardware.PrepFactor
	firstDirectoryTrackAddress kexec.DeviceRelativeWordAddress
	canAllocate                bool                // true if pack is UP, false if it is SU - must be set by fac mgr
	packMask                   uint                // used for calculating blocks from sectors
	freeSpaceTable             *PackFreeSpaceTable // represents all unallocated tracks on the pack
	mfdTrackCount              hardware.TrackCount // number of MFD tracks allocated on the pack
	mfdSectorsUsed             uint64              // number of MFD sectors in use among the MFD tracks
}

func newPackDescriptor(
	nodeId hardware.NodeIdentifier,
	prepFactor hardware.PrepFactor,
	trackCount hardware.TrackCount,
	firstDirectoryTrackAddress kexec.DeviceRelativeWordAddress,
	nodeStatus kexec.FacNodeStatus,
) *packDescriptor {

	recordLength := uint(prepFactor)

	return &packDescriptor{
		nodeId:                     nodeId,
		prepFactor:                 prepFactor,
		firstDirectoryTrackAddress: firstDirectoryTrackAddress,
		canAllocate:                nodeStatus == kexec.FacNodeStatusUp || nodeStatus == kexec.FacNodeStatusReserved,
		packMask:                   (recordLength / 28) - 1,
		freeSpaceTable:             NewPackFreeSpaceTable(trackCount),
	}
}
