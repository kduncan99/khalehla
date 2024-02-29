// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodes

import (
	"io"
	"khalehla/kexec"
)

// Channel manages async communication with the various deviceInfos assigned to it.
// It may also manage caching, automatic mounting, or any other various activities
// on behalf of the exec.
type Channel interface {
	AssignDevice(nodeIdentifier kexec.NodeIdentifier, device Device) error
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() NodeCategoryType
	GetNodeDeviceType() NodeDeviceType
	GetNodeModelType() NodeModelType
	StartIo(ioPacket IoPacket)
}
