// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package channels

import (
	"fmt"
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
	"khalehla/pkg"
)

type ChannelSender interface {
	ChannelProgramComplete(channelProgram *ChannelProgram)
}

type TransferFormat uint

const (
	NoTransferFormat TransferFormat = iota
	TransferPacked
	Transfer8Bit
	Transfer6Bit
)

type TransferDirection uint

const (
	NoTransferDirection TransferDirection = iota
	DirectionForward
	DirectionBackward
	DirectionStatic
	DirectionSkip
)

type ControlWord struct {
	Buffer    []pkg.Word36
	Offset    uint
	Length    uint
	Direction TransferDirection
	Format    TransferFormat
}

// ChannelProgram tells a channel how to perform an IO
type ChannelProgram struct {
	NodeIdentifier   hardware.NodeIdentifier
	IoFunction       ioPackets.IoFunction
	IoStatus         ioPackets.IoStatus
	BlockId          hardware.BlockId       // for disk
	ControlWords     []ControlWord          // for reads and writes
	MountInfo        *ioPackets.IoMountInfo // for mount
	PrepInfo         *ioPackets.IoPrepInfo  // for prep
	BytesTransferred uint
	WordsTransferred uint
	Listener         ChannelSender
}

func (cp *ChannelProgram) GetString() string {
	return fmt.Sprintf("node:%v func:%v blk:%v stat:%v(%v)",
		cp.NodeIdentifier,
		ioPackets.IoFunctionTable[cp.IoFunction],
		cp.BlockId,
		cp.IoStatus,
		ioPackets.IoStatusTable[cp.IoStatus])
}
