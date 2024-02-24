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
	ioFunction     IoFunction
	ioStatus       IoStatus
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

func (pkt *DiskIoPacket) GetIoFunction() IoFunction {
	return pkt.ioFunction
}

func (pkt *DiskIoPacket) GetIoStatus() IoStatus {
	return pkt.ioStatus
}

func (pkt *DiskIoPacket) GetString() string {
	funcStr, ok := IoFunctionTable[pkt.ioFunction]
	if !ok {
		funcStr = fmt.Sprintf("%v", pkt.ioFunction)
	}

	statStr, ok := IoStatusTable[pkt.ioStatus]
	if !ok {
		statStr = fmt.Sprintf("%v", pkt.ioStatus)
	}

	detStr := ""
	if pkt.ioFunction == IofRead || pkt.ioFunction == IofWrite {
		detStr += fmt.Sprintf("blkId:%v ", pkt.blockId)
	}
	if pkt.ioFunction == IofPrep {
		detStr += fmt.Sprintf("packId:%v prep:%v tracks:%v rem:%v ",
			pkt.packName, pkt.prepFactor, pkt.trackCount, pkt.removable)
	}
	if pkt.ioFunction == IofMount {
		detStr += fmt.Sprintf("file:%v writeProt:%v ", pkt.fileName, pkt.writeProtected)
	}

	return fmt.Sprintf("func:%s %sstat:%s", funcStr, detStr, statStr)
}

func (pkt *DiskIoPacket) SetIoStatus(ioStatus IoStatus) {
	pkt.ioStatus = ioStatus
}

func NewDiskIoPacketMount(nodeId types.NodeIdentifier, fileName string, writeProtected bool) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:         nodeId,
		ioFunction:     IofMount,
		ioStatus:       IosNotStarted,
		fileName:       fileName,
		writeProtected: writeProtected,
	}
}

func NewDiskIoPacketPrep(nodeId types.NodeIdentifier, packName string, prepFactor types.PrepFactor, trackCount types.TrackCount, removable bool) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: IofPrep,
		ioStatus:   IosNotStarted,
		packName:   packName,
		prepFactor: prepFactor,
		trackCount: trackCount,
		removable:  removable,
	}
}

func NewDiskIoPacketRead(nodeId types.NodeIdentifier, blockId types.BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: IofRead,
		ioStatus:   IosNotStarted,
		blockId:    blockId,
		buffer:     buffer,
	}
}

func NewDiskIoPacketReadLabel(nodeId types.NodeIdentifier, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: IofReadLabel,
		ioStatus:   IosNotStarted,
		buffer:     buffer,
	}
}

func NewDiskIoPacketReset(nodeId types.NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: IofReset,
		ioStatus:   IosNotStarted,
	}
}

func NewDiskIoPacketUnmount(nodeId types.NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: IofUnmount,
		ioStatus:   IosNotStarted,
	}
}

func NewDiskIoPacketWrite(nodeId types.NodeIdentifier, blockId types.BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: IofWrite,
		ioStatus:   IosNotStarted,
		blockId:    blockId,
		buffer:     buffer,
	}
}

func NewDiskIoPacketWriteLabel(nodeId types.NodeIdentifier, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		nodeId:     nodeId,
		ioFunction: IofWriteLabel,
		ioStatus:   IosNotStarted,
		buffer:     buffer,
	}
}
