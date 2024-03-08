// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type StopCode uint

const (
	StopFacilitiesComplex                                 = StopCode(001)
	StopUseStatementToExecPCTFailed                       = StopCode(031)
	StopFileAssignErrorOccurredDuringSystemInitialization = StopCode(034)
	StopInternalExecIOFailed                              = StopCode(040)
	StopFullCycleReachedForRunids                         = StopCode(044)
	StopExecRequestForMassStorageFailed                   = StopCode(052)
	StopErrorAccessingFacilitiesDataStructure             = StopCode(055)
	StopConsoleResponseRequiresReboot                     = StopCode(055)
	StopTrackToBeReleasedWasNotAllocated                  = StopCode(057)
	StopNoMainItemLink                                    = StopCode(057)
	StopInitializationSystemConfigurationError            = StopCode(064)
	StopInitializationSystemLibrariesCorruptOrMissing     = StopCode(044)
	StopClearTestSetAttemptedWhenNotSet                   = StopCode(066)
	StopResourceReleaseFailure                            = StopCode(067)
	StopActivityIdNoLongerExists                          = StopCode(073)
	StopExecContingencyHandler                            = StopCode(0103)
	StopExecActivityTakenToEMode                          = StopCode(0105)
	StopIOErrorBootTape                                   = StopCode(0145)
	StopOperatorInitiatedRecovery                         = StopCode(0150) // i.e., $!
	StopDirectoryErrors                                   = StopCode(0151)
	StopSectorToBeReleasedWasNotAllocated                 = StopCode(0157)
	StopIOPacketErrorForSystemIO                          = StopCode(0202)
	StopErrorInSystemIOTable                              = StopCode(0205)
	StopInvalidLDAT                                       = StopCode(0253)
)
