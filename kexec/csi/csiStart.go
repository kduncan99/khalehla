// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

func handleStart(pkt *handlerPacket) uint64 {

	/*
		@START runstream-name
		@START[,scheduling-priority/options/processor-dispatching-priority]
			runstream-name[,setc,run-id,acct-id/user-id,project-id,runtime/deadline,pages/cards,start-time]
		scheduling-priority is one character from 'A' to 'Z'
		options include
			B, N, P, R, T, U, V, W, X, Y, Z
		processor-dispatching-priority is one character from 'A' to 'Z'
	*/
	// TODO
	return 0
}
