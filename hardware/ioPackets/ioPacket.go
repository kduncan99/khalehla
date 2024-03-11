// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

import (
	"khalehla/hardware"
)

// IoPacketListener should be implemented by any IO caller which wants to be
// notified when an IO is complete.
type IoPacketListener interface {
	IoComplete(ioPacket IoPacket)
}

// IoPacket contains all the information necessary for a Channel to route an IO operation,
// and for a device to perform that IO operation.
type IoPacket interface {
	GetNodeIdentifier() hardware.NodeIdentifier
	GetNodeDeviceType() hardware.NodeDeviceType
	GetIoFunction() IoFunction
	GetIoStatus() IoStatus
	GetString() string
	GetListener() IoPacketListener
	SetIoStatus(ioStatus IoStatus)
}
