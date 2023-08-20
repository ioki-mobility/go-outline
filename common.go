package outline

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rsjethani/rsling"
)

// badResponse contains details of bad HTTP response returned by the server. For now, we are treating all bad responses
// as plain text just to be able to print/log. Later if you need to examine it further than we can convert it into proper
// concrete type.
type badResponse struct {
	clientErr string // only filled in case of 4XX response
	serverErr string // only filled in case of 5XX response
	status    int
	url       string
}

// request adds failure decoder to req and then makes the request bound by ctx. If everything goes fine then success
// would contain decoded response. If HTTP request did not complete normally then an error is returned. If request did
// complete but response was bad then badResponse would contain details. NOTE: Apart from adding failure decoder the
// req is used as is hence the caller must pass fully prepared req.
func request(ctx context.Context, req *rsling.Sling, success any) (*badResponse, error) {
	buf := &bytes.Buffer{}
	resp, err := req.FailureDecoder(rsling.ByteStreamer{}).ReceiveWithContext(ctx, success, buf)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusBadRequest {
		return nil, nil
	}

	br := &badResponse{
		status: resp.StatusCode,
		url:    resp.Request.URL.String(),
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		br.serverErr = buf.String()
	} else if resp.StatusCode >= http.StatusBadRequest {
		br.clientErr = buf.String()
	}

	return br, nil
}

// apiError wraps a bad HTTP response into an [error]. This allows bad response to be logged or chained to
// another error for more context.
type apiError struct {
	br badResponse
}

func (ae *apiError) Error() string {
	return fmt.Sprintf("%+v", ae.br)
}

// Temporary returns true for 5XX errors. This satisfies the temporary interface and enables usage with
// [outline.IsTemporary].
func (ae *apiError) Temporary() bool {
	return ae.br.status >= http.StatusInternalServerError && ae.br.status != http.StatusNotImplemented
}

// IsTemporary returns true if err is temporary in nature i.e. you can retry the same operation in some time.
// This would usually return true for server side errors.
func IsTemporary(err error) bool {
	var e temporary
	return errors.As(err, &e) && e.Temporary()
}

// temporary is supposed to be implemented by [error]s that want to indicate their temporary nature to the user. The
// user can then use [outline.IsTemporary] to check this. Reference from standard library:
// https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/net/net.go;l=507
type temporary interface {
	Temporary() bool
}

// pagination represents pagination logic related metadata usually part of responses containing list of items.
type pagination struct {
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	NextPath string `json:"nextPath"`
}

// paginationQueryParams contains valid query paramters for pagination logic.
// Reference: https://www.getoutline.com/developers#section/Pagination
type paginationQueryParams struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}
