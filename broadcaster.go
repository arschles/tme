package tme

import "sync"

type ackBroadcaster struct {
	sync.RWMutex
	chans []chan Ack
}

func newAckBroadcaster() *ackBroadcaster {
	return &ackBroadcaster{}
}

func (a *ackBroadcaster) addChan(ch chan Ack) {
	a.Lock()
	defer a.Unlock()
	a.chans = append(a.chans, ch)
}

func (a *ackBroadcaster) broadcast(ack Ack) {
	var wg sync.WaitGroup
	for _, ch := range a.chans {
		wg.Add(1)
		go func(ch chan Ack) {
			defer wg.Done()
			ch <- ack
		}(ch)
	}
	wg.Wait()
}

func (a *ackBroadcaster) closeAllChans() {
	for _, ch := range a.chans {
		close(ch)
	}
}
