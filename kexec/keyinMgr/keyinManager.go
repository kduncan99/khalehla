// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"sync"
)

// KeyinManager handles all things related to unsolicited console keyins
type KeyinManager struct {
	exec  types.IExec
	mutex sync.Mutex
}

func NewKeyinManager(exec types.IExec) *KeyinManager {
	return &KeyinManager{
		exec: exec,
	}
}

// CloseManager is invoked when the exec is stopping
func (mgr *KeyinManager) CloseManager() {
}

func (mgr *KeyinManager) InitializeManager() {
}

// ResetManager clears out any artifacts left over by a previous exec session,
// and prepares the console for normal operations
func (mgr *KeyinManager) ResetManager() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
}

func (mgr *KeyinManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vKeyinManager ----------------------------------------------------\n", indent)
	// TODO
}
