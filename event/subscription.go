package event

type Subscription interface {
	Unsubscribe()             // Remove the target channel, and close err channel
	Chan() <-chan interface{} // Get the channel
}

type muxSub struct {
	mu      *TypeMux
	channel chan interface{}
	closed  bool
}

func newMuxSub(mu *TypeMux) *muxSub {
	return &muxSub{
		mu:      mu,
		channel: make(chan interface{}),
	}
}

func (ms *muxSub) Unsubscribe() {
	ms.mu.del(ms)
}

func (ms *muxSub) Chan() <-chan interface{} {
	return ms.channel
}

func (ms *muxSub) ch() chan interface{} {
	return ms.channel
}

// Notice: the caller should holds the lock,
// and must be called by feed
func (ms *muxSub) close() {
	if ms.closed {
		return
	}
	close(ms.channel)
	ms.closed = true
}
