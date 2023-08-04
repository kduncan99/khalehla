// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

type IndicatorKeyRegister struct {
	//	In an ICS frame this contains interrupt status.
	//	Ignored by UR instructionType
	shortStatusField InterruptShortStatus

	//	True if there is a valid current instructionType in F0.
	instructionInF0 bool

	//	True if we are currently executing an ER (or similar) instructionType
	//	The EXR target is currently in F0
	executeRepeatedInstruction bool

	//	A breakpoint match condition has occurred
	breakpointRegisterMatchCondition bool

	//	Software break condition (set only by UR instructionType)
	softwareBreak bool

	//	In an ICS frame, this is the interrupt class number of the interrupt causing the entry
	interruptClassField InterruptClass

	//	Current access key for the code being executed
	accessKey *AccessKey
}

func NewIndicatorKeyRegister() *IndicatorKeyRegister {
	return &IndicatorKeyRegister{
		shortStatusField:                 0,
		instructionInF0:                  false,
		executeRepeatedInstruction:       false,
		breakpointRegisterMatchCondition: false,
		softwareBreak:                    false,
		interruptClassField:              0,
		accessKey:                        NewAccessKeyFromComposite(0),
	}
}

func (ikr *IndicatorKeyRegister) Clear() {
	ikr.shortStatusField = 0
	ikr.instructionInF0 = false
	ikr.executeRepeatedInstruction = false
	ikr.breakpointRegisterMatchCondition = false
	ikr.softwareBreak = false
	ikr.interruptClassField = 0
	ikr.accessKey.SetRing(0).SetDomain(0)
}

func (ikr *IndicatorKeyRegister) GetAccessKey() *AccessKey {
	return ikr.accessKey
}

func (ikr *IndicatorKeyRegister) GetComposite() Word36 {
	value := Word36(0)
	value.SetS1(uint64(ikr.shortStatusField))
	if ikr.instructionInF0 {
		value |= 0_004000_000000
	}
	if ikr.executeRepeatedInstruction {
		value |= 0_002000_000000
	}
	if ikr.breakpointRegisterMatchCondition {
		value |= 0_000400_000000
	}
	if ikr.softwareBreak {
		value |= 0_000200_000000
	}
	value.SetS3(uint64(ikr.GetInterruptClassField()))
	value.SetH2(ikr.GetAccessKey().GetComposite())

	return value
}

func (ikr *IndicatorKeyRegister) GetInterruptClassField() InterruptClass {
	return ikr.interruptClassField
}

func (ikr *IndicatorKeyRegister) GetShortStatusField() InterruptShortStatus {
	return ikr.shortStatusField
}

func (ikr *IndicatorKeyRegister) IsBreakpointRegisterMatchCondition() bool {
	return ikr.breakpointRegisterMatchCondition
}

func (ikr *IndicatorKeyRegister) IsExecuteRepeatedInstruction() bool {
	return ikr.executeRepeatedInstruction
}

func (ikr *IndicatorKeyRegister) IsInstructionInF0() bool {
	return ikr.instructionInF0
}

func (ikr *IndicatorKeyRegister) IsSoftwareBreak() bool {
	return ikr.softwareBreak
}

func (ikr *IndicatorKeyRegister) SetAccessKey(key *AccessKey) *IndicatorKeyRegister {
	ikr.accessKey.SetRing(key.ring)
	ikr.accessKey.SetDomain(key.domain)
	return ikr
}

func (ikr *IndicatorKeyRegister) SetBreakpointRegisterMatchCondition(value bool) *IndicatorKeyRegister {
	ikr.breakpointRegisterMatchCondition = value
	return ikr
}

func (ikr *IndicatorKeyRegister) SetComposite(value uint64) *IndicatorKeyRegister {
	ikr.shortStatusField = InterruptShortStatus((value >> 30) & 077)
	ikr.instructionInF0 = value&0_004000_000000 != 0
	ikr.executeRepeatedInstruction = value&0_002000_000000 != 0
	ikr.breakpointRegisterMatchCondition = value&0_000400_000000 != 0
	ikr.softwareBreak = value&0_000200_000000 != 0
	ikr.interruptClassField = InterruptClass((value >> 18) & 077)
	ikr.accessKey = NewAccessKeyFromComposite(value & 0_777777)

	return ikr
}

func (ikr *IndicatorKeyRegister) SetInstructionInF0(value bool) *IndicatorKeyRegister {
	ikr.instructionInF0 = value
	return ikr
}

func (ikr *IndicatorKeyRegister) SetInterruptClassField(value InterruptClass) *IndicatorKeyRegister {
	ikr.interruptClassField = value & 077
	return ikr
}

func (ikr *IndicatorKeyRegister) SetShortStatusField(value InterruptShortStatus) *IndicatorKeyRegister {
	ikr.shortStatusField = value & 077
	return ikr
}
