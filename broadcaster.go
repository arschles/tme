package tme

import "sync"

type ackBroadcaster struct {
	sync.RWMutex
	chans []chan Ack
}

func newAckBroadcaster() *ackBroadcaster {
	return &ackBroadcaster{}
}

// addChan registers ch to be sent to on all future calls to broadcast
func (a *ackBroadcaster) addChan(ch chan Ack) {
	a.Lock()
	defer a.Unlock()
	a.chans = append(a.chans, ch)
}

// broadcast asynchronously sends ack on all chans that have been registered with addChan.
func (a *ackBroadcaster) broadcast(ack Ack) {
	a.RLock()
	defer a.RUnlock()
	for _, ch := range a.chans {
		go func(ch chan Ack) {
			ch <- ack
		}(ch)
	}
}

// close closes all channels registered with this broadcaster
func (a *ackBroadcaster) close() {
	a.RLock()
	defer a.RUnlock()
	for _, ch := range a.chans {
		close(ch)
	}
}
