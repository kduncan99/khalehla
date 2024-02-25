// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"io"
	"khalehla/kexec"
)

// TapeChannel routes IOs to the appropriate deviceInfos which it manages.
// Some day in the future we may add caching, perhaps in a CacheTapeChannel.
type TapeChannel struct {
	devices map[kexec.NodeIdentifier]TapeDevice
}

func NewTapeChannel() *TapeChannel {
	return &TapeChannel{
		devices: make(map[kexec.NodeIdentifier]TapeDevice),
	}
}

func (ch *TapeChannel) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryChannel
}

func (ch *TapeChannel) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceTape
}

func (ch *TapeChannel) GetNodeModelType() NodeModelType {
	return NodeModelTapeLibraryChannel
}

func (ch *TapeChannel) AssignDevice(nodeIdentifier kexec.NodeIdentifier, device Device) error {
	if device.GetNodeDeviceType() != NodeDeviceTape {
		return fmt.Errorf("device is not a tape")
	}

	ch.devices[nodeIdentifier] = device.(TapeDevice)
	return nil
}

func (ch *TapeChannel) StartIo(ioPacket IoPacket) {
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

func (ch *TapeChannel) Dump(destination io.Writer, indent string) {
	// TODO
}
