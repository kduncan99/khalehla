// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/hardware"
	hardware2 "khalehla/old/hardware"
)

type ReelAttributes struct {
	ReelNumber string
	IsLabeled  bool
}

type TapeAttributes struct {
	Identifier hardware.NodeIdentifier
	Name       string
	Status     FacNodeStatus
	AssignedTo *RunControlEntry
	ReelAttrs  *ReelAttributes
}

func (ta *TapeAttributes) GetFacNodeStatus() FacNodeStatus {
	return ta.Status
}

func (ta *TapeAttributes) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryDevice
}

func (ta *TapeAttributes) GetNodeDeviceType() hardware2.NodeDeviceType {
	return hardware2.NodeDeviceTape
}

func (ta *TapeAttributes) GetNodeIdentifier() hardware.NodeIdentifier {
	return ta.Identifier
}

func (ta *TapeAttributes) GetNodeName() string {
	return ta.Name
}

func (ta *TapeAttributes) SetFacNodeStatus(status FacNodeStatus) {
	ta.Status = status
}
