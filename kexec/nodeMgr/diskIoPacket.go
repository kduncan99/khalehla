// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"fmt"
	"khalehla/kexec/types"
	"khalehla/pkg"
)

type DiskIoPacket struct {
	nodeId         types.NodeIdentifier
	ioFunction     types.IoFunction
	ioStatus       types.IoStatus
	blockId        types.BlockId    // for read, write
	buffer         []pkg.Word36     // for read, readLabel, write
	packName       string           // for prep
	prepFactor     types.PrepFactor // for prep
	trackCount     types.TrackCount // for prep
	removable      bool             // for prep
	fileName       string           // for mount
	writeProtected bool             // for mount
}

func (pkt *DiskIoPacket) GetNodeIdentifier() types.NodeIdentifier {
	return pkt.nodeId
}

func (pkt *DiskIoPacket) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceDisk
}

func (pkt *DiskIoPacket) GetIoFunction() types.IoFunction {
	return pkt.ioFunction
}

func (pkt *DiskIoPacket) GetIoStatus() types.IoStatus {
	return pkt.ioStatus
}

func (pkt *DiskIoPacket) GetString() string {
	funcStr, ok := types.IoFunctionTable[pkt.ioFunction]
	if !ok {
		funcStr = fmt.Sprintf("%v", pkt.ioFunction)
	}

	statStr, ok := types.IoStatusTable[pkt.ioStatus]
	if !ok {
		statStr = fmt.Sprintf("%v", pkt.ioStatus)
	}

	detStr := ""
	if pkt.ioFunction == types.IofRead || pkt.ioFunction == types.IofWrite {
		detStr += fmt.Sprintf("blkId:%v ", pkt.blockId)
	}
	if pkt.ioFunction == types.IofPrep {
		detStr += fmt.Sprintf("packId:%v prep:%v tracks:%v rem:%v ",
			pkt.packName, pkt.prepFactor, pkt.trackCount, pkt.removable)
	}
	if pkt.ioFunction == types.IofMount {
		detStr += fmt.Sprintf("file:%v writeProt:%v ", pkt.fileName, pkt.writeProtected)
	}

	return fmt.Sprintf("func:%s %sstat:%s", funcStr, detStr, statStr)
}

func (pkt *DiskIoPacket) SetIoStatus(ioStatus types.IoStatus) {
	pkt.ioStatus = ioStatus
}

func NewDiskIoPacketMount(nodeId types.NodeIdentifier, fileName string, writeProtected bool) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:         nodeId,
		ioFunction:     types.IofMount,
		ioStatus:       types.IosNotStarted,
		fileName:       fileName,
		writeProtected: writeProtected,
	}
}

func NewDiskIoPacketPrep(nodeId types.NodeIdentifier, packName string, prepFactor types.PrepFactor, trackCount types.TrackCount, removable bool) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: types.IofPrep,
		ioStatus:   types.IosNotStarted,
		packName:   packName,
		prepFactor: prepFactor,
		trackCount: trackCount,
		removable:  removable,
	}
}

func NewDiskIoPacketRead(nodeId types.NodeIdentifier, blockId types.BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: types.IofRead,
		ioStatus:   types.IosNotStarted,
		blockId:    blockId,
		buffer:     buffer,
	}
}

func NewDiskIoPacketReadLabel(nodeId types.NodeIdentifier, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: types.IofReadLabel,
		ioStatus:   types.IosNotStarted,
		buffer:     buffer,
	}
}

func NewDiskIoPacketReset(nodeId types.NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: types.IofReset,
		ioStatus:   types.IosNotStarted,
	}
}

func NewDiskIoPacketUnmount(nodeId types.NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: types.IofUnmount,
		ioStatus:   types.IosNotStarted,
	}
}

func NewDiskIoPacketWrite(nodeId types.NodeIdentifier, blockId types.BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: types.IofWrite,
		ioStatus:   types.IosNotStarted,
		blockId:    blockId,
		buffer:     buffer,
	}
}

func NewDiskIoPacketWriteLabel(nodeId types.NodeIdentifier, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: types.IofWriteLabel,
		ioStatus:   types.IosNotStarted,
		buffer:     buffer,
	}
}
