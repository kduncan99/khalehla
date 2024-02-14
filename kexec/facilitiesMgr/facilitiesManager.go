// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/types"
	"khalehla/pkg"
	"log"
	"sync"
	"time"
)

type inventory struct {
	disks map[types.DeviceIdentifier]*types.DiskAttributes
	tapes map[types.DeviceIdentifier]*types.TapeAttributes
}

func newInventory() *inventory {
	return &inventory{
		disks: make(map[types.DeviceIdentifier]*types.DiskAttributes),
		tapes: make(map[types.DeviceIdentifier]*types.TapeAttributes),
	}
}

type FacilitiesManager struct {
	exec                         types.IExec
	mutex                        sync.Mutex
	isInitialized                bool
	terminateThread              bool
	threadStarted                bool
	threadStopped                bool
	inventory                    *inventory
	deviceReadyNotificationQueue map[types.DeviceIdentifier]bool
}

func NewFacilitiesManager(exec types.IExec) *FacilitiesManager {
	return &FacilitiesManager{
		exec:                         exec,
		inventory:                    newInventory(),
		deviceReadyNotificationQueue: make(map[types.DeviceIdentifier]bool),
	}
}

// CloseManager is invoked when the exec is stopping
func (mgr *FacilitiesManager) CloseManager() {
	// TODO
	mgr.threadStop()
	mgr.isInitialized = false
}

func (mgr *FacilitiesManager) InitializeManager() error {
	// create inventory based on nodeMgr
	nm := mgr.exec.GetNodeManager()
	for _, devInfo := range nm.GetDeviceInfos() {
		devId := devInfo.GetDeviceIdentifier()
		switch devInfo.GetNodeType() {
		case types.NodeTypeDisk:
			mgr.inventory.disks[devId] = &types.DiskAttributes{}
		case types.NodeTypeTape:
			mgr.inventory.tapes[devId] = &types.TapeAttributes{}
		}

	}

	mgr.threadStart()
	mgr.isInitialized = true
	return nil
}

func (mgr *FacilitiesManager) IsInitialized() bool {
	return mgr.isInitialized
}

// ResetManager clears out any artifacts left over by a previous exec session,
// and prepares the console for normal operations
func (mgr *FacilitiesManager) ResetManager() error {
	// TODO

	mgr.threadStop()
	mgr.threadStart()
	mgr.isInitialized = true
	return nil
}

func (mgr *FacilitiesManager) AssignDiskDeviceToExec(deviceId types.DeviceIdentifier) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	diskAttr, ok := mgr.inventory.disks[deviceId]
	if !ok {
		msg := fmt.Sprintf("Device %v is not known", deviceId)
		log.Println(msg)
		mgr.exec.Stop(types.StopFacilitiesComplex)
		return fmt.Errorf(msg)
	}

	if diskAttr.AssignedTo != nil {
		msg := fmt.Sprintf("Device %v is already assigned to %v", deviceId, diskAttr.AssignedTo.RunId)
		log.Println(msg)
		mgr.exec.Stop(types.StopFacilitiesComplex)
		return fmt.Errorf(msg)
	}

	diskAttr.AssignedTo = mgr.exec.GetRunControlEntry()
	return nil
}

// GetDeviceStatusDetail generates a short string to be used as a suffix to the basic
// disk or tape status for FS and related keyin displays
func (mgr *FacilitiesManager) GetDeviceStatusDetail(deviceId types.DeviceIdentifier) string {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	str := ""
	if mgr.isInitialized {
		da, ok := mgr.inventory.disks[deviceId]
		if ok {
			//	[[*] [R|F] PACKID pack-id]
			if da.AssignedTo != nil {
				str += "* "
			} else {
				str += "  "
			}

			if da.PackAttrs != nil && da.PackAttrs.IsPrepped {
				if da.PackAttrs.IsFixed {
					str += "F "
				} else {
					str += "R "
				}

				str += "PACKID " + da.PackAttrs.PackName
			}
		}

		// ta, ok := mgr.inventory.tapes[deviceId]
		// if ok {
		//	if ta.AssignedTo != nil {
		//		//	[* RUNID run-id REEL reel [RING|NORING] [POS [*]ffff[+|-][*]bbbbbb | POS LOST]]
		//		str += "* RUNID " + ta.AssignedTo.RunId + " REEL " + ta.reelNumber
		//		// TODO RING | NORING
		//		// TODO POS
		//	}
		// }
	}

	return str
}

func (mgr *FacilitiesManager) GetDiskAttributes(deviceId types.DeviceIdentifier) (*types.DiskAttributes, error) {
	attr, ok := mgr.inventory.disks[deviceId]
	if ok {
		return attr, nil
	} else {
		return nil, fmt.Errorf("not found")
	}
}

func (mgr *FacilitiesManager) IsDeviceAssigned(deviceId types.DeviceIdentifier) bool {
	dAttr, ok := mgr.inventory.disks[deviceId]
	if ok {
		return dAttr.AssignedTo != nil
	}

	tAttr, ok := mgr.inventory.tapes[deviceId]
	if ok {
		return tAttr.AssignedTo != nil
	}

	return false
}

func (mgr *FacilitiesManager) NotifyDeviceReady(deviceInfo types.DeviceInfo, isReady bool) {
	// queue this for the thread to pick up
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.deviceReadyNotificationQueue[deviceInfo.GetDeviceIdentifier()] = isReady
}

func (mgr *FacilitiesManager) diskBecameReady(deviceId types.DeviceIdentifier) {
	// Device became ready - any pack attributes we have, are obsolete, so reload them
	log.Printf("FacMgr:Disk %v became ready", deviceId)

	nm := mgr.exec.GetNodeManager()
	ni, err := nm.GetNodeInfoByIdentifier(types.NodeIdentifier(deviceId))
	if err != nil {
		mgr.exec.Stop(types.StopFacilitiesComplex)
	}

	diskAttr := mgr.inventory.disks[deviceId]
	diskAttr.PackAttrs = nil

	// we only care if the unit is UP, SU, or RV (i.e., not DN)
	devStat := ni.GetNodeStatus()
	if devStat != types.NodeStatusDown {
		packAttr := &types.PackAttributes{}
		packAttr.Label = make([]pkg.Word36, 28)
		ioPkt := nodeMgr.NewDiskIoPacketReadLabel(deviceId, packAttr.Label)
		nm.RouteIo(ioPkt)
		ioStat := ioPkt.GetIoStatus()
		if ioStat == types.IosInternalError {
			return
		} else if ioStat != types.IosComplete {
			log.Printf("FacMgr:IO Error reading label disk:%v status:%v", deviceId, ioStat)
			consMsg := fmt.Sprintf("%v IO ERROR Reading Pack Label - Status=%v", ni.GetNodeName(), ioStat)
			mgr.exec.SendExecReadOnlyMessage(consMsg)
			// if unit is UP or SU, tell node manager to DN the unit
			if devStat == types.NodeStatusUp || devStat == types.NodeStatusSuspended {
				_ = nm.SetNodeStatus(types.NodeIdentifier(deviceId), types.NodeStatusDown)
			}
			return
		}

		mgr.mutex.Lock()
		diskAttr := mgr.inventory.disks[deviceId]
		diskAttr.PackAttrs = packAttr
		mgr.mutex.Unlock()

		if packAttr.Label[0].ToStringAsAscii() != "VOL1" {
			consMsg := fmt.Sprintf("%v Pack has no VOL1 label", ni.GetNodeName())
			mgr.exec.SendExecReadOnlyMessage(consMsg)
			// if unit is UP or SU, tell node manager to DN the unit
			if devStat == types.NodeStatusUp || devStat == types.NodeStatusSuspended {
				_ = nm.SetNodeStatus(types.NodeIdentifier(deviceId), types.NodeStatusDown)
			}
			return
		}

		packName := (packAttr.Label[1].ToStringAsAscii() + packAttr.Label[2].ToStringAsAscii())[0:6]
		packAttr.PackName = packName
		if !nodeMgr.IsValidPackName(packName) {
			consMsg := fmt.Sprintf("%v Invalid pack ID in VOL1 label", ni.GetNodeName())
			mgr.exec.SendExecReadOnlyMessage(consMsg)
			// if unit is UP or SU, tell node manager to DN the unit
			if devStat == types.NodeStatusUp || devStat == types.NodeStatusSuspended {
				_ = nm.SetNodeStatus(types.NodeIdentifier(deviceId), types.NodeStatusDown)
			}
			return
		}

		packAttr.IsPrepped = true
	}
}

func (mgr *FacilitiesManager) tapeBecameReady(deviceId types.DeviceIdentifier) {
	// Device became ready
	// what we do here depends upon the current state of the device...
}

func (mgr *FacilitiesManager) thread() {
	mgr.threadStarted = true

	for !mgr.terminateThread {
		time.Sleep(10 * time.Millisecond)

		// any device ready notifications?
		mgr.mutex.Lock()
		queue := mgr.deviceReadyNotificationQueue
		mgr.deviceReadyNotificationQueue = make(map[types.DeviceIdentifier]bool)
		mgr.mutex.Unlock()

		for devId, flag := range queue {
			if flag {
				_, ok := mgr.inventory.disks[devId]
				if ok {
					go mgr.diskBecameReady(devId)
					continue
				}

				_, ok = mgr.inventory.tapes[devId]
				if ok {
					go mgr.tapeBecameReady(devId)
				}
			}
		}
	}

	mgr.threadStopped = true
}

func (mgr *FacilitiesManager) threadStart() {
	mgr.terminateThread = false
	if !mgr.threadStarted {
		go mgr.thread()
		for !mgr.threadStarted {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *FacilitiesManager) threadStop() {
	if mgr.threadStarted {
		mgr.terminateThread = true
		for !mgr.threadStopped {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *FacilitiesManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vFacilitiesManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  initialized:     %v\n", indent, mgr.isInitialized)
	_, _ = fmt.Fprintf(dest, "%v  threadStarted:   %v\n", indent, mgr.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:   %v\n", indent, mgr.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, mgr.terminateThread)

	nm := mgr.exec.GetNodeManager()
	_, _ = fmt.Fprintf(dest, "%v  Disk units:\n", indent)
	for deviceId, diskAttr := range mgr.inventory.disks {
		nodeInfo, _ := nm.GetNodeInfoByIdentifier(types.NodeIdentifier(deviceId))
		str := nodeInfo.GetNodeName()
		if diskAttr.AssignedTo != nil {
			str += "* " + diskAttr.AssignedTo.RunId
		}
		if diskAttr.PackAttrs != nil {
			packAttr := diskAttr.PackAttrs
			str += fmt.Sprintf(" PACK-ID:%v Prepped:%v Fixed:%v",
				packAttr.PackName, packAttr.IsPrepped, packAttr.IsFixed)
		}

		_, _ = fmt.Fprintf(dest, "%s    %s\n", indent, str)
	}

	_, _ = fmt.Fprintf(dest, "%v  Tape units:\n", indent)
	for deviceId, tapeAttr := range mgr.inventory.tapes {
		nodeInfo, _ := nm.GetNodeInfoByIdentifier(types.NodeIdentifier(deviceId))
		str := nodeInfo.GetNodeName()
		if tapeAttr.AssignedTo != nil {
			str += "* " + tapeAttr.AssignedTo.RunId
		}
		if tapeAttr.ReelAttrs != nil {
			str += fmt.Sprintf(" REEL-ID:%v Labeled:%v",
				tapeAttr.ReelAttrs.ReelNumber,
				tapeAttr.ReelAttrs.IsLabeled)
		}

		_, _ = fmt.Fprintf(dest, "%s    %s\n", indent, str)
	}

	_, _ = fmt.Fprintf(dest, "%v  Queued device-ready notifications:\n", indent)
	for devId, ready := range mgr.deviceReadyNotificationQueue {
		wId := pkg.Word36(devId)
		_, _ = fmt.Fprintf(dest, "%v    devId:0%v ready:%v\n", indent, wId.ToStringAsOctal(), ready)
	}
}
