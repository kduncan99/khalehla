// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec/facilitiesMgr"
	"khalehla/logger"
)

func handleStart(pkt *handlerPacket) (facResult *facilitiesMgr.FacStatusResult, resultCode uint64) {
	logger.LogTraceF("CSI", "handleStart:%v", *pkt.pcs)

	/*
		@START runstream-name
		@START[,scheduling-priority/options/processor-dispatching-priority]
			runstream-name[,setc,run-id,acct-id/user-id,project-id,runtime/deadline,pages/cards,start-time]
		scheduling-priority is one character from 'A' to 'Z'
		options include
			B, N, P, R, T, U, V, W, X, Y, Z
		processor-dispatching-priority is one character from 'A' to 'Z'

		resultCode values:
		000 Request processed normally 1 Improper runstream in file
		002 File access denied
		003 Element unobtainable
		004 No file specified
		005 I/O error encountered
		006 File not cataloged on mass storage 7 File freed while processing @START
		010 X option security violation
		011 U option requires SENTRY configured.
		012 User-id and account needed on @RUN statement
		013 Started denied access to target user-id
		014 Target user does not own @START file.
		015 Copy I/O error (plus err code appended to end)
		016 Copy file unobtainable (plus fac status appended to end)
		017 PCT full. Cannot copy started runstream.
		020 Z option requires starter to have Z-option privilege.
		021 Z option requires user-id to have system-high capability.
		022 Contact site administrator and report @START command internal error 01.
		023 STD mass storage tight. Cannot start run.
		024 The user-id is invalid for an Exec-initiated @START.
		025 The user-id does not exist.
	*/

	// TODO implement

	return nil, 0
}
