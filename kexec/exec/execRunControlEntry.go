// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package exec

import (
	"khalehla/kexec"
	"khalehla/pkg"
)

const execAccountId = ""
const execUserId = "EXEC8"
const execRunId = "EXEC-8"
const execProjectId = "SYS$"
const execQualifier = "SYS$"

// ExecRunControlEntry is a special RCE that applies to the exec
type ExecRunControlEntry struct {
	accountId string
}

func NewExecRunControlEntry(
	masterAccount string,
) *ExecRunControlEntry {
	return &ExecRunControlEntry{
		accountId: masterAccount,
	}
}

func (rce *ExecRunControlEntry) GetAccountId() string {
	return rce.accountId
}

func (rce *ExecRunControlEntry) GetDefaultQualifier() string {
	return execQualifier
}

func (rce *ExecRunControlEntry) GetImpliedQualifier() string {
	return execQualifier
}

func (rce *ExecRunControlEntry) GetOriginalRunId() string {
	return execRunId
}

func (rce *ExecRunControlEntry) GetProjectId() string {
	return execProjectId
}

func (rce *ExecRunControlEntry) GetRunConditionWord() uint64 {
	return 0
}

func (rce *ExecRunControlEntry) GetRunId() string {
	return execRunId
}

func (rce *ExecRunControlEntry) GetUserId() string {
	return execUserId
}

func (rce *ExecRunControlEntry) IsBatch() bool {
	return false
}

func (rce *ExecRunControlEntry) IsDemand() bool {
	return false
}

func (rce *ExecRunControlEntry) IsExec() bool {
	return true
}

func (rce *ExecRunControlEntry) IsTIP() bool {
	return false
}

func (rce *ExecRunControlEntry) PostContingency(
	contingencyType kexec.ContingencyType,
	errorType uint,
	errorCode uint,
) {
	// TODO causes a stop
}

func (rce *ExecRunControlEntry) PostContingencyWithAuxiliary(
	contingencyType kexec.ContingencyType,
	errorType uint,
	errorCode uint,
	auxiliaryInfo pkg.Word36,
) {
	// TODO causes a stop
}

func (rce *ExecRunControlEntry) PostToTailSheet(message string) {
	// We don't have a tail sheet, so...
}

func (rce *ExecRunControlEntry) SetDefaultQualifier(string) {}
func (rce *ExecRunControlEntry) SetImpliedQualifier(string) {}
func (rce *ExecRunControlEntry) SetRunConditionWord(uint64) {}
