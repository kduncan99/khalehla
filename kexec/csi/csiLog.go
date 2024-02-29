// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"fmt"
	"khalehla/kexec"
	"log"
)

// handleLog
//
//	@LOG message
//
// Inserts a message into the system log.
// The entire message is in pcs.operandFields[0][0]
func handleLog(pkt *handlerPacket) (facResult *kexec.FacStatusResult, resultCode uint64) {
	facResult = kexec.NewFacResult()
	resultCode = 0

	optWord, ok := cleanOptions(pkt)
	if !ok {
		facResult.PostMessage(kexec.FacStatusSyntaxErrorInImage, nil)
		resultCode = 0_400000_000000
		return
	}

	validMask := uint64(kexec.DOption | kexec.ROption)
	if !checkIllegalOptions(pkt, optWord, validMask, facResult) {
		resultCode = 0_600000_000000
		return
	}

	text := pkt.pcs.operandFields[0][0]
	if len(text) == 0 {
		log.Printf("%v:Missing log message '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
		if pkt.sourceIsExecRequest {
			pkt.rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
		}

		return nil, 0_600000_000000
	}

	msg := fmt.Sprintf("%v:%v", pkt.rce.RunId, text)
	log.Printf(msg)
	return
}
