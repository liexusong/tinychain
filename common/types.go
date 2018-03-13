/*
HyperDCDN License
Copyright (C) 2017 The HyperDCDN Authors.
*/
package common

const (
	HashLength    = 32
	AddressLength = 20
)

type Hash [HashLength]byte

func (h Hash) String() string {
	return string(h)
}

type Address [AddressLength]byte
