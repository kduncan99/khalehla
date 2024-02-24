// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

type NodeModelType uint

const (
	_ NodeModelType = iota

	// channels

	NodeModelDiskChannel
	NodeModelTapeChannel

	// devices

	NodeModelFileSystemDiskDevice
	NodeModelFileSystemTapeDevice
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
	"FSTAPE": {
		categoryType: NodeCategoryDevice,
		deviceType:   NodeDeviceTape,
		modelType:    NodeModelFileSystemTapeDevice,
	},
}
