// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"khalehla/kexec/types"
	"log"
)

// fileAllocationTable is a collection of fileAllocationEntry structs for all files
// which are currently assigned
type fileAllocationTable struct {
	content map[types.MFDRelativeAddress]*fileAllocationEntry
}

func newFileAllocationTable() *fileAllocationTable {
	fat := &fileAllocationTable{}
	fat.content = make(map[types.MFDRelativeAddress]*fileAllocationEntry)
	return fat
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

func newFileAllocationEntry(
	mainItem0Address types.MFDRelativeAddress,
	dadItem0Address types.MFDRelativeAddress) *fileAllocationEntry {
	return &fileAllocationEntry{
		dadItem0Address:       dadItem0Address,
		mainItem0Address:      mainItem0Address,
		isUpdated:             false,
		highestTrackAllocated: 0,
		regionEntries:         make([]*fileRegionEntry, 0),
	}
}

type fileRegionEntry struct {
	fileTrackId types.TrackId
	trackCount  types.TrackCount
	ldatIndex   types.LDATIndex
	packTrackId types.TrackId
}

func newFileRegionEntry(
	fileTrackId types.TrackId,
	trackCount types.TrackCount,
	ldatIndex types.LDATIndex,
	packTrackId types.TrackId) *fileRegionEntry {
	return &fileRegionEntry{
		fileTrackId: fileTrackId,
		trackCount:  trackCount,
		ldatIndex:   ldatIndex,
		packTrackId: packTrackId,
	}
}

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
			return
		}

		// If the new entry is on the same pack as the indexed entry, see if the new entry is contiguous
		// with the end of the indexed entry
		if newEntry.ldatIndex == re.ldatIndex {
			next := types.TrackId(uint64(re.fileTrackId) + uint64(re.trackCount))
			if next == newEntry.fileTrackId {
				re.trackCount += newEntry.trackCount
				return
			}
		}

		// move on to the next entry
		rex++
	}

	// If we get here, the new entry is definitely not contiguous with any existing entry.
	fae.regionEntries = append(fae.regionEntries, newEntry)
}

// allocateTrack allocates a track for the file associated with the given mainItem0Address.
// If provided (preferred != 0), we will try to allocate a track using the preferred ldat index.
// Otherwise:
//
//	If possible we will allocate a track to extend an already-allocated region of the file
//	Else, we will try to allocate a track from the same pack as the first allocation of the file.
//	Finally, we will allocate from any available pack.
//
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) allocateTrack(
	mainItem0Address types.MFDRelativeAddress,
	preferred types.LDATIndex,
	fileTrackId types.TrackId) error {

	// TODO

	return nil
}

// allocateSpecificTrack allocates particular contiguous specified physical tracks
// to be associated with the indicated file-relative tracks.
// If we return an error, we've already stopped the exec
// ONLY FOR VERY SPECIFIC USE-CASES - CALL UNDER LOCK, OR DURING MFD INIT
func (mgr *MFDManager) allocateSpecificTrack(
	mainItem0Address types.MFDRelativeAddress,
	fileTrackId types.TrackId,
	trackCount types.TrackCount,
	ldatIndex types.LDATIndex,
	deviceTrackId types.TrackId) error {

	fae, ok := mgr.fileAllocations.content[mainItem0Address]
	if !ok {
		log.Printf("MFDMgr:allocateSpecificTrack Cannot find fae for address %012o", mainItem0Address)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return fmt.Errorf("fae not loaded")
	}

	re := newFileRegionEntry(fileTrackId, trackCount, ldatIndex, deviceTrackId)
	fae.mergeIntoFileAllocationEntry(re)

	if fileTrackId > fae.highestTrackAllocated {
		fae.highestTrackAllocated = fileTrackId
	}
	fae.isUpdated = true

	return nil
}

// convertFileRelativeAddress takes a file-relative track-id (i.e., word offset from start of file divided by 1792)
// and uses the fae entries in the fat for the given file instance to determine the device LDAT and
// the device-relative track id which contains that file address.
// If the logical track is not allocated, we will return 0_400000 and 0 for those values (since 0 is an invalid LDAT index)
// If the fae is not loaded, we will throw an error - even an empty file has an fae, albeit a puny one.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) convertFileRelativeTrackId(
	mainItem0Address types.MFDRelativeAddress,
	fileTrackId types.TrackId) (types.LDATIndex, types.TrackId, error) {

	fae, ok := mgr.fileAllocations.content[mainItem0Address]
	if !ok {
		log.Printf("MFDMgr:convertFileRelativeTrackId Cannot find fae for address %012o", mainItem0Address)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return 0, 0, fmt.Errorf("fae not loaded")
	}

	ldat := types.LDATIndex(0_400000)
	devTrackId := types.TrackId(0)
	if fileTrackId <= fae.highestTrackAllocated {
		for _, re := range fae.regionEntries {
			if fileTrackId < re.fileTrackId {
				// list is ascending - if we get here, there's no point in continuing
				break
			}
			upperLimit := types.TrackId(uint64(re.fileTrackId) + uint64(re.trackCount))
			if fileTrackId < upperLimit {
				// found a good region - update results and stop looking
				ldat = re.ldatIndex
				devTrackId = re.packTrackId + (fileTrackId - re.fileTrackId)
				return ldat, devTrackId, nil
			}
		}
	}

	return ldat, devTrackId, nil
}

// loadFileAllocationEntry initializes the fae for a particular file instance.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) loadFileAllocationEntry(mainItem0Address types.MFDRelativeAddress) (*fileAllocationEntry, error) {
	_, ok := mgr.fileAllocations.content[mainItem0Address]
	if ok {
		log.Printf("MFDMgr:loadFileAllocationEntry fae already loaded for address %012o", mainItem0Address)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return nil, fmt.Errorf("fae already loaded")
	}

	mainItem0, err := mgr.getMFDSector(mainItem0Address)
	if err != nil {
		return nil, err
	}

	dadAddr := types.MFDRelativeAddress(mainItem0[0])
	fae := &fileAllocationEntry{
		dadItem0Address:  dadAddr,
		mainItem0Address: mainItem0Address,
	}

	for dadAddr&0_400000_000000 == 0 {
		dadItem, err := mgr.getMFDSector(dadAddr)
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
				re := newFileRegionEntry(types.TrackId(fileWordAddress/1792),
					types.TrackCount(words/1792),
					ldat,
					types.TrackId(devAddr/1792))
				fae.mergeIntoFileAllocationEntry(re)
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
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) writeFileAllocationEntryUpdates(mainItem0Address types.MFDRelativeAddress) error {
	fae, ok := mgr.fileAllocations.content[mainItem0Address]
	if !ok {
		log.Printf("MFDMgr:convertFileRelativeTrackId Cannot find fae for address %012o", mainItem0Address)
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
