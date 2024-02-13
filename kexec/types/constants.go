// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

const (
	_ IoFunction = iota
	IofMount
	IofPrep
	IofReset
	IofRead
	IofReadLabel
	IofUnmount
	IofWrite
)

const (
	_ IoStatus = iota
	IosNotStarted
	IosInProgress
	IosComplete
	IosSystemError
	IosInternalError // usually means the Exec fell over

	IosDeviceDoesNotExist
	IosDeviceIsNotAccessible
	IosInvalidFunction
	IosNilBuffer
	IosInvalidBufferSize
	IosInvalidBlockId
	IosInvalidNodeType
	IosInvalidPackName
	IosInvalidPrepFactor
	IosInvalidTrackCount
	IosMediaNotMounted
	IosMediaAlreadyMounted
	IosPackNotPrepped
	IosWriteProtected
)

const (
	_ NodeCategory = iota
	NodeCategoryChannel
	NodeCategoryDevice
)

const (
	_ NodeStatus = iota
	NodeStatusUp
	NodeStatusReserved
	NodeStatusDown
	NodeStatusSuspended
)

const (
	_ NodeType = iota
	NodeTypeDisk
	NodeTypeTape
)

const (
	ExecPhaseNotStarted ExecPhase = iota
	ExecPhaseInitializing
	ExecPhaseRunning
	ExecPhaseStopped
)

const (
	StopFacilitiesComplex                                 = 001
	StopUseStatementToExecPCTFailed                       = 031
	StopFileAssignErrorOccurredDuringSystemInitialization = 034
	StopInternalExecIOFailed                              = 040
	StopFullCycleReachedForRunids                         = 044
	StopExecRequestForMassStorageFailed                   = 052
	StopErrorAccessingFacilitiesDataStructure             = 055
	StopConsoleResponseRequiresReboot                     = 055
	StopTrackToBeReleasedWasNotAllocated                  = 057
	StopNoMainItemLink                                    = 057
	StopInitializationSystemConfigurationError            = 064
	StopInitializationSystemLibrariesCorruptOrMissing     = 044
	StopClearTestSetAttemptedWhenNotSet                   = 066
	StopResourceReleaseFailure                            = 067
	StopActivityIdNoLongerExists                          = 073
	StopExecContingencyHandler                            = 0103
	StopExecActivityTakenToEMode                          = 0105
	StopOperatorInitiatedRecovery                         = 0150 // i.e., $!
	StopDirectoryErrors                                   = 0151
	StopSectorToBeReleasedWasNotAllocated                 = 0157
	StopIOPacketErrorForSystemIO                          = 0202
	StopErrorInSystemIOTable                              = 0205
	StopInvalidLDAT                                       = 0253
)
