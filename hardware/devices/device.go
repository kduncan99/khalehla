// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import (
	"io"
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
)

// Device manages real or pseudo IO operations for a particular virtual device.
// It may do so synchronously or asynchronously
type Device interface {
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware.NodeDeviceType
	GetNodeModelType() hardware.NodeModelType
	IsMounted() bool
	IsReady() bool
	SetVerbose(flag bool)
	StartIo(ioPacket ioPackets.IoPacket)
}
