// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"

	"khalehla/hardware"
	hardware2 "khalehla/old/hardware"
	"khalehla/old/hardware/devices"
)

// DeviceInfo is intended primarily as a means of documenting the use of a more generic NodeInfo
type DeviceInfo interface {
	CreateNode()
	Dump(destination io.Writer, indent string)
	GetChannelInfos() []ChannelInfo
	GetDevice() devices.Device
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware2.NodeDeviceType
	GetNodeIdentifier() hardware.NodeIdentifier
	GetNodeName() string
	IsAccessible() bool
	IsReady() bool
	SetIsAccessible(bool)
	SetIsReady(flag bool)
}
