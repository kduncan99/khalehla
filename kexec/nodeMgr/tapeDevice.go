// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package nodeMgr

import (
	"io"
	"os"
	"sync"
)

// This is a very simple pseudo tape device

type TapeDevice struct {
	fileName         *string
	file             *os.File
	isReady          bool
	isWriteProtected bool
	mutex            sync.Mutex
}

func NewTapeDevice() *TapeDevice {
	return &TapeDevice{
		isWriteProtected: true,
	}
}

func (tape *TapeDevice) GetNodeCategoryType() NodeCategoryType {
	return NodeCategoryDevice
}

func (tape *TapeDevice) GetNodeDeviceType() NodeDeviceType {
	return NodeDeviceTape
}

func (tape *TapeDevice) GetNodeModelType() NodeModelType {
	return NodeModelFileSystemTapeDevice
}

func (tape *TapeDevice) IsMounted() bool {
	return tape.file != nil
}

func (tape *TapeDevice) IsReady() bool {
	return tape.isReady
}

func (tape *TapeDevice) IsWriteProtected() bool {
	return tape.isWriteProtected
}

func (tape *TapeDevice) SetIsReady(flag bool) {
	tape.isReady = flag
}

func (tape *TapeDevice) SetIsWriteProtected(flag bool) {
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

func (tape *TapeDevice) StartIo(pkt IoPacket) {
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

func (tape *TapeDevice) doMount(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

func (tape *TapeDevice) doRead(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

// doReset cancels any pending IO and unmounts the media
func (tape *TapeDevice) doReset(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

func (tape *TapeDevice) doUnmount(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

func (tape *TapeDevice) doWrite(pkt *TapeIoPacket) {
	tape.mutex.Lock()
	defer tape.mutex.Unlock()
	// TODO
	pkt.ioStatus = IosSystemError
}

func (tape *TapeDevice) Dump(destination io.Writer, indent string) {
	// TODO
}
