// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/config"
	"khalehla/kexec/mfdMgr"
	"khalehla/kexec/nodes"
	"log"
)

func (mgr *FacilitiesManager) catalogFixedFile(
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *kexec.MFDFileSetInfo,
	usage config.EquipmentUsage,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	//	For Mass Storage Files
	//		@CAT[,options] filename[,type/reserve/granule/maximum,pack-id-1/.../pack-id-n,,,ACR-name]
	//	maximum of 6 fields in argument
	//	options include
	//		B: save on checkpoint
	//		G: guarded file
	//		P: make the file public (not private)
	//		R: make the file read-only
	//		V: file will not be unloaded
	//		W: make the file write-only
	//		Z: run should not be held (probably only happens on removable when the pack is not mounted)
	//
	allowedOpts := uint64(kexec.BOption | kexec.GOption | kexec.POption |
		kexec.ROption | kexec.VOption | kexec.WOption | kexec.ZOption)
	if !rce.CheckIllegalOptions(optionWord, allowedOpts, facResult, rce.IsExec) {

	}

	// saveOnCheckpoint := optionWord&kexec.BOption != 0
	// guardedFile := optionWord&kexec.GOption != 0
	// publicFile := optionWord&kexec.POption != 0
	// readOnly := optionWord&kexec.ROption != 0
	// inhibitUnload := optionWord&kexec.VOption != 0
	// writeOnly := optionWord&kexec.WOption != 0
	// doNotHold := optionWord&kexec.ZOption != 0
	// wordAddressable := usage == config.EquipmentUsageWordAddressableMassStorage

	// TODO ensure nothing is there that we don't like

	// If removable, ensure the pack list is compatible with the files in the fileset (if there is a fileset)
	// Is it okay to just use the highest cycle?
	// TODO

	// ensure initial reserve <= max allocations (means words or granules, depending on word/sector addressable)
	// TODO

	// If we are removable ensure each pack name is known and mounted.
	// Do not wait for mount if Z option is set
	// TODO
	return nil, 0 // TODO
}

func (mgr *FacilitiesManager) catalogRemovableFile(
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *kexec.MFDFileSetInfo,
	usage config.EquipmentUsage,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	//	For Mass Storage Files
	//		@CAT[,options] filename[,type/reserve/granule/maximum,pack-id-1/.../pack-id-n,,,ACR-name]
	//	maximum of 6 fields in argument
	//	options include
	//		B: save on checkpoint
	//		G: guarded file
	//		P: make the file public (not private)
	//		R: make the file read-only
	//		V: file will not be unloaded
	//		W: make the file write-only
	//		Z: run should not be held (probably only happens on removable when the pack is not mounted)
	//

	if (usage != config.EquipmentUsageWordAddressableMassStorage) &&
		(usage != config.EquipmentUsageSectorAddressableMassStorage) {
		// oops
	}

	return nil, 0 // TODO
}

func (mgr *FacilitiesManager) catalogTapeFile(
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *kexec.MFDFileSetInfo,
	usage config.EquipmentUsage,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	//	For Tape Files
	//		@CAT,options filename,type[/units/log/noise/processor/tape/
	//			format/data-converter/block-numbering/data-compression/
	//			buffered-write/expanded-buffer,reel-1/reel-2/.../reel-n,
	//			expiration-period/mmspec,,ACR-name,CTL-pool]
	//	maximum of 7 fields in argument
	//	options include
	//		E: even parity (not supported)
	//		G: guarded file
	//		H: density selection (not supported)
	//		J: tape is to be unlabeled
	//		L: density selection (not supported)
	//		M: density selection (not supported)
	//		O: odd parity (supported but ignored)
	//		P: make the file public
	//		R: make the file read-only
	//		S: 6250 BPI (only for SCSI 9-track - future)
	//		V: 1600 BPI (only for SCSI 9-track - future)
	//		W: make the file write-only
	//		Z: run should not be held (probably only happens on removable when the pack is not mounted)

	if usage != config.EquipmentUsageTape {
		// oops
	}

	return nil, 0 // TODO
}

func (mgr *FacilitiesManager) fallOver(
	rce *kexec.RunControlEntry,
	logMessage string,
	facResult *kexec.FacStatusResult,
	fsCode kexec.FacStatusCode,
	messageParameters []string,
) {
	log.Printf("%v:%v", rce.RunId, logMessage)
	facResult.PostMessage(fsCode, messageParameters)
}

func (mgr *FacilitiesManager) fallOverWithContingency(
	rce *kexec.RunControlEntry,
	logMessage string,
	facResult *kexec.FacStatusResult,
	fsCode kexec.FacStatusCode,
	messageParameters []string,
	contingencyType kexec.ContingencyType,
	contingencyErrorType uint,
	contingencyErrorCode uint,
) {
	log.Printf("%v:%v", rce.RunId, logMessage)
	rce.PostContingency(contingencyType, contingencyErrorType, contingencyErrorCode)
	facResult.PostMessage(fsCode, messageParameters)
}

// selectEquipmentModel accepts an equipment mnemonic (likely from a control statement)
// and an optional MFDFileSetInfo struct, and returns a list of NodeModel structs
// representing the various equipment models which can be used to satisfy the mnemonic.
// If the mnemonic is an @ASG or @CAT for a file cycle of an existing file set,
// the corresponding MFDFileSetInfo struct must be specified.
// A false return indicates that the mnemonic is not found.
func (mgr *FacilitiesManager) selectEquipmentModel(
	mnemonic string,
	fileSetInfo *kexec.MFDFileSetInfo,
) ([]nodes.NodeModel, config.EquipmentUsage, bool) {

	effectiveMnemonic := mnemonic

	// If we do not have a given mnemonic but we *do* have a fileSetInfo...
	if len(effectiveMnemonic) == 0 && fileSetInfo != nil {
		// Use the equipment type from the highest absolute fcycle entry of a not to-be file cycle
		//	(an existing file cycle which is not to-be-cataloged or to-be-deleted)... if there is one.
		// Otherwise, use the equipment type from the highest absolute fcycle entry of a to-be file cycle
		for _, preventToBe := range []bool{true, false} {
			for _, fsi := range fileSetInfo.CycleInfo {
				if !preventToBe || (!fsi.ToBeCataloged && !fsi.ToBeDropped) {
					mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
					fileInfo, _, _ := mm.GetFileInfo(fileSetInfo.Qualifier, fileSetInfo.Filename, fsi.AbsoluteCycle)
					effectiveMnemonic = fileInfo.GetAssignMnemonic()
				}
			}
		}
	}

	// If we still do not have an effective mnemonic use the default sector-formatted mass storage mnemonic.
	if len(effectiveMnemonic) == 0 {
		effectiveMnemonic = "F"
	}

	// Now go look for the mnemonic in the configured equipment entry table.
	entry, ok := mgr.exec.GetConfiguration().EquipmentTable[mnemonic]
	if !ok {
		return nil, 0, false
	}

	models := make([]nodes.NodeModel, 0)
	usage := entry.Usage
	for _, modelName := range entry.SelectableEquipment {
		model, ok := nodes.NodeModelTable[modelName]
		if ok {
			models = append(models, model)
		}
	}
	return models, usage, true
}
