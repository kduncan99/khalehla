// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import "khalehla/pkg"

type DiskIoPacket struct {
	deviceIdentifier NodeIdentifier
	ioFunction       IoFunction
	ioStatus         IoStatus
	blockId          BlockId      // for read, write
	buffer           []pkg.Word36 // for read, readLabel, write
	packName         string       // for prep
	prepFactor       PrepFactor   // for prep
	trackCount       TrackCount   // for prep
	removable        bool         // for prep
	fileName         string       // for mount
	writeProtected   bool         // for mount
}

func (pkt *DiskIoPacket) GetDeviceIdentifier() NodeIdentifier {
	return pkt.deviceIdentifier
}

func (pkt *DiskIoPacket) GetNodeType() NodeType {
	return NodeTypeDisk
}

func (pkt *DiskIoPacket) GetIoFunction() IoFunction {
	return pkt.ioFunction
}

func (pkt *DiskIoPacket) GetIoStatus() IoStatus {
	return pkt.ioStatus
}

func (pkt *DiskIoPacket) SetIoStatus(ioStatus IoStatus) {
	pkt.ioStatus = ioStatus
}

func NewDiskIoPacketMount(fileName string, writeProtected bool) *DiskIoPacket {
	return &DiskIoPacket{
		ioFunction:     IofMount,
		ioStatus:       IosNotStarted,
		fileName:       fileName,
		writeProtected: writeProtected,
	}
}

func NewDiskIoPacketPrep(deviceIdentifier NodeIdentifier, packName string, prepFactor PrepFactor, trackCount TrackCount, removable bool) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       IofPrep,
		ioStatus:         IosNotStarted,
		packName:         packName,
		prepFactor:       prepFactor,
		trackCount:       trackCount,
		removable:        removable,
	}
}

func NewDiskIoPacketRead(deviceIdentifier NodeIdentifier, blockId BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       IofRead,
		ioStatus:         IosNotStarted,
		blockId:          blockId,
		buffer:           buffer,
	}
}

func NewDiskIoPacketReadLabel(deviceIdentifier NodeIdentifier, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       IofReadLabel,
		ioStatus:         IosNotStarted,
		buffer:           buffer,
	}
}

func NewDiskIoPacketReset(deviceIdentifier NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       IofReset,
		ioStatus:         IosNotStarted,
	}
}

func NewDiskIoPacketUnmount(deviceIdentifier NodeIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       IofUnmount,
		ioStatus:         IosNotStarted,
	}
}

func NewDiskIoPacketWrite(deviceIdentifier NodeIdentifier, blockId BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       IofWrite,
		ioStatus:         IosNotStarted,
		blockId:          blockId,
		buffer:           buffer,
	}
}
