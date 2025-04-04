// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"strings"
	"time"

	"khalehla/kexec"
	kexec2 "khalehla/old/kexec"
)

type DJKeyinHandler struct {
	exec    kexec.IExec
	source  kexec2.ConsoleIdentifier
	options string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewDJKeyinHandler(exec kexec.IExec, source kexec2.ConsoleIdentifier, options string, arguments string) IKeyinHandler {
	return &DJKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *DJKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *DJKeyinHandler) CheckSyntax() bool {
	return len(kh.options) == 0 && len(kh.arguments) == 0
}

func (kh *DJKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *DJKeyinHandler) GetCommand() string {
	return "DJ"
}

func (kh *DJKeyinHandler) GetHelp() []string {
	return []string{
		"DJ",
		"Displays the jump keys which are currently set"}
}

func (kh *DJKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *DJKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *DJKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *DJKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *DJKeyinHandler) IsAllowed() bool {
	return true
}

func displayJumpKeys(exec kexec.IExec, source *kexec2.ConsoleIdentifier) {
	str := ""
	for jk := 1; jk <= 36; jk++ {
		jkSet := exec.GetJumpKey(jk)
		if jkSet {
			if len(str) == 0 {
				str = fmt.Sprintf("Jump Keys Set: %v", jk)
			} else {
				str += fmt.Sprintf(",%v", jk)
				if len(str) > 60 {
					exec.SendExecReadOnlyMessage(str, source)
					str = ""
				}
			}
		}
	}

	if len(str) > 0 {
		exec.SendExecReadOnlyMessage(str, source)
	}
}

func (kh *DJKeyinHandler) thread() {
	kh.threadStarted = true
	displayJumpKeys(kh.exec, &kh.source)
	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
