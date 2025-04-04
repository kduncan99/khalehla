// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tapeUtil

import (
	"khalehla/old/pkg"
)

type TapeReader interface {
	Close() error
	OpenInputFile(fileName string) error
	ReadVolumeHeader() (volHeader *VolumeHeader, err error)
	ReadFileHeader() (fileHeader *FileHeader, err error)
	ReadBlock() (data []pkg.Word36, eof bool, err error)
}
