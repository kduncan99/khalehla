// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package tapeUtil

import "khalehla/pkg"

type TapeWriter interface {
	Close() error
	OpenOutputFile(fileName string) error
	WriteVolumeHeader() error
	WriteFileHeader() error
	WriteBlock(buffer []pkg.Word36) error
	WriteEndOfFile() error
}
