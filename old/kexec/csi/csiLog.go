// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec"
	"khalehla/kexec/facilitiesMgr"
	"khalehla/logger"
	facilitiesMgr2 "khalehla/old/kexec/facilitiesMgr"
)

// handleLog
//
//	@LOG message
//
// Inserts a message into the system log.
// The entire message is in pcs.operandFields[0][0]
func handleLog(pkt *handlerPacket) (facResult *facilitiesMgr.FacStatusResult, resultCode uint64) {
	logger.LogTraceF("CSI", "handleLog:%v", *pkt.pcs)
	facResult = facilitiesMgr.NewFacResult()
	resultCode = 0

	optWord, ok := cleanOptions(pkt)
	if !ok {
		facResult.PostMessage(kexec.FacStatusSyntaxErrorInImage, nil)
		resultCode = 0_400000_000000
		return
	}

	validMask := uint64(kexec.DOption | kexec.ROption)
	if !facilitiesMgr2.CheckIllegalOptions(pkt.rce, optWord, validMask, facResult, pkt.sourceIsExecRequest) {
		resultCode = 0_600000_000000
		return
	}

	text := pkt.pcs.operandFields[0][0]
	if len(text) == 0 {
		logger.LogWarningF("CSILog",
			"%v:Missing log message '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
		if pkt.sourceIsExecRequest {
			pkt.rce.PostContingency(012, 04, 040)
		}

		return nil, 0_600000_000000
	}

	logger.LogInfo(pkt.rce.RunId, text)
	return
}
