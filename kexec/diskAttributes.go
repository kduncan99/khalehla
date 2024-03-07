// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type DiskAttributes struct {
	Identifier    NodeIdentifier
	Name          string
	Status        FacNodeStatus
	AssignedTo    *RunControlEntry
	PackLabelInfo *PackLabelInfo
	IsPrepped     bool
	IsFixed       bool
}

func (da *DiskAttributes) GetFacNodeStatus() FacNodeStatus {
	return da.Status
}

func (da *DiskAttributes) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryDevice
}

func (da *DiskAttributes) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceDisk
}

func (da *DiskAttributes) GetNodeIdentifier() NodeIdentifier {
	return da.Identifier
}

func (da *DiskAttributes) GetNodeName() string {
	return da.Name
}

func (da *DiskAttributes) SetFacNodeStatus(status FacNodeStatus) {
	da.Status = status
}
