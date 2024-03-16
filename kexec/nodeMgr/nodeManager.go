// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/hardware"
	"khalehla/hardware/channels"
	"khalehla/hardware/ioPackets"
	"khalehla/kexec"
	"khalehla/klog"
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
	exec         kexec.IExec
	mutex        sync.Mutex
	threadDone   bool
	threadStop   bool
	nodeInfos    map[hardware.NodeIdentifier]NodeInfo    // all nodes
	channelInfos map[hardware.NodeIdentifier]ChannelInfo // this is loaded from the config
	deviceInfos  map[hardware.NodeIdentifier]DeviceInfo  // this is loaded from the config
	strategy     selectionStrategy                       // strategy used for selecting a channel fo IO
	nextChannel  []hardware.NodeIdentifier               // used for selecting channel to be used for IO for round-robin
}

func NewNodeManager(exec kexec.IExec) *NodeManager {
	return &NodeManager{
		exec:     exec,
		strategy: StrategyRoundRobin, // TODO read this from configuration
	}
}

// Boot is invoked when the exec is booting
func (mgr *NodeManager) Boot() error {
	klog.LogTrace("NodeMgr", "Boot")
	// nothing to do
	return nil
}

// Close is invoked when the application is terminating
func (mgr *NodeManager) Close() {
	klog.LogTrace("NodeMgr", "Close")
	mgr.threadStop = true
	for !mgr.threadDone {
		time.Sleep(25 * time.Millisecond)
	}

	// close devices and channels, if appropriate
}

// Initialize is invoked when the application is starting
func (mgr *NodeManager) Initialize() error {
	klog.LogTrace("NodeMgr", "Initialized")
	mgr.nodeInfos = make(map[hardware.NodeIdentifier]NodeInfo)
	mgr.channelInfos = make(map[hardware.NodeIdentifier]ChannelInfo)
	mgr.deviceInfos = make(map[hardware.NodeIdentifier]DeviceInfo)

	// read configuration
	// TODO from a data file or database or something
	chInfo0 := NewDiskChannelInfo("CHDISK")
	chInfo1 := NewTapeChannelInfo("CHTAPE")
	chInfos := []ChannelInfo{chInfo0, chInfo1}

	fn1 := "media/fix000.pack"
	ddInfo0 := NewDiskDeviceInfo("DISK0", &fn1)
	fn2 := "media/fix001.pack"
	ddInfo1 := NewDiskDeviceInfo("DISK1", &fn2)
	fn3 := "media/fix002.pack"
	ddInfo2 := NewDiskDeviceInfo("DISK2", &fn3)
	fn4 := "media/rem000.pack"
	ddInfo3 := NewDiskDeviceInfo("DISK3", &fn4)

	tdInfo0 := NewTapeDeviceInfo("TAPE0")
	tdInfo1 := NewTapeDeviceInfo("TAPE1")

	devInfos := []DeviceInfo{ddInfo0, ddInfo1, ddInfo2, ddInfo3, tdInfo0, tdInfo1}

	chInfo0.deviceInfos = []*DiskDeviceInfo{ddInfo0, ddInfo1, ddInfo2, ddInfo3}
	chInfo1.deviceInfos = []*TapeDeviceInfo{tdInfo0, tdInfo1}

	ddInfo0.channelInfos = []*DiskChannelInfo{chInfo0}
	ddInfo1.channelInfos = []*DiskChannelInfo{chInfo0}
	ddInfo2.channelInfos = []*DiskChannelInfo{chInfo0}
	ddInfo3.channelInfos = []*DiskChannelInfo{chInfo0}
	tdInfo0.channelInfos = []*TapeChannelInfo{chInfo1}
	tdInfo1.channelInfos = []*TapeChannelInfo{chInfo1}

	// Create channels
	verbose := mgr.exec.GetConfiguration().LogIOs
	for _, cInfo := range chInfos {
		cInfo.CreateNode()
		klog.LogInfoF("NodeMgr", "Created channel %v '%v'", cInfo.GetNodeIdentifier(), cInfo.GetNodeName())
		cInfo.GetChannel().SetVerbose(verbose)
		mgr.channelInfos[cInfo.GetNodeIdentifier()] = cInfo
		mgr.nodeInfos[cInfo.GetNodeIdentifier()] = cInfo
		mgr.nextChannel = append(mgr.nextChannel, cInfo.GetNodeIdentifier())
	}

	// Create devices
	for _, dInfo := range devInfos {
		dInfo.CreateNode()
		klog.LogInfoF("NodeMgr", "Created device %v '%v'", dInfo.GetNodeIdentifier(), dInfo.GetNodeName())
		dInfo.GetDevice().SetVerbose(verbose)
		mgr.deviceInfos[dInfo.GetNodeIdentifier()] = dInfo
		mgr.nodeInfos[dInfo.GetNodeIdentifier()] = dInfo
	}

	// Connect devices to channels
	errors := false
	for cid, cInfo := range mgr.channelInfos {
		switch cInfo.GetNodeDeviceType() {
		case hardware.NodeDeviceDisk:
			dchInfo := cInfo.(*DiskChannelInfo)
			for _, dInfo := range dchInfo.deviceInfos {
				did := dInfo.GetNodeIdentifier()
				klog.LogInfoF("NodeMgr", "assigning %v to %v", dInfo.GetNodeName(), cInfo.GetNodeName())
				err := mgr.channelInfos[cid].GetChannel().AssignDevice(did, mgr.deviceInfos[did].GetDevice())
				if err != nil {
					klog.LogError("NodeMgr", err.Error())
					errors = true
				} else {
					dInfo.SetIsAccessible(true)
				}
			}
		case hardware.NodeDeviceTape:
			tchInfo := cInfo.(*TapeChannelInfo)
			for _, dInfo := range tchInfo.deviceInfos {
				did := dInfo.GetNodeIdentifier()
				klog.LogInfoF("NodeMgr", "assigning %v to %v", dInfo.GetNodeName(), cInfo.GetNodeName())
				err := mgr.channelInfos[cid].GetChannel().AssignDevice(did, mgr.deviceInfos[did].GetDevice())
				if err != nil {
					klog.LogError("NodeMgr", err.Error())
					errors = true
				} else {
					dInfo.SetIsAccessible(true)
				}
			}
		}
	}

	for _, dInfo := range mgr.deviceInfos {
		if !dInfo.IsAccessible() {
			klog.LogInfoF("NodeMgr", "%v is not accessible", dInfo.GetNodeName())
		}
	}

	if errors {
		klog.LogFatal("NodeMgr", "initialization errors exist")
		mgr.exec.Stop(kexec.StopInitializationSystemConfigurationError)
		return fmt.Errorf("init error")
	}

	go mgr.thread()
	return nil
}

// Stop is invoked when the exec is stopping
func (mgr *NodeManager) Stop() {
	klog.LogTrace("NodeMgr", "Stop")
	// Reset all devices and channels
	for _, ci := range mgr.channelInfos {
		ci.GetChannel().Reset()
	}
	for _, di := range mgr.deviceInfos {
		di.GetDevice().Reset()
	}
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

func (mgr *NodeManager) GetNodeInfoByIdentifier(nodeId hardware.NodeIdentifier) (NodeInfo, error) {
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
func (mgr *NodeManager) RouteIo(cp *channels.ChannelProgram) {
	klog.LogTraceF("NodeMgr", "RouteIO %v", cp.GetString())

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	devInfo, ok := mgr.deviceInfos[cp.NodeIdentifier]
	if !ok {
		cp.IoStatus = ioPackets.IosDeviceDoesNotExist
		return
	}

	if !devInfo.IsAccessible() {
		cp.IoStatus = ioPackets.IosDeviceIsNotAccessible
		return
	}

	chInfo, err := mgr.selectChannelForDevice(devInfo)
	if err != nil {
		cp.IoStatus = ioPackets.IosInternalError
		mgr.exec.Stop(kexec.StopErrorInSystemIOTable)
		return
	}

	chInfo.GetChannel().StartIo(cp)
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
	klog.LogError("NodeMgr", "Cannot satisfy channel selection for accessible device")
	mgr.exec.Stop(kexec.StopErrorAccessingFacilitiesDataStructure)
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

		listener := mgr.exec.GetFacilitiesManager().(IDeviceListener)
		for devInfo, isReady := range updates {
			listener.NotifyDeviceReady(devInfo.GetNodeIdentifier(), isReady)
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
