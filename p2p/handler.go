package p2p

import (
	"tinychain/p2p/pb"
	"fmt"
	"errors"
)

var (
	ErrDupHandler = errors.New("p2p handler duplicate")
)

type Handler struct {
	// Name should match the message type
	Name string

	// Run func handles the message from the stream
	Run func(message *pb.Message) error

	// Error func handles the error returned from the stream
	Error func(error)
}

func (h *Handler) String() string {
	return fmt.Sprintf("P2P Handler %s", h.Name)
}

func (p *Peer) AddHandler(h *Handler) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if handlers, exist := p.handlers[h.Name]; exist {
		for _, handler := range handlers {
			if handler == h {
				return ErrDupHandler
			}
		}
		p.handlers[h.Name] = append(handlers, h)
	} else {
		p.handlers[h.Name] = []*Handler{h}
	}
	return nil
}

func (p *Peer) RemoveHandler(h *Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if handlers, exist := p.handlers[h.Name]; exist {
		for i, handler := range handlers {
			if handler == h {
				handlers = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
		p.handlers[h.Name] = handlers
	}
}
