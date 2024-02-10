// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import (
	"fmt"
	"log"
)

// DeviceManager handles the inventory of pseudo-hardware channelInfos and deviceInfos
type DeviceManager struct {
	channelInfos     map[NodeIdentifier]ChannelInfo // this is loaded from the config
	deviceInfos      map[NodeIdentifier]DeviceInfo  // this is loaded from the config
	channelDeviceMap map[ChannelInfo][]DeviceInfo   // this is loaded from the config
	deviceChannelMap map[DeviceInfo][]ChannelInfo   // this is built dynamically from the config
}

// BuildConfiguration reads the configuration with respect to pseudo-hardware deviceInfos,
// instantiating that configuration along the way.
// This must happen very early in Exec startup, before the operator is allowed to modify the config.
// Thus, we can assume all devices are UP.
func (mgr *DeviceManager) BuildConfiguration() error {
	// read configuration
	// TODO from a data file or database or something
	chan0 := &DiskChannelInfo{
		channelName: "CHDISK",
	}

	fn := "fixed0.pack"
	disk0 := &DiskDeviceInfo{
		deviceName:      "DISK0",
		nodeStatus:      NodeStatusUp,
		initialFileName: &fn,
	}

	fn = "fixed1.pack"
	disk1 := &DiskDeviceInfo{
		deviceName:      "DISK1",
		nodeStatus:      NodeStatusUp,
		initialFileName: &fn,
	}

	mgr.channelInfos = make(map[NodeIdentifier]ChannelInfo)
	mgr.deviceInfos = make(map[NodeIdentifier]DeviceInfo)
	mgr.channelDeviceMap = make(map[ChannelInfo][]DeviceInfo)
	mgr.deviceChannelMap = make(map[DeviceInfo][]ChannelInfo)

	mgr.channelInfos[chan0.nodeIdentifier] = chan0
	mgr.deviceInfos[disk0.nodeIdentifier] = disk0
	mgr.deviceInfos[disk1.nodeIdentifier] = disk1
	mgr.channelDeviceMap[chan0] = []DeviceInfo{disk0, disk1}
	mgr.deviceChannelMap[disk0] = []ChannelInfo{chan0}
	mgr.deviceChannelMap[disk1] = []ChannelInfo{chan0}

	// Create channelInfos
	for _, cInfo := range mgr.channelInfos {
		cInfo.CreateNode()
	}

	// Create deviceInfos
	for _, dInfo := range mgr.deviceInfos {
		dInfo.CreateNode()
	}

	// Connect deviceInfos to channelInfos
	errors := false
	for cInfo := range mgr.channelDeviceMap {
		for _, dInfo := range mgr.channelDeviceMap[cInfo] {
			log.Printf("DevMgr:assigning %v to %v", dInfo.GetNodeName(), cInfo.GetNodeName())
			err := cInfo.GetChannel().AssignDevice(dInfo.GetNodeIdentifier(), dInfo.GetDevice())
			if err != nil {
				log.Printf("DevMgr:%v", err)
				errors = true
			} else {
				dInfo.SetIsAccessible(true)
			}
		}
	}

	for dInfo := range mgr.deviceChannelMap {
		if !dInfo.IsAccessible() {
			log.Printf("DevMgr:%v is not accessible", dInfo.GetNodeName())
		}
	}

	if errors {
		return fmt.Errorf("deviceManager encountered 1 or more errors during initialization")
	}

	return nil
}

// Initialize is invoked after the operator has been allowed to modify the config.
// Devices may be UP, DN, RV, or SU.
// We don't mess with tape devices - they were freshly created, thus they are not mounted.
// For disk devices, some (maybe all) are pre-mounted thus we can (if the device is not DN and is accessible)
// probe the device to try to read VOL1, S0, S1, and maybe some other interesting bits.
func (mgr *DeviceManager) Initialize() error {
	// TODO
	return nil
}

// Recover is an alternative to BuildConfiguration, and is used when the exec is re-starting.
// It is expected that the deviceInfos all need to be reset, and that some mountable need to be unmounted.
func (mgr *DeviceManager) Recover() error {
	// Reset all the deviceInfos
	errors := false
	// for cInfo := range mgr.channelDeviceMap {
	// 	for addr := range mgr.channelDeviceMap[cInfo] {
	// 		dInfo := mgr.channelDeviceMap[cInfo][addr]
	// 		if dInfo.GetDeviceStatus() == NodeStatusUp {
	// 			log.Printf("DevMgr:resetting device %v", dInfo.GetNodeName())
	//
	// 			switch cInfo.GetChannelType() {
	// 			case channelInfos.ChannelTypeDisk:
	// 				diskInfo := dInfo.(*DiskDeviceInfo)
	// 				pkt := deviceInfos.NewDiskIoPacketReset()
	// 				diskInfo.GetDevice().StartIo(pkt)
	// 				if pkt.GetIoStatus() != deviceInfos.IosComplete {
	// 					// TODO console messages?
	// 					log.Printf("DevMgr:IO error status %v", pkt.GetIoStatus())
	// 					diskInfo.nodeStatus = NodeStatusDown
	// 					log.Printf("DevMgr:Marking Device %v DN", diskInfo.GetNodeName())
	// 				} else {
	// 					// should we unmount?
	// 					if diskInfo.GetInitialFileName() == nil && diskInfo.isMounted {
	// 						log.Printf("DevMgr:dismounting media from device %v", diskInfo.GetNodeName())
	// 						pkt = deviceInfos.NewDiskIoPacketUnmount()
	// 						diskInfo.GetDevice().StartIo(pkt)
	// 						if pkt.GetIoStatus() != deviceInfos.IosComplete {
	// 							// TODO console messages?
	// 							log.Printf("DevMgr:IO error status %v", pkt.GetIoStatus())
	// 							diskInfo.nodeStatus = NodeStatusDown
	// 							log.Printf("DevMgr:Marking Device %v DN", diskInfo.GetNodeName())
	// 						}
	//
	// 						diskInfo.isMounted = false
	// 					}
	// 				}
	//
	// 				// Clear cached information about the media on the device
	// 				diskInfo.isFixed = false
	// 				diskInfo.isPrepped = false
	// 				// TODO anything else to clear out?
	//
	// 			case channelInfos.ChannelTypeTape:
	// 				tapeInfo := dInfo.(*TapeDeviceInfo)
	// 				pkt := deviceInfos.NewTapeIoPacketReset()
	// 				tapeInfo.GetDevice().StartIo(pkt)
	// 				if pkt.GetIoStatus() != deviceInfos.IosComplete {
	// 					// TODO console messages?
	// 					log.Printf("DevMgr:IO error status %v", pkt.GetIoStatus())
	// 					tapeInfo.nodeStatus = NodeStatusDown
	// 					log.Printf("DevMgr:Marking Device %v DN", tapeInfo.GetNodeName())
	// 				} else {
	// 					// should we unmount?
	// 					if tapeInfo.isMounted {
	// 						log.Printf("DevMgr:dismounting media from device %v", tapeInfo.GetNodeName())
	// 						pkt = deviceInfos.NewTapeIoPacketUnmount()
	// 						tapeInfo.GetDevice().StartIo(pkt)
	// 						if pkt.GetIoStatus() != deviceInfos.IosComplete {
	// 							// TODO console messages?
	// 							log.Printf("DevMgr:IO error status %v", pkt.GetIoStatus())
	// 							tapeInfo.nodeStatus = NodeStatusDown
	// 							log.Printf("DevMgr:Marking Device %v DN", tapeInfo.GetNodeName())
	// 						}
	//
	// 						tapeInfo.isMounted = false
	// 					}
	// 				}
	//
	// 				// Clear cached information about the media on the device
	// 				// TODO anything to clear out?
	// 			}
	// 		}
	// 	}
	// }

	err := mgr.probeMountedDisks()
	if err != nil {
		// TODO
		errors = true
	}

	if errors {
		return fmt.Errorf("deviceManager encountered 1 or more errors during initialization")
	}

	return nil
}

func (mgr *DeviceManager) probeMountedDisks() error {
	return nil
}

func IsValidPackName(name string) bool {
	if len(name) < 1 || len(name) > 6 {
		return false
	}

	if name[0] < 'A' || name[0] > 'Z' {
		return false
	}

	for nx := 1; nx < len(name); nx++ {
		if (name[nx] < 'A' || name[nx] > 'Z') && (name[nx] < '0' || name[nx] > '9') {
			return false
		}
	}

	return true
}

func IsValidPrepFactor(prepFactor PrepFactor) bool {
	return prepFactor == 28 || prepFactor == 56 || prepFactor == 112 || prepFactor == 224 ||
		prepFactor == 448 || prepFactor == 896 || prepFactor == 1792
}
