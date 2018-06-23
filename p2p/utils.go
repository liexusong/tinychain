package p2p

import (
	ma "github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-host"
	"strings"
	"fmt"
	"github.com/libp2p/go-libp2p-crypto"
	"crypto/rand"
	"errors"
)

// Parse ipfs addr like '/ip4/127.0.0.1/tcp/65532/ipfs/QmWxRLJvALbQRqE8ay91e5kzeNuPkNZQ4UkvXRrPxfXFuX'
// Split it to '/ip4/127.0.0.1/tcp/65532' and 'QmWxRLJvALbQRqE8ay91e5kzeNuPkNZQ4UkvXRrPxfXFuX'(peer id)
func ParseFromIPFSAddr(ipfsAddr ma.Multiaddr) (peer.ID, ma.Multiaddr, error) {
	addr, err := ma.NewMultiaddr(strings.Split(ipfsAddr.String(), "/ipfs/")[0])
	if err != nil {
		return "", nil, err
	}

	pid, err := ipfsAddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		return "", nil, err
	}

	id, err := peer.IDB58Decode(pid)
	if err != nil {
		return "", nil, err
	}

	return id, addr, nil
}

// Reverse func of ParseFromIPFSAddr
func GetCompleteAddrs(h host.Host) ([]ma.Multiaddr, error) {
	var comAddrs []ma.Multiaddr
	addrs := h.Addrs()
	for _, addr := range addrs {
		c, err := ma.NewMultiaddr(fmt.Sprintf("%s/ipfs/%s", addr, peer.IDB58Encode(h.ID())))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid multiaddr:%s\n", addr))
		}
		comAddrs = append(comAddrs, c)
	}
	return comAddrs, nil
}

// Base64 encode private_key and return a b64 string
func B64EncodePrivKey(priv crypto.PrivKey) (string, error) {
	privKey, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return "", err
	}
	return crypto.ConfigEncodeKey(privKey), nil
}

// Base64 decode and return private_key
func B64DecodePrivKey(data string) (crypto.PrivKey, error) {
	b64, err := crypto.ConfigDecodeKey(data)
	if err != nil {
		return nil, err
	}
	privKey, err := crypto.UnmarshalPrivateKey(b64)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

// Random generate a peer id
func RandomGeneratePid() (peer.ID, error) {
	_, pub, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return "", err
	}
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return "", err
	}
	return pid, nil
}
