package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-peer"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	libnet "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"tinychain/p2p/pb"
	"sync"
	"time"
	"github.com/pkg/errors"
	"tinychain/common"
)

var (
	TransProtocol = "/chain/1.0.0."
	log           = common.GetLogger("p2p")
)

// Peer stands for a logical peer of p2p layer
type Peer struct {
	host       *bhost.BasicHost // Local peer host
	routeTable *RouteTable      // Local route table
	context    context.Context
	respCh     chan *pb.Message // Response channel. Receive message from stream.
	quitCh     chan struct{}

	timeout time.Duration // Timeout of per connection
}

// Creates new peer struct
func NewPeer(config *Config) (*Peer, error) {
	host, err := NewHost(config.port, config.privKey)
	if err != nil {
		log.Errorf("Cannot create host:%s", err)
		return nil, err
	}

	peer := &Peer{
		host:    host,
		context: context.Background(),
		respCh:  make(chan *pb.Message, 100),
		quitCh:  make(chan struct{}),
		timeout: time.Second * 60,
	}
	peer.routeTable = NewRouteTable(config, peer)

	return peer, nil
}

func (peer *Peer) ID() peer.ID {
	return peer.host.ID()
}

// Link to a unknown peer with its multiaddr
// and send handshake
func (peer *Peer) Link(addr ma.Multiaddr) {
	// TODO
}

// Connect to a peer
func (peer *Peer) Connect(pid peer.ID) error {
	ctx, cancel := context.WithTimeout(peer.context, peer.timeout) //TODO: configurable?
	defer cancel()
	err := peer.host.Connect(ctx, pstore.PeerInfo{ID: pid})
	if err != nil {
		log.Infof("Failed to connect to peer %s\n", pid.Pretty())
		return err
	}
	return nil
}

// Send message to a peer
func (peer *Peer) SendMessage(pid peer.ID, name string, data interface{}) error {
	if pid == peer.ID() {
		log.Info("Cannot send message to peer itself.")
		return errors.New("Send message to self")
	}
	err := peer.Connect(pid)
	if err != nil {
		return err
	}
	stream := NewStreamWithPid(pid, peer)
	//peer.Streams.AddStream(stream)
	return stream.SendMessage(name, data)
}

func (peer *Peer) Start() {
	log.Infof("Peer start with pid %s. Listen on addr: %s.\n", peer.ID().Pretty(), peer.host.Addrs())
	// Listen to stream arriving
	peer.host.SetStreamHandler(TransProtocol, peer.onStreamConnected)

	// Sync route with seeds and neighbor
	peer.routeTable.Start()

	go peer.ListenMsg()
}

func (peer *Peer) Stop() {
	if peer.host != nil {
		peer.host.Network().Close()
		peer.host.Close()
	}
	peer.routeTable.Stop()
	peer.quitCh <- struct{}{}
}

func (peer *Peer) ListenMsg() {
	for {
		select {
		case <-peer.quitCh:
			break
		case message := <-peer.respCh:
			log.Infof("Receive message: Name:%s, data:%s \n", message.Name, message.Data)
		}
	}
}

func (peer *Peer) onStreamConnected(s libnet.Stream) {
	log.Infof("Receive stream. Peer is:%s\n", s.Conn().RemotePeer().String())
	stream := NewStream(s.Conn().RemotePeer(), s.Conn().RemoteMultiaddr(), s, peer)

	//peer.Streams.AddStream(stream)
	stream.StartLoop()
}

func (peer *Peer) Broadcast(pbName string, data interface{}) {
	for pid := range peer.routeTable.Peers() {
		if pid == peer.ID() {
			continue
		}
		stream := NewStreamWithPid(pid, peer)
		go stream.SendMessage(pbName, data)
	}
}
