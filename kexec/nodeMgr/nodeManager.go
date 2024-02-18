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
	exec            types.IExec
	mutex           sync.Mutex
	isInitialized   bool
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	nodeInfos       map[types.NodeIdentifier]types.NodeInfo       // all nodes
	channelInfos    map[types.ChannelIdentifier]types.ChannelInfo // this is loaded from the config
	deviceInfos     map[types.DeviceIdentifier]types.DeviceInfo   // this is loaded from the config
	strategy        selectionStrategy                             // strategy used for selecting a channel fo IO
	nextChannel     []types.ChannelIdentifier                     // used for selecting channel to be used for IO for round-robin
}

func NewNodeManager(exec types.IExec) *NodeManager {
	return &NodeManager{
		exec:     exec,
		strategy: StrategyRoundRobin, // TODO read this from configuration
	}
}

func (mgr *NodeManager) CloseManager() {
	mgr.threadStop()
	mgr.isInitialized = false
}

// InitializeManager reads the configuration with respect to pseudo-hardware deviceInfos,
// instantiating that configuration along the way.
// This must happen very early in Exec startup, before the operator is allowed to modify the config.
// Thus, we can assume all devices are UP.
func (mgr *NodeManager) InitializeManager() error {
	mgr.nodeInfos = make(map[types.NodeIdentifier]types.NodeInfo)
	mgr.channelInfos = make(map[types.ChannelIdentifier]types.ChannelInfo)
	mgr.deviceInfos = make(map[types.DeviceIdentifier]types.DeviceInfo)

	// read configuration
	// TODO from a data file or database or something
	chan0 := NewDiskChannelInfo("CHDISK")
	chan1 := NewTapeChannelInfo("CHTAPE")

	mgr.channelInfos[chan0.channelIdentifier] = chan0
	mgr.channelInfos[chan1.channelIdentifier] = chan1

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

	mgr.deviceInfos[disk0.deviceIdentifier] = disk0
	mgr.deviceInfos[disk1.deviceIdentifier] = disk1
	mgr.deviceInfos[disk2.deviceIdentifier] = disk2
	mgr.deviceInfos[disk3.deviceIdentifier] = disk3
	mgr.deviceInfos[tape0.deviceIdentifier] = tape0
	mgr.deviceInfos[tape1.deviceIdentifier] = tape1

	chan0.deviceInfos = []*DiskDeviceInfo{disk0, disk1, disk2, disk3}
	chan1.deviceInfos = []*TapeDeviceInfo{tape0, tape1}

	disk0.channelInfos = []*DiskChannelInfo{chan0}
	disk1.channelInfos = []*DiskChannelInfo{chan0}
	disk2.channelInfos = []*DiskChannelInfo{chan0}
	disk3.channelInfos = []*DiskChannelInfo{chan0}
	tape0.channelInfos = []*TapeChannelInfo{chan1}
	tape1.channelInfos = []*TapeChannelInfo{chan1}

	for chId, chInfo := range mgr.channelInfos {
		mgr.nodeInfos[types.NodeIdentifier(chId)] = chInfo
	}
	for devId, devInfo := range mgr.deviceInfos {
		mgr.nodeInfos[types.NodeIdentifier(devId)] = devInfo
	}
	// TODO End TODOs

	// Create channels
	for _, cInfo := range mgr.channelInfos {
		cInfo.CreateNode()
		mgr.nextChannel = append(mgr.nextChannel, cInfo.GetChannelIdentifier())
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
		mgr.exec.Stop(types.StopInitializationSystemConfigurationError)
		return fmt.Errorf("init error")
	}

	mgr.threadStart()
	mgr.isInitialized = true
	return nil
}

func (mgr *NodeManager) IsInitialized() bool {
	return mgr.isInitialized
}

func (mgr *NodeManager) ResetManager() error {
	mgr.threadStop()
	mgr.threadStart()

	// TODO should we do anything here?

	mgr.isInitialized = true
	return nil
}

// -----------------------------------------------------------

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

func (mgr *NodeManager) GetChannelInfos() []types.ChannelInfo {
	var result = make([]types.ChannelInfo, len(mgr.channelInfos))
	cx := 0
	for _, chInfo := range mgr.channelInfos {
		result[cx] = chInfo
		cx++
	}
	return result
}

func (mgr *NodeManager) GetDeviceInfos() []types.DeviceInfo {
	var result = make([]types.DeviceInfo, len(mgr.deviceInfos))
	dx := 0
	for _, devInfo := range mgr.deviceInfos {
		result[dx] = devInfo
		dx++
	}
	return result
}

func (mgr *NodeManager) GetNodeInfoByName(nodeName string) (types.NodeInfo, error) {
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

func (mgr *NodeManager) GetNodeInfoByIdentifier(nodeId types.NodeIdentifier) (types.NodeInfo, error) {
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
func (mgr *NodeManager) RouteIo(ioPacket types.IoPacket) {
	if mgr.exec.GetConfiguration().LogIOs {
		devId := pkg.Word36(ioPacket.GetDeviceIdentifier())
		devName := devId.ToStringAsFieldata()
		switch ioPacket.GetNodeType() {
		case types.NodeTypeDisk:
			iop := ioPacket.(*DiskIoPacket)
			log.Printf("NodeMgr:RouteIO %v iof:%v blk:%v", devName, iop.ioFunction, iop.blockId)
		case types.NodeTypeTape:
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

	devInfo, ok := mgr.deviceInfos[ioPacket.GetDeviceIdentifier()]
	if !ok {
		ioPacket.SetIoStatus(types.IosDeviceDoesNotExist)
		return
	}

	if !devInfo.IsAccessible() {
		ioPacket.SetIoStatus(types.IosDeviceIsNotAccessible)
		return
	}

	if devInfo.GetNodeStatus() == types.NodeStatusDown {
		ioPacket.SetIoStatus(types.IosDeviceIsDown)
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

func (mgr *NodeManager) SetNodeStatus(nodeId types.NodeIdentifier, status types.NodeStatus) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	ni, ok := mgr.nodeInfos[nodeId]
	if !ok {
		return fmt.Errorf("node not found")
	}

	// TODO should we be sending NotifyDeviceReady to various listeners?
	//	Only makes sense for disk really, and we have FA keyin for that

	// for now, we do not allow changing status of anything except devices
	if ni.GetNodeCategory() != types.NodeCategoryDevice {
		return fmt.Errorf("not allowed")
	}

	devInfo := ni.(types.DeviceInfo)
	stopExec := false
	switch status {
	case types.NodeStatusDown:
		if ni.GetNodeType() == types.NodeTypeDisk {
			ddInfo := ni.(*DiskDeviceInfo)
			ddInfo.nodeStatus = status
			stopExec = mgr.exec.GetFacilitiesManager().IsDeviceAssigned(devInfo.GetDeviceIdentifier())
			break
		} else if ni.GetNodeType() == types.NodeTypeTape {
			// reset the tape device (unmounts it as part of the process)
			tdInfo := ni.(*TapeDeviceInfo)
			ioPkt := NewTapeIoPacketReset(devInfo.GetDeviceIdentifier())
			chInfo, err := mgr.selectChannelForDevice(tdInfo)
			if err != nil {
				chInfo.GetChannel().StartIo(ioPkt)
			}

			tdInfo.nodeStatus = status

			if mgr.exec.GetFacilitiesManager().IsDeviceAssigned(tdInfo.GetDeviceIdentifier()) {
				// TODO - tell Exec to abort the run to which the thing was assigned
			}
			break
		}

		// anything other than disk or tape
		return fmt.Errorf("not allowed")

	case types.NodeStatusReserved:
		if ni.GetNodeType() == types.NodeTypeDisk {
			ddInfo := ni.(*DiskDeviceInfo)
			ddInfo.nodeStatus = status
		} else if ni.GetNodeType() == types.NodeTypeTape {
			tdInfo := ni.(*TapeDeviceInfo)
			tdInfo.nodeStatus = status
		}

		// anything other than disk or tape
		return fmt.Errorf("not allowed")

	case types.NodeStatusSuspended:
		if ni.GetNodeType() == types.NodeTypeDisk {
			ddInfo := ni.(*DiskDeviceInfo)
			ddInfo.nodeStatus = status
			break
		}

		// anything other than disk or tape
		return fmt.Errorf("not allowed")

	case types.NodeStatusUp:
		if ni.GetNodeType() == types.NodeTypeDisk {
			ddInfo := ni.(*DiskDeviceInfo)
			ddInfo.nodeStatus = status
		} else if ni.GetNodeType() == types.NodeTypeTape {
			tdInfo := ni.(*TapeDeviceInfo)
			tdInfo.nodeStatus = status
		}

		// anything other than disk or tape
		return fmt.Errorf("not allowed")

	default:
		return fmt.Errorf("internal error")
	}

	msg := ni.GetNodeName() + " " + GetNodeStatusString(ni.GetNodeStatus(), ni.IsAccessible())
	mgr.exec.SendExecReadOnlyMessage(msg, nil)
	if stopExec {
		mgr.exec.Stop(types.StopConsoleResponseRequiresReboot)
	}

	return nil
}

// -----------------------------------------------------------

// selectChannelForDevice chooses the *best* channel to be used for accessing the device.
// THIS MUST BE CALLED UNDER LOCK
func (mgr *NodeManager) selectChannelForDevice(devInfo types.DeviceInfo) (types.ChannelInfo, error) {
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
				if cInfo.GetChannelIdentifier() == cid {
					// shuffle the nextChannel array
					for dx := cx; dx < len(mgr.nextChannel)-1; dx++ {
						mgr.nextChannel[dx] = mgr.nextChannel[dx+1]
					}
					mgr.nextChannel[len(mgr.nextChannel)-1] = cInfo.GetChannelIdentifier()

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

// -----------------------------------------------------------

func (mgr *NodeManager) thread() {
	mgr.threadStarted = true

	for !mgr.terminateThread {
		time.Sleep(25 * time.Millisecond)

		// Check devices to see if any have become ready or not ready since our last poll.
		// Make a list while we are under lock, then unlock and notify the appropriate authorities
		// of any such devices.
		updates := make(map[types.DeviceInfo]bool)
		mgr.mutex.Lock()
		for _, devInfo := range mgr.deviceInfos {
			if devInfo.GetDevice().IsReady() != devInfo.IsReady() {
				updates[devInfo] = devInfo.GetDevice().IsReady()
				devInfo.SetIsReady(devInfo.GetDevice().IsReady())
			}
		}
		mgr.mutex.Unlock()

		fm := mgr.exec.GetFacilitiesManager()
		mm := mgr.exec.GetMFDManager()
		for devInfo, isReady := range updates {
			fm.NotifyDeviceReady(devInfo, isReady)
			mm.NotifyDeviceReady(devInfo, isReady)
		}
	}

	mgr.threadStopped = true
}

func (mgr *NodeManager) threadStart() {
	mgr.terminateThread = false
	if !mgr.threadStarted {
		go mgr.thread()
		for !mgr.threadStarted {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *NodeManager) threadStop() {
	if mgr.threadStarted {
		mgr.terminateThread = true
		for !mgr.threadStopped {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *NodeManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vNodeManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  initialized:     %v\n", indent, mgr.isInitialized)
	_, _ = fmt.Fprintf(dest, "%v  threadStarted:   %v\n", indent, mgr.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:   %v\n", indent, mgr.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, mgr.terminateThread)

	for _, ni := range mgr.nodeInfos {
		ni.Dump(dest, indent+"  ")
	}
}
