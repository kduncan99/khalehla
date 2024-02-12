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

type FooKeyinHandler struct {
	exec            types.IExec
	source          types.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
}

func NewFooKeyinHandler(exec types.IExec, source types.ConsoleIdentifier, options string, arguments string) *FooKeyinHandler {
	return &FooKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *FooKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *FooKeyinHandler) CheckSyntax() bool {
	// Accepted:
	//		D
	//		D,UTC
	//		D SHIFT
	// TODO
	return true
}

func (kh *FooKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *FooKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *FooKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *FooKeyinHandler) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vFOO KEYIN ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  threadStarted:  %v\n", indent, kh.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:  %v\n", indent, kh.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, kh.terminateThread)
}

func (kh *FooKeyinHandler) thread() {
	kh.threadStarted = true

	// TODO
	kh.exec.SendExecReadOnlyMessage("YOU ARE A DUMMY")
	// TODO end

	kh.threadStopped = true
}
