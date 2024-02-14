// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package mfdMgr

import "khalehla/kexec/types"

// recoverMassStorage handles MFD recovery for what is NOT a JK13 boot.
// If we return an error, we must previously stop the exec.
func (mgr *MFDManager) recoverMassStorage() error {
	// TODO
	mgr.exec.SendExecReadOnlyMessage("MFD Recovery is not implemented")
	mgr.exec.Stop(types.StopDirectoryErrors)
	return nil
}
