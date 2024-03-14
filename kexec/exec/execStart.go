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
		reply, _ := e.SendExecRestrictedReadReplyMessage(msg, accepted, nil)
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

	// TODO ASG,A SYS$*MFD$$

	// TODO security officer setup

	if e.jumpKeys[kexec.JumpKey4Index] {
		// TODO for JK4, load libraries from library tape
		// libray tape (in lieu of boot tape) contains the following files in copy,g format...
		//   SYS$*LIB$
		//   SYS$*RUN$
		//   SYS$*RLIB$
	}

	if e.jumpKeys[kexec.JumpKey7Index] {
		// TODO ask whether to innitialize TIP, etc
	}

	if e.jumpKeys[kexec.JumpKey9Index] {
		// TODO ask whether to initialiize or recover GENF$
		// (I think this is implied with JK13)
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

	// TODO recovery boot
}
