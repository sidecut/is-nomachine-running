package main

// From https://stackoverflow.com/a/64536953/53107

import "sync"

type Subscriber[Message any] struct {
}

// Broadcaster is the struct which encompasses broadcasting
type Broadcaster[Message any] struct {
	cond        *sync.Cond
	subscribers map[Subscriber[Message]]func(Message)
	message     Message
	running     bool
}

// SetupBroadcaster gives the broadcaster object to be used further in messaging
func SetupBroadcaster[Message any]() *Broadcaster[Message] {

	return &Broadcaster[Message]{
		cond:        sync.NewCond(&sync.RWMutex{}),
		subscribers: map[Subscriber[Message]]func(Message){},
	}
}

// Subscribe let others enroll in broadcast event!
func (b *Broadcaster[Message]) Subscribe(id Subscriber[Message], f func(input Message)) {

	b.subscribers[id] = f
}

// Unsubscribe stop receiving broadcasting
func (b *Broadcaster[Message]) Unsubscribe(id Subscriber[Message]) {
	b.cond.L.Lock()
	delete(b.subscribers, id)
	b.cond.L.Unlock()
}

// Publish publishes the message
func (b *Broadcaster[Message]) Publish(message Message) {
	go func() {
		b.cond.L.Lock()

		b.message = message
		b.cond.Broadcast()
		b.cond.L.Unlock()
	}()
}

// Start the main broadcasting event
func (b *Broadcaster[Message]) Start() {
	b.running = true
	for b.running {
		b.cond.L.Lock()
		b.cond.Wait()
		go func() {
			for _, f := range b.subscribers {
				f(b.message) // publishes the message
			}
		}()
		b.cond.L.Unlock()
	}

}

// Stop broadcasting event
func (b *Broadcaster[Message]) Stop() {
	b.running = false
}
