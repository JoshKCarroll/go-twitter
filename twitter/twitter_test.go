package twitter

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var defaultTestTimeout = time.Second * 1

func assertMethod(t *testing.T, expectedMethod string, req *http.Request) {
	assert.Equal(t, expectedMethod, req.Method)
}

// assertQuery tests that the Request has the expected url query key/val pairs
func assertQuery(t *testing.T, expected map[string]string, req *http.Request) {
	queryValues := req.URL.Query()
	expectedValues := url.Values{}
	for key, value := range expected {
		expectedValues.Add(key, value)
	}
	assert.Equal(t, expectedValues, queryValues)
}

// assertPostForm tests that the Request has the expected key values pairs url
// encoded in its Body
func assertPostForm(t *testing.T, expected map[string]string, req *http.Request) {
	req.ParseForm() // parses request Body to put url.Values in r.Form/r.PostForm
	expectedValues := url.Values{}
	for key, value := range expected {
		expectedValues.Add(key, value)
	}
	assert.Equal(t, expectedValues, req.Form)
}

// assertDone asserts that the empty struct channel is closed before the given
// timeout elapses.
func assertDone(t *testing.T, ch <-chan struct{}, timeout time.Duration) {
	select {
	case <-ch:
		_, more := <-ch
		assert.False(t, more)
	case <-time.After(timeout):
		t.Errorf("expected channel to be closed within timeout %v", timeout)
	}
}

// assertClosed asserts that the channel is closed before the given timeout
// elapses.
func assertClosed(t *testing.T, ch <-chan interface{}, timeout time.Duration) {
	select {
	case <-ch:
		_, more := <-ch
		assert.False(t, more)
	case <-time.After(timeout):
		t.Errorf("expected channel to be closed within timeout %v", timeout)
	}
}
