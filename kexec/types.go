// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/pkg"
)

// ConsoleIdentifier is a 36-bit word implemented as uint64, containing a unique Identifier for the console.
// It *might* be the console Name, but that is merely by convention.
type ConsoleIdentifier pkg.Word36

// pseudo-enumeration pkg

type ContingencyType uint
type DeviceRelativeWordAddress uint
type ExecPhase uint
type FacStatus uint
type Granularity uint
type LDATIndex uint64
type MFDBlockId uint64
type MFDRelativeAddress uint64
type MFDSectorId uint64
type MFDTrackId uint64

// Things related to ConsoleManager which are reference by kexec.interface
// and therefore cannot go into consoleMgr package

type ConsoleReadOnlyMessage struct {
	Source         *RunControlEntry
	Routing        *ConsoleIdentifier // may be nil
	RunId          *string            // for logging purposes, may not match RCE - may be nil
	Text           string             // message to be displayed
	DoNotEmitRunId bool
}

type ConsoleReadReplyMessage struct {
	Source         *RunControlEntry
	Routing        *ConsoleIdentifier // may be nil
	RunId          *string            // for logging purposes, may not match RCE - may be nil
	Text           string             // message to be displayed
	DoNotEmitRunId bool
	DoNotLogReply  bool // reply may contain secure information
	MaxReplyLength int
	Reply          string
}
