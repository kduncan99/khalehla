// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

import (
	"fmt"
)

type TapeIoPacket struct {
	Listener   IoPacketListener
	IoFunction IoFunction
	IoStatus   IoStatus
	Buffer     []byte       // provided by caller on IofWrite, by tape device on IofRead, IofReadBackward
	DataLength uint32       // IofWrite bytes to be written, or IofRead, IofReadBackward bytes read (<= buffer length)
	MountInfo  *IoMountInfo // for IofMount
}

func (pkt *TapeIoPacket) GetListener() IoPacketListener {
	return pkt.Listener
}

func (pkt *TapeIoPacket) GetPacketType() PacketType {
	return TapePacketType
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
	if pkt.IoFunction == IofWrite {
		detStr += fmt.Sprintf("bytes:%v ", pkt.DataLength)
	}
	if pkt.IoFunction == IofMount {
		if pkt.MountInfo == nil {
			detStr += "no MountInfo "
		} else {
			detStr += fmt.Sprintf("file:%v writeProt:%v ", pkt.MountInfo.Filename, pkt.MountInfo.WriteProtect)
		}
	}

	return fmt.Sprintf("func:%s %sstat:%s", funcStr, detStr, statStr)
}

func (pkt *TapeIoPacket) SetIoStatus(ioStatus IoStatus) {
	pkt.IoStatus = ioStatus
}
