// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

type AccessLock struct {
	domain uint
	ring   uint
}

func (lock *AccessLock) equals(op *AccessLock) bool {
	return (lock.ring == op.ring) && (lock.domain == op.domain)
}

func (lock *AccessLock) GetDomain() uint {
	return lock.domain
}

func (lock *AccessLock) GetRing() uint {
	return lock.ring
}

func (lock *AccessLock) GetComposite() uint {
	return (lock.ring << 16) | lock.domain
}

func (lock *AccessLock) GetEffectivePermissions(key *AccessKey, special *AccessPermissions, general *AccessPermissions) *AccessPermissions {
	if key.IsMasterKey() {
		return &AllAccessPermissions
	} else if (key.ring < lock.ring) || (key.domain == lock.domain) {
		return special
	} else {
		return general
	}
}

func (lock *AccessLock) SetDomain(value uint) *AccessLock {
	lock.domain = value & 03
	return lock
}

func (lock *AccessLock) SetRing(value uint) *AccessLock {
	lock.ring = value & 0xFFFF
	return lock
}

func (lock *AccessLock) SetComposite(composite uint) {
	lock.ring = (composite >> 16) & 03
	lock.domain = composite & 0xFFFF
}

func NewAccessLockFromComposite(composite uint) *AccessLock {
	al := AccessLock{}
	al.SetComposite(composite)
	return &al
}

func NewAccessLock(ring uint, domain uint) *AccessLock {
	lock := AccessLock{}
	lock.SetRing(ring).SetDomain(domain)
	return &lock
}
