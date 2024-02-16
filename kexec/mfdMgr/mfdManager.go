// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/types"
	"khalehla/pkg"
	"log"
	"sync"
	"time"
)

// MFDManager is essentially the file system
//
// Fixed Pool things
//
// MFD-Relative addresses refer to a particular MFD-relative sector.
// Such addresses are formatted (octally) as 00llllttttss
// where llll is the LDAT index of the pack which contains the track in which the sector resides,
// tttt is the track number relative to the LDAT index, and ss is the sector number within the track.
// This address does not relate to the location on the pack where the track resides (other than the LDAT portion).
// The first directory track on any fixed pack is locatable via the VOL1 label for the pack, and all
// subsequent tracks exist on a forward-linked list.
//
// Lookup Table
// Since we have the entire MFD in core, we do not persist lookup table entries in the MFD.
// However, we do have a lookup table. We use the language map function to implement the thing,
// with the key being the concatenation of qualifier ':' filename -> the MFDRelativeAddress of the main item
// for the file set.
//
// free track list
// Instead of maintaining HMBT (which is no longer needed) and SMBT (which is annoying) we manage a free space
// list for each individual pack.
//
// For reference:
//  U001 search item
//  U010 lead item 0
//  U100 main item 0
//	U000 lead item 1 (U==0), main item sector {n}, DAD table

type fixedPackDescriptor struct {
	deviceId       types.DeviceIdentifier
	packAttributes *types.PackAttributes // created and maintained by facMgr
	wordsPerBlock  types.PrepFactor
	canAllocate    bool // true if pack is UP, false if it is SU
	packMask       uint64
	fixedFeeSpace  *packFreeSpaceTable
}

func newFixedPackDescriptor(
	deviceId types.DeviceIdentifier,
	packAttrs *types.PackAttributes,
	allocatable bool) *fixedPackDescriptor {
	recordLength := packAttrs.Label[04].GetH2()
	trackCount := packAttrs.Label[016].GetW()
	return &fixedPackDescriptor{
		deviceId:       deviceId,
		packAttributes: packAttrs,
		wordsPerBlock:  types.PrepFactor(recordLength),
		canAllocate:    allocatable,
		packMask:       (recordLength / 28) - 1,
		fixedFeeSpace:  newPackFreeSpaceTable(types.TrackCount(trackCount)),
	}
}

type MFDManager struct {
	exec                         types.IExec
	mutex                        sync.Mutex
	isInitialized                bool
	terminateThread              bool
	threadStarted                bool
	threadStopped                bool
	msInitialize                 bool
	mfdFileMainItem0Address      types.MFDRelativeAddress
	cachedTracks                 map[types.MFDRelativeAddress][]pkg.Word36 // key is MFD addr of first sector in track
	dirtyBlocks                  map[types.MFDRelativeAddress]bool         // 00llllttttss address of block containing the dirty sector
	deviceReadyNotificationQueue map[types.DeviceIdentifier]bool
	fixedPackDescriptors         map[types.LDATIndex]*fixedPackDescriptor
	fixedLookupTable             map[string]map[string]types.MFDRelativeAddress
	fileAllocations              *fileAllocationTable
}

func NewMFDManager(exec types.IExec) *MFDManager {
	return &MFDManager{
		exec:                         exec,
		cachedTracks:                 make(map[types.MFDRelativeAddress][]pkg.Word36),
		dirtyBlocks:                  make(map[types.MFDRelativeAddress]bool),
		deviceReadyNotificationQueue: make(map[types.DeviceIdentifier]bool),
		fixedPackDescriptors:         make(map[types.LDATIndex]*fixedPackDescriptor),
		fixedLookupTable:             make(map[string]map[string]types.MFDRelativeAddress),
		fileAllocations:              newFileAllocationTable(),
	}
}

// CloseManager is invoked when the exec is stopping
func (mgr *MFDManager) CloseManager() {
	mgr.threadStop()
	mgr.isInitialized = false
}

func (mgr *MFDManager) InitializeManager() error {
	var err error
	if mgr.msInitialize {
		err = mgr.initializeMassStorage()
	} else {
		err = mgr.recoverMassStorage()
	}

	if err != nil {
		log.Println("MFDMgr:Cannot continue boot")
		return err
	}

	mgr.threadStart()
	mgr.isInitialized = true
	return nil
}

func (mgr *MFDManager) IsInitialized() bool {
	return mgr.isInitialized
}

// ResetManager clears out any artifacts left over by a previous exec session,
// and prepares the console for normal operations
func (mgr *MFDManager) ResetManager() error {
	mgr.threadStop()
	mgr.threadStart()
	mgr.isInitialized = true
	return nil
}

// ----------------------------------------------------------------

func (mgr *MFDManager) NotifyDeviceReady(deviceInfo types.DeviceInfo, isReady bool) {
	// post it, and let the tread deal with it later
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.deviceReadyNotificationQueue[deviceInfo.GetDeviceIdentifier()] = isReady
}

// SetMSInitialize sets or clears the flag which indicates whether to initialze
// mass-storage upon initialization. Invoke this before calling Initialize.
func (mgr *MFDManager) SetMSInitialize(flag bool) {
	mgr.msInitialize = flag
}

// ----------------------------------------------------------------

func composeMFDAddress(ldatIndex types.LDATIndex, trackId types.MFDTrackId, sectorId types.MFDSectorId) types.MFDRelativeAddress {
	return types.MFDRelativeAddress(uint64(ldatIndex&07777)<<18 | uint64(trackId&07777)<<6 | uint64(sectorId&077))
}

// establishNewMFDTrack puts structures in place for a new cached MFD track.
func (mgr *MFDManager) establishNewMFDTrack(ldatIndex types.LDATIndex, mfdTrackId types.MFDTrackId) {
	mfdAddr := composeMFDAddress(ldatIndex, mfdTrackId, 0)
	mgr.cachedTracks[mfdAddr] = make([]pkg.Word36, 1792)
}

func getLDATIndexFromMFDAddress(address types.MFDRelativeAddress) types.LDATIndex {
	return types.LDATIndex(address>>18) & 07777
}

func getMFDTrackIdFromMFDAddress(address types.MFDRelativeAddress) types.MFDTrackId {
	return types.MFDTrackId(address>>6) & 07777
}

func getMFDSectorIdFromMFDAddress(address types.MFDRelativeAddress) types.MFDSectorId {
	return types.MFDSectorId(address & 077)
}

// getMFDAddressForBlock takes a given MFD-relative sector address and normalizes it to
// the first sector in the block containing the given sector.
func (mgr *MFDManager) getMFDAddressForBlock(address types.MFDRelativeAddress) types.MFDRelativeAddress {
	ldat := getLDATIndexFromMFDAddress(address)
	mask := mgr.fixedPackDescriptors[ldat].packMask
	return types.MFDRelativeAddress(uint64(address) & ^mask)
}

// getMFDBlock returns a slice corresponding to all the sectors in the physical block
// containing the sector represented by the given address. Used for reading/writing MFD blocks.
func (mgr *MFDManager) getMFDBlock(address types.MFDRelativeAddress) ([]pkg.Word36, error) {
	ldatAndTrack := address & 0_007777_777700
	data, ok := mgr.cachedTracks[ldatAndTrack]
	if !ok {
		log.Printf("MFDMgr:getMFDBlock address:%v is not in cache", address)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return nil, fmt.Errorf("internal error")
	}

	ldat := getLDATIndexFromMFDAddress(address)
	sector := getMFDSectorIdFromMFDAddress(address)
	mask := mgr.fixedPackDescriptors[ldat].packMask
	baseSectorId := uint64(sector) & ^mask

	start := 28 * baseSectorId
	end := start + uint64(mgr.fixedPackDescriptors[ldat].wordsPerBlock)
	return data[start:end], nil
}

// getMFDSector returns a slice corresponding to the portion of the MFD block which represents the indicated sector.
// CALL UNDER LOCK
func (mgr *MFDManager) getMFDSector(address types.MFDRelativeAddress) ([]pkg.Word36, error) {
	ldatAndTrack := address & 0_007777_777700
	data, ok := mgr.cachedTracks[ldatAndTrack]
	if !ok {
		log.Printf("MFDMgr:getMFDSector address:%v is not in cache", address)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return nil, fmt.Errorf("internal error")
	}

	sectorId := getMFDSectorIdFromMFDAddress(address)
	start := 28 * sectorId
	end := start + 28
	return data[start:end], nil
}

func (mgr *MFDManager) markSectorDirty(address types.MFDRelativeAddress) {
	blockAddr := mgr.getMFDAddressForBlock(address)
	mgr.dirtyBlocks[blockAddr] = true
}

func (mgr *MFDManager) writeLookupTableEntry(qualifier string, filename string, leadItem0Addr types.MFDRelativeAddress) {
	_, ok := mgr.fixedLookupTable[qualifier]
	if !ok {
		mgr.fixedLookupTable[qualifier] = make(map[string]types.MFDRelativeAddress)
	}
	mgr.fixedLookupTable[qualifier][filename] = leadItem0Addr
}

// ----------------------------------------------------------------

func (mgr *MFDManager) threadPersist() {
	// iterate over dirty block addresses
	mgr.mutex.Lock()
	if len(mgr.dirtyBlocks) > 0 {
		var blockAddr types.MFDRelativeAddress
		for key, _ := range mgr.dirtyBlocks {
			blockAddr = key
			break
		}

		delete(mgr.dirtyBlocks, blockAddr)
		mgr.mutex.Unlock()

		block, err := mgr.getMFDBlock(blockAddr)
		if err != nil {
			log.Printf("MFDMgr:cannot find MFD block for dirty block address:%012o", blockAddr)
			mgr.exec.Stop(types.StopDirectoryErrors)
			return
		}

		mfdTrackId := getMFDTrackIdFromMFDAddress(blockAddr)
		mfdSectorId := getMFDSectorIdFromMFDAddress(blockAddr)

		ldat, devTrackId, err := mgr.convertFileRelativeTrackId(mgr.mfdFileMainItem0Address, types.TrackId(mfdTrackId))
		if err != nil {
			log.Printf("MFDMgr:error converting mfdaddr:%012o trackId:%06v", mgr.mfdFileMainItem0Address, mfdTrackId)
			mgr.exec.Stop(types.StopDirectoryErrors)
			return
		} else if ldat == 0_400000 {
			log.Printf("MFDMgr:error converting mfdaddr:%012o trackId:%06v track not allocated",
				mgr.mfdFileMainItem0Address, mfdTrackId)
			mgr.exec.Stop(types.StopDirectoryErrors)
			return
		}

		packDesc, ok := mgr.fixedPackDescriptors[ldat]
		if !ok {
			log.Printf("MFDMgr:threadPersist cannot find packDesc for ldat:%04v", ldat)
			mgr.exec.Stop(types.StopDirectoryErrors)
			return
		}

		blocksPerTrack := 1792 / packDesc.wordsPerBlock
		sectorsPerBlock := packDesc.wordsPerBlock / 28
		devBlockId := uint64(devTrackId) * uint64(blocksPerTrack)
		devBlockId += uint64(mfdSectorId) / uint64(sectorsPerBlock)
		ioPkt := nodeMgr.NewDiskIoPacketRead(packDesc.deviceId, types.BlockId(devBlockId), block)
		mgr.exec.GetNodeManager().RouteIo(ioPkt)
		ioStat := ioPkt.GetIoStatus()
		if ioStat != types.IosComplete {
			log.Printf("MFDMgr:IO error writing MFD block status=%v", ioStat)
			mgr.exec.Stop(types.StopInternalExecIOFailed)
			return
		}

		mgr.mutex.Lock()
	}
}

func (mgr *MFDManager) thread() {
	mgr.threadStarted = true

	for !mgr.terminateThread {
		time.Sleep(25 * time.Millisecond)
		mgr.threadPersist()
		// TODO - anything else?
	}

	mgr.threadStopped = true
}

func (mgr *MFDManager) threadStart() {
	mgr.terminateThread = false
	if !mgr.threadStarted {
		go mgr.thread()
		for !mgr.threadStarted {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *MFDManager) threadStop() {
	if mgr.threadStarted {
		mgr.terminateThread = true
		for !mgr.threadStopped {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *MFDManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vMFDManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  isInitialized:   %v\n", indent, mgr.isInitialized)
	_, _ = fmt.Fprintf(dest, "%v  threadStarted:   %v\n", indent, mgr.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:   %v\n", indent, mgr.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, mgr.terminateThread)

	_, _ = fmt.Fprintf(dest, "%v  init MS:         %v\n", indent, mgr.msInitialize)

	_, _ = fmt.Fprintf(dest, "%v  Fixed Packs:\n", indent)
	for ldat, packDesc := range mgr.fixedPackDescriptors {
		_, _ = fmt.Fprintf(dest, "%v    ldat=%04o %s alloc=%v mask=%06o\n",
			indent,
			ldat,
			packDesc.packAttributes.PackName,
			packDesc.canAllocate,
			packDesc.packMask)

		_, _ = fmt.Fprintf(dest, "    %v      FreeSpace trackId  trackCount\n", indent)
		for _, fsRegion := range packDesc.fixedFeeSpace.content {
			_, _ = fmt.Fprintf(dest, "%v               %7v  %10v\n", indent, fsRegion.trackId, fsRegion.trackCount)
		}
	}

	_, _ = fmt.Fprintf(dest, "%v  MFD Cache:\n", indent)
	for addr, data := range mgr.cachedTracks {
		_, _ = fmt.Fprintf(dest, "%v    mfdAddr:%012o\n", indent, addr)
		pkg.DumpWord36Buffer(data, 7)
	}

	_, _ = fmt.Fprintf(dest, "%v  Dirty cache blocks:\n", indent)
	for addr, _ := range mgr.dirtyBlocks {
		_, _ = fmt.Fprintf(dest, "%v    %012o\n", indent, addr)
	}

	_, _ = fmt.Fprintf(dest, "%v  Fixed Lookup Table:\n", indent)
	for qual, sub := range mgr.fixedLookupTable {
		for file, addr := range sub {
			qualFile := qual + "*" + file
			_, _ = fmt.Fprintf(dest, "%v    %-25s  %012o\n", indent, qualFile, addr)
		}
	}

	_, _ = fmt.Fprintf(dest, "%v  Queued device-ready notifications:\n", indent)
	for devId, ready := range mgr.deviceReadyNotificationQueue {
		wId := pkg.Word36(devId)
		_, _ = fmt.Fprintf(dest, "%v    devId:0%v ready:%v\n", indent, wId.ToStringAsOctal(), ready)
	}
}
