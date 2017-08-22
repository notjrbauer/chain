package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func AddHeader(key, val string) func(chain Chain) (*http.Response, error) {
	return func(chain Chain) (*http.Response, error) {
		chain.Request().Header.Add(key, val)
		return chain.Proceed(chain.Request())
	}
}

func Abort(chain Chain) (*http.Response, error) {
	header := http.Header{}
	header.Add("test", "ok")
	return &http.Response{Header: header}, nil
}

func TestPrependingRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				get := "foo"
				got := r.Header.Get(get)
				want := "bar"
				if got != want {
					t.Errorf("Header.Get(%s) = %s; want %s", get, got, want)
				}
				get = "fizz"
				got = r.Header.Get(get)
				want = "buzz"
				if got != want {
					t.Errorf("Header.Get(%s) = %s; want %s", get, got, want)
				}
				w.WriteHeader(http.StatusOK)
			},
		),
	)
	transport := Transport(
		CallbackFunc(AddHeader("foo", "bar")),
		CallbackFunc(AddHeader("fizz", "buzz")),
	)
	client := &http.Client{Transport: transport}
	url, _ := url.Parse(ts.URL)
	req := http.Request{
		URL:    url,
		Method: "GET",
	}
	res, err := client.Do(&req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 status code, got %d", res.StatusCode)
	}
	if err != nil {
		t.Error(err)
	}
}

func TestShortCircuit(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r != nil {
					t.Error("Should not recieve a request")
				}
			},
		),
	)
	transport := Transport(
		CallbackFunc(AddHeader("foo", "bar")),
		CallbackFunc(Abort),
	)
	client := &http.Client{Transport: transport}
	url, _ := url.Parse(ts.URL)

	req := http.Request{URL: url, Method: "GET"}
	res, err := client.Do(&req)
	if err != nil {
		t.Error(err)
	}

	// OK
	if res.StatusCode != 0 {
		t.Errorf("Expected 0 status code, got %d", res.StatusCode)
	}

	get := "test"
	got := res.Header.Get(get)
	want := "ok"
	if got != want {
		t.Errorf("Header.Get(%s) = %s; want %s", get, got, want)
	}
}
