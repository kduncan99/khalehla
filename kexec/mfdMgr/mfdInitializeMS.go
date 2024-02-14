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
)

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

	msg := fmt.Sprintf("Fixed Disk Pool = %v Devices", len(disks))
	mgr.exec.SendExecReadOnlyMessage(msg)

	// Go do the work
	mgr.directoryTracks = make(map[uint64][]pkg.Word36)
	mgr.fixedLDAT = make(map[uint]types.DeviceIdentifier)

	err := mgr.initializeFixed(fixedDisks)
	if err != nil {
		return err
	}

	// Make sure we have at least one fixed pack after the previous shenanigans
	if len(mgr.directoryTracks) == 0 {
		mgr.exec.SendExecReadOnlyMessage("No Fixed Disks - Cannot Continue Initialization")
		mgr.exec.Stop(types.StopInitializationSystemConfigurationError)
		return fmt.Errorf("boot canceled")
	}

	err = mgr.initializeRemovable(removableDisks)
	return err
}

func (mgr *MFDManager) initializeFixed(disks map[*nodeMgr.DiskDeviceInfo]*types.DiskAttributes) error {
	// make sure there are no pack name conflicts
	// TODO
	os.Exit(1) // TODO remove

	nextLdatIndex := uint(1)
	totalTracks := uint64(0)
	for diskInfo, diskAttr := range disks {
		// Assign an LDAT to the pack, update the pack label, then rewrite the label
		ldatIndex := nextLdatIndex
		nextLdatIndex++

		// Rewrite first directory track to the pack
		dirTrack := make([]pkg.Word36, 1792)
		availableTracks := diskAttr.PackAttrs.Label[016].GetW() - 2
		recordLength := diskAttr.PackAttrs.Label[04].GetH2()
		blocksPerTrack := diskAttr.PackAttrs.Label[04].GetH1()

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

		err := false
		for blocksLeft > 0 && !err {
			subSet := dirTrack[subSetStart : subSetStart+int(recordLength)]
			pkt := nodeMgr.NewDiskIoPacketWrite(diskInfo.GetDeviceIdentifier(), blockId, subSet)
			mgr.exec.GetNodeManager().RouteIo(pkt)
			ioStat := pkt.GetIoStatus()
			if ioStat != types.IosComplete {
				msg := fmt.Sprintf("IO error reading directory track on device %v", diskInfo.GetDeviceName())
				log.Printf("MFDMgr:%v", msg)
				mgr.exec.SendExecReadOnlyMessage(msg)
				_ = mgr.exec.GetNodeManager().SetNodeStatus(diskInfo.GetNodeIdentifier(), types.NodeStatusDown)
				err = true
			}

			blocksLeft++
			subSetStart += int(recordLength)
		}

		if err {
			continue
		}

		// Now merge information for this pack into our master directory
		dirTrackSectorAddr := uint64(ldatIndex) << 18
		mgr.directoryTracks[dirTrackSectorAddr] = dirTrack

		mgr.fixedLDAT[ldatIndex] = diskInfo.GetDeviceIdentifier()
		totalTracks += availableTracks
	}

	msg := fmt.Sprintf("MS Initialized - %v Tracks Available", totalTracks)
	mgr.exec.SendExecReadOnlyMessage(msg)
	return nil
}

func (mgr *MFDManager) initializeRemovable(disks map[*nodeMgr.DiskDeviceInfo]*types.DiskAttributes) error {
	//for diskInfo, diskAttr := range disks {
	//	// TODO
	//}
	return nil
}
