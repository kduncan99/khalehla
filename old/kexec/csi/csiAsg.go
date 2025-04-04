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

func handleAsg(pkt *handlerPacket) (facResult *facilitiesMgr.FacStatusResult, resultCode uint64) {
	logger.LogTraceF("CSI", "handleAsg:%v", *pkt.pcs)
	/*
		Mass-Storage:
			@ASG[,options] filename[,type/reserve/granule/maximum/placement][,pack-id-1/.../pack-id-n,,,ACR-name]
		options include
			I, Z
			B, C, G, P, R, S, U, V, W for files to be cataloged
			A, D, E, K, M, Q, R, X, Y for existing files
			T for temporary files
		In the absence of any options, default to A if the file exists and is not assigned,
			or to the previous options if it is already assigned,
			or to T if it is not assigned and does not exist

		Tape Files:
			@ASG,options filename,type[/units/log/noise/processor/tape/
				format/data-converter/block-numbering/data-compression/
				buffered-write/expanded-buffer/,reel-1/reel-2/.../
				reel-n,expiration/mmspec,ring-indicator,ACR-name,CTL-pool ]
		options include
			B, I, E, N, R, W, Z
			C, G, P, U for files to be cataloged
			A, D, K, Q, X, Y for existing files
			R, T, W for temporary files
			E, H, L, M, O, S, V are hardware options
			F, J are labeling options

		Absolute Device:
			@ASG[,options] filename,*name,pack-id
		options include
			H, I, T, Z
	*/

	facResult = facilitiesMgr.NewFacResult()
	resultCode = 0

	// basic options validation - we'll do more specific checks later
	optWord, ok := cleanOptions(pkt)
	if !ok {
		facResult.PostMessage(kexec.FacStatusSyntaxErrorInImage, nil)
		resultCode = 0_600000_000000
		logger.LogTraceF("CSI", "handleAsg stat=%012o syntax", resultCode)
		return
	}

	// get the file specification and find the fileset if one exists
	var fsString string
	if len(pkt.pcs.operandFields) > 0 || len(pkt.pcs.operandFields[0]) > 0 {
		fsString = pkt.pcs.operandFields[0][0]
	} else {
		facResult.PostMessage(kexec.FacStatusFilenameIsRequired, nil)
		resultCode = 0_600000_000000
		logger.LogTraceF("CSI", "handleAsg stat=%012o fn required", resultCode)
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
		logger.LogTraceF("CSI", "handleAsg stat=%012o bad parse", resultCode)
		return
	}

	fm := pkt.exec.GetFacilitiesManager().(*facilitiesMgr.FacilitiesManager)
	facResult, resultCode = fm.AssignFile(pkt.rce, pkt.sourceIsExecRequest, fileSpec, optWord, pkt.pcs.operandFields)
	logger.LogTraceF("CSI", "handleAsg stat=%012o", resultCode)
	return
}
