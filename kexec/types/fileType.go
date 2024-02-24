// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

type FileType uint

const (
	FileTypeFixed     = 000
	FileTypeTape      = 001
	FileTypeRemovable = 040
)

func NewFileTypeFromField(field uint64) FileType {
	switch field {
	case 001:
		return FileTypeTape
	case 040:
		return FileTypeRemovable
	default:
		return FileTypeFixed
	}
}
