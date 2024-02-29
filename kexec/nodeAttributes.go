// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type NodeAttributes interface {
	GetFacNodeStatus() FacNodeStatus
	GetNodeCategoryType() NodeCategoryType
	GetNodeDeviceType() NodeDeviceType
	GetNodeIdentifier() NodeIdentifier
	GetNodeName() string
	SetFacNodeStatus(status FacNodeStatus)
}
