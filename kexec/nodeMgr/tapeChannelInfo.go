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

type TapeChannelInfo struct {
	channelName       string
	channelIdentifier types.ChannelIdentifier
	channel           *TapeChannel
	deviceInfos       []*TapeDeviceInfo
}

func NewTapeChannelInfo(channelName string) *TapeChannelInfo {
	return &TapeChannelInfo{
		channelName:       channelName,
		channelIdentifier: types.ChannelIdentifier(pkg.NewFromStringToFieldata(channelName, 1)[0]),
		deviceInfos:       make([]*TapeDeviceInfo, 0),
	}
}

func (tci *TapeChannelInfo) CreateNode() {
	tci.channel = NewTapeChannel()
}

func (tci *TapeChannelInfo) GetChannel() Channel {
	return tci.channel
}

func (tci *TapeChannelInfo) GetChannelIdentifier() types.ChannelIdentifier {
	return tci.channelIdentifier
}

func (tci *TapeChannelInfo) GetChannelName() string {
	return tci.channelName
}

func (tci *TapeChannelInfo) GetDeviceInfos() []DeviceInfo {
	result := make([]DeviceInfo, len(tci.deviceInfos))
	for dx, di := range tci.deviceInfos {
		result[dx] = di
	}
	return result
}

func (tci *TapeChannelInfo) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryChannel
}

func (tci *TapeChannelInfo) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceTape
}

func (tci *TapeChannelInfo) GetNodeIdentifier() types.NodeIdentifier {
	return types.NodeIdentifier(tci.channelIdentifier)
}

func (tci *TapeChannelInfo) GetNodeName() string {
	return tci.channelName
}

func (tci *TapeChannelInfo) GetNodeStatus() types.NodeStatus {
	return types.NodeStatusUp
}

func (tci *TapeChannelInfo) IsAccessible() bool {
	return true
}

func (tci *TapeChannelInfo) Dump(dest io.Writer, indent string) {
	str := fmt.Sprintf("%v %v",
		tci.channelName, GetNodeStatusString(types.NodeStatusUp, true))
	str += " devices:"
	for _, devInfo := range tci.deviceInfos {
		str += " " + devInfo.GetNodeName()
	}

	_, _ = fmt.Fprintf(dest, "%v%v\n", indent, str)

	tci.channel.Dump(dest, indent+"  ")
}
