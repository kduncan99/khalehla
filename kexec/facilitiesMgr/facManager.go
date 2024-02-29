// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

// we can know about nodeMgr, but nodeMgr cannot know about us
import (
	"fmt"
	"io"
	"khalehla/kexec"
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/nodes"
	"khalehla/pkg"
	"log"
	"sync"
	"time"
)

type FacilitiesManager struct {
	exec                         kexec.IExec
	mutex                        sync.Mutex
	threadDone                   bool
	inventory                    *kexec.inventory
	deviceReadyNotificationQueue map[kexec.NodeIdentifier]bool
}

func NewFacilitiesManager(exec kexec.IExec) *FacilitiesManager {
	return &FacilitiesManager{
		exec:                         exec,
		inventory:                    kexec.newInventory(),
		deviceReadyNotificationQueue: make(map[kexec.NodeIdentifier]bool),
	}
}

// Boot is invoked when the exec is booting
func (mgr *FacilitiesManager) Boot() error {
	log.Printf("FacMgr:Boot")

	// clear device ready notifications
	mgr.deviceReadyNotificationQueue = make(map[kexec.NodeIdentifier]bool)

	// (re)build inventory based on nodeMgr
	// this implies that nodeMgr.Boot() MUST be invoked before invoking us.
	// TODO at some point, it might be nice to add the channels in here
	// TODO should we really do this? don't we want to preserve the inventory for the previous session?
	mgr.inventory = kexec.newInventory()
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

func (mgr *FacilitiesManager) AssignDiskDeviceToExec(deviceId kexec.NodeIdentifier) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	diskAttr, ok := mgr.inventory.disks[deviceId]
	if !ok {
		msg := fmt.Sprintf("Device %v is not known", deviceId)
		log.Println(msg)
		mgr.exec.Stop(kexec.StopFacilitiesComplex)
		return fmt.Errorf(msg)
	}

	if diskAttr.AssignedTo != nil {
		msg := fmt.Sprintf("Device %v is already assigned to %v", deviceId, diskAttr.AssignedTo.RunId)
		log.Println(msg)
		mgr.exec.Stop(kexec.StopFacilitiesComplex)
		return fmt.Errorf(msg)
	}

	// TODO Need to update the Exec RCE fac item table, once we have fac item tables

	diskAttr.AssignedTo = mgr.exec.GetRunControlEntry()
	return nil
}

func (mgr *FacilitiesManager) GetDiskAttributes(nodeId kexec.NodeIdentifier) (*kexec.DiskAttributes, bool) {
	attr, ok := mgr.inventory.nodes[nodeId]
	return attr.(*kexec.DiskAttributes), ok
}

func (mgr *FacilitiesManager) GetNodeAttributes(nodeId kexec.NodeIdentifier) (kexec.NodeAttributes, bool) {
	attr, ok := mgr.inventory.nodes[nodeId]
	return attr, ok
}

func (mgr *FacilitiesManager) GetNodeAttributesByName(name string) (kexec.NodeAttributes, bool) {
	for _, nodeAttr := range mgr.inventory.nodes {
		if nodeAttr.GetNodeName() == name {
			return nodeAttr, true
		}
	}

	return nil, false
}

func (mgr *FacilitiesManager) GetNodeStatusString(nodeId kexec.NodeIdentifier) string {
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
	case kexec.FacNodeStatusDown:
		str = "DN" + accStr
	case kexec.FacNodeStatusReserved:
		str = "RV" + accStr
	case kexec.FacNodeStatusSuspended:
		str = "SU" + accStr
	case kexec.FacNodeStatusUp:
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

func (mgr *FacilitiesManager) IsDeviceAssigned(deviceId kexec.NodeIdentifier) bool {
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

func (mgr *FacilitiesManager) NotifyDeviceReady(deviceId kexec.NodeIdentifier, isReady bool) {
	// queue this for the thread to pick up
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.deviceReadyNotificationQueue[deviceId] = isReady
}

func (mgr *FacilitiesManager) SetNodeStatus(nodeId kexec.NodeIdentifier, status kexec.FacNodeStatus) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	nodeAttr, ok := mgr.inventory.nodes[nodeId]
	if !ok {
		return fmt.Errorf("node not found")
	}

	// for now, we do not allow changing status of anything except devices
	if nodeAttr.GetNodeCategoryType() != nodes.NodeCategoryDevice {
		return fmt.Errorf("not allowed")
	}

	stopExec := false
	nodeManager := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)

	switch status {
	case kexec.FacNodeStatusDown:
		if nodeAttr.GetNodeCategoryType() == nodes.NodeCategoryDevice {
			if nodeAttr.GetNodeDeviceType() == nodes.NodeDeviceDisk {
				nodeAttr.SetFacNodeStatus(status)
				stopExec = mgr.IsDeviceAssigned(nodeId)
				break
			} else if nodeAttr.GetNodeDeviceType() == nodes.NodeDeviceDisk {
				// reset the tape device (unmounts it as part of the process)
				ioPkt := nodes.NewTapeIoPacketReset(nodeId)
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

	case kexec.FacNodeStatusReserved:
		if nodeAttr.GetNodeCategoryType() == nodes.NodeCategoryDevice {
			nodeAttr.SetFacNodeStatus(status)
		} else {
			// anything other than disk or tape
			return fmt.Errorf("not allowed")
		}

	case kexec.FacNodeStatusSuspended:
		if nodeAttr.GetNodeCategoryType() == nodes.NodeCategoryDevice &&
			nodeAttr.GetNodeDeviceType() == nodes.NodeDeviceDisk {
			nodeAttr.SetFacNodeStatus(status)
		} else {
			// anything other than disk
			return fmt.Errorf("not allowed")
		}

	case kexec.FacNodeStatusUp:
		if nodeAttr.GetNodeCategoryType() == nodes.NodeCategoryDevice {
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
		mgr.exec.Stop(kexec.StopConsoleResponseRequiresReboot)
	}

	return nil
}
func (mgr *FacilitiesManager) diskBecameReady(nodeId kexec.NodeIdentifier) {
	// Device became ready - any pack attributes we have, are obsolete, so reload them
	log.Printf("FacMgr:Disk %v became ready", nodeId)

	mgr.mutex.Lock()

	diskAttr := mgr.inventory.disks[nodeId]
	diskAttr.PackLabelInfo = nil

	// we only care if the unit is UP, SU, or RV (i.e., not DN)
	facStat := diskAttr.GetFacNodeStatus()
	if facStat != kexec.FacNodeStatusDown {
		label := make([]pkg.Word36, 28)
		nodeManager := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
		ioPkt := nodes.NewDiskIoPacketReadLabel(nodeId, label)
		nodeManager.RouteIo(ioPkt)
		ioStat := ioPkt.GetIoStatus()
		if ioStat == nodes.IosInternalError {
			mgr.mutex.Unlock()
			return
		} else if ioStat != nodes.IosComplete {
			mgr.mutex.Unlock()

			log.Printf("FacMgr:IO Error reading label disk:%v %v", nodeId, ioPkt.GetString())
			consMsg := fmt.Sprintf("%v IO ERROR Reading Pack Label - Status=%v", diskAttr.GetNodeName(), ioStat)
			mgr.exec.SendExecReadOnlyMessage(consMsg, nil)

			// if unit is UP or SU, make it DN
			if facStat == kexec.FacNodeStatusUp || facStat == kexec.FacNodeStatusSuspended {
				_ = mgr.SetNodeStatus(nodeId, kexec.FacNodeStatusDown)
			}
			return
		}

		var ok bool
		diskAttr.PackLabelInfo, ok = kexec.NewPackLabelInfo(label)
		if !ok {
			mgr.mutex.Unlock()

			consMsg := fmt.Sprintf("%v Pack has no VOL1 label", diskAttr.GetNodeName())
			mgr.exec.SendExecReadOnlyMessage(consMsg, nil)

			// if unit is UP or SU, tell node manager to DN the unit
			if facStat == kexec.FacNodeStatusUp || facStat == kexec.FacNodeStatusSuspended {
				_ = mgr.SetNodeStatus(nodeId, kexec.FacNodeStatusDown)
			}
			return
		}

		diskAttr.IsPrepped = true
	}

	mgr.mutex.Unlock()
}

func (mgr *FacilitiesManager) tapeBecameReady(nodeId kexec.NodeIdentifier) {
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
		mgr.deviceReadyNotificationQueue = make(map[kexec.NodeIdentifier]bool)
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
