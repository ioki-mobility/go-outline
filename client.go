package outline

import (
	"net/http"
	"strings"

	"github.com/dghubble/sling"
	"github.com/ioki-mobility/go-outline/internal/common"
)

// Client is per server top level client which acts as entry point and stores common configuration (like base url) for
// resource level clients. It is preferred to reuse same client while communicating to the same server as this makes
// better utilization of resources.
type Client struct {
	// base acts as the 'base' request on which various common properties like HTTP headers, server url etc. are
	// configured. The resource level clients create their own customized request derived from this.
	base *sling.Sling
}

// New creates and returns a new per server client.
func New(baseURL string, hc *http.Client, apiKey string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	sl := sling.New().Client(hc).Base(baseURL)
	sl.Set(common.HdrKeyAuthorization, common.HdrValueAuthorization(apiKey))
	sl.Set(common.HdrKeyContentType, common.HdrValueContentType)
	sl.Set(common.HdrKeyAccept, common.HdrValueAccept)

	return &Client{base: sl}
}

// Documents creates client for operating on documents.
func (cl *Client) Documents() *DocumentsClient { return nil }

// Collections creates client for operating on collections.
func (cl *Client) Collections() *CollectionsClient {
	return &CollectionsClient{sl: cl.base.New()}
}
