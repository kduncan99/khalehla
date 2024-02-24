// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"khalehla/kexec/types"
	"strings"
	"time"
)

type FOOKeyinHandler struct {
	exec            types.IExec
	source          types.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewFOOKeyinHandler(exec types.IExec, source types.ConsoleIdentifier, options string, arguments string) KeyinHandler {
	return &FOOKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *FOOKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *FOOKeyinHandler) CheckSyntax() bool {
	// Accepted:
	//		D
	//		D,UTC
	//		D SHIFT
	// TODO
	return true
}

func (kh *FOOKeyinHandler) GetCommand() string {
	return "FOO"
}

func (kh *FOOKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *FOOKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *FOOKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *FOOKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *FOOKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *FOOKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *FOOKeyinHandler) thread() {
	kh.threadStarted = true

	// TODO
	kh.exec.SendExecReadOnlyMessage("YOU ARE A DUMMY", &kh.source)
	// TODO end

	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
