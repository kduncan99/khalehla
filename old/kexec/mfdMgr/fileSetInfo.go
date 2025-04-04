// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"strings"

	"khalehla/old/pkg"
)

// FileSetInfo contains all the relevant information about a file set which is known to the MFD.
// It is presented to the client as a result of certain queries, and to the MFD as part of some requests.
// Note that CycleInfo is an array of pointers. If any particular entry is nil, then that entry refers
// to a file cycle which does not exist, possibly having been deleted.
// The size of the array will always correspond to the CurrentRange value.
type FileSetInfo struct {
	FileSetIdentifier     FileSetIdentifier
	Qualifier             string
	Filename              string
	ProjectId             string
	ReadKey               string
	WriteKey              string
	FileType              FileType
	Guarded               bool // at least one cycle is guarded, so file set is guarded
	PlusOneExists         bool
	Count                 uint
	MaxCycleRange         uint
	CurrentRange          uint
	HighestAbsolute       uint
	CycleInfo             []*FileSetCycleInfo
	NumberOfSecurityWords uint
}

// FileSetIdentifier is a unique opaque identifier allowing clients to refer to a fileset
// without using qualifier and filename. Internally it is the lead item sector 0 address
// for the file set - but clients should not be concerned with, nor rely on, that.
type FileSetIdentifier uint64

// FileSetCycleInfo contains a few bits of information regarding a particular file cycle.
// It is contained within a FileSetInfo struct.
type FileSetCycleInfo struct {
	ToBeCataloged       bool
	ToBeDropped         bool
	AbsoluteCycle       uint
	FileCycleIdentifier FileCycleIdentifier
}

// FileType describes whether a file is fixed, removable, or tape.
// Clients should refer only to the constant values.
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

// NewFileSetInfo populate a FileSetInfo struct.
// It is intended to be used by clients in preparation for a subsequent call on MFD services.
func NewFileSetInfo(
	qualifier string,
	filename string,
	projectId string,
	readKey string,
	writeKey string,
	fileType FileType,
) *FileSetInfo {
	return &FileSetInfo{
		Qualifier: qualifier,
		Filename:  filename,
		ProjectId: projectId,
		ReadKey:   readKey,
		WriteKey:  writeKey,
		FileType:  fileType,
	}
}

// populateFromLeadItems populates the FileSetInfo object from the Content of the given leadItem0 and
// (optional) leadItem1 sectors. If there is no leadItem1, that argument should be nil.
func (fsi *FileSetInfo) populateFromLeadItems(leadItem0 []pkg.Word36, leadItem1 []pkg.Word36) {
	fsi.Qualifier = strings.TrimRight(leadItem0[1].ToStringAsFieldata()+leadItem0[2].ToStringAsFieldata(), " ")
	fsi.Filename = strings.TrimRight(leadItem0[3].ToStringAsFieldata()+leadItem0[4].ToStringAsFieldata(), " ")
	fsi.ProjectId = strings.TrimRight(leadItem0[5].ToStringAsFieldata()+leadItem0[6].ToStringAsFieldata(), " ")
	fsi.ReadKey = strings.TrimRight(leadItem0[7].ToStringAsFieldata(), " ")
	fsi.WriteKey = strings.TrimRight(leadItem0[010].ToStringAsFieldata(), " ")
	fsi.FileType = NewFileTypeFromField(leadItem0[011].GetS1())
	fsi.Count = uint(leadItem0[011].GetS2())
	fsi.MaxCycleRange = uint(leadItem0[011].GetS3())
	fsi.CurrentRange = uint(leadItem0[011].GetS4())
	fsi.HighestAbsolute = uint(leadItem0[011].GetT3())
	fsi.Guarded = leadItem0[012]&0_400000_000000 != 0
	fsi.PlusOneExists = leadItem0[012]&0_200000_000000 != 0
	fsi.NumberOfSecurityWords = uint(leadItem0[012].GetS4())
	fsi.CycleInfo = make([]*FileSetCycleInfo, fsi.CurrentRange)

	leadItems := [][]pkg.Word36{leadItem0, leadItem1}
	absCycle := fsi.HighestAbsolute
	lx := 0
	wx := 11 + fsi.NumberOfSecurityWords
	for ax := 0; ax < int(fsi.CurrentRange); ax++ {
		if wx == 28 {
			if leadItem1 == nil {
				break
			}
			lx++
			wx = 1
		}
		w := leadItems[lx][wx].GetW()
		link := w & 0_007777_777777
		if link > 0 {
			fsi.CycleInfo[ax] = &FileSetCycleInfo{
				ToBeCataloged:       w&0_200000_000000 != 0,
				ToBeDropped:         w&0_100000_000000 != 0,
				AbsoluteCycle:       absCycle,
				FileCycleIdentifier: FileCycleIdentifier(link),
			}
		}
	}
}
