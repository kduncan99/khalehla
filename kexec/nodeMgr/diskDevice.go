// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"khalehla/kexec/types"
)

type DiskDevice interface {
	Dump(dest io.Writer, indent string)
	GetGeometry() *types.DiskPackGeometry
	GetNodeCategoryType() NodeCategoryType
	GetNodeDeviceType() NodeDeviceType
	GetNodeModelType() NodeModelType
	IsMounted() bool
	IsPrepped() bool
	IsReady() bool
	IsWriteProtected() bool
	SetIsReady(flag bool)
	SetIsWriteProtected(flag bool)
	StartIo(pkt IoPacket)
}
