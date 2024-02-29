// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/exec"
	"khalehla/kexec/mfdMgr"
	"khalehla/kexec/nodeMgr"
	"khalehla/pkg"
)

func (mgr *FacilitiesManager) AssignFile() (FacStatusResult, bool) {
	var facResult FacStatusResult
	// TODO
	return facResult, false
}

func (mgr *FacilitiesManager) CatalogFile(
	rce *exec.RunControlEntry,
	fileSpecification *FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult = NewFacResult()
	resultCode = 0

	// See if there is already a fileset
	mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
	effectiveQual := rce.GetEffectiveQualifier(fileSpecification)
	fsInfo, _, mfdResult := mm.GetFileSetInfo(effectiveQual, fileSpecification.Filename)
	if mfdResult == mfdMgr.MFDInternalError {
		return
	}

	// Resolve the mnemonic (or lack of one) to a short list of acceptable node models
	// and a guide as to the usage of the model (sector, word, tape).
	var mnemonic string
	if len(operandFields) >= 2 && len(operandFields[1]) >= 1 {
		mnemonic = operandFields[1][0]
	}
	models, usage, ok := mgr.selectEquipmentModel(mnemonic, fsInfo)
	if !ok {
		// This isn't going to work for us.
		facResult.PostMessage(FacStatusMnemonicIsNotConfigured, []string{mnemonic})
		return facResult, 0_600000_000000
	}

	// We now know whether we are word-addressable, sector-addressable, or tape.
	// We don't yet know whether we are fixed or removable, and the possibility exists
	// that there is still some conflict between what we know about usage from equipment type,
	// and what the caller is implying based on options, arguments, pack/reel names, etc.
	// If there is a fileset, the answer is already there.
	var fileType kexec.FileType
	if fsInfo != nil {
		fileType = fsInfo.FileType
	} else {
		if models[0].DeviceType == nodeMgr.NodeDeviceTape {
			fileType = kexec.FileTypeTape
		} else {
			// fixed or removable?
			// if there's anything in field 2 (pack names) then we'll assume it's removable
			if len(operandFields) >= 3 && len(operandFields[2]) > 0 {
				fileType = kexec.FileTypeRemovable
			} else {
				fileType = kexec.FileTypeFixed
			}
		}
	}

	if fileType == kexec.FileTypeFixed {
		return mgr.catalogFixedFile(rce, fileSpecification, optionWord, operandFields, fsInfo, usage)
	} else if fileType == kexec.FileTypeRemovable {
		return mgr.catalogRemovableFile(rce, fileSpecification, optionWord, operandFields, fsInfo, usage)
	} else { // fileType == kexec.FileTypeTape
		return mgr.catalogTapeFile(rce, fileSpecification, optionWord, operandFields, fsInfo, usage)
	}
}

func (mgr *FacilitiesManager) FreeFile(
	rce *exec.RunControlEntry,
	fileSpecification FileSpecification,
	options pkg.Word36,
) (FacStatusResult, bool) {
	var facResult FacStatusResult
	// TODO
	return facResult, false
}
