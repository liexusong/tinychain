package common

import (
	"encoding/hex"
	"github.com/libp2p/go-libp2p-crypto"
)

func Hex(b []byte) []byte {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return enc
}

// Generate address by public key
func GenAddrByPubkey(key crypto.PubKey) (Address, error) {
	var addr Address
	pubkey, err := key.Bytes()
	if err != nil {
		return addr, err
	}
	pubkey = pubkey[1:]
	h := Sha256(pubkey)
	hash := h[len(h)-AddressLength:]
	addr = HashToAddr(Sha256(hash))
	return addr, nil
}

// Generate address by private key
func GenAddrByPrivkey(key crypto.PrivKey) (Address, error) {
	pubkey := key.GetPublic()
	return GenAddrByPubkey(pubkey)
}
