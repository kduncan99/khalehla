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
	disks map[types.DeviceIdentifier]*DiskAttributes
	tapes map[types.DeviceIdentifier]*TapeAttributes
}

func newInventory() *inventory {
	return &inventory{
		disks: make(map[types.DeviceIdentifier]*DiskAttributes),
		tapes: make(map[types.DeviceIdentifier]*TapeAttributes),
	}
}

func (i *inventory) injectNode(nodeInfo nodeMgr.NodeInfo) {
	if nodeInfo.GetNodeCategoryType() == nodeMgr.NodeCategoryDevice {
		devInfo := nodeInfo.(nodeMgr.DeviceInfo)
		devId := devInfo.GetDeviceIdentifier()
		switch devInfo.GetNodeDeviceType() {
		case nodeMgr.NodeDeviceDisk:
			attr := &DiskAttributes{
				name:   devInfo.GetNodeName(),
				status: FacNodeStatusUp,
			}
			i.nodes[types.NodeIdentifier(devId)] = attr
			i.disks[devId] = attr
		case nodeMgr.NodeDeviceTape:
			attr := &TapeAttributes{
				name:   devInfo.GetNodeName(),
				status: FacNodeStatusUp,
			}
			i.nodes[types.NodeIdentifier(devId)] = attr
			i.tapes[devId] = attr
		}
	}
}
