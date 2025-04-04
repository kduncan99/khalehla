// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

type RunConditionWord struct {
	inhibitPageEject                   bool // bit 3
	runWasRescheduledAfterRecoveryBoot bool // bit 4
	inhibitRunTermination              bool // bit 5
	secondMostRecentTaskHadError       bool // bit 7
	mostRecentTaskHadError             bool // bit 8

	// If bits 9 and 10 are both clear, then most recent activity termination was normal
	mostRecentActivityTerminationWasAbort bool // bit 9
	mostRecentActivityTerminationWasError bool // bit 10

	// Some activity for the most recent task (or the current task) had an error termination
	someActivityWasErrorTermination bool   // bit 11
	atSetCValue                     uint64 // 12 bits set by @SETC
	erSetCValue                     uint64 // 12 bits set by ER SETC$
}

func (rcw *RunConditionWord) GetT1() uint64 {
	var value uint64
	if rcw.inhibitPageEject {
		value |= 0_0400
	}
	if rcw.runWasRescheduledAfterRecoveryBoot {
		value |= 0_0200
	}
	if rcw.inhibitRunTermination {
		value |= 0_0100
	}
	if rcw.secondMostRecentTaskHadError {
		value |= 0_0020
	}
	if rcw.mostRecentTaskHadError {
		value |= 0_0010
	}
	if rcw.mostRecentActivityTerminationWasAbort {
		value |= 0_0004
	} else if rcw.mostRecentActivityTerminationWasError {
		value |= 0_0002
	}
	if rcw.someActivityWasErrorTermination {
		value |= 0_0001
	}
	return value
}

func (rcw *RunConditionWord) GetT2() uint64 {
	return rcw.atSetCValue
}

func (rcw *RunConditionWord) GetT3() uint64 {
	return rcw.erSetCValue
}

func (rcw *RunConditionWord) GetW() uint64 {
	return (rcw.GetT1() << 24) | (rcw.GetT2() << 12) | rcw.GetT3()
}
