// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facItems

import (
	"io"
	"khalehla/hardware"
	"khalehla/kexec"
)

type InternalName struct {
	IOption      bool
	InternalName string
}

// IFacilitiesItem structs are store in the RCE for all assigned facilities
type IFacilitiesItem interface {
	Dump(dest io.Writer, indent string)
	GetAbsoluteCycle() uint
	GetFilename() string
	GetInternalNames() []InternalName
	GetMnemonic() string   // mnemonic from @ASG statement, if any
	GetOptionWord() uint64 // @ASG options mask
	GetQualifier() string
	GetRelativeCycle() int
	HasAbsoluteCycle() bool
	HasInternalNames() bool
	HasRelativeCycle() bool
	IsAbsoluteDevice() bool
	IsDisk() bool
	IsNameItem() bool // an entry specified by @USE, but not (yet) assigned
	IsTape() bool
	IsTemporary() bool
}

type IAbsoluteFacilityItem interface {
	GetUnitIds() hardware.NodeIdentifier
	IsHeldForUnitAvailability() bool
	IsUnitAssigned() bool
}

type ICatalogedFileFacilitiesItem interface {
	GetMainItem0Address() kexec.MFDRelativeAddress
	IsExclusive() bool // assigned to this run with X option
	IsHeldForExclusiveUseNeed() bool
	IsHeldForExclusiveUseRelease() bool
	IsHeldForRollback() bool
}

type IDiskFileFacilitiesItem interface {
	GetFileAllocationSet() *kexec.FileAllocationSet
	IsSectorAddressable() bool
	IsWordAddressable() bool
	TranslateTrackId(fileTrackId hardware.TrackId) (ldat kexec.LDATIndex, deviceTrackId hardware.TrackId)
}

type ITapeFacilitiesItem interface {
	GetReelIds() []string
	GetUnitIds() []hardware.NodeIdentifier
	IsHeldForReelIdAvailability() bool
	IsHeldForUnitAvailability() bool
	IsUnitAssigned() bool
}

func FacilityItemMatches(
	facItem IFacilitiesItem,
	qualifier string,
	filename string,
	absoluteCycle *uint,
	relativeCycle *int,
) bool {
	if qualifier == facItem.GetQualifier() && filename == facItem.GetFilename() {
		if absoluteCycle == nil && !facItem.HasAbsoluteCycle() && relativeCycle == nil && !facItem.HasRelativeCycle() {
			return true
		}

		if absoluteCycle != nil && facItem.HasAbsoluteCycle() && *absoluteCycle == facItem.GetAbsoluteCycle() {
			return true
		}

		if relativeCycle != nil && facItem.HasRelativeCycle() && *relativeCycle == facItem.GetRelativeCycle() {
			return true
		}
	}
	return false
}

func FacilityItemsMatch(facItem1 IFacilitiesItem, facItem2 IFacilitiesItem) bool {
	if facItem1.GetQualifier() == facItem2.GetQualifier() && facItem1.GetFilename() == facItem2.GetFilename() {
		if !facItem1.HasAbsoluteCycle() && !facItem1.HasRelativeCycle() &&
			!facItem2.HasAbsoluteCycle() && !facItem2.HasRelativeCycle() {
			return true
		}

		if facItem1.HasAbsoluteCycle() && facItem2.HasAbsoluteCycle() &&
			facItem1.GetAbsoluteCycle() == facItem2.GetAbsoluteCycle() {
			return true
		}

		if facItem1.HasRelativeCycle() && facItem2.HasRelativeCycle() &&
			facItem1.GetRelativeCycle() == facItem2.GetRelativeCycle() {
			return true
		}
	}

	return false
}

// TODO move all the following to individual files in the package

// -------------------------------------------------------------------------------------

type AbsoluteDiskFacilitiesItem struct {
	deviceName     string
	isUnitAssigned bool // if false, we are waiting on the requested device
	packId         string
	unitId         hardware.NodeIdentifier
}

// -------------------------------------------------------------------------------------

type AbsoluteTapeFacilitiesItem struct {
	deviceNames    []string
	isUnitAssigned bool // if false, we are waiting on the requested device(s)
	reelIds        []string
	unitIds        []hardware.NodeIdentifier
}

// -------------------------------------------------------------------------------------

type CatalogedFixedDiskFacilitiesItem struct {
	qualifier        string
	filename         string
	absoluteCycle    *uint
	relativeCycle    *int
	optionWord       uint64
	mnemonic         string
	mainItem0Address kexec.MFDRelativeAddress
	internalNames    []InternalName
}

func NewCatalogedFixedDiskFacilitiesItem(
	qualifier string,
	filename string,
	absoluteCycle *uint,
	relativeCycle *int,
	optionWord uint64,
	mnemonic string,
	mainItem0Address kexec.MFDRelativeAddress,
) *CatalogedFixedDiskFacilitiesItem {
	return &CatalogedFixedDiskFacilitiesItem{
		qualifier:        qualifier,
		filename:         filename,
		absoluteCycle:    absoluteCycle,
		relativeCycle:    relativeCycle,
		optionWord:       optionWord,
		mnemonic:         mnemonic,
		mainItem0Address: mainItem0Address,
		internalNames:    make([]InternalName, 0),
	}
}

func (fi *CatalogedFixedDiskFacilitiesItem) Dump(dest io.Writer, indent string) {
	// TODO
}

func (fi *CatalogedFixedDiskFacilitiesItem) GetAbsoluteCycle() uint {
	return *fi.absoluteCycle
}

func (fi *CatalogedFixedDiskFacilitiesItem) GetFilename() string {
	return fi.filename
}

func (fi *CatalogedFixedDiskFacilitiesItem) GetInternalNames() []InternalName {
	return fi.internalNames
}

func (fi *CatalogedFixedDiskFacilitiesItem) GetMnemonic() string {
	return fi.mnemonic
}

func (fi *CatalogedFixedDiskFacilitiesItem) GetOptionWord() uint64 {
	return fi.optionWord
}

func (fi *CatalogedFixedDiskFacilitiesItem) GetQualifier() string {
	return fi.qualifier
}

func (fi *CatalogedFixedDiskFacilitiesItem) GetRelativeCycle() int {
	return *fi.relativeCycle
}

func (fi *CatalogedFixedDiskFacilitiesItem) HasAbsoluteCycle() bool {
	return fi.absoluteCycle != nil
}

func (fi *CatalogedFixedDiskFacilitiesItem) HasInternalNames() bool {
	return len(fi.internalNames) > 0
}

func (fi *CatalogedFixedDiskFacilitiesItem) HasRelativeCycle() bool {
	return fi.relativeCycle != nil
}

func (fi *CatalogedFixedDiskFacilitiesItem) IsAbsoluteDevice() bool {
	return false
}

func (fi *CatalogedFixedDiskFacilitiesItem) IsDisk() bool {
	return false
}

func (fi *CatalogedFixedDiskFacilitiesItem) IsNameItem() bool {
	return false
}

func (fi *CatalogedFixedDiskFacilitiesItem) IsTape() bool {
	return true
}

func (fi *CatalogedFixedDiskFacilitiesItem) IsTemporary() bool {
	return true
}

// -------------------------------------------------------------------------------------

type CatalogedRemovableDiskFacilitiesItem struct {
	mainItem0Address kexec.MFDRelativeAddress
}

// -------------------------------------------------------------------------------------

type CatalogedTapeFacilitiesItem struct {
	mainItem0Address kexec.MFDRelativeAddress
}

// -------------------------------------------------------------------------------------

type NameFacilitiesItem struct {
}

// -------------------------------------------------------------------------------------

type TemporaryFixedDiskFileFacilitiesItem struct {
	fileAllocations kexec.FileAllocationSet
}

// -------------------------------------------------------------------------------------

type TemporaryRemovableDiskFileFacilitiesItem struct {
	fileAllocations kexec.FileAllocationSet
	packIds         []string
}

// -------------------------------------------------------------------------------------

type TemporaryTapeFileFacilitiesItem struct {
	qualifier     string
	filename      string
	absoluteCycle *uint
	relativeCycle *int
	optionWord    uint64
	mnemonic      string
	units         int
	reelIds       []string
	unitIds       []hardware.NodeIdentifier
	internalNames []InternalName
}

func NewTemporaryTapeFileFacilitiesItem(
	qualifier string,
	filename string,
	absoluteCycle *uint,
	relativeCycle *int,
	optionWord uint64,
	mnemonic string,
	units int,
	reelIds []string,
) *TemporaryTapeFileFacilitiesItem {
	return &TemporaryTapeFileFacilitiesItem{
		qualifier:     qualifier,
		filename:      filename,
		absoluteCycle: absoluteCycle,
		relativeCycle: relativeCycle,
		optionWord:    optionWord,
		mnemonic:      mnemonic,
		units:         units,
		reelIds:       reelIds,
		unitIds:       make([]hardware.NodeIdentifier, 0),
	}
}

func (fi *TemporaryTapeFileFacilitiesItem) Dump(dest io.Writer, indent string) {
	// TODO
}

func (fi *TemporaryTapeFileFacilitiesItem) GetAbsoluteCycle() uint {
	return *fi.absoluteCycle
}

func (fi *TemporaryTapeFileFacilitiesItem) GetFilename() string {
	return fi.filename
}

func (fi *TemporaryTapeFileFacilitiesItem) GetInternalNames() []InternalName {
	return fi.internalNames
}

func (fi *TemporaryTapeFileFacilitiesItem) GetMnemonic() string {
	return fi.mnemonic
}

func (fi *TemporaryTapeFileFacilitiesItem) GetOptionWord() uint64 {
	return fi.optionWord
}

func (fi *TemporaryTapeFileFacilitiesItem) GetQualifier() string {
	return fi.qualifier
}

func (fi *TemporaryTapeFileFacilitiesItem) GetRelativeCycle() int {
	return *fi.relativeCycle
}

func (fi *TemporaryTapeFileFacilitiesItem) HasAbsoluteCycle() bool {
	return fi.absoluteCycle != nil
}

func (fi *TemporaryTapeFileFacilitiesItem) HasInternalNames() bool {
	return len(fi.internalNames) > 0
}

func (fi *TemporaryTapeFileFacilitiesItem) HasRelativeCycle() bool {
	return fi.relativeCycle != nil
}

func (fi *TemporaryTapeFileFacilitiesItem) IsAbsoluteDevice() bool {
	return false
}

func (fi *TemporaryTapeFileFacilitiesItem) IsDisk() bool {
	return false
}

func (fi *TemporaryTapeFileFacilitiesItem) IsNameItem() bool {
	return false
}

func (fi *TemporaryTapeFileFacilitiesItem) IsTape() bool {
	return true
}

func (fi *TemporaryTapeFileFacilitiesItem) IsTemporary() bool {
	return true
}
