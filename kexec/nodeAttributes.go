// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/kexec/nodes"
)

type NodeAttributes interface {
	GetFacNodeStatus() FacNodeStatus
	GetNodeCategoryType() nodes.NodeCategoryType
	GetNodeDeviceType() nodes.NodeDeviceType
	GetNodeIdentifier() NodeIdentifier
	GetNodeName() string
	SetFacNodeStatus(status FacNodeStatus)
}
