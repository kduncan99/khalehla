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

type TapeChannelInfo struct {
	nodeName       string
	nodeIdentifier kexec.NodeIdentifier
	channel        *nodes.TapeChannel
	deviceInfos    []*TapeDeviceInfo
}

func NewTapeChannelInfo(nodeName string) *TapeChannelInfo {
	return &TapeChannelInfo{
		nodeName:       nodeName,
		nodeIdentifier: kexec.NodeIdentifier(pkg.NewFromStringToFieldata(nodeName, 1)[0]),
		deviceInfos:    make([]*TapeDeviceInfo, 0),
	}
}

func (tci *TapeChannelInfo) CreateNode() {
	tci.channel = nodes.NewTapeChannel()
}

func (tci *TapeChannelInfo) GetChannel() nodes.Channel {
	return tci.channel
}

func (tci *TapeChannelInfo) GetDeviceInfos() []DeviceInfo {
	result := make([]DeviceInfo, len(tci.deviceInfos))
	for dx, di := range tci.deviceInfos {
		result[dx] = di
	}
	return result
}

func (tci *TapeChannelInfo) GetNodeCategoryType() nodes.NodeCategoryType {
	return nodes.NodeCategoryChannel
}

func (tci *TapeChannelInfo) GetNodeDeviceType() nodes.NodeDeviceType {
	return nodes.NodeDeviceTape
}

func (tci *TapeChannelInfo) GetNodeIdentifier() kexec.NodeIdentifier {
	return kexec.NodeIdentifier(tci.nodeIdentifier)
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
