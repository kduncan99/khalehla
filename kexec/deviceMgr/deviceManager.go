// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"log"
)

// DeviceManager handles the inventory of pseudo-hardware channelInfos and deviceInfos
type DeviceManager struct {
	exec         types.IExec
	channelInfos map[types.ChannelIdentifier]types.ChannelInfo // this is loaded from the config
	deviceInfos  map[types.DeviceIdentifier]types.DeviceInfo   // this is loaded from the config
}

func NewDeviceManager(exec types.IExec) *DeviceManager {
	return &DeviceManager{
		exec: exec,
	}
}

func (mgr *DeviceManager) CloseManager() {
	// nothing to do for now
}

func (mgr *DeviceManager) InitializeManager() {
	mgr.channelInfos = make(map[types.ChannelIdentifier]types.ChannelInfo)
	mgr.deviceInfos = make(map[types.DeviceIdentifier]types.DeviceInfo)
}

func (mgr *DeviceManager) ResetManager() {
	// nothing to do for now
}

// BuildConfiguration reads the configuration with respect to pseudo-hardware deviceInfos,
// instantiating that configuration along the way.
// This must happen very early in Exec startup, before the operator is allowed to modify the config.
// Thus, we can assume all devices are UP.
func (mgr *DeviceManager) BuildConfiguration() error {
	// read configuration
	// TODO from a data file or database or something
	chan0 := NewDiskChannelInfo("CHDISK")
	chan1 := NewTapeChannelInfo("CHTAPE")

	mgr.channelInfos[chan0.channelIdentifier] = chan0
	mgr.channelInfos[chan1.channelIdentifier] = chan1

	fn := "disk0.pack"
	disk0 := NewDiskDeviceInfo("DISK0", &fn)
	fn = "disk1.pack"
	disk1 := NewDiskDeviceInfo("DISK1", &fn)
	fn = "disk2.pack"
	disk2 := NewDiskDeviceInfo("DISK2", &fn)

	tape0 := NewTapeDeviceInfo("TAPE0")
	tape1 := NewTapeDeviceInfo("TAPE1")

	mgr.deviceInfos[disk0.deviceIdentifier] = disk0
	mgr.deviceInfos[disk1.deviceIdentifier] = disk1
	mgr.deviceInfos[disk2.deviceIdentifier] = disk2
	mgr.deviceInfos[tape0.deviceIdentifier] = tape0
	mgr.deviceInfos[tape1.deviceIdentifier] = tape1

	chan0.deviceInfos = []*DiskDeviceInfo{disk0, disk1, disk2}
	chan1.deviceInfos = []*TapeDeviceInfo{tape0, tape1}

	disk0.channelInfos = []*DiskChannelInfo{chan0}
	disk1.channelInfos = []*DiskChannelInfo{chan0}
	disk2.channelInfos = []*DiskChannelInfo{chan0}
	tape0.channelInfos = []*TapeChannelInfo{chan1}
	tape1.channelInfos = []*TapeChannelInfo{chan1}
	// TODO End TODOs

	// Create channels
	for _, cInfo := range mgr.channelInfos {
		cInfo.CreateNode()
	}

	// Create devices
	for _, dInfo := range mgr.deviceInfos {
		dInfo.CreateNode()
	}

	// Connect devices to channels
	errors := false
	for cid, cInfo := range mgr.channelInfos {
		switch cInfo.GetNodeType() {
		case types.NodeTypeDisk:
			dchInfo := cInfo.(*DiskChannelInfo)
			for _, dInfo := range dchInfo.deviceInfos {
				did := dInfo.GetDeviceIdentifier()
				log.Printf("DevMgr:assigning %v to %v", dInfo.GetDeviceName(), cInfo.GetNodeName())
				err := mgr.channelInfos[cid].GetChannel().AssignDevice(did, mgr.deviceInfos[did].GetDevice())
				if err != nil {
					log.Printf("DevMgr:%v", err)
					errors = true
				} else {
					dInfo.SetIsAccessible(true)
				}
			}
		case types.NodeTypeTape:
			tchInfo := cInfo.(*TapeChannelInfo)
			for _, dInfo := range tchInfo.deviceInfos {
				did := dInfo.GetDeviceIdentifier()
				log.Printf("DevMgr:assigning %v to %v", dInfo.GetDeviceName(), cInfo.GetNodeName())
				err := mgr.channelInfos[cid].GetChannel().AssignDevice(did, mgr.deviceInfos[did].GetDevice())
				if err != nil {
					log.Printf("DevMgr:%v", err)
					errors = true
				} else {
					dInfo.SetIsAccessible(true)
				}
			}
		}
	}

	for _, dInfo := range mgr.deviceInfos {
		if !dInfo.IsAccessible() {
			log.Printf("DevMgr:%v is not accessible", dInfo.GetNodeName())
		}
	}

	if errors {
		return fmt.Errorf("deviceManager encountered 1 or more errors during initialization")
	}

	return nil
}

// InitializeDevices is invoked after the operator has been allowed to modify the config.
// Devices may be UP, DN, RV, or SU.
// We don't mess with tape devices - they were freshly created, thus they are not mounted.
// For disk devices, some (maybe all) are pre-mounted thus we can (if the device is not DN and is accessible)
// probe the device to try to read VOL1, S0, S1, and maybe some other interesting bits.
func (mgr *DeviceManager) InitializeDevices() error {
	// TODO
	return nil
}

// RecoverDevices is an alternative to BuildConfiguration, and is used when the exec is re-starting.
// It is expected that the deviceInfos all need to be reset, and that some mountable need to be unmounted.
func (mgr *DeviceManager) RecoverDevices() error {
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

func GetNodeStatusString(status types.NodeStatus, isAccessible bool) string {
	accStr := ""
	if !isAccessible {
		accStr = "/NA"
	}

	switch status {
	case types.NodeStatusDown:
		return "DN" + accStr
	case types.NodeStatusReserved:
		return "RV" + accStr
	case types.NodeStatusSuspended:
		return "SU" + accStr
	case types.NodeStatusUp:
		return "UP" + accStr
	}

	return ""
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

func IsValidPrepFactor(prepFactor types.PrepFactor) bool {
	return prepFactor == 28 || prepFactor == 56 || prepFactor == 112 || prepFactor == 224 ||
		prepFactor == 448 || prepFactor == 896 || prepFactor == 1792
}

func (mgr *DeviceManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vDeviceManager ----------------------------------------------------\n", indent)

	for _, chInfo := range mgr.channelInfos {
		_, _ = fmt.Fprintf(dest, "%v  Channel %v:\n", indent, chInfo.GetChannelName())
		//		chInfo.Dump(dest, indent+"  ")
	}

	for _, devInfo := range mgr.deviceInfos {
		_, _ = fmt.Fprintf(dest, "%v  Device %v:\n", indent, devInfo.GetDeviceName())
		//		devInfo.Dump(dest, indent+"  ")
	}
}
