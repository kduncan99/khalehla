// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"khalehla/kexec"
	"log"
	"strings"
	"time"
)

type DUKeyinHandler struct {
	exec            kexec.IExec
	source          kexec.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewDUKeyinHandler(exec kexec.IExec, source kexec.ConsoleIdentifier, options string, arguments string) KeyinHandler {
	return &DUKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *DUKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *DUKeyinHandler) CheckSyntax() bool {
	return kh.arguments == "MP" && len(kh.options) == 0
}

func (kh *DUKeyinHandler) GetCommand() string {
	return "DU"
}

func (kh *DUKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *DUKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *DUKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *DUKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *DUKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *DUKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *DUKeyinHandler) thread() {
	kh.threadStarted = true

	fileName, err := kh.exec.PerformDump(true)
	if err != nil {
		msg := fmt.Sprintf("DU Keyin - %v", err)
		log.Printf("DU:%s\n", msg)
		kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
		return
	}

	msg := "DU Keyin Wrote dump to " + fileName
	kh.exec.SendExecReadOnlyMessage(msg, &kh.source)

	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
