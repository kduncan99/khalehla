// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// FileFlags contain miscellaneous information about a file cycle
type FileFlags struct {
	IsLargeFile            bool
	AssignmentAcceleration bool
	IsWrittenTo            bool
	StoreThrough           bool
}

func (ff *FileFlags) Compose() uint64 {
	value := uint64(0)
	if ff.IsLargeFile {
		value |= 040
	}
	if ff.AssignmentAcceleration {
		value |= 004
	}
	if ff.IsWrittenTo {
		value |= 002
	}
	if ff.StoreThrough {
		value |= 001
	}
	return value
}

func (ff *FileFlags) ExtractFrom(field uint64) {
	ff.IsLargeFile = field&040 != 0
	ff.AssignmentAcceleration = field&004 != 0
	ff.IsWrittenTo = field&002 != 0
	ff.StoreThrough = field&001 != 0
}
