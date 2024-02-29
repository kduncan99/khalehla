// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import "khalehla/kexec"

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
	CategoryType kexec.NodeCategoryType
	DeviceType   kexec.NodeDeviceType
	ModelType    NodeModelType
}

var NodeModelTable = map[string]NodeModel{
	"FSDISK": {
		CategoryType: kexec.NodeCategoryDevice,
		DeviceType:   kexec.NodeDeviceDisk,
		ModelType:    NodeModelFileSystemDiskDevice,
	},
	"RMDISK": {
		CategoryType: kexec.NodeCategoryDevice,
		DeviceType:   kexec.NodeDeviceDisk,
		ModelType:    NodeModelRAMDiskDevice,
	},
	"SCDISK": {
		CategoryType: kexec.NodeCategoryDevice,
		DeviceType:   kexec.NodeDeviceDisk,
		ModelType:    NodeModelSCSIDiskDevice,
	},
	"FSTAPE": {
		CategoryType: kexec.NodeCategoryDevice,
		DeviceType:   kexec.NodeDeviceTape,
		ModelType:    NodeModelFileSystemTapeDevice,
	},
	"SCTAPE": {
		CategoryType: kexec.NodeCategoryDevice,
		DeviceType:   kexec.NodeDeviceTape,
		ModelType:    NodeModelSCSITapeDevice,
	},
}
