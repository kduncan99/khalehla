package processor

import "kalehla/types"

type IndexRegister types.Word36

func (ir *IndexRegister) DecrementModifier() {
	ir.SetXM(uint(types.AddSimple(ir.GetSignedXM(), types.Negate(ir.GetSignedXI()))))
}

func (ir *IndexRegister) DecrementModifier24() {
	ir.SetXM(uint(types.AddSimple(ir.GetSignedXM24(), types.Negate(ir.GetSignedXI12()))))
}

func (ir *IndexRegister) GetSignedXI() types.Word36 {
	return types.GetSignExtended18(uint64(*ir))
}

func (ir *IndexRegister) GetSignedXI12() types.Word36 {
	return types.GetSignExtended12(uint64(*ir))
}

func (ir *IndexRegister) GetSignedXM() types.Word36 {
	return types.GetSignExtended18(uint64(*ir))
}

func (ir *IndexRegister) GetSignedXM24() types.Word36 {
	return types.GetSignExtended24(uint64(*ir))
}

func (ir *IndexRegister) GetXI() uint {
	return uint((*ir) >> 18)
}

func (ir *IndexRegister) GetXI12() uint {
	return uint((*ir) >> 12)
}

func (ir *IndexRegister) GetXM() uint {
	return uint((*ir) & 0_777777)
}

func (ir *IndexRegister) GetXM24() uint {
	return uint((*ir) & 077_777777)
}

func (ir *IndexRegister) GetW() uint64 {
	return uint64(*ir)
}

func (ir *IndexRegister) IncrementModifier() {
	ir.SetXM(uint(types.AddSimple(ir.GetSignedXM(), ir.GetSignedXI())))
}

func (ir *IndexRegister) IncrementModifier24() {
	ir.SetXM(uint(types.AddSimple(ir.GetSignedXM24(), ir.GetSignedXI12())))
}

func (ir *IndexRegister) SetXI(op uint) {
	*ir = (*ir & 0777777) | ((IndexRegister(op) & 0777777) << 18)
}

func (ir *IndexRegister) SetXI12(op uint) {
	*ir = (*ir & 077777777) | ((IndexRegister(op) & 07777) << 24)
}

func (ir *IndexRegister) SetXM(op uint) {
	*ir = (*ir & 0777777000000) | (IndexRegister(op) & 0777777)
}

func (ir *IndexRegister) SetXM24(op uint) {
	*ir = (*ir & 0777700000000) | (IndexRegister(op) & 077777777)
}
