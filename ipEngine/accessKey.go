// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package ipEngine

//	Any code which sets the discrete values should use the Set* methods
//	to ensure that the non-significant bits are zero.

type AccessKey struct {
	domain uint
	ring   uint
}

func (key *AccessKey) equals(op *AccessKey) bool {
	return (key.ring == op.ring) && (key.domain == op.domain)
}

func (key *AccessKey) GetComposite() uint {
	return (key.ring << 16) | key.domain
}

func (key *AccessKey) IsMasterKey() bool {
	return (key.ring == 0) && (key.domain == 0)
}

func (key *AccessKey) SetDomain(value uint) *AccessKey {
	key.domain = value & 03
	return key
}

func (key *AccessKey) SetRing(value uint) *AccessKey {
	key.ring = value & 0xFFFF
	return key
}

func (key *AccessKey) SetComposite(composite uint) *AccessKey {
	key.ring = (composite >> 16) & 03
	key.domain = composite & 0xFFFF
	return key
}

func NewAccessKeyFromComponents(ring uint, domain uint) *AccessKey {
	ak := AccessKey{}
	ak.SetRing(ring).SetDomain(domain)
	return &ak
}

func NewAccessKeyFromComposite(value uint) *AccessKey {
	ak := AccessKey{}
	ak.SetComposite(value)
	return &ak
}
