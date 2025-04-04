// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/hardware"
)

// FileAllocationSet describes the current allocation of tracks to a particular file instance.
// These exist in-memory for every file which is currently assigned.
type FileAllocationSet struct {
	DadItem0Address  MFDRelativeAddress
	MainItem0Address MFDRelativeAddress
	IsUpdated        bool
	FileAllocations  []*FileAllocation // these are kept in order by TrackRegion.TrackId
}

func NewFileAllocationSet(
	mainItem0Address MFDRelativeAddress,
	dadItem0Address MFDRelativeAddress) *FileAllocationSet {
	return &FileAllocationSet{
		DadItem0Address:  dadItem0Address,
		MainItem0Address: mainItem0Address,
		IsUpdated:        false,
		FileAllocations:  make([]*FileAllocation, 0),
	}
}

// GetHighestTrackAllocated calculates the highest track accounted for by the file allocations.
// returns the trackId and ok==true if there are any allocations, or 0 and ok==false otherwise.
func (fas *FileAllocationSet) GetHighestTrackAllocated() (trackId hardware.TrackId, ok bool) {
	faCount := len(fas.FileAllocations)
	if faCount > 0 {
		lastFa := fas.FileAllocations[faCount-1]
		trackId = hardware.TrackId(uint64(lastFa.FileRegion.TrackId)+uint64(lastFa.FileRegion.TrackCount)) - 1
		ok = true
	} else {
		trackId = 0
		ok = false
	}
	return
}

func (fas *FileAllocationSet) appendEntry(alloc *FileAllocation) {
	fas.FileAllocations = append(fas.FileAllocations, alloc)
}

func (fas *FileAllocationSet) insertEntryAt(alloc *FileAllocation, index int) {
	temp := make([]*FileAllocation, 0)
	temp = append(temp, fas.FileAllocations[:index]...)
	temp = append(temp, alloc)
	temp = append(temp, fas.FileAllocations[index:]...)
	fas.FileAllocations = temp
}

func (fas *FileAllocationSet) removeEntryAt(index int) {
	temp := make([]*FileAllocation, 0)
	temp = append(temp, fas.FileAllocations[:index]...)
	temp = append(temp, fas.FileAllocations[index+1:]...)
	fas.FileAllocations = temp
}

// ExtractRegionFromFileAllocationSet extracts the allocation described by the given region
// from this file allocation set.
// Caller MUST ensure that the requested region is a subset (or a match) of exactly one existing entry.
// Returns LDAT index and first device-relative track ID
// corresponding to the pack which contained the indicated region.
func (fas *FileAllocationSet) ExtractRegionFromFileAllocationSet(
	region *TrackRegion,
) (ldatIndex LDATIndex, deviceTrackId hardware.TrackId) {
	for rex, fileAlloc := range fas.FileAllocations {
		if fileAlloc.FileRegion.TrackId == region.TrackId {
			ldatIndex = fileAlloc.LDATIndex
			deviceTrackId = fileAlloc.DeviceTrackId

			if fileAlloc.FileRegion.TrackCount == region.TrackCount {
				// deallocating the entire file allocation
				fas.removeEntryAt(rex)
			} else {
				// deallocating from the front of the file allocation
				fileAlloc.FileRegion.TrackId += hardware.TrackId(region.TrackCount)
				fileAlloc.DeviceTrackId += hardware.TrackId(region.TrackCount)
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
				newTrackId := hardware.TrackId(entryLimit)
				newTrackCount := hardware.TrackCount(allocLimit - entryLimit)
				newDevTrackId := fileAlloc.DeviceTrackId + (newTrackId - fileAlloc.FileRegion.TrackId)
				newAlloc := NewFileAllocation(newTrackId, newTrackCount, fileAlloc.LDATIndex, newDevTrackId)

				fileAlloc.FileRegion.TrackCount = hardware.TrackCount(region.TrackId - fileAlloc.FileRegion.TrackId)
				fas.insertEntryAt(newAlloc, rex+1)
			}
			fas.IsUpdated = true
			return
		}
	}

	ldatIndex = InvalidLDAT
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
				next := hardware.TrackId(uint64(newEntry.FileRegion.TrackId) + uint64(newEntry.FileRegion.TrackCount))
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
			fas.insertEntryAt(newEntry, rex)
			fas.IsUpdated = true
			return
		}

		// If the new entry is on the same pack as the indexed entry, see if the new entry is contiguous
		// with the end of the indexed entry
		if newEntry.LDATIndex == fileAlloc.LDATIndex {
			next := hardware.TrackId(uint64(fileAlloc.FileRegion.TrackId) + uint64(fileAlloc.FileRegion.TrackCount))
			if next == newEntry.FileRegion.TrackId {
				fileAlloc.FileRegion.TrackCount += newEntry.FileRegion.TrackCount
				fas.IsUpdated = true
				return
			}
		}

		// move on to the next entry
		rex++
	}

	// If we get here, the new entry is definitely not contiguous with any existing entry.
	fas.appendEntry(newEntry)
	fas.IsUpdated = true
}

// FindPrecedingAllocation retrieves the FileAllocation which immediately precedes or contains the
// indicated file-relative track id. If we return nil, there is no such FileAllocation.
func (fas *FileAllocationSet) FindPrecedingAllocation(
	fileTrackId hardware.TrackId,
) (alloc *FileAllocation) {
	alloc = nil
	for _, fa := range fas.FileAllocations {
		if fileTrackId >= fa.FileRegion.TrackId {
			alloc = fa
		} else {
			break
		}
	}
	return
}

// GetHighestTrackAssigned finds that value from the FileAllocationSet
func (fas *FileAllocationSet) GetHighestTrackAssigned() hardware.TrackId {
	entryCount := len(fas.FileAllocations)
	if entryCount == 0 {
		return 0
	} else {
		last := entryCount - 1
		fAlloc := fas.FileAllocations[last]
		return fAlloc.FileRegion.TrackId + hardware.TrackId(fAlloc.FileRegion.TrackCount) - 1
	}
}

// resolveFileRelativeTrackId converts a file-relative track id (file-relative sector address * 28,
// or file-relative word address * 1792) to the LDAT index of the pack which contains that track,
// and to the corresponding device/pack-relative track ID.
// If we return false, no allocation exists (the space has not (yet) been allocated).
func (fas *FileAllocationSet) resolveFileRelativeTrackId(
	fileTrackId hardware.TrackId,
) (LDATIndex, hardware.TrackId, bool) {
	for _, fa := range fas.FileAllocations {
		highestAllocTrack := hardware.TrackId(uint64(fa.FileRegion.TrackId) + uint64(fa.FileRegion.TrackCount) - 1)
		if fileTrackId >= fa.FileRegion.TrackId && fileTrackId <= highestAllocTrack {
			offset := fileTrackId - fa.FileRegion.TrackId
			return fa.LDATIndex, fa.DeviceTrackId + offset, true
		} else if highestAllocTrack < fileTrackId {
			break
		}
	}

	return 0, 0, false
}
