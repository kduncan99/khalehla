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

// CreateFixedFileCycle creates a new file cycle within the file set specified by fsIdentifier.
// fcSpecification is nil if no file cycle was specified.
// Returns
//
//		FileCycleIdentifier of the newly-created file cycle if successful, and
//		MFDResult values which may include
//	     MFDSuccessful if things went well
//	     MFDInternalError if something is badly wrong, and we've stopped the exec
//	     MFDAlreadyExists if user does not specify a cycle, and a file cycle already exists
//		     MFDInvalidRelativeFileCycle if the caller specified a negative relative file cycle
//	      MFDInvalidAbsoluteFileCycle
//	         any file cycle out of range
//	         an absolute file cycle which conflicts with an existing cycle
//	     MFDPlusOneCycleExists if caller specifies +1 and a +1 already exists for the file set
//	     MFDDropOldestCycleRequired is returned if everything else would be fine if the oldest file cycle did not exist
func (mgr *MFDManager) CreateFixedFileCycle(
	fsIdentifier FileSetIdentifier,
	fcSpecification *kexec.FileCycleSpecification,
	accountId string,
	assignMnemonic string,
	descriptorFlags DescriptorFlags,
	pcharFlags PCHARFlags,
	inhibitFlags InhibitFlags,
	initialReserve uint64,
	maxGranules uint64,
	unitSelection UnitSelectionIndicators, // TODO should we/can we do anything with this?
	diskPacks []DiskPackEntry,
) (fcIdentifier FileCycleIdentifier, result MFDResult) {
	result = MFDSuccessful

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	leadItem0Addr := kexec.MFDRelativeAddress(fsIdentifier)
	leadItem0, err := mgr.getMFDSector(leadItem0Addr)
	if err != nil {
		result = MFDInternalError
		return
	}

	leadItem1Addr := kexec.MFDRelativeAddress(leadItem0[0].GetW())
	var leadItem1 []pkg.Word36
	if leadItem1Addr&0_400000_000000 == 0 {
		leadItem1, err = mgr.getMFDSector(leadItem1Addr)
		if err != nil {
			result = MFDInternalError
			return
		}
	}

	// get a FileSetInfo for convenience
	fsInfo := &FileSetInfo{}
	fsInfo.FileSetIdentifier = fsIdentifier
	fsInfo.populateFromLeadItems(leadItem0, leadItem1)

	absCycle, cycIndex, shift, newRange, plusOne, result := mgr.checkCycle(fcSpecification, fsInfo)
	if result != MFDSuccessful {
		return
	}

	// Do we need to allocate a lead item sector 1?
	if leadItem1 == nil && newRange > 28-(11+fsInfo.NumberOfSecurityWords) {
		leadItem1Addr, leadItem1, err = mgr.allocateLeadItem1(leadItem0Addr, leadItem0)
		if err != nil {
			result = MFDInternalError
			return
		}
	}

	// Do we need to shift links?
	if shift > 0 {
		adjustLeadItemLinks(leadItem0, leadItem1, shift)
		mgr.markDirectorySectorDirty(leadItem0Addr)
		mgr.markDirectorySectorDirty(leadItem1Addr)
	}

	if plusOne {
		leadItem0[012].Or(0_200000_000000)
	}

	// Create necessary main items for the new file cycle
	preferredLDAT := getLDATIndexFromMFDAddress(leadItem0Addr)
	mainItem0Addr, mainItem0, err := mgr.allocateDirectorySector(preferredLDAT)
	if err != nil {
		result = MFDInternalError
		return
	}
	mainItem1Addr, mainItem1, err := mgr.allocateDirectorySector(preferredLDAT)
	if err != nil {
		result = MFDInternalError
		return
	}

	packNames := make([]string, 0)
	if diskPacks != nil {
		for _, dp := range diskPacks {
			packNames = append(packNames, dp.PackName)
		}
	}

	populateMassStorageMainItem0(
		mainItem0,
		leadItem0Addr,
		mainItem1Addr,
		fsInfo.Qualifier,
		fsInfo.Filename,
		uint64(absCycle),
		fsInfo.ReadKey,
		fsInfo.WriteKey,
		fsInfo.ProjectId,
		accountId,
		assignMnemonic,
		descriptorFlags,
		pcharFlags,
		inhibitFlags,
		false,
		initialReserve,
		maxGranules,
		packNames)

	populateFixedMainItem1(
		mainItem1,
		fsInfo.Qualifier,
		fsInfo.Filename,
		mainItem0Addr,
		uint64(absCycle),
		packNames)

	// Link the new file cycle into the lead item
	lw := getLeadItemLinkWord(leadItem0, leadItem1, cycIndex)
	lw.SetW(uint64(mainItem0Addr))

	mgr.markDirectorySectorDirty(leadItem0Addr)
	if leadItem1 != nil {
		mgr.markDirectorySectorDirty(leadItem1Addr)
	}

	fsIdentifier = FileSetIdentifier(mainItem0Addr)
	return
}

// DropFileCycle effectively deletes a file cycle.
// It also updates the main item as necessary. It will *not* delete the main items if it is the
// last file cycle - the caller *must* do that.
// This is the service wrapper which locks the manager before going to the core function.
// Caller should NOT invoke this on any file which is still assigned.
func (mgr *MFDManager) DropFileCycle(
	fcIdentifier FileCycleIdentifier,
) MFDResult {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// TODO

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
