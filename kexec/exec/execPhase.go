// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package exec

import "khalehla/kexec"

const (
	ExecPhaseNotStarted kexec.ExecPhase = iota
	ExecPhaseInitializing
	ExecPhaseInitStopped // initialization error
	ExecPhaseRunning
	ExecPhaseStopped
)
