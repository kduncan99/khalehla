// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/nodeMgr"
	"khalehla/pkg"
	"log"
	"strings"
	"time"
)

// mfdServices contains code which provides directory-level services to all other exec code
// (but mostly facilities manager) such as assigning files, cataloging files, and general file allocation.
// It is *PRIMARILY* a service for facilities - it should be rare that any other component uses it
// except via fac manager.

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
	fileType FileType,
	qualifier string,
	filename string,
	projectId string,
	readKey string,
	writeKey string,
) (fsIdentifier FileSetIdentifier, result MFDResult) {
	leadItem0Address := kexec.InvalidLink
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
	fsIdentifier = FileSetIdentifier(leadItem0Address)
	return
}

// CreateFixedDiskFileCycle creates a new file cycle within the file set specified by fsIdentifier.
// Either absoluteCycle or relativeCycle must be specified (but not both).
func (mgr *MFDManager) CreateFixedDiskFileCycle(
	fsIdentifier FileSetIdentifier,
	absoluteCycle *uint,
	relativeCycle *int,
	accountId string,
	assignMnemonic string,
	descriptorFlags DescriptorFlags,
	fileFlags FileFlags,
	pcharFlags PCHARFlags,
	inhibitFlags InhibitFlags,
	initialReserve uint64,
	maxGranules uint64,
	unitSelection UnitSelectionIndicators,
	diskPacks []DiskPackEntry,
) (fcIdentifier FileCycleIdentifier, result MFDResult) {
	mainItem0Address := kexec.InvalidLink
	result = MFDSuccessful

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	/*
		DescriptorFlags          DescriptorFlags
		FileFlags                FileFlags ?
		PCHARFlags               PCHARFlags
		AssignMnemonic           string
		InhibitFlags             InhibitFlags
		AbsoluteFCycle           uint64
		TimeCataloged            uint64
		InitialGranulesReserved  uint64
		MaxGranules              uint64
		HighestGranuleAssigned   uint64
		UnitSelectionIndicators  UnitSelectionIndicators
		DiskPackEntries          []DiskPackEntry
		FileAllocations          []FileAllocation
	*/
	// TODO

	fsIdentifier = FileSetIdentifier(mainItem0Address)
	return
}

// DropFileCycle effectively deletes a file cycle.
// It does *not* delete the fileset, even if the fileset is now empty.
// This is the service wrapper which locks the manager before going to the core function.
// Caller should NOT invoke this on any file which is still assigned.
func (mgr *MFDManager) DropFileCycle(
	fcIdentifier FileCycleIdentifier,
) MFDResult {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	return mgr.dropFileCycle(kexec.MFDRelativeAddress(fcIdentifier))
}

// DropFileSet deletes an entire (possibly -- actually, hopefully -- empty) file set.
// Even though we would rather the caller drop the various cycles individually,
// we'll honor the request on a non-empty file set as well.
// Caller should NOT invoke this on any fileset which still has a file cycle assigned.
func (mgr *MFDManager) DropFileSet(
	fsIdentifier FileSetIdentifier,
) MFDResult {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	leadItem0Addr := kexec.MFDRelativeAddress(fsIdentifier)
	leadItem0, err := mgr.getMFDSector(leadItem0Addr)
	if err != nil {
		return MFDNotFound
	}

	_ = mgr.markDirectorySectorUnallocated(leadItem0Addr)
	leadItem1Addr := kexec.MFDRelativeAddress(leadItem0[0].GetW())
	var leadItem1 []pkg.Word36
	if leadItem1Addr&0_400000_000000 == 0 {
		leadItem1Addr &= 0_007777_777777
		leadItem1, err = mgr.getMFDSector(leadItem1Addr)
		if err != nil {
			return MFDNotFound
		}
		_ = mgr.markDirectorySectorUnallocated(leadItem1Addr)
	}

	links := leadItem0[013:]
	if leadItem1 != nil {
		links = append(links, leadItem1[1:]...)
	}
	currentRange := leadItem0[011].GetS4()
	links = links[:currentRange]

	for lx := 0; lx < len(links); lx++ {
		if links[lx] != 0 {
			mgr.DropFileCycle(FileCycleIdentifier(links[lx] & 0_007777_777777))
		}
	}

	qualifier := strings.Trim(leadItem0[0].ToStringAsFieldata()+leadItem0[1].ToStringAsFieldata(), " ")
	filename := strings.Trim(leadItem0[2].ToStringAsFieldata()+leadItem0[3].ToStringAsFieldata(), " ")
	delete(mgr.fileLeadItemLookupTable[qualifier], filename)

	return MFDSuccessful
}

// GetFileCycleInfo returns a FileCycleInfo struct representing the file cycle corresponding to the given
// file cycle identifier.
// If we return MFDInternalError, the exec has been stopped
// If we return MFDNotFound then there is no such file cycle
func (mgr *MFDManager) GetFileCycleInfo(
	fcIdentifier FileCycleIdentifier,
) (fcInfo *FixedFileCycleInfo, mfdResult MFDResult) {
	fcInfo = nil
	mfdResult = MFDSuccessful

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mainItem0, err := mgr.getMFDSector(kexec.MFDRelativeAddress(fcIdentifier))
	if err != nil {
		mfdResult = MFDNotFound
		return
	}

	mainItem1Address := kexec.MFDRelativeAddress(mainItem0[015].GetW()) & 0_007777_777777
	mainItem1, err := mgr.getMFDSector(mainItem1Address)
	if err != nil {
		mfdResult = MFDInternalError
		return
	}

	mainItems := make([][]pkg.Word36, 2)
	mainItems[0] = mainItem0
	mainItems[1] = mainItem1
	link := mainItem1[0].GetW()
	for link&0_400000_000000 == 0 {
		mi, err := mgr.getMFDSector(kexec.MFDRelativeAddress(link & 0_007777_777777))
		if err != nil {
			mfdResult = MFDInternalError
			return
		}
		mainItems = append(mainItems, mi)
	}

	fcInfo = &FixedFileCycleInfo{}
	fcInfo.setFileCycleIdentifier(fcIdentifier)
	fcInfo.setFileSetIdentifier(FileSetIdentifier(mainItem0[013] & 0_007777_777777))
	fcInfo.populateFromMainItems(mainItems)

	return fcInfo, MFDSuccessful
}

func (mgr *MFDManager) GetFileSetIdentifier(
	qualifier string,
	filename string,
) (fsi FileSetIdentifier, mfdResult MFDResult) {
	leadItem0Address, ok := mgr.fileLeadItemLookupTable[qualifier][filename]
	if !ok {
		return 0, MFDNotFound
	} else {
		return FileSetIdentifier(leadItem0Address), MFDSuccessful
	}
}

// GetFileSetInfo returns a FileSetInfo struct representing the file set corresponding to the given
// file set identifier.
// If we return MFDInternalError, the exec has been stopped
// If we return MFDNotFound then there is no such file set
func (mgr *MFDManager) GetFileSetInfo(
	fsIdentifier FileSetIdentifier,
) (fsInfo *FileSetInfo, mfdResult MFDResult) {
	fsInfo = nil
	mfdResult = MFDSuccessful

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	leadItem0, err := mgr.getMFDSector(kexec.MFDRelativeAddress(fsIdentifier))
	if err != nil {
		mfdResult = MFDNotFound
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

	fsInfo = &FileSetInfo{}
	fsInfo.FileSetIdentifier = fsIdentifier
	fsInfo.populateFromLeadItems(leadItem0, leadItem1)
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
		if dInfo.GetNodeDeviceType() == kexec.NodeDeviceDisk {
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
			pkt := nodeMgr.NewDiskIoPacketRead(ddInfo.GetNodeIdentifier(), dirTrackBlockId, buf)
			nm.RouteIo(pkt)
			ioStat := pkt.GetIoStatus()
			if ioStat != nodeMgr.IosComplete {
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
		mgr.exec.SendExecReadOnlyMessage("No Fixed disks - Cannot Continue Initialization", nil)
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

//func populateRemovableMainItem1(
//	mainItem1 []pkg.Word36,
//	mainItem0Address kexec.MFDRelativeAddress,
//	absoluteCycle uint64,
//	packIds []string,
//) {
//	for wx := 0; wx < 28; wx++ {
//		mainItem1[wx].SetW(0)
//	}
//
//	mainItem1[0].SetW(uint64(kexec.InvalidLink)) // no sector 2 (yet, anyway)
//	mainItem1[6].SetW(uint64(mainItem0Address))
//	mainItem1[7].SetT3(absoluteCycle)
//	mainItem1[17].SetT3(uint64(len(packIds)))
//
//	// TODO note that for >5 pack entries, we need additional main item sectors
//	//  one per 10 additional packs beyond 5
//	mix := 18
//	limit := len(packIds)
//	if limit > 5 {
//		limit = 5
//	}
//	for dpx := 0; dpx < limit; dpx++ {
//		mainItem1[mix].FromStringToFieldata(packIds[dpx])
//		// TODO for isRemovable, we need the main item address for this file on that pack
//		mix += 2
//	}
//}

//func populateTapeMainItem0(
//	mainItem0 []pkg.Word36,
//	qualifier string,
//	filename string,
//	projectId string,
//	accountId string,
//	reelTable0Address kexec.MFDRelativeAddress,
//	leadItem0Address kexec.MFDRelativeAddress,
//	mainItem1Address kexec.MFDRelativeAddress,
//	toBeCataloged bool, // for @ASG,C or @ASG,U
//	isGuarded bool,
//	isPrivate bool,
//	isWriteOnly bool,
//	isReadOnly bool,
//	absoluteCycle uint64,
//	density uint,
//	format uint,
//	features uint,
//	featuresExtension uint,
//	mtapop uint,
//	ctlPool string,
//) {
//	for wx := 0; wx < 28; wx++ {
//		mainItem0[wx].SetW(0)
//	}
//
//	mainItem0[0].SetW(uint64(reelTable0Address))
//	mainItem0[0].Or(0_200000_000000)
//	pkg.FromStringToFieldata(qualifier, mainItem0[1:3])
//	pkg.FromStringToFieldata(filename, mainItem0[3:5])
//	pkg.FromStringToFieldata(projectId, mainItem0[5:7])
//	pkg.FromStringToFieldata(accountId, mainItem0[7:9])
//
//	// TODO
//}

//func populateTapeMainItem1(
//	mainItem1 []pkg.Word36,
//	qualifier string,
//	filename string,
//	mainItem0Address kexec.MFDRelativeAddress,
//	absoluteCycle uint64,
//) {
//	for wx := 0; wx < 28; wx++ {
//		mainItem1[wx].SetW(0)
//	}
//
//	mainItem1[0].SetW(uint64(kexec.InvalidLink)) // no sector 2 (yet, anyway)
//	pkg.FromStringToFieldata(qualifier, mainItem1[1:3])
//	pkg.FromStringToFieldata(filename, mainItem1[3:5])
//	pkg.FromStringToFieldata("*No.1*", mainItem1[5:6])
//	mainItem1[6].SetW(uint64(mainItem0Address))
//	mainItem1[7].SetT3(absoluteCycle)
//}
