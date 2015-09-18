package tme

import (
	"sync"
	"sync/atomic"
	"time"
)

const dur = 50 * time.Millisecond

func numTimerRecv(timer Timer, n int) int32 {
	numDones := int32(0)
	wg := sync.WaitGroup{}
	wg.Add(1)
	for i := 0; i < n; i++ {
		go func() {
			select {
			case <-timer.Done():
				atomic.AddInt32(&numDones, 1)
				wg.Done()
			case <-time.After(dur):
			}
		}()
	}
	wg.Wait()
	return atomic.LoadInt32(&numDones)
}
