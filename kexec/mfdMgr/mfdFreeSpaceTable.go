// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

// TODO are we even using this?
// type mfdFreeSpaceTable struct {
// 	entries map[types.LDATIndex]*packFreeSpaceTable
// }
//
// func newFreeSpaceTable() *mfdFreeSpaceTable {
// 	return &mfdFreeSpaceTable{
// 		entries: make(map[types.LDATIndex]*packFreeSpaceTable),
// 	}
// }
//
// // allocateSpecificTrackRegion is used only when it has been determined by some external means, that a particular
// // track or range of tracks is not to be allocated otherwise (such as for VOL1 or directory tracks).
// func (fst *mfdFreeSpaceTable) allocateSpecificTrackRegion(
// 	ldatIndex types.LDATIndex,
// 	trackId types.TrackId,
// 	trackCount types.TrackCount) error {
//
// 	packTable, ok := fst.entries[ldatIndex]
// 	if !ok {
// 		return fmt.Errorf("pack not found")
// 	}
//
// 	ok = packTable.markTrackRegionAllocated(ldatIndex, trackId, trackCount)
// 	if !ok {
// 		return fmt.Errorf("track not allocated")
// 	}
// 	return nil
// }
//
// // allocateExactContiguous allocates space from any pack such that the allocated space represents exactly one
// // free area on the pack ... we do this for contiguous allocations where-in we'd like to avoid fragmenting
// // free space.
// // Returns true, ldat, trackId if successful, else false, 0, 0
// func (fst *mfdFreeSpaceTable) allocateExactContiguous(trackCount types.TrackCount) (bool, types.LDATIndex, types.TrackId) {
// 	for ldatIndex, packTable := range fst.entries {
// 		for fx, fse := range packTable.content {
// 			if fse.trackCount == trackCount {
// 				packTable.content = append(packTable.content[:fx], packTable.content[fx:]...)
// 				return true, ldatIndex, fse.trackId
// 			}
// 		}
// 	}
//
// 	return false, 0, 0
// }
//
// // establishFreeSpaceTableForPack sets up a free space table for a particular pack, identified by its LDAT index
// func (fst *mfdFreeSpaceTable) establishFreeSpaceTableForPack(
// 	ldatIndex types.LDATIndex,
// 	trackCount types.TrackCount) error {
//
// 	_, ok := fst.entries[ldatIndex]
// 	if ok {
// 		log.Printf("establishFreeSpaceTableForPack ldat:%v trackCount:%v table already exists", ldatIndex, trackCount)
// 		return fmt.Errorf("table already exists")
// 	}
//
// 	fst.entries[ldatIndex] = newPackFreeSpaceTable(trackCount)
// 	return nil
// }
