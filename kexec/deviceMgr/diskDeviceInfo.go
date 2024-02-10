package deviceMgr

import "khalehla/pkg"

// -------------------------------------------------------------------------------------

type DiskDeviceInfo struct {
	deviceName      string
	nodeIdentifier  NodeIdentifier
	initialFileName *string
	device          *DiskDevice
	nodeStatus      NodeStatus
	isAccessible    bool // can only be true if status is UP, RV, or SU and the device is assigned to at least one channel
	isMounted       bool
	isPrepped       bool
	isFixed         bool
}

// NewDiskDeviceInfo creates a new struct
// deviceName is required, but initialFileName can be set to nil if the device is not to be initially mounted
func NewDiskDeviceInfo(deviceName string, initialFileName *string) *DiskDeviceInfo {
	return &DiskDeviceInfo{
		deviceName:      deviceName,
		nodeIdentifier:  NodeIdentifier(pkg.NewFromStringToFieldata(deviceName, 1)[0]),
		nodeStatus:      NodeStatusUp,
		isAccessible:    false,
		initialFileName: initialFileName,
		isMounted:       false,
		isPrepped:       false,
		isFixed:         false,
	}
}

func (ddi *DiskDeviceInfo) CreateNode() {
	ddi.device = NewDiskDevice(ddi.initialFileName)
}

func (ddi *DiskDeviceInfo) GetDevice() Device {
	return ddi.device
}

func (ddi *DiskDeviceInfo) GetInitialFileName() *string {
	return ddi.initialFileName
}

func (ddi *DiskDeviceInfo) GetNodeIdentifier() NodeIdentifier {
	return ddi.nodeIdentifier
}

func (ddi *DiskDeviceInfo) GetNodeName() string {
	return ddi.deviceName
}

func (ddi *DiskDeviceInfo) GetNodeStatus() NodeStatus {
	return ddi.nodeStatus
}

func (ddi *DiskDeviceInfo) GetNodeType() NodeType {
	return NodeTypeDisk
}

func (ddi *DiskDeviceInfo) IsAccessible() bool {
	return ddi.isAccessible
}

func (ddi *DiskDeviceInfo) IsFixed() bool {
	return ddi.isFixed
}

func (ddi *DiskDeviceInfo) IsMounted() bool {
	return ddi.isMounted
}

func (ddi *DiskDeviceInfo) IsPrepped() bool {
	return ddi.isPrepped
}

func (ddi *DiskDeviceInfo) SetIsAccessible(isAccessible bool) {
	ddi.isAccessible = isAccessible
}
