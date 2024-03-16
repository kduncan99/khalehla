// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package hardware

import (
	"io"
	"sync"
)

var nextNodeIdentifier NodeIdentifier = 1
var mutex sync.Mutex

type Node interface {
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() NodeCategoryType
	GetNodeDeviceType() NodeDeviceType
	GetNodeIdentifier() NodeIdentifier
	GetNodeModelType() NodeModelType
	IsReady() bool
	Reset()
}

func GetNextNodeIdentifier() (ni NodeIdentifier) {
	mutex.Lock()
	ni = nextNodeIdentifier
	nextNodeIdentifier++
	mutex.Unlock()
	return
}
