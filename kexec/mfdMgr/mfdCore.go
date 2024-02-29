// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// we are allowed to know about facMgr, nodeMgr, consMgr

import (
	"fmt"
	"khalehla/kexec"
	//	"khalehla/kexec/facilitiesMgr"
	"khalehla/kexec/nodeMgr"
	"khalehla/pkg"
	"log"
)

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

	re := NewFileAllocation(fileTrackId, trackCount, ldatIndex, deviceTrackId)
	fae.MergeIntoFileAllocationEntry(re)

	if fileTrackId > fae.HighestTrackAllocated {
		fae.HighestTrackAllocated = fileTrackId
	}
	fae.IsUpdated = true

	return nil
}

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

	mfdFileName := "MFD$$"
	populateNewLeadItem0(leadItem0, cfg.SystemQualifier, mfdFileName, cfg.SystemProjectId,
		cfg.SystemReadKey, cfg.SystemWriteKey, 0, 1, true, uint64(mainAddr0))
	populateMassStorageMainItem0(mainItem0, cfg.SystemQualifier, mfdFileName, cfg.SystemProjectId,
		cfg.SystemReadKey, cfg.SystemWriteKey, cfg.SystemAccountId, leadAddr0, mainAddr1,
		false, false, false, false, false,
		cfg.AssignMnemonic, true, true, true, false, false,
		1, 0, 262153, []string{})
	populateFixedMainItem1(mainItem1, cfg.SystemQualifier, mfdFileName, mainAddr0, 1, []string{})

	// Before we can play DAD table games, we have to get the MFD$$ in-core structures in place,
	// including *particularly* the file allocation table.
	// We need to create one allocation region for each pack's initial directory track.
	highestMFDTrackId := kexec.TrackId(0)
	fae := NewFileAllocationEntry(mainAddr0, 0_400000_000000)
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
	mgr.writeLookupTableEntry(cfg.SystemQualifier, "MFD$$", leadAddr0)

	// Set file assigned in facmgr, RCE, or wherever it makes sense
	// TODO

	err = mgr.writeMFDCache()
	if err != nil {
		return err
	}

	log.Printf("MFDMgr:bootstrapMFD done")
	return nil
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

// findDASEntryForSector chases the appropriate DAS chain to find the DAS which describes the given sector address,
// and then the entry within that DAS which describes the sector address.
// Returns:
//
//	the address of the containing DAS sector
//	the index of the DAS entry
//	a slice to the 3-word DAS entry itself.
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
	return nil
}

// loadFileAllocationEntry initializes the fae for a particular file instance.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) loadFileAllocationEntry(
	mainItem0Address kexec.MFDRelativeAddress,
) (*FileAllocationEntry, error) {
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
	fae := &FileAllocationEntry{
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
				fae.MergeIntoFileAllocationEntry(re)
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

	if fae.IsUpdated {
		// TODO process
		//	rewrite the entire set of DAD entries, allocate a new one if we need to do so,
		//  and release any left over when we're done.
		//  Don't forget to write hole DADs (see pg 2-63 for Device index field)

		fae.IsUpdated = false
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
