// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"fmt"
	"khalehla/hardware"
	"khalehla/kexec"
	"khalehla/kexec/config"
	"khalehla/kexec/facItems"
	"khalehla/kexec/mfdMgr"
	"khalehla/kexec/nodeMgr"
	"khalehla/klog"
	"strconv"
	"strings"
	"time"
)

/*
TODO Placement topics for fixed disk
placement

  Specifies the placement of the file on a controller or device. The specification can be a logical or absolute request.
  The placement can be any of the following:
    *cu
      Is a controller name (absolute request).
   *device
   	  Is a device name (absolute request).
   i
      i Is a letter, A to Z, representing logical subsystems A through Z
   in
      n is a number, 1 to 15, representing a logical device on controller i.

If an absolute placement cannot be done as specified, it is rejected.
The Exec attempts the following file assignments:
  - Files with the same logical specifications (that is, the same controller and device number) are assigned to the same physical device.
  - Files with the same logical controller but different device numbers are assigned to different physical devices on the same control unit.
  - Files with different logical controllers and device numbers are assigned to physical devices that are reachable by
    completely different CUs within the CONV selection string for the assign type specified.

Since logical placement is driven off of logical CU selections, specifying logical placement causes the Exec to avoid
placing this file into memory space, even if MEMFL is specified for the assign type. Consequently, do not use the
logical placement field if you want a file placed in memory on systems with MEMFLSZ configured. If a logical
specification cannot be honored on a device within the CONV selection string, processing of the assign request ignores
the placement field.

Placement on removable disk is not supported. The pack-ID enables you to specify allocation on a removable disk.

When the file is cataloged, placement information is placed in the directory to indicate the type of selection
(logical or absolute, controller or device) and the device chosen for initial allocation. The directory cell containing
this information is not affected by subsequent assignment of the file.

On reassignment, a placement specification is taken as a placement change overriding the last device used for allocation.
A maximum of six characters, seven in the case of absolute placement, is allowed for specification of this subfield.

Logical placement specifications have meaning only within the current run. For example, logical controller "A" can
represent two different physical controllers in two different runs.

Logical placement applies only to the reserve specified. If the file dynamically expands beyond the initial reserve,
the Exec attempts to allocate on the same device as the last placement specified. If this is not possible,the expansion
occurs on devices other than those specified in the placement field.

If the file is expanded through static expansion (reserve specified) and the file was created with absolute placement,
the Exec allocates mass storage only on the device originally specified, unless a different device is specified in the
placement field on the assignment. If there is not enough room on the specified (either implicitly or explicitly) device
to satisfy the request, the request is rejected.
*/

/*
TODO pack-id topics for removable disk
pack-ID
  Specifies the removable disk packs required for the file. Pack-IDs consist of from one to six characters of the set
  A through Z, 0 through 9. The pack-IDs for cataloged files are recorded in the master file directory and need not be
  specified on reassignments.

A pack's master file directory must have sufficient space available to record the file otherwise a generation request is
rejected with the following message:
  E:207433 File cannot be created due to lack of removable directory space on pack pack-id.

A pack-ID cannot be specified unless the type subfield is also specified. If pack-ID is omitted on requests for
non-cataloged file space, fixed disk is assumed if disk equipment is requested. For cataloged removable disk files,
specification of a pack-ID causes the pack to be mounted and registered. If already registered, the pack information is
compared with the image for compatibility.

All cycles of a removable disk file must be on the same set of packs.
The maximum number of pack-IDs that can be specified on an @ASG statement is 510.
Many jobs can specify the same set of removable disk packs for unique files.
Pack-IDs can be added to unassigned single file cycle files by specifying the original pack-ID list with the new
  pack-IDs appended on an @ASG,A request.
Packs can only be added if the A option is specified (but not the Y option).
Packs cannot be added to currently assigned files or to files that have more than one file cycle.
If packs are being added, you must have delete access to the file.
*/

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

var freeFSIs = newFieldSubfieldIndices().
	add(0, 0)

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

// checkIllegalOptionCombination checks to see if more than one of the exclusive options have been specified.
func checkIllegalOptionCombination(
	rce *kexec.RunControlEntry,
	givenOptions uint64,
	mutuallyExclusiveOptions uint64,
	facResult *FacStatusResult,
	sourceIsExec bool,
) bool {
	leftBit := uint64(kexec.AOption)
	leftOpt := 'A'
	ok := true

	for leftBit != uint64(kexec.ZOption) {
		if leftBit&givenOptions != 0 {
			rightBit := leftBit >> 1
			rightOpt := leftOpt + 1
			for {
				if rightBit&givenOptions != 0 {
					mask := leftBit | rightBit
					if mutuallyExclusiveOptions&mask == mask {
						facResult.PostMessage(
							kexec.FacStatusIllegalOptionCombination,
							[]string{string(leftOpt), string(rightOpt)})
						ok = false
					}
				}
				if rightBit == kexec.ZOption {
					break
				} else {
					rightBit >>= 1
					rightOpt++
				}
			}
		}

		leftBit >>= 1
		leftOpt++
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

// assignWait
// Disk assign should *not* invoke us if E or Y option is set.
// Callers *should* invoke us even with Z-option set, so that we can reject the request if necessary
// with the proper fac status codes.
func (mgr *FacilitiesManager) assignWait(
	rce *kexec.RunControlEntry,
	optionWord uint64,
	fsInfo *mfdMgr.FixedFileCycleInfo, // nil if we are waiting on temporary file
	usage config.EquipmentUsage,
	facResult *FacStatusResult,
) (resultCode uint64) {
	klog.LogTraceF("FacMgr", "[%v] assignWait", rce.RunId)
	resultCode = 0

	xOption := optionWord&kexec.XOption != 0
	zOption := optionWord&kexec.ZOption != 0

	// TODO also, waiting on reel number (already assigned to another run)
	rollbackRequested := false
	waitingForDiskUnit := false
	waitingForTapeUnit := false
	waitingForXUse := false       // we need X-use and file is assigned elsewhere
	waitingOnXUseRelease := false // we need the file, but it is assigned exclusively elsewhere
	waitingOnRollback := false    // file is rolled out
	var waitingForDiskUnitStart time.Time
	var waitingForTapeUnitStart time.Time
	var waitingForXUseStart time.Time
	var waitingOnXUseReleaseStart time.Time
	var waitingOnRollbackStart time.Time
	waitStartTime := time.Now()
	lastMessageSent := waitStartTime.Add(time.Duration(-5) * time.Minute)

	tapeUnits := make([]hardware.NodeIdentifier, 0)

	for {
		// Cataloged file involved?
		if fsInfo != nil {
			// Is the file currently rolled out?
			if fsInfo.DescriptorFlags.Unloaded {
				if !rollbackRequested {
					// we have to do this bit instead of relying on waitingOnRollback
					// because specifying Z option won't wait, but *does* request rollback.
					mgr.exec.RollbackFile(fsInfo.Qualifier, fsInfo.Filename, fsInfo.AbsoluteFileCycle)
					rollbackRequested = true
				}
				if zOption {
					facResult.PostMessage(kexec.FacStatusHoldForRollbackRejected, nil)
					resultCode |= 0_400001_000000
					break
				}
				if !waitingOnRollback {
					waitingOnRollbackStart = time.Now()
					waitingOnRollback = true
					// TODO set wait on Rollback in RCE
				}
			} else {
				if waitingOnRollback {
					// TODO clear wait on Rollback in RCE
					waitingOnRollback = false
				}
			}

			// Is the file exclusively assigned elsewhere? (we wouldn't be here if it was assigned to us)
			if fsInfo.InhibitFlags.IsAssignedExclusive {
				if zOption {
					facResult.PostMessage(kexec.FacStatusHoldForReleaseXUseRejected, nil)
					resultCode |= 0_400001_000000
					break
				}
				if !waitingOnXUseRelease {
					waitingOnXUseReleaseStart = time.Now()
					waitingOnXUseRelease = true
					// TODO set wait on XUse in RCE
				}
			} else {
				if waitingOnXUseRelease {
					// TODO clear wait on XUse in RCE
					waitingOnXUseRelease = false
				}
			}

			// Are we asking for exclusive use of an already-assigned file?
			if xOption && !fsInfo.InhibitFlags.IsReadOnly && fsInfo.AssignedIndicator > 0 {
				if zOption {
					facResult.PostMessage(kexec.FacStatusHoldForXUseRejected, nil)
					resultCode |= 0_400001_000000
					break
				}
				if !waitingForXUse {
					waitingForXUseStart = time.Now()
					waitingForXUse = true
					// TODO set wait for XUse in RCE (could be more than one, if multiple activities)
				}
			} else {
				if waitingForXUse {
					// TODO clear wait for XUse in RCE (could be more than one)
					waitingForXUse = false
				}
			}
		}

		if (usage == config.EquipmentUsageTape) && (len(tapeUnits) == 0) {
			// TODO we're waiting for 1 or 2 tape units
			//   If we can get them, attach them and note them in tapeUnits
			//   Otherwise, set waitingForTapeUnit
		}

		if usage == config.EquipmentUsageWordAddressableMassStorage ||
			usage == config.EquipmentUsageSectorAddressableMassStorage {
			// TODO are we removable and needing a pack mounted?
			//   If so, look for a free disk unit
			//     if we can get it, attach it and note it in (tbd)
			//     Otherwise, set waitingForDiskUnit
		}

		if !waitingForDiskUnit && !waitingForTapeUnit && !waitingOnRollback && !waitingOnXUseRelease && !waitingForXUse {
			break
		}

		mgr.mutex.Unlock()

		// Time to send messages? (only for batch or demand)
		if rce.IsBatch() || rce.IsDemand() {
			now := time.Now()
			if now.After(lastMessageSent.Add(time.Duration(2) * time.Minute)) {
				if waitingForDiskUnit {
					// "Run %v held for disk unit availability for %v min.")
					waitTime := now.Sub(waitingForDiskUnitStart).Minutes()
					msg := fmt.Sprintln(FacStatusMessageTemplates[kexec.FacStatusRunHeldForDiskUnitAvailability], rce.RunId, waitTime)
					rce.PostToPrint(msg, 1)
				}
				if waitingForTapeUnit {
					// "Run %v held for tape unit availability for %v min.")
					waitTime := now.Sub(waitingForTapeUnitStart).Minutes()
					msg := fmt.Sprintln(FacStatusMessageTemplates[kexec.FacStatusRunHeldForTapeUnitAvailability], rce.RunId, waitTime)
					rce.PostToPrint(msg, 1)
				}
				if waitingForXUse {
					// "Run %v held for need of exclusive use for %v min."
					waitTime := now.Sub(waitingForXUseStart).Minutes()
					msg := fmt.Sprintln(FacStatusMessageTemplates[kexec.FacStatusRunHeldForNeedOfExclusiveUse], rce.RunId, waitTime)
					rce.PostToPrint(msg, 1)
				}
				if waitingOnXUseRelease {
					// "Run %v held for exclusive file use release for %v min."
					waitTime := now.Sub(waitingOnXUseReleaseStart).Minutes()
					msg := fmt.Sprintln(FacStatusMessageTemplates[kexec.FacStatusRunHeldForExclusiveFileUseRelease], rce.RunId, waitTime)
					rce.PostToPrint(msg, 1)
				}
				if waitingOnRollback {
					// "Run %v held for rollback of unloaded file for %v min."
					waitTime := now.Sub(waitingOnRollbackStart).Minutes()
					msg := fmt.Sprintln(FacStatusMessageTemplates[kexec.FacStatusRunHeldForRollback], rce.RunId, waitTime)
					rce.PostToPrint(msg, 1)
				}
				lastMessageSent = now
			}
		}

		time.Sleep(50 * time.Millisecond)
		mgr.mutex.Lock()
	}

	klog.LogTraceF("FacMgr", "[%v] returning %012o", rce.RunId, resultCode)
	return
}

// assignCatalogedFixedFileToRun assumes the requested file exists, and that all security checks are complete.
// We are responsible for waiting on or for exclusive use, and for rollback if appropriate.
// We are also responsible for creating the appropriate facility item.
func (mgr *FacilitiesManager) assignCatalogedFixedFileToRun(
	rce *kexec.RunControlEntry,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	fcInfo *mfdMgr.FixedFileCycleInfo,
	facResult *FacStatusResult,
) (resultCode uint64) {
	resultCode = 0
	if optionWord&(kexec.EOption|kexec.YOption) == 0 {
		resultCode |= mgr.assignWait(rce, 0, fcInfo, facResult)
	}

	// Check for hold for fcycle conflict
	// TODO

	// Check for any applicable warning status codes
	// (except file-already-assigned - we should already have taken care of that)
	// TODO

	// Accelerate the fc info
	mfdResult := mgr.exec.GetMFDManager().(*mfdMgr.MFDManager).AccelerateFileCycle(fcInfo.FileCycleIdentifier)
	if mfdResult != mfdMgr.MFDSuccessful {
		mgr.exec.Stop(kexec.StopFacilitiesComplex)
		return
	}

	// Create facility item and add it to the rce
	relCycle := 0
	if fileSpecification.HasFileCycleSpecification() && fileSpecification.FileCycleSpec.IsRelative() {
		relCycle = *fileSpecification.FileCycleSpec.RelativeCycle
	}

	fi := &kexec.SectorAddressableFacilityItem{
		Qualifier:              fcInfo.Qualifier,
		Filename:               fcInfo.Filename,
		EquipmentCode:          036, // sector-addressable mass storage
		RelativeCycle:          relCycle,
		AbsoluteCycle:          fcInfo.AbsoluteFileCycle,
		Attributes:             kexec.FacItemAttributes{},
		FileMode:               kexec.FacItemFileMode{},
		Granularity:            fcInfo.PCHARFlags.Granularity,
		InitialReserve:         fcInfo.InitialGranulesReserved,
		MaxGranules:            fcInfo.MaxGranules,
		HighestTrackReferenced: fcInfo.HighestTrackWritten,
		HighestGranuleAssigned: fcInfo.HighestGranuleAssigned,
		TotalPackCount:         0,
	}
	rce.FacilityItems = append(rce.FacilityItems, fi)

	return
}

func (mgr *FacilitiesManager) assignCatalogedFile(
	rce *kexec.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	facResult *FacStatusResult,
) (resultCode uint64) {
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

	// TODO check read/write keys (we don't know cycle yet, so we cannot check for private access)

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

	if fileSpecification.HasFileCycleSpecification() &&
		fileSpecification.FileCycleSpec.IsRelative() &&
		*fileSpecification.FileCycleSpec.AbsoluteCycle != 0 {
		// TODO E:252033 F-cycle of +1 is illegal with A option.

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
	if fileSpecification.HasFileCycleSpecification() &&
		fileSpecification.FileCycleSpec.IsAbsolute() {
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
			if ci.ToBeCataloged {
				// It is to-be-cataloged, and thus does not yet exist.
			}

			fcInfo, foo := mm.GetFileCycleInfo(ci.FileCycleIdentifier)
			if foo != mfdMgr.MFDSuccessful {
				klog.LogFatalF("FacMgr",
					"MFDMgr could not find file cycle info for cycle which should exist %012o",
					ci.FileCycleIdentifier)
				mgr.exec.Stop(kexec.StopFacilitiesComplex)
				return
			}

			// Private file violation?
			if !mgr.canAccessFileCycle(rce, fcInfo) {
				facResult.PostMessage(kexec.FacStatusIncorrectPrivacyKey, nil)
				resultCode |= 0_400000_020000
				return
			}

			// TODO
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

			resultCode |= mgr.assignCatalogedFixedFileToRun(rce, fileSpecification, optionWord, fcInfo, facResult)
			return
		}
	}

	// Did not find the cycle - complain that it does not exist
	facResult.PostMessage(kexec.FacStatusFileIsNotCataloged, nil)
	resultCode |= 0_400010_000000
	klog.LogTraceF("FacMgr", "[%v] file does not exist", rce.RunId)
	return
}

// TODO For tape files, we assign the unit first, then we have some other subsequent bit
//  which is responsible for deciding how/what reel number(s) to mount.
//  But what does that look like?

// TODO what happens when two activities assign the same rolled-out file?
//  cannot simply return already-assigned...

// holdForDiskFile() waits for the various states of availability for a cataloged fixed disk file assign,
// specific to a particular activity.
// We do not interact with the system console (that is done elsewhere), and we only send
// PRINT$ messages for batch/demand while in control mode (i.e., not doing ER CSF$).
// We specifically do holds here for:
//   rollback of file
//   cataloged file need for exclusive use (someone else has the file assigned)
//   cataloged file exclusive use release (someone else has x-use)
// Nothing else applies.

func holdForDiskFile(
	rce *kexec.RunControlEntry,
	facItem facItems.IFacilitiesItem,
	sourceIsExecRequest bool,
	facResult *FacStatusResult,
) (resultCode uint64) {
	if facItem.GetOptionWord()&kexec.EOption != 0 || facItem.GetOptionWord()&kexec.YOption != 0 {
		if facItem.GetOptionWord()&kexec.ZOption != 0 {
			// TODO if file is rolled out, start rollback.
			//   do not wait - instead return bad status
		} else {
			// todo wait
		}
	}
	return
}

// holdForTape() waits for the various states of availability for a tape or tape unit assign,
// specific to a particular activity.
// We do not interact with the system console (that is done elsewhere), and we only send
// PRINT$ messages for batch/demand while in control mode (i.e., not doing ER CSF$).
// We do not do holds on reel numbers here - that is done subsequent to file and unit assignation.
// We specifically do holds here for:
//   tape unit availability
//   cataloged file need for exclusive use (someone else has the file assigned)
//   cataloged file exclusive use release (someone else has x-use)
// Nothing else applies.

func holdForTape(
	rce *kexec.RunControlEntry,
	facItem facItems.IFacilitiesItem,
	sourceIsExecRequest bool,
	facResult *FacStatusResult,
) (resultCode uint64) {
	if facItem.GetOptionWord()&kexec.ZOption != 0 {
		// TODO do not wait - instead return bad status
	} else {
		// todo wait
	}
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
	usage config.EquipmentUsage, // TODO does this account for unit assigns? I'm not sure...
	facResult *FacStatusResult,
) (resultCode uint64) {
	klog.LogTraceF("FacMgr", "[%v] assignTemporaryFile", rce.RunId)
	defer klog.LogTraceF("FacMgr", "[%v] assignTemporaryFile result %012o", rce.RunId, resultCode)

	// What type of assign are we doing?
	var isDiskFile bool
	var isDiskUnit bool
	var isTapeFile bool
	var isTapeUnit bool

	// For temporary files, we ignore any provided read/write keys.
	// Check fac items to see if a facilities item already exists with the given specification.
	var prevFacItem facItems.IFacilitiesItem
	var absSpec *uint
	var relSpec *int
	if fileSpecification.HasFileCycleSpecification() {
		absSpec = fileSpecification.FileCycleSpec.AbsoluteCycle
		relSpec = fileSpecification.FileCycleSpec.RelativeCycle
	}

	for _, facItem := range rce.FacilityItems {
		if facItems.FacilityItemMatches(facItem, fileSpecification.Qualifier, fileSpecification.Filename, absSpec, relSpec) {
			prevFacItem = facItem
			break
		}
	}

	if prevFacItem != nil {
		// Check whether caller is attempting to change the general file type
		// or to apply a non-conforming equipment type.
		equipSpecified := len(operandFields) >= 2 && len(operandFields[1][0]) > 0
		if equipSpecified {
			if usage == config.EquipmentUsageSectorAddressableMassStorage ||
				usage == config.EquipmentUsageWordAddressableMassStorage {
				if !prevFacItem.IsDisk() {
					facResult.PostMessage(kexec.FacStatusAttemptToChangeGenericType, nil)
					resultCode |= 0_600000_000000
					return
				}
			} else if usage == config.EquipmentUsageTape {
				if !prevFacItem.IsTape() {
					facResult.PostMessage(kexec.FacStatusAttemptToChangeGenericType, nil)
					resultCode |= 0_600000_000000
					return
				}
			}
		}

		// Check hold condition
		if isDiskFile {
			resultCode |= holdForDisk(rce, prevFacItem, sourceIsExecRequest, facResult)
		} else if isDiskUnit {
			resultCode != holdForDiskUnit(rce, prevFacItem, sourceIsExecRequest, facResult)
		} else if isTapeFile || isTapeUnit {
			resultCode |= holdForTape(rce, prevFacItem, sourceIsExecRequest, facResult)
		}

		if resultCode&0_400000_000000 != 0 {
			return
		}

		// If we're still okay
		if resultCode&0_400000_000000 == 0 {
			facResult.PostMessage(kexec.FacStatusFileAlreadyAssigned, nil)
			resultCode |= 0_100000_000000
		}

		return
	}

	// File is not already assigned - is it fixed, removable, tape?
	if isDiskFile {
		allowedOpts := uint64(kexec.IOption | kexec.TOption | kexec.ZOption)
		if !CheckIllegalOptions(rce, optionWord, allowedOpts, facResult, rce.IsExec()) {
			resultCode |= 0_600000_000000
			return
		}

		if !mgr.checkSubFields(operandFields, asgDiskFSIs) {
			facResult.PostMessage(kexec.FacStatusUndefinedFieldOrSubfield, nil)
			resultCode |= 0_600000_000000
			return
		}

		// TODO For Mass Storage, we need to honor allocation requests and set up a local allocation table
		// TODO For Removable, we may need to wait for pack mount(s) (but not here, maybe?)
	} else if isTapeFile {
		allowedOpts := uint64(kexec.EOption | kexec.FOption | kexec.HOption | kexec.IOption | kexec.JOption |
			kexec.LOption | kexec.MOption | kexec.NOption | kexec.OOption | kexec.ROption | kexec.SOption |
			kexec.TOption | kexec.VOption | kexec.WOption | kexec.XOption | kexec.ZOption)
		if !CheckIllegalOptions(rce, optionWord, allowedOpts, facResult, rce.IsExec()) {
			resultCode |= 0_600000_000000
			return
		}

		if !mgr.checkSubFields(operandFields, asgTapeFSIs) {
			facResult.PostMessage(kexec.FacStatusUndefinedFieldOrSubfield, nil)
			resultCode |= 0_600000_000000
			return
		}

		// TODO check duplicate reel numbers specified
		// TODO check for requested tape reel is already assigned by this run.
	}

	// Is filename not unique?
	for _, facItem := range rce.FacilityItems {
		if facItem.GetFilename() == fileSpecification.Filename {
			facResult.PostMessage(kexec.FacStatusFilenameNotUnique, nil)
			resultCode |= 0_004000_000000
			break
		}
	}

	return
}

func (mgr *FacilitiesManager) assignToBeCatalogedFile(
	rce *kexec.RunControlEntry,
	sourceIsExecRequest bool,
	fileSpecification *kexec.FileSpecification,
	optionWord uint64,
	operandFields [][]string,
	facResult *FacStatusResult,
) (resultCode uint64) {
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

// attachUnit attaches a unit (tape or disk) to a particular run.
// Used primarily for attaching fixed/removable disk units to the exec (or rarely to a run),
// or to attach tape units to a run (or sometimes to the exec).
// It is the caller's responsibility to properly detach a unit before attaching it somewhere else.
func (mgr *FacilitiesManager) attachUnit(
	rce *kexec.RunControlEntry,
	nodeId hardware.NodeIdentifier,
) {
	klog.LogTraceF("FacMgr", "attachUnit(nodeId=%v runid=%s)", nodeId, rce.RunId)
	mgr.attachments[nodeId] = rce
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
		fsIdent, result := mm.CreateFileSet(
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
		fileSetInfo, _ = mm.GetFileSetInfo(fsIdent)
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

// detachUnit detaches the unit from the attached run, if any.
// set force to true if this is not a result of the run doing a voluntary @FREE.
// if force is set, the run is aborted.
func (mgr *FacilitiesManager) detachUnit(
	nodeId hardware.NodeIdentifier,
	force bool,
) {
	klog.LogTraceF("FacMgr", "detachUnit(nodeId=%v)", nodeId)
	rce, ok := mgr.attachments[nodeId]
	if ok {
		if force {
			klog.LogInfoF("FacMgr", "detachUnit forcing removal from run %s", rce.RunId)
			// TODO abort (err?) the run unless it is the exec
		}
	}
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
		_, ok := mgr.attachments[diskAttr.Identifier]
		if ok {
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

	tapeAttr, ok := mgr.inventory.tapes[nodeId]
	if ok {
		rce, ok := mgr.attachments[tapeAttr.Identifier]
		if ok {
			// [* RUNID run-id REEL reel [RING|NORING] [POS [*]ffff[+|-][*]bbbbbb | POS LOST]]
			str += "* RUNID " + rce.RunId
			// TODO REEL reel
			// TODO RING | NORING in FS display for tape unit
			// TODO POS in FS display for tape unit
		}
	}

	return str
}
