// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"khalehla/kexec"
	"strconv"
	"strings"
	"time"
)

type SJKeyinHandler struct {
	exec            kexec.IExec
	source          kexec.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewSJKeyinHandler(exec kexec.IExec, source kexec.ConsoleIdentifier, options string, arguments string) KeyinHandler {
	return &SJKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *SJKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *SJKeyinHandler) CheckSyntax() bool {
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

func (kh *SJKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *SJKeyinHandler) GetCommand() string {
	return "SJ"
}

func (kh *SJKeyinHandler) GetHelp() []string {
	return []string{
		"SJ jumpKey,...",
		"Sets the indicated jump keys"}
}

func (kh *SJKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *SJKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *SJKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *SJKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *SJKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *SJKeyinHandler) thread() {
	kh.threadStarted = true
	split := strings.Split(kh.arguments, ",")
	for _, str := range split {
		jk, _ := strconv.Atoi(str)
		kh.exec.SetJumpKey(jk, true)
	}

	displayJumpKeys(kh.exec, &kh.source)
	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
