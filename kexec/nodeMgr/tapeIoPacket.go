// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/pkg"
)

type TapeIoPacket struct {
	nodeId         kexec.NodeIdentifier
	ioFunction     IoFunction
	ioStatus       IoStatus
	buffer         []pkg.Word36
	fileName       string // for mount
	writeProtected bool   // for mount
}

func (pkt *TapeIoPacket) GetBuffer() []pkg.Word36 {
	return pkt.buffer
}

func (pkt *TapeIoPacket) GetNodeIdentifier() kexec.NodeIdentifier {
	return pkt.nodeId
}

func (pkt *TapeIoPacket) GetNodeDeviceType() kexec.NodeDeviceType {
	return kexec.NodeDeviceTape
}

func (pkt *TapeIoPacket) GetIoFunction() IoFunction {
	return pkt.ioFunction
}

func (pkt *TapeIoPacket) GetIoStatus() IoStatus {
	return pkt.ioStatus
}

func (pkt *TapeIoPacket) GetString() string {
	funcStr, ok := IoFunctionTable[pkt.ioFunction]
	if !ok {
		funcStr = fmt.Sprintf("%v", pkt.ioFunction)
	}

	statStr, ok := IoStatusTable[pkt.ioStatus]
	if !ok {
		statStr = fmt.Sprintf("%v", pkt.ioStatus)
	}

	detStr := ""
	// TODO detStr

	return fmt.Sprintf("func:%s %sstat:%s", funcStr, detStr, statStr)
}

func (pkt *TapeIoPacket) SetIoStatus(ioStatus IoStatus) {
	pkt.ioStatus = ioStatus
}

func NewTapeIoPacketMount(nodeId kexec.NodeIdentifier, fileName string, writeProtected bool) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:         nodeId,
		ioFunction:     IofMount,
		ioStatus:       IosNotStarted,
		fileName:       fileName,
		writeProtected: writeProtected,
	}
}

func NewTapeIoPacketRead(nodeId kexec.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:     nodeId,
		ioFunction: IofRead,
		ioStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketReset(nodeId kexec.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:     nodeId,
		ioFunction: IofReset,
		ioStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketRewind(nodeId kexec.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:     nodeId,
		ioFunction: IofRewind,
		ioStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketRewindAndUnload(nodeId kexec.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:     nodeId,
		ioFunction: IofRewindAndUnload,
		ioStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketUnmount(nodeId kexec.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:     nodeId,
		ioFunction: IofUnmount,
		ioStatus:   IosNotStarted,
	}
}

func NewTapeIoPacketWrite(nodeId kexec.NodeIdentifier, buffer []pkg.Word36) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:     nodeId,
		ioFunction: IofWrite,
		ioStatus:   IosNotStarted,
		buffer:     buffer,
	}
}

func NewTapeIoPacketWriteTapeMark(nodeId kexec.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:     nodeId,
		ioFunction: IofWriteTapeMark,
		ioStatus:   IosNotStarted,
	}
}
