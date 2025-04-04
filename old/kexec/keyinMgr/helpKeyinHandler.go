// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"strings"
	"time"

	"khalehla/kexec"
	kexec2 "khalehla/old/kexec"
)

type HELPKeyinHandler struct {
	exec    kexec.IExec
	source  kexec2.ConsoleIdentifier
	options string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewHELPKeyinHandler(exec kexec.IExec, source kexec2.ConsoleIdentifier, options string, arguments string) IKeyinHandler {
	return &HELPKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *HELPKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *HELPKeyinHandler) CheckSyntax() bool {
	return len(kh.options) == 0 && len(kh.arguments) > 0
}

func (kh *HELPKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *HELPKeyinHandler) GetCommand() string {
	return "HELP"
}

func (kh *HELPKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *HELPKeyinHandler) GetHelp() []string {
	return []string{
		"HELP keyin",
		"Displays syntax and brief help for the indicated keyin",
	}
}

func (kh *HELPKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *HELPKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *HELPKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *HELPKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *HELPKeyinHandler) thread() {
	kh.threadStarted = true

	arg := strings.ToUpper(kh.arguments)
	for token, newHandler := range handlerTable {
		if arg == token {
			handler := newHandler(kh.exec, kh.source, "", "")
			for _, str := range handler.GetHelp() {
				msg := "HELP:" + str
				kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
			}
		}
	}

	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
