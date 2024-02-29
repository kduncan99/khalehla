// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/kexec/nodes"
)

type DiskAttributes struct {
	identifier    NodeIdentifier
	name          string
	status        FacNodeStatus
	AssignedTo    *RunControlEntry
	PackLabelInfo *PackLabelInfo
	IsPrepped     bool
	IsFixed       bool
}

func (da *DiskAttributes) GetFacNodeStatus() FacNodeStatus {
	return da.status
}

func (da *DiskAttributes) GetNodeCategoryType() nodes.NodeCategoryType {
	return nodes.NodeCategoryDevice
}

func (da *DiskAttributes) GetNodeDeviceType() nodes.NodeDeviceType {
	return nodes.NodeDeviceDisk
}

func (da *DiskAttributes) GetNodeIdentifier() NodeIdentifier {
	return da.identifier
}

func (da *DiskAttributes) GetNodeName() string {
	return da.name
}

func (da *DiskAttributes) SetFacNodeStatus(status FacNodeStatus) {
	da.status = status
}
