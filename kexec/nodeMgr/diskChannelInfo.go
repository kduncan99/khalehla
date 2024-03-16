// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/hardware"
	"khalehla/hardware/channels"
)

type DiskChannelInfo struct {
	nodeName    string
	channel     *channels.DiskChannel
	deviceInfos []*DiskDeviceInfo
}

func NewDiskChannelInfo(nodeName string) *DiskChannelInfo {
	return &DiskChannelInfo{
		nodeName:    nodeName,
		deviceInfos: make([]*DiskDeviceInfo, 0),
	}
}

func (dci *DiskChannelInfo) CreateNode() {
	dci.channel = channels.NewDiskChannel()
}

func (dci *DiskChannelInfo) GetChannel() channels.Channel {
	return dci.channel
}

func (dci *DiskChannelInfo) GetDeviceInfos() []DeviceInfo {
	result := make([]DeviceInfo, len(dci.deviceInfos))
	for dx, di := range dci.deviceInfos {
		result[dx] = di
	}
	return result
}

func (dci *DiskChannelInfo) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryChannel
}

func (dci *DiskChannelInfo) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceDisk
}

func (dci *DiskChannelInfo) GetNodeIdentifier() hardware.NodeIdentifier {
	return dci.channel.GetNodeIdentifier()
}

func (dci *DiskChannelInfo) GetNodeName() string {
	return dci.nodeName
}

func (dci *DiskChannelInfo) IsAccessible() bool {
	return true
}

func (dci *DiskChannelInfo) Dump(dest io.Writer, indent string) {
	str := fmt.Sprintf("%v", dci.nodeName)
	str += " devices:"
	for _, devInfo := range dci.deviceInfos {
		str += " " + devInfo.GetNodeName()
	}

	_, _ = fmt.Fprintf(dest, "%v%v\n", indent, str)

	dci.channel.Dump(dest, indent+"  ")
}
