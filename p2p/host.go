package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-swarm"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p-peer"
	ps "github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-crypto"

	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	"crypto/rand"
	"fmt"
)

// NewHost construct a host of libp2p
func NewHost(port int, privKey crypto.PrivKey) (*bhost.BasicHost, error) {
	var (
		priv crypto.PrivKey
		pub  crypto.PubKey
	)
	if privKey == nil {
		var err error
		//priv, pub, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
		priv, pub, err = crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		pubB64, _ := peer.IDFromPublicKey(pub)
		privB64, _ := B64EncodePrivKey(priv)

		log.Info("Private key not found in config file.")
		log.Info("Generate new key pair of privkey and pubkey by Ed25519:")
		log.Infof("Pubkey:%s\n", pubB64)
		log.Infof("Privkey:%s\n", privB64)
	} else {
		priv = privKey
		pub = privKey.GetPublic()
	}
	//privKey, _ := crypto.MarshalPrivateKey(priv)
	//data := crypto.ConfigEncodeKey(privKey)
	//log.Info(data)
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	listener, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	if err != nil {
		log.Infof("New multiaddr:%s\n", err)
		log.Error(err)
		return nil, err
	}

	pstore := ps.NewPeerstore()
	pstore.AddPrivKey(pid, priv)
	pstore.AddPubKey(pid, pub)

	ctx := context.Background()
	n, err := swarm.NewNetwork(ctx,
		[]ma.Multiaddr{listener}, pid, pstore, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	opts := &bhost.HostOpts{
		NATManager: bhost.NewNATManager(n),
	}
	return bhost.NewHost(ctx, n, opts)
}

//func NewTcpHost(addr string, port int) (p2phost.Host, error) {
//	addr = log.Sprintf("/ip4/%s/tcp/%d", addr, port)
//	return NewHost(addr)
//}
//
//func NewIpfsHost() (p2phost.Host, error) {
//	return NewHost()
//}
