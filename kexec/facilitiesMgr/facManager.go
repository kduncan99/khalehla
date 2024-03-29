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
	"khalehla/klog"
	"khalehla/pkg"
	"sync"
	"time"
)

// FacilitiesManager does all the facilities-specific work for the exec.
// This file contains the APIs described in the IFacilitiesManager interface.
// All other externalized APIs are implemented in facServices, and internal functions are in facCore.
type FacilitiesManager struct {
	exec                         kexec.IExec
	mutex                        sync.Mutex
	threadDone                   bool
	inventory                    *inventory
	deviceReadyNotificationQueue map[hardware.NodeIdentifier]bool
	attachments                  map[hardware.NodeIdentifier]*kexec.RunControlEntry
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
	klog.LogTrace("FacMgr", "Boot")

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
	klog.LogTrace("FacMgr", "Close")
	// nothing to do
}

// Initialize is invoked when the application starts
func (mgr *FacilitiesManager) Initialize() error {
	klog.LogTrace("FacMgr", "Initialize")
	// nothing to do here
	return nil
}

func (mgr *FacilitiesManager) Stop() {
	klog.LogTrace("FacMgr", "Stop")
	for !mgr.threadDone {
		time.Sleep(25 * time.Millisecond)
	}
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

	return mgr.getNodeStatusString(nodeId)
}

func (mgr *FacilitiesManager) IsDeviceAssigned(deviceId hardware.NodeIdentifier) bool {
	_, ok := mgr.attachments[deviceId]
	return ok
}

func (mgr *FacilitiesManager) NotifyDeviceReady(nodeIdentifier hardware.NodeIdentifier, isReady bool) {
	// queue this for the thread to pick up
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.deviceReadyNotificationQueue[nodeIdentifier] = isReady
}

func (mgr *FacilitiesManager) SetNodeStatus(nodeId hardware.NodeIdentifier, status kexec.FacNodeStatus) bool {
	klog.LogTraceF("FacMgr", "SetNodeStatus nodeId=%v status=%v", nodeId, status)
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	nodeAttr, ok := mgr.inventory.nodes[nodeId]
	if !ok {
		klog.LogErrorF("FacMgr", "SetNodeStatus node %v not found", nodeId)
		return false
	}

	// for now, we do not allow changing status of anything except devices
	if nodeAttr.GetNodeCategoryType() != hardware.NodeCategoryDevice {
		klog.LogErrorF("FacMgr", "SetNodeStatus node %v not a device", nodeId)
		return false
	}

	stopExec := false
	nodeManager := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)

	switch status {
	case kexec.FacNodeStatusDown:
		if nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceDisk {
			nodeAttr.SetFacNodeStatus(status)
			stopExec = mgr.IsDeviceAssigned(nodeId)
			break
		} else if nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceTape {
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
		} else {
			// anything else
			klog.LogErrorF("FacMgr", "SetNodeStatus node %v DN not allowed, not a disk or tape device", nodeId)
			return false
		}

	case kexec.FacNodeStatusReserved:
		if nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceDisk {
			// TODO RV disk
			nodeAttr.SetFacNodeStatus(status)
		} else if nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceTape {
			// TODO RV tape
		} else {
			// anything other than disk or tape
			klog.LogErrorF("FacMgr", "SetNodeStatus node %v RV not allowed, not a disk or tape device", nodeId)
			return false
		}

	case kexec.FacNodeStatusSuspended:
		if nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceDisk {
			// TODO check for SU of last UP disk
			nodeAttr.SetFacNodeStatus(status)
		} else {
			// anything other than disk
			klog.LogErrorF("FacMgr", "SetNodeStatus node %v SU not allowed, not a disk device", nodeId)
			return false
		}

	case kexec.FacNodeStatusUp:
		if nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceDisk {
			// TODO UP disk
			nodeAttr.SetFacNodeStatus(status)
		} else if nodeAttr.GetNodeDeviceType() == hardware.NodeDeviceTape {
			// TODO UP tape
		} else {
			// anything other than disk or tape
			klog.LogErrorF("FacMgr", "SetNodeStatus node %v UP not allowed, not a disk or tape device", nodeId)
			return false
		}

	default:
		klog.LogFatalF("FacMgr", "SetNodeStatus node %v status %v not recognized", nodeId, status)
		mgr.exec.Stop(kexec.StopFacilitiesComplex)
		return false
	}

	msg := nodeAttr.GetNodeName() + " " + mgr.getNodeStatusString(nodeId)
	mgr.exec.SendExecReadOnlyMessage(msg, nil)
	if stopExec {
		mgr.exec.Stop(kexec.StopConsoleResponseRequiresReboot)
	}

	return true
}

// diskBecameReady handles the notification which arrives after a unit attention.
// this waits on IO, so do NOT call it under lock.
func (mgr *FacilitiesManager) diskBecameReady(nodeId hardware.NodeIdentifier) {
	// Device became ready - any pack attributes we have are obsolete, so reload them
	mgr.mutex.Lock()

	diskAttr := mgr.inventory.disks[nodeId]
	diskAttr.PackLabelInfo = nil
	klog.LogInfoF("FacMgr", "Disk %v [%v] became ready", nodeId, diskAttr.GetNodeName())
	msg := fmt.Sprintf("%v Device Ready", diskAttr.GetNodeName())
	mgr.exec.SendExecReadOnlyMessage(msg, nil)

	// we only care if the unit is UP, SU, or RV (i.e., not DN)
	facStat := diskAttr.GetFacNodeStatus()
	if facStat != kexec.FacNodeStatusDown {
		nm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
		ni, _ := nm.GetNodeInfoByIdentifier(nodeId)
		ddi := ni.(*nodeMgr.DiskDeviceInfo)
		dev := ddi.GetDiskDevice()
		_, _, prepFactor, _ := dev.GetDiskGeometry()
		if prepFactor == 0 {
			consMsg := fmt.Sprintf("%v Pack is not prepped", diskAttr.GetNodeName())
			mgr.exec.SendExecReadOnlyMessage(consMsg, nil)
		} else {
			label := make([]pkg.Word36, prepFactor)
			cw := channels.ControlWord{
				Buffer:    label,
				Offset:    0,
				Length:    uint(prepFactor),
				Direction: channels.DirectionForward,
				Format:    channels.TransferPacked,
			}
			cp := &channels.ChannelProgram{
				NodeIdentifier: nodeId,
				IoFunction:     ioPackets.IofRead,
				BlockId:        0,
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

				klog.LogInfoF("FacMgr", "IO Error reading label disk:%v", cp.GetString())
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
	}

	mgr.mutex.Unlock()
}

func (mgr *FacilitiesManager) tapeBecameReady(nodeId hardware.NodeIdentifier) {
	// Device became ready
	// what we do here depends upon the current state of the device...
	mgr.mutex.Lock()

	diskAttr := mgr.inventory.disks[nodeId]
	diskAttr.PackLabelInfo = nil
	klog.LogInfoF("FacMgr", "Disk %v [%v] became ready", nodeId, diskAttr.GetNodeName())
	msg := fmt.Sprintf("%v Device Ready", diskAttr.GetNodeName())
	mgr.exec.SendExecReadOnlyMessage(msg, nil)

	// TODO implmement tapeBecameReady()

	mgr.mutex.Unlock()
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
