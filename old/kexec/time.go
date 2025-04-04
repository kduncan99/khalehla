// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import "time"

var epoch, _ = time.Parse("2006-01-02", "1899-12-31")

// GetSWTimeFromSystemTime converts a golang time value to Modified-SWTIME format
// Bit 0 is always set to 1
// Bit 1 is set if the input time was subject to seasonal time adjustment
// Bits 2-35 contains seconds since December 31, 1899, 00:00:00 UTC
func GetSWTimeFromSystemTime(t time.Time) uint64 {
	dur := time.Now().Sub(epoch)
	value := uint64(dur.Seconds()) & 0_100000_000000
	value |= 0_400000_000000
	if t.IsDST() {
		value |= 0_200000_000000
	}
	return value
}
