// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"khalehla/kexec/types"
)

// fileAllocationTable is a collection of fileAllocationEntry structs for all files
// which are currently assigned
type fileAllocationTable struct {
	content map[types.MFDRelativeAddress]*fileAllocationEntry
}

// fileAllocationEntry describes the current allocation of tracks to a particular file instance.
// These exist in-memory for every file which is currently assigned.
type fileAllocationEntry struct {
	dadItem0Address       types.MFDRelativeAddress
	mainItem0Address      types.MFDRelativeAddress
	isUpdated             bool
	highestTrackAllocated types.TrackId
	regionEntries         []*fileRegionEntry
}

type fileRegionEntry struct {
	fileTrackId types.TrackId
	trackCount  types.TrackCount
	ldatIndex   types.LDATIndex
	packTrackId types.TrackId
}

//func (fae *fileAllocationEntry) coalesceFileAllocationEntry() {
//	// TODO we might not need this one
//	// merges logically-adjacent re's
//	rex := 0
//	for rex < len(fae.regionEntries)-1 {
//		if fae.regionEntries[rex].ldatIndex == fae.regionEntries[rex+1].ldatIndex {
//			next := types.TrackId(uint64(fae.regionEntries[rex].fileTrackId) + uint64(fae.regionEntries[rex].trackCount))
//			if next == fae.regionEntries[rex+1].fileTrackId {
//				fae.regionEntries[rex].trackCount += fae.regionEntries[rex+1].trackCount
//				fae.regionEntries = append(fae.regionEntries[:rex+1], fae.regionEntries[rex+2:]...)
//				continue
//			}
//		}
//		rex++
//	}
//	fae.isUpdated = true
//}

func (fae *fileAllocationEntry) mergeIntoFileAllocationEntry(newEntry *fileRegionEntry) {
	// puts a new re into the fae at the appropriate location.
	// if it appends to an existing re, then just update that re.
	// we are only called by other code in this file, and those callers *MUST* ensure no overlaps occur.
	for rex, re := range fae.regionEntries {
		if newEntry.fileTrackId < re.fileTrackId {
			// the new entry appears before the indexed entry and after the previous entry
			// if they are the same LDAT, see whether we need to merge
			if newEntry.ldatIndex == re.ldatIndex {
				next := types.TrackId(uint64(newEntry.fileTrackId) + uint64(newEntry.trackCount))
				if next == re.fileTrackId {
					// merge them
					re.fileTrackId = newEntry.fileTrackId
					re.packTrackId = newEntry.packTrackId
					re.trackCount += newEntry.trackCount
					fae.isUpdated = true
					return
				}
			}

			// the new entry is not contiguous with the previous, nor with the next. splice it in.
			newTable := fae.regionEntries[:rex]
			newTable = append(newTable, newEntry)
			newTable = append(newTable, fae.regionEntries[rex:]...)
			fae.regionEntries = newTable
			fae.isUpdated = true
			return
		}

		// If the new entry is on the same pack as the indexed entry, see if the new entry is contiguous
		// with the end of the indexed entry
		if newEntry.ldatIndex == re.ldatIndex {
			next := types.TrackId(uint64(re.fileTrackId) + uint64(re.trackCount))
			if next == newEntry.fileTrackId {
				re.trackCount += newEntry.trackCount
				fae.isUpdated = true
				return
			}
		}

		// move on to the next entry
		rex++
	}

	// If we get here, the new entry is definitely not contiguous with any existing entry.
	fae.regionEntries = append(fae.regionEntries, newEntry)
	fae.isUpdated = true
}

// allocateTrack allocates a track for the file associated with the given mainItem0Address.
// If provided (preferred != 0), we will try to allocate a track using the preferred ldat index.
// Otherwise:
//
//	If possible we will allocate a track to extend an already-allocated region of the file
//	Else, we will try to allocate a track from the same pack as the first allocation of the file.
//	Finally, we will allocate from any available pack.
func (mgr *MFDManager) allocateTrack(
	mainItem0Address types.MFDRelativeAddress,
	preferred types.LDATIndex,
	fileTrackId types.TrackId) error {

	// TODO
	//fae, ok := mgr.fileAllocations.content[mainItem0Address]
	//if !ok {
	//	mgr.exec.Stop(types.StopDirectoryErrors)
	//	return fmt.Errorf("fae not loaded")
	//}
	//
	//if preferred != 0 {
	//
	//}

	return nil
}

// allocateSpecificTrack allocates particular contiguous specified physical tracks
// to be associated with the indicated file-relative tracks.
// ONLY FOR VERY SPECIFIC USE-CASES - CALL UNDER LOCK, OR DURING MFD INIT
func (mgr *MFDManager) allocateSpecificTrack(
	mainItem0Address types.MFDRelativeAddress,
	fileTrackId types.TrackId,
	trackCount types.TrackCount,
	ldatIndex types.LDATIndex,
	deviceTrackId types.TrackId) error {

	fae, ok := mgr.fileAllocations.content[mainItem0Address]
	if !ok {
		mgr.exec.Stop(types.StopDirectoryErrors)
		return fmt.Errorf("fae not loaded")
	}

	re := &fileRegionEntry{
		fileTrackId: fileTrackId,
		trackCount:  trackCount,
		ldatIndex:   ldatIndex,
		packTrackId: deviceTrackId,
	}
	fae.mergeIntoFileAllocationEntry(re)

	return nil
}

// convertFileRelativeAddress takes a file-relative track-id (i.e., word offset from start of file divided by 1792)
// and uses the fae entries in the fat for the given file instance to determine the device LDAT and
// the device-relative track id which contains that file address.
// If the logical track is not allocated, we will return 0_400000 and 0 for those values (since 0 is an invalid LDAT index)
// If the fae is not loaded, we will throw an error - even an empty file has an fae, albeit a puny one.
// CALL UNDER LOCK!
func (mgr *MFDManager) convertFileRelativeTrackId(
	mainItem0Address types.MFDRelativeAddress,
	fileTrackId types.TrackId) (types.LDATIndex, types.TrackId, error) {

	fae, ok := mgr.fileAllocations.content[mainItem0Address]
	if !ok {
		mgr.exec.Stop(types.StopDirectoryErrors)
		return 0, 0, fmt.Errorf("fae not loaded")
	}

	ldat := types.LDATIndex(0_400000)
	devTrackId := types.TrackId(0)
	if fileTrackId <= fae.highestTrackAllocated {
		for _, re := range fae.regionEntries {
			if fileTrackId < re.fileTrackId {
				break // list is ascending - if we get here, there's no point in continuing
			}
			upperLimit := types.TrackId(uint64(re.fileTrackId) + uint64(re.trackCount))
			if fileTrackId < upperLimit {
				// found a good region - update results and stop looking
				ldat = re.ldatIndex
				devTrackId = re.packTrackId + (fileTrackId - re.fileTrackId)
				break
			}
		}
	}

	return ldat, devTrackId, nil
}

// loadFileAllocationEntry initializes the fae for a particular file instance.
// CALL UNDER LOCK!
func (mgr *MFDManager) loadFileAllocationEntry(mainItem0Address types.MFDRelativeAddress) (*fileAllocationEntry, error) {
	_, ok := mgr.fileAllocations.content[mainItem0Address]
	if ok {
		mgr.exec.Stop(types.StopDirectoryErrors)
		return nil, fmt.Errorf("fae already loaded")
	}

	mainItem0, err := mgr.getMFDSector(mainItem0Address, false)
	if err != nil {
		return nil, err
	}

	dadAddr := types.MFDRelativeAddress(mainItem0[0])
	fae := &fileAllocationEntry{
		dadItem0Address:  dadAddr,
		mainItem0Address: mainItem0Address,
	}

	for dadAddr&0_400000_000000 == 0 {
		dadItem, err := mgr.getMFDSector(dadAddr, false)
		if err != nil {
			return nil, err
		}

		fileWordAddress := dadItem[02].GetW()
		fileWordLimit := dadItem[03].GetW()
		ex := 0
		dx := 3
		for ex < 8 && fileWordAddress < fileWordLimit {
			devAddr := types.DeviceRelativeWordAddress(dadItem[dx].GetW())
			words := dadItem[dx+1].GetW()
			ldat := types.LDATIndex(dadItem[dx+2].GetH2())
			if ldat != 0_400000 {
				re := &fileRegionEntry{
					fileTrackId: types.TrackId(fileWordAddress / 1792),
					trackCount:  types.TrackCount(words / 1792),
					ldatIndex:   ldat,
					packTrackId: types.TrackId(devAddr / 1792),
				}
				fae.regionEntries = append(fae.regionEntries, re)
			}
			ex++
			dx++
		}

		dadAddr = types.MFDRelativeAddress(dadItem[0].GetW())
	}

	mgr.fileAllocations.content[mainItem0Address] = fae
	return fae, nil
}

// writeFileAllocationEntryUpdates writes an updated fae to the on-disk MFD
// CALL UNDER LOCK!
func (mgr *MFDManager) writeFileAllocationEntryUpdates(mainItem0Address types.MFDRelativeAddress) error {
	fae, ok := mgr.fileAllocations.content[mainItem0Address]
	if !ok {
		mgr.exec.Stop(types.StopDirectoryErrors)
		return fmt.Errorf("fae not loaded")
	}

	if fae.isUpdated {
		// TODO process
		//	rewrite the entire set of DAD entries, allocate a new one if we need to do so,
		//  and release any left over when we're done.
		//  Don't forget to write hole DADs (see pg 2-63 for Device index field)

		fae.isUpdated = false
	}

	return nil
}
