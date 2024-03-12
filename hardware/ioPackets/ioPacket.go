// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

type PacketType uint

const (
	_ PacketType = iota
	DiskPacketType
	TapePacketType
)

// IoPacketListener should be implemented by any IO caller which wants to be
// notified when an IO is complete.
type IoPacketListener interface {
	IoComplete(ioPacket IoPacket)
}

// IoMountInfo contains information necessary for mounting file-system based media
type IoMountInfo struct {
	Filename     string
	WriteProtect bool
}

// IoPacket contains all the information necessary for a Channel to route an IO operation,
// and for a device to perform that IO operation.
type IoPacket interface {
	GetPacketType() PacketType
	GetIoFunction() IoFunction
	GetIoStatus() IoStatus
	GetString() string
	GetListener() IoPacketListener
	SetIoStatus(ioStatus IoStatus)
}
