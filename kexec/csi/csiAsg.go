// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

func handleAsg(pkt *handlerPacket) uint64 {
	/*
		Mass-Storage:
			@ASG[,options] filename[,type/reserve/granule/maximum/placement][,pack-id-1/.../pack-id-n,,,ACR-name]
		options include
			I, Z
			B, C, G, P, R, S, U, V, W for files to be cataloged
			A, D, E, K, M, Q, R, X, Y for existing files
			T for temporary files
		In the absence of any options, default to A if the file exists and is not assigned,
			or to the previous options if it is already assigned,
			or to T if it is not assigned and does not exist

		Tape Files:
			@ASG,options filename,type[/units/log/noise/processor/tape/
				format/data-converter/block-numbering/data-compression/
				buffered-write/expanded-buffer/,reel-1/reel-2/.../
				reel-n,expiration/mmspec,ring-indicator,ACR-name,CTL-pool ]
		options include
			B, I, E, N, R, W, Z
			C, G, P, U for files to be cataloged
			A, D, K, Q, X, Y for existing files
			R, T, W for temporary files
			E, H, L, M, O, S, V are hardware options
			F, J are labeling options

		Absolute Device:
			@ASG[,options] filename,*name,pack-id
		options include
			H, I, T, Z
	*/

	// TODO
	return 0
}
