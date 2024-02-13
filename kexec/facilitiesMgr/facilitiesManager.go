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

type diskAttributes struct {
	assignedTo     *types.RunControlEntry
	packAttributes *packAttributes
}

type packAttributes struct {
	label     []pkg.Word36
	isPrepped bool
	isFixed   bool
	packName  string
}

type tapeAttributes struct {
	assignedTo     *types.RunControlEntry
	reelAttributes *reelAttributes
}

type reelAttributes struct {
	reelNumber string
	isLabeled  bool
}

type inventory struct {
	disks map[types.DeviceIdentifier]diskAttributes
	tapes map[types.DeviceIdentifier]tapeAttributes
}

func newInventory() *inventory {
	return &inventory{
		disks: make(map[types.DeviceIdentifier]diskAttributes),
		tapes: make(map[types.DeviceIdentifier]tapeAttributes),
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
	nm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
	for _, devInfo := range nm.GetDeviceInfos() {
		switch devInfo.GetNodeType() {
		case types.NodeTypeDisk:
			mgr.inventory.disks[devInfo.GetDeviceIdentifier()] = diskAttributes{}
		case types.NodeTypeTape:
			mgr.inventory.tapes[devInfo.GetDeviceIdentifier()] = tapeAttributes{}
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

	if diskAttr.assignedTo != nil {
		msg := fmt.Sprintf("Device %v is already assigned to %v", deviceId, diskAttr.assignedTo.RunId)
		log.Println(msg)
		mgr.exec.Stop(types.StopFacilitiesComplex)
		return fmt.Errorf(msg)
	}

	diskAttr.assignedTo = mgr.exec.GetRunControlEntry()
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
			if da.assignedTo != nil {
				str += "* "
			} else {
				str += "  "
			}

			if da.packAttributes != nil && da.packAttributes.isPrepped {
				if da.packAttributes.isFixed {
					str += "F "
				} else {
					str += "R "
				}

				str += "PACKID " + da.packAttributes.packName
			}
		}

		//ta, ok := mgr.inventory.tapes[deviceId]
		//if ok {
		//	if ta.assignedTo != nil {
		//		//	[* RUNID run-id REEL reel [RING|NORING] [POS [*]ffff[+|-][*]bbbbbb | POS LOST]]
		//		str += "* RUNID " + ta.assignedTo.RunId + " REEL " + ta.reelNumber
		//		// TODO RING | NORING
		//		// TODO POS
		//	}
		//}
	}

	return str
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

	dm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
	ni, err := dm.GetNodeInfoByIdentifier(types.NodeIdentifier(deviceId))
	if err != nil {
		mgr.exec.Stop(types.StopFacilitiesComplex)
	}

	// we only care if the unit is UP, SU, or RV (i.e., not DN)
	if ni.GetNodeStatus() == types.NodeStatusDown {
		return
	}

	label := make([]pkg.Word36, 28)
	ioPkt := nodeMgr.NewDiskIoPacketReadLabel(deviceId, label)
	dm.RouteIo(ioPkt)
	ioStat := ioPkt.GetIoStatus()
	if ioStat == types.IosInternalError {
		return
	} else if ioStat != types.IosComplete {
		log.Printf("FacMgr:IO Error reading label disk:%v status:%v", deviceId, ioStat)
		consMsg := fmt.Sprintf("%v IO ERROR Reading Pack Label - Status=%v", ni.GetNodeName(), ioStat)
		mgr.exec.SendExecReadOnlyMessage(consMsg)
		// TODO tell node manager to DN the unit
		return
	}

	if label[0].ToStringAsAscii() != "VOL1" {
		consMsg := fmt.Sprintf("%v Pack has no VOL1 label", ni.GetNodeName())
		mgr.exec.SendExecReadOnlyMessage(consMsg)
		// TODO tell node manager to DN the unit
		return
	}

	packName := (label[1].ToStringAsAscii() + label[2].ToStringAsAscii())[0:6]
	if !nodeMgr.IsValidPackName(packName) {
		consMsg := fmt.Sprintf("%v Invalid pack ID in VOL1 label", ni.GetNodeName())
		mgr.exec.SendExecReadOnlyMessage(consMsg)
		// TODO tell node manager to DN the unit
		return
	}

	// TODO much more to do here, but we need to stop doing MBTs first...
}

func (mgr *FacilitiesManager) tapeBecameReady(deviceId types.DeviceIdentifier) {
	// Device became ready
	// what we do here depends upon the current state of the device...
}

func (mgr *FacilitiesManager) thread() {
	mgr.threadStarted = true

	for !mgr.terminateThread {
		time.Sleep(10 * time.Millisecond)

		mgr.mutex.Lock()
		// any device ready notifications?
		for devId, flag := range mgr.deviceReadyNotificationQueue {
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
		mgr.mutex.Unlock()
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

	// TODO

}
