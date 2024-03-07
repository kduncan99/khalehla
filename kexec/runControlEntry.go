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

type RunType uint

const (
	RunTypeExec RunType = iota
	RunTypeBatch
	RunTypeDemand
	RunTypeTIP
)

// RunControlEntry is the portion of the canonical PCT which contains information specific to a thread,
// but not to a program or any activities of the program.
type RunControlEntry struct {
	AccountId        string
	DefaultQualifier string
	ImpliedQualifier string
	OriginalRunId    string
	ProjectId        string
	RunConditionWord uint64
	RunId            string
	UserId           string
	RunType          RunType
	UseItems         map[string]*UseItem
}

func newRunControlEntry(
	runType RunType,
	runId string,
	originalRunId string,
	projectId string,
	accountId string,
	userId string,
	runConditionWord uint64,
) *RunControlEntry {
	return &RunControlEntry{
		AccountId:        accountId,
		DefaultQualifier: projectId,
		ImpliedQualifier: projectId,
		OriginalRunId:    originalRunId,
		ProjectId:        projectId,
		RunConditionWord: runConditionWord,
		RunId:            runId,
		UserId:           userId,
		RunType:          runType,
		UseItems:         make(map[string]*UseItem),
	}
}

// TODO need NewBatchRunControlEntry()
// TODO need NewDemandRunControlEntry()
// TODO need NewTIPRunControlEntry()

func NewExecRunControlEntry(
	masterAccount string,
) *RunControlEntry {
	return newRunControlEntry(
		RunTypeExec,
		"EXEC-8",
		"EXEC-8",
		"SYS$",
		masterAccount,
		"EXEC-8",
		0)
}

func (rce *RunControlEntry) IsBatch() bool {
	return rce.RunType == RunTypeBatch
}

func (rce *RunControlEntry) IsDemand() bool {
	return rce.RunType == RunTypeDemand
}

func (rce *RunControlEntry) IsExec() bool {
	return rce.RunType == RunTypeExec
}

func (rce *RunControlEntry) IsTIP() bool {
	return rce.RunType == RunTypeTIP
}

func (rce *RunControlEntry) PostContingency(
	contingencyType ContingencyType,
	errorType uint,
	errorCode uint,
) {
	// TODO
	//  can this be disposed of in preference to the following version?
}

func (rce *RunControlEntry) PostContingencyWithAuxiliary(
	contingencyType ContingencyType,
	errorType uint,
	errorCode uint,
	aux pkg.Word36,
) {
	// TODO
}

func (rce *RunControlEntry) PostToTailSheet(message string) {
	// TODO
}
