// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import "khalehla/kexec"

type FacStatusMessageTemplate struct {
	Category kexec.FacStatusCategory
	Code     kexec.FacStatusCode
	Template string
}

func newFacStatusMessageTemplate(category kexec.FacStatusCategory, code kexec.FacStatusCode, message string) *FacStatusMessageTemplate {
	return &FacStatusMessageTemplate{
		Category: category,
		Code:     code,
		Template: message,
	}
}

var FacStatusMessageTemplates = importTemplates()

func importErrorTemplate(dest map[kexec.FacStatusCode]*FacStatusMessageTemplate, code kexec.FacStatusCode, text string) {
	importTemplate(dest, newFacStatusMessageTemplate(kexec.FacMsgError, code, text))
}

func importInfoTemplate(dest map[kexec.FacStatusCode]*FacStatusMessageTemplate, code kexec.FacStatusCode, text string) {
	importTemplate(dest, newFacStatusMessageTemplate(kexec.FacMsgInfo, code, text))
}

func importWarningTemplate(dest map[kexec.FacStatusCode]*FacStatusMessageTemplate, code kexec.FacStatusCode, text string) {
	importTemplate(dest, newFacStatusMessageTemplate(kexec.FacMsgWarning, code, text))
}

func importTemplate(dest map[kexec.FacStatusCode]*FacStatusMessageTemplate, template *FacStatusMessageTemplate) {
	dest[template.Code] = template
}

func importTemplates() map[kexec.FacStatusCode]*FacStatusMessageTemplate {
	result := make(map[kexec.FacStatusCode]*FacStatusMessageTemplate)

	// in order by info/warning/error, and then alphabetically by code -- some of these are never used
	// info
	importInfoTemplate(result, kexec.FacStatusComplete, "%v complete")
	importInfoTemplate(result, kexec.FacStatusDeviceIsSelected, "%v is selected %v %v %v")
	importInfoTemplate(result, kexec.FacStatusRunHeldForCacheControl, "Run %v held for control of caching for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForComGroup, "Run %v held for com group availability for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForComLine, "Run %v held for com line availability for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForDevice, "Run %v held for unit for %v device assign for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForDisketteMount, "Run %v held for diskette to be mounted for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForDisketteUnitAvailability, "Run %v held for diskette unit availability for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForDiskPackMount, "Run %v held for disk pack to be mounted for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForDiskUnitAvailability, "Run %v held for disk unit availability for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForExclusiveFileUseRelease, "Run %v held for exclusive file use release for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForFileCycleConflict, "Run %v held for file cycle conflict for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForMassStorageSpace, "Run %v held for mass storage space for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForNeedOfExclusiveUse, "Run %v held for need of exclusive use for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForPack, "Run %v held for pack availability for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForReel, "Run %v held for reel availability for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForRemovable, "Run %v held for %v for abs rem disk for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForRollback, "Run %v held for rollback of unloaded file for %v min.")
	importInfoTemplate(result, kexec.FacStatusRunHeldForTapeUnitAvailability, "Run %v held for tape unit availability for %v min.")

	// warning
	importWarningTemplate(result, kexec.FacStatusFileAlreadyAssigned, "File is already assigned.")
	importWarningTemplate(result, kexec.FacStatusFilenameNotUnique, "Filename not unique.")

	// error
	importErrorTemplate(result, kexec.FacStatusAssignMnemonicMustBeSpecifiedWithPackId, "Assign mnemonic must be specified with a packid.")
	importErrorTemplate(result, kexec.FacStatusAssignMnemonicMustBeWordAddressable, "Assign mnemonic must be word addressable.")
	importErrorTemplate(result, kexec.FacStatusAssignMnemonicTooLong, "Assign mnemonic cannot be longer than 6 characters.")
	importErrorTemplate(result, kexec.FacStatusAttemptToChangeGenericType, "Attempt to change generic type of the file.")
	importErrorTemplate(result, kexec.FacStatusDirectoryAndQualifierMayNotAppear, "Directory id and qualifier may not appear on image when R option is used.")
	importErrorTemplate(result, kexec.FacStatusDirectoryOrQualifierMustAppear, "Directory id or qualifier must appear on image.")
	importErrorTemplate(result, kexec.FacStatusFileAlreadyCataloged, "File is already catalogued.")
	importErrorTemplate(result, kexec.FacStatusFileCycleOutOfRange, "File cycle out of range.")
	importErrorTemplate(result, kexec.FacStatusFileIsNotCataloged, "File is not catalogued.")
	importErrorTemplate(result, kexec.FacStatusFilenameIsRequired, "A filename is required on the image.")
	importErrorTemplate(result, kexec.FacStatusFileNotCatalogedWithReadKey, "File is not cataloged with a read key.")
	importErrorTemplate(result, kexec.FacStatusFileNotCatalogedWithWriteKey, "File is not cataloged with a write key.")
	importErrorTemplate(result, kexec.FacStatusIllegalControlStatement, "Illegal control statement type submitted to ER CSI$.")
	importErrorTemplate(result, kexec.FacStatusIllegalDroppingPrivateFile, "Creation of file would require illegal dropping of private file.")
	importErrorTemplate(result, kexec.FacStatusIllegalInitialReserve, "Illegal value specified for initial reserve.")
	importErrorTemplate(result, kexec.FacStatusIllegalMaxGranules, "Illegal value specified for maximum.")
	importErrorTemplate(result, kexec.FacStatusIllegalOption, "Illegal option %v.")
	importErrorTemplate(result, kexec.FacStatusIllegalOptionCombination, "Illegal option combination %v%v.")
	importErrorTemplate(result, kexec.FacStatusIllegalValueForFCycle, "Illegal value specified for F-cycle.")
	importErrorTemplate(result, kexec.FacStatusIllegalValueForGranularity, "Illegal value specified for granularity.")
	importErrorTemplate(result, kexec.FacStatusIncorrectReadKey, "Incorrect read key.")
	importErrorTemplate(result, kexec.FacStatusIncorrectWriteKey, "Incorrect write key.")
	importErrorTemplate(result, kexec.FacStatusInternalNameRequired, "Internal Name is required.")
	importErrorTemplate(result, kexec.FacStatusIOptionOnlyAllowed, "I option is the only legal option on USE.")
	importErrorTemplate(result, kexec.FacStatusMaximumIsLessThanInitialReserve, "Maximum is less than the initial reserve.")
	importErrorTemplate(result, kexec.FacStatusMnemonicIsNotConfigured, "%v is not a configured assign mnemonic.")
	importErrorTemplate(result, kexec.FacStatusPlacementFieldNotAllowed, "Placement field is not allowed with CAT.")
	importErrorTemplate(result, kexec.FacStatusReadWriteKeysNeeded, "Read and/or write keys are needed.")
	importErrorTemplate(result, kexec.FacStatusRelativeFCycleConflict, "Relative F-cycle conflict.")
	importErrorTemplate(result, kexec.FacStatusSyntaxErrorInImage, "Syntax error in image submitted to ER CSI$.")
	importErrorTemplate(result, kexec.FacStatusUndefinedFieldOrSubfield, "Image contains an undefined field or subfield.")

	return result
}
