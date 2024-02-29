// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"khalehla/kexec"
	"log"
)

type PackFreeSpaceTable struct {
	Capacity kexec.TrackCount
	Content  []*kexec.TrackRegion
}

func NewPackFreeSpaceTable(capacity kexec.TrackCount) *PackFreeSpaceTable {
	fst := &PackFreeSpaceTable{}
	fst.Capacity = capacity
	fsr := kexec.NewTrackRegion(0, capacity)
	fst.Content = []*kexec.TrackRegion{fsr}
	return fst
}

// AllocateTrack allocates one track - used primarily for MFD expansion
// an error return does NOT imply an exec stop
func (fst *PackFreeSpaceTable) AllocateTrack() (kexec.MFDTrackId, error) {
	// quick check...
	if len(fst.Content) == 0 {
		return 0, fmt.Errorf("no space")
	}
	// first see if there's a region of just one track
	for _, region := range fst.Content {
		if region.TrackCount == 1 {
			// use this one
			trackId := region.TrackId
			fst.MarkTrackRegionUnallocated(region.TrackId, region.TrackCount)
			return kexec.MFDTrackId(trackId), nil
		}
	}

	// just use the next available
	region := fst.Content[0]
	trackId := kexec.MFDTrackId(region.TrackId)
	fst.MarkTrackRegionUnallocated(region.TrackId, region.TrackCount)
	return trackId, nil
}

// AllocateSpecificTrackRegion is used only when it has been determined by some external means, that a particular
// track or range of tracks is not to be allocated otherwise (such as for VOL1 or directory tracks).
func (fst *PackFreeSpaceTable) AllocateSpecificTrackRegion(
	ldatIndex kexec.LDATIndex,
	trackId kexec.TrackId,
	trackCount kexec.TrackCount,
) error {

	ok := fst.MarkTrackRegionAllocated(ldatIndex, trackId, trackCount)
	if !ok {
		return fmt.Errorf("track not allocated")
	}
	return nil
}

// MarkTrackRegionAllocated is a general-purpose function which manipulates the entries in a free space table
func (fst *PackFreeSpaceTable) MarkTrackRegionAllocated(
	ldatIndex kexec.LDATIndex, // only for logging
	trackId kexec.TrackId,
	trackCount kexec.TrackCount,
) bool {

	if trackCount == 0 {
		log.Printf("MarkTrackRegionAllocated ldat:%v id:%v count:%v requested trackCount is zero",
			ldatIndex, trackId, trackCount)
	}

	// We're looking for a region of free space which contains the requested region
	reqTrackLimit := trackId + kexec.TrackId(trackCount) // track limit from specified id and count
	for fx, fsRegion := range fst.Content {
		// Is requested region less than the current entry? If so, there's no point in continuing
		if trackId < fsRegion.TrackId {
			break
		}

		// Does requested region begin within the current entry?
		entryLimit := kexec.TrackId(uint64(fsRegion.TrackId) + uint64(fsRegion.TrackCount))
		if trackId >= fsRegion.TrackId && trackId <= entryLimit {
			// Quick check to ensure requested region does not exceed this entry.
			// If it does, something is bigly wrong
			if reqTrackLimit > entryLimit {
				log.Printf("MarkTrackRegionAllocated ldat:%v id:%v count:%v region too big",
					ldatIndex, trackId, trackCount)
				return false
			}

			// Does the requested region exactly match the current entry?
			// If so, just remove the current entry
			if fsRegion.TrackCount == trackCount {
				fst.Content = append(fst.Content[:fx], fst.Content[fx:]...)
				return true
			}

			// Is the region to be removed aligned with the front of the current entry?
			if trackId == fsRegion.TrackId {
				fsRegion.TrackId += kexec.TrackId(trackCount)
				fsRegion.TrackCount -= trackCount
				return true
			}

			// Is the region to be removed aligned with the back of the current entry?
			if reqTrackLimit == entryLimit {
				fsRegion.TrackCount -= trackCount
				return true
			}

			// Break the region into two sections. Messy. Don't like it.
			newRegion := kexec.NewTrackRegion(entryLimit, kexec.TrackCount(entryLimit-reqTrackLimit))
			fsRegion.TrackCount = kexec.TrackCount(trackId - fsRegion.TrackId)
			newTable := append(fst.Content[0:fx+1], newRegion)
			fst.Content = append(newTable, fst.Content[fx+1])
			return true
		}
	}

	log.Printf("MarkTrackRegionAllocated ldat:%v id:%v count:%v track not allocated",
		ldatIndex, trackId, trackCount)
	return false
}

// MarkTrackRegionUnallocated is a general-purpose function which manipulates the entries in a free space table
func (fst *PackFreeSpaceTable) MarkTrackRegionUnallocated(
	trackId kexec.TrackId,
	trackCount kexec.TrackCount) bool {

	if trackCount == 0 {
		log.Printf("MarkTrackRegionUnallocated id:%v count:%v requested trackCount is zero",
			trackId, trackCount)
	}

	// We are hoping that we do not find an entry which contains all or part of the requested region
	reqTrackLimit := trackId + kexec.TrackId(trackCount) // track limit from specified id and count
	for fx, fsRegion := range fst.Content {
		// Does requested region overlap with this entry?
		entryTrackLimit := kexec.TrackId(uint64(fsRegion.TrackId) + uint64(fsRegion.TrackCount))
		if trackId >= fsRegion.TrackId && trackId < entryTrackLimit {
			log.Printf("MarkTrackRegionUnallocated id:%v count:%v region overlap", trackId, trackCount)
			return false
		} else if reqTrackLimit > fsRegion.TrackId && reqTrackLimit <= entryTrackLimit {
			log.Printf("MarkTrackRegionUnallocated id:%v count:%v region overlap", trackId, trackCount)
			return false
		}

		// Is requested region between this entry and the next?
		// If so, we need to coalesce this and the next region
		if fx < len(fst.Content)-1 {
			fsNext := fst.Content[fx+1]
			if trackId == entryTrackLimit && reqTrackLimit == fst.Content[fx+1].TrackId {
				fsRegion.TrackCount += trackCount + fsNext.TrackCount
				fst.Content = append(fst.Content[:fx+1], fst.Content[fx+1:]...)
				return true
			}
		}

		// Is requested region aligned with the front of this entry?
		if reqTrackLimit == fsRegion.TrackId {
			fsRegion.TrackId = trackId
			fsRegion.TrackCount += trackCount
			return true
		}

		// Is requested region aligned with the back of this entry?
		if trackId == entryTrackLimit {
			fsRegion.TrackCount += trackCount
			return true
		}

		// Region is not aligned with the front or back of this entry, nor does it overlap.
		// If it is ahead of this entry, then we just need to insert a new entry for the requested region.
		if trackId < fsRegion.TrackId {
			re := kexec.NewTrackRegion(trackId, trackCount)
			newTable := append(fst.Content[:fx], re)
			fst.Content = append(newTable, fst.Content[fx:]...)
			return true
		}
	}

	// Region is somewhere at the end of the pack. Create a new entry.
	re := kexec.NewTrackRegion(trackId, trackCount)
	fst.Content = append(fst.Content, re)
	return true
}

// GetFreeTrackCount retrieves a sum of all the free tracks
func (fst *PackFreeSpaceTable) GetFreeTrackCount() kexec.TrackCount {
	count := kexec.TrackCount(0)
	for _, entry := range fst.Content {
		count += entry.TrackCount
	}
	return count
}
