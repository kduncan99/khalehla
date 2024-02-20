// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import (
	"khalehla/kexec/facilitiesMgr"
	"khalehla/kexec/types"
)

// mfdServices contains code which provides directory-level services to all other exec code
// such as assigning files, cataloging files, and general file allocation.
func checkOptions(
	optionsGiven uint64,
	optionsAllowed uint64,
	facResult *facilitiesMgr.FacResult,
) bool {
	opt := 'A'
	mask := uint64(types.AOption)
	result := true
	for opt <= 'Z' {
		if (optionsGiven&mask != 0) && (optionsAllowed&mask == 0) {
			str := string(opt)
			facResult.PostMessage(facilitiesMgr.FacStatusIllegalOption, []string{str})
			result = false
		}
		opt++
		mask >>= 1
	}
	return result
}

// CatalogFixedFile attempts to catalog a file on fixed mass storage according to the given parameters.
// Caller must ensure the qualifier, filename, readKey, and writeKey values are upper-case
// Caller must ensure the qualifier, filename, readKey, and writeKey are valid
//
//	caller must resolve default/implied qualifier and provide the correct value
//	readKey and writeKey are empty strings if they were not specified
//
// Caller must ensure absolute file cycle is greater than zero and less than or equal to 999
// Caller must ensure relative file cycle is greater than -32 and less than +2
// Caller must ensure that only absolute or relative is specified (or none) - but not both
//
//	If neither file cycle parameter is specified, we assume relative cycle zero.
//
// Caller must ensure the equipment value is valid and known (which it should be anyway, or we can't be here)
//
//	If equipment was unspecified, caller should leave this as an empty string
//
// Caller must ensure granularity is valid
// If reserve and/or maximum were unspecified, the caller needs to determine what values to use
func (mgr *MFDManager) CatalogFixedFile(
	options uint64,
	qualifier string,
	filename string,
	absoluteFileCycle *uint,
	relativeFileCycle *int,
	readKey string,
	writeKey string,
	equipment string,
	granularity types.Granularity,
	reserve uint64,
	maximum uint64,
) *facilitiesMgr.FacResult {

	// TODO a very lot of this applies equally to fixed and removable
	//  - should we collapse the two into this one function?

	facResult := facilitiesMgr.NewFacResult()

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	// TODO check options
	//  Allow B, G, P, R, V, W, Z
	//	 G applies to the fileset AND to the file cycle... ?
	//	 to-be-cataloged and to-be-dropped is a fileset attribute indicating at least one of the cycles is in to-be state
	//	 to-be-cataloged and to-be-dropped is also a file cycle attribute
	//	 B option (save on checkpoint) applies to file cycle
	//	 to-be-write-only, to-be-read-only, to-be-dropped is file cycle
	//	 V, P, W, R apply to the file cycle
	allowed := uint64(types.BOption | types.GOption | types.POption | types.ROption | types.VOption | types.WOption | types.ZOption)
	if !checkOptions(options, allowed, facResult) {
		return facResult
	}

	// TODO check read/write keys
	//	 I suspect we have to meet write key verification, but probably not read key
	_ /*mainItem0*/, fileSetExists := mgr.fileLeadItemLookupTable[qualifier][filename]

	// TODO check equipment string
	// 	Should specify word-addressable or sector-addressable mass storage - applies to the file cycle, not the file set
	// Most type field entries supplied by Unisys start with either D or F, but sites can configure any mnemonics they wish.
	// Typically, D represents word- addressable and F represents sector-addressable mass storage.
	// If you omit type, the Exec uses sector-addressable mass storage.
	// If you are creating a new cycle for an existing file cycle set, and you omit the type,
	// the type is taken from the latest existing cycle of the file.
	// If all cycles of the file are in a to-be-cataloged or to-be-dropped state, then the type from the highest cycle of this type is used.

	// TODO check file cycles - zero? +1? absolute? negative relative cycles are not allowed
	// 	Absolute cycles are 1 to 999, Relative cycles are -{n}, 0, or +1
	//  if a +1 is currently assigned, we cannot catalog anything in this file set (a +1 can never exist if it is not assigned)
	//  if we are +0, and a fileset already exists, this is an error
	//  if we try to catalog a cycle which would delete the lowest-number cycle
	//  	and it is assigned, we have an f-cycle conflict
	//		and there are more than one such cycles (see formula below) we have an f-cycle conflict
	//		and the cycle to be deleted has a write key other than what we specificed, we have an f-cycle conflict
	//  For a new cycle to be created, its absolute F-cycle number must be within the following range:
	// 		(x-w) < z ≤ (x-y+w+1) where:
	// 		x is T3 of word 9 of the lead item (cycle number of latest F-cycle).
	// 		w is S3 of word 9 of the lead item (maximum number of F-cycles).
	// 		z is the absolute F-cycle number requested.
	// 		y is S4 of word 9 of the lead item (current range of F-cycles)
	//			this is the highest-absolute-f-cycle - lowest-absolute-f-cycle + 1
	//			*or* highest-absolute-f-cycle + 1000 - lowest-absolute-f-cycle
	//
	// e.g.: we have cycles 10, 11, 20, 21, 30, 31 and max cycles is 25
	// so x-w = 31 - 25 = 6
	// x-y+w+1 = 31 - (31 - 10 + 1) + 25 + 1 = 35
	// thus 6 < z <= 35
	if fileSetExists {
		// TODO
	}

	// TODO make it happen
	if !fileSetExists {
		// TODO create lead item(s)
	}

	return facResult
}

func (mgr *MFDManager) CatalogRemovableFile(
	options uint64,
	qualifier string,
	filename string,
	fileCycle int,
	readKey string,
	writeKey string,
	equipment string,
	granularity types.Granularity,
	reserve uint64,
	maximum uint64,
	packNames []string,
) *facilitiesMgr.FacResult {

	// TODO
	// @CAT[,options] filename[,type/reserve/granule/maximum,pack-id-1/... /pack-id-n,,,ACR-name]
	facResult := facilitiesMgr.NewFacResult()

	return facResult
}

func (mgr *MFDManager) CatalogTapeFile(
	options uint64,
	qualifier string,
	filename string,
	fileCycle int,
	readKey string,
	writeKey string,
	equipment string,
	// TODO other parameters go here
) *facilitiesMgr.FacResult {

	// TODO
	// @CAT,options filename,type[/units/log/noise/processor/tape/ format/data-converter/block-numbering/data-compression/ buffered-write/expanded-buffer,reel-1/reel-2/.../reel-n, expiration-period/mmspec,,ACR-name,CTL-pool]
	facResult := facilitiesMgr.NewFacResult()

	return facResult
}
