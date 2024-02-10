package deviceMgr

type TapeIoPacket struct {
	deviceIdentifier NodeIdentifier
	ioFunction       IoFunction
	ioStatus         IoStatus
	fileName         string // for mount
	writeProtected   bool   // for mount
}

func (pkt *TapeIoPacket) GetDeviceIdentifier() NodeIdentifier {
	return pkt.deviceIdentifier
}

func (pkt *TapeIoPacket) GetNodeType() NodeType {
	return NodeTypeTape
}

func (pkt *TapeIoPacket) GetIoFunction() IoFunction {
	return pkt.ioFunction
}

func (pkt *TapeIoPacket) GetIoStatus() IoStatus {
	return pkt.ioStatus
}

func (pkt *TapeIoPacket) SetIoStatus(ioStatus IoStatus) {
	pkt.ioStatus = ioStatus
}

func NewTapeIoPacketMount(deviceIdentifier NodeIdentifier, fileName string, writeProtected bool) *TapeIoPacket {
	return &TapeIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       IofMount,
		ioStatus:         IosNotStarted,
		fileName:         fileName,
		writeProtected:   writeProtected,
	}
}

func NewTapeIoPacketReset(deviceIdentifier NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       IofReset,
		ioStatus:         IosNotStarted,
	}
}

func NewTapeIoPacketUnmount(deviceIdentifier NodeIdentifier) *TapeIoPacket {
	return &TapeIoPacket{
		deviceIdentifier: deviceIdentifier,
		ioFunction:       IofUnmount,
		ioStatus:         IosNotStarted,
	}
}
