// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "fmt"

type ActiveBaseTableEntry struct {
	bankLevel           uint64 // the top 3 bits of the extended mode L.BDI
	bankDescriptorIndex uint64 // only the bottom 15 bits of the extended mode L.BDI
	subsetSpecification uint64
}

func (abte *ActiveBaseTableEntry) GetComposite() uint64 {
	return (abte.bankLevel << 33) | (abte.bankDescriptorIndex << 18) | abte.subsetSpecification
}

func (abte *ActiveBaseTableEntry) GetString() string {
	return fmt.Sprintf(
		"level=%0o bdi=%05o subset=%06o",
		abte.bankLevel,
		abte.bankDescriptorIndex,
		abte.subsetSpecification)
}

func (abte *ActiveBaseTableEntry) SetComposite(value uint64) *ActiveBaseTableEntry {
	return abte.SetBankLevel(value >> 33).SetBankDescriptorIndex(value >> 18).SetSubsetSpecification(value)
}

func (abte *ActiveBaseTableEntry) SetBankLevel(value uint64) *ActiveBaseTableEntry {
	abte.bankLevel = value & 07
	return abte
}

func (abte *ActiveBaseTableEntry) SetBankDescriptorIndex(value uint64) *ActiveBaseTableEntry {
	abte.bankDescriptorIndex = value & 077777
	return abte
}

func (abte *ActiveBaseTableEntry) SetSubsetSpecification(value uint64) *ActiveBaseTableEntry {
	abte.subsetSpecification = value & 0777777
	return abte
}

func NewActiveBaseTableEntry(bankLevel uint64, bankDescriptorIndex uint64, subsetSpecification uint64) *ActiveBaseTableEntry {
	abte := ActiveBaseTableEntry{}
	abte.SetBankLevel(bankLevel).SetBankDescriptorIndex(bankDescriptorIndex).SetSubsetSpecification(subsetSpecification)
	return &abte
}

func NewActiveBaseTableEntryFromComposite(value uint64) *ActiveBaseTableEntry {
	abte := ActiveBaseTableEntry{}
	abte.SetComposite(value)
	return &abte
}
