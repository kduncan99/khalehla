// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"fmt"
	"io"
	"khalehla/kexec"
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
// However, we do have a lookup table with the key being the qualifier, the subkey being the filename,
// and the value being the lead item sector 0 address for the file set.
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

// Keeps track of things which pertain to a specific disk pack (an internal MFD struct)
type packDescriptor struct {
	nodeId                     kexec.NodeIdentifier
	prepFactor                 kexec.PrepFactor
	firstDirectoryTrackAddress kexec.DeviceRelativeWordAddress
	canAllocate                bool                         // true if pack is UP, false if it is SU - must be set by fac mgr
	packMask                   uint                         // used for calculating blocks from sectors
	freeSpaceTable             *kexec.MFDPackFreeSpaceTable // represents all unallocated tracks on the pack
	mfdTrackCount              kexec.TrackCount             // number of MFD tracks allocated on the pack
	mfdSectorsUsed             uint64                       // number of MFD sectors in use among the MFD tracks
}

func newPackDescriptor(
	nodeId kexec.NodeIdentifier,
	diskAttributes *kexec.DiskAttributes,
) *packDescriptor {

	recordLength := diskAttributes.PackLabelInfo.WordsPerRecord
	trackCount := diskAttributes.PackLabelInfo.TrackCount
	facStatus := diskAttributes.GetFacNodeStatus()

	return &packDescriptor{
		nodeId:                     nodeId,
		prepFactor:                 diskAttributes.PackLabelInfo.PrepFactor,
		firstDirectoryTrackAddress: diskAttributes.PackLabelInfo.FirstDirectoryTrackAddress,
		canAllocate:                facStatus == kexec.FacNodeStatusUp || facStatus == kexec.FacNodeStatusReserved,
		packMask:                   (recordLength / 28) - 1,
		freeSpaceTable:             kexec.NewMFDPackFreeSpaceTable(trackCount),
	}
}

type MFDManager struct {
	exec                    kexec.IExec
	mutex                   sync.Mutex
	threadDone              bool
	mfdFileMainItem0Address kexec.MFDRelativeAddress                                   // MFD address of MFD$$ main file item 0
	cachedTracks            map[kexec.MFDRelativeAddress][]pkg.Word36                  // key is MFD addr of first sector in track
	dirtyBlocks             map[kexec.MFDRelativeAddress]bool                          // MFD addresses of blocks containing dirty sectors
	freeMFDSectors          []kexec.MFDRelativeAddress                                 // MFD addresses of existing but unused MFD sectors
	fixedPackDescriptors    map[kexec.LDATIndex]*packDescriptor                        // packDescriptors of all known fixed packs
	fileLeadItemLookupTable map[string]map[string]kexec.MFDRelativeAddress             // MFD address of lead item 0 of all cataloged files
	assignedFileAllocations map[kexec.MFDRelativeAddress]*kexec.MFDFileAllocationEntry // key is main item sector 0 address of file
}

func NewMFDManager(exec kexec.IExec) *MFDManager {
	return &MFDManager{
		exec:                    exec,
		cachedTracks:            make(map[kexec.MFDRelativeAddress][]pkg.Word36),
		dirtyBlocks:             make(map[kexec.MFDRelativeAddress]bool),
		freeMFDSectors:          make([]kexec.MFDRelativeAddress, 0),
		fixedPackDescriptors:    make(map[kexec.LDATIndex]*packDescriptor),
		fileLeadItemLookupTable: make(map[string]map[string]kexec.MFDRelativeAddress),
		assignedFileAllocations: make(map[kexec.MFDRelativeAddress]*kexec.MFDFileAllocationEntry),
	}
}

// Boot is invoked when the exec is booting
func (mgr *MFDManager) Boot() error {
	log.Printf("MFDMgr:Boot")

	// reset tables
	mgr.cachedTracks = make(map[kexec.MFDRelativeAddress][]pkg.Word36)
	mgr.dirtyBlocks = make(map[kexec.MFDRelativeAddress]bool)
	mgr.freeMFDSectors = make([]kexec.MFDRelativeAddress, 0)
	mgr.fixedPackDescriptors = make(map[kexec.LDATIndex]*packDescriptor)
	mgr.fileLeadItemLookupTable = make(map[string]map[string]kexec.MFDRelativeAddress)
	mgr.assignedFileAllocations = make(map[kexec.MFDRelativeAddress]*kexec.MFDFileAllocationEntry)

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
		_, _ = fmt.Fprintf(dest, "%v    ldat=%04o prep=%v alloc=%v mask=%06o\n",
			indent,
			ldat,
			packDesc.prepFactor,
			packDesc.canAllocate,
			packDesc.packMask)

		_, _ = fmt.Fprintf(dest, "    %v      FreeSpace TrackId  TrackCount\n", indent)
		for _, fsRegion := range packDesc.freeSpaceTable.Content {
			_, _ = fmt.Fprintf(dest, "%v               %7v  %10v\n", indent, fsRegion.TrackId, fsRegion.TrackCount)
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
	for _, fae := range mgr.assignedFileAllocations {
		_, _ = fmt.Fprintf(dest, "%v    mainItem:%012o 1stDAD:%012o upd:%v highest:%v\n",
			indent, fae.MainItem0Address, fae.DadItem0Address, fae.IsUpdated, fae.HighestTrackAllocated)
		for _, fileAlloc := range fae.FileAllocations {
			_, _ = fmt.Fprintf(dest, "%v      fileTrkId:%v trkCount:%v ldat:%v devTrkId:%v\n",
				indent, fileAlloc.FileRegion.TrackId, fileAlloc.FileRegion.TrackCount, fileAlloc.LDATIndex, fileAlloc.DeviceTrackId)
		}
	}
}
