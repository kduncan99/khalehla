// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"khalehla/hardware"
	"khalehla/kexec"
	"khalehla/kexec/config"
	"khalehla/kexec/mfdMgr"
	"khalehla/kexec/nodeMgr"
	"khalehla/klog"
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

var asgDiskFSIs = newFieldSubfieldIndices().
	add(0, 0).
	add(1, 0).
	add(1, 1).
	add(1, 2).
	add(1, 4).
	addAll(2)

var asgTapeFSIs = newFieldSubfieldIndices().
	add(0, 0).
	add(1, 0).
	add(1, 1).
	add(1, 2).
	add(1, 3).
	add(1, 4).
	add(1, 5).
	add(1, 6).
	add(1, 7).
	add(1, 8).
	add(1, 9).
	add(1, 10).
	add(1, 11).
	addAll(2).
	add(3, 0).
	add(3, 1).
	add(4, 0).
	add(6, 0)

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

var useFSIs = newFieldSubfieldIndices().
	add(0, 0).
	add(1, 0)

// -----------------------------------------------------------------------------

// canAccessFileCycle determines whether the current run can access the file cycle.
func (mgr *FacilitiesManager) canAccessFileCycle(
	rce *kexec.RunControlEntry,
	fileCycleInfo mfdMgr.FileCycleInfo,
) bool {
	var match bool
	if mgr.exec.GetConfiguration().FilesPrivateByAccount {
		match = rce.AccountId == fileCycleInfo.GetAccountId()
	} else {
		match = rce.ProjectId == fileCycleInfo.GetProjectId()
	}

	return match || (rce.IsPrivileged() && !fileCycleInfo.GetInhibitFlags().IsGuarded)
}

func canCatalogFile(
	rce *kexec.RunControlEntry,
	fileSetInfo *mfdMgr.FileSetInfo,
	fileSpecification *kexec.FileSpecification,
	sourceIsExecRequest bool,
	facResult *FacStatusResult,
	resultCode *uint64,
) bool {
	return checkReadKey(rce, fileSetInfo, fileSpecification.ReadKey, sourceIsExecRequest, facResult, resultCode) &&
		checkWriteKey(rce, fileSetInfo, fileSpecification.WriteKey, sourceIsExecRequest, facResult, resultCode)
}

func (mgr *FacilitiesManager) canDropFileCycle(
	rce *kexec.RunControlEntry,
	fileSetInfo *mfdMgr.FileSetInfo,
	fileCycleInfo mfdMgr.FileCycleInfo,
	fileSpecification *kexec.FileSpecification,
	sourceIsExecRequest bool,
	facResult *FacStatusResult,
	resultCode *uint64,
) bool {
	if !mgr.canAccessFileCycle(rce, fileCycleInfo) {
		facResult.PostMessage(kexec.FacStatusIllegalDroppingPrivateFile, nil)
		*resultCode |= 0_400000_020000
		return false
	}

	ok := checkReadKey(rce, fileSetInfo, fileSpecification.ReadKey, sourceIsExecRequest, facResult, resultCode) &&
		checkWriteKey(rce, fileSetInfo, fileSpecification.WriteKey, sourceIsExecRequest, facResult, resultCode)
	if !ok {
		return false
	}

	if fileCycleInfo.IsAssigned() {
		facResult.PostMessage(kexec.FacStatusRelativeFCycleConflict, nil)
		*resultCode |= 0_400000_000040
		return false
	}

	return true
}

// checkIllegalOptionCombination checks to see if more than one of the given options are provided.
func checkIllegalOptionCombination(
	rce *kexec.RunControlEntry,
	givenOptions uint64,
	mutuallyExclusiveOptions uint64,
	facResult *FacStatusResult,
	sourceIsExec bool,
) bool {
	bit := uint64(kexec.AOption)
	letter := 'A'
	ok := true

	var firstOption string
	var secondOption string
	for {
		if bit&givenOptions != 0 && bit&mutuallyExclusiveOptions == 0 {
			if len(firstOption) > 0 {
				secondOption = string(letter)
				facResult.PostMessage(kexec.FacStatusIllegalOptionCombination, []string{firstOption, secondOption})
				ok = false
				break
			} else {
				firstOption = string(letter)
			}
		}
	}

	if !ok {
		if sourceIsExec {
			rce.PostContingency(012, 04, 040)
		}
	}

	return ok
}

// checkReadKey verifies the given read key against the key in the provided fileSetInfo.
// If the file check fails, we update facResult and resultCode, and maybe post a contingency.
func checkReadKey(
	rce *kexec.RunControlEntry,
	fileSetInfo *mfdMgr.FileSetInfo,
	readKey string,
	sourceIsExecRequest bool,
	facResult *FacStatusResult,
	resultCode *uint64,
) bool {
	gaveReadKey := len(readKey) > 0
	hasReadKey := len(fileSetInfo.ReadKey) > 0
	if hasReadKey {
		if !gaveReadKey && (!rce.IsPrivileged() || fileSetInfo.Guarded) {
			facResult.PostMessage(kexec.FacStatusReadWriteKeysNeeded, nil)
			*resultCode |= 0_600000_000000
			return false
		} else if fileSetInfo.ReadKey != readKey {
			facResult.PostMessage(kexec.FacStatusIncorrectReadKey, nil)
			*resultCode |= 0_401000_000000
			if sourceIsExecRequest {
				rce.PostContingencyWithAuxiliary(017, 0, 0, 015)
			}
			return false
		}
	} else {
		if gaveReadKey {
			facResult.PostMessage(kexec.FacStatusFileNotCatalogedWithReadKey, nil)
			*resultCode |= 0_400040_000000
			if sourceIsExecRequest {
				rce.PostContingencyWithAuxiliary(017, 0, 0, 015)
			}
			return false
		}
	}

	return true
}

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

// checkWriteKey verifies the given read key against the key in the provided fileSetInfo.
// If the file check fails, we update facResult and resultCode, and maybe post a contingency.
func checkWriteKey(
	rce *kexec.RunControlEntry,
	fileSetInfo *mfdMgr.FileSetInfo,
	writeKey string,
	sourceIsExecRequest bool,
	facResult *FacStatusResult,
	resultCode *uint64,
) bool {
	gaveWriteKey := len(writeKey) > 0
	hasWriteKey := len(fileSetInfo.ReadKey) > 0
	if hasWriteKey {
		if !gaveWriteKey && (!rce.IsPrivileged() || fileSetInfo.Guarded) {
			facResult.PostMessage(kexec.FacStatusReadWriteKeysNeeded, nil)
			*resultCode |= 0_600000_000000
			return false
		} else if fileSetInfo.WriteKey != writeKey {
			facResult.PostMessage(kexec.FacStatusIncorrectWriteKey, nil)
			*resultCode |= 0_400400_000000
			if sourceIsExecRequest {
				rce.PostContingencyWithAuxiliary(017, 0, 0, 015)
			}
			return false
		}
	} else {
		if gaveWriteKey {
			facResult.PostMessage(kexec.FacStatusFileNotCatalogedWithWriteKey, nil)
			*resultCode |= 0_400020_000000
			if sourceIsExecRequest {
				rce.PostContingencyWithAuxiliary(017, 0, 0, 015)
			}
			return false
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

// resolveFileSpecification follows use item table to find the final external file name
// entry which applies to the caller, and fills in an effective qualifier if necessary.
func (mgr *FacilitiesManager) resolveFileSpecification(
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
) *kexec.FileSpecification {
	result := fileSpecification
	for result.CouldBeInternalName() {
		useItem, ok := rce.UseItems[result.Filename]
		if !ok {
			break
		}

		result = useItem.FileSpecification
	}

	if len(result.Qualifier) == 0 {
		var qual string
		if result.HasAsterisk {
			qual = rce.ImpliedQualifier
		} else {
			qual = rce.DefaultQualifier
		}
		result = kexec.CopyFileSpecification(result)
		result.Qualifier = qual
	}

	return result
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
) ([]hardware.NodeModel, config.EquipmentUsage, bool) {

	effectiveMnemonic := mnemonic

	// If we do not have a given mnemonic but we *do* have a fileSetInfo...
	if len(effectiveMnemonic) == 0 && fileSetInfo != nil {
		// Use the equipment type from the highest absolute fcycle entry of a not to-be file cycle
		//	(an existing file cycle which is not to-be-cataloged or to-be-deleted)... if there is one.
		// Otherwise, use the equipment type from the highest absolute fcycle entry of a to-be file cycle
		for _, preventToBe := range []bool{true, false} {
			for _, fsCycleInfo := range fileSetInfo.CycleInfo {
				if !preventToBe || (!fsCycleInfo.ToBeCataloged && !fsCycleInfo.ToBeDropped) {
					mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
					fcInfo, mfdResult := mm.GetFileCycleInfo(fsCycleInfo.FileCycleIdentifier)
					if mfdResult != mfdMgr.MFDSuccessful {
						mgr.exec.Stop(kexec.StopFacilitiesComplex)
						return nil, 0, false
					}
					effectiveMnemonic = fcInfo.GetAssignMnemonic()
				}
			}
		}
	}

	// If we still do not have an effective mnemonic use the default sector-formatted mass storage mnemonic.
	if len(effectiveMnemonic) == 0 {
		effectiveMnemonic = mgr.exec.GetConfiguration().MassStorageDefaultMnemonic
	}

	// Now go look for the mnemonic in the configured equipment entry table.
	entry, ok := mgr.exec.GetConfiguration().EquipmentTable[mnemonic]
	if !ok {
		return nil, 0, false
	}

	models := make([]hardware.NodeModel, 0)
	usage := entry.Usage
	for _, modelName := range entry.SelectableEquipment {
		model, ok := hardware.NodeModelTable[modelName]
		if ok {
			models = append(models, model)
		}
	}

	return models, usage, true
}

// -----------------------------------------------------------------------------

func (mgr *FacilitiesManager) assignCatalogedFile(
	rce *kexec.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	klog.LogTraceF("FacMgr", "assignCatalogedFile [%v]", rce.RunId)

	mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
	fsIdent, mfdResult := mm.GetFileSetIdentifier(fileSpecification.Qualifier, fileSpecification.Filename)
	if mfdResult == mfdMgr.MFDNotFound {
		klog.LogInfoF("FacMgr", "[%v] file not cataloged", rce.RunId)
		facResult.PostMessage(kexec.FacStatusFileIsNotCataloged, nil)
		resultCode |= 0_400010_000000
		return
	}

	fsInfo, mfdResult := mm.GetFileSetInfo(fsIdent)
	if mfdResult == mfdMgr.MFDNotFound {
		klog.LogInfoF("FacMgr", "[%v] internal error", rce.RunId)
		facResult.PostMessage(kexec.FacStatusFileIsNotCataloged, nil)
		resultCode |= 0_400010_000000
		return
	}

	// for already assigned files, see ECL 4.3 for security access validation
	//   see ECL 7.2.4 for changing certain fields

	/*
		E:241433 Attempt to change assign mnemonic.
		E:241533 Illegal attempt to change assignment type.
		E:241633 Attempt to change generic type of the file.
		E:241733 Attempt to change granularity.
		E:242033 Attempt to change initial reserve of write inhibited file.
		E:242133 Attempt to change maximum granules of a file cataloged under a different account.
		E:242233 Attempt to change maximum granules on a write inhibited file.
		E:242333 Assignment of units of the requested equipment type is not allowed.
	*/

	if fileSpecification.FileCycleSpec.IsRelative() && *fileSpecification.FileCycleSpec.AbsoluteCycle != 0 {
		// TODO assignCatalogedFile() handle relative file cycle
		// The following needs to be excised into a separate function so other top-level assign functions can use it.
		//
		// "When a cataloged file is assigned with a relative F-cycle, the list of currently assigned files is searched
		// to determine whether that relative F-cycle is already assigned to the run. If it is, that F-cycle is used.
		// If the specified relative F-cycle cannot be found, the master file directories are searched to determine
		// whether that relative F-cycle exists. If the relative F-cycle does not exist, the assignment request is rejected.
		//
		// "If the relative F-cycle does exist, it is converted to the corresponding absolute F-cycle number and the list
		// of currently assigned files is searched again to determine whether that absolute F-cycle is already assigned
		// to the run.
		//
		// 'If it is assigned and it already has a relative F-cycle number associated with it, the request is rejected
		// because a naming conflict exists when one absolute F-cycle-has two relative F-cycle values associated with it.
		//
		// 'If it is assigned and no relative F-cycle number is associated with that absolute F-cycle, then the absolute
		// F-cycle is associated with the relative F-cycle.
		//
		// "If it is not already assigned to the run, that relative F-cycle/absolute F-cycle is assigned to the run.
		//
		// "When you specify a relative F-cycle, the file is known to the assigning run by both its relative F-cycle
		// and the appropriate absolute F-cycle number. However, if you specify an absolute F-cycle, the Exec does not
		// attempt to convert it to a relative F-cycle. In this case, the file is known to the assigning run only by its
		// absolute F-cycle."

		facResult.PostMessage(kexec.FacStatusFileIsNotCataloged, nil)
		resultCode |= 0_400010_000000
		klog.LogTraceF("FacMgr", "[%v] file does not exist", rce.RunId)
		return
	}

	var absCycle uint
	if fileSpecification.FileCycleSpec.IsAbsolute() {
		absCycle = *fileSpecification.FileCycleSpec.AbsoluteCycle
	} else {
		// caller either specified relative cycle 0, or did not specify any cycle.
		// either way, we assume highest absolute cycle.
		absCycle = fsInfo.CycleInfo[0].AbsoluteCycle
	}

	// Is it already assigned?
	for _, facItem := range rce.FacilityItems {
		if fileSpecification.Filename == facItem.GetFilename() &&
			fileSpecification.Qualifier == facItem.GetQualifier() &&
			absCycle == facItem.GetAbsoluteCycle() {
			// TODO assignCatalogedFile() check for attempt to change settings
			facResult.PostMessage(kexec.FacStatusFileAlreadyAssigned, nil)
			resultCode |= 0_100000_000000
			klog.LogTraceF("FacMgr", "[%v] already assigned", rce.RunId)
			return
		}
	}

	// See if the thing exists, and if it does, assign it
	for _, ci := range fsInfo.CycleInfo {
		if ci.AbsoluteCycle == absCycle {
			// TODO  assignCatalogedFile()
			//  What if it is to-be-cataloged? pretend it does not yet exist? try experiment
			//  W:121433 File is cataloged as a read-only file.
			//  W:122433 File is cataloged write-only.
			//  E:241233 File is being dropped.
			//  6	That portion of the file Name used as the internal Name for I/O packets is not unique.
			//  7	X option specified; file already in exclusive use.
			//  8*†	Incorrect read key for cataloged file
			//  9*†	Incorrect write key for cataloged file
			//  10	Write key that exists in the master file directory is not specified in the @ASG control statement (file assigned in the read-only mode).
			//  11	Read key that exists in the master file directory is not specified in the @ASG control statement (file assigned in the write-only mode).
			//  12*†Read key specified in the @ASG control statement; none exists in the master file directory.
			//  13*†Write key specified in the @ASG control statement; none exists in the master file directory.
			//  16*	Mass storage file has been rolled out (only if the Z option is used;
			//  otherwise, the run is held until the file is rolled in).
			//  17*	Request on wait Status for facilities. For a tape file, this usually means a tape unit is not currently
			//  available. For a disk file, this usually is caused by an exclusive use conflict with another run
			//  (only if the Z option is used; otherwise, the run is held).
			//  18*	For cataloged files, an option conflict occurred:
			//  The D and K options were specified.
			//  	C or U, or P, R, or W in combination with C or U, was specified for a file that already exists in the directory.
			//  	C was specified on a @CAT image.
			//      For a tape, an option conflict occurred on tape assignment (for example, FJ without Media Manager installed).
			//  19*	File assigned exclusively to another run
			//  20	Find was made on a cataloged file request and the file was already assigned to another run.
			//  21*	File to be decataloged when no run has file assigned
			//  22*	Project-id incorrect for cataloged private file
			//  24	Read-only file cataloged with an R option
			//  25	Write-only file cataloged with a W option
			//  28	File specified on the @ASG control statement has been disabled because the file was assigned during a system failure.
			//  29*	File specified on the @ASG control statement has been disabled because the file has been rolled out and the
			//      backup copy is unrecoverable, unless an @ENABLE command, followed by an @ASG,A command, is used to retry the loading operation.
		}
	}

	facResult.PostMessage(kexec.FacStatusFileIsNotCataloged, nil)
	resultCode |= 0_400010_000000
	klog.LogTraceF("FacMgr", "[%v] file does not exist", rce.RunId)
	return
}

// assignTemporaryFile is invoked for any @ASG,T situation.
func (mgr *FacilitiesManager) assignTemporaryFile(
	rce *kexec.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	models []hardware.NodeModel,
	usage config.EquipmentUsage,
) (facResult *FacStatusResult, resultCode uint64) {
	klog.LogTraceF("FacMgr", "assignTemporaryFile [%v]", rce.RunId)

	// For temporary files, we ignore any provided read/write keys.
	// Check fac items to see if an item already exists with the given specification.
	alreadyAssigned := false
	var prevFacItem kexec.FacilitiesItem
	for _, facItem := range rce.FacilityItems {
		if facItem.GetQualifier() == fileSpecification.Qualifier &&
			facItem.GetFilename() == fileSpecification.Filename {
			if fileSpecification.FileCycleSpec == nil {
				if facItem.GetRelativeCycle() == 0 && facItem.GetAbsoluteCycle() == 0 {
					alreadyAssigned = true
				}
			} else if fileSpecification.FileCycleSpec.IsRelative() {
				alreadyAssigned = *fileSpecification.FileCycleSpec.RelativeCycle == facItem.GetRelativeCycle()
			} else if fileSpecification.FileCycleSpec.IsAbsolute() {
				alreadyAssigned = *fileSpecification.FileCycleSpec.AbsoluteCycle == facItem.GetAbsoluteCycle()
			}

			if alreadyAssigned {
				prevFacItem = facItem
				break
			}
		}
	}

	if alreadyAssigned {
		// Check whether caller is attempting to change the general file type
		// or to apply a non-conforming equipment type.
		equipSpecified := len(operandFields) >= 2 && len(operandFields[1][0]) > 0
		if equipSpecified {
			if usage == config.EquipmentUsageSectorAddressableMassStorage ||
				usage == config.EquipmentUsageWordAddressableMassStorage {
				if !prevFacItem.IsDisk() {
					facResult.PostMessage(kexec.FacStatusAttemptToChangeGenericType, nil)
					resultCode |= 0_600000_000000
					klog.LogInfoF("FacMgr", "assignTemporaryFile exit resultCode %012o", resultCode)
					return
				}
			} else if usage == config.EquipmentUsageTape {
				if !prevFacItem.IsTape() {
					facResult.PostMessage(kexec.FacStatusAttemptToChangeGenericType, nil)
					resultCode |= 0_600000_000000
					klog.LogInfoF("FacMgr", "assignTemporaryFile exit resultCode %012o", resultCode)
					return
				}
			}

			// TODO incompatible equipment? (among the matching general type)
		}

		facResult.PostMessage(kexec.FacStatusFileAlreadyAssigned, nil)
		resultCode |= 0_100000_000000
		klog.LogInfoF("FacMgr", "assignTemporaryFile exit resultCode %012o", resultCode)
		return
	}

	if usage == config.EquipmentUsageSectorAddressableMassStorage ||
		usage == config.EquipmentUsageWordAddressableMassStorage {

		allowedOpts := uint64(kexec.IOption | kexec.TOption | kexec.ZOption)
		if !CheckIllegalOptions(rce, optionWord, allowedOpts, facResult, rce.IsExec()) {
			resultCode |= 0_600000_000000
			klog.LogInfoF("FacMgr", "assignTemporaryFile exit resultCode %012o", resultCode)
			return
		}

		if !mgr.checkSubFields(operandFields, asgDiskFSIs) {
			facResult.PostMessage(kexec.FacStatusUndefinedFieldOrSubfield, nil)
			resultCode |= 0_600000_000000
			klog.LogInfoF("FacMgr", "assignTemporaryFile exit resultCode %012o", resultCode)
			return
		}
	} else if usage == config.EquipmentUsageTape {

		allowedOpts := uint64(kexec.EOption | kexec.FOption | kexec.HOption | kexec.IOption | kexec.JOption |
			kexec.LOption | kexec.MOption | kexec.NOption | kexec.OOption | kexec.ROption | kexec.SOption |
			kexec.TOption | kexec.VOption | kexec.WOption | kexec.XOption | kexec.ZOption)
		if !CheckIllegalOptions(rce, optionWord, allowedOpts, facResult, rce.IsExec()) {
			resultCode |= 0_600000_000000
			klog.LogInfoF("FacMgr", "assignTemporaryFile exit resultCode %012o", resultCode)
			return
		}

		if !mgr.checkSubFields(operandFields, asgTapeFSIs) {
			facResult.PostMessage(kexec.FacStatusUndefinedFieldOrSubfield, nil)
			resultCode |= 0_600000_000000
			klog.LogInfoF("FacMgr", "assignTemporaryFile exit resultCode %012o", resultCode)
			return
		}
	}

	// Is filename not unique?
	for _, facItem := range rce.FacilityItems {
		if facItem.GetFilename() == fileSpecification.Filename {
			facResult.PostMessage(kexec.FacStatusFilenameNotUnique, nil)
			resultCode |= 0_004000_000000
			klog.LogInfoF("FacMgr", "assignTemporaryFile exit resultCode %012o", resultCode)
			break
		}
	}

	// TODO check duplicate reel numbers specified
	// TODO check for requested tape reel is already assigned by this run.

	// Do usage-specific stuff
	if usage == config.EquipmentUsageSectorAddressableMassStorage ||
		usage == config.EquipmentUsageWordAddressableMassStorage {
		// TODO For Mass Storage, we need to honor allocation requests and set up a local allocation table
		// TODO For Removable, we may need to wait for pack mount(s)
	} else if usage == config.EquipmentUsageTape {
		// TODO check for the reel already assigned to some other run (hold condition)
		// TODO check for tape unit availability (hold condition)
	}

	klog.LogInfoF("FacMgr", "assignTemporaryFile exit resultCode %012o", resultCode)
	return
}

func (mgr *FacilitiesManager) assignToBeCatalogedFile(
	rce *kexec.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
) (facResult *FacStatusResult, resultCode uint64) {
	klog.LogTraceF("FacMgr", "assignToBeCatalogedFile [%v]", rce.RunId)

	// TODO implement assignToBeCatalogedFile()
	mm := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager)
	_ /*fsIdent*/, mfdResult := mm.GetFileSetIdentifier(fileSpecification.Qualifier, fileSpecification.Filename)
	if mfdResult == mfdMgr.MFDSuccessful {
		facResult.PostMessage(kexec.FacStatusFileAlreadyCataloged, nil)
		resultCode |= 0_500000_000000
		return
	}

	/*
		E:241433 Attempt to change assign mnemonic.
		E:241533 Illegal attempt to change assignment type.
		E:241633 Attempt to change generic type of the file.
		E:241733 Attempt to change granularity.
		E:242033 Attempt to change initial reserve of write inhibited file.
		E:242133 Attempt to change maximum granules of a file cataloged under a different account.
		E:242233 Attempt to change maximum granules on a write inhibited file.
		E:242333 Assignment of units of the requested equipment type is not allowed.
	*/

	return
}

func (mgr *FacilitiesManager) catalogCommon(
	exec kexec.IExec,
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *mfdMgr.FileSetInfo,
	mnemonic string,
	usage config.EquipmentUsage,
	isRemovable bool,
	sourceIsExecRequest bool,
	facResult *FacStatusResult,
	resultCode *uint64,
) bool {
	saveOnCheckpoint := optionWord&kexec.BOption != 0
	guardedFile := optionWord&kexec.GOption != 0
	publicFile := optionWord&kexec.POption != 0
	readOnly := optionWord&kexec.ROption != 0
	inhibitUnload := optionWord&kexec.VOption != 0
	writeOnly := optionWord&kexec.WOption != 0
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
		*resultCode |= 0_600000_000000
		return false
	}

	var initReserve uint64
	if len(initStr) > 12 {
		facResult.PostMessage(kexec.FacStatusIllegalInitialReserve, nil)
		*resultCode |= 0_600000_000000
		return false
	} else if len(initStr) > 0 {
		initReserve, err := strconv.Atoi(initStr)
		if err != nil || initReserve < 0 {
			facResult.PostMessage(kexec.FacStatusIllegalInitialReserve, nil)
			*resultCode |= 0_600000_000000
			return false
		}
	}

	maxGranules := exec.GetConfiguration().MaxGranules
	if len(maxStr) > 12 {
		facResult.PostMessage(kexec.FacStatusIllegalMaxGranules, nil)
		*resultCode |= 0_600000_000000
		return false
	} else if len(maxStr) > 0 {
		iMaxGran, err := strconv.Atoi(maxStr)
		maxGranules = uint64(iMaxGran)
		if err != nil || maxGranules < 0 || maxGranules > 262143 {
			facResult.PostMessage(kexec.FacStatusIllegalMaxGranules, nil)
			*resultCode |= 0_600000_000000
			return false
		} else if maxGranules < initReserve {
			facResult.PostMessage(kexec.FacStatusMaximumIsLessThanInitialReserve, nil)
			*resultCode |= 0_600000_000000
			return false
		}
	}

	// If there isn't an existing fileset, create one.
	mm := exec.GetMFDManager().(*mfdMgr.MFDManager)
	if fileSetInfo == nil {
		_, result := mm.CreateFileSet(
			mfdMgr.FileTypeFixed,
			fileSpecification.Qualifier,
			fileSpecification.Filename,
			rce.ProjectId,
			fileSpecification.ReadKey,
			fileSpecification.WriteKey)
		if result == mfdMgr.MFDInternalError {
			return false
		} else if result != mfdMgr.MFDSuccessful {
			klog.LogFatal("FacMgr", "MFD failed to create file set")
			exec.Stop(kexec.StopFacilitiesComplex)
			return false
		}
	} else {
		if !canCatalogFile(rce, fileSetInfo, fileSpecification, sourceIsExecRequest, facResult, resultCode) {
			return false
		}
	}

	descriptorFlags := mfdMgr.DescriptorFlags{
		SaveOnCheckpoint:    saveOnCheckpoint,
		IsTapeFile:          false,
		IsRemovableDiskFile: false,
	}
	pcharFlags := mfdMgr.PCHARFlags{
		Granularity:       granularity,
		IsWordAddressable: wordAddressable,
	}
	inhibitFlags := mfdMgr.InhibitFlags{
		IsGuarded:         guardedFile,
		IsUnloadInhibited: inhibitUnload,
		IsPrivate:         publicFile,
		IsWriteOnly:       writeOnly,
		IsReadOnly:        readOnly,
	}

	retry := true
	for retry {
		var mfdResult mfdMgr.MFDResult
		if isRemovable {
			_, mfdResult = mm.CreateRemovableFileCycle(
				fileSetInfo.FileSetIdentifier,
				fileSpecification.FileCycleSpec,
				rce.AccountId,
				mnemonic,
				descriptorFlags,
				pcharFlags,
				inhibitFlags,
				initReserve,
				maxGranules,
				make([]mfdMgr.DiskPackEntry, 0))
		} else {
			_, mfdResult = mm.CreateFixedFileCycle(
				fileSetInfo.FileSetIdentifier,
				fileSpecification.FileCycleSpec,
				rce.AccountId,
				mnemonic,
				descriptorFlags,
				pcharFlags,
				inhibitFlags,
				initReserve,
				maxGranules,
				make([]mfdMgr.DiskPackEntry, 0))
		}

		retry = false
		switch mfdResult {
		case mfdMgr.MFDSuccessful: // nothing to do
		case mfdMgr.MFDInternalError: // nothing to do, we're already dead in the water
		case mfdMgr.MFDAlreadyExists:
			facResult.PostMessage(kexec.FacStatusFileAlreadyCataloged, nil)
			*resultCode |= 0_500000_000000
		case mfdMgr.MFDInvalidAbsoluteFileCycle:
			facResult.PostMessage(kexec.FacStatusFileCycleOutOfRange, nil)
			*resultCode |= 0_600000_000040
		case mfdMgr.MFDInvalidRelativeFileCycle:
			facResult.PostMessage(kexec.FacStatusRelativeFCycleConflict, nil)
			*resultCode |= 0_600000_000040
		case mfdMgr.MFDPlusOneCycleExists:
			facResult.PostMessage(kexec.FacStatusRelativeFCycleConflict, nil)
			*resultCode |= 0_600000_000040
		case mfdMgr.MFDDropOldestCycleRequired:
			cx := len(fileSetInfo.CycleInfo) - 1
			fcIdentifier := fileSetInfo.CycleInfo[cx].FileCycleIdentifier
			fcInfo, _ := mm.GetFileCycleInfo(fcIdentifier)
			if !mgr.canDropFileCycle(rce, fileSetInfo, fcInfo, fileSpecification, sourceIsExecRequest, facResult, resultCode) {
				return false
			}

			result := mm.DropFileCycle(fcIdentifier)
			if result == mfdMgr.MFDInternalError {
				return false
			} else if result != mfdMgr.MFDSuccessful {
				klog.LogFatal("FacMGR", "Cannot delete oldest file cycle")
				mgr.exec.Stop(kexec.StopFacilitiesComplex)
				*resultCode |= 0_400000_000000
				return false
			}
			retry = true
		}
	}

	return true
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
	mnemonic string,
	usage config.EquipmentUsage,
	sourceIsExecRequest bool,
) (facResult *FacStatusResult, resultCode uint64) {
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
	//			I'm unaware of any situation where cataloging a fixed file would result in a hold.
	facResult = NewFacResult()
	resultCode = 0

	allowedOpts := uint64(kexec.BOption | kexec.GOption | kexec.POption |
		kexec.ROption | kexec.VOption | kexec.WOption | kexec.ZOption)
	if !CheckIllegalOptions(rce, optionWord, allowedOpts, facResult, rce.IsExec()) {
		resultCode |= 0_600000_000000
		return
	}

	if !mgr.checkSubFields(operandFields, catFixedFSIs) {
		if len(mgr.getSubField(operandFields, 1, 4)) > 0 {
			facResult.PostMessage(kexec.FacStatusPlacementFieldNotAllowed, nil)
		}
		facResult.PostMessage(kexec.FacStatusUndefinedFieldOrSubfield, nil)
		resultCode |= 0_600000_000000
		return
	}

	mgr.catalogCommon(
		exec,
		rce,
		fileSpecification,
		optionWord,
		operandFields,
		fileSetInfo,
		mnemonic,
		usage,
		false,
		sourceIsExecRequest,
		facResult,
		&resultCode)
	return
}

func (mgr *FacilitiesManager) catalogRemovableFile(
	exec kexec.IExec,
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *mfdMgr.FileSetInfo,
	mnemonic string,
	usage config.EquipmentUsage,
	sourceIsExecRequest bool,
) (facResult *FacStatusResult, resultCode uint64) {
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
	facResult = NewFacResult()
	resultCode = 0

	allowedOpts := uint64(kexec.BOption | kexec.GOption | kexec.POption |
		kexec.ROption | kexec.VOption | kexec.WOption | kexec.ZOption)
	if !CheckIllegalOptions(rce, optionWord, allowedOpts, facResult, rce.IsExec()) {
		resultCode |= 0_600000_000000
		return
	}

	if !mgr.checkSubFields(operandFields, catRemovableFSIs) {
		if len(mgr.getSubField(operandFields, 1, 4)) > 0 {
			facResult.PostMessage(kexec.FacStatusPlacementFieldNotAllowed, nil)
		}
		facResult.PostMessage(kexec.FacStatusUndefinedFieldOrSubfield, nil)
		resultCode |= 0_600000_000000
		return
	}

	mgr.catalogCommon(
		exec,
		rce,
		fileSpecification,
		optionWord,
		operandFields,
		fileSetInfo,
		mnemonic,
		usage,
		false,
		sourceIsExecRequest,
		facResult,
		&resultCode)
	return
}

func (mgr *FacilitiesManager) catalogTapeFile(
	exec kexec.IExec,
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	fileSetInfo *mfdMgr.FileSetInfo,
	mnemonic string,
	usage config.EquipmentUsage,
	sourceIsExecRequest bool,
) (facResult *FacStatusResult, resultCode uint64) {
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
	if !CheckIllegalOptions(rce, optionWord, allowedOpts, facResult, rce.IsExec()) {
		// TODO implement CatalogTapeFile()
	}

	return nil, 0
}

func (mgr *FacilitiesManager) getNodeStatusString(nodeId hardware.NodeIdentifier) string {
	accStr := "   "
	nm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
	ni, _ := nm.GetNodeInfoByIdentifier(nodeId)
	if !ni.IsAccessible() {
		accStr = " NA"
	}

	str := ""
	facStat := mgr.inventory.nodes[nodeId].GetFacNodeStatus()
	switch facStat {
	case kexec.FacNodeStatusDown:
		str = "DN" + accStr
	case kexec.FacNodeStatusReserved:
		str = "RV" + accStr
	case kexec.FacNodeStatusSuspended:
		str = "SU" + accStr
	case kexec.FacNodeStatusUp:
		str = "UP" + accStr
	}

	diskAttr, ok := mgr.inventory.disks[nodeId]
	if ok {
		//	[[*] [R|F] PACKID pack-id]
		if diskAttr.AssignedTo != nil {
			str += " * "
		} else {
			str += "   "
		}

		if diskAttr.PackLabelInfo != nil {
			if diskAttr.IsFixed {
				str += "F "
			} else if diskAttr.IsRemovable {
				str += "R "
			} else {
				str += "  "
			}

			str += "PACKID " + diskAttr.PackLabelInfo.PackId
		}
	}

	// ta, ok := mgr.inventory.tapes[deviceId]
	// if ok {
	//	if ta.AssignedTo != nil {
	//		//	[* RUNID run-id REEL reel [RING|NORING] [POS [*]ffff[+|-][*]bbbbbb | POS LOST]]
	//		str += "* RUNID " + ta.AssignedTo.RunId + " REEL " + ta.reelNumber
	//		// TODO RING | NORING in FS display for tape unit
	//		// TODO POS in FS display for tape unit
	//	}
	// }

	return str
}
