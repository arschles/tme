package tme

import (
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
	const dur = 20 * time.Millisecond
	chans := []chan Ack{make(chan Ack), make(chan Ack), make(chan Ack)}
	b := newAckBroadcaster()
	for _, ch := range chans {
		b.addChan(ch)
	}

	go func() {
		b.broadcast(Ack{Time: time.Now(), Fn: func() {}})
	}()

	for i, ch := range chans {
		select {
		case <-ch:
		case <-time.After(dur):
			t.Errorf("chan %d didn't receive within %s", i, dur)
		}
	}
}

func TestCloseAllChans(t *testing.T) {
	const dur = 10 * time.Millisecond
	chans := []chan Ack{make(chan Ack), make(chan Ack), make(chan Ack)}
	b := newAckBroadcaster()
	for _, ch := range chans {
		b.addChan(ch)
	}
	b.closeAllChans()
	for i, ch := range chans {
		select {
		case <-ch:
		case <-time.After(dur):
			t.Errorf("chan %d didn't receive after %s", i, dur)
		}
	}
}
