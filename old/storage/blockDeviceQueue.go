package storage

type blockDeviceQueue struct {
	channel chan *BlockIORequest
	device  *BlockDevice
}

func (dq *blockDeviceQueue) Close() AggregatorResult {
	res := (*dq.device).Close()
	if res.status != DeviceStatusSuccessful {
		return AggregatorResult{AggregatorStatusDeviceError, &res}
	} else {
		return AggregatorResult{AggregatorStatusSuccessful, nil}
	}
}

func (dq *blockDeviceQueue) Open(writeProtected bool, writeThrough bool) AggregatorResult {
	res := (*dq.device).Open(writeProtected, writeThrough)
	if res.status != DeviceStatusSuccessful {
		return AggregatorResult{AggregatorStatusDeviceError, &res}
	} else {
		return AggregatorResult{AggregatorStatusSuccessful, nil}
	}
}

func (dq *blockDeviceQueue) routine() {
	ok := true
	for ok {
		ioreq, ok := <-dq.channel
		if ok {
			if ioreq.function == AggregatorFunctionAllocate {
				// TODO
			} else if ioreq.function == AggregatorFunctionRelease {
				// TODO
			} else if ioreq.function == AggregatorFunctionRead {
				// TODO
			} else if ioreq.function == AggregatorFunctionWrite {
				// TODO
			} else {
				ioreq.aggregatorStatus = AggregatorStatusInvalidFunction
			}
		}
	}
}

func NewBlockDeviceQueue(device *BlockDevice) blockDeviceQueue {
	return blockDeviceQueue{
		channel: make(chan *BlockIORequest),
		device:  device,
	}
}
