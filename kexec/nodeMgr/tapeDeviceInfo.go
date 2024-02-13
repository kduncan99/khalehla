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

type TapeDeviceInfo struct {
	deviceName       string
	deviceIdentifier types.DeviceIdentifier
	device           *TapeDevice
	nodeStatus       types.NodeStatus
	channelInfos     []*TapeChannelInfo
	isAccessible     bool // can only be true if status is UP, RV, or SU and the device is assigned to at least one channel
	isReady          bool // cached version of device.IsReady() - when there is a mismatch, we need to do something
}

// NewTapeDeviceInfo creates a new struct
func NewTapeDeviceInfo(deviceName string) *TapeDeviceInfo {
	return &TapeDeviceInfo{
		deviceName:       deviceName,
		deviceIdentifier: types.DeviceIdentifier(pkg.NewFromStringToFieldata(deviceName, 1)[0]),
		nodeStatus:       types.NodeStatusUp,
		channelInfos:     make([]*TapeChannelInfo, 0),
	}
}

func (tdi *TapeDeviceInfo) CreateNode() {
	tdi.device = NewTapeDevice()
}

func (tdi *TapeDeviceInfo) GetChannelInfos() []types.ChannelInfo {
	result := make([]types.ChannelInfo, len(tdi.channelInfos))
	for cx, ci := range tdi.channelInfos {
		result[cx] = ci
	}
	return result
}

func (tdi *TapeDeviceInfo) GetDevice() types.Device {
	return tdi.device
}

func (tdi *TapeDeviceInfo) GetDeviceIdentifier() types.DeviceIdentifier {
	return tdi.deviceIdentifier
}

func (tdi *TapeDeviceInfo) GetDeviceName() string {
	return tdi.deviceName
}

func (tdi *TapeDeviceInfo) GetNodeIdentifier() types.NodeIdentifier {
	return types.NodeIdentifier(tdi.deviceIdentifier)
}

func (tdi *TapeDeviceInfo) GetNodeName() string {
	return tdi.deviceName
}

func (tdi *TapeDeviceInfo) GetNodeStatus() types.NodeStatus {
	return tdi.nodeStatus
}

func (tdi *TapeDeviceInfo) GetNodeType() types.NodeType {
	return types.NodeTypeTape
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
	did := pkg.Word36(tdi.deviceIdentifier)
	str := fmt.Sprintf("%v id:%v %v ready:%v\n",
		tdi.deviceName, did.ToStringAsOctal(), GetNodeStatusString(tdi.nodeStatus, tdi.isAccessible), tdi.isReady)
	str += " channels:"
	for _, chInfo := range tdi.channelInfos {
		chId := pkg.Word36(chInfo.channelIdentifier)
		str += " " + chId.ToStringAsFieldata()
	}

	_, _ = fmt.Fprintf(dest, "%v%v", indent, str)

	tdi.device.Dump(dest, indent+"  ")
}
