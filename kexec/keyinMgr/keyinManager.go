// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"log"
	"strings"
	"sync"
	"time"
)

var handlerTable = map[string]func(types.IExec, types.ConsoleIdentifier, string, string) KeyinHandler{
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
	// "FA"
	// "FB"
	// "FC" ?
	// "FF"
	"FS": NewFSKeyinHandler,
	// "II"
	// "IN"
	// "IT"
	// "LB"
	// "LC"
	// "LG"
	// "MR"
	// "MS"
	// "PM" ?
	// "PR" ?
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
	source types.ConsoleIdentifier
	text   string
}

// KeyinManager handles all things related to unsolicited console keyins
type KeyinManager struct {
	exec            types.IExec
	mutex           sync.Mutex
	threadDone      bool
	postedKeyins    []*keyinInfo
	pendingHandlers []KeyinHandler
}

func NewKeyinManager(exec types.IExec) *KeyinManager {
	return &KeyinManager{
		exec: exec,
	}
}

// Boot is invoked when the exec is booting - return an error to stop the boot
func (mgr *KeyinManager) Boot() error {
	log.Printf("KeyinMgr:Boot")
	mgr.postedKeyins = make([]*keyinInfo, 0)
	mgr.pendingHandlers = make([]KeyinHandler, 0)
	go mgr.thread()
	return nil
}

// Close is invoked when the application is shutting down
func (mgr *KeyinManager) Close() {
	log.Printf("KeyinMgr:Close")
	// nothing to do
}

// Initialize is invoked when the application is starting up
func (mgr *KeyinManager) Initialize() error {
	log.Printf("KeyinMgr:Initialized")
	return nil
}

// Stop is invoked when the exec is stopping
func (mgr *KeyinManager) Stop() {
	log.Printf("KeyinMgr:Stop")
	for !mgr.threadDone {
		time.Sleep(25 * time.Millisecond)
	}
}

// PostKeyin queues the given keyin info and returns immediately to avoid deadlocks
func (mgr *KeyinManager) PostKeyin(source types.ConsoleIdentifier, text string) {
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
			mgr.pendingHandlers = append(mgr.pendingHandlers[:phx], mgr.pendingHandlers[phx+1:]...)
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
