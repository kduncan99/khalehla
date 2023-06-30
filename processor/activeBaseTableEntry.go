package processor

type ActiveBaseTableEntry struct {
	bankLevel           uint
	bankDescriptorIndex uint
	subsetSpecification uint
}

func (abte *ActiveBaseTableEntry) SetBankLevel(value uint) *ActiveBaseTableEntry {
	abte.bankLevel = value & 07
	return abte
}

func (abte *ActiveBaseTableEntry) SetBankDescriptorIndex(value uint) *ActiveBaseTableEntry {
	abte.bankDescriptorIndex = value & 077777
	return abte
}

func (abte *ActiveBaseTableEntry) SetSubsetSpecification(value uint) *ActiveBaseTableEntry {
	abte.subsetSpecification = value & 0777777
	return abte
}

func (abte *ActiveBaseTableEntry) GetComposite() uint64 {
	return (uint64(abte.bankLevel) << 33) | (uint64(abte.bankDescriptorIndex) << 18) | uint64(abte.subsetSpecification)
}

func (abte *ActiveBaseTableEntry) SetComposite(value uint64) *ActiveBaseTableEntry {
	return abte.SetBankLevel(uint(value >> 33)).SetBankDescriptorIndex(uint(value >> 18)).SetSubsetSpecification(uint(value))
}

func NewActiveBaseTableEntryFromComposite(value uint64) *ActiveBaseTableEntry {
	abte := ActiveBaseTableEntry{}
	abte.SetComposite(value)
	return &abte
}

func NewActiveBaseTableEntryFromComponents(bankLevel uint, bankDescriptorIndex uint, subsetSpecification uint) *ActiveBaseTableEntry {
	abte := ActiveBaseTableEntry{}
	abte.SetBankLevel(bankLevel).SetBankDescriptorIndex(bankDescriptorIndex).SetSubsetSpecification(subsetSpecification)
	return &abte
}
