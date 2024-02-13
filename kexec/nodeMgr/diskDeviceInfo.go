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

type DiskDeviceInfo struct {
	deviceName       string
	deviceIdentifier types.DeviceIdentifier
	initialFileName  *string
	device           *DiskDevice
	nodeStatus       types.NodeStatus
	channelInfos     []*DiskChannelInfo
	isAccessible     bool // can only be true if status is UP, RV, or SU and the device is assigned to at least one channel
	isReady          bool // cached version of device.IsReady() - when there is a mismatch, we need to do something
}

// NewDiskDeviceInfo creates a new struct
// deviceName is required, but initialFileName can be set to nil if the device is not to be initially mounted
func NewDiskDeviceInfo(deviceName string, initialFileName *string) *DiskDeviceInfo {
	return &DiskDeviceInfo{
		deviceName:       deviceName,
		deviceIdentifier: types.DeviceIdentifier(pkg.NewFromStringToFieldata(deviceName, 1)[0]),
		nodeStatus:       types.NodeStatusUp,
		initialFileName:  initialFileName,
		channelInfos:     make([]*DiskChannelInfo, 0),
	}
}

func (ddi *DiskDeviceInfo) CreateNode() {
	ddi.device = NewDiskDevice(ddi.initialFileName)
}

func (ddi *DiskDeviceInfo) GetChannelInfos() []types.ChannelInfo {
	result := make([]types.ChannelInfo, len(ddi.channelInfos))
	for cx, ci := range ddi.channelInfos {
		result[cx] = ci
	}
	return result
}

func (ddi *DiskDeviceInfo) GetDevice() types.Device {
	return ddi.device
}

func (ddi *DiskDeviceInfo) GetDeviceIdentifier() types.DeviceIdentifier {
	return ddi.deviceIdentifier
}

func (ddi *DiskDeviceInfo) GetDeviceName() string {
	return ddi.deviceName
}

func (ddi *DiskDeviceInfo) GetInitialFileName() *string {
	return ddi.initialFileName
}

func (ddi *DiskDeviceInfo) GetNodeCategory() types.NodeCategory {
	return types.NodeCategoryDevice
}

func (ddi *DiskDeviceInfo) GetNodeIdentifier() types.NodeIdentifier {
	return types.NodeIdentifier(ddi.deviceIdentifier)
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

func (ddi *DiskDeviceInfo) IsReady() bool {
	return ddi.isReady
}

func (ddi *DiskDeviceInfo) SetIsAccessible(flag bool) {
	ddi.isAccessible = flag
}

func (ddi *DiskDeviceInfo) SetIsReady(flag bool) {
	ddi.isReady = flag
}

func (ddi *DiskDeviceInfo) Dump(dest io.Writer, indent string) {
	did := pkg.Word36(ddi.deviceIdentifier)
	str := fmt.Sprintf("%v id:0%v %v ready:%v",
		ddi.deviceName, did.ToStringAsOctal(), GetNodeStatusString(ddi.nodeStatus, ddi.isAccessible), ddi.isReady)

	str += " channels:"
	for _, chInfo := range ddi.channelInfos {
		str += " " + chInfo.GetChannelName()
	}

	_, _ = fmt.Fprintf(dest, "%v%v\n", indent, str)

	ddi.device.Dump(dest, indent+"  ")
}
