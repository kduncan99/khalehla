// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/types"
)

type ReelAttributes struct {
	ReelNumber string
	IsLabeled  bool
}

type TapeAttributes struct {
	identifier types.NodeIdentifier
	name       string
	status     FacNodeStatus
	AssignedTo *types.RunControlEntry
	ReelAttrs  *ReelAttributes
}

func (ta *TapeAttributes) GetFacNodeStatus() FacNodeStatus {
	return ta.status
}

func (ta *TapeAttributes) GetNodeCategoryType() nodeMgr.NodeCategoryType {
	return nodeMgr.NodeCategoryDevice
}

func (ta *TapeAttributes) GetNodeDeviceType() nodeMgr.NodeDeviceType {
	return nodeMgr.NodeDeviceTape
}

func (ta *TapeAttributes) GetNodeIdentifier() types.NodeIdentifier {
	return ta.identifier
}

func (ta *TapeAttributes) GetNodeName() string {
	return ta.name
}

func (ta *TapeAttributes) SetFacNodeStatus(status FacNodeStatus) {
	ta.status = status
}
