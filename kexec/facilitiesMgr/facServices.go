// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/mfdMgr"
	"khalehla/pkg"
)

func (mgr *FacilitiesManager) AssignFile() (kexec.FacStatusResult, bool) {
	var facResult kexec.FacStatusResult
	// TODO
	return facResult, false
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
//	[2][0..] pack-ids or reel-numbers
//	[3][0] blank or expiration-period
//	[3][1] blank or mmspec
//	[4]    blank
//	[5][0] ACR-name
//	[6][0] blank or CTL-pool
func (mgr *FacilitiesManager) CatalogFile(
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult = kexec.NewFacResult()
	resultCode = 0

	// See if there is already a fileset
	mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
	effectiveQual := rce.GetEffectiveQualifier(fileSpecification)
	fsIdent, mfdResult := mm.GetFileSetIdentifier(effectiveQual, fileSpecification.Filename)
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
		mgr.fallOver(rce,
			fmt.Sprintf("Mnemonic %v too long", mnemonic),
			facResult,
			kexec.FacStatusAssignMnemonicTooLong,
			[]string{mnemonic})
		return
	}

	models, usage, ok := mgr.selectEquipmentModel(mnemonic, fsInfo)
	if !ok {
		// This isn't going to work for us.
		resultCode = 0_600000_000000
		mgr.fallOver(rce,
			fmt.Sprintf("Mnemonic %v not configured", mnemonic),
			facResult,
			kexec.FacStatusMnemonicIsNotConfigured,
			[]string{mnemonic})
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
		return mgr.catalogFixedFile(mgr.exec, rce, fileSpecification, optionWord, operandFields, fsInfo, mnemonic, usage)
	} else if fileType == mfdMgr.FileTypeRemovable {
		return mgr.catalogRemovableFile(mgr.exec, rce, fileSpecification, optionWord, operandFields, fsInfo, mnemonic, usage)
	} else { // fileType == kexec.FileTypeTape
		return mgr.catalogTapeFile(mgr.exec, rce, fileSpecification, optionWord, operandFields, fsInfo, mnemonic, usage)
	}
}

func (mgr *FacilitiesManager) FreeFile(
	rce *kexec.RunControlEntry,
	fileSpecification kexec.FileSpecification,
	options pkg.Word36,
) (kexec.FacStatusResult, bool) {
	var facResult kexec.FacStatusResult
	// TODO
	return facResult, false
}
