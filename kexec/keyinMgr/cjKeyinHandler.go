// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"khalehla/kexec/types"
	"strconv"
	"strings"
	"time"
)

type CJKeyinHandler struct {
	exec            types.IExec
	source          types.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewCJKeyinHandler(exec types.IExec, source types.ConsoleIdentifier, options string, arguments string) KeyinHandler {
	return &CJKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *CJKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *CJKeyinHandler) CheckSyntax() bool {
	if len(kh.options) > 0 || len(kh.arguments) == 0 {
		return false
	}

	split := strings.Split(kh.arguments, ",")
	for _, str := range split {
		jk, err := strconv.Atoi(str)
		if err != nil || jk < 1 || jk > 36 {
			return false
		}
	}
	return true
}

func (kh *CJKeyinHandler) GetCommand() string {
	return "CJ"
}

func (kh *CJKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *CJKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *CJKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *CJKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *CJKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *CJKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *CJKeyinHandler) thread() {
	kh.threadStarted = true
	split := strings.Split(kh.arguments, ",")
	for _, str := range split {
		jk, _ := strconv.Atoi(str)
		kh.exec.SetJumpKey(jk, false)
	}

	displayJumpKeys(kh.exec, &kh.source)
	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
