// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/hardware"
	hardware2 "khalehla/old/hardware"
)

type DiskAttributes struct {
	Identifier    hardware.NodeIdentifier
	Name          string
	Status        FacNodeStatus
	PackLabelInfo *PackLabelInfo
	IsPrepped     bool
	IsFixed       bool
	IsRemovable   bool
}

func (da *DiskAttributes) GetFacNodeStatus() FacNodeStatus {
	return da.Status
}

func (da *DiskAttributes) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryDevice
}

func (da *DiskAttributes) GetNodeDeviceType() hardware2.NodeDeviceType {
	return hardware2.NodeDeviceDisk
}

func (da *DiskAttributes) GetNodeIdentifier() hardware.NodeIdentifier {
	return da.Identifier
}

func (da *DiskAttributes) GetNodeName() string {
	return da.Name
}

func (da *DiskAttributes) SetFacNodeStatus(status FacNodeStatus) {
	da.Status = status
}
