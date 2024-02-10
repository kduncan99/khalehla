// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import (
	"fmt"
)

// DiskChannel routes IOs to the appropriate deviceInfos which it manages.
// Some day in the future we may add caching, perhaps in a CacheDiskChannel.
type DiskChannel struct {
	devices map[NodeIdentifier]*DiskDevice
}

func NewDiskChannel() *DiskChannel {
	return &DiskChannel{
		devices: make(map[NodeIdentifier]*DiskDevice),
	}
}

func (ch *DiskChannel) getNodeType() NodeType {
	return NodeTypeDisk
}

func (ch *DiskChannel) AssignDevice(deviceIdentifier NodeIdentifier, device Device) error {
	if device.getNodeType() != NodeTypeDisk {
		return fmt.Errorf("device is not a disk")
	}

	ch.devices[deviceIdentifier] = device.(*DiskDevice)
	return nil
}

func (ch *DiskChannel) StartIo(ioPacket IoPacket) {
	ioPacket.SetIoStatus(IosInProgress)
	if ioPacket.GetNodeType() != ch.getNodeType() {
		ioPacket.SetIoStatus(IosInvalidNodeType)
		return
	}

	dev, ok := ch.devices[ioPacket.GetDeviceIdentifier()]
	if !ok {
		ioPacket.SetIoStatus(IosDeviceNotAttached)
		return
	}

	dev.startIo(ioPacket)
}
