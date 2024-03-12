// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

import (
	"fmt"
	"khalehla/hardware"
)

type DiskIoPacket struct {
	Listener       IoPacketListener
	NodeId         hardware.NodeIdentifier
	IoFunction     IoFunction
	IoStatus       IoStatus
	BlockId        hardware.BlockId    // for read, write
	Buffer         []byte              // for read, readLabel, write
	PackName       string              // for prep
	PrepFactor     hardware.PrepFactor // for prep
	TrackCount     hardware.TrackCount // for prep
	Removable      bool                // for prep
	Filename       string              // for mount
	WriteProtected bool                // for mount
}

func (pkt *DiskIoPacket) GetBlockId() hardware.BlockId {
	return pkt.BlockId
}

func (pkt *DiskIoPacket) GetListener() IoPacketListener {
	return pkt.Listener
}

func (pkt *DiskIoPacket) GetNodeIdentifier() hardware.NodeIdentifier {
	return pkt.NodeId
}

func (pkt *DiskIoPacket) GetNodeDeviceType() hardware.NodeDeviceType {
	return hardware.NodeDeviceDisk
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
		detStr += fmt.Sprintf("packId:%v prep:%v tracks:%v rem:%v ",
			pkt.PackName, pkt.PrepFactor, pkt.TrackCount, pkt.Removable)
	}
	if pkt.IoFunction == IofMount {
		detStr += fmt.Sprintf("file:%v writeProt:%v ", pkt.Filename, pkt.WriteProtected)
	}

	return fmt.Sprintf("func:%s %sstat:%s", funcStr, detStr, statStr)
}

func (pkt *DiskIoPacket) SetIoStatus(ioStatus IoStatus) {
	pkt.IoStatus = ioStatus
}

func NewDiskIoPacketMount(nodeId hardware.NodeIdentifier, fileName string, writeProtected bool) *DiskIoPacket {
	return &DiskIoPacket{
		NodeId:         nodeId,
		IoFunction:     IofMount,
		IoStatus:       IosNotStarted,
		Filename:       fileName,
		WriteProtected: writeProtected,
	}
}

func NewDiskIoPacketPrep(nodeId hardware.NodeIdentifier, packName string, prepFactor hardware.PrepFactor, trackCount hardware.TrackCount, removable bool) *DiskIoPacket {
	return &DiskIoPacket{
		NodeId:     nodeId,
		IoFunction: IofPrep,
		IoStatus:   IosNotStarted,
		PackName:   packName,
		PrepFactor: prepFactor,
		TrackCount: trackCount,
		Removable:  removable,
	}
}

func NewDiskIoPacketRead(nodeId hardware.NodeIdentifier, blockId hardware.BlockId, buffer []byte) *DiskIoPacket {
	return &DiskIoPacket{
		NodeId:     nodeId,
		IoFunction: IofRead,
		IoStatus:   IosNotStarted,
		BlockId:    blockId,
		Buffer:     buffer,
	}
}

func NewDiskIoPacketReadLabel(nodeId hardware.NodeIdentifier, buffer []byte) *DiskIoPacket {
	return &DiskIoPacket{
		NodeId:     nodeId,
		IoFunction: IofReadLabel,
		IoStatus:   IosNotStarted,
		Buffer:     buffer,
	}
}

func NewDiskIoPacketReset(nodeId hardware.NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		NodeId:     nodeId,
		IoFunction: IofReset,
		IoStatus:   IosNotStarted,
	}
}

func NewDiskIoPacketUnmount(nodeId hardware.NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		NodeId:     nodeId,
		IoFunction: IofUnmount,
		IoStatus:   IosNotStarted,
	}
}

func NewDiskIoPacketWrite(nodeId hardware.NodeIdentifier, blockId hardware.BlockId, buffer []byte) *DiskIoPacket {
	return &DiskIoPacket{
		NodeId:     nodeId,
		IoFunction: IofWrite,
		IoStatus:   IosNotStarted,
		BlockId:    blockId,
		Buffer:     buffer,
	}
}

func NewDiskIoPacketWriteLabel(nodeId hardware.NodeIdentifier, buffer []byte) *DiskIoPacket {
	return &DiskIoPacket{
		NodeId:     nodeId,
		IoFunction: IofWriteLabel,
		IoStatus:   IosNotStarted,
		Buffer:     buffer,
	}
}
