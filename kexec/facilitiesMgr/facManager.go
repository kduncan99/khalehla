// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

// we can know about nodeMgr, but nodeMgr cannot know about us
import (
	"fmt"
	"io"
	"khalehla/hardware"
	"khalehla/hardware/channels"
	"khalehla/hardware/ioPackets"
	"khalehla/kexec"
	"khalehla/kexec/nodeMgr"
	"khalehla/pkg"
	"log"
	"sync"
	"time"
)

type FacilitiesManager struct {
	exec                         kexec.IExec
	mutex                        sync.Mutex
	threadDone                   bool
	inventory                    *inventory
	deviceReadyNotificationQueue map[hardware.NodeIdentifier]bool
}

func NewFacilitiesManager(exec kexec.IExec) *FacilitiesManager {
	return &FacilitiesManager{
		exec:                         exec,
		inventory:                    newInventory(),
		deviceReadyNotificationQueue: make(map[hardware.NodeIdentifier]bool),
	}
}

// Boot is invoked when the exec is booting
func (mgr *FacilitiesManager) Boot() error {
	log.Printf("FacMgr:Boot")

	// clear device ready notifications
	mgr.deviceReadyNotificationQueue = make(map[hardware.NodeIdentifier]bool)

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

func (mgr *FacilitiesManager) AssignDiskDeviceToExec(deviceId hardware.NodeIdentifier) error {
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

func (mgr *FacilitiesManager) GetDiskAttributes(nodeId hardware.NodeIdentifier) (*kexec.DiskAttributes, bool) {
	attr, ok := mgr.inventory.nodes[nodeId]
	return attr.(*kexec.DiskAttributes), ok
}

func (mgr *FacilitiesManager) GetNodeAttributes(nodeId hardware.NodeIdentifier) (kexec.NodeAttributes, bool) {
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

func (mgr *FacilitiesManager) GetNodeStatusString(nodeId hardware.NodeIdentifier) string {
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
	//		// TODO RING | NORING in FS display for tape unit
	//		// TODO POS in FS display for tape unit
	//	}
	// }

	return str
}

func (mgr *FacilitiesManager) IsDeviceAssigned(deviceId hardware.NodeIdentifier) bool {
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

func (mgr *FacilitiesManager) NotifyDeviceReady(deviceId hardware.NodeIdentifier, isReady bool) {
	// queue this for the thread to pick up
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.deviceReadyNotificationQueue[deviceId] = isReady
}

func (mgr *FacilitiesManager) SetNodeStatus(nodeId hardware.NodeIdentifier, status kexec.FacNodeStatus) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	nodeAttr, ok := mgr.inventory.nodes[nodeId]
	if !ok {
		return fmt.Errorf("node not found")
	}

	// for now, we do not allow changing status of anything except devices
	if nodeAttr.GetNodeCategoryType() != hardware.NodeCategoryDevice {
		return fmt.Errorf("not allowed")
	}

	stopExec := false
	nodeManager := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)

	switch status {
	case kexec.FacNodeStatusDown:
		if nodeAttr.GetNodeCategoryType() == hardware.NodeCategoryDevice {
			if nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceDisk {
				nodeAttr.SetFacNodeStatus(status)
				stopExec = mgr.IsDeviceAssigned(nodeId)
				break
			} else if nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceDisk {
				// Reset the tape device (unmounts it as part of the process)
				// We don't need to wait for IO to complete.
				cp := &channels.ChannelProgram{
					NodeIdentifier: nodeId,
					IoFunction:     ioPackets.IofReset,
				}
				nodeManager.RouteIo(cp)
				nodeAttr.SetFacNodeStatus(status)
				if mgr.IsDeviceAssigned(nodeId) {
					// TODO - un-assign the device from the run
					// TODO - tell Exec to abort the run to which the thing was assigned
				}
				break
			}
		}
		// anything else
		return fmt.Errorf("not allowed")

	case kexec.FacNodeStatusReserved:
		if nodeAttr.GetNodeCategoryType() == hardware.NodeCategoryDevice {
			nodeAttr.SetFacNodeStatus(status)
		} else {
			// anything other than disk or tape
			return fmt.Errorf("not allowed")
		}

	case kexec.FacNodeStatusSuspended:
		if nodeAttr.GetNodeCategoryType() == hardware.NodeCategoryDevice &&
			nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceDisk {
			nodeAttr.SetFacNodeStatus(status)
		} else {
			// anything other than disk
			return fmt.Errorf("not allowed")
		}

	case kexec.FacNodeStatusUp:
		if nodeAttr.GetNodeCategoryType() == hardware.NodeCategoryDevice {
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

// diskBecameReady handles the notification which arrives after a unit attention.
// this waits on IO, so do NOT call it under lock.
func (mgr *FacilitiesManager) diskBecameReady(nodeId hardware.NodeIdentifier) {
	// Device became ready - any pack attributes we have, are obsolete, so reload them
	log.Printf("FacMgr:Disk %v became ready", nodeId)

	mgr.mutex.Lock()

	diskAttr := mgr.inventory.disks[nodeId]
	diskAttr.PackLabelInfo = nil

	// we only care if the unit is UP, SU, or RV (i.e., not DN)
	facStat := diskAttr.GetFacNodeStatus()
	if facStat != kexec.FacNodeStatusDown {
		nm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
		ni, _ := nm.GetNodeInfoByIdentifier(nodeId)
		ddi := ni.(*nodeMgr.DiskDeviceInfo)
		dev := ddi.GetDiskDevice()
		blockSize, _, _ := dev.GetDiskGeometry()
		prepFactor := hardware.PrepFactorFromBlockSize[blockSize]

		label := make([]pkg.Word36, prepFactor)
		cw := channels.ControlWord{
			Buffer:    label,
			Offset:    0,
			Length:    0,
			Direction: channels.DirectionForward,
			Format:    channels.TransferPacked,
		}
		cp := &channels.ChannelProgram{
			NodeIdentifier: nodeId,
			IoFunction:     ioPackets.IofRead,
			BlockId:        2,
			ControlWords:   []channels.ControlWord{cw},
		}

		nm.RouteIo(cp)
		for cp.IoStatus == ioPackets.IosInProgress || cp.IoStatus == ioPackets.IosNotStarted {
			time.Sleep(10 * time.Millisecond)
		}

		if cp.IoStatus == ioPackets.IosInternalError {
			mgr.mutex.Unlock()
			return
		} else if cp.IoStatus != ioPackets.IosComplete {
			mgr.mutex.Unlock()

			log.Printf("FacMgr:IO Error reading label disk:%v", cp.GetString())
			consMsg := fmt.Sprintf("%v IO ERROR Reading Pack Label - Status=%v",
				diskAttr.GetNodeName(), ioPackets.IoStatusTable[cp.IoStatus])
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

func (mgr *FacilitiesManager) tapeBecameReady(nodeId hardware.NodeIdentifier) {
	// Device became ready
	// what we do here depends upon the current state of the device...
	// TODO implmement tapeBecameReady()
}

func (mgr *FacilitiesManager) thread() {
	mgr.threadDone = false

	for !mgr.exec.GetStopFlag() {
		time.Sleep(10 * time.Millisecond)

		// any device ready notifications?
		mgr.mutex.Lock()
		queue := mgr.deviceReadyNotificationQueue
		mgr.deviceReadyNotificationQueue = make(map[hardware.NodeIdentifier]bool)
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

	_, _ = fmt.Fprintf(dest, "%v  inventory:\n", indent)
	for _, nodeInfo := range mgr.inventory.nodes {
		_, _ = fmt.Fprintf(dest, "%v    %s id:%v stat:%v cat:%v type:%v\n",
			indent,
			nodeInfo.GetNodeName(),
			nodeInfo.GetNodeIdentifier(),
			nodeInfo.GetFacNodeStatus(),
			nodeInfo.GetNodeCategoryType(),
			nodeInfo.GetNodeDeviceType())
	}

	_, _ = fmt.Fprintf(dest, "%v  Queued device-ready notifications:\n", indent)
	for devId, ready := range mgr.deviceReadyNotificationQueue {
		wId := pkg.Word36(devId)
		_, _ = fmt.Fprintf(dest, "%v    devId:0%v ready:%v\n", indent, wId.ToStringAsOctal(), ready)
	}
}
