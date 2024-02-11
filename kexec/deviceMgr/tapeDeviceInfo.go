// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import "khalehla/kexec/types"

type TapeDeviceInfo struct {
	deviceName     string
	nodeIdentifier types.NodeIdentifier
	device         *TapeDevice
	nodeStatus     types.NodeStatus
	isMounted      bool
}

// NewTapeDeviceInfo creates a new struct
func NewTapeDeviceInfo(deviceName string) *TapeDeviceInfo {
	return &TapeDeviceInfo{
		deviceName: deviceName,
		nodeStatus: types.NodeStatusUp,
	}
}

func (tdi *TapeDeviceInfo) CreateNode() {
	tdi.device = NewTapeDevice()
}

func (tdi *TapeDeviceInfo) GetDevice() types.Device {
	return tdi.device
}

func (tdi *TapeDeviceInfo) GetNodeIdentifier() types.NodeIdentifier {
	return tdi.nodeIdentifier
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

func (tdi *TapeDeviceInfo) IsMounted() bool {
	return tdi.isMounted
}
