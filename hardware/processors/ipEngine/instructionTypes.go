// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

// Assigns instructionType types to the various instructions.
//
//	This value is NOT related to anything in the architecture.
const (
	InvalidInstruction = iota
	CALLInstruction
	GOTOInstruction
	LAEInstruction
	LBEInstruction
	LBJInstruction
	LBUInstruction
	LDJInstruction
	LIJInstruction
	LOCLInstruction
	RTNInstruction
	URInstruction
)
