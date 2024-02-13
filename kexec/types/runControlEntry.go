// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

// RunControlEntry is the portion of the canonical PCT which contains information specific to a thread,
// but not to a program or any activities of the program.
type RunControlEntry struct {
	IsExec           bool
	RunId            string
	OriginalRunId    string
	AccountId        string
	ProjectId        string
	Userid           string
	DefaultQualifier string
	ImpliedQualifier string
	RunConditionWord RunConditionWord
	FacilityItems    []FacilitiesItem

	// TODO @USE table
	// TODO Program Control Entry
}

func NewRunControlEntry(
	runId string,
	originalRunId string,
	accountId string,
	projectId string,
	userId string) *RunControlEntry {
	return &RunControlEntry{
		IsExec:           false,
		RunId:            runId,
		OriginalRunId:    originalRunId,
		AccountId:        accountId,
		ProjectId:        projectId,
		Userid:           userId,
		DefaultQualifier: projectId,
		ImpliedQualifier: projectId,
		RunConditionWord: RunConditionWord{},
		FacilityItems:    make([]FacilitiesItem, 0),
	}
}

func (rce *RunControlEntry) PrintToTailSheet(message string) {
	// TODO
}
