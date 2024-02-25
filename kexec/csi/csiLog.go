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
// Inserts a message into the system log
func handleLog(pkt *handlerPacket) uint64 {
	if len(pkt.options) > 0 {
		log.Printf("%v:Invalid options '%v'", pkt.rce.RunId, pkt.statement)
		if pkt.sourceIsExecRequest {
			pkt.rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
		}
		return 0_600000_000000
	}

	if len(pkt.arguments) == 0 {
		log.Printf("%v:Missing log message '%v'", pkt.rce.RunId, pkt.statement)
		if pkt.sourceIsExecRequest {
			pkt.rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
		}
		return 0_600000_000000
	}

	msg := fmt.Sprintf("%v:%v", pkt.rce.RunId, pkt.arguments)
	log.Printf(msg)
	return 0
}
