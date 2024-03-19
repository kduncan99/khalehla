// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec/facilitiesMgr"
	"khalehla/klog"
)

func handleMode(pkt *handlerPacket) (*facilitiesMgr.FacStatusResult, uint64) {
	klog.LogTraceF("CSI", "handleMode:%v", *pkt.pcs)

	/*
		@MODE,options filename [,noise/processor/tape/format/
			data-converter/block-numbering/data-compression/buffered-write]
		options include
			E, H, L, M, O, S, V
	*/
	// TODO implement
	return nil, 0
}
