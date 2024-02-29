// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/nodeMgr"
)

type inventory struct {
	nodes map[kexec.NodeIdentifier]kexec.NodeAttributes
	disks map[kexec.NodeIdentifier]*kexec.DiskAttributes
	tapes map[kexec.NodeIdentifier]*kexec.TapeAttributes
}

func newInventory() *inventory {
	i := &inventory{
		nodes: make(map[kexec.NodeIdentifier]kexec.NodeAttributes),
		disks: make(map[kexec.NodeIdentifier]*kexec.DiskAttributes),
		tapes: make(map[kexec.NodeIdentifier]*kexec.TapeAttributes),
	}
	return i
}

func (i *inventory) injectNode(nodeInfo nodeMgr.NodeInfo) {
	if nodeInfo.GetNodeCategoryType() == kexec.NodeCategoryDevice {
		devInfo := nodeInfo.(nodeMgr.DeviceInfo)
		devId := devInfo.GetNodeIdentifier()
		switch devInfo.GetNodeDeviceType() {
		case kexec.NodeDeviceDisk:
			attr := &kexec.DiskAttributes{
				Name:   devInfo.GetNodeName(),
				Status: kexec.FacNodeStatusUp,
			}
			i.nodes[devId] = attr
			i.disks[devId] = attr
		case kexec.NodeDeviceTape:
			attr := &kexec.TapeAttributes{
				Name:   devInfo.GetNodeName(),
				Status: kexec.FacNodeStatusUp,
			}
			i.nodes[devId] = attr
			i.tapes[devId] = attr
		}
	}
}
