// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package channels

import (
	"fmt"
	"io"
	"khalehla/hardware"
	"khalehla/hardware/devices"
	"khalehla/hardware/ioPackets"
)

// DiskChannel routes IOs to the appropriate deviceInfos which it manages.
// Some day in the future we may add caching, perhaps in a CacheDiskChannel.
type DiskChannel struct {
	devices map[hardware.NodeIdentifier]devices.DiskDevice
}

func NewDiskChannel() *DiskChannel {
	return &DiskChannel{
		devices: make(map[hardware.NodeIdentifier]devices.DiskDevice),
	}
}

func (ch *DiskChannel) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryChannel
}

func (ch *DiskChannel) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceDisk
}

func (ch *DiskChannel) GetNodeModelType() hardware.NodeModelType {
	return hardware.NodeModelDiskChannel
}

func (ch *DiskChannel) AssignDevice(nodeIdentifier hardware.NodeIdentifier, device devices.Device) error {
	if device.GetNodeDeviceType() != hardware.NodeDeviceDisk {
		return fmt.Errorf("device is not a disk")
	}

	ch.devices[nodeIdentifier] = device.(*devices.FileSystemDiskDevice)
	return nil
}

func (ch *DiskChannel) StartIo(ioPacket ioPackets.IoPacket) {
	ioPacket.SetIoStatus(ioPackets.IosInProgress)
	if ioPacket.GetNodeDeviceType() != ch.GetNodeDeviceType() {
		ioPacket.SetIoStatus(ioPackets.IosInvalidNodeType)
		return
	}

	dev, ok := ch.devices[ioPacket.GetNodeIdentifier()]
	if !ok {
		ioPacket.SetIoStatus(ioPackets.IosDeviceIsNotAccessible)
		return
	}

	dev.StartIo(ioPacket)
}

func (ch *DiskChannel) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vDiskChannel connections\n", indent)
	for id := range ch.devices {
		_, _ = fmt.Fprintf(dest, "%v  %v\n", indent, id)
	}
}
