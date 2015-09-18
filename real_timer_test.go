package tme

import (
	"testing"
	"time"
)

func TestRealTimerDone(t *testing.T) {
	now := time.Now()
	timer := NewRealTimer(dur)
	select {
	case done := <-timer.Done():
		if done.Time.After(now.Add(dur * 2)) {
			t.Fatalf("timer ended after %s", dur)
		}
	case <-time.After(dur * 2):
		t.Fatalf("timer didn't end before %s", dur*2)
	}
}

func TestMultipleStop(t *testing.T) {
	timer := NewRealTimer(dur)
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

func TestDoneAfterStop(t *testing.T) {
	timer := NewRealTimer(dur)
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
