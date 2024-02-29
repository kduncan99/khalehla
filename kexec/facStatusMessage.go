// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type FacStatusMessageTemplate struct {
	Category FacStatusCategory
	Code     FacStatusCode
	Template string
}

func newFacStatusMessageTemplate(category FacStatusCategory, code FacStatusCode, message string) *FacStatusMessageTemplate {
	return &FacStatusMessageTemplate{
		Category: category,
		Code:     code,
		Template: message,
	}
}

var FacStatusMessageTemplates = importTemplates()

func importErrorTemplate(dest map[FacStatusCode]*FacStatusMessageTemplate, code FacStatusCode, text string) {
	importTemplate(dest, newFacStatusMessageTemplate(FacMsgError, code, text))
}

func importInfoTemplate(dest map[FacStatusCode]*FacStatusMessageTemplate, code FacStatusCode, text string) {
	importTemplate(dest, newFacStatusMessageTemplate(FacMsgInfo, code, text))
}

func importWarningTemplate(dest map[FacStatusCode]*FacStatusMessageTemplate, code FacStatusCode, text string) {
	importTemplate(dest, newFacStatusMessageTemplate(FacMsgWarning, code, text))
}

func importTemplate(dest map[FacStatusCode]*FacStatusMessageTemplate, template *FacStatusMessageTemplate) {
	dest[template.Code] = template
}

func importTemplates() map[FacStatusCode]*FacStatusMessageTemplate {
	result := make(map[FacStatusCode]*FacStatusMessageTemplate)

	// in order by info/warning/error, and then alphabetically by code -- some of these are never used
	importInfoTemplate(result, FacStatusComplete, "%v complete")
	importInfoTemplate(result, FacStatusDeviceIsSelected, "%v is selected %v %v %v")
	importInfoTemplate(result, FacStatusRunHeldForCacheControl, "Run %v held for control of caching for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForComGroup, "Run %v held for com group availability for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForComLine, "Run %v held for com line availability for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForDevice, "Run %v held for unit for %v device assign for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForDisketteMount, "Run %v held for diskette to be mounted for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForDisketteUnitAvailability, "Run %v held for diskette unit availability for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForDiskPackMount, "Run %v held for disk pack to be mounted for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForDiskUnitAvailability, "Run %v held for disk unit availability for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForExclusiveFileUseRelease, "Run %v held for exclusive file use release for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForFileCycleConflict, "Run %v held for file cycle conflict for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForMassStorageSpace, "Run %v held for mass storage space for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForNeedOfExclusiveUse, "Run %v held for need of exclusive use for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForPack, "Run %v held for pack availability for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForReel, "Run %v held for reel availability for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForRemovable, "Run %v held for %v for abs rem disk for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForRollback, "Run %v held for rollback of unloaded file for %v min.")
	importInfoTemplate(result, FacStatusRunHeldForTapeUnitAvailability, "Run %v held for tape unit availability for %v min.")

	importErrorTemplate(result, FacStatusAssignMnemonicMustBeSpecifiedWithPackId, "Assign mnemonic must be specified with a packid.")
	importErrorTemplate(result, FacStatusAssignMnemonicMustBeWordAddressable, "Assign mnemonic must be word addressable.")
	importErrorTemplate(result, FacStatusAssignMnemonicTooLong, "Assign mnemonic cannot be longer than 6 characters.")
	importErrorTemplate(result, FacStatusDirectoryAndQualifierMayNotAppear, "Directory id and qualifier may not appear on image when R option is used.")
	importErrorTemplate(result, FacStatusDirectoryOrQualifierMustAppear, "Directory id or qualifier must appear on image.")
	importErrorTemplate(result, FacStatusFilenameIsRequired, "A filename is required on the image.")
	importErrorTemplate(result, FacStatusIllegalControlStatement, "Illegal control statement type submitted to ER CSI$.")
	importErrorTemplate(result, FacStatusIllegalOption, "Illegal option %v.")
	importErrorTemplate(result, FacStatusIllegalOptionCombination, "Illegal option combination %v%v.")
	importErrorTemplate(result, FacStatusIllegalValueForFCycle, "Illegal value specified for F-cycle.")
	importErrorTemplate(result, FacStatusMnemonicIsNotConfigured, "%v is not a configured assign mnemonic.")
	importErrorTemplate(result, FacStatusSyntaxErrorInImage, "Syntax error in image submitted to ER CSI$.")

	return result
}
