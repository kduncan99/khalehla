// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import (
	"io"
	"khalehla/hardware"
	"khalehla/hardware/ioPackets"
	"os"
)

// Device manages real or pseudo IO operations for a particular virtual device.
// It may do so synchronously or asynchronously
type Device interface {
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() hardware.NodeCategoryType
	GetNodeDeviceType() hardware.NodeDeviceType
	GetNodeModelType() hardware.NodeModelType
	IsMounted() bool
	IsReady() bool
	Reset()
	SetVerbose(flag bool)
	StartIo(ioPacket ioPackets.IoPacket)
}

type FileSystemDevice interface {
	GetFile() *os.File
}

// readExact reads exactly the requested number of bytes from the device file.
// This is generally for file-system based devices. Presumes reading is actually allowed.
// Uses the device's buffer, but does not update the device's current offset.
func readExact(
	device FileSystemDevice,
	buffer []byte,
	length uint32,
	offset int64,
) error {
	// do the read - loop to make sure we read all we're asked to read
	index := 0
	remaining := length
	for remaining > 0 {
		bytesRead, err := device.GetFile().ReadAt(buffer[index:length], offset)
		if err != nil {
			return err
		}

		index += bytesRead
		remaining -= uint32(bytesRead)
		offset += int64(bytesRead)
	}

	return nil
}

// writeExact writes exactly the requested number of bytes to the device file.
// This is generally for file-system based devices. Presumes writing is actually allowed.
// Uses the provided buffer (which might not be the device buffer),
// but does not update the device's current offset.
func writeExact(
	device FileSystemDevice,
	buffer []byte,
	length uint32,
	offset int64,
) error {
	index := 0
	remaining := length
	for remaining > 0 {
		bytesWritten, err := device.GetFile().WriteAt(buffer[index:length], offset)
		if err != nil {
			return err
		}

		index += bytesWritten
		remaining -= uint32(bytesWritten)
		offset += int64(bytesWritten)
	}

	return nil
}
