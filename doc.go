// tme is a package that provides unified interfaces around tickers and timers for Go
// programs. It provides a Timer and Ticker interface, each with an implementation
// that uses real time and an implementation that can be manually controlled.
//
// Usage
//
// Generally speaking, the manual implementations are used for testing purposes and the real
// implementations for production purposes. Code should operate on the interfaces so that you can
// later implement and swap new implementations in for various different purposes
//
// Ack
//
// Each of the Timer and Ticker implementations contained herein use an Ack struct to notify
// of events (e.g. the timer finished or a tick happened).
//
// Ack contains both the time that the event occurred and a func() that you, the receiver of events can call
// to acknowledge (ack) the receipt of the event. The creator of the Ticker / Timer implements
// the func().
package tme
