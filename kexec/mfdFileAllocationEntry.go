// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

// MFDFileAllocationEntry describes the current allocation of tracks to a particular file instance.
// These exist in-memory for every file which is currently assigned.
type MFDFileAllocationEntry struct {
	DadItem0Address       MFDRelativeAddress
	MainItem0Address      MFDRelativeAddress
	IsUpdated             bool
	HighestTrackAllocated TrackId
	FileAllocations       []*MFDFileAllocation
}

func NewMFDFileAllocationEntry(
	mainItem0Address MFDRelativeAddress,
	dadItem0Address MFDRelativeAddress) *MFDFileAllocationEntry {
	return &MFDFileAllocationEntry{
		DadItem0Address:       dadItem0Address,
		MainItem0Address:      mainItem0Address,
		IsUpdated:             false,
		HighestTrackAllocated: 0,
		FileAllocations:       make([]*MFDFileAllocation, 0),
	}
}

func (fae *MFDFileAllocationEntry) MergeIntoFileAllocationEntry(newEntry *MFDFileAllocation) {
	// puts a new fileAlloc into the fae at the appropriate location.
	// if it appends to an existing fileAlloc, then just update that fileAlloc.
	// we are only called by other code in this file, and those callers *MUST* ensure no overlaps occur.
	for rex, fileAlloc := range fae.FileAllocations {
		if newEntry.FileRegion.TrackId < fileAlloc.FileRegion.TrackId {
			// the new entry appears before the indexed entry and after the previous entry
			// if they are the same LDAT, see whether we need to merge
			if newEntry.LDATIndex == fileAlloc.LDATIndex {
				next := TrackId(uint64(newEntry.FileRegion.TrackId) + uint64(newEntry.FileRegion.TrackCount))
				if next == fileAlloc.FileRegion.TrackId {
					// merge them
					fileAlloc.FileRegion = newEntry.FileRegion
					fileAlloc.DeviceTrackId = newEntry.DeviceTrackId
					fileAlloc.FileRegion.TrackCount += newEntry.FileRegion.TrackCount
					fae.IsUpdated = true
					return
				}
			}

			// the new entry is not contiguous with the previous, nor with the next. splice it in.
			newTable := fae.FileAllocations[:rex]
			newTable = append(newTable, newEntry)
			newTable = append(newTable, fae.FileAllocations[rex:]...)
			fae.FileAllocations = newTable
			return
		}

		// If the new entry is on the same pack as the indexed entry, see if the new entry is contiguous
		// with the end of the indexed entry
		if newEntry.LDATIndex == fileAlloc.LDATIndex {
			next := TrackId(uint64(fileAlloc.FileRegion.TrackId) + uint64(fileAlloc.FileRegion.TrackCount))
			if next == newEntry.FileRegion.TrackId {
				fileAlloc.FileRegion.TrackCount += newEntry.FileRegion.TrackCount
				return
			}
		}

		// move on to the next entry
		rex++
	}

	// If we get here, the new entry is definitely not contiguous with any existing entry.
	fae.FileAllocations = append(fae.FileAllocations, newEntry)
}
