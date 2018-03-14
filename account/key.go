package account

import (
	"github.com/libp2p/go-libp2p-crypto"
	"crypto/rand"
)

type Key struct {
	priKey crypto.PrivKey
	pubKey crypto.PubKey
}

func NewKeyPairs() (*Key, error) {
	priv, pub, err := crypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Key{priv, pub}, nil
}
