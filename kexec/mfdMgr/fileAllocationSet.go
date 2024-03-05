// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import "khalehla/kexec"

// FileAllocationSet describes the current allocation of tracks to a particular file instance.
// These exist in-memory for every file which is currently assigned.
type FileAllocationSet struct {
	DadItem0Address       kexec.MFDRelativeAddress
	MainItem0Address      kexec.MFDRelativeAddress
	IsUpdated             bool
	HighestTrackAllocated kexec.TrackId
	FileAllocations       []*FileAllocation
}

func NewFileAllocationSet(
	mainItem0Address kexec.MFDRelativeAddress,
	dadItem0Address kexec.MFDRelativeAddress) *FileAllocationSet {
	return &FileAllocationSet{
		DadItem0Address:       dadItem0Address,
		MainItem0Address:      mainItem0Address,
		IsUpdated:             false,
		HighestTrackAllocated: 0,
		FileAllocations:       make([]*FileAllocation, 0),
	}
}

// ExtractRegionFromFileAllocationSet extracts the allocation described by the given region
// from this file allocation set.
// Caller MUST ensure that the requested region is a subset (or a match) of exactly one existing entry.
// Returns LDAT index and first device-relative track ID
// corresponding to the pack which contained the indicated region.
func (fas *FileAllocationSet) ExtractRegionFromFileAllocationSet(
	region *kexec.TrackRegion,
) (ldatIndex kexec.LDATIndex, deviceTrackId kexec.TrackId) {
	for rex, fileAlloc := range fas.FileAllocations {
		if fileAlloc.FileRegion.TrackId == region.TrackId {
			ldatIndex = fileAlloc.LDATIndex
			deviceTrackId = fileAlloc.DeviceTrackId

			if fileAlloc.FileRegion.TrackCount == region.TrackCount {
				// deallocating the entire file allocation
				fas.FileAllocations = append(fas.FileAllocations[:rex], fas.FileAllocations[rex+1:]...)
			} else {
				// deallocating from the front of the file allocation
				fileAlloc.FileRegion.TrackId += kexec.TrackId(region.TrackCount)
				fileAlloc.DeviceTrackId += kexec.TrackId(region.TrackCount)
				fileAlloc.FileRegion.TrackCount -= region.TrackCount
			}

			fas.IsUpdated = true
			return
		} else if fileAlloc.FileRegion.TrackId < region.TrackId {
			ldatIndex = fileAlloc.LDATIndex
			deviceTrackId = fileAlloc.DeviceTrackId + (region.TrackId - fileAlloc.FileRegion.TrackId)

			entryLimit := uint64(region.TrackId) + uint64(region.TrackCount)
			allocLimit := uint64(fileAlloc.FileRegion.TrackId) + uint64(fileAlloc.FileRegion.TrackCount)
			if entryLimit == allocLimit {
				// deallocating from the back of the file allocation
				fileAlloc.FileRegion.TrackCount -= region.TrackCount
			} else {
				// deallocating from inside the file allocation with tracks remaining ahead and behind
				newTrackId := kexec.TrackId(entryLimit)
				newTrackCount := kexec.TrackCount(allocLimit - entryLimit)
				newDevTrackId := fileAlloc.DeviceTrackId + (newTrackId - fileAlloc.FileRegion.TrackId)
				newAlloc := NewFileAllocation(newTrackId, newTrackCount, fileAlloc.LDATIndex, newDevTrackId)

				fileAlloc.FileRegion.TrackCount = kexec.TrackCount(region.TrackId - fileAlloc.FileRegion.TrackId)

				temp := append(fas.FileAllocations[:rex+1], newAlloc)
				fas.FileAllocations = append(temp, fas.FileAllocations[rex+1:]...)
			}
			fas.IsUpdated = true
			return
		}
	}

	ldatIndex = kexec.InvalidLDAT
	return
}

// MergeIntoFileAllocationSet puts a new fileAlloc into the fas at the appropriate location.
// if it appends to an existing fileAlloc, then just update that fileAlloc.
// we are only called by other code in this file, and those callers *MUST* ensure no overlaps occur.
// In practical terms, this means do NOT allocate a track which is already allocated.
func (fas *FileAllocationSet) MergeIntoFileAllocationSet(newEntry *FileAllocation) {
	for rex, fileAlloc := range fas.FileAllocations {
		if newEntry.FileRegion.TrackId < fileAlloc.FileRegion.TrackId {
			// the new entry appears before the indexed entry and after the previous entry
			// if they are the same LDAT, see whether we need to merge
			if newEntry.LDATIndex == fileAlloc.LDATIndex {
				next := kexec.TrackId(uint64(newEntry.FileRegion.TrackId) + uint64(newEntry.FileRegion.TrackCount))
				if next == fileAlloc.FileRegion.TrackId {
					// merge them
					fileAlloc.FileRegion = newEntry.FileRegion
					fileAlloc.DeviceTrackId = newEntry.DeviceTrackId
					fileAlloc.FileRegion.TrackCount += newEntry.FileRegion.TrackCount
					fas.IsUpdated = true
					return
				}
			}

			// the new entry is not contiguous with the previous, nor with the next. splice it in.
			newTable := fas.FileAllocations[:rex]
			newTable = append(newTable, newEntry)
			newTable = append(newTable, fas.FileAllocations[rex:]...)
			fas.FileAllocations = newTable
			return
		}

		// If the new entry is on the same pack as the indexed entry, see if the new entry is contiguous
		// with the end of the indexed entry
		if newEntry.LDATIndex == fileAlloc.LDATIndex {
			next := kexec.TrackId(uint64(fileAlloc.FileRegion.TrackId) + uint64(fileAlloc.FileRegion.TrackCount))
			if next == newEntry.FileRegion.TrackId {
				fileAlloc.FileRegion.TrackCount += newEntry.FileRegion.TrackCount
				return
			}
		}

		// move on to the next entry
		rex++
	}

	// If we get here, the new entry is definitely not contiguous with any existing entry.
	fas.FileAllocations = append(fas.FileAllocations, newEntry)
}
