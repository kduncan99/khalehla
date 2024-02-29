// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type MFDFileType uint

const (
	FileTypeFixed     = 000
	FileTypeTape      = 001
	FileTypeRemovable = 040
)

func NewMFDFileTypeFromField(field uint64) MFDFileType {
	switch field {
	case 001:
		return FileTypeTape
	case 040:
		return FileTypeRemovable
	default:
		return FileTypeFixed
	}
}
