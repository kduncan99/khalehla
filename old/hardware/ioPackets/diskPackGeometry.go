// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ioPackets

import "khalehla/hardware"

// DiskPackGeometry describes various useful attributes of a particular prepped disk pack
type DiskPackGeometry struct {
	PrepFactor           hardware.PrepFactor // number of words contained in a physical record, packed
	TrackCount           hardware.TrackCount // number of tracks on the pack
	BlockCount           hardware.BlockCount // number of blocks on the pack (may not be track-aligned)
	SectorsPerBlock      uint                // number of software sectors (28 words per) in a physical block
	BlocksPerTrack       uint                // physical blocks required to contain one software track (1792 words)
	BytesPerBlock        uint                // bytes consumed strictly by the word 36 Buffer
	PaddedBytesPerBlock  uint                // bytes required for a block containing packed word36 structs, rounded to power of 2
	FirstDirTrackBlockId uint                // block ID of the block containing the first sector of the first directory track
}
