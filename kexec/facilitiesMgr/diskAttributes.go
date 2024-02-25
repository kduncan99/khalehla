// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/exec"
	"khalehla/kexec/nodeMgr"
)

type DiskAttributes struct {
	identifier    kexec.NodeIdentifier
	name          string
	status        FacNodeStatus
	AssignedTo    *exec.RunControlEntry
	PackLabelInfo *kexec.PackLabelInfo
	IsPrepped     bool
	IsFixed       bool
}

func (da *DiskAttributes) GetFacNodeStatus() FacNodeStatus {
	return da.status
}

func (da *DiskAttributes) GetNodeCategoryType() nodeMgr.NodeCategoryType {
	return nodeMgr.NodeCategoryDevice
}

func (da *DiskAttributes) GetNodeDeviceType() nodeMgr.NodeDeviceType {
	return nodeMgr.NodeDeviceDisk
}

func (da *DiskAttributes) GetNodeIdentifier() kexec.NodeIdentifier {
	return da.identifier
}

func (da *DiskAttributes) GetNodeName() string {
	return da.name
}

func (da *DiskAttributes) SetFacNodeStatus(status FacNodeStatus) {
	da.status = status
}
