// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package exec

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/csi"
	"khalehla/klog"
)

// assignSystemLibraryFiles assigns the system library files to the Exec, although NOT exclusively.
func (e *Exec) assignSystemLibraryFiles() {
	// TODO
}

func (e *Exec) callCSI(statement string) uint64 {
	pcs, _, stat := csi.ParseControlStatement(e.runControlEntry, statement)
	if stat&0_400000_000000 == 0 {
		_, stat = csi.HandleControlStatement(e, e.runControlEntry, csi.CSTSourceERCSF, pcs)
	}
	if stat&0_400000_000000 != 0 {
		klog.LogErrorF("Exec", "Status %012o : %v", stat, statement)
	}
	return stat
}

func (e *Exec) copyTapeFile(tapeFileName string, systemFileName string) {
	// TODO
}

// loadSystemLibrary loads the three system library files from a library tape.
// The library tape (in lieu of boot tape) contains the following files in copy,g format
//
//	SYS$*LIB$
//	SYS$*RUN$
//	SYS$*RLIB$
func (e *Exec) loadSystemLibrary() {
	klog.LogTrace("Exec", "loadSystemLibrary")
	reelNumber := ""
	for !kexec.IsValidReelNumber(reelNumber) {
		msg := "Enter reel number of system library tape"
		reelNumber, _ := e.SendExecReadReplyMessage(msg, 6, nil)
		if e.stopFlag {
			return
		}
		if kexec.IsValidReelNumber(reelNumber) {
			break
		}
		msg = "Invalid reel number entered"
		e.SendExecReadOnlyMessage(msg, nil)
	}

	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	}

	stat := e.callCSI(fmt.Sprintf("@ASG,TJ LIB.,%v,%v", e.configuration.TapeDefaultMnemonic, reelNumber))
	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	} else if stat&0_400000_000000 != 0 {
		klog.LogFatalF("Exec", "Cannot assign tape unit for library tape:%012o", stat)
		e.Stop(kexec.StopFileAssignErrorOccurredDuringSystemInitialization)
		return
	}

	stat = e.callCSI(fmt.Sprintf("@CAT,G SYS$*LIB$(+1),%v/%v/TRK/%v",
		e.configuration.LibAssignMnemonic,
		e.configuration.LibInitialReserve,
		e.configuration.LibMaximumSize))
	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	} else if stat&0_400000_000000 != 0 {
		klog.LogFatalF("Exec", "Cannot catalog SYS$*LIB$ file:%012o", stat)
		e.Stop(kexec.StopFileAssignErrorOccurredDuringSystemInitialization)
		return
	}

	e.copyTapeFile("LIB", "SYS$*LIB$")
	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	}

	stat = e.callCSI(fmt.Sprintf("@CAT,G SYS$*RUN$(+1),%v/%v/TRK/%v",
		e.configuration.RunAssignMnemonic,
		e.configuration.RunInitialReserve,
		e.configuration.RunMaximumSize))
	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	} else if stat&0_400000_000000 != 0 {
		klog.LogFatalF("Exec", "Cannot catalog SYS$*RUN$ file:%012o", stat)
		e.Stop(kexec.StopFileAssignErrorOccurredDuringSystemInitialization)
		return
	}

	stat = e.callCSI("@ASG,AX SYS$*RUN$")
	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	} else if stat&0_400000_000000 != 0 {
		klog.LogFatalF("Exec", "Cannot assign SYS$*RUN$ file:%012o", stat)
		e.Stop(kexec.StopFileAssignErrorOccurredDuringSystemInitialization)
		return
	}

	e.copyTapeFile("LIB", "SYS$*RUN$")
	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	}

	stat = e.callCSI(fmt.Sprintf("@CAT,G SYS$*RLIB$(+1),%v/128/TRK/9999",
		e.configuration.MassStorageDefaultMnemonic))
	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	} else if stat&0_400000_000000 != 0 {
		klog.LogFatalF("Exec", "Cannot catalog SYS$*RLIB$ file:%012o", stat)
		e.Stop(kexec.StopFileAssignErrorOccurredDuringSystemInitialization)
		return
	}

	stat = e.callCSI("@ASG,AX SYS$*RLIB$")
	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	} else if stat&0_400000_000000 != 0 {
		klog.LogFatalF("Exec", "Cannot assign SYS$*RLIB$ file:%012o", stat)
		e.Stop(kexec.StopFileAssignErrorOccurredDuringSystemInitialization)
		return
	}

	e.copyTapeFile("LIB", "SYS$*RLIB$")
	if e.stopFlag {
		klog.LogTrace("Exec", "loadSystemLibrary early exit")
		return
	}

	e.callCSI("@FREE LIB")
	e.callCSI("@FREE,X SYS$*LIB$")
	e.callCSI("@FREE,X SYS$*RUN$")
	e.callCSI("@FREE,X SYS$*RLIB$")
	klog.LogTrace("Exec", "loadSystemLibrary normal exit")
}

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

		e.mfdMgr.InitializeMassStorage()
		if e.stopFlag {
			return
		}

		// Catalog SYS$*DLOC$
		// TODO uncomment
		//stat := e.callCSI(fmt.Sprintf("@CAT,G SYS$*DLOC$(+1)/RDKDLC/WRKDLC,%v",
		//	e.configuration.DLOCAssignMnemonic))
		//if e.stopFlag {
		//	return
		//} else if stat&0_400000_000000 != 0 {
		//	klog.LogFatalF("Exec", "Cannot catalog SYS$*DLOC$ file:%012o", stat)
		//	e.Stop(kexec.StopFileAssignErrorOccurredDuringSystemInitialization)
		//	return
		//}

	} else {
		e.mfdMgr.RecoverMassStorage()
		if e.stopFlag {
			return
		}
	}

	stat := e.callCSI("@ASG,AX SYS$*MFD$$")
	if stat&0_400000_000000 != 0 {
		e.Stop(kexec.StopFileAssignErrorOccurredDuringSystemInitialization)
		return
	}

	// TODO security officer setup, user and account file creation or assignment

	if e.jumpKeys[kexec.JumpKey4Index] || e.jumpKeys[kexec.JumpKey13Index] {
		e.loadSystemLibrary()
		if e.stopFlag {
			return
		}
	}

	e.assignSystemLibraryFiles()
	if e.stopFlag {
		return
	}

	if e.jumpKeys[kexec.JumpKey7Index] {
		// TODO ask whether to initialize TIP, etc
	}

	if e.jumpKeys[kexec.JumpKey9Index] {
		// TODO ask whether to initialize or recover GENF$
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
