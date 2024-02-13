// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"khalehla/pkg"
	"log"
)

// NodeManager handles the inventory of pseudo-hardware channelInfos and deviceInfos
type NodeManager struct {
	exec         types.IExec
	channelInfos map[types.ChannelIdentifier]types.ChannelInfo // this is loaded from the config
	deviceInfos  map[types.DeviceIdentifier]types.DeviceInfo   // this is loaded from the config
}

func NewNodeManager(exec types.IExec) *NodeManager {
	return &NodeManager{
		exec: exec,
	}
}

func (mgr *NodeManager) CloseManager() {
	// nothing to do for now
}

func (mgr *NodeManager) InitializeManager() {
	mgr.channelInfos = make(map[types.ChannelIdentifier]types.ChannelInfo)
	mgr.deviceInfos = make(map[types.DeviceIdentifier]types.DeviceInfo)
}

func (mgr *NodeManager) ResetManager() {
	// nothing to do for now
}

// BuildConfiguration reads the configuration with respect to pseudo-hardware deviceInfos,
// instantiating that configuration along the way.
// This must happen very early in Exec startup, before the operator is allowed to modify the config.
// Thus, we can assume all devices are UP.
func (mgr *NodeManager) BuildConfiguration() error {
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
				log.Printf("NodeMgr:assigning %v to %v", dInfo.GetDeviceName(), cInfo.GetNodeName())
				err := mgr.channelInfos[cid].GetChannel().AssignDevice(did, mgr.deviceInfos[did].GetDevice())
				if err != nil {
					log.Printf("NodeMgr:%v", err)
					errors = true
				} else {
					dInfo.SetIsAccessible(true)
				}
			}
		case types.NodeTypeTape:
			tchInfo := cInfo.(*TapeChannelInfo)
			for _, dInfo := range tchInfo.deviceInfos {
				did := dInfo.GetDeviceIdentifier()
				log.Printf("NodeMgr:assigning %v to %v", dInfo.GetDeviceName(), cInfo.GetNodeName())
				err := mgr.channelInfos[cid].GetChannel().AssignDevice(did, mgr.deviceInfos[did].GetDevice())
				if err != nil {
					log.Printf("NodeMgr:%v", err)
					errors = true
				} else {
					dInfo.SetIsAccessible(true)
				}
			}
		}
	}

	for _, dInfo := range mgr.deviceInfos {
		if !dInfo.IsAccessible() {
			log.Printf("NodeMgr:%v is not accessible", dInfo.GetNodeName())
		}
	}

	if errors {
		return fmt.Errorf("errors during initialization")
	}

	return nil
}

func (mgr *NodeManager) ForceUnitAttention() {
	// TODO
}

func (mgr *NodeManager) GetChannelInfos() []types.ChannelInfo {
	var result = make([]types.ChannelInfo, len(mgr.channelInfos))
	for cx, chInfo := range mgr.channelInfos {
		result[cx] = chInfo
	}
	return result
}

func (mgr *NodeManager) GetDeviceInfos() []types.DeviceInfo {
	var result = make([]types.DeviceInfo, len(mgr.deviceInfos))
	for dx, devInfo := range mgr.deviceInfos {
		result[dx] = devInfo
	}
	return result
}

func (mgr *NodeManager) getNodeStatusStringForChannel(chInfo types.ChannelInfo) string {
	return chInfo.GetChannelName() + " " + GetNodeStatusString(chInfo.GetNodeStatus(), true)
}

func (mgr *NodeManager) getNodeStatusStringForDevice(devInfo types.DeviceInfo) string {
	str := devInfo.GetDeviceName() + " " + GetNodeStatusString(devInfo.GetNodeStatus(), devInfo.IsAccessible())

	switch devInfo.GetNodeType() {
	case types.NodeTypeDisk:
		diskInfo := devInfo.(*DiskDeviceInfo)
		if diskInfo.isReady {
			// TODO
			//	DISK0 UP [NA] [* [F|R] PACKID packName
			//	So we need a lot of additional information in devInfo
		}
	case types.NodeTypeTape:
		tapeInfo := devInfo.(*TapeDeviceInfo)
		if tapeInfo.IsMounted() {
			// TODO
			//  TAPE0 UP[,ACS][,CTL][,PM] [NA] [* RUNID runid REEL reel [RING|NORING] [POS [LOST|j[+|-]k]]]
			//	reel can be L-BLNK for labeled blank or U-BLNK for unlabeled blank
			//	j is number of files extended
			//	k is number of blocks extended + forward, or - backward
			//	So we need a lot of additional information in devInfo
		}
	}

	return str
}

func (mgr *NodeManager) GetNodeStatusStringForNode(nodeName string) (string, error) {
	var nodeInfo types.NodeInfo
	for _, chInfo := range mgr.channelInfos {
		if nodeName == chInfo.GetNodeName() {
			return mgr.getNodeStatusStringForChannel(chInfo), nil
		}
	}

	if nodeInfo == nil {
		for _, devInfo := range mgr.deviceInfos {
			if nodeName == devInfo.GetNodeName() {
				return mgr.getNodeStatusStringForDevice(devInfo), nil
			}
		}
	}

	return "", fmt.Errorf("not found")
}

// InitializeDevices is invoked after the operator has been allowed to modify the config.
// Devices may be UP, DN, RV, or SU.
// We don't mess with tape devices - they were freshly created, thus they are not mounted.
// For disk devices, some (maybe all) are pre-mounted thus we can (if the device is not DN and is accessible)
// probe the device to try to read VOL1, S0, S1, and maybe some other interesting bits.
func (mgr *NodeManager) InitializeDevices() error {
	// TODO
	return nil
}

// RecoverDevices is an alternative to BuildConfiguration, and is used when the exec is re-starting.
// It is expected that the deviceInfos all need to be reset, and that some mountable need to be unmounted.
func (mgr *NodeManager) RecoverDevices() error {
	// Reset all the deviceInfos
	errors := false
	// for cInfo := range mgr.channelDeviceMap {
	// 	for addr := range mgr.channelDeviceMap[cInfo] {
	// 		dInfo := mgr.channelDeviceMap[cInfo][addr]
	// 		if dInfo.GetDeviceStatus() == NodeStatusUp {
	// 			log.Printf("NodeMgr:resetting device %v", dInfo.GetNodeName())
	//
	// 			switch cInfo.GetChannelType() {
	// 			case channelInfos.ChannelTypeDisk:
	// 				diskInfo := dInfo.(*DiskDeviceInfo)
	// 				pkt := deviceInfos.NewDiskIoPacketReset()
	// 				diskInfo.GetDevice().StartIo(pkt)
	// 				if pkt.GetIoStatus() != deviceInfos.IosComplete {
	// 					// TODO console messages?
	// 					log.Printf("NodeMgr:IO error status %v", pkt.GetIoStatus())
	// 					diskInfo.nodeStatus = NodeStatusDown
	// 					log.Printf("NodeMgr:Marking Device %v DN", diskInfo.GetNodeName())
	// 				} else {
	// 					// should we unmount?
	// 					if diskInfo.GetInitialFileName() == nil && diskInfo.isMounted {
	// 						log.Printf("NodeMgr:dismounting media from device %v", diskInfo.GetNodeName())
	// 						pkt = deviceInfos.NewDiskIoPacketUnmount()
	// 						diskInfo.GetDevice().StartIo(pkt)
	// 						if pkt.GetIoStatus() != deviceInfos.IosComplete {
	// 							// TODO console messages?
	// 							log.Printf("NodeMgr:IO error status %v", pkt.GetIoStatus())
	// 							diskInfo.nodeStatus = NodeStatusDown
	// 							log.Printf("NodeMgr:Marking Device %v DN", diskInfo.GetNodeName())
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
	// 					log.Printf("NodeMgr:IO error status %v", pkt.GetIoStatus())
	// 					tapeInfo.nodeStatus = NodeStatusDown
	// 					log.Printf("NodeMgr:Marking Device %v DN", tapeInfo.GetNodeName())
	// 				} else {
	// 					// should we unmount?
	// 					if tapeInfo.isMounted {
	// 						log.Printf("NodeMgr:dismounting media from device %v", tapeInfo.GetNodeName())
	// 						pkt = deviceInfos.NewTapeIoPacketUnmount()
	// 						tapeInfo.GetDevice().StartIo(pkt)
	// 						if pkt.GetIoStatus() != deviceInfos.IosComplete {
	// 							// TODO console messages?
	// 							log.Printf("NodeMgr:IO error status %v", pkt.GetIoStatus())
	// 							tapeInfo.nodeStatus = NodeStatusDown
	// 							log.Printf("NodeMgr:Marking Device %v DN", tapeInfo.GetNodeName())
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
		return fmt.Errorf("errors during initialization")
	}

	return nil
}

// probeDiskDevice reads VOL1 from the disk device, updating the various pack-specific
// elements of the device info struct as appropriate.
// If we return an error, it means the label could not be read or was invalid,
// and some higher-level code should mark this device DN.
func (mgr *NodeManager) probeDiskDevice(info *DiskDeviceInfo) error {
	// read label and s1
	// For uninitialized packs, s1:05,H1 0 means removable, 0400000 means fixed
	// For initialized packs, s1:05,H1 for fixed packs, contains the (non-zero) LDAT
	info.isReady = false
	info.packName = ""
	info.isPrepped = false
	info.isFixed = false
	info.geometry = nil

	if !info.device.isReady {
		log.Printf("NodeMgr:Cannot probe %v device not ready", info.GetDeviceName())
		return fmt.Errorf("not ready")
	}

	if !info.isAccessible {
		log.Printf("NodeMgr:Cannot probe %v device not accessible", info.GetDeviceName())
		return fmt.Errorf("not accessible")
	}

	// read label
	var label = make([]pkg.Word36, 28)
	chInfo, _ := mgr.selectChannelForDevice(info)
	pkt := NewDiskIoPacketReadLabel(info.deviceIdentifier, label)
	chInfo.GetChannel().StartIo(pkt)
	if pkt.ioStatus != types.IosComplete {
		log.Printf("NodeMgr:Failed to probe %v IO status = %v", info.GetDeviceName(), pkt.ioStatus)
		return fmt.Errorf("io error")
	}

	if label[0].ToStringAsAscii() != "VOL1" {
		log.Printf("NodeMgr:No VOL1 label on device %v", info.GetDeviceName())
		return fmt.Errorf("no VOL1")
	}

	// read block containing S1, then set s1 variable to refer to the slice of the resulting block
	// which actually represents S1.
	geometry := info.device.GetGeometry()
	dirBlock := make([]pkg.Word36, geometry.PrepFactor)
	dirBlockId := types.BlockId(geometry.FirstDirTrackBlockId)
	if geometry.PrepFactor == 28 {
		dirBlockId++
	}

	pkt = NewDiskIoPacketRead(info.deviceIdentifier, dirBlockId, dirBlock)
	chInfo.GetChannel().StartIo(pkt)
	if pkt.ioStatus != types.IosComplete {
		log.Printf("NodeMgr:Failed to probe %v IO status = %v", info.GetDeviceName(), pkt.ioStatus)
		return fmt.Errorf("io error")
	}

	var s1 []pkg.Word36
	if geometry.PrepFactor == 28 {
		s1 = dirBlock
	} else {
		s1 = dirBlock[28:56]
	}

	ldatIndex := int(s1[5].GetH1())
	if ldatIndex == 0 {
		info.isFixed = false
		info.ldatIndex = ldatIndex
	} else if ldatIndex == 0_400000 {
		info.isFixed = true
		info.ldatIndex = 0
	} else {
		info.isFixed = true
		info.ldatIndex = ldatIndex
	}

	info.packName = label[1].ToStringAsAscii() + label[2].ToStringAsAscii()
	info.packName = info.packName[0:6]
	info.isPrepped = true
	info.geometry = info.device.GetGeometry()
	info.isReady = true

	return nil
}

func (mgr *NodeManager) probeTapeDevice(info *TapeDeviceInfo) {
	// TODO
}

func (mgr *NodeManager) probeMountedDisks() error {
	// TODO
	return nil
}

func (mgr *NodeManager) selectChannelForDevice(devInfo types.DeviceInfo) (types.ChannelInfo, error) {
	// TODO implement round-robin strategy... will require some intelligence
	//	for now, just pick the first in the list
	cInfos := devInfo.GetChannelInfos()
	if len(cInfos) == 0 {
		return nil, fmt.Errorf("not accessible")
	} else {
		return cInfos[0], nil
	}
}

func GetNodeStatusString(status types.NodeStatus, isAccessible bool) string {
	accStr := ""
	if !isAccessible {
		accStr = " NA"
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

func IsValidNodeName(name string) bool {
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

func (mgr *NodeManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vNodeManager ----------------------------------------------------\n", indent)

	for _, chInfo := range mgr.channelInfos {
		_, _ = fmt.Fprintf(dest, "%v  Channel %v:\n", indent, chInfo.GetChannelName())
		//		chInfo.Dump(dest, indent+"  ")
	}

	for _, devInfo := range mgr.deviceInfos {
		_, _ = fmt.Fprintf(dest, "%v  Device %v:\n", indent, devInfo.GetDeviceName())
		//		devInfo.Dump(dest, indent+"  ")
	}
}
