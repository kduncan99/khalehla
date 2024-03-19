// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"khalehla/hardware"
	"khalehla/kexec"
	"log"
)

type PackFreeSpaceTable struct {
	Capacity hardware.TrackCount
	Content  []*kexec.TrackRegion
}

func NewPackFreeSpaceTable(capacity hardware.TrackCount) *PackFreeSpaceTable {
	fst := &PackFreeSpaceTable{}
	fst.Capacity = capacity
	fsr := kexec.NewTrackRegion(0, capacity)
	fst.Content = []*kexec.TrackRegion{fsr}
	return fst
}

func (fst *PackFreeSpaceTable) appendEntry(region *kexec.TrackRegion) {
	fst.Content = append(fst.Content, region)
}

func (fst *PackFreeSpaceTable) prependEntry(region *kexec.TrackRegion) {
	temp := []*kexec.TrackRegion{region}
	temp = append(temp, fst.Content...)
	fst.Content = temp
}

func (fst *PackFreeSpaceTable) removeEntryAt(index int) {
	temp := make([]*kexec.TrackRegion, 0)
	temp = append(temp, fst.Content[:index]...)
	if index < len(fst.Content)-1 {
		temp = append(temp, fst.Content[index+1:]...)
	}
	fst.Content = temp
}

func (fst *PackFreeSpaceTable) spliceNewEntryAt(region *kexec.TrackRegion, index int) {
	temp := make([]*kexec.TrackRegion, 0)
	temp = append(temp, fst.Content[:index]...)
	temp = append(temp, region)
	temp = append(temp, fst.Content[index:]...)
	fst.Content = temp
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
	trackId hardware.TrackId,
	trackCount hardware.TrackCount,
) error {
	ok := fst.MarkTrackRegionAllocated(ldatIndex, trackId, trackCount)
	if !ok {
		return fmt.Errorf("track not allocated")
	}
	return nil
}

// AllocateTrackRegion will allocate one region of up to the indicated size.
// If there is no free space at all, we return nil.
func (fst *PackFreeSpaceTable) AllocateTrackRegion(trackCount hardware.TrackCount) *kexec.TrackRegion {
	if len(fst.Content) == 0 {
		return nil
	}

	// Try to satisfy the request with a single region of the exact requested size
	for cx, region := range fst.Content {
		if region.TrackCount == trackCount {
			fst.removeEntryAt(cx)
			return region
		}
	}

	var largest *kexec.TrackRegion
	var largestIndex int
	for rx, region := range fst.Content {
		if region.TrackCount > trackCount {
			result := kexec.NewTrackRegion(region.TrackId, trackCount)
			region.TrackId += hardware.TrackId(trackCount)
			region.TrackCount -= trackCount
			return result
		}

		if largest == nil || region.TrackCount > largest.TrackCount {
			largest = region
			largestIndex = rx
		}
	}

	fst.removeEntryAt(largestIndex)
	return largest
}

// AllocateTracksFromTrackId will attempt to allocate as much contiguous space as possible up to
// trackCount, beginning at the indicated firstTrackId. It may allocate less than trackCount,
// and may not allocate anything at all. The result is the number of tracks allocated.
func (fst *PackFreeSpaceTable) AllocateTracksFromTrackId(
	firstTrackId hardware.TrackId,
	trackCount hardware.TrackCount,
) (result hardware.TrackCount) {
	result = 0

	for _, fsRegion := range fst.Content {
		highest := fsRegion.TrackId + hardware.TrackId(fsRegion.TrackCount) - 1
		if firstTrackId >= fsRegion.TrackId && firstTrackId <= highest {
			available := hardware.TrackCount(highest-firstTrackId) + 1
			if available < trackCount {
				result = available
			} else {
				result = trackCount
			}
			fst.MarkTrackRegionUnallocated(firstTrackId, result)
			return
		}

		if highest < fsRegion.TrackId {
			break
		}
	}

	return
}

// MarkTrackRegionAllocated is a general-purpose function which manipulates the entries in a free space table
func (fst *PackFreeSpaceTable) MarkTrackRegionAllocated(
	ldatIndex kexec.LDATIndex, // only for logging
	trackId hardware.TrackId,
	trackCount hardware.TrackCount,
) bool {
	if trackCount == 0 {
		log.Printf("MarkTrackRegionAllocated ldat:%v id:%v count:%v requested trackCount is zero",
			ldatIndex, trackId, trackCount)
	}

	// We're looking for a region of free space which contains the requested region
	reqTrackLimit := trackId + hardware.TrackId(trackCount) // track limit from specified id and count
	for fx, fsRegion := range fst.Content {
		// Is requested region less than the current entry? If so, there's no point in continuing
		if trackId < fsRegion.TrackId {
			break
		}

		// Does requested region begin within the current entry?
		entryLimit := hardware.TrackId(uint64(fsRegion.TrackId) + uint64(fsRegion.TrackCount))
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
				fst.removeEntryAt(fx)
				return true
			}

			// Is the region to be removed aligned with the front of the current entry?
			if trackId == fsRegion.TrackId {
				fsRegion.TrackId += hardware.TrackId(trackCount)
				fsRegion.TrackCount -= trackCount
				return true
			}

			// Is the region to be removed aligned with the back of the current entry?
			if reqTrackLimit == entryLimit {
				fsRegion.TrackCount -= trackCount
				return true
			}

			// Break the region into two sections. Messy. Don't like it.
			newRegion := kexec.NewTrackRegion(entryLimit, hardware.TrackCount(entryLimit-reqTrackLimit))
			fsRegion.TrackCount = hardware.TrackCount(trackId - fsRegion.TrackId)
			fst.spliceNewEntryAt(newRegion, fx+1)
			return true
		}
	}

	log.Printf("MarkTrackRegionAllocated ldat:%v id:%v count:%v track not allocated",
		ldatIndex, trackId, trackCount)
	return false
}

// MarkTrackRegionUnallocated is a general-purpose function which manipulates the entries in a free space table
func (fst *PackFreeSpaceTable) MarkTrackRegionUnallocated(
	trackId hardware.TrackId,
	trackCount hardware.TrackCount,
) bool {

	if trackCount == 0 {
		log.Printf("MarkTrackRegionUnallocated id:%v count:%v requested trackCount is zero",
			trackId, trackCount)
		return true
	}

	// We are hoping that we do not find an entry which contains all or part of the requested region
	reqTrackLimit := trackId + hardware.TrackId(trackCount) // track limit from specified id and count
	for fx, fsRegion := range fst.Content {
		// Does requested region overlap with this entry?
		entryTrackLimit := hardware.TrackId(uint64(fsRegion.TrackId) + uint64(fsRegion.TrackCount))
		if trackId >= fsRegion.TrackId && trackId < entryTrackLimit {
			log.Printf("MarkTrackRegionUnallocated id:%v count:%v region overlap", trackId, trackCount)
			return false
		} else if reqTrackLimit > fsRegion.TrackId && reqTrackLimit <= entryTrackLimit {
			log.Printf("MarkTrackRegionUnallocated id:%v count:%v region overlap", trackId, trackCount)
			return false
		}

		// Is requested region exactly between this entry and the next?
		// If so, we need to coalesce this region with the next
		if fx < len(fst.Content)-1 {
			fsNext := fst.Content[fx+1]
			if trackId == entryTrackLimit && reqTrackLimit == fst.Content[fx+1].TrackId {
				fsRegion.TrackCount += trackCount + fsNext.TrackCount
				fst.removeEntryAt(fx + 1)
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
			fst.appendEntry(re)
			return true
		}
	}

	// Region is somewhere at the end of the pack. Create a new entry.
	re := kexec.NewTrackRegion(trackId, trackCount)
	fst.appendEntry(re)
	return true
}

// GetFreeTrackCount retrieves a sum of all the free tracks
func (fst *PackFreeSpaceTable) GetFreeTrackCount() hardware.TrackCount {
	count := hardware.TrackCount(0)
	for _, entry := range fst.Content {
		count += entry.TrackCount
	}
	return count
}
