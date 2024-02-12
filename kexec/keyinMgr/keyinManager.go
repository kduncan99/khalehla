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
	postedKeyins    []*keyinInfo
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
	keyinHandlers   []types.KeyinHandler
}

func NewKeyinManager(exec types.IExec) *KeyinManager {
	return &KeyinManager{
		exec: exec,
	}
}

// CloseManager is invoked when the exec is stopping
func (mgr *KeyinManager) CloseManager() {
	mgr.threadStop()
}

func (mgr *KeyinManager) InitializeManager() {
	mgr.threadStart()
}

// ResetManager clears out any artifacts left over by a previous exec session,
// and prepares the console for normal operations
func (mgr *KeyinManager) ResetManager() {
	mgr.threadStop()
	mgr.threadStart()
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

func (mgr *KeyinManager) handleKeyin(source types.ConsoleIdentifier, text string) {
	split := strings.SplitN(text, " ", 2)
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

	var kh types.KeyinHandler
	switch command {
	case "D":
		kh = NewDKeyinHandler(mgr.exec, source, options, args)
	}

	if kh != nil {
		if !kh.CheckSyntax() {
			mgr.exec.SendExecReadOnlyMessage(fmt.Sprintf("%v KEYIN HAS SYNTAX ERROR, INPUT IGNORED", command))
			return
		}

		if !kh.IsAllowed() {
			mgr.exec.SendExecReadOnlyMessage(fmt.Sprintf("%v KEYIN NOT ALLOWED", command))
			return
		}

		mgr.keyinHandlers = append(mgr.keyinHandlers, kh)
		kh.Invoke()
		return
	}

	// TODO check for registered keyins

	mgr.exec.SendExecReadOnlyMessage(fmt.Sprintf("KEYIN NOT REGISTERED*%v", command))
}

func (mgr *KeyinManager) checkPosted() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if len(mgr.postedKeyins) > 0 {
		top := mgr.postedKeyins[0]
		mgr.postedKeyins = mgr.postedKeyins[1:]
		mgr.handleKeyin(top.source, top.text)
	}
}

// thread simply prunes completed keyins from the list of handlers
func (mgr *KeyinManager) prune() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for len(mgr.keyinHandlers) > 0 {
		front := mgr.keyinHandlers[0]
		if !front.IsDone() {
			break
		}

		mgr.keyinHandlers = mgr.keyinHandlers[1:]
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

	_, _ = fmt.Fprintf(dest, "%v  threadStarted:  %v\n", indent, mgr.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:  %v\n", indent, mgr.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, mgr.terminateThread)
}
