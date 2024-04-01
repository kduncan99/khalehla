// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
	"io"
	"khalehla/kexec/facItems"
	"khalehla/pkg"
)

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
	Privileges       map[Privilege]bool
	UseItems         map[string]*UseItem
	FacilityItems    []facItems.IFacilitiesItem
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
		Privileges:       make(map[Privilege]bool),
		UseItems:         make(map[string]*UseItem),
		FacilityItems:    make([]facItems.IFacilitiesItem, 0),
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

func (rce *RunControlEntry) DeleteUseItem(
	filename string,
) {
	delete(rce.UseItems, filename)
}

// TODO obsolete?
// // FindFacilitiesItem takes a qualifier-resolved FileSpecification and attempts to find an associated facilities item.
// // If checkUseItems is true, we loop through UseItems if and as appropriate as part of the search.
// // If return nil for facItem, there is not associated assigned file.
// // If we return false for foundUseItem, we did not find a use item (or were not asked to do so).
// // Any combination of nil or facItem, and true or false, are possible.
// func (rce *RunControlEntry) FindFacilitiesItem(
//	fileSpec *FileSpecification,
//	checkUseItems bool,
// ) (facItem IFacilitiesItem, foundUseItem bool) {
//	fileSpec = nil
//	foundUseItem = false
//	effectiveSpec := fileSpec
//
//	if checkUseItems {
//		for effectiveSpec.CouldBeInternalName() {
//			useItem, ok := rce.UseItems[facItem.GetFilename()]
//			if ok {
//				foundUseItem = true
//				effectiveSpec = useItem.FileSpecification
//			} else {
//				break
//			}
//		}
//	}
//
//	for _, fi := range rce.FacilityItems {
//		if effectiveSpec.MatchesFacilitiesItem(fi) {
//			facItem = fi
//			return
//		}
//	}
//
//	return
// }

// FindUseItem checks the given fileSpec, and if it refers to a use item, we return that use item.
// Otherwise, we return nil.
func (rce *RunControlEntry) FindUseItem(
	fileSpec *FileSpecification,
) (useItem *UseItem) {
	if fileSpec.CouldBeInternalName() {
		useItem, ok := rce.UseItems[fileSpec.Filename]
		if ok {
			return useItem
		}
	}
	return nil
}

// HasPrivilege indicates whether (for fundamental security) the run has a particular privilege.
// The exec always has all privileges.
func (rce *RunControlEntry) HasPrivilege(privilege Privilege) bool {
	if rce.IsExec() {
		return true
	} else {
		_, ok := rce.Privileges[privilege]
		return ok
	}
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

func (rce *RunControlEntry) IsPrivileged() bool {
	return rce.HasPrivilege(DLOCPrivilege)
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

func (rce *RunControlEntry) PostToPrint(text string, lineSkip uint) {
	// TODO
}

func (rce *RunControlEntry) PostToTailSheet(message string) {
	// TODO
}

// ResolveFileSpecification follows use item table to find the final external file name
// entry which applies to the caller, and fills in an effective qualifier if necessary.
func (rce *RunControlEntry) ResolveFileSpecification(
	fileSpecification *FileSpecification,
	checkUseItems bool,
) *FileSpecification {
	result := fileSpecification
	if checkUseItems {
		for result.CouldBeInternalName() {
			useItem, ok := rce.UseItems[result.Filename]
			if !ok {
				break
			}

			result = useItem.FileSpecification
		}
	}

	if len(result.Qualifier) == 0 {
		var qual string
		if result.HasAsterisk {
			qual = rce.ImpliedQualifier
		} else {
			qual = rce.DefaultQualifier
		}
		result = CopyFileSpecification(result)
		result.Qualifier = qual
	}

	return result
}

func (rce *RunControlEntry) Dump(dest io.Writer, indent string) {
	runType := ""
	if rce.IsBatch() {
		runType = "BATCH"
	} else if rce.IsDemand() {
		runType = "DEMAND"
	} else if rce.IsTIP() {
		runType = "TIP"
	} else if rce.IsExec() {
		runType = "EXEC"
	}

	_, _ = fmt.Fprintf(dest, "%v%v (%v) Acct:%v Proj:%v User:%v %v\n",
		indent, rce.RunId, rce.OriginalRunId, rce.AccountId, rce.ProjectId, rce.UserId, runType)
	_, _ = fmt.Fprintf(dest, "%v  rcw:%012o isPriv:%v priv:%v\n",
		indent, rce.RunConditionWord, rce.IsPrivileged(), rce.Privileges)
	_, _ = fmt.Fprintf(dest, "%v  defQual:%v impQual:%v\n",
		indent, rce.ImpliedQualifier, rce.DefaultQualifier)
	_, _ = fmt.Fprintf(dest, "%v  Facility Items:\n", indent)
	for _, fi := range rce.FacilityItems {
		fi.Dump(dest, indent+"    ")
	}
}
