// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import (
	"io"
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
)

type DiskDevice interface {
	Dump(dest io.Writer, indent string)
	GetDiskGeometry() (hardware.BlockSize, hardware.BlockCount, hardware.TrackCount)
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware.NodeDeviceType
	GetNodeModelType() hardware.NodeModelType
	IsMounted() bool
	IsReady() bool
	IsWriteProtected() bool
	SetIsReady(flag bool)
	SetIsWriteProtected(flag bool)
	SetVerbose(flag bool)
	StartIo(pkt ioPackets.IoPacket)
}
