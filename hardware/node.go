// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package hardware

import (
	"io"
	"khalehla/pkg"
)

// NodeIdentifier uniquely identifies a particular device or channel (or anything else identifiable which we manage)
// It is currently implemented as the 1-6 character device Name, all caps alphas and/or digits LJSF
// stored as Fieldata in a Word36 struct
type NodeIdentifier pkg.Word36

// BlockCount represents a number of pseudo-physical blocks.
// For disk nodes, a block contains a fixed number of words which corresponds to the relevant medium's prep factor.
// For tape nodes, a block contains a variable number of words.
type BlockCount uint64

// BlockId uniquely identifies a particular pseudo-physical block on a disk medium
type BlockId uint64

// PrepFactor indicates the number of words stored in a block of data for disk media.
// Current valid values include 28, 56, 112, 224, 448, 896, and 1792.
type PrepFactor uint

// TrackCount represents a number of software tracks, each of which contain 1792 words of storage
type TrackCount uint

// TrackId represents a software track Identifier, relative to the start of a particular pack
type TrackId uint

type Node interface {
	Dump(destination io.Writer, indent string)
	GetNodeCategoryType() NodeCategoryType
	GetNodeDeviceType() NodeDeviceType
	GetNodeModelType() NodeModelType
	IsReady() bool
}
