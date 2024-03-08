// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/mfdMgr"
	"log"
	"strings"
)

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
	rce kexec.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult = NewFacResult()
	resultCode = 0

	// TODO AssignFile() See ECL 7.69 @USE
	// TODO AssignFile() search internal names

	// TODO AssignFile() if not found, search external names

	// TODO AssignFile() for already assigned files, see ECL 4.3 for security access validation
	//   see ECL 7.2.4 for changing certain fields
	if resultCode&0_400000_000000 == 0 {
		facResult.PostMessage(kexec.FacStatusComplete, []string{"ASG"})
	}
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
	rce kexec.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult = NewFacResult()
	resultCode = 0

	// TODO Need to account for @USE name

	// See if there is already a fileset
	mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
	effectiveQualifier := getEffectiveQualifier(rce, fileSpecification)
	fsIdent, mfdResult := mm.GetFileSetIdentifier(effectiveQualifier, fileSpecification.Filename)
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
		log.Printf("%v:Mnemonic %v too long", rce.RunId, mnemonic)
		facResult.PostMessage(kexec.FacStatusAssignMnemonicTooLong, []string{mnemonic})
		resultCode |= 0_600000_000000
		return
	}

	models, usage, ok := mgr.selectEquipmentModel(mnemonic, fsInfo)
	if !ok {
		// This isn't going to work for us.
		log.Printf("%v:Mnemonic %v not configured", rce.RunId, mnemonic)
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
		if models[0].DeviceType == kexec.NodeDeviceTape {
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
			mgr.exec, rce, fileSpecification, optionWord, operandFields, fsInfo, mnemonic, usage, sourceIsExecRequest)
	} else if fileType == mfdMgr.FileTypeRemovable {
		facResult, resultCode = mgr.catalogRemovableFile(
			mgr.exec, rce, fileSpecification, optionWord, operandFields, fsInfo, mnemonic, usage, sourceIsExecRequest)
	} else { // fileType == kexec.FileTypeTape
		facResult, resultCode = mgr.catalogTapeFile(
			mgr.exec, rce, fileSpecification, optionWord, operandFields, fsInfo, mnemonic, usage, sourceIsExecRequest)
	}

	if resultCode&0_400000_000000 == 0 {
		facResult.PostMessage(kexec.FacStatusComplete, []string{"CAT"})
	}
	return
}

func (mgr *FacilitiesManager) FreeFile(
	rce kexec.RunControlEntry,
	fileSpecification kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
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

func (mgr *FacilitiesManager) Use(
	rce kexec.RunControlEntry,
	internalName string,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
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
