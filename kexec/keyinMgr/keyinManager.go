// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"io"
	"khalehla/kexec"
	"khalehla/klog"
	"strings"
	"sync"
	"time"
)

var handlerTable = map[string]func(kexec.IExec, kexec.ConsoleIdentifier, string, string) IKeyinHandler{
	"$!": NewStopKeyinHandler,
	// "AP"
	// "AT"
	// "B"
	// "BL"
	"CJ": NewCJKeyinHandler,
	// "CS"
	"D":  NewDKeyinHandler,
	"DJ": NewDJKeyinHandler,
	// "DN"
	"DU": NewDUKeyinHandler,
	// "E"
	"FA": NewFAKeyinHandler,
	// "FB"
	// "FC" ?
	// "FF"
	"FS":   NewFSKeyinHandler,
	"HELP": NewHELPKeyinHandler,
	// "II"
	// "IN"
	// "IT"
	// "LB"
	// "LC"
	// "LG"
	// "MR"
	"MS": NewMSKeyinHandler,
	// "PM" ?
	// "PR" ?
	"PREP": NewPREPKeyinHandler,
	// "RC"
	// "RD" ?
	// "RE"
	// "RL"
	// "RM"
	// "RP" ?
	// "RS"
	// "RV"
	// "SEC"
	"SJ": NewSJKeyinHandler,
	// "SM"
	// "SP"
	// "SQ"
	// "SR"
	// "SS"
	// "ST"
	// "SU"
	// "SX"
	// "T"
	// "TB"
	// "TF"
	// "TP"
	// "TS"
	// "TU" ?
	// "UL"
	// "UP"
	// "X"
}

type keyinInfo struct {
	source kexec.ConsoleIdentifier
	text   string
}

// KeyinManager handles all things related to unsolicited console keyins
type KeyinManager struct {
	exec            kexec.IExec
	mutex           sync.Mutex
	threadDone      bool
	postedKeyins    []*keyinInfo
	pendingHandlers []IKeyinHandler
}

func NewKeyinManager(exec kexec.IExec) *KeyinManager {
	return &KeyinManager{
		exec: exec,
	}
}

// Boot is invoked when the exec is booting - return an error to stop the boot
func (mgr *KeyinManager) Boot() error {
	klog.LogTrace("KeyinMgr", "Boot")
	mgr.postedKeyins = make([]*keyinInfo, 0)
	mgr.pendingHandlers = make([]IKeyinHandler, 0)
	go mgr.thread()
	return nil
}

// Close is invoked when the application is shutting down
func (mgr *KeyinManager) Close() {
	klog.LogTrace("KeyinMgr", "Close")
	// nothing to do
}

// Initialize is invoked when the application is starting up
func (mgr *KeyinManager) Initialize() error {
	klog.LogTrace("KeyinMgr", "Initialized")
	return nil
}

// Stop is invoked when the exec is stopping
func (mgr *KeyinManager) Stop() {
	klog.LogTrace("KeyinMgr", "Stop")
	for !mgr.threadDone {
		time.Sleep(25 * time.Millisecond)
	}
}

// PostKeyin queues the given keyin info and returns immediately to avoid deadlocks
func (mgr *KeyinManager) PostKeyin(source kexec.ConsoleIdentifier, text string) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if !mgr.threadDone {
		ki := &keyinInfo{
			source: source,
			text:   text,
		}
		mgr.postedKeyins = append(mgr.postedKeyins, ki)
	}
}

func (mgr *KeyinManager) scheduleKeyinHandler(ki *keyinInfo) {
	split := strings.SplitN(ki.text, " ", 2)
	subSplit := strings.SplitN(split[0], ",", 2)
	command := strings.ToUpper(subSplit[0])
	options := ""
	if len(subSplit) == 2 {
		options = subSplit[1]
	}
	args := ""
	if len(split) == 2 {
		args = split[1]
	}

	newHandler, ok := handlerTable[command]
	if ok {
		handler := newHandler(mgr.exec, ki.source, options, args)
		if !handler.CheckSyntax() {
			mgr.exec.SendExecReadOnlyMessage(fmt.Sprintf("%v KEYIN HAS SYNTAX ERROR, INPUT IGNORED", command), nil)
			return
		}

		if !handler.IsAllowed() {
			mgr.exec.SendExecReadOnlyMessage(fmt.Sprintf("%v KEYIN NOT ALLOWED", command), nil)
			return
		}

		mgr.pendingHandlers = append(mgr.pendingHandlers, handler)
		handler.Invoke()
		return
	}

	// TODO check for registered keyins

	mgr.exec.SendExecReadOnlyMessage(fmt.Sprintf("KEYIN NOT REGISTERED*%v", command), nil)
}

func (mgr *KeyinManager) checkPosted() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if len(mgr.postedKeyins) > 0 {
		ki := mgr.postedKeyins[0]
		mgr.postedKeyins = mgr.postedKeyins[1:]
		mgr.scheduleKeyinHandler(ki)
	}
}

// thread simply prunes completed keyins from the list of handlers
func (mgr *KeyinManager) prune() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	now := time.Now()
	for phx, handler := range mgr.pendingHandlers {
		if handler.IsDone() && now.Sub(handler.GetTimeFinished()).Seconds() > 60 {
			temp := make([]IKeyinHandler, 0)
			temp = append(temp, mgr.pendingHandlers[:phx]...)
			temp = append(temp, mgr.pendingHandlers[phx+1:]...)
			mgr.pendingHandlers = temp
		} else {
			phx++
		}
	}
}

func (mgr *KeyinManager) thread() {
	mgr.threadDone = false

	counter := 0
	for !mgr.exec.GetStopFlag() {
		time.Sleep(25 * time.Millisecond)
		mgr.checkPosted()
		counter++
		if counter > 400 {
			counter = 0
			mgr.prune()
		}
	}

	mgr.threadDone = true
}

func (mgr *KeyinManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vKeyinManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  threadDone: %v\n", indent, mgr.threadDone)

	_, _ = fmt.Fprintf(dest, "%v  Recent or Pending Keyins:\n", indent)
	for _, kh := range mgr.pendingHandlers {
		var msg string
		if kh.IsDone() {
			msg = "DONE: "
		} else {
			msg = "PEND: "
		}

		msg += kh.GetCommand()
		if len(kh.GetOptions()) > 0 {
			msg += "," + kh.GetOptions()
		}
		if len(kh.GetArguments()) > 0 {
			msg += " " + kh.GetArguments()
		}
		_, _ = fmt.Fprintf(dest, "%v    %v\n", indent, msg)
	}
}
