// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// FileCycleInfo is an interface for the more-specific fixed, removable, and tape file cycle info structs
type FileCycleInfo interface {
	GetFileSetIdentifier() FileSetIdentifier
	GetFileCycleIdentifier() FileCycleIdentifier
	GetQualifier() string
	GetFilename() string
	GetInhibitFlags() InhibitFlags
	GetProjectId() string
	GetAccountId() string
	GetAbsoluteFileCycle() uint
	GetAssignMnemonic() string
	setFileCycleIdentifier(fcIdentifier FileCycleIdentifier)
	setFileSetIdentifier(fsIdentifier FileSetIdentifier)
}

// FileCycleIdentifier is a unique opaque identifier allowing clients to refer to a file cycle
// without using qualifier, filename, and file cycle. Internally it is the main item sector 0 address
// for the file cycle - but clients should not be concerned with, nor rely on, that.
type FileCycleIdentifier uint64

// DiskPackEntry describes a particular disk pack as it is referred to by a particular file cycle entity.
// It is contained within one or more of the more-specific file cycle info structs.
type DiskPackEntry struct {
	PackName     string
	MainItemLink uint64
}
