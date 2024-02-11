// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import "khalehla/kexec/types"

type TapeIoPacket struct {
	deviceIdentifier types.NodeIdentifier
	ioFunction       types.IoFunction
	ioStatus         types.IoStatus
	fileName         string // for mount
	writeProtected   bool   // for mount
}

func (pkt *TapeIoPacket) GetDeviceIdentifier() types.NodeIdentifier {
	return pkt.deviceIdentifier
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

func NewTapeIoPacketMount(deviceIdentifier types.NodeIdentifier, fileName string, writeProtected bool) *TapeIoPacket {
	return &TapeIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofMount,
		ioStatus:         types.IosNotStarted,
		fileName:         fileName,
		writeProtected:   writeProtected,
	}
}

func NewTapeIoPacketReset(deviceIdentifier types.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofReset,
		ioStatus:         types.IosNotStarted,
	}
}

func NewTapeIoPacketUnmount(deviceIdentifier types.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofUnmount,
		ioStatus:         types.IosNotStarted,
	}
}
