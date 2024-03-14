// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package exec

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/config"
	"khalehla/kexec/consoleMgr"
	"khalehla/kexec/facilitiesMgr"
	"khalehla/kexec/keyinMgr"
	"khalehla/kexec/mfdMgr"
	"khalehla/kexec/nodeMgr"
	"log"
	"os"
	"strings"
	"time"
)

const Version = "v1.0.0"

type Exec struct {
	configuration *config.Configuration
	consoleMgr    kexec.IConsoleManager
	facMgr        kexec.IFacilitiesManager
	keyinMgr      kexec.IKeyinManager
	mfdMgr        kexec.IMFDManager
	nodeMgr       kexec.INodeManager

	runControlEntry *kexec.RunControlEntry
	runControlTable map[string]*kexec.RunControlEntry // indexed by runid

	jumpKeys []bool
	phase    kexec.ExecPhase
	stopCode kexec.StopCode
	stopFlag bool
}

func NewExec(cfg *config.Configuration) *Exec {
	e := &Exec{}
	e.configuration = cfg

	e.consoleMgr = consoleMgr.NewConsoleManager(e)
	e.facMgr = facilitiesMgr.NewFacilitiesManager(e)
	e.keyinMgr = keyinMgr.NewKeyinManager(e)
	e.mfdMgr = mfdMgr.NewMFDManager(e)
	e.nodeMgr = nodeMgr.NewNodeManager(e)
	e.phase = kexec.ExecPhaseNotStarted

	return e
}

// Boot starts and runs the system.
// It returns only when we are completely done, not just rebooting.
func (e *Exec) Boot(session uint, jumpKeys []bool, invokerChannel chan kexec.StopCode) {
	e.jumpKeys = jumpKeys
	e.stopFlag = false
	e.phase = kexec.ExecPhaseInitializing

	// ExecRunControlEntry is the RCE for the EXEC - it always exists and is always (or should always be) in the RCT
	e.runControlEntry = kexec.NewExecRunControlEntry(e.configuration.MasterAccountId)
	e.runControlTable = make(map[string]*kexec.RunControlEntry)
	e.runControlTable[e.runControlEntry.RunId] = e.runControlEntry

	managers := []kexec.IManager{
		e.consoleMgr,
		e.keyinMgr,
		e.nodeMgr,
		e.facMgr,
		e.mfdMgr,
	}

	// Boot the various managers
	for _, m := range managers {
		_ = m.Boot()
		if e.stopFlag {
			invokerChannel <- e.stopCode
			return
		}
	}

	// Begin the real boot process
	e.SendExecReadOnlyMessage("KEXEC Startup - Version "+Version, nil)

	if session == 0 || e.jumpKeys[kexec.JumpKey1Index] {
		// Let the operator adjust the configuration
		accepted := []string{"DONE"}
		_, _ = e.SendExecRestrictedReadReplyMessage("Modify Config then answer DONE", accepted, nil)
		if e.stopFlag {
			invokerChannel <- e.stopCode
			return
		}
	}

	if session == 0 {
		e.performInitialBoot()
	} else {
		e.performRecoveryBoot()
	}

	// Now wait for someone to stop us
	for !e.stopFlag {
		time.Sleep(25 * time.Millisecond)
	}

	// Stop the manager, then tell the invoker we are done
	for _, m := range managers {
		m.Stop()
	}

	invokerChannel <- e.stopCode
}

// Close invokes the Close method on each of the managers in a particular order.
func (e *Exec) Close() {
	log.Printf("Exec:Close")
	managers := []kexec.IManager{
		e.mfdMgr,
		e.facMgr,
		e.nodeMgr,
		e.keyinMgr,
		e.consoleMgr,
	}

	for _, m := range managers {
		m.Close()
	}
}

func (e *Exec) GetConfiguration() *config.Configuration {
	return e.configuration
}

func (e *Exec) GetConsoleManager() kexec.IConsoleManager {
	return e.consoleMgr
}

func (e *Exec) GetFacilitiesManager() kexec.IFacilitiesManager {
	return e.facMgr
}

func (e *Exec) GetJumpKey(jkNumber int) bool {
	return (jkNumber >= 1 && jkNumber <= 36) && e.jumpKeys[jkNumber-1]
}

func (e *Exec) GetKeyinManager() kexec.IKeyinManager {
	return e.keyinMgr
}

func (e *Exec) GetMFDManager() kexec.IMFDManager {
	return e.mfdMgr
}

func (e *Exec) GetNodeManager() kexec.INodeManager {
	return e.nodeMgr
}

func (e *Exec) GetPhase() kexec.ExecPhase {
	return e.phase
}

func (e *Exec) GetRunControlEntry() *kexec.RunControlEntry {
	return e.runControlEntry
}

func (e *Exec) GetStopCode() kexec.StopCode {
	return e.stopCode
}

func (e *Exec) GetStopFlag() bool {
	return e.stopFlag
}

// Initialize invokes the Initialize method on each of the managers in a particular order.
// If any of them return an error, we pass that error back the the caller which should Close() us and terminate.
// Should be invoked after calling NewExec(), but before calling Boot()
func (e *Exec) Initialize() error {
	managers := []kexec.IManager{
		e.consoleMgr,
		e.keyinMgr,
		e.nodeMgr,
		e.facMgr,
		e.mfdMgr,
	}

	for _, m := range managers {
		err := m.Initialize()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Exec) SetConfiguration(config *config.Configuration) {
	e.configuration = config
}

func (e *Exec) SendExecReadOnlyMessage(message string, routing *kexec.ConsoleIdentifier) {
	consMsg := kexec.ConsoleReadOnlyMessage{
		Source:         e.runControlEntry,
		Routing:        routing,
		Text:           message,
		DoNotEmitRunId: true,
	}
	e.consoleMgr.SendReadOnlyMessage(&consMsg)
}

func (e *Exec) SendExecReadReplyMessage(
	message string,
	maxReplyChars int,
	routing *kexec.ConsoleIdentifier,
) (string, error) {
	consMsg := kexec.ConsoleReadReplyMessage{
		Source:         e.runControlEntry,
		Routing:        routing,
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

func (e *Exec) SendExecRestrictedReadReplyMessage(
	message string,
	accepted []string,
	routing *kexec.ConsoleIdentifier,
) (string, error) {
	if len(accepted) == 0 {
		return "", fmt.Errorf("bad accepted list")
	}

	maxReplyLen := 0
	for _, acceptString := range accepted {
		if maxReplyLen < len(acceptString) {
			maxReplyLen = len(acceptString)
		}
	}

	consMsg := kexec.ConsoleReadReplyMessage{
		Source:         e.runControlEntry,
		Routing:        routing,
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

func (e *Exec) SetJumpKey(jkNumber int, value bool) {
	if (jkNumber >= 1) && (jkNumber <= 36) {
		e.jumpKeys[jkNumber-1] = value
	}
}

func (e *Exec) Stop(code kexec.StopCode) {
	// TODO need to set contingency in the Exec RCE
	e.stopFlag = true
	e.stopCode = code
	e.phase = kexec.ExecPhaseStopped
}

func (e *Exec) PerformDump(fullFlag bool) (string, error) {
	now := time.Now()
	fileName := fmt.Sprintf("kexec-%04v%02v%02v-%02v%02v%02v.dump",
		now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	dumpFile, err := os.Create(fileName)
	if err != nil {
		err := fmt.Errorf("cannot create dump file:%v\n", err)
		return "", err
	}

	_, _ = fmt.Fprintf(dumpFile, "Exec Dump ----------------------------------------------------\n")

	_, _ = fmt.Fprintf(dumpFile, "  Phase:         %v\n", e.phase)
	_, _ = fmt.Fprintf(dumpFile, "  Stopped:       %v\n", e.stopFlag)
	_, _ = fmt.Fprintf(dumpFile, "  StopCode:      %03o\n", e.stopCode)

	str := "Jump Keys Set:"
	for jk := 1; jk <= 36; jk++ {
		if e.jumpKeys[jk-1] {
			str += fmt.Sprintf(" %v", jk)
		}
	}
	_, _ = fmt.Fprintf(dumpFile, "  %v\n", str)

	// TODO something different when fullFlag is set

	e.consoleMgr.Dump(dumpFile, "")
	e.keyinMgr.Dump(dumpFile, "")
	e.nodeMgr.Dump(dumpFile, "")
	e.facMgr.Dump(dumpFile, "")
	e.mfdMgr.Dump(dumpFile, "")

	// TODO run control table, etc

	err = dumpFile.Close()
	if err != nil {
		err := fmt.Errorf("cannot close dump file:%v\n", err)
		return "", err
	}

	return fileName, nil
}
