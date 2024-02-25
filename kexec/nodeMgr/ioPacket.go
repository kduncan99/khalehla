// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"khalehla/kexec"
)

// IoPacket contains all the information necessary for a Channel to route an IO operation,
// and for a device to perform that IO operation.
type IoPacket interface {
	GetNodeIdentifier() kexec.NodeIdentifier
	GetNodeDeviceType() NodeDeviceType
	GetIoFunction() IoFunction
	GetIoStatus() IoStatus
	GetString() string
	SetIoStatus(ioStatus IoStatus)
}
