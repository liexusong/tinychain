package common

import "encoding/hex"

func Hex(b []byte) []byte {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return enc
}