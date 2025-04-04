// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec/facilitiesMgr"
	"khalehla/logger"
)

func handleFree(pkt *handlerPacket) (*facilitiesMgr.FacStatusResult, uint64) {
	logger.LogTraceF("CSI", "handleFree:%v", *pkt.pcs)

	/*
		@FREE[,options] filename
		options include
			A, B, D, I, R, S, X
	*/
	// TODO implement
	return nil, 0
}
