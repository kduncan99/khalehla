// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/hardware"
	"khalehla/kexec"
	"khalehla/kexec/mfdMgr"
	"khalehla/klog"
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
	rce *kexec.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	klog.LogTraceF("FacMgr", "AssignFile [%v] %v", rce.RunId, *fileSpecification)
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult = NewFacResult()
	resultCode = 0

	mutuallyExclusiveOptions := uint64(kexec.AOption | kexec.COption | kexec.POption | kexec.TOption)
	if !checkIllegalOptionCombination(rce, optionWord, mutuallyExclusiveOptions, facResult, sourceIsExecRequest) {
		klog.LogInfoF("FacMgr", "[%v] Illegal option combination %012o", rce.RunId, optionWord)
		return
	}

	var mnemonic string
	if len(operandFields) >= 2 {
		mnemonic = operandFields[1][0]
		if len(mnemonic) > 6 {
			klog.LogInfoF("FacMgr", "[%v] Mnemonic %v too long", rce.RunId, mnemonic)
			facResult.PostMessage(kexec.FacStatusAssignMnemonicTooLong, []string{mnemonic})
			resultCode |= 0_600000_000000
			return
		}
	}

	models, usage, ok := mgr.selectEquipmentModel(mnemonic, nil)
	if !ok {
		// This isn't going to work for us.
		klog.LogInfoF(rce.RunId, "[%v] Mnemonic %v not configured", rce.RunId, mnemonic)
		facResult.PostMessage(kexec.FacStatusMnemonicIsNotConfigured, []string{mnemonic})
		resultCode |= 0_600000_000000
		return
	}

	effectiveFSpec := mgr.resolveFileSpecification(rce, fileSpecification)

	// check for option - branch on A, C/P, T, or none of the previous
	optSubset := optionWord & (kexec.AOption | kexec.COption | kexec.POption | kexec.TOption)
	if optSubset == kexec.AOption {
		facResult, resultCode = mgr.assignCatalogedFile(rce, sourceIsExecRequest, effectiveFSpec, optionWord, operandFields)
	} else if optSubset == kexec.COption || optSubset == kexec.POption {
		facResult, resultCode = mgr.assignToBeCatalogedFile(rce, sourceIsExecRequest, effectiveFSpec, optionWord, operandFields)
	} else if optSubset == kexec.TOption {
		facResult, resultCode =
			mgr.assignTemporaryFile(rce, sourceIsExecRequest, effectiveFSpec, optionWord, operandFields, models, usage)
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

	klog.LogTraceF(rce.RunId, "AssignFile resultCode %012o", resultCode)
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
	rce *kexec.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	klog.LogTraceF("FacMgr", "CatalogFile [%v] %v", rce.RunId, *fileSpecification)
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// TODO @CAT f(0). is legal and is the same as @CAT f.

	facResult = NewFacResult()
	resultCode = 0

	effectiveFSpec := mgr.resolveFileSpecification(rce, fileSpecification)

	// See if there is already a fileset
	mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
	fsIdent, mfdResult := mm.GetFileSetIdentifier(effectiveFSpec.Qualifier, effectiveFSpec.Filename)
	if mfdResult == mfdMgr.MFDInternalError {
		return
	}

	var fsInfo *mfdMgr.FileSetInfo
	fsInfo, mfdResult = mm.GetFileSetInfo(fsIdent)

	// Resolve the mnemonic (or lack of one) to a short list of acceptable node models
	// and a guide as to the usage of the model (sector, word, tape).
	var mnemonic string
	if len(operandFields) >= 2 && len(operandFields[1]) >= 1 {
		mnemonic = operandFields[1][0]
	}

	if len(mnemonic) > 6 {
		klog.LogInfoF(rce.RunId, "Mnemonic %v too long", mnemonic)
		facResult.PostMessage(kexec.FacStatusAssignMnemonicTooLong, []string{mnemonic})
		resultCode |= 0_600000_000000
		return
	}

	models, usage, ok := mgr.selectEquipmentModel(mnemonic, fsInfo)
	if !ok {
		// This isn't going to work for us.
		klog.LogInfoF(rce.RunId, "Mnemonic %v not configured", mnemonic)
		facResult.PostMessage(kexec.FacStatusMnemonicIsNotConfigured, []string{mnemonic})
		resultCode |= 0_600000_000000
		return
	}

	// We now know whether we are word-addressable, sector-addressable, or tape.
	// We don't yet know whether we are fixed or removable.
	var fileType mfdMgr.FileType
	if fsInfo != nil {
		fileType = fsInfo.FileType
	} else {
		if models[0].DeviceType == hardware.NodeDeviceTape {
			fileType = mfdMgr.FileTypeTape
		} else {
			// fixed or removable?
			// if there's anything in field 2 (pack names) then we'll assume it's removable
			if len(operandFields) >= 3 && len(operandFields[2]) > 0 {
				fileType = mfdMgr.FileTypeRemovable
			} else {
				fileType = mfdMgr.FileTypeFixed
			}
		}
	}

	if fileType == mfdMgr.FileTypeFixed {
		facResult, resultCode = mgr.catalogFixedFile(
			mgr.exec, rce, effectiveFSpec, optionWord, operandFields, fsInfo, mnemonic, usage, sourceIsExecRequest)
	} else if fileType == mfdMgr.FileTypeRemovable {
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
	return
}

// CheckIllegalOptions compares the given options word to the allowed options word,
// producing a fac message for each option set in the given word which does not appear in the allowed word.
// Returns true if no such instances were found, else false
// If not ok and the source is an ER CSF$/ACSF$/CSI$, we post a contingency
// This function does not have a good place to live - we put it here because it is mostly used by
// facmgr, and also by csi which is a client of facmgr.
func CheckIllegalOptions(
	rce *kexec.RunControlEntry,
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

func (mgr *FacilitiesManager) FreeFile(
	rce *kexec.RunControlEntry,
	fileSpecification kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	klog.LogTraceF("FacMgr", "FreeFile [%v] %v", rce.RunId, fileSpecification)
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult = NewFacResult()
	resultCode = 0

	// TODO implement FreeFile()

	if resultCode&0_400000_000000 == 0 {
		facResult.PostMessage(kexec.FacStatusComplete, []string{"FREE"})
	}
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

	mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
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
	rce *kexec.RunControlEntry,
	internalName string,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	klog.LogTraceF("FacMgr", "Use [%v] %v", rce.RunId, *fileSpecification)
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

	rce.UseItems[strings.ToUpper(internalName)] = &kexec.UseItem{
		InternalFilename:  internalName,
		FileSpecification: fileSpecification,
		ReleaseFlag:       optionWord&kexec.IOption != 0,
	}

	facResult.PostMessage(kexec.FacStatusComplete, []string{"USE"})
	return
}
