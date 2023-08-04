// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

type ProgramAddressRegister struct {
	level               uint64
	bankDescriptorIndex uint64
	programCounter      uint64
}

func (par *ProgramAddressRegister) GetLevel() uint64 {
	return par.level
}

func (par *ProgramAddressRegister) GetBankDescriptorIndex() uint64 {
	return par.bankDescriptorIndex
}

func (par *ProgramAddressRegister) GetProgramCounter() uint64 {
	return par.programCounter
}

func (par *ProgramAddressRegister) GetComposite() Word36 {
	return (Word36(par.level) << 33) |
		(Word36(par.bankDescriptorIndex) << 18) |
		Word36(par.programCounter)
}

func (par *ProgramAddressRegister) IncrementProgramCounter() {
	par.programCounter++
	if par.programCounter > 0777777 {
		par.programCounter = 0
	}
}

func (par *ProgramAddressRegister) SetLevel(value uint64) *ProgramAddressRegister {
	par.level = value & 07
	return par
}

func (par *ProgramAddressRegister) SetBankDescriptorIndex(value uint64) *ProgramAddressRegister {
	par.bankDescriptorIndex = value & 077777
	return par
}

func (par *ProgramAddressRegister) SetProgramCounter(value uint64) *ProgramAddressRegister {
	par.programCounter = value & 0777777
	return par
}

func (par *ProgramAddressRegister) SetComposite(value uint64) *ProgramAddressRegister {
	par.level = (value >> 33) & 07
	par.bankDescriptorIndex = (value >> 18) & 077777
	par.programCounter = value & 0777777
	return par
}

func NewProgramAddressRegister(level uint64, bankDescriptorIndex uint64, programCounter uint64) *ProgramAddressRegister {
	return &ProgramAddressRegister{
		level:               level & 07,
		bankDescriptorIndex: bankDescriptorIndex & 077777,
		programCounter:      programCounter & 0777777,
	}
}
