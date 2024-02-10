// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package deviceMgr

import "khalehla/pkg"

// basic types -----------------------------------------------------------------------------

// BlockCount represents a number of pseudo-physical blocks.
// For disk nodes, a block contains a fixed number of words which corresponds to the relevant medium's prep factor.
// For tape nodes, a block contains a variable number of words.
type BlockCount uint64

// BlockId uniquely identifies a particular pseudo-physical block on a disk medium
type BlockId uint64

// NodeIdentifier uniquely identifies a particular device or channel (or anything else identifiable which we manage)
// It is currently implemented as the 1-6 character device name, all caps alphas and/or digits LJSF
// stored as fieldata in a Word36 struct
type NodeIdentifier pkg.Word36

// PrepFactor indicates the number of words stored in a block of data for disk media.
// Current valid values include 28, 56, 112, 224, 448, 896, and 1792.
type PrepFactor uint

// TrackCount represents a number of software tracks, each of which contain 1792 words of storage
type TrackCount uint

// pseudo enumerations ---------------------------------------------------------------------

type IoFunction uint

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

type IoStatus uint

const (
	_ IoStatus = iota
	IosNotStarted
	IosInProgress
	IosComplete
	IosSystemError

	IosDeviceNotAttached
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

type NodeStatus uint

const (
	_ NodeStatus = iota
	NodeStatusUp
	NodeStatusReserved
	NodeStatusDown
	NodeStatusSuspended
)

type NodeType uint

const (
	_ NodeType = iota
	NodeTypeDisk
	NodeTypeTape
)

// interfaces ------------------------------------------------------------------------------

// Channel manages async communication with the various deviceInfos assigned to it.
// It may also manage caching, automatic mounting, or any other various activities
// on behalf of the exec.
type Channel interface {
	assignDevice(deviceIdentifier NodeIdentifier, device Device) error
	getNodeType() NodeType
	startIo(ioPacket IoPacket)
}

// Device manages real or pseudo IO operations for a particular virtual device.
// It may do so synchronously or asynchronously
type Device interface {
	getNodeType() NodeType
	startIo(ioPacket IoPacket)
}

// IoPacket contains all the information necessary for a Channel to route an IO operation,
// and for a device to perform that IO operation.
type IoPacket interface {
	GetDeviceIdentifier() NodeIdentifier
	GetNodeType() NodeType
	GetIoFunction() IoFunction
	GetIoStatus() IoStatus
	SetIoStatus(ioStatus IoStatus)
}

// NodeInfo contains all the exec-managed information regarding a particular node
type NodeInfo interface {
	CreateNode()
	GetNodeIdentifier() NodeIdentifier
	GetNodeName() string
	GetNodeStatus() NodeStatus
	GetNodeType() NodeType
}

// ChannelInfo is intended primarily as a means of documenting the use of a more generic NodeInfo
type ChannelInfo interface {
	CreateNode()
	GetChannel() Channel
	GetNodeIdentifier() NodeIdentifier
	GetNodeName() string
	GetNodeStatus() NodeStatus
	GetNodeType() NodeType
}

// DeviceInfo is intended primarily as a means of documenting the use of a more generic NodeInfo
type DeviceInfo interface {
	CreateNode()
	GetDevice() Device
	GetNodeIdentifier() NodeIdentifier
	GetNodeName() string
	GetNodeStatus() NodeStatus
	GetNodeType() NodeType
	IsAccessible() bool
	IsMounted() bool
	SetIsAccessible(bool)
}

// simple structs --------------------------------------------------------------------------

// DiskPackGeometry describes various useful attributes of a particular prepped disk pack
type DiskPackGeometry struct {
	PrepFactor      PrepFactor // number of words contained in a physical record, packed
	TrackCount      TrackCount // number of tracks on the pack
	BlockCount      BlockCount // number of blocks on the pack (may not be track-aligned)
	SectorsPerBlock uint       // number of software sectors (28 words per) in a physical block
	BlocksPerTrack  uint       // physical blocks required to contain one software track (1792 words)
	BytesPerBlock   uint       // bytes required for a block containing packed word36 structs, rounded to power of 2
}
