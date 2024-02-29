// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"khalehla/kexec"
)

type TapeDevice interface {
	Dump(dest io.Writer, indent string)
	GetNodeCategoryType() kexec.NodeCategoryType
	GetNodeDeviceType() kexec.NodeDeviceType
	GetNodeModelType() NodeModelType
	IsMounted() bool
	IsReady() bool
	IsWriteProtected() bool
	SetIsReady(flag bool)
	SetIsWriteProtected(flag bool)
	StartIo(pkt IoPacket)
}
