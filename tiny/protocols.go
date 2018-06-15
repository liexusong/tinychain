package tiny

import (
	"tinychain/common"
	"tinychain/core/types"
	"tinychain/event"
	"tinychain/p2p"
)

type TxPool interface {
	// AddRemotes adds remote transactions to queue tx list
	AddRemotes(txs types.Transactions) error

	// Pending returns all valid and processable transactions
	Pending() map[common.Address]types.Transactions

}

// ProtocolManager manages the subProtocols of p2p layer, and collects
// txs from remote peers
type ProtocolManager struct {
	subProtocols []p2p.Protocol

	net Network

	eventHub *event.TypeMux

	protocolSub event.Subscription

	quitCh chan struct{}
}

func NewProtocolManager(net Network) *ProtocolManager {
	return &ProtocolManager{
		net:      net,
		eventHub: event.GetEventhub(),
		quitCh:   make(chan struct{}),
	}
}

func (pm *ProtocolManager) Start() {
	pm.protocolSub = pm.eventHub.Subscribe()
	go pm.listen()
}

func (pm *ProtocolManager) listen() {
	for {
		select {
		case ev := <-pm.protocolSub.Chan():
			protoEv := ev.(*event.ProtocolEvent)
			if protoEv.Typ == "add" {
				pm.addProtocol(protoEv.Protocol)
			} else {
				pm.delProtocol(protoEv.Protocol)
			}
		case <-pm.quitCh:
			pm.protocolSub.Unsubscribe()
			break
		}
	}
}

func (pm *ProtocolManager) Stop() {
	close(pm.quitCh)
}

func (pm *ProtocolManager) Protocols() []p2p.Protocol {
	return pm.subProtocols
}

func (pm *ProtocolManager) addProtocol(proto p2p.Protocol) error {
	return pm.net.AddProtocol(proto)
}

func (pm *ProtocolManager) delProtocol(proto p2p.Protocol) {
	pm.net.DelProtocol(proto)
}
