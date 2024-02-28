// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/config"
	"khalehla/kexec/exec"
	"khalehla/kexec/mfdMgr"
	"khalehla/pkg"
)

func (mgr *FacilitiesManager) AssignFile() (FacResult, bool) {
	var facResult FacResult
	// TODO
	return facResult, false
}

func (mgr *FacilitiesManager) CatalogFile(
	rce *exec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	options pkg.Word36,
	mnemonic string,
	granularity kexec.Granularity,
	initialReserve uint64,
	maxGranules uint64,
	packNames []string,
) (*FacResult, bool) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	facResult := NewFacResult()

	// See if there is already a fileset
	mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
	effectiveQual := rce.GetEffectiveQualifier(fileSpecification)
	fsInfo, leadItem0Address, mfdResult := mm.GetFileSetInfo(effectiveQual, fileSpecification.Filename)
	if mfdResult == mfdMgr.MFDNotFound {
		// TODO
	} else if mfdResult != mfdMgr.MFDSuccessful {
		// TODO oops
	}

	// Resolve the mnemonic (or lack of one) to a short list of acceptable node models
	// and a guide as to the usage of the model (sector, word, tape).
	models, usage, ok := mgr.selectEquipmentModel(mnemonic, fsInfo)
	if !ok {
		// This isn't going to work for us.
		facResult.PostMessage(FacStatusMnemonicIsNotConfigured, []string{mnemonic})
		return facResult, false
	}

	// Are we fixed, removable, or tape? If there is a fileset, we know the answer.
	// Otherwise, we have to figure it out.
	// If the pack name list is empty, assume fixed.
	// If it is *not* empty, ensure the packs are either all fixed, or all removable.
	// While we're here, make sure the usage we got from the equipment selection matches
	// the usage we've decided to use here.
	// TODO (not sure the algorithm above is quite right... keep in mind that we have usage from equipment selection)

	if usage == config.EquipmentUsageSectorAddressableMassStorage ||
		usage == config.EquipmentUsageWordAddressableMassStorage {
		// check options for disk
		//		B: save on checkpoint
		//		G: guarded file
		//		P: make the file public (not private)
		//		R: make the file read-only
		//		V: file will not be unloaded
		//		W: make the file write-only
		//		Z: run should not be held (probably only happens on removable when the pack is not mounted)
		// we ignore any inapplicable options, and there are no mutually exclusive ones
		// TODO ensure nothing isn't there that we don't like
		saveOnCheckpoint := options&kexec.BOption != 0
		guardedFile := options&kexec.GOption != 0
		makePublic := options&kexec.POption != 0
		makeReadOnly := options&kexec.ROption != 0
		inhibitUnload := options&kexec.VOption != 0
		makeWriteOnly := options&kexec.WOption != 0
		doNotHold := options&kexec.ZOption != 0

		// If removable, ensure the pack list is compatible with the files in the fileset (if there is a fileset)
		// Is it okay to just use the highest cycle?
		// TODO

		// ensure initial reserve <= max allocations (means words or granules, depending on word/sector addressable)
		// TODO

		// If we are removable ensure each pack name is known and mounted.
		// Do not wait for mount if Z option is set
		// TODO

	} else {
		// check options for tape
		// TODO

		// TODO do we need to validate the reel names? that they're not used elsewhere or something?
	}

	// Now deal with MFDManager
	// TODO

	return facResult, false
}

func (mgr *FacilitiesManager) CatalogTapeFile(
	rce *exec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	options pkg.Word36,
	mnemonic string,
	// TODO other parameters
) (FacResult, bool) {
	var facResult FacResult
	// TODO
	return facResult, false
}

func (mgr *FacilitiesManager) FreeFile(
	rce *exec.RunControlEntry,
	fileSpecification kexec.FileSpecification,
	options pkg.Word36,
) (FacResult, bool) {
	var facResult FacResult
	// TODO
	return facResult, false
}
