// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import "khalehla/pkg"

type Gate struct {
	generalAccessPermissions *pkg.AccessPermissions
	specialAccessPermissions *pkg.AccessPermissions
	libFlag                  bool
	gotoInhibit              bool
	designatorInhibit        bool
	accessKeyInhibit         bool
	latentParameter0Inhibit  bool
	latentParameter1Inhibit  bool
	accessLock               *pkg.AccessLock
	targetLevel              uint64
	targetBDI                uint64
	targetOffset             uint64
	basicModeBaseRegister    uint64                  // basic mode register is actually this field + 12
	designatorRegisterValue  *pkg.DesignatorRegister //	only bits 12-17 are significant
	newAccessKey             *pkg.AccessKey
	latentParameterValue0    uint64
	latentParameterValue1    uint64
}

func NewGateFromStorage(buffer []pkg.Word36) *Gate {
	g := Gate{}
	g.generalAccessPermissions = pkg.NewAccessPermissions(buffer[0]&0_400000_000000 != 0, false, false)
	g.specialAccessPermissions = pkg.NewAccessPermissions(buffer[0]&0_040000_000000 != 0, false, false)
	g.libFlag = buffer[0]&0_040_000000 != 0
	g.gotoInhibit = buffer[0]&0_020_000000 != 0
	g.designatorInhibit = buffer[0]&0_010_000000 != 0
	g.accessKeyInhibit = buffer[0]&0_004_000000 != 0
	g.latentParameter1Inhibit = buffer[0]&0_002_000000 != 0
	g.latentParameter1Inhibit = buffer[0]&0_001_000000 != 0
	g.accessLock = pkg.NewAccessLockFromComposite(buffer[0].GetW() & 0777777)
	g.targetLevel = buffer[1].GetW() >> 33
	g.targetBDI = buffer[1].GetH1() & 077777
	g.targetOffset = buffer[1].GetH2()
	g.basicModeBaseRegister = (buffer[2].GetW() >> 24) & 03
	g.designatorRegisterValue = pkg.NewDesignatorRegisterFromComposite(buffer[2].GetW() & 0_000077_000000)
	g.newAccessKey = pkg.NewAccessKeyFromComposite(buffer[2].GetH2())
	g.latentParameterValue0 = buffer[3].GetW()
	g.latentParameterValue1 = buffer[4].GetW()
	return &g
}
