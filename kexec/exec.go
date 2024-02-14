// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
	"io"
	"khalehla/kexec/config"
	"khalehla/kexec/consoleMgr"
	"khalehla/kexec/facilitiesMgr"
	"khalehla/kexec/keyinMgr"
	"khalehla/kexec/mfdMgr"
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/types"
	"strings"
)

const Version = "v1.0.0"

type Exec struct {
	consoleMgr types.IConsoleManager
	facMgr     types.IFacilitiesManager
	keyinMgr   types.IKeyinManager
	mfdMgr     types.IMFDManager
	nodeMgr    types.INodeManager

	runControlEntry *types.RunControlEntry
	runControlTable map[string]*types.RunControlEntry // indexed by runid

	allowRestart bool
	phase        types.ExecPhase
	stopCode     types.StopCode
	stopFlag     bool
}

func NewExec(cfg *config.Configuration) *Exec {
	e := &Exec{}
	e.consoleMgr = consoleMgr.NewConsoleManager(e)
	e.facMgr = facilitiesMgr.NewFacilitiesManager(e)
	e.keyinMgr = keyinMgr.NewKeyinManager(e)
	e.mfdMgr = mfdMgr.NewMFDManager(e)
	e.nodeMgr = nodeMgr.NewNodeManager(e)
	e.phase = types.ExecPhaseNotStarted

	// ExecRunControlEntry is the RCE for the EXEC - it always exists and is always (or should always be) in the RCT
	e.runControlEntry = types.NewRunControlEntry(
		cfg.SystemRunId,
		cfg.SystemRunId,
		cfg.SystemAccountId,
		cfg.SystemProjectId,
		cfg.SystemUserId)
	e.runControlEntry.DefaultQualifier = cfg.SystemQualifier
	e.runControlEntry.ImpliedQualifier = cfg.SystemQualifier
	e.runControlEntry.IsExec = true

	e.runControlTable = make(map[string]*types.RunControlEntry)
	e.runControlTable[e.runControlEntry.RunId] = e.runControlEntry

	return e
}

func (e *Exec) Close() {
	e.keyinMgr.CloseManager()
	e.facMgr.CloseManager()
	e.mfdMgr.CloseManager()
	e.nodeMgr.CloseManager()
	e.consoleMgr.CloseManager()
}

func (e *Exec) GetConsoleManager() types.IConsoleManager {
	return e.consoleMgr
}

func (e *Exec) GetFacilitiesManager() types.IFacilitiesManager {
	return e.facMgr
}

func (e *Exec) GetKeyinManager() types.IKeyinManager {
	return e.keyinMgr
}

func (e *Exec) GetMFDManager() types.IMFDManager {
	return e.mfdMgr
}

func (e *Exec) GetNodeManager() types.INodeManager {
	return e.nodeMgr
}

func (e *Exec) GetPhase() types.ExecPhase {
	return e.phase
}

func (e *Exec) GetRunControlEntry() *types.RunControlEntry {
	return e.runControlEntry
}

func (e *Exec) GetStopCode() types.StopCode {
	return e.stopCode
}

func (e *Exec) GetStopFlag() bool {
	return e.stopFlag
}

func (e *Exec) HandleKeyIn(source types.ConsoleIdentifier, text string) {
	e.keyinMgr.PostKeyin(source, text)
}

func (e *Exec) InitialBoot(initMassStorage bool) error {
	e.phase = types.ExecPhaseInitializing

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
	accepted := []string{"DONE"}
	_, err = e.SendExecRestrictedReadReplyMessage("Modify Config then answer DONE", accepted)
	if err != nil {
		return err
	}

	// spin up facilities, then the MFD
	err = e.facMgr.InitializeManager()
	if err != nil {
		return err
	}

	e.mfdMgr.SetMSInitialize(true)
	err = e.mfdMgr.InitializeManager()
	if err != nil {
		return err
	}

	e.allowRestart = false // TODO temporary
	return nil
}

func (e *Exec) RecoveryBoot(initMassStorage bool) error {
	// TODO
	return nil
}

func (e *Exec) SendExecReadOnlyMessage(message string) {
	consMsg := types.ConsoleReadOnlyMessage{
		Source:         e.runControlEntry,
		Text:           message,
		DoNotEmitRunId: true,
	}
	e.consoleMgr.SendReadOnlyMessage(&consMsg)
}

func (e *Exec) SendExecReadReplyMessage(message string, maxReplyChars int) (string, error) {
	consMsg := types.ConsoleReadReplyMessage{
		Source:         e.runControlEntry,
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

func (e *Exec) SendExecRestrictedReadReplyMessage(message string, accepted []string) (string, error) {
	if len(accepted) == 0 {
		return "", fmt.Errorf("bad accepted list")
	}

	maxReplyLen := 0
	for _, acceptString := range accepted {
		if maxReplyLen < len(acceptString) {
			maxReplyLen = len(acceptString)
		}
	}

	consMsg := types.ConsoleReadReplyMessage{
		Source:         e.runControlEntry,
		Text:           message,
		DoNotEmitRunId: true,
		MaxReplyLength: maxReplyLen,
	}

	done := false
	for !done {
		err := e.consoleMgr.SendReadReplyMessage(&consMsg)
		if err != nil {
			return "", err
		}

		resp := strings.ToUpper(consMsg.Reply)
		for _, acceptString := range accepted {
			if acceptString == resp {
				done = true
				break
			}
		}
	}

	return strings.ToUpper(consMsg.Reply), nil
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
	e.phase = types.ExecPhaseStopped
}

func (e *Exec) Dump(dest io.Writer) {
	_, _ = fmt.Fprintf(dest, "Exec Dump ----------------------------------------------------\n")

	_, _ = fmt.Fprintf(dest, "  Phase:         %v\n", e.phase)
	_, _ = fmt.Fprintf(dest, "  Stopped:       %v\n", e.stopFlag)
	_, _ = fmt.Fprintf(dest, "  StopCode:      %03o\n", e.stopCode)
	_, _ = fmt.Fprintf(dest, "  Allow Restart: %v\n", e.allowRestart)

	e.consoleMgr.Dump(dest, "")
	e.keyinMgr.Dump(dest, "")
	e.nodeMgr.Dump(dest, "")
	e.facMgr.Dump(dest, "")
	e.mfdMgr.Dump(dest, "")

	// TODO run control table, etc
}
