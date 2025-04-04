package storage

import "khalehla/pkg"

// SimpleAggregator is a simple coordinator for a set of block devices.
// It manages async IO across the multiple devices in a quasi-efficient manner.
type SimpleAggregator struct {
	deviceQueues map[pkg.DeviceIndex]blockDeviceQueue
	isOpen       bool
}

func (agg *SimpleAggregator) Close() AggregatorResult {
	if !agg.IsOpen() {
		return AggregatorResult{AggregatorStatusNotOpen, nil}
	}

	for _, dq := range agg.deviceQueues {
		dq.Close()
	}
	agg.isOpen = false

	return AggregatorResult{AggregatorStatusSuccessful, nil}
}

func (agg *SimpleAggregator) Open() AggregatorResult {
	if agg.IsOpen() {
		return AggregatorResult{AggregatorStatusAlreadyOpen, nil}
	}

	for _, dq := range agg.deviceQueues {
		res := dq.Open(false, true)
		if res.aggregatorStatus != AggregatorStatusSuccessful {
			for _, dq2 := range agg.deviceQueues {
				_ = dq2.Close()
			}
			return res
		}
	}

	return AggregatorResult{AggregatorStatusSuccessful, nil}
}

func (agg *SimpleAggregator) GetDevice(deviceIndex pkg.DeviceIndex) (*BlockDevice, AggregatorResult) {
	dev, ok := agg.deviceQueues[deviceIndex]
	if ok {
		return dev.device, AggregatorResult{AggregatorStatusSuccessful, nil}
	} else {
		return nil, AggregatorResult{AggregatorStatusInvalidDeviceIndex, nil}
	}
}

func (agg *SimpleAggregator) IsOpen() bool {
	return agg.isOpen
}

func (agg *SimpleAggregator) RegisterDevice(deviceIndex pkg.DeviceIndex, device *BlockDevice) AggregatorResult {
	_, ok := agg.deviceQueues[deviceIndex]
	if ok {
		return AggregatorResult{AggregatorStatusInvalidDeviceIndex, nil}
	}

	dq := NewBlockDeviceQueue(device)
	agg.deviceQueues[deviceIndex] = dq
	go dq.routine()

	return AggregatorResult{AggregatorStatusSuccessful, nil}
}

func (agg *SimpleAggregator) StartIO(request *BlockIORequest) {
	dq, ok := agg.deviceQueues[request.deviceIndex]
	if ok {
		request.aggregatorStatus = AggregatorStatusInProgress
		dq.channel <- request
	} else {
		request.aggregatorStatus = AggregatorStatusInvalidDeviceIndex
	}
}

func NewSimpleAggregator() *SimpleAggregator {
	return &SimpleAggregator{
		deviceQueues: make(map[pkg.DeviceIndex]blockDeviceQueue),
	}
}
