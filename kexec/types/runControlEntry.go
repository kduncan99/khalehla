// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

import (
	"khalehla/pkg"
)

// RunControlEntry is the portion of the canonical PCT which contains information specific to a thread,
// but not to a program or any activities of the program.
type RunControlEntry struct {
	IsExec           bool
	RunId            pkg.Word36
	OriginalRunId    pkg.Word36
	AccountId        []pkg.Word36
	ProjectId        []pkg.Word36
	Userid           []pkg.Word36
	DefaultQualifier []pkg.Word36
	ImpliedQualifier []pkg.Word36
	RunConditionWord RunConditionWord

	// TODO Facility Item Table
	// TODO @USE table
	// TODO Program Control Entry
}

// TODO where should these constants live?

var SystemRunId = "EXEC-8"
var SystemAccountId = "SYSTEM"
var OverheadAccountId = "INSTALLATION"
var MasterAccountId = ""
var PrivilegedAccountId = "123456"
var SystemProjectId = "SYSTEM"
var SystemUserId = "EXEC-8"
var OverheadUserId = "INSTALLATION"
var SecurityOfficerUserId = ""
var SystemQualifier = "SYS$"

// TODO Need to implement a logging mechanism, and see logging configuration Exec Install/Config 8.3

// ExecRunControlEntry is the RCE for the EXEC - it always exists and is always (or should always be) in the RCT
var ExecRunControlEntry = RunControlEntry{
	IsExec:           true,
	RunId:            pkg.NewFromStringToFieldata(SystemRunId, 1)[0],
	OriginalRunId:    pkg.NewFromStringToAscii(SystemRunId, 1)[0],
	AccountId:        pkg.NewFromStringToAscii(SystemAccountId, 2),
	ProjectId:        pkg.NewFromStringToAscii(SystemProjectId, 2),
	Userid:           pkg.NewFromStringToAscii(SystemUserId, 2),
	DefaultQualifier: pkg.NewFromStringToAscii(SystemQualifier, 2),
	ImpliedQualifier: pkg.NewFromStringToAscii(SystemQualifier, 2),
	RunConditionWord: RunConditionWord{},
}
