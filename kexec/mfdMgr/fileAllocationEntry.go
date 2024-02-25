// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"khalehla/kexec"
	"log"
)

// fileAllocationEntry describes the current allocation of tracks to a particular file instance.
// These exist in-memory for every file which is currently assigned.
type fileAllocationEntry struct {
	dadItem0Address       kexec.MFDRelativeAddress
	mainItem0Address      kexec.MFDRelativeAddress
	isUpdated             bool
	highestTrackAllocated kexec.TrackId
	regionEntries         []*FileAllocation
}

func newFileAllocationEntry(
	mainItem0Address kexec.MFDRelativeAddress,
	dadItem0Address kexec.MFDRelativeAddress) *fileAllocationEntry {
	return &fileAllocationEntry{
		dadItem0Address:       dadItem0Address,
		mainItem0Address:      mainItem0Address,
		isUpdated:             false,
		highestTrackAllocated: 0,
		regionEntries:         make([]*FileAllocation, 0),
	}
}

func (fae *fileAllocationEntry) mergeIntoFileAllocationEntry(newEntry *FileAllocation) {
	// puts a new re into the fae at the appropriate location.
	// if it appends to an existing re, then just update that re.
	// we are only called by other code in this file, and those callers *MUST* ensure no overlaps occur.
	for rex, re := range fae.regionEntries {
		if newEntry.fileRegion.trackId < re.fileRegion.trackId {
			// the new entry appears before the indexed entry and after the previous entry
			// if they are the same LDAT, see whether we need to merge
			if newEntry.ldatIndex == re.ldatIndex {
				next := kexec.TrackId(uint64(newEntry.fileRegion.trackId) + uint64(newEntry.fileRegion.trackCount))
				if next == re.fileRegion.trackId {
					// merge them
					re.fileRegion = newEntry.fileRegion
					re.deviceTrackId = newEntry.deviceTrackId
					re.fileRegion.trackCount += newEntry.fileRegion.trackCount
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
			next := kexec.TrackId(uint64(re.fileRegion.trackId) + uint64(re.fileRegion.trackCount))
			if next == newEntry.fileRegion.trackId {
				re.fileRegion.trackCount += newEntry.fileRegion.trackCount
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
	mainItem0Address kexec.MFDRelativeAddress,
	preferred kexec.LDATIndex,
	fileTrackId kexec.TrackId) error {

	// TODO

	return nil
}

// allocateSpecificTrack allocates particular contiguous specified physical tracks
// to be associated with the indicated file-relative tracks.
// If we return an error, we've already stopped the exec
// ONLY FOR VERY SPECIFIC USE-CASES - CALL UNDER LOCK, OR DURING MFD INIT
func (mgr *MFDManager) allocateSpecificTrack(
	mainItem0Address kexec.MFDRelativeAddress,
	fileTrackId kexec.TrackId,
	trackCount kexec.TrackCount,
	ldatIndex kexec.LDATIndex,
	deviceTrackId kexec.TrackId) error {

	fae, ok := mgr.assignedFileAllocations[mainItem0Address]
	if !ok {
		log.Printf("MFDMgr:allocateSpecificTrack Cannot find fae for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return fmt.Errorf("fae not loaded")
	}

	re := newFileAllocation(fileTrackId, trackCount, ldatIndex, deviceTrackId)
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
	mainItem0Address kexec.MFDRelativeAddress,
	fileTrackId kexec.TrackId) (kexec.LDATIndex, kexec.TrackId, error) {

	fae, ok := mgr.assignedFileAllocations[mainItem0Address]
	if !ok {
		log.Printf("MFDMgr:convertFileRelativeTrackId Cannot find fae for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return 0, 0, fmt.Errorf("fae not loaded")
	}

	ldat := kexec.LDATIndex(0_400000)
	devTrackId := kexec.TrackId(0)
	if fileTrackId <= fae.highestTrackAllocated {
		for _, re := range fae.regionEntries {
			if fileTrackId < re.fileRegion.trackId {
				// list is ascending - if we get here, there's no point in continuing
				break
			}
			upperLimit := kexec.TrackId(uint64(re.fileRegion.trackId) + uint64(re.fileRegion.trackCount))
			if fileTrackId < upperLimit {
				// found a good region - update results and stop looking
				ldat = re.ldatIndex
				devTrackId = re.deviceTrackId + (fileTrackId - re.fileRegion.trackId)
				return ldat, devTrackId, nil
			}
		}
	}

	return ldat, devTrackId, nil
}

// loadFileAllocationEntry initializes the fae for a particular file instance.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) loadFileAllocationEntry(mainItem0Address kexec.MFDRelativeAddress) (*fileAllocationEntry, error) {
	_, ok := mgr.assignedFileAllocations[mainItem0Address]
	if ok {
		log.Printf("MFDMgr:loadFileAllocationEntry fae already loaded for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return nil, fmt.Errorf("fae already loaded")
	}

	mainItem0, err := mgr.getMFDSector(mainItem0Address)
	if err != nil {
		return nil, err
	}

	dadAddr := kexec.MFDRelativeAddress(mainItem0[0])
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
			devAddr := kexec.DeviceRelativeWordAddress(dadItem[dx].GetW())
			words := dadItem[dx+1].GetW()
			ldat := kexec.LDATIndex(dadItem[dx+2].GetH2())
			if ldat != 0_400000 {
				re := newFileAllocation(kexec.TrackId(fileWordAddress/1792),
					kexec.TrackCount(words/1792),
					ldat,
					kexec.TrackId(devAddr/1792))
				fae.mergeIntoFileAllocationEntry(re)
			}
			ex++
			dx++
		}

		dadAddr = kexec.MFDRelativeAddress(dadItem[0].GetW())
	}

	mgr.assignedFileAllocations[mainItem0Address] = fae
	return fae, nil
}

// writeFileAllocationEntryUpdates writes an updated fae to the on-disk MFD
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) writeFileAllocationEntryUpdates(mainItem0Address kexec.MFDRelativeAddress) error {
	fae, ok := mgr.assignedFileAllocations[mainItem0Address]
	if !ok {
		log.Printf("MFDMgr:convertFileRelativeTrackId Cannot find fae for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
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
