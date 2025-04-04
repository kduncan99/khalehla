// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import (
	"io"

	"khalehla/hardware"
	hardware2 "khalehla/old/hardware"
	"khalehla/old/hardware/ioPackets"
)

type TapeDevice interface {
	Dump(dest io.Writer, indent string)
	GetBlocksExtended() int
	GetFilesExtended() uint
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware2.NodeDeviceType
	GetNodeIdentifier() hardware.NodeIdentifier
	GetNodeModelType() hardware2.NodeModelType
	IsAtLoadPoint() bool
	IsMounted() bool
	IsReady() bool
	IsWriteProtected() bool
	Reset()
	SetIsReady(flag bool)
	SetIsWriteProtected(flag bool)
	SetVerbose(flag bool)
	StartIo(pkt ioPackets.IoPacket)
}
