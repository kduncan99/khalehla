// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

import (
	"fmt"
	"khalehla/hardware"
	"khalehla/pkg"
)

type TapeIoPacket struct {
	NodeId         hardware.NodeIdentifier
	IoFunction     IoFunction
	IoStatus       IoStatus
	Buffer         []pkg.Word36
	Filename       string // for mount
	WriteProtected bool   // for mount
}

// TODO Buffer needs to be []byte

func (pkt *TapeIoPacket) GetBuffer() []pkg.Word36 {
	return pkt.Buffer
}

func (pkt *TapeIoPacket) GetNodeIdentifier() hardware.NodeIdentifier {
	return pkt.NodeId
}

func (pkt *TapeIoPacket) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceTape
}

func (pkt *TapeIoPacket) GetIoFunction() IoFunction {
	return pkt.IoFunction
}

func (pkt *TapeIoPacket) GetIoStatus() IoStatus {
	return pkt.IoStatus
}

func (pkt *TapeIoPacket) GetString() string {
	funcStr, ok := IoFunctionTable[pkt.IoFunction]
	if !ok {
		funcStr = fmt.Sprintf("%v", pkt.IoFunction)
	}

	statStr, ok := IoStatusTable[pkt.IoStatus]
	if !ok {
		statStr = fmt.Sprintf("%v", pkt.IoStatus)
	}

	detStr := ""
	// TODO construct detail string

	return fmt.Sprintf("func:%s %sstat:%s", funcStr, detStr, statStr)
}

func (pkt *TapeIoPacket) SetIoStatus(ioStatus IoStatus) {
	pkt.IoStatus = ioStatus
}

// TODO Need NewTapeIoPacket**** for MoveBack, MoveFwd, ReadBack

func NewTapeIoPacketMount(nodeId hardware.NodeIdentifier, fileName string, writeProtected bool) *TapeIoPacket {
	return &TapeIoPacket{
		NodeId:         nodeId,
		IoFunction:     IofMount,
		IoStatus:       IosNotStarted,
		Filename:       fileName,
		WriteProtected: writeProtected,
	}
}

func NewTapeIoPacketRead(nodeId hardware.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		NodeId:     nodeId,
		IoFunction: IofRead,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketReset(nodeId hardware.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		NodeId:     nodeId,
		IoFunction: IofReset,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketRewind(nodeId hardware.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		NodeId:     nodeId,
		IoFunction: IofRewind,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketRewindAndUnload(nodeId hardware.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		NodeId:     nodeId,
		IoFunction: IofRewindAndUnload,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketUnmount(nodeId hardware.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		NodeId:     nodeId,
		IoFunction: IofUnmount,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketWrite(nodeId hardware.NodeIdentifier, buffer []pkg.Word36) *TapeIoPacket {
	return &TapeIoPacket{
		NodeId:     nodeId,
		IoFunction: IofWrite,
		IoStatus:   IosNotStarted,
		Buffer:     buffer,
	}
}

func NewTapeIoPacketWriteTapeMark(nodeId hardware.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		NodeId:     nodeId,
		IoFunction: IofWriteTapeMark,
		IoStatus:   IosNotStarted,
	}
}
