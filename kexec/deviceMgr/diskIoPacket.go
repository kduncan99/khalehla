// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import (
	"khalehla/kexec/types"
	"khalehla/pkg"
)

type DiskIoPacket struct {
	deviceIdentifier types.NodeIdentifier
	ioFunction       types.IoFunction
	ioStatus         types.IoStatus
	blockId          types.BlockId    // for read, write
	buffer           []pkg.Word36     // for read, readLabel, write
	packName         string           // for prep
	prepFactor       types.PrepFactor // for prep
	trackCount       types.TrackCount // for prep
	removable        bool             // for prep
	fileName         string           // for mount
	writeProtected   bool             // for mount
}

func (pkt *DiskIoPacket) GetDeviceIdentifier() types.NodeIdentifier {
	return pkt.deviceIdentifier
}

func (pkt *DiskIoPacket) GetNodeType() types.NodeType {
	return types.NodeTypeDisk
}

func (pkt *DiskIoPacket) GetIoFunction() types.IoFunction {
	return pkt.ioFunction
}

func (pkt *DiskIoPacket) GetIoStatus() types.IoStatus {
	return pkt.ioStatus
}

func (pkt *DiskIoPacket) SetIoStatus(ioStatus types.IoStatus) {
	pkt.ioStatus = ioStatus
}

func NewDiskIoPacketMount(deviceIdentifier types.NodeIdentifier, fileName string, writeProtected bool) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofMount,
		ioStatus:         types.IosNotStarted,
		fileName:         fileName,
		writeProtected:   writeProtected,
	}
}

func NewDiskIoPacketPrep(deviceIdentifier types.NodeIdentifier, packName string, prepFactor types.PrepFactor, trackCount types.TrackCount, removable bool) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofPrep,
		ioStatus:         types.IosNotStarted,
		packName:         packName,
		prepFactor:       prepFactor,
		trackCount:       trackCount,
		removable:        removable,
	}
}

func NewDiskIoPacketRead(deviceIdentifier types.NodeIdentifier, blockId types.BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofRead,
		ioStatus:         types.IosNotStarted,
		blockId:          blockId,
		buffer:           buffer,
	}
}

func NewDiskIoPacketReadLabel(deviceIdentifier types.NodeIdentifier, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofReadLabel,
		ioStatus:         types.IosNotStarted,
		buffer:           buffer,
	}
}

func NewDiskIoPacketReset(deviceIdentifier types.NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofReset,
		ioStatus:         types.IosNotStarted,
	}
}

func NewDiskIoPacketUnmount(deviceIdentifier types.NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofUnmount,
		ioStatus:         types.IosNotStarted,
	}
}

func NewDiskIoPacketWrite(deviceIdentifier types.NodeIdentifier, blockId types.BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       types.IofWrite,
		ioStatus:         types.IosNotStarted,
		blockId:          blockId,
		buffer:           buffer,
	}
}
