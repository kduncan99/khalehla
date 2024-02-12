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
	isAccessible     bool // can only be true if status is UP, RV, or SU and the device is assigned to at least one channel
	isMounted        bool
	reelName         string // only if isMounted
	channelInfos     []*TapeChannelInfo
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

func (tdi *TapeDeviceInfo) IsMounted() bool {
	return tdi.isMounted
}

func (tdi *TapeDeviceInfo) SetIsAccessible(isAccessible bool) {
	tdi.isAccessible = isAccessible
}

func (tdi *TapeDeviceInfo) Dump(dest io.Writer, indent string) {
	did := pkg.Word36(tdi.deviceIdentifier)
	str := fmt.Sprintf("%v id:%v %v\n",
		tdi.deviceName, did.ToStringAsOctal(), GetNodeStatusString(tdi.nodeStatus, tdi.isAccessible))
	if tdi.isMounted {
		str += " volume:" + tdi.reelName
	}
	str += " channels:"
	for _, chInfo := range tdi.channelInfos {
		chId := pkg.Word36(chInfo.channelIdentifier)
		str += " " + chId.ToStringAsFieldata()
	}

	_, _ = fmt.Fprintf(dest, "%v%v", indent, str)

	tdi.device.Dump(dest, indent+"  ")
}
