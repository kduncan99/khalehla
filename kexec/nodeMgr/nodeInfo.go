// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"khalehla/kexec"
	"khalehla/kexec/nodes"
)

// NodeInfo contains all the exec-managed information regarding a particular node
type NodeInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() nodes.NodeCategoryType
	GetNodeDeviceType() nodes.NodeDeviceType
	GetNodeIdentifier() kexec.NodeIdentifier
	GetNodeName() string
	IsAccessible() bool
}
