package outline_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
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
		// Assert request method and URL.
		assert.Equal(t, http.MethodPost, r.Method)
		u, err := url.JoinPath(testBaseURL, common.CollectionsGetEndpoint())
		require.NoError(t, err)
		assert.Equal(t, u, r.URL.String())

		testAssertHeaders(t, r.Header)
		testAssertBody(t, r, fmt.Sprintf(`{"id":"%s"}`, "collection id"))

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

func TestClientCollectionsList(t *testing.T) {
	requestCount := atomic.Uint32{}
	hc := &http.Client{}
	hc.Transport = &testutils.MockRoundTripper{RoundTripFn: func(r *http.Request) (*http.Response, error) {
		requestCount.Add(1)

		assert.Equal(t, http.MethodPost, r.Method)
		testAssertHeaders(t, r.Header)

		if requestCount.Load() == 1 {
			// Assert URL when asking first page.
			u, err := url.JoinPath(testBaseURL, common.CollectionsListEndpoint())
			require.NoError(t, err)
			assert.Equal(t, u, r.URL.String())

			return &http.Response{
				Request:       r,
				StatusCode:    http.StatusOK,
				ContentLength: -1,
				Body:          io.NopCloser(strings.NewReader(exampleCollectionsListResponse_2collections)),
			}, nil
		}

		// Assert URL when asking second page (first page had 2 items).
		// NOTE: There is some hard coding here but that is okay, no need to over-engineer as of now.
		u, err := url.JoinPath(testBaseURL, common.CollectionsListEndpoint())
		require.NoError(t, err)
		assert.Equal(t, u+"?offset=2", r.URL.String())

		return &http.Response{
			Request:       r,
			StatusCode:    http.StatusOK,
			ContentLength: -1,
			Body:          io.NopCloser(strings.NewReader(exampleCollectionsListResponse_1collection)),
		}, nil
	}}

	cl := outline.New(testBaseURL, hc, testApiKey)

	collectionsListFnCalled := atomic.Uint32{}
	err := cl.Collections().List().Do(context.Background(), func(c *outline.Collection, err error) (bool, error) {
		collectionsListFnCalled.Add(1)
		return true, nil
	})
	require.NoError(t, err)
	assert.Equal(t, uint32(3), collectionsListFnCalled.Load())
}

func TestClientCollectionsCreate(t *testing.T) {
	testResponse := exampleCollectionsGetResponse

	hc := &http.Client{}
	hc.Transport = &testutils.MockRoundTripper{RoundTripFn: func(r *http.Request) (*http.Response, error) {
		// Assert request method and URL.
		assert.Equal(t, http.MethodPost, r.Method)
		u, err := url.JoinPath(testBaseURL, common.CollectionsCreateEndpoint())
		require.NoError(t, err)
		assert.Equal(t, u, r.URL.String())

		testAssertHeaders(t, r.Header)
		testAssertBody(t, r, fmt.Sprintf(`{"name":"%s", "permission":"%s", "description":"%s"}`, "new collection", "read", "desc"))

		return &http.Response{
			Request:       r,
			StatusCode:    http.StatusOK,
			ContentLength: -1,
			Body:          io.NopCloser(strings.NewReader(exampleCollectionsGetResponse)),
		}, nil
	}}

	cl := outline.New(testBaseURL, hc, testApiKey)
	got, err := cl.Collections().
		Create("new collection").
		Permission(outline.PermissionRead).
		Description("desc").
		Do(context.Background())
	require.NoError(t, err)

	// Manually unmarshal test response and see if we get same object via the API.
	expected := &struct {
		Data outline.Collection `json:"data"`
	}{}
	require.NoError(t, json.Unmarshal([]byte(testResponse), expected))
	assert.Equal(t, &expected.Data, got)
}

func testAssertHeaders(t *testing.T, headers http.Header) {
	t.Helper()
	assert.Equal(t, headers.Get(common.HdrKeyAccept), common.HdrValueAccept)
	assert.Equal(t, headers.Get(common.HdrKeyContentType), common.HdrValueContentType)
	assert.Equal(t, fmt.Sprintf("Bearer %s", testApiKey), headers.Get(common.HdrKeyAuthorization))
}

func testAssertBody(t *testing.T, r *http.Request, expected string) {
	t.Helper()
	require.NotNil(t, r.Body)
	b, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	assert.JSONEq(t, expected, string(b))
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

const exampleCollectionsListResponse_2collections string = `
	{
  "data": [
    {
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
    },
    {
      "id": "111f6eca-6276-4993-bfeb-53cbbbba6f08",
      "name": "Human Resources 2",
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
  ],
  "pagination": {
    "offset": 0,
    "limit": 25
  }
}`

const exampleCollectionsListResponse_1collection string = `
	{
  "data": [
    {
      "id": "111f6eca-6276-4993-bfeb-53cbbbba6f08",
      "name": "Human Resources 3",
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
  ],
  "pagination": {
    "offset": 0,
    "limit": 25
  }
}`
