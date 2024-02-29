// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

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
	UseItems         map[string]UseItem // key is internal file Name
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
		UseItems:         make(map[string]UseItem),
	}
}

// CheckIllegalOptions compares the given options word to the allowed options word,
// producing a fac message for each option set in the given word which does not appear in the allowed word.
// Returns true if no such instances were found, else false
// If not ok and the source is an ER CSF$/ACSF$/CSI$, we post a contingency
func (rce *RunControlEntry) CheckIllegalOptions(
	givenOptions uint64,
	allowedOptions uint64,
	facResult *FacStatusResult,
	sourceIsExec bool,
) bool {
	bit := uint64(AOption)
	letter := 'A'
	ok := true

	for {
		if bit&givenOptions != 0 && bit&allowedOptions == 0 {
			param := string(letter)
			facResult.PostMessage(FacStatusIllegalOption, []string{param})
			ok = false
		}

		if bit == ZOption {
			break
		} else {
			letter++
			bit >>= 1
		}
	}

	if !ok {
		if sourceIsExec {
			rce.PostContingency(ContingencyErrorMode, 04, 040)
		}
	}

	return ok
}

func (rce *RunControlEntry) GetEffectiveQualifier(fileSpec *FileSpecification) string {
	if fileSpec.HasAsterisk {
		if len(fileSpec.Qualifier) == 0 {
			return fileSpec.Qualifier
		} else {
			return rce.ImpliedQualifier
		}
	} else {
		return rce.DefaultQualifier
	}
}

func (rce *RunControlEntry) IsTIPTransaction() bool {
	return false // TODO
}

func (rce *RunControlEntry) PostContingency(contingencyType ContingencyType, errorType uint, errorCode uint) {
	// TODO
}

func (rce *RunControlEntry) PrintToTailSheet(message string) {
	// TODO
}
