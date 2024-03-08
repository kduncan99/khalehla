// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// FileFlags contain miscellaneous information about a file cycle
// It is contained within one or more of the more-specific file cycle info structs.
type FileFlags struct {
	IsLargeFile            bool
	AssignmentAcceleration bool
	IsWrittenTo            bool
	StoreThrough           bool
}

func (fif *FileFlags) Compose() uint64 {
	value := uint64(0)
	if fif.IsLargeFile {
		value |= 040
	}
	if fif.AssignmentAcceleration {
		value |= 004
	}
	if fif.IsWrittenTo {
		value |= 002
	}
	if fif.StoreThrough {
		value |= 001
	}
	return value
}

func (fif *FileFlags) ExtractFrom(field uint64) {
	fif.IsLargeFile = field&040 != 0
	fif.AssignmentAcceleration = field&004 != 0
	fif.IsWrittenTo = field&002 != 0
	fif.StoreThrough = field&001 != 0
}

func ExtractNewFileFlags(field uint64) *FileFlags {
	ff := &FileFlags{}
	ff.ExtractFrom(field)
	return ff
}
