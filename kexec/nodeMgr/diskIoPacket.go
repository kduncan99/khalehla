// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"khalehla/kexec/types"
	"khalehla/pkg"
)

type DiskIoPacket struct {
	deviceId       types.DeviceIdentifier
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

func (pkt *DiskIoPacket) GetDeviceIdentifier() types.DeviceIdentifier {
	return pkt.deviceId
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

func (pkt *DiskIoPacket) SetIoStatus(ioStatus types.IoStatus) {
	pkt.ioStatus = ioStatus
}

func NewDiskIoPacketMount(deviceId types.DeviceIdentifier, fileName string, writeProtected bool) *DiskIoPacket {
	return &DiskIoPacket{
		deviceId:       deviceId,
		ioFunction:     types.IofMount,
		ioStatus:       types.IosNotStarted,
		fileName:       fileName,
		writeProtected: writeProtected,
	}
}

func NewDiskIoPacketPrep(deviceId types.DeviceIdentifier, packName string, prepFactor types.PrepFactor, trackCount types.TrackCount, removable bool) *DiskIoPacket {
	return &DiskIoPacket{
		deviceId:   deviceId,
		ioFunction: types.IofPrep,
		ioStatus:   types.IosNotStarted,
		packName:   packName,
		prepFactor: prepFactor,
		trackCount: trackCount,
		removable:  removable,
	}
}

func NewDiskIoPacketRead(deviceId types.DeviceIdentifier, blockId types.BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceId:   deviceId,
		ioFunction: types.IofRead,
		ioStatus:   types.IosNotStarted,
		blockId:    blockId,
		buffer:     buffer,
	}
}

func NewDiskIoPacketReadLabel(deviceId types.DeviceIdentifier, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceId:   deviceId,
		ioFunction: types.IofReadLabel,
		ioStatus:   types.IosNotStarted,
		buffer:     buffer,
	}
}

func NewDiskIoPacketReset(deviceId types.DeviceIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		deviceId:   deviceId,
		ioFunction: types.IofReset,
		ioStatus:   types.IosNotStarted,
	}
}

func NewDiskIoPacketUnmount(deviceId types.DeviceIdentifier) *DiskIoPacket {
	return &DiskIoPacket{
		deviceId:   deviceId,
		ioFunction: types.IofUnmount,
		ioStatus:   types.IosNotStarted,
	}
}

func NewDiskIoPacketWrite(deviceId types.DeviceIdentifier, blockId types.BlockId, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceId:   deviceId,
		ioFunction: types.IofWrite,
		ioStatus:   types.IosNotStarted,
		blockId:    blockId,
		buffer:     buffer,
	}
}

func NewDiskIoPacketWriteLabel(deviceId types.DeviceIdentifier, buffer []pkg.Word36) *DiskIoPacket {
	return &DiskIoPacket{
		deviceId:   deviceId,
		ioFunction: types.IofWriteLabel,
		ioStatus:   types.IosNotStarted,
		buffer:     buffer,
	}
}
