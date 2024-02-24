// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"strings"
	"time"
)

type DKeyinHandler struct {
	exec            types.IExec
	source          types.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewDKeyinHandler(exec types.IExec, source types.ConsoleIdentifier, options string, arguments string) KeyinHandler {
	return &DKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *DKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *DKeyinHandler) CheckSyntax() bool {
	if len(kh.options) != 0 || len(kh.arguments) != 0 {
		return false
	}

	if len(kh.arguments) != 0 {
		return false
	}

	return true
}

func (kh *DKeyinHandler) GetCommand() string {
	return "D"
}

func (kh *DKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *DKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *DKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *DKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *DKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *DKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *DKeyinHandler) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vD KEYIN ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  threadStarted:  %v\n", indent, kh.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:  %v\n", indent, kh.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, kh.terminateThread)
}

func (kh *DKeyinHandler) thread() {
	kh.threadStarted = true

	str := "The current date and time is "
	t := time.Now()
	zoneName, _ := t.Zone()
	str += fmt.Sprintf("%v %02v %03v %04v %02v:%02v:%02v %v",
		t.Weekday(), t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), zoneName)
	kh.exec.SendExecReadOnlyMessage(str, &kh.source)

	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
