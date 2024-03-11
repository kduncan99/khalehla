// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/hardware"
	"khalehla/hardware/channels"
	"khalehla/pkg"
)

type TapeChannelInfo struct {
	nodeName       string
	nodeIdentifier hardware.NodeIdentifier
	channel        *channels.TapeChannel
	deviceInfos    []*TapeDeviceInfo
}

func NewTapeChannelInfo(nodeName string) *TapeChannelInfo {
	return &TapeChannelInfo{
		nodeName:       nodeName,
		nodeIdentifier: hardware.NodeIdentifier(pkg.NewFromStringToFieldata(nodeName, 1)[0]),
		deviceInfos:    make([]*TapeDeviceInfo, 0),
	}
}

func (tci *TapeChannelInfo) CreateNode() {
	tci.channel = channels.NewTapeChannel()
}

func (tci *TapeChannelInfo) GetChannel() channels.Channel {
	return tci.channel
}

func (tci *TapeChannelInfo) GetDeviceInfos() []DeviceInfo {
	result := make([]DeviceInfo, len(tci.deviceInfos))
	for dx, di := range tci.deviceInfos {
		result[dx] = di
	}
	return result
}

func (tci *TapeChannelInfo) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryChannel
}

func (tci *TapeChannelInfo) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceTape
}

func (tci *TapeChannelInfo) GetNodeIdentifier() hardware.NodeIdentifier {
	return hardware.NodeIdentifier(tci.nodeIdentifier)
}

func (tci *TapeChannelInfo) GetNodeName() string {
	return tci.nodeName
}

func (tci *TapeChannelInfo) IsAccessible() bool {
	return true
}

func (tci *TapeChannelInfo) Dump(dest io.Writer, indent string) {
	str := fmt.Sprintf("%v", tci.nodeName)
	str += " devices:"
	for _, devInfo := range tci.deviceInfos {
		str += " " + devInfo.GetNodeName()
	}

	_, _ = fmt.Fprintf(dest, "%v%v\n", indent, str)

	tci.channel.Dump(dest, indent+"  ")
}
