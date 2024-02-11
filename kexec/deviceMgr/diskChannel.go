// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import (
	"fmt"
	"khalehla/kexec/types"
)

// DiskChannel routes IOs to the appropriate deviceInfos which it manages.
// Some day in the future we may add caching, perhaps in a CacheDiskChannel.
type DiskChannel struct {
	devices map[types.NodeIdentifier]*DiskDevice
}

func NewDiskChannel() *DiskChannel {
	return &DiskChannel{
		devices: make(map[types.NodeIdentifier]*DiskDevice),
	}
}

func (ch *DiskChannel) GetNodeType() types.NodeType {
	return types.NodeTypeDisk
}

func (ch *DiskChannel) AssignDevice(deviceIdentifier types.NodeIdentifier, device types.Device) error {
	if device.GetNodeType() != types.NodeTypeDisk {
		return fmt.Errorf("device is not a disk")
	}

	ch.devices[deviceIdentifier] = device.(*DiskDevice)
	return nil
}

func (ch *DiskChannel) StartIo(ioPacket types.IoPacket) {
	ioPacket.SetIoStatus(types.IosInProgress)
	if ioPacket.GetNodeType() != ch.GetNodeType() {
		ioPacket.SetIoStatus(types.IosInvalidNodeType)
		return
	}

	dev, ok := ch.devices[ioPacket.GetDeviceIdentifier()]
	if !ok {
		ioPacket.SetIoStatus(types.IosDeviceNotAttached)
		return
	}

	dev.StartIo(ioPacket)
}
