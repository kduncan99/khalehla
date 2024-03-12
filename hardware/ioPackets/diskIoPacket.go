// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

import (
	"fmt"
	"khalehla/hardware"
)

// DiskIoPacket
// Buffer has a couple of overload usages
type DiskIoPacket struct {
	Listener   IoPacketListener
	IoFunction IoFunction
	IoStatus   IoStatus
	BlockId    hardware.BlockId // for IofRead, IofWrite
	Buffer     []byte           // for IofRead, IofWrite
	MountInfo  *IoMountInfo     // for IofMount
	PrepInfo   *IoPrepInfo      // for IofPrep
}

type IoPrepInfo struct {
	PrepFactor  hardware.PrepFactor
	TrackCount  hardware.TrackCount
	PackName    string
	IsRemovable bool
}

func (pkt *DiskIoPacket) GetBlockId() hardware.BlockId {
	return pkt.BlockId
}

func (pkt *DiskIoPacket) GetListener() IoPacketListener {
	return pkt.Listener
}

func (pkt *DiskIoPacket) GetPacketType() PacketType {
	return DiskPacketType
}

func (pkt *DiskIoPacket) GetIoFunction() IoFunction {
	return pkt.IoFunction
}

func (pkt *DiskIoPacket) GetIoStatus() IoStatus {
	return pkt.IoStatus
}

func (pkt *DiskIoPacket) GetString() string {
	funcStr, ok := IoFunctionTable[pkt.IoFunction]
	if !ok {
		funcStr = fmt.Sprintf("%v", pkt.IoFunction)
	}

	statStr, ok := IoStatusTable[pkt.IoStatus]
	if !ok {
		statStr = fmt.Sprintf("%v", pkt.IoStatus)
	}

	detStr := ""
	if pkt.IoFunction == IofRead || pkt.IoFunction == IofWrite {
		detStr += fmt.Sprintf("blkId:%v ", pkt.BlockId)
	}
	if pkt.IoFunction == IofPrep {
		detStr += fmt.Sprintf("prep:%v tracks:%v packName:%v rem:%v",
			pkt.PrepInfo.PrepFactor, pkt.PrepInfo.TrackCount, pkt.PrepInfo.PackName, pkt.PrepInfo.IsRemovable)
	}
	if pkt.IoFunction == IofMount {
		detStr += fmt.Sprintf("file:%v writeProt:%v ", pkt.MountInfo.Filename, pkt.MountInfo.WriteProtect)
	}

	return fmt.Sprintf("func:%s %sstat:%s", funcStr, detStr, statStr)
}

func (pkt *DiskIoPacket) SetIoStatus(ioStatus IoStatus) {
	pkt.IoStatus = ioStatus
}
