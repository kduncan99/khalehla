// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/nodes"
)

type FacInventory struct {
	Nodes map[NodeIdentifier]NodeAttributes
	Disks map[NodeIdentifier]*DiskAttributes
	Tapes map[NodeIdentifier]*TapeAttributes
}

func NewFacInventory() *FacInventory {
	i := &FacInventory{
		Nodes: make(map[NodeIdentifier]NodeAttributes),
		Disks: make(map[NodeIdentifier]*DiskAttributes),
		Tapes: make(map[NodeIdentifier]*TapeAttributes),
	}
	return i
}

func (i *FacInventory) InjectNode(nodeInfo nodeMgr.NodeInfo) {
	if nodeInfo.GetNodeCategoryType() == nodes.NodeCategoryDevice {
		devInfo := nodeInfo.(nodeMgr.DeviceInfo)
		devId := devInfo.GetNodeIdentifier()
		switch devInfo.GetNodeDeviceType() {
		case nodes.NodeDeviceDisk:
			attr := &DiskAttributes{
				name:   devInfo.GetNodeName(),
				status: FacNodeStatusUp,
			}
			i.Nodes[devId] = attr
			i.Disks[devId] = attr
		case nodes.NodeDeviceTape:
			attr := &TapeAttributes{
				name:   devInfo.GetNodeName(),
				status: FacNodeStatusUp,
			}
			i.Nodes[devId] = attr
			i.Tapes[devId] = attr
		}
	}
}
