// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/hardware"
	"khalehla/kexec"
	"khalehla/kexec/nodeMgr"
	hardware2 "khalehla/old/hardware"
	kexec2 "khalehla/old/kexec"
	nodeMgr2 "khalehla/old/kexec/nodeMgr"
)

type inventory struct {
	nodes map[hardware.NodeIdentifier]kexec2.INodeAttributes
	disks map[hardware.NodeIdentifier]*kexec2.DiskAttributes
	tapes map[hardware.NodeIdentifier]*kexec2.TapeAttributes
}

func newInventory() *inventory {
	i := &inventory{
		nodes: make(map[hardware.NodeIdentifier]kexec2.INodeAttributes),
		disks: make(map[hardware.NodeIdentifier]*kexec2.DiskAttributes),
		tapes: make(map[hardware.NodeIdentifier]*kexec2.TapeAttributes),
	}
	return i
}

func (i *inventory) injectNode(nodeInfo nodeMgr.NodeInfo) {
	if nodeInfo.GetNodeCategoryType() == hardware.NodeCategoryDevice {
		devInfo := nodeInfo.(nodeMgr2.DeviceInfo)
		devId := devInfo.GetNodeIdentifier()
		switch devInfo.GetNodeDeviceType() {
		case hardware2.NodeDeviceDisk:
			attr := &kexec2.DiskAttributes{
				Identifier: devInfo.GetNodeIdentifier(),
				Name:       devInfo.GetNodeName(),
				Status:     kexec.FacNodeStatusUp,
			}
			i.nodes[devId] = attr
			i.disks[devId] = attr
		case hardware2.NodeDeviceTape:
			attr := &kexec2.TapeAttributes{
				Identifier: devInfo.GetNodeIdentifier(),
				Name:       devInfo.GetNodeName(),
				Status:     kexec.FacNodeStatusUp,
			}
			i.nodes[devId] = attr
			i.tapes[devId] = attr
		}
	}
}
