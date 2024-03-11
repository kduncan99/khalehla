// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/hardware"
	"khalehla/kexec"
)

type ReelAttributes struct {
	ReelNumber string
	IsLabeled  bool
}

type TapeAttributes struct {
	Identifier hardware.NodeIdentifier
	Name       string
	Status     kexec.FacNodeStatus
	AssignedTo *kexec.RunControlEntry
	ReelAttrs  *ReelAttributes
}

func (ta *TapeAttributes) GetFacNodeStatus() kexec.FacNodeStatus {
	return ta.Status
}

func (ta *TapeAttributes) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryDevice
}

func (ta *TapeAttributes) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceTape
}

func (ta *TapeAttributes) GetNodeIdentifier() hardware.NodeIdentifier {
	return ta.Identifier
}

func (ta *TapeAttributes) GetNodeName() string {
	return ta.Name
}

func (ta *TapeAttributes) SetFacNodeStatus(status kexec.FacNodeStatus) {
	ta.Status = status
}
