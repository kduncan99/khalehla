// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
	"khalehla/kexec/consoleMgr"
	"khalehla/kexec/deviceMgr"
	"khalehla/kexec/types"
	"khalehla/pkg"
)

const Version = "v1.0.0"

type Exec struct {
	consoleMgr      *consoleMgr.ConsoleManager
	deviceMgr       *deviceMgr.DeviceManager
	runControlTable map[pkg.Word36]*types.RunControlEntry

	allowRestart bool
	stopCode     int
	stopFlag     bool
}

func (e *Exec) InitialBoot(initMassStorage bool) error {
	e.consoleMgr = &consoleMgr.ConsoleManager{}
	e.deviceMgr = &deviceMgr.DeviceManager{}

	e.consoleMgr.Reset()
	e.SendExecReadOnlyMessage("KEXEC Startup - Version " + Version)
	e.SendExecReadOnlyMessage("Building Configuration...")
	err := e.deviceMgr.BuildConfiguration()
	if err != nil {
		e.SendExecReadOnlyMessage("Error:" + err.Error())
		e.Stop(1) // TODO put a real stop code here - error in configuration
		return fmt.Errorf("boot failed")
	}

	// TODO

	e.allowRestart = false // TODO temporary
	e.Stop(063)            // TODO temporary
	return nil
}

func (e *Exec) RecoveryBoot(initMassStorage bool) error {
	// TODO
	return nil
}

func (e *Exec) SendExecReadOnlyMessage(message string) {
	e.consoleMgr.SendReadOnlyMessage(&types.ExecRunControlEntry, message)
}

func (e *Exec) Stop(code int) {
	if e.allowRestart {
		e.SendExecReadOnlyMessage(fmt.Sprintf("Restarting Exec: Status Code %03o", code))
	} else {
		e.SendExecReadOnlyMessage(fmt.Sprintf("Stopping Exec: Status Code %03o", code))
	}

	e.stopFlag = true
	e.stopCode = code
}
