// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
	"io"
	"khalehla/kexec/consoleMgr"
	"khalehla/kexec/deviceMgr"
	"khalehla/kexec/keyinMgr"
	"khalehla/kexec/types"
	"khalehla/pkg"
	"strings"
)

const Version = "v1.0.0"

type Exec struct {
	consoleMgr *consoleMgr.ConsoleManager
	deviceMgr  *deviceMgr.DeviceManager
	keyinMgr   *keyinMgr.KeyinManager

	runControlTable map[pkg.Word36]*types.RunControlEntry

	allowRestart bool
	stopCode     types.StopCode
	stopFlag     bool
}

func NewExec() *Exec {
	e := &Exec{}
	e.consoleMgr = consoleMgr.NewConsoleManager(e)
	e.deviceMgr = deviceMgr.NewDeviceManager(e)
	e.keyinMgr = keyinMgr.NewKeyinManager(e)

	return e
}

func (e *Exec) Close() {
	e.deviceMgr.CloseManager()
	e.consoleMgr.CloseManager()
}

func (e *Exec) GetConsoleManager() types.Manager {
	return e.consoleMgr
}

func (e *Exec) GetDeviceManager() types.Manager {
	return e.deviceMgr
}

func (e *Exec) GetKeyinManager() types.Manager {
	return e.keyinMgr
}

func (e *Exec) GetStopFlag() bool {
	return e.stopFlag
}

func (e *Exec) HandleKeyIn(source types.ConsoleIdentifier, text string) {
	e.keyinMgr.PostKeyin(source, text)
}

func (e *Exec) InitialBoot(initMassStorage bool) error {
	e.consoleMgr.InitializeManager()
	e.deviceMgr.InitializeManager()
	e.keyinMgr.InitializeManager()

	e.SendExecReadOnlyMessage("KEXEC Startup - Version " + Version)
	e.SendExecReadOnlyMessage("Building Configuration...")
	err := e.deviceMgr.BuildConfiguration()
	if err != nil {
		e.SendExecReadOnlyMessage("Error:" + err.Error())
		e.Stop(1) // TODO put a real stop code here - error in configuration
		return fmt.Errorf("boot failed")
	}

	reply := ""
	err = nil
	for strings.ToUpper(reply) != "DONE" && err == nil {
		reply, err = e.SendExecReadReplyMessage("Modify Config then answer DONE", 4)
	}

	if err != nil {
		// TODO DIE HORRIBLY
	}

	e.allowRestart = false // TODO temporary
	e.Stop(063)            // TODO temporary
	return nil
}

func (e *Exec) RecoveryBoot(initMassStorage bool) error {
	// TODO
	return nil
}

func (e *Exec) SendExecReadOnlyMessage(message string) {
	consMsg := types.ConsoleReadOnlyMessage{
		Source:         &types.ExecRunControlEntry,
		Text:           message,
		DoNotEmitRunId: true,
	}
	e.consoleMgr.SendReadOnlyMessage(&consMsg)
}

func (e *Exec) SendExecReadReplyMessage(message string, maxReplyChars int) (string, error) {
	consMsg := types.ConsoleReadReplyMessage{
		Source:         &types.ExecRunControlEntry,
		Text:           message,
		DoNotEmitRunId: true,
		MaxReplyLength: maxReplyChars,
	}

	err := e.consoleMgr.SendReadReplyMessage(&consMsg)
	if err != nil {
		return "", err
	}

	return consMsg.Reply, nil
}

func (e *Exec) Stop(code types.StopCode) {
	// TODO need to set contingency in the Exec RCE
	if e.allowRestart {
		e.SendExecReadOnlyMessage(fmt.Sprintf("Restarting Exec: Status Code %03o", code))
	} else {
		e.SendExecReadOnlyMessage(fmt.Sprintf("Stopping Exec: Status Code %03o", code))
	}

	e.stopFlag = true
	e.stopCode = code
}

func (e *Exec) Dump(dest io.Writer) {
	_, _ = fmt.Fprintf(dest, "Exec Dump ----------------------------------------------------\n")

	_, _ = fmt.Fprintf(dest, "  Stopped:       %v\n", e.stopFlag)
	_, _ = fmt.Fprintf(dest, "  StopCode:      %03o\n", e.stopCode)
	_, _ = fmt.Fprintf(dest, "  Allow Restart: %v\n", e.allowRestart)

	e.consoleMgr.Dump(dest, "")
	e.deviceMgr.Dump(dest, "")
	e.keyinMgr.Dump(dest, "")

	// TODO run control table, etc
}
