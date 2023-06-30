package processor

import "kalehla/types"

type IndicatorKeyRegister struct {
	//	In an ICS frame this contains interrupt status.
	//	Ignored by UR instructionType
	shortStatusField uint

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
	interruptClassField uint

	//	Current access key for the code being executed
	accessKey *AccessKey
}

func (ikr *IndicatorKeyRegister) GetComposite() types.Word36 {
	value := types.Word36(0)
	value.SetS1(ikr.shortStatusField)
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
	value.SetS3(ikr.interruptClassField)
	value.SetH2(ikr.accessKey.GetComposite())

	return value
}

func (ikr *IndicatorKeyRegister) SetComposite(value uint64) *IndicatorKeyRegister {
	ikr.shortStatusField = uint(value>>30) & 077
	ikr.instructionInF0 = value&0_004000_000000 != 0
	ikr.executeRepeatedInstruction = value&0_002000_000000 != 0
	ikr.breakpointRegisterMatchCondition = value&0_000400_000000 != 0
	ikr.softwareBreak = value&0_000200_000000 != 0
	ikr.interruptClassField = uint(value>>18) & 077
	ikr.accessKey = NewAccessKeyFromComposite(uint(value & 0_777777))

	return ikr
}

func (ikr *IndicatorKeyRegister) clear() {
	ikr.shortStatusField = 0
	ikr.instructionInF0 = false
	ikr.executeRepeatedInstruction = false
	ikr.breakpointRegisterMatchCondition = false
	ikr.softwareBreak = false
	ikr.interruptClassField = 0
	ikr.accessKey.SetRing(0).SetDomain(0)
}

func (ikr *IndicatorKeyRegister) SetInterruptClassField(value uint) {
	ikr.interruptClassField = value & 077
}

func (ikr *IndicatorKeyRegister) SetShortStatusField(value uint) {
	ikr.shortStatusField = value & 077
}
