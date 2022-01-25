package main

// From https://stackoverflow.com/a/64536953/53107

import "sync"

// Broadcaster is the struct which encompasses broadcasting
type Broadcaster struct {
	cond        *sync.Cond
	subscribers map[interface{}]func(interface{})
	message     interface{}
	running     bool
}

// SetupBroadcaster gives the broadcaster object to be used further in messaging
func SetupBroadcaster() *Broadcaster {

	return &Broadcaster{
		cond:        sync.NewCond(&sync.RWMutex{}),
		subscribers: map[interface{}]func(interface{}){},
	}
}

// Subscribe let others enroll in broadcast event!
func (b *Broadcaster) Subscribe(id interface{}, f func(input interface{})) {

	b.subscribers[id] = f
}

// Unsubscribe stop receiving broadcasting
func (b *Broadcaster) Unsubscribe(id interface{}) {
	b.cond.L.Lock()
	delete(b.subscribers, id)
	b.cond.L.Unlock()
}

// Publish publishes the message
func (b *Broadcaster) Publish(message interface{}) {
	go func() {
		b.cond.L.Lock()

		b.message = message
		b.cond.Broadcast()
		b.cond.L.Unlock()
	}()
}

// Start the main broadcasting event
func (b *Broadcaster) Start() {
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
func (b *Broadcaster) Stop() {
	b.running = false
}
