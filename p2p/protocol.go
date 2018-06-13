package p2p

import (
	"tinychain/p2p/pb"
	"fmt"
	"errors"
)

var (
	ErrDupHandler = errors.New("p2p handler duplicate")
)

type Protocol struct {
	// Typ should match the message type
	Typ string

	// Run func handles the message from the stream
	Run func(message *pb.Message) error

	// Error func handles the error returned from the stream
	Error func(error)
}

func (h *Protocol) String() string {
	return fmt.Sprintf("P2P Handler %s", h.Typ)
}

func (p *Peer) AddProtocol(h *Protocol) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if handlers, exist := p.protocols[h.Typ]; exist {
		for _, handler := range handlers {
			if handler == h {
				return ErrDupHandler
			}
		}
		p.protocols[h.Typ] = append(handlers, h)
	} else {
		p.protocols[h.Typ] = []*Protocol{h}
	}
	return nil
}

func (p *Peer) DelProtocol(h *Protocol) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if handlers, exist := p.protocols[h.Typ]; exist {
		for i, handler := range handlers {
			if handler == h {
				handlers = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
		p.protocols[h.Typ] = handlers
	}
}
