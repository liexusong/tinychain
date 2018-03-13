package event

import (
	"sync"
	"errors"
)

var (
	eb               *EventBus
	ErrEventNotFound = errors.New("event not found")
)

type EventBus struct {
	subs sync.Map
}

type slot struct {
	channel chan interface{}
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (eb *EventBus) Sub(event string, channel chan interface{}) {
	slot := &slot{channel}
	eb.subs.Store(event, slot)
}

func (eb *EventBus) Unsub(event string, channel chan interface{}) error {
	slots, ok := eb.subs.Load(event)
	if !ok {
		return ErrEventNotFound
	}
	slotArr := slots.([]*slot)
	for i, ch := range slotArr {
		if ch.channel == channel {
			slotArr = append(slotArr[:i], slotArr[i+1:]...)
			break
		}
	}
	return nil
}

func (eb *EventBus) Notify(event string) error {
	slots, ok := eb.subs.Load(event)
	if !ok {
		return ErrEventNotFound
	}
	for _, ch := range slots.([]*slot) {
		ch.channel <- struct{}{}
	}
	return nil
}

func GetEventBus() *EventBus {
	if eb != nil {
		return eb
	}
	eb = NewEventBus()
	return eb
}
