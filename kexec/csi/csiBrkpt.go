// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

func handleBrkpt(pkt *handlerPacket) uint64 {

	/*
		For primary print files
			@BRKPT,L PRINT$[/part-name]
			@BRKPT PRINT$[,filename]
		For alternate print files
			@BRKPT[,E] internal-fname
			@BRKPT[,E] ,fname
		E inhibits EOF positioning for alt read files on magtape
	*/
	// TODO
	return 0
}
