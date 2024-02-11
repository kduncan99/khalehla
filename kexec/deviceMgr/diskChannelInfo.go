// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import (
	"khalehla/kexec/types"
	"khalehla/pkg"
)

type DiskChannelInfo struct {
	channelName    string
	nodeIdentifier types.NodeIdentifier
	channel        *DiskChannel
}

func NewDiskChannelInfo(channelName string) *DiskChannelInfo {
	return &DiskChannelInfo{
		channelName:    channelName,
		nodeIdentifier: types.NodeIdentifier(pkg.NewFromStringToFieldata(channelName, 1)[0]),
	}
}

func (dci *DiskChannelInfo) CreateNode() {
	dci.channel = NewDiskChannel()
}

func (dci *DiskChannelInfo) GetChannel() types.Channel {
	return dci.channel
}

func (dci *DiskChannelInfo) GetNodeIdentifier() types.NodeIdentifier {
	return dci.nodeIdentifier
}

func (dci *DiskChannelInfo) GetNodeName() string {
	return dci.channelName
}

func (dci *DiskChannelInfo) GetNodeStatus() types.NodeStatus {
	return types.NodeStatusUp
}

func (dci *DiskChannelInfo) GetNodeType() types.NodeType {
	return types.NodeTypeDisk
}
