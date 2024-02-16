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
	"os"
	"time"
)

// bootstrapMFDFile creates the MFD entries for the SYS$*MFD$$ file as part of MFD initialization
// lead item always goes to LDAT 1, blkId 0, sectorId 2
// main items always go to LDAT 1, blkId 0, sectorId 3 and 4
// we do NOT store a DAD table for SYS$*MFD$$ - we don't need it, and it would be full of holes anyway.
func (mgr *MFDManager) bootstrapMFDFile() error {
	cfg := mgr.exec.GetConfiguration()

	dasAddr0 := composeMFDAddress(1, 0, 0)
	leadAddr0 := composeMFDAddress(1, 0, 2)
	mainAddr0 := composeMFDAddress(1, 0, 3)
	mainAddr1 := composeMFDAddress(1, 0, 4)

	dasSector0, _ := mgr.getMFDSector(dasAddr0, true)
	leadItem0, _ := mgr.getMFDSector(leadAddr0, true)
	mainItem0, _ := mgr.getMFDSector(mainAddr0, true)
	mainItem1, _ := mgr.getMFDSector(mainAddr1, true)

	// allocate sectors 2, 3, and 4 in the DAS
	dasSector0[01].Or(0_160000_000000)

	swTimeNow := types.GetSWTimeFromSystemTime(time.Now())

	// populate lead item
	leadItem0[0].SetW(0_500000_000000)
	pkg.FromStringToFieldataWithOffset(cfg.SystemQualifier, leadItem0, 1, 2)
	pkg.FromStringToFieldataWithOffset("MFD$$", leadItem0, 3, 2)
	pkg.FromStringToFieldataWithOffset(cfg.SystemProjectId, leadItem0, 5, 2)
	pkg.FromStringToFieldataWithOffset(cfg.SystemReadKey, leadItem0, 7, 1)
	pkg.FromStringToFieldataWithOffset(cfg.SystemWriteKey, leadItem0, 010, 1)

	leadItem0[011].SetS1(0)   // file type == mass stroage
	leadItem0[011].SetS2(1)   // number of f-cycles which exist
	leadItem0[011].SetS3(1)   // max number of f-cycles for this file
	leadItem0[011].SetS4(1)   // current number of f-cycles for this file
	leadItem0[011].SetT3(1)   // highest absolute f-cycle
	leadItem0[012].SetS1(040) // status bits - Guarded file
	leadItem0[012].SetS4(0)   // no security words
	leadItem0[013].SetW(uint64(mainAddr0))

	// populate main items
	mainItem0[0].SetW(0_400000_000000) // no DAD table
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

	// Update lookup table
	lookupKey := cfg.SystemQualifier + ":MFD$$"
	mgr.fixedLookupTable[lookupKey] = leadAddr0

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

	// Go do the work
	mgr.fixedLDAT = make(map[types.LDATIndex]fixedPackDescriptor)
	mgr.fixedLookupTable = make(map[string]types.MFDRelativeAddress)

	err := mgr.initializeFixed(fixedDisks)
	if err != nil {
		return err
	}

	// Make sure we have at least one fixed pack after the previous shenanigans
	if len(mgr.fixedLDAT) == 0 {
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
	// TODO

	// iterate over the fixed packs
	nextLdatIndex := types.LDATIndex(1)
	totalTracks := uint64(0)
	for diskInfo, diskAttr := range disks {
		// Assign an LDAT to the pack, update the pack label, then rewrite the label
		ldatIndex := nextLdatIndex
		nextLdatIndex++

		// Assign the unit
		_ = mgr.exec.GetFacilitiesManager().AssignDiskDeviceToExec(diskInfo.GetDeviceIdentifier())

		// Rewrite first directory track to the pack
		dirTrack := make([]pkg.Word36, 1792)
		availableTracks := diskAttr.PackAttrs.Label[016].GetW() - 2
		recordLength := diskAttr.PackAttrs.Label[04].GetH2()
		blocksPerTrack := diskAttr.PackAttrs.Label[04].GetH1()

		fpDesc := fixedPackDescriptor{
			deviceId:         diskInfo.GetDeviceIdentifier(),
			packAttributes:   diskAttr.PackAttrs,
			wordsPerBlock:    types.PrepFactor(recordLength),
			canAllocate:      diskInfo.GetNodeStatus() == types.NodeStatusUp,
			packMask:         (recordLength / 28) - 1,
			trackDescriptors: make(map[types.MFDTrackId]fixedTrackDescriptor),
			freeSpace:        make(map[types.TrackId]types.TrackCount),
		}
		trackCount := diskAttr.PackAttrs.Label[016].GetW() - 2
		fpDesc.freeSpace[2] = types.TrackCount(trackCount - 2)
		mgr.fixedLDAT[ldatIndex] = fpDesc

		devTrackAddr := types.DeviceRelativeWordAddress(1792)
		_ = mgr.establishNewMFDTrack(ldatIndex, 0, devTrackAddr)

		// sector 0
		das := dirTrack[0:28]
		das[1].SetW(0_600000_000000) // first 2 sectors are allocated
		for dx := 3; dx < 27; dx += 3 {
			das[dx].SetW(0_400000_000000)
		}
		das[27].SetW(0_400000_000000)

		// sector 1
		s1 := dirTrack[28:56]
		// leave +0 and +1 alone (We aren't doing HMBT/SMBT so we don't need the addresses)
		s1[2].SetW(availableTracks)
		s1[3].SetW(availableTracks)
		s1[4].FromStringToFieldata(diskAttr.PackAttrs.PackName)
		s1[5].SetH1(uint64(ldatIndex))
		s1[010].SetT1(blocksPerTrack)
		s1[010].SetS3(1) // Sector 1 version
		s1[010].SetT3(recordLength)

		// Write the entire directory track to storage
		dirTrackWordAddr := diskAttr.PackAttrs.Label[03].GetW()
		blockId := types.BlockId(dirTrackWordAddr / blocksPerTrack)
		blocksLeft := blocksPerTrack
		subSetStart := 0

		foundError := false
		for blocksLeft > 0 && !foundError {
			subSet := dirTrack[subSetStart : subSetStart+int(recordLength)]
			pkt := nodeMgr.NewDiskIoPacketWrite(diskInfo.GetDeviceIdentifier(), blockId, subSet)
			mgr.exec.GetNodeManager().RouteIo(pkt)
			ioStat := pkt.GetIoStatus()
			if ioStat != types.IosComplete {
				msg := fmt.Sprintf("IO error reading directory track on device %v", diskInfo.GetDeviceName())
				log.Printf("MFDMgr:%v", msg)
				mgr.exec.SendExecReadOnlyMessage(msg)
				_ = mgr.exec.GetNodeManager().SetNodeStatus(diskInfo.GetNodeIdentifier(), types.NodeStatusDown)
				foundError = true
			}

			blocksLeft--
			subSetStart += int(recordLength)
		}

		if foundError {
			os.Exit(1)
		}

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
