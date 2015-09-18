package tme

import (
	"sync"
	"sync/atomic"
	"time"
)

const dur = 50 * time.Millisecond

func numTimerRecv(timer Timer, numThreads int) int32 {
	numDones := int32(0)
	wg := sync.WaitGroup{}
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()
			select {
			case <-timer.Done():
				atomic.AddInt32(&numDones, 1)
			case <-time.After(dur):
			}
		}()
	}
	wg.Wait()
	return atomic.LoadInt32(&numDones)
}

func numTickerRecv(ticker Ticker, numThreads int) int32 {
	numRecvs := int32(0)
	wg := sync.WaitGroup{}
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func(i int) {
			defer wg.Done()
			select {
			case <-ticker.Chan():
				atomic.AddInt32(&numRecvs, 1)
			case <-time.After(dur):
			}
		}(i)
	}
	wg.Wait()
	return numRecvs
}
