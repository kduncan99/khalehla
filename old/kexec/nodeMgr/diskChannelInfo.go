// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"

	"khalehla/hardware"
	hardware2 "khalehla/old/hardware"
	channels2 "khalehla/old/hardware/channels"
)

type DiskChannelInfo struct {
	nodeName    string
	channel     *channels2.DiskChannel
	deviceInfos []*DiskDeviceInfo
}

func NewDiskChannelInfo(nodeName string) *DiskChannelInfo {
	return &DiskChannelInfo{
		nodeName:    nodeName,
		deviceInfos: make([]*DiskDeviceInfo, 0),
	}
}

func (dci *DiskChannelInfo) CreateNode() {
	dci.channel = channels2.NewDiskChannel()
}

func (dci *DiskChannelInfo) GetChannel() channels2.Channel {
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

func (dci *DiskChannelInfo) GetNodeDeviceType() hardware2.NodeDeviceType {
	return hardware2.NodeDeviceDisk
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
