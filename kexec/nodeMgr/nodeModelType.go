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
	categoryType NodeCategoryType
	deviceType   NodeDeviceType
	modelType    NodeModelType
}

var NodeModelTable = map[string]NodeModel{
	"FSDISK": {
		categoryType: NodeCategoryDevice,
		deviceType:   NodeDeviceDisk,
		modelType:    NodeModelFileSystemDiskDevice,
	},
	"RMDISK": {
		categoryType: NodeCategoryDevice,
		deviceType:   NodeDeviceDisk,
		modelType:    NodeModelRAMDiskDevice,
	},
	"SCDISK": {
		categoryType: NodeCategoryDevice,
		deviceType:   NodeDeviceDisk,
		modelType:    NodeModelSCSIDiskDevice,
	},
	"FSTAPE": {
		categoryType: NodeCategoryDevice,
		deviceType:   NodeDeviceTape,
		modelType:    NodeModelFileSystemTapeDevice,
	},
	"SCTAPE": {
		categoryType: NodeCategoryDevice,
		deviceType:   NodeDeviceTape,
		modelType:    NodeModelSCSITapeDevice,
	},
}
