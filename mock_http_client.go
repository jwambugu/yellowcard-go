package yellowcard

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

// mockResponseFunc simulates HTTP responses
type mockResponseFunc func() (status int, body string)

type mockResponse struct {
	fn mockResponseFunc
}

type mockHttpClient struct {
	mu        sync.Mutex
	responses map[string]mockResponse
	requests  []*http.Request
}

func mockHttpResponse(status int, body string) *http.Response {
	return &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       io.NopCloser(bytes.NewBuffer([]byte(body))),
	}
}

func newMockHttpClient() *mockHttpClient {
	return &mockHttpClient{
		responses: make(map[string]mockResponse),
	}
}

// MockRequest appends the given response for the provided url.
func (m *mockHttpClient) MockRequest(url string, fn mockResponseFunc) {
	m.responses[url] = mockResponse{fn: fn}
}

// Do checks if the given req.URL exists in the available responses lists and returns the stored response.
// If none exists, it returns status http.StatusNotFound
func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requests = append(m.requests, req.Clone(req.Context()))

	if mock, ok := m.responses[req.URL.String()]; ok {
		if mock.fn != nil {
			status, body := mock.fn()
			return mockHttpResponse(status, body), nil
		}
	}

	return mockHttpResponse(http.StatusNotFound, `{"code": "The error code","message": "The error message"}`), nil
}
