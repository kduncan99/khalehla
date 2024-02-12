// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"sync"
	"time"
)

type FooManager struct {
	exec            types.IExec
	mutex           sync.Mutex
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
}

func NewFooManager(exec types.IExec) *FooManager {
	return &FooManager{
		exec: exec,
	}
}

// CloseManager is invoked when the exec is stopping
func (mgr *FooManager) CloseManager() {
	mgr.threadStop()
}

func (mgr *FooManager) InitializeManager() {
	mgr.threadStart()
}

// ResetManager clears out any artifacts left over by a previous exec session,
// and prepares the console for normal operations
func (mgr *FooManager) ResetManager() {
	mgr.threadStop()
	mgr.threadStart()
}

func (mgr *FooManager) thread() {
	mgr.threadStarted = true

	for !mgr.terminateThread {
		time.Sleep(25 * time.Millisecond)
		// TODO
	}

	mgr.threadStopped = true
}

func (mgr *FooManager) threadStart() {
	mgr.terminateThread = false
	if !mgr.threadStarted {
		go mgr.thread()
		for !mgr.threadStarted {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *FooManager) threadStop() {
	if mgr.threadStarted {
		mgr.terminateThread = true
		for !mgr.threadStopped {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *FooManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vFooManager ----------------------------------------------------\n", indent)

	// TODO

	_, _ = fmt.Fprintf(dest, "%v  threadStarted:  %v\n", indent, mgr.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:  %v\n", indent, mgr.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, mgr.terminateThread)
}
