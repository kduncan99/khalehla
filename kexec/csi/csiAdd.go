// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec/facilitiesMgr"
	"khalehla/klog"
)

// handleAdd
func handleAdd(pkt *handlerPacket) (*facilitiesMgr.FacStatusResult, uint64) {
	klog.LogTraceF("CSI", "handleAdd:%v", *pkt.pcs)
	/*
		@ADD[,options] name
		options:
			D Allows insertion of files or elements when operating under DATA or ELT,D
			E Implies @EOF at the end of the added file/element
			F Frees the file containing the added runstream upon completion of @ADD
			L Lists all ctl statements in added file/element *at the demand terminal*
				(always listed in batch mode)
			P Prints the @ADD stmt in the program listing
			R Prioritize added runstream over deferred ctrl stmt execution statck
	*/
	// TODO implement
	return nil, 0
}
