// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"khalehla/kexec/types"
)

type TapeIoPacket struct {
	nodeId         types.NodeIdentifier
	ioFunction     types.IoFunction
	ioStatus       types.IoStatus
	fileName       string // for mount
	writeProtected bool   // for mount
}

func (pkt *TapeIoPacket) GetNodeIdentifier() types.NodeIdentifier {
	return pkt.nodeId
}

func (pkt *TapeIoPacket) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceTape
}

func (pkt *TapeIoPacket) GetIoFunction() types.IoFunction {
	return pkt.ioFunction
}

func (pkt *TapeIoPacket) GetIoStatus() types.IoStatus {
	return pkt.ioStatus
}

func (pkt *TapeIoPacket) GetString() string {
	funcStr, ok := types.IoFunctionTable[pkt.ioFunction]
	if !ok {
		funcStr = fmt.Sprintf("%v", pkt.ioFunction)
	}

	statStr, ok := types.IoStatusTable[pkt.ioStatus]
	if !ok {
		statStr = fmt.Sprintf("%v", pkt.ioStatus)
	}

	detStr := ""
	// TODO detStr

	return fmt.Sprintf("func:%s %sstat:%s", funcStr, detStr, statStr)
}

func (pkt *TapeIoPacket) SetIoStatus(ioStatus types.IoStatus) {
	pkt.ioStatus = ioStatus
}

func NewTapeIoPacketMount(nodeId types.NodeIdentifier, fileName string, writeProtected bool) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:         nodeId,
		ioFunction:     types.IofMount,
		ioStatus:       types.IosNotStarted,
		fileName:       fileName,
		writeProtected: writeProtected,
	}
}

func NewTapeIoPacketReset(nodeId types.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:     nodeId,
		ioFunction: types.IofReset,
		ioStatus:   types.IosNotStarted,
	}
}

func NewTapeIoPacketUnmount(nodeId types.NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		nodeId:     nodeId,
		ioFunction: types.IofUnmount,
		ioStatus:   types.IosNotStarted,
	}
}
