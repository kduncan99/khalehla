// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

import (
	"fmt"
	"khalehla/hardware"
)

type TapeIoPacket struct {
	Listener       IoPacketListener
	NodeId         hardware.NodeIdentifier
	IoFunction     IoFunction
	IoStatus       IoStatus
	Buffer         []byte // provided by caller on write, by tape device on read
	PayloadLength  uint32 // bytes to be written, or bytes read (<= buffer length)
	Filename       string // for mount
	WriteProtected bool   // for mount
}

func (pkt *TapeIoPacket) GetListener() IoPacketListener {
	return pkt.Listener
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

func NewTapeIoPacketMount(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
	fileName string,
	writeProtected bool,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:       listener,
		NodeId:         nodeId,
		IoFunction:     IofMount,
		IoStatus:       IosNotStarted,
		Filename:       fileName,
		WriteProtected: writeProtected,
	}
}

func NewTapeIoPacketMoveBackward(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofMoveBackward,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketMoveForward(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofMoveForward,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketRead(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofRead,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketReadBackward(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofReadBackward,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketReset(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofReset,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketRewind(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofRewind,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketRewindAndUnload(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofRewindAndUnload,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketUnmount(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofUnmount,
		IoStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketWrite(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
	buffer []byte,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofWrite,
		IoStatus:   IosNotStarted,
		Buffer:     buffer,
	}
}

func NewTapeIoPacketWriteTapeMark(
	listener IoPacketListener,
	nodeId hardware.NodeIdentifier,
) *TapeIoPacket {
	return &TapeIoPacket{
		Listener:   listener,
		NodeId:     nodeId,
		IoFunction: IofWriteTapeMark,
		IoStatus:   IosNotStarted,
	}
}
