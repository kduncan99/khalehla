// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

type IoStatus uint

const (
	_ IoStatus = iota
	IosNotStarted
	IosComplete
	IosInProgress

	IosDeviceDoesNotExist
	IosDeviceIsDown
	IosDeviceIsNotAccessible
	IosDeviceIsNotReady
	IosEndOfFile
	IosInternalError // usually means the Exec fell over
	IosInvalidBlockId
	IosInvalidBufferSize
	IosInvalidFunction
	IosInvalidNodeType
	IosInvalidPackName
	IosInvalidPrepFactor
	IosInvalidTapeBlock
	IosInvalidTrackCount
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
	IosDeviceDoesNotExist:    "DeviceDoesNotExist",
	IosDeviceIsDown:          "DeviceIsDown",
	IosDeviceIsNotAccessible: "DeviceNotAccessible",
	IosDeviceIsNotReady:      "DeviceNotReady",
	IosEndOfFile:             "EndOfFile",
	IosInternalError:         "InternalError",
	IosInvalidBlockId:        "InvalidBlockId",
	IosInvalidBufferSize:     "InvalidBufferSize",
	IosInvalidFunction:       "InvalidFunction",
	IosInvalidNodeType:       "InvalidNodeType",
	IosInvalidPackName:       "InvalidPackName",
	IosInvalidPrepFactor:     "InvalidPrepFactor",
	IosInvalidTapeBlock:      "InvalidTapeBlock",
	IosInvalidTrackCount:     "InvalidTrackCount",
	IosMediaAlreadyMounted:   "MediaAlreadyMounted",
	IosMediaNotMounted:       "MediaNotMounted",
	IosNilBuffer:             "NilBuffer",
	IosPackNotPrepped:        "PackNotPrepped",
	IosReadNotAllowed:        "ReadNotAllowed",
	IosSystemError:           "SystemError",
	IosWriteProtected:        "WriteProtected",
}
