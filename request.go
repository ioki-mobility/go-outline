package outline

import (
	"bytes"
	"context"
	"net/http"

	"github.com/dghubble/sling"
)

// badResponse contains details of bad HTTP response returned by the server. For now we are treating all bad responses
// as plain text just to be able to print/log. Later if need to examine it further then we can convert it into proper
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
func request(ctx context.Context, req *sling.Sling, success any) (*badResponse, error) {
	buf := &bytes.Buffer{}
	resp, err := req.FailureDecoder(sling.ByteStreamer{}).ReceiveWithContext(ctx, success, buf)
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
