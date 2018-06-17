package tiny

import (
	"tinychain/event"
	"tinychain/p2p"
	"sync"
)

// ProtocolManager manages the subProtocols of p2p layer, and collects
// txs from remote peers
type ProtocolManager struct {
	mu           sync.RWMutex
	subProtocols []p2p.Protocol
	net          Network
	eventHub     *event.TypeMux
}

func NewProtocolManager(net Network) *ProtocolManager {
	return &ProtocolManager{
		net: net,
	}
}

func (pm *ProtocolManager) Init(protocols []p2p.Protocol) error {
	for _, protocol := range protocols {
		err := pm.AddProtocol(protocol)
		if err != nil {
			log.Errorf("faild to register protocol %s", protocol.Type())
			return nil
		}
	}
	return nil
}

func (pm *ProtocolManager) Protocols() []p2p.Protocol {
	return pm.subProtocols
}

func (pm *ProtocolManager) AddProtocol(proto p2p.Protocol) error {
	err := pm.net.AddProtocol(proto)
	if err != nil {
		return err
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.subProtocols = append(pm.subProtocols, proto)
	return nil
}

func (pm *ProtocolManager) DelProtocol(proto p2p.Protocol) {
	pm.net.DelProtocol(proto)
	pm.mu.Lock()
	defer pm.mu.Unlock()

}
