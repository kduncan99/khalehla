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

// TODO move this information somewhere more useful
//

/*
All facility items:
+00,W   internal file name - Fieldata LJSF
+01,W   (internal file name cont)
+02,W   file name - Fieldata LJSF
+03,W   (file name cont)
+04,W   qualifier - Fieldata LJSF
+05,W   (qualifier cont)
+06,S1  equipment code
         000 file has not been assigned (@USE exists, but @ASG has not been done)
         015 9-track tape
         016 virtual tape handler
         017 cartridge tape, DVD tape
         024 word-addressable mass storage
         036 sector-addressable mass storage
         077 arbitrary device
+07,S1  attributes
+07,b10:b35 @ASG options

Unit record and non-standard peripherals
+07,S1  attributes
         040 tape labeling is supported
         020 file is temporary
         010 internal name is a use name

Sector-formatted mass storage
+06,S2  file mode
         040 exclusively assigned
         020 read key needed
         010 write key needed
         004 file is write inhibited
         002 file is read inhibited
         001 word-addressable (always clear)
+06,S3  granularity
         zero -> track, nonzero -> position
+06,S4  relative file-cycle
+06,T3  absolute file-cycle
+07,S1  attributes
         020 file is temporary
         010 internal name is a use name
         004 shared file
         002 large file
+010,H1 initial granule count (initial reserve)
+010,H2 max granule count
+011,H1 highest track referenced
+011,H2 highest granule assigned
+012,S4 total pack count if removable (63 -> 63 or greater)
+012,S5 equipment code - same as +06,S1
+012,S6 subcode - zero

Magnetic tape peripherals
+06,S2  file mode
         040 exclusively assigned
         020 read key needed
         010 write key needed
         004 file is write inhibited
         002 file is read inhibited
+06,S3  unit count (I presume, the docs are not helpful)
		number of units assigned (0?, 1, 2)
+07,S1  attributes
         040 tape labeling is supported
         020 file is temporary
         010 internal name is a use name
         004 file is a shared file
+010,S1 total reel count
+010,S2 logical channel
+010,S3 noise constant
+012,T1 expiration period
+012,S3 reel index
+012,S4 files extended
+012,T3 blocks extended
+013,W  current reel number
+014,W  next reel number

Word addressable
+06,S2  file mode
         040 exclusively assigned
         020 read key needed
         010 write key needed
         004 write inhibited
         002 read inhibited
         001 word-addressable (always set)
+06,S3  granularity
         zero -> track, nonzero -> position
+06,S4  relative file-cycle
+06,T3  absolute file-cycle
+07,S1  attributes
         020 file is temporary
         010 internal name is a use name
         004 shared file
+010,W  length of file in words
+011,W  maximum file length in words
+012,S4 total pack count if removable (63 -> 63 or greater)
+012,S5 equipment code - same as +06,S1
+012,S6 subcode - zero
*/

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
	threadDone                   bool
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

// Boot is invoked when the exec is booting
func (mgr *FacilitiesManager) Boot() error {
	log.Printf("FacMgr:Boot")

	// clear device ready notifications
	mgr.deviceReadyNotificationQueue = make(map[types.DeviceIdentifier]bool)

	// (re)build inventory based on nodeMgr
	// this implies that nodeMgr.Boot() MUST be invoked before invoking us.
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

	// TODO Need to update the Exec RCE fac item table, once we have fac item tables

	diskAttr.AssignedTo = mgr.exec.GetRunControlEntry()
	return nil
}

// GetDeviceStatusDetail generates a short string to be used as a suffix to the basic
// disk or tape status for FS and related keyin displays
func (mgr *FacilitiesManager) GetDeviceStatusDetail(deviceId types.DeviceIdentifier) string {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	str := ""
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
			mgr.exec.SendExecReadOnlyMessage(consMsg, nil)
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
			mgr.exec.SendExecReadOnlyMessage(consMsg, nil)
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
			mgr.exec.SendExecReadOnlyMessage(consMsg, nil)
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
	mgr.threadDone = false

	for !mgr.exec.GetStopFlag() {
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

	mgr.threadDone = true
}

func (mgr *FacilitiesManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vFacilitiesManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  threadDone: %v\n", indent, mgr.threadDone)

	nm := mgr.exec.GetNodeManager()
	_, _ = fmt.Fprintf(dest, "%v  Disk units:\n", indent)
	for deviceId, diskAttr := range mgr.inventory.disks {
		nodeInfo, _ := nm.GetNodeInfoByIdentifier(types.NodeIdentifier(deviceId))
		str := nodeInfo.GetNodeName()
		if diskAttr.AssignedTo != nil {
			str += " * " + diskAttr.AssignedTo.RunId
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
