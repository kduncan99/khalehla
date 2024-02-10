// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

// func DiskMgrLoadDevices(config *deviceMgr.DeviceManager) error {
// 	for _, dd := range config.diskDevices {
// 		dd.deviceStatus = deviceMgr.NodeStatusUp
// 		dd.device = deviceMgr.NewDiskDevice()
// 		pkt := deviceMgr.NewDiskIoPacketMount(dd.fileName, false)
// 		dd.device.StartIo(pkt)
// 		if pkt.GetIoStatus() != deviceMgr.IosComplete {
// 			log.Printf("%v Cannot mount pack - IOStatus=%v", dd.deviceName, pkt.GetIoStatus())
// 			ConsoleMgrWriteRaw(fmt.Sprintf("%s Pack cannot be mounted", dd.deviceName))
// 			ConsoleMgrWriteRaw(fmt.Sprintf("%s DN", dd.deviceName))
// 			dd.deviceStatus = deviceMgr.NodeStatusDown
// 			continue
// 		}
//
// 		if !dd.device.IsPrepped() {
// 			log.Printf("%v pack is not prepped", dd.deviceName)
// 			ConsoleMgrWriteRaw(fmt.Sprintf("%s Pack is not prepped", dd.deviceName))
// 			ConsoleMgrWriteRaw(fmt.Sprintf("%s RV", dd.deviceName))
// 			dd.isPrepped = false
// 			dd.isFixed = false
// 		} else {
// 			dd.isPrepped = true
// 			// TODO we *could* read the label and s0 and s1 at this point,
// 			//  storing them in the device struct - it might help later on
// 		}
// 	}
//
// 	return nil
// }

// DiskMgrInitializeFixed iterates over the fixed packs, building up an empty fixed storage pool.
// func DiskMgrInitializeFixed(config *deviceMgr.DeviceManager) error {
// 	devs := 0
// 	tracks := 0
// 	for _, dd := range config.diskDevices {
// 		if dd.deviceStatus == deviceMgr.NodeStatusUp {
// 			buffer := make([]pkg.Word36, 28)
// 			pkt := deviceMgr.NewDiskIoPacketReadLabel(buffer)
// 			dd.device.StartIo(pkt)
// 			if pkt.GetIoStatus() != deviceMgr.IosComplete {
// 				log.Printf("%v cannot read pack label", dd.deviceName)
// 				ConsoleMgrWriteRaw(fmt.Sprintf("%s Cannot read pack label", dd.deviceName))
// 				ConsoleMgrWriteRaw(fmt.Sprintf("%s DN", dd.deviceName))
// 				dd.deviceStatus = deviceMgr.NodeStatusDown
// 				continue
// 			}
//
// 			// TODO
// 			// read label, s0, s1, first DAS
// 			// Look at first DAS - If DAS:0,H2 is zero, this pack is uninitialized
// 			// For uninitialized packs, s1:05,H1 bit35 == 0 means removable, == 1 means fixed
// 			// For initialized packs, s1:05,H1 for fixed packs, contains the LDAT
//
// 			// For uninitialized fixed, ask console pack_name TO BECOME FIXED Y/N
// 			//   If Y, initialize the thing, if N, RV it
// 			// For initialized fixed
// 			// 	copy HMBT to SMBT
// 			// 	clear out the extra directory tracks
// 			// 	add this device to the LDAT
// 			// 	Inject this device into the MFD
// 		}
// 	}
//
// 	// Tell the MFD to do general initialization
// 	// Print number of MS devices and tracks initialized - if zero, STOP
// 	ConsoleMgrWriteRaw(fmt.Sprintf("MS Initialized %v devices, %v tracks available", devs, tracks))
// 	return nil
// }

// func DiskMgrRecoverFixed(config *deviceMgr.DeviceManager) error {
// 	devs := 0
// 	allocated := 0
// 	recovered := 0
// 	for _, dd := range config.diskDevices {
// 		if dd.deviceStatus == deviceMgr.NodeStatusUp {
// 			buffer := make([]pkg.Word36, 28)
// 			pkt := deviceMgr.NewDiskIoPacketReadLabel(buffer)
// 			dd.device.StartIo(pkt)
// 			if pkt.GetIoStatus() != deviceMgr.IosComplete {
// 				log.Printf("%v cannot read pack label", dd.deviceName)
// 				ConsoleMgrWriteRaw(fmt.Sprintf("%s Cannot read pack label", dd.deviceName))
// 				ConsoleMgrWriteRaw(fmt.Sprintf("%s DN", dd.deviceName))
// 				dd.deviceStatus = deviceMgr.NodeStatusDown
// 				continue
// 			}
//
// 			// TODO
// 			// read label, s0, s1, first DAS
// 			// Look at first DAS - If DAS:0,H2 is zero, this pack is uninitialized
// 			// For uninitialized packs, s1:05,H1 bit35 == 0 means removable, == 1 means fixed
// 			// For initialized packs, s1:05,H1 for fixed packs, contains the LDAT
//
// 			// For uninitialized fixed, ask console pack_name TO BECOME FIXED Y/N
// 			//   If Y, initialize the thing, if N, RV it
// 			// For initialized fixed
// 			// 	Add this device to the LDAT
// 			// 	Inject this device into the MFD
// 		}
// 	}
//
// 	// Tell the MFD to do general recovery
// 	// Print number of MS devices and tracks recovered - if zero, STOP
// 	ConsoleMgrWriteRaw(fmt.Sprintf("MS Initialized %v devices, %v tracks allocated, %v available",
// 		devs, allocated, recovered))
// 	return nil
// }
