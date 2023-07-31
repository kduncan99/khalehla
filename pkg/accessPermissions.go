// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

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

func (perm *AccessPermissions) GetComposite() uint {
	var value uint
	if perm.canEnter {
		value |= 04
	}
	if perm.canRead {
		value |= 02
	}
	if perm.canWrite {
		value |= 01
	}
	return value
}

func (perm *AccessPermissions) GetString() string {
	str := "Permissions:"

	if perm.canEnter {
		str += "E"
	} else {
		str += "-"
	}

	if perm.canRead {
		str += "R"
	} else {
		str += "-"
	}

	if perm.canWrite {
		str += "W"
	} else {
		str += "-"
	}

	return str
}

func NewAccessPermissions(canEnter bool, canRead bool, canWrite bool) *AccessPermissions {
	perm := AccessPermissions{
		canEnter: canEnter,
		canRead:  canRead,
		canWrite: canWrite,
	}
	return &perm
}
