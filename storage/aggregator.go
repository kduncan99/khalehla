package storage

import "kalehla/types"

type AggregatorFunction int
type AggregatorStatus int
type AggregatorType int

const (
	AggregatorFunctionAllocate = iota
	AggregatorFunctionRead
	AggregatorFunctionRelease
	AggregatorFunctionWrite
	//	TODO others?
)

const (
	AggregatorTypeSimple = iota
	AggregatorTypeCache
	AggregatorTypeDedupe
	//	TODO others?
)

const (
	AggregatorStatusSuccessful = iota
	AggregatorStatusNotOpen
	AggregatorStatusAlreadyOpen
	AggregatorStatusInvalidDeviceIndex
	AggregatorStatusInvalidFunction
	AggregatorStatusInProgress
	AggregatorStatusDeviceError
	AggregatorStatusSystemError
	//	TODO others
)

type BlockIORequest struct {
	function         AggregatorFunction
	deviceIndex      types.DeviceIndex
	blockId          types.BlockId
	blockCount       types.BlockCount
	buffer           []types.Word36
	aggregatorStatus AggregatorStatus
	deviceStatus     types.DeviceStatus
	systemError      error
}

type AggregatorResult struct {
	aggregatorStatus AggregatorStatus
	deviceResult     *DeviceResult
}

// Aggregator collects a set of block devices and accepts IO on their behalf.
// Aggregators are generally expected to do something useful besides just acting as a bottleneck.
// The I/O contract assumes requests are handled asynchronously, and indeed, that would be one
// good use of an aggregator. Other uses might be caching, load or capacity balancing, dedupe, etc.
// There is no requirement (although there might be an expectation) that the I/O address as
// understood by the client (the OS, generally) actually corresponds to a real device and block...
// In particular, dedupe will completely change this.
type Aggregator interface {
	Close() AggregatorResult
	Open() AggregatorResult
	GetDevice(deviceIndex types.DeviceIndex) (*BlockDevice, AggregatorResult)
	IsOpen() bool
	RegisterDevice(deviceIndex types.DeviceIndex, device *BlockDevice) AggregatorResult
	StartIO(request *BlockIORequest)
}

func GetBlockGeometry(agg Aggregator, deviceIndex types.DeviceIndex) (BlockGeometry, AggregatorResult) {
	dev, res := agg.GetDevice(deviceIndex)
	if res.aggregatorStatus == AggregatorStatusSuccessful {
		geo, res := (*dev).GetGeometry()
		if res.status == DeviceStatusSuccessful {
			return geo, AggregatorResult{AggregatorStatusSuccessful, &res}
		} else {
			return BlockGeometry{}, AggregatorResult{AggregatorStatusDeviceError, &res}
		}
	} else {
		return BlockGeometry{}, AggregatorResult{AggregatorStatusInvalidDeviceIndex, nil}
	}
}
