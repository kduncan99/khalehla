// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

type AccessPermissions struct {
	canEnter bool
	canRead  bool
	canWrite bool
}

var AllAccessPermissions = AccessPermissions{
	canEnter: true,
	canRead:  true,
	canWrite: true,
}

func (perm *AccessPermissions) CanEnter() bool {
	return perm.canEnter
}

func (perm *AccessPermissions) CanRead() bool {
	return perm.canRead
}

func (perm *AccessPermissions) CanWrite() bool {
	return perm.canWrite
}

func NewAccessPermissions(canEnter bool, canRead bool, canWrite bool) *AccessPermissions {
	perm := AccessPermissions{
		canEnter: canEnter,
		canRead:  canRead,
		canWrite: canWrite,
	}
	return &perm
}
