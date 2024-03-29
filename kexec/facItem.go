// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
	"io"
)

// IFacilitiesItem structs are store in the RCE for all assigned facilities
type IFacilitiesItem interface {

	// TODO does a fac item have an attached internal name?

	Dump(dest io.Writer, indent string)
	GetQualifier() string
	GetFilename() string
	GetEquipmentCode() uint
	GetAttributes() FacItemAttributes
	GetRelativeCycle() int  // zero if none specified
	GetAbsoluteCycle() uint // zero if none specified
	IsDisk() bool
	IsTape() bool
}

/*
All facility items:
+00,W   internal file Name - Fieldata LJSF
+01,W   (internal file Name cont)
+02,W   file Name - Fieldata LJSF
+03,W   (file Name cont)
+04,W   qualifier - Fieldata LJSF
+05,W   (qualifier cont)
+06,S1  equipment code
         000 file has not been assigned (@USE exists, but @ASG has not been done)
         015 9-track tape
         016 virtual tape handler
         017 cartridge tape, DVD tape
         024 word-addressable mass storage
         036 sector-addressable mass storage
         077 arbitrary device
+07,S1  attributes
+07,b10:b35 @ASG options

Unit record and non-standard peripherals
Currently we do not recognize nor support non-standard / unit record peripherals,
so kexec will not consider them in facilities code.
+07,S1  attributes
         040 tape labeling is supported
         020 file is temporary
         010 internal Name is a use Name


*/

type FacItemAttributes struct {
	TapeLabelingIsSupported   bool
	FileIsTemporary           bool
	InternalFileNameIsUseName bool
	IsLargeFile               bool
}

type FacItemFileMode struct {
	IsExclusivelyAssigned bool
	IsReadKeyNeeded       bool
	IsWriteKeyNeeded      bool
	IsWriteInhibited      bool
	IsReadInhibited       bool
	IsWordAddressable     bool
}

// -----------------------------------------------------------------------

/*
Sector-formatted mass storage
+06,S2  file mode
         040 exclusively assigned
         020 read key needed
         010 write key needed
         004 file is write inhibited
         002 file is read inhibited
         001 word-addressable (always clear)
+06,S3  granularity
         zero -> track, nonzero -> position
+06,S4  relative file-cycle
+06,T3  absolute file-cycle
+07,S1  attributes
         020 file is temporary
         010 internal Name is a use Name
         004 shared file
         002 large file
+010,H1 initial granule count (initial reserve)
+010,H2 max granule count
+011,H1 highest track referenced
+011,H2 highest granule assigned
+012,S4 total pack count if removable (63 -> 63 or greater)
+012,S5 equipment code - same as +06,S1
+012,S6 subcode - zero
*/

type SectorAddressableFacilityItem struct {
	Qualifier     string
	Filename      string
	EquipmentCode uint
	RelativeCycle int
	AbsoluteCycle uint
	Attributes    FacItemAttributes
	FileMode      FacItemFileMode
	// TODO I think the following should not be here - they're in the MFD already, and only needed for ER FITEM$
	Granularity            Granularity
	InitialReserve         uint64 // in granules
	MaxGranules            uint64
	HighestTrackReferenced uint64
	HighestGranuleAssigned uint64
	TotalPackCount         uint // if removable
}

func (fi *SectorAddressableFacilityItem) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vSector %v*%v abs=%v rel=%v equip=%03o\n",
		indent, fi.Qualifier, fi.Filename, fi.AbsoluteCycle, fi.RelativeCycle, fi.EquipmentCode)

	_, _ = fmt.Fprintf(dest, "%v  Attributes: tpLbl:%v temp:%v intUse:%v large:%v\n",
		indent,
		fi.Attributes.TapeLabelingIsSupported,
		fi.Attributes.FileIsTemporary,
		fi.Attributes.InternalFileNameIsUseName,
		fi.Attributes.IsLargeFile)

	_, _ = fmt.Fprintf(dest, "%v  Mode: xAsg:%v rKey:%v wKey:%v wInh:%v rInh:%v wad:%v\n",
		indent,
		fi.FileMode.IsExclusivelyAssigned,
		fi.FileMode.IsReadKeyNeeded,
		fi.FileMode.IsWriteKeyNeeded,
		fi.FileMode.IsWriteInhibited,
		fi.FileMode.IsReadInhibited,
		fi.FileMode.IsWordAddressable)

	_, _ = fmt.Fprintf(dest, "%v  gran:%v rsv:%v max:%v hTrk:%v hGrn:%v packs:%v\n",
		indent, fi.Granularity, fi.InitialReserve, fi.MaxGranules,
		fi.HighestTrackReferenced, fi.HighestGranuleAssigned, fi.TotalPackCount)
}

func (fi *SectorAddressableFacilityItem) GetQualifier() string {
	return fi.Qualifier
}

func (fi *SectorAddressableFacilityItem) GetFilename() string {
	return fi.Filename
}

func (fi *SectorAddressableFacilityItem) GetEquipmentCode() uint {
	return fi.EquipmentCode
}

func (fi *SectorAddressableFacilityItem) GetAttributes() FacItemAttributes {
	return fi.Attributes
}

func (fi *SectorAddressableFacilityItem) GetRelativeCycle() int {
	return fi.RelativeCycle
}

func (fi *SectorAddressableFacilityItem) GetAbsoluteCycle() uint {
	return fi.AbsoluteCycle
}

func (fi *SectorAddressableFacilityItem) IsDisk() bool {
	return true
}

func (fi *SectorAddressableFacilityItem) IsTape() bool {
	return false
}

// -----------------------------------------------------------------------

/*
Word addressable
+06,S2  file mode
         040 exclusively assigned
         020 read key needed
         010 write key needed
         004 write inhibited
         002 read inhibited
         001 word-addressable (always set)
+06,S3  granularity
         zero -> track, nonzero -> position
+06,S4  relative file-cycle
+06,T3  absolute file-cycle
+07,S1  attributes
         020 file is temporary
         010 internal Name is a use Name
         004 shared file
+010,W  length of file in words
+011,W  maximum file length in words
+012,S4 total pack count if removable (63 -> 63 or greater)
+012,S5 equipment code - same as +06,S1
+012,S6 subcode - zero
*/

type WordAddressableFacilityItem struct {
	Qualifier     string
	Filename      string
	EquipmentCode uint
	RelativeCycle int
	AbsoluteCycle uint
	Attributes    FacItemAttributes
	FileMode      FacItemFileMode
	// TODO I think the following should not be here - they're in the MFD already, and only needed for ER FITEM$
	Granularity       Granularity
	LengthOfFile      uint64 // in words
	MaximumFileLength uint64 // in words
	TotalPackCount    uint   // if removable
}

func (fi *WordAddressableFacilityItem) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vSector %v*%v abs=%v rel=%v equip=%03o\n",
		indent, fi.Qualifier, fi.Filename, fi.AbsoluteCycle, fi.RelativeCycle, fi.EquipmentCode)

	_, _ = fmt.Fprintf(dest, "%v  Attributes: [tpLbl:%v temp:%v intUse:%v large:%v\n",
		indent,
		fi.Attributes.TapeLabelingIsSupported,
		fi.Attributes.FileIsTemporary,
		fi.Attributes.InternalFileNameIsUseName,
		fi.Attributes.IsLargeFile)

	_, _ = fmt.Fprintf(dest, "%v  Mode: xAsg:%v rKey:%v wKey:%v wInh:%v rInh:%v wad:%v\n",
		indent,
		fi.FileMode.IsExclusivelyAssigned,
		fi.FileMode.IsReadKeyNeeded,
		fi.FileMode.IsWriteKeyNeeded,
		fi.FileMode.IsWriteInhibited,
		fi.FileMode.IsReadInhibited,
		fi.FileMode.IsWordAddressable)

	_, _ = fmt.Fprintf(dest, "%v  gran:%v len:%v max:%v packs:%v\n",
		indent, fi.Granularity, fi.LengthOfFile, fi.MaximumFileLength, fi.TotalPackCount)
}

func (fi *WordAddressableFacilityItem) GetQualifier() string {
	return fi.Qualifier
}

func (fi *WordAddressableFacilityItem) GetFilename() string {
	return fi.Filename
}

func (fi *WordAddressableFacilityItem) GetEquipmentCode() uint {
	return fi.EquipmentCode
}

func (fi *WordAddressableFacilityItem) GetAttributes() FacItemAttributes {
	return fi.Attributes
}

func (fi *WordAddressableFacilityItem) GetRelativeCycle() int {
	return fi.RelativeCycle
}

func (fi *WordAddressableFacilityItem) GetAbsoluteCycle() uint {
	return fi.AbsoluteCycle
}

func (fi *WordAddressableFacilityItem) IsDisk() bool {
	return true
}

func (fi *WordAddressableFacilityItem) IsTape() bool {
	return false
}

// -----------------------------------------------------------------------

/*
Magnetic tape peripherals
+06,S2  file mode
         040 exclusively assigned
         020 read key needed
         010 write key needed
         004 file is write inhibited
         002 file is read inhibited
+06,S3  unit count (I presume, the docs are not helpful)
		number of units assigned (0?, 1, 2)
+06,S4  relative file-cycle
+06,T3  absolute file-cycle
+07,S1  attributes
         040 tape labeling is supported
         020 file is temporary
         010 internal Name is a use Name
         004 file is a shared file
+010,S1 total reel count
+010,S2 logical channel
+010,S3 noise constant
+012,T1 expiration period
+012,S3 reel index
+012,S4 files extended
+012,T3 blocks extended
+013,W  current reel number
+014,W  next reel number
*/

type TapeFacilityItem struct {
	Qualifier     string
	Filename      string
	EquipmentCode uint
	RelativeCycle int
	AbsoluteCycle uint
	Attributes    FacItemAttributes
	// TODO I think the following should not be here - they're elsewhere already, and only needed for ER FITEM$
	TotalReelCount    uint
	LogicalChannel    uint
	NoiseConstant     uint
	ExpirationPeriod  uint
	ReelIndex         uint
	FilesExtended     uint
	BlocksExtended    uint
	CurrentReelNumber string
	NextReelNumber    string
}

func (fi *TapeFacilityItem) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vTape %v*%v abs=%v rel=%v equip=%03o\n",
		indent, fi.Qualifier, fi.Filename, fi.AbsoluteCycle, fi.RelativeCycle, fi.EquipmentCode)

	_, _ = fmt.Fprintf(dest, "%v  Attributes: tpLbl:%v temp:%v intUse:%v large:%v\n",
		indent,
		fi.Attributes.TapeLabelingIsSupported,
		fi.Attributes.FileIsTemporary,
		fi.Attributes.InternalFileNameIsUseName,
		fi.Attributes.IsLargeFile)
}

func (fi *TapeFacilityItem) GetQualifier() string {
	return fi.Qualifier
}

func (fi *TapeFacilityItem) GetFilename() string {
	return fi.Filename
}

func (fi *TapeFacilityItem) GetEquipmentCode() uint {
	return fi.EquipmentCode
}

func (fi *TapeFacilityItem) GetAttributes() FacItemAttributes {
	return fi.Attributes
}

func (fi *TapeFacilityItem) GetRelativeCycle() int {
	return fi.RelativeCycle
}

func (fi *TapeFacilityItem) GetAbsoluteCycle() uint {
	return fi.AbsoluteCycle
}

func (fi *TapeFacilityItem) IsDisk() bool {
	return true
}

func (fi *TapeFacilityItem) IsTape() bool {
	return false
}
