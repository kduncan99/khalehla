// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tapeUtil

import (
	"fmt"
	"khalehla/kexec/nodeMgr"
	"khalehla/pkg"
)

// FSTapeDevice is a wrapper around the standard file-system tape reader/writer
type FSTapeDevice struct {
	device *nodeMgr.FileSystemTapeDevice
}

func NewFSTapeDevice() *FSTapeDevice {
	return &FSTapeDevice{
		device: nodeMgr.NewFileSystemTapeDevice(),
	}
}

func (dev *FSTapeDevice) Close() error {
	if !dev.device.IsMounted() {
		return fmt.Errorf("device is not mounted")
	} else if !dev.device.IsReady() {
		return fmt.Errorf("device is not ready")
	}

	ioPkt := nodeMgr.NewTapeIoPacketUnmount(0)
	dev.device.StartIo(ioPkt)
	if ioPkt.GetIoStatus() != nodeMgr.IosComplete {
		return fmt.Errorf(nodeMgr.IoStatusTable[ioPkt.GetIoStatus()])
	}

	return nil
}

func (dev *FSTapeDevice) OpenInputFile(fileName string) error {
	if dev.device.IsMounted() {
		return fmt.Errorf("device is already mounted")
	}

	ioPkt := nodeMgr.NewTapeIoPacketMount(0, fileName, false)
	dev.device.StartIo(ioPkt)
	if ioPkt.GetIoStatus() != nodeMgr.IosComplete {
		return fmt.Errorf(nodeMgr.IoStatusTable[ioPkt.GetIoStatus()])
	}
	dev.device.SetIsReady(true)

	return nil
}

func (dev *FSTapeDevice) OpenOutputFile(fileName string) error {
	if dev.device.IsMounted() {
		return fmt.Errorf("device is already mounted")
	}

	ioPkt := nodeMgr.NewTapeIoPacketMount(0, fileName, false)
	dev.device.StartIo(ioPkt)
	if ioPkt.GetIoStatus() != nodeMgr.IosComplete {
		return fmt.Errorf(nodeMgr.IoStatusTable[ioPkt.GetIoStatus()])
	}
	dev.device.SetIsReady(true)

	return nil
}

func (dev *FSTapeDevice) ReadVolumeHeader() (volumeHeader *VolumeHeader, err error) {
	if !dev.device.IsMounted() {
		return nil, fmt.Errorf("device is not mounted")
	} else if !dev.device.IsReady() {
		return nil, fmt.Errorf("device is not ready")
	}

	// This format does not store individual volume headers,
	// therefore we cannot read a header, and must produce a default VolumeHeader struct.
	return NewVolumeHeader(), nil
}

func (dev *FSTapeDevice) ReadFileHeader() (fileHeader *FileHeader, err error) {
	if !dev.device.IsMounted() {
		return nil, fmt.Errorf("device is not mounted")
	} else if !dev.device.IsReady() {
		return nil, fmt.Errorf("device is not ready")
	}

	// This format does not store individual file headers,
	// therefore we cannot read a header, and must produce a default FileHeader struct.
	fileHeader = NewFileHeader()
	return
}

func (dev *FSTapeDevice) ReadBlock() (data []pkg.Word36, eof bool, err error) {
	data = nil
	eof = false
	err = nil

	if !dev.device.IsMounted() {
		err = fmt.Errorf("device is not mounted")
		return
	} else if !dev.device.IsReady() {
		err = fmt.Errorf("device is not ready")
		return
	}

	ioPkt := nodeMgr.NewTapeIoPacketRead(0)
	dev.device.StartIo(ioPkt)
	if ioPkt.GetIoStatus() == nodeMgr.IosEndOfFile {
		return nil, true, nil
	} else if ioPkt.GetIoStatus() != nodeMgr.IosComplete {
		err = fmt.Errorf(nodeMgr.IoStatusTable[ioPkt.GetIoStatus()])
		return
	}

	data = ioPkt.GetBuffer()
	return
}

func (dev *FSTapeDevice) WriteVolumeHeader() error {
	if !dev.device.IsMounted() {
		return fmt.Errorf("device is not mounted")
	} else if !dev.device.IsReady() {
		return fmt.Errorf("device is not ready")
	}

	return nil
}

func (dev *FSTapeDevice) WriteFileHeader() error {
	if !dev.device.IsMounted() {
		return fmt.Errorf("device is not mounted")
	} else if !dev.device.IsReady() {
		return fmt.Errorf("device is not ready")
	}

	return nil
}

func (dev *FSTapeDevice) WriteBlock(buffer []pkg.Word36) error {
	if !dev.device.IsMounted() {
		return fmt.Errorf("device is not mounted")
	} else if !dev.device.IsReady() {
		return fmt.Errorf("device is not ready")
	}

	ioPkt := nodeMgr.NewTapeIoPacketWrite(0, buffer)
	dev.device.StartIo(ioPkt)
	if ioPkt.GetIoStatus() != nodeMgr.IosComplete {
		return fmt.Errorf(nodeMgr.IoStatusTable[ioPkt.GetIoStatus()])
	}

	return nil
}

func (dev *FSTapeDevice) WriteEndOfFile() error {
	if !dev.device.IsMounted() {
		return fmt.Errorf("device is not mounted")
	} else if !dev.device.IsReady() {
		return fmt.Errorf("device is not ready")
	}

	ioPkt := nodeMgr.NewTapeIoPacketWriteTapeMark(0)
	dev.device.StartIo(ioPkt)
	if ioPkt.GetIoStatus() != nodeMgr.IosComplete {
		return fmt.Errorf(nodeMgr.IoStatusTable[ioPkt.GetIoStatus()])
	}

	return nil
}
