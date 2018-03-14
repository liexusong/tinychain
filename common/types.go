/*
HyperDCDN License
Copyright (C) 2017 The HyperDCDN Authors.
*/
package common

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
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

func (h Hash) Hex() string {
	enc := make([]byte, len(h)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], h[:])
	return string(enc)
}

func DecodeHash(data []byte) Hash {
	var h Hash
	dec := make([]byte, HashLength)
	hex.Decode(dec, data[2:])
	return h.SetBytes(dec)
}

func (h Hash) SetBytes(d []byte) Hash {
	if len(d) > len(h) {
		d = d[:HashLength]
	}
	copy(h[:], d)
	return h
}

func Sha256(d []byte) Hash {
	var h Hash
	sha := sha3.New256()
	h.SetBytes(sha.Sum(d))
	return h
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
	sha := sha3.New256()
	hash := sha.Sum(enc)
	return "0x" + string(hash)
}

func (addr Address) SetBytes(b []byte) Address {
	if len(b) > len(addr) {
		b = b[:AddressLength]
	}
	copy(addr[:], b)
	return addr
}

func HashToAddr(hash Hash) Address {
	var addr Address
	addr.SetBytes(hash[:AddressLength])
	return addr
}

func DecodeAddr(d []byte) Address {
	var addr Address
	dec := make([]byte, AddressLength)
	hex.Decode(dec, d[2:])
	return addr.SetBytes(dec)
}
