// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import "khalehla/kexec/types"

type TapeIoPacket struct {
	deviceId       types.DeviceIdentifier
	ioFunction     types.IoFunction
	ioStatus       types.IoStatus
	fileName       string // for mount
	writeProtected bool   // for mount
}

func (pkt *TapeIoPacket) GetDeviceIdentifier() types.DeviceIdentifier {
	return pkt.deviceId
}

func (pkt *TapeIoPacket) GetNodeType() types.NodeType {
	return types.NodeTypeTape
}

func (pkt *TapeIoPacket) GetIoFunction() types.IoFunction {
	return pkt.ioFunction
}

func (pkt *TapeIoPacket) GetIoStatus() types.IoStatus {
	return pkt.ioStatus
}

func (pkt *TapeIoPacket) SetIoStatus(ioStatus types.IoStatus) {
	pkt.ioStatus = ioStatus
}

func NewTapeIoPacketMount(deviceId types.DeviceIdentifier, fileName string, writeProtected bool) *TapeIoPacket {
	return &TapeIoPacket{
		deviceId:       deviceId,
		ioFunction:     types.IofMount,
		ioStatus:       types.IosNotStarted,
		fileName:       fileName,
		writeProtected: writeProtected,
	}
}

func NewTapeIoPacketReset(deviceId types.DeviceIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		deviceId:   deviceId,
		ioFunction: types.IofReset,
		ioStatus:   types.IosNotStarted,
	}
}

func NewTapeIoPacketUnmount(deviceId types.DeviceIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		deviceId:   deviceId,
		ioFunction: types.IofUnmount,
		ioStatus:   types.IosNotStarted,
	}
}
