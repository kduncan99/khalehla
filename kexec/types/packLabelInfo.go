// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

import (
	"fmt"
	"khalehla/pkg"
	"strings"
)

// PackLabelInfo contains all the information available in a disk pack label, extracted into discrete data fields
type PackLabelInfo struct {
	PackId                     string
	FirstDirectoryTrackAddress DeviceRelativeWordAddress
	RecordsPerTrack            uint
	WordsPerRecord             uint       // aka PrepFactor
	SystemReserveSize          uint       // for DRS packs (which we don't do) in words
	HMBTPaddedSize             uint       // size of s0+s1+HMBT+pad to next physical record boundary (we don't do MBTs)
	MBTSize                    uint       // size of MBT in words, including control word and checksum (we don't do MBts)
	PreppedByLevel             string     // level of whoever prepped the thing
	PreppedBy                  uint       // 040 DPREP, 020 TPREP(?), 010 workstation utility
	VOL1Version                uint       // we always do version 1
	HeadsPerCylinder           uint       // mostly used for diagnostics - we might not do this ever, at all
	TrackCount                 TrackCount // may or may not include initial allocation (depends on things)
	PrepFactor                 PrepFactor
	Attributes                 uint64 // we don't use this, but it must be 0_000000_000006
}

func NewPackLabelInfo(buffer []pkg.Word36) (*PackLabelInfo, bool) {
	pl := &PackLabelInfo{}
	if !pl.PopulateFrom(buffer) {
		return nil, false
	} else {
		return pl, true
	}
}

// PopulateFrom returns false if the buffer does not have "VOL1" in word 0
func (pl *PackLabelInfo) PopulateFrom(buffer []pkg.Word36) bool {
	vol1 := buffer[0].ToStringAsAscii()
	if vol1 != "VOL1" {
		return false
	}

	pl.PackId = strings.TrimRight(buffer[1].ToStringAsAscii()+buffer[2].ToStringAsAscii(), " ")
	pl.FirstDirectoryTrackAddress = DeviceRelativeWordAddress(buffer[3].GetW())
	pl.RecordsPerTrack = uint(buffer[4].GetH1())
	pl.WordsPerRecord = uint(buffer[4].GetH2())
	pl.SystemReserveSize = uint(buffer[5].GetH2())
	pl.HMBTPaddedSize = uint(buffer[011].GetH1())
	pl.MBTSize = uint(buffer[011].GetH2())
	pl.PreppedByLevel = strings.TrimRight(buffer[012].ToStringAsAscii(), " ")
	pl.PreppedBy = uint(buffer[014].GetS1())
	pl.VOL1Version = uint(buffer[014].GetS2())
	pl.HeadsPerCylinder = uint(buffer[014].GetH2())
	pl.TrackCount = TrackCount(buffer[016].GetW())
	pl.PrepFactor = PrepFactor(buffer[017].GetH1())
	pl.Attributes = buffer[020].GetW()

	return true
}

func (pl *PackLabelInfo) WriteTo(buffer []pkg.Word36) {
	for wx := 0; wx < 28; wx++ {
		buffer[wx].SetW(0)
	}

	buffer[0].FromStringToAscii("VOL1")

	tempPackId := fmt.Sprintf("-%6s", pl.PackId[:6])
	buffer[1].FromStringToAscii(tempPackId[:4])
	buffer[2].FromStringToAscii(tempPackId[4:])
	buffer[2].SetH2(0)

	buffer[3].SetW(uint64(pl.FirstDirectoryTrackAddress))
	buffer[4].SetH1(uint64(pl.RecordsPerTrack))
	buffer[4].SetH2(uint64(pl.WordsPerRecord))
	buffer[5].SetH2(uint64(pl.SystemReserveSize))
	buffer[011].SetH1(uint64(pl.HMBTPaddedSize))
	buffer[011].SetH2(uint64(pl.MBTSize))
	buffer[012].FromStringToAscii(fmt.Sprintf("-%4s", pl.PreppedByLevel[:4]))
	buffer[014].SetS1(uint64(pl.PreppedBy))
	buffer[014].SetS2(uint64(pl.VOL1Version))
	buffer[014].SetH2(uint64(pl.HeadsPerCylinder))
	buffer[016].SetW(uint64(pl.TrackCount))
	buffer[017].SetH1(uint64(pl.PrepFactor))
	buffer[020].SetW(pl.Attributes)
}
