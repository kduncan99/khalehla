// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec"
	"khalehla/kexec/facilitiesMgr"
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
	exec                kexec.IExec
	rce                 kexec.RunControlEntry
	isTip               bool
	source              controlStatementSource
	sourceIsExecRequest bool
	pcs                 *ParsedControlStatement
}

type ParsedControlStatement struct {
	originalStatement string
	label             string
	mnemonic          string
	optionsFields     []string
	operandFields     [][]string
}

// HandleControlStatement parses a purported control statement,
// and if it seems initially valid the statement is handed off to a more-specific handler.
// We are invoked by the ER handlers for CSF$, ACSF$, and CSI$, by the RSI manager,
// and by whatever handler processes ECL for batch and demand runs.
// Returns a FacStatusResult suitable for use in output messages and return values from ER CSI$,
// as well as a more generic resultCode, suitable for return in A0 from ER CSF$/ACSF$.
func HandleControlStatement(
	exec kexec.IExec,
	rce kexec.RunControlEntry,
	source controlStatementSource,
	pcs *ParsedControlStatement,
) (facResult *facilitiesMgr.FacStatusResult, resultCode uint64) {

	pkt := handlerPacket{
		exec:                exec,
		rce:                 rce,
		isTip:               rce.IsTIP(),
		source:              source,
		sourceIsExecRequest: source == CSTSourceERCSF || source == CSTSourceERCSI,
	}

	switch pcs.mnemonic {
	case "ADD":
		if pkt.isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.GetRunId(), pcs.originalStatement)
			rce.PostContingency(012, 04, 042)

			facResult = facilitiesMgr.NewFacResult()
			facResult.PostMessage(facilitiesMgr.FacStatusIllegalControlStatement, []string{})
			resultCode = 0_400000_000000
			return
		}
		return handleAdd(&pkt)

	case "ASG":
		return handleAsg(&pkt)

	case "BRKPT":
		if pkt.isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.GetRunId(), pcs.originalStatement)
			rce.PostContingency(012, 04, 042)

			facResult = facilitiesMgr.NewFacResult()
			facResult.PostMessage(facilitiesMgr.FacStatusIllegalControlStatement, []string{})
			resultCode = 0_400000_000000
			return
		}
		return handleBrkpt(&pkt)

	case "@CAT":
		return handleCat(&pkt)

	case "@CKPT":
		if pkt.isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.GetRunId(), pcs.originalStatement)
			rce.PostContingency(012, 04, 042)

			facResult = facilitiesMgr.NewFacResult()
			facResult.PostMessage(facilitiesMgr.FacStatusIllegalControlStatement, []string{})
			resultCode = 0_400000_000000
			return
		}
		return handleCkpt(&pkt)

	case "@FREE":
		return handleFree(&pkt)

	case "@LOG":
		if source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for ER CSI$", rce.GetRunId(), pcs.originalStatement)
			rce.PostContingency(012, 04, 042)

			facResult = facilitiesMgr.NewFacResult()
			facResult.PostMessage(facilitiesMgr.FacStatusIllegalControlStatement, []string{})
			resultCode = 0_400000_000000
			return
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
			log.Printf("%v:CS '%v' not allowed for ER CSI$", rce.GetRunId(), pcs.originalStatement)

			rce.PostContingency(012, 04, 042)
			facResult = facilitiesMgr.NewFacResult()
			facResult.PostMessage(facilitiesMgr.FacStatusIllegalControlStatement, []string{})
			resultCode = 0_400000_000000
			return
		}
		return handleStart(&pkt)

	case "@SYM":
		if source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for ER CSI$", rce.GetRunId(), pcs.originalStatement)

			rce.PostContingency(012, 04, 042)
			facResult = facilitiesMgr.NewFacResult()
			facResult.PostMessage(facilitiesMgr.FacStatusIllegalControlStatement, []string{})
			resultCode = 0_400000_000000
			return
		}
		return handleSym(&pkt)

	case "@SYMCN":
		if pkt.isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.GetRunId(), pcs.originalStatement)

			rce.PostContingency(012, 04, 042)
			facResult = facilitiesMgr.NewFacResult()
			facResult.PostMessage(facilitiesMgr.FacStatusIllegalControlStatement, []string{})
			resultCode = 0_400000_000000
			return
		}
		return handleSymcn(&pkt)

	case "@USE":
		return handleUse(&pkt)
	}

	// syntax error
	log.Printf("%v:CS '%v' invalid control statement", rce.GetRunId(), pcs.originalStatement)
	rce.PostContingency(012, 04, 040)

	facResult = facilitiesMgr.NewFacResult()
	facResult.PostMessage(facilitiesMgr.FacStatusSyntaxErrorInImage, []string{})
	resultCode = 0_400000_000000
	return
}

// ParseControlStatement takes raw text input and splits it into its particular parts with little regard
// as to the semantics of any particular portion.
//
// format:
//
//	'@' [wsp] [ label ':' ]
//		[wsp] mnemonic [',' [wsp] option_field]
//		[wsp] field [',' [ws] field ]...
//		[' ' '.' ' ' command]
//
// option_field:
//
//	subfield [ '/' [wsp] subfield ]...
//
// field:
//
//	subfield [ '/' [wsp] subfield ]...
//
// wsp is whitespace, and is usually optional
//
// We do not handle continuation characters here - that must be dealt with at a higher level, with the
// fully-composed multi-line statement passed to us as a single string.
// We do not check image length here - that must be dealt with at a higher level.
func ParseControlStatement(
	rce kexec.RunControlEntry,
	statement string,
) (pcs *ParsedControlStatement, facResult *facilitiesMgr.FacStatusResult, resultCode uint64) {
	pcs = &ParsedControlStatement{
		label:             "",
		mnemonic:          "",
		originalStatement: statement,
		optionsFields:     make([]string, 0),
		operandFields:     make([][]string, 0),
	}
	facResult = facilitiesMgr.NewFacResult()
	resultCode = 0

	// trim comment
	working := statement + " "
	ix := strings.Index(working, " . ")
	if ix > 0 {
		working = working[:ix]
	}

	// find masterspace
	p := kexec.NewParser(working)
	if p.IsAtEnd() || !p.ParseSpecificCharacter('@') {
		// does not start with '@' - invalid statement
		log.Printf("%v:CS Syntax Error '%v' does not start with @", rce.GetRunId(), statement)
		rce.PostContingency(012, 04, 040)

		facResult.PostMessage(facilitiesMgr.FacStatusSyntaxErrorInImage, []string{})
		resultCode = 0_400000_000000
		return
	}

	// look for optional label
	p.SkipSpaces()
	ident, found, ok := p.ParseIdentifier()
	if !ok {
		facResult.PostMessage(facilitiesMgr.FacStatusSyntaxErrorInImage, []string{})
		resultCode = 0_400000_000000
		return
	} else if !found {
		// the only thing we can accept here is trailing whitespace, indicating an empty statement.
		// anything else constitutes a syntax error.
		p.SkipSpaces()
		if !p.IsAtEnd() {
			facResult.PostMessage(facilitiesMgr.FacStatusSyntaxErrorInImage, []string{})
			resultCode = 0_400000_000000
			return
		} else {
			return
		}
	}

	// We either found an identifier or a mnemonic - is there a ':' ?
	isLabel := p.ParseSpecificCharacter(':')
	if isLabel {
		// store the label and read another identifier - this time it will definitely be a mnemonic.
		pcs.label = ident

		p.SkipSpaces()
		ident, found, ok = p.ParseIdentifier()
		if !ok {
			// error in mnemonic
			facResult.PostMessage(facilitiesMgr.FacStatusSyntaxErrorInImage, []string{})
			resultCode = 0_400000_000000
			return
		} else if !found {
			// either an empty labeled statement, or a syntax error
			p.SkipSpaces()
			if !p.IsAtEnd() {
				facResult.PostMessage(facilitiesMgr.FacStatusSyntaxErrorInImage, []string{})
				resultCode = 0_400000_000000
				return
			} else {
				return
			}
		}
	}
	pcs.mnemonic = strings.ToUpper(ident)

	// Do we have subfields? (at this point, they'd be option subfields...)
	if p.ParseSpecificCharacter(',') {
		p.SkipSpaces()
		for !p.IsAtEnd() {
			sub, term := p.ParseUntil(" /,")
			pcs.optionsFields = append(pcs.optionsFields, sub)
			if term == 0 || term == ' ' {
				break
			} else if term == ',' {
				facResult.PostMessage(facilitiesMgr.FacStatusSyntaxErrorInImage, []string{})
				resultCode = 0_400000_000000
				return
			} else {
				p.SkipSpaces()
			}
		}
	}

	// Now do the operand fields - @LOG gets special treatment
	p.SkipSpaces()

	if pcs.mnemonic == "LOG" {
		str := p.GetRemainder()
		pcs.operandFields = make([][]string, 1)
		pcs.operandFields[0] = []string{str}
		return
	}

	// Anyone needing a file name with potential read/write keys gets special treatment...
	// in that we do not split the first operand field on forward slashes.
	cutSet := " /,"
	if pcs.mnemonic == "ASG" || pcs.mnemonic == "CAT" {
		cutSet = " ,"
	}

	opx := 0
	opy := 0
	for !p.IsAtEnd() {
		sub, term := p.ParseUntil(cutSet)
		for len(pcs.operandFields) <= opx {
			pcs.operandFields = append(pcs.operandFields, make([]string, 0))
		}
		for len(pcs.operandFields[opx]) <= opy {
			pcs.operandFields[opx] = append(pcs.operandFields[opx], "")
		}
		pcs.operandFields[opx][opy] = sub

		if term == 0 || term == ' ' {
			break
		} else if term == ',' {
			opx++
			opy = 0
			cutSet = " /,"
			p.SkipSpaces()
		} else {
			opy++
			p.SkipSpaces()
		}
	}

	return
}

// cleanOptions scans the options field and produces a corresponding options word.
// It will also post a contingency if there is an error and the source is ER CSF$/ACSF$/CSI$
func cleanOptions(pkt *handlerPacket) (result uint64, ok bool) {
	result = 0
	ok = true

	for _, opt := range strings.ToUpper(pkt.pcs.optionsFields[0]) {
		if opt < 'A' || opt > 'Z' {
			log.Printf("%v:CS Syntax Error in options field", pkt.rce.GetRunId())
			if pkt.sourceIsExecRequest {
				pkt.rce.PostContingency(012, 04, 040)
			}

			ok = false
			return
		}

		shift := opt - 'A' // A==0, Z==25
		result = result | (kexec.AOption >> shift)
	}

	return
}
