// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

// MFDInhibitFlags is a 6-bit field found in main item sector 0.
// It contains certain flags regarding inhibits which apply to a file.
type MFDInhibitFlags struct {
	isGuarded           bool
	isUnloadInhibited   bool
	isPrivate           bool
	isAssignedExclusive bool
	isWriteOnly         bool
	isReadOnly          bool
}

func (inf *MFDInhibitFlags) Compose() uint64 {
	value := uint64(0)
	if inf.isGuarded {
		value |= 040
	}
	if inf.isUnloadInhibited {
		value |= 020
	}
	if inf.isPrivate {
		value |= 010
	}
	if inf.isAssignedExclusive {
		value |= 004
	}
	if inf.isWriteOnly {
		value |= 002
	}
	if inf.isReadOnly {
		value |= 001
	}
	return value
}

func (inf *MFDInhibitFlags) ExtractFrom(field uint64) {
	inf.isGuarded = field&040 != 0
	inf.isUnloadInhibited = field&020 != 0
	inf.isPrivate = field&010 != 0
	inf.isAssignedExclusive = field&004 != 0
	inf.isWriteOnly = field&002 != 0
	inf.isReadOnly = field&001 != 0
}
