// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"io"
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
	freeSpaceTable *packFreeSpaceTable
	mfdTrackCount  uint64
	mfdSectorsUsed uint64
}

func newFixedPackDescriptor(
	deviceId types.DeviceIdentifier,
	packAttrs *types.PackAttributes,
	allocatable bool,
) *fixedPackDescriptor {

	recordLength := packAttrs.Label[04].GetH2()
	trackCount := packAttrs.Label[016].GetW()
	return &fixedPackDescriptor{
		deviceId:       deviceId,
		packAttributes: packAttrs,
		wordsPerBlock:  types.PrepFactor(recordLength),
		canAllocate:    allocatable,
		packMask:       (recordLength / 28) - 1,
		freeSpaceTable: newPackFreeSpaceTable(types.TrackCount(trackCount)),
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
	freeMFDSectors               []types.MFDRelativeAddress
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
		freeMFDSectors:               make([]types.MFDRelativeAddress, 0),
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

func (mgr *MFDManager) thread() {
	mgr.threadStarted = true

	for !mgr.terminateThread {
		time.Sleep(25 * time.Millisecond)
		// TODO - anything?
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
		for _, fsRegion := range packDesc.freeSpaceTable.content {
			_, _ = fmt.Fprintf(dest, "%v               %7v  %10v\n", indent, fsRegion.trackId, fsRegion.trackCount)
		}
	}

	_, _ = fmt.Fprintf(dest, "%v  MFD Cache:\n", indent)
	for addr, data := range mgr.cachedTracks {
		secAddr := int(addr)
		for sx := 0; sx < 64; sx++ {
			_, _ = fmt.Fprintf(dest, "%v    MFD sector %012o:\n", indent, secAddr+sx)
			offset := sx * 28
			for wx := 0; wx < 28; wx += 7 {
				str := "      "
				for wy := 0; wy < 7; wy++ {
					str += fmt.Sprintf("%012o ", data[offset+wx+wy])
				}

				str += "  "
				for wy := 0; wy < 7; wy++ {
					str += data[offset+wx+wy].ToStringAsFieldata() + " "
				}

				_, _ = fmt.Fprintf(dest, "%s%s\n", indent, str)
			}
		}
	}

	_, _ = fmt.Fprintf(dest, "%v  Dirty cache blocks:\n", indent)
	for addr := range mgr.dirtyBlocks {
		_, _ = fmt.Fprintf(dest, "%v    %012o\n", indent, addr)
	}

	_, _ = fmt.Fprintf(dest, "%v  Fixed Lookup Table:\n", indent)
	for qual, sub := range mgr.fixedLookupTable {
		for file, addr := range sub {
			qualFile := qual + "*" + file
			_, _ = fmt.Fprintf(dest, "%v    %-25s  %012o\n", indent, qualFile, addr)
		}
	}

	_, _ = fmt.Fprintf(dest, "%v  Assigned file allocations:\n", indent)
	fat := mgr.fileAllocations
	for _, fae := range fat.content {
		_, _ = fmt.Fprintf(dest, "%v    mainItem:%012o 1stDAD:%012o upd:%v highest:%v\n",
			indent, fae.mainItem0Address, fae.dadItem0Address, fae.isUpdated, fae.highestTrackAllocated)
		for _, re := range fae.regionEntries {
			_, _ = fmt.Fprintf(dest, "%v      fileTrkId:%v trkCount:%v ldat:%v devTrkId:%v\n",
				indent, re.fileTrackId, re.trackCount, re.ldatIndex, re.packTrackId)
		}
	}

	_, _ = fmt.Fprintf(dest, "%v  Queued device-ready notifications:\n", indent)
	for devId, ready := range mgr.deviceReadyNotificationQueue {
		wId := pkg.Word36(devId)
		_, _ = fmt.Fprintf(dest, "%v    devId:0%v ready:%v\n", indent, wId.ToStringAsOctal(), ready)
	}
}
