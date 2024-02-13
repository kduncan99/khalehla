// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
	"io"
	"khalehla/kexec/consoleMgr"
	"khalehla/kexec/facilitiesMgr"
	"khalehla/kexec/keyinMgr"
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/types"
	"khalehla/pkg"
	"strings"
)

const Version = "v1.0.0"

type Exec struct {
	consoleMgr *consoleMgr.ConsoleManager
	facMgr     *facilitiesMgr.FacilitiesManager
	keyinMgr   *keyinMgr.KeyinManager
	nodeMgr    *nodeMgr.NodeManager

	runControlTable map[pkg.Word36]*types.RunControlEntry

	allowRestart bool
	stopCode     types.StopCode
	stopFlag     bool
}

func NewExec() *Exec {
	e := &Exec{}
	e.consoleMgr = consoleMgr.NewConsoleManager(e)
	e.keyinMgr = keyinMgr.NewKeyinManager(e)
	e.nodeMgr = nodeMgr.NewNodeManager(e)

	return e
}

func (e *Exec) Close() {
	e.keyinMgr.CloseManager()
	e.nodeMgr.CloseManager()
	e.consoleMgr.CloseManager()
}

func (e *Exec) GetConsoleManager() types.Manager {
	return e.consoleMgr
}

func (e *Exec) GetFacilitiesManager() types.Manager {
	return e.facMgr
}

func (e *Exec) GetKeyinManager() types.Manager {
	return e.keyinMgr
}

func (e *Exec) GetNodeManager() types.Manager {
	return e.nodeMgr
}

func (e *Exec) GetStopFlag() bool {
	return e.stopFlag
}

func (e *Exec) HandleKeyIn(source types.ConsoleIdentifier, text string) {
	e.keyinMgr.PostKeyin(source, text)
}

func (e *Exec) InitialBoot(initMassStorage bool) error {
	// we need the console before anything else, and then the keyin manager right after that
	err := e.consoleMgr.InitializeManager()
	if err != nil {
		return err
	}

	err = e.keyinMgr.InitializeManager()
	if err != nil {
		return err
	}

	e.SendExecReadOnlyMessage("KEXEC Startup - Version " + Version)

	// now let's have the disks and tapes
	e.SendExecReadOnlyMessage("Building Configuration...")
	err = e.nodeMgr.InitializeManager()
	if err != nil {
		return err
	}

	// Let the operator adjust the configuration
	reply := ""
	for strings.ToUpper(reply) != "DONE" && err == nil {
		reply, err = e.SendExecReadReplyMessage("Modify Config then answer DONE", 4)
		if err != nil {
			return err
		}
	}

	// spin up facilities
	err = e.facMgr.InitializeManager()
	if err != nil {
		return err
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
	e.keyinMgr.Dump(dest, "")
	e.nodeMgr.Dump(dest, "")
	e.facMgr.Dump(dest, "")

	// TODO run control table, etc
}
