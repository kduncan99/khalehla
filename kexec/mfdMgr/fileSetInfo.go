// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/pkg"
	"strings"
)

type FileSetCycleInfo struct {
	ToBeCataloged bool
	ToBeDropped   bool
	AbsoluteCycle uint
}

type FileSetInfo struct {
	Qualifier       string
	Filename        string
	ProjectId       string
	ReadKey         string
	WriteKey        string
	FileType        FileType
	PlusOneExists   bool
	Count           uint
	MaxCycleRange   uint
	CurrentRange    uint
	HighestAbsolute uint
	CycleInfo       []FileSetCycleInfo
}

// PopulateFromLeadItems populates the FileSetInfo object from the Content of the
// given leadItem0 and (optional) leadItem1 sectors.
// If there is no leadItem1, that argument should be nil.
func (fsi *FileSetInfo) PopulateFromLeadItems(leadItem0 []pkg.Word36, leadItem1 []pkg.Word36) {
	fsi.Qualifier = strings.TrimRight(leadItem0[1].ToStringAsFieldata()+leadItem0[2].ToStringAsFieldata(), " ")
	fsi.Filename = strings.TrimRight(leadItem0[3].ToStringAsFieldata()+leadItem0[4].ToStringAsFieldata(), " ")
	fsi.ProjectId = strings.TrimRight(leadItem0[5].ToStringAsFieldata()+leadItem0[6].ToStringAsFieldata(), " ")
	fsi.ReadKey = strings.TrimRight(leadItem0[7].ToStringAsFieldata(), " ")
	fsi.WriteKey = strings.TrimRight(leadItem0[010].ToStringAsFieldata(), " ")
	fsi.FileType = NewFileTypeFromField(leadItem0[011].GetS1())
	fsi.PlusOneExists = false
	fsi.Count = uint(leadItem0[011].GetS2())
	fsi.MaxCycleRange = uint(leadItem0[011].GetS3())
	fsi.CurrentRange = uint(leadItem0[011].GetS4())
	fsi.HighestAbsolute = uint(leadItem0[011].GetT3())
	fsi.CycleInfo = make([]FileSetCycleInfo, fsi.MaxCycleRange)

	leadItems := [][]pkg.Word36{leadItem0, leadItem1}
	absCycle := fsi.HighestAbsolute
	lx := 0
	wx := 11 + leadItem0[0].GetS4()
	for ax := 0; ax < int(fsi.MaxCycleRange); ax++ {
		if wx == 28 {
			if leadItem1 == nil {
				break
			}
			lx++
			wx = 1
		}
		w := leadItems[lx][wx].GetW()
		link := w & 0_077777_777777
		if link > 0 {
			fsi.CycleInfo[ax] = FileSetCycleInfo{
				ToBeCataloged: w&0_200000_000000 != 0,
				ToBeDropped:   w&0_100000_000000 != 0,
				AbsoluteCycle: absCycle,
			}
		}
	}
}
