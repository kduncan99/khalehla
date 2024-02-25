// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package exec

import (
	"khalehla/kexec"
)

// performInitialBoot is invoked on the first session initiated by the application.
// being here does NOT imply JK13, although that is a possibility
func (e *Exec) performInitialBoot() {
	if e.jumpKeys[kexec.JumpKey13Index] {
		msg := "Jump key 13 is set on initial boot - Continue? Y/N"
		accepted := []string{"Y", "N"}
		reply, _ := e.SendExecRestrictedReadReplyMessage(msg, accepted)
		if e.stopFlag {
			return
		} else if reply == "N" {
			e.Stop(kexec.StopConsoleResponseRequiresReboot)
		}
	}

	e.mfdMgr.InitializeMassStorage()
	if e.stopFlag {
		return
	}

	// TODO security officer setup

	if e.jumpKeys[kexec.JumpKey7Index] {
		// TODO
	}

	if e.jumpKeys[kexec.JumpKey9Index] {
		// TODO
	}
}

// performRecoveryBoot is invoked on subsequent sessions initiated by the application...
// i.e., recovery boots. We try to pick up where we left off after clearing up some tables and structures.
func (e *Exec) performRecoveryBoot() {
	if e.jumpKeys[kexec.JumpKey13Index] {
		msg := "Jump key 13 is set on recovery boot and will be ignored"
		e.SendExecReadOnlyMessage(msg, nil)
		if e.stopFlag {
			return
		}
	}

	e.mfdMgr.RecoverMassStorage()
	if e.stopFlag {
		return
	}

	// TODO
}
