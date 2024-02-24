// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
)

// Device manages real or pseudo IO operations for a particular virtual device.
// It may do so synchronously or asynchronously
type Device interface {
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() NodeCategoryType
	GetNodeDeviceType() NodeDeviceType
	GetNodeModelType() NodeModelType
	IsMounted() bool
	IsReady() bool
	StartIo(ioPacket IoPacket)
}
