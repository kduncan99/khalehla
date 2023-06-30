package processor

import "kalehla/types"

type ProgramAddressRegister struct {
	level               uint
	bankDescriptorIndex uint
	programCounter      uint
}

func (par *ProgramAddressRegister) GetLevel() uint {
	return par.level
}

func (par *ProgramAddressRegister) GetBankDescriptorIndex() uint {
	return par.bankDescriptorIndex
}

func (par *ProgramAddressRegister) GetProgramCounter() uint {
	return par.programCounter
}

func (par *ProgramAddressRegister) GetComposite() types.Word36 {
	return (types.Word36(par.level) << 33) |
		(types.Word36(par.bankDescriptorIndex) << 18) |
		types.Word36(par.programCounter)
}

func (par *ProgramAddressRegister) IncrementProgramCounter() {
	par.programCounter++
	if par.programCounter > 0777777 {
		par.programCounter = 0
	}
}

func (par *ProgramAddressRegister) SetLevel(value uint) *ProgramAddressRegister {
	par.level = value & 07
	return par
}

func (par *ProgramAddressRegister) SetBankDescriptorIndex(value uint) *ProgramAddressRegister {
	par.bankDescriptorIndex = value & 077777
	return par
}

func (par *ProgramAddressRegister) SetProgramCounter(value uint) *ProgramAddressRegister {
	par.programCounter = value & 0777777
	return par
}

func (par *ProgramAddressRegister) SetComposite(value uint64) *ProgramAddressRegister {
	par.level = uint((value >> 33) & 07)
	par.bankDescriptorIndex = uint((value >> 18) & 077777)
	par.programCounter = uint(value & 0777777)
	return par
}

func NewProgramAddressRegister(level uint, bankDescriptorIndex uint, programCounter uint) *ProgramAddressRegister {
	return &ProgramAddressRegister{
		level:               level & 07,
		bankDescriptorIndex: bankDescriptorIndex & 077777,
		programCounter:      programCounter & 0777777,
	}
}
