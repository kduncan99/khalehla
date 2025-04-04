// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"

	"khalehla/hardware"
	hardware2 "khalehla/old/hardware"
	"khalehla/old/hardware/channels"
)

// ChannelInfo is intended primarily as a means of documenting the use of a more generic NodeInfo
type ChannelInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetChannel() channels.Channel
	GetDeviceInfos() []DeviceInfo
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware2.NodeDeviceType
	GetNodeIdentifier() hardware.NodeIdentifier
	GetNodeName() string
	IsAccessible() bool
}
