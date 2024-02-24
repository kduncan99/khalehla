// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"khalehla/pkg"
)

type DiskChannelInfo struct {
	channelName       string
	channelIdentifier types.ChannelIdentifier
	channel           *DiskChannel
	deviceInfos       []*DiskDeviceInfo
}

func NewDiskChannelInfo(channelName string) *DiskChannelInfo {
	return &DiskChannelInfo{
		channelName:       channelName,
		channelIdentifier: types.ChannelIdentifier(pkg.NewFromStringToFieldata(channelName, 1)[0]),
		deviceInfos:       make([]*DiskDeviceInfo, 0),
	}
}

func (dci *DiskChannelInfo) CreateNode() {
	dci.channel = NewDiskChannel()
}

func (dci *DiskChannelInfo) GetChannel() Channel {
	return dci.channel
}

func (dci *DiskChannelInfo) GetChannelIdentifier() types.ChannelIdentifier {
	return dci.channelIdentifier
}

func (dci *DiskChannelInfo) GetChannelName() string {
	return dci.channelName
}

func (dci *DiskChannelInfo) GetDeviceInfos() []DeviceInfo {
	result := make([]DeviceInfo, len(dci.deviceInfos))
	for dx, di := range dci.deviceInfos {
		result[dx] = di
	}
	return result
}

func (dci *DiskChannelInfo) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryChannel
}

func (dci *DiskChannelInfo) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceDisk
}

func (dci *DiskChannelInfo) GetNodeIdentifier() types.NodeIdentifier {
	return types.NodeIdentifier(dci.channelIdentifier)
}

func (dci *DiskChannelInfo) GetNodeName() string {
	return dci.channelName
}

func (dci *DiskChannelInfo) IsAccessible() bool {
	return true
}

func (dci *DiskChannelInfo) Dump(dest io.Writer, indent string) {
	str := fmt.Sprintf("%v", dci.channelName)
	str += " devices:"
	for _, devInfo := range dci.deviceInfos {
		str += " " + devInfo.GetNodeName()
	}

	_, _ = fmt.Fprintf(dest, "%v%v\n", indent, str)

	dci.channel.Dump(dest, indent+"  ")
}
