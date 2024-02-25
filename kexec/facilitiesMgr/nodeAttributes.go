// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/nodeMgr"
)

type NodeAttributes interface {
	GetFacNodeStatus() FacNodeStatus
	GetNodeCategoryType() nodeMgr.NodeCategoryType
	GetNodeDeviceType() nodeMgr.NodeDeviceType
	GetNodeIdentifier() kexec.NodeIdentifier
	GetNodeName() string
	SetFacNodeStatus(status FacNodeStatus)
}
