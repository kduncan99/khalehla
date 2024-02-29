// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"khalehla/kexec"
)

type Node interface {
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() kexec.NodeCategoryType
	GetNodeDeviceType() kexec.NodeDeviceType
	GetNodeModelType() NodeModelType
	IsReady() bool
}
