// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import "strings"

// handleCat updates the default and/or implied qualifier for the run control entry.
//
//	For Mass Storage Files
//		@CAT[,options] filename[,type/reserve/granule/maximum,pack-id-1/.../pack-id-n,,,ACR-name]
//	maximum of 6 fields in argument
//	options include
//		B: save on checkpoint
//		G: guarded file
//		P: make the file public (not private)
//		R: make the file read-only
//		V: file will not be unloaded
//		W: make the file write-only
//		Z: run should not be held (probably only happens on removable when the pack is not mounted)
//
//	For Tape Files
//		@CAT,options filename,type[/units/log/noise/processor/tape/
//			format/data-converter/block-numbering/data-compression/
//			buffered-write/expanded-buffer,reel-1/reel-2/.../reel-n,
//			expiration-period/mmspec,,ACR-name,CTL-pool]
//	maximum of 7 fields in argument
//	options include
//		E: even parity (not supported)
//		G: guarded file
//		H: density selection (not supported)
//		J: tape is to be unlabeled
//		L: density selection (not supported)
//		M: density selection (not supported)
//		O: odd parity (supported but ignored)
//		P: make the file public
//		R: make the file read-only
//		S: 6250 BPI (only for SCSI 9-track - future)
//		V: 1600 BPI (only for SCSI 9-track - future)
//		W: make the file write-only
//		Z: run should not be held (probably only happens on removable when the pack is not mounted)
func handleCat(pkt *handlerPacket) uint64 {
	cleanOpts := ""
	if len(pkt.options) > 0 {
		var ok bool
		cleanOpts, ok = cleanOptions(pkt)
		if !ok {
			return 0_400000_000000
		}
	}

	fields := strings.Split(pkt.arguments, ",")
	subfields := make([][]string, 0)
	for fx := 0; fx < len(fields); fx++ {
		subfields[fx] = strings.Split(fields[fx], "/")
	}

	fileSpec := fields[0]
	// TODO ensure filespec is valid - parse it into a FileSpecification struct
	// 	[ [qualifier] '*' ] filename [ '/' [read_key] [ '/' [write_key] ] ] ['.']

	// If mnemonic is not specified...
	//	If this is a new fileset, use 'F', otherwise
	//		If there is any existing (i.e., not to-be) cycle, use the equipment type of the latest existing cycle
	//		Else use the equipment type of the latest to-be cycle
	// TODO need mnemonic->equip-type in config, and equip-type in node manager.
	mnemonic := ""
	if len(fields) >= 2 {
		mnemonic = subfields[1][0]
		// TODO ensure mnemonic is either empty or valid
		//  [ ' ' | 'D' | 'F' | 'T' ]
	}

	// TODO
	//  for mass storage:
	//    validate options
	//    reserve, granule, max
	//      are subsequent subfields of fields[1]
	//    pack-names are subfields of fields[2]
	//    fields[3] and [4] are to be blank
	//    fields[5] is for ACR name (which we don't do)

	// TODO
	//  for tape:
	//    validate options
	//    units, logical, noise, processor, tape, format, dataConv, blockNum, compression, buffered, expanded-buf
	//      are subsequent subfields of fields[1]
	//    reel-numbers are subfields of fields[2]
	//    fields[3] and [4] are to be blank
	//    fields[5] is for ACR name (which we don't do)
	//    fields[6] is for CTL-pool... which *maybe* we do?

	// TODO hand off to fac mgr

	return 0
}
