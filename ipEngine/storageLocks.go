// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

import (
	"fmt"
	"khalehla/pkg"
	"sync"
	"time"
)

var waitTime = time.Millisecond

type StorageLockClient interface {
	GetStorageLockClientName() string
}

type storageLockKey uint64

func newStorageLockKey(address pkg.VirtualAddress) storageLockKey {
	return storageLockKey(address.GetComposite())
}

type StorageLocks struct {
	locks map[storageLockKey]StorageLockClient
	mutex sync.Mutex
}

// NewStorageLocks creates one of these things.
// There should be exactly one, shared among all the InstructionEngine instances.
// These are used to protect the integrity of the following instructions:
// ADD1, CR, DEC, DEC2, ENZ, INC, INC2, SUB1, TCS, TS, TSS
func NewStorageLocks() *StorageLocks {
	return &StorageLocks{
		locks: make(map[storageLockKey]StorageLockClient),
	}
}

func (sl *StorageLocks) Dump() {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	for value, client := range sl.locks {
		fmt.Printf("    %012o:%s\n", value, client.GetStorageLockClientName())
	}
}

func (sl *StorageLocks) Lock(address pkg.VirtualAddress, client StorageLockClient) bool {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	key := newStorageLockKey(address)
	_, ok := sl.locks[key]
	if ok {
		return false
	}

	sl.locks[key] = client
	return true
}

func (sl *StorageLocks) LockWait(address pkg.VirtualAddress, client StorageLockClient) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	for true {
		key := newStorageLockKey(address)
		_, ok := sl.locks[key]
		if ok {
			sl.mutex.Unlock()
			time.Sleep(waitTime)
			sl.mutex.Lock()
		} else {
			sl.locks[key] = client
			return
		}
	}
}

func (sl *StorageLocks) Release(address pkg.VirtualAddress, client StorageLockClient) bool {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	key := newStorageLockKey(address)
	lockClient, ok := sl.locks[key]
	if ok && lockClient == client {
		delete(sl.locks, key)
		return true
	} else {
		return false
	}
}

func (sl *StorageLocks) ReleaseAll(client StorageLockClient) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	for key, lockClient := range sl.locks {
		if lockClient == client {
			delete(sl.locks, key)
		}
	}
}
