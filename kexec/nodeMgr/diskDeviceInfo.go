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
	isReady          bool
	packName         string // only if device is ready, read by probeDisk
	isPrepped        bool   // as above
	isFixed          bool   // as above
	ldatIndex        int
	geometry         *types.DiskPackGeometry // as above
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

func (ddi *DiskDeviceInfo) GetLDATIndex() int {
	return ddi.ldatIndex
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

func (ddi *DiskDeviceInfo) IsFixed() bool {
	return ddi.isFixed
}

func (ddi *DiskDeviceInfo) IsPrepped() bool {
	return ddi.isPrepped
}

func (ddi *DiskDeviceInfo) SetIsAccessible(isAccessible bool) {
	ddi.isAccessible = isAccessible
}

func (ddi *DiskDeviceInfo) Dump(dest io.Writer, indent string) {
	did := pkg.Word36(ddi.deviceIdentifier)
	str := fmt.Sprintf("%v id:%v %v\n",
		ddi.deviceName, did.ToStringAsOctal(), GetNodeStatusString(ddi.nodeStatus, ddi.isAccessible))

	if ddi.isPrepped {
		str += "PREPPED"
		if ddi.isFixed {
			str += " FIXED"
		} else {
			str += " REM"
		}
	}

	str += " channels:"
	for _, chInfo := range ddi.channelInfos {
		chId := pkg.Word36(chInfo.channelIdentifier)
		str += " " + chId.ToStringAsFieldata()
	}

	_, _ = fmt.Fprintf(dest, "%v%v", indent, str)

	ddi.device.Dump(dest, indent+"  ")
}
