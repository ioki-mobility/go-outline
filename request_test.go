package outline

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dghubble/sling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_request_failed(t *testing.T) {
	client := &http.Client{}
	client.Transport = &mockRoundTripper{
		RoundTripFn: func(req *http.Request) (*http.Response, error) {
			return nil, &net.DNSError{} // simulate irrecoverable error
		},
	}

	sl := sling.New().Client(client)

	ed, err := request(context.Background(), sl, nil)
	assert.NotNil(t, err)
	assert.Nil(t, ed)
}

func Test_request_returns_bad_response(t *testing.T) {
	u, err := url.Parse("https://test.url")
	require.NoError(t, err)

	tests := map[string]struct {
		body     string
		expected badResponse
	}{
		"4XX response": {
			body: "HTTP 400",
			expected: badResponse{
				status:    http.StatusBadRequest,
				url:       u.String(),
				clientErr: "HTTP 400",
				serverErr: "",
			},
		},
		"5XX response": {
			body: "HTTP 503",
			expected: badResponse{
				status:    http.StatusServiceUnavailable,
				url:       u.String(),
				clientErr: "",
				serverErr: "HTTP 503",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &http.Client{}
			client.Transport = &mockRoundTripper{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Request:       req,
						StatusCode:    test.expected.status,
						ContentLength: -1,
						Body:          io.NopCloser(strings.NewReader(test.body)),
					}, nil
				},
			}

			sl := sling.New().Client(client).Get(test.expected.url)
			got, err := request(context.Background(), sl, nil)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, *got)
		})
	}
}

func Test_makeRequest(t *testing.T) {
	type testModel struct {
		Field1 string  `json:"field_1"`
		Field2 float64 `json:"field_2"`
	}

	// Create HTTP client and override its transport to simulate server returning a good response with valid JSON data
	// of testModel type.
	client := &http.Client{}
	client.Transport = &mockRoundTripper{
		RoundTripFn: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Request:       req,
				StatusCode:    http.StatusOK,
				ContentLength: -1,
				Body:          io.NopCloser(strings.NewReader(`{"field_1":"hello","field_2":5.5}`)),
			}, nil
		},
	}

	expected := &testModel{
		Field1: "hello",
		Field2: 5.5,
	}
	got := &testModel{}

	sl := sling.New().Client(client).Get("https://some.url")
	ed, err := request(context.Background(), sl, got)
	assert.Nil(t, err)
	assert.Nil(t, ed)
	assert.Equal(t, expected, got)
}

// mockRoundTripper implements [http.RoundTripper] interface. The tests can use this to mock HTTP transaction without
// creating a test HTTP server. Just override the RoundTripFn to shortcircuit and return whatever response you want.
type mockRoundTripper struct {
	RoundTripFn func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.RoundTripFn != nil {
		return m.RoundTripFn(req)
	}

	return &http.Response{
		Request:    req,
		StatusCode: http.StatusNoContent,
	}, nil
}
