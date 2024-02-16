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
	"strings"
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
// The entire fixed MFD is kept loaded in core, arranged in blocks.
// A sector address is converted to pack/unit-relative-block-id by developing a mask for each pack.
// The pack-mask has bits set corresponding to the record length indicating the number of sectors contained in a
// pack record, and is equal to sectors-per-record minus 1.
// e.g., for prep factor 28, there is 1 sector per physical record, and the pack-mask is 000.
//       for prep factor 112, there are 4 sectors per physical record, and the pack-mask is 003.
//       for prep factor 1792, there are 64 sectors per physical record, and the pack-mask is 077.
// There are three components to looking up an MFD sector in the cache table.
// 		The MFD-relative-LDAT (mrLDATIndex) is equal to MFD-relative-sector >> 18
//  	The MFD-relative-block (mrBlockId) is equal to MFD-relative-LDAT & 0777777 & ^pack-mask
//			Note that these values are *not* monotonically increasing - there may be gaps depending upon
//			the pack prep factor.
// 		The MFD-relative-sector (mrSectorId) is equal to MFD-relative-sector & pack-mask
//
// TODO the following is not entirely correct - fix it
// The MFD blocks are accessible first in a map keyed by mrLDATIndex -> pack descriptor.
// The pack descriptor contains information such as whether allocation is allowed on the pack, and the block size
// of the pack - and it contains a map of blocks keyed by MFD-relative block ID to a block descriptor.
// The block descriptors are accessible via a map keyed by mrBlockId -> block descriptor.
// The block descriptor contains the corresponding physical pack-relative block-id,
// a flag indicating whether the block needs to be persisted to disk,
// and an array of Word36 entities containing the data for the MFD block.
// Finally, the 28-word MFD sector is located within that block, at the offset indicated by 28 * mrSectorId.
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
	deviceId         types.DeviceIdentifier
	packAttributes   *types.PackAttributes // created and maintained by facMgr
	wordsPerBlock    types.PrepFactor
	canAllocate      bool // true if pack is UP, false if it is SU
	packMask         uint64
	trackDescriptors map[types.MFDTrackId]fixedTrackDescriptor
	fixedFeeSpace    *packFreeSpaceTable
}

func newFixedPackDescriptor(
	deviceId types.DeviceIdentifier,
	packAttrs *types.PackAttributes,
	allocatable bool) *fixedPackDescriptor {
	recordLength := packAttrs.Label[04].GetH2()
	trackCount := packAttrs.Label[016].GetW()
	return &fixedPackDescriptor{
		deviceId:         deviceId,
		packAttributes:   packAttrs,
		wordsPerBlock:    types.PrepFactor(recordLength),
		canAllocate:      allocatable,
		packMask:         (recordLength / 28) - 1,
		trackDescriptors: make(map[types.MFDTrackId]fixedTrackDescriptor),
		fixedFeeSpace:    newPackFreeSpaceTable(types.TrackCount(trackCount)),
	}
}

type fixedTrackDescriptor struct {
	blockDescriptors map[types.MFDBlockId]fixedBlockDescriptor
}

type fixedBlockDescriptor struct {
	packRelativeBlockId types.BlockId
	needToPersist       bool
	data                []pkg.Word36
}

type MFDManager struct {
	exec                         types.IExec
	mutex                        sync.Mutex
	isInitialized                bool
	terminateThread              bool
	threadStarted                bool
	threadStopped                bool
	msInitialize                 bool
	deviceReadyNotificationQueue map[types.DeviceIdentifier]bool
	fixedLDAT                    map[types.LDATIndex]*fixedPackDescriptor
	needsPersist                 bool // true if any block in fixedLDAT needs to be persisted
	fixedLookupTable             map[string]types.MFDRelativeAddress
	fileAllocations              fileAllocationTable
}

func NewMFDManager(exec types.IExec) *MFDManager {
	return &MFDManager{
		exec:                         exec,
		deviceReadyNotificationQueue: make(map[types.DeviceIdentifier]bool),
		fixedLDAT:                    make(map[types.LDATIndex]*fixedPackDescriptor),
		fixedLookupTable:             make(map[string]types.MFDRelativeAddress),
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
// Marks the blocks in the track to be persisted, but does NOT update any other portion of the MFD.
// mfdAddress is the MFD-relative address of the first sector of the track.
// trackAddress is the device-relative-word-address of the physical block corresponding to the first
// sector of the track.
func (mgr *MFDManager) establishNewMFDTrack(ldatIndex types.LDATIndex, mfdTrackId types.MFDTrackId, trackAddress types.DeviceRelativeWordAddress) error {
	pDesc, ok := mgr.fixedLDAT[ldatIndex]
	if !ok {
		log.Printf("MFDMgr:internal error establishNewMFDTrack ldatIndex:%v is unknown", ldatIndex)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return fmt.Errorf("internal error")
	}

	tDesc := fixedTrackDescriptor{
		blockDescriptors: make(map[types.MFDBlockId]fixedBlockDescriptor),
	}

	blocksPerTrack := 1792 / pDesc.wordsPerBlock
	blockId := types.MFDBlockId(0)
	packRelativeBlockId := types.BlockId(uint64(trackAddress) / uint64(pDesc.wordsPerBlock))
	for bx := 0; bx < int(blocksPerTrack); bx++ {
		bd := fixedBlockDescriptor{}
		bd.needToPersist = true
		bd.packRelativeBlockId = packRelativeBlockId
		bd.data = make([]pkg.Word36, pDesc.wordsPerBlock)
		tDesc.blockDescriptors[blockId] = bd
		blockId++
		packRelativeBlockId++
	}

	pDesc.trackDescriptors[mfdTrackId] = tDesc
	return nil
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

// getMFDSector returns a slice corresponding to the portion of the MFD block which represents the
// indicated sector.
func (mgr *MFDManager) getMFDSector(address types.MFDRelativeAddress, markPersist bool) ([]pkg.Word36, error) {
	ldatIndex := getLDATIndexFromMFDAddress(address)
	pDesc, ok := mgr.fixedLDAT[ldatIndex]
	if !ok {
		log.Printf("MFDMgr:internal error establishNewMFDTrack ldatIndex:%v is unknown", ldatIndex)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return nil, fmt.Errorf("internal error")
	}

	trackId := getMFDTrackIdFromMFDAddress(address)
	tDesc := pDesc.trackDescriptors[trackId]

	sectorId := getMFDSectorIdFromMFDAddress(address)
	sectorsPerBlock := pDesc.wordsPerBlock / 28
	blockId := types.MFDBlockId(uint(sectorId) / uint(sectorsPerBlock))
	sectorOffset := uint(sectorId) % uint(sectorsPerBlock)

	start := sectorOffset * 28
	end := start + 28
	bDesc := tDesc.blockDescriptors[blockId]
	sector := bDesc.data[start:end]
	bDesc.needToPersist = markPersist

	return sector, nil
}

// markMFDSectorForPersist marks the block which contains the indicated sector to be persisted
func (mgr *MFDManager) markMFDSectorForPersist(address types.MFDRelativeAddress) error {
	ldatIndex := getLDATIndexFromMFDAddress(address)
	pDesc, ok := mgr.fixedLDAT[ldatIndex]
	if !ok {
		log.Printf("MFDMgr:internal error establishNewMFDTrack ldatIndex:%v is unknown", ldatIndex)
		mgr.exec.Stop(types.StopDirectoryErrors)
		return fmt.Errorf("internal error")
	}

	trackId := getMFDTrackIdFromMFDAddress(address)
	tDesc := pDesc.trackDescriptors[trackId]

	sectorId := getMFDSectorIdFromMFDAddress(address)
	sectorsPerBlock := pDesc.wordsPerBlock / 28
	blockId := types.MFDBlockId(uint(sectorId) / uint(sectorsPerBlock))

	bDesc := tDesc.blockDescriptors[blockId]
	bDesc.needToPersist = true
	return nil
}

// ----------------------------------------------------------------

func (mgr *MFDManager) threadPersist() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	nm := mgr.exec.GetNodeManager()
	for ldat, packDesc := range mgr.fixedLDAT {
		for _, trackDesc := range packDesc.trackDescriptors {
			for blkId, blockDesc := range trackDesc.blockDescriptors {
				if blockDesc.needToPersist {
					pkt := nodeMgr.NewDiskIoPacketWrite(packDesc.deviceId, blockDesc.packRelativeBlockId, blockDesc.data)
					nm.RouteIo(pkt)
					ioStat := pkt.GetIoStatus()
					if ioStat != types.IosComplete {
						log.Printf("MFDMgr:Cannot write directory block LDAT:%v mfdBlkId:%v packBlkId:%v ioStat:%v",
							ldat, blkId, blockDesc.packRelativeBlockId, ioStat)
						mgr.exec.Stop(types.StopInternalExecIOFailed)
						return
					}

					blockDesc.needToPersist = false
				}
			}
		}
	}
}

func (mgr *MFDManager) thread() {
	mgr.threadStarted = true

	for !mgr.terminateThread {
		time.Sleep(25 * time.Millisecond)

		if mgr.needsPersist {
			// TODO
			mgr.needsPersist = false
		}

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
	_, _ = fmt.Fprintf(dest, "%v  needs Persist:   %v\n", indent, mgr.needsPersist)

	_, _ = fmt.Fprintf(dest, "%v  Fixed Packs:\n", indent)
	for ldat, packDesc := range mgr.fixedLDAT {
		_, _ = fmt.Fprintf(dest, "%v    ldat=%04o %s alloc=%v mask=%06o\n",
			indent,
			ldat,
			packDesc.packAttributes.PackName,
			packDesc.canAllocate,
			packDesc.packMask)

		_, _ = fmt.Fprintf(dest, "%v      Block Descriptors:\n", indent)
		for trkId, trackDesc := range packDesc.trackDescriptors {
			for blkId, blockDesc := range trackDesc.blockDescriptors {
				_, _ = fmt.Fprintf(dest, "%v        mfdTrkId:%04o mfdBlkId:%04o packBlkId:%v dirty:%v\n",
					indent,
					trkId,
					blkId,
					blockDesc.packRelativeBlockId,
					blockDesc.needToPersist)
			}
		}

		_, _ = fmt.Fprintf(dest, "    %v      FreeSpace trackId  trackCount\n", indent)
		for _, fsRegion := range packDesc.fixedFeeSpace.content {
			_, _ = fmt.Fprintf(dest, "%v               %7v  %10v\n", indent, fsRegion.trackId, fsRegion.trackCount)
		}
	}

	_, _ = fmt.Fprintf(dest, "%v  Fixed Lookup Table:\n", indent)
	for str, addr := range mgr.fixedLookupTable {
		split := strings.Split(str, ":")
		qualFile := split[0] + ":" + split[1]
		_, _ = fmt.Fprintf(dest, "%v    %-25s  %012o\n", indent, qualFile, addr)
	}

	_, _ = fmt.Fprintf(dest, "%v  Queued device-ready notifications:\n", indent)
	for devId, ready := range mgr.deviceReadyNotificationQueue {
		wId := pkg.Word36(devId)
		_, _ = fmt.Fprintf(dest, "%v    devId:0%v ready:%v\n", indent, wId.ToStringAsOctal(), ready)
	}
}
