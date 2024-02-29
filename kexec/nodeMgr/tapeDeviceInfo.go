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

type TapeDeviceInfo struct {
	nodeName       string
	nodeIdentifier kexec.NodeIdentifier
	device         *nodes.FileSystemTapeDevice
	channelInfos   []*TapeChannelInfo
	isAccessible   bool // can only be true if status is UP, RV, or SU and the device is assigned to at least one channel
	isReady        bool // cached version of device.IsReady() - when there is a mismatch, we need to do something
}

// NewTapeDeviceInfo creates a new struct
func NewTapeDeviceInfo(nodeName string) *TapeDeviceInfo {
	return &TapeDeviceInfo{
		nodeName:       nodeName,
		nodeIdentifier: kexec.NodeIdentifier(pkg.NewFromStringToFieldata(nodeName, 1)[0]),
		channelInfos:   make([]*TapeChannelInfo, 0),
	}
}

func (tdi *TapeDeviceInfo) CreateNode() {
	tdi.device = nodes.NewFileSystemTapeDevice()
}

func (tdi *TapeDeviceInfo) GetChannelInfos() []ChannelInfo {
	result := make([]ChannelInfo, len(tdi.channelInfos))
	for cx, ci := range tdi.channelInfos {
		result[cx] = ci
	}
	return result
}

func (tdi *TapeDeviceInfo) GetDevice() nodes.Device {
	return tdi.device
}

func (tdi *TapeDeviceInfo) GetNodeCategoryType() nodes.NodeCategoryType {
	return nodes.NodeCategoryDevice
}

func (tdi *TapeDeviceInfo) GetNodeDeviceType() nodes.NodeDeviceType {
	return nodes.NodeDeviceTape
}

func (tdi *TapeDeviceInfo) GetNodeIdentifier() kexec.NodeIdentifier {
	return tdi.nodeIdentifier
}

func (tdi *TapeDeviceInfo) GetNodeName() string {
	return tdi.nodeName
}

func (tdi *TapeDeviceInfo) IsAccessible() bool {
	return tdi.isAccessible
}

func (tdi *TapeDeviceInfo) IsReady() bool {
	return tdi.isReady
}

func (tdi *TapeDeviceInfo) SetIsAccessible(flag bool) {
	tdi.isAccessible = flag
}

func (tdi *TapeDeviceInfo) SetIsReady(flag bool) {
	tdi.isReady = flag
}

func (tdi *TapeDeviceInfo) Dump(dest io.Writer, indent string) {
	did := pkg.Word36(tdi.nodeIdentifier)
	str := fmt.Sprintf("%v id:0%v ready:%v", tdi.nodeName, did.ToStringAsOctal(), tdi.isReady)
	str += " channels:"
	for _, chInfo := range tdi.channelInfos {
		str += " " + chInfo.GetNodeName()
	}

	_, _ = fmt.Fprintf(dest, "%v%v\n", indent, str)

	tdi.device.Dump(dest, indent+"  ")
}
