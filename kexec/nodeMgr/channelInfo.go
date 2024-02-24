// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"khalehla/kexec/types"
)

// ChannelInfo is intended primarily as a means of documenting the use of a more generic NodeInfo
type ChannelInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetChannel() Channel
	GetChannelName() string
	GetChannelIdentifier() types.ChannelIdentifier
	GetDeviceInfos() []DeviceInfo
	GetNodeCategoryType() NodeCategoryType
	GetNodeDeviceType() NodeDeviceType
	GetNodeIdentifier() types.NodeIdentifier
	GetNodeName() string
	GetNodeStatus() types.NodeStatus
	IsAccessible() bool
}
