// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec/facilitiesMgr"
)

func handleRstrt(pkt *handlerPacket) (*facilitiesMgr.FacStatusResult, uint64) {

	/*
		@RSTRT[,scheduling-priority/options/processor-dispatching-priority]
			run-id,[acct-id/user-id,project-id,
			filename.element,ckpt-number,eqpmnt-type,reel-1/... reel-n]
		options include
			A, B, J, M, N, P, Q, R, S, U, V, Y, Z
	*/
	// TODO
	return nil, 0
}
