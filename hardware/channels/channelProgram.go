// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package channels

import (
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
	"khalehla/pkg"
)

type ChannelSender interface {
	ChannelProgramComplete(channelProgram *ChannelProgram)
}

type TransferFormat uint

const (
	TransferPacked TransferFormat = iota
	Transfer8Bit
	Transfer6Bit
)

type TransferDirection uint

const (
	DirectionForward TransferDirection = iota
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
	ControlWords     []ControlWord // for reads and writes
	Filename         string        // for mount
	WriteProtected   bool          // for mount
	BytesTransferred uint
	WordsTransferred uint
	Listener         ChannelSender
}
