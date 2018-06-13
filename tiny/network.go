package tiny

import (
	"tinychain/p2p"
	"tinychain/event"
)

// Network is the wrapper of physical p2p network layer
type Network struct {
	config  *Config
	network *p2p.Peer
	event   *event.TypeMux

	// Send message event subscription
	sendSub event.Subscription
	// Add and delete protocol event subscription
	protocolSub event.Subscription

	quitCh chan struct{}
}

func NewNetwork(config *Config) *Network {
	network, err := p2p.New(config.p2p)
	if err != nil {
		return nil
	}
	return &Network{
		config:  config,
		network: network,
		event:   event.GetEventhub(),
		quitCh:  make(chan struct{}),
	}
}

func (p *Network) Start() error {
	p.sendSub = p.event.Subscribe(&event.SendMsgEvent{})
	p.protocolSub = p.event.Subscribe(&event.ProtocolEvent{})
	go p.listen()
	return nil
}

func (p *Network) listen() {
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

func (p *Network) Stop() error {
	close(p.quitCh)
	return nil
}
