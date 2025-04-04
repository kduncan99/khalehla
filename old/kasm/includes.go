// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

var axrIncludes = []string{
	" X0   $EQU 0",
	" X1   $EQU 1",
	" X2   $EQU 2",
	" X3   $EQU 3",
	" X4   $EQU 4",
	" X5   $EQU 5",
	" X6   $EQU 6",
	" X7   $EQU 7",
	" X8   $EQU 8",
	" X9   $EQU 9",
	" X10  $EQU 10",
	" X11  $EQU 11",
	" X12  $EQU 12",
	" X13  $EQU 13",
	" X14  $EQU 14",
	" X15  $EQU 15",

	" A0   $EQU 12",
	" A1   $EQU 13",
	" A2   $EQU 14",
	" A3   $EQU 15",
	" A4   $EQU 16",
	" A5   $EQU 17",
	" A6   $EQU 18",
	" A7   $EQU 19",
	" A8   $EQU 20",
	" A9   $EQU 21",
	" A10  $EQU 22",
	" A11  $EQU 23",
	" A12  $EQU 24",
	" A13  $EQU 25",
	" A14  $EQU 26",
	" A15  $EQU 27",

	" R0   $EQU 64",
	" R1   $EQU 65",
	" R2   $EQU 66",
	" R3   $EQU 67",
	" R4   $EQU 68",
	" R5   $EQU 69",
	" R6   $EQU 70",
	" R7   $EQU 71",
	" R8   $EQU 72",
	" R9   $EQU 73",
	" R10  $EQU 74",
	" R11  $EQU 75",
	" R12  $EQU 76",
	" R13  $EQU 77",
	" R14  $EQU 78",
	" R15  $EQU 79",

	" W    $EQU 0",
	" H2   $EQU 1",
	" H1   $EQU 2",
	" XH2  $EQU 3",
	" XH1  $EQU 4",
	" T3   $EQU 5",
	" T2   $EQU 6",
	" T1   $EQU 7",
	" Q2   $EQU 4",
	" Q4   $EQU 5",
	" Q3   $EQU 6",
	" Q1   $EQU 7",
	" S6   $EQU 10",
	" S5   $EQU 11",
	" S4   $EQU 12",
	" S3   $EQU 13",
	" S2   $EQU 14",
	" S1   $EQU 15",
}

var formIncludes = []string{
	"I$*       $FORM     6,4,4,4,2,16",
	"EI$*      $FORM     6,4,4,4,2,4,12",
}

var instructionIncludes = []string{
	// f,j a,u,x
	"P         $PROC     1,1   $BASIC",
	"SA*       $NAME     001",
	"LA*       $NAME     010",
	"hi        $EQU      P(1,*2) << 1 | P(1,*1)",
	"          I$        P(0,0), P(0,1), P(1,0)-A0, P(1,2), hi, P(1,1)",
	"          $END",

	// f,j a,d,x,b
	"P         $PROC     1,1   $EXTEND",
	"SA*       $NAME     001",
	"LA*       $NAME     010",
	"hi        $EQU      P(1,*1)",
	"          EI$       P(0,0), P(0,1), P(1,0)-A0, P(1,2), hi, P(1,1), P(1,3)",
	"          $END",
}
