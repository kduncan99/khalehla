package deviceMgr

type TapeDeviceInfo struct {
	deviceName     string
	nodeIdentifier NodeIdentifier
	device         *TapeDevice
	nodeStatus     NodeStatus
	isMounted      bool
}

// NewTapeDeviceInfo creates a new struct
func NewTapeDeviceInfo(deviceName string) *TapeDeviceInfo {
	return &TapeDeviceInfo{
		deviceName: deviceName,
		nodeStatus: NodeStatusUp,
	}
}

func (tdi *TapeDeviceInfo) CreateNode() {
	tdi.device = NewTapeDevice()
}

func (tdi *TapeDeviceInfo) GetDevice() Device {
	return tdi.device
}

func (tdi *TapeDeviceInfo) GetNodeIdentifier() NodeIdentifier {
	return tdi.nodeIdentifier
}

func (tdi *TapeDeviceInfo) GetNodeName() string {
	return tdi.deviceName
}

func (tdi *TapeDeviceInfo) GetNodeStatus() NodeStatus {
	return tdi.nodeStatus
}

func (tdi *TapeDeviceInfo) GetNodeType() NodeType {
	return NodeTypeTape
}

func (tdi *TapeDeviceInfo) IsMounted() bool {
	return tdi.isMounted
}
