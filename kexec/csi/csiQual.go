// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec/exec"
	"khalehla/kexec/types"
	"log"
	"strings"
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
func handleQual(pkt *handlerPacket) uint64 {
	cleanOpts := ""
	if len(pkt.options) > 0 {
		var ok bool
		cleanOpts, ok = cleanOptions(pkt)
		if !ok {
			return 0_400000_000000
		}
	}

	var cleanQualifier string
	if len(pkt.arguments) > 0 {
		temp := strings.ToUpper(pkt.arguments)
		if !exec.IsValidQualifier(temp) {
			log.Printf("%v:Invalid qualifier '%v'", pkt.rce.RunId, pkt.statement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(types.ContingencyErrorMode, 04, 040)
			}
			return 0_600000_000000
		}
	}

	if cleanOpts == "" {
		// set implied qualifier
		if len(cleanQualifier) == 0 {
			log.Printf("%v: Missing qualifier '%v'", pkt.rce.RunId, pkt.statement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(types.ContingencyErrorMode, 04, 040)
			}
			return 0_600000_000000
		}
		pkt.rce.ImpliedQualifier = cleanQualifier
		return 0
	} else if cleanOpts == "D" {
		// set default qualifier
		if len(cleanQualifier) == 0 {
			log.Printf("%v: Missing qualifier '%v'", pkt.rce.RunId, pkt.statement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(types.ContingencyErrorMode, 04, 040)
			}
			return 0_600000_000000
		}
		pkt.rce.DefaultQualifier = cleanQualifier
		return 0
	} else if cleanOpts == "R" {
		// revert default and implied qualifiers
		if len(cleanQualifier) > 0 {
			log.Printf("%v: Should not specify qualifier with R option '%v'", pkt.rce.RunId, pkt.statement)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(types.ContingencyErrorMode, 04, 040)
			}
			return 0_600000_000000
		}
		pkt.rce.DefaultQualifier = pkt.rce.ProjectId
		pkt.rce.ImpliedQualifier = pkt.rce.ProjectId
		return 0
	} else {
		// option error
		log.Printf("%v: Conflicting options '%v'", pkt.rce.RunId, pkt.statement)
		if pkt.sourceIsExecRequest {
			pkt.rce.PostContingency(types.ContingencyErrorMode, 04, 040)
		}
		return 0_600000_000000
	}
}
