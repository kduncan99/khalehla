// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/kexec"
	"khalehla/kexec/config"
	"khalehla/kexec/mfdMgr"
	"khalehla/kexec/nodeMgr"
	"log"
	"strconv"
	"strings"
)

type fieldSubfieldIndex struct {
	fieldIndex    int
	subFieldIndex int
	allSubfields  bool
}

type fieldSubfieldIndices struct {
	content []fieldSubfieldIndex
}

func newFieldSubfieldIndices() *fieldSubfieldIndices {
	return &fieldSubfieldIndices{
		content: make([]fieldSubfieldIndex, 0),
	}
}

func (fsi *fieldSubfieldIndices) add(fieldIndex int, subfieldIndex int) *fieldSubfieldIndices {
	index := fieldSubfieldIndex{
		fieldIndex:    fieldIndex,
		subFieldIndex: subfieldIndex,
	}
	fsi.content = append(fsi.content, index)
	return fsi
}

func (fsi *fieldSubfieldIndices) addAll(fieldIndex int) *fieldSubfieldIndices {
	index := fieldSubfieldIndex{
		fieldIndex:   fieldIndex,
		allSubfields: true,
	}
	fsi.content = append(fsi.content, index)
	return fsi
}

func (fsi *fieldSubfieldIndices) contains(fieldIndex int, subfieldIndex int) bool {
	for _, fsx := range fsi.content {
		if fieldIndex == fsx.fieldIndex && subfieldIndex == fsx.subFieldIndex {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------------

var catFixedFSIs = newFieldSubfieldIndices().
	add(0, 0).
	add(1, 0).
	add(1, 1).
	add(1, 2).
	add(1, 3)

var catRemovableFSIs = newFieldSubfieldIndices().
	add(0, 0).
	add(1, 0).
	add(1, 1).
	add(1, 2).
	add(1, 3).
	addAll(2)

// -----------------------------------------------------------------------------

// checkSubFields
// Checks the user-provided operation fields against a list of accepted field/subfield combinations
// to see whether the user provided a subfield which is not acceptable.
// Returns true if all is well, else false
func (mgr *FacilitiesManager) checkSubFields(operandFields [][]string, accepted *fieldSubfieldIndices) bool {
	for fx := 0; fx < len(operandFields); fx++ {
		for fy := 0; fy < len(operandFields[fx]); fy++ {
			if len(operandFields[fx][fy]) > 0 && !accepted.contains(fx, fy) {
				return false
			}
		}
	}
	return true
}

// getField
// Retrieves the field indicated by the given field index as an array of strings.
// If the field was not specified, we return an empty array.
func (mgr *FacilitiesManager) getField(operandFields [][]string, fieldIndex int) []string {
	if fieldIndex < len(operandFields) {
		return operandFields[fieldIndex]
	} else {
		return []string{}
	}
}

// getSubField
// Retrieves the subfield indicated by the given field and subfield indicies.
// If the subfield was not specified, we return a blank string.
func (mgr *FacilitiesManager) getSubField(operandFields [][]string, fieldIndex int, subfieldIndex int) string {
	if fieldIndex < len(operandFields) && subfieldIndex < len(operandFields[fieldIndex]) {
		return operandFields[fieldIndex][subfieldIndex]
	} else {
		return ""
	}
}

// -----------------------------------------------------------------------------

func (mgr *FacilitiesManager) assignFixedFile(
	exec kexec.IExec,
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *mfdMgr.FileSetInfo,
	usage config.EquipmentUsage,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	return nil, 0 // TODO
}

func (mgr *FacilitiesManager) assignRemovableFile(
	exec kexec.IExec,
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *mfdMgr.FileSetInfo,
	usage config.EquipmentUsage,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	return nil, 0 // TODO
}

func (mgr *FacilitiesManager) assignTapeFile(
	exec kexec.IExec,
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *mfdMgr.FileSetInfo,
	usage config.EquipmentUsage,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	return nil, 0 // TODO
}

// catalogFixedFile takes the various inputs, validates them, and then invokes
// mfd services to create the appropriate file cycle (and file set, if necessary).
// Caller should immediately check whether the exec has stopped upon return.
func (mgr *FacilitiesManager) catalogFixedFile(
	exec kexec.IExec,
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *mfdMgr.FileSetInfo,
	usage config.EquipmentUsage,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	//	For Mass Storage Files
	//		@CAT[,options] filename[,type/reserve/granule/maximum,pack-id-1/.../pack-id-n,,,ACR-name]
	//	options include
	//		B: save on checkpoint
	//		G: guarded file
	//		P: make the file public (not private)
	//		R: make the file read-only
	//		V: file will not be unloaded
	//		W: make the file write-only
	//		Z: run should not be held (probably only happens on removable when the pack is not mounted)
	resultCode = 0

	allowedOpts := uint64(kexec.BOption | kexec.GOption | kexec.POption |
		kexec.ROption | kexec.VOption | kexec.WOption | kexec.ZOption)
	if !rce.CheckIllegalOptions(optionWord, allowedOpts, facResult, rce.IsExec) {
		resultCode |= 0_600000_000000
	}

	if !mgr.checkSubFields(operandFields, catFixedFSIs) {
		if len(mgr.getSubField(operandFields, 1, 4)) > 0 {
			facResult.PostMessage(kexec.FacStatusPlacementFieldNotAllowed, nil)
		}
		facResult.PostMessage(kexec.FacStatusUndefinedFieldOrSubfield, nil)
		resultCode |= 0_600000_000000
	}

	saveOnCheckpoint := optionWord&kexec.BOption != 0
	guardedFile := optionWord&kexec.GOption != 0
	publicFile := optionWord&kexec.POption != 0
	readOnly := optionWord&kexec.ROption != 0
	inhibitUnload := optionWord&kexec.VOption != 0
	writeOnly := optionWord&kexec.WOption != 0
	doNotHold := optionWord&kexec.ZOption != 0
	wordAddressable := usage == config.EquipmentUsageWordAddressableMassStorage

	// ensure initial reserve <= max allocations (means words or granules, depending on word/sector addressable)
	initStr := mgr.getSubField(operandFields, 1, 1)
	granStr := strings.ToUpper(mgr.getSubField(operandFields, 1, 2))
	maxStr := mgr.getSubField(operandFields, 1, 3)

	var granularity kexec.Granularity
	if len(granStr) == 0 || granStr == "TRK" {
		granularity = kexec.TrackGranularity
	} else if granStr == "POS" {
		granularity = kexec.PositionGranularity
	} else {
		facResult.PostMessage(kexec.FacStatusIllegalValueForGranularity, nil)
		resultCode |= 0_600000_000000
	}

	var initReserve uint64
	if len(initStr) > 12 {
		facResult.PostMessage(kexec.FacStatusIllegalInitialReserve, nil)
		resultCode |= 0_600000_000000
	} else if len(initStr) > 0 {
		initReserve, err := strconv.Atoi(initStr)
		if err != nil || initReserve < 0 {
			facResult.PostMessage(kexec.FacStatusIllegalInitialReserve, nil)
			resultCode |= 0_600000_000000
		}
	}

	maxGranules := exec.GetConfiguration().MaxGranules
	if len(maxStr) > 12 {
		facResult.PostMessage(kexec.FacStatusIllegalMaxGranules, nil)
		resultCode |= 0_600000_000000
	} else if len(maxStr) > 0 {
		iMaxGran, err := strconv.Atoi(maxStr)
		maxGranules = uint64(iMaxGran)
		if err != nil || maxGranules < 0 || maxGranules > 262143 {
			facResult.PostMessage(kexec.FacStatusIllegalMaxGranules, nil)
			resultCode |= 0_600000_000000
		} else if maxGranules < initReserve {
			facResult.PostMessage(kexec.FacStatusMaximumIsLessThanInitialReserve, nil)
			resultCode |= 0_600000_000000
		}
	}

	if resultCode&0_400000_000000 != 0 {
		return
	}

	// If there isn't an existing fileset, create one.
	// Otherwise, do sanity checking on the requested file cycle.
	mm := exec.GetMFDManager().(*mfdMgr.MFDManager)
	if fileSetInfo == nil {
		fileSetInfo = mfdMgr.NewFileSetInfo(
			fileSpecification.Qualifier,
			fileSpecification.Filename,
			rce.ProjectId,
			fileSpecification.ReadKey,
			fileSpecification.WriteKey,
			mfdMgr.FileTypeFixed)
		_, result := mm.CreateFileSet(fileSetInfo)
		if result == mfdMgr.MFDInternalError {
			return
		} else if result != mfdMgr.MFDSuccessful {
			log.Printf("FacMgr:MFD failed to create file set")
			exec.Stop(kexec.StopFacilitiesComplex)
			return
		}
	} else {
		gaveReadKey := len(fileSpecification.ReadKey) > 0
		gaveWriteKey := len(fileSpecification.WriteKey) > 0
		hasReadKey := len(fileSetInfo.ReadKey) > 0
		hasWriteKey := len(fileSetInfo.WriteKey) > 0
		if hasReadKey {
			if !gaveReadKey {
				// TODO E:246433 Read and/or write keys are needed.
				resultCode |= 0_600000_000000
			} else if fileSetInfo.ReadKey != fileSpecification.ReadKey {
				// TODO E:253333 Incorrect read key.
				resultCode |= 0_401000_000000
				// TODO if ER CSF$, abort the run
			}
		} else {
			if gaveReadKey {
				// TODO E:253433 File is not cataloged with read key.
				resultCode |= 0_400040_000000
				// TODO if ER CSF$, abort the run
			}
		}

		if hasWriteKey {
			if !gaveWriteKey {
				// TODO E:246433 Read and/or write keys are needed.
				resultCode |= 0_600000_000000
			} else if fileSetInfo.WriteKey != fileSpecification.WriteKey {
				// TODO E:256633 Incorrect write key.
				resultCode |= 0_400400_000000
				// TODO if ER CSF$, abort the run
			}
		} else {
			if gaveWriteKey {
				// E:256733 File is not cataloged with write key.
				resultCode |= 0_400020_000000
				// TODO if ER CSF$, abort the run
			}
		}

		if resultCode&0_400000_000000 != 0 {
			return
		}

		/*
			E:242533 File cycle out of range.
			E:242633 Cannot catalog file because read or write access not allowed.
			E:243233 Creation of file would require illegal dropping of private file.
			E:244433 File is already catalogued.
			E:253733 Relative F-cycle conflict.
		*/

		if fileSpecification.AbsoluteCycle != nil {
			// TODO
			// If a file is created by its absolute cycle and the absolute cycle is not the next numerically sequential
			// absolute F-cycle available, the sequence of F-cycles is updated to point at the newly created F-cycle.
			// In addition, a buffer of noncataloged F-cycles exists between the newly cataloged F-cycle and the
			// previously cataloged F-cycles. This increases the F-cycle range by more than 1.
		} else if fileSpecification.RelativeCycle != nil {
			// TODO
			// To create the next sequential absolute F-cycle, you can use the relative specification +1.
			// If files have been deleted, the relative specification +1 creates the highest numbered deleted F-cycle.
			// When this +1 file is cataloged (by freeing the file or by run termination), its relative F-cycle number
			// is set to 0 and other existing files of the set have their relative F-cycle numbers decreased by 1,
			// thus maintaining consecutive relative numbering.
		} else {
			// We're here with a file set but no cycle spec on a @CAT request. That won't fly.
			facResult.PostMessage(kexec.FacStatusFileAlreadyCataloged, nil)
			resultCode |= 0_500000_000000
		}
	}

	if resultCode&0_400000_000000 != 0 {
		return
	}

	// TODO call mfd services to create the file

	return nil, 0 // TODO
}

func (mgr *FacilitiesManager) catalogRemovableFile(
	exec kexec.IExec,
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *mfdMgr.FileSetInfo,
	usage config.EquipmentUsage,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	//	For Mass Storage Files
	//		@CAT[,options] filename[,type/reserve/granule/maximum,pack-id-1/.../pack-id-n,,,ACR-name]
	//	options include
	//		B: save on checkpoint
	//		G: guarded file
	//		P: make the file public (not private)
	//		R: make the file read-only
	//		V: file will not be unloaded
	//		W: make the file write-only
	//		Z: run should not be held (probably only happens on removable when the pack is not mounted)
	allowedOpts := uint64(kexec.BOption | kexec.GOption | kexec.POption |
		kexec.ROption | kexec.VOption | kexec.WOption | kexec.ZOption)
	if !rce.CheckIllegalOptions(optionWord, allowedOpts, facResult, rce.IsExec) {
		// TODO
	}

	if !mgr.checkSubFields(operandFields, catRemovableFSIs) {
		// TODO
	}

	saveOnCheckpoint := optionWord&kexec.BOption != 0
	guardedFile := optionWord&kexec.GOption != 0
	publicFile := optionWord&kexec.POption != 0
	readOnly := optionWord&kexec.ROption != 0
	inhibitUnload := optionWord&kexec.VOption != 0
	writeOnly := optionWord&kexec.WOption != 0
	doNotHold := optionWord&kexec.ZOption != 0
	wordAddressable := usage == config.EquipmentUsageWordAddressableMassStorage

	// TODO granularity, initial-reserve, max-granules

	// Ensure the pack list is compatible with the files in the fileset (if there is a fileset)
	// Is it okay to just use the highest cycle?
	// TODO

	// If we are removable ensure each pack name is known and mounted.
	// Do not wait for mount if Z option is set
	// TODO

	return nil, 0 // TODO
}

func (mgr *FacilitiesManager) catalogTapeFile(
	exec kexec.IExec,
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *mfdMgr.FileSetInfo,
	usage config.EquipmentUsage,
) (facResult *kexec.FacStatusResult, resultCode uint64) {
	//	For Tape Files
	//		@CAT,options filename,type[/units/log/noise/processor/tape/
	//			format/data-converter/block-numbering/data-compression/
	//			buffered-write/expanded-buffer,reel-1/reel-2/.../reel-n,
	//			expiration-period/mmspec,,ACR-name,CTL-pool]
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
	allowedOpts := uint64(kexec.EOption|kexec.GOption|kexec.HOption|kexec.JOption|
		kexec.LOption|kexec.MOption|kexec.OOption) | kexec.POption | kexec.ROption |
		kexec.SOption | kexec.VOption | kexec.WOption | kexec.ZOption
	if !rce.CheckIllegalOptions(optionWord, allowedOpts, facResult, rce.IsExec) {
		// TODO
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
// and an optional FileSetInfo struct, and returns a list of NodeModel structs
// representing the various equipment models which can be used to satisfy the mnemonic.
// If the mnemonic is an @ASG or @CAT for a file cycle of an existing file set,
// the corresponding FileSetInfo struct must be specified.
// A false return indicates that the mnemonic is not found.
func (mgr *FacilitiesManager) selectEquipmentModel(
	mnemonic string,
	fileSetInfo *mfdMgr.FileSetInfo,
) ([]nodeMgr.NodeModel, config.EquipmentUsage, bool) {

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

	models := make([]nodeMgr.NodeModel, 0)
	usage := entry.Usage
	for _, modelName := range entry.SelectableEquipment {
		model, ok := nodeMgr.NodeModelTable[modelName]
		if ok {
			models = append(models, model)
		}
	}
	return models, usage, true
}
