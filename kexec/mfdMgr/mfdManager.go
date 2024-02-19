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
	exec                    types.IExec
	mutex                   sync.Mutex
	threadDone              bool
	mfdFileMainItem0Address types.MFDRelativeAddress                       // MFD address of MFD$$ main file item 0
	cachedTracks            map[types.MFDRelativeAddress][]pkg.Word36      // key is MFD addr of first sector in track
	dirtyBlocks             map[types.MFDRelativeAddress]bool              // MFD addresses of blocks containing dirty sectors
	freeMFDSectors          []types.MFDRelativeAddress                     // MFD addresses of existing but unused MFD sectors
	fixedPackDescriptors    map[types.LDATIndex]*fixedPackDescriptor       // fpDesc's of all known fixed packs
	fileLeadItemLookupTable map[string]map[string]types.MFDRelativeAddress // MFD address of lead item 0 of all cataloged files
	fileAllocations         *fileAllocationTable                           // Tracks file allocations of all assigned files
	// TODO make fileAllocations just a map of fae's (like it is now, but not in a separate struct)
}

func NewMFDManager(exec types.IExec) *MFDManager {
	return &MFDManager{
		exec:                    exec,
		cachedTracks:            make(map[types.MFDRelativeAddress][]pkg.Word36),
		dirtyBlocks:             make(map[types.MFDRelativeAddress]bool),
		freeMFDSectors:          make([]types.MFDRelativeAddress, 0),
		fixedPackDescriptors:    make(map[types.LDATIndex]*fixedPackDescriptor),
		fileLeadItemLookupTable: make(map[string]map[string]types.MFDRelativeAddress),
		fileAllocations:         newFileAllocationTable(),
	}
}

// Boot is invoked when the exec is booting
func (mgr *MFDManager) Boot() error {
	log.Printf("MFDMgr:Boot")

	// reset tables
	mgr.cachedTracks = make(map[types.MFDRelativeAddress][]pkg.Word36)
	mgr.dirtyBlocks = make(map[types.MFDRelativeAddress]bool)
	mgr.freeMFDSectors = make([]types.MFDRelativeAddress, 0)
	mgr.fixedPackDescriptors = make(map[types.LDATIndex]*fixedPackDescriptor)
	mgr.fileLeadItemLookupTable = make(map[string]map[string]types.MFDRelativeAddress)
	mgr.fileAllocations = newFileAllocationTable()

	return nil
}

// Close is invoked when the application is shutting down
func (mgr *MFDManager) Close() {
	log.Printf("MFDMgr:Close")
	// nothing to do
}

// Initialize is invoked when the application is starting up
// If we encounter any trouble, we return an error and the application stops.
func (mgr *MFDManager) Initialize() error {
	log.Printf("MFDMgr:Initialize")
	// nothing much to do here
	return nil
}

// Stop is invoked when the exec is stopping
func (mgr *MFDManager) Stop() {
	log.Printf("MFDMgr:Stop")
	// nothing to do
}

func (mgr *MFDManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vMFDManager ----------------------------------------------------\n", indent)

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
	for qual, sub := range mgr.fileLeadItemLookupTable {
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
}
