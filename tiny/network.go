package tiny

import (
	"tinychain/p2p"
	"tinychain/event"
)

type Network interface {
	Start() error
	Stop() error
	AddProtocol(p2p.Protocol) error
	DelProtocol(p2p.Protocol)
}

// Network is the wrapper of physical p2p network layer
type Peer struct {
	config  *p2p.Config
	network *p2p.Peer
	event   *event.TypeMux

	// Send message event subscription
	sendSub event.Subscription
	// Add and delete protocol event subscription
	protocolSub event.Subscription

	quitCh chan struct{}
}

func NewNetwork(config *p2p.Config) Network {
	network, err := p2p.New(config)
	if err != nil {
		log.Error("Failed to create p2p peer")
		return nil
	}
	return &Peer{
		config:  config,
		network: network,
		event:   event.GetEventhub(),
		quitCh:  make(chan struct{}),
	}
}

func (p *Peer) Start() error {
	p.sendSub = p.event.Subscribe(&event.SendMsgEvent{})
	p.protocolSub = p.event.Subscribe(&event.ProtocolEvent{})
	go p.listen()
	return nil
}

func (p *Peer) listen() {
	for {
		select {
		case ev := <-p.sendSub.Chan():
			msg := ev.(*event.SendMsgEvent)
			err := p.network.Send(msg.Target, msg.Typ, msg.Data)
			if err != nil {
				log.Errorf("Failed to send message to %s with type %s", msg.Target, msg.Typ)
			}
		case ev := <-p.protocolSub.Chan():
			action := ev.(*event.ProtocolEvent)
			if action.Typ == "add" {
				p.network.AddProtocol(action.Protocol)
			} else if action.Typ == "del" {
				p.network.DelProtocol(action.Protocol)
			} else {
				log.Errorf("Unknown protocol action type %s", action.Typ)
			}
		case p.quitCh:
			p.sendSub.Unsubscribe()
			p.protocolSub.Unsubscribe()
		}
	}
}

func (p *Peer) Stop() error {
	close(p.quitCh)
	return nil
}

func (p *Peer) AddProtocol(proto p2p.Protocol) error {
	return p.network.AddProtocol(proto)
}

func (p *Peer) DelProtocol(proto p2p.Protocol) {
	p.network.DelProtocol(proto)
}
