// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodes

import (
	"io"
)

type TapeDevice interface {
	Dump(dest io.Writer, indent string)
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
