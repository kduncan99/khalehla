// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

const InvalidLink MFDRelativeAddress = 0_400000_000000
const InvalidLDAT LDATIndex = 0_400000

const (
	_ Granularity = iota
	TrackGranularity
	PositionGranularity
)

const (
	JumpKey1Index  = 0
	JumpKey2Index  = 1
	JumpKey3Index  = 2
	JumpKey4Index  = 3
	JumpKey6Index  = 5
	JumpKey7Index  = 6
	JumpKey9Index  = 8
	JumpKey13Index = 12
)

const (
	AOption = 0_000200_000000
	BOption = 0_000100_000000
	COption = 0_000040_000000
	DOption = 0_000020_000000
	EOption = 0_000010_000000
	FOption = 0_000004_000000
	GOption = 0_000002_000000
	HOption = 0_000001_000000
	IOption = 0_000000_400000
	JOption = 0_000000_200000
	KOption = 0_000000_100000
	LOption = 0_000000_040000
	MOption = 0_000000_020000
	NOption = 0_000000_010000
	OOption = 0_000000_004000
	POption = 0_000000_002000
	QOption = 0_000000_001000
	ROption = 0_000000_000400
	SOption = 0_000000_000200
	TOption = 0_000000_000100
	UOption = 0_000000_000040
	VOption = 0_000000_000020
	WOption = 0_000000_000010
	XOption = 0_000000_000004
	YOption = 0_000000_000002
	ZOption = 0_000000_000001
)

const (
	ExecPhaseNotStarted ExecPhase = iota
	ExecPhaseInitializing
	ExecPhaseRunning
	ExecPhaseStopped
)
