// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"khalehla/kexec/types"
	"log"
	"os"
	"strings"
	"time"
)

type DUKeyinHandler struct {
	exec            types.IExec
	source          types.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewDUKeyinHandler(exec types.IExec, source types.ConsoleIdentifier, options string, arguments string) *DUKeyinHandler {
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
	return kh.options == "MP" && len(kh.arguments) == 0
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

	now := time.Now()
	fileName := fmt.Sprintf("kexecDump-%04v%02v%02v-%02v%02v%02v.log",
		now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	dumpFile, err := os.Create(fileName)
	if err != nil {
		msg := "DU,MP Failed - Cannot create log file"
		log.Printf("DU:%s\n", msg)
		kh.exec.SendExecReadOnlyMessage(msg)
		return
	}

	defer func() {
		if err := dumpFile.Close(); err != nil {
			msg := "DU,MP Failed - Error closing log file"
			log.Printf("DU:%s\n", msg)
			kh.exec.SendExecReadOnlyMessage(msg)
			return
		}
	}()

	kh.exec.Dump(dumpFile)
	msg := "DU,MP Wrote dump to " + fileName
	kh.exec.SendExecReadOnlyMessage(msg)

	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
