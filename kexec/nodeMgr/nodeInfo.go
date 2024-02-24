// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"khalehla/kexec/types"
)

// NodeInfo contains all the exec-managed information regarding a particular node
type NodeInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() NodeCategoryType
	GetNodeDeviceType() NodeDeviceType
	GetNodeIdentifier() types.NodeIdentifier
	GetNodeName() string
	GetNodeStatus() types.NodeStatus
	IsAccessible() bool
}
