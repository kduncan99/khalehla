// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"

	"khalehla/hardware"
	hardware2 "khalehla/old/hardware"
)

// NodeInfo contains all the exec-managed information regarding a particular node
type NodeInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware2.NodeDeviceType
	GetNodeIdentifier() hardware.NodeIdentifier
	GetNodeName() string
	IsAccessible() bool
}
