package tme

// Timer sends once on the Done() channel unless Stop() is called beforehand
type Timer interface {
	// Done returns a channel that receives and closes after the timer is done.
	// It should return the same underlying channel on every call, and the channel
	// returned by Done will never receive nor stop after stop is called.
	Done() <-chan Ack
	// Stop cancels the Timer. If this timer is not done when Stop is called, the timer
	// will not send on any channels returned by Done (but it won't close any of those chans).
	// returns true if the timer was not previously done, false otherwise.
	Stop() bool
}
