// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package hardware

// BlockCount represents a number of pseudo-physical blocks.
// For disk nodes, a block contains a fixed number of words which corresponds to the relevant medium's prep factor.
// For tape nodes, a block contains a variable number of words.
type BlockCount uint64

// BlockId uniquely identifies a particular pseudo-physical block on a disk medium
type BlockId uint64

// BlockSize describes the size of some type of block.
// Usually, it is intended for describing the number of bytes in a block of bytes.
type BlockSize uint32

// PrepFactor indicates the number of words stored in a block of data for disk media.
// Current valid values include 28, 56, 112, 224, 448, 896, and 1792.
type PrepFactor uint

// TrackCount represents a number of software tracks, each of which contain 1792 words of storage
type TrackCount uint

// TrackId represents a software track Identifier, relative to the start of a particular pack
type TrackId uint

var BlockSizeFromPrepFactor = map[PrepFactor]BlockSize{
	28:   128,  // slop 2 bytes
	56:   256,  // slop 4 bytes
	112:  512,  // slop 8 bytes
	224:  1024, // slop 16 bytes
	448:  2048, // slop 32 bytes
	896:  4096, // slop 64 bytes
	1792: 8192, // slop 128 bytes
}

var PrepFactorFromBlockSize = map[BlockSize]PrepFactor{
	28:   128,  // slop 2 bytes
	56:   256,  // slop 4 bytes
	112:  512,  // slop 8 bytes
	224:  1024, // slop 16 bytes
	448:  2048, // slop 32 bytes
	896:  4096, // slop 64 bytes
	1792: 8192, // slop 128 bytes
}
