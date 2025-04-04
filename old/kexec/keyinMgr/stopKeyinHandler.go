// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"khalehla/kexec"
	kexec2 "khalehla/old/kexec"

	"strings"
	"time"
)

type StopKeyinHandler struct {
	exec    kexec.IExec
	source  kexec2.ConsoleIdentifier
	options string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewStopKeyinHandler(exec kexec.IExec, source kexec2.ConsoleIdentifier, options string, arguments string) IKeyinHandler {
	return &StopKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *StopKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *StopKeyinHandler) CheckSyntax() bool {
	return len(kh.options) == 0 && len(kh.arguments) == 0
}

func (kh *StopKeyinHandler) GetCommand() string {
	return "$!"
}

func (kh *StopKeyinHandler) GetHelp() []string {
	return []string{
		"$!",
		"Initiates auto-recovery of the operating system"}
}

func (kh *StopKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *StopKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *StopKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *StopKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *StopKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *StopKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *StopKeyinHandler) thread() {
	kh.threadStarted = true
	kh.exec.Stop(kexec2.StopOperatorInitiatedRecovery)
	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
