// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package common

import (
	"fmt"
)

type AccessLock struct {
	domain uint64
	ring   uint64
}

func (lock *AccessLock) equals(op *AccessLock) bool {
	return (lock.ring == op.ring) && (lock.domain == op.domain)
}

func (lock *AccessLock) GetDomain() uint64 {
	return lock.domain
}

func (lock *AccessLock) GetRing() uint64 {
	return lock.ring
}

func (lock *AccessLock) GetComposite() uint64 {
	return (lock.ring << 16) | lock.domain
}

func (lock *AccessLock) GetEffectivePermissions(
	key *AccessKey,
	special *AccessPermissions,
	general *AccessPermissions) *AccessPermissions {

	if key.IsMasterKey() {
		return &AllAccessPermissions
	} else if (key.ring < lock.ring) || (key.domain == lock.domain) {
		return special
	} else {
		return general
	}
}

func (lock *AccessLock) GetString() string {
	return fmt.Sprintf("Ring:%v Domain:%06o", lock.ring, lock.domain)
}

func (lock *AccessLock) SetDomain(value uint64) *AccessLock {
	lock.domain = value & 03
	return lock
}

func (lock *AccessLock) SetRing(value uint64) *AccessLock {
	lock.ring = value & 0xFFFF
	return lock
}

func (lock *AccessLock) SetComposite(composite uint64) {
	lock.ring = (composite >> 16) & 03
	lock.domain = composite & 0xFFFF
}

func NewAccessLockFromComposite(composite uint64) *AccessLock {
	al := AccessLock{}
	al.SetComposite(composite)
	return &al
}

func NewAccessLock(ring uint64, domain uint64) *AccessLock {
	lock := AccessLock{}
	lock.SetRing(ring).SetDomain(domain)
	return &lock
}
