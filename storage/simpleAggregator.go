package storage

import "kalehla/types"

// SimpleAggregator is a simple coordinator for a set of block devices.
// It manages async IO across the multiple devices in a quasi-efficient manner.
type SimpleAggregator struct {
	devices map[types.DeviceIndex]*BlockDevice
	isOpen  bool
}

func (agg *SimpleAggregator) Close() AggregatorResult {
	if !agg.IsOpen() {
		return AggregatorResult{AggregatorStatusNotOpen, nil}
	}

	for _, dev := range agg.devices {
		(*dev).Close()
	}
	agg.isOpen = false

	return AggregatorResult{AggregatorStatusSuccessful, nil}
}

func (agg *SimpleAggregator) Open() AggregatorResult {
	if agg.IsOpen() {
		return AggregatorResult{AggregatorStatusAlreadyOpen, nil}
	}

	for _, dev := range agg.devices {
		res := (*dev).Open(false, true)
		if res.status != DeviceStatusSuccessful {
			for _, dev := range agg.devices {
				_ = (*dev).Close()
			}
			return AggregatorResult{AggregatorStatusDeviceError, &res}
		}
	}

	return AggregatorResult{AggregatorStatusSuccessful, nil}
}

func (agg *SimpleAggregator) GetDevice(deviceIndex types.DeviceIndex) (*BlockDevice, AggregatorResult) {
	dev, ok := agg.devices[deviceIndex]
	if ok {
		return dev, AggregatorResult{AggregatorStatusSuccessful, nil}
	} else {
		return nil, AggregatorResult{AggregatorStatusInvalidDeviceIndex, nil}
	}
}

func (agg *SimpleAggregator) IsOpen() bool {
	return agg.isOpen
}

func (agg *SimpleAggregator) RegisterDevice(deviceIndex types.DeviceIndex, device *BlockDevice) AggregatorResult {
	_, ok := agg.devices[deviceIndex]
	if ok {
		return AggregatorResult{AggregatorStatusInvalidDeviceIndex, nil}
	}

	agg.devices[deviceIndex] = device
	return AggregatorResult{AggregatorStatusSuccessful, nil}
}

func (agg *SimpleAggregator) StartIO(request *BlockIORequest) {
	//	TODO
}

func NewSimpleAggregator() *SimpleAggregator {
	return &SimpleAggregator{
		devices: make(map[types.DeviceIndex]*BlockDevice),
	}
}
