// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"

	"khalehla/hardware"
	"khalehla/hardware/devices"
	hardware2 "khalehla/old/hardware"
	devices2 "khalehla/old/hardware/devices"
)

type TapeDeviceInfo struct {
	nodeName     string
	device       *devices.FileSystemTapeDevice
	channelInfos []*TapeChannelInfo
	isAccessible bool // can only be true if status is UP, RV, or SU and the device is assigned to at least one channel
	isReady      bool // cached version of device.IsReady() - when there is a mismatch, we need to do something
}

// NewTapeDeviceInfo creates a new struct
func NewTapeDeviceInfo(nodeName string) *TapeDeviceInfo {
	return &TapeDeviceInfo{
		nodeName:     nodeName,
		channelInfos: make([]*TapeChannelInfo, 0),
	}
}

func (tdi *TapeDeviceInfo) CreateNode() {
	tdi.device = devices.NewFileSystemTapeDevice()
}

func (tdi *TapeDeviceInfo) GetChannelInfos() []ChannelInfo {
	result := make([]ChannelInfo, len(tdi.channelInfos))
	for cx, ci := range tdi.channelInfos {
		result[cx] = ci
	}
	return result
}

func (tdi *TapeDeviceInfo) GetDevice() devices2.Device {
	return tdi.device
}

func (tdi *TapeDeviceInfo) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryDevice
}

func (tdi *TapeDeviceInfo) GetNodeDeviceType() hardware2.NodeDeviceType {
	return hardware2.NodeDeviceTape
}

func (tdi *TapeDeviceInfo) GetNodeIdentifier() hardware.NodeIdentifier {
	return tdi.device.GetNodeIdentifier()
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
	str := fmt.Sprintf("%v ready:%v", tdi.nodeName, tdi.isReady)
	str += " channels:"
	for _, chInfo := range tdi.channelInfos {
		str += " " + chInfo.GetNodeName()
	}

	_, _ = fmt.Fprintf(dest, "%v%v\n", indent, str)

	tdi.device.Dump(dest, indent+"  ")
}
