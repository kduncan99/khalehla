// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

type FacStatusCategory uint

const (
	FacMsgInfo FacStatusCategory = iota
	FacMsgWarning
	FacMsgError
)

type FacStatusCode uint

const (
	FacStatusRunHeldForDevice                   = 000133
	FacStatusRunHeldForRemovable                = 000233
	FacStatusRunHeldForPack                     = 000333
	FacStatusRunHeldForReel                     = 000433
	FacStatusRunHeldForComLine                  = 000533
	FacStatusRunHeldForComGroup                 = 000633
	FacStatusRunHeldForMassStorageSpace         = 000733
	FacStatusRunHeldForTapeUnitAvailability     = 001033
	FacStatusRunHeldForExclusiveFileUseRelease  = 001133
	FacStatusRunHeldForNeedOfExclusiveUse       = 001233
	FacStatusRunHeldForDiskUnitAvailability     = 001333
	FacStatusRunHeldForRollback                 = 001433
	FacStatusRunHeldForFileCycleConflict        = 001533
	FacStatusRunHeldForDiskPackMount            = 001633
	FacStatusDeviceIsSelected                   = 001733
	FacStatusRunHeldForCacheControl             = 002033
	FacStatusRunHeldForDisketteUnitAvailability = 002133
	FacStatusRunHeldForDisketteMount            = 002233
	FacStatusComplete                           = 002333
	FacStatusMnemonicIsNotConfigured            = 0201033
	FacStatusIllegalOptionCombination           = 0201433
	FacStatusIllegalOption                      = 0201533
)

type FacStatusMessage struct {
	Category FacStatusCategory
	Code     FacStatusCode
	Message  string
}

func newMessage(category FacStatusCategory, code FacStatusCode, message string) *FacStatusMessage {
	return &FacStatusMessage{
		Category: category,
		Code:     code,
		Message:  message,
	}
}

var FacMessageTemplates = map[FacStatusCode]*FacStatusMessage{
	FacStatusRunHeldForDevice:                   newMessage(FacMsgInfo, FacStatusRunHeldForDevice, "Run %v held for unit for %v device assign for %v min."),
	FacStatusRunHeldForRemovable:                newMessage(FacMsgInfo, FacStatusRunHeldForRemovable, "Run %v held for %v for abs rem disk for %v min."),
	FacStatusRunHeldForPack:                     newMessage(FacMsgInfo, FacStatusRunHeldForPack, "Run %v held for pack availability for %v min."),
	FacStatusRunHeldForReel:                     newMessage(FacMsgInfo, FacStatusRunHeldForReel, "Run %v held for reel availability for %v min."),
	FacStatusRunHeldForComLine:                  newMessage(FacMsgInfo, FacStatusRunHeldForComLine, "Run %v held for com line availability for %v min."),
	FacStatusRunHeldForComGroup:                 newMessage(FacMsgInfo, FacStatusRunHeldForComGroup, "Run %v held for com group availability for %v min."),
	FacStatusRunHeldForMassStorageSpace:         newMessage(FacMsgInfo, FacStatusRunHeldForMassStorageSpace, "Run %v held for mass storage space for %v min."),
	FacStatusRunHeldForTapeUnitAvailability:     newMessage(FacMsgInfo, FacStatusRunHeldForTapeUnitAvailability, "Run %v held for tape unit availability for %v min."),
	FacStatusRunHeldForExclusiveFileUseRelease:  newMessage(FacMsgInfo, FacStatusRunHeldForExclusiveFileUseRelease, "Run %v held for exclusive file use release for %v min."),
	FacStatusRunHeldForNeedOfExclusiveUse:       newMessage(FacMsgInfo, FacStatusRunHeldForNeedOfExclusiveUse, "Run %v held for need of exclusive use for %v min."),
	FacStatusRunHeldForDiskUnitAvailability:     newMessage(FacMsgInfo, FacStatusRunHeldForDiskUnitAvailability, "Run %v held for disk unit availability for %v min."),
	FacStatusRunHeldForRollback:                 newMessage(FacMsgInfo, FacStatusRunHeldForRollback, "Run %v held for rollback of unloaded file for %v min."),
	FacStatusRunHeldForFileCycleConflict:        newMessage(FacMsgInfo, FacStatusRunHeldForFileCycleConflict, "Run %v held for file cycle conflict for %v min."),
	FacStatusRunHeldForDiskPackMount:            newMessage(FacMsgInfo, FacStatusRunHeldForDisketteMount, "Run %v held for disk pack to be mounted for %v min."),
	FacStatusDeviceIsSelected:                   newMessage(FacMsgInfo, FacStatusDeviceIsSelected, "%v is selected %v %v %v"),
	FacStatusRunHeldForCacheControl:             newMessage(FacMsgInfo, FacStatusRunHeldForCacheControl, "Run %v held for control of caching for %v min."),
	FacStatusRunHeldForDisketteUnitAvailability: newMessage(FacMsgInfo, FacStatusRunHeldForDisketteUnitAvailability, "Run %v held for diskette unit availability for %v min."),
	FacStatusRunHeldForDisketteMount:            newMessage(FacMsgInfo, FacStatusRunHeldForDisketteMount, "Run %v held for diskette to be mounted for %v min."),
	FacStatusComplete:                           newMessage(FacMsgInfo, FacStatusComplete, "%v complete"),
	FacStatusMnemonicIsNotConfigured:            newMessage(FacMsgError, FacStatusMnemonicIsNotConfigured, "%v is not a configured assign mnemonic."),
	FacStatusIllegalOptionCombination:           newMessage(FacMsgError, FacStatusIllegalOptionCombination, "Illegal option combination %v%v."),
	FacStatusIllegalOption:                      newMessage(FacMsgError, FacStatusIllegalOption, "Illegal option %v."),
}

type FacMessage struct {
	code   FacStatusCode
	values []string
}

type FacResult struct {
	Infos    []*FacMessage
	Warnings []*FacMessage
	Errors   []*FacMessage
}

func NewFacResult() *FacResult {
	return &FacResult{
		Infos:    make([]*FacMessage, 0),
		Warnings: make([]*FacMessage, 0),
		Errors:   make([]*FacMessage, 0),
	}
}

func (fr *FacResult) HasInformationalMessages() bool {
	return len(fr.Infos) > 0
}

func (fr *FacResult) HasWarningMessages() bool {
	return len(fr.Warnings) > 0
}

func (fr *FacResult) HasErrorMessages() bool {
	return len(fr.Errors) > 0
}

func (fr *FacResult) PostMessage(code FacStatusCode, values []string) {
	msg := &FacMessage{
		code:   code,
		values: values,
	}
	temp := FacMessageTemplates[code]
	switch temp.Category {
	case FacMsgInfo:
		fr.Infos = append(fr.Infos, msg)
	case FacMsgWarning:
		fr.Warnings = append(fr.Warnings, msg)
	case FacMsgError:
		fr.Errors = append(fr.Errors, msg)
	}
}
