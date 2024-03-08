// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import "khalehla/kexec"

// Keeps track of things which pertain to a specific disk pack (an internal MFD struct)
type packDescriptor struct {
	nodeId                     kexec.NodeIdentifier
	prepFactor                 kexec.PrepFactor
	firstDirectoryTrackAddress kexec.DeviceRelativeWordAddress
	canAllocate                bool                // true if pack is UP, false if it is SU - must be set by fac mgr
	packMask                   uint                // used for calculating blocks from sectors
	freeSpaceTable             *PackFreeSpaceTable // represents all unallocated tracks on the pack
	mfdTrackCount              kexec.TrackCount    // number of MFD tracks allocated on the pack
	mfdSectorsUsed             uint64              // number of MFD sectors in use among the MFD tracks
}

func newPackDescriptor(
	nodeId kexec.NodeIdentifier,
	diskAttributes *kexec.DiskAttributes,
) *packDescriptor {

	recordLength := diskAttributes.PackLabelInfo.WordsPerRecord
	trackCount := diskAttributes.PackLabelInfo.TrackCount
	facStatus := diskAttributes.GetFacNodeStatus()

	return &packDescriptor{
		nodeId:                     nodeId,
		prepFactor:                 diskAttributes.PackLabelInfo.PrepFactor,
		firstDirectoryTrackAddress: diskAttributes.PackLabelInfo.FirstDirectoryTrackAddress,
		canAllocate:                facStatus == kexec.FacNodeStatusUp || facStatus == kexec.FacNodeStatusReserved,
		packMask:                   (recordLength / 28) - 1,
		freeSpaceTable:             NewPackFreeSpaceTable(trackCount),
	}
}
