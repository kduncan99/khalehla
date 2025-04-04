// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package common

type IndexRegister Word36

func (ir *IndexRegister) DecrementModifier() {
	ir.SetXM(AddSimple(ir.GetSignedXM(), Negate(ir.GetSignedXI())))
}

func (ir *IndexRegister) DecrementModifier24() {
	ir.SetXM24(AddSimple(ir.GetSignedXM24(), Negate(ir.GetSignedXI12())))
}

func (ir *IndexRegister) GetSignedXI() uint64 {
	return GetSignExtended18(ir.GetXI())
}

func (ir *IndexRegister) GetSignedXI12() uint64 {
	return GetSignExtended12(ir.GetXI12())
}

func (ir *IndexRegister) GetSignedXM() uint64 {
	return GetSignExtended18(ir.GetXM())
}

func (ir *IndexRegister) GetSignedXM24() uint64 {
	return GetSignExtended24(ir.GetXM24())
}

func (ir *IndexRegister) GetXI() uint64 {
	return uint64((*ir) >> 18)
}

func (ir *IndexRegister) GetXI12() uint64 {
	return uint64((*ir) >> 24)
}

func (ir *IndexRegister) GetXM() uint64 {
	return uint64((*ir) & 0_777777)
}

func (ir *IndexRegister) GetXM24() uint64 {
	return uint64((*ir) & 077_777777)
}

func (ir *IndexRegister) GetW() uint64 {
	return uint64(*ir)
}

func (ir *IndexRegister) IncrementModifier() {
	ir.SetXM(AddSimple(ir.GetSignedXM(), ir.GetSignedXI()))
}

func (ir *IndexRegister) IncrementModifier24() {
	ir.SetXM24(AddSimple(ir.GetSignedXM24(), ir.GetSignedXI12()))
}

func (ir *IndexRegister) SetXI(op uint64) {
	*ir = (*ir & 0777777) | ((IndexRegister(op) & 0777777) << 18)
}

func (ir *IndexRegister) SetXI12(op uint64) {
	*ir = (*ir & 077777777) | ((IndexRegister(op) & 07777) << 24)
}

func (ir *IndexRegister) SetXM(op uint64) {
	*ir = (*ir & 0777777000000) | (IndexRegister(op) & 0777777)
}

func (ir *IndexRegister) SetXM24(op uint64) {
	*ir = (*ir & 0777700000000) | (IndexRegister(op) & 077777777)
}

func (ir *IndexRegister) SetW(op uint64) {
	*ir = IndexRegister(op & 0_777777_777777)
}
