// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

type NodeModelType uint

const (
	_ NodeModelType = iota

	// channels

	NodeModelDiskChannel
	NodeModelCacheDiskChannel
	NodeModelTapeLibraryChannel

	// devices

	NodeModelFileSystemDiskDevice
	NodeModelRAMDiskDevice
	NodeModelSCSIDiskDevice
	NodeModelFileSystemTapeDevice
	NodeModelSCSITapeDevice
)

type NodeModel struct {
	CategoryType NodeCategoryType
	DeviceType   NodeDeviceType
	ModelType    NodeModelType
}

var NodeModelTable = map[string]NodeModel{
	"FSDISK": {
		CategoryType: NodeCategoryDevice,
		DeviceType:   NodeDeviceDisk,
		ModelType:    NodeModelFileSystemDiskDevice,
	},
	"RMDISK": {
		CategoryType: NodeCategoryDevice,
		DeviceType:   NodeDeviceDisk,
		ModelType:    NodeModelRAMDiskDevice,
	},
	"SCDISK": {
		CategoryType: NodeCategoryDevice,
		DeviceType:   NodeDeviceDisk,
		ModelType:    NodeModelSCSIDiskDevice,
	},
	"FSTAPE": {
		CategoryType: NodeCategoryDevice,
		DeviceType:   NodeDeviceTape,
		ModelType:    NodeModelFileSystemTapeDevice,
	},
	"SCTAPE": {
		CategoryType: NodeCategoryDevice,
		DeviceType:   NodeDeviceTape,
		ModelType:    NodeModelSCSITapeDevice,
	},
}
