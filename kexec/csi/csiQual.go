// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec"
	"log"
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
func handleQual(pkt *handlerPacket) (facResult *kexec.FacStatusResult, resultCode uint64) {
	facResult = kexec.NewFacResult()
	resultCode = 0

	optWord, ok := cleanOptions(pkt)
	if !ok {
		facResult.PostMessage(kexec.FacStatusSyntaxErrorInImage, nil)
		resultCode = 0_600000_000000
		return
	}

	validMask := uint64(kexec.DOption | kexec.ROption)
	if !pkt.rce.CheckIllegalOptions(optWord, validMask, facResult, pkt.sourceIsExecRequest) {
		resultCode = 0_600000_000000
		return
	}

	var qualifier string
	if len(pkt.pcs.operandFields) > 1 {
		log.Printf("%v:Too many op fields '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
		goto syntaxError
	} else if len(pkt.pcs.operandFields) == 1 {
		if len(pkt.pcs.operandFields[0]) != 1 {
			log.Printf("%v:Wrong # op subfields '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
			goto syntaxError
		} else {
			qualifier = pkt.pcs.operandFields[0][0]
			if !kexec.IsValidQualifier(qualifier) {
				log.Printf("%v:Invalid qualifier '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
				goto syntaxError
			}
		}
	}

	if optWord == 0 {
		// set implied qualifier
		if len(qualifier) == 0 {
			log.Printf("%v: Missing qualifier '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
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
			log.Printf("%v: Missing qualifier '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
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
			log.Printf("%v: Should not specify qualifier with R option '%v'",
				pkt.rce.RunId, pkt.pcs.originalStatement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
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
		log.Printf("%v: Conflicting options '%v'", pkt.rce.RunId, pkt.pcs.originalStatement)
		if pkt.sourceIsExecRequest {
			pkt.rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
		}

		facResult.PostMessage(kexec.FacStatusIllegalOptionCombination, []string{"D", "R"})
		return nil, 0_600000_000000
	}

syntaxError:
	if pkt.sourceIsExecRequest {
		pkt.rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
	}
	facResult.PostMessage(kexec.FacStatusSyntaxErrorInImage, nil)
	resultCode = 0_600000_000000
	return
}
