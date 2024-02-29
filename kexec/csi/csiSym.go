// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

import (
	"khalehla/kexec"
)

func handleSym(pkt *handlerPacket) (*kexec.FacStatusResult, uint64) {

	/*
		@SYM,options filename,[number,device,pname-1/pname-2/pname-3/.../pname-n,banner]
		@SYM,options filename,[number,user-id/U,pname-1/pname/pname-3/.../pname-n,banner]
		options include
			A, D, F, J, K, L, N, U

		values for resultCode:
			001 First field required for L option BRKPT
			002 No file specification for breakpoint
			003 Excess fields specified on breakpoint
			004 Max BRKPTS exceeded for PRINT$/PUNCH$
			005 BRKPT PRINT$ is illegal for demand.
			006 File specified on BRKPT not assigned
			007 Alternate file being breakpointed has not been referenced.
			010 I/O status writing EOF
			011 File being breakpointed to is in use.
			012 Illegal syntax for L option breakpoint
			013 No file name specified on SYM image
			014 SYM of PRINT$ illegal from demand mode
			015 SYM before BRKPT on user defined file is illegal.
			016 File name error
			017 Fac status n on attempted ASG/USE
			020 @SYM of temporary file is illegal.
			021 D option illegal for user file
			022 C option illegal for nonpunch type device
			023 Device name specified is not configured.
			024 SYM of system file is illegal.
			025 SYM of program file is illegal.
			026 Read of @SYM file illegal
			027 Drop of file that was queued by @SYM illegal
			030 SYM of file in +1 state illegal
			031 File type does not match configured devices for this group.
			032 File marked to be dropped - SYM illegal
			033 PRINT$/PUNCH$ illegal as second file name
			034 File was not queued - PR queue is not available.
			035 Number of copies specified contains illegal character.
			036 Illegal user-id specified
			037 No user-id configured
			040 PCT expansion error
			041 BRKPT to write-inhibited file not allowed
			042 File sharing environment not available: @SYM of shared file not completed.
			043 Unsolicited data transmission to a terminal not allowed.
			044 Invalid or conflicting @SYM options. File not queued.
			045 Exec was unable to complete the @SYM; main item was not found.
			046 Exec was unable to complete @SYM of file with changed name.
			047 @SYM of PRINT$ or PUNCH$ is not allowed from an impersonated activity.
	*/
	// TODO
	return nil, 0
}
