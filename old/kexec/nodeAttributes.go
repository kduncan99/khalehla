// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/hardware"
	hardware2 "khalehla/old/hardware"
)

type INodeAttributes interface {
	GetFacNodeStatus() FacNodeStatus
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware2.NodeDeviceType
	GetNodeIdentifier() hardware.NodeIdentifier
	GetNodeName() string
	SetFacNodeStatus(status FacNodeStatus)
}
