// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/hardware"
	"khalehla/kexec"
	"khalehla/kexec/mfdMgr"
	"khalehla/logger"
	hardware2 "khalehla/old/hardware"
	kexec2 "khalehla/old/kexec"
	"khalehla/old/kexec/config"
	mfdMgr2 "khalehla/old/kexec/mfdMgr"

	"strings"
)

// TODO Need a Boot() routine which assigns fixed packs to the exec

// AssignFile is the front end which all code should invoke when asking for a file to be assigned
// (and possibly cataloged).
// The operand fields are:
//
//	[0][0][0] all:file specification string including qual, file, cycle, read key, and write key.
//	[1][0][1] disk:assignMnemonic / tape:assignMnemonic / absolute:'*'{deviceName}
//	[1][1] disk:reserve / tape:units
//	[1][2] disk:[ TRK | POS ] / tape:log
//	[1][3] disk:maxGranules / tape:noiseConstant
//	[1][4] disk:placement / tape:processor [ASCII | EBCDIC | FLDATA ]
//	[1][5] tape:tape [ EBCDIC | FLATA | ASCII | BCD ](tape)
//	[1][6] tape:format [ Q | 8 ]
//	[1][7] tape:data-converter (n/a)
//	[1][8] tape:block-numbering
//	[1][9] tape:data-compression
//	[1][10] tape:buffered-write
//	[1][11] tape:expanded-buffer
//	[2][x] disk:packIds / tape:reelNumbers / absolute:packName (only one entry)
//	[3][0] tape:expirationPeriod
//	[3][1] tape:mmSpec
//	[4]    tape:ringIndicator
//	[5][0] all:ACR-name (we don't do ACRs yet)
//	[6][0] tape:ctlPool
func (mgr *FacilitiesManager) AssignFile(
	rce *kexec2.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec2.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	logger.LogTraceF("FacMgr", "AssignFile [%v] %v", rce.RunId, fileSpecification.ToString())
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult = NewFacResult()
	resultCode = 0

	mutuallyExclusiveOptions := uint64(kexec.AOption | kexec.COption | kexec.POption | kexec.TOption)
	if !checkIllegalOptionCombination(rce, optionWord, mutuallyExclusiveOptions, facResult, sourceIsExecRequest) {
		logger.LogInfoF("FacMgr", "[%v] Illegal option combination %012o", rce.RunId, optionWord)
		return
	}

	var mnemonic string
	if len(operandFields) >= 2 {
		mnemonic = operandFields[1][0]
		if len(mnemonic) > 6 {
			logger.LogInfoF("FacMgr", "[%v] Mnemonic %v too long", rce.RunId, mnemonic)
			facResult.PostMessage(kexec.FacStatusAssignMnemonicTooLong, []string{mnemonic})
			resultCode |= 0_600000_000000
			return
		}
	}

	// TODO does this make more sense in the more-specific routines in facCore?
	var models []hardware2.NodeModel
	var usage config.EquipmentUsage
	if mnemonic != "" {
		var ok bool
		models, usage, ok = mgr.selectEquipmentModel(mnemonic, nil)
		if !ok {
			// This isn't going to work for us.
			logger.LogInfoF(rce.RunId, "[%v] Mnemonic %v not configured", rce.RunId, mnemonic)
			facResult.PostMessage(kexec.FacStatusMnemonicIsNotConfigured, []string{mnemonic})
			resultCode |= 0_600000_000000
			return
		}
	}

	effectiveFSpec := rce.ResolveFileSpecification(fileSpecification, true)

	// check for option - branch on A, C/P, T, or none of the previous
	aOpt := optionWord&kexec.AOption != 0
	cOpt := optionWord&kexec.COption != 0
	pOpt := optionWord&kexec.POption != 0
	tOpt := optionWord&kexec.TOption != 0
	if aOpt {
		resultCode |= mgr.assignCatalogedFile(rce, sourceIsExecRequest, effectiveFSpec, optionWord, operandFields, facResult)
	} else if cOpt || pOpt {
		resultCode |= mgr.assignToBeCatalogedFile(rce, sourceIsExecRequest, effectiveFSpec, optionWord, operandFields, facResult)
	} else if tOpt {
		resultCode |= mgr.assignTemporaryFile(rce, sourceIsExecRequest, effectiveFSpec, optionWord, operandFields, models, usage, facResult)
	} else {
		// TODO No options among A,C,P,T - figure out what to do based on the rce's facility items
		// "The options field, if left blank, defaults to the A option unless the file is already assigned,
		// in which case the options on the previous assign are used. If the file specified is not cataloged,
		// the blank options field defaults to a T option, thus creating a temporary file.
		// Any option not allowed on a statement causes the statement to be rejected."
	}

	if resultCode&0_400000_000000 == 0 {
		facResult.PostMessage(kexec.FacStatusComplete, []string{"ASG"})
	}

	mgr.exec.GetMFDManager().(*mfdMgr2.MFDManager).PurgeDirectory()
	logger.LogTraceF(rce.RunId, "AssignFile resultCode %012o", resultCode)
	return
}

// CatalogFile is the front end which all code should invoke when asking for a file to be cataloged.
// The values allowed in optionWord depend upon whether we are cataloging a disk or a tape file.
// The operand fields are:
//
//	[0][0] file specification string including qual, file, cycle, read key, and write key.
//	[1][0] assign mnemonic
//	[1][1] reserve or number-of-units (n/a for @CAT)
//	[1][2] [ TRK | POS ] or logical-control-unit (n/a for @CAT)
//	[1][3] max grans or noise-constant (1 to 63)
//	[1][4] blank or processor [ASCII | EBCDIC | FLDATA ]
//	[1][5] blank or tape [ EBCDIC | FLATA | ASCII | BCD ]
//	[1][6] blank or format [ Q | 8 ]
//	[1][7] blank or data-converter (n/a)
//	[1][8] blank or block-numbering
//	[1][9] blank or data-compression
//	[1][10] blank or buffered-write
//	[1][11] blank or expanded-buffer
//	[2][x] pack-ids or reel-numbers
//	[3][0] blank or expiration-period
//	[3][1] blank or mmspec
//	[4]    blank
//	[5][0] ACR-name (we don't do ACRs yet)
//	[6][0] blank or CTL-pool
func (mgr *FacilitiesManager) CatalogFile(
	rce *kexec2.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec2.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	logger.LogTraceF("FacMgr", "CatalogFile [%v] %v", rce.RunId, *fileSpecification)
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// TODO @CAT f(0). is legal and is the same as @CAT f.

	facResult = NewFacResult()
	resultCode = 0

	effectiveFSpec := rce.ResolveFileSpecification(fileSpecification, true)

	// See if there is already a fileset
	var fsInfo *mfdMgr2.FileSetInfo
	mm := mgr.exec.GetMFDManager().(*mfdMgr2.MFDManager)
	fsIdent, mfdResult := mm.GetFileSetIdentifier(effectiveFSpec.Qualifier, effectiveFSpec.Filename)
	if mfdResult == mfdMgr.MFDInternalError {
		logger.LogTrace("FacMgr", "CatalogFile early exit")
		return
	} else if mfdResult == mfdMgr.MFDSuccessful {
		fsInfo, mfdResult = mm.GetFileSetInfo(fsIdent)
		if mfdResult == mfdMgr.MFDInternalError {
			logger.LogTrace("FacMgr", "CatalogFile early exit")
			return
		}
	}

	// Resolve the mnemonic (or lack of one) to a short list of acceptable node models
	// and a guide as to the usage of the model (sector, word, tape).
	var mnemonic string
	if len(operandFields) >= 2 && len(operandFields[1]) >= 1 {
		mnemonic = operandFields[1][0]
	}

	if len(mnemonic) > 6 {
		logger.LogInfoF(rce.RunId, "Mnemonic %v too long", mnemonic)
		facResult.PostMessage(kexec.FacStatusAssignMnemonicTooLong, []string{mnemonic})
		resultCode |= 0_600000_000000
		return
	}

	models, usage, ok := mgr.selectEquipmentModel(mnemonic, fsInfo)
	if !ok {
		// This isn't going to work for us.
		logger.LogInfoF(rce.RunId, "Mnemonic %v not configured", mnemonic)
		facResult.PostMessage(kexec.FacStatusMnemonicIsNotConfigured, []string{mnemonic})
		resultCode |= 0_600000_000000
		return
	}

	// We now know whether we are word-addressable, sector-addressable, or tape.
	// We don't yet know whether we are fixed or removable.
	var fileType mfdMgr2.FileType
	if fsInfo != nil {
		fileType = fsInfo.FileType
	} else {
		if models[0].DeviceType == hardware2.NodeDeviceTape {
			fileType = mfdMgr2.FileTypeTape
		} else {
			// fixed or removable?
			// if there's anything in field 2 (pack names) then we'll assume it's removable
			if len(operandFields) >= 3 && len(operandFields[2]) > 0 {
				fileType = mfdMgr2.FileTypeRemovable
			} else {
				fileType = mfdMgr2.FileTypeFixed
			}
		}
	}

	if fileType == mfdMgr2.FileTypeFixed {
		facResult, resultCode = mgr.catalogFixedFile(
			mgr.exec, rce, effectiveFSpec, optionWord, operandFields, fsInfo, mnemonic, usage, sourceIsExecRequest)
	} else if fileType == mfdMgr2.FileTypeRemovable {
		facResult, resultCode = mgr.catalogRemovableFile(
			mgr.exec, rce, effectiveFSpec, optionWord, operandFields, fsInfo, mnemonic, usage, sourceIsExecRequest)
	} else { // fileType == kexec.FileTypeTape
		facResult, resultCode = mgr.catalogTapeFile(
			mgr.exec, rce, effectiveFSpec, optionWord, operandFields, fsInfo, mnemonic, usage, sourceIsExecRequest)
	}

	if resultCode&0_400000_000000 == 0 {
		facResult.PostMessage(kexec.FacStatusComplete, []string{"CAT"})
	}

	mm.PurgeDirectory()
	logger.LogTraceF(rce.RunId, "CatalogFile resultCode %012o", resultCode)
	return
}

// CheckIllegalOptions compares the given options word to the allowed options word,
// producing a fac message for each option set in the given word which does not appear in the allowed word.
// Returns true if no such instances were found, else false
// If not ok and the source is an ER CSF$/ACSF$/CSI$, we post a contingency
// This function does not have a good place to live - we put it here because it is mostly used by
// facmgr, and also by csi which is a client of facmgr.
func CheckIllegalOptions(
	rce *kexec2.RunControlEntry,
	givenOptions uint64,
	allowedOptions uint64,
	facResult *FacStatusResult,
	sourceIsExec bool,
) bool {
	bit := uint64(kexec.AOption)
	letter := 'A'
	ok := true

	for {
		if bit&givenOptions != 0 && bit&allowedOptions == 0 {
			param := string(letter)
			facResult.PostMessage(kexec.FacStatusIllegalOption, []string{param})
			ok = false
		}

		if bit == kexec.ZOption {
			break
		} else {
			letter++
			bit >>= 1
		}
	}

	if !ok {
		if sourceIsExec {
			rce.PostContingency(012, 04, 040)
		}
	}

	return ok
}

// FreeFile releases a file or a use item (or some combination thereof) from the current run.
// Behavior (assume cycles 1, 2, and 3 exist):
//
//	asg t. free t(0), t., t(-0), t(+0), t(3) all work
//	asg t(3). ONLY free t(3) works
//	asg t(2). ONLY free t(2) works
//	asg t(-1). free t(-1) and t(2) works
//
// Options
//
//	A - releases the internal use filename, but leaves the attached file assigned
//	B - same as A, but if no other use names are attached, the external file is free'd
//	D - deletes and frees the file if cataloged, frees it if temporary, prevents cataloging if C or U option assign
//	I - inhibits cataloging for C or U option assign
//	R - releases the file, but maintains all use names (C/U options cause catalog to take effect here)
//	S - releases the file from the run, but keeps the tape unit assigned
//	X - releases exclusive use if assigned A or X; otherwise releases the file
//
// Mass storage files to be cataloged:    Blank, A, B, D, I, R, X*
// Already cataloged mass storage files:  Blank, A, B, D, R, X
// Temporary mass storage files:          Blank, A, B, D, R, X*
// Tape files to be cataloged:            Blank, A, B, D, I, R, S
// Already cataloged tape files:          Blank, A, B, D, R, S
// Temporary tape files:                  Blank, A, B, D, R, S
// Tape devices:                          Blank, A, B, D, R, X*
// Internal file name:                    Blank, A, B, R, S, X*
// * Allowed, but has no effect
func (mgr *FacilitiesManager) FreeFile(
	rce *kexec2.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec2.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	logger.LogTraceF("FacMgr", "FreeFile [%v] %v", rce.RunId, fileSpecification)
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult = NewFacResult()
	resultCode = 0

	if !mgr.checkSubFields(operandFields, asgDiskFSIs) {
		facResult.PostMessage(kexec.FacStatusUndefinedFieldOrSubfield, nil)
		resultCode |= 0_600000_000000
		logger.LogTraceF("FacMgr", "FreeFile exit resultCode %012o", resultCode)
		return
	}

	validMask := uint64(kexec.AOption | kexec.BOption | kexec.DOption |
		kexec.IOption | kexec.ROption | kexec.SOption | kexec.XOption)
	if !CheckIllegalOptions(rce, optionWord, validMask, facResult, sourceIsExecRequest) {
		resultCode = 0_600000_000000
		return
	}

	effectiveSpec := fileSpecification
	aOpt := optionWord&kexec.AOption != 0
	bOpt := optionWord&kexec.BOption != 0

	// Did the caller specify a use name?
	useItem := rce.FindUseItem(fileSpecification)
	if useItem != nil {
		if aOpt || bOpt {
			// release the use item only
			rce.DeleteUseItem(fileSpecification.Filename)
			if bOpt {
				// If there is a fac item, and if it has no other use items, release that as well.
				// TODO releaseFacItem(useItem.FileSpecification) (do not resolve)
			}

			facResult.PostMessage(kexec.FacStatusComplete, []string{"FREE"})
			return
		}

		// chase use items to find an effective fileSpec then drop through
		effectiveSpec = rce.ResolveFileSpecification(effectiveSpec, true)
	}

	if aOpt || bOpt {
		logger.LogInfoF("FacMgr", "A or B option on non-internal name [%v]", fileSpecification.ToString())
		facResult.PostMessage(kexec.FacStatusInternalNameRequired, nil)
		resultCode |= 0_400000_000000
		return
	}

	// TODO releaseFacItem(effectiveSpec) (do not resolve)

	return
}

func (mgr *FacilitiesManager) GetTrackCounts() (
	msAccessible hardware.TrackCount,
	msAvailable hardware.TrackCount,
	mfdAccessible hardware.TrackCount,
	mfdAvailable hardware.TrackCount,
) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	msAccessible = 0
	msAvailable = 0
	mfdAccessible = 0
	mfdAvailable = 0

	mm := mgr.exec.GetMFDManager().(*mfdMgr2.MFDManager)
	for nodeId, diskAttr := range mgr.inventory.disks {
		stat := diskAttr.GetFacNodeStatus()
		if stat == kexec.FacNodeStatusUp || stat == kexec.FacNodeStatusSuspended {
			msAcc, msAvail, mfdAcc, mfdAvail := mm.GetTrackCountsForPack(nodeId)
			msAccessible += msAcc
			if stat == kexec.FacNodeStatusUp {
				msAvailable += msAvail
				mfdAccessible += mfdAcc
				mfdAvailable += mfdAvail
			}
		}
	}
	return
}

func (mgr *FacilitiesManager) Use(
	rce *kexec2.RunControlEntry,
	internalName string,
	fileSpecification *kexec2.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	logger.LogTraceF("FacMgr", "Use [%v] %v", rce.RunId, *fileSpecification)
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult = NewFacResult()
	resultCode = 0

	if optionWord != 0 && optionWord != kexec.IOption {
		facResult.PostMessage(kexec.FacStatusIOptionOnlyAllowed, nil)
		resultCode |= 0_400000_400000
		return
	}

	if !mgr.checkSubFields(operandFields, useFSIs) {
		facResult.PostMessage(kexec.FacStatusUndefinedFieldOrSubfield, nil)
		resultCode |= 0_600000_000000
		return
	}

	rce.UseItems[strings.ToUpper(internalName)] = &kexec2.UseItem{
		InternalFilename:  internalName,
		FileSpecification: fileSpecification,
		ReleaseFlag:       optionWord&kexec.IOption != 0,
	}

	facResult.PostMessage(kexec.FacStatusComplete, []string{"USE"})
	return
}
