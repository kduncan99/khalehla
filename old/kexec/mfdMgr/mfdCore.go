// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// we are allowed to know about facMgr, nodeMgr, consMgr

import (
	"fmt"
	"math/rand"
	"time"

	"khalehla/hardware"
	"khalehla/hardware/channels"
	"khalehla/hardware/ioPackets"
	"khalehla/kexec"
	"khalehla/logger"
	ioPackets2 "khalehla/old/hardware/ioPackets"
	kexec2 "khalehla/old/kexec"
	"khalehla/old/kexec/nodeMgr"
	"khalehla/old/pkg"
)

// adjustLeadItemLinks shifts the links in the given lead item(s) downward by the indicated
// shift value (presumably to accommodate a new higher-cycled link).
// Adjusts the current file cycle count in leadItem0 accordingly.
// Make darn sure there is a leadItem1 if it is necessary to contain the expanded set of links.
// Otherwise, it can be nil.
// Caller must ensure both leadItem0 and leadItem1 (if it exists) are marked dirty - we don't do it.
// Caller must update current range, count, and highest fcycle fields in the lead item.
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
// Note that, apart from observing the LDAT portion, we do not select the sectors in any particular order.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) allocateDirectorySector(
	preferredLDAT kexec2.LDATIndex,
) (kexec2.MFDRelativeAddress, []pkg.Word36, error) {
	logger.LogTraceF("MFDMgr", "allocateDirectorySector(ldat=%06o)", preferredLDAT)

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
	chosenPackSectorsUsed := uint64(0)
	for addr := range mgr.freeMFDSectors {
		ldat := getLDATIndexFromMFDAddress(addr)
		desc := mgr.fixedPackDescriptors[ldat]
		freeCount := (uint64(desc.mfdTrackCount) * 64) - desc.mfdSectorsUsed

		if freeCount > 0 {
			if ldat == preferredLDAT {
				chosenAddress = addr
				break
			}

			if chosenAddress == kexec.InvalidLink || desc.mfdSectorsUsed < chosenPackSectorsUsed {
				chosenAddress = addr
				chosenPackSectorsUsed = desc.mfdSectorsUsed
			}
		}
	}

	// When we get here, we *will* have valid chosen elements.
	// If we had no free sectors anywhere, we would have allocated a new track
	// (and if that failed, we would already have crashed and returned)
	delete(mgr.freeMFDSectors, chosenAddress)
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

	logger.LogTraceF("MFDMgr", "allocateDirectorySector returns %012o", chosenAddress)
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
	preferredLDAT kexec2.LDATIndex,
) (kexec2.LDATIndex, kexec2.MFDTrackId, error) {
	logger.LogTraceF("MFDMgr", "allocateDirectoryTrack(ldat=%v)", preferredLDAT)

	chosenLDAT := kexec.InvalidLDAT
	var chosenDesc *packDescriptor
	chosenAvailableTracks := hardware.TrackCount(0)

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
		logger.LogFatal("MFDMgr", "No space available for directory track allocation")
		mgr.exec.Stop(kexec2.StopExecRequestForMassStorageFailed)
		return 0, 0, fmt.Errorf("no disk")
	}

	// First, find the MFD relative address of the first unused MFD track
	trackId := kexec2.MFDTrackId(0)
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

	logger.LogInfoF("MFDMgr", "allocateDirectorySector returns ldat=%06o track=%012o", chosenLDAT, trackId)
	return chosenLDAT, trackId, nil
}

// allocateFileSpace attempts to allocate tracks for the given file cycle,
// at the specific file-relative track id, for the specified number of tracks.
// Caller must ensure that the requested space does not overlap already-allocated
// file-relative space.
// Caller must ensure that the file allocations are accelerated into memory.
// Caller must allocate space according to granularity -- we do not impose this requirement here.
func (mgr *MFDManager) allocateFileSpace(
	mainItem0Address kexec2.MFDRelativeAddress,
	fileRelativeTrackId hardware.TrackId,
	trackCount hardware.TrackCount,
) (result MFDResult) {
	result = MFDSuccessful

	faSet, ok := mgr.acceleratedFileAllocations[mainItem0Address]
	if !ok {
		logger.LogFatalF("MFDMgr", "allocateFileSpace main item %12o is not accelerated into memory", mainItem0Address)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return MFDInternalError
	}

	highestTrack := faSet.GetHighestTrackAssigned()

	// Find the existing allocation immediately preceding this one.
	// If there is no such allocation (we are a sparse file or an empty one)
	// just find space preferably on the pack containing the main item.
	// Otherwise, try to extend the allocation.
	frTrackId := fileRelativeTrackId
	remaining := trackCount
	fa := faSet.FindPrecedingAllocation(fileRelativeTrackId)
	if fa != nil {
		pDesc, ok := mgr.fixedPackDescriptors[fa.LDATIndex]
		if !ok {
			logger.LogFatalF("MFDMgr", "allocateFileSpace no packDesc for LDAT %06o", fa.LDATIndex)
			mgr.exec.Stop(kexec2.StopDirectoryErrors)
			return MFDInternalError
		}

		baseTrackId := fa.DeviceTrackId + hardware.TrackId(fa.FileRegion.TrackCount)
		allocated := pDesc.freeSpaceTable.AllocateTracksFromTrackId(baseTrackId, remaining)
		if allocated > 0 {
			alloc := kexec2.NewFileAllocation(frTrackId, allocated, fa.LDATIndex, baseTrackId)
			faSet.MergeIntoFileAllocationSet(alloc)
			frTrackId += hardware.TrackId(allocated)
			remaining -= allocated
		}
	}

	if remaining > 0 {
		preferred := getLDATIndexFromMFDAddress(mainItem0Address)
		var packRegions []*kexec2.PackRegion
		packRegions, result = mgr.allocateSpace(preferred, remaining)
		if result != MFDSuccessful {
			return
		}

		for _, region := range packRegions {
			alloc := kexec2.NewFileAllocation(frTrackId, region.TrackCount, region.LDATIndex, region.TrackId)
			faSet.MergeIntoFileAllocationSet(alloc)
			frTrackId += hardware.TrackId(region.TrackCount)
		}
	}

	// Update highest granule assigned in main item
	allocHighestTrack := fileRelativeTrackId + hardware.TrackId(trackCount) - 1
	if allocHighestTrack > highestTrack {
		mainItem0, err := mgr.getMFDSector(mainItem0Address)
		if err != nil {
			return MFDInternalError
		}
		posGran := mainItem0[015].GetS1()&040 != 0
		if posGran {
			highestGranule := allocHighestTrack >> 6
			mainItem0[026].SetH1(uint64(highestGranule))
		} else {
			mainItem0[026].SetH1(uint64(allocHighestTrack))
		}
		mgr.markDirectorySectorDirty(mainItem0Address)
	}

	faSet.IsUpdated = true
	return
}

// allocateLeadItem1 allocates a lead item directory sector which extends a
// currently-not-extended lead item 0 such that it can refer to a larger file cycle range.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) allocateLeadItem1(
	leadItem0Address kexec2.MFDRelativeAddress,
	leadItem0 []pkg.Word36,
) (leadItem1Address kexec2.MFDRelativeAddress, leadItem1 []pkg.Word36, err error) {
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

// allocateSpace attempts to allocate the indicated number of tracks, preferring the pack identified
// by the given LDAT index. We will allocate as much of the request as possible in one block, but we will
// allocate in multiple blocks if necessary, and from multiple packs if necessary.
// We prioritize keeping content for a file on one particular pack, over keeping the various packs de-fragmented.
// Due to restricted mass storage space, we may not be able to satisfy the entire (or any of the) request.
func (mgr *MFDManager) allocateSpace(
	preferred kexec2.LDATIndex,
	trackCount hardware.TrackCount,
) (regions []*kexec2.PackRegion, result MFDResult) {
	regions = make([]*kexec2.PackRegion, 0)
	result = MFDSuccessful

	remaining := trackCount
	for remaining > 0 {
		pDesc, ok := mgr.fixedPackDescriptors[preferred]
		if ok {
			region := pDesc.freeSpaceTable.AllocateTrackRegion(remaining)
			if region == nil {
				break
			}

			regions = append(regions, kexec2.NewPackRegion(preferred, region.TrackId, region.TrackCount))
			remaining -= region.TrackCount
			continue
		}
	}

	// randomize the pack descriptors, then allocate from them, as much from each individual pack as can be done.
	randomMap := make(map[int]*packDescriptor)
	for _, pDesc := range mgr.fixedPackDescriptors {
		randomMap[rand.Int()] = pDesc
	}

	for _, pDesc := range randomMap {
		for remaining > 0 {
			region := pDesc.freeSpaceTable.AllocateTrackRegion(remaining)
			if region == nil {
				break
			}

			regions = append(regions, kexec2.NewPackRegion(preferred, region.TrackId, region.TrackCount))
			remaining -= region.TrackCount
			continue
		}
	}

	if remaining > 0 {
		for _, region := range regions {
			pDesc, _ := mgr.fixedPackDescriptors[region.LDATIndex]
			pDesc.freeSpaceTable.MarkTrackRegionUnallocated(region.TrackId, region.TrackCount)
		}
		result = MFDOutOfSpace
		regions = nil
	}

	return
}

// allocateSpecificTrack allocates particular contiguous specified physical tracks
// to be associated with the indicated file-relative tracks.
// If we return an error, we've already stopped the exec
// ONLY FOR VERY SPECIFIC USE-CASES - CALL UNDER LOCK, OR DURING MFD INIT
func (mgr *MFDManager) allocateSpecificTrack(
	mainItem0Address kexec2.MFDRelativeAddress,
	fileTrackId hardware.TrackId,
	trackCount hardware.TrackCount,
	ldatIndex kexec2.LDATIndex,
	deviceTrackId hardware.TrackId) error {

	fas, ok := mgr.acceleratedFileAllocations[mainItem0Address]
	if !ok {
		logger.LogFatalF("MFDMgr", "allocateSpecificTrack Cannot find fas for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return fmt.Errorf("fas not loaded")
	}

	fileAlloc := kexec2.NewFileAllocation(fileTrackId, trackCount, ldatIndex, deviceTrackId)
	fas.MergeIntoFileAllocationSet(fileAlloc)
	fas.IsUpdated = true

	return nil
}

// bootstrapMFD creates the various MFD structures as part of MFD initialization.
// One consequence is the cataloging of SYS$*MFD$$.
// Since this is used during initialization we do not call it under lock.
func (mgr *MFDManager) bootstrapMFD() error {
	logger.LogTrace("MFDMgr", "bootstrapMFD start")

	cfg := mgr.exec.GetConfiguration()

	// Find the highest and lowest LDAT indices.
	// While we're here, pre-populate the free MFD sectors map for the first directory track on each pack
	var lowestLDAT = kexec.InvalidLDAT
	var highestLDAT = kexec2.LDATIndex(0)
	for ldat := range mgr.fixedPackDescriptors {
		if lowestLDAT == kexec.InvalidLDAT && ldat < lowestLDAT {
			lowestLDAT = ldat
		}
		if ldat > highestLDAT {
			highestLDAT = ldat
		}

		for sx := 2; sx < 64; sx++ {
			sectorAddr := composeMFDAddress(ldat, 0, kexec2.MFDSectorId(sx))
			mgr.freeMFDSectors[sectorAddr] = true
		}
	}

	// Allocate MFD sectors for MFD$$ file items not including DAD tables (we do those separately)
	leadItem0Addr, leadItem0, _ := mgr.allocateDirectorySector(lowestLDAT)
	mainItem0Addr, mainItem0, _ := mgr.allocateDirectorySector(lowestLDAT)
	mainItem1Addr, mainItem1, _ := mgr.allocateDirectorySector(lowestLDAT)

	mgr.mfdFileMainItem0Address = mainItem0Addr // we'll need this later

	// Manually catalog the MFD file
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
		uint64(mainItem0Addr))

	populateMassStorageMainItem0(mainItem0,
		leadItem0Addr,
		mainItem1Addr,
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
			IsAssignedExclusive: false,
			IsWriteOnly:         false,
			IsReadOnly:          false,
		},
		false,
		0,
		262153)

	populateFixedMainItem1(mainItem1, mfdFileQualifier, mfdFileName, mainItem0Addr, 1)

	// Before we can play DAD table games, we have to get the MFD$$ in-core structures in place,
	// including *particularly* the file allocation table.
	// We need to create one allocation region for each pack's initial directory track.
	highestMFDTrackId := hardware.TrackId(0)
	fas := kexec2.NewFileAllocationSet(mainItem0Addr, 0_400000_000000)
	mgr.acceleratedFileAllocations[mainItem0Addr] = fas

	for ldat, desc := range mgr.fixedPackDescriptors {
		mfdTrackId := hardware.TrackId(ldat << 12)
		if mfdTrackId > highestMFDTrackId {
			highestMFDTrackId = mfdTrackId
		}

		packTrackId := desc.firstDirectoryTrackAddress / 1792
		err := mgr.allocateSpecificTrack(mainItem0Addr, mfdTrackId, 1, ldat, hardware.TrackId(packTrackId))
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
	err := mgr.writeFileAllocationEntryUpdatesForFileCycle(mainItem0Addr)
	if err != nil {
		return err
	}

	// mark the sectors we've updated to be written
	mgr.markDirectorySectorDirty(leadItem0Addr)
	mgr.markDirectorySectorDirty(mainItem0Addr)
	mgr.markDirectorySectorDirty(mainItem1Addr)

	// Update lookup table
	mgr.writeLookupTableEntry("SYS$", "MFD$$", leadItem0Addr)

	err = mgr.writeMFDCache()
	if err != nil {
		return err
	}

	logger.LogTrace("MFDMgr", "bootstrapMFD done")
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
	fcSpecification *kexec2.FileCycleSpecification,
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
	logger.LogFatalF("MFDMgr", "newRange is %v which is more than max range %v", newCycleRange, fsInfo.MaxCycleRange)
	mgr.exec.Stop(kexec2.StopDirectoryErrors)
	result = MFDInternalError
	return
}

// composeMFDAddress creates an MFDRelativeAddress from its component parts
func composeMFDAddress(
	ldatIndex kexec2.LDATIndex,
	trackId kexec2.MFDTrackId,
	sectorId kexec2.MFDSectorId,
) kexec2.MFDRelativeAddress {

	return kexec2.MFDRelativeAddress(uint64(ldatIndex&07777)<<18 | uint64(trackId&07777)<<6 | uint64(sectorId&077))
}

// convertFileRelativeAddress takes a file-relative track-id (i.e., word offset from start of file divided by 1792)
// and uses the fae entries in the fat for the given file instance to determine the device LDAT and
// the device-relative track id which contains that file address.
// If the logical track is not allocated, we will return 0_400000 and 0 for those values (since 0 is an invalid LDAT index)
// If the fae is not loaded, we will throw an error - even an empty file has an fae, albeit a puny one.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) convertFileRelativeTrackId(
	mainItem0Address kexec2.MFDRelativeAddress,
	fileTrackId hardware.TrackId,
) (kexec2.LDATIndex, hardware.TrackId, error) {
	logger.LogTraceF("MFDMgr", "convertFileRelativeTrackId(mainItem0Addr=%012o fileTid=%012o", mainItem0Address, fileTrackId)

	fas, ok := mgr.acceleratedFileAllocations[mainItem0Address]
	if !ok {
		logger.LogFatalF("MFDMgr", "convertFileRelativeTrackId Cannot find fas for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return 0, 0, fmt.Errorf("fas not loaded")
	}

	ldat := kexec2.LDATIndex(0_400000)
	devTrackId := hardware.TrackId(0)
	highestAllocated, hasHighest := fas.GetHighestTrackAllocated()
	if hasHighest && fileTrackId <= highestAllocated {
		for _, fileAlloc := range fas.FileAllocations {
			if fileTrackId < fileAlloc.FileRegion.TrackId {
				// list is ascending - if we get here, there's no point in continuing
				break
			}
			upperLimit := hardware.TrackId(uint64(fileAlloc.FileRegion.TrackId) + uint64(fileAlloc.FileRegion.TrackCount))
			if fileTrackId < upperLimit {
				// found a good region - update results and stop looking
				ldat = fileAlloc.LDATIndex
				devTrackId = fileAlloc.DeviceTrackId + (fileTrackId - fileAlloc.FileRegion.TrackId)
				logger.LogTraceF("MFDMgr", "convertFileRelativeTrackId returning ldat=%06o devTid=%012o", ldat, devTrackId)
				return ldat, devTrackId, nil
			}
		}
	}

	logger.LogTraceF("MFDMgr", "convertFileRelativeTrackId returning ldat=%06o devTid=%012o", ldat, devTrackId)
	return ldat, devTrackId, nil
}

// dropFileCycle this is where the hard work gets done for dropping a file cycle.
// Caller must be very sure that this file is accelerated into our table.
// If we return anything other than MFDSuccessful, the exec will already be stopped.
func (mgr *MFDManager) dropFileCycle(
	mainItem0Address kexec2.MFDRelativeAddress,
) MFDResult {
	fas, ok := mgr.acceleratedFileAllocations[mainItem0Address]
	if !ok {
		logger.LogFatalF("MFDMgr", "Attempt to drop non-accelerated file cycle mainItem0: %012o", mainItem0Address)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return MFDInternalError
	}

	mainItem0, err := mgr.getMFDSector(mainItem0Address)
	if err != nil {
		return MFDInternalError
	}

	// We'll need the lead item(s) sooner, and later.
	leadItem0Addr, leadItem1Addr, leadItem0, leadItem1, err := mgr.getLeadItemsForMainItem(mainItem0)
	if err != nil {
		return MFDInternalError
	}

	fileType := NewFileTypeFromField(leadItem0[011].GetS1())
	if fileType == FileTypeFixed || fileType == FileTypeRemovable {
		// For mass storage, we need to release the space occupied by the file...
		for _, fa := range fas.FileAllocations {
			tr := kexec2.NewTrackRegion(fa.DeviceTrackId, fa.FileRegion.TrackCount)
			mgr.releaseTrackRegion(fa.LDATIndex, tr)
		}

		// ... and then the DAD table entries themselves.
		err = mgr.releaseDADChain(mainItem0Address)
		if err != nil {
			return MFDInternalError
		}
	} else if fileType == FileTypeTape {
		// For tape, just release the reel tables (if any)
		err = mgr.releaseReelNumberChain(mainItem0Address)
		if err != nil {
			return MFDInternalError
		}
	}

	// Release main items
	addresses := []kexec2.MFDRelativeAddress{mainItem0Address}
	nextAddr := kexec2.MFDRelativeAddress(mainItem0[015].GetW() & 0_007777_777777)
	for nextAddr != 0 {
		addresses = append(addresses, nextAddr)
		sector, err := mgr.getMFDSector(nextAddr)
		if err != nil {
			return MFDInternalError
		}

		nextAddr = kexec2.MFDRelativeAddress(sector[0].GetW() & 0_007777_777777)
	}

	// Was this the only cycle in the file set?
	cycleCount := uint(leadItem0[011].GetS2())
	if cycleCount == 1 {
		// Yes - drop the file set and un-accelerate it
		_ = mgr.markDirectorySectorUnallocated(leadItem0Addr)
		_ = mgr.markDirectorySectorUnallocated(leadItem1Addr)
	} else {
		// No - update the file cycle links in the lead item(s)
		leadItems := [][]pkg.Word36{leadItem0, leadItem1}
		highestCycle := uint(leadItem0[011].GetT3())
		mainItemCycle := uint(mainItem0[021].GetT3())
		newCycleCount := 0
		cycleRange := uint(leadItem0[011].GetS4())
		linkIndex := highestCycle - mainItemCycle
		ix, wx := getLeadItemLinkIndices(leadItem0, linkIndex)
		leadItems[ix][wx].SetW(0)
		for chkLinkIndex := uint(0); chkLinkIndex < cycleRange; chkLinkIndex++ {
			link := getLeadItemLinkWord(leadItem0, leadItem1, chkLinkIndex)
			if !link.IsNegative() {
				newCycleCount++
			}
			if link.GetW()&0_007777_777777 != 0 {
				cycleRange = chkLinkIndex + 1
			}
		}
		leadItem0[011].SetS2(uint64(newCycleCount))
		leadItem0[011].SetS4(uint64(cycleRange))

		mgr.markDirectorySectorDirty(leadItem0Addr)
		if leadItem1 != nil {
			mgr.markDirectorySectorDirty(leadItem1Addr)
		}
	}

	// decelerate the file cycle
	delete(mgr.acceleratedFileAllocations, mainItem0Address)
	return MFDSuccessful
}

// findDASEntryForSector chases the appropriate DAS chain to find the DAS which describes the given sector address,
// and then the entry within that DAS which describes the sector address.
// Returns
//
//	the address of the containing DAS sector
//	the index of the DAS entry
//	a slice to the 3-word DAS entry itself.
//
// If we return an error, we have already stopped the exec
func (mgr *MFDManager) findDASEntryForSector(
	sectorAddr kexec2.MFDRelativeAddress,
) (kexec2.MFDRelativeAddress, int, []pkg.Word36, error) {

	// what are we looking for?
	ldat := getLDATIndexFromMFDAddress(sectorAddr)
	trackId := getMFDTrackIdFromMFDAddress(sectorAddr)

	dasAddr := composeMFDAddress(ldat, 0, 0)
	for dasAddr != kexec.InvalidLink {
		das, err := mgr.getMFDSector(dasAddr)
		if err != nil {
			logger.LogFatalF("MFDMgr", "findDASEntrySector cannot compose MFD address for %012o", sectorAddr)
			mgr.exec.Stop(kexec2.StopDirectoryErrors)
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
			entryAddr := kexec2.MFDRelativeAddress(das[ey].GetW())
			if entryAddr != kexec.InvalidLink {
				entryTrackId := getMFDTrackIdFromMFDAddress(entryAddr)
				if entryTrackId == trackId {
					// found it.
					return dasAddr, ex, das[ey : ey+3], nil
				}
			}
		}

		// So it is not this DAS - move on to the next
		dasAddr = kexec2.MFDRelativeAddress(das[033].GetW())
	}

	// We did not find the DAS entry - complain and crash
	logger.LogFatalF("MFDMgr", "Cannot find DAS for sector %012o", sectorAddr)
	mgr.exec.Stop(kexec2.StopDirectoryErrors)
	return 0, 0, nil, fmt.Errorf("cannot find DAS")
}

// getLDATIndexFromMFDAddress extracts the LDAT value of an MFD relative address
func getLDATIndexFromMFDAddress(address kexec2.MFDRelativeAddress) kexec2.LDATIndex {
	return kexec2.LDATIndex(address>>18) & 07777
}

// getLeadItems retrieves and returns the lead item(s) for which the address for
// lead item 0 has been provided.
// If there is no lead item 1, then its returned address will be InvalidLink and the sector will be nil.
// If we return an error, the exec is already stopped
func (mgr *MFDManager) getLeadItems(leadItem0Addr kexec2.MFDRelativeAddress) (
	leadItem1Address kexec2.MFDRelativeAddress,
	leadItem0 []pkg.Word36,
	leadItem1 []pkg.Word36,
	err error,
) {
	leadItem1Address = kexec.InvalidLink

	leadItem0, err = mgr.getMFDSector(leadItem0Addr)
	if err != nil {
		return
	}

	if !leadItem0[0].IsNegative() {
		leadItem1Address = kexec2.MFDRelativeAddress(leadItem0[0].GetW() & 0_007777_777777)
		leadItem1, err = mgr.getMFDSector(leadItem1Address)
	}

	return
}

// getLeadItemsForMainItem retrieves and returns the lead item(s) associated with a given main item,
// along with their MFD-relative addresses.
// If there is no lead item 1, then the address will be InvalidLink and the sector will be nil.
// If we return an error, the exec is already stopped
func (mgr *MFDManager) getLeadItemsForMainItem(mainItem []pkg.Word36) (
	leadItem0Address kexec2.MFDRelativeAddress,
	leadItem1Address kexec2.MFDRelativeAddress,
	leadItem0 []pkg.Word36,
	leadItem1 []pkg.Word36,
	err error,
) {
	leadItem1Address = kexec.InvalidLink

	leadItem0Address = kexec2.MFDRelativeAddress(mainItem[013].GetW() & 0_007777_777777)
	leadItem0, err = mgr.getMFDSector(leadItem0Address)
	if err != nil {
		return
	}

	if !leadItem0[0].IsNegative() {
		leadItem1Address = kexec2.MFDRelativeAddress(leadItem0[0].GetW() & 0_007777_777777)
		leadItem1, err = mgr.getMFDSector(leadItem1Address)
	}

	return
}

// getLeadItemsForMainItemAddress - as above, but accepts the MFD-relative address of the
// main item instead of the item itself.
func (mgr *MFDManager) getLeadItemsForMainItemAddress(mainItem0Address kexec2.MFDRelativeAddress) (
	leadItem0Address kexec2.MFDRelativeAddress,
	leadItem1Address kexec2.MFDRelativeAddress,
	leadItem0 []pkg.Word36,
	leadItem1 []pkg.Word36,
	err error,
) {
	mainItem0, err := mgr.getMFDSector(mainItem0Address)
	if err != nil {
		return
	} else {
		return mgr.getLeadItemsForMainItem(mainItem0)
	}
}

// getLeadItemLinkIndices takes an index into the file cycle table (where 0 corresponds to the highest absolute cycle)
// and returns an index to the lead item which contains the corresponding file cycle link (0 for lead item 0,
// 1 for lead item 1), and a word index indicating the word offset from the start of the lead item which contains
// the link.
func getLeadItemLinkIndices(leadItem0 []pkg.Word36, cycleIndex uint) (itemIndex uint, wordIndex uint) {
	secWords := uint(leadItem0[012].GetS4())
	entriesInItem0 := 28 - (11 + secWords)
	if cycleIndex < entriesInItem0 {
		itemIndex = 0
		wordIndex = cycleIndex
	} else {
		itemIndex = 1
		wordIndex = cycleIndex - entriesInItem0
	}
	return
}

// getLeadItemLinkIndicesForCycle -- as above, but given an absolute cycle instead of a cycle index
func getLeadItemLinkIndicesForCycle(
	leadItem0 []pkg.Word36,
	absoluteCycle uint,
) (itemIndex uint, wordIndex uint, ok bool) {
	itemIndex = 0
	wordIndex = 0
	ok = true

	highest := uint(leadItem0[011].GetT3())
	var cycIndex uint
	if absoluteCycle < highest {
		cycIndex = highest - (absoluteCycle + 999)
	} else {
		cycIndex = highest - absoluteCycle
	}

	cycRange := uint(leadItem0[011].GetS4())
	if cycIndex >= cycRange {
		ok = false
	} else {
		itemIndex, wordIndex = getLeadItemLinkIndices(leadItem0, cycIndex)
	}

	return
}

// getLeadItemLinkWord retrieves a pointer to a file cycle link from among the given lead item(s)
// corresponding to the file cycle's position indicated by the linkIndex value.
// linkIndex == 0 corresponds to the highest absolute cycle link.
// If there is no lead item 1, pass nil for that value, and make dang sure you don't specify
// an invalid link index. Returns false if something is obviously amiss.
func getLeadItemLinkWord(
	leadItem0 []pkg.Word36,
	leadItem1 []pkg.Word36,
	linkIndex uint,
) *pkg.Word36 {
	ix, wx := getLeadItemLinkIndices(leadItem0, linkIndex)
	if ix == 0 {
		return &leadItem0[wx]
	} else {
		return &leadItem1[wx]
	}
}

// getLeadItemLinkWordForCycle -- as above, but given an absolute file cycle instead of a link index
func getLeadItemLinkWordForCycle(
	leadItem0 []pkg.Word36,
	leadItem1 []pkg.Word36,
	absoluteCycle uint,
) (word *pkg.Word36, ok bool) {
	ix, wx, ok := getLeadItemLinkIndicesForCycle(leadItem0, absoluteCycle)
	if ok {
		if ix == 0 {
			word = &leadItem0[wx]
		} else {
			word = &leadItem1[wx]
		}
	}
	return
}

// getMFDTrackIdFromMFDAddress extracts the MFD-relative track-id portion of an MFD relative address
func getMFDTrackIdFromMFDAddress(address kexec2.MFDRelativeAddress) kexec2.MFDTrackId {
	return kexec2.MFDTrackId(address>>6) & 07777
}

// getMFDSectorIdFromMFDAddress extracts the MFD-relative sector-id portion of an MFD relative address
func getMFDSectorIdFromMFDAddress(address kexec2.MFDRelativeAddress) kexec2.MFDSectorId {
	return kexec2.MFDSectorId(address & 077)
}

// getMFDAddressForBlock takes a given MFD-relative sector address and normalizes it to
// the first sector in the block containing the given sector.
// CALL UNDER LOCK
func (mgr *MFDManager) getMFDAddressForBlock(address kexec2.MFDRelativeAddress) kexec2.MFDRelativeAddress {
	ldat := getLDATIndexFromMFDAddress(address)
	mask := uint64(mgr.fixedPackDescriptors[ldat].packMask)
	return kexec2.MFDRelativeAddress(uint64(address) & ^mask)
}

// getMFDBlock returns a slice corresponding to all the sectors in the physical block
// containing the sector represented by the given address. Used for reading/writing MFD blocks.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) getMFDBlock(address kexec2.MFDRelativeAddress) ([]pkg.Word36, error) {
	ldatAndTrack := address & 0_007777_777700
	data, ok := mgr.cachedTracks[ldatAndTrack]
	if !ok {
		logger.LogFatalF("MFDMgr", "getMFDBlock address:%012o is not in cache", address)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
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
func (mgr *MFDManager) getMFDSector(address kexec2.MFDRelativeAddress) ([]pkg.Word36, error) {
	ldatAndTrack := address & 0_007777_777700
	data, ok := mgr.cachedTracks[ldatAndTrack]
	if !ok {
		logger.LogFatalF("MFDMgr", "getMFDSector address:%012o is not in cache", address)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return nil, fmt.Errorf("internal error")
	}

	sectorId := getMFDSectorIdFromMFDAddress(address)
	start := 28 * sectorId
	end := start + 28
	return data[start:end], nil
}

func (mgr *MFDManager) getPackDescriptorForNodeIdentifier(
	nodeIdentifier hardware.NodeIdentifier,
) (kexec2.LDATIndex, *packDescriptor, bool) {
	for ldat, packDesc := range mgr.fixedPackDescriptors {
		if packDesc.nodeId == nodeIdentifier {
			return ldat, packDesc, true
		}
	}
	return 0, nil, false
}

// initializeFixed initializes the fixed pool for a jk13 boot
func (mgr *MFDManager) initializeFixed(disks map[*nodeMgr.DiskDeviceInfo]*kexec2.DiskAttributes) error {
	msg := fmt.Sprintf("Fixed Disk Pool = %v Devices", len(disks))
	mgr.exec.SendExecReadOnlyMessage(msg, nil)

	if len(disks) == 0 {
		return nil
	}

	replies := []string{"Y", "N"}
	msg = "Mass Storage will be Initialized - Do You Want To Continue? Y/N"
	reply, err := mgr.exec.SendExecRestrictedReadReplyMessage(msg, replies, nil)
	if err != nil {
		return err
	} else if reply != "Y" {
		mgr.exec.Stop(kexec2.StopConsoleResponseRequiresReboot)
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
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return fmt.Errorf("packid conflict")
	}

	// iterate over the fixed packs - we start at 1, which may not be conventional, but it works
	nextLdatIndex := kexec2.LDATIndex(1)
	totalTracks := uint64(0)
	for diskInfo, diskAttr := range disks {
		// Assign an LDAT to the pack, update the pack label, then rewrite the label
		ldatIndex := nextLdatIndex
		nextLdatIndex++

		// Set up fixed pack descriptor
		nodeId := diskInfo.GetNodeIdentifier()
		fpDesc := newPackDescriptor(
			nodeId,
			diskAttr.PackLabelInfo.PrepFactor,
			diskAttr.PackLabelInfo.TrackCount,
			diskAttr.PackLabelInfo.FirstDirectoryTrackAddress,
			diskAttr.GetFacNodeStatus())

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
		mfdTrackId := kexec2.MFDTrackId(ldatIndex << 12)
		mfdAddr := composeMFDAddress(ldatIndex, mfdTrackId, 0)
		data := make([]pkg.Word36, 1792)
		mgr.cachedTracks[mfdAddr] = data

		// read the directory track into cache
		blockId := hardware.BlockId(blocksPerTrack)
		wx := uint(0)
		for bx := 0; bx < int(blocksPerTrack); bx++ {
			sub := data[wx : wx+wordsPerBlock]
			ioStat := mgr.readBlockFromDisk(fpDesc.nodeId, sub, blockId)
			if ioStat == ioPackets2.IosInternalError {
				return fmt.Errorf("init stopped")
			} else if ioStat != ioPackets2.IosComplete {
				logger.LogFatalF("MFDMgr", "initializeFixed cannot read directory track dev:%v blockId:%v",
					fpDesc.nodeId, blockId)
				mgr.exec.Stop(kexec2.StopInternalExecIOFailed)
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
func (mgr *MFDManager) initializeRemovable(disks map[*nodeMgr.DiskDeviceInfo]*kexec2.DiskAttributes) error {
	// TODO implement initializeRemovable()
	return nil
}

// loadFileAllocationSet initializes the fas for a particular file instance.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) loadFileAllocationSet(
	mainItem0Address kexec2.MFDRelativeAddress,
) (*kexec2.FileAllocationSet, error) {
	mainItem0, err := mgr.getMFDSector(mainItem0Address)
	if err != nil {
		return nil, err
	}

	dadAddr := kexec2.MFDRelativeAddress(mainItem0[0])
	fae := &kexec2.FileAllocationSet{
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
		dx := 4
		for ex < 8 && fileWordAddress < fileWordLimit {
			devAddr := kexec2.DeviceRelativeWordAddress(dadItem[dx].GetW())
			words := dadItem[dx+1].GetW()
			ldat := kexec2.LDATIndex(dadItem[dx+2].GetH2())
			if ldat != 0_400000 {
				re := kexec2.NewFileAllocation(
					hardware.TrackId(fileWordAddress/1792),
					hardware.TrackCount(words/1792),
					ldat,
					hardware.TrackId(devAddr/1792))
				fae.MergeIntoFileAllocationSet(re)
			}

			fileWordAddress += words
			ex++
			dx += 3
		}

		dadAddr = kexec2.MFDRelativeAddress(dadItem[0].GetW())
	}

	mgr.acceleratedFileAllocations[mainItem0Address] = fae
	return fae, nil
}

// markDirectorySectorAllocated finds the appropriate DAS entry for the given sector address
// and marks the sector as allocated, as well as marking the DAS entry as updated.
// If we return an error, we have already stopped the exec
func (mgr *MFDManager) markDirectorySectorAllocated(sectorAddr kexec2.MFDRelativeAddress) error {
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

	ldat := getLDATIndexFromMFDAddress(sectorAddr)
	mgr.fixedPackDescriptors[ldat].mfdSectorsUsed++

	mgr.markDirectorySectorDirty(dasAddr)
	return nil
}

// markDirectorySectorDirty marks the block which contains the given sector as dirty,
// so that it can subsequently be written to storage.
func (mgr *MFDManager) markDirectorySectorDirty(address kexec2.MFDRelativeAddress) {
	blockAddr := mgr.getMFDAddressForBlock(address)
	mgr.dirtyBlocks[blockAddr] = true
}

// markDirectorySectorUnallocated finds the appropriate DAS entry for the given sector address
// and marks the sector as unallocated, as well as marking the DAS entry as updated.
// If we return an error, we have already stopped the exec
func (mgr *MFDManager) markDirectorySectorUnallocated(sectorAddr kexec2.MFDRelativeAddress) error {
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

	ldat := getLDATIndexFromMFDAddress(sectorAddr)
	mgr.fixedPackDescriptors[ldat].mfdSectorsUsed--

	mgr.markDirectorySectorDirty(dasAddr)
	return nil
}

func populateFixedMainItem1(
	mainItem1 []pkg.Word36,
	qualifier string,
	filename string,
	mainItem0Address kexec2.MFDRelativeAddress,
	absoluteCycle uint64,
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
}

func populateMassStorageMainItem0(
	mainItem0 []pkg.Word36,
	leadItem0Address kexec2.MFDRelativeAddress,
	mainItem1Address kexec2.MFDRelativeAddress,
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
	mainItem0[013].SetW(uint64(leadItem0Address))
	mainItem0[013].SetS1(0) // disable flags

	mainItem0[014].SetT1(descriptorFlags.Compose())

	mainItem0[015].SetW(uint64(mainItem1Address))
	mainItem0[015].SetS1(pcharFlags.Compose())

	mainItem0[016].FromStringToFieldata(mnemonic)

	mainItem0[021].SetS2(inhibitFlags.Compose())
	mainItem0[021].SetT3(absoluteCycle)

	swTimeNow := kexec.GetSWTimeFromSystemTime(time.Now())
	mainItem0[023].SetW(swTimeNow)
	mainItem0[024].SetH1(reserve)
	mainItem0[025].SetH1(maximum)

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

// readBlockFromDisk reads a single block from a particular disk device
// and sleeps until the IO is complete.
func (mgr *MFDManager) readBlockFromDisk(
	nodeId hardware.NodeIdentifier,
	buffer []pkg.Word36,
	blockId hardware.BlockId,
) ioPackets2.IoStatus {
	cw := channels.ControlWord{
		Buffer:    buffer,
		Offset:    0,
		Length:    uint(len(buffer)),
		Direction: channels.DirectionForward,
		Format:    channels.TransferPacked,
	}
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeId,
		IoFunction:     ioPackets.IofRead,
		BlockId:        blockId,
		ControlWords:   []channels.ControlWord{cw},
	}

	mgr.exec.GetNodeManager().RouteIo(cp)
	for cp.IoStatus == ioPackets2.IosInProgress || cp.IoStatus == ioPackets2.IosNotStarted {
		time.Sleep(10 * time.Millisecond)
	}

	return cp.IoStatus
}

// releaseDADChain releases all the DAD entry sectors attached to a particular main item
func (mgr *MFDManager) releaseDADChain(
	mainItem0Address kexec2.MFDRelativeAddress,
) error {
	mainItem0, err := mgr.getMFDSector(mainItem0Address)
	if err != nil {
		return err
	}

	// only do the work if there actually is a DAD chain
	if !mainItem0[0].IsNegative() {
		// walk the chain, collecting the sector addresses of all the entries in the chain
		addresses := make([]kexec2.MFDRelativeAddress, 0)
		dadAddr := kexec2.MFDRelativeAddress(mainItem0[0].GetW() & 0_007777_777777)
		for dadAddr != 0 {
			addresses = append(addresses, dadAddr)
			dad, err := mgr.getMFDSector(dadAddr)
			if err != nil {
				return err
			}

			dadAddr = 0
			if !dad[0].IsNegative() {
				dadAddr = kexec2.MFDRelativeAddress(dad[0].GetW() & 0_007777_777777)
			}
		}

		// clear the DAD link in the main item sector, and mark it dirty
		mainItem0[0].SetW((mainItem0[0].GetW() & 0_340000_000000) | 0_400000_000000)
		mgr.markDirectorySectorDirty(mainItem0Address)

		// release the DAD sectors
		for _, address := range addresses {
			err = mgr.markDirectorySectorUnallocated(address)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// releaseReelNumberChain releases all the reel number sectors attached to a particular tape main item
func (mgr *MFDManager) releaseReelNumberChain(
	mainItem0Address kexec2.MFDRelativeAddress,
) error {
	mainItem0, err := mgr.getMFDSector(mainItem0Address)
	if err != nil {
		return err
	}

	// only do the work if there actually is a reel number chain
	if !mainItem0[0].IsNegative() {
		// walk the chain, collecting the sector addresses of all the entries in the chain
		addresses := make([]kexec2.MFDRelativeAddress, 0)
		reelAddr := kexec2.MFDRelativeAddress(mainItem0[0].GetW() & 0_007777_777777)
		for reelAddr != 0 {
			addresses = append(addresses, reelAddr)
			entry, err := mgr.getMFDSector(reelAddr)
			if err != nil {
				return err
			}

			reelAddr = 0
			if !entry[0].IsNegative() {
				reelAddr = kexec2.MFDRelativeAddress(entry[0].GetW() & 0_007777_777777)
			}
		}

		// clear the reel table link in the main item sector, and mark it dirty
		mainItem0[0].SetW((mainItem0[0].GetW() & 0_340000_000000) | 0_400000_000000)
		mgr.markDirectorySectorDirty(mainItem0Address)

		// release the reel number table sectors
		for _, address := range addresses {
			err = mgr.markDirectorySectorUnallocated(address)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// releaseTrackRegion releases the indicated tracks from the file cycle's DAD entries,
// and adds the corresponding device tracks to the free space table.
// Only intended to be invoked on files which are currently assigned.
// If we return anything other than MFDSuccessful, the exec is stopped.
func (mgr *MFDManager) releaseFileCycleTrackRegion(
	mainItem0Address kexec2.MFDRelativeAddress,
	region *kexec2.TrackRegion,
) MFDResult {
	fas, ok := mgr.acceleratedFileAllocations[mainItem0Address]
	if !ok {
		logger.LogFatalF("MFDMgr", "convertFileRelativeTrackId Cannot find alloc for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return MFDInternalError
	}

	ldat, devTrackId := fas.ExtractRegionFromFileAllocationSet(region)
	if ldat == kexec.InvalidLDAT {
		logger.LogFatalF("MFDMgr", "convertFileRelativeTrackId Cannot extract alloc for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return MFDInternalError
	}

	return mgr.releaseTrackRegion(ldat, kexec2.NewTrackRegion(devTrackId, region.TrackCount))
}

// releaseTrackRegion adds tracks for a particular pack to the free space table.
// Only intended to be invoked on files which are *not* currently assigned, *OR* by releaseFileCycleTrackRegion().
// If we return anything other than MFDSuccessful, the exec is stopped.
func (mgr *MFDManager) releaseTrackRegion(
	ldatIndex kexec2.LDATIndex,
	region *kexec2.TrackRegion,
) MFDResult {
	packDesc, ok := mgr.fixedPackDescriptors[ldatIndex]
	if !ok {
		logger.LogFatalF("MFDMgr", "Cannot find pack descriptor for LDAT %v", ldatIndex)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return MFDInternalError
	}

	if !packDesc.freeSpaceTable.MarkTrackRegionUnallocated(region.TrackId, region.TrackCount) {
		logger.LogFatalF("MFDMgr", "Cannot un-allocate region LDAT %v TrkId %v Count %v",
			ldatIndex, region.TrackId, region.TrackCount)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return MFDInternalError
	}

	return MFDSuccessful
}

// writeBlockToDisk reads a single block from a particular disk device
// and sleeps until the IO is complete.
func (mgr *MFDManager) writeBlockToDisk(
	nodeId hardware.NodeIdentifier,
	buffer []pkg.Word36,
	blockId hardware.BlockId,
) ioPackets2.IoStatus {
	cw := channels.ControlWord{
		Buffer:    buffer,
		Offset:    0,
		Length:    uint(len(buffer)),
		Direction: channels.DirectionForward,
		Format:    channels.TransferPacked,
	}
	cp := &channels.ChannelProgram{
		NodeIdentifier: nodeId,
		IoFunction:     ioPackets.IofWrite,
		BlockId:        blockId,
		ControlWords:   []channels.ControlWord{cw},
	}

	mgr.exec.GetNodeManager().RouteIo(cp)
	for cp.IoStatus == ioPackets2.IosInProgress || cp.IoStatus == ioPackets2.IosNotStarted {
		time.Sleep(10 * time.Millisecond)
	}

	return cp.IoStatus
}

// writeFileAllocationEntryUpdatesForFileCycle writes an updated fae to the on-disk MFD
// as a series of one or more DAD tables.
// MUST be invoked when a file is free'd.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) writeFileAllocationEntryUpdatesForFileCycle(mainItem0Address kexec2.MFDRelativeAddress) error {
	fas, ok := mgr.acceleratedFileAllocations[mainItem0Address]
	if !ok {
		logger.LogFatalF("MFDMgr", "convertFileRelativeTrackId Cannot find alloc for address %012o", mainItem0Address)
		mgr.exec.Stop(kexec2.StopDirectoryErrors)
		return fmt.Errorf("fas not loaded")
	}

	if fas.IsUpdated {
		//	rewrite the entire set of DAD entries, allocate a new one if we need to do so,
		//  and release any left over when we're done.
		//  Don't forget to write hole DADs (see pg 2-63 for Device index field)
		mainItem0, err := mgr.getMFDSector(mainItem0Address)
		if err != nil {
			return err
		}

		// release all the current DAD entries
		err = mgr.releaseDADChain(mainItem0Address)
		if err != nil {
			return err
		}

		current := mainItem0
		prevAddr := mainItem0Address
		for fax := 0; fax < len(fas.FileAllocations); {
			preferred := getLDATIndexFromMFDAddress(prevAddr)
			newAddr, newEntry, err := mgr.allocateDirectorySector(preferred)
			if err != nil {
				return err
			}

			// link to next entry
			newEntry[0].SetW(0_400000_000000)
			// link to previous entry (or to main item if this is the first)
			newEntry[1].SetW(uint64(prevAddr))
			// file-relative word address of first entry
			newEntry[2].SetW(uint64(fas.FileAllocations[fax].FileRegion.TrackId) * 1792)

			// Link the new entry to the current entry, then make the new entry current.
			// This works even if the current entry is a main item.
			current[0].SetW((current[0].GetW() & 0_340000_000000) | uint64(newAddr))
			mgr.markDirectorySectorDirty(prevAddr)
			current = newEntry

			var ex int                       // index of next-to-be-used DAD table entry
			var nextTrackId hardware.TrackId // expected next file-relative track ID - only valid for ex > 0
			for fax < len(fas.FileAllocations) && ex <= 7 {
				this := fas.FileAllocations[fax]

				// Is this entry file-contiguous with the previous entry?
				if ex > 0 && this.FileRegion.TrackId != nextTrackId {
					// no - if this is the last DAD entry in the current table, the table is done.
					// otherwise, we need to write a hole entry.
					if ex < 7 {
						holeTracks := this.FileRegion.TrackId - nextTrackId
						ey := ex*3 + 4

						current[ey].SetW(0)
						current[ey+1].SetW(uint64(holeTracks) * 1792)
						// TODO if removable, set bit 16 in +02,H1
						current[ey+2].SetH2(uint64(kexec.InvalidLDAT))

						nextTrackId += holeTracks
						current[03].SetW(uint64(nextTrackId * 1792))
						ex++
					} else {
						// stop here with ex still set to 7 (so we can mark ex-1 as the last entry in the table)
						break
					}
				} else {
					ey := ex*3 + 4
					current[ey].SetW(uint64(this.DeviceTrackId) * 1792)
					current[ey+1].SetW(uint64(this.FileRegion.TrackCount * 1792))
					// TODO if removable, set bit 16 in +02,H1
					current[ey+2].SetH2(uint64(this.LDATIndex))

					nextTrackId = this.FileRegion.TrackId + hardware.TrackId(this.FileRegion.TrackCount)
					current[03].SetW(uint64(nextTrackId * 1792))
					ex++
					fax++
				}
			}

			// mark the last entry as ... well, the last entry. Word +02,H1 bit 15
			ey := (ex-1)*3 + 4
			current[ey+2].SetH1(current[ey+2].GetW() & 0_000004)
			mgr.markDirectorySectorDirty(newAddr)

			prevAddr = newAddr
		}

		fas.IsUpdated = false
	}

	return nil
}

func (mgr *MFDManager) writeLookupTableEntry(
	qualifier string,
	filename string,
	leadItem0Addr kexec2.MFDRelativeAddress) {

	key := qualifier + "*" + filename
	mgr.fileLeadItemLookupTable[key] = leadItem0Addr
}

// writeMFDCache writes all the dirty cache blocks to storage.
// If we return error, we've already stopped the exec.
// Currently, we do our own resolution of file-relative address to disk-relative.
// CALL UNDER LOCK
func (mgr *MFDManager) writeMFDCache() error {
	logger.LogTraceF("MFDMgr", "writeMFDCache (%v dirty blocks)", len(mgr.dirtyBlocks))

	for blockAddr := range mgr.dirtyBlocks {
		block, err := mgr.getMFDBlock(blockAddr)
		if err != nil {
			logger.LogFatalF("MFDMgr", "writeMFDCache cannot find MFD block for dirty block address:%012o", blockAddr)
			mgr.exec.Stop(kexec2.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		}

		mfdTrackId := (blockAddr >> 6) & 077777777
		mfdSectorId := getMFDSectorIdFromMFDAddress(blockAddr)

		ldat, devTrackId, err := mgr.convertFileRelativeTrackId(mgr.mfdFileMainItem0Address, hardware.TrackId(mfdTrackId))
		if err != nil {
			logger.LogFatalF("MFDMgr", "writeMFDCache error converting mfdaddr:%012o TrackId:%06o",
				mgr.mfdFileMainItem0Address, mfdTrackId)
			mgr.exec.Stop(kexec2.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		} else if ldat == 0_400000 {
			logger.LogFatalF("MFDMgr", "writeMFDCache error converting mfdaddr:%012o TrackId:%06o track not allocated",
				mgr.mfdFileMainItem0Address, mfdTrackId)
			mgr.exec.Stop(kexec2.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		}

		packDesc, ok := mgr.fixedPackDescriptors[ldat]
		if !ok {
			logger.LogFatalF("MFDMgr", "writeMFDCache cannot find packDesc for ldat:%04v", ldat)
			mgr.exec.Stop(kexec2.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		}

		blocksPerTrack := 1792 / packDesc.prepFactor
		sectorsPerBlock := packDesc.prepFactor / 28
		devBlockId := uint64(devTrackId) * uint64(blocksPerTrack)
		devBlockId += uint64(mfdSectorId) / uint64(sectorsPerBlock)

		ioStat := mgr.writeBlockToDisk(packDesc.nodeId, block, hardware.BlockId(devBlockId))
		if ioStat == ioPackets2.IosInternalError {
			return fmt.Errorf("internal error")
		} else if ioStat != ioPackets2.IosComplete {
			logger.LogFatalF("MFDMgr", "writeMFDCache error writing MFD block status=%v", ioStat)
			mgr.exec.Stop(kexec2.StopInternalExecIOFailed)
			return fmt.Errorf("error draining MFD cache")
		}
	}

	mgr.dirtyBlocks = make(map[kexec2.MFDRelativeAddress]bool)
	return nil
}
