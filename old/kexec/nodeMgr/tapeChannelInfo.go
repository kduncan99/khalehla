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

type TapeChannelInfo struct {
	nodeName    string
	channel     *channels2.TapeChannel
	deviceInfos []*TapeDeviceInfo
}

func NewTapeChannelInfo(nodeName string) *TapeChannelInfo {
	return &TapeChannelInfo{
		nodeName:    nodeName,
		deviceInfos: make([]*TapeDeviceInfo, 0),
	}
}

func (tci *TapeChannelInfo) CreateNode() {
	tci.channel = channels2.NewTapeChannel()
}

func (tci *TapeChannelInfo) GetChannel() channels2.Channel {
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

func (tci *TapeChannelInfo) GetNodeDeviceType() hardware2.NodeDeviceType {
	return hardware2.NodeDeviceTape
}

func (tci *TapeChannelInfo) GetNodeIdentifier() hardware.NodeIdentifier {
	return tci.channel.GetNodeIdentifier()
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
