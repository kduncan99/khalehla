// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/kexec"
	"khalehla/kexec/nodes"
	"khalehla/pkg"
)

type DiskChannelInfo struct {
	nodeName       string
	nodeIdentifier kexec.NodeIdentifier
	channel        *nodes.DiskChannel
	deviceInfos    []*DiskDeviceInfo
}

func NewDiskChannelInfo(nodeName string) *DiskChannelInfo {
	return &DiskChannelInfo{
		nodeName:       nodeName,
		nodeIdentifier: kexec.NodeIdentifier(pkg.NewFromStringToFieldata(nodeName, 1)[0]),
		deviceInfos:    make([]*DiskDeviceInfo, 0),
	}
}

func (dci *DiskChannelInfo) CreateNode() {
	dci.channel = nodes.NewDiskChannel()
}

func (dci *DiskChannelInfo) GetChannel() nodes.Channel {
	return dci.channel
}

func (dci *DiskChannelInfo) GetDeviceInfos() []DeviceInfo {
	result := make([]DeviceInfo, len(dci.deviceInfos))
	for dx, di := range dci.deviceInfos {
		result[dx] = di
	}
	return result
}

func (dci *DiskChannelInfo) GetNodeCategoryType() nodes.NodeCategoryType {
	return nodes.NodeCategoryChannel
}

func (dci *DiskChannelInfo) GetNodeDeviceType() nodes.NodeDeviceType {
	return nodes.NodeDeviceDisk
}

func (dci *DiskChannelInfo) GetNodeIdentifier() kexec.NodeIdentifier {
	return kexec.NodeIdentifier(dci.nodeIdentifier)
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
