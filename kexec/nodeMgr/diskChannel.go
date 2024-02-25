// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
)

// DiskChannel routes IOs to the appropriate deviceInfos which it manages.
// Some day in the future we may add caching, perhaps in a CacheDiskChannel.
type DiskChannel struct {
	devices map[types.NodeIdentifier]DiskDevice
}

func NewDiskChannel() *DiskChannel {
	return &DiskChannel{
		devices: make(map[types.NodeIdentifier]DiskDevice),
	}
}

func (ch *DiskChannel) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryChannel
}

func (ch *DiskChannel) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceDisk
}

func (ch *DiskChannel) GetNodeModelType() NodeModelType {
	return NodeModelDiskChannel
}

func (ch *DiskChannel) AssignDevice(nodeIdentifier types.NodeIdentifier, device Device) error {
	if device.GetNodeDeviceType() != NodeDeviceDisk {
		return fmt.Errorf("device is not a disk")
	}

	ch.devices[nodeIdentifier] = device.(*FileSystemDiskDevice)
	return nil
}

func (ch *DiskChannel) StartIo(ioPacket IoPacket) {
	ioPacket.SetIoStatus(IosInProgress)
	if ioPacket.GetNodeDeviceType() != ch.GetNodeDeviceType() {
		ioPacket.SetIoStatus(IosInvalidNodeType)
		return
	}

	dev, ok := ch.devices[ioPacket.GetNodeIdentifier()]
	if !ok {
		ioPacket.SetIoStatus(IosDeviceIsNotAccessible)
		return
	}

	dev.StartIo(ioPacket)
}

func (ch *DiskChannel) Dump(destination io.Writer, indent string) {
	// TODO
}
