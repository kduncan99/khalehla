// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/nodes"
	"khalehla/pkg"
	"log"
	"time"
)

// mfdServices contains code which provides directory-level services to all other exec code
// (but mostly facilities manager) such as assigning files, cataloging files, and general file allocation.

func (mgr *MFDManager) ChangeFileSetName(
	leadItem0Address kexec.MFDRelativeAddress,
	newQualifier string,
	newFilename string,
) error {
	// TODO
	return nil
}

// CreateFileSet creates lead items to establish an empty file set.
// When the MFD is in a normalized state, no empty file sets should exist -
// hence we expect the client to subsequently create a file as part of the file set.
// If we return MFDInternalError, the exec has been stopped
// If we return MFDFileNameConflict, a file set already exists with this file name.
func (mgr *MFDManager) CreateFileSet(
	qualifier string,
	filename string,
	projectId string,
	readKey string,
	writeKey string,
	fileType kexec.MFDFileType,
) (leadItem0Address kexec.MFDRelativeAddress, result MFDResult) {
	leadItem0Address = kexec.InvalidLink
	result = MFDSuccessful

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	_, ok := mgr.fileLeadItemLookupTable[qualifier][filename]
	if ok {
		result = MFDFileNameConflict
		return
	}

	leadItem0Address, leadItem0, err := mgr.allocateDirectorySector(kexec.InvalidLDAT)
	if err != nil {
		result = MFDInternalError
		return
	}

	leadItem0[0].SetW(uint64(kexec.InvalidLink))
	leadItem0[0].Or(0_500000_000000)

	pkg.FromStringToFieldata(qualifier, leadItem0[1:3])
	pkg.FromStringToFieldata(filename, leadItem0[3:5])
	pkg.FromStringToFieldata(projectId, leadItem0[5:7])
	leadItem0[7].FromStringToFieldata(readKey)
	leadItem0[8].FromStringToAscii(writeKey)
	leadItem0[9].SetS1(uint64(fileType))

	mgr.markDirectorySectorDirty(leadItem0Address)
	return
}

func (mgr *MFDManager) CreateFileCycle(
	leadItem0Address kexec.MFDRelativeAddress,
) (kexec.MFDRelativeAddress, error) {
	// TODO
	return 0, nil
}

func (mgr *MFDManager) GetFileInfo(
	qualifier string,
	filename string,
	absoluteCycle uint,
) (fi kexec.MFDFileInfo, mainItem0Address kexec.MFDRelativeAddress, err error) {
	// TODO
	return nil, 0, nil
}

// GetFileSetInfo returns a MFDFileSetInfo struct representing the file set
// and the file set's lead item 0 address.
// corresponding to the given qualifier and filename.
// If we return MFDInternalError, the exec has been stopped
// If we return MFDNotFound then there is no such file set
func (mgr *MFDManager) GetFileSetInfo(
	qualifier string,
	filename string,
) (fsInfo *kexec.MFDFileSetInfo, leadItem0Address kexec.MFDRelativeAddress, mfdResult MFDResult) {
	fsInfo = nil
	leadItem0Address = kexec.InvalidLink
	mfdResult = MFDSuccessful

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	leadItem0Address, ok := mgr.fileLeadItemLookupTable[qualifier][filename]
	if !ok {
		mfdResult = MFDNotFound
		return
	}

	leadItem0, err := mgr.getMFDSector(leadItem0Address)
	if err != nil {
		mfdResult = MFDInternalError
		return
	}

	leadItem1Address := kexec.MFDRelativeAddress(leadItem0[0].GetW())
	var leadItem1 []pkg.Word36
	if leadItem1Address != kexec.InvalidLink {
		leadItem1, err = mgr.getMFDSector(leadItem1Address)
		if err != nil {
			mfdResult = MFDInternalError
			return
		}
	}

	fsInfo = &kexec.MFDFileSetInfo{}
	fsInfo.PopulateFromLeadItems(leadItem0, leadItem1)
	return
}

// InitializeMassStorage handles MFD initialization for what is effectively a JK13 boot.
// If we return an error, we must previously stop the exec.
func (mgr *MFDManager) InitializeMassStorage() {
	// Get the list of disks from the node manager
	disks := make([]*nodeMgr.DiskDeviceInfo, 0)
	fm := mgr.exec.GetFacilitiesManager()
	nm := mgr.exec.GetNodeManager().(*nodeMgr.NodeManager)
	for _, dInfo := range nm.GetDeviceInfos() {
		if dInfo.GetNodeDeviceType() == nodes.NodeDeviceDisk {
			disks = append(disks, dInfo.(*nodeMgr.DiskDeviceInfo))
		}
	}

	// Check the labels on the disks so that we may segregate them into fixed and isRemovable lists.
	// Any problems at this point will lead us to DN the unit.
	// At this point, FacMgr should know about all the disks.
	fixedDisks := make(map[*nodeMgr.DiskDeviceInfo]*kexec.DiskAttributes)
	removableDisks := make(map[*nodeMgr.DiskDeviceInfo]*kexec.DiskAttributes)
	for _, ddInfo := range disks {
		nodeAttr, _ := fm.GetNodeAttributes(ddInfo.GetNodeIdentifier())
		if nodeAttr.GetFacNodeStatus() == kexec.FacNodeStatusUp {
			// Get the pack attributes from fac mgr
			diskAttr, ok := fm.GetDiskAttributes(ddInfo.GetNodeIdentifier())
			if !ok {
				mgr.exec.SendExecReadOnlyMessage("Internal configuration error", nil)
				mgr.exec.Stop(kexec.StopInitializationSystemConfigurationError)
				return
			}

			if diskAttr.PackLabelInfo == nil {
				msg := fmt.Sprintf("No valid label exists for pack on device %v", ddInfo.GetNodeName())
				mgr.exec.SendExecReadOnlyMessage(msg, nil)
				_ = fm.SetNodeStatus(ddInfo.GetNodeIdentifier(), kexec.FacNodeStatusDown)
				continue
			}

			// Read sector 1 of the initial directory track.
			// This is a little messy due to the potential of problematic block sizes.
			wordsPerBlock := uint64(diskAttr.PackLabelInfo.WordsPerRecord)
			dirTrackWordAddr := uint64(diskAttr.PackLabelInfo.FirstDirectoryTrackAddress)
			dirTrackBlockId := kexec.BlockId(dirTrackWordAddr / wordsPerBlock)
			if wordsPerBlock == 28 {
				dirTrackBlockId++
			}

			buf := make([]pkg.Word36, wordsPerBlock)
			pkt := nodes.NewDiskIoPacketRead(ddInfo.GetNodeIdentifier(), dirTrackBlockId, buf)
			nm.RouteIo(pkt)
			ioStat := pkt.GetIoStatus()
			if ioStat != nodes.IosComplete {
				msg := fmt.Sprintf("IO error reading directory track on device %v", ddInfo.GetNodeName())
				log.Printf("MFDMgr:%v", msg)
				mgr.exec.SendExecReadOnlyMessage(msg, nil)
				_ = fm.SetNodeStatus(ddInfo.GetNodeIdentifier(), kexec.FacNodeStatusDown)
				continue
			}

			var sector1 []pkg.Word36
			if wordsPerBlock == 28 {
				sector1 = buf
			} else {
				sector1 = buf[28:56]
			}

			// get the LDAT field from sector 1
			// If it is 0, it is a isRemovable pack
			// 0400000, it is an uninitialized fixed pack
			// anything else, it is a pre-used fixed pack which we're going to initialize
			ldat := sector1[5].GetH1()
			if ldat == 0 {
				removableDisks[ddInfo] = diskAttr
			} else {
				fixedDisks[ddInfo] = diskAttr
				diskAttr.IsFixed = true
			}
		}
	}

	err := mgr.initializeFixed(fixedDisks)
	if err != nil {
		return
	}

	// Make sure we have at least one fixed pack after the previous shenanigans
	if len(mgr.fixedPackDescriptors) == 0 {
		mgr.exec.SendExecReadOnlyMessage("No Fixed Disks - Cannot Continue Initialization", nil)
		mgr.exec.Stop(kexec.StopInitializationSystemConfigurationError)
		return
	}

	err = mgr.initializeRemovable(removableDisks)
	return
}

// RecoverMassStorage handles MFD recovery for what is NOT a JK13 boot.
// If we return an error, we must previously stop the exec.
func (mgr *MFDManager) RecoverMassStorage() {
	// TODO
	mgr.exec.SendExecReadOnlyMessage("MFD Recovery is not implemented", nil)
	mgr.exec.Stop(kexec.StopDirectoryErrors)
	return
}

func (mgr *MFDManager) SetFileCycleRange(
	leadItem0Address kexec.MFDRelativeAddress,
	cycleRange uint,
) error {
	// TODO
	return nil
}

func (mgr *MFDManager) SetFileToBeDeleted(
	leadItem0Address kexec.MFDRelativeAddress,
	absoluteCycle uint,
) error {
	// TODO
	return nil
}

// ----- mostly obsolete below here -----

// type MFDCatalogMassStorageFileRequest struct {
//	qualifier         string
//	filename          string
//	absoluteFileCycle *uint
//	relativeFileCycle *int
//	readKey           string
//	writeKey          string
//	projectId         string
//	accountId         string
//	granularity       pkg.Granularity
//	isRemovable       bool
//	wordAddressable   bool
//	saveOnCheckpoint  bool
//	equipment         string
//	isGuarded         bool
//	inhibitUnloadFlag bool
//	isPrivate         bool
//	isWriteOnly       bool
//	isReadOnly        bool
//	// for disk
//	reserve uint64
//	maximum uint64
//	packIds []string
//	// for tape
//	isTape            bool
//	density           uint
//	format            uint
//	features          uint
//	featuresExtension uint
//	mtapop            uint
//	ctlPool           string
//	reelNumbers       []string
// }
//
// type MFDDeleteFileRequest struct {
//	qualifier         string
//	filename          string
//	absoluteFileCycle *uint
//	relativeFileCycle *int
//	readKey           string
//	writeKey          string
// }

// populateNewLeadItem0 sets up a lead item sector 0 in the given buffer,
// assuming we are cataloging a new file, will have one cycle, and the absolute cycle is given to us.
// Implied is that there will be no sector 1 (since there aren't enough cycles to warrant it).
func populateNewLeadItem0(
	leadItem0 []pkg.Word36,
	qualifier string,
	filename string,
	projectId string,
	readKey string,
	writeKey string,
	fileType uint64, // 000=Fixed, 001=Tape, 040=Removable
	absoluteCycle uint64,
	guardedFlag bool,
	mainItem0Address uint64,
) {
	for wx := 0; wx < 28; wx++ {
		leadItem0[wx].SetW(0)
	}

	leadItem0[0].SetW(uint64(kexec.InvalidLink))
	leadItem0[0].Or(0_500000_000000)

	pkg.FromStringToFieldata(qualifier, leadItem0[1:3])
	pkg.FromStringToFieldata(filename, leadItem0[3:5])
	pkg.FromStringToFieldata(projectId, leadItem0[5:7])
	if len(readKey) > 0 {
		leadItem0[7].FromStringToFieldata(readKey)
	}
	if len(writeKey) > 0 {
		leadItem0[8].FromStringToAscii(writeKey)
	}

	leadItem0[9].SetS1(fileType)
	leadItem0[9].SetS2(1)  // number of cycles
	leadItem0[9].SetS3(31) // max range of cycles (default is 31)
	leadItem0[9].SetS4(1)  // current range
	leadItem0[9].SetT3(absoluteCycle)

	var statusBits uint64
	if guardedFlag {
		statusBits |= 01000
	}
	leadItem0[10].SetT1(statusBits)
	leadItem0[11].SetW(mainItem0Address)
}

func populateMassStorageMainItem0(
	mainItem0 []pkg.Word36,
	qualifier string,
	filename string,
	projectId string,
	readKey string,
	writeKey string,
	accountId string,
	leadItem0Address kexec.MFDRelativeAddress,
	mainItem1Address kexec.MFDRelativeAddress,
	saveOnCheckpoint bool,
	toBeCataloged bool, // for @ASG,C or @ASG,U
	isRemovable bool,
	isPosition bool,
	isWordAddressable bool,
	mnemonic string,
	isGuarded bool,
	inhibitUnload bool,
	isPrivate bool,
	isWriteOnly bool,
	isReadOnly bool,
	absoluteCycle uint64,
	reserve uint64,
	maximum uint64,
	packIds []string,
) {
	for wx := 0; wx < 28; wx++ {
		mainItem0[wx].SetW(0)
	}

	mainItem0[0].SetW(uint64(kexec.InvalidLink)) // no DAD table (yet, anyway)
	mainItem0[0].Or(0_200000_000000)

	pkg.FromStringToFieldata(qualifier, mainItem0[1:3])
	pkg.FromStringToFieldata(filename, mainItem0[3:5])
	pkg.FromStringToFieldata(projectId, mainItem0[5:7])
	pkg.FromStringToFieldata(accountId, mainItem0[7:9])
	mainItem0[11].SetW(uint64(leadItem0Address))
	mainItem0[11].SetS1(0) // disable flags

	var descriptorFlags uint64
	if saveOnCheckpoint {
		descriptorFlags |= 01000
	}
	if toBeCataloged {
		descriptorFlags |= 00100
	}
	if isRemovable {
		descriptorFlags |= 00010
	}
	mainItem0[12].SetT1(descriptorFlags)

	var pcharFlags uint64
	if isPosition {
		pcharFlags |= 040
	}
	if isWordAddressable {
		pcharFlags |= 010
	}
	mainItem0[13].SetW(uint64(mainItem1Address))
	mainItem0[13].SetS1(pcharFlags)

	mainItem0[14].FromStringToFieldata(mnemonic)

	var inhibitFlags uint64
	if isGuarded {
		inhibitFlags |= 040
	}
	if inhibitUnload {
		inhibitFlags |= 020
	}
	if isPrivate {
		inhibitFlags |= 010
	}
	if isWriteOnly {
		inhibitFlags |= 002
	}
	if isReadOnly {
		inhibitFlags |= 001
	}
	mainItem0[17].SetH1(inhibitFlags)
	mainItem0[17].SetT3(absoluteCycle)

	swTimeNow := kexec.GetSWTimeFromSystemTime(time.Now())
	mainItem0[19].SetW(swTimeNow)
	mainItem0[20].SetH1(reserve)
	mainItem0[21].SetH1(maximum)

	if isRemovable {
		var rKey pkg.Word36
		if len(readKey) > 0 {
			rKey.FromStringToFieldata(readKey)
		}
		var wKey pkg.Word36
		if len(writeKey) > 0 {
			wKey.FromStringToAscii(writeKey)
		}

		mainItem0[24].SetH1(rKey.GetH1())
		mainItem0[25].SetH1(rKey.GetH2())
		mainItem0[26].SetH1(wKey.GetH1())
		mainItem0[27].SetH1(wKey.GetH2())
	} else {
		// initially selected LDAT and optional device placement flag
		// TODO if there is at least one pack-id, then go find its LDAT and use that,
		//  and mask in 0_400000_000000 to indicate device placement.
		var ldat uint64
		if len(packIds) > 0 {

		} else {
			ldat = uint64(getLDATIndexFromMFDAddress(leadItem0Address))
		}
		mainItem0[27].SetH1(ldat)
	}
}

func populateFixedMainItem1(
	mainItem1 []pkg.Word36,
	qualifier string,
	filename string,
	mainItem0Address kexec.MFDRelativeAddress,
	absoluteCycle uint64,
	packIds []string,
) {
	for wx := 0; wx < 28; wx++ {
		mainItem1[wx].SetW(0)
	}

	mainItem1[0].SetW(uint64(kexec.InvalidLink)) // no sector 2 (yet, anyway)
	pkg.FromStringToFieldata(qualifier, mainItem1[1:3])
	pkg.FromStringToFieldata(filename, mainItem1[3:5])
	pkg.FromStringToFieldata("*No.1*", mainItem1[5:6])
	mainItem1[6].SetW(uint64(mainItem0Address))
	mainItem1[7].SetT3(absoluteCycle)

	// TODO note that for >5 pack entries, we need additional main item sectors
	//  one per 10 additional packs beyond 5
	mix := 18
	limit := len(packIds)
	if limit > 5 {
		limit = 5
	}
	for dpx := 0; dpx < limit; dpx++ {
		mainItem1[mix].FromStringToFieldata(packIds[dpx])
		mix += 2
	}
}

func populateRemovableMainItem1(
	mainItem1 []pkg.Word36,
	mainItem0Address kexec.MFDRelativeAddress,
	absoluteCycle uint64,
	packIds []string,
) {
	for wx := 0; wx < 28; wx++ {
		mainItem1[wx].SetW(0)
	}

	mainItem1[0].SetW(uint64(kexec.InvalidLink)) // no sector 2 (yet, anyway)
	mainItem1[6].SetW(uint64(mainItem0Address))
	mainItem1[7].SetT3(absoluteCycle)
	mainItem1[17].SetT3(uint64(len(packIds)))

	// TODO note that for >5 pack entries, we need additional main item sectors
	//  one per 10 additional packs beyond 5
	mix := 18
	limit := len(packIds)
	if limit > 5 {
		limit = 5
	}
	for dpx := 0; dpx < limit; dpx++ {
		mainItem1[mix].FromStringToFieldata(packIds[dpx])
		// TODO for isRemovable, we need the main item address for this file on that pack
		mix += 2
	}
}

func populateTapeMainItem0(
	mainItem0 []pkg.Word36,
	qualifier string,
	filename string,
	projectId string,
	accountId string,
	reelTable0Address kexec.MFDRelativeAddress,
	leadItem0Address kexec.MFDRelativeAddress,
	mainItem1Address kexec.MFDRelativeAddress,
	toBeCataloged bool, // for @ASG,C or @ASG,U
	isGuarded bool,
	isPrivate bool,
	isWriteOnly bool,
	isReadOnly bool,
	absoluteCycle uint64,
	density uint,
	format uint,
	features uint,
	featuresExtension uint,
	mtapop uint,
	ctlPool string,
) {
	for wx := 0; wx < 28; wx++ {
		mainItem0[wx].SetW(0)
	}

	mainItem0[0].SetW(uint64(reelTable0Address))
	mainItem0[0].Or(0_200000_000000)
	pkg.FromStringToFieldata(qualifier, mainItem0[1:3])
	pkg.FromStringToFieldata(filename, mainItem0[3:5])
	pkg.FromStringToFieldata(projectId, mainItem0[5:7])
	pkg.FromStringToFieldata(accountId, mainItem0[7:9])

	// TODO
}

func populateTapeMainItem1(
	mainItem1 []pkg.Word36,
	qualifier string,
	filename string,
	mainItem0Address kexec.MFDRelativeAddress,
	absoluteCycle uint64,
) {
	for wx := 0; wx < 28; wx++ {
		mainItem1[wx].SetW(0)
	}

	mainItem1[0].SetW(uint64(kexec.InvalidLink)) // no sector 2 (yet, anyway)
	pkg.FromStringToFieldata(qualifier, mainItem1[1:3])
	pkg.FromStringToFieldata(filename, mainItem1[3:5])
	pkg.FromStringToFieldata("*No.1*", mainItem1[5:6])
	mainItem1[6].SetW(uint64(mainItem0Address))
	mainItem1[7].SetT3(absoluteCycle)
}

// // CatalogFile attempts to catalog a file on mass storage according to the given parameters.
// // We really only expect to be invoked via fac mgr, and only for @CAT of word and sector addressable disk files.
// // If we return err, we've stopped the exec
// func (mgr *MFDManager) CatalogFile(parameters *MFDCatalogMassStorageFileRequest) (*facilitiesMgr.FacResult, error) {
//
//	facResult := facilitiesMgr.NewFacResult()
//
//	mgr.mutex.Lock()
//	defer mgr.mutex.Unlock()
//
//	// get lead item(s) for the fileset if the fileset exists
//	_ /*leadAddr0*/, fileSetExists := mgr.fileLeadItemLookupTable[parameters.qualifier][parameters.filename]
//
//	if !fileSetExists {
//		// Create a lead and main items and mark them dirty
//		leadAddr0, leadItem0, err := mgr.allocateDirectorySector(pkg.InvalidLDAT)
//		if err != nil {
//			return nil, err
//		}
//
//		mainAddr0, mainItem0, err := mgr.allocateDirectorySector(pkg.InvalidLDAT)
//		if err != nil {
//			return nil, err
//		}
//
//		mainAddr1, mainItem1, err := mgr.allocateDirectorySector(pkg.InvalidLDAT)
//		if err != nil {
//			return nil, err
//		}
//
//		mgr.markDirectorySectorDirty(leadAddr0)
//		mgr.markDirectorySectorDirty(mainAddr0)
//		mgr.markDirectorySectorDirty(mainAddr1)
//
//		var effectiveAbsolute uint64
//		if parameters.relativeFileCycle != nil {
//			effectiveAbsolute = 1
//		} else {
//			effectiveAbsolute = uint64(*parameters.absoluteFileCycle)
//		}
//
//		fileType := uint64(0)
//		if parameters.isTape {
//			fileType = 001
//		} else if parameters.isRemovable {
//			fileType = 040
//		}
//
//		populateNewLeadItem0(
//			leadItem0, parameters.qualifier, parameters.filename, parameters.projectId,
//			parameters.readKey, parameters.writeKey, fileType,
//			effectiveAbsolute, parameters.isGuarded, uint64(mainAddr0))
//
//		equip := parameters.equipment
//		if len(equip) == 0 {
//			if parameters.wordAddressable {
//				equip = "D"
//			} else {
//				equip = "F"
//			}
//		}
//
//		if !parameters.isTape {
//			populateMassStorageMainItem0(mainItem0, parameters.qualifier, parameters.filename, parameters.projectId,
//				parameters.readKey, parameters.writeKey, parameters.accountId, leadAddr0, mainAddr1,
//				parameters.saveOnCheckpoint, false, parameters.isRemovable,
//				parameters.granularity == pkg.PositionGranularity, parameters.wordAddressable,
//				equip, parameters.isGuarded, parameters.inhibitUnloadFlag,
//				parameters.isPrivate, parameters.isWriteOnly, parameters.isReadOnly, effectiveAbsolute,
//				parameters.reserve, parameters.maximum, parameters.packIds)
//
//			if !parameters.isRemovable {
//				populateFixedMainItem1(mainItem1, parameters.qualifier, parameters.filename,
//					mainAddr0, effectiveAbsolute, parameters.packIds)
//				// TODO possibly more
//			} else {
//				populateRemovableMainItem1(mainItem1, mainAddr0, effectiveAbsolute, parameters.packIds)
//				// TODO possibly more
//			}
//		} else {
//			reelAddr0 := pkg.InvalidLink // TODO REEL table address or InvalidLink
//			populateTapeMainItem0(mainItem0, parameters.qualifier, parameters.filename, parameters.projectId,
//				parameters.accountId, reelAddr0, leadAddr0, mainAddr1, false,
//				parameters.isGuarded, parameters.isPrivate, parameters.isWriteOnly, parameters.isReadOnly,
//				effectiveAbsolute, parameters.density, parameters.format, parameters.features,
//				parameters.featuresExtension, parameters.mtapop, parameters.ctlPool)
//			populateTapeMainItem1(mainItem1, parameters.qualifier, parameters.filename, mainAddr0, effectiveAbsolute)
//			// TODO tape file reel table(s)
//		}
//	} else {
//		// TODO check read/write keys
//		//	 I suspect we have to meet write key verification, but probably not read key
//
//		// TODO check file cycles - zero? +1? absolute? negative relative cycles are not allowed
//		// 	Absolute cycles are 1 to 999, Relative cycles are -{n}, 0, or +1
//		//  if a +1 is currently assigned, we cannot catalog anything in this file set (a +1 can never exist if it is not assigned)
//		//  if we are +0, and a fileset already exists, this is an error
//		//  if we try to catalog a cycle which would delete the lowest-number cycle
//		//  	and it is assigned, we have an f-cycle conflict
//		//		and there are more than one such cycles (see formula below) we have an f-cycle conflict
//		//		and the cycle to be deleted has a write key other than what we specificed, we have an f-cycle conflict
//		//  For a new cycle to be created, its absolute F-cycle number must be within the following range:
//		// 		(x-w) < z ≤ (x-y+w+1) where:
//		// 		x is T3 of word 9 of the lead item (cycle number of latest F-cycle).
//		// 		w is S3 of word 9 of the lead item (maximum number of F-cycles).
//		// 		z is the absolute F-cycle number requested.
//		// 		y is S4 of word 9 of the lead item (current range of F-cycles)
//		//			this is the highest-absolute-f-cycle - lowest-absolute-f-cycle + 1
//		//			*or* highest-absolute-f-cycle + 1000 - lowest-absolute-f-cycle
//		//
//		// e.g.: we have cycles 10, 11, 20, 21, 30, 31 and max cycles is 25
//		// so x-w = 31 - 25 = 6
//		// x-y+w+1 = 31 - (31 - 10 + 1) + 25 + 1 = 35
//		// thus 6 < z <= 35
//
//		// leadItem0, _ := mgr.getMFDSector(leadAddr0)
//		// var leadAddr1 = pkg.MFDRelativeAddress(leadItem0[0].GetW())
//		// var leadItem1 []pkg.Word36
//		// if leadAddr1 != pkg.InvalidLink {
//		// 	leadItem1, _ = mgr.getMFDSector(leadAddr1)
//		// }
//		//
//		// maxRange := leadItem0[9].GetS3()
//		// currentRange := leadItem0[9].GetS4()
//		// highestAbsolute := leadItem0[9].GetT3()
//
//		// TODO make sure the requested file cycle is valid
//
//		var effectiveEquipmentType = parameters.equipment
//		if len(effectiveEquipmentType) == 0 {
//			// Find the main item for the highest f-cycle.
//			// To do this, we have to walk the list in the lead item(s).
//			// TODO
//		}
//
//		// TODO re-use whatever we can from the above code...
//	}
//
//	facResult.PostMessage(facilitiesMgr.FacStatusComplete, []string{"CAT"})
//	return facResult, nil
// }
