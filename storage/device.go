package storage

import "khalehla/pkg"

// Stuff which potentially relates to any kind of device

const (
	DeviceTypeRawBlock       pkg.DeviceType = 010
	DeviceTypeFileBlock      pkg.DeviceType = 011
	DeviceTypePackedBlock    pkg.DeviceType = 012
	DeviceTypeTemporaryBlock pkg.DeviceType = 014
)

const (
	DeviceStatusSuccessful pkg.DeviceStatus = iota
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
	status      pkg.DeviceStatus
	systemError error
}

type Device interface {
	Close() DeviceResult
	GetDeviceType() pkg.DeviceType
	IsOpen() bool
	IsWriteProtected() bool
	Open(writeProtected bool, writeThrough bool) DeviceResult
}
