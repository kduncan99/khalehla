// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"sync"
	"time"
)

type MFDManager struct {
	exec            types.IExec
	mutex           sync.Mutex
	isInitialized   bool
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
}

func NewMFDManager(exec types.IExec) *MFDManager {
	return &MFDManager{
		exec: exec,
	}
}

// CloseManager is invoked when the exec is stopping
func (mgr *MFDManager) CloseManager() {
	mgr.threadStop()
	mgr.isInitialized = false
}

func (mgr *MFDManager) InitializeManager() error {
	// TODO

	mgr.threadStart()
	mgr.isInitialized = true
	return nil
}

func (mgr *MFDManager) IsInitialized() bool {
	return mgr.isInitialized
}

// ResetManager clears out any artifacts left over by a previous exec session,
// and prepares the console for normal operations
func (mgr *MFDManager) ResetManager() error {
	mgr.threadStop()
	mgr.threadStart()
	mgr.isInitialized = true
	return nil
}

func (mgr *MFDManager) NotifyDeviceReady(deviceInfo types.DeviceInfo, isReady bool) {
	// TODO
}

func (mgr *MFDManager) thread() {
	mgr.threadStarted = true

	for !mgr.terminateThread {
		time.Sleep(25 * time.Millisecond)
		// TODO
	}

	mgr.threadStopped = true
}

func (mgr *MFDManager) threadStart() {
	mgr.terminateThread = false
	if !mgr.threadStarted {
		go mgr.thread()
		for !mgr.threadStarted {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *MFDManager) threadStop() {
	if mgr.threadStarted {
		mgr.terminateThread = true
		for !mgr.threadStopped {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *MFDManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vMFDManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  isInitialized:   %v\n", indent, mgr.isInitialized)
	_, _ = fmt.Fprintf(dest, "%v  threadStarted:   %v\n", indent, mgr.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:   %v\n", indent, mgr.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, mgr.terminateThread)

	// TODO

}
