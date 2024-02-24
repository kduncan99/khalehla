// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import "khalehla/kexec/types"

// IoPacket contains all the information necessary for a Channel to route an IO operation,
// and for a device to perform that IO operation.
type IoPacket interface {
	GetNodeIdentifier() types.NodeIdentifier
	GetNodeDeviceType() NodeDeviceType
	GetIoFunction() types.IoFunction
	GetIoStatus() types.IoStatus
	GetString() string
	SetIoStatus(ioStatus types.IoStatus)
}
