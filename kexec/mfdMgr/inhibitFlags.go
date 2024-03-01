// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// InhibitFlags is a 6-bit field found in main item sector 0.
// It contains certain flags regarding inhibits which apply to a file.
type InhibitFlags struct {
	isGuarded           bool
	isUnloadInhibited   bool
	isPrivate           bool
	isAssignedExclusive bool
	isWriteOnly         bool
	isReadOnly          bool
}

func (pf *InhibitFlags) Compose() uint64 {
	value := uint64(0)
	if pf.isGuarded {
		value |= 040
	}
	if pf.isUnloadInhibited {
		value |= 020
	}
	if pf.isPrivate {
		value |= 010
	}
	if pf.isAssignedExclusive {
		value |= 004
	}
	if pf.isWriteOnly {
		value |= 002
	}
	if pf.isReadOnly {
		value |= 001
	}
	return value
}

func (pf *InhibitFlags) ExtractFrom(field uint64) {
	pf.isGuarded = field&040 != 0
	pf.isUnloadInhibited = field&020 != 0
	pf.isPrivate = field&010 != 0
	pf.isAssignedExclusive = field&004 != 0
	pf.isWriteOnly = field&002 != 0
	pf.isReadOnly = field&001 != 0
}
