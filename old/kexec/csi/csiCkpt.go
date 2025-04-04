// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec/facilitiesMgr"
	"khalehla/logger"
)

func handleCkpt(pkt *handlerPacket) (*facilitiesMgr.FacStatusResult, uint64) {
	logger.LogTraceF("CSI", "handleCkpt:%v", *pkt.pcs)

	/*
		@CKPT[,options] [filename.element,,eqpmnt-type,reel-1/.../reel-n,,,,CTL-pool]
		options include
			A, B, C, J, M, N, P, Q, R, S, T, Z
	*/
	// TODO implement
	return nil, 0
}
