// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"strings"
)

type DKeyinHandler struct {
	exec            types.IExec
	source          types.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
}

func NewDKeyinHandler(exec types.IExec, source types.ConsoleIdentifier, options string, arguments string) *DKeyinHandler {
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
	// Accepted:
	//		D
	//		D,UTC
	//		D SHIFT
	if len(kh.options) != 0 {
		if kh.options != "UTC" || len(kh.arguments) != 0 {
			return false
		}
	}

	if len(kh.arguments) != 0 && kh.arguments != "SHIFT" {
		return false
	}

	return true
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

	// TODO
	kh.exec.SendExecReadOnlyMessage("YOU ARE A DUMMY")
	// TODO end

	kh.threadStopped = true
}
