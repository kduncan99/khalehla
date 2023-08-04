// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package pkg

import "fmt"

//	Any code which sets the discrete values should use the Set* methods
//	to ensure that the non-significant bits are zero.

type AccessKey struct {
	domain uint64
	ring   uint64
}

func (key *AccessKey) equals(op *AccessKey) bool {
	return (key.ring == op.ring) && (key.domain == op.domain)
}

func (key *AccessKey) GetComposite() uint64 {
	return (key.ring << 16) | key.domain
}

func (key *AccessKey) GetDomain() uint64 {
	return key.domain
}

func (key *AccessKey) GetRing() uint64 {
	return key.ring
}

func (key *AccessKey) GetString() string {
	return fmt.Sprintf("Ring:%v Domain:%06o", key.ring, key.domain)
}

func (key *AccessKey) IsMasterKey() bool {
	return (key.ring == 0) && (key.domain == 0)
}

func (key *AccessKey) SetDomain(value uint64) *AccessKey {
	key.domain = value & 03
	return key
}

func (key *AccessKey) SetRing(value uint64) *AccessKey {
	key.ring = value & 0xFFFF
	return key
}

func (key *AccessKey) SetComposite(composite uint64) *AccessKey {
	key.ring = (composite >> 16) & 03
	key.domain = composite & 0xFFFF
	return key
}

func NewAccessKey() *AccessKey {
	return &AccessKey{
		ring:   0,
		domain: 0,
	}
}

func NewAccessKeyFromComponents(ring uint64, domain uint64) *AccessKey {
	return &AccessKey{
		ring:   ring,
		domain: domain,
	}
}

func NewAccessKeyFromComposite(value uint64) *AccessKey {
	ak := AccessKey{}
	ak.SetComposite(value)
	return &ak
}
