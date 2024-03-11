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

	IosAtLoadPoint
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
	IosInvalidPackName
	IosInvalidPrepFactor
	IosInvalidTapeBlock
	IosInvalidTrackCount
	IosLostPosition
	IosMediaAlreadyMounted
	IosMediaNotMounted
	IosNilBuffer
	IosPackNotPrepped
	IosReadNotAllowed
	IosSystemError
	IosWriteProtected
)

var IoStatusTable = map[IoStatus]string{
	IosNotStarted:            "NotStarted",
	IosComplete:              "Complete",
	IosInProgress:            "InProgress",
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
	IosInvalidFunction:       "InvalidFunction",
	IosInvalidNodeType:       "InvalidNodeType",
	IosInvalidPackName:       "InvalidPackName",
	IosInvalidPrepFactor:     "InvalidPrepFactor",
	IosInvalidTapeBlock:      "InvalidTapeBlock",
	IosInvalidTrackCount:     "InvalidTrackCount",
	IosLostPosition:          "IosLostPosition",
	IosMediaAlreadyMounted:   "MediaAlreadyMounted",
	IosMediaNotMounted:       "MediaNotMounted",
	IosNilBuffer:             "NilBuffer",
	IosPackNotPrepped:        "PackNotPrepped",
	IosReadNotAllowed:        "ReadNotAllowed",
	IosSystemError:           "SystemError",
	IosWriteProtected:        "WriteProtected",
}
