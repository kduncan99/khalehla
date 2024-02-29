// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"khalehla/kexec"
	"khalehla/kexec/nodes"
)

// ChannelInfo is intended primarily as a means of documenting the use of a more generic NodeInfo
type ChannelInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetChannel() nodes.Channel
	GetDeviceInfos() []DeviceInfo
	GetNodeCategoryType() nodes.NodeCategoryType
	GetNodeDeviceType() nodes.NodeDeviceType
	GetNodeIdentifier() kexec.NodeIdentifier
	GetNodeName() string
	IsAccessible() bool
}
