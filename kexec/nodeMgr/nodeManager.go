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
	"sync"
	"time"
)

type selectionStrategy uint

const (
	StrategyFirst selectionStrategy = iota
	StrategyRoundRobin
)

// NodeManager handles the inventory of pseudo-hardware channelInfos and deviceInfos
type NodeManager struct {
	exec         types.IExec
	mutex        sync.Mutex
	threadDone   bool
	threadStop   bool
	nodeInfos    map[types.NodeIdentifier]NodeInfo    // all nodes
	channelInfos map[types.NodeIdentifier]ChannelInfo // this is loaded from the config
	deviceInfos  map[types.NodeIdentifier]DeviceInfo  // this is loaded from the config
	strategy     selectionStrategy                    // strategy used for selecting a channel fo IO
	nextChannel  []types.NodeIdentifier               // used for selecting channel to be used for IO for round-robin
}

func NewNodeManager(exec types.IExec) *NodeManager {
	return &NodeManager{
		exec:     exec,
		strategy: StrategyRoundRobin, // TODO read this from configuration
	}
}

// Boot is invoked when the exec is booting
func (mgr *NodeManager) Boot() error {
	log.Printf("NodeMgr:Boot")
	// nothing to do
	return nil
}

// Close is invoked when the application is terminating
func (mgr *NodeManager) Close() {
	log.Printf("NodeMgr:Close")
	mgr.threadStop = true
	for !mgr.threadDone {
		time.Sleep(25 * time.Millisecond)
	}

	// close devices and channels, if appropriate
}

// Initialize is invoked when the application is starting
func (mgr *NodeManager) Initialize() error {
	log.Printf("NodeMgr:Initialized")
	mgr.nodeInfos = make(map[types.NodeIdentifier]NodeInfo)
	mgr.channelInfos = make(map[types.NodeIdentifier]ChannelInfo)
	mgr.deviceInfos = make(map[types.NodeIdentifier]DeviceInfo)

	// read configuration
	// TODO from a data file or database or something
	chan0 := NewDiskChannelInfo("CHDISK")
	chan1 := NewTapeChannelInfo("CHTAPE")

	mgr.channelInfos[chan0.GetNodeIdentifier()] = chan0
	mgr.channelInfos[chan1.GetNodeIdentifier()] = chan1

	fn1 := "resources/fix000.pack"
	disk0 := NewDiskDeviceInfo("DISK0", &fn1)
	fn2 := "resources/fix001.pack"
	disk1 := NewDiskDeviceInfo("DISK1", &fn2)
	fn3 := "resources/fix002.pack"
	disk2 := NewDiskDeviceInfo("DISK2", &fn3)
	fn4 := "resources/rem000.pack"
	disk3 := NewDiskDeviceInfo("DISK3", &fn4)

	tape0 := NewTapeDeviceInfo("TAPE0")
	tape1 := NewTapeDeviceInfo("TAPE1")

	mgr.deviceInfos[disk0.nodeIdentifier] = disk0
	mgr.deviceInfos[disk1.nodeIdentifier] = disk1
	mgr.deviceInfos[disk2.nodeIdentifier] = disk2
	mgr.deviceInfos[disk3.nodeIdentifier] = disk3
	mgr.deviceInfos[tape0.nodeIdentifier] = tape0
	mgr.deviceInfos[tape1.nodeIdentifier] = tape1

	chan0.deviceInfos = []*DiskDeviceInfo{disk0, disk1, disk2, disk3}
	chan1.deviceInfos = []*TapeDeviceInfo{tape0, tape1}

	disk0.channelInfos = []*DiskChannelInfo{chan0}
	disk1.channelInfos = []*DiskChannelInfo{chan0}
	disk2.channelInfos = []*DiskChannelInfo{chan0}
	disk3.channelInfos = []*DiskChannelInfo{chan0}
	tape0.channelInfos = []*TapeChannelInfo{chan1}
	tape1.channelInfos = []*TapeChannelInfo{chan1}

	for nodeId, chInfo := range mgr.channelInfos {
		mgr.nodeInfos[nodeId] = chInfo
	}
	for nodeId, devInfo := range mgr.deviceInfos {
		mgr.nodeInfos[nodeId] = devInfo
	}
	// TODO End TODOs

	// Create channels
	for _, cInfo := range mgr.channelInfos {
		cInfo.CreateNode()
		mgr.nextChannel = append(mgr.nextChannel, cInfo.GetNodeIdentifier())
	}

	// Create devices
	for _, dInfo := range mgr.deviceInfos {
		dInfo.CreateNode()
	}

	// Connect devices to channels
	errors := false
	for cid, cInfo := range mgr.channelInfos {
		switch cInfo.GetNodeDeviceType() {
		case NodeDeviceDisk:
			dchInfo := cInfo.(*DiskChannelInfo)
			for _, dInfo := range dchInfo.deviceInfos {
				did := dInfo.GetNodeIdentifier()
				log.Printf("NodeMgr:assigning %v to %v", dInfo.GetNodeName(), cInfo.GetNodeName())
				err := mgr.channelInfos[cid].GetChannel().AssignDevice(did, mgr.deviceInfos[did].GetDevice())
				if err != nil {
					log.Printf("NodeMgr:%v", err)
					errors = true
				} else {
					dInfo.SetIsAccessible(true)
				}
			}
		case NodeDeviceTape:
			tchInfo := cInfo.(*TapeChannelInfo)
			for _, dInfo := range tchInfo.deviceInfos {
				did := dInfo.GetNodeIdentifier()
				log.Printf("NodeMgr:assigning %v to %v", dInfo.GetNodeName(), cInfo.GetNodeName())
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
		mgr.exec.Stop(types.StopInitializationSystemConfigurationError)
		return fmt.Errorf("init error")
	}

	go mgr.thread()
	return nil
}

// Stop is invoked when the exec is stopping
func (mgr *NodeManager) Stop() {
	log.Printf("NodeMgr:Stop")
	// nothing to do
}

func (mgr *NodeManager) GetChannelInfos() []ChannelInfo {
	var result = make([]ChannelInfo, len(mgr.channelInfos))
	cx := 0
	for _, chInfo := range mgr.channelInfos {
		result[cx] = chInfo
		cx++
	}
	return result
}

func (mgr *NodeManager) GetDeviceInfos() []DeviceInfo {
	var result = make([]DeviceInfo, len(mgr.deviceInfos))
	dx := 0
	for _, devInfo := range mgr.deviceInfos {
		result[dx] = devInfo
		dx++
	}
	return result
}

func (mgr *NodeManager) GetNodeInfoByName(nodeName string) (NodeInfo, error) {
	for _, chInfo := range mgr.channelInfos {
		if nodeName == chInfo.GetNodeName() {
			return chInfo, nil
		}
	}

	for _, devInfo := range mgr.deviceInfos {
		if nodeName == devInfo.GetNodeName() {
			return devInfo, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (mgr *NodeManager) GetNodeInfoByIdentifier(nodeId types.NodeIdentifier) (NodeInfo, error) {
	for _, chInfo := range mgr.channelInfos {
		if nodeId == chInfo.GetNodeIdentifier() {
			return chInfo, nil
		}
	}

	for _, devInfo := range mgr.deviceInfos {
		if nodeId == devInfo.GetNodeIdentifier() {
			return devInfo, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

// RouteIo handles all disk and tape IO for the exec
func (mgr *NodeManager) RouteIo(ioPacket IoPacket) {
	if mgr.exec.GetConfiguration().LogIOs {
		devId := pkg.Word36(ioPacket.GetNodeIdentifier())
		devName := devId.ToStringAsFieldata()
		switch ioPacket.GetNodeDeviceType() {
		case NodeDeviceDisk:
			iop := ioPacket.(*DiskIoPacket)
			log.Printf("NodeMgr:RouteIO %v iof:%v blk:%v", devName, iop.ioFunction, iop.blockId)
		case NodeDeviceTape:
			iop := ioPacket.(*TapeIoPacket)
			log.Printf("NodeMgr:RouteIO %v iof:%v", devName, iop.ioFunction)
		}
	}

	if ioPacket == nil {
		ioPacket.SetIoStatus(types.IosInternalError)
		mgr.exec.Stop(types.StopErrorInSystemIOTable)
		return
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	devInfo, ok := mgr.deviceInfos[ioPacket.GetNodeIdentifier()]
	if !ok {
		ioPacket.SetIoStatus(types.IosDeviceDoesNotExist)
		return
	}

	if !devInfo.IsAccessible() {
		ioPacket.SetIoStatus(types.IosDeviceIsNotAccessible)
		return
	}

	chInfo, err := mgr.selectChannelForDevice(devInfo)
	if err != nil {
		ioPacket.SetIoStatus(types.IosInternalError)
		mgr.exec.Stop(types.StopErrorInSystemIOTable)
		return
	}

	ioPacket.SetIoStatus(types.IosInProgress)
	chInfo.GetChannel().StartIo(ioPacket)
}

// -----------------------------------------------------------

// selectChannelForDevice chooses the *best* channel to be used for accessing the device.
// THIS MUST BE CALLED UNDER LOCK
func (mgr *NodeManager) selectChannelForDevice(devInfo DeviceInfo) (ChannelInfo, error) {
	cInfos := devInfo.GetChannelInfos()
	if len(cInfos) == 0 {
		return nil, fmt.Errorf("not accessible")
	}

	switch mgr.strategy {
	case StrategyFirst:
		// We choose the first controller in the device list
		return cInfos[0], nil

	case StrategyRoundRobin:
		// We choose the first channel in the next-channel list which is attached to the device,
		// then move that entry down to the bottom.
		// This implements a round-robin strategy which is aware of the possibility that we might have more than one
		// channel for a device, that we will likely have more than one device per channel, and that the assignation of
		// devices to channels can get a bit... messy.
		for cx, cid := range mgr.nextChannel {
			for _, cInfo := range cInfos {
				if cInfo.GetNodeIdentifier() == cid {
					// shuffle the nextChannel array
					for dx := cx; dx < len(mgr.nextChannel)-1; dx++ {
						mgr.nextChannel[dx] = mgr.nextChannel[dx+1]
					}
					mgr.nextChannel[len(mgr.nextChannel)-1] = cInfo.GetNodeIdentifier()

					// done
					return cInfo, nil
				}
			}
		}
	}

	// if we get here, something is badly wrong
	log.Printf("NodeMgr: Cannot satisfy channel selection for accessible device")
	return nil, fmt.Errorf("internal error")
}

// thread runs as long as the application is running.
// Its job is to monitor the hardware, and to send device ready notifications to the facilities manager.
// The facilities manager *may not* be up and running - it is part of the exec, so it does not run if
// the exec is not running. SO... those notifications might get dropped on the floor. That is okay.
func (mgr *NodeManager) thread() {
	mgr.threadDone = false

	for !mgr.threadStop {
		time.Sleep(25 * time.Millisecond)

		// Check devices to see if any have become ready or not ready since our last poll.
		// Make a list while we are under lock, then unlock and notify the appropriate authorities
		// of any such devices.
		updates := make(map[DeviceInfo]bool)
		mgr.mutex.Lock()
		for _, devInfo := range mgr.deviceInfos {
			if devInfo.GetDevice().IsReady() != devInfo.IsReady() {
				updates[devInfo] = devInfo.GetDevice().IsReady()
				devInfo.SetIsReady(devInfo.GetDevice().IsReady())
			}
		}
		mgr.mutex.Unlock()

		fm := mgr.exec.GetFacilitiesManager().(types.IFacilitiesManager)
		for devInfo, isReady := range updates {
			fm.NotifyDeviceReady(devInfo.GetNodeIdentifier(), isReady)
		}
	}

	mgr.threadDone = true
}

func (mgr *NodeManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vNodeManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  threadDone: %v\n", indent, mgr.threadDone)

	for _, ni := range mgr.nodeInfos {
		ni.Dump(dest, indent+"  ")
	}
}
