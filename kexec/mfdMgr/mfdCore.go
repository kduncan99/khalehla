// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// we are allowed to know about facMgr, nodeMgr, consMgr

import (
	"fmt"
	"khalehla/kexec"
	"time"

	//	"khalehla/kexec/facilitiesMgr"
	"khalehla/kexec/nodeMgr"
	"khalehla/pkg"
	"log"
)

// adjustLeadItemLinks shifts the links in the given lead item(s) downward by the indicated
// shift value (presumably to accommodate a new higher-cycled link).
// Adjusts the current file cycle count in leadItem0 accordingly.
// Make darn sure there is a leadItem1 if it is necessary to contain the expanded set of links.
// Otherwise, it can be nil.
// Caller must ensure both leadItem0 and leadItem1 (if it exists) are marked dirty - we don't do it.
func adjustLeadItemLinks(leadItem0 []pkg.Word36, leadItem1 []pkg.Word36, shift uint) {
	currentRange := uint(leadItem0[011].GetS4())
	newRange := currentRange + shift
	sourceIndex := currentRange - 1
	destIndex := newRange - 1

	for {
		sourceWord := getLeadItemLinkWord(leadItem0, leadItem1, sourceIndex)
		destWord := getLeadItemLinkWord(leadItem0, leadItem1, destIndex)
		destWord.SetW(sourceWord.GetW())
		if sourceIndex == 0 {
			break
		}
		sourceIndex--
		destIndex--
	}

	for {
		destWord := getLeadItemLinkWord(leadItem0, leadItem1, destIndex)
		destWord.SetW(0)
		if destIndex == 0 {
			break
		}
		destIndex--
	}

	leadItem0[011].SetS4(uint64(newRange))
}

// allocateDirectorySector allocates an MFD directory sector for the caller.
// If preferredLDAT is not InvalidLDAT we will try to allocate a sector from this pack first.
// Apart from this, we prefer packs with the least number of allocated sectors, to balance the allocations.
// If there is no free sector, we allocate a new track (again using the preferredLDAT), then
// allocate the first free sector from that track.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) allocateDirectorySector(
	preferredLDAT kexec.LDATIndex,
) (kexec.MFDRelativeAddress, []pkg.Word36, error) {
	// Are there any free sectors? If not, allocate a directory track
	if len(mgr.freeMFDSectors) == 0 {
		_, _, err := mgr.allocateDirectoryTrack(preferredLDAT)
		if err != nil {
			return 0, nil, err
		}
	}

	// Choose a pack among the following priorities:
	//		preferred (if not InvalidLDAT)
	//		pack with the least number of sectors
	//		anything available
	chosenAddress := kexec.InvalidLink
	chosenSectorsUsed := uint64(0)
	chosenIndex := 0
	for fsx, addr := range mgr.freeMFDSectors {
		ldat := getLDATIndexFromMFDAddress(addr)
		desc := mgr.fixedPackDescriptors[ldat]
		freeCount := (uint64(desc.mfdTrackCount) * 64) - desc.mfdSectorsUsed

		if freeCount > 0 {
			if ldat == preferredLDAT {
				chosenAddress = addr
				chosenIndex = fsx
				break
			}

			if desc.mfdSectorsUsed < chosenSectorsUsed {
				chosenAddress = addr
				chosenSectorsUsed = desc.mfdSectorsUsed
				chosenIndex = fsx
			}
		}
	}

	// When we get here, we *will* have valid chosen elements.
	// If we had no free sectors anywhere, we would have allocated a new track
	// (and if that failed, we would already have crashed and returned)
	mgr.freeMFDSectors = append(mgr.freeMFDSectors[:chosenIndex], mgr.freeMFDSectors[chosenIndex+1:]...)
	data, err := mgr.getMFDSector(chosenAddress)
	if err != nil {
		return 0, nil, err
	}

	err = mgr.markDirectorySectorAllocated(chosenAddress)
	if err != nil {
		return 0, nil, err
	}

	// clean it out for the caller
	for wx := 0; wx < 28; wx++ {
		data[wx] = 0
	}

	return chosenAddress, data, nil
}

// allocateDirectoryTrack allocates a new directory track for MFD purposes.
// If preferredLDAT is not InvalidLDAT, we try to allocate from that pack first.
// Apart from that, we allocate from the pack which has the least number of MFD tracks to preserve balance.
// If the track is not a ninth track, we update the appropriate DAS track.
// Otherwise, the new track becomes a DAS track, and we *still* have to update the previous DAS track.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) allocateDirectoryTrack(
	preferredLDAT kexec.LDATIndex,
) (kexec.LDATIndex, kexec.MFDTrackId, error) {
	chosenLDAT := kexec.InvalidLDAT
	var chosenDesc *packDescriptor
	chosenAvailableTracks := kexec.TrackCount(0)

	if preferredLDAT != kexec.InvalidLDAT {
		packDesc, ok := mgr.fixedPackDescriptors[preferredLDAT]
		if ok {
			availMFDTracks := 07777 - packDesc.mfdTrackCount
			availTracks := packDesc.freeSpaceTable.GetFreeTrackCount()
			if availMFDTracks > 0 && availTracks > 0 {
				chosenLDAT = preferredLDAT
				chosenDesc = packDesc
				chosenAvailableTracks = availTracks
			}
		}
	}

	if chosenLDAT == kexec.InvalidLDAT {
		for ldat, packDesc := range mgr.fixedPackDescriptors {
			availMFDTracks := 07777 - packDesc.mfdTrackCount
			availTracks := packDesc.freeSpaceTable.GetFreeTrackCount()
			if availMFDTracks > 0 && availTracks > chosenAvailableTracks {
				chosenLDAT = ldat
				chosenDesc = packDesc
				chosenAvailableTracks = availTracks
			}
		}
	}

	if chosenLDAT == kexec.InvalidLDAT {
		log.Printf("MFDMgr:No space available for directory track allocation")
		mgr.exec.Stop(kexec.StopExecRequestForMassStorageFailed)
		return 0, 0, fmt.Errorf("no disk")
	}

	// First, find the MFD relative address of the first unused MFD track
	trackId := kexec.MFDTrackId(0)
	for {
		mfdAddr := composeMFDAddress(chosenLDAT, trackId, 0)
		_, ok := mgr.cachedTracks[mfdAddr]
		if !ok {
			break
		}
		trackId++
	}

	// Now allocate a track (any track) from the pack for the chosen LDAT.
	// Make sure we update the MFD track count.
	trackId, _ = chosenDesc.freeSpaceTable.AllocateTrack()
	chosenDesc.mfdTrackCount++
	return chosenLDAT, trackId, nil
}

// allocateLeadItem1 allocates a lead item directory sector which extends a
// currently-not-extended lead item 0 such that it can refer to a larger file cycle range.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) allocateLeadItem1(
	leadItem0Address kexec.MFDRelativeAddress,
	leadItem0 []pkg.Word36,
) (leadItem1Address kexec.MFDRelativeAddress, leadItem1 []pkg.Word36, err error) {
	preferredLDAT := getLDATIndexFromMFDAddress(leadItem0Address)
	leadItem1Address, leadItem1, err = mgr.allocateDirectorySector(preferredLDAT)
	if err == nil {
		leadItem0[0].SetW(0_100000_000000 | uint64(leadItem1Address))
		leadItem1[0].SetW(0_400000_000000)
		mgr.markDirectorySectorDirty(leadItem0Address)
		mgr.markDirectorySectorDirty(leadItem1Address)
	}
	return
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

	re := NewFileAllocation(fileTrackId, trackCount, ldatIndex, deviceTrackId)
	fae.MergeIntoFileAllocationSet(re)

	if fileTrackId > fae.HighestTrackAllocated {
		fae.HighestTrackAllocated = fileTrackId
	}
	fae.IsUpdated = true

	return nil
}

// TODO allocateTrackRegion - params include mainItem0Address, region
//   will return an array of file allocations (?), possibly one encompassing the whole request, else several which
//   when combined, encompass the whole request. Maybe instead of returning, we just do it behind the scenes?

// bootstrapMFD creates the various MFD structures as part of MFD initialization.
// One consequence is the cataloging of SYS$*MFD$$.
// Since this is used during initialization we do not call it under lock.
func (mgr *MFDManager) bootstrapMFD() error {
	log.Printf("MFDMgr:bootstrapMFD start")

	cfg := mgr.exec.GetConfiguration()

	// Find the highest and lowest LDAT indices
	var lowestLDAT = kexec.InvalidLDAT
	var highestLDAT = kexec.LDATIndex(0)
	for ldat := range mgr.fixedPackDescriptors {
		if lowestLDAT == kexec.InvalidLDAT && ldat < lowestLDAT {
			lowestLDAT = ldat
		}
		if ldat > highestLDAT {
			highestLDAT = ldat
		}
	}

	// We've already initialized the DAS sector - put the expected free sectors (2 through 63)
	// into the free sector list.
	for sx := 2; sx < 64; sx++ {
		sectorAddr := composeMFDAddress(lowestLDAT, 0, kexec.MFDSectorId(sx))
		mgr.freeMFDSectors = append(mgr.freeMFDSectors, sectorAddr)
	}

	// Allocate MFD sectors for MFD$$ file items not including DAD tables (we do those separately)
	leadAddr0, leadItem0, _ := mgr.allocateDirectorySector(lowestLDAT)
	mainAddr0, mainItem0, _ := mgr.allocateDirectorySector(lowestLDAT)
	mainAddr1, mainItem1, _ := mgr.allocateDirectorySector(lowestLDAT)

	mgr.mfdFileMainItem0Address = mainAddr0 // we'll need this later

	// TODO can these items use a services or core routine other than populate*** as below?
	mfdFileQualifier := "SYS$"
	mfdFileName := "MFD$$"
	mfdProjectId := "EXEC-8"
	populateNewLeadItem0(
		leadItem0,
		mfdFileQualifier,
		mfdFileName,
		1,
		"",
		"",
		mfdProjectId,
		0,
		true,
		uint64(mainAddr0))

	populateMassStorageMainItem0(mainItem0,
		leadAddr0,
		mainAddr1,
		mfdFileQualifier,
		mfdFileName,
		1,
		"",
		"",
		mfdProjectId,
		cfg.MasterAccountId,
		cfg.MassStorageDefaultMnemonic,
		DescriptorFlags{},
		PCHARFlags{
			Granularity:       kexec.TrackGranularity,
			IsWordAddressable: false,
		},
		InhibitFlags{
			IsGuarded:           true,
			IsUnloadInhibited:   true,
			IsPrivate:           true,
			isAssignedExclusive: false,
			IsWriteOnly:         false,
			IsReadOnly:          false,
		},
		false,
		0,
		262153,
		nil)

	populateFixedMainItem1(mainItem1, mfdFileQualifier, mfdFileName, mainAddr0, 1, nil)

	// Before we can play DAD table games, we have to get the MFD$$ in-core structures in place,
	// including *particularly* the file allocation table.
	// We need to create one allocation region for each pack's initial directory track.
	highestMFDTrackId := kexec.TrackId(0)
	fae := NewFileAllocationSet(mainAddr0, 0_400000_000000)
	mgr.assignedFileAllocations[mainAddr0] = fae

	for ldat, desc := range mgr.fixedPackDescriptors {
		mfdTrackId := kexec.TrackId(ldat << 12)
		if mfdTrackId > highestMFDTrackId {
			highestMFDTrackId = mfdTrackId
		}

		packTrackId := desc.firstDirectoryTrackAddress / 1792
		err := mgr.allocateSpecificTrack(mainAddr0, mfdTrackId, 1, ldat, kexec.TrackId(packTrackId))
		if err != nil {
			return err
		}
	}

	// update main item
	// we are TRK granularity, but we are a large file
	highestGranuleAssigned := uint64(highestMFDTrackId)
	mod := highestGranuleAssigned & 077
	highestPositionWritten := highestGranuleAssigned >> 6
	if mod > 0 {
		highestPositionWritten++
	}
	mainItem0[026].SetH1(highestGranuleAssigned) // highest granule assigned
	mainItem0[027].SetH1(highestPositionWritten) // highest track written

	// Now populate DAD items for MFD
	err := mgr.writeFileAllocationEntryUpdates(mainAddr0)
	if err != nil {
		return err
	}

	// mark the sectors we've updated to be written
	mgr.markDirectorySectorDirty(leadAddr0)
	mgr.markDirectorySectorDirty(mainAddr0)
	mgr.markDirectorySectorDirty(mainAddr1)

	// Update lookup table
	mgr.writeLookupTableEntry("SYS$", "MFD$$", leadAddr0)

	// Set file assigned in facmgr, RCE, or wherever it makes sense
	// TODO - we should do this probably in the exec startup code, where we catalog or assign other system files.

	err = mgr.writeMFDCache()
	if err != nil {
		return err
	}

	log.Printf("MFDMgr:bootstrapMFD done")
	return nil
}

// checkCycle checks the given file cycle specification (if any) against the provided file set info
// to determine if a file cycle of the given spec can be cataloged.
// If fcSpecification is nil, the caller did not specify a cycle (which has its own meaning).
// Returned values include
//
//		  absoluteCycle the effective absolute cycle, which is either the given absolute cycle, or the
//			    absolute cycle deduced from the given relative cycle, or from the FileSetInfo if no cycle was given.
//		  cycleIndex the index into the lead item links where the new file cycle is to be placed
//			    zero corresponds to the highest absolute cycle position
//		  shiftAmount indicates whether, and by how many entries, the existing links in the lead item should be
//			    shifted (downward) to accommodate the new link
//		  newCycleRange indicates the resulting current cycle range presuming the file cycle is created.
//			    this may be greater than the current cycle range, but never less than, nor greater than the max.
//		  plusOne indicates that the newly-created file cycle is to be a plus-one cycle
//		  result indicates whether the check is successful, or if not, why not:
//		     MFDSuccessful if the checks succeed
//		     MFDInternalError if something is badly wrong, and we've stopped the exec
//		     MFDAlreadyExists if user does not specify a cycle, and a file cycle already exists
//		     MFDInvalidRelativeFileCycle if the caller specified a negative relative file cycle
//	      MFDInvalidAbsoluteFileCycle
//		        any file cycle out of range
//		        an absolute file cycle which conflicts with an existing cycle
//		     MFDPlusOneCycleExists if caller specifies +1 and a +1 already exists for the file set
//		     MFDDropOldestCycleRequired is returned if everything else would be fine if the oldest file cycle did not exist
func (mgr *MFDManager) checkCycle(
	fcSpecification *FileCycleSpecification,
	fsInfo *FileSetInfo,
) (absoluteCycle uint, cycleIndex uint, shiftAmount uint, newCycleRange uint, plusOne bool, result MFDResult) {
	absoluteCycle = 1
	cycleIndex = 0
	shiftAmount = 0
	newCycleRange = fsInfo.CurrentRange
	plusOne = false
	result = MFDSuccessful

	if fcSpecification == nil {
		// caller did not specify any file cycle.
		// either the file set is empty (caller just created it) and we're okay,
		// or it is not empty which requires us to reject the attempt.
		if fsInfo.CurrentRange != 0 {
			result = MFDAlreadyExists
		} else {
			newCycleRange = 1
		}
		return
	}

	if fcSpecification.IsRelative() {
		// We reject all negative cycles - if a negative cycle actually refers to a file cycle,
		// then we can't catalog it. If it doesn't, then we don't know which absolute cycle
		// it should refer to.
		// If the request is for +1, and it already exits, we reject that as well.
		if *fcSpecification.RelativeCycle < 0 {
			result = MFDInvalidRelativeFileCycle
			return
		}

		if fsInfo.PlusOneExists {
			result = MFDPlusOneCycleExists
			return
		}

		// If the file set is empty, the default values mostly suffice.
		if fsInfo.CurrentRange == 0 {
			newCycleRange = 1
			return
		}

		// If the highest file cycle is not actually there, the plus-one takes its place.
		// If it *is* there, we need to check whether the oldest cycle is in the way
		// (evidenced by current range == max range).
		// If all is well, set the absolute cycle to one more than the current highest value
		// then we're done (although calling code needs to be aware of the need to shift
		// the cycle links downward by one).
		if fsInfo.CycleInfo[0] == nil {
			absoluteCycle = fsInfo.HighestAbsolute
		} else if fsInfo.CurrentRange == fsInfo.MaxCycleRange {
			result = MFDDropOldestCycleRequired
		} else {
			absoluteCycle = fsInfo.HighestAbsolute + 1
			if absoluteCycle == 1000 {
				absoluteCycle = 1
			}
			shiftAmount = 1
			newCycleRange++
		}
		return
	}

	// The request is for an absolute file cycle.
	// If the default set is empty, use the requested absolute, and the other defaults are good.
	absoluteCycle = *fcSpecification.AbsoluteCycle
	if fsInfo.CurrentRange == 0 {
		newCycleRange = 1
		return
	}

	// Check whether the requested cycle is within the current range,
	// and if so, whether there is already a file cycle with that absolute cycle.
	chkCycle := fsInfo.HighestAbsolute
	for cx, cycleInfo := range fsInfo.CycleInfo {
		if chkCycle == absoluteCycle {
			if cycleInfo != nil {
				// there's already a cycle there - fail.
				result = MFDAlreadyExists
				return
			}

			// There is a hole here which we fit into nicely.
			cycleIndex = uint(cx)
			return
		}
	}

	// Is the request for a cycle above the highest file cycle?
	// If so, ensure that it is not so far above as to be out of range.
	// This test actually depends upon the lowest file cycle's position.
	if absoluteCycle > fsInfo.HighestAbsolute {
		var lowestAbsoluteCycle uint
		if fsInfo.HighestAbsolute < fsInfo.CurrentRange {
			lowestAbsoluteCycle = fsInfo.HighestAbsolute + 1000 - fsInfo.CurrentRange
		} else {
			lowestAbsoluteCycle = fsInfo.HighestAbsolute + 1 - fsInfo.CurrentRange
		}

		var requestedRange uint
		if absoluteCycle > lowestAbsoluteCycle {
			requestedRange = absoluteCycle - lowestAbsoluteCycle + 1
		} else {
			requestedRange = absoluteCycle + 1000 - lowestAbsoluteCycle
		}

		if requestedRange == fsInfo.MaxCycleRange+1 {
			result = MFDDropOldestCycleRequired
			return
		} else if requestedRange > fsInfo.MaxCycleRange {
			result = MFDInvalidAbsoluteFileCycle
			return
		}

		// The requested absolute is acceptable, however a shift is required
		shiftAmount = absoluteCycle - fsInfo.HighestAbsolute
		newCycleRange += shiftAmount

		// Temporary code - newRange should never be > maxRange...
		// but we'll check it just to make sure.
		if newCycleRange > fsInfo.MaxCycleRange {
			goto oops
		}

		return
	}

	// Note that cycles 1 through 31 are considered to be above cycles 968 through 999.
	// Check that as well.
	if absoluteCycle <= 31 && fsInfo.HighestAbsolute >= 968 {
		chkCycle := absoluteCycle + 999
		lowestAbsoluteCycle := fsInfo.HighestAbsolute + 1 - fsInfo.CurrentRange
		requestedRange := chkCycle - lowestAbsoluteCycle + 1

		if requestedRange == fsInfo.MaxCycleRange+1 {
			result = MFDDropOldestCycleRequired
			return
		} else if requestedRange > fsInfo.MaxCycleRange {
			result = MFDInvalidAbsoluteFileCycle
			return
		}

		// The requested absolute is acceptable, however a shift is required
		shiftAmount = chkCycle - fsInfo.HighestAbsolute
		newCycleRange += shiftAmount

		// Temporary code - newRange should never be > maxRange...
		// but we'll check it just to make sure.
		if newCycleRange > fsInfo.MaxCycleRange {
			goto oops
		}

		return
	}

	// Is the request for a cycle below the highest file cycle?
	// If so, ensure it is not so far below as to be out of range.
	if absoluteCycle < fsInfo.HighestAbsolute {
		cycleIndex = fsInfo.HighestAbsolute - absoluteCycle
		if cycleIndex >= fsInfo.MaxCycleRange {
			result = MFDInvalidAbsoluteFileCycle
			return
		}

		newCycleRange = cycleIndex + 1

		// Temporary code - newRange should never be > maxRange...
		// but we'll check it just to make sure.
		if newCycleRange > fsInfo.MaxCycleRange {
			goto oops
		}

		return
	}

	// Note that cycles 968 through 999 are below cycles 1 through 31.
	// Check that as well.
	if absoluteCycle >= 968 && fsInfo.HighestAbsolute <= 31 {
		cycleIndex = fsInfo.HighestAbsolute + 999 - absoluteCycle
		if cycleIndex >= fsInfo.MaxCycleRange {
			result = MFDInvalidAbsoluteFileCycle
			return
		}

		newCycleRange = cycleIndex + 1

		// Temporary code - newRange should never be > maxRange...
		// but we'll check it just to make sure.
		if newCycleRange > fsInfo.MaxCycleRange {
			goto oops
		}

		return
	}

	// The requested absolute is out of range
	result = MFDInvalidAbsoluteFileCycle
	return

oops:
	log.Printf("MFDMgr:(a)newRange is %v which is more than max range %v", newCycleRange, fsInfo.MaxCycleRange)
	mgr.exec.Stop(kexec.StopDirectoryErrors)
	result = MFDInternalError
	return
}

func composeMFDAddress(
	ldatIndex kexec.LDATIndex,
	trackId kexec.MFDTrackId,
	sectorId kexec.MFDSectorId,
) kexec.MFDRelativeAddress {

	return kexec.MFDRelativeAddress(uint64(ldatIndex&07777)<<18 | uint64(trackId&07777)<<6 | uint64(sectorId&077))
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
	fileTrackId kexec.TrackId,
) (kexec.LDATIndex, kexec.TrackId, error) {

	fae, ok := mgr.assignedFileAllocations[mainItem0Address]
	if !ok {
		log.Printf("MFDMgr:convertFileRelativeTrackId Cannot find fae for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return 0, 0, fmt.Errorf("fae not loaded")
	}

	ldat := kexec.LDATIndex(0_400000)
	devTrackId := kexec.TrackId(0)
	if fileTrackId <= fae.HighestTrackAllocated {
		for _, fileAlloc := range fae.FileAllocations {
			if fileTrackId < fileAlloc.FileRegion.TrackId {
				// list is ascending - if we get here, there's no point in continuing
				break
			}
			upperLimit := kexec.TrackId(uint64(fileAlloc.FileRegion.TrackId) + uint64(fileAlloc.FileRegion.TrackCount))
			if fileTrackId < upperLimit {
				// found a good region - update results and stop looking
				ldat = fileAlloc.LDATIndex
				devTrackId = fileAlloc.DeviceTrackId + (fileTrackId - fileAlloc.FileRegion.TrackId)
				return ldat, devTrackId, nil
			}
		}
	}

	return ldat, devTrackId, nil
}

// dropFileCycle this is where the hard work gets done for dropping a file cycle.
// We put it here because it is invoked by both DropFileCycle() and DropFileSet().
// Caller must be very sure that this file is not assigned.
// If we return anything other than MFDSuccessful, the exec will already be stopped.
func (mgr *MFDManager) dropFileCycle(
	mainItem0Address kexec.MFDRelativeAddress,
) MFDResult {
	_, ok := mgr.assignedFileAllocations[mainItem0Address]
	if ok {
		log.Printf("MFDMgr:Attempt to drop assigned file cycle mainItem0: %012o", mainItem0Address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return MFDInternalError
	}

	mainItem0, err := mgr.getMFDSector(mainItem0Address)
	if err != nil {
		return MFDInternalError
	}

	if mainItem0[014].GetT1()&0040 == 0 {
		// This is mass storage. Release the DAD tables, marking the allocated tracks as unallocated.
		dadAddress := kexec.MFDRelativeAddress(mainItem0[0])
		for dadAddress&0_400000_000000 == 0 {
			dadAddress &= 0_007777_777777
			dadSector, err := mgr.getMFDSector(dadAddress)
			if err != nil {
				return MFDInternalError
			}

			for wx := 07; wx < 28; wx += 3 {
				ldat := kexec.LDATIndex(dadSector[wx+2].GetH2())
				if ldat != kexec.InvalidLDAT {
					pd, ok := mgr.fixedPackDescriptors[ldat]
					if ok {
						devTrackId := kexec.TrackId(uint64(dadSector[wx].GetW()) / 1792)
						devTrackCount := kexec.TrackCount(uint64(dadSector[wx+1].GetW()) / 1792)
						pd.freeSpaceTable.MarkTrackRegionUnallocated(devTrackId, devTrackCount)
					}
				}
			}

			err = mgr.markDirectorySectorUnallocated(dadAddress)
			if err != nil {
				return MFDInternalError
			}

			dadAddress = kexec.MFDRelativeAddress(dadSector[0])
		}
	} else {
		// This is a tape file - release the reel number tables
		rtAddress := kexec.MFDRelativeAddress(mainItem0[0])
		for rtAddress&0_400000_000000 == 0 {
			rtAddress &= 0_007777_777777
			rtSector, err := mgr.getMFDSector(rtAddress)
			if err != nil {
				return MFDInternalError
			}

			err = mgr.markDirectorySectorUnallocated(rtAddress)
			if err != nil {
				return MFDInternalError
			}

			rtAddress = kexec.MFDRelativeAddress(rtSector[0])
		}
	}

	// Release main items
	err = mgr.markDirectorySectorUnallocated(mainItem0Address)
	if err != nil {
		return MFDInternalError
	}

	mainItemAddress := kexec.MFDRelativeAddress(mainItem0[015] & 0_007777_777777)
	for mainItemAddress&0_400000_000000 == 0 {
		mainItemAddress &= 0_007777_777777
		mainItem, err := mgr.getMFDSector(mainItemAddress)
		if err != nil {
			return MFDInternalError
		}

		err = mgr.markDirectorySectorUnallocated(mainItemAddress)
		if err != nil {
			return MFDInternalError
		}

		mainItemAddress = kexec.MFDRelativeAddress(mainItem[0].GetW())
	}

	return MFDSuccessful
}

// findDASEntryForSector chases the appropriate DAS chain to find the DAS which describes the given sector address,
// and then the entry within that DAS which describes the sector address.
// Returns:
//
//	the address of the containing DAS sector
//	the index of the DAS entry
//	a slice to the 3-word DAS entry itself.
//
// If we return an error, we have already stopped the exec
func (mgr *MFDManager) findDASEntryForSector(
	sectorAddr kexec.MFDRelativeAddress,
) (kexec.MFDRelativeAddress, int, []pkg.Word36, error) {

	// what are we looking for?
	ldat := getLDATIndexFromMFDAddress(sectorAddr)
	trackId := getMFDTrackIdFromMFDAddress(sectorAddr)

	dasAddr := composeMFDAddress(ldat, 0, 0)
	for dasAddr != kexec.InvalidLink {
		das, err := mgr.getMFDSector(dasAddr)
		if err != nil {
			log.Printf("MFDMgr:findDASEntrySector cannot compose MFD address for %012o", sectorAddr)
			mgr.exec.Stop(kexec.StopDirectoryErrors)
			return 0, 0, nil, err
		}

		// The first entry of a DAS describes the track which contains the DAS.
		// Is it this one?
		firstTrackId := getMFDTrackIdFromMFDAddress(dasAddr)
		if firstTrackId == trackId {
			return dasAddr, 0, das[0:3], nil
		}

		// Look at the other entries
		for ex := 1; ex < 8; ex++ {
			ey := ex * 3
			entryAddr := kexec.MFDRelativeAddress(das[ey].GetW())
			if entryAddr != kexec.InvalidLink {
				entryTrackId := getMFDTrackIdFromMFDAddress(entryAddr)
				if entryTrackId == trackId {
					// found it.
					return dasAddr, ex, das[ey : ey+3], nil
				}
			}
		}

		// So it is not this DAS - move on to the next
		dasAddr = kexec.MFDRelativeAddress(das[033].GetW())
	}

	// We did not find the DAS entry - complain and crash
	log.Printf("MFDMgr:Cannot find DAS for sector %012o", sectorAddr)
	mgr.exec.Stop(kexec.StopDirectoryErrors)
	return 0, 0, nil, fmt.Errorf("cannot find DAS")
}

func getLDATIndexFromMFDAddress(address kexec.MFDRelativeAddress) kexec.LDATIndex {
	return kexec.LDATIndex(address>>18) & 07777
}

// getLeadItemLinkWord retrieves a pointer to a file cycle link from among the given lead item(s)
// corresponding to the file cycle's position indicated by the linkIndex value.
// linkIndex == 0 corresponds to the highest absolute cycle link.
// If there is no lead item 1, pass nil for that value, and make dang sure you don't specify
// an invalid link index.
func getLeadItemLinkWord(leadItem0 []pkg.Word36, leadItem1 []pkg.Word36, linkIndex uint) *pkg.Word36 {
	headerCount := uint(11 + leadItem0[012].GetS4())
	if linkIndex+headerCount < 28 {
		return &leadItem0[linkIndex+headerCount]
	} else {
		return &leadItem1[(linkIndex+headerCount)-28]
	}
}

func getMFDTrackIdFromMFDAddress(address kexec.MFDRelativeAddress) kexec.MFDTrackId {
	return kexec.MFDTrackId(address>>6) & 07777
}

func getMFDSectorIdFromMFDAddress(address kexec.MFDRelativeAddress) kexec.MFDSectorId {
	return kexec.MFDSectorId(address & 077)
}

// getMFDAddressForBlock takes a given MFD-relative sector address and normalizes it to
// the first sector in the block containing the given sector.
// CALL UNDER LOCK
func (mgr *MFDManager) getMFDAddressForBlock(address kexec.MFDRelativeAddress) kexec.MFDRelativeAddress {
	ldat := getLDATIndexFromMFDAddress(address)
	mask := uint64(mgr.fixedPackDescriptors[ldat].packMask)
	return kexec.MFDRelativeAddress(uint64(address) & ^mask)
}

// getMFDBlock returns a slice corresponding to all the sectors in the physical block
// containing the sector represented by the given address. Used for reading/writing MFD blocks.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) getMFDBlock(address kexec.MFDRelativeAddress) ([]pkg.Word36, error) {
	ldatAndTrack := address & 0_007777_777700
	data, ok := mgr.cachedTracks[ldatAndTrack]
	if !ok {
		log.Printf("MFDMgr:getMFDBlock address:%v is not in cache", address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return nil, fmt.Errorf("internal error")
	}

	ldat := getLDATIndexFromMFDAddress(address)
	sectorId := getMFDSectorIdFromMFDAddress(address)
	mask := uint64(mgr.fixedPackDescriptors[ldat].packMask)
	baseSectorId := uint64(sectorId) & ^mask

	start := 28 * baseSectorId
	end := start + uint64(mgr.fixedPackDescriptors[ldat].prepFactor)
	return data[start:end], nil
}

// getMFDSector returns a slice corresponding to the portion of the MFD block which represents the indicated sector.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) getMFDSector(address kexec.MFDRelativeAddress) ([]pkg.Word36, error) {
	ldatAndTrack := address & 0_007777_777700
	data, ok := mgr.cachedTracks[ldatAndTrack]
	if !ok {
		log.Printf("MFDMgr:getMFDSector address:%v is not in cache", address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return nil, fmt.Errorf("internal error")
	}

	sectorId := getMFDSectorIdFromMFDAddress(address)
	start := 28 * sectorId
	end := start + 28
	return data[start:end], nil
}

// initializeFixed initializes the fixed pool for a jk13 boot
func (mgr *MFDManager) initializeFixed(disks map[*nodeMgr.DiskDeviceInfo]*kexec.DiskAttributes) error {
	msg := fmt.Sprintf("Fixed Disk Pool = %v Devices", len(disks))
	mgr.exec.SendExecReadOnlyMessage(msg, nil)

	if len(disks) == 0 {
		return nil
	}

	replies := []string{"Y", "N"}
	msg = "Mass Storage will be Initialized - Do You Want To Continue? Y/N"
	reply, err := mgr.exec.SendExecRestrictedReadReplyMessage(msg, replies)
	if err != nil {
		return err
	} else if reply != "Y" {
		mgr.exec.Stop(kexec.StopConsoleResponseRequiresReboot)
		return fmt.Errorf("boot canceled")
	}

	// make sure there are no pack name conflicts
	conflicts := false
	packNames := make([]string, 0)
	for _, diskAttr := range disks {
		for _, existing := range packNames {
			if diskAttr.PackLabelInfo.PackId == existing {
				msg := fmt.Sprintf("Fixed pack name conflict - %v", existing)
				mgr.exec.SendExecReadOnlyMessage(msg, nil)
				conflicts = true
			} else {
				packNames = append(packNames, diskAttr.PackLabelInfo.PackId)
			}
		}
	}

	if conflicts {
		mgr.exec.SendExecReadOnlyMessage("Resolve pack name conflicts and reboot", nil)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return fmt.Errorf("packid conflict")
	}

	nm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)

	// iterate over the fixed packs - we start at 1, which may not be conventional, but it works
	nextLdatIndex := kexec.LDATIndex(1)
	totalTracks := uint64(0)
	for diskInfo, diskAttr := range disks {
		// Assign an LDAT to the pack, update the pack label, then rewrite the label
		ldatIndex := nextLdatIndex
		nextLdatIndex++

		// Assign the unit
		_ = mgr.exec.GetFacilitiesManager().AssignDiskDeviceToExec(diskInfo.GetNodeIdentifier())

		// Set up fixed pack descriptor
		fpDesc := newPackDescriptor(diskInfo.GetNodeIdentifier(), diskAttr)

		// Mark VOL1 and first directory track as allocated
		fpDesc.mfdTrackCount = 1
		fpDesc.mfdSectorsUsed = 2
		_ = fpDesc.freeSpaceTable.AllocateSpecificTrackRegion(ldatIndex, 0, 2)

		// Set up first directory track for the pack within our cache to be eventually rewritten
		mgr.fixedPackDescriptors[ldatIndex] = fpDesc

		availableTracks := diskAttr.PackLabelInfo.TrackCount
		recordLength := diskAttr.PackLabelInfo.WordsPerRecord
		blocksPerTrack := diskAttr.PackLabelInfo.RecordsPerTrack
		wordsPerBlock := diskAttr.PackLabelInfo.WordsPerRecord

		// We need to read the first directory track into cache
		// so we can update sector 0 and sector 1 appropriately.
		mfdTrackId := kexec.MFDTrackId(ldatIndex << 12)
		mfdAddr := composeMFDAddress(ldatIndex, mfdTrackId, 0)
		data := make([]pkg.Word36, 1792)
		mgr.cachedTracks[mfdAddr] = data

		// read the directory track into cache
		blockId := kexec.BlockId(blocksPerTrack)
		wx := uint(0)
		for bx := 0; bx < int(blocksPerTrack); bx++ {
			sub := data[wx : wx+wordsPerBlock]
			ioPkt := nodeMgr.NewDiskIoPacketRead(fpDesc.nodeId, blockId, sub)
			nm.RouteIo(ioPkt)
			ioStat := ioPkt.GetIoStatus()
			if ioStat != nodeMgr.IosComplete {
				log.Printf("MFDMgr:initializeFixed cannot read directory track dev:%v blockId:%v",
					fpDesc.nodeId, blockId)
				mgr.exec.Stop(kexec.StopInternalExecIOFailed)
				return fmt.Errorf("IO error")
			}

			blockId++
			wx += wordsPerBlock
		}

		// sector 0
		sector0Addr := composeMFDAddress(ldatIndex, 0, 0)
		sector0, err := mgr.getMFDSector(sector0Addr)
		if err != nil {
			return err
		}

		sector0[0].SetH1(uint64(ldatIndex))
		sector0[0].SetH2(0)
		sector0[1].SetW(0_600000_000000) // first 2 sectors are allocated
		for dx := 3; dx < 27; dx += 3 {
			sector0[dx].SetW(uint64(kexec.InvalidLink))
			sector0[dx+1].SetW(0)
			sector0[dx+2].SetW(0)
		}
		sector0[27].SetW(0_400000_000000)
		mgr.markDirectorySectorDirty(sector0Addr)

		// sector 1
		sector1Addr := composeMFDAddress(ldatIndex, 0, 1)
		sector1, err := mgr.getMFDSector(sector1Addr)
		if err != nil {
			return err
		}

		// leave +0 and +1 alone (We aren't doing HMBT/SMBT so we don't need the addresses)
		sector1[2].SetW(uint64(availableTracks))
		sector1[3].SetW(uint64(availableTracks))
		sector1[4].FromStringToFieldata(diskAttr.PackLabelInfo.PackId)
		sector1[5].SetH1(uint64(ldatIndex))
		sector1[010].SetT1(uint64(blocksPerTrack))
		sector1[010].SetS3(1) // Sector 1 version
		sector1[010].SetT3(uint64(recordLength))
		mgr.markDirectorySectorDirty(sector1Addr)

		totalTracks += uint64(availableTracks)
	}

	err = mgr.bootstrapMFD()
	if err != nil {
		return err
	}

	msg = fmt.Sprintf("MS Initialized - %v Tracks Available", totalTracks)
	mgr.exec.SendExecReadOnlyMessage(msg, nil)
	return nil
}

// initializeRemovable registers the isRemovable packs (if any) as part of a JK13 boot.
func (mgr *MFDManager) initializeRemovable(disks map[*nodeMgr.DiskDeviceInfo]*kexec.DiskAttributes) error {
	// TODO
	return nil
}

// loadFileAllocationEntry initializes the fae for a particular file instance.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) loadFileAllocationEntry(
	mainItem0Address kexec.MFDRelativeAddress,
) (*FileAllocationSet, error) {
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
	fae := &FileAllocationSet{
		DadItem0Address:  dadAddr,
		MainItem0Address: mainItem0Address,
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
				re := NewFileAllocation(kexec.TrackId(fileWordAddress/1792),
					kexec.TrackCount(words/1792),
					ldat,
					kexec.TrackId(devAddr/1792))
				fae.MergeIntoFileAllocationSet(re)
			}
			ex++
			dx++
		}

		dadAddr = kexec.MFDRelativeAddress(dadItem[0].GetW())
	}

	mgr.assignedFileAllocations[mainItem0Address] = fae
	return fae, nil
}

// markDirectorySectorAllocated finds the appropriate DAS entry for the given sector address
// and marks the sector as allocated, as well as marking the DAS entry as updated.
// If we return an error, we have already stopped the exec
func (mgr *MFDManager) markDirectorySectorAllocated(sectorAddr kexec.MFDRelativeAddress) error {
	dasAddr, _, entry, err := mgr.findDASEntryForSector(sectorAddr)
	if err != nil {
		return err
	}

	sectorId := getMFDSectorIdFromMFDAddress(sectorAddr)
	if sectorId < 32 {
		mask := uint64(0_400000_000000) >> sectorId
		entry[1].Or(mask)
	} else {
		mask := uint64(0_400000_000000) >> (sectorId - 32)
		entry[2].Or(mask)
	}

	mgr.markDirectorySectorDirty(dasAddr)
	return nil
}

// markDirectorySectorDirty marks the block which contains the given sector as dirty,
// so that it can subsequently be written to storage.
func (mgr *MFDManager) markDirectorySectorDirty(address kexec.MFDRelativeAddress) {
	blockAddr := mgr.getMFDAddressForBlock(address)
	mgr.dirtyBlocks[blockAddr] = true
}

// markDirectorySectorUnallocated finds the appropriate DAS entry for the given sector address
// and marks the sector as unallocated, as well as marking the DAS entry as updated.
// If we return an error, we have already stopped the exec
func (mgr *MFDManager) markDirectorySectorUnallocated(sectorAddr kexec.MFDRelativeAddress) error {
	dasAddr, _, entry, err := mgr.findDASEntryForSector(sectorAddr)
	if err != nil {
		return err
	}

	sectorId := getMFDSectorIdFromMFDAddress(sectorAddr)
	if sectorId < 32 {
		mask := uint64(0_400000_000000) >> sectorId
		entry[1].And(^mask)
	} else {
		mask := uint64(0_400000_000000) >> (sectorId - 32)
		entry[2].And(^mask)
	}

	mgr.markDirectorySectorDirty(dasAddr)
	return nil
}

func populateFixedMainItem1(
	mainItem1 []pkg.Word36,
	qualifier string,
	filename string,
	mainItem0Address kexec.MFDRelativeAddress,
	absoluteCycle uint64,
	packIds []string,
) {
	for wx := 0; wx < 28; wx++ {
		mainItem1[wx].SetW(0)
	}

	mainItem1[0].SetW(uint64(kexec.InvalidLink)) // no sector 2 (yet, anyway)
	pkg.FromStringToFieldata(qualifier, mainItem1[1:3])
	pkg.FromStringToFieldata(filename, mainItem1[3:5])
	pkg.FromStringToFieldata("*No.1*", mainItem1[5:6])
	mainItem1[6].SetW(uint64(mainItem0Address))
	mainItem1[7].SetT3(absoluteCycle)

	// TODO note that for >5 pack entries, we need additional main item sectors
	//  one per 10 additional packs beyond 5
	mix := 18
	limit := len(packIds)
	if limit > 5 {
		limit = 5
	}
	for dpx := 0; dpx < limit; dpx++ {
		mainItem1[mix].FromStringToFieldata(packIds[dpx])
		mix += 2
	}
}

func populateMassStorageMainItem0(
	mainItem0 []pkg.Word36,
	leadItem0Address kexec.MFDRelativeAddress,
	mainItem1Address kexec.MFDRelativeAddress,
	qualifier string,
	filename string,
	absoluteCycle uint64,
	readKey string,
	writeKey string,
	projectId string,
	accountId string,
	mnemonic string,
	descriptorFlags DescriptorFlags,
	pcharFlags PCHARFlags,
	inhibitFlags InhibitFlags,
	isRemovable bool,
	reserve uint64,
	maximum uint64,
	packIds []string,
) {
	for wx := 0; wx < 28; wx++ {
		mainItem0[wx].SetW(0)
	}

	mainItem0[0].SetW(uint64(kexec.InvalidLink)) // no DAD table (yet, anyway)
	mainItem0[0].Or(0_200000_000000)

	pkg.FromStringToFieldata(qualifier, mainItem0[1:3])
	pkg.FromStringToFieldata(filename, mainItem0[3:5])
	pkg.FromStringToFieldata(projectId, mainItem0[5:7])
	pkg.FromStringToFieldata(accountId, mainItem0[7:9])
	mainItem0[11].SetW(uint64(leadItem0Address))
	mainItem0[11].SetS1(0) // disable flags

	mainItem0[12].SetT1(descriptorFlags.Compose())

	mainItem0[13].SetW(uint64(mainItem1Address))
	mainItem0[13].SetS1(pcharFlags.Compose())

	mainItem0[14].FromStringToFieldata(mnemonic)

	mainItem0[17].SetH1(inhibitFlags.Compose())
	mainItem0[17].SetT3(absoluteCycle)

	swTimeNow := kexec.GetSWTimeFromSystemTime(time.Now())
	mainItem0[19].SetW(swTimeNow)
	mainItem0[20].SetH1(reserve)
	mainItem0[21].SetH1(maximum)

	if isRemovable {
		var rKey pkg.Word36
		if len(readKey) > 0 {
			rKey.FromStringToFieldata(readKey)
		}
		var wKey pkg.Word36
		if len(writeKey) > 0 {
			wKey.FromStringToAscii(writeKey)
		}

		mainItem0[24].SetH1(rKey.GetH1())
		mainItem0[25].SetH1(rKey.GetH2())
		mainItem0[26].SetH1(wKey.GetH1())
		mainItem0[27].SetH1(wKey.GetH2())
	} else {
		// initially selected LDAT and optional device placement flag
		// TODO if there is at least one pack-id, then go find its LDAT and use that,
		//  and mask in 0_400000_000000 to indicate device placement.
		var ldat uint64
		if len(packIds) > 0 {

		} else {
			ldat = uint64(getLDATIndexFromMFDAddress(leadItem0Address))
		}
		mainItem0[27].SetH1(ldat)
	}
}

// populateNewLeadItem0 sets up a lead item sector 0 in the given buffer,
// assuming we are cataloging a new file, will have one cycle, and the absolute cycle is given to us.
// Implied is that there will be no sector 1 (since there aren't enough cycles to warrant it).
func populateNewLeadItem0(
	leadItem0 []pkg.Word36,
	qualifier string,
	filename string,
	absoluteCycle uint64,
	readKey string,
	writeKey string,
	projectId string,
	fileType uint64, // 000=Fixed, 001=Tape, 040=Removable
	guardedFlag bool,
	mainItem0Address uint64,
) {
	for wx := 0; wx < 28; wx++ {
		leadItem0[wx].SetW(0)
	}

	leadItem0[0].SetW(uint64(kexec.InvalidLink))
	leadItem0[0].Or(0_500000_000000)

	pkg.FromStringToFieldata(qualifier, leadItem0[1:3])
	pkg.FromStringToFieldata(filename, leadItem0[3:5])
	pkg.FromStringToFieldata(projectId, leadItem0[5:7])
	if len(readKey) > 0 {
		leadItem0[7].FromStringToFieldata(readKey)
	}
	if len(writeKey) > 0 {
		leadItem0[8].FromStringToAscii(writeKey)
	}

	leadItem0[9].SetS1(fileType)
	leadItem0[9].SetS2(1)  // number of cycles
	leadItem0[9].SetS3(31) // max range of cycles (default is 31)
	leadItem0[9].SetS4(1)  // current range
	leadItem0[9].SetT3(absoluteCycle)

	var statusBits uint64
	if guardedFlag {
		statusBits |= 01000
	}
	leadItem0[10].SetT1(statusBits)
	leadItem0[11].SetW(mainItem0Address)
}

// TODO
//func populateRemovableMainItem1(
//	mainItem1 []pkg.Word36,
//	mainItem0Address kexec.MFDRelativeAddress,
//	absoluteCycle uint64,
//	packIds []string,
//) {
//	for wx := 0; wx < 28; wx++ {
//		mainItem1[wx].SetW(0)
//	}
//
//	mainItem1[0].SetW(uint64(kexec.InvalidLink)) // no sector 2 (yet, anyway)
//	mainItem1[6].SetW(uint64(mainItem0Address))
//	mainItem1[7].SetT3(absoluteCycle)
//	mainItem1[17].SetT3(uint64(len(packIds)))
//
//	// TODO note that for >5 pack entries, we need additional main item sectors
//	//  one per 10 additional packs beyond 5
//	mix := 18
//	limit := len(packIds)
//	if limit > 5 {
//		limit = 5
//	}
//	for dpx := 0; dpx < limit; dpx++ {
//		mainItem1[mix].FromStringToFieldata(packIds[dpx])
//		// TODO for isRemovable, we need the main item address for this file on that pack
//		mix += 2
//	}
//}

// TODO
//func populateTapeMainItem0(
//	mainItem0 []pkg.Word36,
//	qualifier string,
//	filename string,
//	projectId string,
//	accountId string,
//	reelTable0Address kexec.MFDRelativeAddress,
//	leadItem0Address kexec.MFDRelativeAddress,
//	mainItem1Address kexec.MFDRelativeAddress,
//	toBeCataloged bool, // for @ASG,C or @ASG,U
//	isGuarded bool,
//	isPrivate bool,
//	isWriteOnly bool,
//	isReadOnly bool,
//	absoluteCycle uint64,
//	density uint,
//	format uint,
//	features uint,
//	featuresExtension uint,
//	mtapop uint,
//	ctlPool string,
//) {
//	for wx := 0; wx < 28; wx++ {
//		mainItem0[wx].SetW(0)
//	}
//
//	mainItem0[0].SetW(uint64(reelTable0Address))
//	mainItem0[0].Or(0_200000_000000)
//	pkg.FromStringToFieldata(qualifier, mainItem0[1:3])
//	pkg.FromStringToFieldata(filename, mainItem0[3:5])
//	pkg.FromStringToFieldata(projectId, mainItem0[5:7])
//	pkg.FromStringToFieldata(accountId, mainItem0[7:9])
//
//	// TODO
//}

// TODO
//func populateTapeMainItem1(
//	mainItem1 []pkg.Word36,
//	qualifier string,
//	filename string,
//	mainItem0Address kexec.MFDRelativeAddress,
//	absoluteCycle uint64,
//) {
//	for wx := 0; wx < 28; wx++ {
//		mainItem1[wx].SetW(0)
//	}
//
//	mainItem1[0].SetW(uint64(kexec.InvalidLink)) // no sector 2 (yet, anyway)
//	pkg.FromStringToFieldata(qualifier, mainItem1[1:3])
//	pkg.FromStringToFieldata(filename, mainItem1[3:5])
//	pkg.FromStringToFieldata("*No.1*", mainItem1[5:6])
//	mainItem1[6].SetW(uint64(mainItem0Address))
//	mainItem1[7].SetT3(absoluteCycle)
//}

// releaseTrackRegion releases the indicated tracks from the file cycle's DAD entries,
// and adds the corresponding device tracks to the free space table.
// Only intended to be invoked on files which are currently assigned.
// If we return anything other than MFDSuccessful, the exec is stopped.
func (mgr *MFDManager) releaseFileCycleTrackRegion(
	mainItem0Address kexec.MFDRelativeAddress,
	region *kexec.TrackRegion,
) MFDResult {
	fas, ok := mgr.assignedFileAllocations[mainItem0Address]
	if !ok {
		log.Printf("MFDMgr:convertFileRelativeTrackId Cannot find alloc for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return MFDInternalError
	}

	ldat, devTrackId := fas.ExtractRegionFromFileAllocationSet(region)
	if ldat == kexec.InvalidLDAT {
		log.Printf("MFDMgr:convertFileRelativeTrackId Cannot extract alloc for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return MFDInternalError
	}

	return mgr.releaseTrackRegion(ldat, kexec.NewTrackRegion(devTrackId, region.TrackCount))
}

// releaseTrackRegion adds tracks for a particular pack to the free space table.
// Only intended to be invoked on files which are *not* currently assigned, *OR* by releaseFileCycleTrackRegion().
// If we return anything other than MFDSuccessful, the exec is stopped.
func (mgr *MFDManager) releaseTrackRegion(
	ldatIndex kexec.LDATIndex,
	region *kexec.TrackRegion,
) MFDResult {
	packDesc, ok := mgr.fixedPackDescriptors[ldatIndex]
	if !ok {
		log.Printf("MFDMgr:Cannot find pack descriptor for LDAT %v", ldatIndex)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return MFDInternalError
	}

	if !packDesc.freeSpaceTable.MarkTrackRegionUnallocated(region.TrackId, region.TrackCount) {
		log.Printf("MFDMgr:Cannot unallocation region LDAT %v TrkId %v Count %v",
			ldatIndex, region.TrackId, region.TrackCount)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return MFDInternalError
	}

	return MFDSuccessful
}

// writeFileAllocationEntryUpdates writes an updated fae to the on-disk MFD
// as a series of one or more DAD tables.
// MUST be invoked when a file is free'd.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) writeFileAllocationEntryUpdates(mainItem0Address kexec.MFDRelativeAddress) error {
	fas, ok := mgr.assignedFileAllocations[mainItem0Address]
	if !ok {
		log.Printf("MFDMgr:convertFileRelativeTrackId Cannot find alloc for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec.StopDirectoryErrors)
		return fmt.Errorf("fas not loaded")
	}

	if fas.IsUpdated {
		// TODO process
		//	rewrite the entire set of DAD entries, allocate a new one if we need to do so,
		//  and release any left over when we're done.
		//  Don't forget to write hole DADs (see pg 2-63 for Device index field)

		fas.IsUpdated = false
	}

	return nil
}

func (mgr *MFDManager) writeLookupTableEntry(
	qualifier string,
	filename string,
	leadItem0Addr kexec.MFDRelativeAddress) {

	_, ok := mgr.fileLeadItemLookupTable[qualifier]
	if !ok {
		mgr.fileLeadItemLookupTable[qualifier] = make(map[string]kexec.MFDRelativeAddress)
	}
	mgr.fileLeadItemLookupTable[qualifier][filename] = leadItem0Addr
}

// writeMFDCache writes all the dirty cache blocks to storage.
// If we return error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) writeMFDCache() error {
	for blockAddr := range mgr.dirtyBlocks {
		block, err := mgr.getMFDBlock(blockAddr)
		if err != nil {
			log.Printf("MFDMgr:writeMFDCache cannot find MFD block for dirty block address:%012o", blockAddr)
			mgr.exec.Stop(kexec.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		}

		mfdTrackId := (blockAddr >> 6) & 077777777
		mfdSectorId := getMFDSectorIdFromMFDAddress(blockAddr)

		ldat, devTrackId, err := mgr.convertFileRelativeTrackId(mgr.mfdFileMainItem0Address, kexec.TrackId(mfdTrackId))
		if err != nil {
			log.Printf("MFDMgr:writeMFDCache error converting mfdaddr:%012o TrackId:%06v", mgr.mfdFileMainItem0Address, mfdTrackId)
			mgr.exec.Stop(kexec.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		} else if ldat == 0_400000 {
			log.Printf("MFDMgr:writeMFDCache error converting mfdaddr:%012o TrackId:%06v track not allocated",
				mgr.mfdFileMainItem0Address, mfdTrackId)
			mgr.exec.Stop(kexec.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		}

		packDesc, ok := mgr.fixedPackDescriptors[ldat]
		if !ok {
			log.Printf("MFDMgr:writeMFDCache cannot find packDesc for ldat:%04v", ldat)
			mgr.exec.Stop(kexec.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		}

		blocksPerTrack := 1792 / packDesc.prepFactor
		sectorsPerBlock := packDesc.prepFactor / 28
		devBlockId := uint64(devTrackId) * uint64(blocksPerTrack)
		devBlockId += uint64(mfdSectorId) / uint64(sectorsPerBlock)
		ioPkt := nodeMgr.NewDiskIoPacketWrite(packDesc.nodeId, kexec.BlockId(devBlockId), block)
		nm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
		nm.RouteIo(ioPkt)
		ioStat := ioPkt.GetIoStatus()
		if ioStat != nodeMgr.IosComplete {
			log.Printf("MFDMgr:writeMFDCache error writing MFD block status=%v", ioStat)
			mgr.exec.Stop(kexec.StopInternalExecIOFailed)
			return fmt.Errorf("error draining MFD cache")
		}
	}

	mgr.dirtyBlocks = make(map[kexec.MFDRelativeAddress]bool)
	return nil
}
