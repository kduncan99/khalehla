package storage

import "kalehla/types"

// Stuff which potentially relates to any kind of device

const (
	DeviceTypeRawBlock       types.DeviceType = 010
	DeviceTypeFileBlock      types.DeviceType = 011
	DeviceTypePackedBlock    types.DeviceType = 012
	DeviceTypeTemporaryBlock types.DeviceType = 014
)

const (
	DeviceStatusSuccessful types.DeviceStatus = iota
	DeviceStatusNotOpen
	DeviceStatusAlreadyOpen
	DeviceStatusSystemError
	DeviceStatusCannotSetWriteProtect
	DeviceStatusWriteProtected
	DeviceStatusInvalidBlockSize
	DeviceStatusInvalidBlockId
	DeviceStatusInvalidBufferSize
	DeviceStatusMaxBlocksExceeded
	DeviceStatusInvalidLabel
	DeviceStatusInvalidIdentifierConstant
)

type DeviceResult struct {
	status      types.DeviceStatus
	systemError error
}

type Device interface {
	Close() DeviceResult
	GetDeviceType() types.DeviceType
	IsOpen() bool
	IsWriteProtected() bool
	Open(writeProtected bool, writeThrough bool) DeviceResult
}
