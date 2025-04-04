// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package channels

// TODO I'm not entirely sure we want this -
//	I think it is better to just have the channel program exist in some main storage bank,
//	accessible to the IOP and to the Channel.
/*
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
	Buffer    []uint64
	Offset    uint
	Length    uint
	Direction TransferDirection
	Format    TransferFormat
}

// ChannelProgram tells a channel how to perform an IO
type ChannelProgram struct {
	NodeIdentifier   hardware.NodeIdentifier
	IoFunction       ioPackets.IoFunction
	IoStatus         ioPackets2.IoStatus
	BlockId          hardware.BlockId        // for disk
	ControlWords     []ControlWord           // for reads and writes
	MountInfo        *ioPackets2.IoMountInfo // for mount
	PrepInfo         *ioPackets2.IoPrepInfo  // for prep
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
		ioPackets2.IoStatusTable[cp.IoStatus])
}
*/
