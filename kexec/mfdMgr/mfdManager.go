// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"khalehla/pkg"
	"log"
	"sync"
	"time"
)

type MFDManager struct {
	exec                         types.IExec
	mutex                        sync.Mutex
	isInitialized                bool
	terminateThread              bool
	threadStarted                bool
	threadStopped                bool
	msInitialize                 bool
	deviceReadyNotificationQueue map[types.DeviceIdentifier]bool
	directoryTracks              map[uint64][]pkg.Word36 // key is MFD-relative sector address
	fixedLDAT                    map[uint]types.DeviceIdentifier
}

func NewMFDManager(exec types.IExec) *MFDManager {
	return &MFDManager{
		exec:                         exec,
		deviceReadyNotificationQueue: make(map[types.DeviceIdentifier]bool),
		directoryTracks:              make(map[uint64][]pkg.Word36),
		fixedLDAT:                    make(map[uint]types.DeviceIdentifier),
	}
}

// CloseManager is invoked when the exec is stopping
func (mgr *MFDManager) CloseManager() {
	mgr.threadStop()
	mgr.isInitialized = false
}

func (mgr *MFDManager) InitializeManager() error {
	var err error
	if mgr.msInitialize {
		replies := []string{"Y", "N"}
		msg := "Mass Storage will be Initialized - Do You Want To Continue? Y/N"
		reply, err := mgr.exec.SendExecRestrictedReadReplyMessage(msg, replies)
		if err != nil {
			return err
		} else if reply != "Y" {
			mgr.exec.Stop(types.StopConsoleResponseRequiresReboot)
			return fmt.Errorf("boot canceled")
		}

		err = mgr.initializeMassStorage()
	} else {
		err = mgr.recoverMassStorage()
	}

	if err != nil {
		log.Println("MFDMgr:Cannot continue boot")
		return err
	}

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
	// post it, and let the tread deal with it later
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.deviceReadyNotificationQueue[deviceInfo.GetDeviceIdentifier()] = isReady
}

// SetMSInitialize sets or clears the flag which indicates whether to initialze
// mass-storage upon initialization. Invoke this before calling Initialize.
func (mgr *MFDManager) SetMSInitialize(flag bool) {
	mgr.msInitialize = flag
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

	_, _ = fmt.Fprintf(dest, "%v  Queued device-ready notifications:\n", indent)
	for devId, ready := range mgr.deviceReadyNotificationQueue {
		wId := pkg.Word36(devId)
		_, _ = fmt.Fprintf(dest, "%v    devId:0%v ready:%v\n", indent, wId.ToStringAsOctal(), ready)
	}
}
