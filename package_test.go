package outline_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/ioki-mobility/go-outline"
	"github.com/ioki-mobility/go-outline/internal/common"
	"github.com/ioki-mobility/go-outline/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testApiKey  string = "api key"
	testBaseURL string = "https://localhost.123"
)

func TestClientCollectionsGet_failed(t *testing.T) {
	tests := map[string]struct {
		isTemporary bool
		rt          http.RoundTripper
	}{
		"HTTP request failed": {
			isTemporary: false,
			rt: &testutils.MockRoundTripper{
				RoundTripFn: func(r *http.Request) (*http.Response, error) {
					return nil, &net.DNSError{}
				},
			},
		},
		"server side error": {
			isTemporary: true,
			rt: &testutils.MockRoundTripper{
				RoundTripFn: func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						Request:       r,
						StatusCode:    http.StatusServiceUnavailable,
						ContentLength: -1,
						Body:          io.NopCloser(strings.NewReader("service unavailable")),
					}, nil
				},
			},
		},
		"client side error": {
			isTemporary: false,
			rt: &testutils.MockRoundTripper{
				RoundTripFn: func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						Request:       r,
						ContentLength: -1,
						StatusCode:    http.StatusUnauthorized,
						Body:          io.NopCloser(strings.NewReader("unauthorized key")),
					}, nil
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hc := &http.Client{}
			hc.Transport = test.rt
			cl := outline.New(testBaseURL, hc, testApiKey)
			col, err := cl.Collections().Get("collection id").Do(context.Background())
			assert.Nil(t, col)
			require.NotNil(t, err)
			assert.Equal(t, test.isTemporary, outline.IsTemporary(err))
		})
	}
}

func TestClientCollectionsGet(t *testing.T) {
	testResponse := exampleCollectionsGetResponse

	// Prepare HTTP client with mocked transport.
	hc := &http.Client{}
	hc.Transport = &testutils.MockRoundTripper{RoundTripFn: func(r *http.Request) (*http.Response, error) {
		testHeaders(t, r.Header)
		return &http.Response{
			Request:       r,
			ContentLength: -1,
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(strings.NewReader(testResponse)),
		}, nil
	}}

	cl := outline.New(testBaseURL, hc, testApiKey)
	got, err := cl.Collections().Get("collection id").Do(context.Background())
	require.NoError(t, err)

	// Manually unmarshal test response and see if we get same object via the API.
	expected := &struct {
		Data outline.Collection `json:"data"`
	}{}
	require.NoError(t, json.Unmarshal([]byte(testResponse), expected))
	assert.Equal(t, &expected.Data, got)
}

func testHeaders(t *testing.T, headers http.Header) {
	t.Helper()
	assert.Equal(t, headers.Get(common.HdrKeyAccept), common.HdrValueAccept)
	assert.Equal(t, headers.Get(common.HdrKeyContentType), common.HdrValueContentType)
	assert.Equal(t, fmt.Sprintf("Bearer %s", testApiKey), headers.Get(common.HdrKeyAuthorization))
}

const exampleCollectionsGetResponse string = `{
  "data": {
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "name": "Human Resources",
    "description": "",
    "sort": {
      "field": "string",
      "direction": "asc"
    },
    "index": "P",
    "color": "#123123",
    "icon": "string",
    "permission": "read",
    "createdAt": "2019-08-24T14:15:22Z",
    "updatedAt": "2019-08-24T14:15:22Z",
    "deletedAt": "2019-08-24T14:15:22Z"
  }
}`
