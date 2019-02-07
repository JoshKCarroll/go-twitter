package twitter

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

// NewTestStream creates a Stream for testing with a provided input Messages channel.
// It is safe to call Stop() once on the provided *Stream, as with a normal Stream.
func NewTestStream(messages chan interface{}) *Stream {
	return &Stream{
		Messages: messages,
		done:     make(chan struct{}),
	}
}

// NewTestServer exposes testServer for test scaffolding in libraries that use go-twitter
// it takes a map of path:functions to set the ServeMux.
func NewTestServer(handlers map[string]func(w http.ResponseWriter, r *http.Request)) (*http.Client, *httptest.Server) {
	client, mux, server := testServer()
	for path, handler := range handlers {
		mux.HandleFunc(path, handler)
	}
	return client, server
}

// testServer returns an http Client, ServeMux, and Server. The client proxies
// requests to the server and handlers can be registered on the mux to handle
// requests. The caller must close the test server.
func testServer() (*http.Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	transport := &RewriteTransport{&http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}}
	client := &http.Client{Transport: transport}
	return client, mux, server
}

// RewriteTransport rewrites https requests to http to avoid TLS cert issues
// during testing.
type RewriteTransport struct {
	Transport http.RoundTripper
}

// RoundTrip rewrites the request scheme to http and calls through to the
// composed RoundTripper or if it is nil, to the http.DefaultTransport.
func (t *RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	if t.Transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.Transport.RoundTrip(req)
}
