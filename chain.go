package main

import (
	"context"
	"net/http"
)

type Chain interface {
	Proceed(*http.Request) (*http.Response, error)
	Request() *http.Request
	Context() context.Context
}

type Callback interface {
	Call(Chain) (*http.Response, error)
}

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

func (rt *roundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ch := &chain{
		index:     0,
		transport: rt.transport,
		callbacks: rt.callbacks,
		req:       req,
	}
	return ch.Proceed(req)
}

func Transport(callbacks ...Callback) http.RoundTripper {
	return WithTransport(http.DefaultTransport, callbacks...)
}

func WithTransport(rt http.RoundTripper, callbacks ...Callback) http.RoundTripper {
	return &roundtripper{
		transport: rt,
		callbacks: append([]Callback{}, callbacks...),
	}
}
