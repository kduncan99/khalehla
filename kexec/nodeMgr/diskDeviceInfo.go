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
	nodeName        string
	nodeIdentifier  types.NodeIdentifier
	initialFileName *string
	device          *FileSystemDiskDevice
	channelInfos    []*DiskChannelInfo
	isAccessible    bool // can only be true if status is UP, RV, or SU and the device is assigned to at least one channel
	isReady         bool // cached version of device.IsReady() - when there is a mismatch, we need to do something
}

// NewDiskDeviceInfo creates a new struct
// nodeName is required, but initialFileName can be set to nil if the device is not to be initially mounted
func NewDiskDeviceInfo(nodeName string, initialFileName *string) *DiskDeviceInfo {
	return &DiskDeviceInfo{
		nodeName:        nodeName,
		nodeIdentifier:  types.NodeIdentifier(pkg.NewFromStringToFieldata(nodeName, 1)[0]),
		initialFileName: initialFileName,
		channelInfos:    make([]*DiskChannelInfo, 0),
	}
}

func (ddi *DiskDeviceInfo) CreateNode() {
	ddi.device = NewFileSystemDiskDevice(ddi.initialFileName)
}

func (ddi *DiskDeviceInfo) GetChannelInfos() []ChannelInfo {
	result := make([]ChannelInfo, len(ddi.channelInfos))
	for cx, ci := range ddi.channelInfos {
		result[cx] = ci
	}
	return result
}

func (ddi *DiskDeviceInfo) GetDevice() Device {
	return ddi.device
}

func (ddi *DiskDeviceInfo) GetInitialFileName() *string {
	return ddi.initialFileName
}

func (ddi *DiskDeviceInfo) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryDevice
}

func (ddi *DiskDeviceInfo) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceDisk
}

func (ddi *DiskDeviceInfo) GetNodeIdentifier() types.NodeIdentifier {
	return ddi.nodeIdentifier
}

func (ddi *DiskDeviceInfo) GetNodeName() string {
	return ddi.nodeName
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
	did := pkg.Word36(ddi.nodeIdentifier)
	str := fmt.Sprintf("%v id:0%v ready:%v", ddi.nodeName, did.ToStringAsOctal(), ddi.isReady)

	str += " channels:"
	for _, chInfo := range ddi.channelInfos {
		str += " " + chInfo.GetNodeName()
	}

	_, _ = fmt.Fprintf(dest, "%v%v\n", indent, str)

	ddi.device.Dump(dest, indent+"  ")
}
