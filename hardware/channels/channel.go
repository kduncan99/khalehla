// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package channels

import (
	"io"
	"khalehla/hardware"
	"khalehla/hardware/devices"
	"khalehla/pkg"
)

// Channel manages async communication with the various deviceInfos assigned to it.
// It may also manage caching, automatic mounting, or any other various activities
// on behalf of the exec.
type Channel interface {
	AssignDevice(nodeIdentifier hardware.NodeIdentifier, device devices.Device) error
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware.NodeDeviceType
	GetNodeModelType() hardware.NodeModelType
	StartIo(program *ChannelProgram)
}

// transferFromBytes translates a byte buffer to a word buffer
// while observing transfer direction and format.
// returns word count and an indication whether we detected a non-integral transfer.
func transferFromBytes(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []pkg.Word36,
	destinationOffset uint,
	direction TransferDirection,
	format TransferFormat,
) (nonIntegral bool, wordCount uint) {
	switch format {
	case Transfer6Bit:
		return transferFromBytes6Bit(source, sourceOffset, sourceLength, destination, destinationOffset, direction)
	case Transfer8Bit:
		return transferFromBytes8Bit(source, sourceOffset, sourceLength, destination, destinationOffset, direction)
	case TransferPacked:
		return transferFromBytesPacked(source, sourceOffset, sourceLength, destination, destinationOffset, direction)
	}

	return
}

func transferFromBytes6Bit(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []pkg.Word36,
	destinationOffset uint,
	direction TransferDirection,
) (nonIntegral bool, wordCount uint) {
	switch direction {
	case DirectionForward:
		return pkg.ByteArray6BitToWord36(source, sourceOffset, sourceLength, destination, destinationOffset)
	case DirectionBackward:
		return pkg.ByteArray6BitToWord36Reversed(source, sourceOffset, sourceLength, destination, destinationOffset)
	default:
		return
	}
}

func transferFromBytes8Bit(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []pkg.Word36,
	destinationOffset uint,
	direction TransferDirection,
) (nonIntegral bool, wordCount uint) {
	switch direction {
	case DirectionForward:
		return pkg.ByteArray8BitToWord36(source, sourceOffset, sourceLength, destination, destinationOffset)
	case DirectionBackward:
		return pkg.ByteArray8BitToWord36Reversed(source, sourceOffset, sourceLength, destination, destinationOffset)
	default:
		return
	}
}

func transferFromBytesPacked(
	source []byte,
	sourceOffset uint,
	sourceLength uint,
	destination []pkg.Word36,
	destinationOffset uint,
	direction TransferDirection,
) (nonIntegral bool, wordCount uint) {
	switch direction {
	case DirectionForward:
		return pkg.ByteArrayPackedToWord36(source, sourceOffset, sourceLength, destination, destinationOffset)
	case DirectionBackward:
		return pkg.ByteArrayPackedToWord36Reversed(source, sourceOffset, sourceLength, destination, destinationOffset)
	default:
		return
	}
}

func transferFromWords(
	source []pkg.Word36,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
	direction TransferDirection,
	format TransferFormat,
) {
	switch format {
	case Transfer6Bit:
		transferFromWords6Bit(source, sourceOffset, sourceLength, destination, destinationOffset, direction)
	case Transfer8Bit:
		transferFromWords8Bit(source, sourceOffset, sourceLength, destination, destinationOffset, direction)
	case TransferPacked:
		transferFromWordsPacked(source, sourceOffset, sourceLength, destination, destinationOffset, direction)
	}
}

func transferFromWords6Bit(
	source []pkg.Word36,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
	direction TransferDirection,
) {
	switch direction {
	case DirectionForward:
		pkg.Word36ToByteArray6Bit(source, sourceOffset, sourceLength, destination, destinationOffset)
		break

	case DirectionBackward:
		pkg.Word36ToByteArray6BitReversed(source, sourceOffset, sourceLength, destination, destinationOffset)
		break

	case DirectionStatic:
		sx := sourceOffset
		dx := destinationOffset
		for cx := uint(0); cx < sourceLength; cx++ {
			destination[dx] = byte(source[sx].GetS1())
			destination[dx+1] = byte(source[sx].GetS2())
			destination[dx+2] = byte(source[sx].GetS3())
			destination[dx+3] = byte(source[sx].GetS4())
			destination[dx+4] = byte(source[sx].GetS5())
			destination[dx+5] = byte(source[sx].GetS6())
			dx += 6
		}

	case DirectionSkip:
		destBytes := 6 * sourceLength
		dy := destinationOffset
		for dx := uint(0); dx < destBytes; dx++ {
			destination[dy] = 0
			dy++
		}
	}
}

func transferFromWords8Bit(
	source []pkg.Word36,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
	direction TransferDirection,
) {
	switch direction {
	case DirectionForward:
		pkg.Word36ToByteArray8Bit(source, sourceOffset, sourceLength, destination, destinationOffset)
		break

	case DirectionBackward:
		pkg.Word36ToByteArray8BitReversed(source, sourceOffset, sourceLength, destination, destinationOffset)
		break

	case DirectionStatic:
		sx := sourceOffset
		dx := destinationOffset
		for cx := uint(0); cx < sourceLength; cx++ {
			destination[dx] = byte(source[sx].GetQ1())
			destination[dx+1] = byte(source[sx].GetQ2())
			destination[dx+2] = byte(source[sx].GetQ3())
			destination[dx+3] = byte(source[sx].GetQ4())
			dx += 4
		}

	case DirectionSkip:
		destBytes := 4 * sourceLength
		dy := destinationOffset
		for dx := uint(0); dx < destBytes; dx++ {
			destination[dy] = 0
			dy++
		}
	}
}

func transferFromWordsPacked(
	source []pkg.Word36,
	sourceOffset uint,
	sourceLength uint,
	destination []byte,
	destinationOffset uint,
	direction TransferDirection,
) {
	switch direction {
	case DirectionForward:
		pkg.Word36ToByteArrayPacked(source, sourceOffset, sourceLength, destination, destinationOffset)
		break

	case DirectionBackward:
		pkg.Word36ToByteArrayPackedReversed(source, sourceOffset, sourceLength, destination, destinationOffset)
		break

	case DirectionStatic:
		temp := []pkg.Word36{source[sourceOffset], source[sourceOffset]}
		for dx := uint(0); dx < sourceLength*9/2; dx += 9 {
			pkg.PackWord36Strict(temp, destination[dx:dx+9])
		}

	case DirectionSkip:
		destBytes := sourceLength * 9 / 2
		dy := destinationOffset
		for dx := uint(0); dx < destBytes; dx++ {
			destination[dy] = 0
			dy++
		}
	}
}
