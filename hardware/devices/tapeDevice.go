// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import (
	"io"
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
)

type TapeDevice interface {
	Dump(dest io.Writer, indent string)
	GetBlocksExtended() int
	GetFilesExtended() uint
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware.NodeDeviceType
	GetNodeModelType() hardware.NodeModelType
	IsAtLoadPoint() bool
	IsMounted() bool
	IsReady() bool
	IsWriteProtected() bool
	SetIsReady(flag bool)
	SetIsWriteProtected(flag bool)
	SetVerbose(flag bool)
	StartIo(pkt ioPackets.IoPacket)
}
