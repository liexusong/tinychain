/*
HyperDCDN License
Copyright (C) 2017 The HyperDCDN Authors.
*/
package common

import (
	"encoding/hex"
	"crypto/sha256"
)

const (
	HashLength    = 32
	AddressLength = 20
)

type Hash [HashLength]byte

func (h Hash) String() string {
	return string(h[:])
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) Hex() []byte {
	return Hex(h[:])
}

// Decode hash string with "0x...." format to Hash type
func DecodeHash(data []byte) Hash {
	dec := make([]byte, HashLength)
	hex.Decode(dec, data[2:])
	return BytesToHash(dec)
}

func BytesToHash(d []byte) Hash {
	var h Hash
	if len(d) > HashLength {
		d = d[:HashLength]
	}
	copy(h[:], d)
	return h
}

func (h Hash) Nil() bool {
	return h == Hash{}
}

func Sha256(d []byte) Hash {
	return sha256.Sum256(d)
}

type Address [AddressLength]byte

func (addr Address) String() string {
	return string(addr[:])
}

func (addr Address) Bytes() []byte {
	return addr[:]
}

func (addr Address) Hex() string {
	enc := make([]byte, len(addr)*2)
	hex.Encode(enc, addr[:])
	hash := Sha256(enc)
	return "0x" + hash.String()
}

func (addr Address) BytesToHash(b []byte) Address {
	if len(b) > len(addr) {
		b = b[:AddressLength]
	}
	copy(addr[:], b)
	return addr
}

func HashToAddr(hash Hash) Address {
	var addr Address
	addr.BytesToHash(hash[:AddressLength])
	return addr
}

func DecodeAddr(d []byte) Address {
	var addr Address
	dec := make([]byte, AddressLength)
	hex.Decode(dec, d[2:])
	return addr.BytesToHash(dec)
}