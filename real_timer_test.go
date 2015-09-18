package tme

import (
	"testing"
	"time"
)

func TestRealTimerDone(t *testing.T) {
	now := time.Now()
	timer := NewRealTimer(dur, func() {})
	defer timer.Stop()
	select {
	case done := <-timer.Done():
		if done.Time.After(now.Add(dur * 2)) {
			t.Fatalf("timer ended after %s", dur)
		}
	case <-time.After(dur * 2):
		t.Fatalf("timer didn't end before %s", dur*2)
	}
}

func TestRealTimerMultipleStop(t *testing.T) {
	timer := NewRealTimer(dur, func() {})
	if !timer.Stop() {
		t.Fatal("first stop call returned false")
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			if timer.Stop() {
				t.Errorf("stop # %d returned true", i)
			}
		}(i)
	}
}

func TestRealTimerDoneAfterStop(t *testing.T) {
	timer := NewRealTimer(dur, func() {})
	if !timer.Stop() {
		t.Fatal("stop returned false")
	}
	ch := timer.Done()
	select {
	case <-time.After(dur):
	case <-ch:
		t.Errorf("channel returned after stop called")
	}
}

func TestRealTimerAck(t *testing.T) {
	ackCh := make(chan struct{})
	timer := NewRealTimer(dur, func() {
		ackCh <- struct{}{}
	})
	defer timer.Stop()
	ch := timer.Done()
	select {
	case ack := <-ch:
		go func() { ack.Fn() }()
	case <-time.After(dur * 2):
		t.Errorf("timer didn't stop after %s", dur*2)
	}

	select {
	case <-ackCh:
	case <-time.After(dur):
		t.Errorf("ack chan didn't receive after %s", dur)
	}
}
