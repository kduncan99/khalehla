// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package hardware

import (
	"io"
	"khalehla/pkg"
)

// NodeIdentifier uniquely identifies a particular device or channel (or anything else identifiable which we manage)
// It is currently implemented as the 1-6 character device Name, all caps alphas and/or digits LJSF
// stored as Fieldata in a Word36 struct
type NodeIdentifier pkg.Word36

type Node interface {
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() NodeCategoryType
	GetNodeDeviceType() NodeDeviceType
	GetNodeModelType() NodeModelType
	IsReady() bool
}
