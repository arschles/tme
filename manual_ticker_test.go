package tme

import (
	"testing"
	"time"
)

func TestManualTickerTick(t *testing.T) {
	ticker := NewManualTicker(func() {})
	defer ticker.Stop()
	go func() {
		ticker.Tick()
	}()
	select {
	case <-ticker.Chan():
	case <-time.After(dur):
		t.Errorf("didn't tick within %s", dur)
	}
}

func TestManualTickerTickAfterStop(t *testing.T) {
	ticker := NewManualTicker(func() {})
	ticker.Stop()
	go func() {
		ticker.Tick()
	}()
	select {
	case <-ticker.Chan():
		t.Errorf("ticker ticked after stop")
	case <-time.After(dur):
	}
}

func TestManualTickerAck(t *testing.T) {
	ackCh := make(chan struct{})
	ticker := NewManualTicker(func() {
		ackCh <- struct{}{}
	})
	defer ticker.Stop()
	go func() {
		ticker.Tick()
	}()
	select {
	case ack := <-ticker.Chan():
		go func() { ack.Fn() }()
	case <-time.After(dur):
		t.Errorf("ticker didn't tick within %s", dur)
	}

	select {
	case <-ackCh:
	case <-time.After(dur):
		t.Errorf("ack wasn't received within %s", dur)
	}
}

func TestManualTickerMultipleRecv(t *testing.T) {
	const n = 10
	ticker := NewManualTicker(func() {})
	go func() {
		ticker.Tick()
	}()
	defer ticker.Stop()
	numRecvs := numTickerRecv(ticker, 10)
	if numRecvs != 1 {
		t.Errorf("%d ticker receives", numRecvs)
	}
}
