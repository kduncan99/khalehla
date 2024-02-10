package deviceMgr

import (
	"os"
	"sync"
)

// This is a very simple pseudo tape device

type TapeDevice struct {
	fileName     *string
	file         *os.File
	writeProtect bool
	mutex        sync.Mutex
}

func NewTapeDevice() *TapeDevice {
	return &TapeDevice{}
}

func (tape *TapeDevice) getNodeType() NodeType {
	return NodeTypeTape
}

func (tape *TapeDevice) IsMounted() bool {
	return tape.file != nil
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

func (tape *TapeDevice) startIo(pkt IoPacket) {
	pkt.SetIoStatus(IosInProgress)

	if pkt.GetNodeType() != tape.getNodeType() {
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
