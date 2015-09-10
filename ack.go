package tme

import "time"

// Ack is sent by Ticker and Timer on each tick and when done (respectively
type Ack struct {
	// Time is the time the tick started or the timer was done
	Time time.Time
	// Fn is the function that the receiver of the tick or done signal may call to acknowledge receipt.
	// Some Ticker/Timer implementations provide a feedback mechanism for other goroutines
	// to be notified of the acks, which can be used for backoff or rate limiting purposes.
	// This func must never block.
	Fn func()
}
