package testutils

import "net/http"

// MockRoundTripper implements [http.RoundTripper] interface. The tests can use this to mock HTTP transaction without
// creating a test HTTP server. Just override the RoundTripFn to shortcircuit and return whatever response you want.
type MockRoundTripper struct {
	RoundTripFn func(*http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.RoundTripFn != nil {
		return m.RoundTripFn(req)
	}

	return &http.Response{
		Request:    req,
		StatusCode: http.StatusNoContent,
	}, nil
}
