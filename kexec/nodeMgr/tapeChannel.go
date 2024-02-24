// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

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

func (ch *TapeChannel) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryChannel
}

func (ch *TapeChannel) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceTape
}

func (ch *TapeChannel) GetNodeModelType() NodeModelType {
	return NodeModelTapeChannel
}

func (ch *TapeChannel) AssignDevice(deviceIdentifier types.DeviceIdentifier, device Device) error {
	if device.GetNodeDeviceType() != NodeDeviceTape {
		return fmt.Errorf("device is not a tape")
	}

	ch.devices[deviceIdentifier] = device.(*TapeDevice)
	return nil
}

func (ch *TapeChannel) StartIo(ioPacket IoPacket) {
	ioPacket.SetIoStatus(types.IosInProgress)
	if ioPacket.GetNodeDeviceType() != ch.GetNodeDeviceType() {
		ioPacket.SetIoStatus(types.IosInvalidNodeType)
		return
	}

	dev, ok := ch.devices[ioPacket.GetDeviceIdentifier()]
	if !ok {
		ioPacket.SetIoStatus(types.IosDeviceIsNotAccessible)
		return
	}

	dev.StartIo(ioPacket)
}

func (ch *TapeChannel) Dump(destination io.Writer, indent string) {
	// TODO
}
