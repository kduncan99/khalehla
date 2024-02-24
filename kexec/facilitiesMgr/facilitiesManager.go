// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

// we can know about nodeMgr, but nodeMgr cannot know about us
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

type FacilitiesManager struct {
	exec                         types.IExec
	mutex                        sync.Mutex
	threadDone                   bool
	inventory                    *inventory
	deviceReadyNotificationQueue map[types.NodeIdentifier]bool
}

func NewFacilitiesManager(exec types.IExec) *FacilitiesManager {
	return &FacilitiesManager{
		exec:                         exec,
		inventory:                    newInventory(),
		deviceReadyNotificationQueue: make(map[types.NodeIdentifier]bool),
	}
}

// Boot is invoked when the exec is booting
func (mgr *FacilitiesManager) Boot() error {
	log.Printf("FacMgr:Boot")

	// clear device ready notifications
	mgr.deviceReadyNotificationQueue = make(map[types.NodeIdentifier]bool)

	// (re)build inventory based on nodeMgr
	// this implies that nodeMgr.Boot() MUST be invoked before invoking us.
	// TODO at some point, it might be nice to add the channels in here
	// TODO should we really do this? don't we want to preserve the inventory for the previous session?
	mgr.inventory = newInventory()
	nm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
	for _, devInfo := range nm.GetDeviceInfos() {
		mgr.inventory.injectNode(devInfo)
	}

	go mgr.thread()
	return nil
}

// Close is invoked when the application is terminating
func (mgr *FacilitiesManager) Close() {
	log.Printf("FacMgr:Close")
	// nothing to do
}

// Initialize is invoked when the application starts
func (mgr *FacilitiesManager) Initialize() error {
	log.Printf("FacMgr:Initialize")
	// nothing to do here
	return nil
}

func (mgr *FacilitiesManager) Stop() {
	log.Printf("FacMgr:Stop")
	for !mgr.threadDone {
		time.Sleep(25 * time.Millisecond)
	}
}

func (mgr *FacilitiesManager) AssignDiskDeviceToExec(deviceId types.NodeIdentifier) error {
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

	// TODO Need to update the Exec RCE fac item table, once we have fac item tables

	diskAttr.AssignedTo = mgr.exec.GetRunControlEntry()
	return nil
}

func (mgr *FacilitiesManager) GetDiskAttributes(nodeId types.NodeIdentifier) (*DiskAttributes, bool) {
	attr, ok := mgr.inventory.nodes[nodeId]
	return attr.(*DiskAttributes), ok
}

func (mgr *FacilitiesManager) GetNodeAttributes(nodeId types.NodeIdentifier) (NodeAttributes, bool) {
	attr, ok := mgr.inventory.nodes[nodeId]
	return attr, ok
}

func (mgr *FacilitiesManager) GetNodeAttributesByName(name string) (NodeAttributes, bool) {
	for _, nodeAttr := range mgr.inventory.nodes {
		if nodeAttr.GetNodeName() == name {
			return nodeAttr, true
		}
	}

	return nil, false
}

func (mgr *FacilitiesManager) GetNodeStatusString(nodeId types.NodeIdentifier) string {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	accStr := "   "
	nm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
	ni, _ := nm.GetNodeInfoByIdentifier(nodeId)
	if !ni.IsAccessible() {
		accStr = " NA"
	}

	str := ""
	facStat := mgr.inventory.nodes[nodeId].GetFacNodeStatus()
	switch facStat {
	case FacNodeStatusDown:
		str = "DN" + accStr
	case FacNodeStatusReserved:
		str = "RV" + accStr
	case FacNodeStatusSuspended:
		str = "SU" + accStr
	case FacNodeStatusUp:
		str = "UP" + accStr
	}

	da, ok := mgr.inventory.disks[nodeId]
	if ok {
		//	[[*] [R|F] PACKID pack-id]
		if da.AssignedTo != nil {
			str += " * "
		} else {
			str += "   "
		}

		if da.PackLabelInfo != nil {
			if da.IsFixed {
				str += "F "
			} else {
				str += "R "
			}

			str += "PACKID " + da.PackLabelInfo.PackId
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

	return str
}

func (mgr *FacilitiesManager) IsDeviceAssigned(deviceId types.NodeIdentifier) bool {
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

func (mgr *FacilitiesManager) NotifyDeviceReady(deviceId types.NodeIdentifier, isReady bool) {
	// queue this for the thread to pick up
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.deviceReadyNotificationQueue[deviceId] = isReady
}

func (mgr *FacilitiesManager) SetNodeStatus(nodeId types.NodeIdentifier, status FacNodeStatus) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	nodeAttr, ok := mgr.inventory.nodes[nodeId]
	if !ok {
		return fmt.Errorf("node not found")
	}

	// for now, we do not allow changing status of anything except devices
	if nodeAttr.GetNodeCategoryType() != nodeMgr.NodeCategoryDevice {
		return fmt.Errorf("not allowed")
	}

	stopExec := false
	nodeManager := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)

	switch status {
	case FacNodeStatusDown:
		if nodeAttr.GetNodeCategoryType() == nodeMgr.NodeCategoryDevice {
			if nodeAttr.GetNodeDeviceType() == nodeMgr.NodeDeviceDisk {
				nodeAttr.SetFacNodeStatus(status)
				stopExec = mgr.IsDeviceAssigned(nodeId)
				break
			} else if nodeAttr.GetNodeDeviceType() == nodeMgr.NodeDeviceDisk {
				// reset the tape device (unmounts it as part of the process)
				ioPkt := nodeMgr.NewTapeIoPacketReset(nodeId)
				nodeManager.RouteIo(ioPkt)
				nodeAttr.SetFacNodeStatus(status)
				if mgr.IsDeviceAssigned(nodeId) {
					// TODO - unassign the device from the run
					// TODO - tell Exec to abort the run to which the thing was assigned
				}
				break
			}
		}
		// anything else
		return fmt.Errorf("not allowed")

	case FacNodeStatusReserved:
		if nodeAttr.GetNodeCategoryType() == nodeMgr.NodeCategoryDevice {
			nodeAttr.SetFacNodeStatus(status)
		} else {
			// anything other than disk or tape
			return fmt.Errorf("not allowed")
		}

	case FacNodeStatusSuspended:
		if nodeAttr.GetNodeCategoryType() == nodeMgr.NodeCategoryDevice &&
			nodeAttr.GetNodeDeviceType() == nodeMgr.NodeDeviceDisk {
			nodeAttr.SetFacNodeStatus(status)
		} else {
			// anything other than disk
			return fmt.Errorf("not allowed")
		}

	case FacNodeStatusUp:
		if nodeAttr.GetNodeCategoryType() == nodeMgr.NodeCategoryDevice {
			nodeAttr.SetFacNodeStatus(status)
		} else {
			// anything other than disk or tape
			return fmt.Errorf("not allowed")
		}

	default:
		return fmt.Errorf("internal error")
	}

	msg := nodeAttr.GetNodeName() + " " + mgr.GetNodeStatusString(nodeId)
	mgr.exec.SendExecReadOnlyMessage(msg, nil)
	if stopExec {
		mgr.exec.Stop(types.StopConsoleResponseRequiresReboot)
	}

	return nil
}
func (mgr *FacilitiesManager) diskBecameReady(nodeId types.NodeIdentifier) {
	// Device became ready - any pack attributes we have, are obsolete, so reload them
	log.Printf("FacMgr:Disk %v became ready", nodeId)

	mgr.mutex.Lock()

	diskAttr := mgr.inventory.disks[nodeId]
	diskAttr.PackLabelInfo = nil

	// we only care if the unit is UP, SU, or RV (i.e., not DN)
	facStat := diskAttr.GetFacNodeStatus()
	if facStat != FacNodeStatusDown {
		label := make([]pkg.Word36, 28)
		nodeManager := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
		ioPkt := nodeMgr.NewDiskIoPacketReadLabel(nodeId, label)
		nodeManager.RouteIo(ioPkt)
		ioStat := ioPkt.GetIoStatus()
		if ioStat == types.IosInternalError {
			mgr.mutex.Unlock()
			return
		} else if ioStat != types.IosComplete {
			mgr.mutex.Unlock()

			log.Printf("FacMgr:IO Error reading label disk:%v %v", nodeId, ioPkt.GetString())
			consMsg := fmt.Sprintf("%v IO ERROR Reading Pack Label - Status=%v", diskAttr.GetNodeName(), ioStat)
			mgr.exec.SendExecReadOnlyMessage(consMsg, nil)

			// if unit is UP or SU, make it DN
			if facStat == FacNodeStatusUp || facStat == FacNodeStatusSuspended {
				_ = mgr.SetNodeStatus(nodeId, FacNodeStatusDown)
			}
			return
		}

		var ok bool
		diskAttr.PackLabelInfo, ok = types.NewPackLabelInfo(label)
		if !ok {
			mgr.mutex.Unlock()

			consMsg := fmt.Sprintf("%v Pack has no VOL1 label", diskAttr.GetNodeName())
			mgr.exec.SendExecReadOnlyMessage(consMsg, nil)

			// if unit is UP or SU, tell node manager to DN the unit
			if facStat == FacNodeStatusUp || facStat == FacNodeStatusSuspended {
				_ = mgr.SetNodeStatus(nodeId, FacNodeStatusDown)
			}
			return
		}

		diskAttr.IsPrepped = true
	}

	mgr.mutex.Unlock()
}

func (mgr *FacilitiesManager) tapeBecameReady(nodeId types.NodeIdentifier) {
	// Device became ready
	// what we do here depends upon the current state of the device...
	// TODO
}

func (mgr *FacilitiesManager) thread() {
	mgr.threadDone = false

	for !mgr.exec.GetStopFlag() {
		time.Sleep(10 * time.Millisecond)

		// any device ready notifications?
		mgr.mutex.Lock()
		queue := mgr.deviceReadyNotificationQueue
		mgr.deviceReadyNotificationQueue = make(map[types.NodeIdentifier]bool)
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

	mgr.threadDone = true
}

func (mgr *FacilitiesManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vFacilitiesManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  threadDone: %v\n", indent, mgr.threadDone)

	// TODO dump inventory

	_, _ = fmt.Fprintf(dest, "%v  Queued device-ready notifications:\n", indent)
	for devId, ready := range mgr.deviceReadyNotificationQueue {
		wId := pkg.Word36(devId)
		_, _ = fmt.Fprintf(dest, "%v    devId:0%v ready:%v\n", indent, wId.ToStringAsOctal(), ready)
	}
}
