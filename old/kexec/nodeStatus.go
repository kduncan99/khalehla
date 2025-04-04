// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type FacNodeStatus uint

const (
	_ FacNodeStatus = iota
	FacNodeStatusUp
	FacNodeStatusReserved
	FacNodeStatusDown
	FacNodeStatusSuspended
)
