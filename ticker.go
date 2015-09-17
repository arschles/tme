package tme

// Ticker continuously sends on Chan() channel until Stop is called
type Ticker interface {
	// Chan returns a channel that receives every time the ticker ticks. The chan will
	// not receieve and stay open when Stop is called. Multiple calls to Chan return
	// the same underlying channel.
	//
	// The receiver of each tick may call ack.Fn() to acknowledge its receipt of the tick,
	// but that is not a requirement.
	Chan() <-chan Ack
	// Stop stops the ticker so that it doesn't send on Chan() anymore. Multiple calls to Stop
	// are allowed. Returns true if the ticker was running beforehand, false otherwise.
	// Stopping a ticker doesn't close any of the channels returned by Chan().
	//
	// Always call Stop() on a ticker when you're done with it, as this func may release resources that
	// can then be garbage collected.
	Stop() bool
}
