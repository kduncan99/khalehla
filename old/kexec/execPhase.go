// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type ExecPhase uint

const (
	ExecPhaseNotStarted ExecPhase = iota
	ExecPhaseInitializing
	ExecPhaseInitStopped // initialization error
	ExecPhaseRunning
	ExecPhaseStopped
)
