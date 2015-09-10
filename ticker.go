package tme

// Ticker continuously sends on Chan() channel until Stop is called
type Ticker interface {
	// Chan returns a channel that receives every time the ticker ticks. The chan will be
	// closed when ticking is done (e.g. when Stop() is called).
	//
	// The receiver of the tick should call ack.Fn() to acknowledge its receipt of the tick.
	Chan() <-chan Ack
	// Stop stops the ticker so that it doesn't send on Chan() anymore. Multiple calls to Stop
	// are allowed. Returns true if the ticker was running beforehand, false otherwise.
	// Stopping a ticker doesn't close any of the channels returned by Chan().
	//
	// Always call Stop() on a ticker when you're done with it, as this func may release resources that
	// can then be garbage collected.
	Stop() bool
}
