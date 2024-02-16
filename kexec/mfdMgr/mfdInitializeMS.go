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
	"time"
)

// bootstrapMFDFile creates the MFD entries for the SYS$*MFD$$ file as part of MFD initialization
// lead item always goes to LDAT 1, blkId 0, sectorId 2
// main items always go to LDAT 1, blkId 0, sectorId 3 and 4
func (mgr *MFDManager) bootstrapMFDFile() error {
	cfg := mgr.exec.GetConfiguration()

	dasAddr0 := composeMFDAddress(1, 0, 0)
	leadAddr0 := composeMFDAddress(1, 0, 2)
	mainAddr0 := composeMFDAddress(1, 0, 3)
	mainAddr1 := composeMFDAddress(1, 0, 4)

	dasSector0, _ := mgr.getMFDSector(dasAddr0)
	leadItem0, _ := mgr.getMFDSector(leadAddr0)
	mainItem0, _ := mgr.getMFDSector(mainAddr0)
	mainItem1, _ := mgr.getMFDSector(mainAddr1)

	// allocate sectors 2, 3, 4, and 5 in the DAS
	dasSector0[01].Or(0_170000_000000)

	swTimeNow := types.GetSWTimeFromSystemTime(time.Now())

	// populate lead item
	leadItem0[0].SetW(0_500000_000000)
	pkg.FromStringToFieldataWithOffset(cfg.SystemQualifier, leadItem0, 1, 2)
	pkg.FromStringToFieldataWithOffset("MFD$$", leadItem0, 3, 2)
	pkg.FromStringToFieldataWithOffset(cfg.SystemProjectId, leadItem0, 5, 2)
	pkg.FromStringToFieldataWithOffset(cfg.SystemReadKey, leadItem0, 7, 1)
	pkg.FromStringToFieldataWithOffset(cfg.SystemWriteKey, leadItem0, 010, 1)

	leadItem0[011].SetS1(0)   // file type == mass storage
	leadItem0[011].SetS2(1)   // number of f-cycles which exist
	leadItem0[011].SetS3(1)   // max number of f-cycles for this file
	leadItem0[011].SetS4(1)   // current number of f-cycles for this file
	leadItem0[011].SetT3(1)   // highest absolute f-cycle
	leadItem0[012].SetS1(040) // status bits - Guarded file
	leadItem0[012].SetS4(0)   // no security words
	leadItem0[013].SetW(uint64(mainAddr0))

	// populate main items
	mainItem0[0].SetW(0_200000_000000)
	pkg.FromStringToFieldataWithOffset(cfg.SystemQualifier, mainItem0, 1, 2)
	pkg.FromStringToFieldataWithOffset("MFD$$", mainItem0, 3, 2)
	pkg.FromStringToFieldataWithOffset(cfg.SystemProjectId, mainItem0, 5, 2)
	pkg.FromStringToFieldataWithOffset(cfg.SystemAccountId, mainItem0, 7, 2)
	mainItem0[012].SetW(swTimeNow)
	mainItem0[013].SetW(uint64(leadAddr0))
	mainItem0[014].SetT1(0)                // descriptor flags
	mainItem0[014].SetS3(02)               // file flags - written-to
	mainItem0[015].SetW(uint64(mainAddr1)) // link to main item sector 1
	mainItem0[015].SetS1(0)                // PCHAR flags
	pkg.FromStringToFieldataWithOffset(cfg.AssignMnemonic, mainItem0, 016, 1)
	mainItem0[021].SetS2(070) // guarded, inhibit unload, private
	mainItem0[021].SetT2(1)   // absolute f-cycle number
	mainItem0[022].SetW(swTimeNow)
	mainItem0[023].SetW(swTimeNow)
	mainItem0[024].SetH1(1)                      // initial granules
	mainItem0[025].SetH1(0777777)                // max granules
	mainItem0[026].SetH1(uint64(leadAddr0) >> 6) // highest granule assigned
	mainItem0[027].SetH1(uint64(leadAddr0) >> 6) // highest track written

	mainItem1[1].SetW(0_400000_000000) // no sector 2
	pkg.FromStringToFieldataWithOffset(cfg.SystemQualifier, mainItem1, 1, 2)
	pkg.FromStringToFieldataWithOffset("MFD$$", mainItem1, 3, 2)
	pkg.FromStringToFieldataWithOffset("*NO.1*", mainItem1, 5, 1)
	mainItem1[06].SetW(uint64(mainAddr0))
	mainItem1[07].SetT3(1) // abs f-cycle
	// TODO what are disk pack entries for? do we need them?

	// populate DAD items for MFD
	// TODO
	//  actually, can we populate this from the existing file allocation table?
	//  if so, maybe we can tolerate more fixed packs...
	//  note that our limit for 1 DAD entry is 4 packs, not 8, since we'd have to do hole DAD entries...

	mgr.mfdFileMainItem0Address = mainAddr0

	mgr.markSectorDirty(dasAddr0)
	mgr.markSectorDirty(leadAddr0)
	mgr.markSectorDirty(mainAddr0)
	mgr.markSectorDirty(mainAddr1)

	// Update lookup table
	mgr.writeLookupTableEntry(cfg.SystemQualifier, "MFD$$", leadAddr0)

	// Create FAT
	_, err := mgr.loadFileAllocationEntry(mainAddr0)
	if err != nil {
		return fmt.Errorf("cannot load FAT for MFD")
	}

	// Set file assigned here, in MFD item, in facmgr, wherever it makes sense
	//TODO

	return nil
}

// initializeMassStorage handles MFD initialization for what is effectively a JK13 boot.
// If we return an error, we must previously stop the exec.
func (mgr *MFDManager) initializeMassStorage() error {
	// drain the device ready notification queue, and use that to build up our initial list of disk packs.
	// We should only have notifications for disks, and the ready flag should always be true.
	// However, we'll filter out any nonsense anyway.
	nm := mgr.exec.GetNodeManager()
	queue := mgr.deviceReadyNotificationQueue
	mgr.deviceReadyNotificationQueue = make(map[types.DeviceIdentifier]bool)
	disks := make([]*nodeMgr.DiskDeviceInfo, 0)
	for devId, ready := range queue {
		if ready {
			devInfo, err := nm.GetNodeInfoByIdentifier(types.NodeIdentifier(devId))
			if err == nil && devInfo.GetNodeType() == types.NodeTypeDisk {
				disks = append(disks, devInfo.(*nodeMgr.DiskDeviceInfo))
			}
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
				mgr.exec.SendExecReadOnlyMessage("Internal configuration error")
				mgr.exec.Stop(types.StopInitializationSystemConfigurationError)
				return fmt.Errorf("boot canceled")
			}

			if attr.PackAttrs == nil {
				msg := fmt.Sprintf("No label exists for pack on device %v", ddInfo.GetDeviceName())
				mgr.exec.SendExecReadOnlyMessage(msg)
				_ = mgr.exec.GetNodeManager().SetNodeStatus(ddInfo.GetNodeIdentifier(), types.NodeStatusDown)
				continue
			}

			if !attr.PackAttrs.IsPrepped {
				msg := fmt.Sprintf("Pack is not prepped on device %v", ddInfo.GetDeviceName())
				mgr.exec.SendExecReadOnlyMessage(msg)
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
				mgr.exec.SendExecReadOnlyMessage(msg)
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
			ldat := sector1[5].GetS1()
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
		return err
	}

	// Make sure we have at least one fixed pack after the previous shenanigans
	if len(mgr.fixedPackDescriptors) == 0 {
		mgr.exec.SendExecReadOnlyMessage("No Fixed Disks - Cannot Continue Initialization")
		mgr.exec.Stop(types.StopInitializationSystemConfigurationError)
		return fmt.Errorf("boot canceled")
	}

	err = mgr.initializeRemovable(removableDisks)
	return err
}

func (mgr *MFDManager) initializeFixed(disks map[*nodeMgr.DiskDeviceInfo]*types.DiskAttributes) error {
	msg := fmt.Sprintf("Fixed Disk Pool = %v Devices", len(disks))
	mgr.exec.SendExecReadOnlyMessage(msg)

	if len(disks) == 0 {
		return nil
	} else if len(disks) > 8 {
		mgr.exec.SendExecReadOnlyMessage("Max of 8 fixed packs allowed for initial boot")
		mgr.exec.Stop(types.StopInitializationSystemConfigurationError)
		return fmt.Errorf("MFDMgr too many fixed packs")
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
				mgr.exec.SendExecReadOnlyMessage(msg)
				conflicts = true
			} else {
				packNames = append(packNames, diskAttr.PackAttrs.PackName)
			}
		}
	}

	if conflicts {
		mgr.exec.SendExecReadOnlyMessage("Resolve pack name conflicts and reboot")
		mgr.exec.Stop(types.StopDirectoryErrors)
		return fmt.Errorf("packid conflict")
	}

	// iterate over the fixed packs
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
		_ = fpDesc.fixedFeeSpace.allocateSpecificTrackRegion(ldatIndex, 0, 2)

		// Set up first directory track for the pack within our cache to be eventually rewritten
		mgr.fixedPackDescriptors[ldatIndex] = fpDesc
		mgr.establishNewMFDTrack(ldatIndex, 0)

		availableTracks := diskAttr.PackAttrs.Label[016].GetW() - 2
		recordLength := diskAttr.PackAttrs.Label[04].GetH2()
		blocksPerTrack := diskAttr.PackAttrs.Label[04].GetH1()

		// sector 0
		sector0Addr := composeMFDAddress(ldatIndex, 0, 0)
		sector0, _ := mgr.getMFDSector(sector0Addr)
		sector0[1].SetW(0_600000_000000) // first 2 sectors are allocated
		for dx := 3; dx < 27; dx += 3 {
			sector0[dx].SetW(0_400000_000000)
		}
		sector0[27].SetW(0_400000_000000)
		mgr.markSectorDirty(sector0Addr)

		// sector 1
		sector1Addr := composeMFDAddress(ldatIndex, 0, 1)
		sector1, _ := mgr.getMFDSector(sector1Addr)
		// leave +0 and +1 alone (We aren't doing HMBT/SMBT so we don't need the addresses)
		sector1[2].SetW(availableTracks)
		sector1[3].SetW(availableTracks)
		sector1[4].FromStringToFieldata(diskAttr.PackAttrs.PackName)
		sector1[5].SetH1(uint64(ldatIndex))
		sector1[010].SetT1(blocksPerTrack)
		sector1[010].SetS3(1) // Sector 1 version
		sector1[010].SetT3(recordLength)
		mgr.markSectorDirty(sector1Addr)

		totalTracks += availableTracks
	}

	err = mgr.bootstrapMFDFile()
	if err != nil {
		return err
	}

	msg = fmt.Sprintf("MS Initialized - %v Tracks Available", totalTracks)
	mgr.exec.SendExecReadOnlyMessage(msg)
	return nil
}

func (mgr *MFDManager) initializeRemovable(disks map[*nodeMgr.DiskDeviceInfo]*types.DiskAttributes) error {
	//for diskInfo, diskAttr := range disks {
	//	// TODO
	//}
	return nil
}
