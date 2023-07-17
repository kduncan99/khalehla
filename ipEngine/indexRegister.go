// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"khalehla/pkg"
)

type IndexRegister pkg.Word36

func (ir *IndexRegister) DecrementModifier() {
	ir.SetXM(pkg.AddSimple(ir.GetSignedXM(), pkg.Negate(ir.GetSignedXI())))
}

func (ir *IndexRegister) DecrementModifier24() {
	ir.SetXM24(pkg.AddSimple(ir.GetSignedXM24(), pkg.Negate(ir.GetSignedXI12())))
}

func (ir *IndexRegister) GetSignedXI() uint64 {
	return pkg.GetSignExtended18(ir.GetXI())
}

func (ir *IndexRegister) GetSignedXI12() uint64 {
	return pkg.GetSignExtended12(ir.GetXI12())
}

func (ir *IndexRegister) GetSignedXM() uint64 {
	return pkg.GetSignExtended18(ir.GetXM())
}

func (ir *IndexRegister) GetSignedXM24() uint64 {
	return pkg.GetSignExtended24(ir.GetXM24())
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
	ir.SetXM(pkg.AddSimple(ir.GetSignedXM(), ir.GetSignedXI()))
}

func (ir *IndexRegister) IncrementModifier24() {
	ir.SetXM24(pkg.AddSimple(ir.GetSignedXM24(), ir.GetSignedXI12()))
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
