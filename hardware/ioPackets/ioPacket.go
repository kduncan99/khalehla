// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

import (
	"khalehla/hardware"
)

// IoPacket contains all the information necessary for a Channel to route an IO operation,
// and for a device to perform that IO operation.
type IoPacket interface {
	GetNodeIdentifier() hardware.NodeIdentifier
	GetNodeDeviceType() hardware.NodeDeviceType
	GetIoFunction() IoFunction
	GetIoStatus() IoStatus
	GetString() string
	SetIoStatus(ioStatus IoStatus)
}
