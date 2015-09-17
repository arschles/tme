package tme

import "sync"

type signals struct {
	sync.RWMutex
	chans []chan struct{}
}

func newSignals() *signals {
	return &signals{}
}

// add registers ch to receive on every broadcast call
func (a *signals) add(ch chan struct{}) {
	a.Lock()
	a.chans = append(a.chans, ch)
	a.Unlock()
}

// broadcast asynchronously sends ack on all chans that have been registered with add
func (a *signals) broadcast() {
	a.RLock()
	defer a.RUnlock()
	for _, ch := range a.chans {
		go func(ch chan<- struct{}) {
			ch <- struct{}{}
		}(ch)
	}
}
