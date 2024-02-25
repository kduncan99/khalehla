// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/exec"
	"khalehla/kexec/mfdMgr"
	"khalehla/pkg"
)

func (mgr *FacilitiesManager) AssignMassStorageFile() (FacResult, bool) {
	var facResult FacResult
	// TODO
	return facResult, false
}

func (mgr *FacilitiesManager) AssignTapeFile() (FacResult, bool) {
	var facResult FacResult
	// TODO
	return facResult, false
}

func (mgr *FacilitiesManager) CatalogMassStorageFile(
	rce *exec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	options pkg.Word36,
	mnemonic string,
	granularity kexec.Granularity,
	initialReserve uint64,
	maxGranules uint64,
	packNames []string,
) (FacResult, bool) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	var facResult FacResult
	// check options:
	//		B: save on checkpoint
	//		G: guarded file
	//		P: make the file public (not private)
	//		R: make the file read-only
	//		V: file will not be unloaded
	//		W: make the file write-only
	//		Z: run should not be held (probably only happens on removable when the pack is not mounted)
	// we ignore any inapplicable options, and there are no mutually exclusive ones
	saveOnCheckpoint := options&kexec.BOption != 0
	guardedFile := options&kexec.GOption != 0
	makePublic := options&kexec.POption != 0
	makeReadOnly := options&kexec.ROption != 0
	inhibitUnload := options&kexec.VOption != 0
	makeWriteOnly := options&kexec.WOption != 0
	doNotHold := options&kexec.ZOption != 0

	// See if there is already a fileset
	mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
	effectiveQual := rce.GetEffectiveQualifier(fileSpecification)
	fsInfo, leadItem0Address, mfdResult := mm.GetFileSetInfo(effectiveQual, fileSpecification.Filename)
	if mfdResult == mfdMgr.MFDNotFound {
		// TODO
	} else if mfdResult != mfdMgr.MFDSuccessful {
		// TODO oops
	}

	// If we do not have a mnemonic, determine what to use.
	// If the fileset does not already exist, use the system default...
	// Otherwise, use the highest f-cycle which is not to-be, or if none, then use the highest f-cycle of to-be
	// TODO

	// Are we fixed or removable? Should we build this into the equipment table?
	// If fixed, ensure pack name list is empty. If removable, ensure pack name list is *not* empty.
	// TODO

	// resolve list of potential disk units based on mnemonic, as well as word or sector addressable
	// TODO

	// ensure initial reserve <= max allocations (means words or granules, depending on word/sector addressable)
	// TODO

	// For removable, deal with pack names.
	// If fileset already exists, pack name list should match what is listed for the fileset.
	// Ensure each pack name is known and mounted... Do not wait for mount if Z option is set
	// TODO

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
