// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/facilitiesMgr"
	"strings"
	"time"
)

type MSKeyinHandler struct {
	exec            kexec.IExec
	source          kexec.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewMSKeyinHandler(exec kexec.IExec, source kexec.ConsoleIdentifier, options string, arguments string) IKeyinHandler {
	return &MSKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *MSKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *MSKeyinHandler) CheckSyntax() bool {
	if len(kh.options) != 0 || len(kh.arguments) != 0 {
		return false
	}

	if len(kh.arguments) != 0 {
		return false
	}

	return true
}

func (kh *MSKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *MSKeyinHandler) GetCommand() string {
	return "MS"
}

func (kh *MSKeyinHandler) GetHelp() []string {
	return []string{
		"MS",
		"Displays mass storage availability"}
}

func (kh *MSKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *MSKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *MSKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *MSKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *MSKeyinHandler) IsAllowed() bool {
	return kh.exec.GetPhase() == kexec.ExecPhaseRunning
}

func (kh *MSKeyinHandler) thread() {
	kh.threadStarted = true

	fm := kh.exec.GetFacilitiesManager().(*facilitiesMgr.FacilitiesManager)
	msAcc, msAvail, mfdAcc, mfdAvail := fm.GetTrackCounts()
	cfg := kh.exec.GetConfiguration()

	messages := []string{
		fmt.Sprintf("SUMMARY: FIXED TRACKS ACCESSIBLE = %v", msAcc),
		fmt.Sprintf("FIXED TRACKS AVAILABLE = %v", msAvail),
		fmt.Sprintf("STD ROLOUT START THRESHOLD = %1.2f", cfg.StandardRoloutStartThreshold),
		fmt.Sprintf("STD ROLOUT AVAILABILITY GOAL = %1.2f", cfg.StandardRoloutAvailabilityGoal),
		fmt.Sprintf("FIXED DIRECTORY TRACKS ACCESSIBLE = %v", mfdAcc),
		fmt.Sprintf("FIXED DIRECTORY TRACKS AVAILABLE = %v", mfdAvail),
	}

	for _, msg := range messages {
		kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
	}

	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
