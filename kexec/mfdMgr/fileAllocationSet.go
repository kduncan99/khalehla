// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/hardware"
	"khalehla/kexec"
)

// FileAllocationSet describes the current allocation of tracks to a particular file instance.
// These exist in-memory for every file which is currently assigned.
type FileAllocationSet struct {
	DadItem0Address       kexec.MFDRelativeAddress
	MainItem0Address      kexec.MFDRelativeAddress
	IsUpdated             bool
	HighestTrackAllocated hardware.TrackId
	FileAllocations       []*FileAllocation // these are kept in order by TrackRegion.TrackId
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

// extractRegionFromFileAllocationSet extracts the allocation described by the given region
// from this file allocation set.
// Caller MUST ensure that the requested region is a subset (or a match) of exactly one existing entry.
// Returns LDAT index and first device-relative track ID
// corresponding to the pack which contained the indicated region.
func (fas *FileAllocationSet) extractRegionFromFileAllocationSet(
	region *kexec.TrackRegion,
) (ldatIndex kexec.LDATIndex, deviceTrackId hardware.TrackId) {
	for rex, fileAlloc := range fas.FileAllocations {
		if fileAlloc.FileRegion.TrackId == region.TrackId {
			ldatIndex = fileAlloc.LDATIndex
			deviceTrackId = fileAlloc.DeviceTrackId

			if fileAlloc.FileRegion.TrackCount == region.TrackCount {
				// deallocating the entire file allocation
				fas.FileAllocations = append(fas.FileAllocations[:rex], fas.FileAllocations[rex+1:]...)
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

// mergeIntoFileAllocationSet puts a new fileAlloc into the fas at the appropriate location.
// if it appends to an existing fileAlloc, then just update that fileAlloc.
// we are only called by other code in this file, and those callers *MUST* ensure no overlaps occur.
// In practical terms, this means do NOT allocate a track which is already allocated.
func (fas *FileAllocationSet) mergeIntoFileAllocationSet(newEntry *FileAllocation) {
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
			newTable := make([]*FileAllocation, 0)
			newTable = append(newTable, fas.FileAllocations[:rex]...)
			newTable = append(newTable, newEntry)
			newTable = append(newTable, fas.FileAllocations[rex:]...)
			fas.FileAllocations = newTable
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
	fas.FileAllocations = append(fas.FileAllocations, newEntry)
	fas.IsUpdated = true
}

// findPrecedingAllocation retrieves the FileAllocation which immediately precedes or contains the
// indicated file-relative track id. If we return nil, there is no such FileAllocation.
func (fas *FileAllocationSet) findPrecedingAllocation(
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

// getHighestTrackAssigned finds that value from the FileAllocationSet
func (fas *FileAllocationSet) getHighestTrackAssigned() hardware.TrackId {
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
) (kexec.LDATIndex, hardware.TrackId, bool) {
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
