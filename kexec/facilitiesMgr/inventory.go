// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/nodeMgr"
)

type inventory struct {
	nodes map[kexec.NodeIdentifier]NodeAttributes
	disks map[kexec.NodeIdentifier]*DiskAttributes
	tapes map[kexec.NodeIdentifier]*TapeAttributes
}

func newInventory() *inventory {
	i := &inventory{
		nodes: make(map[kexec.NodeIdentifier]NodeAttributes),
		disks: make(map[kexec.NodeIdentifier]*DiskAttributes),
		tapes: make(map[kexec.NodeIdentifier]*TapeAttributes),
	}
	return i
}

func (i *inventory) injectNode(nodeInfo nodeMgr.NodeInfo) {
	if nodeInfo.GetNodeCategoryType() == nodeMgr.NodeCategoryDevice {
		devInfo := nodeInfo.(nodeMgr.DeviceInfo)
		devId := devInfo.GetNodeIdentifier()
		switch devInfo.GetNodeDeviceType() {
		case nodeMgr.NodeDeviceDisk:
			attr := &DiskAttributes{
				name:   devInfo.GetNodeName(),
				status: FacNodeStatusUp,
			}
			i.nodes[devId] = attr
			i.disks[devId] = attr
		case nodeMgr.NodeDeviceTape:
			attr := &TapeAttributes{
				name:   devInfo.GetNodeName(),
				status: FacNodeStatusUp,
			}
			i.nodes[devId] = attr
			i.tapes[devId] = attr
		}
	}
}
