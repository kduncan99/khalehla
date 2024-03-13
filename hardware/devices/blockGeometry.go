// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package devices

import "khalehla/hardware"

// A BlockGeometry struct allows a particular block device to report the included information
// regarding the actual storage medium which it controls. The implementation must be able to
// derive the information from whatever underlying storage medium it uses, such that the
// client needs only construct the medium once, providing the necessary parameters, and then
// can cause the medium to be opened and read without knowing this information ahead of time.
// Generally, the device stores this information in some convenient format in block zero of
// the storage - and would then reject any attempt to overwrite block 0.
type BlockGeometry struct {
	BytesPerBlock  hardware.BlockSize
	WordsPerBlock  hardware.PrepFactor
	BlocksPerTrack hardware.BlockCount
	BlockCount     hardware.BlockCount
	TrackCount     hardware.TrackCount
	Label          string
}
