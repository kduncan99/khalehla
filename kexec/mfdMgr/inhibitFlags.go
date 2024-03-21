// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// InhibitFlags is a 6-bit field found in main item sector 0.
// It contains certain flags regarding inhibits which apply to a file.
type InhibitFlags struct {
	IsGuarded           bool
	IsUnloadInhibited   bool
	IsPrivate           bool
	IsAssignedExclusive bool
	IsWriteOnly         bool
	IsReadOnly          bool
}

func (inf *InhibitFlags) Compose() uint64 {
	value := uint64(0)
	if inf.IsGuarded {
		value |= 040
	}
	if inf.IsUnloadInhibited {
		value |= 020
	}
	if inf.IsPrivate {
		value |= 010
	}
	if inf.IsAssignedExclusive {
		value |= 004
	}
	if inf.IsWriteOnly {
		value |= 002
	}
	if inf.IsReadOnly {
		value |= 001
	}
	return value
}

func (inf *InhibitFlags) ExtractFrom(field uint64) {
	inf.IsGuarded = field&040 != 0
	inf.IsUnloadInhibited = field&020 != 0
	inf.IsPrivate = field&010 != 0
	inf.IsAssignedExclusive = field&004 != 0
	inf.IsWriteOnly = field&002 != 0
	inf.IsReadOnly = field&001 != 0
}

func ExtractNewInhibitFlags(field uint64) *InhibitFlags {
	inf := &InhibitFlags{}
	inf.ExtractFrom(field)
	return inf
}
