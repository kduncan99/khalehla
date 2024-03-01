// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import "khalehla/pkg"

/*
TIP:
When TIP initializes a TIP program or a program connected to TIP, control parameters are inserted into the
program control table (PCT). These parameters are inserted in that portion of the PCT called the Status list of the
operating program (SLOP) table. The SLOP control word (currently, word 0337) of the PCT indicates whether the SLOP table
is being used. The SLOP table starts immediately after the SLOP control word (currently, word 0340 of the PCT).

SLOP control word:
scw,H1 TL$SCH offset to the scheduling packet
scw,S4 PRINT$ file assignment flag
scw,T3 transaction program state
		00000 Not TIP (has no SLOP table)
		00001 Batch/Demand connected to TIP
		00050 TIP transaction
		05000 Online batch program not connected to TIP
		05001 Online batch connected to TIP
*/

// RunControlEntry is the portion of the canonical PCT which contains information specific to a thread,
// but not to a program or any activities of the program.
type RunControlEntry interface {
	GetAccountId() string
	GetDefaultQualifier() string
	GetImpliedQualifier() string
	GetOriginalRunId() string
	GetProjectId() string
	GetRunConditionWord() uint64
	GetRunId() string
	GetUserId() string
	IsBatch() bool
	IsDemand() bool
	IsExec() bool
	IsTIP() bool
	PostContingency(contingencyType ContingencyType, errorType uint, errorCode uint)
	PostContingencyWithAuxiliary(contingencyType ContingencyType, errorType uint, errorCode uint, aux pkg.Word36)
	PostToTailSheet(message string)
	SetDefaultQualifier(string)
	SetImpliedQualifier(string)
	SetRunConditionWord(uint64)
}

func GetEffectiveQualifier(rce RunControlEntry, fileSpec *FileSpecification) string {
	if fileSpec.HasAsterisk {
		if len(fileSpec.Qualifier) == 0 {
			return fileSpec.Qualifier
		} else {
			return rce.GetImpliedQualifier()
		}
	} else {
		return rce.GetDefaultQualifier()
	}
}
