// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/types"
	"khalehla/pkg"
	"log"
)

// InitializeMassStorage handles MFD initialization for what is effectively a JK13 boot.
// If we return an error, we must previously stop the exec.
func (mgr *MFDManager) InitializeMassStorage() {
	// Get the list of disks from the node manager
	disks := make([]*nodeMgr.DiskDeviceInfo, 0)
	nm := mgr.exec.GetNodeManager()
	for _, dInfo := range nm.GetDeviceInfos() {
		if dInfo.GetNodeType() == types.NodeTypeDisk {
			disks = append(disks, dInfo.(*nodeMgr.DiskDeviceInfo))
		}
	}

	// Check the labels on the disks so that we may segregate them into fixed and removable lists.
	// Any problems at this point will lead us to DN the unit.
	fixedDisks := make(map[*nodeMgr.DiskDeviceInfo]*types.DiskAttributes)
	removableDisks := make(map[*nodeMgr.DiskDeviceInfo]*types.DiskAttributes)
	for _, ddInfo := range disks {
		if ddInfo.GetNodeStatus() == types.NodeStatusUp {
			// Get the pack label from fac mgr
			attr, err := mgr.exec.GetFacilitiesManager().GetDiskAttributes(ddInfo.GetDeviceIdentifier())
			if err != nil {
				mgr.exec.SendExecReadOnlyMessage("Internal configuration error", nil)
				mgr.exec.Stop(types.StopInitializationSystemConfigurationError)
				return
			}

			if attr.PackAttrs == nil {
				msg := fmt.Sprintf("No label exists for pack on device %v", ddInfo.GetDeviceName())
				mgr.exec.SendExecReadOnlyMessage(msg, nil)
				_ = mgr.exec.GetNodeManager().SetNodeStatus(ddInfo.GetNodeIdentifier(), types.NodeStatusDown)
				continue
			}

			if !attr.PackAttrs.IsPrepped {
				msg := fmt.Sprintf("Pack is not prepped on device %v", ddInfo.GetDeviceName())
				mgr.exec.SendExecReadOnlyMessage(msg, nil)
				_ = mgr.exec.GetNodeManager().SetNodeStatus(ddInfo.GetNodeIdentifier(), types.NodeStatusDown)
				continue
			}

			// Read sector 1 of the initial directory track.
			// This is a little messy due to the potential of problematic block sizes.
			wordsPerBlock := attr.PackAttrs.Label[4].GetH2()
			dirTrackWordAddr := attr.PackAttrs.Label[03].GetW()
			dirTrackBlockId := types.BlockId(dirTrackWordAddr / wordsPerBlock)
			if wordsPerBlock == 28 {
				dirTrackBlockId++
			}

			buf := make([]pkg.Word36, wordsPerBlock)
			pkt := nodeMgr.NewDiskIoPacketRead(ddInfo.GetDeviceIdentifier(), dirTrackBlockId, buf)
			mgr.exec.GetNodeManager().RouteIo(pkt)
			ioStat := pkt.GetIoStatus()
			if ioStat != types.IosComplete {
				msg := fmt.Sprintf("IO error reading directory track on device %v", ddInfo.GetDeviceName())
				log.Printf("MFDMgr:%v", msg)
				mgr.exec.SendExecReadOnlyMessage(msg, nil)
				_ = mgr.exec.GetNodeManager().SetNodeStatus(ddInfo.GetNodeIdentifier(), types.NodeStatusDown)
				continue
			}

			var sector1 []pkg.Word36
			if wordsPerBlock == 28 {
				sector1 = buf
			} else {
				sector1 = buf[28:56]
			}

			// get the LDAT field from sector 1
			// If it is 0, it is a removable pack
			// 0400000, it is an uninitialized fixed pack
			// anything else, it is a pre-used fixed pack which we're going to initialize
			ldat := sector1[5].GetH1()
			if ldat == 0 {
				removableDisks[ddInfo] = attr
			} else {
				fixedDisks[ddInfo] = attr
				attr.PackAttrs.IsFixed = true
			}
		}
	}

	err := mgr.initializeFixed(fixedDisks)
	if err != nil {
		return
	}

	// Make sure we have at least one fixed pack after the previous shenanigans
	if len(mgr.fixedPackDescriptors) == 0 {
		mgr.exec.SendExecReadOnlyMessage("No Fixed Disks - Cannot Continue Initialization", nil)
		mgr.exec.Stop(types.StopInitializationSystemConfigurationError)
		return
	}

	err = mgr.initializeRemovable(removableDisks)
	return
}

// RecoverMassStorage handles MFD recovery for what is NOT a JK13 boot.
// If we return an error, we must previously stop the exec.
func (mgr *MFDManager) RecoverMassStorage() {
	// TODO
	mgr.exec.SendExecReadOnlyMessage("MFD Recovery is not implemented", nil)
	mgr.exec.Stop(types.StopDirectoryErrors)
	return
}

// ------------------------------------------------------------------------------------------------

// allocateDirectorySector allocates an MFD directory sector for the caller.
// If preferredLDAT is not InvalidLDAT we will try to allocate a sector from this pack first.
// Apart from this, we prefer packs with the least number of allocated sectors, to balance the allocations.
// If there is no free sector, we allocate a new track (again using the preferredLDAT), then
// allocate the first free sector from that track.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) allocateDirectorySector(
	preferredLDAT types.LDATIndex,
) (types.MFDRelativeAddress, []pkg.Word36, error) {
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
	chosenAddress := types.InvalidLink
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
	preferredLDAT types.LDATIndex,
) (types.LDATIndex, types.MFDTrackId, error) {
	chosenLDAT := types.InvalidLDAT
	var chosenDesc *fixedPackDescriptor
	chosenAvailableTracks := types.TrackCount(0)

	if preferredLDAT != types.InvalidLDAT {
		packDesc, ok := mgr.fixedPackDescriptors[preferredLDAT]
		if ok {
			availMFDTracks := 07777 - packDesc.mfdTrackCount
			availTracks := packDesc.freeSpaceTable.getFreeTrackCount()
			if availMFDTracks > 0 && availTracks > 0 {
				chosenLDAT = preferredLDAT
				chosenDesc = packDesc
				chosenAvailableTracks = availTracks
			}
		}
	}

	if chosenLDAT == types.InvalidLDAT {
		for ldat, packDesc := range mgr.fixedPackDescriptors {
			availMFDTracks := 07777 - packDesc.mfdTrackCount
			availTracks := packDesc.freeSpaceTable.getFreeTrackCount()
			if availMFDTracks > 0 && availTracks > chosenAvailableTracks {
				chosenLDAT = ldat
				chosenDesc = packDesc
				chosenAvailableTracks = availTracks
			}
		}
	}

	if chosenLDAT == types.InvalidLDAT {
		log.Printf("MFDMgr:No space available for directory track allocation")
		mgr.exec.Stop(types.StopExecRequestForMassStorageFailed)
		return 0, 0, fmt.Errorf("no disk")
	}

	// First, find the MFD relative address of the first unused MFD track
	trackId := types.MFDTrackId(0)
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
	trackId, _ = chosenDesc.freeSpaceTable.allocateTrack()
	chosenDesc.mfdTrackCount++
	return chosenLDAT, trackId, nil
}

// bootstrapMFD creates the various MFD structures as part of MFD initialization.
// One consequence is the cataloging of SYS$*MFD$$.
// Since this is used during initialization we do not call it under lock.
func (mgr *MFDManager) bootstrapMFD() error {
	log.Printf("MFDMgr:bootstrapMFD start")

	cfg := mgr.exec.GetConfiguration()

	// Find the highest and lowest LDAT indices
	var lowestLDAT = types.InvalidLDAT
	var highestLDAT = types.LDATIndex(0)
	for ldat := range mgr.fixedPackDescriptors {
		if lowestLDAT == types.InvalidLDAT && ldat < lowestLDAT {
			lowestLDAT = ldat
		}
		if ldat > highestLDAT {
			highestLDAT = ldat
		}
	}

	// We've already initialized the DAS sector - put the expected free sectors (2 through 63)
	// into the free sector list.
	for sx := 2; sx < 64; sx++ {
		sectorAddr := composeMFDAddress(lowestLDAT, 0, types.MFDSectorId(sx))
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
	populateMainItem0(mainItem0, cfg.SystemQualifier, mfdFileName, cfg.SystemProjectId,
		cfg.SystemReadKey, cfg.SystemWriteKey, cfg.SystemAccountId, leadAddr0, mainAddr1,
		false, false, false, false, true,
		false, false, cfg.AssignMnemonic, true, true,
		true, false, false, 1, 0, 262153, []string{})
	populateFixedMainItem1(mainItem1, cfg.SystemQualifier, mfdFileName, mainAddr0, 1, []string{})

	// Before we can play DAD table games, we have to get the MFD$$ in-core structures in place,
	// including *particularly* the file allocation table.
	// We need to create one allocation region for each pack's initial directory track.
	highestMFDTrackId := types.TrackId(0)
	fae := newFileAllocationEntry(mainAddr0, 0_400000_000000)
	mgr.fileAllocations.content[mainAddr0] = fae

	for ldat, desc := range mgr.fixedPackDescriptors {
		mfdTrackId := types.TrackId(ldat << 12)
		if mfdTrackId > highestMFDTrackId {
			highestMFDTrackId = mfdTrackId
		}
		packTrackId := desc.packAttributes.Label[03] / 1792
		err := mgr.allocateSpecificTrack(mainAddr0, mfdTrackId, 1, ldat, types.TrackId(packTrackId))
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
	ldatIndex types.LDATIndex,
	trackId types.MFDTrackId,
	sectorId types.MFDSectorId,
) types.MFDRelativeAddress {

	return types.MFDRelativeAddress(uint64(ldatIndex&07777)<<18 | uint64(trackId&07777)<<6 | uint64(sectorId&077))
}

// findDASEntryForSector chases the appropriate DAS chain to find the DAS which describes the given sector address,
// and then the entry within that DAS which describes the sector address.
// Returns:
//
//	the address of the containing DAS sector
//	the index of the DAS entry
//	a slice to the 3-word DAS entry itself.
func (mgr *MFDManager) findDASEntryForSector(
	sectorAddr types.MFDRelativeAddress,
) (types.MFDRelativeAddress, int, []pkg.Word36, error) {

	// what are we looking for?
	ldat := getLDATIndexFromMFDAddress(sectorAddr)
	trackId := getMFDTrackIdFromMFDAddress(sectorAddr)

	dasAddr := composeMFDAddress(ldat, 0, 0)
	for dasAddr != types.InvalidLink {
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
			entryAddr := types.MFDRelativeAddress(das[ey].GetW())
			if entryAddr != types.InvalidLink {
				entryTrackId := getMFDTrackIdFromMFDAddress(entryAddr)
				if entryTrackId == trackId {
					// found it.
					return dasAddr, ex, das[ey : ey+3], nil
				}
			}
		}

		// So it is not this DAS - move on to the next
		dasAddr = types.MFDRelativeAddress(das[033].GetW())
	}

	// We did not find the DAS entry - complain and crash
	log.Printf("MFDMgr:Cannot find DAS for sector %012o", sectorAddr)
	mgr.exec.Stop(types.StopDirectoryErrors)
	return 0, 0, nil, fmt.Errorf("cannot find DAS")
}

func getLDATIndexFromMFDAddress(address types.MFDRelativeAddress) types.LDATIndex {
	return types.LDATIndex(address>>18) & 07777
}

func getMFDTrackIdFromMFDAddress(address types.MFDRelativeAddress) types.MFDTrackId {
	return types.MFDTrackId(address>>6) & 07777
}

func getMFDSectorIdFromMFDAddress(address types.MFDRelativeAddress) types.MFDSectorId {
	return types.MFDSectorId(address & 077)
}

// getMFDAddressForBlock takes a given MFD-relative sector address and normalizes it to
// the first sector in the block containing the given sector.
// CALL UNDER LOCK
func (mgr *MFDManager) getMFDAddressForBlock(address types.MFDRelativeAddress) types.MFDRelativeAddress {
	ldat := getLDATIndexFromMFDAddress(address)
	mask := mgr.fixedPackDescriptors[ldat].packMask
	return types.MFDRelativeAddress(uint64(address) & ^mask)
}

// getMFDBlock returns a slice corresponding to all the sectors in the physical block
// containing the sector represented by the given address. Used for reading/writing MFD blocks.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) getMFDBlock(address types.MFDRelativeAddress) ([]pkg.Word36, error) {
	ldatAndTrack := address & 0_007777_777700
	data, ok := mgr.cachedTracks[ldatAndTrack]
	if !ok {
		log.Printf("MFDMgr:getMFDBlock address:%v is not in cache", address)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return nil, fmt.Errorf("internal error")
	}

	ldat := getLDATIndexFromMFDAddress(address)
	sector := getMFDSectorIdFromMFDAddress(address)
	mask := mgr.fixedPackDescriptors[ldat].packMask
	baseSectorId := uint64(sector) & ^mask

	start := 28 * baseSectorId
	end := start + uint64(mgr.fixedPackDescriptors[ldat].wordsPerBlock)
	return data[start:end], nil
}

// getMFDSector returns a slice corresponding to the portion of the MFD block which represents the indicated sector.
// If we return an error, we've already stopped the exec
// CALL UNDER LOCK
func (mgr *MFDManager) getMFDSector(address types.MFDRelativeAddress) ([]pkg.Word36, error) {
	ldatAndTrack := address & 0_007777_777700
	data, ok := mgr.cachedTracks[ldatAndTrack]
	if !ok {
		log.Printf("MFDMgr:getMFDSector address:%v is not in cache", address)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return nil, fmt.Errorf("internal error")
	}

	sectorId := getMFDSectorIdFromMFDAddress(address)
	start := 28 * sectorId
	end := start + 28
	return data[start:end], nil
}

// initializeFixed initializes the fixed pool for a jk13 boot
func (mgr *MFDManager) initializeFixed(disks map[*nodeMgr.DiskDeviceInfo]*types.DiskAttributes) error {
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
		mgr.exec.Stop(types.StopConsoleResponseRequiresReboot)
		return fmt.Errorf("boot canceled")
	}

	// make sure there are no pack name conflicts
	conflicts := false
	packNames := make([]string, 0)
	for _, diskAttr := range disks {
		for _, existing := range packNames {
			if diskAttr.PackAttrs.PackName == existing {
				msg := fmt.Sprintf("Fixed pack name conflict - %v", existing)
				mgr.exec.SendExecReadOnlyMessage(msg, nil)
				conflicts = true
			} else {
				packNames = append(packNames, diskAttr.PackAttrs.PackName)
			}
		}
	}

	if conflicts {
		mgr.exec.SendExecReadOnlyMessage("Resolve pack name conflicts and reboot", nil)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return fmt.Errorf("packid conflict")
	}

	// iterate over the fixed packs - we start at 1, which may not be conventional, but it works
	nextLdatIndex := types.LDATIndex(1)
	totalTracks := uint64(0)
	for diskInfo, diskAttr := range disks {
		// Assign an LDAT to the pack, update the pack label, then rewrite the label
		ldatIndex := nextLdatIndex
		nextLdatIndex++

		// Assign the unit
		_ = mgr.exec.GetFacilitiesManager().AssignDiskDeviceToExec(diskInfo.GetDeviceIdentifier())

		// Set up fixed pack descriptor
		fpDesc := newFixedPackDescriptor(diskInfo.GetDeviceIdentifier(),
			diskAttr.PackAttrs,
			diskInfo.GetNodeStatus() == types.NodeStatusUp)

		// Mark VOL1 and first directory track as allocated
		fpDesc.mfdTrackCount = 1
		fpDesc.mfdSectorsUsed = 2
		_ = fpDesc.freeSpaceTable.allocateSpecificTrackRegion(ldatIndex, 0, 2)

		// Set up first directory track for the pack within our cache to be eventually rewritten
		mgr.fixedPackDescriptors[ldatIndex] = fpDesc

		availableTracks := diskAttr.PackAttrs.Label[016].GetW() - 2
		recordLength := diskAttr.PackAttrs.Label[04].GetH2()
		blocksPerTrack := diskAttr.PackAttrs.Label[04].GetH1()
		wordsPerBlock := uint64(fpDesc.wordsPerBlock)

		// We need to read the first directory track into cache
		// so we can update sector 0 and sector 1 appropriately.
		mfdTrackId := types.MFDTrackId(ldatIndex << 12)
		mfdAddr := composeMFDAddress(ldatIndex, mfdTrackId, 0)
		data := make([]pkg.Word36, 1792)
		mgr.cachedTracks[mfdAddr] = data

		// read the directory track into cache
		blockId := types.BlockId(blocksPerTrack)
		wx := uint64(0)
		for bx := 0; bx < int(blocksPerTrack); bx++ {
			sub := data[wx : wx+wordsPerBlock]
			ioPkt := nodeMgr.NewDiskIoPacketRead(fpDesc.deviceId, blockId, sub)
			mgr.exec.GetNodeManager().RouteIo(ioPkt)
			ioStat := ioPkt.GetIoStatus()
			if ioStat != types.IosComplete {
				log.Printf("MFDMgr:initializeFixed cannot read directory track dev:%v blockId:%v",
					fpDesc.deviceId, blockId)
				mgr.exec.Stop(types.StopInternalExecIOFailed)
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
			sector0[dx].SetW(uint64(types.InvalidLink))
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
		sector1[2].SetW(availableTracks)
		sector1[3].SetW(availableTracks)
		sector1[4].FromStringToFieldata(diskAttr.PackAttrs.PackName)
		sector1[5].SetH1(uint64(ldatIndex))
		sector1[010].SetT1(blocksPerTrack)
		sector1[010].SetS3(1) // Sector 1 version
		sector1[010].SetT3(recordLength)
		mgr.markDirectorySectorDirty(sector1Addr)

		totalTracks += availableTracks
	}

	err = mgr.bootstrapMFD()
	if err != nil {
		return err
	}

	msg = fmt.Sprintf("MS Initialized - %v Tracks Available", totalTracks)
	mgr.exec.SendExecReadOnlyMessage(msg, nil)
	return nil
}

// initializeRemovable registers the removable packs (if any) as part of a JK13 boot.
func (mgr *MFDManager) initializeRemovable(disks map[*nodeMgr.DiskDeviceInfo]*types.DiskAttributes) error {
	return nil
}

// markDirectorySectorAllocated finds the appropriate DAS entry for the given sector address
// and marks the sector as allocated, as well as marking the DAS entry as updated.
func (mgr *MFDManager) markDirectorySectorAllocated(sectorAddr types.MFDRelativeAddress) error {
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
func (mgr *MFDManager) markDirectorySectorDirty(address types.MFDRelativeAddress) {
	blockAddr := mgr.getMFDAddressForBlock(address)
	mgr.dirtyBlocks[blockAddr] = true
}

// markDirectorySectorUnallocated finds the appropriate DAS entry for the given sector address
// and marks the sector as unallocated, as well as marking the DAS entry as updated.
func (mgr *MFDManager) markDirectorySectorUnallocated(sectorAddr types.MFDRelativeAddress) error {
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

func (mgr *MFDManager) writeLookupTableEntry(
	qualifier string,
	filename string,
	leadItem0Addr types.MFDRelativeAddress) {

	_, ok := mgr.fileLeadItemLookupTable[qualifier]
	if !ok {
		mgr.fileLeadItemLookupTable[qualifier] = make(map[string]types.MFDRelativeAddress)
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
			mgr.exec.Stop(types.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		}

		mfdTrackId := (blockAddr >> 6) & 077777777
		mfdSectorId := getMFDSectorIdFromMFDAddress(blockAddr)

		ldat, devTrackId, err := mgr.convertFileRelativeTrackId(mgr.mfdFileMainItem0Address, types.TrackId(mfdTrackId))
		if err != nil {
			log.Printf("MFDMgr:writeMFDCache error converting mfdaddr:%012o trackId:%06v", mgr.mfdFileMainItem0Address, mfdTrackId)
			mgr.exec.Stop(types.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		} else if ldat == 0_400000 {
			log.Printf("MFDMgr:writeMFDCache error converting mfdaddr:%012o trackId:%06v track not allocated",
				mgr.mfdFileMainItem0Address, mfdTrackId)
			mgr.exec.Stop(types.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		}

		packDesc, ok := mgr.fixedPackDescriptors[ldat]
		if !ok {
			log.Printf("MFDMgr:writeMFDCache cannot find packDesc for ldat:%04v", ldat)
			mgr.exec.Stop(types.StopDirectoryErrors)
			return fmt.Errorf("error draining MFD cache")
		}

		blocksPerTrack := 1792 / packDesc.wordsPerBlock
		sectorsPerBlock := packDesc.wordsPerBlock / 28
		devBlockId := uint64(devTrackId) * uint64(blocksPerTrack)
		devBlockId += uint64(mfdSectorId) / uint64(sectorsPerBlock)
		ioPkt := nodeMgr.NewDiskIoPacketWrite(packDesc.deviceId, types.BlockId(devBlockId), block)
		mgr.exec.GetNodeManager().RouteIo(ioPkt)
		ioStat := ioPkt.GetIoStatus()
		if ioStat != types.IosComplete {
			log.Printf("MFDMgr:writeMFDCache error writing MFD block status=%v", ioStat)
			mgr.exec.Stop(types.StopInternalExecIOFailed)
			return fmt.Errorf("error draining MFD cache")
		}
	}

	mgr.dirtyBlocks = make(map[types.MFDRelativeAddress]bool)
	return nil
}
