// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type ReelAttributes struct {
	ReelNumber string
	IsLabeled  bool
}

type TapeAttributes struct {
	Identifier NodeIdentifier
	Name       string
	Status     FacNodeStatus
	AssignedTo *RunControlEntry
	ReelAttrs  *ReelAttributes
}

func (ta *TapeAttributes) GetFacNodeStatus() FacNodeStatus {
	return ta.Status
}

func (ta *TapeAttributes) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryDevice
}

func (ta *TapeAttributes) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceTape
}

func (ta *TapeAttributes) GetNodeIdentifier() NodeIdentifier {
	return ta.Identifier
}

func (ta *TapeAttributes) GetNodeName() string {
	return ta.Name
}

func (ta *TapeAttributes) SetFacNodeStatus(status FacNodeStatus) {
	ta.Status = status
}
