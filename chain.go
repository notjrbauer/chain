package chain

import (
	"context"
	"net/http"
)

// The contract of being a chain.
type Chain interface {
	// Moves the chain to and executes the next callback.
	// If no other callbacks are defined, the request is dispatched out.
	Proceed(*http.Request) (*http.Response, error)
	// Returns the request to be sent out.
	Request() *http.Request
	// Returns the request's context.
	Context() context.Context
}

// Observes, modifies, and potentially short-ciruits requests going out and the corresponding
// responses coming back in. Typically callbacks add, remove or transport headers on the request
// or response.
type Callback interface {
	Call(Chain) (*http.Response, error)
}

// Callback adapter function to wrap functions that accept a chain.
type CallbackFunc func(Chain) (*http.Response, error)

func (cb CallbackFunc) Call(ch Chain) (*http.Response, error) {
	return cb(ch)
}

type chain struct {
	index     int
	req       *http.Request
	transport http.RoundTripper
	callbacks []Callback
}

func (c *chain) Proceed(req *http.Request) (*http.Response, error) {
	if c.index < len(c.callbacks) {
		link := &chain{
			index:     c.index + 1,
			req:       req,
			transport: c.transport,
			callbacks: c.callbacks,
		}
		cb := c.callbacks[c.index]
		return cb.Call(link)
	}
	return c.transport.RoundTrip(req)
}

func (c *chain) Context() context.Context {
	return c.Request().Context()
}

func (c *chain) Request() *http.Request {
	return c.req
}

type roundtripper struct {
	transport http.RoundTripper
	callbacks []Callback
}

// Called by the http client to initiate the request.
// A chain is initialized and kicked off with the provided request.
func (rt *roundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ch := &chain{
		index:     0,
		transport: rt.transport,
		callbacks: rt.callbacks,
		req:       req,
	}
	return ch.Proceed(req)
}

// Returns a RoundTripper with n arity of callbacks passed in.
// http.DefaultTransport is used as the default RoundTripper.
func Transport(callbacks ...Callback) http.RoundTripper {
	return WithTransport(http.DefaultTransport, callbacks...)
}

// Returns a wrapped RoundTripper with the provided transport, and n arity of callbacks.
func WithTransport(rt http.RoundTripper, callbacks ...Callback) http.RoundTripper {
	return &roundtripper{
		transport: rt,
		callbacks: append([]Callback{}, callbacks...),
	}
}
