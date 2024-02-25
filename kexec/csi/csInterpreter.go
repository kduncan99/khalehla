// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/exec"
	"log"
	"strings"
)

type controlStatementSource uint

const (
	CSTSourceECL   controlStatementSource = iota // batch or demand
	CSTSourceERCSF                               // ER CSF$ or ER ACSF$
	CSTSourceERCSI
	CSTSourceTransparent // from RSI complex
)

type handlerPacket struct {
	rce                 *exec.RunControlEntry
	isTip               bool
	source              controlStatementSource
	sourceIsExecRequest bool
	statement           string
	options             string
	arguments           string
}

// HandleControlStatement parses a purported control statement,
// and if it seems initially valid the statement is handed off to a more-specific handler.
// We are invoked by the ER handlers for CSF$, ACSF$, and CSI$, by the RSI manager,
// and by whatever handler processes ECL for batch and demand runs.
// TODO we should return a status packet including CSF status word and list of facility status entries
func HandleControlStatement(
	rce *exec.RunControlEntry,
	source controlStatementSource,
	statement string,
) uint64 {

	split := strings.SplitN(strings.TrimSpace(statement), " ", 1)
	if len(split) == 0 || split[0][0] != '@' {
		// does not start with '@' - invalid statement
		log.Printf("%v:CS Syntax Error '%v' does not start with @", rce.RunId, statement)
		rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
		return 0_400000_000000
	}

	subSplit := strings.SplitN(split[0], ",", 1)
	command := strings.ToUpper(subSplit[0])
	var options string
	if len(subSplit) == 2 {
		options = subSplit[1]
	}

	var arguments string
	if len(split) == 2 {
		arguments = split[1]
	}

	pkt := handlerPacket{
		rce:                 rce,
		isTip:               rce.IsTIPTransaction(),
		source:              source,
		sourceIsExecRequest: source == CSTSourceERCSF || source == CSTSourceERCSI,
		statement:           statement,
		options:             options,
		arguments:           arguments,
	}

	switch command {
	case "ADD":
		if pkt.isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.RunId, statement)
			rce.PostContingency(kexec.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return handleAdd(&pkt)

	case "ASG":
		return handleAsg(&pkt)

	case "BRKPT":
		if pkt.isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.RunId, statement)
			rce.PostContingency(kexec.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return handleBrkpt(&pkt)

	case "@CAT":
		return handleCat(&pkt)

	case "@CKPT":
		if pkt.isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.RunId, statement)
			rce.PostContingency(kexec.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return handleCkpt(&pkt)

	case "@FREE":
		return handleFree(&pkt)

	case "@LOG":
		if source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for ER CSI$", rce.RunId, statement)
			rce.PostContingency(kexec.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return handleLog(&pkt)

	case "@MODE":
		return handleMode(&pkt)

	case "@QUAL":
		return handleQual(&pkt)

	case "@RSTRT":
		return handleRstrt(&pkt)

	case "@START":
		if source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for ER CSI$", rce.RunId, statement)
			rce.PostContingency(kexec.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return handleStart(&pkt)

	case "@SYM":
		if source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for ER CSI$", rce.RunId, statement)
			rce.PostContingency(kexec.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return handleSym(&pkt)

	case "@SYMCN":
		if pkt.isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.RunId, statement)
			rce.PostContingency(kexec.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return handleSymcn(&pkt)

	case "@USE":
		return handleUse(&pkt)
	}

	// syntax error
	log.Printf("%v:CS '%v' invalid control statement", rce.RunId, statement)
	rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
	return 0_400000_000000
}

// cleanOptions scans the options field and produces an alternate field with
// all option letters made uppercase and duplicates removed.
// It also posts a contingency if the source is ER CSF$/ACSF$/CSI$
func cleanOptions(pkt *handlerPacket) (string, bool) {
	upperOptions := strings.ToUpper(pkt.options)
	optFlags := make([]bool, 26)
	for _, opt := range upperOptions {
		if opt < 'A' || opt > 'Z' {
			log.Printf("%v:CS Syntax Error in options field", pkt.rce.RunId)
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(kexec.ContingencyErrorMode, 04, 040)
			}
			return "", false
		}
		optFlags[opt-'A'] = true
	}

	str := ""
	for ox, flag := range optFlags {
		if flag {
			str += fmt.Sprintf("%c", ox+'A')
		}
	}

	return str, true
}
