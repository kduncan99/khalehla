// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"khalehla/kexec"
	"os"
	"sync"
)

// This is a very simple pseudo tape device

type FileSystemTapeDevice struct {
	fileName         *string
	file             *os.File
	isReady          bool
	isWriteProtected bool
	mutex            sync.Mutex
}

func NewFileSystemTapeDevice() *FileSystemTapeDevice {
	return &FileSystemTapeDevice{
		isWriteProtected: true,
	}
}

func (tape *FileSystemTapeDevice) GetNodeCategoryType() kexec.NodeCategoryType {
	return kexec.NodeCategoryDevice
}

func (tape *FileSystemTapeDevice) GetNodeDeviceType() kexec.NodeDeviceType {
	return kexec.NodeDeviceTape
}

func (tape *FileSystemTapeDevice) GetNodeModelType() NodeModelType {
	return NodeModelFileSystemTapeDevice
}

func (tape *FileSystemTapeDevice) IsMounted() bool {
	return tape.file != nil
}

func (tape *FileSystemTapeDevice) IsReady() bool {
	return tape.isReady
}

func (tape *FileSystemTapeDevice) IsWriteProtected() bool {
	return tape.isWriteProtected
}

func (tape *FileSystemTapeDevice) SetIsReady(flag bool) {
	tape.isReady = flag
}

func (tape *FileSystemTapeDevice) SetIsWriteProtected(flag bool) {
	tape.isWriteProtected = flag
}

func IsValidReelName(name string) bool {
	if len(name) < 1 || len(name) > 6 {
		return false
	}

	for nx := 0; nx < len(name); nx++ {
		if (name[nx] < 'A' || name[nx] > 'Z') && (name[nx] < '0' || name[nx] > '9') {
			return false
		}
	}

	return true
}

func (tape *FileSystemTapeDevice) StartIo(pkt IoPacket) {
	pkt.SetIoStatus(IosInProgress)

	if pkt.GetNodeDeviceType() != tape.GetNodeDeviceType() {
		pkt.SetIoStatus(IosInvalidNodeType)
	}

	switch pkt.GetIoFunction() {
	case IofMount:
		tape.doMount(pkt.(*TapeIoPacket))
	case IofRead:
		tape.doRead(pkt.(*TapeIoPacket))
	case IofReset:
		tape.doReset(pkt.(*TapeIoPacket))
	case IofUnmount:
		tape.doUnmount(pkt.(*TapeIoPacket))
	case IofWrite:
		tape.doWrite(pkt.(*TapeIoPacket))
	default:
		pkt.SetIoStatus(IosInvalidFunction)
	}
}

func (tape *FileSystemTapeDevice) doMount(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

func (tape *FileSystemTapeDevice) doRead(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

// doReset cancels any pending IO and unmounts the media
func (tape *FileSystemTapeDevice) doReset(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

func (tape *FileSystemTapeDevice) doUnmount(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

func (tape *FileSystemTapeDevice) doWrite(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

func (tape *FileSystemTapeDevice) Dump(destination io.Writer, indent string) {
	// TODO
}
