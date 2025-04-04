// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec"
	"khalehla/kexec/facilitiesMgr"
	"khalehla/logger"
	kexec2 "khalehla/old/kexec"
	facilitiesMgr2 "khalehla/old/kexec/facilitiesMgr"
)

// handleQual updates the default and/or implied qualifier for the run control entry.
//
//	@QUAL[,options] [directory-id#][qualifier]
//
// With no options, the implied qualifier is set as specified
// With 'D' option, the default qualifier is set as specified
// With 'R' option, the default and implied qualifiers are reset to the project-id
// D and R are mutually exclusive.
// We do not support directory-id at the moment
func handleQual(pkt *handlerPacket) (facResult *facilitiesMgr.FacStatusResult, resultCode uint64) {
	logger.LogTraceF("CSI", "handleQual:%v", *pkt.pcs)
	facResult = facilitiesMgr.NewFacResult()
	resultCode = 0

	optWord, ok := cleanOptions(pkt)
	if !ok {
		facResult.PostMessage(kexec.FacStatusSyntaxErrorInImage, nil)
		resultCode = 0_600000_000000
		return
	}

	validMask := uint64(kexec.DOption | kexec.ROption)
	if !facilitiesMgr2.CheckIllegalOptions(pkt.rce, optWord, validMask, facResult, pkt.sourceIsExecRequest) {
		resultCode = 0_600000_000000
		return
	}

	var qualifier string
	if len(pkt.pcs.operandFields) > 1 {
		logger.LogInfoF("CSIQual", "%v:Too many op fields '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
		goto syntaxError
	} else if len(pkt.pcs.operandFields) == 1 {
		if len(pkt.pcs.operandFields[0]) != 1 {
			logger.LogInfoF("CSIQual", "%v:Wrong # op subfields '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
			goto syntaxError
		} else {
			qualifier = pkt.pcs.operandFields[0][0]
			if !kexec2.IsValidQualifier(qualifier) {
				logger.LogInfoF("CSIQual", "%v:Invalid qualifier '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
				goto syntaxError
			}
		}
	}

	if optWord == 0 {
		// set implied qualifier
		if len(qualifier) == 0 {
			logger.LogInfoF("CSIQual", "%v: Missing qualifier '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(012, 04, 040)
			}
			facResult.PostMessage(kexec.FacStatusDirectoryOrQualifierMustAppear, nil)
			resultCode = 0_600000_000000
			return
		}

		pkt.rce.ImpliedQualifier = qualifier
		facResult.PostMessage(kexec.FacStatusComplete, []string{"QUAL"})
		return
	} else if optWord == kexec.DOption {
		// set default qualifier
		if len(qualifier) == 0 {
			logger.LogInfoF("CSIQual", "%v: Missing qualifier '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(012, 04, 040)
			}
			facResult.PostMessage(kexec.FacStatusDirectoryOrQualifierMustAppear, nil)
			resultCode = 0_600000_000000
			return
		}

		pkt.rce.DefaultQualifier = qualifier
		facResult.PostMessage(kexec.FacStatusComplete, []string{"QUAL"})
		return
	} else if optWord == kexec.ROption {
		// revert default and implied qualifiers
		if len(qualifier) > 0 {
			logger.LogInfoF("CSIQual", "%v: Should not specify qualifier with R option '%v'",
				pkt.rce.RunId, pkt.pcs.originalStatement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(012, 04, 040)
			}
			facResult.PostMessage(kexec.FacStatusDirectoryAndQualifierMayNotAppear, nil)
			resultCode = 0_600000_000000
			return
		}

		pkt.rce.DefaultQualifier = pkt.rce.ProjectId
		pkt.rce.ImpliedQualifier = pkt.rce.ProjectId
		facResult.PostMessage(kexec.FacStatusComplete, []string{"QUAL"})
		return
	} else {
		logger.LogInfoF("CSIQual", "%v: Conflicting options '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
		if pkt.sourceIsExecRequest {
			pkt.rce.PostContingency(012, 04, 040)
		}

		facResult.PostMessage(kexec.FacStatusIllegalOptionCombination, []string{"D", "R"})
		return nil, 0_600000_000000
	}

syntaxError:
	if pkt.sourceIsExecRequest {
		pkt.rce.PostContingency(012, 04, 040)
	}
	facResult.PostMessage(kexec.FacStatusSyntaxErrorInImage, nil)
	resultCode = 0_600000_000000
	return
}
