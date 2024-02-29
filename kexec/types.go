// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"khalehla/pkg"
)

// BlockCount represents a number of pseudo-physical blocks.
// For disk Nodes, a block contains a fixed number of words which corresponds to the relevant medium's prep factor.
// For tape Nodes, a block contains a variable number of words.
type BlockCount uint64

// BlockId uniquely identifies a particular pseudo-physical block on a disk medium
type BlockId uint64

// ConsoleIdentifier is a 36-bit word implemented as uint64, containing a unique identifier for the console.
// It *might* be the console name, but that is merely by convention.
type ConsoleIdentifier pkg.Word36

// NodeIdentifier uniquely identifies a particular device or channel (or anything else identifiable which we manage)
// It is currently implemented as the 1-6 character device name, all caps alphas and/or digits LJSF
// stored as Fieldata in a Word36 struct
type NodeIdentifier pkg.Word36

// PrepFactor indicates the number of words stored in a block of data for disk media.
// Current valid values include 28, 56, 112, 224, 448, 896, and 1792.
type PrepFactor uint

// TrackCount represents a number of software tracks, each of which contain 1792 words of storage
type TrackCount uint

// TrackId represents a software track identifier, relative to the start of a particular pack
type TrackId uint

// pseudo-enumeration pkg

type ContingencyType uint
type DeviceRelativeWordAddress uint
type ExecPhase uint
type FacStatus uint
type Granularity uint
type LDATIndex uint64
type MFDBlockId uint64
type MFDRelativeAddress uint64
type MFDSectorId uint64
type MFDTrackId uint64
type StopCode uint

// DiskPackGeometry describes various useful attributes of a particular prepped disk pack
type DiskPackGeometry struct {
	PrepFactor           PrepFactor // number of words contained in a physical record, packed
	TrackCount           TrackCount // number of tracks on the pack
	BlockCount           BlockCount // number of blocks on the pack (may not be track-aligned)
	SectorsPerBlock      uint       // number of software sectors (28 words per) in a physical block
	BlocksPerTrack       uint       // physical blocks required to contain one software track (1792 words)
	BytesPerBlock        uint       // bytes consumed strictly by the word 36 buffer
	PaddedBytesPerBlock  uint       // bytes required for a block containing packed word36 structs, rounded to power of 2
	FirstDirTrackBlockId uint       // block ID of the block containing the first sector of the first directory track
}

// Things related to ConsoleManager
// TODO move them to ConsoleManager maybe?

type ConsoleReadOnlyMessage struct {
	Source         *RunControlEntry
	Routing        *ConsoleIdentifier // may be nil
	RunId          *string            // for logging purposes, may not match RCE - may be nil
	Text           string             // message to be displayed
	DoNotEmitRunId bool
}

type ConsoleReadReplyMessage struct {
	Source         *RunControlEntry
	Routing        *ConsoleIdentifier // may be nil
	RunId          *string            // for logging purposes, may not match RCE - may be nil
	Text           string             // message to be displayed
	DoNotEmitRunId bool
	DoNotLogReply  bool // reply may contain secure information
	MaxReplyLength int
	Reply          string
}
