// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import "khalehla/hardware"

type IDeviceListener interface {
	NotifyDeviceReady(nodeId hardware.NodeIdentifier, isReady bool)
}
