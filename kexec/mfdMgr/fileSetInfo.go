// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/pkg"
	"strings"
)

type FileSetCycleInfo struct {
	toBeCataloged bool
	toBeDropped   bool
	absoluteCycle uint
}

type FileSetInfo struct {
	qualifier       string
	filename        string
	projectId       string
	readKey         string
	writeKey        string
	fileType        FileType
	plusOneExists   bool
	count           uint
	maxCycleRange   uint
	currentRange    uint
	highestAbsolute uint
	cycleInfo       []FileSetCycleInfo
}

// PopulateFromLeadItems populates the FileSetInfo object from the content of the
// given leadItem0 and (optional) leadItem1 sectors.
// If there is no leadItem1, that argument should be nil.
func (fsi *FileSetInfo) PopulateFromLeadItems(leadItem0 []pkg.Word36, leadItem1 []pkg.Word36) {
	fsi.qualifier = strings.TrimRight(leadItem0[1].ToStringAsFieldata()+leadItem0[2].ToStringAsFieldata(), " ")
	fsi.filename = strings.TrimRight(leadItem0[3].ToStringAsFieldata()+leadItem0[4].ToStringAsFieldata(), " ")
	fsi.projectId = strings.TrimRight(leadItem0[5].ToStringAsFieldata()+leadItem0[6].ToStringAsFieldata(), " ")
	fsi.readKey = strings.TrimRight(leadItem0[7].ToStringAsFieldata(), " ")
	fsi.writeKey = strings.TrimRight(leadItem0[010].ToStringAsFieldata(), " ")
	fsi.fileType = NewFileTypeFromField(leadItem0[011].GetS1())
	fsi.plusOneExists = false
	fsi.count = uint(leadItem0[011].GetS2())
	fsi.maxCycleRange = uint(leadItem0[011].GetS3())
	fsi.currentRange = uint(leadItem0[011].GetS4())
	fsi.highestAbsolute = uint(leadItem0[011].GetT3())
	fsi.cycleInfo = make([]FileSetCycleInfo, fsi.maxCycleRange)

	leadItems := [][]pkg.Word36{leadItem0, leadItem1}
	absCycle := fsi.highestAbsolute
	lx := 0
	wx := 11 + leadItem0[0].GetS4()
	for ax := 0; ax < int(fsi.maxCycleRange); ax++ {
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
			fsi.cycleInfo[ax] = FileSetCycleInfo{
				toBeCataloged: w&0_200000_000000 != 0,
				toBeDropped:   w&0_100000_000000 != 0,
				absoluteCycle: absCycle,
			}
		}
	}
}
