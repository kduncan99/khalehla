// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"strings"
	"sync"
	"time"
)

type keyinInfo struct {
	source types.ConsoleIdentifier
	text   string
}

// KeyinManager handles all things related to unsolicited console keyins
type KeyinManager struct {
	exec            types.IExec
	mutex           sync.Mutex
	isInitialized   bool
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	postedKeyins    []*keyinInfo
	pendingHandlers []types.KeyinHandler
}

func NewKeyinManager(exec types.IExec) *KeyinManager {
	return &KeyinManager{
		exec: exec,
	}
}

// CloseManager is invoked when the exec is stopping
func (mgr *KeyinManager) CloseManager() {
	mgr.threadStop()
	mgr.isInitialized = false
}

func (mgr *KeyinManager) InitializeManager() error {
	mgr.threadStart()
	mgr.isInitialized = true
	return nil
}

func (mgr *KeyinManager) IsInitialized() bool {
	return mgr.isInitialized
}

// ResetManager clears out any artifacts left over by a previous exec session,
// and prepares the console for normal operations
func (mgr *KeyinManager) ResetManager() error {
	mgr.threadStop()
	mgr.threadStart()
	mgr.isInitialized = true
	return nil
}

// PostKeyin queues the given keyin info and returns immediately to avoid deadlocks
func (mgr *KeyinManager) PostKeyin(source types.ConsoleIdentifier, text string) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	ki := &keyinInfo{
		source: source,
		text:   text,
	}
	mgr.postedKeyins = append(mgr.postedKeyins, ki)
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

	var handler types.KeyinHandler
	switch command {
	case "D":
		handler = NewDKeyinHandler(mgr.exec, ki.source, options, args)
	case "DU":
		handler = NewDUKeyinHandler(mgr.exec, ki.source, options, args)
	case "FS":
		handler = NewFSKeyinHandler(mgr.exec, ki.source, options, args)
	}

	if handler != nil {
		if !handler.CheckSyntax() {
			mgr.exec.SendExecReadOnlyMessage(fmt.Sprintf("%v KEYIN HAS SYNTAX ERROR, INPUT IGNORED", command))
			return
		}

		if !handler.IsAllowed() {
			mgr.exec.SendExecReadOnlyMessage(fmt.Sprintf("%v KEYIN NOT ALLOWED", command))
			return
		}

		mgr.pendingHandlers = append(mgr.pendingHandlers, handler)
		handler.Invoke()
		return
	}

	// TODO check for registered keyins

	mgr.exec.SendExecReadOnlyMessage(fmt.Sprintf("KEYIN NOT REGISTERED*%v", command))
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
	mgr.threadStarted = true

	counter := 0
	for !mgr.terminateThread {
		time.Sleep(25 * time.Millisecond)
		mgr.checkPosted()
		counter++
		if counter > 400 {
			counter = 0
			mgr.prune()
		}
	}

	mgr.threadStopped = true
}

func (mgr *KeyinManager) threadStart() {
	mgr.terminateThread = false
	if !mgr.threadStarted {
		go mgr.thread()
		for !mgr.threadStarted {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *KeyinManager) threadStop() {
	if mgr.threadStarted {
		mgr.terminateThread = true
		for !mgr.threadStopped {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *KeyinManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vKeyinManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  initialized:     %v\n", indent, mgr.isInitialized)
	_, _ = fmt.Fprintf(dest, "%v  threadStarted:   %v\n", indent, mgr.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:   %v\n", indent, mgr.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, mgr.terminateThread)

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
