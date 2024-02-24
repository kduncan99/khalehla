// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/types"
)

type inventory struct {
	nodes map[types.NodeIdentifier]NodeAttributes
	disks map[types.NodeIdentifier]*DiskAttributes
	tapes map[types.NodeIdentifier]*TapeAttributes
}

func newInventory() *inventory {
	i := &inventory{
		nodes: make(map[types.NodeIdentifier]NodeAttributes),
		disks: make(map[types.NodeIdentifier]*DiskAttributes),
		tapes: make(map[types.NodeIdentifier]*TapeAttributes),
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
