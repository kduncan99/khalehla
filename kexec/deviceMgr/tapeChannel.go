// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
)

// TapeChannel routes IOs to the appropriate deviceInfos which it manages.
// Some day in the future we may add caching, perhaps in a CacheTapeChannel.
type TapeChannel struct {
	devices map[types.DeviceIdentifier]*TapeDevice
}

func NewTapeChannel() *TapeChannel {
	return &TapeChannel{
		devices: make(map[types.DeviceIdentifier]*TapeDevice),
	}
}

func (ch *TapeChannel) GetNodeType() types.NodeType {
	return types.NodeTypeTape
}

func (ch *TapeChannel) AssignDevice(deviceIdentifier types.DeviceIdentifier, device types.Device) error {
	if device.GetNodeType() != types.NodeTypeTape {
		return fmt.Errorf("device is not a tape")
	}

	ch.devices[deviceIdentifier] = device.(*TapeDevice)
	return nil
}

func (ch *TapeChannel) StartIo(ioPacket types.IoPacket) {
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

func (ch *TapeChannel) Dump(destination io.Writer, indent string) {
	// TODO
}
