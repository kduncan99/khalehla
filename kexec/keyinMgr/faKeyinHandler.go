// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"khalehla/hardware"
	"khalehla/kexec"
	"khalehla/kexec/facilitiesMgr"
	"khalehla/kexec/nodeMgr"
	"strings"
	"time"
)

type FAKeyinHandler struct {
	exec            kexec.IExec
	source          kexec.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	timeFinished    time.Time
}

func NewFAKeyinHandler(exec kexec.IExec, source kexec.ConsoleIdentifier, options string, arguments string) IKeyinHandler {
	return &FAKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *FAKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *FAKeyinHandler) CheckSyntax() bool {
	return len(kh.options) == 0 && len(kh.arguments) > 0
}

func (kh *FAKeyinHandler) GetArguments() string {
	return kh.arguments
}

func (kh *FAKeyinHandler) GetCommand() string {
	return "FA"
}

func (kh *FAKeyinHandler) GetHelp() []string {
	return []string{
		"FA device",
		"Forces unit attention for the given device",
		"Used primarily after online prep of a disk unit"}
}

func (kh *FAKeyinHandler) GetOptions() string {
	return kh.options
}

func (kh *FAKeyinHandler) GetTimeFinished() time.Time {
	return kh.timeFinished
}

func (kh *FAKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *FAKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *FAKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *FAKeyinHandler) process() {
	nm := kh.exec.GetNodeManager().(*nodeMgr.NodeManager)
	devName := strings.ToUpper(kh.GetArguments())
	nodeInfo, err := nm.GetNodeInfoByName(devName)
	if err != nil {
		msg := fmt.Sprintf("%v not found", devName)
		kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
		return
	}

	if nodeInfo.GetNodeCategoryType() != hardware.NodeCategoryDevice {
		msg := fmt.Sprintf("%v is not a Device", devName)
		kh.exec.SendExecReadOnlyMessage(msg, &kh.source)
		return
	}

	fm := kh.exec.GetFacilitiesManager().(*facilitiesMgr.FacilitiesManager)
	fm.NotifyDeviceReady(nodeInfo.GetNodeIdentifier(), true)
}

func (kh *FAKeyinHandler) thread() {
	kh.threadStarted = true
	kh.process()
	kh.threadStopped = true
	kh.timeFinished = time.Now()
}
