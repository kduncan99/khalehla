// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

type IoStatus uint

const (
	_ IoStatus = iota
	IosNotStarted
	IosComplete
	IosInProgress
	IosCanceled // channel or device was reset

	IosAtLoadPoint
	IosInvalidChannelProgram
	IosDeviceDoesNotExist
	IosDeviceIsDown
	IosDeviceIsNotAccessible
	IosDeviceIsNotReady
	IosEndOfFile
	IosEndOfTape
	IosInternalError // usually means the Exec fell over
	IosInvalidBlockId
	IosInvalidBufferSize
	IosInvalidFunction
	IosInvalidNodeType
	IosInvalidPacket
	IosInvalidPackName
	IosInvalidPrepFactor
	IosInvalidTapeBlock
	IosInvalidTrackCount
	IosLostPosition
	IosMediaAlreadyMounted
	IosMediaNotMounted
	IosNonIntegralRead
	IosPackNotPrepped
	IosReadNotAllowed
	IosReadOverrun
	IosSystemError
	IosWriteProtected
)

var IoStatusTable = map[IoStatus]string{
	IosNotStarted:            "NotStarted",
	IosComplete:              "Complete",
	IosInProgress:            "InProgress",
	IosCanceled:              "Canceled",
	IosAtLoadPoint:           "AtLoadPoint",
	IosDeviceDoesNotExist:    "DeviceDoesNotExist",
	IosDeviceIsDown:          "DeviceIsDown",
	IosDeviceIsNotAccessible: "DeviceNotAccessible",
	IosDeviceIsNotReady:      "DeviceNotReady",
	IosEndOfFile:             "EndOfFile",
	IosEndOfTape:             "EndOfTape",
	IosInternalError:         "InternalError",
	IosInvalidBlockId:        "InvalidBlockId",
	IosInvalidBufferSize:     "InvalidBufferSize",
	IosInvalidChannelProgram: "InvalidChannelProgram",
	IosInvalidFunction:       "InvalidFunction",
	IosInvalidNodeType:       "InvalidNodeType",
	IosInvalidPacket:         "InvalidPacket",
	IosInvalidPackName:       "InvalidPackName",
	IosInvalidPrepFactor:     "InvalidPrepFactor",
	IosInvalidTapeBlock:      "InvalidTapeBlock",
	IosInvalidTrackCount:     "InvalidTrackCount",
	IosLostPosition:          "IosLostPosition",
	IosMediaAlreadyMounted:   "MediaAlreadyMounted",
	IosMediaNotMounted:       "MediaNotMounted",
	IosNonIntegralRead:       "NonIntegralRead",
	IosPackNotPrepped:        "PackNotPrepped",
	IosReadNotAllowed:        "ReadNotAllowed",
	IosReadOverrun:           "ReadOverrun",
	IosSystemError:           "SystemError",
	IosWriteProtected:        "WriteProtected",
}
