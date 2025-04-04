// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec"
	"khalehla/kexec/facilitiesMgr"
	"khalehla/logger"
	kexec2 "khalehla/old/kexec"
)

// handleCat updates the default and/or implied qualifier for the run control entry.
// pcs.operandFields[0] contains the full file specification in a single subfield
// pcs.operandFields[1..n] may contain multiple subfields
func handleCat(pkt *handlerPacket) (facResult *facilitiesMgr.FacStatusResult, resultCode uint64) {
	logger.LogTraceF("CSI", "handleCat:%v", *pkt.pcs)
	facResult = facilitiesMgr.NewFacResult()
	resultCode = 0

	// basic options validation - we'll do more specific checks later
	optWord, ok := cleanOptions(pkt)
	if !ok {
		facResult.PostMessage(kexec.FacStatusSyntaxErrorInImage, nil)
		resultCode = 0_600000_000000
		return
	}

	// get the file specification and find the fileset if one exists
	var fsString string
	if len(pkt.pcs.operandFields) > 0 || len(pkt.pcs.operandFields[0]) > 0 {
		fsString = pkt.pcs.operandFields[0][0]
	} else {
		facResult.PostMessage(kexec.FacStatusFilenameIsRequired, nil)
		resultCode = 0_600000_000000
		logger.LogTraceF("CSI", "handleAsg stat=%012o", resultCode)
		return
	}

	p := kexec2.NewParser(fsString)
	fileSpec, fsCode, ok := kexec2.ParseFileSpecification(p)
	if !ok {
		if pkt.sourceIsExecRequest {
			pkt.rce.PostContingency(012, 04, 040)
		}
		facResult.PostMessage(fsCode, []string{})
		resultCode = 0_600000_000000
		return
	}

	fm := pkt.exec.GetFacilitiesManager().(*facilitiesMgr.FacilitiesManager)
	return fm.CatalogFile(pkt.rce, pkt.sourceIsExecRequest, fileSpec, optWord, pkt.pcs.operandFields)
}
