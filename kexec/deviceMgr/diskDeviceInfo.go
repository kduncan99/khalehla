// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import (
	"khalehla/kexec/types"
	"khalehla/pkg"
)

// -------------------------------------------------------------------------------------

type DiskDeviceInfo struct {
	deviceName      string
	nodeIdentifier  types.NodeIdentifier
	initialFileName *string
	device          *DiskDevice
	nodeStatus      types.NodeStatus
	isAccessible    bool // can only be true if status is UP, RV, or SU and the device is assigned to at least one channel
	isMounted       bool
	isPrepped       bool
	isFixed         bool
}

// NewDiskDeviceInfo creates a new struct
// deviceName is required, but initialFileName can be set to nil if the device is not to be initially mounted
func NewDiskDeviceInfo(deviceName string, initialFileName *string) *DiskDeviceInfo {
	return &DiskDeviceInfo{
		deviceName:      deviceName,
		nodeIdentifier:  types.NodeIdentifier(pkg.NewFromStringToFieldata(deviceName, 1)[0]),
		nodeStatus:      types.NodeStatusUp,
		isAccessible:    false,
		initialFileName: initialFileName,
		isMounted:       false,
		isPrepped:       false,
		isFixed:         false,
	}
}

func (ddi *DiskDeviceInfo) CreateNode() {
	ddi.device = NewDiskDevice(ddi.initialFileName)
}

func (ddi *DiskDeviceInfo) GetDevice() types.Device {
	return ddi.device
}

func (ddi *DiskDeviceInfo) GetInitialFileName() *string {
	return ddi.initialFileName
}

func (ddi *DiskDeviceInfo) GetNodeIdentifier() types.NodeIdentifier {
	return ddi.nodeIdentifier
}

func (ddi *DiskDeviceInfo) GetNodeName() string {
	return ddi.deviceName
}

func (ddi *DiskDeviceInfo) GetNodeStatus() types.NodeStatus {
	return ddi.nodeStatus
}

func (ddi *DiskDeviceInfo) GetNodeType() types.NodeType {
	return types.NodeTypeDisk
}

func (ddi *DiskDeviceInfo) IsAccessible() bool {
	return ddi.isAccessible
}

func (ddi *DiskDeviceInfo) IsFixed() bool {
	return ddi.isFixed
}

func (ddi *DiskDeviceInfo) IsMounted() bool {
	return ddi.isMounted
}

func (ddi *DiskDeviceInfo) IsPrepped() bool {
	return ddi.isPrepped
}

func (ddi *DiskDeviceInfo) SetIsAccessible(isAccessible bool) {
	ddi.isAccessible = isAccessible
}
