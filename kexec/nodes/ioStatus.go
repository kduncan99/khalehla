// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodes

type IoStatus uint

const (
	_ IoStatus = iota
	IosNotStarted
	IosInProgress
	IosComplete
	IosSystemError
	IosInternalError // usually means the Exec fell over

	IosDeviceDoesNotExist
	IosDeviceIsDown
	IosDeviceIsNotAccessible
	IosInvalidFunction
	IosNilBuffer
	IosInvalidBufferSize
	IosInvalidBlockId
	IosInvalidNodeType
	IosInvalidPackName
	IosInvalidPrepFactor
	IosInvalidTrackCount
	IosMediaAlreadyMounted
	IosMediaNotMounted
	IosPackNotPrepped
	IosWriteProtected
)

var IoStatusTable = map[IoStatus]string{
	IosComplete:              "Complete",
	IosDeviceDoesNotExist:    "DeviceDoesNotExist",
	IosDeviceIsDown:          "DeviceIsDown",
	IosDeviceIsNotAccessible: "DeviceNotAccessible",
	IosInProgress:            "InProgress",
	IosInternalError:         "InternalError",
	IosInvalidBlockId:        "InvalidBlockId",
	IosInvalidBufferSize:     "InvalidBufferSize",
	IosInvalidFunction:       "InvalidFunction",
	IosInvalidNodeType:       "InvalidNodeType",
	IosInvalidPackName:       "InvalidPackName",
	IosInvalidPrepFactor:     "InvalidPrepFactor",
	IosInvalidTrackCount:     "InvalidTrackCount",
	IosMediaAlreadyMounted:   "MediaAlreadyMounted",
	IosMediaNotMounted:       "MediaNotMounted",
	IosNilBuffer:             "NilBuffer",
	IosNotStarted:            "NotStarted",
	IosPackNotPrepped:        "PackNotPrepped",
	IosSystemError:           "SystemError",
	IosWriteProtected:        "WriteProtected",
}
