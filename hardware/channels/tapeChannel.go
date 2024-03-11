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

// TapeChannel routes IOs to the appropriate deviceInfos which it manages.
// Some day in the future we may add caching, perhaps in a CacheTapeChannel.
type TapeChannel struct {
	devices map[hardware.NodeIdentifier]devices.TapeDevice
}

func NewTapeChannel() *TapeChannel {
	return &TapeChannel{
		devices: make(map[hardware.NodeIdentifier]devices.TapeDevice),
	}
}

func (ch *TapeChannel) GetNodeCategoryType() hardware.NodeCategoryType {
	return hardware.NodeCategoryChannel
}

func (ch *TapeChannel) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceTape
}

func (ch *TapeChannel) GetNodeModelType() hardware.NodeModelType {
	return hardware.NodeModelTapeLibraryChannel
}

func (ch *TapeChannel) AssignDevice(nodeIdentifier hardware.NodeIdentifier, device devices.Device) error {
	if device.GetNodeDeviceType() != hardware.NodeDeviceTape {
		return fmt.Errorf("device is not a tape")
	}

	ch.devices[nodeIdentifier] = device.(devices.TapeDevice)
	return nil
}

func (ch *TapeChannel) StartIo(ioPacket ioPackets.IoPacket) {
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

func (ch *TapeChannel) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vTapeChannel connections\n", indent)
	for id := range ch.devices {
		_, _ = fmt.Fprintf(dest, "%v  %v\n", indent, id)
	}
}
