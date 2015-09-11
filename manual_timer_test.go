package tme

import (
	"sync"
	"testing"
	"time"
)

func TestManualTimerMarkDone(t *testing.T) {
	const dur = 100 * time.Millisecond
	timer := NewManualTimer()
	bch := timer.Done()
	select {
	case <-bch:
		t.Errorf("timer was marked done when it wasn't done")
	case <-time.After(dur):
	}

	timer.MarkDone()
	ach := timer.Done()

	// bch should receive and close
	select {
	case <-bch:
	case <-time.After(dur):
		t.Errorf("timer was not marked done after %s", dur)
	}

	// ach should receive and close
	select {
	case <-ach:
	case <-time.After(dur):
		t.Errorf("Done didn't receive within %s", dur)
	}
}

func TestManualTimerStop(t *testing.T) {
	const dur = 100 * time.Millisecond
	timer := NewManualTimer()
	bch := timer.Done()
	timer.Stop()
	ach := timer.Done()
	select {
	case <-bch:
		t.Errorf("Done received when it shouldn't have")
	case <-time.After(dur):
	}
	select {
	case <-ach:
		t.Errorf("Done received when it shouldn't have")
	case <-time.After(dur):
	}
}

func TestManualTimerMultipleStop(t *testing.T) {
	const n = 10
	const dur = 100 * time.Millisecond
	timer := NewManualTimer()
	if !timer.Stop() {
		t.Fatalf("first call to stop didn't return true")
	}
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if timer.Stop() {
				t.Errorf("call # %d to stop returned true", i)
			}
		}(i)
	}

	wg.Wait()
	select {
	case <-timer.Done():
		t.Errorf("Done returned")
	case <-time.After(dur):
	}
}
