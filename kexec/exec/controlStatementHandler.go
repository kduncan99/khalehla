// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package exec

import (
	"fmt"
	"khalehla/kexec/types"
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

// HandleControlStatement parses a purported control statement,
// and if it seems initially valid the statement is handed off to a more-specific handler.
// We are invoked by the ER handlers for CSF$, ACSF$, and CSI$, by the RSI manager,
// and by whatever handler processes ECL for batch and demand runs.
// TODO we should return a status packet including CSF status word and list of facility status entries
func (e *Exec) HandleControlStatement(
	rce *types.RunControlEntry,
	source controlStatementSource,
	statement string,
) uint64 {

	split := strings.SplitN(strings.TrimSpace(statement), " ", 1)
	if len(split) == 0 || split[0][0] != '@' {
		// does not start with '@' - invalid statement
		log.Printf("%v:CS Syntax Error '%v' does not start with @", rce.RunId, statement)
		rce.PostContingency(types.ContingencyErrorMode, 04, 040)
		return 0_400000_000000
	}

	subSplit := strings.SplitN(split[0], ",", 1)
	command := strings.ToUpper(subSplit[0])
	var options *string
	if len(subSplit) == 2 {
		options = &subSplit[1]
	}

	var arguments *string
	if len(split) == 2 {
		arguments = &split[1]
	}

	isTip := rce.IsTIPTransaction()
	switch command {
	case "ADD":
		if isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return e.handleAdd(rce, source, options, arguments)
	case "ASG":
		return e.handleAsg(rce, statement, options, arguments)
	case "BRKPT":
		if isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return e.handleBrkpt(rce, statement, options, arguments)
	case "@CAT":
		return e.handleCat(rce, statement, options, arguments)
	case "@CKPT":
		if isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return e.handleCkpt(rce, statement, options, arguments)
	case "@FREE":
		return e.handleFree(rce, statement, options, arguments)
	case "@LOG":
		if source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for ER CSI$", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return e.handleLog(rce, statement, options, arguments)
	case "@MODE":
		return e.handleMode(rce, statement, options, arguments)
	case "@QUAL":
		return e.handleQual(rce, statement, options, arguments)
	case "@RSTRT":
		return e.handleRstrt(rce, statement, options, arguments)
	case "@START":
		if source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for ER CSI$", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return e.handleStart(rce, statement, options, arguments)
	case "@SYM":
		if source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for ER CSI$", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return e.handleSym(rce, statement, options, arguments)
	case "@SYMCN":
		if isTip || source == CSTSourceERCSI {
			log.Printf("%v:CS '%v' not allowed for TIP or ER CSI$", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 042)
			return 0_400000_000000
		}
		return e.handleSymcn(rce, statement, options, arguments)
	case "@USE":
		return e.handleUse(rce, statement, options, arguments)
	}

	// syntax error
	log.Printf("%v:CS '%v' invalid control statement", rce.RunId, statement)
	rce.PostContingency(types.ContingencyErrorMode, 04, 040)
	return 0_400000_000000
}

func cleanOptions(options string) (string, bool) {
	upperOptions := strings.ToUpper(options)
	optFlags := make([]bool, 26)
	for _, opt := range upperOptions {
		if opt < 'A' || opt > 'Z' {
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

// handleAdd
func (e *Exec) handleAdd(
	rce *types.RunControlEntry,
	source controlStatementSource,
	options *string,
	arguments *string,
) uint64 {

	/*
		@ADD[,options] name
		options:
			D Allows insertion of files or elements when operating under DATA or ELT,D
			E Implies @EOF at the end of the added file/element
			F Frees the file containing the added runstream upon completion of @ADD
			L Lists all ctl statements in added file/element *at the demand terminal*
				(always listed in batch mode)
			P Prints the @ADD stmt in the program listing
			R Prioritize added runstream over deferred ctrl stmt execution statck
	*/
	// TODO
	return 0
}

func (e *Exec) handleAsg(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {
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

	// TODO
	return 0
}

func (e *Exec) handleBrkpt(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {

	/*
		For primary print files
			@BRKPT,L PRINT$[/part-name]
			@BRKPT PRINT$[,filename]
		For alternate print files
			@BRKPT[,E] internal-fname
			@BRKPT[,E] ,fname
		E inhibits EOF positioning for alt read files on magtape
	*/
	// TODO
	return 0
}

func (e *Exec) handleCat(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {
	/*
		For Mass Storage Files
			@CAT[,options] filename[,type/reserve/granule/maximum,pack-id-1/.../pack-id-n,,,ACR-name]
		options include
			B, G, P, R, S, V, W, Z

		For Tape Files
			@CAT,options filename,type[/units/log/noise/processor/tape/
				format/data-converter/block-numbering/data-compression/
				buffered-write/expanded-buffer,reel-1/reel-2/.../reel-n,
				expiration-period/mmspec,,ACR-name,CTL-pool]
		options include
			E, G, H, J, L, M, O, P, R, S, V, W, Z
	*/
	// TODO
	return 0
}

func (e *Exec) handleCkpt(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {

	/*
		@CKPT[,options] [filename.element,,eqpmnt-type,reel-1/.../reel-n,,,,CTL-pool]
		options include
			A, B, C, J, M, N, P, Q, R, S, T, Z
	*/
	// TODO
	return 0
}

func (e *Exec) handleFree(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {

	/*
		@FREE[,options] filename
		options include
			A, B, D, I, R, S, X
	*/
	// TODO
	return 0
}

// handleLog
//
//	@LOG message
//
// Inserts a message into the system log
func (e *Exec) handleLog(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {
	if options != nil {
		log.Printf("%v:Invalid options '%v'", rce.RunId, statement)
		rce.PostContingency(types.ContingencyErrorMode, 04, 040)
		return 0_600000_000000
	}

	if arguments == nil {
		log.Printf("%v:Missing log message '%v'", rce.RunId, statement)
		rce.PostContingency(types.ContingencyErrorMode, 04, 040)
		return 0_600000_000000
	}

	msg := fmt.Sprintf("%v:%v", rce.RunId, *arguments)
	log.Printf(msg)
	return 0
}

func (e *Exec) handleMode(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {

	/*
		@MODE,options filename [,noise/processor/tape/format/
			data-converter/block-numbering/data-compression/buffered-write]
		options include
			E, H, L, M, O, S, V
	*/
	// TODO
	return 0
}

// handleQual updates the default and/or implied qualifier for the run control entry.
//
//	@QUAL[,options] [directory-id#][qualifier]
//
// With no options, the implied qualifier is set as specified
// With 'D' option, the default qualifier is set as specified
// With 'R' option, the default and implied qualifiers are reset to the project-id
// D and R are mutually exclusive.
// We do not support directory-id at the moment
func (e *Exec) handleQual(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {
	cleanOpts := ""
	if options != nil {
		var ok bool
		cleanOpts, ok = cleanOptions(*options)
		if !ok {
			// invalid option, syntax error
			log.Printf("%v:CS Syntax Error '%v' error in options field", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 040)
			return 0_400000_000000
		}
	}

	var cleanQual *string = nil
	if arguments != nil {
		temp := strings.ToUpper(*arguments)
		if !e.IsValidQualifier(temp) {
			log.Printf("%v:Invalid qualifier '%v'", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 040)
			return 0_600000_000000
		}
	}

	if cleanOpts == "" {
		// set implied qualifier
		if cleanQual == nil {
			log.Printf("%v: Missing qualifier '%v'", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 040)
			return 0_600000_000000
		}
		rce.ImpliedQualifier = *cleanQual
		return 0
	} else if cleanOpts == "D" {
		// set default qualifier
		if cleanQual == nil {
			log.Printf("%v: Missing qualifier '%v'", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 040)
			return 0_600000_000000
		}
		rce.DefaultQualifier = *cleanQual
		return 0
	} else if cleanOpts == "R" {
		// revert default and implied qualifiers
		if cleanQual != nil {
			log.Printf("%v: Should not specify qualifier with R option '%v'", rce.RunId, statement)
			rce.PostContingency(types.ContingencyErrorMode, 04, 040)
			return 0_600000_000000
		}
		rce.DefaultQualifier = rce.ProjectId
		rce.ImpliedQualifier = rce.ProjectId
		return 0
	} else {
		// option error
		log.Printf("%v: Conflicting options '%v'", rce.RunId, statement)
		rce.PostContingency(types.ContingencyErrorMode, 04, 040)
		return 0_600000_000000
	}
}

func (e *Exec) handleRstrt(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {

	/*
		@RSTRT[,scheduling-priority/options/processor-dispatching-priority]
			run-id,[acct-id/user-id,project-id,
			filename.element,ckpt-number,eqpmnt-type,reel-1/... reel-n]
		options include
			A, B, J, M, N, P, Q, R, S, U, V, Y, Z
	*/
	// TODO
	return 0
}

func (e *Exec) handleStart(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {

	/*
		@START runstream-name
		@START[,scheduling-priority/options/processor-dispatching-priority]
			runstream-name[,setc,run-id,acct-id/user-id,project-id,runtime/deadline,pages/cards,start-time]
		scheduling-priority is one character from 'A' to 'Z'
		options include
			B, N, P, R, T, U, V, W, X, Y, Z
		processor-dispatching-priority is one character from 'A' to 'Z'
	*/
	// TODO
	return 0
}

func (e *Exec) handleSym(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {

	/*
		@SYM,options filename,[number,device,pname-1/pname-2/pname-3/.../pname-n,banner]
		@SYM,options filename,[number,user-id/U,pname-1/pname/pname-3/.../pname-n,banner]
		options include
			A, D, F, J, K, L, N, U
	*/
	// TODO
	return 0
}

func (e *Exec) handleSymcn(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {

	/*
		@SYMCN,L
		@SYMCN,N [{n}]
	*/
	// TODO
	return 0
}

func (e *Exec) handleUse(
	rce *types.RunControlEntry,
	statement string,
	options *string,
	arguments *string,
) uint64 {

	/*
		@USE[,I] fname-1,fname-2
	*/
	// TODO
	return 0
}
