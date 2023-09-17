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
	testApiKey    string = "api key"
	testServerURL string = "https://localhost.123"
)

func TestClientCollectionsStructure_failed(t *testing.T) {
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
			cl := outline.New(testServerURL, hc, testApiKey)
			col, err := cl.Collections().DocumentStructure("collection id").Do(context.Background())
			assert.Nil(t, col)
			require.NotNil(t, err)
			assert.Equal(t, test.isTemporary, outline.IsTemporary(err))
		})
	}
}

func TestClientCollectionsStructure(t *testing.T) {
	testResponse := exampleCollectionsDocumentStructureResponse

	// Prepare HTTP client with mocked transport.
	hc := &http.Client{}
	hc.Transport = &testutils.MockRoundTripper{RoundTripFn: func(r *http.Request) (*http.Response, error) {
		// Assert request method and URL.
		assert.Equal(t, http.MethodPost, r.Method)
		u, err := url.JoinPath(common.BaseURL(testServerURL), common.CollectionsStructureEndpoint())
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

	cl := outline.New(testServerURL, hc, testApiKey)
	got, err := cl.Collections().DocumentStructure("collection id").Do(context.Background())
	require.NoError(t, err)

	// Manually unmarshal test response and see if we get same object via the API.
	expected := struct {
		Data outline.DocumentStructure `json:"data"`
	}{}
	require.NoError(t, json.Unmarshal([]byte(testResponse), &expected))
	assert.Equal(t, expected.Data, got)
}

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
			cl := outline.New(testServerURL, hc, testApiKey)
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
		u, err := url.JoinPath(common.BaseURL(testServerURL), common.CollectionsGetEndpoint())
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

	cl := outline.New(testServerURL, hc, testApiKey)
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
			u, err := url.JoinPath(common.BaseURL(testServerURL), common.CollectionsListEndpoint())
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
		u, err := url.JoinPath(common.BaseURL(testServerURL), common.CollectionsListEndpoint())
		require.NoError(t, err)
		assert.Equal(t, u+"?offset=2", r.URL.String())

		return &http.Response{
			Request:       r,
			StatusCode:    http.StatusOK,
			ContentLength: -1,
			Body:          io.NopCloser(strings.NewReader(exampleCollectionsListResponse_1collection)),
		}, nil
	}}

	cl := outline.New(testServerURL, hc, testApiKey)

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
		u, err := url.JoinPath(common.BaseURL(testServerURL), common.CollectionsCreateEndpoint())
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

	cl := outline.New(testServerURL, hc, testApiKey)
	got, err := cl.Collections().
		Create("new collection").
		PermissionRead().
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

func TestDocumentsClientCreate(t *testing.T) {
	testResponse := exampleDocumentResponse

	// Prepare HTTP client with mocked transport.
	hc := &http.Client{}
	hc.Transport = &testutils.MockRoundTripper{RoundTripFn: func(r *http.Request) (*http.Response, error) {
		// Assert request method and URL.
		assert.Equal(t, http.MethodPost, r.Method)
		u, err := url.JoinPath(common.BaseURL(testServerURL), common.DocumentsCreateEndpoint())
		require.NoError(t, err)
		assert.Equal(t, u, r.URL.String())

		testAssertHeaders(t, r.Header)
		testAssertBody(
			t,
			r,
			fmt.Sprintf(
				`{"collectionId":"%s", "title":"%s", "text":"%s", "publish":%t}`,
				"collection id",
				"🎉 Welcome to Acme Inc",
				"Some text",
				true,
			),
		)

		return &http.Response{
			Request:       r,
			ContentLength: -1,
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(strings.NewReader(testResponse)),
		}, nil
	}}

	cl := outline.New(testServerURL, hc, testApiKey)
	var id outline.CollectionID = "collection id"
	got, err := cl.Documents().Create("🎉 Welcome to Acme Inc", id).Text("Some text").Publish(true).Do(context.Background())
	require.NoError(t, err)

	// Manually unmarshal test response and see if we get same object via the API.
	expected := &struct {
		Data outline.Document `json:"data"`
	}{}
	require.NoError(t, json.Unmarshal([]byte(testResponse), expected))
	assert.Equal(t, &expected.Data, got)
}

func TestDocumentsClientUpdate(t *testing.T) {
	testResponse := exampleDocumentResponse

	// Prepare HTTP client with mocked transport.
	hc := &http.Client{}
	hc.Transport = &testutils.MockRoundTripper{RoundTripFn: func(r *http.Request) (*http.Response, error) {
		// Assert request method and URL.
		assert.Equal(t, http.MethodPost, r.Method)
		u, err := url.JoinPath(common.BaseURL(testServerURL), common.DocumentsUpdateEndpoint())
		require.NoError(t, err)
		assert.Equal(t, u, r.URL.String())

		testAssertHeaders(t, r.Header)
		testAssertBody(
			t,
			r,
			fmt.Sprintf(
				`{"id":"%s", "title":"%s", "text":"%s", "publish":%t}`,
				"497f6eca-6276-4993-bfeb-53cbbbba6f08",
				"🎉 Welcome to Acme Inc",
				"Updated text!",
				true,
			),
		)

		return &http.Response{
			Request:       r,
			ContentLength: -1,
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(strings.NewReader(testResponse)),
		}, nil
	}}

	cl := outline.New(testServerURL, hc, testApiKey)

	got, err := cl.Documents().Update("497f6eca-6276-4993-bfeb-53cbbbba6f08").
		Title("🎉 Welcome to Acme Inc").Text("Updated text!").Publish(true).
		Do(context.Background())
	require.NoError(t, err)

	// Manually unmarshal test response and see if we get same object via the API.
	expected := &struct {
		Data outline.Document `json:"data"`
	}{}
	require.NoError(t, json.Unmarshal([]byte(testResponse), expected))
	assert.Equal(t, &expected.Data, got)
}

func TestDocumentsClientUpdate_failed(t *testing.T) {
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
		"resource not found": {
			isTemporary: false,
			rt: &testutils.MockRoundTripper{
				RoundTripFn: func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						Request:       r,
						ContentLength: -1,
						StatusCode:    http.StatusUnauthorized,
						Body:          io.NopCloser(strings.NewReader("the document with the id doesn't exist")),
					}, nil
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hc := &http.Client{}
			hc.Transport = test.rt
			cl := outline.New(testServerURL, hc, testApiKey)
			col, err := cl.Documents().Update("id").Do(context.Background())
			assert.Nil(t, col)
			require.NotNil(t, err)
			assert.Equal(t, test.isTemporary, outline.IsTemporary(err))
		})
	}
}

func TestAttachmentsClientCreate(t *testing.T) {
	testResponse := exampleAttachmentsCreateResponse

	// Prepare HTTP client with mocked transport.
	hc := &http.Client{}
	hc.Transport = &testutils.MockRoundTripper{RoundTripFn: func(r *http.Request) (*http.Response, error) {
		// Assert request method and URL.
		assert.Equal(t, http.MethodPost, r.Method)
		u, err := url.JoinPath(common.BaseURL(testServerURL), common.AttachmentsCreateEndpoint())
		require.NoError(t, err)
		assert.Equal(t, u, r.URL.String())

		testAssertHeaders(t, r.Header)
		testAssertBody(
			t,
			r,
			fmt.Sprintf(
				`{"name":"%s", "contentType":"%s", "size":%d, "documentId":"%s"}`,
				"My Image",
				"image/png",
				42,
				"4704590c-004e-410d-adf7-acb7ca0a7052",
			),
		)

		return &http.Response{
			Request:       r,
			ContentLength: -1,
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(strings.NewReader(testResponse)),
		}, nil
	}}

	cl := outline.New(testServerURL, hc, testApiKey)
	got, err := cl.Attachments().Create("My Image", "image/png", 42).
		DocumentID("4704590c-004e-410d-adf7-acb7ca0a7052").
		Do(context.Background())
	require.NoError(t, err)

	// Manually unmarshal test response and see if we get same object via the API.
	expected := &struct {
		Data outline.Attachment `json:"data"`
	}{}
	require.NoError(t, json.Unmarshal([]byte(testResponse), expected))
	assert.Equal(t, &expected.Data, got)
}

func TestAttachmentsClientCreate_failed(t *testing.T) {
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
		"request failed": {
			isTemporary: false,
			rt: &testutils.MockRoundTripper{
				RoundTripFn: func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						Request:       r,
						ContentLength: -1,
						StatusCode:    http.StatusBadRequest,
						Body:          io.NopCloser(strings.NewReader("validation failed")),
					}, nil
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hc := &http.Client{}
			hc.Transport = test.rt
			cl := outline.New(testServerURL, hc, testApiKey)
			col, err := cl.Attachments().Create("name", "image/png", 42).Do(context.Background())
			assert.Nil(t, col)
			require.NotNil(t, err)
			assert.Equal(t, test.isTemporary, outline.IsTemporary(err))
		})
	}
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

const exampleDocumentResponse string = `{
	"data": {
		"id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
		"collectionId": "collection id",
		"parentDocumentId": "ce8a7254-3ff2-448e-a302-0033b010f00b",
		"title": "🎉 Welcome to Acme Inc",
		"fullWidth": true,
		"emoji": "🎉",
		"text": "Some text",
		"urlId": "hDYep1TPAM",
		"collaborators": [],
		"pinned": true,
		"template": true,
		"templateId": "196100ac-4eec-4fb6-a7f7-86c8b584771d",
		"revision": 0,
		"createdAt": "2019-08-24T14:15:22Z",
		"createdBy": {},
		"updatedAt": "2019-08-24T14:15:22Z",
		"updatedBy": {},
		"publishedAt": "2019-08-24T14:15:22Z",
		"archivedAt": "2019-08-24T14:15:22Z",
		"deletedAt": "2019-08-24T14:15:22Z"
	}
}`

const exampleCollectionsDocumentStructureResponse string = `
{
  "data": [
	{
		"id": "doc1",
		"title": "Doc 1",
		"url": "https://doc1.url"
	},
	{
		"id": "doc2",
		"title": "Doc 2",
		"url": "https://doc2.url",
		"children": [
			{
				"id": "doc2-1",
				"title": "Doc 2-1",
				"url": "https://doc2-1.url"
			}
		]
	}
  ]
}
`

const exampleAttachmentsCreateResponse string = `{
	"data": {
		"maxUploadSize": 0,
		"uploadUrl": "https://s3.ioki.com",
		"form": { },
		"attachment": {
			"contentType": "image/png",
			"size": 42,
			"name": "My Image",
			"url": "https://ioki.com",
			"documentId": "4704590c-004e-410d-adf7-acb7ca0a7052"
		}
	}
}`
